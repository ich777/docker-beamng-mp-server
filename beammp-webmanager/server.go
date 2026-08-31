package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const configName = "ServerConfig.toml"

type MapEntry struct {
	Title  string `json:"title"`
	Path   string `json:"path"`
	Zip    string `json:"zip"`
	Source string `json:"source"`
	Active bool   `json:"active"`
	Img    string `json:"img"`
}

type Item struct {
	Title  string `json:"title"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
	Img    string `json:"img"`
}

type State struct {
	Dir            string     `json:"dir"`
	ConfigFound    bool       `json:"configFound"`
	ResourceFolder string     `json:"resourceFolder"`
	ServerName     string     `json:"serverName"`
	Port           int        `json:"port"`
	CurrentMap     string     `json:"currentMap"`
	Maps           []MapEntry `json:"maps"`
	Mods           []Item     `json:"mods"`
	Plugins        []Item     `json:"plugins"`
	LogFile        string     `json:"logFile"`
	Warnings       []Msg      `json:"warnings"`
}

type Manager struct{ dir string }

func (m *Manager) resourceFolder(cfg string) string {
	res := tomlString(cfg, "General", "ResourceFolder")
	if res == "" {
		res = "Resources"
	}
	res = strings.Trim(strings.ReplaceAll(res, "\\", "/"), "/")
	if st, err := os.Stat(filepath.Join(m.dir, filepath.FromSlash(res))); err != nil || !st.IsDir() {
		return "Resources"
	}
	return res
}

func (m *Manager) paths(cfg string) (client, deactMods, srv, deactPlugins, customMaps string) {
	res := filepath.Join(m.dir, filepath.FromSlash(m.resourceFolder(cfg)))
	return filepath.Join(res, "Client"),
		filepath.Join(res, "Client", "deactivated_mods"),
		filepath.Join(res, "Server"),
		filepath.Join(res, "Server", "deactivated_plugins"),
		filepath.Join(m.dir, "custom_maps")
}

func (m *Manager) readConfig() (string, error) {
	b, err := os.ReadFile(filepath.Join(m.dir, configName))
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (m *Manager) writeConfig(content string) error {
	dst := filepath.Join(m.dir, configName)
	tmp := dst + ".tmp"
	if err := os.WriteFile(tmp, []byte(content), 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, dst); err != nil {
		os.Remove(tmp)
		return err
	}
	return nil
}

type zipInfo struct {
	hasLevels bool
	mapPath   string
	mapTitle  string
	modTitle  string
	mapImage  string
	modImage  string
	imageFrom string
	images    []string
}

func inspectZip(file string) (zipInfo, error) {
	var zi zipInfo
	r, err := zip.OpenReader(file)
	if err != nil {
		return zi, err
	}
	defer r.Close()

	base := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	var levelInfo, modInfo, anyInfo, icon *zip.File
	misFound, misPath := false, ""
	for _, f := range r.File {
		name := strings.ReplaceAll(f.Name, "\\", "/")
		leaf := path.Base(name)
		inLevels := strings.Contains(name, "levels/")
		if inLevels {
			zi.hasLevels = true
		}
		switch {
		case leaf == "info.json" && inLevels:

			if levelInfo == nil || strings.Count(name, "/") < strings.Count(levelInfo.Name, "/") {
				levelInfo = f
			}
		case path.Ext(leaf) == ".mis" && inLevels:
			if !misFound || strings.Count(name, "/") < strings.Count(misPath, "/") {
				misPath = name
			}
			misFound = true
		}
		if leaf == "info.json" {
			anyInfo = f
			if strings.Contains(name, "mod_info/") {
				modInfo = f
			}
		}
		if icon == nil && (leaf == "icon.png" || leaf == "icon.jpg") {
			icon = f
		}
		if isImage(leaf) {
			zi.images = append(zi.images, f.Name)
		}
	}

	if levelInfo != nil {
		zi.mapPath = "/" + strings.ReplaceAll(levelInfo.Name, "\\", "/")
		if misFound {
			zi.mapPath = "/" + strings.ReplaceAll(misPath, "\\", "/")
		}
	}

	zi.mapTitle = base
	if t, _ := jsonInfo(modInfo); t != "" {
		zi.mapTitle = t
	}
	zi.modTitle = base
	if t, _ := jsonInfo(anyInfo); t != "" {
		zi.modTitle = t
	} else if t, _ := jsonInfo(modInfo); t != "" {
		zi.modTitle = t
	}

	if icon != nil {
		zi.modImage = icon.Name
	}
	if levelInfo != nil {
		_, previews := jsonInfo(levelInfo)
		zi.mapImage, zi.imageFrom = pickMapImage(zi.images, previews,
			path.Dir(strings.ReplaceAll(levelInfo.Name, "\\", "/"))+"/", icon)
	}
	return zi, nil
}

func isImage(leaf string) bool {
	switch strings.ToLower(path.Ext(leaf)) {
	case ".jpg", ".jpeg", ".png", ".webp", ".dds", ".bmp":
		return true
	}
	return false
}

func browserSafe(name string) bool {
	switch strings.ToLower(path.Ext(name)) {
	case ".jpg", ".jpeg", ".png", ".webp", ".bmp":
		return true
	}
	return false
}

func pickMapImage(images, previews []string, levelDir string, icon *zip.File) (string, string) {
	norm := func(s string) string { return strings.ToLower(strings.ReplaceAll(s, "\\", "/")) }

	for _, want := range previews {
		w := norm(want)
		if w == "" {
			continue
		}
		for _, img := range images {
			if norm(img) == strings.TrimSuffix(norm(levelDir), "/")+"/"+strings.TrimPrefix(w, "/") &&
				browserSafe(img) {
				return img, "previews"
			}
		}
		for _, img := range images {
			if norm(path.Base(img)) == norm(path.Base(w)) && browserSafe(img) {
				return img, "previews"
			}
		}
	}

	for _, img := range images {
		n := norm(img)
		if strings.HasPrefix(n, norm(levelDir)) && browserSafe(img) &&
			(strings.Contains(path.Base(n), "preview") || strings.Contains(path.Base(n), "screenshot")) {
			return img, "preview-name"
		}
	}

	if icon != nil && browserSafe(icon.Name) {
		return icon.Name, "icon"
	}

	for _, img := range images {
		n := norm(img)
		if browserSafe(img) && path.Dir(n)+"/" == norm(levelDir) {
			return img, "level-folder"
		}
	}
	return "", ""
}

func jsonInfo(f *zip.File) (string, []string) {
	if f == nil {
		return "", nil
	}
	rc, err := f.Open()
	if err != nil {
		return "", nil
	}
	defer rc.Close()
	b, err := io.ReadAll(io.LimitReader(rc, 1<<20))
	if err != nil {
		return "", nil
	}
	b = bytes.TrimPrefix(b, []byte{0xEF, 0xBB, 0xBF})

	var v struct {
		Title    string          `json:"title"`
		Previews json.RawMessage `json:"previews"`
	}
	if json.Unmarshal(b, &v) != nil {

		return scrapeInfo(string(b))
	}
	var previews []string
	if len(v.Previews) > 0 {
		if json.Unmarshal(v.Previews, &previews) != nil {
			var one string
			if json.Unmarshal(v.Previews, &one) == nil && one != "" {
				previews = []string{one}
			}
		}
	}
	return strings.TrimSpace(v.Title), previews
}

var (
	reTitle    = regexp.MustCompile(`(?i)"title"\s*:\s*"((?:[^"\\]|\\.)*)"`)
	rePreviews = regexp.MustCompile(`(?is)"previews"\s*:\s*(\[[^\]]*\]|"(?:[^"\\]|\\.)*")`)
	reString   = regexp.MustCompile(`"((?:[^"\\]|\\.)*)"`)
)

