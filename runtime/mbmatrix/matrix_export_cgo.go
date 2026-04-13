//go:build cgo || (windows && !cgo)

package mbmatrix

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

func AllocMatrix(s *heap.Store, m rl.Matrix) (heap.Handle, error) {
	return s.Alloc(&matObj{m: fromM(m)})
}

// MatrixRaylib returns the matrix for a TagMatrix handle, or identity if h==0.
func MatrixRaylib(s *heap.Store, h heap.Handle) (rl.Matrix, error) {
	if h == 0 {
		return rl.MatrixIdentity(), nil
	}
	o, err := heap.Cast[*matObj](s, h)
	if err != nil {
		return rl.Matrix{}, err
	}
	return toM(o.m), nil
}
