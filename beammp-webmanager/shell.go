package main

import (
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const defaultRestartCmd = "kill $(pidof BeamMP-Server)"

type CmdState struct {
	Enabled bool   `json:"enabled"`
	Running bool   `json:"running"`
	At      int64  `json:"at"`
	Code    int    `json:"code"`
	Out     string `json:"out"`
	Failed  bool   `json:"failed"`
}

func (a *app) cmdState() CmdState {
	cs := CmdState{Enabled: a.cmdOn, Running: a.cmdRunning}
	if !a.cmdAt.IsZero() {
		cs.At = a.cmdAt.Unix()
		cs.Code = a.cmdCode
		cs.Out = a.cmdOut
		cs.Failed = a.cmdCode != 0
	}
	return cs
}

func runShell(line string) (string, int) {
	out, err := exec.Command("/bin/sh", "-c", line).CombinedOutput()
	code := 0
	if ee, ok := err.(*exec.ExitError); ok {
		code = ee.ExitCode()
	} else if err != nil {
		return strings.TrimSpace(string(out)) + " " + err.Error(), -1
	}
	return strings.TrimSpace(string(out)), code
}

func (a *app) handleCmd(w http.ResponseWriter) {
	if !a.cmdOn {
		a.respond(w, errMsg("cmd_disabled"), "")
		return
	}
	if a.cmdRunning {
		a.respond(w, errMsg("cmd_running"), "")
		return
	}
	a.cmdRunning = true
	a.cmdAt = time.Time{}
	line := a.cmdLine

	a.respond(w, nil, "cmd_started", line)
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	go func() {
		out, code := runShell(line)
		a.mu.Lock()
		a.cmdRunning = false
		a.cmdAt = time.Now()
		a.cmdCode = code
		a.cmdOut = trunc(strings.TrimSpace(out), 2000)
		a.mu.Unlock()
	}()
}
