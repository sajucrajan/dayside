// Dayside — frontend logic (observation-only build)
// Wails v2 exposes Go methods on window.go.main.App.*

// ============================================================
// Constants
// ============================================================

const FLAG_LABELS = {
    KNOWN_COPILOT: 'Known AI interview assistant',
    REMOTE_ACCESS: 'Remote access tool',
    SYSTEM_IMPERSONATION: 'System impersonation',
    HIDDEN_FROM_SCREEN_SHARE: 'Hidden from screen share',
    COPILOT_TITLE_MATCH: 'Title matches known assistant',
    CLOAKED_WINDOW: 'Cloaked window',
    TRANSPARENT_OVERLAY: 'Transparent overlay',
    UNSIGNED: 'Unsigned binary',
    RUNS_FROM_TEMP_OR_DOWNLOADS: 'Runs from Temp/Downloads',
};

const FLAG_REASON = {
    KNOWN_COPILOT: 'known:copilot',
    REMOTE_ACCESS: 'remote-access',
    SYSTEM_IMPERSONATION: 'path:mismatch',
    HIDDEN_FROM_SCREEN_SHARE: 'win32:affinity',
    COPILOT_TITLE_MATCH: 'title:copilot',
    CLOAKED_WINDOW: 'win32:cloaked',
    TRANSPARENT_OVERLAY: 'ui:overlay',
    UNSIGNED: 'sig:missing',
    RUNS_FROM_TEMP_OR_DOWNLOADS: 'path:temp',
};

const FLAG_HINT = {
    KNOWN_COPILOT: 'This process matches a known AI interview-assistant tool.',
    REMOTE_ACCESS: 'Screen-sharing / remote-control software. Close before the interview unless the interviewer is running this.',
    SYSTEM_IMPERSONATION: 'Process name matches a system binary but runs from a non-system path.',
    HIDDEN_FROM_SCREEN_SHARE: 'Window is set to hide from screen-sharing (WDA_EXCLUDEFROMCAPTURE). Designed to stay invisible in share.',
    COPILOT_TITLE_MATCH: 'Window title contains strings associated with known AI assistants.',
    CLOAKED_WINDOW: 'Windows reports the window as cloaked — visible to the user but hidden from some capture APIs.',
    TRANSPARENT_OVERLAY: 'Layered + topmost window — commonly used by transparent overlays.',
    UNSIGNED: 'Binary has no digital signature. Legitimate software usually is signed.',
    RUNS_FROM_TEMP_OR_DOWNLOADS: 'Executable is running from a temporary / downloads folder.',
};

// Severity normalization (backend uses red/yellow/green)
const SEV_TO_DS = { red: 'crit', yellow: 'warn', green: 'ok' };

const BROWSER_INFO = {
    'chrome.exe':   { short: 'Chrome',   private: 'Incognito',        key: 'chrome'  },
    'msedge.exe':   { short: 'Edge',     private: 'InPrivate',        key: 'edge'    },
    'firefox.exe':  { short: 'Firefox',  private: 'Private Browsing', key: 'firefox' },
    'brave.exe':    { short: 'Brave',    private: 'Private',          key: 'brave'   },
    'opera.exe':    { short: 'Opera',    private: 'Private',          key: 'opera'   },
    'vivaldi.exe':  { short: 'Vivaldi',  private: 'Private',          key: 'vivaldi' },
    'iexplore.exe': { short: 'IE',       private: 'InPrivate',        key: 'ie'      },
    'arc.exe':      { short: 'Arc',      private: 'Incognito',        key: 'arc'     },
    'safari':       { short: 'Safari',   private: 'Private',          key: 'safari'  },
};

function browserInfo(name) {
    if (!name) return { short: '—', private: 'Private', key: 'generic' };
    const n = name.toLowerCase();
    if (BROWSER_INFO[n]) return BROWSER_INFO[n];
    for (const key of Object.keys(BROWSER_INFO)) {
        if (n.includes(key.replace('.exe', ''))) return BROWSER_INFO[key];
    }
    return { short: name, private: 'Private', key: 'generic' };
}

const CATEGORY_META = {
    processes: {
        title: 'Processes',
        icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="4" y="4" width="16" height="16" rx="2"/><line x1="9" y1="9" x2="15" y2="9"/><line x1="9" y1="13" x2="15" y2="13"/></svg>',
    },
    windows: {
        title: 'Visible windows',
        icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="3" width="18" height="18" rx="2"/><line x1="3" y1="9" x2="21" y2="9"/></svg>',
    },
    tabs: {
        title: 'Browser tabs',
        icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="9"/><path d="M3 12h18M12 3a15 15 0 010 18M12 3a15 15 0 000 18"/></svg>',
    },
    devices: {
        title: 'Audio & video devices',
        icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="7" width="14" height="10" rx="2"/><polyline points="16 11 22 7 22 17 16 13"/></svg>',
    },
    monitors: {
        title: 'Monitors',
        icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="3" width="20" height="14" rx="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>',
    },
    remote: {
        title: 'Remote access',
        icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M5 12a7 7 0 0114 0"/><path d="M8 12a4 4 0 018 0"/><circle cx="12" cy="12" r="1"/></svg>',
    },
    system: {
        title: 'System state',
        icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41"/><circle cx="12" cy="12" r="4"/></svg>',
    },
};

const DEVICE_ICON = {
    video: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="7" width="14" height="10" rx="2"/><polyline points="16 11 22 7 22 17 16 13"/></svg>',
    audio: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 1v22M5 7v10M19 7v10M1 10v4M23 10v4"/></svg>',
};

const PERIODIC_INTERVAL_MS = 5 * 60 * 1000;
const HISTORY_MAX = 8;

// ============================================================
// State
// ============================================================

