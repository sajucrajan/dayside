# What You Can't See on the Screen-Share

**Invisible AI assistants have quietly changed what a remote interview actually measures. Here's a small, read-only, open-source tool to help get back to measuring the candidate.**

*By [your handle] · April 2026 · ~8 min read*

---

A friend of mine runs technical hiring for a mid-sized software company. Last month she told me she'd stopped trusting any remote interview. Not the candidates — the format. She'd been on enough calls recently where the candidate's answers were slightly too polished, slightly too fast, slightly too correct on problems that should have produced some visible thinking. She couldn't prove anything. She couldn't see anything. But she could feel it.

She's not imagining it. There's a category of software — a growing one — whose entire value proposition is making that exact experience invisible to her.

This is a piece about what that category is, why it works, why nothing in your existing interview stack catches it, and the small open-source tool I built to do something about it. I want to be honest upfront: this isn't a solved problem. The tool I'm publishing is a first line, not a last line. But "nothing" is worse than "something," and nothing is roughly what interviewers have today.

---

## The shape of the problem

A few years ago, the frontier of interview cheating was a second laptop. Someone off-camera with a ChatGPT tab open, whispering or typing answers into a shared document. It worked, sort of, but it was clumsy. The candidate had to split their attention. The delay was obvious. Interviewers learned to look for the tells — the pauses, the eye flicks, the too-formal phrasing.

What's changed in the last eighteen months is that the whole apparatus moved onto the candidate's primary machine and rendered itself invisible to the screen-share.

The category now includes, broadly:

- **Browser-based assistants** that listen to the interview audio in a tab, transcribe questions in real time, and display suggested answers in the same browser window. They don't need to install anything. They bill themselves as "no download required."
- **Desktop overlays that set themselves as "topmost"** — always on top of other windows — so the candidate can see them while working. Early versions were easily caught because a topmost window also showed up in the screen-share. Later versions fixed that.
- **Desktop overlays that exclude themselves from screen capture.** This is the generation that matters. Every modern operating system has APIs that let a window flag itself as "please do not include me when the OS captures the screen." Those APIs exist for legitimate reasons — DRM, password managers that don't want to leak into a shared video call. Invisible AI assistants use them to render answers in a floating window that the candidate can read while the interviewer, staring at the same screen over Zoom, sees nothing.
- **Second-device tools** that still exist, now more polished — phones propped up out of camera view, a smartwatch with scrollable transcripts, a tablet mirrored to a hidden display.
- **Open-source clones** that replicate each of the above, released within days of the commercial products they copy, removing any cost barrier.

The commercial leaders in this space are not hiding. They have marketing pages. They list "invisible to screen sharing" as a feature. They have testimonials. One widely-covered product raised fifteen million dollars in venture funding after its founder was suspended from his university for building it. Another markets itself with the tagline "cheat on everything." These are not shadowy tools — they are well-funded companies with designers and growth teams.

To give this category a name for the rest of this piece without singling out any real product, I'll call the archetype **GhostPrompt**. It's fictional. No product by that name currently exists, as of this writing. But every capability I'll describe exists in real shipped software you can buy today.

GhostPrompt runs on the candidate's machine. It listens to the interview audio through the system microphone. When the interviewer asks a question, GhostPrompt transcribes it, sends it to a large language model, and displays the answer on a floating window that's excluded from screen capture. The candidate reads the answer, rephrases it in their own words, and delivers it with the cadence of someone thinking. The interviewer, screen-sharing over Zoom, sees the candidate's code editor, the candidate's browser, the candidate's Slack. Nothing suggests anything else is happening. The latency from question to on-screen answer is under two seconds.

This is not theoretical. This is shipping software, used by tens of thousands of candidates a month, advertising itself on the platforms interviewers use.

---

## Why your current defenses don't work

Here's what interviewers typically have in their toolbox, and what each one actually catches.

**Screen sharing.** Zoom, Teams, Meet, Slack Huddles, Around — they all share what the OS tells them is visible. If a window sets itself as "not for screen capture," the conferencing tool never sees it. There is nothing the interviewer can click to override this. The conferencing tool itself doesn't know the window exists. Capture-excluded windows are a feature of the operating system, not a bug in Zoom.

**Camera monitoring.** Watching the candidate's eyes only tells you they're looking at their screen, which they have to be anyway to code. "Eye drift" used to be a signal when the answers were on a second device; now the answers are on the same screen, a few pixels from where the candidate's eyes should legitimately be. The signal is gone.

**Asking the candidate to show their desktop.** They show it. The capture-excluded overlay remains invisible. You see a clean desktop. Nothing is hidden because nothing is shown.

