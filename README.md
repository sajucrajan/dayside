# Dayside

![Dayside — Candidate-operated pre-interview machine check](docs/hero.png)

> A candidate-operated pre-interview machine check. See what's running, what's hidden, and what doesn't belong — before and during the interview.

[![Live Site](https://img.shields.io/badge/Live%20Site-dayside-blue?style=flat-square&logo=github)](https://sajucrajan.github.io/dayside/)
[![Releases](https://img.shields.io/github/v/release/sajucrajan/dayside?style=flat-square&label=Download)](https://github.com/sajucrajan/dayside/releases/latest)
[![License](https://img.shields.io/badge/license-Apache%202.0-green?style=flat-square)](LICENSE)
[![Platform](https://img.shields.io/badge/platform-Windows%20%7C%20macOS-lightgrey?style=flat-square)](#getting-it)

<p align="center">
  <a href="https://sajucrajan.github.io/dayside/"><strong>Live Site</strong></a> &nbsp;·&nbsp;
  <a href="https://github.com/sajucrajan/dayside/releases/latest"><strong>Download</strong></a> &nbsp;·&nbsp;
  <a href="https://sajucrajan.github.io/dayside/employers.html"><strong>For Employers</strong></a> &nbsp;·&nbsp;
  <a href="https://sajucrajan.github.io/dayside/privacy.html"><strong>Privacy</strong></a> &nbsp;·&nbsp;
  <a href="LICENSE"><strong>License</strong></a>
</p>

---

## Table of Contents

- [What this is](#what-this-is)
- [How it works](#how-it-works)
- [What it checks](#what-it-checks)
- [What it does not do](#what-it-does-not-do)
- [Getting it](#getting-it)
- [For Employers](#for-interviewers-considering-this-tool)
- [Blocklist contributions](#blocklist-contributions)
- [Privacy and consent](#privacy-and-consent)
- [Security](#security)
- [License](#license)

---

## What this is

Remote interviews have a new problem. A generation of invisible AI interview assistants now renders answers in a window that your screen-share cannot see. The candidate reads them live. The interviewer sees nothing. For interviewers, the camera view and the Zoom tile are no longer enough.

Dayside is a small, local, open tool the candidate runs before the interview begins. It enumerates what's running on the machine, flags known AI interview assistants and stealth patterns, and shows the results on the candidate's own screen. The candidate shares that screen with the interviewer, and they discuss anything flagged together. Dayside stays visible during the interview — in its compact pill mode — so the interviewer can see at a glance that the scan is current and clean.

It's intentionally simple. Explicit consent. An explicit scan button. Optional periodic re-scans if the candidate opts in. No cloud. No recording. No account. **The tool only observes** — it never terminates processes or modifies the candidate's machine. When action is needed, the candidate takes it themselves through their operating system's normal tools (Task Manager on Windows, Activity Monitor on macOS).

> The name "Dayside" comes from the astronomical term for the lit side of a planet — the part where everything is visible. That's the posture of the tool: everything in daylight, nothing hidden.

---

## How it works

1. **Candidate downloads and runs Dayside.** A consent screen appears explaining exactly what the tool will do. The candidate either agrees or declines. Declining closes the app.

2. **Main window opens in "Ready" state.** Nothing scans yet. The candidate clicks **Start scan** when they're ready — typically at the start of the interview, during the screen-share.

3. **One scan runs (~10 seconds).** Results appear on the candidate's screen. The candidate and interviewer discuss anything flagged.

4. **Optional: the candidate can enable "Keep checking during interview"** — a checkbox that, if ticked, re-runs the scan every 5 minutes to catch anything launched mid-interview. **This is off by default.** The candidate chooses whether to enable it.

5. **Compact mode (pill) is available.** The main window can shrink to a small always-on-top status pill that stays visible during the interview without interfering with coding or reading. It shows the current scan status and time of the last scan.

6. **When the interview ends, the candidate closes the app.** Nothing persists. No data leaves the machine. Ever.

---

## What it checks

| Category | Details |
|---|---|
| **Running processes** | Names and signatures matched against a maintained blocklist of known AI interview assistants and stealth tools |
| **Visible & cloaked windows** | Including windows that set themselves "excluded from screen capture" |
| **Virtual audio/video devices** | Virtual cameras (deepfake vectors) and virtual audio cables (voice-routing vectors) |
| **Remote-access & tunneling tools** | TeamViewer, AnyDesk, Parsec, RustDesk, ngrok, cloudflared, tailscale, etc. |
| **Browser tabs** | Titles on Windows, titles plus URLs on macOS — checked against known AI-assistant domains and LLM chat sites |
| **System state** | Monitor count, RDP/screen-sharing session status, architecture |

---

## What it does not do

- **Nothing is uploaded.** Everything runs locally. No network calls. No telemetry.
- **No recording.** Screen, camera, and microphone are not accessed.
- **No biometrics.** No face scan, no voiceprint, no keystroke capture.
- **No persistent data.** When the app closes, nothing remains.
- **No termination of processes.** Dayside is read-only — it shows you what's running, but you close things yourself through your OS's own task manager.
- **No automatic actions.** Every decision is the candidate's: when to scan, whether to enable periodic re-scans, whether to close flagged programs.
- **No hiding.** Dayside itself is always visible to the interviewer during the screen-share. It does not exclude itself from screen capture. It does not run in the system tray.

> This is a blocklist and heuristic scanner. It is not a proctoring tool. It is not a background check. It will not catch a second device sitting off-camera. It will not catch kernel-level tools designed to evade detection. It is one layer in a thoughtful interview process, not a silver bullet.

---

## Getting it

Pre-built binaries for Windows and macOS are published on the [Releases page](../../releases). Download, extract, run. No install, no server, no account.

To build from source, see [BUILD.md](BUILD.md). Short version: install Go 1.22+, then run `./build.ps1` (Windows) or `./build.sh` (macOS). The build script will check dependencies and install any missing Wails/WebView2 pieces with your permission.

---

## For interviewers considering this tool

See [FOR_EMPLOYERS.md](FOR_EMPLOYERS.md) or the [online guide](https://sajucrajan.github.io/dayside/employers.html) for a short guide on responsible deployment.

Short version: give candidates advance notice, offer an alternative for those who can't or won't run it, train your interviewers to treat flags as conversation starters rather than verdicts, and accommodate candidates who rely on assistive software. In some jurisdictions — notably the EU, UK, and several US states — additional compliance steps apply. **You are responsible for your own legal compliance. This project is a tool, not legal advice.**

---

## Blocklist contributions

The tool is only as current as its blocklist. New commercial AI interview assistants and open-source clones launch every few weeks. If you spot one that Dayside doesn't catch, open a pull request against `internal/detect/allowlist.go` and `internal/detect/tabs_score.go`. Keep entries alphabetical within each section. A brief citation in the PR description (link to the product's site or a reputable news article) helps the review go faster.

The project also welcomes additions to the **accessibility allowlist** — legitimate assistive-technology tools (screen readers, voice control, captioning, magnifiers, etc.) that should never be flagged as suspicious. This list protects candidates with disabilities and we want it to be thorough.

---

## Privacy and consent

Every scan session begins with a consent screen showing the user exactly what the tool will do and giving them an unambiguous option to decline. No scan runs without explicit consent followed by an explicit click of the Start scan button. Periodic re-scanning during the interview is off by default and requires the candidate to opt in. Nothing is uploaded, now or ever.

See [PRIVACY.md](PRIVACY.md) or the [online privacy policy](https://sajucrajan.github.io/dayside/privacy.html) for full details.

---

## Security

Dayside takes supply-chain and dependency hygiene seriously.

**Dependency scanning.** The build script runs [`govulncheck`](https://go.dev/blog/govulncheck) on every build. If any dependency has a known vulnerability that affects code paths Dayside actually uses, the build fails until it's resolved.

**DLL-loading hardening (Windows).** Early versions of Wails v2 were vulnerable to a DLL-hijacking pattern where `uxtheme.dll` could be loaded from the application directory before the process could restrict its DLL search path. This was [reported and fixed by the Wails team](https://github.com/wailsapp/wails/pull/4207) and Dayside builds against the patched Wails version (v2.10.0 or later). On top of that, `main.go` explicitly calls `windows.SetDefaultDllDirectories(LOAD_LIBRARY_SEARCH_SYSTEM32)` at startup as defense in depth.

**Minimal dependency surface.** Dayside deliberately avoids pulling in cryptographic or network libraries — the tool does not do any cryptography, does not open any sockets, and does not speak any network protocol. This keeps the dependency graph small and the attack surface tiny.

**Reporting issues.** If you find a security issue — particularly anything that could cause Dayside to mis-identify a legitimate process as malicious, or anything that could be exploited to access data beyond what's documented — please report it privately via [GitHub's security advisory feature](../../security/advisories/new) rather than a public issue.

---

## License

[Apache License 2.0](LICENSE). Dayside is provided as-is with no warranty. You use it at your own risk.

---

*Dayside is a defensive tool. It does not facilitate, endorse, or enable the behavior it detects. Its purpose is to help interviewers conduct fair interviews, and to help honest candidates demonstrate their work on a clean machine.*
