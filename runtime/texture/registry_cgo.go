//go:build cgo || (windows && !cgo)

package texture

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Asset struct {
	Tex  rl.Texture2D
	Path string
}

type Registry struct {
	mu     sync.RWMutex
	assets map[int32]Asset
	hashes map[string]int32
	nextID int32
}

func NewRegistry() *Registry {
	return &Registry{
		assets: make(map[int32]Asset),
		hashes: make(map[string]int32),
		nextID: 1,
	}
}

func (r *Registry) HashData(data []byte) string {
	h := sha256.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func (r *Registry) Register(tex rl.Texture2D, path, hash string) int32 {
	r.mu.Lock()
	defer r.mu.Unlock()
	if hash != "" {
		if id, exists := r.hashes[hash]; exists {
			return id
		}
	}
	id := r.nextID
	r.nextID++
	r.assets[id] = Asset{Tex: tex, Path: path}
	if hash != "" {
		r.hashes[hash] = id
	}
	return id
}

func (r *Registry) Get(id int32) (Asset, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	a, ok := r.assets[id]
	return a, ok
}
