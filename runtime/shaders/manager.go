//go:build cgo || (windows && !cgo)

package shaders

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Manager handles the caching of shaders and uniform auto-injection.
type Manager struct {
	cache map[int32]rl.Shader
}

// NewManager creates a new Shader registry.
func NewManager() *Manager {
	return &Manager{cache: make(map[int32]rl.Shader)}
}

// LoadEmbedded loads a predefined standard library shader.
func (m *Manager) LoadEmbedded(id int32, vsName, fsName string) error {
	var vsCode, fsCode string
	
	if vsName != "" {
		b, err := EmbeddedShaders.ReadFile("shd/" + vsName)
		if err != nil { return err }
		vsCode = string(b)
	}
	if fsName != "" {
		b, err := EmbeddedShaders.ReadFile("shd/" + fsName)
		if err != nil { return err }
		fsCode = string(b)
	}

	sh := rl.LoadShaderFromMemory(vsCode, fsCode)
	if sh.ID == 0 {
		return fmt.Errorf("failed to compile embedded shader")
	}

	// Validate Auto-Bindings locations
	locTime := rl.GetShaderLocation(sh, "uTime")
	locRes  := rl.GetShaderLocation(sh, "uResolution")
	if locTime >= 0 || locRes >= 0 {
		// Log mapping success
	}

	m.cache[id] = sh
	return nil
}

// Get returns the cached shader by ID.
func (m *Manager) Get(id int32) (rl.Shader, bool) {
	s, ok := m.cache[id]
	return s, ok
}

// InjectGlobalUniforms injects system variables into the specified shader id.
func (m *Manager) InjectGlobalUniforms(id int32, totalTime float32, rx, ry float32) {
	sh, ok := m.cache[id]
	if !ok { return }

	locTime := rl.GetShaderLocation(sh, "uTime")
	if locTime >= 0 {
		rl.SetShaderValue(sh, locTime, []float32{totalTime}, rl.ShaderUniformFloat)
	}

	locRes := rl.GetShaderLocation(sh, "uResolution")
	if locRes >= 0 {
		rl.SetShaderValue(sh, locRes, []float32{rx, ry}, rl.ShaderUniformVec2)
	}
}

// LoadCustom securely compiles a user shader, falling back if it fails.
func (m *Manager) LoadCustom(id int32, vsPath, fsPath string) error {
	sh := rl.LoadShader(vsPath, fsPath)
	if sh.ID == 0 {
		// Fallback to PBR_LIT logic
		return m.LoadEmbedded(id, "", "pbr_lit.fs")
	}
	m.cache[id] = sh
	return nil
}

// Free flushes all caches.
func (m *Manager) Free() {
	for _, sh := range m.cache {
		rl.UnloadShader(sh)
	}
	m.cache = make(map[int32]rl.Shader)
}

// Standard Library ID Defines
const (
	SHADER_PBR_LIT          int32 = 1
	SHADER_PS1_RETRO        int32 = 2
	SHADER_CEL_STYLED       int32 = 3
	SHADER_WATER_PROCEDURAL int32 = 4
	
	PP_BLOOM          int32 = 101
	PP_CRT_SCANLINES  int32 = 102
	PP_PIXELATE       int32 = 103
)
