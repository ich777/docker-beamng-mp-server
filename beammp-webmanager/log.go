package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	defaultTailBytes = 96 << 10
	maxTailBytes     = 2 << 20
)

func (m *Manager) findLog() string {
	if p := filepath.Join(m.dir, "Server.log"); fileExists(p) {
		return p
	}
	entries, err := os.ReadDir(m.dir)
	if err != nil {
		return ""
	}
	type cand struct {
		path string
		mod  time.Time
	}
	var cands []cand
	for _, e := range entries {
		if e.IsDir() || !strings.EqualFold(filepath.Ext(e.Name()), ".log") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		cands = append(cands, cand{filepath.Join(m.dir, e.Name()), info.ModTime()})
	}
	if len(cands) == 0 {
		return ""
	}
	sort.Slice(cands, func(i, j int) bool { return cands[i].mod.After(cands[j].mod) })
	return cands[0].path
}

func fileExists(p string) bool {
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func tail(path string, n int64) (string, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return "", 0, err
	}
	size := st.Size()
	start := size - n
	if start < 0 {
		start = 0
	}
	if _, err := f.Seek(start, io.SeekStart); err != nil {
		return "", size, err
	}
	buf, err := io.ReadAll(io.LimitReader(f, n))
	if err != nil {
		return "", size, err
	}
	text := string(buf)
	if start > 0 {
		if i := strings.IndexByte(text, '\n'); i >= 0 {
			text = text[i+1:]
		}
	}
	if !utf8.ValidString(text) {
		text = strings.ToValidUTF8(text, "�")
	}
	return strings.ReplaceAll(text, "\r\n", "\n"), size, nil
}

func (a *app) handleLog(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	mgr := a.mgr
	a.mu.Unlock()

	n := int64(defaultTailBytes)
	if v, err := strconv.ParseInt(r.URL.Query().Get("bytes"), 10, 64); err == nil && v > 0 {
		n = min(v, maxTailBytes)
	}

	res := map[string]any{}
	path := mgr.findLog()
	if path == "" {
		res["error"] = msg("log_missing", mgr.dir)
		writeJSON(w, res)
		return
	}
	res["file"] = filepath.Base(path)
	text, size, err := tail(path, n)
	if err != nil {
		res["error"] = msg("log_unreadable", filepath.Base(path), err.Error())
		writeJSON(w, res)
		return
	}
	res["text"] = text
	res["size"] = size
	res["truncated"] = size > int64(len(text))
	writeJSON(w, res)
}
