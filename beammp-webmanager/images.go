package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

const (
	maxImageBytes = 16 << 20
	maxCacheBytes = 64 << 20
)

type cachedImage struct {
	data  []byte
	ctype string
}

type imageCache struct {
	mu    sync.Mutex
	items map[string]cachedImage
	size  int
}

var imgCache = imageCache{items: map[string]cachedImage{}}

func (c *imageCache) get(key string) (cachedImage, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.items[key]
	return v, ok
}

func (c *imageCache) put(key string, v cachedImage) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.size+len(v.data) > maxCacheBytes {
		c.items = map[string]cachedImage{}
		c.size = 0
	}
	c.items[key] = v
	c.size += len(v.data)
}

func (m *Manager) locateZip(name, kind string) (string, error) {
	if name == "" || name != filepath.Base(name) {
		return "", errors.New("invalid file name")
	}
	cfg, err := m.readConfig()
	if err != nil {
		return "", err
	}
	client, deactMods, _, _, customMaps := m.paths(cfg)
	dirs := []string{client, deactMods}
	if kind == "map" {
		dirs = []string{client, customMaps}
	}
	for _, d := range dirs {
		p := filepath.Join(d, name)
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p, nil
		}
	}
	return "", fmt.Errorf("%s not found", name)
}

func zipEntryBytes(file, entry string) ([]byte, error) {
	r, err := zip.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	for _, f := range r.File {
		if f.Name != entry {
			continue
		}
		if f.UncompressedSize64 > maxImageBytes {
			return nil, errors.New("image too large")
		}
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		defer rc.Close()
		return io.ReadAll(io.LimitReader(rc, maxImageBytes))
	}
	return nil, errors.New("image not present in the archive")
}

func contentType(name string, data []byte) string {
	if ct := http.DetectContentType(data); strings.HasPrefix(ct, "image/") {
		return ct
	}
	switch strings.ToLower(path.Ext(name)) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	}
	return "application/octet-stream"
}

func (a *app) handleImg(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	a.mu.Lock()
	mgr := a.mgr
	a.mu.Unlock()

	kind := q.Get("kind")
	if kind != "map" && kind != "mod" {
		http.Error(w, "kind missing", http.StatusBadRequest)
		return
	}
	src, err := mgr.locateZip(q.Get("zip"), kind)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	zi, err := inspectZip(src)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	entry := zi.modImage
	if kind == "map" {
		entry = zi.mapImage
	}
	if entry == "" {
		http.NotFound(w, r)
		return
	}

	st, err := os.Stat(src)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	key := fmt.Sprintf("%s|%d|%d|%s", src, st.ModTime().UnixNano(), st.Size(), entry)
	if img, ok := imgCache.get(key); ok {
		serveImage(w, img)
		return
	}

	data, err := zipEntryBytes(src, entry)
	if err != nil || len(data) == 0 {
		http.NotFound(w, r)
		return
	}
	img := cachedImage{data: data, ctype: contentType(entry, data)}
	imgCache.put(key, img)
	serveImage(w, img)
}

func serveImage(w http.ResponseWriter, img cachedImage) {
	w.Header().Set("Content-Type", img.ctype)
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(img.data)
}
