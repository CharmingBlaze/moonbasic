package mbdebug

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func strPool() []string {
	if r := runtime.ActiveRegistry(); r != nil && r.Prog != nil {
		return r.Prog.StringTable
	}
	return nil
}

func formatMBValueRT(rt *runtime.Runtime, v value.Value) string {
	if v.Kind == value.KindString {
		s, err := rt.ArgString([]value.Value{v}, 0)
		if err == nil {
			return s
		}
	}
	return v.String()
}

func (m *Module) Register(r runtime.Registrar) {
	r.Register("DEBUG.ENABLE", "debug", runtime.AdaptLegacy(m.debugEnable))
	r.Register("DEBUG.DISABLE", "debug", runtime.AdaptLegacy(m.debugDisable))
	r.Register("DEBUG.ISENABLED", "debug", runtime.AdaptLegacy(m.debugIsEnabled))
	r.Register("DEBUG.PRINT", "debug", m.debugPrint1)
	r.Register("DEBUG.PRINTL", "debug", m.debugPrintLabeled)
	r.Register("DEBUG.WATCH", "debug", m.debugWatch)
	r.Register("DEBUG.WATCHCLEAR", "debug", runtime.AdaptLegacy(m.debugWatchClear))
	r.Register("DEBUG.ASSERT", "debug", m.debugAssert)
	r.Register("ASSERT", "debug", m.debugAssert)
	r.Register("DEBUG.BREAKPOINT", "debug", runtime.AdaptLegacy(m.debugBreakpoint))
	r.Register("DEBUG.LOG", "debug", m.debugLog)
	r.Register("DEBUG.LOGFILE", "debug", m.debugLogFile)
	r.Register("DEBUG.PROFILESTART", "debug", m.profileStart)
	r.Register("DEBUG.PROFILEEND", "debug", m.profileEnd)
	r.Register("DEBUG.PROFILEREPORT", "debug", runtime.AdaptLegacy(m.profileReport))
	r.Register("DEBUG.STACKTRACE", "debug", runtime.AdaptLegacy(m.debugStackTrace))
	r.Register("DEBUG.HEAPSTATS", "debug", runtime.AdaptLegacy(m.debugHeapStats))
	r.Register("DEBUG.GCSTATS", "debug", runtime.AdaptLegacy(m.debugGCStats))
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.overlayUser = false
	m.watches = nil
	m.profStack = nil
	m.profSum = nil
	m.profN = nil
}

func (m *Module) debugEnable(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("DEBUG.ENABLE expects 0 arguments")
	}
	m.mu.Lock()
	m.overlayUser = true
	m.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) debugDisable(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("DEBUG.DISABLE expects 0 arguments")
	}
	m.mu.Lock()
	m.overlayUser = false
	m.mu.Unlock()
	return value.Nil, nil
}

func (m *Module) debugIsEnabled(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("DEBUG.ISENABLED expects 0 arguments")
	}
	reg := runtime.ActiveRegistry()
	m.mu.Lock()
	u := m.overlayUser
	m.mu.Unlock()
	on := u || (reg != nil && reg.DebugMode)
	return value.FromBool(on), nil
}

func (m *Module) debugPrint1(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, runtime.Errorf("DEBUG.PRINT expects 1 argument (use DEBUG.PRINTL for label + value)")
	}
	fmt.Fprintln(runtime.DiagWriter(), formatMBValueRT(rt, args[0]))
	return value.Nil, nil
}

func (m *Module) debugPrintLabeled(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("DEBUG.PRINTL expects (label$, value)")
	}
	label, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	fmt.Fprintf(runtime.DiagWriter(), "%s: %s\n", label, formatMBValueRT(rt, args[1]))
	return value.Nil, nil
}

func (m *Module) debugWatch(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("DEBUG.WATCH expects (label$, value)")
	}
	label, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	text := formatMBValueRT(rt, args[1])
	m.mu.Lock()
	defer m.mu.Unlock()
	for i := range m.watches {
		if m.watches[i].label == label {
			m.watches[i].text = text
			return value.Nil, nil
		}
	}
	m.watches = append(m.watches, watchEntry{label: label, text: text})
	return value.Nil, nil
}

func (m *Module) debugWatchClear(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("DEBUG.WATCHCLEAR expects 0 arguments")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.watches = nil
	return value.Nil, nil
}

func (m *Module) debugAssert(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("ASSERT expects (cond?, msg$)")
	}
	pool := strPool()
	if !value.Truthy(args[0], pool, rt.Heap) {
		msg, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
		return value.Nil, runtime.Errorf("ASSERT: %s", msg)
	}
	return value.Nil, nil
}

func (m *Module) debugBreakpoint(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("DEBUG.BREAKPOINT expects 0 arguments")
	}
	reg := runtime.ActiveRegistry()
	var stack string
	if reg != nil && reg.StackTraceFn != nil {
		stack = reg.StackTraceFn()
	} else {
		stack = "(stack trace unavailable)\n"
	}
	return value.Nil, fmt.Errorf("DEBUG.BREAKPOINT\n%s", stack)
}

