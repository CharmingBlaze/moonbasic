package mbgame

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/runtime/mbmatrix"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerCollisionVecBuiltins(r runtime.Registrar) {
	r.Register("COLLISION.BOXOVERLAP2D", "game", runtime.AdaptLegacy(m.collisionBoxOverlap2D))
	r.Register("COLLISION.CIRCLEOVERLAP2D", "game", runtime.AdaptLegacy(m.collisionCircleOverlap2D))
	r.Register("COLLISION.POINTINBOX2D", "game", runtime.AdaptLegacy(m.collisionPointInBox2D))
	r.Register("COLLISION.CIRCLEBOX2D", "game", runtime.AdaptLegacy(m.collisionCircleBox2D))
	r.Register("COLLISION.LINESEGINTERSECT2D", "game", runtime.AdaptLegacy(m.collisionLineSegIntersect2D))
	r.Register("COLLISION.POINTONSEG2D", "game", runtime.AdaptLegacy(m.collisionPointOnSeg2D))
	r.Register("COLLISION.SPHEREOVERLAP3D", "game", runtime.AdaptLegacy(m.collisionSphereOverlap3D))
	r.Register("COLLISION.AABBOVERLAP3D", "game", runtime.AdaptLegacy(m.collisionAABBOverlap3D))
	r.Register("COLLISION.SPHEREBOX3D", "game", runtime.AdaptLegacy(m.collisionSphereBox3D))
	r.Register("COLLISION.POINTINAABB3D", "game", runtime.AdaptLegacy(m.collisionPointInAABB3D))
}

func (m *Module) vec2Arg(args []value.Value, ix int, op string) (heap.Handle, error) {
	if err := m.requireHeap(op); err != nil {
		return 0, err
	}
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return 0, fmt.Errorf("%s: argument %d must be vec2 handle", op, ix+1)
	}
	return heap.Handle(args[ix].IVal), nil
}

func (m *Module) vec3Arg(args []value.Value, ix int, op string) (heap.Handle, error) {
	if err := m.requireHeap(op); err != nil {
		return 0, err
	}
	if ix >= len(args) || args[ix].Kind != value.KindHandle {
		return 0, fmt.Errorf("%s: argument %d must be vec3 handle", op, ix+1)
	}
	return heap.Handle(args[ix].IVal), nil
}

func (m *Module) collisionBoxOverlap2D(args []value.Value) (value.Value, error) {
	const op = "COLLISION.BOXOVERLAP2D"
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("%s expects 4 arguments (posA, sizeA, posB, sizeB as VEC2 handles)", op)
	}
	ha, err := m.vec2Arg(args, 0, op)
	if err != nil {
		return value.Nil, err
	}
	hb, err := m.vec2Arg(args, 1, op)
	if err != nil {
		return value.Nil, err
	}
	hc, err := m.vec2Arg(args, 2, op)
	if err != nil {
		return value.Nil, err
	}
	hd, err := m.vec2Arg(args, 3, op)
	if err != nil {
		return value.Nil, err
	}
	pa, err := mbmatrix.Vec2FromHandle(m.h, ha)
	if err != nil {
		return value.Nil, err
	}
	sa, err := mbmatrix.Vec2FromHandle(m.h, hb)
	if err != nil {
		return value.Nil, err
	}
	pb, err := mbmatrix.Vec2FromHandle(m.h, hc)
	if err != nil {
		return value.Nil, err
	}
	sb, err := mbmatrix.Vec2FromHandle(m.h, hd)
	if err != nil {
		return value.Nil, err
	}
	ok := boxCollide2D(float64(pa.X), float64(pa.Y), float64(sa.X), float64(sa.Y),
		float64(pb.X), float64(pb.Y), float64(sb.X), float64(sb.Y))
	return value.FromBool(ok), nil
}

func (m *Module) collisionCircleOverlap2D(args []value.Value) (value.Value, error) {
	const op = "COLLISION.CIRCLEOVERLAP2D"
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("%s expects 4 arguments (center1, radius1, center2, radius2)", op)
	}
	h1, err := m.vec2Arg(args, 0, op)
	if err != nil {
		return value.Nil, err
	}
	r1, ok1 := argF(args[1])
	h2, err := m.vec2Arg(args, 2, op)
	if err != nil {
		return value.Nil, err
	}
	r2, ok2 := argF(args[3])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("%s: radii must be numeric", op)
	}
	c1, err := mbmatrix.Vec2FromHandle(m.h, h1)
	if err != nil {
		return value.Nil, err
	}
	c2, err := mbmatrix.Vec2FromHandle(m.h, h2)
	if err != nil {
		return value.Nil, err
	}
	ok := circleCollide2D(float64(c1.X), float64(c1.Y), r1, float64(c2.X), float64(c2.Y), r2)
	return value.FromBool(ok), nil
}

