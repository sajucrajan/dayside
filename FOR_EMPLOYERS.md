# For Employers Considering Dayside

*This guide is for hiring teams, recruiters, and interviewers evaluating whether to use Dayside in their interview process. It is practical guidance, not legal advice. Before deploying the tool in real interviews, consult employment counsel for your jurisdiction.*

---

## What Dayside is — and isn't — from your perspective

**It is** a candidate-operated, read-only diagnostic tool. The candidate downloads and runs the executable on their own machine, voluntarily, before the interview starts. The candidate explicitly clicks **Start scan** when they're ready — the tool does not scan automatically after consent. Results appear on the candidate's screen. During a screen-share, the interviewer can see the results and discuss anything flagged. The candidate may optionally enable a "Keep checking every 5 minutes" checkbox to catch anything launched during the interview, but this is off by default.

The tool stays visible during the interview in a small always-on-top pill ("compact mode") so both parties can see at a glance that the scan is current. If the candidate decides to close a flagged program, they do so themselves through their operating system's standard tools (Task Manager on Windows, Activity Monitor on macOS) — Dayside never takes any action on the machine.

**It is not** a proctoring tool, a background check, a biometric system, or a hiring algorithm. It does not make decisions. It does not send data to any server. It does not record. It does not terminate processes. It does not modify anything on the candidate's machine. It does not keep anything after it closes. It does not hide from the screen-share — it is designed to be visible to both parties throughout the interview.

**It will not** detect:
- AI interview assistants running on a second device (phone, tablet, another laptop)
- Kernel-level tools designed to evade process enumeration
- Anything the candidate chooses not to run the scanner on
- Anything launched after the last scan if the candidate did not enable periodic re-scanning

**It may** produce false positives. Legitimate software — especially assistive technology — can share characteristics with invisible AI assistants (topmost windows, unsigned binaries, system-wide hooks). Treat every flag as a conversation starter, not a verdict.

---

## The short responsible-deployment checklist

Before using Dayside in a single interview, make sure you've done all of these:

- [ ] **Given candidates advance written notice** that you'll ask them to run a pre-interview machine check
- [ ] **Explained that the scan is candidate-initiated** — they consent, then they click Start scan when they're ready
- [ ] **Mentioned that periodic re-scanning is optional** — the candidate can choose whether to enable it
- [ ] **Mentioned that the tool stays visible** during the interview as an always-on-top pill
- [ ] **Offered an alternative** for candidates who cannot or will not run the tool (e.g., manually share Task Manager or the macOS Activity Monitor)
- [ ] **Established an accommodations process** for candidates who rely on assistive technology
- [ ] **Trained your interviewers** on how to interpret flags and what to do with false positives
- [ ] **Decided who reviews the results** and how they're documented (if at all) in the hiring file
- [ ] **Reviewed your jurisdiction's requirements** for pre-employment screening and electronic monitoring
- [ ] **Reviewed your own privacy notices** to confirm the scan is described accurately

Skip any of these and you're increasing legal and ethical risk. None is burdensome, and together they take maybe a few hours to set up organization-wide.

---

## Pre-interview notice template

Send something like this with the interview invitation, at least 48 hours before the interview:

> **About your upcoming interview**
>
> To ensure a fair interview environment for all candidates, we'll ask you to run a short pre-interview machine check at the start of the call. The tool is called Dayside; it's an open-source, candidate-operated diagnostic that looks for known AI interview assistants on the machine you'll use for the interview.
>
> - Download: [URL]
> - Initial scan: roughly 10 seconds
> - Data sent off your machine: none
> - Recording: none
> - The tool is read-only — it shows you what's running but never closes, modifies, or changes anything on your machine
>
> **How it will work during our interview:**
>
> 1. Before the interview, you open Dayside. A consent screen appears describing what the tool will do. If you decline, the tool closes.
> 2. If you agree, the main Dayside window opens in a "Ready" state. **Nothing scans yet.** When we start the screen-share at the beginning of our call, you click "Start scan" and the first scan runs.
> 3. We'll look at the results together and discuss anything flagged.
> 4. You can optionally tick "Keep checking every 5 minutes" if you'd like the tool to re-scan during the interview. This is your choice — it's off by default and we won't ask you to turn it on.
> 5. The tool has a compact "pill" mode that stays visible in a small window during the interview so you can focus on the interview itself. You can drag it anywhere on your screen.
>
> **If you cannot run the tool**, for any reason, please reply to this email and we'll arrange an alternative — typically a brief walk-through of Task Manager (Windows) or Activity Monitor (macOS).
>
> **If you rely on assistive software** (screen readers, voice control, captioning, magnification, etc.), please let us know in advance so we can interpret any results correctly. Running these tools will never affect your candidacy.
>
> **Questions?** Reply to this email.