const state = {
    consented: false,
    lastScan: null,
    lastScanAt: null,
    activeCategory: null,
    categoryCounts: {},
    history: [], // [{ at, sev, total, crit, warn }]
    scanInProgress: false,
    compactMode: false,
    periodicOn: false,
    compactOnStart: true,
    periodicTimer: null,
    countdownTimer: null,
    nextScanAt: null,
};

// ============================================================
// Helpers
// ============================================================

function $(id) { return document.getElementById(id); }
function $$(sel) { return document.querySelectorAll(sel); }

function getApi() {
    if (window.go && window.go.main && window.go.main.App) {
        return window.go.main.App;
    }
    return null;
}

function escapeHtml(s) {
    if (s === null || s === undefined) return '';
    return String(s)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}

function fmtTime(d) {
    if (!d) return '—';
    const h = d.getHours();
    const m = d.getMinutes().toString().padStart(2, '0');
    return `${h}:${m}`;
}

function fmtTimestamp(d) {
    if (!d) return '—';
    const h = d.getHours().toString().padStart(2, '0');
    const m = d.getMinutes().toString().padStart(2, '0');
    const s = d.getSeconds().toString().padStart(2, '0');
    return `${h}:${m}:${s}`;
}

function showToast(msg, ms = 2500) {
    const t = $('ds-toast');
    t.textContent = msg;
    t.classList.remove('hidden');
    clearTimeout(showToast._t);
    showToast._t = setTimeout(() => t.classList.add('hidden'), ms);
}

// ============================================================
// Theme
// ============================================================

const THEME_KEY = 'dayside.theme';

function applyTheme(mode) {
    const root = document.documentElement;
    root.classList.remove('theme-dark');
    if (mode === 'dark') root.classList.add('theme-dark');
    $$('.theme-btn').forEach(b => {
        b.classList.toggle('active', b.dataset.theme === mode);
    });
}

function initTheme() {
    const stored = localStorage.getItem(THEME_KEY) === 'dark' ? 'dark' : 'light';
    applyTheme(stored);
    $$('.theme-btn').forEach(b => {
        b.addEventListener('click', () => {
            const mode = b.dataset.theme;
            localStorage.setItem(THEME_KEY, mode);
            applyTheme(mode);
        });
    });
}

// ============================================================
// State machine (ready / scanning / results)
// ============================================================

function showState(name) {
    const ready = $('state-ready');
    const scanning = $('state-scanning');
    const results = $('state-results');
    ready.classList.toggle('active', name === 'ready');
    scanning.classList.toggle('active', name === 'scanning');
    results.classList.toggle('active', name === 'results');
}

// ============================================================
// Consent gate
// ============================================================

function initConsent() {
    const modal    = $('ds-consent-modal');
    const checkbox = $('ds-consent-checkbox');
    const btnAgree = $('ds-consent-agree');
    const btnDecl  = $('ds-consent-decline');

    function show() {
        checkbox.checked = false;
        btnAgree.disabled = true;
        modal.style.display = 'flex';
    }
    function hide() { modal.style.display = 'none'; }

    checkbox.addEventListener('change', () => {
        btnAgree.disabled = !checkbox.checked;
    });

    btnAgree.addEventListener('click', () => {
        if (!checkbox.checked) return;
        hide();
        state.consented = true;
        $('ds-start-btn').disabled = false;
        window.dispatchEvent(new CustomEvent('dayside:consent-granted'));
    });

    btnDecl.addEventListener('click', () => {
        hide();
        const api = getApi();
        if (api && typeof api.Quit === 'function') {
            api.Quit();
        } else {
            document.body.innerHTML =
                '<div style="padding:40px;font-family:system-ui,sans-serif;color:#57534e;">' +
                'Dayside closed. No scan was performed. You can close this window.' +
                '</div>';
        }
    });

    window.addEventListener('DOMContentLoaded', () => {
        setTimeout(show, 100);
    });
    // If DOMContentLoaded already fired, show now.
    if (document.readyState !== 'loading') {
        setTimeout(show, 100);
    }
}

// ============================================================
// Options (toggles)
// ============================================================

function bindToggle(el, onChange) {
    el.addEventListener('click', () => {
        const next = !el.classList.contains('on');
        el.classList.toggle('on', next);
        el.setAttribute('aria-checked', next ? 'true' : 'false');
        onChange(next);
    });
    el.addEventListener('keydown', (e) => {
        if (e.key === ' ' || e.key === 'Enter') {
            e.preventDefault();
            el.click();
        }
    });
}

function initOptions() {
    bindToggle($('ds-toggle-periodic'), (on) => {
        state.periodicOn = on;
        if (on) startPeriodic();
        else stopPeriodic();
    });
    bindToggle($('ds-toggle-compact-start'), (on) => {
        state.compactOnStart = on;
    });
}

// ============================================================
// Periodic scanning
// ============================================================

function startPeriodic() {
    stopPeriodic();
    state.nextScanAt = Date.now() + PERIODIC_INTERVAL_MS;
    state.periodicTimer = setInterval(() => {
        if (state.scanInProgress) return;
        scan({ manual: false });
    }, PERIODIC_INTERVAL_MS);
    startCountdown();
    updateCompactPeriodic();
}

function stopPeriodic() {
    if (state.periodicTimer) { clearInterval(state.periodicTimer); state.periodicTimer = null; }
    if (state.countdownTimer) { clearInterval(state.countdownTimer); state.countdownTimer = null; }
    state.nextScanAt = null;
    updateCompactPeriodic();
}

function startCountdown() {
    if (state.countdownTimer) clearInterval(state.countdownTimer);
    renderCountdown();
    state.countdownTimer = setInterval(renderCountdown, 1000);
}

