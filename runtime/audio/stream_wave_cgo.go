//go:build cgo || (windows && !cgo)

package mbaudio

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// stream/wave/sound implementation below

func (m *Module) getStream(args []value.Value, ix int, op string) (*audioStreamObj, error) {
	if err := m.requireHeap(); err != nil {
		return nil, err
	}
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be audio stream handle", op, ix+1)
	}
	return heap.Cast[*audioStreamObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) streamMake(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.MAKE expects 3 arguments (sampleRate, bitDepth, channels)")
	}
	sr, ok1 := argUInt32(args[0])
	bd, ok2 := argUInt32(args[1])
	ch, ok3 := argUInt32(args[2])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.MAKE: arguments must be non-negative integers")
	}
	if bd != 8 && bd != 16 && bd != 32 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.MAKE: bitDepth must be 8, 16, or 32")
	}
	if ch < 1 || ch > 2 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.MAKE: channels must be 1 or 2")
	}
	s := rl.LoadAudioStream(sr, bd, ch)
	id, err := m.h.Alloc(&audioStreamObj{s: s})
	if err != nil {
		rl.UnloadAudioStream(s)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) streamUpdate(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.UPDATE expects (stream, arrayHandle) — 1-D numeric array of PCM samples")
	}
	o, err := m.getStream(args, 0, "AUDIOSTREAM.UPDATE")
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.UPDATE: second argument must be array handle")
	}
	arr, err := heap.Cast[*heap.Array](m.h, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.UPDATE: %w", err)
	}
	buf, err := pcmSliceFrom1DArray(arr, o.s.SampleSize, "AUDIOSTREAM.UPDATE")
	if err != nil {
		return value.Nil, err
	}
	rl.UpdateAudioStream(o.s, buf)
	return value.Nil, nil
}

func (m *Module) streamIsReady(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.ISREADY expects stream handle")
	}
	o, err := m.getStream(args, 0, "AUDIOSTREAM.ISREADY")
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsAudioStreamValid(o.s)), nil
}

func (m *Module) streamIsPlaying(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.ISPLAYING expects stream handle")
	}
	o, err := m.getStream(args, 0, "AUDIOSTREAM.ISPLAYING")
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsAudioStreamPlaying(o.s)), nil
}

func (m *Module) streamPlay(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.PLAY expects stream handle")
	}
	o, err := m.getStream(args, 0, "AUDIOSTREAM.PLAY")
	if err != nil {
		return value.Nil, err
	}
	rl.PlayAudioStream(o.s)
	return value.Nil, nil
}

func (m *Module) streamPause(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.PAUSE expects stream handle")
	}
	o, err := m.getStream(args, 0, "AUDIOSTREAM.PAUSE")
	if err != nil {
		return value.Nil, err
	}
	rl.PauseAudioStream(o.s)
	return value.Nil, nil
}

func (m *Module) streamResume(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.RESUME expects stream handle")
	}
	o, err := m.getStream(args, 0, "AUDIOSTREAM.RESUME")
	if err != nil {
		return value.Nil, err
	}
	rl.ResumeAudioStream(o.s)
	return value.Nil, nil
}

func (m *Module) streamStop(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.STOP expects stream handle")
	}
	o, err := m.getStream(args, 0, "AUDIOSTREAM.STOP")
	if err != nil {
		return value.Nil, err
	}
	rl.StopAudioStream(o.s)
	return value.Nil, nil
}

func (m *Module) streamSetVolume(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.SETVOLUME expects (stream, volume)")
	}
	o, err := m.getStream(args, 0, "AUDIOSTREAM.SETVOLUME")
	if err != nil {
		return value.Nil, err
	}
	v, ok := argFloat32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.SETVOLUME: volume must be numeric")
	}
	rl.SetAudioStreamVolume(o.s, v)
	return value.Nil, nil
}

func (m *Module) streamSetPitch(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.SETPITCH expects (stream, pitch)")
	}
	o, err := m.getStream(args, 0, "AUDIOSTREAM.SETPITCH")
	if err != nil {
		return value.Nil, err
	}
	v, ok := argFloat32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.SETPITCH: pitch must be numeric")
	}
	rl.SetAudioStreamPitch(o.s, v)
	return value.Nil, nil
}

func (m *Module) streamSetPan(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.SETPAN expects (stream, pan)")
	}
	o, err := m.getStream(args, 0, "AUDIOSTREAM.SETPAN")
	if err != nil {
		return value.Nil, err
	}
	v, ok := argFloat32(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.SETPAN: pan must be numeric")
	}
	rl.SetAudioStreamPan(o.s, v)
	return value.Nil, nil
}

func (m *Module) streamFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("AUDIOSTREAM.FREE expects stream handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) getWave(args []value.Value, ix int, op string) (*waveObj, error) {
	if err := m.requireHeap(); err != nil {
		return nil, err
	}
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: argument %d must be wave handle", op, ix+1)
	}
	return heap.Cast[*waveObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) waveLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("WAVE.LOAD expects file path string")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	w := rl.LoadWave(path)
	id, err := m.h.Alloc(&waveObj{w: w})
	if err != nil {
		rl.UnloadWave(w)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) waveCopy(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WAVE.COPY expects wave handle")
	}
	o, err := m.getWave(args, 0, "WAVE.COPY")
	if err != nil {
		return value.Nil, err
	}
	c := rl.WaveCopy(o.w)
	id, err := m.h.Alloc(&waveObj{w: c})
	if err != nil {
		rl.UnloadWave(c)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) waveCrop(args []value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("WAVE.CROP expects (wave, startFrame, endFrame)")
	}
	o, err := m.getWave(args, 0, "WAVE.CROP")
	if err != nil {
		return value.Nil, err
	}
	start, ok1 := argInt32(args[1])
	end, ok2 := argInt32(args[2])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("WAVE.CROP: frame indices must be numeric")
	}
	rl.WaveCrop(&o.w, start, end)
	return value.Nil, nil
}

func (m *Module) waveFormat(args []value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("WAVE.FORMAT expects (wave, sampleRate, bitDepth, channels)")
	}
	o, err := m.getWave(args, 0, "WAVE.FORMAT")
	if err != nil {
		return value.Nil, err
	}
	sr, ok1 := argInt32(args[1])
	bd, ok2 := argInt32(args[2])
	ch, ok3 := argInt32(args[3])
	if !ok1 || !ok2 || !ok3 {
		return value.Nil, fmt.Errorf("WAVE.FORMAT: format parameters must be numeric")
	}
	rl.WaveFormat(&o.w, sr, bd, ch)
	return value.Nil, nil
}

func (m *Module) waveExport(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WAVE.EXPORT expects (wave, path$)")
	}
	o, err := m.getWave(args, 0, "WAVE.EXPORT")
	if err != nil {
		return value.Nil, err
	}
	if args[1].Kind != value.KindString {
		return value.Nil, fmt.Errorf("WAVE.EXPORT: path must be string")
	}
	path, err := rt.ArgString(args, 1)
	if err != nil {
		return value.Nil, err
	}
	rl.ExportWave(o.w, path)
	return value.Nil, nil
}

func (m *Module) waveFree(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("WAVE.FREE expects wave handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) soundFromWave(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("SOUND.FROMWAVE expects wave handle")
	}
	o, err := m.getWave(args, 0, "SOUND.FROMWAVE")
	if err != nil {
		return value.Nil, err
	}
	snd := rl.LoadSoundFromWave(o.w)
	id, err := m.h.Alloc(&soundObj{snd: snd, gain: 1, pan: 0})
	if err != nil {
		rl.UnloadSound(snd)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
