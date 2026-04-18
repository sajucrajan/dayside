//go:build windows

package detect

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"time"
)

// powershellProcess is the JSON shape we parse from the Get-Process + Get-AuthenticodeSignature call.
type powershellProcess struct {
	PID      uint32 `json:"PID"`
	Name     string `json:"Name"`
	Path     string `json:"Path"`
	Signed   string `json:"Signed"`
	Signer   string `json:"Signer"`
}

// processMetadataScript runs once per Scan() and returns path + signature for
// every user-visible process. Historically we called Get-AuthenticodeSignature
// here, but that cmdlet walks the Windows certificate chain and performs
// OCSP/CRL network checks per file -- multiplied across ~100 processes, with
// any revocation endpoint slow or firewalled, the scan would stall for a
// minute or more. Instead we read the Authenticode blob directly from the PE
// file via X509Certificate.CreateFromSignedFile: no chain walk, no network,
// ~100x faster. Trade-off: we only detect "has an embedded signature at all"
// rather than "has a currently-valid signature" -- fine for our use case,
// which is flagging unsigned binaries as a yellow warning.
const processMetadataScript = `
$ErrorActionPreference = 'SilentlyContinue'
$results = Get-Process | Where-Object { $_.Path } | ForEach-Object {
    $signed = 'unsigned'
    $signer = ''
    try {
        $cert = [System.Security.Cryptography.X509Certificates.X509Certificate]::CreateFromSignedFile($_.Path)
        if ($cert) {
            $signed = 'signed'
            $signer = $cert.Subject
        }
    } catch {
        # CreateFromSignedFile throws if the file has no embedded signature.
        # That's the normal "unsigned" outcome, not an error.
        $signed = 'unsigned'
    }
    [PSCustomObject]@{
        PID    = $_.Id
        Name   = $_.ProcessName
        Path   = $_.Path
        Signed = $signed
        Signer = $signer
    }
}
$results | ConvertTo-Json -Depth 3 -Compress
`

// fetchProcessMetadata returns a PID -> metadata map by invoking PowerShell
// once. The command is wrapped in a 15-second context timeout so that a
// wedged PowerShell (antivirus interference, stuck cmdlet, etc.) can't hang
// the whole scan indefinitely -- the user sees a warning instead.
func fetchProcessMetadata() (map[uint32]powershellProcess, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "powershell.exe",
		"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass",
		"-Command", processMetadataScript)

	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		return nil, "process metadata timed out after 15s; results may be incomplete"
	}
	if err != nil {
		return nil, "process metadata: " + err.Error()
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return map[uint32]powershellProcess{}, ""
	}

	// PowerShell returns a single object instead of an array when there's only one
	// (rare in practice). Wrap it defensively.
	if !strings.HasPrefix(trimmed, "[") {
		trimmed = "[" + trimmed + "]"
	}

	var procs []powershellProcess
	if err := json.Unmarshal([]byte(trimmed), &procs); err != nil {
		return nil, "process metadata parse: " + err.Error()
	}

	m := make(map[uint32]powershellProcess, len(procs))
	for _, p := range procs {
		m[p.PID] = p
	}
	return m, ""
}

// EnumerateProcesses is the main entry point: combines window enumeration
// with process metadata and produces the UI-ready process list.
func EnumerateProcesses() ([]ProcessInfo, string) {
	wins := enumerateWindows()

	byPID := map[uint32][]rawWindow{}
	for _, w := range wins {
		byPID[w.PID] = append(byPID[w.PID], w)
	}

	meta, warn := fetchProcessMetadata()

	// Build one ProcessInfo per PID that either has metadata OR has at least one window.
	seen := map[uint32]bool{}
	var result []ProcessInfo

	addProc := func(pid uint32) {
		if seen[pid] {
			return
		}
		seen[pid] = true

		m, hasMeta := meta[pid]

		var name, path string
		if hasMeta {
			name = m.Name
			if !strings.HasSuffix(strings.ToLower(name), ".exe") {
				name += ".exe"
			}
			path = m.Path
		} else {
			// Fallback: we have windows but no PowerShell metadata.
			// Try to derive name from window info - but since we only have PID/HWND,
			// the name will be empty. Skip in that case.
			return
		}

		// Skip if no windows AND not known copilot AND protected (noise).
		winList := byPID[pid]
		if len(winList) == 0 && !isKnownCopilot(name) && isProtected(name) {
			return
		}

		// Skip pure background system processes with no windows that aren't known copilots.
		if len(winList) == 0 && !isKnownCopilot(name) && !isRemoteAccessTool(name) {
			return
		}

		var winEntries []WindowEntry
		for _, w := range winList {
			winEntries = append(winEntries, WindowEntry{
				HWND:     uint64(w.HWND),
				Title:    w.Title,
				Topmost:  w.Topmost,
				Layered:  w.Layered,
				Affinity: w.Affinity,
				Cloaked:  w.Cloaked,
			})
		}

		signed := "unknown"
		signer := ""
		if hasMeta {
			signed = m.Signed
			signer = m.Signer
		}

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

	// First pass: every PID we have metadata for.
	for pid := range meta {
		addProc(pid)
	}
	// Second pass: PIDs that have windows but somehow not in metadata (rare).
	for pid := range byPID {
		addProc(pid)
	}

	return result, warn
}