func scrapeInfo(s string) (string, []string) {
	title := ""
	if m := reTitle.FindStringSubmatch(s); m != nil {
		title = strings.TrimSpace(unescapeJSON(m[1]))
	}
	var previews []string
	if m := rePreviews.FindStringSubmatch(s); m != nil {
		for _, p := range reString.FindAllStringSubmatch(m[1], -1) {
			if v := unescapeJSON(p[1]); v != "" {
				previews = append(previews, v)
			}
		}
	}
	return title, previews
}

func unescapeJSON(s string) string {
	var out string
	if json.Unmarshal([]byte(`"`+s+`"`), &out) == nil {
		return out
	}
	return s
}

func prettify(s string) string {
	return strings.TrimSpace(strings.NewReplacer("_", " ", ".", " ").Replace(s))
}

func zipsIn(dir string) []string {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".zip") {
			out = append(out, filepath.Join(dir, e.Name()))
		}
	}
	return out
}

func (m *Manager) ensureFolders(cfg string) {
	client, deactMods, srv, deactPlugins, customMaps := m.paths(cfg)
	for _, d := range []string{client, srv, customMaps, deactMods, deactPlugins} {
		os.MkdirAll(d, 0o755)
	}
}

func (m *Manager) State() State {
	st := State{Dir: m.dir}
	cfg, err := m.readConfig()
	if err != nil {
		st.Warnings = append(st.Warnings, msg("config_missing", m.dir))
		return st
	}
	st.ConfigFound = true
	st.ResourceFolder = m.resourceFolder(cfg)
	st.ServerName = tomlString(cfg, "General", "Name")
	st.Port, _ = tomlInt(cfg, "General", "Port")
	st.CurrentMap = strings.TrimSpace(tomlString(cfg, "General", "Map"))
	if p := m.findLog(); p != "" {
		st.LogFile = filepath.Base(p)
	}
	m.ensureFolders(cfg)

	client, deactMods, srv, deactPlugins, customMaps := m.paths(cfg)

	for _, dir := range []string{client, customMaps} {
		source := "Client (aktiv)"
		if dir == customMaps {
			source = "custom_maps"
		}
		for _, f := range zipsIn(dir) {
			zi, err := inspectZip(f)
			if err != nil {
				st.Warnings = append(st.Warnings, msg("zip_unreadable", filepath.Base(f)))
				continue
			}
			if zi.mapPath == "" {
				continue
			}
			e := MapEntry{
				Title:  prettify(zi.mapTitle),
				Path:   zi.mapPath,
				Zip:    filepath.Base(f),
				Source: source,
			}
			if zi.mapImage != "" {
				e.Img = "/img?kind=map&zip=" + url.QueryEscape(e.Zip)
			}
			st.Maps = append(st.Maps, e)
		}
	}
	for i := range st.Maps {
		if st.CurrentMap != "" && st.Maps[i].Path == st.CurrentMap {
			st.Maps[i].Active = true
		}
	}

	for _, dir := range []string{client, deactMods} {
		active := dir == client
		for _, f := range zipsIn(dir) {
			zi, err := inspectZip(f)
			if err != nil {
				st.Warnings = append(st.Warnings, msg("zip_unreadable", filepath.Base(f)))
				continue
			}
			if zi.hasLevels {
				continue
			}
			it := Item{Title: prettify(zi.modTitle), Name: filepath.Base(f), Active: active}
			if zi.modImage != "" {
				it.Img = "/img?kind=mod&zip=" + url.QueryEscape(it.Name)
			}
			st.Mods = append(st.Mods, it)
		}
	}
	sort.Slice(st.Mods, func(i, j int) bool {
		return strings.ToLower(st.Mods[i].Title) < strings.ToLower(st.Mods[j].Title)
	})

	if entries, err := os.ReadDir(srv); err == nil {
		for _, e := range entries {
			if !e.IsDir() || e.Name() == "deactivated_plugins" || !hasLua(filepath.Join(srv, e.Name())) {
				continue
			}
			st.Plugins = append(st.Plugins, Item{Title: e.Name(), Name: e.Name(), Active: true})
		}
	}
	if entries, err := os.ReadDir(deactPlugins); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				st.Plugins = append(st.Plugins, Item{Title: e.Name(), Name: e.Name()})
			}
		}
	}
	sort.Slice(st.Plugins, func(i, j int) bool {
		return strings.ToLower(st.Plugins[i].Title) < strings.ToLower(st.Plugins[j].Title)
	})
	return st
}

