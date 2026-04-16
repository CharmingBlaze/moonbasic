package mbsystem

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	rt "runtime"
	"runtime/debug"
	"strings"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) registerHost(reg runtime.Registrar) {
	reg.Register("SystemProperty", "system", m.sysSystemProperty)
	reg.Register("SYSTEM.CPUNAME", "system", m.sysCpuName)
	reg.Register("SYSTEM.GPUNAME", "system", m.sysGpuName)
	reg.Register("GpuName", "system", m.sysGpuName)
	reg.Register("SYSTEM.TOTALMEMORY", "system", m.sysTotalMemory)
	reg.Register("SYSTEM.FREEMEMORY", "system", m.sysFreeMemory)
	reg.Register("SYSTEM.GETENV", "system", m.sysGetenv)
	reg.Register("ENVIRON", "system", m.sysGetenv)
	reg.Register("ENVIRON$", "system", m.sysGetenv)
	reg.Register("SYSTEM.SETENV", "system", m.sysSetenv)
	reg.Register("SYSTEM.EXECUTE", "system", m.sysExecute)
	reg.Register("SYSTEM.OPENURL", "system", m.sysOpenURL)
	reg.Register("SYSTEM.LOCALE", "system", m.sysLocale)
	reg.Register("SYSTEM.USERNAME", "system", m.sysUsername)
	reg.Register("SYSTEM.ISDEBUGBUILD", "system", m.sysIsDebugBuild)
	reg.Register("SYSTEM.VERSION", "system", m.sysVersion)
	m.registerClipboard(reg)
}

// sysSystemProperty returns a small set of OS/runtime facts keyed by name (MoonBasic 3D parity).
func (m *Module) sysSystemProperty(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("SystemProperty expects (key)")
	}
	key, err := rt2.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "os", "os_name":
		return rt2.RetString(rt.GOOS), nil
	case "arch":
		return rt2.RetString(rt.GOARCH), nil
	case "os_version":
		return rt2.RetString(rt.GOOS + " " + rt.GOARCH), nil
	case "cpu_cores":
		return value.FromInt(int64(rt.NumCPU())), nil
	case "compiler":
		return rt2.RetString(rt.Version()), nil
	default:
		return rt2.RetString(""), nil
	}
}

func (m *Module) sysVersion(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = m
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("SYSTEM.VERSION expects 0 arguments")
	}
	return rt2.RetString("1.0.0-GOLD"), nil
}

func (m *Module) sysCpuName(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("SYSTEM.CPUNAME expects 0 arguments")
	}
	infos, err := cpu.Info()
	if err == nil && len(infos) > 0 && infos[0].ModelName != "" {
		return rt2.RetString(strings.TrimSpace(infos[0].ModelName)), nil
	}
	if s := os.Getenv("PROCESSOR_IDENTIFIER"); s != "" && rt.GOOS == "windows" {
		return rt2.RetString(s), nil
	}
	if rt.GOOS == "darwin" {
		out, e := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output()
		if e == nil {
			return rt2.RetString(strings.TrimSpace(string(out))), nil
		}
	}
	if rt.GOOS == "linux" {
		if b, e := os.ReadFile("/proc/cpuinfo"); e == nil {
			for _, line := range strings.Split(string(b), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "model name") {
					if i := strings.IndexByte(line, ':'); i >= 0 {
						return rt2.RetString(strings.TrimSpace(line[i+1:])), nil
					}
				}
			}
		}
	}
	return rt2.RetString(rt.GOARCH), nil
}

