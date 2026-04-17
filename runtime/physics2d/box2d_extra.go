package mbphysics2d

import (
	"fmt"
	"math"

	"github.com/ByteArena/box2d"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// Per-body collision state after the last PHYSICS2D.STEP (queried until the next STEP).
type coll2d struct {
	hit       bool
	other     heap.Handle
	nx, ny    float64
	px, py    float64
}

var (
	contact2d      map[heap.Handle]coll2d
	debug2dMode    int
	debug2dSegs    []float64
)

func clearPhysics2dAux() {
	contact2d = nil
	debug2dSegs = nil
	debug2dMode = 0
}

func bodyUserDataToHandle(b *box2d.B2Body) heap.Handle {
	if b == nil {
		return 0
	}
	u := b.GetUserData()
	if u == nil {
		return 0
	}
	switch v := u.(type) {
	case int64:
		return heap.Handle(v)
	case int:
		return heap.Handle(v)
	default:
		return 0
	}
}

func syncContactsAfterStep(m *Module) {
	if globalWorld == nil {
		return
	}
	contact2d = make(map[heap.Handle]coll2d)
	for c := globalWorld.world.GetContactList(); c != nil; c = c.GetNext() {
		if !c.IsTouching() {
			continue
		}
		wm := box2d.MakeB2WorldManifold()
		c.GetWorldManifold(&wm)
		ba := c.GetFixtureA().GetBody()
		bb := c.GetFixtureB().GetBody()
		ha := bodyUserDataToHandle(ba)
		hb := bodyUserDataToHandle(bb)
		if ha == 0 || hb == 0 {
			continue
		}
		nx, ny := wm.Normal.X, wm.Normal.Y
		px, py := 0.0, 0.0
		if c.GetManifold().PointCount > 0 {
			px = wm.Points[0].X
			py = wm.Points[0].Y
		}
		contact2d[ha] = coll2d{hit: true, other: hb, nx: nx, ny: ny, px: px, py: py}
		contact2d[hb] = coll2d{hit: true, other: ha, nx: -nx, ny: -ny, px: px, py: py}
	}
	if debug2dMode > 0 {
		collectDebugSegments()
	}
}

func collectDebugSegments() {
	debug2dSegs = debug2dSegs[:0]
	if globalWorld == nil {
		return
	}
	w := globalWorld.world
	for b := w.GetBodyList(); b != nil; b = b.GetNext() {
		xf := b.GetTransform()
		for f := b.GetFixtureList(); f != nil; f = f.GetNext() {
			sh := f.GetShape()
			switch sh.GetType() {
			case box2d.B2Shape_Type.E_polygon:
				poly, ok := sh.(*box2d.B2PolygonShape)
				if !ok {
					continue
				}
				n := poly.M_count
				if n < 2 {
					continue
				}
				for i := 0; i < n; i++ {
					i2 := (i + 1) % n
					v1 := box2d.B2TransformVec2Mul(xf, poly.M_vertices[i])
					v2 := box2d.B2TransformVec2Mul(xf, poly.M_vertices[i2])
					debug2dSegs = append(debug2dSegs, v1.X, v1.Y, v2.X, v2.Y)
				}
			case box2d.B2Shape_Type.E_circle:
				circ, ok := sh.(*box2d.B2CircleShape)
				if !ok {
					continue
				}
				center := b.GetWorldPoint(circ.M_p)
				r := circ.M_radius
				const seg = 16
				for i := 0; i < seg; i++ {
					a0 := float64(i) / float64(seg) * 2 * math.Pi
					a1 := float64(i+1) / float64(seg) * 2 * math.Pi
					x0 := center.X + r*math.Cos(a0)
					y0 := center.Y + r*math.Sin(a0)
					x1 := center.X + r*math.Cos(a1)
					y1 := center.Y + r*math.Sin(a1)
					debug2dSegs = append(debug2dSegs, x0, y0, x1, y1)
				}
			}
		}
	}
}

type joint2dObj struct {
	joint   box2d.B2JointInterface
	release heap.ReleaseOnce
}

func (o *joint2dObj) TypeName() string { return "Joint2D" }
func (o *joint2dObj) TypeTag() uint16  { return heap.TagJoint2D }
func (o *joint2dObj) Free() {
	o.release.Do(func() {
		if globalWorld != nil && o.joint != nil {
			globalWorld.world.DestroyJoint(o.joint)
		}
		o.joint = nil
	})
}

func (m *Module) phDebugDraw(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("PHYSICS2D.DEBUGDRAW expects (mode)")
	}
	mode, _ := args[0].ToFloat()
	debug2dMode = int(mode)
	return value.Nil, nil
}

