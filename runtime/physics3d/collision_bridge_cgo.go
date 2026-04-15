//go:build (linux || windows) && cgo

package mbphysics3d

import (
	"fmt"
	"math"
	"sync"

	"github.com/bbitechnologies/jolt-go/jolt"

	"moonbasic/vm/heap"
)

// ContactHitData holds the last matched collision normal, contact point, and a force proxy for BASIC getters.
type ContactHitData struct {
	NX, NY, NZ float64
	PX, PY, PZ float64
	Force      float64
}

var (
	collisionMu sync.Mutex

	bufIdxToBodyHeap map[int]heap.Handle
	bufIdxToEntity   map[int]int64
	entityToBufIdx   map[int64]int
	handleToEntity   map[heap.Handle]int64

	// Frame state (drained at end of PHYSICS3D.STEP).
	collisionFrame []contactFrameEvent
	lastHit        ContactHitData
	lastHitValid   bool

	deadEntities map[int64]struct{}
)

type contactFrameEvent struct {
	E1, E2 int64
	Hit    ContactHitData
}

func resetCollisionBridgeState() {
	collisionMu.Lock()
	defer collisionMu.Unlock()
	bufIdxToBodyHeap = make(map[int]heap.Handle)
	bufIdxToEntity = make(map[int]int64)
	entityToBufIdx = make(map[int64]int)
	handleToEntity = make(map[heap.Handle]int64)
	collisionFrame = collisionFrame[:0]
	lastHitValid = false
	if deadEntities == nil {
		deadEntities = make(map[int64]struct{})
	} else {
		for k := range deadEntities {
			delete(deadEntities, k)
		}
	}
}

// RegisterEntityBufferLink binds an entity id to a physics matrix buffer index for collision queries.
// Call after BODY3D.COMMIT (buffer index = BODY3D.BUFFERINDEX(body)) and before relying on EntityCollided.
func RegisterEntityBufferLink(entityID int64, bufIdx int) {
	collisionMu.Lock()
	defer collisionMu.Unlock()
	if deadEntities != nil {
		delete(deadEntities, entityID)
	}
	if bufIdxToEntity == nil {
		bufIdxToEntity = make(map[int]int64)
		entityToBufIdx = make(map[int64]int)
		handleToEntity = make(map[heap.Handle]int64)
		bufIdxToBodyHeap = make(map[int]heap.Handle)
	}
	// Replace mapping if entity moved buffer index
	if old, ok := entityToBufIdx[entityID]; ok && old != bufIdx {
		delete(bufIdxToEntity, old)
	}
	bufIdxToEntity[bufIdx] = entityID
	entityToBufIdx[entityID] = bufIdx
	if h, ok := bufIdxToBodyHeap[bufIdx]; ok && h != 0 {
		handleToEntity[h] = entityID
	}
}

// UnregisterEntityCollision removes an entity from the bridge (e.g. ENTITY.FREE or CLEARPHYSBUFFER).
func UnregisterEntityCollision(entityID int64) {
	collisionMu.Lock()
	defer collisionMu.Unlock()
	if deadEntities == nil {
		deadEntities = make(map[int64]struct{})
	}
	deadEntities[entityID] = struct{}{}
	if entityToBufIdx == nil {
		return
	}
	idx, ok := entityToBufIdx[entityID]
	if !ok {
		return
	}
	delete(entityToBufIdx, entityID)
	delete(bufIdxToEntity, idx)
	if h, ok := bufIdxToBodyHeap[idx]; ok && h != 0 {
		delete(handleToEntity, h)
	}
}

func registerBufferBodyForCollision(bufIdx int, bodyHandle heap.Handle) {
	collisionMu.Lock()
	defer collisionMu.Unlock()
	if bufIdxToBodyHeap == nil {
		bufIdxToBodyHeap = make(map[int]heap.Handle)
	}
	bufIdxToBodyHeap[bufIdx] = bodyHandle
	if entityToBufIdx != nil {
		// If entity already linked this buffer index, bind handle.
		if eid, ok := bufIdxToEntity[bufIdx]; ok {
			if handleToEntity == nil {
				handleToEntity = make(map[heap.Handle]int64)
			}
			handleToEntity[bodyHandle] = eid
		}
	}
}

func unregisterBufferBodyForCollision(bufIdx int) {
	collisionMu.Lock()
	defer collisionMu.Unlock()
	if bufIdxToBodyHeap == nil {
		return
	}
	if h, ok := bufIdxToBodyHeap[bufIdx]; ok {
		delete(handleToEntity, h)
	}
	delete(bufIdxToBodyHeap, bufIdx)
}

// matrixTranslationForBodyID reads the last-published translation from the matrix buffer when registered.
func matrixTranslationForBodyID(id *jolt.BodyID) (x, y, z float32, ok bool) {
	if id == nil {
		return 0, 0, 0, false
	}
	joltBodyMu.Lock()
	idx, has := bufferIndexMap[id]
	joltBodyMu.Unlock()
	if !has {
		return 0, 0, 0, false
	}
	off := idx * 16
	joltMu.Lock()
	buf := matrixBuffer
	joltMu.Unlock()
	if off+15 >= len(buf) {
		return 0, 0, 0, false
	}
	return buf[off+12], buf[off+13], buf[off+14], true
}

