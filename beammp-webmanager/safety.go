package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type fileOp struct{ from, to string }

func undo(ops []fileOp) {
	for i := len(ops) - 1; i >= 0; i-- {
		os.Rename(ops[i].to, ops[i].from)
	}
}

func uniquePath(dir, name string) string {
	ext := filepath.Ext(name)
	stem := strings.TrimSuffix(name, ext)
	for i := 2; i < 1000; i++ {
		p := filepath.Join(dir, fmt.Sprintf("%s (%d)%s", stem, i, ext))
		if _, err := os.Stat(p); os.IsNotExist(err) {
			return p
		}
	}
	return ""
}

func sameFile(a, b string) bool {
	sa, err := os.Stat(a)
	if err != nil {
		return false
	}
	sb, err := os.Stat(b)
	if err != nil || sa.Size() != sb.Size() {
		return false
	}
	ha, err := hashFile(a)
	if err != nil {
		return false
	}
	hb, err := hashFile(b)
	return err == nil && ha == hb
}

func hashFile(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

var (
	backupMu   sync.Mutex
	backedUp   = map[string]bool{}
	backupSuff = ".bak"
)

func backupOnce(path string) error {
	backupMu.Lock()
	defer backupMu.Unlock()
	if backedUp[path] {
		return nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path+backupSuff, data, 0o644); err != nil {
		return err
	}
	backedUp[path] = true
	return nil
}
