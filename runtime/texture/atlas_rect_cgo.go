//go:build cgo || (windows && !cgo)

package texture

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

// FrameSourceRect returns the source rectangle for DrawTexture / DrawBillboardRec (atlas cell + scroll offsets).
func (t *TextureObject) FrameSourceRect() rl.Rectangle {
	t.mu.RLock()
	defer t.mu.RUnlock()
	w := float32(t.Tex.Width)
	h := float32(t.Tex.Height)
	if t.AtlasCols <= 0 || t.AtlasRows <= 0 {
		ww := w * t.UScl
		if ww <= 0 {
			ww = w
		}
		hh := h * t.VScl
		if hh <= 0 {
			hh = h
		}
		return rl.Rectangle{
			X:      t.UPos + t.ScrollAccumU,
			Y:      t.VPos + t.ScrollAccumV,
			Width:  ww,
			Height: hh,
		}
	}
	cw := w / float32(t.AtlasCols)
	ch := h / float32(t.AtlasRows)
	n := int(t.AtlasCols * t.AtlasRows)
	fi := int(t.FrameIndex)
	if n > 0 {
		fi = ((fi % n) + n) % n
	}
	fx := fi % int(t.AtlasCols)
	fy := fi / int(t.AtlasCols)
	return rl.Rectangle{
		X:      float32(fx)*cw + t.UPos + t.ScrollAccumU,
		Y:      float32(fy)*ch + t.VPos + t.ScrollAccumV,
		Width:  cw,
		Height: ch,
	}
}

// Tick advances automatic frame playback and UV scroll; call once per frame via TEXTURE.TICKALL.
func (t *TextureObject) Tick(dt float32) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.AnimPlaying && t.AtlasCols > 0 && t.AtlasRows > 0 && t.AnimFPS > 0 {
		t.animTime += dt
		step := float32(1.0 / float64(t.AnimFPS))
		if step <= 0 {
			step = 1.0 / 30.0
		}
		n := int(t.AtlasCols * t.AtlasRows)
		for t.animTime >= step && n > 0 {
			t.animTime -= step
			t.FrameIndex++
			if int(t.FrameIndex) >= n {
				if t.AnimLoop {
					t.FrameIndex = 0
				} else {
					t.FrameIndex = int32(n - 1)
					t.AnimPlaying = false
					break
				}
			}
		}
	}
	t.ScrollAccumU += t.ScrollSpeedU * dt
	t.ScrollAccumV += t.ScrollSpeedV * dt
}