func collectContactsAfterStep(m *Module) {
	joltMu.Lock()
	ps := joltSys
	bi := joltBi
	joltMu.Unlock()
	if ps == nil || bi == nil || m == nil || m.h == nil {
		return
	}

	collisionMu.Lock()
	lastHitValid = false
	if len(bufIdxToEntity) == 0 {
		collisionFrame = collisionFrame[:0]
		collisionMu.Unlock()
		return
	}

	pairSeen := make(map[string]struct{}, 32)
	collisionFrame = collisionFrame[:0]

	for bufIdx, entA := range bufIdxToEntity {
		if isDead(entA) {
			continue
		}
		bh, ok := bufIdxToBodyHeap[bufIdx]
		if !ok || bh == 0 {
			continue
		}
		bo, err := heap.Cast[*body3dObj](m.h, bh)
		if err != nil || bo.queryShape == nil || bo.id == nil {
			continue
		}
		var pos jolt.Vec3
		if px, py, pz, ok := matrixTranslationForBodyID(bo.id); ok {
			pos = jolt.Vec3{X: px, Y: py, Z: pz}
		} else {
			pos = bi.GetPosition(bo.id)
		}
		hits := ps.CollideShapeGetHits(bo.queryShape, pos, 16, 1e-3)
		for _, hit := range hits {
			if hit.BodyID == nil {
				continue
			}
			if hit.BodyID == bo.id {
				continue
			}
			otherH, ok := joltLookupHandle(hit.BodyID)
			if !ok {
				continue
			}
			entB, ok := handleToEntity[otherH]
			if !ok || entB == 0 {
				continue
			}
			if isDead(entB) {
				continue
			}
			k := pairKey(entA, entB)
			if _, dup := pairSeen[k]; dup {
				continue
			}
			pairSeen[k] = struct{}{}

			var otherPos jolt.Vec3
			if ox, oy, oz, ok2 := matrixTranslationForBodyID(hit.BodyID); ok2 {
				otherPos = jolt.Vec3{X: ox, Y: oy, Z: oz}
			} else {
				otherPos = bi.GetPosition(hit.BodyID)
			}
			dx := float64(otherPos.X - pos.X)
			dy := float64(otherPos.Y - pos.Y)
			dz := float64(otherPos.Z - pos.Z)
			lenv := math.Sqrt(dx*dx + dy*dy + dz*dz)
			var nx, ny, nz float64 = 0, 1, 0
			if lenv > 1e-8 {
				nx, ny, nz = dx/lenv, dy/lenv, dz/lenv
			}
			px := float64(hit.ContactPoint.X)
			py := float64(hit.ContactPoint.Y)
			pz := float64(hit.ContactPoint.Z)
			force := math.Abs(float64(hit.PenetrationDepth))

			hd := ContactHitData{NX: nx, NY: ny, NZ: nz, PX: px, PY: py, PZ: pz, Force: force}
			collisionFrame = append(collisionFrame, contactFrameEvent{E1: entA, E2: entB, Hit: hd})
		}
	}
	collisionMu.Unlock()
}

func isDead(id int64) bool {
	if deadEntities == nil {
		return false
	}
	_, ok := deadEntities[id]
	return ok
}

func pairKey(a, b int64) string {
	if a > b {
		a, b = b, a
	}
	return fmt.Sprintf("%d:%d", a, b)
}

// PairCollidedThisFrame reports whether two entities had a contact this frame and returns hit data.
func PairCollidedThisFrame(a, b int64) (ContactHitData, bool) {
	collisionMu.Lock()
	defer collisionMu.Unlock()
	if isDead(a) || isDead(b) {
		return ContactHitData{}, false
	}
	want := pairKey(a, b)
	for _, ev := range collisionFrame {
		if pairKey(ev.E1, ev.E2) == want {
			lastHit = ev.Hit
			lastHitValid = true
			return ev.Hit, true
		}
	}
	return ContactHitData{}, false
}

// CountCollisionsForEntity counts distinct contact pairs involving entity e this frame.
func CountCollisionsForEntity(e int64) int {
	collisionMu.Lock()
	defer collisionMu.Unlock()
	if isDead(e) {
		return 0
	}
	n := 0
	for _, ev := range collisionFrame {
		if ev.E1 == e || ev.E2 == e {
			n++
		}
	}
	return n
}

// LastCollisionData returns data from the last successful EntityCollided match (see mbentity).
func LastCollisionData() ContactHitData {
	collisionMu.Lock()
	defer collisionMu.Unlock()
	if !lastHitValid {
		return ContactHitData{}
	}
	return lastHit
}

// EntityIDForBodyHandle returns the entity# linked to a Body3D heap handle, if registered via LINKPHYSBUFFER.
func EntityIDForBodyHandle(bodyH heap.Handle) (int64, bool) {
	collisionMu.Lock()
	defer collisionMu.Unlock()
	if bodyH == 0 {
		return 0, false
	}
	id, ok := handleToEntity[bodyH]
	return id, ok
}

// entityIDForCollisionRuleHandle resolves PHYSICS3D.ONCOLLISION rule handles: BODY3D handles with LINKPHYSBUFFER, or EntityRef handles.
func entityIDForCollisionRuleHandle(m *Module, h heap.Handle) (int64, bool) {
	if h == 0 {
		return 0, false
	}
	if eid, ok := EntityIDForBodyHandle(h); ok {
		return eid, true
	}
	if m == nil || m.h == nil {
		return 0, false
	}
	obj, ok := m.h.Get(h)
	if !ok {
		return 0, false
	}
	if obj.TypeTag() == heap.TagEntityRef {
		if ref, ok := obj.(interface{ GetID() int64 }); ok {
			return ref.GetID(), true
		}
	}
	return 0, false
}

// ReviveEntity clears a tombstone so a reused id is not treated as dead (defensive).
func ReviveEntity(id int64) {
	collisionMu.Lock()
	defer collisionMu.Unlock()
	if deadEntities != nil {
		delete(deadEntities, id)
	}
}
