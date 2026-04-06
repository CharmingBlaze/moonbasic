//go:build cgo || (windows && !cgo)

package mbaudio

import (
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var audioMu sync.Mutex
var audioInited bool

func initAudioOnce() {
	audioMu.Lock()
	defer audioMu.Unlock()
	if audioInited {
		return
	}
	rl.InitAudioDevice()
	audioInited = true
}

func closeAudioOnce() {
	audioMu.Lock()
	defer audioMu.Unlock()
	if !audioInited {
		return
	}
	rl.CloseAudioDevice()
	audioInited = false
}
