package detect

import "strings"

// Protected processes are exempted from most noise-reducing checks (cloaked
// windows, unsigned binaries, layered overlays). They're either OS-critical,
// meeting clients the candidate must keep running, or this app itself.
// Used to suppress false positives during observation.
var protectedProcesses = []string{
	// OS-critical (killing these breaks the machine)
	"system", "registry", "smss.exe", "csrss.exe", "wininit.exe", "services.exe",
	"lsass.exe", "winlogon.exe", "fontdrvhost.exe", "svchost.exe", "dwm.exe",
	"spoolsv.exe", "taskhostw.exe", "explorer.exe", "sihost.exe", "runtimebroker.exe",
	"searchindexer.exe", "conhost.exe", "ctfmon.exe", "ntoskrnl.exe",
	// Meeting clients (must remain open during the interview)
	"teams.exe", "ms-teams.exe", "zoom.exe", "zoomruntime.exe",
	"webex.exe", "webexmta.exe", "googlemeet.exe",
	// Security tooling
	"msmpeng.exe", "nissrv.exe", "securityhealthservice.exe",
	// This app
	"interviewsweep.exe", "spectrecheck.exe", "dayside.exe",
}

// Known copilot / interview-cheating executables.
// Names appear in the process list even when window capture is hidden.
// Sources: vendor sites, GitHub repos, reverse-engineering articles, and
// reviews indexed in the v1.2 research dossier.
//
// Some entries are deliberately broad substrings (e.g., "cluely", "wizard")
// because vendors frequently rebrand or release variant builds. Match logic
// is substring-based and case-insensitive.
var knownCopilotNames = []string{
	// === Tier 3 - Native desktop overlay apps (the dangerous ones) ===
	// Cluely (Roy Lee / Chungin Lee) - the category leader
	"cluely",
	// Interview Coder / Interview Coder 2.0 (Roy Lee's earlier product)
	"interviewcoder", "interview-coder", "interview_coder",
	// Interview Solver
	"interviewsolver", "interview-solver", "interview_solver",
	// ShadeCoder (launched Dec 2025)
	"shadecoder", "shade-coder", "shade_coder",
	// LockedIn AI (Bright Data backend; "True Stealth Mode")
	"lockedin", "lockedin-ai", "lockedinai", "locked-in-ai",
	// Final Round AI (10M+ users; Stealth Mode in God-tier plans)
	"finalround", "final-round", "finalroundai", "final-round-ai",
	// UltraCode AI ($899 lifetime; OpenAI o3/o4-mini)
	"ultracode", "ultra-code", "ultracode-ai",
	// Leetcode Wizard (includes "humanizer" that mimics nervous typing)
	"leetcodewizard", "leetcode-wizard", "leetcode_wizard",
	"lc-wizard", "lcwizard",
	// Linkjob.ai (Chinese-market focus)
	"linkjob", "link-job", "linkjob-ai",
	// ParakeetAI (system-audio capture; voice-output variant)
	"parakeet", "parakeetai", "parakeet-ai", "parakeet_ai",
	// Sensei Copilot
	"sensei", "senseiai", "sensei-ai", "sensei-copilot", "senseicopilot",
	// Interview Pilot (mobile-only by design)
	"interview-pilot", "interviewpilot", "interview_pilot",
	// Verve AI / Verve Copilot
	"verve-ai", "vervecopilot", "verve-copilot",
	// MetaView
	"metaview", "meta-view",

	// === Open-source tools (BYOK; no central server to take down) ===
	// Glass (Pickle, YC; pixel-clone of Cluely launched 4 days after)
	"glass-app", "pickle-glass",
	// Pluely (Tauri/Rust; ~10MB; Linux support)
	"pluely",
	// Natively (claims 9k+ users; local RAG)
	"natively", "natively-ai", "natively-cluely",
	// free-cluely / OpenCluely / Cheating Daddy / others on GitHub
	"free-cluely", "freecluely", "opencluely", "open-cluely",
	"cheating-daddy", "cheatingdaddy",
	"mindwhisper", "mind-whisper",
	"phantomlens", "phantom-lens",

	// === Tier 1/2 - Browser-based and extension-based ===
	// (caught primarily via tab title/URL, but list process names where applicable)
	"interviewcopilot", "interview-copilot",
	"interviews-chat", "interviewschat",
	"interviewprep", "interview-prep",
	"ntro-io", "ntroio", "ntro.io",

	// === Generic giveaway substrings (NOT normal app names) ===
	// These catch deliberately-renamed builds that swap their exe name
	// to something innocuous-sounding but still containing common terms.
	"copilot-ai", "ai-copilot",
	"interview-ai", "ai-interview",
	"interview-helper", "interview-assistant",
	"coding-copilot", "coding-helper",
	"meeting-copilot", "meeting-assistant",
	"stealth-copilot", "stealth-assistant",

	// === Sometimes-renamed-to-blend-in patterns ===
	// LockedIn AI's docs literally suggest aliasing as "Ghost.exe";
	// flag obvious stealth-themed names that aren't standard processes.
	"ghost-app", "ghost-copilot",
	"shadow-ai", "shadowai",
	"invisible-ai", "invisibleai",
}

