//go:build !cgo && windows

package mbmodel3d

import rl "github.com/gen2brain/raylib-go/raylib"

func drawMeshInstancedMO(mo *meshObj, mato *materialObj, mats []rl.Matrix, cnt int32) {
	rl.DrawMeshInstanced(mo.m, mato.mat, mats, cnt)
}
