//go:build cgo

package mbgui

import (
	"fmt"
	"strings"

	"github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// rayguiOfficialStyleKeys lists bundled .rgs names (raysan5/raygui styles/).
var rayguiOfficialStyleKeys = []string{
	"amber", "ashes", "bluish", "candy", "cherry", "cyber", "dark",
	"enefete", "genesis", "jungle", "lavanda", "rltech", "sunny", "terminal",
}

var rayguiOfficialThemeUpper = map[string]struct{}{
	"AMBER": {}, "ASHES": {}, "BLUISH": {}, "CANDY": {}, "CHERRY": {}, "CYBER": {}, "DARK": {},
	"ENEFETE": {}, "GENESIS": {}, "JUNGLE": {}, "LAVANDA": {}, "RLTECH": {}, "SUNNY": {}, "TERMINAL": {},
}

func loadEmbeddedRayguiStyle(lowerName string) error {
	path := "raygui_styles/" + lowerName + ".rgs"
	data, err := embeddedRayguiStyles.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read embedded %s: %w", path, err)
	}
	if len(data) == 0 {
		return fmt.Errorf("embedded style %q is empty", lowerName)
	}
	raygui.LoadStyleFromMemory(data)
	return nil
}

func registerThemeCommands(m *Module, reg runtime.Registrar) {
	reg.Register("GUI.SETCOLOR", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 7 {
			return value.Nil, fmt.Errorf("GUI.SETCOLOR expects 7 args: control, property, r, g, b, a")
		}
		c, ok1 := argI32(args[0])
		p, ok2 := argI32(args[1])
		if !ok1 || !ok2 {
			return value.Nil, fmt.Errorf("GUI.SETCOLOR: control and property must be numeric")
		}
		col, err := colorArgs(args, 2)
		if err != nil {
			return value.Nil, err
		}
		raygui.SetStyle(raygui.ControlID(c), raygui.PropertyID(p), raygui.NewColorPropertyValue(col))
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTSIZE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTSIZE expects (size) — sets global gui font pixel height")
		}
		sz, ok := argI32(args[0])
		if !ok || sz < 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTSIZE: size must be a positive integer")
		}
		raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_SIZE, raygui.PropertyValue(sz))
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTSPACING", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTSPACING expects (spacing)")
		}
		s, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.SETTEXTSPACING: spacing must be numeric")
		}
		raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_SPACING, raygui.PropertyValue(s))
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTLINEHEIGHT", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTLINEHEIGHT expects (extraSpacing)")
		}
		s, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.SETTEXTLINEHEIGHT: value must be numeric")
		}
		raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_LINE_SPACING, raygui.PropertyValue(s))
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTWRAP", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTWRAP expects (mode) use GUI_TEXT_WRAP_*")
		}
		mo, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.SETTEXTWRAP: mode must be numeric")
		}
		raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_WRAP_MODE, raygui.PropertyValue(mo))
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTALIGN", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTALIGN expects (mode) use GUI_TEXT_ALIGN_LEFT/CENTER/RIGHT")
		}
		mo, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.SETTEXTALIGN: mode must be numeric")
		}
		raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_ALIGNMENT, raygui.PropertyValue(mo))
		return value.Nil, nil
	})
	reg.Register("GUI.SETTEXTALIGNVERT", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.SETTEXTALIGNVERT expects (mode) use GUI_TEXT_ALIGN_TOP/MIDDLE/BOTTOM")
		}
		mo, ok := argI32(args[0])
		if !ok {
			return value.Nil, fmt.Errorf("GUI.SETTEXTALIGNVERT: mode must be numeric")
		}
		raygui.SetStyle(raygui.DEFAULT, raygui.TEXT_ALIGNMENT_VERTICAL, raygui.PropertyValue(mo))
		return value.Nil, nil
	})
	reg.Register("GUI.GETTEXTSIZE", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.GETTEXTSIZE expects 0 arguments")
		}
		v := raygui.GetStyle(raygui.DEFAULT, raygui.TEXT_SIZE)
		return rt.RetInt(int64(v)), nil
	})
	reg.Register("GUI.THEMEAPPLY", "gui", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 {
			return value.Nil, fmt.Errorf("GUI.THEMEAPPLY expects (name) — see docs/reference/GUI.md (built-in + raygui styles)")
		}
		name, err := rt.ArgString(args, 0)
		if err != nil {
			return value.Nil, err
		}
		key := strings.ToUpper(strings.TrimSpace(name))
		switch key {
		case "DEFAULT", "RESET":
			raygui.LoadStyleDefault()
		case "LIGHT":
			applyGUILightTheme()
		case "BUILTIN_DARK":
			applyGUIDarkTheme()
		default:
			if _, ok := rayguiOfficialThemeUpper[key]; ok {
				low := strings.ToLower(key)
				if err := loadEmbeddedRayguiStyle(low); err != nil {
					return value.Nil, fmt.Errorf("GUI.THEMEAPPLY(%q): %w", name, err)
				}
				break
			}
			return value.Nil, fmt.Errorf("GUI.THEMEAPPLY: unknown theme %q — use DEFAULT, RESET, LIGHT, BUILTIN_DARK, or a raygui style: %s",
				name, strings.Join(rayguiOfficialStyleKeys, ", "))
		}
		return value.Nil, nil
	})
	guiThemeNames := func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, fmt.Errorf("GUI.THEMENAMES expects 0 arguments")
		}
		return rt.RetString(joinThemeNamesForDocs()), nil
	}
	reg.Register("GUI.THEMENAMES", "gui", guiThemeNames)
	reg.Register("GUI.THEMENAMES$", "gui", guiThemeNames)
}

