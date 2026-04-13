# MoonBASIC Physics Synchronization

MoonBASIC uses a **dual-path physics architecture** to ensure stability across platforms while leveraging high-performance native solvers where available.

## 1. Dual-Path Architecture

### Path A: Linux + CGO (Jolt Physics)
On Linux with CGo enabled, MoonBASIC links to the **Jolt Physics** engine. This is the "gold standard" for the engine and supports:
- **`CharacterVirtual`**: Native kinematic character controller (KCC) with robust sweep tests, stair stepping, and slope management.
- **Rigid Body Dynamics**: Full multi-threaded solver for cubes, spheres, and complex meshes.
- **Shared Memory Sync**: Physics state is written to a shared buffer and synced back to entities each frame via `syncEntitiesFromPhysics`.

### Path B: Windows / Non-CGO (Host Solver)
On Windows (by default) or targets without Jolt, MoonBASIC falls back to an **Iterative Host Solver**.
- **Script-Driven Interpolation**: Most "Easy Mode" physics (like `Character.Create`) use a high-level Go implementation that mimics the Jolt KCC behavior through simpler bounding-box or raycast tests.
- **Visual Parity**: The goal of the Host Solver is to provide identical behavior to Jolt for typical "Mario 64" style gameplay, ensuring code written on Windows works perfectly when deployed to Linux servers.

---

## 2. Visual Synchronization

### The "Visual Snap" Logic
Internal physics integration often results in "micro-bounce" (solver slop of ~0.02 m). To ensure a premium aesthetic, `runtime/mbentity/entity_phys_sync_cgo.go` implements a **Visual Snap Band**:

> [!TIP]
> If an entity is grounded and its vertical distance to the floor is within the `joltGroundVisualSnapBand` (default 0.14m), its *displayed* position is snapped to the floor, even if the physics body is technically hovering slightly.

### Zeroing Vertical Velocity
When a character is detected as grounded by the `RaycastDownGroundProbe`, tiny vertical velocities (jitter) are zeroed out while horizontal (XZ) momentum is preserved. This prevents "ping-ponging" between grounded and falling states.

---

## 3. Contributor Rules for Physics

### Update Jolt Directly
Commands that modify physical properties (e.g., `ENTITY.SETBOUNCE`, `ENTITY.SETFRICTION`) must check if the entity is linked to a physics body (`physBufIndex >= 0`). If it is, you **must** update the Jolt index directly using `mbphysics3d.SetRestitutionToIndex` / `SetFrictionToIndex`. Failure to do this will result in the visual property diverging from the physical behavior.

### Platform Parity
When adding a new physics feature:
1. Implement the **Jolt** path in `*_linux.go` or `*_cgo.go`.
2. Implement an **identical Go signature stub** in `*_stub.go` for the Host Solver target.
3. Ensure both paths expose the **same manifest keys** in `compiler/builtinmanifest/commands.json`.

---

## 4. Troubleshooting

- **Flickering Grounded State**: Check the `joltGroundRayStartLift` and `joltGroundPastFeetSkin` constants in `entity_phys_sync_cgo.go`. The ray must start slightly above the pivot to avoid self-hits.
- **Micro-Jitter on Slopes**: Ensure the `groundNormal` check in the probe logic properly filters out wall hits (normal Y < 0.28).
