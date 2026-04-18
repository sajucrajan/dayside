//go:build darwin

package detect

import (
	"encoding/json"
	"os/exec"
	"strings"
)

// AppleScript returns tabs as a JSON array. Works for Chrome, Safari, Edge,
// Brave, Opera, Vivaldi because they all expose a roughly-compatible
// scripting dictionary. Firefox does NOT support AppleScript for tabs; we
// fall back to window titles for Firefox.
const tabsAppleScript = `
on escapeJson(s)
	set out to ""
	repeat with i from 1 to length of s
		set c to character i of s
		if c = "\"" then
			set out to out & "\\\""
		else if c = "\\" then
			set out to out & "\\\\"
		else if (ASCII number c) < 32 then
			set out to out & " "
		else
			set out to out & c
		end if
	end repeat
	return out
end escapeJson

on processBrowser(browserName, browserKey)
	set output to ""
	try
		tell application "System Events"
			if not (exists process browserKey) then return ""
		end tell
		tell application browserName
			repeat with w in windows
				try
					repeat with t in tabs of w
						set tabTitle to ""
						set tabURL to ""
						try
							set tabTitle to (title of t) as text
						end try
						try
							set tabURL to (URL of t) as text
						end try
						set output to output & "{\"browser\":\"" & browserName & "\",\"title\":\"" & my escapeJson(tabTitle) & "\",\"url\":\"" & my escapeJson(tabURL) & "\"},"
					end repeat
				end try
			end repeat
		end tell
	end try
	return output
end processBrowser

set allTabs to ""
set allTabs to allTabs & my processBrowser("Google Chrome", "Google Chrome")
set allTabs to allTabs & my processBrowser("Microsoft Edge", "Microsoft Edge")
set allTabs to allTabs & my processBrowser("Safari", "Safari")
set allTabs to allTabs & my processBrowser("Brave Browser", "Brave Browser")
set allTabs to allTabs & my processBrowser("Opera", "Opera")
set allTabs to allTabs & my processBrowser("Vivaldi", "Vivaldi")
set allTabs to allTabs & my processBrowser("Arc", "Arc")

if allTabs ends with "," then
	set allTabs to text 1 thru -2 of allTabs
end if
return "[" & allTabs & "]"
`

type tabRaw struct {
	Browser string `json:"browser"`
	Title   string `json:"title"`
	URL     string `json:"url"`
}

// EnumerateBrowserTabs returns all tabs across supported macOS browsers.
// First run may prompt the user for "System Events" / individual app
// automation permissions - this is a macOS privacy requirement.
func EnumerateBrowserTabs() ([]BrowserTab, string) {
	cmd := exec.Command("osascript", "-e", tabsAppleScript)
	out, err := cmd.Output()
	if err != nil {
		return nil, "browser tabs: " + err.Error()
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" || trimmed == "[]" {
		return []BrowserTab{}, ""
	}

	var raw []tabRaw
	if err := json.Unmarshal([]byte(trimmed), &raw); err != nil {
		return nil, "browser tabs parse: " + err.Error()
	}

	// Map browser names to PIDs via ps (best effort - PID isn't critical for display).
	pidByName := make(map[string]uint32)
	psOut, _ := exec.Command("ps", "-axo", "pid,comm").Output()
	for _, line := range strings.Split(string(psOut), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "PID") {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		if len(parts) != 2 {
			continue
		}
		var pid uint32
		if _, err := parseUint(parts[0], &pid); err != nil {
			continue
		}
		exePath := strings.TrimSpace(parts[1])
		base := exePath
		if idx := strings.LastIndex(exePath, "/"); idx >= 0 {
			base = exePath[idx+1:]
		}
		if _, exists := pidByName[base]; !exists {
			pidByName[base] = pid
		}
	}

	tabs := make([]BrowserTab, 0, len(raw))
	for _, r := range raw {
		tab := BrowserTab{
			BrowserName: r.Browser,
			Title:       r.Title,
			URL:         r.URL,
			BrowserPID:  pidByName[r.Browser],
			Kind:        "tab",
		}
		scoreTab(&tab)
		tabs = append(tabs, tab)
	}
	return tabs, ""
}

// parseUint is a tiny helper to avoid importing strconv in a hot path.
func parseUint(s string, out *uint32) (int, error) {
	var n uint32
	consumed := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			break
		}
		n = n*10 + uint32(c-'0')
		consumed++
	}
	*out = n
	return consumed, nil
}