func (m *Module) phGetDebugSegments(args []value.Value) (value.Value, error) {
	if m.h == nil {
		return value.Nil, runtime.Errorf("PHYSICS2D.GETDEBUGSEGMENTS: heap not bound")
	}
	if len(args) != 0 {
		return value.Nil, fmt.Errorf("PHYSICS2D.GETDEBUGSEGMENTS expects 0 arguments")
	}
	n := int64(len(debug2dSegs))
	if n == 0 {
		arr, err := heap.NewArray([]int64{1})
		if err != nil {
			return value.Nil, err
		}
		_ = arr.Set([]int64{0}, 0)
		id, err := m.h.Alloc(arr)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	arr, err := heap.NewArray([]int64{n})
	if err != nil {
		return value.Nil, err
	}
	for i := int64(0); i < n; i++ {
		_ = arr.Set([]int64{i}, debug2dSegs[i])
	}
	id, err := m.h.Alloc(arr)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) bdCollided(args []value.Value) (value.Value, error) {
	if _, err := m.getBody(args, 0, "BODY2D.COLLIDED"); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("BODY2D.COLLIDED expects (body)")
	}
	h := heap.Handle(args[0].IVal)
	c, ok := contact2d[h]
	if !ok || !c.hit {
		return value.FromInt(0), nil
	}
	return value.FromInt(1), nil
}

func (m *Module) bdCollisionOther(args []value.Value) (value.Value, error) {
	if _, err := m.getBody(args, 0, "BODY2D.COLLISIONOTHER"); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("BODY2D.COLLISIONOTHER expects (body)")
	}
	h := heap.Handle(args[0].IVal)
	c, ok := contact2d[h]
	if !ok || !c.hit {
		return value.FromHandle(0), nil
	}
	return value.FromHandle(c.other), nil
}

func (m *Module) bdCollisionNormal(args []value.Value) (value.Value, error) {
	if _, err := m.getBody(args, 0, "BODY2D.COLLISIONNORMAL"); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("BODY2D.COLLISIONNORMAL expects (body)")
	}
	h := heap.Handle(args[0].IVal)
	c, ok := contact2d[h]
	if !ok || !c.hit {
		p := heap.NewInstance("Point2D")
		p.SetField("x", value.FromFloat(0))
		p.SetField("y", value.FromFloat(0))
		id, err := m.h.Alloc(p)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	p := heap.NewInstance("Point2D")
	p.SetField("x", value.FromFloat(c.nx))
	p.SetField("y", value.FromFloat(c.ny))
	id, err := m.h.Alloc(p)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) bdCollisionPoint(args []value.Value) (value.Value, error) {
	if _, err := m.getBody(args, 0, "BODY2D.COLLISIONPOINT"); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("BODY2D.COLLISIONPOINT expects (body)")
	}
	h := heap.Handle(args[0].IVal)
	c, ok := contact2d[h]
	if !ok || !c.hit {
		p := heap.NewInstance("Point2D")
		p.SetField("x", value.FromFloat(0))
		p.SetField("y", value.FromFloat(0))
		id, err := m.h.Alloc(p)
		if err != nil {
			return value.Nil, err
		}
		return value.FromHandle(id), nil
	}
	p := heap.NewInstance("Point2D")
	p.SetField("x", value.FromFloat(c.px))
	p.SetField("y", value.FromFloat(c.py))
	id, err := m.h.Alloc(p)
	if err != nil {
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) bodyFromArg(args []value.Value, ix int, op string) (*box2d.B2Body, error) {
	o, err := m.getBody(args, ix, op)
	if err != nil {
		return nil, err
	}
	return o.body, nil
}