// joinThemeNamesForDocs returns names accepted by GUI.THEMEAPPLY, semicolon-separated (for scripts / UI).
func joinThemeNamesForDocs() string {
	parts := []string{"DEFAULT", "RESET", "LIGHT", "BUILTIN_DARK"}
	for _, k := range rayguiOfficialStyleKeys {
		parts = append(parts, strings.ToUpper(k))
	}
	return strings.Join(parts, ";")
}

func applyGUIDarkTheme() {
	d := raygui.DEFAULT
	set := func(p raygui.PropertyID, r, g, b, a uint8) {
		raygui.SetStyle(d, p, raygui.NewColorPropertyValue(rl.Color{R: r, G: g, B: b, A: a}))
	}
	set(raygui.BORDER_COLOR_NORMAL, 72, 78, 88, 255)
	set(raygui.BASE_COLOR_NORMAL, 48, 54, 64, 255)
	set(raygui.TEXT_COLOR_NORMAL, 228, 232, 242, 255)
	set(raygui.BORDER_COLOR_FOCUSED, 110, 155, 235, 255)
	set(raygui.BASE_COLOR_FOCUSED, 58, 66, 80, 255)
	set(raygui.TEXT_COLOR_FOCUSED, 248, 250, 255, 255)
	set(raygui.BORDER_COLOR_PRESSED, 90, 130, 210, 255)
	set(raygui.BASE_COLOR_PRESSED, 70, 82, 100, 255)
	set(raygui.TEXT_COLOR_PRESSED, 255, 255, 255, 255)
	set(raygui.BORDER_COLOR_DISABLED, 48, 52, 58, 255)
	set(raygui.BASE_COLOR_DISABLED, 38, 42, 48, 255)
	set(raygui.TEXT_COLOR_DISABLED, 118, 124, 132, 255)
	raygui.SetStyle(d, raygui.BORDER_WIDTH, raygui.PropertyValue(1))
	raygui.SetStyle(d, raygui.TEXT_PADDING, raygui.PropertyValue(4))
	raygui.SetStyle(d, raygui.TEXT_SIZE, raygui.PropertyValue(12))
	raygui.SetStyle(d, raygui.TEXT_SPACING, raygui.PropertyValue(1))
	set(raygui.BACKGROUND_COLOR, 30, 34, 40, 255)
	set(raygui.LINE_COLOR, 88, 94, 104, 255)
}

func applyGUILightTheme() {
	d := raygui.DEFAULT
	set := func(p raygui.PropertyID, r, g, b, a uint8) {
		raygui.SetStyle(d, p, raygui.NewColorPropertyValue(rl.Color{R: r, G: g, B: b, A: a}))
	}
	set(raygui.BORDER_COLOR_NORMAL, 190, 195, 205, 255)
	set(raygui.BASE_COLOR_NORMAL, 248, 249, 252, 255)
	set(raygui.TEXT_COLOR_NORMAL, 32, 38, 52, 255)
	set(raygui.BORDER_COLOR_FOCUSED, 80, 130, 230, 255)
	set(raygui.BASE_COLOR_FOCUSED, 230, 238, 255, 255)
	set(raygui.TEXT_COLOR_FOCUSED, 20, 28, 48, 255)
	set(raygui.BORDER_COLOR_PRESSED, 60, 110, 200, 255)
	set(raygui.BASE_COLOR_PRESSED, 210, 225, 250, 255)
	set(raygui.TEXT_COLOR_PRESSED, 16, 24, 40, 255)
	set(raygui.BORDER_COLOR_DISABLED, 210, 212, 218, 255)
	set(raygui.BASE_COLOR_DISABLED, 235, 236, 240, 255)
	set(raygui.TEXT_COLOR_DISABLED, 150, 155, 165, 255)
	raygui.SetStyle(d, raygui.BORDER_WIDTH, raygui.PropertyValue(1))
	raygui.SetStyle(d, raygui.TEXT_PADDING, raygui.PropertyValue(4))
	raygui.SetStyle(d, raygui.TEXT_SIZE, raygui.PropertyValue(12))
	raygui.SetStyle(d, raygui.TEXT_SPACING, raygui.PropertyValue(1))
	set(raygui.BACKGROUND_COLOR, 236, 238, 244, 255)
	set(raygui.LINE_COLOR, 180, 186, 198, 255)
}
