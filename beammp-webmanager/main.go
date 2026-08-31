package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type app struct {
	mu  sync.Mutex
	mgr Manager

	restartOn      bool
	restartPath    string
	restartPending bool
	restartDoneAt  time.Time

	dockerOn bool

	instance string

	sigOn      bool
	procName   string
	sigOldPIDs []int
	sigPending bool
	sigNewPID  int
	sigDoneAt  time.Time
	sigGoneAt  time.Time
	sigSentAt  time.Time

	tmuxOn      bool
	tmuxSession string
	tmuxWide    bool

	cmdOn      bool
	cmdLine    string
	cmdRunning bool
	cmdAt      time.Time
	cmdCode    int
	cmdOut     string
	sigFailed  bool
}

type optionalString struct {
	set   bool
	value string
}

func (o *optionalString) String() string   { return o.value }
func (o *optionalString) IsBoolFlag() bool { return true }
func (o *optionalString) Set(v string) error {
	o.set = true
	if v != "true" {
		o.value = v
	}
	return nil
}

type prefs struct {
	Dir string `json:"dir"`
}

func prefsPath() string {
	base, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(base, "beammp-tool", "config.json")
}

func loadPrefs() prefs {
	var p prefs
	if path := prefsPath(); path != "" {
		if b, err := os.ReadFile(path); err == nil {
			json.Unmarshal(b, &p)
		}
	}
	return p
}

func savePrefs(p prefs) {
	path := prefsPath()
	if path == "" {
		return
	}
	os.MkdirAll(filepath.Dir(path), 0o755)
	if b, err := json.MarshalIndent(p, "", "  "); err == nil {
		os.WriteFile(path, b, 0o644)
	}
}

func main() {
	dir := flag.String("dir", "", "BeamMP server directory (default: last used, otherwise the current one)")
	addr := flag.String("addr", "127.0.0.1:8477", "address to serve the web interface on")
	noBrowser := flag.Bool("no-browser", false, "do not open a browser")
	inspect := flag.Bool("inspect", false, "inspect the zip files and print the result instead of serving the interface")
	ps := flag.Bool("ps", false, "list visible processes and show which ones are detected as the server")
	var restartFlag optionalString
	flag.Var(&restartFlag, "restartfile",
		"show a restart button that creates a flag file (default: <dir>/restart, or -restartfile=/custom/path)")
	var sigFlag optionalString
	flag.Var(&sigFlag, "restartsignal",
		"show a button that sends SIGTERM to the running server (default: "+defaultProcName+", or -restartsignal=OwnName)")
	var dockerFlag optionalString
	flag.Var(&dockerFlag, "dockerrestart",
		"restart button that runs "+defaultRestartCmd+" (or -dockerrestart=\"own command\"); also hides the directory input")
	var cmdFlag optionalString
	flag.Var(&cmdFlag, "restartcmd",
		"button that runs a shell command (default: "+defaultRestartCmd+", or -restartcmd=\"own command\")")
	var tmuxFlag optionalString
	flag.Var(&tmuxFlag, "tmux",
		"console tab attached to a tmux session (default: "+defaultTmuxSession+", or -tmux=OwnSession)")
	flag.Parse()

	start := *dir
	if start == "" {
		start = loadPrefs().Dir
	}
	if start == "" {
		start, _ = os.Getwd()
	}
	abs, err := filepath.Abs(start)
	if err == nil {
		start = abs
	}

	if *ps {
		name := defaultProcName
		if sigFlag.value != "" {
			name = sigFlag.value
		}
		runPS(name)
		return
	}

	if *inspect {
		mgr := Manager{dir: start}
		runInspect(&mgr)
		return
	}

	a := &app{mgr: Manager{dir: start}, restartOn: restartFlag.set, sigOn: sigFlag.set}
	a.instance = newInstanceID()
	a.tmuxOn = tmuxFlag.set
	a.tmuxSession = defaultTmuxSession
	if tmuxFlag.value != "" {
		a.tmuxSession = tmuxFlag.value
	}
	a.dockerOn = dockerFlag.set
	a.cmdOn = cmdFlag.set || dockerFlag.set
	a.cmdLine = defaultRestartCmd
	if cmdFlag.value != "" {
		a.cmdLine = cmdFlag.value
	} else if dockerFlag.value != "" {
		a.cmdLine = dockerFlag.value
	}

	a.procName = defaultProcName
	if sigFlag.value != "" {
		a.procName = sigFlag.value
	}
	if restartFlag.value != "" {
		if abs, err := filepath.Abs(restartFlag.value); err == nil {
			a.restartPath = abs
		} else {
			a.restartPath = restartFlag.value
		}
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", a.handleIndex)
	mux.HandleFunc("/img", a.handleImg)
	mux.HandleFunc("/api/log", a.handleLog)
	mux.HandleFunc("/api/config", a.handleConfig)
	mux.HandleFunc("/api/console", a.handleConsole)
	mux.HandleFunc("/api/state", a.handleState)
	mux.HandleFunc("/api/dir", a.handleDir)
	mux.HandleFunc("/api/map", a.handleMap)
	mux.HandleFunc("/api/mod", a.handleToggle(false))
	mux.HandleFunc("/api/plugin", a.handleToggle(true))
	mux.HandleFunc("/api/restart", a.handleRestart)

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("cannot listen on %s: %v", *addr, err)
	}
	url := "http://" + ln.Addr().String() + "/"
	if !*noBrowser {
		openBrowser(url)
	}
	log.Fatal(http.Serve(ln, mux))
}