func (m *Module) jtDistance(args []value.Value) (value.Value, error) {
	if globalWorld == nil {
		return value.Nil, fmt.Errorf("JOINT2D.DISTANCE: PHYSICS2D not started")
	}
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("JOINT2D.DISTANCE expects (bodyA, bodyB, ax#, ay#, bx#, by#)")
	}
	ba, err := m.bodyFromArg(args, 0, "JOINT2D.DISTANCE")
	if err != nil {
		return value.Nil, err
	}
	bb, err := m.bodyFromArg(args, 1, "JOINT2D.DISTANCE")
	if err != nil {
		return value.Nil, err
	}
	ax, _ := args[2].ToFloat()
	ay, _ := args[3].ToFloat()
	bx, _ := args[4].ToFloat()
	by, _ := args[5].ToFloat()
	def := box2d.MakeB2DistanceJointDef()
	def.Initialize(ba, bb, box2d.MakeB2Vec2(ax, ay), box2d.MakeB2Vec2(bx, by))
	j := globalWorld.world.CreateJoint(&def)
	if j == nil {
		return value.Nil, fmt.Errorf("JOINT2D.DISTANCE: CreateJoint failed")
	}
	id, err := m.h.Alloc(&joint2dObj{joint: j})
	if err != nil {
		globalWorld.world.DestroyJoint(j)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) jtRevolute(args []value.Value) (value.Value, error) {
	if globalWorld == nil {
		return value.Nil, fmt.Errorf("JOINT2D.REVOLUTE: PHYSICS2D not started")
	}
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("JOINT2D.REVOLUTE expects (bodyA, bodyB, x#, y#)")
	}
	ba, err := m.bodyFromArg(args, 0, "JOINT2D.REVOLUTE")
	if err != nil {
		return value.Nil, err
	}
	bb, err := m.bodyFromArg(args, 1, "JOINT2D.REVOLUTE")
	if err != nil {
		return value.Nil, err
	}
	x, _ := args[2].ToFloat()
	y, _ := args[3].ToFloat()
	def := box2d.MakeB2RevoluteJointDef()
	def.Initialize(ba, bb, box2d.MakeB2Vec2(x, y))
	j := globalWorld.world.CreateJoint(&def)
	if j == nil {
		return value.Nil, fmt.Errorf("JOINT2D.REVOLUTE: CreateJoint failed")
	}
	id, err := m.h.Alloc(&joint2dObj{joint: j})
	if err != nil {
		globalWorld.world.DestroyJoint(j)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) jtPrismatic(args []value.Value) (value.Value, error) {
	if globalWorld == nil {
		return value.Nil, fmt.Errorf("JOINT2D.PRISMATIC: PHYSICS2D not started")
	}
	if len(args) != 6 {
		return value.Nil, fmt.Errorf("JOINT2D.PRISMATIC expects (bodyA, bodyB, x#, y#, ax#, ay#)")
	}
	ba, err := m.bodyFromArg(args, 0, "JOINT2D.PRISMATIC")
	if err != nil {
		return value.Nil, err
	}
	bb, err := m.bodyFromArg(args, 1, "JOINT2D.PRISMATIC")
	if err != nil {
		return value.Nil, err
	}
	x, _ := args[2].ToFloat()
	y, _ := args[3].ToFloat()
	axx, _ := args[4].ToFloat()
	axy, _ := args[5].ToFloat()
	def := box2d.MakeB2PrismaticJointDef()
	def.Initialize(ba, bb, box2d.MakeB2Vec2(x, y), box2d.MakeB2Vec2(axx, axy))
	j := globalWorld.world.CreateJoint(&def)
	if j == nil {
		return value.Nil, fmt.Errorf("JOINT2D.PRISMATIC: CreateJoint failed")
	}
	id, err := m.h.Alloc(&joint2dObj{joint: j})
	if err != nil {
		globalWorld.world.DestroyJoint(j)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}

func (m *Module) jtFree(args []value.Value) (value.Value, error) {
	if len(args) != 1 || args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("JOINT2D.FREE expects joint handle")
	}
	if err := m.h.Free(heap.Handle(args[0].IVal)); err != nil {
		return value.Nil, err
	}
	return value.Nil, nil
}

