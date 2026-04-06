//go:build cgo || (windows && !cgo)

package mbimage

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func registerImageQuery(m *Module, reg runtime.Registrar) {
	reg.Register("IMAGE.GETCOLORR", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("IMAGE.GETCOLORR expects handle, x, y")
		}
		img, err := m.getImage(args, 0, "IMAGE.GETCOLORR")
		if err != nil {
			return value.Nil, err
		}
		x, ok1 := argInt(args[1])
		y, ok2 := argInt(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("IMAGE.GETCOLORR: x, y must be numeric")
		}
		c := rl.GetImageColor(*img, x, y)
		return value.FromInt(int64(c.R)), nil
	}))

	reg.Register("IMAGE.GETCOLORG", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("IMAGE.GETCOLORG expects handle, x, y")
		}
		img, err := m.getImage(args, 0, "IMAGE.GETCOLORG")
		if err != nil {
			return value.Nil, err
		}
		x, ok1 := argInt(args[1])
		y, ok2 := argInt(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("IMAGE.GETCOLORG: x, y must be numeric")
		}
		c := rl.GetImageColor(*img, x, y)
		return value.FromInt(int64(c.G)), nil
	}))

	reg.Register("IMAGE.GETCOLORB", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("IMAGE.GETCOLORB expects handle, x, y")
		}
		img, err := m.getImage(args, 0, "IMAGE.GETCOLORB")
		if err != nil {
			return value.Nil, err
		}
		x, ok1 := argInt(args[1])
		y, ok2 := argInt(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("IMAGE.GETCOLORB: x, y must be numeric")
		}
		c := rl.GetImageColor(*img, x, y)
		return value.FromInt(int64(c.B)), nil
	}))

	reg.Register("IMAGE.GETCOLORA", "image", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 3 {
			return value.Nil, fmt.Errorf("IMAGE.GETCOLORA expects handle, x, y")
		}
		img, err := m.getImage(args, 0, "IMAGE.GETCOLORA")
		if err != nil {
			return value.Nil, err
		}
		x, ok1 := argInt(args[1])
		y, ok2 := argInt(args[2])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("IMAGE.GETCOLORA: x, y must be numeric")
		}
		c := rl.GetImageColor(*img, x, y)
		return value.FromInt(int64(c.A)), nil
	}))

	bboxComponent := func(which int) func([]value.Value) (value.Value, error) {
		return func(args []value.Value) (value.Value, error) {
			if err := m.requireHeap(); err != nil {
				return value.Nil, err
			}
			if len(args) != 2 {
				return value.Nil, fmt.Errorf("IMAGE.GETBBOX*: expects handle, alphaThreshold")
			}
			img, err := m.getImage(args, 0, "IMAGE.GETBBOX")
			if err != nil {
				return value.Nil, err
			}
			th, ok := argFloat(args[1])
			if !ok {
				return value.Nil, fmt.Errorf("IMAGE.GETBBOX*: threshold must be numeric")
			}
			r := imageAlphaBBox(img, th)
			switch which {
			case 0:
				return value.FromInt(int64(r.X)), nil
			case 1:
				return value.FromInt(int64(r.Y)), nil
			case 2:
				return value.FromInt(int64(r.Width)), nil
			default:
				return value.FromInt(int64(r.Height)), nil
			}
		}
	}

	reg.Register("IMAGE.GETBBOXX", "image", runtime.AdaptLegacy(bboxComponent(0)))
	reg.Register("IMAGE.GETBBOXY", "image", runtime.AdaptLegacy(bboxComponent(1)))
	reg.Register("IMAGE.GETBBOXW", "image", runtime.AdaptLegacy(bboxComponent(2)))
	reg.Register("IMAGE.GETBBOXH", "image", runtime.AdaptLegacy(bboxComponent(3)))
}