**Asking the candidate to turn off all other applications.** They do. Some of them do. The honest ones do. The ones using GhostPrompt click the minimize button, which on these tools means "hide the floating window" — the process keeps running, still listening, ready to reappear the moment they alt-tab back. Closing the window and closing the process are different things, and one of them is easy to fake.

**Enterprise proctoring software.** These exist — ProctorU, Proctorio, Honorlock, Respondus Monitor, Examity. They're built for the academic market and they're aggressive: they record the candidate's room with the webcam, they track eye movement with facial analysis, they flag background noise, they scan the desktop for forbidden processes. They also have a long, ugly track record. Class-action lawsuits over biometric data collection. Settlements. Reporting on algorithmic bias against darker-skinned and disabled students. A federal judge ruling that a live room-scan violated the Fourth Amendment at a public university. This is not a template for corporate hiring. Using academic proctoring tools to vet job candidates is a liability invitation and a candidate-experience disaster. It's also wildly disproportionate for a 45-minute coding interview.

**Asking the candidate to promise they won't cheat.** The honor system works on honest people. It was always going to.

---

## What would work, actually

The thing that would actually catch GhostPrompt is someone enumerating the candidate's running processes, visible and hidden windows, connected audio and video devices, open browser tabs, and remote-access tools — and comparing what they find to a list of known invisible AI assistants and suspicious patterns.

In principle, a candidate could do this themselves. Open Task Manager, walk the interviewer through every running process, alt-tab to every window. In practice, this is clumsy, unreliable, and nobody wants to do it at the start of a high-pressure call. More importantly, it won't catch windows that are capture-excluded, because the candidate may not know which of their windows are. They know they installed GhostPrompt. They may not know precisely which API it uses to hide from Zoom.

What I wanted was a tool that would do the enumeration properly, in about ten seconds, and show the results to the candidate themselves. Not to the interviewer. Not to a server. Just to the candidate, on the candidate's screen. The candidate shares that screen, both people look at the results, they discuss anything flagged, and — if anything needs closing — the candidate closes it themselves using their operating system's own task manager. The tool never touches the machine; the candidate's agency is preserved end-to-end.

That's what I built.

---

## Meet Dayside

Dayside is a small, open-source executable for Windows and macOS. The name comes from astronomy — the *dayside* of a planet is the lit side, the face turned toward its star, the part where everything is visible. That's the posture of the tool: everything in daylight, nothing hidden. I wanted a name that communicated the tool's philosophy, not a name that communicated its aggressiveness. "Spectre" and "Sweep" were earlier working names; both felt adversarial in a way that didn't match what the tool actually does.

Here's how using Dayside actually works.

The candidate downloads it before the interview and runs it. A consent screen appears describing exactly what the tool will do. If they agree, the main window opens in a **Ready** state. Nothing has scanned yet. The candidate sees a prominent "Start scan" button, a checkbox for optional periodic re-scanning (unchecked by default), and a toggle for compact mode. That's it.

When the interview starts and the screen-share begins, the candidate clicks **Start scan**. The scan runs for roughly ten seconds and displays results locally. Both the candidate and the interviewer look at the results together over the screen-share and discuss anything flagged.

If the candidate wants, they can tick "Keep checking every 5 minutes." Dayside will re-run the same scan periodically to catch anything launched after the interview starts. **This is off by default** — it's a candidate choice, not a default behavior. Some candidates will want it on because it reassures them that the machine remained clean throughout; others will leave it off because a single pre-interview check is enough.

The candidate can shrink the main window to a compact pill — small, draggable, always-on-top — so the tool doesn't take up working space while they're coding. The pill shows the current scan status ("clean," "2 flags," whatever) and the time of the last scan. Both parties can glance at it any time. If it's green, the last scan was clean; if it's yellow or red, there's something to talk about.

No data is transmitted. No account is required. No server is contacted. Ever.

