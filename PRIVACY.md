# Privacy Policy

Dayside is a local, offline tool. This document describes exactly what data it accesses and what it does with it.

## What Dayside reads

During a scan, Dayside reads:

- **Running processes** — names and paths, to match against a blocklist
- **Open windows** — titles and screen-capture exclusion flags
- **Audio/video devices** — device names, to detect virtual cameras and audio cables
- **Browser tab titles** — on Windows; titles and URLs on macOS
- **System state** — monitor count, RDP/screen-sharing session status

## What Dayside does NOT do

- **No network calls.** Dayside does not connect to the internet. Ever. There is no telemetry, no analytics, no update check, no license server.
- **No data storage.** Nothing is written to disk. When the app closes, nothing remains.
- **No recording.** Screen, camera, and microphone are not accessed.
- **No biometrics.** No face scan, no voiceprint, no keystroke logging.
- **No process termination.** Dayside is read-only. By design, it does not kill, suspend, or modify any process.
- **No background operation.** The app does not start on login, does not run as a service, and does not operate after it is closed.

## Consent

Every session begins with a consent screen that explains what the scan will do. No scan runs unless the user explicitly clicks "Start scan" after reading and accepting the consent screen. Periodic re-scanning is off by default and requires the user to opt in.

## Data shared with the interviewer

The only data shared is what appears on the candidate's screen during their screen-share session — the scan results panel that both parties look at together. No data is transmitted electronically to the interviewer or to any third party.

## Open source

Dayside is fully open source under the Apache License 2.0. Anyone can read the code, build it from source, and verify these claims independently.
