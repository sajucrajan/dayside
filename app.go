package main

import (
	"context"

	"dayside/internal/detect"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Window sizes for the two display modes. Compact mode is a small always-on-top
// window the user can leave visible during a screen-shared interview so the
// interviewer can see (and the candidate can glance at) the machine's status
// at any moment. Full mode is the normal application window.
const (
	compactWidth  = 360
	compactHeight = 360
	fullWidth     = 1040
	fullHeight    = 760
	edgeMargin    = 24
)

// App is the bridge between the frontend and the Go detection logic.
// Every exported method on this struct is callable from JavaScript.
// Dayside is an observation-only tool — it does not terminate processes
// or close windows. The candidate remediates detected items themselves.
type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// ScanResult is the full snapshot returned to the UI.
type ScanResult struct {
	Processes []detect.ProcessInfo `json:"processes"`
	Tabs      []detect.BrowserTab  `json:"tabs"`
	Devices   detect.DeviceReport  `json:"devices"`
	System    detect.SystemReport  `json:"system"`
	Warnings  []string             `json:"warnings"`
}

// Scan performs a full machine sweep and returns everything the UI needs.
func (a *App) Scan() ScanResult {
	result := ScanResult{}

	processes, warn := detect.EnumerateProcesses()
	result.Processes = processes
	if warn != "" {
		result.Warnings = append(result.Warnings, warn)
	}

	tabs, warn := detect.EnumerateBrowserTabs()
	result.Tabs = tabs
	if warn != "" {
		result.Warnings = append(result.Warnings, warn)
	}

	devices, warn := detect.EnumerateDevices()
	result.Devices = devices
	if warn != "" {
		result.Warnings = append(result.Warnings, warn)
	}

	system, warn := detect.EnumerateSystem()
	result.System = system
	if warn != "" {
		result.Warnings = append(result.Warnings, warn)
	}

	return result
}

// EnterCompactMode shrinks the window to a small always-on-top status window
// in the top-right corner of the current screen. Always-on-top is important
// because the use case is the candidate screen-sharing during an interview:
// a normal window would get hidden behind the code editor / meeting app, but
// always-on-top stays visible to the interviewer in a full-desktop share.
func (a *App) EnterCompactMode() {
	runtime.WindowSetAlwaysOnTop(a.ctx, true)
	runtime.WindowSetMinSize(a.ctx, compactWidth, compactHeight)
	runtime.WindowSetMaxSize(a.ctx, compactWidth, compactHeight)
	runtime.WindowSetSize(a.ctx, compactWidth, compactHeight)

	screens, err := runtime.ScreenGetAll(a.ctx)
	if err != nil || len(screens) == 0 {
		return
	}
	target := screens[0]
	for _, s := range screens {
		if s.IsCurrent {
			target = s
			break
		}
		if s.IsPrimary {
			target = s
		}
	}
	x := target.Width - compactWidth - edgeMargin
	if x < 0 {
		x = edgeMargin
	}
	runtime.WindowSetPosition(a.ctx, x, edgeMargin)
}

// ExitCompactMode restores the full-size window and turns off always-on-top.
func (a *App) ExitCompactMode() {
	runtime.WindowSetAlwaysOnTop(a.ctx, false)
	runtime.WindowSetMaxSize(a.ctx, 0, 0)
	runtime.WindowSetMinSize(a.ctx, fullWidth, fullHeight)
	runtime.WindowSetSize(a.ctx, fullWidth, fullHeight)
	runtime.WindowCenter(a.ctx)
}

// Quit cleanly shuts down the application. Called by the consent modal when
// the user declines — consent is always required before the main UI becomes
// interactive.
func (a *App) Quit() {
	runtime.Quit(a.ctx)
}
