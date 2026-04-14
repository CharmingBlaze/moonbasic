package mbphysics3d

import (
	"fmt"
	"math"
	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type AeroState struct {
	LiftCoeff   float32
	ThrustPower float32
	DragCoeff   float32
}

var bodyAero = make(map[heap.Handle]*AeroState)

func registerAeroCommands(m *Module, reg runtime.Registrar) {
	reg.Register("AERO.SETLIFT", "physics3d", runtime.AdaptLegacy(m.ARSetLift))
	reg.Register("AERO.SETTHRUST", "physics3d", runtime.AdaptLegacy(m.ARSetThrust))
	reg.Register("AERO.SETDRAG", "physics3d", runtime.AdaptLegacy(m.ARSetDrag))
}

func (m *Module) ARSetLift(args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, fmt.Errorf("AERO.SETLIFT expects (body, coeff#)")
	}
	h := heap.Handle(args[0].IVal)
	f, _ := args[1].ToFloat()
	getAero(h).LiftCoeff = float32(f)
	return value.Nil, nil
}

func (m *Module) ARSetThrust(args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, fmt.Errorf("AERO.SETTHRUST expects (body, power#)")
	}
	h := heap.Handle(args[0].IVal)
	f, _ := args[1].ToFloat()
	getAero(h).ThrustPower = float32(f)
	return value.Nil, nil
}

func (m *Module) ARSetDrag(args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, fmt.Errorf("AERO.SETDRAG expects (body, coeff#)")
	}
	h := heap.Handle(args[0].IVal)
	f, _ := args[1].ToFloat()
	getAero(h).DragCoeff = float32(f)
	return value.Nil, nil
}

func getAero(h heap.Handle) *AeroState {
	s := bodyAero[h]
	if s == nil {
		s = &AeroState{}
		bodyAero[h] = s
	}
	return s
}

// ProcessAeroDynamics is called from the tick loop (bridged via mbentity or Step).
func (m *Module) ProcessAeroDynamics(dt float32) {
	for h, s := range bodyAero {
		b, err := heap.Cast[*body3dObj](m.h, h)
		if err != nil {
			delete(bodyAero, h)
			continue
		}
		
		// 1. Get Transform and Velocity (Platform Abstraction)
		pos, rot, ok := m.getBodyTransform(b)
		if !ok { continue }
		vel := m.getBodyVelocity(b)
		
		// Forward Vector (Local Z)
		fwd := rl.Vector3RotateByQuaternion(rl.Vector3{X: 0, Y: 0, Z: 1}, rot)
		up := rl.Vector3RotateByQuaternion(rl.Vector3{X: 0, Y: 1, Z: 0}, rot)
		
		// 2. Thrust
		if s.ThrustPower > 0 {
			thrust := rl.Vector3Scale(fwd, s.ThrustPower)
			m.applyBodyForce(b, thrust)
		}
		
		// 3. Lift (Arcade model: Upward force based on speed and wing alignment)
		speedSq := rl.Vector3LengthSqr(vel)
		if speedSq > 0.1 && s.LiftCoeff > 0 {
			// Angle of Attack (how much we are pointed into the wind)
			normVel := rl.Vector3Normalize(vel)
			aoa := rl.Vector3DotProduct(fwd, normVel)
			if aoa > 0 {
				liftAmount := s.LiftCoeff * speedSq * aoa
				lift := rl.Vector3Scale(up, liftAmount)
				m.applyBodyForce(b, lift)
			}
		}
		
		// 4. Drag (Air resistance)
		if s.DragCoeff > 0 && speedSq > 0.1 {
			drag := rl.Vector3Scale(vel, -s.DragCoeff * float32(math.Sqrt(float64(speedSq))))
			m.applyBodyForce(b, drag)
		}

		_ = pos // Pos unused here but retrieved for full state helper completeness
	}
}

// getBodyTransform / getBodyVelocity / applyBodyForce are implemented in jolt_world_cgo.go ((linux||windows)+cgo)
// and stub.go ((!linux&&!windows)||!cgo). Do not add duplicate methods here.