function renderCountdown() {
    const el = $('ds-compact-countdown');
    const text = $('ds-compact-countdown-text');
    if (!state.periodicOn || !state.nextScanAt) {
        el.classList.add('off');
        text.textContent = 'Periodic scanning off';
        return;
    }
    el.classList.remove('off');
    const secs = Math.max(0, Math.floor((state.nextScanAt - Date.now()) / 1000));
    const m = Math.floor(secs / 60);
    const s = (secs % 60).toString().padStart(2, '0');
    text.textContent = `Next scan in ${m}:${s}`;
}

function updateCompactPeriodic() {
    const root = $('ds-compact-inner');
    if (!root) return;
    root.classList.toggle('periodic', state.periodicOn);
    const fill = $('ds-compact-progress-fill');
    if (!fill) return;
    if (!state.periodicOn || !state.nextScanAt) {
        fill.style.width = '0%';
        return;
    }
    const pct = Math.max(0, Math.min(100,
        ((state.nextScanAt - Date.now()) / PERIODIC_INTERVAL_MS) * 100));
    fill.style.width = `${pct}%`;
}

setInterval(updateCompactPeriodic, 1000);

// ============================================================
// Scan
// ============================================================

async function scan(opts = {}) {
    const manual = opts.manual !== false;
    if (state.scanInProgress) return;
    if (!state.consented) return;

    const api = getApi();
    if (!api) {
        if (manual) showToast('Wails runtime not available.');
        return;
    }

    state.scanInProgress = true;
    showState('scanning');
    cycleScanStages();

    try {
        const result = await api.Scan();
        state.lastScan = result;
        state.lastScanAt = new Date();
        recordHistory(result);
        renderResults(result);
        showState('results');

        // After the first successful scan, if compactOnStart is enabled, auto-enter compact
        if (state.compactOnStart && !state.compactMode && state.history.length === 1) {
            setTimeout(enterCompactMode, 400);
        }

        // Reset periodic timer so it starts from now
        if (state.periodicOn) {
            state.nextScanAt = Date.now() + PERIODIC_INTERVAL_MS;
            startCountdown();
        }

        updateCompact();

        if (result.warnings && result.warnings.length > 0) {
            showToast('Scan finished with warnings');
        }
    } catch (err) {
        console.error(err);
        showToast('Scan failed: ' + err);
        showState('ready');
    } finally {
        state.scanInProgress = false;
    }
}

const SCAN_STAGES = [
    { t: 0,    stage: 'Checking processes & windows…', sub: 'Enumerating running processes' },
    { t: 900,  stage: 'Reading browser tabs…',          sub: 'Checking browsers for open tabs' },
    { t: 1900, stage: 'Checking audio & video…',        sub: 'Enumerating connected devices' },
    { t: 2700, stage: 'Checking system state…',         sub: 'Monitor count and remote session' },
];

let scanStageTimers = [];
function cycleScanStages() {
    scanStageTimers.forEach(clearTimeout);
    scanStageTimers = [];
    const stageEl = $('ds-scan-stage');
    const subEl   = $('ds-scan-sublabel');
    stageEl.textContent = SCAN_STAGES[0].stage;
    subEl.textContent = SCAN_STAGES[0].sub;
    for (let i = 1; i < SCAN_STAGES.length; i++) {
        scanStageTimers.push(setTimeout(() => {
            stageEl.textContent = SCAN_STAGES[i].stage;
            subEl.textContent = SCAN_STAGES[i].sub;
        }, SCAN_STAGES[i].t));
    }
}

// ============================================================
// Aggregate / severity
// ============================================================

function aggregateCounts(s) {
    const counts = { crit: 0, warn: 0, ok: 0 };
    (s.processes || []).forEach(p => {
        const k = SEV_TO_DS[p.severity] || 'ok';
        counts[k]++;
    });
    ((s.devices && s.devices.video) || []).forEach(d => {
        const k = SEV_TO_DS[d.severity] || 'ok';
        counts[k]++;
    });
    ((s.devices && s.devices.audio) || []).forEach(d => {
        const k = SEV_TO_DS[d.severity] || 'ok';
        counts[k]++;
    });
    (s.tabs || []).forEach(t => {
        const k = SEV_TO_DS[t.severity] || 'ok';
        counts[k]++;
    });
    // Remote session is a critical system signal
    if (s.system && s.system.remoteSession) counts.crit++;
    return counts;
}

function categorySeverities(s) {
    const out = {
        processes: { crit: 0, warn: 0, total: 0 },
        windows:   { crit: 0, warn: 0, total: 0 },
        tabs:      { crit: 0, warn: 0, total: 0 },
        devices:   { crit: 0, warn: 0, total: 0 },
        monitors:  { crit: 0, warn: 0, total: 1 },
        remote:    { crit: 0, warn: 0, total: 0 },
        system:    { crit: 0, warn: 0, total: 0 },
    };

    (s.processes || []).forEach(p => {
        out.processes.total++;
        const k = SEV_TO_DS[p.severity];
        if (k === 'crit') out.processes.crit++;
        else if (k === 'warn') out.processes.warn++;

        (p.windows || []).forEach(w => {
            out.windows.total++;
            if (w.cloaked || w.affinity === 0x11) out.windows.crit++;
            else if (w.layered && w.topmost) out.windows.warn++;
        });
    });

    (s.tabs || []).forEach(t => {
        out.tabs.total++;
        const k = SEV_TO_DS[t.severity];
        if (k === 'crit') out.tabs.crit++;
        else if (k === 'warn') out.tabs.warn++;
    });

    const vids = (s.devices && s.devices.video) || [];
    const auds = (s.devices && s.devices.audio) || [];
    out.devices.total = vids.length + auds.length;
    [...vids, ...auds].forEach(d => {
        const k = SEV_TO_DS[d.severity];
        if (k === 'crit') out.devices.crit++;
        else if (k === 'warn') out.devices.warn++;
    });

    out.monitors.total = (s.system && s.system.monitorCount) || 1;

    out.remote.total = (s.system && s.system.remoteSession) ? 1 : 0;
    if (s.system && s.system.remoteSession) out.remote.crit = 1;

    // System has a few data points worth surfacing.
    out.system.total = 4;

    return out;
}

