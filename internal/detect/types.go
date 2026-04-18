package detect

// Window display-affinity constants. Windows defines these natively via
// GetWindowDisplayAffinity. On macOS we map kCGWindowSharingState into the
// same namespace so scoring logic doesn't need to know the platform.
const (
	WDA_NONE               uint32 = 0x00
	WDA_MONITOR            uint32 = 0x01
	WDA_EXCLUDEFROMCAPTURE uint32 = 0x11
)

// ProcessInfo is one row in the main process/window table.
type ProcessInfo struct {
	PID       uint32        `json:"pid"`
	Name      string        `json:"name"`
	Path      string        `json:"path"`
	Windows   []WindowEntry `json:"windows"`
	Signed    string        `json:"signed"` // "signed", "unsigned", "unknown"
	Signer    string        `json:"signer"`
	Protected bool          `json:"protected"`
	Severity  string        `json:"severity"` // "red", "yellow", "green"
	Flags     []string      `json:"flags"`
}

// WindowEntry is a single top-level window owned by a process.
// Affinity values: Windows uses WDA_NONE=0 / WDA_MONITOR=1 / WDA_EXCLUDEFROMCAPTURE=0x11.
// macOS translation: CGWindowSharing kCGWindowSharingNone=0, ReadOnly=1, ReadWrite=2.
// We use Affinity=0x11 to mean "hidden from capture" on both platforms.
type WindowEntry struct {
	HWND     uint64 `json:"hwnd"`
	Title    string `json:"title"`
	Topmost  bool   `json:"topmost"`
	Layered  bool   `json:"layered"`
	Affinity uint32 `json:"affinity"`
	Cloaked  bool   `json:"cloaked"`
}

// BrowserTab represents one tab, side-panel, or pop-out window discovered in
// a running browser. Kind distinguishes the three:
//   - "tab":    a regular browser tab (TabItem in UIA / tabs collection on macOS)
//   - "panel":  an embedded side panel inside a browser window (e.g. Chrome's
//               "Ask Gemini" slide-out). Detected by walking UIA descendants.
//   - "popout": a separate top-level window owned by the browser process that
//               has no tab strip (e.g. a popped-out Gemini / Copilot chat).
type BrowserTab struct {
	BrowserPID  uint32 `json:"browserPid"`
	BrowserName string `json:"browserName"` // "chrome.exe", "Google Chrome", etc.
	Title       string `json:"title"`
	URL         string `json:"url"`       // populated on macOS; empty on Windows (UIA doesn't expose URL)
	Incognito   bool   `json:"incognito"` // true for Chrome Incognito, Edge InPrivate, Firefox Private Browsing
	Kind        string `json:"kind"`      // "tab", "panel", "popout"
	Severity    string `json:"severity"`  // "red" if title/url matches a copilot, else "green"
	Reason      string `json:"reason"`    // why flagged
}

// DeviceReport covers audio and video devices.
type DeviceReport struct {
	Audio []DeviceEntry `json:"audio"`
	Video []DeviceEntry `json:"video"`
}

// DeviceEntry is one audio or video device.
type DeviceEntry struct {
	Name     string `json:"name"`
	Virtual  bool   `json:"virtual"`
	Severity string `json:"severity"`
	Reason   string `json:"reason"`
}

// SystemReport covers system-level signals.
type SystemReport struct {
	MonitorCount  int      `json:"monitorCount"`
	RemoteSession bool     `json:"remoteSession"`
	RemoteTools   []string `json:"remoteTools"`
	HostName      string   `json:"hostName"`
	UserName      string   `json:"userName"`
	OSVersion     string   `json:"osVersion"`
	Platform      string   `json:"platform"` // "windows" or "darwin"
}
