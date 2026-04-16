//go:build !cgo && !windows

package mbaudio

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "AUDIO.* requires CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

func raylibAudioOpen()  {}
func raylibAudioClose() {}

func (m *Module) audioInit(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("AUDIO.INIT: %s", hint)
}

func (m *Module) audioClose(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) audioPlay(args []value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("AUDIO.PLAY: %s", hint)
}

func (m *Module) audioStop(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) audioPause(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) audioResume(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) audioSetVolume(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) soundLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("AUDIO.LOADSOUND: %s", hint)
}

func (m *Module) soundFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) musicLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("AUDIO.LOADMUSIC: %s", hint)
}

func (m *Module) musicUpdate(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) musicFree(args []value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) setSoundVolume(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) setSoundPitch(args []value.Value) (value.Value, error)  { return value.Nil, nil }
func (m *Module) setSoundPan(args []value.Value) (value.Value, error)    { return value.Nil, nil }
func (m *Module) setMusicVolume(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) setMusicPitch(args []value.Value) (value.Value, error)  { return value.Nil, nil }
func (m *Module) setMasterVolume(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) isSoundPlaying(args []value.Value) (value.Value, error)  { return value.FromBool(false), nil }
func (m *Module) isMusicPlaying(args []value.Value) (value.Value, error)  { return value.FromBool(false), nil }
func (m *Module) getMusicLength(args []value.Value) (value.Value, error)  { return value.FromFloat(0), nil }
func (m *Module) getMusicTime(args []value.Value) (value.Value, error)    { return value.FromFloat(0), nil }
func (m *Module) seekMusic(args []value.Value) (value.Value, error)      { return value.Nil, nil }

func (m *Module) audioPlayVarySound(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) audioPlayRndSound(args []value.Value) (value.Value, error)  { return value.Nil, nil }

func (m *Module) audioListenerCamera(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}

func (m *Module) soundLoad3D(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("Load3DSound: %s", hint)
}

func (m *Module) streamMake(args []value.Value) (value.Value, error)      { return value.Nil, nil }
func (m *Module) streamUpdate(args []value.Value) (value.Value, error)    { return value.Nil, nil }
func (m *Module) streamIsReady(args []value.Value) (value.Value, error)   { return value.FromBool(false), nil }
func (m *Module) streamIsPlaying(args []value.Value) (value.Value, error) { return value.FromBool(false), nil }
func (m *Module) streamPlay(args []value.Value) (value.Value, error)      { return value.Nil, nil }
func (m *Module) streamPause(args []value.Value) (value.Value, error)     { return value.Nil, nil }
func (m *Module) streamResume(args []value.Value) (value.Value, error)    { return value.Nil, nil }
func (m *Module) streamStop(args []value.Value) (value.Value, error)      { return value.Nil, nil }
func (m *Module) streamSetVolume(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) streamSetPitch(args []value.Value) (value.Value, error)  { return value.Nil, nil }
func (m *Module) streamSetPan(args []value.Value) (value.Value, error)    { return value.Nil, nil }
func (m *Module) streamFree(args []value.Value) (value.Value, error)      { return value.Nil, nil }

func (m *Module) waveLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, fmt.Errorf("WAVE.LOAD: %s", hint)
}
func (m *Module) waveCopy(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) waveCrop(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) waveFormat(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) waveExport(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	return value.Nil, nil
}
func (m *Module) waveFree(args []value.Value) (value.Value, error) { return value.Nil, nil }

func (m *Module) soundFromWave(args []value.Value) (value.Value, error) { return value.Nil, nil }

func (m *Module) soundPlay3D(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) soundAttach(args []value.Value) (value.Value, error)   { return value.Nil, nil }
func (m *Module) worldSetAmbience(args []value.Value) (value.Value, error) { return value.Nil, nil }
func (m *Module) worldSetReverb(args []value.Value) (value.Value, error) { return value.Nil, nil }

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
