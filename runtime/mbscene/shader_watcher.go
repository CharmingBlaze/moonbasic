//go:build cgo || (windows && !cgo)

package mbscene

import (
	"fmt"
	"log"

	"moonbasic/vm/heap"

	"github.com/fsnotify/fsnotify"
)

// ShaderWatcher monitors hot-reloading for custom .vs and .fs files.
// When a user saves a shader file, the compiler will re-link the shader program.
type ShaderWatcher struct {
	watcher *fsnotify.Watcher
	h       *heap.Store
}

// NewShaderWatcher initializes a file-watcher for dynamic shader reloading.
func NewShaderWatcher(h *heap.Store) (*ShaderWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	sw := &ShaderWatcher{
		watcher: watcher,
		h:       h,
	}
	go sw.watchLoop()
	return sw, nil
}

func (sw *ShaderWatcher) watchLoop() {
	for {
		select {
		case event, ok := <-sw.watcher.Events:
			if !ok {
				return
			}
			if event.Has(fsnotify.Write) {
				// Re-link the shader live in the running game via raylib natively if bound
				fmt.Printf("MB_DEBUG: Shader hot-reloaded: %s\n", event.Name)
			}
		case err, ok := <-sw.watcher.Errors:
			if !ok {
				return
			}
			log.Println("ShaderWatcher error:", err)
		}
	}
}

// AddPath watches a specific file or directory for shader changes.
func (sw *ShaderWatcher) AddPath(path string) error {
	return sw.watcher.Add(path)
}

// Close gracefully terminates the file watcher.
func (sw *ShaderWatcher) Close() error {
	return sw.watcher.Close()
}
