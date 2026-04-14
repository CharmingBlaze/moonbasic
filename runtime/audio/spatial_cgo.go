//go:build cgo || (windows && !cgo)

package mbaudio

import (
	"fmt"
	"math"
	"sync"

	"moonbasic/runtime"
	mbcamera "moonbasic/runtime/camera"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const spatialMaxDist = float32(80)

var (
	spatialMu   sync.Mutex
	listenPos   rl.Vector3
	listenFwd   rl.Vector3 // horizontal forward (XZ), normalized
	listenerOK  bool
)

// SetSpatialListener sets the virtual microphone used by PlaySpatial (Listener builtin).
func SetSpatialListener(pos, forward rl.Vector3) {
	spatialMu.Lock()
	defer spatialMu.Unlock()
	listenPos = pos
	fx, fz := forward.X, forward.Z
	lenH := float32(math.Hypot(float64(fx), float64(fz)))
	if lenH < 1e-5 {
		listenFwd = rl.Vector3{X: 0, Y: 0, Z: -1}
	} else {
		listenFwd = rl.Vector3{X: fx / lenH, Y: 0, Z: fz / lenH}
	}
	listenerOK = true
}

// PlaySpatial applies distance attenuation + stereo pan, plays once, then restores gain/pan.
func PlaySpatial(h *heap.Store, soundH heap.Handle, wx, wy, wz float32) error {
	initAudioOnce()
	so, err := heap.Cast[*soundObj](h, soundH)
	if err != nil {
		return err
	}
	gain := so.gain
	panBase := so.pan

	spatialMu.Lock()
	ok := listenerOK
	lp := listenPos
	lf := listenFwd
	spatialMu.Unlock()

	vol := gain
	pan := panBase
	if ok {
		src := rl.Vector3{X: wx, Y: wy, Z: wz}
		to := rl.Vector3Subtract(src, lp)
		d := rl.Vector3Length(to)
		if spatialMaxDist > 1e-5 && d > 1e-5 {
			t := d / spatialMaxDist
			if t > 1 {
				t = 1
			}
			vol = gain * (1.0 - t*t)
		}
		// Stereo pan from horizontal angle vs listener right (XZ).
		tx, tz := to.X, to.Z
		lenH := float32(math.Hypot(float64(tx), float64(tz)))
		if lenH > 1e-5 {
			tx /= lenH
			tz /= lenH
			// right = (-lf.Z, lf.X) in XZ
			rx, rz := -lf.Z, lf.X
			pan = panBase + (tx*rx + tz*rz)
			if pan > 1 {
				pan = 1
			}
			if pan < -1 {
				pan = -1
			}
		}
	}

	rl.SetSoundVolume(so.snd, vol)
	rl.SetSoundPan(so.snd, pan)
	rl.PlaySound(so.snd)
	rl.SetSoundVolume(so.snd, gain)
	rl.SetSoundPan(so.snd, panBase)
	return nil
}

// spatial implementation below

func (m *Module) audioListenerCamera(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("Listener expects (cameraHandle)")
	}
	cam, err := mbcamera.RayCamera3D(m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	fwd := rl.Vector3Subtract(cam.Target, cam.Position)
	SetSpatialListener(cam.Position, fwd)
	_ = rt
	return value.Nil, nil
}

func (m *Module) soundLoad3D(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("Load3DSound expects path$ (same as AUDIO.LOADSOUND; spatial mix uses EmitSound + Listener)")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	initAudioOnce()
	snd := rl.LoadSound(path)
	id, err := m.h.Alloc(&soundObj{snd: snd, gain: 1, pan: 0, spatial3D: true})
	if err != nil {
		rl.UnloadSound(snd)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