func (m *Module) debugLog(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("DEBUG.LOG expects (msg$)")
	}
	msg, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	fmt.Fprintln(runtime.DiagWriter(), msg)
	return value.Nil, nil
}

func (m *Module) debugLogFile(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
		return value.Nil, runtime.Errorf("DEBUG.LOGFILE expects (path$, msg$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	msg, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return value.Nil, err
	}
	defer f.Close()
	_, err = io.WriteString(f, msg+"\n")
	return value.Nil, err
}

func (m *Module) profileStart(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("DEBUG.PROFILESTART expects (label$)")
	}
	label, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.profSum == nil {
		m.profSum = make(map[string]time.Duration)
		m.profN = make(map[string]int64)
	}
	m.profStack = append(m.profStack, profileFrame{label: label, start: time.Now()})
	return value.Nil, nil
}

func (m *Module) profileEnd(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, runtime.Errorf("DEBUG.PROFILEEND expects (label$)")
	}
	want, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.profStack) == 0 {
		return value.Nil, runtime.Errorf("DEBUG.PROFILEEND: no matching PROFILESTART")
	}
	top := m.profStack[len(m.profStack)-1]
	if top.label != want {
		return value.Nil, runtime.Errorf("DEBUG.PROFILEEND: expected label %q, stack has %q", want, top.label)
	}
	m.profStack = m.profStack[:len(m.profStack)-1]
	d := time.Since(top.start)
	m.profSum[want] += d
	m.profN[want]++
	return value.Nil, nil
}

func (m *Module) profileReport(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("DEBUG.PROFILEREPORT expects 0 arguments")
	}
	m.mu.Lock()
	labels := make([]string, 0, len(m.profSum))
	for k := range m.profSum {
		labels = append(labels, k)
	}
	sort.Strings(labels)
	var b strings.Builder
	for _, k := range labels {
		fmt.Fprintf(&b, "  %s: total %s (%d samples)\n", k, m.profSum[k].Round(time.Microsecond), m.profN[k])
	}
	m.mu.Unlock()
	if b.Len() == 0 {
		fmt.Fprintln(runtime.DiagWriter(), "[DEBUG.PROFILEREPORT] (no samples)")
	} else {
		fmt.Fprintf(runtime.DiagWriter(), "[DEBUG.PROFILEREPORT]\n%s", b.String())
	}
	return value.Nil, nil
}

func (m *Module) debugStackTrace(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("DEBUG.STACKTRACE expects 0 arguments")
	}
	reg := runtime.ActiveRegistry()
	var stack string
	if reg != nil && reg.StackTraceFn != nil {
		stack = reg.StackTraceFn()
	} else {
		stack = "(stack trace unavailable)\n"
	}
	fmt.Fprintf(runtime.DiagWriter(), "[DEBUG.STACKTRACE]\n%s", stack)
	return value.Nil, nil
}

func (m *Module) debugHeapStats(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("DEBUG.HEAPSTATS expects 0 arguments")
	}
	reg := runtime.ActiveRegistry()
	if reg == nil || reg.Heap == nil {
		fmt.Fprintln(runtime.DiagWriter(), "[DEBUG.HEAPSTATS] (heap unavailable)")
		return value.Nil, nil
	}
	st := reg.Heap.Stats()
	fmt.Fprintf(runtime.DiagWriter(), "[DEBUG.HEAPSTATS] live=%d free_slots=%d peak_slots=%d\n",
		st.LiveCount, st.FreeSlots, st.PeakSlots)
	return value.Nil, nil
}

func (m *Module) debugGCStats(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("DEBUG.GCSTATS expects 0 arguments")
	}
	var s debug.GCStats
	debug.ReadGCStats(&s)
	w := runtime.DiagWriter()
	lastGCS := "(none)"
	if !s.LastGC.IsZero() {
		lastGCS = s.LastGC.Format(time.RFC3339)
	}
	fmt.Fprintf(w, "[DEBUG.GCSTATS] num_gc=%d pause_total=%s last_pause=%s last_gc=%s\n",
		s.NumGC, s.PauseTotal.Round(time.Microsecond),
		lastPause(&s).Round(time.Microsecond), lastGCS)
	if len(s.Pause) > 0 && len(s.Pause) <= 256 {
		fmt.Fprintf(w, "  recent pauses (us): ")
		for i, p := range s.Pause {
			if i > 0 {
				fmt.Fprint(w, ", ")
			}
			fmt.Fprintf(w, "%d", p.Microseconds())
		}
		fmt.Fprintln(w)
	}
	return value.Nil, nil
}

func lastPause(s *debug.GCStats) time.Duration {
	if len(s.Pause) == 0 {
		return 0
	}
	return s.Pause[len(s.Pause)-1]
}