// Known copilot window title fragments (case-insensitive).
// Catches both browser-tab copilots (titles in tab strips) and any
// desktop tool that didn't bother to randomize its window title.
var knownCopilotTitles = []string{
	// === Brand names that appear in window/tab titles ===
	"cluely",
	"interview copilot", "interviewcopilot",
	"interview coder", "interviewcoder",
	"interview solver", "interviewsolver",
	"interview pilot", "interviewpilot",
	"interview helper", "interview assistant",
	"leetcode wizard", "leetcodewizard",
	"shadecoder", "shade coder",
	"linkjob", "link job",
	"parakeet",
	"sensei copilot", "sensei ai",
	"final round ai", "finalround", "final round",
	"lockedin ai", "lockedinai", "locked in ai",
	"ultracode", "ultra code",
	"verve copilot", "verve ai",
	"metaview", "meta view",
	"interviews.chat", "interviews chat",
	"ntro.io", "ntro",
	"interviewprep", "interview prep ai",

	// === Open-source clones ===
	"glass - cluely",
	"pluely",
	"natively",
	"free cluely", "free-cluely",
	"opencluely", "open cluely",
	"cheating daddy",
	"mindwhisper", "mind whisper",
	"phantom lens",

	// === Common copilot UI strings observed across multiple tools ===
	// Many copilots render answers with these labels in the window title
	// or tab title; flagging them catches generically-branded variants.
	"ai answer", "ai response",
	"interview answer", "interview response",
	"coding interview helper",
	"real-time interview",
	"undetectable interview", "undetectable copilot",
	"stealth interview", "stealth copilot",
	"hidden ai", "invisible ai",
}

// Virtual audio/video device name patterns. These are legitimate drivers,
// but their presence during an interview is suspicious because they're
// commonly used to pipe AI-generated audio or video into conference apps.
//
// Two specific concerns relevant to v1.2 research:
// 1. ParakeetAI generates SPOKEN answers via voice models - candidates may
//    use a virtual audio cable to route synthesized speech as their mic.
// 2. North Korean Famous Chollima actors use real-time deepfake face-swap
//    over virtual webcams (research shows working identities in 70 minutes
//    on a 5-year-old GPU).
var virtualDevicePatterns = []string{
	// === Virtual cameras (deepfake / proxy interview vector) ===
	"obs virtual camera", "obs-camera", "obs virtualcam",
	"snap camera", "snapcamera",
	"manycam",
	"xsplit vcam", "xsplit broadcaster", "xsplit",
	"droidcam", "droid-cam",
	"iriun webcam", "iriun",
	"epoccam", "epoc-cam",
	"nvidia broadcast camera", "nvidia broadcast",
	"ndi virtual", "ndi-camera",
	"reincubate camo", "camo studio",
	"e2esoft vcam", "vcam",
	"unity capture",
	"avatarify", "avatar-ify",
	"deepfacelive", "deepface-live", "deep-face-live",
	"facerig", "face-rig",
	"live3d", "live-3d",
	"webcamoid",
	"v4l2loopback", "v4l2-loopback",
	// Specific deepfake-software camera labels seen in DPRK investigations
	"reflect", "reface", "rope", "facefusion",

	// === Virtual audio (AI-voice routing / system-audio capture vector) ===
	"vb-audio", "vb-cable", "vb_cable", "vb audio",
	"voicemeeter",
	"virtual audio cable", "virtual-audio-cable",
	"cable input", "cable output",
	"vac", // VAC = Virtual Audio Cable; common shorthand
	"voicemod virtual", "voicemod",
	"clownfish",
	"morphvox",
	"obs-monitor", "obs monitor",
	"loopback audio", "loopback-audio", "rogue amoeba loopback",
	"blackhole", "blackhole-2ch", "blackhole-16ch",
	"soundflower",
	"jack audio", "jackaudio", "jack-audio",
	"pulseaudio virtual",
	"stereo mix",
	"resemble ai", "elevenlabs voice", // AI voice synthesis tools
	"xvc", "rvc", // realtime voice conversion tools
	"krisp virtual",
}

