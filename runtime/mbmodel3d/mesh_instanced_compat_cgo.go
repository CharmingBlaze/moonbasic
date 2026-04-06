//go:build cgo

package mbmodel3d

import rl "github.com/gen2brain/raylib-go/raylib"

func drawMeshInstancedCompat(mesh rl.Mesh, material rl.Material, transforms []rl.Matrix, n int) {
	rl.DrawMeshInstanced(mesh, material, transforms, n)
}
