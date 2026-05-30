# Joint Commands

Unified joint shorthand namespace — wrappers over `JOINT3D.*` for hinge and point constraints. Prefer `JOINT3D.*` for new code; `JOINT.*` is kept for backward compatibility.

## Commands

### `JOINT.CREATEHINGE(bodyA, bodyB, ax, ay, az, pivotBx, pivotBy, pivotBz)` 

Creates a hinge constraint. Returns a joint handle. See [JOINT3D.md](JOINT3D.md) for full parameter documentation.

---

### `JOINT.CREATEPOINT(bodyA, bodyB, px, py, pz)` 

Creates a point (ball-and-socket) constraint at world point `(px, py, pz)`. Returns a joint handle.

---

### `JOINT.FREE(jointHandle)` 

Destroys the joint constraint.

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `JOINT.MAKEHINGE(bodyA, bodyB, ax,ay,az, bx,by,bz)` | Deprecated alias of `JOINT.HINGE`. |
| `JOINT.MAKEPOINT(bodyA, bodyB, ax,ay,az)` | Deprecated alias of `JOINT.POINT`. |

---

## See also

- [JOINT3D.md](JOINT3D.md) — full 3D joint API (hinge, slider, cone, fixed)
- [JOINT2D.md](JOINT2D.md) — 2D joints (distance, revolute, prismatic)
- [PHYSICS_ADVANCED.md](PHYSICS_ADVANCED.md) — joint usage examples
