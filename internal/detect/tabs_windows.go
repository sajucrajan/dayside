//go:build windows

package detect

import (
	"encoding/json"
	"os/exec"
	"strings"
	"time"
)

// tabsScript uses .NET's UIAutomation framework to enumerate three things
// inside every running browser:
//
//  1. Regular tabs (TabItem controls) -- Kind="tab"
//  2. Embedded side panels (e.g. Chrome's "Ask Gemini" slide-out) by scanning
//     the top-level browser window for descendant Pane/Group/Document
//     elements whose Name matches known copilot panel labels. Kind="panel"
//  3. Pop-out AI chat windows -- top-level browser windows whose title
//     matches a copilot pattern. Kind="popout"
//
// Works across Chrome, Edge, Firefox, Brave, Opera, Vivaldi. Does not require
// any browser extension.
//
// Important: we do NOT iterate Get-Process and use $proc.MainWindowHandle
// because that returns only ONE window handle per process - and Chrome,
// Edge, etc. put Incognito / InPrivate windows and popped-out chat windows
// inside the same browser process as the normal window. Relying on
// MainWindowHandle would hide every private and pop-out window entirely.
//
// Instead we walk the UIAutomation desktop root and find every top-level
// Window element, filter to ones owned by a browser PID, then for each:
//   - enumerate TabItem descendants (regular tabs)
//   - if zero tabs and the window title matches a copilot, emit a popout row
//   - walk Pane/Group/Document descendants for Names that match side-panel
//     patterns (e.g. "Ask Gemini", "Copilot", "Claude for Chrome")
//
// Try/catch wraps the tight loops because UIA can throw COMException for
// windows in transient states - best-effort, not all-or-nothing.
const tabsScript = `
$ErrorActionPreference = 'SilentlyContinue'
try {
    Add-Type -AssemblyName UIAutomationClient
    Add-Type -AssemblyName UIAutomationTypes
} catch {
    Write-Output "[]"
    exit
}

$browsers = @('chrome','msedge','firefox','brave','opera','vivaldi','iexplore','arc')

# ClassName fragments (case-insensitive) that identify a browser side-panel
# container. Chrome/Edge Chromium-based browsers use ClassName='SidePanel'
# for the resizable slide-out pane that hosts Ask Gemini / Copilot / etc.
$sidePanelClassFragments = @('sidepanel')
$browserPidToName = @{}
foreach ($proc in Get-Process) {
    if ($browsers -contains $proc.ProcessName.ToLower()) {
        $browserPidToName[[int]$proc.Id] = $proc.ProcessName + ".exe"
    }
}
if ($browserPidToName.Count -eq 0) {
    Write-Output "[]"
    exit
}

# Private-browsing markers we recognize in a top-level browser window title
# OR in any descendant element's Name. Chrome/Edge expose an "Incognito" /
# "InPrivate" toggle button inside the window chrome whose accessible name
# contains the marker even when the window title is just the page title
# (which happens on some Chrome builds where the " - Incognito" suffix is
# only shown in the taskbar, not the UIA Name).
$privateMarkers = @(
    'incognito',
    'inprivate',
    'in-private',
    'in private',
    'private browsing',
    'private window',
    'private mode',
    '(private)',
    ' - private',
    'privat',
    'privado',
    'privee', 'privée',
    'privata',
    'anonim',
    'guest'
)

# AI-assistant markers covering the major providers and their in-browser
# embeddings. Used to filter BOTH side-panel and pop-out candidates so a
# benign side panel (Reading List, Bookmarks, History, Lens) or a benign
# browser pop-out (a PWA, a DevTools window) is not flagged.
#
# Design: match a short list of brand + product words in the panel/popout
# title. Structural detection (ClassName='SidePanel', zero tabs on a
# browser-owned top-level window) has already narrowed the candidate set,
# so these markers don't have to carry the whole false-positive load.
$aiChatMarkers = @(
    # Major LLM vendors (brand words)
    'gemini',
    'copilot',
    'chatgpt',
    'openai',
    'claude',
    'perplexity',
    'grok',
    'deepseek',
    'mistral',
    'le chat',

    # Generic AI-chat UI phrasing that shows up in panel/popout titles
    'ai chat',
    'ai assistant',
    'ai sidebar',
    'ask ai',
    'chat with',

    # Popular third-party AI side-panel extensions
    'monica',
    'sider',
    'merlin',
    'maxai',
    'harpa',
    'wiseone',
    'readergpt',
    'chathub',

    # Known interview-copilot products
    'cluely',
    'interview copilot',
    'interview coder',
    'leetcode wizard',
    'sensei copilot',
    'final round ai',
    'lockedin ai',
    'ultracode',
    'verve copilot',
    'metaview',
    'ntro.io',
    'parakeet ai'
)

function IsAIText($s) {
    if (-not $s) { return $false }
    $l = $s.ToLower()
    foreach ($m in $aiChatMarkers) {
        if ($l.Contains($m)) { return $true }
    }
    return $false
}

$results = @()

$tabCondition = New-Object System.Windows.Automation.PropertyCondition(
    [System.Windows.Automation.AutomationElement]::ControlTypeProperty,
    [System.Windows.Automation.ControlType]::TabItem
)
$windowCondition = New-Object System.Windows.Automation.PropertyCondition(
    [System.Windows.Automation.AutomationElement]::ControlTypeProperty,
    [System.Windows.Automation.ControlType]::Window
)
# Side-panel hosts are usually Pane or Group. We also include Document for
# the rare case where a panel is hosted as an embedded document element.
$paneCondition = New-Object System.Windows.Automation.PropertyCondition(
    [System.Windows.Automation.AutomationElement]::ControlTypeProperty,
    [System.Windows.Automation.ControlType]::Pane
)
$groupCondition = New-Object System.Windows.Automation.PropertyCondition(
    [System.Windows.Automation.AutomationElement]::ControlTypeProperty,
    [System.Windows.Automation.ControlType]::Group
)
$documentCondition = New-Object System.Windows.Automation.PropertyCondition(
    [System.Windows.Automation.AutomationElement]::ControlTypeProperty,
    [System.Windows.Automation.ControlType]::Document
)
$panelConditions = [System.Windows.Automation.Condition[]] @($paneCondition, $groupCondition, $documentCondition)
$panelOr = New-Object System.Windows.Automation.OrCondition($panelConditions)

$root = [System.Windows.Automation.AutomationElement]::RootElement
$topLevel = $root.FindAll([System.Windows.Automation.TreeScope]::Children, $windowCondition)

foreach ($win in $topLevel) {
    try {
        $ownerPid = [int]$win.Current.ProcessId
        if (-not $browserPidToName.ContainsKey($ownerPid)) { continue }

        $browserName = $browserPidToName[$ownerPid]
        $winTitle = $win.Current.Name
        $winTitleLower = ''
        if ($winTitle) { $winTitleLower = $winTitle.ToLower() }
        $incognito = $false
        foreach ($marker in $privateMarkers) {
            if ($winTitleLower.Contains($marker)) { $incognito = $true; break }
        }

        # Structural fallback: Chrome/Edge/Brave expose an "Incognito" or
        # "InPrivate" indicator button in the top-right of the window chrome.
        # Its accessible Name contains the marker even on builds where the
        # UIA window title is stripped of the suffix. We scan Button and
        # Pane descendants whose Name matches any private marker.
        if (-not $incognito) {
            try {
                $buttonCondition = New-Object System.Windows.Automation.PropertyCondition(
                    [System.Windows.Automation.AutomationElement]::ControlTypeProperty,
                    [System.Windows.Automation.ControlType]::Button
                )
                $btns = $win.FindAll([System.Windows.Automation.TreeScope]::Descendants, $buttonCondition)
                foreach ($b in $btns) {
                    try {
                        $bn = $b.Current.Name
                        if (-not $bn) { continue }
                        $bnl = $bn.ToLower()
                        foreach ($marker in $privateMarkers) {
                            if ($bnl.Contains($marker)) { $incognito = $true; break }
                        }
                        if ($incognito) { break }
                    } catch { }
                }
            } catch { }
        }

        $tabs = $win.FindAll([System.Windows.Automation.TreeScope]::Descendants, $tabCondition)
        $tabCount = 0
        if ($tabs) { $tabCount = $tabs.Count }

        # Collect tab names for later de-dup of side-panel candidates: the
        # content Pane inside a tab mirrors the tab title, so we skip any
        # panel whose Name matches a tab we already emitted from this
        # window. Otherwise every Claude-docs tab would also be reported
        # as a "panel".
        $tabNamesThisWin = @{}

        foreach ($tab in $tabs) {
            try {
                $name = $tab.Current.Name
                if ($name) {
                    $tabNamesThisWin[$name] = $true
                    $results += [PSCustomObject]@{
                        BrowserPID  = $ownerPid
                        BrowserName = $browserName
                        Title       = $name
                        URL         = ""
                        Incognito   = $incognito
                        Kind        = "tab"
                    }
                }
            } catch { }
        }

        # Pop-out detection: a browser-owned top-level window with no tab
        # strip whose title matches a copilot pattern. Typical Chrome
        # pop-out chat windows have zero TabItems because they're the
        # chrome-less "App" window mode.
        if ($tabCount -eq 0 -and (IsAIText $winTitle)) {
            $results += [PSCustomObject]@{
                BrowserPID  = $ownerPid
                BrowserName = $browserName
                Title       = $winTitle
                URL         = ""
                Incognito   = $incognito
                Kind        = "popout"
            }
        }

        # Side-panel detection: walk Pane/Group/Document descendants whose
        # Name matches a copilot pattern. We de-dup against the window
        # title so we don't double-emit when the side-panel's Name is just
        # the window title repeated.
        # Structural side-panel detection: Chrome/Edge host the AI slide-out
        # inside a Pane whose ClassName contains "SidePanel". Find those
        # containers, then read the Name of the embedded Document
        # (AutomationId=RootWebArea) to identify WHICH AI tool is active --
        # e.g. "Gemini Chrome :: New Conversation". This is far more reliable
        # than substring-matching on free-form UI text, which picked up
        # window frame titles and tab content names as false positives.
        try {
            $panes = $win.FindAll([System.Windows.Automation.TreeScope]::Descendants, $panelOr)
            $seenPanelNames = @{}
            foreach ($pane in $panes) {
                try {
                    $cls = $pane.Current.ClassName
                    if (-not $cls) { continue }
                    $clsLower = $cls.ToLower()
                    $isSidePanel = $false
                    foreach ($frag in $sidePanelClassFragments) {
                        if ($clsLower.Contains($frag)) { $isSidePanel = $true; break }
                    }
                    if (-not $isSidePanel) { continue }

                    # Found a side-panel container. Read the first
                    # descendant Document's Name (Chrome exposes the web
                    # content's <title> there) as the panel identity.
                    # ClassName fragment matches several helper panes too
                    # (SidePanelResizeArea, SidePanel::BorderView) -- those
                    # have no Document child, so requiring a named Document
                    # skips them cleanly and only real side-panel hosts
                    # make it through.
                    $panelTitle = $null
                    $docs = $pane.FindAll(
                        [System.Windows.Automation.TreeScope]::Descendants,
                        $documentCondition)
                    foreach ($doc in $docs) {
                        try {
                            $dn = $doc.Current.Name
                            if ($dn) { $panelTitle = $dn; break }
                        } catch { }
                    }
                    if (-not $panelTitle) { continue }

                    # Drop benign built-in side panels (Reading List,
                    # Bookmarks, History, Lens, Shopping insights, etc.).
                    # Only emit if the panel's Document title matches an
                    # AI/copilot marker. This keeps panel detection generic
                    # across every AI vendor while filtering out
                    # non-suspicious panels.
                    if (-not (IsAIText $panelTitle)) { continue }

                    if ($seenPanelNames.ContainsKey($panelTitle)) { continue }
                    $seenPanelNames[$panelTitle] = $true
                    $results += [PSCustomObject]@{
                        BrowserPID  = $ownerPid
                        BrowserName = $browserName
                        Title       = $panelTitle
                        URL         = ""
                        Incognito   = $incognito
                        Kind        = "panel"
                    }
                } catch { }
            }
        } catch { }
    } catch { }
}

if ($results.Count -eq 0) {
    Write-Output "[]"
} else {
    $results | ConvertTo-Json -Depth 3 -Compress
}
`

