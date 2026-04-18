//go:build darwin

package detect

/*
#cgo LDFLAGS: -framework CoreGraphics -framework CoreFoundation
#include <CoreGraphics/CoreGraphics.h>
#include <CoreFoundation/CoreFoundation.h>
#include <stdlib.h>
#include <string.h>

// Helpers to extract values from a CFDictionary (per-window info dict).

int cfDictGetInt(CFDictionaryRef dict, CFStringRef key) {
    CFNumberRef num = (CFNumberRef)CFDictionaryGetValue(dict, key);
    if (!num) return 0;
    int result = 0;
    CFNumberGetValue(num, kCFNumberIntType, &result);
    return result;
}

double cfDictGetDouble(CFDictionaryRef dict, CFStringRef key) {
    CFNumberRef num = (CFNumberRef)CFDictionaryGetValue(dict, key);
    if (!num) return 0.0;
    double result = 0.0;
    CFNumberGetValue(num, kCFNumberDoubleType, &result);
    return result;
}

// Returns 1 if the key exists and the string was copied into buf.
int cfDictGetString(CFDictionaryRef dict, CFStringRef key, char *buf, int bufSize) {
    CFStringRef str = (CFStringRef)CFDictionaryGetValue(dict, key);
    if (!str) { buf[0] = 0; return 0; }
    Boolean ok = CFStringGetCString(str, buf, bufSize, kCFStringEncodingUTF8);
    return ok ? 1 : 0;
}

// Static CFString keys - declared once, reused.
static CFStringRef kOwnerPID = CFSTR("kCGWindowOwnerPID");
static CFStringRef kOwnerName = CFSTR("kCGWindowOwnerName");
static CFStringRef kWindowName = CFSTR("kCGWindowName");
static CFStringRef kWindowLayer = CFSTR("kCGWindowLayer");
static CFStringRef kWindowAlpha = CFSTR("kCGWindowAlpha");
static CFStringRef kWindowSharingState = CFSTR("kCGWindowSharingState");
static CFStringRef kWindowNumber = CFSTR("kCGWindowNumber");

// Wrapper accessors exported to Go.
CFStringRef k_OwnerPID() { return kOwnerPID; }
CFStringRef k_OwnerName() { return kOwnerName; }
CFStringRef k_WindowName() { return kWindowName; }
CFStringRef k_WindowLayer() { return kWindowLayer; }
CFStringRef k_WindowAlpha() { return kWindowAlpha; }
CFStringRef k_WindowSharingState() { return kWindowSharingState; }
CFStringRef k_WindowNumber() { return kWindowNumber; }

CFArrayRef getWindowList() {
    return CGWindowListCopyWindowInfo(
        kCGWindowListOptionOnScreenOnly | kCGWindowListExcludeDesktopElements,
        kCGNullWindowID
    );
}

int arrayCount(CFArrayRef arr) {
    if (!arr) return 0;
    return (int)CFArrayGetCount(arr);
}

CFDictionaryRef arrayGet(CFArrayRef arr, int i) {
    return (CFDictionaryRef)CFArrayGetValueAtIndex(arr, i);
}

void releaseArray(CFArrayRef arr) {
    if (arr) CFRelease(arr);
}

int displayCount() {
    uint32_t count = 0;
    CGGetActiveDisplayList(0, NULL, &count);
    return (int)count;
}
*/
import "C"
import (
	"unsafe"
)

// macOS CGWindowSharingState values:
//   kCGWindowSharingNone      = 0  <- equivalent to Windows WDA_EXCLUDEFROMCAPTURE
//   kCGWindowSharingReadOnly  = 1  <- normal (visible in screen capture)
//   kCGWindowSharingReadWrite = 2  <- normal
const (
	sharingNone      = 0
	sharingReadOnly  = 1
	sharingReadWrite = 2
)

type rawWindow struct {
	HWND     uintptr
	PID      uint32
	Title    string
	Owner    string
	Topmost  bool
	Layered  bool
	Affinity uint32
	Cloaked  bool
}

// enumerateWindows returns on-screen windows with sharing state mapped to the
// same Affinity constants used by the Windows code, so scoring logic is shared.
func enumerateWindows() []rawWindow {
	arr := C.getWindowList()
	if arr == 0 {
		return nil
	}
	defer C.releaseArray(arr)

	count := int(C.arrayCount(arr))
	out := make([]rawWindow, 0, count)

	buf := make([]C.char, 1024)
	bufPtr := (*C.char)(unsafe.Pointer(&buf[0]))

	for i := 0; i < count; i++ {
		d := C.arrayGet(arr, C.int(i))
		if d == 0 {
			continue
		}

		pid := uint32(C.cfDictGetInt(d, C.k_OwnerPID()))
		layer := int(C.cfDictGetInt(d, C.k_WindowLayer()))
		alpha := float64(C.cfDictGetDouble(d, C.k_WindowAlpha()))
		sharing := int(C.cfDictGetInt(d, C.k_WindowSharingState()))
		windowNum := uint64(C.cfDictGetInt(d, C.k_WindowNumber()))

		C.cfDictGetString(d, C.k_WindowName(), bufPtr, 1024)
		title := C.GoString(bufPtr)

		C.cfDictGetString(d, C.k_OwnerName(), bufPtr, 1024)
		owner := C.GoString(bufPtr)

		// Skip layer-25 items like menu bar, dock, OS overlays with no title
		if title == "" && owner == "" {
			continue
		}

		var affinity uint32 = WDA_NONE
		if sharing == sharingNone {
			// Hidden from screen capture - the Tier 3 equivalent of WDA_EXCLUDEFROMCAPTURE
			affinity = WDA_EXCLUDEFROMCAPTURE
		}

		out = append(out, rawWindow{
			HWND:     uintptr(windowNum),
			PID:      pid,
			Title:    title,
			Owner:    owner,
			Topmost:  layer > 0,        // CG layer > 0 means floats above normal content
			Layered:  alpha < 0.95,     // semi-transparent
			Affinity: affinity,
			Cloaked:  false, // macOS has no exact equivalent
		})
	}

	return out
}

// countMonitors returns the number of active displays.
func countMonitors() int {
	return int(C.displayCount())
}
