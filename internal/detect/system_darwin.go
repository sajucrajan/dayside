//go:build darwin

package detect

import (
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// EnumerateSystem produces the macOS system report.
func EnumerateSystem() (SystemReport, string) {
	report := SystemReport{
		MonitorCount:  countMonitors(),
		RemoteSession: isScreenSharingActive(),
		Platform:      runtime.GOOS,
	}

	if hn, err := os.Hostname(); err == nil {
		report.HostName = hn
	}
	if u := os.Getenv("USER"); u != "" {
		report.UserName = u
	}

	// OS version
	if out, err := exec.Command("sw_vers", "-productVersion").Output(); err == nil {
		report.OSVersion = "macOS " + strings.TrimSpace(string(out))
	} else {
		report.OSVersion = "macOS"
	}

	return report, ""
}

// isScreenSharingActive returns true if Apple's Screen Sharing (VNC-based)
// is currently serving a session. This is macOS's closest analog to an
// active RDP session.
func isScreenSharingActive() bool {
	out, err := exec.Command("launchctl", "list", "com.apple.screensharing").Output()
	if err != nil {
		return false
	}
	// If the service is loaded and has a PID, someone may be connected.
	// This is not perfectly accurate - launchctl shows loaded state, not
	// active session - but it's a useful signal.
	return strings.Contains(string(out), "\"PID\"")
}
