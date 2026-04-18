package detect

import "strings"

// scoreTab flags a BrowserTab if its title or URL matches known copilot sites.
//
// Severity model:
//   - "red"    = a copilot product or known-cheating service
//   - "yellow" = a generic LLM chat (legitimate uses exist; flag for review)
//   - "green"  = no signal
func scoreTab(t *BrowserTab) {
	hayTitle := strings.ToLower(t.Title)
	hayURL := strings.ToLower(t.URL)

	// Known copilot URL domains.
	// Sources: vendor sites + GitHub repos + reviews indexed in v1.2 dossier.
	copilotDomains := []string{
		// === Tier 3 native desktop tools (still have marketing sites) ===
		"cluely.com", "cluely.cc",
		"interviewcoder.co", "interview-coder.com",
		"interviewsolver.com",
		"shadecoder.com",
		"linkjob.ai",
		"leetcodewizard.io", "leetcode-wizard.com",
		"parakeet-ai.com", "parakeetai.com",
		"senseicopilot.com", "sensei-copilot.com",
		"final-round.ai", "finalroundai.com", "finalround.ai",
		"interviewpilot.app", "interviewpilot.io",
		"lockedin-ai.com", "lockedinai.com", "support.lockedinai.com",
		"ultracode.ai",
		"vervecopilot.com", "verve-ai.com",
		"metaview.ai", "metaview.app",

		// === Tier 1 browser-based copilots ===
		"interviewcopilot.io",
		"interviews.chat",
		"ntro.io",
		"alinterviewprep.com",
		"interviewman.com",
		"hyring.com",

		// === Tier 2 browser extensions (also have web dashboards) ===
		// (Sensei AI's extension lives on chromewebstore - flagged via title)
		// (Ntro.io extension - flagged via domain above)

		// === Open-source clone landing pages ===
		"pickle.com/glass",   // Glass by Pickle
		"pluely.app",
		"natively.ai",        // Natively
		// (free-cluely, OpenCluely, Cheating Daddy live on github.com -
		// not blocked at domain level since github is legitimate; caught
		// via repo-name in URL path below)

		// === GitHub paths for known cheating repos ===
		// Substring match on full URL catches /Prat011/free-cluely, etc.
		"/free-cluely",
		"/opencluely",
		"/cheating-daddy",
		"/natively-cluely",
		"/pluely",
		"/cluely-alternative",

		// === DPRK / fraud-adjacent infrastructure (research dossier) ===
		// Real-time deepfake services often used in interview fraud
		"deepfacelive.com",
		"avatarify.ai",
		"reflect.tech", "reflect-ai.com",
		"reface.ai",
		"facefusion.io",
		"swapface.org", "swap-face.com",
	}

	for _, d := range copilotDomains {
		if strings.Contains(hayURL, d) {
			t.Severity = "red"
			t.Reason = "Matches copilot/cheating domain: " + d
			return
		}
	}

	for _, p := range knownCopilotTitles {
		if strings.Contains(hayTitle, p) {
			t.Severity = "red"
			t.Reason = "Title matches copilot: " + p
			return
		}
	}

	// === Generic LLM chats - yellow flag (legitimate uses exist) ===
	llmDomains := []string{
		// OpenAI
		"chat.openai.com", "chatgpt.com", "platform.openai.com",
		// Anthropic
		"claude.ai", "console.anthropic.com",
		// Google
		"gemini.google.com", "bard.google.com", "aistudio.google.com",
		"makersuite.google.com",
		// Microsoft
		"copilot.microsoft.com", "bing.com/chat",
		// xAI
		"grok.com", "x.com/i/grok", "twitter.com/i/grok",
		// Perplexity
		"perplexity.ai",
		// Meta
		"meta.ai",
		// DeepSeek
		"deepseek.com", "chat.deepseek.com",
		// Mistral
		"chat.mistral.ai", "le-chat.mistral.ai",
		// Aggregators
		"poe.com",
		"you.com",
		"phind.com",
		"hugging.chat", "huggingface.co/chat",
		// Coding-specific assistants
		"codeium.com", "windsurf.com",
		"cursor.com", "cursor.sh",
		"tabnine.com",
		"continue.dev",
		// Notebook / RAG
		"notebooklm.google.com",
		"app.wordtune.com",
	}
	for _, d := range llmDomains {
		if strings.Contains(hayURL, d) {
			t.Severity = "yellow"
			t.Reason = "Open LLM chat: " + d
			return
		}
	}

	// === Suspicious paths suggesting interview-cheat content ===
	// Even on legitimate sites, certain URL patterns are revealing
	// (e.g., a Reddit thread on "how to cheat in interview").
	suspiciousPaths := []string{
		"how-to-cheat",
		"interview-cheat",
		"cheat-on-interview",
		"undetectable-interview",
		"bypass-proctoring",
		"defeat-screen-share",
		"hide-from-zoom",
	}
	for _, p := range suspiciousPaths {
		if strings.Contains(hayURL, p) {
			t.Severity = "yellow"
			t.Reason = "Suspicious URL path: " + p
			return
		}
	}

	t.Severity = "green"
}