func (a *app) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func (a *app) handleState(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	writeJSON(w, a.fullState())
}

func (a *app) fullState() any {
	return struct {
		State
		Restart   RestartState `json:"restart"`
		Cmd       CmdState     `json:"cmd"`
		Console   ConsoleState `json:"console"`
		Instance  string       `json:"instance"`
		DirLocked bool         `json:"dirLocked"`
	}{a.mgr.State(), a.restartState(), a.cmdState(), a.consoleState(), a.instance, a.dockerOn}
}

func newInstanceID() string {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(b)
}

func (a *app) decode(w http.ResponseWriter, r *http.Request, v any) bool {
	if r.Method != http.MethodPost {
		http.Error(w, "POST expected", http.StatusMethodNotAllowed)
		return false
	}
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		writeJSON(w, map[string]any{"error": msg("bad_request")})
		return false
	}
	return true
}

func (a *app) respond(w http.ResponseWriter, err error, okCode string, okArgs ...string) {
	a.respondNotes(w, err, nil, okCode, okArgs...)
}

func (a *app) respondNotes(w http.ResponseWriter, err error, notes []Msg, okCode string, okArgs ...string) {
	res := map[string]any{"state": a.fullState()}
	if err != nil {
		res["error"] = toMsg(err)
	} else {
		res["ok"] = msg(okCode, okArgs...)
	}
	if len(notes) > 0 {
		res["notes"] = notes
	}
	writeJSON(w, res)
}

func (a *app) handleDir(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Dir string `json:"dir"`
	}
	if !a.decode(w, r, &req) {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.dockerOn {
		a.respond(w, errMsg("dir_locked"), "")
		return
	}
	dir := filepath.Clean(req.Dir)
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}
	st, err := os.Stat(dir)
	if err != nil || !st.IsDir() {
		a.respond(w, errMsg("dir_not_found", dir), "")
		return
	}
	a.mgr.dir = dir
	savePrefs(prefs{Dir: dir})
	a.respond(w, nil, "dir_changed")
}

func (a *app) handleMap(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path string `json:"path"`
		Zip  string `json:"zip"`
	}
	if !a.decode(w, r, &req) {
		return
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if req.Path == "" {
		a.respond(w, errMsg("no_map"), "")
		return
	}
	if req.Zip != "" && req.Zip != filepath.Base(req.Zip) {
		a.respond(w, errMsg("invalid_name"), "")
		return
	}
	notes, err := a.mgr.SetMap(req.Path, req.Zip)
	a.respondNotes(w, err, notes, "map_set")
}

func (a *app) handleToggle(plugin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Name     string `json:"name"`
			Activate bool   `json:"activate"`
		}
		if !a.decode(w, r, &req) {
			return
		}
		a.mu.Lock()
		defer a.mu.Unlock()
		var err error
		code := "mod_off"
		if plugin {
			err = a.mgr.TogglePlugin(req.Name, req.Activate)
			code = "plugin_off"
			if req.Activate {
				code = "plugin_on"
			}
		} else {
			err = a.mgr.ToggleMod(req.Name, req.Activate)
			if req.Activate {
				code = "mod_on"
			}
		}
		a.respond(w, err, code, req.Name)
	}
}

func openBrowser(url string) {
	cmd := exec.Command("xdg-open", url)
	if err := cmd.Start(); err != nil {
		return
	}
	go cmd.Wait()
}
