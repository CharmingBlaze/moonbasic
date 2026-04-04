package input

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type bindKind uint8

const (
	bkKey bindKind = iota
	bkGamepadBtn
	bkGamepadAxis
)

type inputBind struct {
	Kind        bindKind
	Pad         int32
	Code        int32
	KeyAxisSign float64 // keyboard contribution to ActionAxis (0 = none)
}

var (
	actMu  sync.RWMutex
	actMap map[string][]inputBind
)

func normAction(a string) string { return strings.ToLower(strings.TrimSpace(a)) }

func ensureActMap() {
	if actMap == nil {
		actMap = make(map[string][]inputBind)
	}
}

func clearAllBindings() {
	actMap = make(map[string][]inputBind)
}

// defaultKeyAxisSign maps common movement keys to -1 / +1 for ActionAxis (Raylib KeyboardKey values).
func defaultKeyAxisSign(k int32) float64 {
	switch k {
	case 263, 65: // KeyLeft, A
		return -1
	case 262, 68: // KeyRight, D
		return 1
	default:
		return 0
	}
}

func appendBind(action string, b inputBind) {
	ensureActMap()
	a := normAction(action)
	if a == "" {
		return
	}
	actMap[a] = append(actMap[a], b)
}

func actionPressedAny(action string) bool {
	q := actionQueries()
	actMu.RLock()
	binds := actMap[normAction(action)]
	actMu.RUnlock()
	for _, b := range binds {
		switch b.Kind {
		case bkKey:
			if q.keyPressed(b.Code) {
				return true
			}
		case bkGamepadBtn:
			if q.gamepadBtnPressed(b.Pad, b.Code) {
				return true
			}
		}
	}
	return false
}

func actionDownAny(action string) bool {
	q := actionQueries()
	actMu.RLock()
	binds := actMap[normAction(action)]
	actMu.RUnlock()
	for _, b := range binds {
		switch b.Kind {
		case bkKey:
			if q.keyDown(b.Code) {
				return true
			}
		case bkGamepadBtn:
			if q.gamepadBtnDown(b.Pad, b.Code) {
				return true
			}
		}
	}
	return false
}

func actionReleasedAny(action string) bool {
	q := actionQueries()
	actMu.RLock()
	binds := actMap[normAction(action)]
	actMu.RUnlock()
	for _, b := range binds {
		switch b.Kind {
		case bkKey:
			if q.keyReleased(b.Code) {
				return true
			}
		case bkGamepadBtn:
			if q.gamepadBtnReleased(b.Pad, b.Code) {
				return true
			}
		}
	}
	return false
}

func actionAxisSum(action string) float64 {
	q := actionQueries()
	var sum float64
	actMu.RLock()
	binds := actMap[normAction(action)]
	actMu.RUnlock()
	for _, b := range binds {
		switch b.Kind {
		case bkKey:
			if b.KeyAxisSign != 0 && q.keyDown(b.Code) {
				sum += b.KeyAxisSign
			}
		case bkGamepadAxis:
			sum += float64(q.gamepadAxis(b.Pad, b.Code))
		}
	}
	if sum > 1 {
		return 1
	}
	if sum < -1 {
		return -1
	}
	return sum
}

func saveMappingsFile(path string) error {
	actMu.RLock()
	ensureActMap()
	snap := make(map[string][]inputBind, len(actMap))
	for k, v := range actMap {
		cp := make([]inputBind, len(v))
		copy(cp, v)
		snap[k] = cp
	}
	actMu.RUnlock()

	keys := make([]string, 0, len(snap))
	for k := range snap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var sb strings.Builder
	sb.WriteString("# moonBASIC INPUT mappings v1\n")
	for _, act := range keys {
		for _, b := range snap[act] {
			switch b.Kind {
			case bkKey:
				fmt.Fprintf(&sb, "%s key %d %.6g\n", act, b.Code, b.KeyAxisSign)
			case bkGamepadBtn:
				fmt.Fprintf(&sb, "%s gpb %d %d\n", act, b.Pad, b.Code)
			case bkGamepadAxis:
				fmt.Fprintf(&sb, "%s gpa %d %d\n", act, b.Pad, b.Code)
			}
		}
	}
	return os.WriteFile(path, []byte(sb.String()), 0o644)
}

func loadMappingsFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	clearAllBindings()
	sc := bufio.NewScanner(strings.NewReader(string(data)))
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		act := normAction(parts[0])
		if act == "" {
			continue
		}
		switch strings.ToLower(parts[1]) {
		case "key":
			if len(parts) < 4 {
				continue
			}
			kc, err1 := strconv.ParseInt(parts[2], 10, 32)
			sign, err2 := strconv.ParseFloat(parts[3], 64)
			if err1 != nil || err2 != nil {
				continue
			}
			appendBind(act, inputBind{Kind: bkKey, Code: int32(kc), KeyAxisSign: sign})
		case "gpb":
			if len(parts) < 4 {
				continue
			}
			pad, err1 := strconv.ParseInt(parts[2], 10, 32)
			btn, err2 := strconv.ParseInt(parts[3], 10, 32)
			if err1 != nil || err2 != nil {
				continue
			}
			appendBind(act, inputBind{Kind: bkGamepadBtn, Pad: int32(pad), Code: int32(btn)})
		case "gpa":
			if len(parts) < 4 {
				continue
			}
			pad, err1 := strconv.ParseInt(parts[2], 10, 32)
			ax, err2 := strconv.ParseInt(parts[3], 10, 32)
			if err1 != nil || err2 != nil {
				continue
			}
			appendBind(act, inputBind{Kind: bkGamepadAxis, Pad: int32(pad), Code: int32(ax)})
		}
	}
	return sc.Err()
}