**Dayside is strictly read-only.** It looks at what's running on the machine. It does not close, start, or modify any process. It cannot change a single file or setting. If the candidate decides to close something the scan flagged, they do so through Windows Task Manager or macOS Activity Monitor — the same way they'd close any other program. The tool is an information display, nothing more. That matters, both ethically (the candidate's machine remains entirely under their control) and practically (there's no class of bug in Dayside that can damage a candidate's system, because the code paths to do so do not exist).

What the scan checks, in rough order of what catches the most real-world AI assistants:

**Running processes.** Names and signatures compared to a maintained blocklist of commercial AI interview assistants, open-source clones of those products, tunneling tools frequently used to bridge a second device, and remote-access software. If the candidate has one of these running, it's flagged red.

**Visible and cloaked windows.** The tool enumerates windows using the platform's native window API and separately asks each window whether it has flagged itself as excluded from screen capture. If a window exists but claims it's excluded, that's noteworthy. Legitimate password managers do this. So do some DRM-protected video players. So do interview-assistance tools. The tool doesn't assume intent — it surfaces the signal so the candidate and interviewer can discuss it.

**Virtual audio and video devices.** A virtual camera (OBS Virtual Camera, Snap Camera, ManyCam) is a deepfake vector. A virtual audio cable (VB-Audio, Blackhole, Loopback) is a voice-routing vector — a way to feed a second audio stream into the call for coaching. These are flagged; candidates with legitimate reasons to have them (streamers, podcasters, conference presenters) can explain.

**Browser tabs.** On Windows, the tool reads window titles. On macOS, it uses AppleScript to read both titles and URLs. It compares both against a list of known AI-assistant domains and popular LLM chat sites. A ChatGPT tab open during an interview isn't conclusive evidence of anything — people use it for perfectly normal work — but it's worth a sentence of conversation.

**Remote access and session state.** TeamViewer, AnyDesk, Parsec, RustDesk, and half a dozen others are enumerated. The tool also checks whether the current session is a Remote Desktop session, which would mean someone else might be watching or controlling. This is rarely an issue, but when it is, it's a big one.

**System state.** Monitor count, architecture, virtual-machine indicators, whether this looks like a work laptop that someone might be managing remotely.

Every flag comes with a severity (critical / warning / clean), a reason, and a label. The UI is optimized to be readable in a screen-share at normal resolution: large type, clear color coding, a concise "how to close this" hint per flagged item pointing to Task Manager or Activity Monitor. No drill-downs, no hidden settings.

---

## About the always-on-top thing

Someone will inevitably point out the irony: the tool built to catch always-on-top invisible AI assistants is itself an always-on-top window. I want to address this directly because the difference matters.

Always-on-top is a UI pattern, not a cheating pattern. Password managers are sometimes always-on-top. Color pickers are always-on-top. Sticky notes are always-on-top. The pattern itself is neutral.

What makes an invisible AI assistant a cheating tool isn't that it stays on top — it's that it **stays on top while excluding itself from screen capture**, so only the candidate sees it and the interviewer sees nothing. It is the combination of visible-to-one-party plus invisible-to-the-other that defines the category.

Dayside is the opposite on the second axis. It is visible to both parties. It does not exclude itself from screen capture. Its compact-mode pill is literally designed to be seen during the screen-share — by both the candidate and the interviewer simultaneously. The interviewer watching the screen-share sees the same pill the candidate sees, in the same position, showing the same status. That's the whole point. It's always-on-top in the service of transparency, not in the service of secrecy.

If Dayside ever excluded itself from screen capture, it would be a cheating tool disguised as a defensive tool, and you shouldn't use it. The fact that both parties can see it at all times is what makes the design defensible.

---

## The things I decided not to do

A lot of the design of Dayside is what it *doesn't* include, and I think those choices matter more than what it does.

**Nothing is uploaded.** No server. No telemetry. No phone-home analytics. The binary makes zero network calls. If you run it with the network disconnected, it behaves identically to running it online. This isn't a performance optimization — it's an architectural choice about what kind of tool this should be.

**Nothing persists.** When the application closes, nothing remains on the machine. No log files. No cached scan results. No preferences. Running it is a single, self-contained session. The next candidate who borrows the same machine starts from zero.

**No biometrics. Ever.** No webcam access. No microphone access. No face-matching, no voice-printing, no behavioral analysis of typing patterns. This is the line that keeps the tool out of an entire class of privacy laws and lawsuits. It's also just the right line, independent of the legal consequences.

**No ability to modify the candidate's machine.** This is the read-only point again, and I want to be explicit: the tool has no code path that can terminate a process, kill a connection, modify a file, or change a setting. The design simply does not include those capabilities. If you reverse-engineer the binary you will not find an "are you sure you want to close this?" hidden under a flag — those code paths do not exist. This is deliberate. It removes an entire class of "what if the tool does something bad?" concerns by making them architecturally impossible.

**No automatic scans.** Consenting to use Dayside does not start a scan. The candidate has to explicitly click **Start scan** in the main UI. The same pattern applies to periodic re-scanning: the candidate has to explicitly tick the "Keep checking" checkbox to enable it. The tool is never running without the candidate making an affirmative choice to have it run. This matters because it preserves the candidate's control over *when* the interviewer sees results — a screen-share moment that happens on the candidate's timing, not on the app's.

**No automatic decisions.** The tool flags. It does not score. It does not produce a pass/fail verdict. It does not rank candidates. It does not integrate with applicant tracking systems. The outputs are signals for a conversation between two humans, and any attempt to automate the conversation would undermine what makes the approach defensible.

**An explicit accessibility allowlist.** This was actually the first thing I built after the core detection logic. There is a class of software — screen readers, voice-control tools, captioning software, magnifiers, reading assistants — that shares technical characteristics with invisible AI assistants. They run topmost. They hook every window. Some of them are unsigned. Flagging them as suspicious creates a disability-discrimination problem the moment the tool is used in hiring. Dayside ships with these on a green-labeled allowlist. If something is missed, the fix is a pull request, not a manual review.

**Consent is affirmative and repeated.** Every launch of the tool starts with a consent screen the candidate has to actively agree to. The "Agree" button is disabled until a checkbox is ticked. The "Decline" button is visually equivalent to the "Agree" button — no dark patterns. If they decline, the tool closes without scanning anything. Consent does not persist across launches, because a shared machine may be used by a different person next time.

None of this makes Dayside bulletproof. A sophisticated attacker — a state-sponsored actor, a well-funded individual with kernel-level access, a candidate using an entirely separate device out of camera range — can defeat the tool. That's true of every defensive tool, ever. The question is not whether the tool is perfect but whether it's better than nothing, and I think the answer to that is clearly yes.

---

## Who this is for and how to deploy it responsibly

This tool exists for interviewers who care about fairness, candidates who want to demonstrate they're not cheating, and hiring processes that currently have no answer to the invisible-overlay problem. It's open-source, permissively licensed, and free for any use consistent with the Apache 2.0 license.

If you're an interviewer thinking about using it, I want to be direct about what responsible use looks like. I've written a longer deployment guide — linked from the repo — but the short version is: give candidates advance notice, offer an alternative for anyone who can't or won't run it, train your interviewers to treat flags as questions rather than verdicts, accommodate candidates who rely on assistive software, and don't retain the results.

If you're in a jurisdiction with significant data-protection or employment-monitoring law — the EU, UK, parts of Canada, several US states — there are additional compliance steps that apply. The tool's local-only, read-only architecture makes most of these lighter than they'd be for a cloud-based proctoring product, but lighter is not zero. Talk to your employment counsel before you deploy it at scale.

And don't use it for surveillance. This is a pre-interview and during-interview check, not an ongoing monitoring tool. If you find yourself running it on current employees, or using it at every stage of every interview for every role, you've drifted away from the proportionate use case.

---

## Why open-source

I thought about shipping this as a closed binary. Legal advice I looked at suggested closed-source has slightly worse liability characteristics for an individual developer, and obscured detection logic makes it hard for security-conscious candidates to verify what the tool actually does before running it on their machines. More importantly: the blocklist of AI interview assistants needs to stay current, and new ones launch every few weeks. One person can't maintain that list alone. Open-source means that when a security researcher spots a new tool in the wild, they can send a pull request with the new executable name, and the whole user base gets the update.

There's also something appropriate about defensive tools being more transparent than the offensive ones they detect. The commercial invisible-AI-assistant products I'm not naming are, mostly, closed-source. Their open-source clones are on GitHub. If the offensive side of this arms race is half-and-half open, the defensive side should at least be fully open.

The repository is at [link]. The Medium article you're reading is the long-form pitch; the README is the short version. Contributions are welcome, especially additions to the blocklist (new tools to catch) and the accessibility allowlist (legitimate tools being flagged by mistake). I'm maintaining it under a pseudonym because the commercial invisible-assistant founders are unlikely to appreciate the attention and I'd rather not find out. If you want to reach me, GitHub issues are the channel.

---

## One more thing

I don't think this problem gets solved by a better scanner alone. The deeper fix is interview formats that are less vulnerable to the specific failure mode. Pair programming on problems the candidate picks. Take-home projects with follow-up technical conversations. Work-sample tests where the signal is how someone thinks through a problem over hours, not how quickly they answer a LeetCode question. Behavioral interviews where what matters is the specificity and texture of the candidate's own stories, which are hard for any AI to fabricate convincingly.

Dayside is a tool for the interviews you have today, not the interviews we should be running tomorrow. I hope it's useful. I hope, eventually, it becomes unnecessary.

If you use it and it finds something, have the conversation. Continue with the interview if it makes sense to continue, and end it if it doesn't. Hire based on what you learn, not what the tool scored. That's the whole point.

---

*Dayside is licensed under Apache 2.0. Source, binaries, and documentation at [repository URL]. No account, no telemetry, no upload, no ability to modify your machine — the tool is strictly read-only. Bug reports and blocklist contributions welcome. This article describes a category of software without singling out specific commercial products; "GhostPrompt" is a fictional placeholder used to illustrate the archetype, and any resemblance to a real product of that name is coincidental.*
