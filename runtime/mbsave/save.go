//go:build cgo || (windows && !cgo)

package mbsave

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

var (
	saveData = make(map[string]string)
	saveMu   sync.RWMutex
	savePath = "save.json"
)

type Module struct{}

func NewModule() *Module { return &Module{} }

func (m *Module) Register(r runtime.Registrar) {
	r.Register("SAVE.DATA", "save", m.saveData)
	r.Register("SAVE.GET", "save", m.saveGet)
	loadSaveData()
}

func (m *Module) Shutdown() {
	flushSaveData()
}

func loadSaveData() {
	saveMu.Lock()
	defer saveMu.Unlock()
	b, err := os.ReadFile(savePath)
	if err == nil {
		json.Unmarshal(b, &saveData)
	}
}

func flushSaveData() {
	saveMu.RLock()
	defer saveMu.RUnlock()
	b, err := json.MarshalIndent(saveData, "", "  ")
	if err == nil {
		os.WriteFile(savePath, b, 0644)
	}
}

func (m *Module) saveData(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 { return value.Nil, fmt.Errorf("SAVE.DATA expects (key$, value$)") }
	k, _ := rt.ArgString(args, 0)
	v, _ := rt.ArgString(args, 1)

	saveMu.Lock()
	saveData[k] = v
	saveMu.Unlock()
	
	return value.Nil, nil
}

func (m *Module) saveGet(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 { return value.Nil, fmt.Errorf("SAVE.GET expects (key$)") }
	k, _ := rt.ArgString(args, 0)
	
	saveMu.RLock()
	v, ok := saveData[k]
	saveMu.RUnlock()
	
	if ok {
		return rt.RetString(v), nil
	}
	return rt.RetString(""), nil
}
