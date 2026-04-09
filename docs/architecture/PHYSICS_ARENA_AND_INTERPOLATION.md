# Physics arena, fixed timestep, and render interpolation

Design target for high–refresh-rate games: **Jolt** (and Box2D) run at a **fixed simulation rate** while **rendering** uses the display’s refresh rate, with **interpolation** between physics states.

## Goals

| Goal | Approach |
|------|----------|
| Reduce CGO crossings | Batch body transforms into a **Structure-of-Arrays** buffer (positions, rotations) updated from native code once per physics step. |
| Stable simulation | Fixed **`dt`** (e.g. 1/120 s or 1/60 s); accumulate real time and run **N** steps per frame when behind. |
| Smooth visuals | Store **previous** and **current** transforms; at draw time use **`alpha = (now - lastPhysicsTime) / fixedDt`** (clamped to [0,1]) to interpolate. |

## Threading model (non-negotiable with Raylib)

- **Raylib** and typical **OpenGL** contexts are **main-thread only**.
- A **worker goroutine** may run **`PhysicsSystem.Update`** only if the Jolt integration is documented thread-safe for that path; otherwise **step physics on the main thread** before **`Render.Frame`**.
- **Never** call **`rl.*`** draw from a background goroutine.

Recommended pattern:

```text
main thread:  poll input → net → fixedStepPhysics() → interpolateTransforms(alpha) → draw → Render.Frame
```

## Data layout (future implementation)

- **Arena:** one contiguous block (C `malloc` or Go slice with stable address) holding `N` × `(vec3 pos, quat rot)` **double-buffered** (prev / current) per body index.
- **Mapping:** `heap.Handle` / `BodyID` → row index in the SoA (maintained in [runtime/physics3d/](../../runtime/physics3d/)).

## Alpha calculation

```text
alpha = min(1, (tReal - tPhysics) / fixedDt)
```

Use **`alpha`** only for **rendering**; **gameplay** reads should use the **current** physics state after the last completed step.

## References

- [ENGINE_IR_V3.md](../../ENGINE_IR_V3.md) — VM performance context.
- Jolt integration: [runtime/physics3d/](../../runtime/physics3d/).
