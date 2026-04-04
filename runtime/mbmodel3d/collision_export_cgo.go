//go:build cgo

package mbmodel3d

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

// MeshRaylib returns a copy of the Raylib mesh for a TagMesh handle (ray collision, etc.).
func MeshRaylib(s *heap.Store, h heap.Handle) (rl.Mesh, error) {
	o, err := heap.Cast[*meshObj](s, h)
	if err != nil {
		return rl.Mesh{}, err
	}
	return o.m, nil
}

// ModelRaylib returns a copy of the Raylib model for a TagModel handle.
func ModelRaylib(s *heap.Store, h heap.Handle) (rl.Model, error) {
	o, err := heap.Cast[*modelObj](s, h)
	if err != nil {
		return rl.Model{}, err
	}
	return o.model, nil
}

// ShaderRaylib returns the Raylib shader for a TagShader handle.
func ShaderRaylib(s *heap.Store, h heap.Handle) (rl.Shader, error) {
	o, err := heap.Cast[*shaderObj](s, h)
	if err != nil {
		return rl.Shader{}, err
	}
	return o.sh, nil
}