function pickDefaultCategory(cats) {
    const order = ['processes', 'tabs', 'windows', 'devices', 'remote', 'monitors', 'system'];
    for (const k of order) if (cats[k] && cats[k].crit > 0) return k;
    for (const k of order) if (cats[k] && cats[k].warn > 0) return k;
    return 'processes';
}

// ============================================================
// Render results
// ============================================================

function renderResults(s) {
    const counts = aggregateCounts(s);
    const cats = categorySeverities(s);
    state.categoryCounts = cats;

    // Timestamp + KPIs
    const totalItems =
        cats.processes.total + cats.tabs.total + cats.devices.total +
        cats.windows.total + cats.remote.total + cats.system.total;
    const findings = counts.crit + counts.warn;
    $('ds-results-timestamp').textContent =
        `${fmtTimestamp(state.lastScanAt)} · ${findings} finding${findings === 1 ? '' : 's'} from ${totalItems} items checked`;

    $('ds-kpi-crit').textContent = counts.crit;
    $('ds-kpi-warn').textContent = counts.warn;
    $('ds-kpi-ok').textContent   = Math.max(0, totalItems - counts.crit - counts.warn);

    $('ds-kpi-crit-sub').innerHTML = counts.crit === 0
        ? 'No critical findings.'
        : 'Flagged items need review. <strong>Discuss with interviewer.</strong>';
    $('ds-kpi-warn-sub').innerHTML = counts.warn === 0
        ? 'No warnings.'
        : `${counts.warn} item${counts.warn === 1 ? '' : 's'} worth a closer look.`;
    $('ds-kpi-ok-sub').textContent = 'Items checked across processes, tabs, devices, and system.';

    // Title dot + summary
    const worstSev = counts.crit > 0 ? 'crit' : (counts.warn > 0 ? 'warn' : 'ok');
    const titleDot = $('ds-title-dot');
    titleDot.classList.remove('crit', 'warn', 'muted');
    if (worstSev !== 'ok') titleDot.classList.add(worstSev);

    const okItems = Math.max(0, totalItems - counts.crit - counts.warn);
    $('ds-summary-big').textContent = okItems;
    $('ds-summary-of').textContent = `of ${totalItems} clean`;

    const flaggedCats = Object.entries(cats).filter(([k, v]) => v.crit + v.warn > 0).length;
    const totalCats = Object.keys(cats).length;
    $('ds-summary-tagline').textContent = flaggedCats === 0
        ? `All ${totalCats} categories clean.`
        : `${totalCats - flaggedCats} categories clean · ${flaggedCats} flagged.`;

    // Sidebar counts + dots
    Object.entries(cats).forEach(([cat, v]) => {
        const dot = document.querySelector(`[data-dot="${cat}"]`);
        const count = document.querySelector(`[data-count="${cat}"]`);
        if (dot) {
            dot.classList.remove('crit', 'warn', 'ok', 'muted');
            if (v.crit > 0) dot.classList.add('crit');
            else if (v.warn > 0) dot.classList.add('warn');
            else if (v.total > 0) dot.classList.add('ok');
            else dot.classList.add('muted');
        }
        if (count) {
            if (cat === 'monitors') count.textContent = v.total;
            else if (cat === 'system') count.textContent = '—';
            else count.textContent = (v.crit + v.warn) || v.total;
        }
    });

    // Render each category pane
    renderProcessesPane(s, cats.processes);
    renderWindowsPane(s, cats.windows);
    renderTabsPane(s, cats.tabs);
    renderDevicesPane(s, cats.devices);
    renderMonitorsPane(s);
    renderRemotePane(s);
    renderSystemPane(s);

    // Auto-activate the most-flagged category
    activateCategory(pickDefaultCategory(cats));
}

function activateCategory(catId) {
    state.activeCategory = catId;
    $$('.ds-category').forEach(sec => {
        sec.classList.remove('active', 'crit', 'warn', 'ok');
        if (sec.dataset.cat === catId) {
            sec.classList.add('active');
            const v = state.categoryCounts[catId] || {};
            if (v.crit > 0) sec.classList.add('crit');
            else if (v.warn > 0) sec.classList.add('warn');
            else sec.classList.add('ok');
        }
    });
    $$('.ds-nav-item').forEach(item => {
        item.classList.remove('active', 'crit-active', 'warn-active');
        if (item.dataset.cat === catId) {
            item.classList.add('active');
            const v = state.categoryCounts[catId] || {};
            if (v.crit > 0) item.classList.add('crit-active');
            else if (v.warn > 0) item.classList.add('warn-active');
        }
    });
}

function initSidebarNav() {
    $$('.ds-nav-item').forEach(item => {
        item.addEventListener('click', (e) => {
            e.preventDefault();
            activateCategory(item.dataset.cat);
        });
    });
}

// ============================================================
// Pane renderers
// ============================================================

function paneHeaderHTML(cat, counts, metaStr) {
    const meta = CATEGORY_META[cat];
    let statusLabel, statusClass;
    if (counts.crit > 0) {
        statusLabel = `${counts.crit} critical`;
        statusClass = 'crit';
    } else if (counts.warn > 0) {
        statusLabel = `${counts.warn} warning${counts.warn === 1 ? '' : 's'}`;
        statusClass = 'warn';
    } else {
        statusLabel = 'Clean';
        statusClass = 'ok';
    }
    return `
        <div class="ds-pane-header">
            <div class="ds-pane-icon">${meta.icon}</div>
            <div class="ds-pane-title">${meta.title}</div>
            <div class="ds-pane-meta">${escapeHtml(metaStr || '')}</div>
            <div class="ds-pane-status">${escapeHtml(statusLabel)}</div>
        </div>
    `;
}

