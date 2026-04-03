package mbmatrix

import (
	"image/color"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// SeedColorGlobals allocates preset COLOR_* handles (Raylib-aligned RGBA).
func SeedColorGlobals(s *heap.Store, globals map[string]value.Value) {
	if s == nil || globals == nil {
		return
	}
	presets := map[string]color.RGBA{
		"COLOR_WHITE":       {255, 255, 255, 255},
		"COLOR_BLACK":       {0, 0, 0, 255},
		"COLOR_RED":         {230, 41, 55, 255},
		"COLOR_GREEN":       {0, 228, 48, 255},
		"COLOR_BLUE":        {0, 121, 241, 255},
		"COLOR_YELLOW":      {253, 249, 0, 255},
		"COLOR_MAGENTA":     {255, 0, 255, 255},
		"COLOR_CYAN":        {0, 255, 255, 255},
		"COLOR_ORANGE":      {255, 161, 0, 255},
		"COLOR_PINK":        {255, 109, 194, 255},
		"COLOR_GRAY":        {130, 130, 130, 255},
		"COLOR_DARKGRAY":    {80, 80, 80, 255},
		"COLOR_TRANSPARENT": {0, 0, 0, 0},
	}
	for name, c := range presets {
		id, err := s.Alloc(&colorObj{c: c})
		if err != nil {
			continue
		}
		globals[name] = value.FromHandle(id)
	}
}
