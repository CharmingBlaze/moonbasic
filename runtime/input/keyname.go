//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"

	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) inGetKeyName(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.GETKEYNAME expects (key)")
	}
	k, err := KeyCodeFromValue(args[0])
	if err != nil {
		return value.Nil, fmt.Errorf("INPUT.GETKEYNAME: %w", err)
	}
	idx := m.requireHeap().Intern(keyboardKeyName(k))
	return value.FromStringIndex(idx), nil
}

// keyboardKeyName returns an English label for a Raylib keyboard key (Raylib 5.x scancodes).
// This replaces GLFW's GetKeyName (CGO); layout-specific names are not available without the OS API.
func keyboardKeyName(k int32) string {
	if s, ok := raylibKeyNames[k]; ok {
		return s
	}
	if k >= 'A' && k <= 'Z' {
		return string(rune(k))
	}
	if k >= '0' && k <= '9' {
		return string(rune(k))
	}
	return ""
}

// raylibKeyNames maps non-alphanumeric Raylib keyboard constants to stable names.
var raylibKeyNames = map[int32]string{
	rl.KeyNull: "",

	rl.KeyApostrophe:   "'",
	rl.KeyComma:        ",",
	rl.KeyMinus:        "-",
	rl.KeyPeriod:       ".",
	rl.KeySlash:        "/",
	rl.KeySemicolon:    ";",
	rl.KeyEqual:        "=",
	rl.KeyLeftBracket:  "[",
	rl.KeyBackSlash:    "\\",
	rl.KeyRightBracket: "]",
	rl.KeyGrave:        "`",

	rl.KeySpace:       "SPACE",
	rl.KeyEscape:      "ESCAPE",
	rl.KeyEnter:       "ENTER",
	rl.KeyTab:         "TAB",
	rl.KeyBackspace:   "BACKSPACE",
	rl.KeyInsert:      "INSERT",
	rl.KeyDelete:      "DELETE",
	rl.KeyRight:       "RIGHT",
	rl.KeyLeft:        "LEFT",
	rl.KeyDown:        "DOWN",
	rl.KeyUp:          "UP",
	rl.KeyPageUp:      "PAGE UP",
	rl.KeyPageDown:    "PAGE DOWN",
	rl.KeyHome:        "HOME",
	rl.KeyEnd:         "END",
	rl.KeyCapsLock:    "CAPS LOCK",
	rl.KeyScrollLock:  "SCROLL LOCK",
	rl.KeyNumLock:     "NUM LOCK",
	rl.KeyPrintScreen: "PRINT SCREEN",
	rl.KeyPause:       "PAUSE",

	rl.KeyF1: "F1", rl.KeyF2: "F2", rl.KeyF3: "F3", rl.KeyF4: "F4",
	rl.KeyF5: "F5", rl.KeyF6: "F6", rl.KeyF7: "F7", rl.KeyF8: "F8",
	rl.KeyF9: "F9", rl.KeyF10: "F10", rl.KeyF11: "F11", rl.KeyF12: "F12",

	rl.KeyLeftShift:    "LEFT SHIFT",
	rl.KeyLeftControl:  "LEFT CONTROL",
	rl.KeyLeftAlt:      "LEFT ALT",
	rl.KeyLeftSuper:    "LEFT SUPER",
	rl.KeyRightShift:   "RIGHT SHIFT",
	rl.KeyRightControl: "RIGHT CONTROL",
	rl.KeyRightAlt:     "RIGHT ALT",
	rl.KeyRightSuper:   "RIGHT SUPER",
	rl.KeyKbMenu:       "MENU",

	rl.KeyKp0: "KEYPAD 0", rl.KeyKp1: "KEYPAD 1", rl.KeyKp2: "KEYPAD 2",
	rl.KeyKp3: "KEYPAD 3", rl.KeyKp4: "KEYPAD 4", rl.KeyKp5: "KEYPAD 5",
	rl.KeyKp6: "KEYPAD 6", rl.KeyKp7: "KEYPAD 7", rl.KeyKp8: "KEYPAD 8",
	rl.KeyKp9:        "KEYPAD 9",
	rl.KeyKpDecimal:  "KEYPAD .",
	rl.KeyKpDivide:   "KEYPAD /",
	rl.KeyKpMultiply: "KEYPAD *",
	rl.KeyKpSubtract: "KEYPAD -",
	rl.KeyKpAdd:      "KEYPAD +",
	rl.KeyKpEnter:    "KEYPAD ENTER",
	rl.KeyKpEqual:    "KEYPAD =",

	rl.KeyBack:       "BACK",
	rl.KeyMenu:       "MENU",
	rl.KeyVolumeUp:   "VOLUME UP",
	rl.KeyVolumeDown: "VOLUME DOWN",
}
