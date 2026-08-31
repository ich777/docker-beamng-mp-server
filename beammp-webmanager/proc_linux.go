package main

import (
	"os"
	"strconv"
	"strings"
	"syscall"
)

func listProcesses(want string) []ProcInfo {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil
	}
	self := os.Getpid()
	var out []ProcInfo
	for _, e := range entries {
		pid, err := strconv.Atoi(e.Name())
		if err != nil || pid == self {
			continue
		}
		comm, err := os.ReadFile("/proc/" + e.Name() + "/comm")
		if err != nil {
			continue
		}
		raw, _ := os.ReadFile("/proc/" + e.Name() + "/cmdline")
		args := splitCmdline(raw)
		exe, _ := os.Readlink("/proc/" + e.Name() + "/exe")
		out = append(out, newProcInfo(pid, strings.TrimSpace(string(comm)), exe, args, want))
	}
	return out
}

func splitCmdline(raw []byte) []string {
	parts := strings.Split(strings.TrimRight(string(raw), "\x00"), "\x00")
	if len(parts) == 1 && parts[0] == "" {
		return nil
	}
	return parts
}

func terminate(pid int) error {
	p, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return p.Signal(syscall.SIGTERM)
}
