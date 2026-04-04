//go:build cgo

package mbaudio

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerMusic(r runtime.Registrar) {
	r.Register("AUDIO.LOADMUSIC", "audio", m.musicLoad)
	r.Register("AUDIO.UPDATEMUSIC", "audio", runtime.AdaptLegacy(m.musicUpdate))
	r.Register("MUSIC.FREE", "audio", runtime.AdaptLegacy(m.musicFree))
}

func (m *Module) musicLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("AUDIO.LOADMUSIC expects 1 string path")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	initAudioOnce()
	mu := rl.LoadMusicStream(path)
	id, err := m.h.Alloc(&musicObj{m: mu})
	if err != nil {
		rl.UnloadMusicStream(mu)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) musicUpdate(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("AUDIO.UPDATEMUSIC expects music handle")
	}
	o, err := heap.Cast[*musicObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	rl.UpdateMusicStream(o.m)
	return value.Nil, nil
}

func (m *Module) musicFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("MUSIC.FREE expects music handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