function renderProcessesPane(s, counts) {
    const root = $('cat-processes');
    const procs = [...(s.processes || [])].sort((a, b) => {
        const order = { red: 0, yellow: 1, green: 2 };
        if (order[a.severity] !== order[b.severity]) return order[a.severity] - order[b.severity];
        return (a.name || '').localeCompare(b.name || '');
    });

    const header = paneHeaderHTML('processes', counts, `${procs.length} total`);

    const flagged = procs.filter(p => p.severity !== 'green');
    const clean = procs.filter(p => p.severity === 'green');

    let body = '';
    if (flagged.length > 0) {
        body += `<div class="ds-list-subheader">Flagged<span class="ds-list-subheader-count">${flagged.length}</span></div>`;
        body += flagged.map(renderProcessRow).join('');
    }
    if (clean.length > 0) {
        body += `<div class="ds-list-subheader">Clean processes<span class="ds-list-subheader-count">${clean.length}</span></div>`;
        body += clean.map(p => renderProcessRow(p, true)).join('');
    }
    if (procs.length === 0) {
        body = '<div class="ds-category-empty">No processes enumerated.</div>';
    }

    root.innerHTML = header + `<div class="ds-category-body">${body}</div>`;
}

function renderProcessRow(p, dense = false) {
    const sev = SEV_TO_DS[p.severity] || 'ok';
    const flags = p.flags || [];
    const primaryFlag = flags[0];
    const label = primaryFlag ? (FLAG_REASON[primaryFlag] || primaryFlag.toLowerCase()) : 'clean';
    const badge = primaryFlag
        ? `<span class="ds-list-badge">${escapeHtml(label)}</span>`
        : `<span class="ds-list-badge">ok</span>`;

    const meta = `PID ${p.pid}`;
    const rowCls = `ds-list-row ${sev}${dense ? ' dense' : ''}`;
    let out = `
        <div class="${rowCls}">
            <div class="ds-list-dot"></div>
            <div class="ds-list-primary">${escapeHtml(p.name || '(unknown)')}</div>
            <div class="ds-list-meta">${escapeHtml(meta)}</div>
            ${badge}
        </div>
    `;

    if (flags.length > 0) {
        const labels = flags.map(f => FLAG_LABELS[f] || f);
        const hint = FLAG_HINT[flags[0]];
        const pathLine = p.path ? `<span class="hint">Path: <code>${escapeHtml(p.path)}</code></span>` : '';
        out += `
            <div class="ds-list-reason ${sev}">
                <strong>${escapeHtml(labels.join(' · '))}</strong>
                ${hint ? ` — ${escapeHtml(hint)}` : ''}
                ${pathLine}
            </div>
        `;
    }
    return out;
}

function renderWindowsPane(s, counts) {
    const root = $('cat-windows');
    const entries = [];
    (s.processes || []).forEach(p => {
        (p.windows || []).forEach(w => {
            entries.push({ proc: p, window: w });
        });
    });

    const flagged = entries.filter(({ window: w }) =>
        w.cloaked || w.affinity === 0x11 || w.affinity === 0x01 || (w.topmost && w.layered));
    const clean = entries.filter(e => !flagged.includes(e));

    const header = paneHeaderHTML('windows', counts, `${entries.length} visible`);
    let body = '';

    if (flagged.length > 0) {
        body += `<div class="ds-list-subheader">Flagged<span class="ds-list-subheader-count">${flagged.length}</span></div>`;
        body += flagged.map(renderWindowRow).join('');
    }
    if (clean.length > 0) {
        body += `<div class="ds-list-subheader">Visible windows<span class="ds-list-subheader-count">${clean.length}</span></div>`;
        body += clean.map(renderWindowRow).join('');
    }
    if (entries.length === 0) {
        body = '<div class="ds-category-empty">No visible windows enumerated.</div>';
    }

    root.innerHTML = header + `<div class="ds-category-body">${body}</div>`;
}

function renderWindowRow({ proc, window: w }) {
    let sev = 'ok';
    const flags = [];
    if (w.affinity === 0x11) { sev = 'crit'; flags.push('excluded from capture'); }
    else if (w.affinity === 0x01) { sev = 'warn'; flags.push('monitor-only'); }
    if (w.cloaked) { sev = 'crit'; flags.push('cloaked'); }
    if (w.layered && w.topmost) { if (sev === 'ok') sev = 'warn'; flags.push('layered + topmost'); }
    else if (w.topmost) flags.push('topmost');

    const title = w.title || '(untitled)';
    const badge = flags.length > 0
        ? `<span class="ds-list-badge">${escapeHtml(flags[0])}</span>`
        : '<span class="ds-list-badge">ok</span>';

    let out = `
        <div class="ds-list-row ${sev}">
            <div class="ds-list-dot"></div>
            <div class="ds-list-primary">${escapeHtml(title)}</div>
            <div class="ds-list-meta">${escapeHtml(proc.name || '')}</div>
            ${badge}
        </div>
    `;
    if (sev !== 'ok' && flags.length > 0) {
        const hint = w.affinity === 0x11
            ? 'Windows flag WDA_EXCLUDEFROMCAPTURE is set — window is hidden from screen-share.'
            : w.cloaked ? 'Window is reported as cloaked by the OS.'
            : 'Layered + topmost window is commonly used by overlay tools.';
        out += `
            <div class="ds-list-reason ${sev}">
                <strong>${escapeHtml(flags.join(' · '))}</strong> — ${escapeHtml(hint)}
            </div>
        `;
    }
    return out;
}

