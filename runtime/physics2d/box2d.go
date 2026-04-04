package mbphysics2d

import (
	"fmt"

	"strings"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	"github.com/ByteArena/box2d"
)

type physics2dObj struct {
	world    *box2d.B2World
	stepRate float64
	gravityX float64
	gravityY float64
	release  heap.ReleaseOnce
}

func (o *physics2dObj) TypeName() string { return "Physics2D" }
func (o *physics2dObj) TypeTag() uint16  { return heap.TagPhysics2D }
func (o *physics2dObj) Free() {
	o.release.Do(func() { globalWorld = nil })
}

type body2dTemplate struct {
	kind     string // "static", "dynamic", "kinematic"
	fixtures []fixtureDef
}

type fixtureDef struct {
	isRect bool
	w, h   float64
	radius float64
}

func (o *body2dTemplate) TypeName() string { return "Body2DTemplate" }
func (o *body2dTemplate) TypeTag() uint16  { return heap.TagBody2D } // Re-using tag
func (o *body2dTemplate) Free()            {}

type body2dObj struct {
	body    *box2d.B2Body
	release heap.ReleaseOnce
}

func (o *body2dObj) TypeName() string { return "Body2D" }
func (o *body2dObj) TypeTag() uint16  { return heap.TagBody2D }
func (o *body2dObj) Free() {
	o.release.Do(func() {
		if globalWorld != nil && o.body != nil {
			globalWorld.world.DestroyBody(o.body)
		}
		o.body = nil
	})
}

func (m *Module) Register(r runtime.Registrar) {
	r.Register("PHYSICS2D.START", "physics2d", m.phStart)
	r.Register("PHYSICS2D.STOP", "physics2d", m.phStop)
	r.Register("PHYSICS2D.SETGRAVITY", "physics2d", m.phSetGravity)
	r.Register("PHYSICS2D.STEP", "physics2d", m.phStep)

	r.Register("BODY2D.MAKE", "physics2d", m.bdMake)
	r.Register("BODY2D.ADDRECT", "physics2d", m.bdAddRect)
	r.Register("BODY2D.ADDCIRCLE", "physics2d", m.bdAddCircle)
	r.Register("BODY2D.COMMIT", "physics2d", m.bdCommit)
	r.Register("BODY2D.X", "physics2d", m.bdX)
	r.Register("BODY2D.Y", "physics2d", m.bdY)
	r.Register("BODY2D.ROT", "physics2d", m.bdRot)
	r.Register("BODY2D.FREE", "physics2d", m.bdFree)
}

var globalWorld *physics2dObj

func (m *Module) phStart(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if globalWorld != nil {
		return value.Nil, nil
	}
	world := box2d.MakeB2World(box2d.MakeB2Vec2(0, -9.81))
	globalWorld = &physics2dObj{
		world:    &world,
		stepRate: 1.0 / 60.0,
		gravityX: 0,
		gravityY: -9.81,
	}
	return value.Nil, nil
}

func (m *Module) phStop(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	globalWorld = nil
	return value.Nil, nil
}

func (m *Module) phSetGravity(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if globalWorld == nil {
		return value.Nil, fmt.Errorf("PHYSICS2D not started")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PHYSICS2D.SETGRAVITY expects (gx, gy)")
	}
	gx, _ := args[0].ToFloat()
	gy, _ := args[1].ToFloat()
	globalWorld.gravityX = gx
	globalWorld.gravityY = gy
	globalWorld.world.SetGravity(box2d.MakeB2Vec2(gx, gy))
	return value.Nil, nil
}

func (m *Module) phStep(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if globalWorld == nil {
		return value.Nil, nil
	}
	globalWorld.world.Step(globalWorld.stepRate, 8, 3)
	return value.Nil, nil
}

func (m *Module) bdMake(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("BODY2D.MAKE expects (kind$)")
	}
	kind, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	kind = strings.ToLower(kind)
	id, err := m.h.Alloc(&body2dTemplate{kind: kind})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) bdAddRect(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("BODY2D.ADDRECT expects (handle, w, h)")
	}
	tmp, err := heap.Cast[*body2dTemplate](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	w, _ := args[1].ToFloat()
	h, _ := args[2].ToFloat()
	tmp.fixtures = append(tmp.fixtures, fixtureDef{isRect: true, w: w, h: h})
	return value.Nil, nil
}

func (m *Module) bdAddCircle(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("BODY2D.ADDCIRCLE expects (handle, radius)")
	}
	tmp, err := heap.Cast[*body2dTemplate](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	r, _ := args[1].ToFloat()
	tmp.fixtures = append(tmp.fixtures, fixtureDef{isRect: false, radius: r})
	return value.Nil, nil
}

func (m *Module) bdCommit(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if globalWorld == nil {
		return value.Nil, fmt.Errorf("PHYSICS2D not started")
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("BODY2D.COMMIT expects (handle, x, y)")
	}
	handle := heap.Handle(args[0].IVal)
	tmp, err := heap.Cast[*body2dTemplate](m.h, handle)
	if err != nil {
		return value.Nil, err
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()

	bd := box2d.MakeB2BodyDef()
	switch tmp.kind {
	case "static":
		bd.Type = box2d.B2BodyType.B2_staticBody
	case "kinematic":
		bd.Type = box2d.B2BodyType.B2_kinematicBody
	default:
		bd.Type = box2d.B2BodyType.B2_dynamicBody
	}
	bd.Position.Set(x, y)

	body := globalWorld.world.CreateBody(&bd)
	for _, f := range tmp.fixtures {
		if f.isRect {
			shape := box2d.MakeB2PolygonShape()
			shape.SetAsBox(f.w/2, f.h/2)
			body.CreateFixture(&shape, 1.0)
		} else {
			shape := box2d.MakeB2CircleShape()
			shape.M_radius = f.radius
			body.CreateFixture(&shape, 1.0)
		}
	}

	// Replace template with real body in place or new handle?
	// The commands.json says it returns a handle.
	id, err := m.h.Alloc(&body2dObj{body: body})
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) bdX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := heap.Cast[*body2dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(o.body.GetPosition().X), nil
}

func (m *Module) bdY(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := heap.Cast[*body2dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(o.body.GetPosition().Y), nil
}

func (m *Module) bdRot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := heap.Cast[*body2dObj](m.h, heap.Handle(args[0].IVal))
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(o.body.GetAngle()), nil
}

func (m *Module) bdFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}
