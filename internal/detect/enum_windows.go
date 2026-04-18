//go:build windows

package detect

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

// Win32 API bindings for window enumeration.
var (
	user32 = windows.NewLazySystemDLL("user32.dll")
	dwmapi = windows.NewLazySystemDLL("dwmapi.dll")

	procEnumWindows              = user32.NewProc("EnumWindows")
	procGetWindowTextW           = user32.NewProc("GetWindowTextW")
	procGetWindowTextLengthW     = user32.NewProc("GetWindowTextLengthW")
	procIsWindowVisible          = user32.NewProc("IsWindowVisible")
	procGetWindowThreadProcessId = user32.NewProc("GetWindowThreadProcessId")
	procGetWindowLongW           = user32.NewProc("GetWindowLongW")
	procGetWindowDisplayAffinity = user32.NewProc("GetWindowDisplayAffinity")

	procDwmGetWindowAttribute = dwmapi.NewProc("DwmGetWindowAttribute")
)

const (
	GWL_EXSTYLE      = -20
	WS_EX_TOPMOST    = 0x00000008
	WS_EX_LAYERED    = 0x00080000
	WS_EX_TOOLWINDOW = 0x00000080

	DWMWA_CLOAKED = 14
)

type rawWindow struct {
	HWND     uintptr
	PID      uint32
	Title    string
	Topmost  bool
	Layered  bool
	Affinity uint32
	Cloaked  bool
}

// enumerateWindows returns every visible top-level window on the desktop,
// along with the flags that matter for copilot detection.
func enumerateWindows() []rawWindow {
	var out []rawWindow

	cb := syscall.NewCallback(func(hwnd uintptr, lparam uintptr) uintptr {
		// Filter invisible windows
		visible, _, _ := procIsWindowVisible.Call(hwnd)
		if visible == 0 {
			return 1
		}

		// Title
		length, _, _ := procGetWindowTextLengthW.Call(hwnd)
		title := ""
		if length > 0 {
			buf := make([]uint16, int(length)+1)
			procGetWindowTextW.Call(hwnd, uintptr(unsafe.Pointer(&buf[0])), length+1)
			title = windows.UTF16ToString(buf)
		}

		// Owning PID
		var pid uint32
		procGetWindowThreadProcessId.Call(hwnd, uintptr(unsafe.Pointer(&pid)))

		// Extended style (topmost, layered, toolwindow).
		// GWL_EXSTYLE is -20; a negative constant can't be converted to uintptr
		// at compile time, so route it through a runtime int32 variable.
		nIndex := int32(GWL_EXSTYLE)
		exStyle, _, _ := procGetWindowLongW.Call(hwnd, uintptr(nIndex))

		// Display affinity - the Tier 3 signal
		var affinity uint32
		procGetWindowDisplayAffinity.Call(hwnd, uintptr(unsafe.Pointer(&affinity)))

		// Cloaked state (hidden from Alt-Tab)
		var cloaked uint32
		procDwmGetWindowAttribute.Call(
			hwnd,
			uintptr(DWMWA_CLOAKED),
			uintptr(unsafe.Pointer(&cloaked)),
			unsafe.Sizeof(cloaked),
		)

		// Skip invisible toolwindows with no title (background helpers)
		if title == "" && (exStyle&WS_EX_TOOLWINDOW != 0) && cloaked == 0 {
			return 1
		}

		out = append(out, rawWindow{
			HWND:     hwnd,
			PID:      pid,
			Title:    title,
			Topmost:  exStyle&WS_EX_TOPMOST != 0,
			Layered:  exStyle&WS_EX_LAYERED != 0,
			Affinity: affinity,
			Cloaked:  cloaked != 0,
		})

		return 1
	})

	procEnumWindows.Call(cb, 0)
	return out
}