function renderTabsPane(s, counts) {
    const root = $('cat-tabs');
    const tabs = [...(s.tabs || [])].sort((a, b) => {
        const order = { red: 0, yellow: 1, green: 2 };
        if (order[a.severity] !== order[b.severity]) return order[a.severity] - order[b.severity];
        return (a.title || '').localeCompare(b.title || '');
    });

    const perBrowser = {};
    tabs.forEach(t => {
        const k = browserInfo(t.browserName).short;
        perBrowser[k] = (perBrowser[k] || 0) + 1;
    });
    const summary = Object.entries(perBrowser).map(([k, v]) => `${k}:${v}`).join(' · ');

    const header = paneHeaderHTML('tabs', counts, summary || `${tabs.length} total`);

    const flagged = tabs.filter(t => t.severity !== 'green');
    const clean = tabs.filter(t => t.severity === 'green');

    let body = '';
    if (flagged.length > 0) {
        body += `<div class="ds-list-subheader">Flagged<span class="ds-list-subheader-count">${flagged.length}</span></div>`;
        body += flagged.map(renderTabRow).join('');
    }
    if (clean.length > 0) {
        body += `<div class="ds-list-subheader">Other open tabs<span class="ds-list-subheader-count">${clean.length}</span></div>`;
        body += clean.map(renderTabRow).join('');
    }
    if (tabs.length === 0) {
        body = '<div class="ds-category-empty">No browser tabs enumerated. On macOS, first run prompts for automation permission — grant it and re-scan.</div>';
    }

    root.innerHTML = header + `<div class="ds-category-body">${body}</div>`;
}

function renderTabRow(t) {
    const sev = SEV_TO_DS[t.severity] || 'ok';
    const info = browserInfo(t.browserName);
    const kind = t.kind || 'tab';
    const title = t.title || '(untitled)';

    const badges = [];
    if (kind === 'panel')
        badges.push(`<span class="ds-list-badge kind-panel" title="Embedded AI side-panel inside the browser window">side panel</span>`);
    else if (kind === 'popout')
        badges.push(`<span class="ds-list-badge kind-popout" title="Popped-out browser chat window (no tab strip)">pop-out</span>`);
    if (sev === 'crit') badges.push(`<span class="ds-list-badge">flagged</span>`);
    else if (sev === 'warn') badges.push(`<span class="ds-list-badge">watch</span>`);
    if (t.incognito) {
        badges.push(`<span class="ds-list-badge incognito browser-${info.key}" title="${escapeHtml(info.short)} private browsing window">${escapeHtml(info.private)}</span>`);
    }
    if (badges.length === 0) badges.push(`<span class="ds-list-badge">ok</span>`);

    let out = `
        <div class="ds-list-row ${sev}">
            <div class="ds-list-dot"></div>
            <div class="ds-list-primary">${escapeHtml(title)}</div>
            <div class="ds-list-meta">${escapeHtml(info.short)}</div>
            <div style="display:flex; gap:4px;">${badges.join('')}</div>
        </div>
    `;

    const lines = [];
    if (t.url) lines.push(`<span class="hint"><code>${escapeHtml(t.url)}</code></span>`);
    if (t.reason) lines.push(`<strong>${escapeHtml(t.reason)}</strong>`);
    if (kind === 'panel')
        lines.unshift('<strong>Side panel</strong> — this is an embedded panel inside the browser, not a normal tab.');
    else if (kind === 'popout')
        lines.unshift('<strong>Pop-out window</strong> — a separate browser chat window with no tab strip.');

    if (sev !== 'ok' && lines.length > 0) {
        out += `<div class="ds-list-reason ${sev}">${lines.join('<br>')}</div>`;
    }
    return out;
}

function renderDevicesPane(s, counts) {
    const root = $('cat-devices');
    const vids = (s.devices && s.devices.video) || [];
    const auds = (s.devices && s.devices.audio) || [];

    const header = paneHeaderHTML('devices', counts, `${vids.length} video · ${auds.length} audio`);
    let body = '';

    if (vids.length > 0) {
        body += `<div class="ds-list-subheader">Video devices<span class="ds-list-subheader-count">${vids.length}</span></div>`;
        body += vids.map(d => renderDeviceRow(d, 'video')).join('');
    }
    if (auds.length > 0) {
        body += `<div class="ds-list-subheader">Audio devices<span class="ds-list-subheader-count">${auds.length}</span></div>`;
        body += auds.map(d => renderDeviceRow(d, 'audio')).join('');
    }
    if (vids.length === 0 && auds.length === 0) {
        body = '<div class="ds-category-empty">No audio or video devices detected.</div>';
    }

    root.innerHTML = header + `<div class="ds-category-body">${body}</div>`;
}

function renderDeviceRow(d, kind) {
    const sev = SEV_TO_DS[d.severity] || 'ok';
    const badge = sev === 'ok'
        ? '<span class="ds-detail-badge">ok</span>'
        : `<span class="ds-detail-badge ${sev}">${sev === 'crit' ? 'flagged' : 'watch'}</span>`;
    let out = `
        <div class="ds-detail-row">
            <div class="ds-detail-icon">${DEVICE_ICON[kind] || ''}</div>
            <div class="ds-detail-body">
                <div class="ds-detail-name">${escapeHtml(d.name || '(unknown)')}</div>
                <div class="ds-detail-meta">${escapeHtml(kind)}${d.reason ? ' · ' + escapeHtml(d.reason) : ''}</div>
            </div>
            ${badge}
        </div>
    `;
    return out;
}

