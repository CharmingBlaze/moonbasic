# Physics Architecture & Internal Implementation

This document describes the low-level implementation of the MoonBASIC physics system. For usage instructions and beginner guides, see [docs/PHYSICS.md](../PHYSICS.md).

## 1. Dual-Path Architecture

MoonBASIC uses a **dual-path physics architecture** to ensure stability across platforms while leveraging high-performance native solvers where available.

### Path A: Native Jolt Physics (CGO, Linux and Windows)

On **Linux or Windows** with **CGO enabled** and **Jolt static libraries** available, MoonBASIC links to the **Jolt Physics** engine. This is the "gold standard" and supports:
- **`CharacterVirtual`**: Native kinematic character controller (KCC) with robust sweep tests, stair stepping, and slope management.
- **Rigid Body Dynamics**: Full multi-threaded solver for primitives and complex meshes.
- **Shared Memory Sync**: Physics state is written to a shared buffer and synced back to entities each frame via `syncEntitiesFromPhysics`.

### Path B: Stub / no-native Jolt

When **CGO is disabled**, or the OS is not desktop Linux/Windows, MoonBASIC uses **`physics3d` stub files**. There is **no second Go "host" character solver**: desktop builds with **CGO + Jolt** are the only path for high-fidelity physics.

---

## 2. Vendored Jolt Go API (`third_party/jolt-go`)

The module `github.com/bbitechnologies/jolt-go` is replaced in this repository by a local version under `third_party/jolt-go`. It intentionally exposes a small C wrapper surface (bodies, shapes, raycasts, linear velocity, etc.).

### Binding Limitations

| Area | Internal Handling |
|------|-------------------|
| **Joints** | Allocates a **placeholder** handle. Real hinge/point constraints are not yet created in Jolt. |
| **ApplyForce** | Implemented as **`AddImpulse(F * dt)`** because the wrapper lacks `AddForce`. |
| **Angular Vel** | Setters return "not implemented"; getters return zeros (not exposed in wrapper). |
| **Mass** | Currently hardcoded to `1.0` in the binding unless manually overridden after commit. |

---

## 3. Build Tag Contract

To avoid duplicate symbols between stubs and native implementations, use these mutually exclusive tags:

| Role | `//go:build` |
|------|--------------|
| Native Jolt | `(linux || windows) && cgo` |
| Stubs | `(!linux && !windows) || !cgo` |

**Windows Note:** Requires `libJolt.a` and `libjolt_wrapper.a` in `third_party/jolt-go/jolt/lib/windows_amd64/`.

---

## 4. Visual Synchronization Logic

### The "Visual Snap" Band
Physics solvers often have a "slop" or "allowed penetration" (approx 0.02m). To prevent jitter:
- If an entity is grounded and its distance to the floor is within `joltGroundVisualSnapBand` (default 0.14m), the **visual** model is snapped to the floor even if the physics body is technically hovering.

### Velocity Zeroing
When grounded, tiny vertical velocities (jitter) are zeroed in the sync loop while horizontal (XZ) momentum is preserved. This ensures characters stay firmly planted on moving platforms or slopes.

---

## 5. Contact Listeners & Fan-In

Jolt's `CharacterVirtual` uses a separate `CharacterContactListener` pipeline from rigid bodies.
- **KCC Contacts**: Drained via `CHARACTERREF.DRAINCONTACTS`.
- **Rigid Body Collisions**: Processed via `PHYSICS3D.PROCESSCOLLISIONS` (feeds the script callback queue).

The engine provides internal "Fan-In" hooks to merge these events when necessary for player-object interactions.
