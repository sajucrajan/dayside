//go:build darwin

package detect

import (
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// EnumerateProcesses returns the full process list on macOS, combining
// window data (from CGWindowList via enum_darwin.go) with process paths
// from `ps` and code signatures from `codesign`.
func EnumerateProcesses() ([]ProcessInfo, string) {
	wins := enumerateWindows()
	byPID := map[uint32][]rawWindow{}
	for _, w := range wins {
		byPID[w.PID] = append(byPID[w.PID], w)
	}

	procs, warn := psProcessList()
	if warn != "" {
		// keep going with what we have
	}

	seen := map[uint32]bool{}
	var result []ProcessInfo

	add := func(pid uint32, name, path string) {
		if seen[pid] {
			return
		}
		seen[pid] = true

		if name == "" && path != "" {
			name = filepath.Base(path)
		}

		winList := byPID[pid]
		if len(winList) == 0 && !isKnownCopilot(name) && !isRemoteAccessTool(name) {
			// Skip pure background processes with no windows and no red flags
			return
		}

		var winEntries []WindowEntry
		for _, w := range winList {
			title := w.Title
			if title == "" && w.Owner != "" {
				title = "(" + w.Owner + ")"
			}
			winEntries = append(winEntries, WindowEntry{
				HWND:     uint64(w.HWND),
				Title:    title,
				Topmost:  w.Topmost,
				Layered:  w.Layered,
				Affinity: w.Affinity,
				Cloaked:  w.Cloaked,
			})
		}

		signed, signer := codesignCheck(path)

		info := ProcessInfo{
			PID:       pid,
			Name:      name,
			Path:      path,
			Windows:   winEntries,
			Signed:    signed,
			Signer:    signer,
			Protected: isProtected(name) && !isSystemImpersonation(name, path),
			Flags:     []string{},
		}
		scoreProcess(&info)
		result = append(result, info)
	}

	// First pass: processes from ps
	for _, p := range procs {
		add(p.pid, p.name, p.path)
	}
	// Second pass: windows that reference PIDs we didn't get from ps (rare)
	for pid, winList := range byPID {
		if seen[pid] || len(winList) == 0 {
			continue
		}
		add(pid, winList[0].Owner, "")
	}

	return result, warn
}

type psEntry struct {
	pid  uint32
	name string
	path string
}

// psProcessList returns every user-owned process with its executable path.
func psProcessList() ([]psEntry, string) {
	// -A: all users, -o: custom columns, -ww: no line truncation, comm=: full path.
	cmd := exec.Command("ps", "-Aww", "-o", "pid=,comm=")
	out, err := cmd.Output()
	if err != nil {
		return nil, "ps: " + err.Error()
	}

	var result []psEntry
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		pid64, err := strconv.ParseUint(strings.TrimSpace(parts[0]), 10, 32)
		if err != nil {
			continue
		}
		path := strings.TrimSpace(parts[1])
		if path == "" {
			continue
		}
		name := filepath.Base(path)
		result = append(result, psEntry{
			pid:  uint32(pid64),
			name: name,
			path: path,
		})
	}

	return result, ""
}

// codesignCheck returns ("signed", signer) / ("unsigned", "") / ("unknown", "").
// Shelling to codesign per-process is slow; we accept that - scan runs once.
func codesignCheck(path string) (string, string) {
	if path == "" {
		return "unknown", ""
	}
	// `codesign -dv --verbose=1 <path>` prints metadata to stderr.
	cmd := exec.Command("codesign", "-dv", "--verbose=1", path)
	out, err := cmd.CombinedOutput()
	text := string(out)
	if err != nil {
		if strings.Contains(text, "code object is not signed") ||
			strings.Contains(text, "not signed at all") {
			return "unsigned", ""
		}
		return "unknown", ""
	}
	// Look for "Authority=" lines; first one is the leaf certificate.
	signer := ""
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "Authority=") {
			signer = strings.TrimPrefix(line, "Authority=")
			break
		}
	}
	return "signed", signer
}