function renderMonitorsPane(s) {
    const root = $('cat-monitors');
    const count = (s.system && s.system.monitorCount) || 1;
    const counts = { crit: 0, warn: 0 };
    const header = paneHeaderHTML('monitors', counts, `${count} connected`);
    const body = `
        <div class="ds-detail-list">
            <div class="ds-detail-row">
                <div class="ds-detail-icon">${CATEGORY_META.monitors.icon}</div>
                <div class="ds-detail-body">
                    <div class="ds-detail-name">${count} monitor${count === 1 ? '' : 's'} detected</div>
                    <div class="ds-detail-meta">${count > 1 ? 'Multiple displays — ensure the interviewer sees the one you intend.' : 'Single display configuration.'}</div>
                </div>
                <span class="ds-detail-badge">ok</span>
            </div>
        </div>
    `;
    root.innerHTML = header + `<div class="ds-category-body">${body}</div>`;
}

function renderRemotePane(s) {
    const root = $('cat-remote');
    const rdp = !!(s.system && s.system.remoteSession);
    const counts = { crit: rdp ? 1 : 0, warn: 0 };
    const header = paneHeaderHTML('remote', counts, rdp ? 'active' : 'none');
    let body;
    if (rdp) {
        body = `
            <div class="ds-list-row crit">
                <div class="ds-list-dot"></div>
                <div class="ds-list-primary">Remote desktop session active</div>
                <div class="ds-list-meta">RDP</div>
                <span class="ds-list-badge">remote-session</span>
            </div>
            <div class="ds-list-reason crit">
                <strong>You appear to be connected to this machine over RDP.</strong>
                Remote input is possible. Confirm with the interviewer whether this is expected.
            </div>
        `;
    } else {
        body = `
            <div class="ds-list-row ok">
                <div class="ds-list-dot"></div>
                <div class="ds-list-primary">No remote desktop session detected</div>
                <div class="ds-list-meta">local</div>
                <span class="ds-list-badge">ok</span>
            </div>
            <div class="ds-category-empty">
                Dayside also flags known remote-access tools (TeamViewer, AnyDesk, etc.) in the Processes category.
            </div>
        `;
    }
    root.innerHTML = header + `<div class="ds-category-body">${body}</div>`;
}

function renderSystemPane(s) {
    const root = $('cat-system');
    const sys = s.system || {};
    const counts = { crit: 0, warn: 0 };
    const header = paneHeaderHTML('system', counts, sys.hostName || 'host');
    const body = `
        <div class="ds-detail-list">
            <div class="ds-detail-row">
                <div class="ds-detail-icon">${CATEGORY_META.system.icon}</div>
                <div class="ds-detail-body">
                    <div class="ds-detail-name">Hostname</div>
                    <div class="ds-detail-meta">${escapeHtml(sys.hostName || '(unknown)')}</div>
                </div>
            </div>
            <div class="ds-detail-row">
                <div class="ds-detail-icon">${CATEGORY_META.system.icon}</div>
                <div class="ds-detail-body">
                    <div class="ds-detail-name">User</div>
                    <div class="ds-detail-meta">${escapeHtml(sys.userName || '(unknown)')}</div>
                </div>
            </div>
            <div class="ds-detail-row">
                <div class="ds-detail-icon">${CATEGORY_META.monitors.icon}</div>
                <div class="ds-detail-body">
                    <div class="ds-detail-name">Monitors</div>
                    <div class="ds-detail-meta">${sys.monitorCount || 1} connected</div>
                </div>
            </div>
            <div class="ds-detail-row">
                <div class="ds-detail-icon">${CATEGORY_META.remote.icon}</div>
                <div class="ds-detail-body">
                    <div class="ds-detail-name">Remote session</div>
                    <div class="ds-detail-meta">${sys.remoteSession ? 'YES — running inside an RDP session' : 'No'}</div>
                </div>
                <span class="ds-detail-badge ${sys.remoteSession ? 'crit' : ''}">${sys.remoteSession ? 'flagged' : 'ok'}</span>
            </div>
        </div>
    `;
    root.innerHTML = header + `<div class="ds-category-body">${body}</div>`;
}

// ============================================================
// History (in-memory, lost on quit)
// ============================================================

function recordHistory(s) {
    const counts = aggregateCounts(s);
    const sev = counts.crit > 0 ? 'crit' : (counts.warn > 0 ? 'warn' : 'ok');
    const total = (s.processes || []).length + (s.tabs || []).length +
        ((s.devices && s.devices.video) || []).length +
        ((s.devices && s.devices.audio) || []).length;
    state.history.push({
        at: new Date(),
        sev, total,
        crit: counts.crit,
        warn: counts.warn,
        firstPreview: pickFirstFinding(s),
    });
    if (state.history.length > HISTORY_MAX) {
        state.history = state.history.slice(-HISTORY_MAX);
    }
}

function pickFirstFinding(s) {
    const flaggedProc = (s.processes || []).find(p => p.severity === 'red')
        || (s.processes || []).find(p => p.severity === 'yellow');
    if (flaggedProc) {
        const f = (flaggedProc.flags || [])[0];
        return {
            name: flaggedProc.name,
            reason: f ? (FLAG_LABELS[f] || f) : 'Flagged process',
            kind: 'process',
        };
    }
    const flaggedTab = (s.tabs || []).find(t => t.severity === 'red')
        || (s.tabs || []).find(t => t.severity === 'yellow');
    if (flaggedTab) {
        return {
            name: flaggedTab.title || '(untitled tab)',
            reason: flaggedTab.reason || (flaggedTab.kind === 'panel' ? 'AI side panel' : 'Flagged tab'),
            kind: 'tab',
        };
    }
    if (s.system && s.system.remoteSession) {
        return { name: 'Remote desktop session', reason: 'Running inside RDP', kind: 'system' };
    }
    return null;
}

// ============================================================
// Compact view
// ============================================================