func (m *Module) sysGpuName(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("SYSTEM.GPUNAME expects 0 arguments")
	}
	switch rt.GOOS {
	case "windows":
		out, err := exec.Command("wmic", "path", "win32_VideoController", "get", "Name").Output()
		if err == nil {
			if name := parseWMICName(string(out)); name != "" {
				return rt2.RetString(name), nil
			}
		}
	case "darwin":
		out, err := exec.Command("system_profiler", "SPDisplaysDataType").Output()
		if err == nil {
			for _, line := range strings.Split(string(out), "\n") {
				line = strings.TrimSpace(line)
				const prefix = "Chipset Model:"
				if idx := strings.Index(line, prefix); idx >= 0 {
					s := strings.TrimSpace(line[idx+len(prefix):])
					if s != "" {
						return rt2.RetString(s), nil
					}
				}
			}
		}
	}
	return rt2.RetString("(unavailable)"), nil
}

func parseWMICName(out string) string {
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.EqualFold(line, "Name") {
			continue
		}
		return line
	}
	return ""
}

func (m *Module) sysTotalMemory(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt2
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("SYSTEM.TOTALMEMORY expects 0 arguments")
	}
	v, err := mem.VirtualMemory()
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(v.Total)), nil
}

func (m *Module) sysFreeMemory(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt2
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("SYSTEM.FREEMEMORY expects 0 arguments")
	}
	v, err := mem.VirtualMemory()
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(v.Available)), nil
}

func (m *Module) sysGetenv(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("SYSTEM.GETENV / ENVIRON expects (key)")
	}
	key, err := rt2.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	return rt2.RetString(os.Getenv(key)), nil
}

func (m *Module) sysSetenv(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("SYSTEM.SETENV expects (key, val)")
	}
	key, err := rt2.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	val, err := rt2.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	err = os.Setenv(key, val)
	return value.Nil, err
}

func (m *Module) sysExecute(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("SYSTEM.EXECUTE expects (cmd)")
	}
	cmdline, err := rt2.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	cmd := shellExec(cmdline)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err == nil {
		return value.FromInt(0), nil
	}
	var ee *exec.ExitError
	if errors.As(err, &ee) {
		return value.FromInt(int64(ee.ExitCode())), nil
	}
	return value.Nil, err
}

func shellExec(cmdline string) *exec.Cmd {
	if rt.GOOS == "windows" {
		return exec.Command("cmd", "/C", cmdline)
	}
	return exec.Command("sh", "-c", cmdline)
}

func (m *Module) sysOpenURL(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("SYSTEM.OPENURL expects (url)")
	}
	url, err := rt2.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	url = strings.TrimSpace(url)
	if url == "" {
		return value.Nil, fmt.Errorf("SYSTEM.OPENURL: empty url")
	}
	var cmd *exec.Cmd
	switch rt.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return value.Nil, err
	}
	go func() { _ = cmd.Wait() }()
	return value.Nil, nil
}

func (m *Module) sysLocale(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("SYSTEM.LOCALE expects 0 arguments")
	}
	for _, k := range []string{"LC_ALL", "LC_MESSAGES", "LANG", "LANGUAGE"} {
		if s := os.Getenv(k); s != "" {
			return rt2.RetString(s), nil
		}
	}
	return rt2.RetString("C"), nil
}

func (m *Module) sysUsername(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("SYSTEM.USERNAME expects 0 arguments")
	}
	if u, err := user.Current(); err == nil && u.Username != "" {
		return rt2.RetString(u.Username), nil
	}
	if s := os.Getenv("USER"); s != "" {
		return rt2.RetString(s), nil
	}
	if s := os.Getenv("USERNAME"); s != "" {
		return rt2.RetString(s), nil
	}
	return rt2.RetString(""), nil
}

func (m *Module) sysIsDebugBuild(rt2 *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt2
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("SYSTEM.ISDEBUGBUILD expects 0 arguments")
	}
	if v := os.Getenv("MOONBASIC_DEBUG"); v == "1" || strings.EqualFold(v, "true") {
		return value.FromBool(true), nil
	}
	info, ok := debug.ReadBuildInfo()
	if ok && info.Main.Version == "(devel)" {
		return value.FromBool(true), nil
	}
	return value.FromBool(false), nil
}
