//go:build !linux || !cgo

package mbcharcontroller

import (
	"math"
	"moonbasic/runtime"
	mbphysics3d "moonbasic/runtime/physics3d"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

type Vec3 struct {
	X, Y, Z float32
}

type charObj struct {
	pos      Vec3
	vel      Vec3
	grounded bool
	radius   float32
	height   float32
	stepH    float32
	snapD    float32
	gravityG float32
	maxSlope float32
	friction float32
	release  heap.ReleaseOnce
}

func (c *charObj) TypeName() string { return "CharController" }
func (c *charObj) TypeTag() uint16  { return heap.TagCharController }
func (c *charObj) Free()            {}

func registerCharControllerCommands(m *Module, reg runtime.Registrar) {
	// Redundant CHARACTER.* and CHARACTERREF.* commands are now handled 
	// by the high-level player module to ensure API consistency and 
	// support for both standalone and entity-bound characters.
}

func hGet(m *Module, v value.Value) (*charObj, error) {
	return heap.Cast[*charObj](m.h, heap.Handle(v.IVal))
}

func (c *charObj) hostUpdate(dt float32) {
	if c.grounded {
		c.vel.Y = -0.1
	} else {
		c.vel.Y -= 32.0 * c.gravityG * dt
	}

	// 1. Horizontal Phase (Iterative Slide)
	moveX := c.vel.X * dt
	moveZ := c.vel.Z * dt
	
	c.pos.X += moveX
	c.pos.Z += moveZ

	statics := mbphysics3d.GetStaticBodyRegistry()
	
	// Iterative Slide (Simplified for AABB-stubs)
	for i := 0; i < 3; i++ {
		hit := false
		var nx, nz float32
		
		for _, b := range statics {
			if b.Shape == nil || b.Shape.Kind != 1 { continue }
			hx, hy, hz := b.Shape.F1, b.Shape.F2, b.Shape.F3
			
			// Y-range check for body
			if c.pos.Y+c.height*0.5 < b.Pos.Y-hy || c.pos.Y-c.height*0.5 > b.Pos.Y+hy {
				continue
			}

			// AABB overlap check
			if c.pos.X > b.Pos.X-hx-c.radius && c.pos.X < b.Pos.X+hx+c.radius &&
			   c.pos.Z > b.Pos.Z-hz-c.radius && c.pos.Z < b.Pos.Z+hz+c.radius {
				
				hit = true
				// Calculate push-out normal
				dx := c.pos.X - b.Pos.X
				dz := c.pos.Z - b.Pos.Z
				
				if math.Abs(float64(dx))/(float64(hx)+float64(c.radius)) > math.Abs(float64(dz))/(float64(hz)+float64(c.radius)) {
					if dx > 0 { nx = 1; c.pos.X = b.Pos.X + hx + c.radius + 0.001 } else { nx = -1; c.pos.X = b.Pos.X - hx - c.radius - 0.001 }
				} else {
					if dz > 0 { nz = 1; c.pos.Z = b.Pos.Z + hz + c.radius + 0.001 } else { nz = -1; c.pos.Z = b.Pos.Z - hz - c.radius - 0.001 }
				}
				break
			}
		}
		if !hit { break }
		
		// Project velocity onto plane
		dot := c.vel.X*nx + c.vel.Z*nz
		if dot < 0 {
			c.vel.X -= nx * dot
			c.vel.Z -= nz * dot
		}
	}

	// 2. Vertical Phase
	c.pos.Y += c.vel.Y * dt

	// 3. Ground Snapping
	c.grounded = false
	feetY := c.pos.Y - c.height*0.5
	
	for _, b := range statics {
		if b.Shape == nil || b.Shape.Kind != 1 { continue }
		hx, hy, hz := b.Shape.F1, b.Shape.F2, b.Shape.F3
		if c.pos.X > b.Pos.X-hx-c.radius && c.pos.X < b.Pos.X+hx+c.radius &&
		   c.pos.Z > b.Pos.Z-hz-c.radius && c.pos.Z < b.Pos.Z+hz+c.radius {
			
			topY := b.Pos.Y + hy
			if feetY <= topY + c.snapD && feetY >= topY - 0.5 {
				c.pos.Y = topY + c.height*0.5
				c.vel.Y = 0
				c.grounded = true
				break
			}
		}
	}
}

func shutdownCharController(m *Module) { _ = m }
