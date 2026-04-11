//go:build cgo || (windows && !cgo)

package mbimage

import (
	"fmt"
	"image"
	"image/draw"
	gifpkg "image/gif"
	"os"
	"path/filepath"
	"sort"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// ImageSequence holds CPU images for ENTITY.SETANIMATION / GIF / frame sequences.
type ImageSequence struct {
	Frames []*rl.Image
}

func (s *ImageSequence) TypeName() string { return "ImageSequence" }

// TypeTag implements heap.HeapObject.
func (s *ImageSequence) TypeTag() uint16 { return heap.TagImageSequence }

// Free implements heap.HeapObject.
func (s *ImageSequence) Free() {
	for _, im := range s.Frames {
		if im != nil {
			rl.UnloadImage(im)
		}
	}
	s.Frames = nil
}

func registerImageSequence(m *Module, reg runtime.Registrar) {
	reg.Register("IMAGE.LOADSEQUENCE", "image", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("IMAGE.LOADSEQUENCE expects pathPrefix (e.g. assets/water_)")
		}
		prefix, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		dir, base := filepath.Split(prefix)
		if base == "" {
			return value.Nil, fmt.Errorf("IMAGE.LOADSEQUENCE: need file prefix after directory")
		}
		pattern := filepath.Join(dir, base+"*")
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return value.Nil, err
		}
		var files []string
		for _, p := range matches {
			low := strings.ToLower(p)
			if strings.HasSuffix(low, ".png") || strings.HasSuffix(low, ".jpg") || strings.HasSuffix(low, ".jpeg") || strings.HasSuffix(low, ".bmp") {
				files = append(files, p)
			}
		}
		sort.Strings(files)
		if len(files) == 0 {
			return value.Nil, fmt.Errorf("IMAGE.LOADSEQUENCE: no images matched %q", pattern)
		}
		seq := &ImageSequence{Frames: make([]*rl.Image, 0, len(files))}
		for _, p := range files {
			im := rl.LoadImage(p)
			if im == nil || im.Data == nil {
				for _, x := range seq.Frames {
					rl.UnloadImage(x)
				}
				return value.Nil, fmt.Errorf("IMAGE.LOADSEQUENCE: failed to load %q", p)
			}
			seq.Frames = append(seq.Frames, im)
		}
		id, err := m.h.Alloc(seq)
		if err != nil {
			seq.Free()
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})

	reg.Register("IMAGE.LOADGIF", "image", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if err := m.requireHeap(); err != nil {
			return value.Nil, err
		}
		if len(args) != 1 || args[0].Kind != value.KindString {
			return value.Nil, fmt.Errorf("IMAGE.LOADGIF expects path")
		}
		path, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		seq, err := decodeGIFSequence(path)
		if err != nil {
			return value.Nil, err
		}
		id, err := m.h.Alloc(seq)
		if err != nil {
			seq.Free()
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	})
}

func decodeGIFSequence(path string) (*ImageSequence, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	g, err := gifpkg.DecodeAll(f)
	if err != nil {
		return nil, fmt.Errorf("IMAGE.LOADGIF: %w", err)
	}
	if len(g.Image) == 0 {
		return nil, fmt.Errorf("IMAGE.LOADGIF: no frames")
	}
	rect := image.Rect(0, 0, g.Config.Width, g.Config.Height)
	canvas := image.NewRGBA(rect)
	seq := &ImageSequence{Frames: make([]*rl.Image, 0, len(g.Image))}
	for _, paletted := range g.Image {
		draw.Draw(canvas, paletted.Bounds(), paletted, paletted.Bounds().Min, draw.Over)
		cp := image.NewRGBA(rect)
		draw.Draw(cp, rect, canvas, image.Point{}, draw.Src)
		im := stdRGBAToRaylib(cp)
		if im == nil {
			seq.Free()
			return nil, fmt.Errorf("IMAGE.LOADGIF: frame allocation failed")
		}
		seq.Frames = append(seq.Frames, im)
	}
	return seq, nil
}

func stdRGBAToRaylib(rgba *image.RGBA) *rl.Image {
	b := rgba.Bounds()
	w, h := b.Dx(), b.Dy()
	if w < 1 || h < 1 {
		return nil
	}
	im := rl.GenImageColor(int(w), int(h), rl.Blank)
	if im == nil || im.Data == nil {
		return nil
	}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r16, g16, b16, a16 := rgba.At(b.Min.X+x, b.Min.Y+y).RGBA()
			c := rl.Color{R: uint8(r16 >> 8), G: uint8(g16 >> 8), B: uint8(b16 >> 8), A: uint8(a16 >> 8)}
			rl.ImageDrawPixel(im, int32(x), int32(y), c)
		}
	}
	return im
}

// FrameCount returns the number of frames (for runtime bridges).
func (s *ImageSequence) FrameCount() int {
	if s == nil {
		return 0
	}
	return len(s.Frames)
}

// Frame returns frame i or nil.
func (s *ImageSequence) Frame(i int) *rl.Image {
	if s == nil || i < 0 || i >= len(s.Frames) {
		return nil
	}
	return s.Frames[i]
}
