package mbphysics3d

import (
	"fmt"
	"math"
	rl "github.com/gen2brain/raylib-go/raylib"
	"moonbasic/vm/value"
)

type Wheel struct {
	Offset      rl.Vector3
	Radius      float32
	SuspensionH float32
	Friction    float32
	
	// State
	Compress    float32
	IsGrounded  bool
	GroundPos   rl.Vector3
	GroundNorm  rl.Vector3
	WorldPos    rl.Vector3
}

type Vehicle struct {
	EntityID int64
	Wheels   []Wheel
	
	// Input State
	Throttle float32 // -1 to 1
	Steering float32 // -1 to 1
	Brake    float32 // 0 to 1
	
	// Tuning
	SpringStrength float32
	SpringDamping  float32
	MaxSpeed       float32
	SteerSpeed     float32
	
	// Dynamic State
	Velocity rl.Vector3
	AngVel   rl.Vector3
}

func (m *Module) VHCreate(args []value.Value) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, fmt.Errorf("VEHICLE.CREATE(entity, wheelCount)")
	}
	eid, _ := args[0].ToInt()
	wCount, _ := args[1].ToInt()
	
	v := &Vehicle{
		EntityID:       eid,
		Wheels:         make([]Wheel, wCount),
		SpringStrength: 15.0,
		SpringDamping:  0.8,
		MaxSpeed:       45.0,
		SteerSpeed:     1.5,
	}
	m.vehicles[eid] = v
	return value.FromInt(eid), nil
}

func (m *Module) VHSetTuning(args []value.Value) (value.Value, error) {
	if len(args) < 5 {
		return value.Nil, fmt.Errorf("VEHICLE.SETTUNING(v, spring#, damp#, maxSpeed#, steerSpeed#)")
	}
	vid, _ := args[0].ToInt()
	v := m.vehicles[vid]
	if v == nil { return value.Nil, nil }
	
	s, _ := args[1].ToFloat()
	d, _ := args[2].ToFloat()
	ms, _ := args[3].ToFloat()
	ss, _ := args[4].ToFloat()
	
	v.SpringStrength = float32(s)
	v.SpringDamping = float32(d)
	v.MaxSpeed = float32(ms)
	v.SteerSpeed = float32(ss)
	return value.Nil, nil
}

func (m *Module) VHWheelAxis(args []value.Value, axis int) (value.Value, error) {
	if len(args) < 2 {
		return value.Nil, fmt.Errorf("VEHICLE.WHEELX/Y/Z(v, idx)")
	}
	vid, _ := args[0].ToInt()
	idx, _ := args[1].ToInt()
	v := m.vehicles[vid]
	if v == nil || idx < 0 || idx >= int64(len(v.Wheels)) {
		return value.FromFloat(0), nil
	}
	
	w := v.Wheels[idx]
	switch axis {
	case 0: return value.FromFloat(float64(w.WorldPos.X)), nil
	case 1: return value.FromFloat(float64(w.WorldPos.Y)), nil
	case 2: return value.FromFloat(float64(w.WorldPos.Z)), nil
	}
	return value.FromFloat(0), nil
}

func (m *Module) VHSetWheel(args []value.Value) (value.Value, error) {
	if len(args) < 6 {
		return value.Nil, fmt.Errorf("VEHICLE.SETWHEEL(vehicle, idx, ox#, oy#, oz#, radius#)")
	}
	vid, _ := args[0].ToInt()
	idx, _ := args[1].ToInt()
	v := m.vehicles[vid]
	if v == nil || idx < 0 || idx >= int64(len(v.Wheels)) {
		return value.Nil, fmt.Errorf("VEHICLE.SETWHEEL: invalid vehicle or wheel index")
	}
	
	ox, _ := args[2].ToFloat()
	oy, _ := args[3].ToFloat()
	oz, _ := args[4].ToFloat()
	rad, _ := args[5].ToFloat()
	
	v.Wheels[idx] = Wheel{
		Offset:      rl.Vector3{X: float32(ox), Y: float32(oy), Z: float32(oz)},
		Radius:      float32(rad),
		SuspensionH: float32(rad) * 2.0,
		Friction:    0.8,
	}
	return value.Nil, nil
}

func (m *Module) VHControl(args []value.Value) (value.Value, error) {
	if len(args) < 4 {
		return value.Nil, fmt.Errorf("VEHICLE.CONTROL(vehicle, throttle#, steer#, brake#)")
	}
	vid, _ := args[0].ToInt()
	v := m.vehicles[vid]
	if v == nil { return value.Nil, nil }
	
	tf, _ := args[1].ToFloat()
	v.Throttle = float32(tf)
	sf, _ := args[2].ToFloat()
	v.Steering = float32(sf)
	bf, _ := args[3].ToFloat()
	v.Brake = float32(bf)
	return value.Nil, nil
}