// Remote-access tool processes. Running any of these during an interview
// is highly suspicious — common pattern is a confederate connecting in
// remotely to type answers (DPRK laptop-farm pattern).
var remoteAccessTools = []string{
	// Mainstream remote-desktop
	"teamviewer.exe", "tv_w32.exe", "tv_x64.exe", "teamviewer_service.exe",
	"anydesk.exe", "anydesk-service.exe",
	"chrome_remote_desktop_host.exe", "remoting_host.exe",
	"parsec.exe", "parsecd.exe", "parsec-vdd.exe",
	"splashtop.exe", "srserver.exe", "srmanager.exe", "splashtopstreamer.exe",
	"logmein.exe", "lmiguardiansvc.exe", "logmeinsystray.exe",
	"dwservice.exe", "dwagent.exe", "dwagsvc.exe",
	"rustdesk.exe",
	"nomachine.exe", "nxnode.exe", "nxd.exe",
	"radmin.exe", "radmin3-server.exe",
	"goto.exe", "gotomypc.exe", "g2comm.exe",
	"connectwisecontrol.exe", "screenconnect.exe", "screenconnect.client.exe",
	"bomgar.exe", "beyondtrust.exe",
	"supremo.exe",

	// Windows Remote Desktop
	"mstsc.exe", "rdpclip.exe",

	// VNC variants
	"vncserver.exe", "winvnc.exe", "tvnserver.exe", "ultravnc.exe",
	"realvnc-server.exe", "vncviewer.exe", "tightvnc.exe",
	"x11vnc.exe",

	// Tunnel / proxy tools commonly used by DPRK laptop farms
	// (NB: ngrok and ssh have legitimate uses; flag yellow not red in scoring)
	"ngrok.exe",
	"frpc.exe", "frps.exe",
	"localtunnel.exe", "lt.exe",
	"cloudflared.exe",
	"tailscale.exe", "tailscale-ipn.exe",
	"zerotier-one.exe", "zerotier_one.exe",
}