func hasLua(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, e := range entries {
		if !e.IsDir() && strings.EqualFold(filepath.Ext(e.Name()), ".lua") {
			return true
		}
	}
	return false
}

func (m *Manager) SetMap(mapPath, zipName string) ([]Msg, error) {
	cfg, err := m.readConfig()
	if err != nil {
		return nil, errMsg("config_unreadable", err.Error())
	}
	m.ensureFolders(cfg)
	client, _, _, _, customMaps := m.paths(cfg)

	var (
		ops   []fileOp
		notes []Msg
	)

	fail := func(err error) ([]Msg, error) {
		undo(ops)
		return nil, err
	}

	for _, f := range zipsIn(client) {
		name := filepath.Base(f)
		if name == zipName {
			continue
		}
		zi, err := inspectZip(f)
		if err != nil || zi.mapPath == "" {
			continue
		}
		target := filepath.Join(customMaps, name)
		if _, err := os.Stat(target); err == nil {

			if sameFile(f, target) {
				if err := os.Remove(f); err != nil {
					return fail(errMsg("remove_failed", name, err.Error()))
				}
				continue
			}
			alt := uniquePath(customMaps, name)
			if alt == "" {
				return fail(errMsg("no_free_name", name))
			}
			if err := os.Rename(f, alt); err != nil {
				return fail(errMsg("map_park_failed", name, err.Error()))
			}
			ops = append(ops, fileOp{f, alt})
			notes = append(notes, msg("map_parked_renamed", name, filepath.Base(alt)))
			continue
		}
		if err := os.Rename(f, target); err != nil {
			return fail(errMsg("map_park_failed", name, err.Error()))
		}
		ops = append(ops, fileOp{f, target})
	}

	if zipName != "" {
		dst := filepath.Join(client, zipName)
		if _, err := os.Stat(dst); errors.Is(err, os.ErrNotExist) {
			src := filepath.Join(customMaps, zipName)
			if err := os.Rename(src, dst); err != nil {
				return fail(errMsg("map_activate_failed", zipName, err.Error()))
			}
			ops = append(ops, fileOp{src, dst})
		}
	}

	out, err := tomlSetString(cfg, "General", "Map", mapPath)
	if err != nil {
		return fail(err)
	}
	if err := backupOnce(filepath.Join(m.dir, configName)); err != nil {
		notes = append(notes, msg("backup_failed", err.Error()))
	}
	if err := m.writeConfig(out); err != nil {
		return fail(errMsg("config_write_failed", err.Error()))
	}
	return notes, nil
}