Adapt to your tone and brand. The core points are: advance notice, clear description of the flow (consent → Ready → candidate clicks Start scan), read-only nature of the tool, optional periodic re-scanning, compact mode stays visible, alternative offered, accommodations welcomed, contact provided.

---

## Training interviewers — a one-page brief

Spend ten minutes with anyone who'll be on the interviewing side of the scan. Cover:

**A flag is not a verdict.** Every row on the scan results is a signal that deserves a question, not a conclusion. "I see Otter.ai running — tell me how you use it?" leads to useful information. "You have Otter.ai, so you're disqualified" leads to a lawsuit.

**Expect false positives.** Legitimate apps often look suspicious to a tool like this. Screen recorders, meeting notetakers, password managers, streaming software, virtual cameras for good reasons (OBS for conference presenters), remote-access tools used by legitimate IT support — all can flag. The candidate's explanation is data. Accept it unless something about the explanation is itself off.

**Never ask about disabilities.** If the flagged item turns out to be assistive technology, the candidate may volunteer the disability context. You do not need it. Accept "that's an accessibility tool I use" as a complete answer, note it was explained, move on. Under the ADA you are specifically prohibited from asking about disabilities before extending an offer, even if the candidate has opened the door themselves.

**The tool is observational. Any action belongs to the candidate.** Dayside does not close programs — it only shows you what's running. If the candidate decides to close something, they do it themselves through Task Manager or Activity Monitor. You never direct action on the candidate's machine. If the candidate declines to close something, that's their choice — note it and discuss whether the interview should continue.

**Periodic re-scanning is the candidate's choice.** The candidate may or may not enable the "Keep checking every 5 minutes" option. Both choices are legitimate. Do not pressure the candidate to enable it. Some candidates will prefer the single initial scan; others will want continuous monitoring for their own peace of mind. Either is fine.

**The compact pill stays visible — that's intentional.** During the interview you'll see a small Dayside pill hovering on the candidate's screen showing scan status and timestamp. This is not evasion; this is the design. It means the tool is honestly continuing to be present and visible during the interview rather than disappearing into the background. If the pill disappears, it's because the candidate closed the app entirely, which would be unusual and worth asking about.