func (m *Module) VHSetSteering(args []value.Value) (value.Value, error) {
	if len(args) < 2 { return value.Nil, fmt.Errorf("VEHICLE.SETSTEER(v, val#)") }
	vid, _ := args[0].ToInt()
	v := m.vehicles[vid]
	if v == nil { return value.Nil, nil }
	f, _ := args[1].ToFloat()
	v.Steering = float32(f)
	return value.Nil, nil
}

func (m *Module) VHSetThrottle(args []value.Value) (value.Value, error) {
	if len(args) < 2 { return value.Nil, fmt.Errorf("VEHICLE.SETTHROTTLE(v, val#)") }
	vid, _ := args[0].ToInt()
	v := m.vehicles[vid]
	if v == nil { return value.Nil, nil }
	f, _ := args[1].ToFloat()
	v.Throttle = float32(f)
	return value.Nil, nil
}

func (m *Module) VHStep(args []value.Value) (value.Value, error) {
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("VEHICLE.STEP(dt#)")
	}
	df, _ := args[0].ToFloat()
	dt := float32(df)
	if dt <= 0 { return value.Nil, nil }
	
	// Raycast provider (bridged from mbentity via Module)
	// For now, this is a simplified solver.
	for _, v := range m.vehicles {
		m.stepVehicle(v, dt)
	}
	return value.Nil, nil
}

func (m *Module) stepVehicle(v *Vehicle, dt float32) {
	if m.xformLookup == nil || m.xformUpdate == nil {
		return
	}
	pos, yaw, ok := m.xformLookup(v.EntityID)
	if !ok {
		return
	}

	// 1. Suspension & Grounding
	groundCount := 0
	rad := float32(math.Pi / 180.0)
	cosY := float32(math.Cos(float64(yaw * rad)))
	sinY := float32(math.Sin(float64(yaw * rad)))

	// Gravity
	v.Velocity.Y -= 28.0 * dt

	for i := range v.Wheels {
		w := &v.Wheels[i]
		// World wheel offset (rotate by yaw)
		wx := w.Offset.X*cosY + w.Offset.Z*sinY
		wz := w.Offset.Z*cosY - w.Offset.X*sinY
		worldWPos := rl.Vector3{X: pos.X + wx, Y: pos.Y, Z: pos.Z + wz}

		w.WorldPos = worldWPos

		// Raycast down
		_, _, _, hitY, ok := RaycastDownGroundProbe(float64(worldWPos.X), float64(worldWPos.Y), float64(worldWPos.Z), float64(w.SuspensionH))
		if ok {
			dist := worldWPos.Y - float32(hitY)
			if dist < w.SuspensionH {
				w.IsGrounded = true
				w.Compress = (w.SuspensionH - dist) / w.SuspensionH
				
				// Spring Force
				spring := w.Compress * v.SpringStrength
				// Damping (approximate)
				damp := (v.SpringDamping * v.Velocity.Y)
				
				v.Velocity.Y += (spring - damp) * dt
				groundCount++
				
				// Update world pos to snapped Y for mesh-tracking visual parity
				w.WorldPos.Y = float32(hitY) + w.Radius
			} else {
				w.IsGrounded = false
				w.Compress = 0
			}
		} else {
			w.IsGrounded = false
			w.Compress = 0
		}
	}

	// 2. Traction & Movement
	if groundCount > 0 {
		// Damping
		v.Velocity.X *= (1.0 - 0.2*dt)
		v.Velocity.Z *= (1.0 - 0.2*dt)

		// Steering (Yaw Change)
		if math.Abs(float64(v.Velocity.X)) > 0.1 || math.Abs(float64(v.Velocity.Z)) > 0.1 {
			yaw += v.Steering * 120.0 * dt
		}

		// Forward Drive
		fwd := rl.Vector3{X: sinY, Y: 0, Z: cosY}
		accel := v.Throttle * 60.0
		if v.Brake > 0.1 {
			accel = -v.Velocity.Z * v.Brake * 5.0 // Braking force
		}
		
		v.Velocity = rl.Vector3Add(v.Velocity, rl.Vector3Scale(fwd, accel*dt))
	}

	// 3. Integrate & Update Entity
	pos = rl.Vector3Add(pos, rl.Vector3Scale(v.Velocity, dt))
	m.xformUpdate(v.EntityID, pos)
	
	// Bonus: Tilt slightly based on velocity (arcade feel)
	// [Future expansion: actual pitch/roll from suspension]
}