async function enterCompactMode() {
    state.compactMode = true;
    document.body.classList.add('compact-mode');
    updateCompact();
    const api = getApi();
    if (api && api.EnterCompactMode) {
        try { await api.EnterCompactMode(); } catch (e) { console.error(e); }
    }
}

async function exitCompactMode() {
    state.compactMode = false;
    document.body.classList.remove('compact-mode');
    const api = getApi();
    if (api && api.ExitCompactMode) {
        try { await api.ExitCompactMode(); } catch (e) { console.error(e); }
    }
}

function updateCompact() {
    const inner = $('ds-compact-inner');
    if (!inner) return;

    inner.classList.remove('ok', 'warn', 'crit', 'neutral');

    const s = state.lastScan;
    if (!s) {
        inner.classList.add('neutral');
        $('ds-compact-title').textContent = 'Ready';
        $('ds-compact-sub').textContent = 'Click Scan now to check';
    } else {
        const counts = aggregateCounts(s);
        let sev = 'ok';
        if (counts.crit > 0) sev = 'crit';
        else if (counts.warn > 0) sev = 'warn';
        inner.classList.add(sev);

        if (sev === 'crit') {
            $('ds-compact-title').textContent = 'Needs review';
            $('ds-compact-sub').textContent = '';
        } else if (sev === 'warn') {
            $('ds-compact-title').textContent = 'Some warnings';
            $('ds-compact-sub').textContent = '';
        } else {
            $('ds-compact-title').textContent = 'All clear';
            $('ds-compact-sub').textContent = 'No flags found';
        }
    }

    updateCompactHeroes();
    renderHistoryTimeline();
    renderCountdown();
    updateCompactPeriodic();
}

function updateCompactHeroes() {
    const critCard = $('ds-compact-hero-crit');
    const warnCard = $('ds-compact-hero-warn');
    const okCard = $('ds-compact-hero-ok');
    const critV = $('ds-compact-hero-crit-value');
    const warnV = $('ds-compact-hero-warn-value');
    const okV = $('ds-compact-hero-ok-value');

    const s = state.lastScan;
    if (!s) {
        critV.textContent = '0';
        warnV.textContent = '0';
        okV.textContent = '0';
        critCard.classList.remove('crit');
        warnCard.classList.remove('warn');
        okCard.classList.remove('ok');
        return;
    }

    const counts = aggregateCounts(s);
    const cats = categorySeverities(s);
    const totalItems =
        cats.processes.total + cats.tabs.total + cats.devices.total +
        cats.windows.total + cats.remote.total + cats.system.total;
    const okCount = Math.max(0, totalItems - counts.crit - counts.warn);

    critV.textContent = String(counts.crit);
    warnV.textContent = String(counts.warn);
    okV.textContent = String(okCount);

    critCard.classList.toggle('crit', counts.crit > 0);
    warnCard.classList.toggle('warn', counts.warn > 0);
    okCard.classList.toggle('ok', okCount > 0);
}

function renderHistoryTimeline() {
    const timeline = $('ds-history-timeline');
    const count = $('ds-history-count');
    const first = $('ds-history-first');

    if (state.history.length === 0) {
        // Render empty placeholder slots for visual consistency.
        timeline.innerHTML = Array.from({ length: HISTORY_MAX }).map(() =>
            `<div class="ds-history-tick" style="height: 35%;"><div class="ds-history-tick-bar"></div></div>`
        ).join('');
        count.textContent = `no scans yet`;
        first.textContent = '—';
        return;
    }

    const allClean = state.history.every(h => h.sev === 'ok');
    count.textContent = allClean
        ? `last ${state.history.length} · all clean`
        : `last ${state.history.length}`;

    first.textContent = fmtTime(state.history[0].at);

    const emptySlots = Math.max(0, HISTORY_MAX - state.history.length);
    const slots = [];

    for (let i = 0; i < emptySlots; i++) {
        slots.push(`<div class="ds-history-tick" style="height: 35%;"><div class="ds-history-tick-bar"></div></div>`);
    }

    state.history.forEach((h, i) => {
        const isCurrent = i === state.history.length - 1;
        const color = h.sev === 'crit' ? 'var(--crit)' : h.sev === 'warn' ? 'var(--warn)' : 'var(--ok)';
        const heightPct = h.sev === 'crit' ? 90 : h.sev === 'warn' ? 70 : 55;
        const tip = `${fmtTime(h.at)} · ${h.sev === 'ok' ? 'clean' : h.sev === 'warn' ? 'warnings' : 'critical'}${isCurrent ? ' · now' : ''}`;
        slots.push(`
            <div class="ds-history-tick ${h.sev} ${isCurrent ? 'current' : ''}" style="height: ${heightPct}%; color: ${color};">
                <div class="ds-history-tick-bar"></div>
                <div class="ds-history-tick-tip">${escapeHtml(tip)}</div>
            </div>
        `);
    });

    timeline.innerHTML = slots.join('');
}

// ============================================================
// Wiring
// ============================================================

document.addEventListener('DOMContentLoaded', () => {
    initTheme();
    initConsent();
    initOptions();
    initSidebarNav();

    $('ds-start-btn').addEventListener('click', () => scan({ manual: true }));
    $('ds-rescan-btn').addEventListener('click', () => scan({ manual: true }));
    $('ds-compact-btn').addEventListener('click', enterCompactMode);
    $('ds-compact-expand').addEventListener('click', exitCompactMode);
    $('ds-compact-scan').addEventListener('click', () => scan({ manual: true }));

    document.addEventListener('keydown', (e) => {
        if (e.target && (e.target.tagName === 'INPUT' || e.target.tagName === 'TEXTAREA')) return;
        if (e.key === 'r' || e.key === 'R') {
            if (state.consented) {
                e.preventDefault();
                scan({ manual: true });
            }
        }
    });

    showState('ready');
    updateCompact();
});
