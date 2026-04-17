//go:build cgo || (windows && !cgo)

package mbaudio

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// playback implementation below

func (m *Module) audioPlay(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("AUDIO.PLAY expects 1 handle (sound or music)")
	}
	initAudioOnce()
	hid := heap.Handle(args[0].IVal)
	if so, err := heap.Cast[*soundObj](m.h, hid); err == nil {
		rl.PlaySound(so.snd)
		return args[0], nil
	}
	if mo, err := heap.Cast[*musicObj](m.h, hid); err == nil {
		rl.PlayMusicStream(mo.m)
		return args[0], nil
	}
	return value.Nil, fmt.Errorf("AUDIO.PLAY: handle must be sound or music")
}

func (m *Module) audioStop(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("AUDIO.STOP expects 1 handle (sound or music)")
	}
	hid := heap.Handle(args[0].IVal)
	if so, err := heap.Cast[*soundObj](m.h, hid); err == nil {
		rl.StopSound(so.snd)
		return args[0], nil
	}
	if mo, err := heap.Cast[*musicObj](m.h, hid); err == nil {
		rl.StopMusicStream(mo.m)
		return args[0], nil
	}
	return value.Nil, fmt.Errorf("AUDIO.STOP: handle must be sound or music")
}

func (m *Module) audioPause(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("AUDIO.PAUSE expects 1 handle (sound or music)")
	}
	hid := heap.Handle(args[0].IVal)
	if _, err := heap.Cast[*soundObj](m.h, hid); err == nil {
		return value.Nil, fmt.Errorf("AUDIO.PAUSE: sounds do not support pause; use AUDIO.STOP")
	}
	if mo, err := heap.Cast[*musicObj](m.h, hid); err == nil {
		rl.PauseMusicStream(mo.m)
		return args[0], nil
	}
	return value.Nil, fmt.Errorf("AUDIO.PAUSE: handle must be music")
}

func (m *Module) audioResume(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("AUDIO.RESUME expects 1 handle (sound or music)")
	}
	hid := heap.Handle(args[0].IVal)
	if _, err := heap.Cast[*soundObj](m.h, hid); err == nil {
		return value.Nil, fmt.Errorf("AUDIO.RESUME: use AUDIO.PLAY for sounds")
	}
	if mo, err := heap.Cast[*musicObj](m.h, hid); err == nil {
		rl.ResumeMusicStream(mo.m)
		return args[0], nil
	}
	return value.Nil, fmt.Errorf("AUDIO.RESUME: handle must be music")
}
func (m *Module) audioSetVolume(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("SOUNDVOLUME expects (handle, volume#)")
	}
	hid := heap.Handle(args[0].IVal)
	vol, ok := args[1].ToFloat()
	if !ok {
		return value.Nil, fmt.Errorf("SOUNDVOLUME: volume must be numeric")
	}
	if so, err := heap.Cast[*soundObj](m.h, hid); err == nil {
		rl.SetSoundVolume(so.snd, float32(vol))
		return args[0], nil
	}
	if mo, err := heap.Cast[*musicObj](m.h, hid); err == nil {
		rl.SetMusicVolume(mo.m, float32(vol))
		return args[0], nil
	}
	return value.Nil, fmt.Errorf("SOUNDVOLUME: handle must be sound or music")
}
