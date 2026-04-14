//go:build cgo || (windows && !cgo)

package mbaudio

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// audio properties implementation below

func (m *Module) setSoundVolume(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("AUDIO.SETSOUNDVOLUME expects (sound, volume#)")
	}
	so, err := heap.Cast[*soundObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v, ok := argFloat32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("AUDIO.SETSOUNDVOLUME: volume must be numeric")
	}
	so.gain = v
	rl.SetSoundVolume(so.snd, v)
	return value.Nil, nil
}

func (m *Module) setSoundPitch(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("AUDIO.SETSOUNDPITCH expects (sound, pitch#)")
	}
	so, err := heap.Cast[*soundObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v, ok := argFloat32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("AUDIO.SETSOUNDPITCH: pitch must be numeric")
	}
	rl.SetSoundPitch(so.snd, v)
	return value.Nil, nil
}

func (m *Module) setSoundPan(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("AUDIO.SETSOUNDPAN expects (sound, pan#)")
	}
	so, err := heap.Cast[*soundObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v, ok := argFloat32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("AUDIO.SETSOUNDPAN: pan must be numeric")
	}
	so.pan = v
	rl.SetSoundPan(so.snd, v)
	return value.Nil, nil
}

func (m *Module) setMusicVolume(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("AUDIO.SETMUSICVOLUME expects (music, volume#)")
	}
	mo, err := heap.Cast[*musicObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v, ok := argFloat32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("AUDIO.SETMUSICVOLUME: volume must be numeric")
	}
	rl.SetMusicVolume(mo.m, v)
	return value.Nil, nil
}

func (m *Module) setMusicPitch(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("AUDIO.SETMUSICPITCH expects (music, pitch#)")
	}
	mo, err := heap.Cast[*musicObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	v, ok := argFloat32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("AUDIO.SETMUSICPITCH: pitch must be numeric")
	}
	rl.SetMusicPitch(mo.m, v)
	return value.Nil, nil
}

func (m *Module) setMasterVolume(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("AUDIO.SETMASTERVOLUME expects (volume#)")
	}
	v, ok := argFloat32(args[0])
	if !ok {
		return value.Nil, fmt.Errorf("AUDIO.SETMASTERVOLUME: volume must be numeric")
	}
	initAudioOnce()
	rl.SetMasterVolume(v)
	return value.Nil, nil
}

func (m *Module) isSoundPlaying(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("AUDIO.ISSOUNDPLAYING expects sound handle")
	}
	so, err := heap.Cast[*soundObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsSoundPlaying(so.snd)), nil
}

func (m *Module) isMusicPlaying(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("AUDIO.ISMUSICPLAYING expects music handle")
	}
	mo, err := heap.Cast[*musicObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsMusicStreamPlaying(mo.m)), nil
}

func (m *Module) getMusicLength(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("AUDIO.GETMUSICLENGTH expects music handle")
	}
	mo, err := heap.Cast[*musicObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.GetMusicTimeLength(mo.m))), nil
}

func (m *Module) getMusicTime(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("AUDIO.GETMUSICTIME expects music handle")
	}
	mo, err := heap.Cast[*musicObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(float64(rl.GetMusicTimePlayed(mo.m))), nil
}

func (m *Module) seekMusic(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("AUDIO.SEEKMUSIC expects (music, seconds#)")
	}
	mo, err := heap.Cast[*musicObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	t, ok := argFloat32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("AUDIO.SEEKMUSIC: position must be numeric")
	}
	rl.SeekMusicStream(mo.m, t)
	return value.Nil, nil
}
