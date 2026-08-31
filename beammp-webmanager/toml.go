package main

import (
	"fmt"
	"strconv"
	"strings"
)

func tomlGet(content, section, key string) (string, bool) {
	cur := ""
	for _, line := range strings.Split(content, "\n") {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "#") || t == "" {
			continue
		}
		if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
			cur = strings.TrimSpace(t[1 : len(t)-1])
			continue
		}
		if cur != section {
			continue
		}
		eq := strings.Index(t, "=")
		if eq < 0 {
			continue
		}
		if strings.TrimSpace(t[:eq]) != key {
			continue
		}
		val, _ := splitValue(strings.TrimSpace(t[eq+1:]))
		return val, true
	}
	return "", false
}

func tomlString(content, section, key string) string {
	raw, ok := tomlGet(content, section, key)
	if !ok {
		return ""
	}
	return unquote(raw)
}

func splitValue(rest string) (value, trailing string) {
	if strings.HasPrefix(rest, `"`) {
		esc := false
		for i := 1; i < len(rest); i++ {
			if esc {
				esc = false
				continue
			}
			switch rest[i] {
			case '\\':
				esc = true
			case '"':
				return rest[:i+1], rest[i+1:]
			}
		}
		return rest, ""
	}
	if h := strings.Index(rest, "#"); h >= 0 {
		return strings.TrimRight(rest[:h], " \t"), rest[h:]
	}
	return strings.TrimRight(rest, " \t\r"), ""
}

func unquote(raw string) string {
	if len(raw) >= 2 && strings.HasPrefix(raw, `"`) && strings.HasSuffix(raw, `"`) {
		if s, err := strconv.Unquote(raw); err == nil {
			return s
		}
		return raw[1 : len(raw)-1]
	}
	return raw
}

func quote(s string) string {
	r := strings.NewReplacer(`\`, `\\`, `"`, `\"`, "\n", `\n`, "\r", `\r`, "\t", `\t`)
	return `"` + r.Replace(s) + `"`
}

func tomlSetString(content, section, key, value string) (string, error) {
	lines := strings.Split(content, "\n")
	cur := ""
	sectionStart, sectionEnd := -1, -1
	for i, line := range lines {
		t := strings.TrimSpace(line)
		if strings.HasPrefix(t, "[") && strings.HasSuffix(t, "]") {
			if cur == section {
				sectionEnd = i
			}
			cur = strings.TrimSpace(t[1 : len(t)-1])
			if cur == section {
				sectionStart = i
			}
			continue
		}
		if cur != section || strings.HasPrefix(t, "#") || t == "" {
			continue
		}
		eq := strings.Index(t, "=")
		if eq < 0 || strings.TrimSpace(t[:eq]) != key {
			continue
		}
		indent := line[:len(line)-len(strings.TrimLeft(line, " \t"))]
		_, trailing := splitValue(strings.TrimSpace(t[eq+1:]))
		lines[i] = indent + key + " = " + quote(value) + trailing
		return strings.Join(lines, "\n"), nil
	}

	newLine := key + " = " + quote(value)
	switch {
	case sectionStart < 0:
		return content + "\n[" + section + "]\n" + newLine + "\n", nil
	case sectionEnd < 0:
		return strings.Join(append(lines, newLine), "\n"), nil
	default:
		out := append([]string{}, lines[:sectionEnd]...)
		out = append(out, newLine)
		out = append(out, lines[sectionEnd:]...)
		return strings.Join(out, "\n"), nil
	}
}

func tomlInt(content, section, key string) (int, error) {
	raw, ok := tomlGet(content, section, key)
	if !ok {
		return 0, fmt.Errorf("%s.%s not found", section, key)
	}
	return strconv.Atoi(strings.TrimSpace(raw))
}