// Canonical install locations for system-protected process names. A process
// claiming one of these names but running from some other directory is
// almost certainly an impersonator (a renamed copilot trying to blend in
// with OS infrastructure). The real binaries live in directories ACL'd to
// TrustedInstaller, so writing a fake into the expected path requires a
// privilege escalation the user would notice.
//
// Path prefixes are lowercase; the caller lowercases the candidate path
// before comparing. Windows-only - macOS uses codesign's team identifier
// for equivalent verification.
var expectedSystemPaths = map[string][]string{
	"svchost.exe":       {`c:\windows\system32\`, `c:\windows\syswow64\`},
	"explorer.exe":      {`c:\windows\`},
	"dwm.exe":           {`c:\windows\system32\`},
	"winlogon.exe":      {`c:\windows\system32\`},
	"lsass.exe":         {`c:\windows\system32\`},
	"services.exe":      {`c:\windows\system32\`},
	"csrss.exe":         {`c:\windows\system32\`},
	"smss.exe":          {`c:\windows\system32\`},
	"wininit.exe":       {`c:\windows\system32\`},
	"spoolsv.exe":       {`c:\windows\system32\`},
	"taskhostw.exe":     {`c:\windows\system32\`},
	"sihost.exe":        {`c:\windows\system32\`},
	"runtimebroker.exe": {`c:\windows\system32\`},
	"searchindexer.exe": {`c:\windows\system32\`},
	"conhost.exe":       {`c:\windows\system32\`},
	"ctfmon.exe":        {`c:\windows\system32\`},
	"fontdrvhost.exe":   {`c:\windows\system32\`},
}

// Allowlist of apps that legitimately use WDA_EXCLUDEFROMCAPTURE
// (or kCGWindowSharingNone on macOS). A window with capture-excluded
// state from one of these is not a copilot.
var affinityAllowlist = []string{
	// Password managers
	"1password.exe", "1password-7.exe",
	"bitwarden.exe",
	"keepass.exe", "keepassxc.exe",
	"lastpass.exe",
	"dashlane.exe",
	"enpass.exe",
	"keeper.exe", "keeperpasswordmanager.exe",
	"protonpass.exe",
	"nordpass.exe",
	// Banking / crypto clients
	"ledgerlive.exe", "ledger-live.exe",
	"trezor-suite.exe", "trezorsuite.exe",
	"exodus.exe",
	// DRM video / music players
	"netflix.exe",
	"spotify.exe",
	"appletv.exe", "tv.exe",
	"primevideo.exe",
	// Authenticator / 2FA apps
	"authy.exe", "authy-desktop.exe",
	"microsoftauthenticator.exe",
	// Windows itself
	"logonui.exe", "credentialuibroker.exe",
	"windowssecurityhealth.exe",
}

// ---------- helpers (unchanged from v1.1) ----------

func isProtected(name string) bool {
	n := strings.ToLower(name)
	for _, p := range protectedProcesses {
		if n == p {
			return true
		}
	}
	return false
}

func isKnownCopilot(name string) bool {
	n := strings.ToLower(name)
	n = strings.TrimSuffix(n, ".exe")
	for _, p := range knownCopilotNames {
		if strings.Contains(n, p) {
			return true
		}
	}
	return false
}

func titleMatchesCopilot(title string) bool {
	t := strings.ToLower(title)
	for _, p := range knownCopilotTitles {
		if strings.Contains(t, p) {
			return true
		}
	}
	return false
}

func isVirtualDevice(name string) (bool, string) {
	n := strings.ToLower(name)
	for _, p := range virtualDevicePatterns {
		if strings.Contains(n, p) {
			return true, "matches virtual device pattern: " + p
		}
	}
	return false, ""
}

func isRemoteAccessTool(name string) bool {
	n := strings.ToLower(name)
	for _, p := range remoteAccessTools {
		if n == p {
			return true
		}
	}
	return false
}

func isBrowser(name string) bool {
	n := strings.ToLower(name)
	for _, b := range browserNames {
		if n == b {
			return true
		}
	}
	return false
}

// isSystemImpersonation flags processes whose name matches a core Windows
// binary but whose on-disk path is outside the canonical system directory.
// Renaming a copilot to svchost.exe is documented tradecraft (LockedIn AI
// explicitly suggests aliasing as Ghost.exe); the path check is what
// catches it. Returns false when path is empty so that missing metadata
// doesn't produce spurious hits.
func isSystemImpersonation(name, path string) bool {
	if path == "" {
		return false
	}
	n := strings.ToLower(name)
	expected, ok := expectedSystemPaths[n]
	if !ok {
		return false
	}
	pl := strings.ToLower(path)
	for _, prefix := range expected {
		if strings.HasPrefix(pl, prefix) {
			return false
		}
	}
	return true
}

func isAffinityAllowlisted(name string) bool {
	n := strings.ToLower(name)
	for _, p := range affinityAllowlist {
		if n == p {
			return true
		}
	}
	return false
}

// browserNames is referenced by isBrowser.
var browserNames = []string{
	"chrome.exe", "msedge.exe", "firefox.exe", "opera.exe",
	"brave.exe", "vivaldi.exe", "iexplore.exe", "safari.exe",
	"arc.exe", "thorium.exe", "zen.exe", "librewolf.exe",
}