**Document the conversation briefly.** For each flag that mattered: what was flagged, what the candidate said, what happened. Two sentences. This protects both the candidate (evidence you didn't make unfounded inferences) and the company (evidence of interactive dialogue if the hire or non-hire is later questioned).

**You may see things you don't need to know.** Browser tabs sometimes reveal personal context — a medical portal, a job-search site, a family emergency tab. Ignore everything that isn't directly relevant to AI-assistant detection. What you see during a screen-share is not yours to remember, share, or use in decisioning.

---

## What to do with the results

**Do not retain them.** No screenshots, no notes, no copy-pasted lists of running processes into the hiring file. The scan is a momentary verification, not an audit artifact. The only thing that should land in the hiring file is a short summary like "pre-interview machine check completed, no critical flags" or "pre-interview machine check completed, flag for Otter.ai explained by candidate as meeting-notes tool, allowed to continue."

**Do not use them as a hiring signal on their own.** A clean scan does not mean the candidate is qualified. A flagged scan does not mean they're disqualified. Results inform the interview; they don't replace it.

**Do not share them outside the immediate hiring team.** The results may contain information that's sensitive in ways you don't realize (research tabs revealing competitive intelligence, app choices suggesting religion or health status, etc.). Treat them like any other private information gathered during interview: need-to-know within the hiring decision, then forgotten.

---

## Refusal is allowed — and how to handle it

Some candidates will decline to run Dayside. Reasons vary: they're using a work laptop whose IT policy prohibits unapproved software, they're privacy-conscious, they're uncomfortable with the request, they don't trust the binary. All of these are legitimate.

**Refusal is not evidence of cheating.** A candidate who refuses may be entirely honest and simply uncomfortable with the premise. The tool is opt-in; the consent is real only if refusing is a genuine option.

**Offer an alternative.** The most common one: ask the candidate to open Task Manager (Windows) or Activity Monitor (macOS) and walk through the running processes with you verbally. This covers the same ground with less friction. You lose the automated blocklist matching; you gain a candidate who's comfortable with the conversation.

**Document the refusal and the alternative used.** One line in the hiring file.

**Do not let refusal alone drive a decision.** If you end up rejecting the candidate, the reason should be something independently defensible — performance in the interview, qualifications, culture fit based on substantive signals — not "they declined the scan." Decisions that can be traced to the refusal are vulnerable to discrimination claims if the refuser happens to be in a protected class.

---

## Jurisdiction-specific notes

**This is a starting summary, not comprehensive. Verify with local counsel.**

### United States

- **ADA / disability accommodation (federal).** You must accommodate candidates who use assistive technology. Dayside ships with an accessibility allowlist that should catch most common tools, but plan for the possibility that something legitimate is flagged. The accommodation process is interactive — talk to the candidate, find a workable path, document it.
- **Illinois BIPA.** Dayside does not collect biometric data, so BIPA does not apply in its current form. Do not add any biometric verification on top of this tool unless you've worked through BIPA compliance.
- **California CCPA / CPRA.** Job applicants are covered under the 2023 amendments. Include Dayside in your candidate privacy notice. Mention that the scan runs locally and nothing is uploaded.
- **New York City Local Law 144.** If you deploy Dayside as a substantial factor in hiring decisions, it may qualify as an "automated employment decision tool" requiring an annual bias audit and candidate notice. Pattern-matching against a blocklist is arguably not an AEDT (it doesn't "substantially assist or replace discretionary decision-making"), but the question is unsettled. Err on the side of disclosure.
- **State electronic monitoring laws.** New York, Connecticut, and Delaware require advance notice of electronic monitoring of employees. These apply to employees rather than candidates, but candidates become employees on day one — your pre-interview notice should align with your ongoing monitoring notice. Note that if candidates enable the periodic re-scanning option, the session becomes closer to "ongoing monitoring" for the duration of the interview than a single check; describe this accurately in your notice.

### European Union and United Kingdom

**Higher bar. Plan for more compliance work.**

- **GDPR lawful basis.** Consent in an employment context is generally not valid because of the power imbalance. Rely on **legitimate interests** (Article 6(1)(f)) and document a Legitimate Interests Assessment before deployment. Note that if candidates enable periodic re-scanning, the assessment should cover both the single-scan and periodic-scan modes.
- **Privacy notice.** Required. Describe the processing in plain language before the interview. What data is processed, on what basis, for what purpose, what rights the candidate has. Include the periodic re-scanning option and make clear it's opt-in.
- **Data Protection Impact Assessment (DPIA).** Almost certainly required under Article 35 because you're running software on a candidate's machine to inspect their activity. This is a document, not a process — but skip it and you've committed a procedural violation that carries its own fines regardless of any substantive harm.
- **EU AI Act.** The Act classifies AI used in recruitment as "high-risk" with onerous obligations effective August 2, 2026. Dayside as currently designed is rule-based pattern matching, which likely does not qualify as "AI" under the Act's definition. If that changes (if you add ML-based flagging), the high-risk obligations apply.
- **Works councils (Germany especially).** German Works Councils have co-determination rights over "technical equipment designed to monitor employees" (BetrVG §87(1) No. 6). Consult the Works Council before rolling the tool out across your hiring process, even though candidates aren't employees — the interpretation leans toward consultation being required.
- **CNIL guidance (France).** French Labor Code Article L. 1221-8 requires candidates to be informed in advance of any assessment method used in recruitment. Send the pre-interview notice; don't surprise candidates.

### Canada

- **PIPEDA (federal).** Requires knowledge and consent for collection of personal information.
- **Ontario Working for Workers Act.** Employers with 25+ employees must have a written electronic monitoring policy covering applicants.
- **Quebec Law 25.** Explicit consent standard for processing personal data.

### Australia

- **Privacy Act 1988 and Australian Privacy Principles.** Notice requirements for collection.
- **NSW/ACT/Victoria workplace surveillance.** Advance notice required for employees; candidate protection is thinner but align with the same standard.

### Other

- **Brazil (LGPD), India (DPDP Act), Singapore (PDPA), South Africa (POPIA)** — all broadly GDPR-influenced. Treat as requiring explicit notice and lawful basis; legitimate interest typically available.
- **China (PIPL)** — explicit consent for cross-border data transfers; Dayside's local-only design avoids most PIPL complications but notice remains required.

---

## What NOT to do

- **Do not require candidates to install the tool on company-issued devices from their current employer.** That's a violation of their current employer's IT policy and puts the candidate in an impossible position.
- **Do not run the tool yourself on the candidate's behalf.** The consent architecture depends on the candidate operating it. If you send IT to "help," you've converted a candidate-operated diagnostic into employer monitoring, which is a much higher bar.
- **Do not pressure candidates to enable periodic re-scanning.** The checkbox is opt-in for a reason. Asking "can you turn on continuous scanning?" shifts the power dynamic and undermines the consent architecture that makes this tool defensible. If you want continuous monitoring as a requirement, use a proctoring tool designed for that purpose — and accept the higher legal burden that comes with it.
- **Do not keep the results.** No screenshots, no video of the scan, no log files copied to your hiring system. What you saw, you saw; that's the limit.
- **Do not use Dayside on candidates who have not received the pre-interview notice.** Surprise deployments create consent problems in every jurisdiction.
- **Do not treat refusal as dispositive.** See the refusal section above.
- **Do not rebrand or modify the tool to misrepresent what it does.** Apache 2.0 permits modification; it does not permit deceiving candidates about what they're running.
- **Do not use the tool for employee monitoring.** It is designed for a single interview session, with consent, on a candidate's own machine. Turning it into an ongoing surveillance tool for current employees crosses into territory this project does not support.

---

## A word on proportionality

Dayside is appropriate when the role involves significant coding, reasoning, or communication that would be materially affected by AI-assistant use during the interview itself — engineering interviews, technical screens, real-time problem-solving sessions. It is **not** appropriate for casual first-round conversations, behavioral interviews with no technical content, or junior roles where the interview is primarily about motivation and fit.

Running a machine scan on a candidate you haven't yet decided to take seriously is disproportionate and creates a bad candidate experience. The best use is: a strong candidate who has passed earlier rounds, a technical interview where real-time AI assistance would materially distort the signal, and a brief opt-in scan at the start that either passes cleanly or opens a short conversation and then passes.

If the tool finds itself being used in every interview, at every stage, for every role, something has gone wrong. It is a targeted instrument, not standard-issue surveillance.

---

## Questions, feedback, bug reports

- **Bugs**: open an issue on GitHub
- **False-positive reports**: open a PR against the allowlist with a link to the product's site
- **Missing AI interview assistants**: open a PR against the blocklist
- **Security issues**: use GitHub's private security advisory feature
- **Legal questions about your specific situation**: consult your own counsel; the project cannot give you legal advice

The project is maintained by volunteers. Response times are best-effort. If you need commercial-grade support, you should use a commercial tool.
