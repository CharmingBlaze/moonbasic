package mbgame

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

type configStore struct {
	path  string
	order []string
	kv    map[string]string
}

func (m *Module) registerConfig(r runtime.Registrar) {
	r.Register("CONFIG.LOAD", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("CONFIG.LOAD expects path$")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return value.Nil, err
		}
		m.config = &configStore{path: path, kv: make(map[string]string)}
		sc := bufio.NewScanner(strings.NewReader(string(data)))
		for sc.Scan() {
			line := strings.TrimSpace(sc.Text())
			if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
				continue
			}
			i := strings.IndexByte(line, '=')
			if i < 0 {
				continue
			}
			k := strings.TrimSpace(line[:i])
			v := strings.TrimSpace(line[i+1:])
			if k == "" {
				continue
			}
			if _, ok := m.config.kv[k]; !ok {
				m.config.order = append(m.config.order, k)
			}
			m.config.kv[k] = v
		}
		return value.Nil, sc.Err()
	})
	r.Register("CONFIG.SAVE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.config == nil {
			return value.Nil, fmt.Errorf("CONFIG.SAVE: no CONFIG.LOAD yet")
		}
		path := m.config.path
		if len(args) >= 1 && args[0].Kind == value.KindString {
			p, err := rt.ArgString(args, 0)
			if err != nil {
				return value.Nil, err
			}
			path = p
		}
		var b strings.Builder
		for _, k := range m.config.order {
			b.WriteString(k)
			b.WriteByte('=')
			b.WriteString(m.config.kv[k])
			b.WriteByte('\n')
		}
		return value.Nil, os.WriteFile(path, []byte(b.String()), 0644)
	})
	r.Register("CONFIG.SETINT", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.config == nil {
			m.config = &configStore{kv: make(map[string]string)}
		}
		if len(args) != 2 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("CONFIG.SETINT expects (key$, value)")
		}
		k, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		n, ok := argI(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("CONFIG.SETINT: integer value")
		}
		m.setCfg(k, fmt.Sprintf("%d", n))
		return value.Nil, nil
	})
	r.Register("CONFIG.SETFLOAT", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.config == nil {
			m.config = &configStore{kv: make(map[string]string)}
		}
		if len(args) != 2 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("CONFIG.SETFLOAT expects (key$, value)")
		}
		k, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		x, ok := argF(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("CONFIG.SETFLOAT: float value")
		}
		m.setCfg(k, fmt.Sprintf("%g", x))
		return value.Nil, nil
	})
	r.Register("CONFIG.SETSTRING", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.config == nil {
			m.config = &configStore{kv: make(map[string]string)}
		}
		if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
			return value.Nil, fmt.Errorf("CONFIG.SETSTRING expects (key$, value$)")
		}
		k, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		v, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
		m.setCfg(k, v)
		return value.Nil, nil
	})
	r.Register("CONFIG.SETBOOL", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if m.config == nil {
			m.config = &configStore{kv: make(map[string]string)}
		}
		if len(args) != 2 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("CONFIG.SETBOOL expects (key$, value?)")
		}
		k, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		var s string
		switch args[1].Kind {
		case value.KindBool:
			if args[1].IVal != 0 {
				s = "true"
			} else {
				s = "false"
			}
		default:
			return value.Nil, fmt.Errorf("CONFIG.SETBOOL: bool value")
		}
		m.setCfg(k, s)
		return value.Nil, nil
	})
	r.Register("CONFIG.GETINT", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("CONFIG.GETINT expects (key$, default)")
		}
		k, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		def, ok := argI(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("CONFIG.GETINT: default int")
		}
		if m.config == nil {
			return value.FromInt(def), nil
		}
		v, ok := m.config.kv[k]
		if !ok {
			return value.FromInt(def), nil
		}
		var n int64
		_, err = fmt.Sscan(v, &n)
		if err != nil {
			return value.FromInt(def), nil
		}
		return value.FromInt(n), nil
	})
	r.Register("CONFIG.GETFLOAT", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("CONFIG.GETFLOAT expects (key$, default)")
		}
		k, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		def, ok := argF(args[1])
		if !ok {
			return value.Nil, fmt.Errorf("CONFIG.GETFLOAT: default float")
		}
		if m.config == nil {
			return value.FromFloat(def), nil
		}
		v, ok := m.config.kv[k]
		if !ok {
			return value.FromFloat(def), nil
		}
		var x float64
		_, err = fmt.Sscan(v, &x)
		if err != nil {
			return value.FromFloat(def), nil
		}
		return value.FromFloat(x), nil
	})
	r.Register("CONFIG.GETSTRING", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 || args[0].Kind != value.KindString || args[1].Kind != value.KindString {
			return value.Nil, fmt.Errorf("CONFIG.GETSTRING expects (key$, default$)")
		}
		k, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		def, err := rt.ArgString(args, 1)
		if err != nil {
			return value.Nil, err
		}
		if m.config == nil {
			return rt.RetString(def), nil
		}
		v, ok := m.config.kv[k]
		if !ok {
			return rt.RetString(def), nil
		}
		return rt.RetString(v), nil
	})
	r.Register("CONFIG.GETBOOL", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 2 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("CONFIG.GETBOOL expects (key$, default?)")
		}
		k, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		def := args[1].Kind == value.KindBool && args[1].IVal != 0
		if m.config == nil {
			return value.FromBool(def), nil
		}
		v, ok := m.config.kv[k]
		if !ok {
			return value.FromBool(def), nil
		}
		v = strings.ToLower(strings.TrimSpace(v))
		return value.FromBool(v == "1" || v == "true" || v == "yes" || v == "on"), nil
	})
	r.Register("CONFIG.HAS", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("CONFIG.HAS expects key$")
		}
		k, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		if m.config == nil {
			return value.FromBool(false), nil
		}
		_, ok := m.config.kv[k]
		return value.FromBool(ok), nil
	})
	r.Register("CONFIG.DELETE", "game", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("CONFIG.DELETE expects key$")
		}
		k, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		if m.config == nil {
			return value.Nil, nil
		}
		delete(m.config.kv, k)
		for i, x := range m.config.order {
			if x == k {
				m.config.order = append(m.config.order[:i], m.config.order[i+1:]...)
				break
			}
		}
		return value.Nil, nil
	})
}

func (m *Module) setCfg(k, v string) {
	if m.config.kv == nil {
		m.config.kv = make(map[string]string)
	}
	if _, ok := m.config.kv[k]; !ok {
		m.config.order = append(m.config.order, k)
	}
	m.config.kv[k] = v
}