func (m *Module) collisionPointInBox2D(args []value.Value) (value.Value, error) {
	const op = "COLLISION.POINTINBOX2D"
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("%s expects 3 arguments (point, boxPos, boxSize as VEC2 handles)", op)
	}
	hp, err := m.vec2Arg(args, 0, op)
	if err != nil {
		return value.Nil, err
	}
	hq, err := m.vec2Arg(args, 1, op)
	if err != nil {
		return value.Nil, err
	}
	hs, err := m.vec2Arg(args, 2, op)
	if err != nil {
		return value.Nil, err
	}
	pt, err := mbmatrix.Vec2FromHandle(m.h, hp)
	if err != nil {
		return value.Nil, err
	}
	bp, err := mbmatrix.Vec2FromHandle(m.h, hq)
	if err != nil {
		return value.Nil, err
	}
	bs, err := mbmatrix.Vec2FromHandle(m.h, hs)
	if err != nil {
		return value.Nil, err
	}
	ok := pointInBox2D(float64(pt.X), float64(pt.Y), float64(bp.X), float64(bp.Y), float64(bs.X), float64(bs.Y))
	return value.FromBool(ok), nil
}

func (m *Module) collisionCircleBox2D(args []value.Value) (value.Value, error) {
	const op = "COLLISION.CIRCLEBOX2D"
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("%s expects 4 arguments (circleCenter, radius, boxPos, boxSize)", op)
	}
	hc, err := m.vec2Arg(args, 0, op)
	if err != nil {
		return value.Nil, err
	}
	cr, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("%s: radius must be numeric", op)
	}
	hb, err := m.vec2Arg(args, 2, op)
	if err != nil {
		return value.Nil, err
	}
	hs, err := m.vec2Arg(args, 3, op)
	if err != nil {
		return value.Nil, err
	}
	c, err := mbmatrix.Vec2FromHandle(m.h, hc)
	if err != nil {
		return value.Nil, err
	}
	bp, err := mbmatrix.Vec2FromHandle(m.h, hb)
	if err != nil {
		return value.Nil, err
	}
	bz, err := mbmatrix.Vec2FromHandle(m.h, hs)
	if err != nil {
		return value.Nil, err
	}
	okb := circleBoxCollide2D(float64(c.X), float64(c.Y), cr, float64(bp.X), float64(bp.Y), float64(bz.X), float64(bz.Y))
	return value.FromBool(okb), nil
}

func (m *Module) collisionLineSegIntersect2D(args []value.Value) (value.Value, error) {
	const op = "COLLISION.LINESEGINTERSECT2D"
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("%s expects 4 arguments (a1, a2, b1, b2 as VEC2 segment endpoints)", op)
	}
	var hs [4]heap.Handle
	for i := 0; i < 4; i++ {
		h, err := m.vec2Arg(args, i, op)
		if err != nil {
			return value.Nil, err
		}
		hs[i] = h
	}
	var vs [4]struct{ x, y float64 }
	for i := 0; i < 4; i++ {
		v, err := mbmatrix.Vec2FromHandle(m.h, hs[i])
		if err != nil {
			return value.Nil, err
		}
		vs[i].x, vs[i].y = float64(v.X), float64(v.Y)
	}
	ok := lineCollide2D(vs[0].x, vs[0].y, vs[1].x, vs[1].y, vs[2].x, vs[2].y, vs[3].x, vs[3].y)
	return value.FromBool(ok), nil
}

func (m *Module) collisionPointOnSeg2D(args []value.Value) (value.Value, error) {
	const op = "COLLISION.POINTONSEG2D"
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("%s expects 4 arguments (point, segA, segB, threshold)", op)
	}
	hp, err := m.vec2Arg(args, 0, op)
	if err != nil {
		return value.Nil, err
	}
	ha, err := m.vec2Arg(args, 1, op)
	if err != nil {
		return value.Nil, err
	}
	hb, err := m.vec2Arg(args, 2, op)
	if err != nil {
		return value.Nil, err
	}
	th, ok := argF(args[3])
	if !ok {
		return value.Nil, fmt.Errorf("%s: threshold must be numeric", op)
	}
	pt, err := mbmatrix.Vec2FromHandle(m.h, hp)
	if err != nil {
		return value.Nil, err
	}
	a, err := mbmatrix.Vec2FromHandle(m.h, ha)
	if err != nil {
		return value.Nil, err
	}
	b, err := mbmatrix.Vec2FromHandle(m.h, hb)
	if err != nil {
		return value.Nil, err
	}
	okb := pointOnLine2D(float64(pt.X), float64(pt.Y), float64(a.X), float64(a.Y), float64(b.X), float64(b.Y), th)
	return value.FromBool(okb), nil
}