type tabRaw struct {
	BrowserPID  uint32 `json:"BrowserPID"`
	BrowserName string `json:"BrowserName"`
	Title       string `json:"Title"`
	URL         string `json:"URL"`
	Incognito   bool   `json:"Incognito"`
	Kind        string `json:"Kind"`
}

// EnumerateBrowserTabs returns all tabs, side-panels, and pop-out AI chat
// windows across all running browsers. On Windows this uses UIAutomation
// which can take 1-3 seconds for machines with many browser windows open.
func EnumerateBrowserTabs() ([]BrowserTab, string) {
	cmd := exec.Command("powershell.exe",
		"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass",
		"-Command", tabsScript)

	// UIA queries can hang on misbehaving windows - give it a hard ceiling.
	done := make(chan struct{}, 1)
	var out []byte
	var err error

	go func() {
		out, err = cmd.Output()
		done <- struct{}{}
	}()

	select {
	case <-done:
		// complete
	case <-time.After(8 * time.Second):
		_ = cmd.Process.Kill()
		return nil, "browser tab enumeration timed out after 8s"
	}

	if err != nil {
		return nil, "browser tabs: " + err.Error()
	}

	trimmed := strings.TrimSpace(string(out))
	if trimmed == "" || trimmed == "[]" {
		return []BrowserTab{}, ""
	}

	if !strings.HasPrefix(trimmed, "[") {
		trimmed = "[" + trimmed + "]"
	}

	var raw []tabRaw
	if err := json.Unmarshal([]byte(trimmed), &raw); err != nil {
		return nil, "browser tabs parse: " + err.Error()
	}

	tabs := make([]BrowserTab, 0, len(raw))
	for _, r := range raw {
		kind := r.Kind
		if kind == "" {
			kind = "tab"
		}
		tab := BrowserTab{
			BrowserPID:  r.BrowserPID,
			BrowserName: r.BrowserName,
			Title:       r.Title,
			URL:         r.URL,
			Incognito:   r.Incognito,
			Kind:        kind,
		}
		scoreTab(&tab)
		// Side-panels and pop-outs that matched copilot patterns are
		// always suspicious even if scoreTab didn't tag them red (e.g.
		// a panel called "Gemini" with no known domain in URL).
		if (tab.Kind == "panel" || tab.Kind == "popout") && tab.Severity != "red" {
			tab.Severity = "red"
			if tab.Kind == "panel" {
				tab.Reason = "AI side-panel in browser: " + tab.Title
			} else {
				tab.Reason = "Browser pop-out AI chat window: " + tab.Title
			}
		}
		tabs = append(tabs, tab)
	}
	return tabs, ""
}
