package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

type ProcInfo struct {
	PID   int    `json:"pid"`
	Comm  string `json:"comm"`
	Exe   string `json:"exe"`
	Cmd   string `json:"cmd"`
	Match bool   `json:"match"`
	Why   string `json:"why"`
}

func findProcesses(name string) []ProcInfo {
	var out []ProcInfo
	for _, p := range listProcesses(name) {
		if p.Match {
			out = append(out, p)
		}
	}
	return out
}

func newProcInfo(pid int, comm, exe string, args []string, want string) ProcInfo {
	p := ProcInfo{PID: pid, Comm: comm, Exe: exe, Cmd: strings.Join(args, " ")}
	argv0 := ""
	if len(args) > 0 {
		argv0 = filepath.Base(strings.ReplaceAll(args[0], "\\", "/"))
	}
	switch {
	case nameMatches(comm, want):
		p.Match, p.Why = true, "comm"
	case argv0 != "" && nameMatches(argv0, want):
		p.Match, p.Why = true, "argv0"
	case exe != "" && nameMatches(filepath.Base(exe), want):
		p.Match, p.Why = true, "exe"
	}
	return p
}

func nameMatches(got, want string) bool {
	got = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(got)), ".exe")
	want = strings.TrimSuffix(strings.ToLower(strings.TrimSpace(want)), ".exe")
	if got == "" || want == "" {
		return false
	}
	return got == want || strings.HasPrefix(got, want) ||
		(len(got) == 15 && strings.HasPrefix(want, got))
}

func runPS(procName string) {
	all := listProcesses(procName)
	fmt.Println("Looking for:", procName)
	fmt.Println("Visible processes:", len(all))
	if len(all) <= 1 {
		fmt.Println("\nAlmost nothing visible - most likely a separate PID namespace.")
		fmt.Println("In a container use --pid=host, otherwise the tool only sees itself.")
	}
	var hits int
	fmt.Printf("\n%-8s %-18s %-11s %s\n", "PID", "COMM", "MATCH", "COMMAND LINE")
	for _, p := range all {
		mark := ""
		if p.Match {
			mark = "yes (" + p.Why + ")"
			hits++
		}
		cmd := p.Cmd
		if cmd == "" {
			cmd = p.Exe
		}
		if len(cmd) > 90 {
			cmd = cmd[:90] + "…"
		}
		fmt.Printf("%-8d %-18s %-11s %s\n", p.PID, trunc(p.Comm, 18), mark, cmd)
	}
	fmt.Printf("\n%d match(es).\n", hits)
	if hits == 0 {
		fmt.Println("No match. Take the right name from the COMM column and pass it as")
		fmt.Println("-restartsignal=NAME.")
	}
}

func trunc(s string, n int) string {
	if len(s) > n {
		return s[:n-1] + "…"
	}
	return s
}
