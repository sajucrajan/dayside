//go:build windows

package detect

import (
	"encoding/json"
	"os/exec"
	"strings"
)

// deviceScript enumerates audio and video capture devices. We use
// Get-PnpDevice which covers both audio (Class=AudioEndpoint, MEDIA) and
// video (Class=Camera, Image) categories in a single call.
const deviceScript = `
$ErrorActionPreference = 'SilentlyContinue'
$audio = Get-PnpDevice -Class AudioEndpoint,MEDIA -PresentOnly -Status OK |
    Where-Object { $_.FriendlyName } |
    ForEach-Object {
        [PSCustomObject]@{
            Kind = 'audio'
            Name = $_.FriendlyName
        }
    }
$video = Get-PnpDevice -Class Camera,Image -PresentOnly -Status OK |
    Where-Object { $_.FriendlyName } |
    ForEach-Object {
        [PSCustomObject]@{
            Kind = 'video'
            Name = $_.FriendlyName
        }
    }
$all = @()
if ($audio) { $all += $audio }
if ($video) { $all += $video }
$all | ConvertTo-Json -Depth 2 -Compress
`

type deviceRaw struct {
	Kind string `json:"Kind"`
	Name string `json:"Name"`
}

// EnumerateDevices returns audio and video devices with virtual-device flagging.
func EnumerateDevices() (DeviceReport, string) {
	report := DeviceReport{
		Audio: []DeviceEntry{},
		Video: []DeviceEntry{},
	}

	cmd := exec.Command("powershell.exe",
		"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass",
		"-Command", deviceScript)

	out, err := cmd.Output()
	if err != nil {
		return report, "device enumeration: " + err.Error()
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" {
		return report, ""
	}

	if !strings.HasPrefix(trimmed, "[") {
		trimmed = "[" + trimmed + "]"
	}

	var devices []deviceRaw
	if err := json.Unmarshal([]byte(trimmed), &devices); err != nil {
		return report, "device parse: " + err.Error()
	}

	for _, d := range devices {
		virtual, reason := isVirtualDevice(d.Name)
		entry := DeviceEntry{
			Name:    d.Name,
			Virtual: virtual,
		}
		if virtual {
			entry.Severity = "red"
			entry.Reason = reason
		} else {
			entry.Severity = "green"
		}

		switch d.Kind {
		case "audio":
			report.Audio = append(report.Audio, entry)
		case "video":
			report.Video = append(report.Video, entry)
		}
	}

	return report, ""
}
