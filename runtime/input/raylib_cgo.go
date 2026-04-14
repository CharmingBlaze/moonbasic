//go:build cgo || (windows && !cgo)

package input

import (
	"fmt"

	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (m *Module) inKeyDown(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.KEYDOWN expects 1 argument")
	}
	kc, err := KeyCodeFromValue(args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsKeyDown(kc)), nil
}

func (m *Module) inKeyPressed(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.KEYPRESSED expects 1 argument")
	}
	kc, err := KeyCodeFromValue(args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsKeyPressed(kc)), nil
}

func (m *Module) inKeyReleased(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("INPUT.KEYUP expects 1 argument")
	}
	kc, err := KeyCodeFromValue(args[0])
	if err != nil {
		return value.Nil, err
	}
	return value.FromBool(rl.IsKeyReleased(kc)), nil
}

func (m *Module) Shutdown() {}

func (m *Module) inGetInactivity(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("INPUT.GETINACTIVITY expects 0 arguments")
	}
	
	// Heuristic for interaction
	interacted := false
	if rl.GetMouseDelta().X != 0 || rl.GetMouseDelta().Y != 0 { interacted = true }
	if rl.GetMouseWheelMove() != 0 { interacted = true }
	for b := int32(0); b < 3; b++ {
		if rl.IsMouseButtonDown(rl.MouseButton(b)) { interacted = true; break }
	}
	if !interacted {
		// Looping keys is slow, but we only need to find ONE
		for k := int32(32); k < 348; k++ {
			if rl.IsKeyDown(k) { interacted = true; break }
		}
	}

	now := rl.GetTime()
	if interacted || m.lastInteraction == 0 {
		m.lastInteraction = now
	}
	
	return value.FromFloat(now - m.lastInteraction), nil
}
