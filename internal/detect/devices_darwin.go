//go:build darwin

package detect

import (
	"encoding/json"
	"os/exec"
)

// system_profiler -json gives clean, stable output for device enumeration.

type spAudio struct {
	Audio []spAudioSection `json:"SPAudioDataType"`
}
type spAudioSection struct {
	Items []spAudioItem `json:"_items"`
}
type spAudioItem struct {
	Name     string        `json:"_name"`
	Items    []spAudioItem `json:"_items"`
	Input    string        `json:"coreaudio_input_source"`
	IsInput  string        `json:"coreaudio_device_input"`
	IsOutput string        `json:"coreaudio_device_output"`
}

type spCamera struct {
	Cameras []spCameraItem `json:"SPCameraDataType"`
}
type spCameraItem struct {
	Name string `json:"_name"`
}

// EnumerateDevices returns audio inputs and video devices for macOS.
func EnumerateDevices() (DeviceReport, string) {
	report := DeviceReport{Audio: []DeviceEntry{}, Video: []DeviceEntry{}}

	// Audio
	audioOut, err := exec.Command("system_profiler", "-json", "SPAudioDataType").Output()
	if err == nil {
		var parsed spAudio
		if json.Unmarshal(audioOut, &parsed) == nil {
			names := map[string]bool{}
			for _, sec := range parsed.Audio {
				collectAudioNames(sec.Items, names)
			}
			for n := range names {
				virtual, reason := isVirtualDevice(n)
				e := DeviceEntry{Name: n, Virtual: virtual}
				if virtual {
					e.Severity = "red"
					e.Reason = reason
				} else {
					e.Severity = "green"
				}
				report.Audio = append(report.Audio, e)
			}
		}
	}

	// Video / cameras
	videoOut, err := exec.Command("system_profiler", "-json", "SPCameraDataType").Output()
	if err == nil {
		var parsed spCamera
		if json.Unmarshal(videoOut, &parsed) == nil {
			for _, c := range parsed.Cameras {
				virtual, reason := isVirtualDevice(c.Name)
				e := DeviceEntry{Name: c.Name, Virtual: virtual}
				if virtual {
					e.Severity = "red"
					e.Reason = reason
				} else {
					e.Severity = "green"
				}
				report.Video = append(report.Video, e)
			}
		}
	}

	return report, ""
}

func collectAudioNames(items []spAudioItem, out map[string]bool) {
	for _, it := range items {
		if it.Name != "" {
			out[it.Name] = true
		}
		if len(it.Items) > 0 {
			collectAudioNames(it.Items, out)
		}
	}
}
