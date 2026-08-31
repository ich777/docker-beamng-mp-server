package main

import (
	"net/http"
	"os/exec"
	"strings"
)

const (
	defaultTmuxSession = "BeamMP-Server"
	tmuxCols           = "500"
	tmuxRows           = "50"
)

type ConsoleState struct {
	Enabled bool   `json:"enabled"`
	Session string `json:"session"`
}

func (a *app) consoleState() ConsoleState {
	return ConsoleState{Enabled: a.tmuxOn, Session: a.tmuxSession}
}

func tmuxRun(args ...string) (string, error) {
	out, err := exec.Command("tmux", args...).CombinedOutput()
	return strings.TrimRight(string(out), " \t\r\n"), err
}

func (a *app) widenOnce() {
	a.mu.Lock()
	done := a.tmuxWide
	a.tmuxWide = true
	a.mu.Unlock()
	if done {
		return
	}
	tmuxRun("set-option", "-t", a.tmuxSession, "window-size", "manual")
	tmuxRun("resize-window", "-t", a.tmuxSession, "-x", tmuxCols, "-y", tmuxRows)
}

func (a *app) handleConsole(w http.ResponseWriter, r *http.Request) {
	if !a.tmuxOn {
		writeJSON(w, map[string]any{"error": msg("con_disabled")})
		return
	}

	if r.Method == http.MethodGet {
		a.widenOnce()
		out, err := tmuxRun("capture-pane", "-p", "-S", "-500", "-t", a.tmuxSession)
		if err != nil {
			a.mu.Lock()
			a.tmuxWide = false
			a.mu.Unlock()
			writeJSON(w, map[string]any{"error": msg("con_no_session", a.tmuxSession, out)})
			return
		}
		writeJSON(w, map[string]any{"text": out})
		return
	}

	var req struct {
		Line string `json:"line"`
	}
	if !a.decode(w, r, &req) {
		return
	}
	if req.Line != "" {
		if out, err := tmuxRun("send-keys", "-t", a.tmuxSession, "-l", "--", req.Line); err != nil {
			writeJSON(w, map[string]any{"error": msg("con_send_failed", out)})
			return
		}
	}
	if out, err := tmuxRun("send-keys", "-t", a.tmuxSession, "C-m"); err != nil {
		writeJSON(w, map[string]any{"error": msg("con_send_failed", out)})
		return
	}
	writeJSON(w, map[string]any{"ok": msg("con_sent")})
}