func (m *Manager) ToggleMod(name string, activate bool) error {
	if name != filepath.Base(name) || name == "" {
		return errMsg("invalid_name")
	}
	cfg, err := m.readConfig()
	if err != nil {
		return errMsg("config_unreadable", err.Error())
	}
	m.ensureFolders(cfg)
	client, deactMods, _, _, _ := m.paths(cfg)
	src, dst := filepath.Join(client, name), filepath.Join(deactMods, name)
	if activate {
		src, dst = dst, src
	}
	return move(src, dst, name)
}

func (m *Manager) TogglePlugin(name string, activate bool) error {
	if name != filepath.Base(name) || name == "" {
		return errMsg("invalid_name")
	}
	cfg, err := m.readConfig()
	if err != nil {
		return errMsg("config_unreadable", err.Error())
	}
	m.ensureFolders(cfg)
	_, _, srv, deactPlugins, _ := m.paths(cfg)
	src, dst := filepath.Join(srv, name), filepath.Join(deactPlugins, name)
	if activate {
		src, dst = dst, src
	}
	return move(src, dst, name)
}

func move(src, dst, name string) error {
	if _, err := os.Stat(dst); err == nil {
		return errMsg("duplicate", name)
	}
	if err := os.Rename(src, dst); err != nil {
		return errMsg("move_failed", name, err.Error())
	}
	return nil
}
