# Physics 3D (Jolt)

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **Physics3D.Start / Stop / Step** | **`PHYSICS3D.START`**, **`STOP`**, **`STEP`** | **Linux + CGO + jolt-go** only; other platforms use stubs (see runtime message). |
| **Physics3D.SetGravity** | **`PHYSICS3D.SETGRAVITY`** `(gx, gy, gz)` | |
| **Physics3D.DebugDraw** | **`PHYSICS3D.DEBUGDRAW`** | Currently returns an error (not wired in this runtime). |
| **Body3D.Create** | **`BODY3D.MAKE`** → **`ADDBOX`** / **`ADDSPHERE`** / **`ADDCAPSULE`** → **`COMMIT`** `(builder, x, y, z)` | Builder freed on commit; body handle returned. |
| **Shapes** | **`BODY3D.ADDBOX`**, **`ADDSPHERE`**, **`ADDCAPSULE`**, **`ADDMESH`** (mesh not implemented) | A duplicate query shape is kept for overlap tests. |
| **Body3D.SetMass / friction / restitution / forces / velocity** | **`BODY3D.SETMASS`**, **`SETFRICTION`**, **`SETRESTITUTION`**, **`APPLYFORCE`**, **`APPLYIMPULSE`**, **`SETLINEARVEL`**, **`SETANGULARVEL`**, **`SETROT`** | **No-op** in the current **`jolt-go`** binding (minimal C wrapper). Do not rely on dynamics until the binding exposes `BodyInterface` force APIs. |
| **Collision** | **`BODY3D.COLLIDED`**, **`COLLISIONOTHER`** | Uses **`CollideShapeGetHits`** with a duplicate shape at the body position. **`COLLISIONOTHER`** resolves the other body’s VM handle when it was registered at commit. |
| **Collision point / normal** | **`BODY3D.COLLISIONPOINT`** (first hit contact point), **`BODY3D.COLLISIONNORMAL`** | Normal is **not** available from the hit query; **`COLLISIONNORMAL`** returns placeholder `(0,0,1)` — use **`PHYSICS3D.RAYCAST`** for a surface normal. |
| **Joint3D.*** | **`JOINT3D.FIXED`**, **`HINGE`**, **`SLIDER`**, **`CONE`**, **`DELETE`** | Return errors: Jolt constraints are not exposed in **`jolt-go`**. |

**Aliases:** **`PHYSICS.*`** mirrors several **`PHYSICS3D.*`** names.

See also: [PHYSICS3D.md](../PHYSICS3D.md).
