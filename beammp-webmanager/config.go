package main

import (
	"net/http"
	"path/filepath"
	"strings"
)

func (a *app) handleConfig(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	if r.Method == http.MethodGet {
		text, err := a.mgr.readConfig()
		if err != nil {
			writeJSON(w, map[string]any{"error": msg("config_unreadable", err.Error())})
			return
		}
		writeJSON(w, map[string]any{"text": text, "file": configName})
		return
	}

	var req struct {
		Text string `json:"text"`
	}
	if !a.decode(w, r, &req) {
		return
	}
	if strings.TrimSpace(req.Text) == "" {
		a.respond(w, errMsg("config_empty"), "")
		return
	}
	var notes []Msg
	if !strings.Contains(req.Text, "[General]") {
		notes = append(notes, msg("config_no_general"))
	}
	path := filepath.Join(a.mgr.dir, configName)
	if err := backupOnce(path); err != nil {
		notes = append(notes, msg("backup_failed", err.Error()))
	}
	if err := a.mgr.writeConfig(req.Text); err != nil {
		a.respond(w, errMsg("config_write_failed", err.Error()), "")
		return
	}
	a.respondNotes(w, nil, notes, "config_saved")
}
