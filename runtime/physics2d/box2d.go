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
	velIters int
	posIters int
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
	r.Register("PHYSICS2D.SETSTEP", "physics2d", m.phSetStep)
	r.Register("PHYSICS2D.SETITERATIONS", "physics2d", m.phSetIterations)
	r.Register("PHYSICS2D.STEP", "physics2d", m.phStep)

	r.Register("BODY2D.MAKE", "physics2d", m.bdMake)
	r.Register("BODY2D.ADDRECT", "physics2d", m.bdAddRect)
	r.Register("BODY2D.ADDCIRCLE", "physics2d", m.bdAddCircle)
	r.Register("BODY2D.COMMIT", "physics2d", m.bdCommit)
	r.Register("BODY2D.X", "physics2d", m.bdX)
	r.Register("BODY2D.Y", "physics2d", m.bdY)
	r.Register("BODY2D.ROT", "physics2d", m.bdRot)
	r.Register("BODY2D.FREE", "physics2d", m.bdFree)

	r.Register("BODY2D.SETPOS", "physics2d", m.bdSetPos)
	r.Register("BODY2D.GETPOS", "physics2d", m.bdGetPos)
	r.Register("BODY2D.SETROT", "physics2d", m.bdSetRot)
	r.Register("BODY2D.GETROT", "physics2d", m.bdGetRot)
	r.Register("BODY2D.SETMASS", "physics2d", m.bdSetMass)
	r.Register("BODY2D.SETFRICTION", "physics2d", m.bdSetFriction)
	r.Register("BODY2D.SETRESTITUTION", "physics2d", m.bdSetRestitution)
	r.Register("BODY2D.APPLYFORCE", "physics2d", m.bdApplyForce)
	r.Register("BODY2D.APPLYIMPULSE", "physics2d", m.bdApplyImpulse)

	// BOX2D aliases (legacy compatible names)
	r.Register("BOX2D.WORLDCREATE", "physics2d", m.phStart)
	r.Register("BOX2D.BODYCREATE", "physics2d", m.bdMake)
	r.Register("BOX2D.FIXTUREBOX", "physics2d", m.bdAddRect)
	r.Register("BOX2D.FIXTURECIRCLE", "physics2d", m.bdAddCircle)
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
		velIters: 8,
		posIters: 3,
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

func (m *Module) phSetStep(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if globalWorld == nil {
		return value.Nil, fmt.Errorf("PHYSICS2D not started")
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PHYSICS2D.SETSTEP expects (dt#)")
	}
	dt, err := rt.ArgFloat(args, 0)
	if err != nil {
		return value.Nil, err
	}
	if dt <= 0 || dt > 1.0 {
		return value.Nil, fmt.Errorf("PHYSICS2D.SETSTEP: dt must be in (0, 1] seconds")
	}
	globalWorld.stepRate = dt
	return value.Nil, nil
}

func (m *Module) phSetIterations(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if globalWorld == nil {
		return value.Nil, fmt.Errorf("PHYSICS2D not started")
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("PHYSICS2D.SETITERATIONS expects (velocityIters, positionIters)")
	}
	vi, err := rt.ArgInt(args, 0)
	if err != nil {
		return value.Nil, err
	}
	pi, err := rt.ArgInt(args, 1)
	if err != nil {
		return value.Nil, err
	}
	if vi < 1 || vi > 64 || pi < 1 || pi > 32 {
		return value.Nil, fmt.Errorf("PHYSICS2D.SETITERATIONS: velocity 1–64, position 1–32")
	}
	globalWorld.velIters = int(vi)
	globalWorld.posIters = int(pi)
	return value.Nil, nil
}

func (m *Module) phStep(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if globalWorld == nil {
		return value.Nil, nil
	}
	globalWorld.world.Step(globalWorld.stepRate, globalWorld.velIters, globalWorld.posIters)
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

	id, err := m.h.Alloc(&body2dObj{body: body})
	if err != nil {
		globalWorld.world.DestroyBody(body)
		return value.Nil, err
	}
	if err := m.h.Free(handle); err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) bdX(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.X")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(o.body.GetPosition().X), nil
}

func (m *Module) bdY(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.Y")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(o.body.GetPosition().Y), nil
}

func (m *Module) bdRot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.ROT")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(o.body.GetAngle()), nil
}

func (m *Module) bdFree(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("BODY2D.FREE expects handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

func (m *Module) getBody(args []value.Value, ix int, op string) (*body2dObj, error) {
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return nil, fmt.Errorf("%s: expected Body2D handle", op)
	}
	return heap.Cast[*body2dObj](m.h, heap.Handle(args[ix].IVal))
}

func (m *Module) bdSetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.SETPOS")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("BODY2D.SETPOS expects (handle, x, y)")
	}
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	o.body.SetTransform(box2d.MakeB2Vec2(x, y), o.body.GetAngle())
	return value.Nil, nil
}

func (m *Module) bdGetPos(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.GETPOS")
	if err != nil {
		return value.Nil, err
	}
	pos := o.body.GetPosition()
	p := heap.NewInstance("Point2D")
	p.SetField("x", value.FromFloat(pos.X))
	p.SetField("y", value.FromFloat(pos.Y))
	id, err := m.h.Alloc(p)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) bdSetRot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.SETROT")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("BODY2D.SETROT expects (handle, angle#)")
	}
	a, _ := args[1].ToFloat()
	o.body.SetTransform(o.body.GetPosition(), a)
	return value.Nil, nil
}

func (m *Module) bdGetRot(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.GETROT")
	if err != nil {
		return value.Nil, err
	}
	return value.FromFloat(o.body.GetAngle()), nil
}

func (m *Module) bdSetMass(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.SETMASS")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("BODY2D.SETMASS expects (handle, mass#)")
	}
	mval, _ := args[1].ToFloat()
	var data box2d.B2MassData
	o.body.GetMassData(&data)
	data.Mass = float64(mval)
	o.body.SetMassData(&data)
	return value.Nil, nil
}

func (m *Module) bdSetFriction(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.SETFRICTION")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("BODY2D.SETFRICTION expects (handle, friction#)")
	}
	fval, _ := args[1].ToFloat()
	for f := o.body.GetFixtureList(); f != nil; f = f.GetNext() {
		f.SetFriction(fval)
	}
	return value.Nil, nil
}

func (m *Module) bdSetRestitution(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.SETRESTITUTION")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("BODY2D.SETRESTITUTION expects (handle, bouncy#)")
	}
	rval, _ := args[1].ToFloat()
	for f := o.body.GetFixtureList(); f != nil; f = f.GetNext() {
		f.SetRestitution(rval)
	}
	return value.Nil, nil
}

func (m *Module) bdApplyForce(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.APPLYFORCE")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("BODY2D.APPLYFORCE expects (handle, fx#, fy#)")
	}
	fx, _ := args[1].ToFloat()
	fy, _ := args[2].ToFloat()
	o.body.ApplyForceToCenter(box2d.MakeB2Vec2(fx, fy), true)
	return value.Nil, nil
}

func (m *Module) bdApplyImpulse(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	o, err := m.getBody(args, 0, "BODY2D.APPLYIMPULSE")
	if err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("BODY2D.APPLYIMPULSE expects (handle, ix#, iy#)")
	}
	ix, _ := args[1].ToFloat()
	iy, _ := args[2].ToFloat()
	o.body.ApplyLinearImpulseToCenter(box2d.MakeB2Vec2(ix, iy), true)
	return value.Nil, nil
}
