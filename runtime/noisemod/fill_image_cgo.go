//go:build cgo || (windows && !cgo)

package noisemod

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/runtime/mbimage"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) noiseFillImage(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("NOISE.FILLIMAGE expects (noise, img, offsetX#, offsetY#)")
	}
	n, err := m.getNoise(rt, args, 0)
	if err != nil {
		return value.Nil, err
	}
	if err := n.assertLive(); err != nil {
		return value.Nil, err
	}
	hs, err := fixAllocHeap(rt, m)
	if err != nil {
		return value.Nil, err
	}
	img, err := mbimage.RayImageForTexture(hs, heap.Handle(args[1].IVal))
	if err != nil {
		return value.Nil, err
	}
	offx, err := rt.ArgFloat(args, 2)
	if err != nil {
		return value.Nil, err
	}
	offy, err := rt.ArgFloat(args, 3)
	if err != nil {
		return value.Nil, err
	}
	w := int(img.Width)
	h := int(img.Height)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			v := normNoise(n.Sample2D(float64(x)+offx, float64(y)+offy))
			g := uint8(v * 255)
			rl.ImageDrawPixel(img, int32(x), int32(y), rl.Color{R: g, G: g, B: g, A: 255})
		}
	}
	return args[0], nil
}
