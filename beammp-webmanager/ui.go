package main

const indexHTML = `<!doctype html>
<html lang="de">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>BeamMP Tool</title>
<style>
  :root {
    --bg: #151a22; --bar: #1d2531; --line: #2c3646; --card: #1b232f;
    --fg: #e6e9ee; --dim: #93a1b5; --dimmer: #7b899e; --sub: #8c9ab0;
    --btn: #2b3a52; --btnline: #3d4f6d; --btnhov: #35486a;
    --sel: #35486a; --selline: #4a5f85; --field: #111721; --fieldline: #38445a;
    --green: #2f6f47; --greenline: #3d8c5b; --greenhov: #388554;
    --red: #6a2f2f; --redline: #8c3d3d; --redhov: #853838;
    --rowline: #232c3a; --tag: #26313f; --tagfg: #9fb0c6;
    --on: #24462f; --onfg: #8fdca8; --off: #45272a; --offfg: #e59a9a;
    --okbg: #1d3a28; --okfg: #a9e6bf; --errbg: #3d2226; --errfg: #f0b0b0;
    --thumb: #121822; --shadow: #0008;
  }
  :root[data-theme="light"] {
    --bg: #f2f4f8; --bar: #ffffff; --line: #d8dee8; --card: #ffffff;
    --fg: #1b2430; --dim: #5a6879; --dimmer: #78859a; --sub: #63718a;
    --btn: #e3e8f0; --btnline: #c3ccdb; --btnhov: #d3dae6;
    --sel: #d3dcec; --selline: #a9b8d0; --field: #ffffff; --fieldline: #c3ccdb;
    --green: #2f8552; --greenline: #2f8552; --greenhov: #276e44;
    --red: #b0463f; --redline: #b0463f; --redhov: #983b35;
    --rowline: #e5eaf2; --tag: #e7ecf4; --tagfg: #55637a;
    --on: #d8f0e0; --onfg: #1e6b3c; --off: #f7dedd; --offfg: #99332c;
    --okbg: #dcf2e4; --okfg: #1d6b3d; --errbg: #fbe0e0; --errfg: #96332c;
    --thumb: #e8ecf3; --shadow: #00000026;
  }
  * { box-sizing: border-box; }
  body { margin: 0; font: 14px/1.45 system-ui, "Segoe UI", sans-serif;
         background: var(--bg); color: var(--fg); }
  header { padding: 14px 18px; background: var(--bar); border-bottom: 1px solid var(--line); }
  .top { display: flex; align-items: center; gap: 10px; margin-bottom: 10px; }
  h1 { margin: 0; font-size: 17px; font-weight: 600; flex: 1; }
  .row { display: flex; gap: 8px; flex-wrap: wrap; align-items: center; }
  input[type=text] { flex: 1 1 320px; min-width: 200px; padding: 7px 10px;
    background: var(--field); border: 1px solid var(--fieldline); border-radius: 5px;
    color: inherit; }
  button { padding: 7px 14px; background: var(--btn); color: var(--fg); cursor: pointer;
    border: 1px solid var(--btnline); border-radius: 5px; font-size: 13px; }
  button:hover { background: var(--btnhov); }
  button.primary { background: var(--green); border-color: var(--greenline); color: #fff; }
  button.primary:hover { background: var(--greenhov); }
  button.off { background: var(--red); border-color: var(--redline); color: #fff; }
  button.off:hover { background: var(--redhov); }
  button.small { padding: 5px 10px; font-size: 12.5px; }
  .meta { margin-top: 8px; color: var(--dim); font-size: 12.5px; }
  .meta b { color: var(--fg); font-weight: 600; }
  nav { display: flex; gap: 4px; padding: 10px 18px 0; flex-wrap: wrap; }
  nav button { border-bottom-left-radius: 0; border-bottom-right-radius: 0; }
  nav button.sel { background: var(--sel); border-color: var(--selline); }
  main { padding: 12px 18px 40px; }
  table { width: 100%; border-collapse: collapse; }
  th, td { text-align: left; padding: 7px 10px; border-bottom: 1px solid var(--rowline); }
  th { color: var(--sub); font-size: 12px; text-transform: uppercase; letter-spacing: .04em; }
  td.act { width: 1%; white-space: nowrap; }
  .tag { display: inline-block; padding: 1px 7px; border-radius: 20px; font-size: 11.5px;
    background: var(--tag); color: var(--tagfg); }
  .tag.on { background: var(--on); color: var(--onfg); }
  .tag.off { background: var(--off); color: var(--offfg); }
  .sub { color: var(--dimmer); font-size: 12px; }
  #msg { position: fixed; left: 0; right: 0; bottom: 0; z-index: 20;
    padding: 10px 18px; display: none; white-space: pre-line;
    border-top: 1px solid var(--line); box-shadow: 0 -2px 10px var(--shadow); }
  #msg.show { display: block; }
  body { padding-bottom: 52px; }
  #msg.ok { background: var(--okbg); color: var(--okfg); }
  #msg.err { background: var(--errbg); color: var(--errfg); }
  .empty { padding: 24px 4px; color: var(--dimmer); }
  .grid { display: grid; gap: 12px;
    grid-template-columns: repeat(auto-fill, minmax(210px, 1fr)); }
  .card { background: var(--card); border: 1px solid var(--line); border-radius: 8px;
    overflow: hidden; display: flex; flex-direction: column; }
  .card.on { border-color: var(--greenline); box-shadow: 0 0 0 1px var(--greenline) inset; }
  .thumb { position: relative; aspect-ratio: 16 / 9; background: var(--thumb); }
  .pic { position: relative; width: 100%; height: 100%; }
  .pic img { position: absolute; inset: 0; width: 100%; height: 100%;
    object-fit: cover; display: block; }
  .ph { width: 100%; height: 100%; display: flex; align-items: center;
    justify-content: center; text-align: center; padding: 8px; font-weight: 600;
    font-size: 15px; color: #f0f4fa; text-shadow: 0 1px 3px var(--shadow); }
  .badge { position: absolute; top: 6px; right: 6px; padding: 2px 8px; z-index: 1;
    border-radius: 20px; font-size: 11.5px; background: var(--on); color: var(--onfg); }
  .card .body { padding: 9px 10px; display: flex; flex-direction: column; gap: 6px; }
  .card .name { font-weight: 600; line-height: 1.25; }
  .card button { width: 100%; }
  td.ico { width: 44px; padding-right: 0; }
  .ico .pic { width: 36px; height: 36px; border-radius: 6px; overflow: hidden; }
  .ico .ph { font-size: 13px; padding: 0; }
  .logbar { display: flex; gap: 10px; align-items: center; flex-wrap: wrap;
    margin-bottom: 8px; color: var(--dim); font-size: 12.5px; }
  .logbar label { display: flex; gap: 5px; align-items: center; cursor: pointer; }
  #reconnect { position: fixed; inset: 0; display: flex; align-items: center;
    justify-content: center; background: var(--bg); z-index: 10; }
  #reconnect[hidden] { display: none; }
  #reconnect .box { text-align: center; max-width: 460px; padding: 24px; }
  #reconnect .sub { margin-top: 8px; }
  .spin { width: 26px; height: 26px; margin: 0 auto 14px; border-radius: 50%;
    border: 3px solid var(--btnline); border-top-color: var(--greenline);
    animation: spin 0.9s linear infinite; }
  @keyframes spin { to { transform: rotate(360deg); } }
  #cfgbox { width: 100%; height: 62vh; padding: 12px; resize: vertical;
    background: var(--field); color: var(--fg); border: 1px solid var(--fieldline);
    border-radius: 6px; white-space: pre; overflow: auto;
    font: 12.5px/1.5 ui-monospace, Menlo, Consolas, monospace; }
  .conrow { display: flex; gap: 8px; margin-top: 8px; }
  .conrow input { flex: 1 1 auto; font-family: ui-monospace, Menlo, Consolas, monospace; }
  pre#logtext, pre#context { margin: 0; padding: 12px; background: var(--field); color: var(--fg);
    border: 1px solid var(--fieldline); border-radius: 6px; height: 62vh;
    overflow: auto; white-space: pre-wrap; word-break: break-word;
    font: 12.5px/1.5 ui-monospace, Menlo, Consolas, monospace; }
  pre#context { height: 56vh; white-space: pre; word-break: normal; overflow-wrap: normal; }
</style>
</head>
<body>
<header>
  <div class="top">
    <h1>BeamMP Server Tool</h1>
    <button class="small" id="langBtn" onclick="toggleLang()"></button>
    <button class="small" id="themeBtn" onclick="toggleTheme()"></button>
  </div>
  <div class="row" id="toolbar">
    <input type="text" id="dir" spellcheck="false">
    <button id="dirBtn" onclick="setDir()"></button>
    <button class="primary" id="reloadBtn" onclick="load()"></button>
    <button id="restartBtn" onclick="doRestart()" hidden></button>
    <button id="restartCancelBtn" onclick="doRestartCancel()" hidden></button>
    <button id="sigBtn" onclick="doSignal()" hidden></button>
    <button id="cmdBtn" onclick="doCmd()" hidden></button>
  </div>
  <div class="meta" id="meta"></div>
  <div class="meta" id="restartInfo" hidden></div>
  <div class="meta" id="sigInfo" hidden></div>
  <div class="meta" id="cmdInfo" hidden></div>
</header>
<div id="msg"></div>
<nav id="tabs"></nav>
<main id="view"></main>
<div id="reconnect" hidden>
  <div class="box">
    <div class="spin"></div>
    <div id="reconnectText"></div>
    <div class="sub" id="reconnectSub"></div>
  </div>
</div>

<script>
const T = {
  de: {
    dirPlaceholder: 'Pfad zum Serververzeichnis', setDir: 'Verzeichnis setzen',
    reload: 'Aktualisieren', maps: 'Maps', mods: 'Mods', plugins: 'Plugins', log: 'Log',
    theme_dark: '\u{1F319} Dunkel', theme_light: '☀ Hell', langName: 'EN',
    server: 'Server', port: 'Port', resources: 'Resources', currentMap: 'Aktuelle Map',
    noConfig: 'ServerConfig.toml nicht gefunden. Pfad prüfen oder Server einmal starten, damit die Datei angelegt wird.',
    set: 'Setzen', running: 'läuft', active: 'aktiv', inactive: 'inaktiv',
    activate: 'Aktivieren', deactivate: 'Deaktivieren',
    colMod: 'Mod', colPlugin: 'Plugin', colStatus: 'Status',
    srcClient: 'Client (aktiv)', srcCustom: 'custom_maps',
    noMaps: 'Keine Maps gefunden.',
    noMods: 'Keine Mods gefunden (Resources/Client bzw. deactivated_mods).',
    noPlugins: 'Keine Plugins gefunden (Resources/Server bzw. deactivated_plugins).',
    logAuto: 'Automatisch aktualisieren', logFile: 'Datei', logSize: 'Größe',
    logTail: 'gezeigt', logBottom: 'Ans Ende', logEmpty: 'Logdatei ist leer.',
    reloaded: 'Neu eingelesen',
    restart: 'Server neu starten', restartCancel: 'Anforderung zurücknehmen',
    reconnectSub: 'Die Oberfläche ist in ein paar Sekunden von selbst wieder da. Dieses Fenster offen lassen.',
    sigBtn: 'Serverprozess beenden (SIGTERM)',
    sigIdle: 'Prozess %s: %s',
    sigNone: 'läuft nicht',
    sigFound: 'PID %s',
    sigWaiting: 'SIGTERM gesendet – warte auf das Ende von PID %s …',
    sigWaitingNew: 'Prozess beendet – warte darauf, dass der Server wieder hochkommt …',
    sigLastNew: 'Neu gestartet um %s, jetzt PID %s',
    sigLastStopped: 'Beendet um %s – es läuft kein Serverprozess',
    proc_signaled: 'SIGTERM an %s gesendet (PID %s)',
    proc_none: 'Kein laufender Prozess namens „%s" gefunden',
    reconnectGone: 'Server beendet – die Oberfläche ist mit runtergegangen …',
    cmdBtn: 'Docker Neustart',
    cmdFailedInfo: 'Letzter Versuch um %s fehlgeschlagen (Exitcode %s): %s',
    cmd_started: 'Neustart eingeleitet',
    cfgSave: 'Speichern', config: 'Config', console: 'Konsole',
    conSend: 'Senden', conPlaceholder: 'Befehl eingeben und Enter',
    conOff: 'Konsole ist nicht aktiv (mit -tmux starten)',
    con_disabled: 'Konsole ist nicht aktiv (mit -tmux starten)',
    con_no_session: 'tmux-Session „%s" nicht gefunden: %s',
    con_send_failed: 'Befehl konnte nicht gesendet werden: %s',
    con_sent: 'Gesendet',
    config_saved: 'ServerConfig.toml gespeichert – Server neu starten, damit es greift',
    config_empty: 'Die Datei darf nicht leer sein',
    config_no_general: 'Hinweis: Der Abschnitt [General] fehlt – der Server kann die Datei so vermutlich nicht lesen.',
    restartWait: 'Neustart eingeleitet – warte, bis der Server wieder da ist …',
    restartWaitSub: 'Die Seite lädt sich von selbst neu, sobald er antwortet.',
    restartDoneMsg: 'Neustart eingeleitet',
    cmd_running: 'Der Befehl läuft bereits',
    cmd_disabled: 'Kein Neustart-Befehl hinterlegt (mit -restartcmd starten)',
    sigFailed: 'PID %s läuft noch – SIGTERM wurde ignoriert oder der Prozess hängt',
    signal_disabled: 'SIGTERM-Funktion ist nicht aktiv (mit -restartsignal starten)',
    restartWaiting: 'Neustart angefordert – warte darauf, dass „%s" verschwindet …',
    restartIdle: 'Neustart-Datei: %s',
    restartLast: 'Zuletzt neu gestartet: %s',
    restart_requested: 'Neustart angefordert – Datei „%s" wurde angelegt',
    restart_already_pending: 'Es liegt bereits eine Neustart-Anforderung vor',
    restart_cancelled: 'Neustart-Anforderung zurückgenommen',
    restart_done: 'Server wurde neu gestartet',
    restart_disabled: 'Neustart-Funktion ist nicht aktiv (mit -restartfile starten)',
    restart_failed: 'Neustart-Datei konnte nicht angelegt werden: %s',
    restart_cancel_failed: 'Neustart-Datei konnte nicht entfernt werden: %s',
    map_set: 'Map gesetzt – Server neu starten, damit es greift',
    no_free_name: 'Für „%s" wurde kein freier Name in custom_maps gefunden',
    map_parked_renamed: 'Hinweis: „%s" lag doppelt vor und wurde als „%s" in custom_maps abgelegt – nichts gelöscht.',
    backup_failed: 'Hinweis: Sicherungskopie der ServerConfig.toml konnte nicht angelegt werden: %s',
    config_write_failed: 'ServerConfig.toml konnte nicht geschrieben werden: %s (Dateien wurden zurückgeschoben)',
    mod_on: 'Mod aktiviert: %s', mod_off: 'Mod deaktiviert: %s',
    plugin_on: 'Plugin aktiviert: %s', plugin_off: 'Plugin deaktiviert: %s',
    dir_changed: 'Verzeichnis gewechselt',
    config_missing: 'ServerConfig.toml nicht gefunden in %s',
    config_unreadable: 'ServerConfig.toml nicht lesbar: %s',
    zip_unreadable: 'Zip nicht lesbar: %s',
    duplicate: '„%s" liegt bereits im Zielordner – bitte eine der beiden Kopien löschen',
    move_failed: '„%s" konnte nicht verschoben werden: %s',
    remove_failed: '%s konnte nicht entfernt werden: %s',
    map_park_failed: '%s konnte nicht nach custom_maps verschoben werden: %s',
    map_activate_failed: '%s konnte nicht aktiviert werden: %s',
    invalid_name: 'Ungültiger Name', no_map: 'Keine Map angegeben',
    dir_not_found: 'Verzeichnis nicht gefunden: %s',
    dir_locked: 'Das Verzeichnis ist fest vorgegeben und lässt sich hier nicht ändern',
    dirFixed: 'Verzeichnis',
    bad_request: 'Ungültige Anfrage',
    log_missing: 'Keine Logdatei in %s gefunden (erwartet: Server.log)',
    log_unreadable: '%s nicht lesbar: %s',
    generic: '%s'
  },
  en: {
    dirPlaceholder: 'Path to the server directory', setDir: 'Set directory',
    reload: 'Refresh', maps: 'Maps', mods: 'Mods', plugins: 'Plugins', log: 'Log',
    theme_dark: '\u{1F319} Dark', theme_light: '☀ Light', langName: 'DE',
    server: 'Server', port: 'Port', resources: 'Resources', currentMap: 'Current map',
    noConfig: 'ServerConfig.toml not found. Check the path, or start the server once so the file gets created.',
    set: 'Set', running: 'running', active: 'active', inactive: 'inactive',
    activate: 'Activate', deactivate: 'Deactivate',
    colMod: 'Mod', colPlugin: 'Plugin', colStatus: 'Status',
    srcClient: 'Client (active)', srcCustom: 'custom_maps',
    noMaps: 'No maps found.',
    noMods: 'No mods found (Resources/Client or deactivated_mods).',
    noPlugins: 'No plugins found (Resources/Server or deactivated_plugins).',
    logAuto: 'Auto refresh', logFile: 'File', logSize: 'Size',
    logTail: 'shown', logBottom: 'To bottom', logEmpty: 'Log file is empty.',
    reloaded: 'Reloaded',
    restart: 'Restart server', restartCancel: 'Withdraw request',
    reconnectSub: 'The interface will come back on its own in a few seconds. Leave this window open.',
    sigBtn: 'Stop server process (SIGTERM)',
    sigIdle: 'Process %s: %s',
    sigNone: 'not running',
    sigFound: 'PID %s',
    sigWaiting: 'SIGTERM sent – waiting for PID %s to exit …',
    sigWaitingNew: 'Process stopped – waiting for the server to come back up …',
    sigLastNew: 'Restarted at %s, now PID %s',
    sigLastStopped: 'Stopped at %s – no server process running',
    proc_signaled: 'SIGTERM sent to %s (PID %s)',
    proc_none: 'No running process named "%s" found',
    reconnectGone: 'Server stopped – the interface went down with it …',
    cmdBtn: 'Docker restart',
    cmdFailedInfo: 'Last attempt at %s failed (exit code %s): %s',
    cmd_started: 'Restart initiated',
    cfgSave: 'Save', config: 'Config', console: 'Console',
    conSend: 'Send', conPlaceholder: 'Type a command and press Enter',
    conOff: 'Console is off (start with -tmux)',
    con_disabled: 'Console is off (start with -tmux)',
    con_no_session: 'tmux session "%s" not found: %s',
    con_send_failed: 'Could not send the command: %s',
    con_sent: 'Sent',
    config_saved: 'ServerConfig.toml saved – restart the server for it to take effect',
    config_empty: 'The file must not be empty',
    config_no_general: 'Note: the [General] section is missing – the server probably cannot read this file.',
    restartWait: 'Restart initiated – waiting for the server to come back …',
    restartWaitSub: 'The page reloads by itself as soon as it answers.',
    restartDoneMsg: 'Restart initiated',
    cmd_running: 'The command is already running',
    cmd_disabled: 'No restart command configured (start with -restartcmd)',
    sigFailed: 'PID %s is still running – SIGTERM was ignored or the process is stuck',
    signal_disabled: 'SIGTERM feature is off (start with -restartsignal)',
    restartWaiting: 'Restart requested – waiting for "%s" to disappear …',
    restartIdle: 'Restart file: %s',
    restartLast: 'Last restarted: %s',
    restart_requested: 'Restart requested – file "%s" created',
    restart_already_pending: 'A restart request is already pending',
    restart_cancelled: 'Restart request withdrawn',
    restart_done: 'Server was restarted',
    restart_disabled: 'Restart feature is off (start with -restartfile)',
    restart_failed: 'Could not create the restart file: %s',
    restart_cancel_failed: 'Could not remove the restart file: %s',
    map_set: 'Map set – restart the server for it to take effect',
    no_free_name: 'No free name found in custom_maps for "%s"',
    map_parked_renamed: 'Note: "%s" already existed, parked as "%s" in custom_maps – nothing deleted.',
    backup_failed: 'Note: could not create a backup of ServerConfig.toml: %s',
    config_write_failed: 'ServerConfig.toml could not be written: %s (files were moved back)',
    mod_on: 'Mod activated: %s', mod_off: 'Mod deactivated: %s',
    plugin_on: 'Plugin activated: %s', plugin_off: 'Plugin deactivated: %s',
    dir_changed: 'Directory changed',
    config_missing: 'ServerConfig.toml not found in %s',
    config_unreadable: 'ServerConfig.toml could not be read: %s',
    zip_unreadable: 'Zip could not be read: %s',
    duplicate: '"%s" already exists in the target folder – delete one of the two copies',
    move_failed: '"%s" could not be moved: %s',
    remove_failed: '%s could not be removed: %s',
    map_park_failed: '%s could not be moved back to custom_maps: %s',
    map_activate_failed: '%s could not be activated: %s',
    invalid_name: 'Invalid name', no_map: 'No map given',
    dir_not_found: 'Directory not found: %s',
    dir_locked: 'The directory is fixed and cannot be changed here',
    dirFixed: 'Directory',
    bad_request: 'Bad request',
    log_missing: 'No log file found in %s (expected: Server.log)',
    log_unreadable: '%s could not be read: %s',
    generic: '%s'
  }
};

let state = null, tab = 'maps', busy = false, logTimer = null;
let lang = localStorage.getItem('lang');
if (!T[lang]) lang = (navigator.language || '').toLowerCase().startsWith('de') ? 'de' : 'en';
let theme = localStorage.getItem('theme');
if (theme !== 'light' && theme !== 'dark')
  theme = matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';

const esc = s => String(s ?? '').replace(/[&<>"']/g, c =>
  ({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'}[c]));

function t(m, ...rest) {
  const code = typeof m === 'string' ? m : (m && m.code);
  const args = typeof m === 'string' ? rest : ((m && m.args) || []);
  let s = T[lang][code];
  if (s === undefined) return code + (args.length ? ': ' + args.join(' ') : '');
  for (const a of args) s = s.replace('%s', a);
  return s;
}

function applyTheme() {
  document.documentElement.dataset.theme = theme;
  document.getElementById('themeBtn').textContent =
    t(theme === 'dark' ? 'theme_light' : 'theme_dark');
}
function toggleTheme() {
  theme = theme === 'dark' ? 'light' : 'dark';
  localStorage.setItem('theme', theme);
  applyTheme();
}
function toggleLang() {
  lang = lang === 'de' ? 'en' : 'de';
  localStorage.setItem('lang', lang);
  document.documentElement.lang = lang;
  applyStrings();
  render();
}

function applyStrings() {
  document.getElementById('langBtn').textContent = t('langName');
  document.getElementById('dir').placeholder = t('dirPlaceholder');
  document.getElementById('dirBtn').textContent = t('setDir');
  document.getElementById('reloadBtn').textContent = t('reload');
  document.getElementById('tabs').innerHTML = ['maps','mods','plugins','console','config','log'].map(k =>
    '<button id="tab-' + k + '" class="' + (k === tab ? 'sel' : '') +
    '" onclick="show(&quot;' + k + '&quot;)">' + esc(t(k)) + '</button>').join('');
  applyTheme();
}

async function post(url, body) {
  if (busy) return;
  busy = true;
  try {
    const r = await fetch(url, {method:'POST', headers:{'Content-Type':'application/json'},
                               body: JSON.stringify(body)});
    const d = await r.json();
    if (d.state) { state = d.state; document.getElementById('dir').value = state.dir; render(); }
    const notes = (d.notes || []).map(n => t(n));
    msg([d.error ? t(d.error) : t(d.ok)].concat(notes).join('\n'), !!d.error, notes.length > 0);
  } catch (e) { msg(String(e), true); }
  finally { busy = false; }
}

async function load() {
  try {
    state = await (await fetch('/api/state')).json();
    document.getElementById('dir').value = state.dir;
    render();
    msg(t('reloaded'), false);
    if (tab === 'log') loadLog();
    if (tab === 'config') loadConfig();
    if (tab === 'console') startConsole();
  } catch (e) { msg(String(e), true); }
}

let msgTimer = null;
function msg(text, isErr, sticky) {
  const el = document.getElementById('msg');
  clearTimeout(msgTimer);
  if (!text) { el.className = ''; return; }
  el.textContent = text;
  el.className = 'show ' + (isErr ? 'err' : 'ok');
  if (!isErr && !sticky) msgTimer = setTimeout(() => { el.className = ''; }, 4000);
}

const setDir    = () => post('/api/dir', {dir: document.getElementById('dir').value.trim()});
const setMap    = (path, zip) => post('/api/map', {path, zip});
const toggleMod = (name, activate) => post('/api/mod', {name, activate});
const togglePlg = (name, activate) => post('/api/plugin', {name, activate});

function show(k) {
  tab = k;
  if (location.hash.slice(1) !== k) history.replaceState(null, '', '#' + k);
  for (const x of ['maps','mods','plugins','console','config','log']) {
    const b = document.getElementById('tab-' + x);
    if (b) b.className = x === k ? 'sel' : '';
  }
  stopLogTimer();
  stopConsole();
  render();
  if (k === 'log') loadLog();
  if (k === 'config') loadConfig();
  if (k === 'console') startConsole();
}

let restartTimer = null, restartWasPending = false;

const doRestart       = () => post('/api/restart', {});
const doRestartCancel = () => post('/api/restart', {cancel: true});
async function doCmd() {
  await post('/api/restart', {cmd: true});
  watchRestart();
}

function watchRestart() {
  const wasInstance = state && state.instance;
  showReconnect(t('restartWait'), t('restartWaitSub'));
  let lost = false, settled = 0;
  const tick = setInterval(async () => {
    try {
      const r = await fetch('/api/state', {cache: 'no-store'});
      if (!r.ok) throw new Error('down');
      const d = await r.json();
      if (d.instance !== wasInstance) { clearInterval(tick); location.reload(); return; }
      if (lost) { clearInterval(tick); location.reload(); return; }
      if (!d.cmd.running && ++settled > 4) {
        clearInterval(tick);
        document.getElementById('reconnect').hidden = true;
        state = d;
        render();
        msg(d.cmd.failed ? t('cmdFailedInfo', new Date(d.cmd.at * 1000).toLocaleTimeString(),
                             String(d.cmd.code), d.cmd.out || '')
                         : t('restartDoneMsg'), !!d.cmd.failed);
      }
    } catch (e) { lost = true; }
  }, 1000);
}

function renderCmd() {
  const c = (state && state.cmd) || {};
  const btn = document.getElementById('cmdBtn');
  const info = document.getElementById('cmdInfo');
  btn.hidden = info.hidden = !c.enabled;
  if (!c.enabled) return;
  btn.textContent = t('cmdBtn');
  btn.disabled = !!c.running;
  info.hidden = !c.failed;
  if (c.failed)
    info.textContent = t('cmdFailedInfo', new Date(c.at * 1000).toLocaleTimeString(),
                         String(c.code), c.out || '');
}

async function doSignal() {
  await post('/api/restart', {signal: true});
  watchAfterSignal();
}

function watchAfterSignal() {
  const wasInstance = state && state.instance;
  let lost = false;
  const tick = setInterval(async () => {
    try {
      const r = await fetch('/api/state', {cache: 'no-store'});
      if (!r.ok) throw new Error('not ok');
      const d = await r.json();
      if (lost && d.instance !== wasInstance) { clearInterval(tick); location.reload(); return; }
      state = d;
      renderRestart(); renderSignal(); renderCmd();
      if (!d.restart.sigPending && !d.cmd.running) { clearInterval(tick); render(); }
    } catch (e) {
      if (!lost) { lost = true; showReconnect(t('reconnectGone'), t('reconnectSub')); }
    }
  }, 1000);
}


function showReconnect(title, sub) {
  document.getElementById('reconnectText').textContent = title;
  document.getElementById('reconnectSub').textContent = sub;
  document.getElementById('reconnect').hidden = false;
  stopLogTimer();
  stopRestartTimer();
}


function renderRestart() {
  const r = (state && state.restart) || {enabled: false};
  const btn = document.getElementById('restartBtn');
  const cancel = document.getElementById('restartCancelBtn');
  const info = document.getElementById('restartInfo');
  btn.hidden = cancel.hidden = info.hidden = !r.enabled;
  if (!r.enabled) { stopRestartTimer(); return; }

  btn.textContent = t('restart');
  btn.disabled = r.pending;
  cancel.textContent = t('restartCancel');
  cancel.hidden = !r.pending;

  const bits = [r.pending ? t('restartWaiting', shortName(r.file)) : t('restartIdle', r.file)];
  if (r.doneAt) bits.push(t('restartLast', new Date(r.doneAt * 1000).toLocaleString()));
  info.textContent = bits.join(' · ');

  if (r.pending) { restartWasPending = true; startRestartTimer(); }
  else if (!r.sigPending) {
    stopRestartTimer();
    if (restartWasPending) { restartWasPending = false; msg(t('restart_done'), false); }
  }
}

let sigWasPending = false;


function renderSignal() {
  const r = (state && state.restart) || {};
  const btn = document.getElementById('sigBtn');
  const info = document.getElementById('sigInfo');
  btn.hidden = info.hidden = !r.signal;
  if (!r.signal) return;

  const procs = r.procs || [];
  btn.textContent = t('sigBtn');
  btn.disabled = r.sigPending || procs.length === 0;

  const running = procs.length
    ? t('sigFound', procs.map(p => p.pid).join(', '))
    : t('sigNone');
  const dying = procs.map(p => p.pid).join(', ');
  const bits = [r.sigFailed ? t('sigFailed', dying)
              : r.sigStopped ? t('sigWaitingNew')
              : r.sigPending ? t('sigWaiting', dying)
              : t('sigIdle', r.procName, running)];
  if (r.sigDoneAt) {
    const when = new Date(r.sigDoneAt * 1000).toLocaleString();
    bits.push(r.sigNewPid ? t('sigLastNew', when, String(r.sigNewPid))
                          : t('sigLastStopped', when));
  }
  info.textContent = bits.filter(Boolean).join(' · ');

  if (r.sigPending) { sigWasPending = true; startRestartTimer(); }
  else if (sigWasPending) {
    sigWasPending = false;
    stopRestartTimer();
    msg(r.sigNewPid ? t('sigLastNew', new Date((r.sigDoneAt || 0) * 1000).toLocaleString(),
                        String(r.sigNewPid))
                    : t('sigLastStopped', new Date((r.sigDoneAt || 0) * 1000).toLocaleString()), false);
  }
}

const shortName = p => String(p || '').split(/[\\/]/).pop();

function startRestartTimer() {
  if (restartTimer) return;
  restartTimer = setInterval(refreshRestart, 2000);
}
function stopRestartTimer() {
  if (restartTimer) { clearInterval(restartTimer); restartTimer = null; }
}

async function refreshRestart() {
  try {
    const d = await (await fetch('/api/state')).json();
    state = d;
    renderRestart();
    renderSignal();
   
    renderCmd();
  } catch (e) { }
}

function hue(s) {
  let h = 0;
  for (const c of String(s)) h = (h * 31 + c.charCodeAt(0)) % 360;
  return h;
}

function ph(title, small) {
  const h = hue(title);
  const bg = 'linear-gradient(140deg, hsl(' + h + ' 32% 34%), hsl(' + ((h + 40) % 360) + ' 30% 22%))';
  const text = small ? esc(String(title).trim().slice(0, 2).toUpperCase()) : esc(title);
  return '<div class="ph" style="background:' + bg + '">' + text + '</div>';
}

function pic(src, title, small) {
  const p = ph(title, small);
  if (!src) return '<div class="pic">' + p + '</div>';
  return '<div class="pic">' + p +
    '<img src="' + esc(src) + '" alt="" loading="lazy" decoding="async" onerror="this.remove()">' +
    '</div>';
}

const srcLabel = s => s === 'custom_maps' ? t('srcCustom') : t('srcClient');

function render() {
  const m = document.getElementById('meta');
  if (!state) { m.textContent = ''; return; }
  const locked = !!state.dirLocked;
  document.getElementById('dir').hidden = locked;
  document.getElementById('dirBtn').hidden = locked;
  m.innerHTML = (locked ? esc(t('dirFixed')) + ': <b>' + esc(state.dir) + '</b><br>' : '') +
    (state.configFound
    ? esc(t('server')) + ': <b>' + esc(state.serverName || '—') + '</b> &nbsp;·&nbsp; ' +
      esc(t('port')) + ' <b>' + esc(state.port) + '</b> &nbsp;·&nbsp; ' +
      esc(t('resources')) + ': <b>' + esc(state.resourceFolder) + '</b><br>' +
      esc(t('currentMap')) + ': <b>' + esc(state.currentMap || '—') + '</b>'
    : '<b>' + esc(t('noConfig')) + '</b>');
  if (state.warnings && state.warnings.length)
    m.innerHTML += '<br>' + state.warnings.map(w => '⚠ ' + esc(t(w))).join('<br>');

  renderRestart();
  renderSignal();
 
  renderCmd();
  document.getElementById('view').innerHTML =
    tab === 'maps' ? renderMaps() :
    tab === 'mods' ? renderMods() :
    tab === 'plugins' ? renderPlugins() :
    tab === 'console' ? renderConsole() :
    tab === 'config' ? renderConfig() : renderLog();
}

function table(head, rows, emptyText) {
  if (!rows) return '<div class="empty">' + esc(emptyText) + '</div>';
  return '<table><thead><tr>' + head.map(h => '<th>' + esc(h) + '</th>').join('') +
         '</tr></thead><tbody>' + rows + '</tbody></table>';
}

function renderMaps() {
  const maps = (state && state.maps) || [];
  if (!maps.length) return '<div class="empty">' + esc(t('noMaps')) + '</div>';
  return '<div class="grid">' + maps.map(mp => {
    const args = JSON.stringify(mp.path).replace(/"/g, '&quot;') + ',' +
                 JSON.stringify(mp.zip).replace(/"/g, '&quot;');
    return '<div class="card' + (mp.active ? ' on' : '') + '">' +
      '<div class="thumb">' +
        (mp.active ? '<span class="badge">' + esc(t('running')) + '</span>' : '') +
        pic(mp.img, mp.title, false) + '</div>' +
      '<div class="body">' +
        '<div class="name">' + esc(mp.title) + '</div>' +
        '<div class="sub">' + esc(srcLabel(mp.source)) + '</div>' +
        (mp.active ? '' : '<button class="primary" onclick="setMap(' + args + ')">' +
          esc(t('set')) + '</button>') +
      '</div></div>';
  }).join('') + '</div>';
}

function itemRows(items, fn) {
  return (items || []).map(it =>
    '<tr>' +
      '<td class="ico">' + pic(it.img, it.title, true) + '</td>' +
      '<td>' + esc(it.title) + '<div class="sub">' + esc(it.name) + '</div></td>' +
      '<td><span class="tag ' + (it.active ? 'on">' + esc(t('active'))
                                           : 'off">' + esc(t('inactive'))) + '</span></td>' +
      '<td class="act"><button class="' + (it.active ? 'off' : 'primary') + '" onclick="' +
        fn + '(' + JSON.stringify(it.name).replace(/"/g, '&quot;') + ',' + !it.active + ')">' +
        esc(it.active ? t('deactivate') : t('activate')) + '</button></td>' +
    '</tr>').join('');
}

const renderMods = () => table(['', t('colMod'), t('colStatus'), ''],
  itemRows(state.mods, 'toggleMod'), t('noMods'));
const renderPlugins = () => table(['', t('colPlugin'), t('colStatus'), ''],
  itemRows(state.plugins, 'togglePlg'), t('noPlugins'));

let conTimer = null;

function renderConsole() {
  const c = (state && state.console) || {};
  if (!c.enabled) return '<div class="empty">' + esc(t('conOff')) + '</div>';
  return '<div class="logbar"><span class="sub">tmux: ' + esc(c.session) + '</span></div>' +
    '<pre id="context"></pre>' +
    '<div class="conrow">' +
      '<input type="text" id="conin" spellcheck="false" placeholder="' + esc(t('conPlaceholder')) + '">' +
      '<button class="primary" onclick="conSend()">' + esc(t('conSend')) + '</button>' +
    '</div>';
}

function startConsole() {
  const inp = document.getElementById('conin');
  if (inp) {
    inp.addEventListener('keydown', e => { if (e.key === 'Enter') conSend(); });
    inp.focus();
  }
  pollConsole();
  if (!conTimer) conTimer = setInterval(pollConsole, 1000);
}

function stopConsole() {
  if (conTimer) { clearInterval(conTimer); conTimer = null; }
}

async function pollConsole() {
  const el = document.getElementById('context');
  if (!el || tab !== 'console') { stopConsole(); return; }
  try {
    const d = await (await fetch('/api/console', {cache: 'no-store'})).json();
    if (d.error) { el.textContent = t(d.error); stopConsole(); return; }
    if (d.text === el.textContent) return;
    const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 40;
    el.textContent = d.text;
    if (atBottom) el.scrollTop = el.scrollHeight;
  } catch (e) { }
}

async function conSend() {
  const inp = document.getElementById('conin');
  if (!inp) return;
  const line = inp.value;
  inp.value = '';
  try {
    const r = await fetch('/api/console', {method: 'POST',
      headers: {'Content-Type': 'application/json'}, body: JSON.stringify({line: line})});
    const d = await r.json();
    if (d.error) msg(t(d.error), true);
  } catch (e) { msg(String(e), true); }
  setTimeout(pollConsole, 150);
}

function renderConfig() {
  return '<div class="logbar">' +
      '<button class="small primary" onclick="saveConfig()">' + esc(t('cfgSave')) + '</button>' +
      '<button class="small" onclick="loadConfig()">' + esc(t('reload')) + '</button>' +
      '<span id="cfginfo"></span>' +
    '</div><textarea id="cfgbox" spellcheck="false"></textarea>';
}

async function loadConfig() {
  const box = document.getElementById('cfgbox');
  if (!box) return;
  try {
    const d = await (await fetch('/api/config', {cache: 'no-store'})).json();
    if (d.error) { msg(t(d.error), true); return; }
    box.value = d.text;
    const info = document.getElementById('cfginfo');
    if (info) info.textContent = d.file;
  } catch (e) { msg(String(e), true); }
}

async function saveConfig() {
  const box = document.getElementById('cfgbox');
  if (!box) return;
  await post('/api/config', {text: box.value});
}

function renderLog() {
  const auto = localStorage.getItem('logauto') !== '0';
  return '<div class="logbar">' +
      '<label><input type="checkbox" id="logauto"' + (auto ? ' checked' : '') +
        ' onchange="logAutoChanged()">' + esc(t('logAuto')) + '</label>' +
      '<button class="small" onclick="loadLog()">' + esc(t('reload')) + '</button>' +
      '<button class="small" onclick="logBottom()">' + esc(t('logBottom')) + '</button>' +
      '<span id="loginfo"></span>' +
    '</div><pre id="logtext"></pre>';
}

function fmtBytes(n) {
  if (n < 1024) return n + ' B';
  if (n < 1024 * 1024) return (n / 1024).toFixed(1) + ' KB';
  return (n / 1024 / 1024).toFixed(1) + ' MB';
}

function logBottom() {
  const el = document.getElementById('logtext');
  if (el) el.scrollTop = el.scrollHeight;
}

function logAutoChanged() {
  const on = document.getElementById('logauto').checked;
  localStorage.setItem('logauto', on ? '1' : '0');
  stopLogTimer();
  if (on) startLogTimer();
}

function startLogTimer() {
  if (logTimer) return;
  logTimer = setInterval(() => { if (tab === 'log') loadLog(true); }, 3000);
}
function stopLogTimer() {
  if (logTimer) { clearInterval(logTimer); logTimer = null; }
}

async function loadLog(quiet) {
  const el = document.getElementById('logtext');
  if (!el) return;
  try {
    const d = await (await fetch('/api/log')).json();
    const info = document.getElementById('loginfo');
    if (d.error) {
      el.textContent = '';
      if (info) info.textContent = '';
      if (!quiet) msg(t(d.error), true);
      stopLogTimer();
      return;
    }
    const atBottom = el.scrollHeight - el.scrollTop - el.clientHeight < 40;
    el.textContent = d.text || t('logEmpty');
    if (info) info.textContent = t('logFile') + ': ' + d.file + ' · ' +
      t('logSize') + ': ' + fmtBytes(d.size || 0) +
      (d.truncated ? ' · ' + t('logTail') + ' ' + fmtBytes((d.text || '').length) : '');
    if (atBottom) logBottom();
    const box = document.getElementById('logauto');
    if (box && box.checked) startLogTimer();
  } catch (e) { if (!quiet) msg(String(e), true); }
}

document.getElementById('dir').addEventListener('keydown', e => { if (e.key === 'Enter') setDir(); });
addEventListener('hashchange', () => {
  const k = location.hash.slice(1);
  if (k && k !== tab) show(k);
});
const startTab = location.hash.slice(1);
if (['maps','mods','plugins','console','config','log'].includes(startTab)) tab = startTab;
document.documentElement.lang = lang;
applyStrings();
load();
</script>
</body>
</html>
`
