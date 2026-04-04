//go:build cgo

package window

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type automationListObj struct {
	mod     *Module
	hid     heap.Handle
	list    rl.AutomationEventList
	release heap.ReleaseOnce
}

func (o *automationListObj) TypeName() string { return "AutomationEventList" }

func (o *automationListObj) TypeTag() uint16 { return heap.TagAutomationList }

func (o *automationListObj) Free() {
	o.release.Do(func() {
		o.mod.detachAutomationList(o.hid)
		rl.UnloadAutomationEventList(&o.list)
	})
}

func (m *Module) detachAutomationList(h heap.Handle) {
	m.autoMu.Lock()
	defer m.autoMu.Unlock()
	if m.activeAutoHandle == h {
		m.activeAutoHandle = 0
		if m.automationRec {
			rl.StopAutomationEventRecording()
			m.automationRec = false
		}
	}
}

func (m *Module) registerAutomationCommands(reg runtime.Registrar) {
	reg.Register("EVENT.LISTMAKE", "event", m.evListMake)
	reg.Register("EVENT.LISTLOAD", "event", m.evListLoad)
	reg.Register("EVENT.LISTEXPORT", "event", m.evListExport)
	reg.Register("EVENT.SETACTIVELIST", "event", m.evSetActiveList)
	reg.Register("EVENT.RECSTART", "event", m.evRecStart)
	reg.Register("EVENT.RECSTOP", "event", m.evRecStop)
	reg.Register("EVENT.REPLAY", "event", m.evReplay)
	recording := m.evRecPlaying
	reg.Register("EVENT.RECPLAYING", "event", recording)
	reg.Register("EVENT.ISPLAYING", "event", recording)
	reg.Register("EVENT.LISTCLEAR", "event", m.evListClear)
	reg.Register("EVENT.LISTCOUNT", "event", m.evListCount)
	reg.Register("EVENT.LISTFREE", "event", m.evListFree)
}

func (m *Module) allocAutomationList(list rl.AutomationEventList) (heap.Handle, error) {
	if m.h == nil {
		return 0, runtime.Errorf("EVENT.*: heap not bound")
	}
	obj := &automationListObj{mod: m, list: list}
	h, err := m.h.Alloc(obj)
	if err != nil {
		return 0, err
	}
	obj.hid = h
	return h, nil
}

func (m *Module) getAutoList(args []value.Value, name string) (*automationListObj, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s expects automation list handle", name)
	}
	o, err := heap.Cast[*automationListObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	return o, nil
}

// EVENT.LISTMAKE — empty list (optional path reserved for default export target; use LISTEXPORT).
func (m *Module) evListMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("EVENT.LISTMAKE expects 1 argument (path$; may be empty for new list)")
	}
	_, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	list := rl.LoadAutomationEventList("")
	h, err := m.allocAutomationList(list)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) evListLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("EVENT.LISTLOAD expects 1 argument (path$)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	list := rl.LoadAutomationEventList(path)
	h, err := m.allocAutomationList(list)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(h), nil
}

func (m *Module) evListExport(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 || args[0].Kind != value.KindHandle || args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("EVENT.LISTEXPORT expects (list, path$)")
	}
	o, err := m.getAutoList(args[:1], "EVENT.LISTEXPORT")
	if err != nil {
		return value.Nil, err
	}
	path, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	if !rl.ExportAutomationEventList(o.list, path) {
		return value.Nil, runtime.Errorf("EVENT.LISTEXPORT failed for %q", path)
	}
	return value.Nil, nil
}

func (m *Module) evSetActiveList(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("EVENT.SETACTIVELIST expects automation list handle")
	}
	o, err := m.getAutoList(args, "EVENT.SETACTIVELIST")
	if err != nil {
		return value.Nil, err
	}
	m.autoMu.Lock()
	defer m.autoMu.Unlock()
	m.activeAutoHandle = o.hid
	rl.SetAutomationEventList(&o.list)
	return value.Nil, nil
}

func (m *Module) evRecStart(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("EVENT.RECSTART expects 0 arguments")
	}
	m.autoMu.Lock()
	defer m.autoMu.Unlock()
	if m.activeAutoHandle == 0 {
		return value.Nil, runtime.Errorf("EVENT.RECSTART: call EVENT.SETACTIVELIST first")
	}
	rl.StartAutomationEventRecording()
	m.automationRec = true
	return value.Nil, nil
}

func (m *Module) evRecStop(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("EVENT.RECSTOP expects 0 arguments")
	}
	m.autoMu.Lock()
	defer m.autoMu.Unlock()
	rl.StopAutomationEventRecording()
	m.automationRec = false
	return value.Nil, nil
}

func (m *Module) evReplay(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	o, err := m.getAutoList(args, "EVENT.REPLAY")
	if err != nil {
		return value.Nil, err
	}
	for _, ev := range o.list.GetEvents() {
		rl.PlayAutomationEvent(ev)
	}
	return value.Nil, nil
}

func (m *Module) evRecPlaying(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("EVENT.RECPLAYING expects 0 arguments")
	}
	m.autoMu.Lock()
	defer m.autoMu.Unlock()
	return value.FromBool(m.automationRec), nil
}

func (m *Module) evListClear(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	o, err := m.getAutoList(args, "EVENT.LISTCLEAR")
	if err != nil {
		return value.Nil, err
	}
	m.detachAutomationList(o.hid)
	rl.UnloadAutomationEventList(&o.list)
	o.list = rl.LoadAutomationEventList("")
	return value.Nil, nil
}

func (m *Module) evListCount(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	o, err := m.getAutoList(args, "EVENT.LISTCOUNT")
	if err != nil {
		return value.Nil, err
	}
	return value.FromInt(int64(o.list.Count)), nil
}

func (m *Module) evListFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	_ = rt
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("EVENT.LISTFREE expects automation list handle")
	}
	h := heap.Handle(args[0].IVal)
	if m.h == nil {
		return value.Nil, runtime.Errorf("EVENT.LISTFREE: heap not bound")
	}
	return value.Nil, m.h.Free(h)
}

// shutdownAutomation stops recording and forgets the active list pointer (window closing).
func (m *Module) shutdownAutomation() {
	m.autoMu.Lock()
	defer m.autoMu.Unlock()
	if m.automationRec {
		rl.StopAutomationEventRecording()
		m.automationRec = false
	}
	m.activeAutoHandle = 0
}
