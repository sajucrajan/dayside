//go:build windows

package detect

import (
	"os"
	"syscall"
)

var (
	procEnumDisplayMonitors = user32.NewProc("EnumDisplayMonitors")
	procGetSystemMetrics    = user32.NewProc("GetSystemMetrics")
)

const (
	// GetSystemMetrics index for remote session.
	SM_REMOTESESSION = 0x1000
)

// countMonitors returns the number of attached monitors.
func countMonitors() int {
	count := 0
	cb := syscall.NewCallback(func(hMonitor, hdc, rect, data uintptr) uintptr {
		count++
		return 1
	})
	procEnumDisplayMonitors.Call(0, 0, cb, 0)
	return count
}

// isRemoteSession returns true if the process is running inside RDP.
func isRemoteSession() bool {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(SM_REMOTESESSION))
	return ret != 0
}

// EnumerateSystem produces the system-level report.
func EnumerateSystem() (SystemReport, string) {
	report := SystemReport{
		MonitorCount:  countMonitors(),
		RemoteSession: isRemoteSession(),
		Platform:      "windows",
	}

	if hn, err := os.Hostname(); err == nil {
		report.HostName = hn
	}

	if u := os.Getenv("USERNAME"); u != "" {
		report.UserName = u
	}

	osName := "Windows"
	if v := os.Getenv("OS"); v != "" {
		osName = v
	}
	report.OSVersion = osName

	return report, ""
}
