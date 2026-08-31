package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	defaultRestartFile = "restart"
	defaultProcName    = "BeamMP-Server"

	restartGrace = 20 * time.Second

	signalTimeout = 15 * time.Second
)

type RestartState struct {
	Enabled bool   `json:"enabled"`
	File    string `json:"file"`
	Pending bool   `json:"pending"`
	DoneAt  int64  `json:"doneAt"`

	Signal     bool       `json:"signal"`
	ProcName   string     `json:"procName"`
	Procs      []ProcInfo `json:"procs"`
	SigPending bool       `json:"sigPending"`
	SigFailed  bool       `json:"sigFailed"`
	SigStopped bool       `json:"sigStopped"`
	SigDoneAt  int64      `json:"sigDoneAt"`
	SigNewPID  int        `json:"sigNewPid"`
}

func (a *app) restartFile() string {
	if a.restartPath != "" {
		return a.restartPath
	}
	return filepath.Join(a.mgr.dir, defaultRestartFile)
}

func (a *app) restartState() RestartState {
	rs := RestartState{Enabled: a.restartOn, Signal: a.sigOn}
	if a.sigOn {
		rs.ProcName = a.procName
		rs.Procs = findProcesses(a.procName)
		if a.sigPending && !anyPID(rs.Procs, a.sigOldPIDs) {

			if a.sigGoneAt.IsZero() {
				a.sigGoneAt = time.Now()
			}
			switch {
			case len(rs.Procs) > 0:
				a.sigPending, a.sigNewPID, a.sigDoneAt = false, rs.Procs[0].PID, time.Now()
			case time.Since(a.sigGoneAt) > restartGrace:
				a.sigPending, a.sigNewPID, a.sigDoneAt = false, 0, a.sigGoneAt
			}
		}

		if a.sigPending && a.sigGoneAt.IsZero() && time.Since(a.sigSentAt) > signalTimeout {
			a.sigPending, a.sigFailed = false, true
		}
		rs.SigPending = a.sigPending
		rs.SigFailed = a.sigFailed
		rs.SigStopped = a.sigPending && !a.sigGoneAt.IsZero()
		rs.SigNewPID = a.sigNewPID
		if !a.sigDoneAt.IsZero() {
			rs.SigDoneAt = a.sigDoneAt.Unix()
		}
	}
	if !a.restartOn {
		return rs
	}
	rs.File = a.restartFile()
	_, err := os.Stat(rs.File)
	exists := err == nil
	if a.restartPending && !exists {
		a.restartPending = false
		a.restartDoneAt = time.Now()
	}
	rs.Pending = a.restartPending && exists
	if !a.restartDoneAt.IsZero() {
		rs.DoneAt = a.restartDoneAt.Unix()
	}
	return rs
}

func anyPID(procs []ProcInfo, pids []int) bool {
	for _, p := range procs {
		for _, old := range pids {
			if p.PID == old {
				return true
			}
		}
	}
	return false
}

func (a *app) handleRestart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Cancel bool `json:"cancel"`
		Signal bool `json:"signal"`
		Cmd    bool `json:"cmd"`
	}
	if !a.decode(w, r, &req) {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	if req.Cmd {
		a.handleCmd(w)
		return
	}
	if req.Signal {
		a.signalServer(w)
		return
	}

	if !a.restartOn {
		a.respond(w, errMsg("restart_disabled"), "")
		return
	}
	file := a.restartFile()

	if req.Cancel {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			a.respond(w, errMsg("restart_cancel_failed", err.Error()), "")
			return
		}
		a.restartPending = false
		a.respond(w, nil, "restart_cancelled")
		return
	}

	if _, err := os.Stat(file); err == nil {
		a.restartPending = true
		a.respond(w, nil, "restart_already_pending")
		return
	}
	body := "restart requested by beammp-tool at " + time.Now().Format(time.RFC3339) + "\n"
	if err := os.WriteFile(file, []byte(body), 0o644); err != nil {
		a.respond(w, errMsg("restart_failed", err.Error()), "")
		return
	}
	a.restartPending = true
	a.restartDoneAt = time.Time{}
	a.respond(w, nil, "restart_requested", filepath.Base(file))
}

func (a *app) signalServer(w http.ResponseWriter) {
	if !a.sigOn {
		a.respond(w, errMsg("signal_disabled"), "")
		return
	}
	procs := findProcesses(a.procName)
	if len(procs) == 0 {
		a.respond(w, errMsg("proc_none", a.procName), "")
		return
	}

	var pids []int
	for _, p := range procs {
		pids = append(pids, p.PID)
	}
	a.sigOldPIDs = pids
	a.sigPending = true
	a.sigNewPID = 0
	a.sigDoneAt = time.Time{}
	a.sigGoneAt = time.Time{}
	a.sigSentAt = time.Now()
	a.sigFailed = false
	a.respond(w, nil, "proc_signaled", a.procName, joinPIDs(pids))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	go func() {
		time.Sleep(300 * time.Millisecond)
		for _, pid := range pids {
			terminate(pid)
		}
	}()
}

func joinPIDs(pids []int) string {
	parts := make([]string, len(pids))
	for i, p := range pids {
		parts[i] = strconv.Itoa(p)
	}
	return strings.Join(parts, ", ")
}
