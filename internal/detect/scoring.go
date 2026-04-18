package detect

import "strings"

// scoreProcess fills in Flags and Severity on a ProcessInfo based on everything
// we know about it. Severity is what drives UI color; Flags are the
// human-readable reasons shown in the expanded row.
func scoreProcess(p *ProcessInfo) {
	flags := []string{}

	// Critical: known copilot by name
	if isKnownCopilot(p.Name) {
		flags = append(flags, "KNOWN_COPILOT")
	}

	// Critical: remote access tool
	if isRemoteAccessTool(p.Name) {
		flags = append(flags, "REMOTE_ACCESS")
	}

	// Critical: process claiming a system-binary name but running from the
	// wrong location. When this fires we also bypass the usual "don't flag
	// protected processes" exemptions below - the whole point of the check
	// is that this one is wearing the name as a disguise.
	impersonator := isSystemImpersonation(p.Name, p.Path)
	if impersonator {
		flags = append(flags, "SYSTEM_IMPERSONATION")
	}

	// isProtected() normally suppresses noisy warnings on core OS processes.
	// An impersonator deliberately picks a protected name, so it must not
	// inherit that shield.
	nameIsProtected := isProtected(p.Name) && !impersonator

	// Critical: any window is hidden from screen capture (the Tier 3 signal)
	hasHiddenWindow := false
	hasTopmost := false
	hasLayered := false
	hasCloaked := false
	hasCopilotTitle := false

	for _, w := range p.Windows {
		if w.Affinity != WDA_NONE {
			hasHiddenWindow = true
		}
		if w.Topmost {
			hasTopmost = true
		}
		if w.Layered {
			hasLayered = true
		}
		if w.Cloaked {
			hasCloaked = true
		}
		if titleMatchesCopilot(w.Title) {
			hasCopilotTitle = true
		}
	}

	if hasHiddenWindow && !isAffinityAllowlisted(p.Name) {
		flags = append(flags, "HIDDEN_FROM_SCREEN_SHARE")
	}
	if hasCopilotTitle {
		flags = append(flags, "COPILOT_TITLE_MATCH")
	}
	// Cloaked windows are normal UWP lifecycle behavior for system apps
	// (Settings, Media Player, etc.). Only flag when the owner isn't a
	// protected/system process or is an impersonator.
	if hasCloaked && !isAffinityAllowlisted(p.Name) && !nameIsProtected {
		flags = append(flags, "CLOAKED_WINDOW")
	}

	// Warning: topmost + layered from non-system process (classic overlay pattern)
	if hasTopmost && hasLayered && !nameIsProtected && !isAffinityAllowlisted(p.Name) {
		flags = append(flags, "TRANSPARENT_OVERLAY")
	}

	// Warning: unsigned process
	if p.Signed == "unsigned" && !nameIsProtected {
		flags = append(flags, "UNSIGNED")
	}

	// Warning: unusual path (user temp, downloads, appdata roaming from unusual sources)
	if p.Path != "" {
		pl := strings.ToLower(p.Path)
		if strings.Contains(pl, `\temp\`) || strings.Contains(pl, `\downloads\`) {
			flags = append(flags, "RUNS_FROM_TEMP_OR_DOWNLOADS")
		}
	}

	p.Flags = flags
	p.Severity = severityFromFlags(flags)
}

func severityFromFlags(flags []string) string {
	redFlags := map[string]bool{
		"KNOWN_COPILOT":            true,
		"REMOTE_ACCESS":            true,
		"SYSTEM_IMPERSONATION":     true,
		"HIDDEN_FROM_SCREEN_SHARE": true,
		"COPILOT_TITLE_MATCH":      true,
		"CLOAKED_WINDOW":           true,
	}
	yellowFlags := map[string]bool{
		"TRANSPARENT_OVERLAY":         true,
		"UNSIGNED":                    true,
		"RUNS_FROM_TEMP_OR_DOWNLOADS": true,
	}

	for _, f := range flags {
		if redFlags[f] {
			return "red"
		}
	}
	for _, f := range flags {
		if yellowFlags[f] {
			return "yellow"
		}
	}
	return "green"
}