func (m *Module) collisionSphereOverlap3D(args []value.Value) (value.Value, error) {
	const op = "COLLISION.SPHEREOVERLAP3D"
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("%s expects 4 arguments (center1, radius1, center2, radius2)", op)
	}
	h1, err := m.vec3Arg(args, 0, op)
	if err != nil {
		return value.Nil, err
	}
	r1, ok1 := argF(args[1])
	h2, err := m.vec3Arg(args, 2, op)
	if err != nil {
		return value.Nil, err
	}
	r2, ok2 := argF(args[3])
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("%s: radii must be numeric", op)
	}
	c1, err := mbmatrix.Vec3FromHandle(m.h, h1)
	if err != nil {
		return value.Nil, err
	}
	c2, err := mbmatrix.Vec3FromHandle(m.h, h2)
	if err != nil {
		return value.Nil, err
	}
	ok := sphereCollide3D(float64(c1.X), float64(c1.Y), float64(c1.Z), r1,
		float64(c2.X), float64(c2.Y), float64(c2.Z), r2)
	return value.FromBool(ok), nil
}

func (m *Module) collisionAABBOverlap3D(args []value.Value) (value.Value, error) {
	const op = "COLLISION.AABBOVERLAP3D"
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("%s expects 4 arguments (minA, maxA, minB, maxB as VEC3 handles)", op)
	}
	var hs [4]heap.Handle
	for i := 0; i < 4; i++ {
		h, err := m.vec3Arg(args, i, op)
		if err != nil {
			return value.Nil, err
		}
		hs[i] = h
	}
	var vs [4]struct{ x, y, z float64 }
	for i := 0; i < 4; i++ {
		v, err := mbmatrix.Vec3FromHandle(m.h, hs[i])
		if err != nil {
			return value.Nil, err
		}
		vs[i].x, vs[i].y, vs[i].z = float64(v.X), float64(v.Y), float64(v.Z)
	}
	ok := aabbCollide3D(vs[0].x, vs[0].y, vs[0].z, vs[1].x, vs[1].y, vs[1].z,
		vs[2].x, vs[2].y, vs[2].z, vs[3].x, vs[3].y, vs[3].z)
	return value.FromBool(ok), nil
}

func (m *Module) collisionSphereBox3D(args []value.Value) (value.Value, error) {
	const op = "COLLISION.SPHEREBOX3D"
	if len(args) != 4 {
		return value.Nil, fmt.Errorf("%s expects 4 arguments (sphereCenter, radius, boxMinCorner, boxSize)", op)
	}
	hs, err := m.vec3Arg(args, 0, op)
	if err != nil {
		return value.Nil, err
	}
	sr, ok := argF(args[1])
	if !ok {
		return value.Nil, fmt.Errorf("%s: radius must be numeric", op)
	}
	hm, err := m.vec3Arg(args, 2, op)
	if err != nil {
		return value.Nil, err
	}
	hz, err := m.vec3Arg(args, 3, op)
	if err != nil {
		return value.Nil, err
	}
	c, err := mbmatrix.Vec3FromHandle(m.h, hs)
	if err != nil {
		return value.Nil, err
	}
	bmin, err := mbmatrix.Vec3FromHandle(m.h, hm)
	if err != nil {
		return value.Nil, err
	}
	bs, err := mbmatrix.Vec3FromHandle(m.h, hz)
	if err != nil {
		return value.Nil, err
	}
	okb := sphereBoxCollide3D(float64(c.X), float64(c.Y), float64(c.Z), sr,
		float64(bmin.X), float64(bmin.Y), float64(bmin.Z), float64(bs.X), float64(bs.Y), float64(bs.Z))
	return value.FromBool(okb), nil
}

func (m *Module) collisionPointInAABB3D(args []value.Value) (value.Value, error) {
	const op = "COLLISION.POINTINAABB3D"
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("%s expects 3 arguments (point, boxMinCorner, boxSize as VEC3 handles)", op)
	}
	hp, err := m.vec3Arg(args, 0, op)
	if err != nil {
		return value.Nil, err
	}
	hm, err := m.vec3Arg(args, 1, op)
	if err != nil {
		return value.Nil, err
	}
	hs, err := m.vec3Arg(args, 2, op)
	if err != nil {
		return value.Nil, err
	}
	pt, err := mbmatrix.Vec3FromHandle(m.h, hp)
	if err != nil {
		return value.Nil, err
	}
	bmin, err := mbmatrix.Vec3FromHandle(m.h, hm)
	if err != nil {
		return value.Nil, err
	}
	bs, err := mbmatrix.Vec3FromHandle(m.h, hs)
	if err != nil {
		return value.Nil, err
	}
	ok := pointInAABB3D(float64(pt.X), float64(pt.Y), float64(pt.Z),
		float64(bmin.X), float64(bmin.Y), float64(bmin.Z), float64(bs.X), float64(bs.Y), float64(bs.Z))
	return value.FromBool(ok), nil
}
