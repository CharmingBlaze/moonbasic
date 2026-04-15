# MoonBASIC Physics Synchronization

MoonBASIC uses a **dual-path physics architecture** to ensure stability across platforms while leveraging high-performance native solvers where available.

## 1. Dual-Path Architecture

### Path A: Native Jolt Physics (CGO, Linux and Windows)

On **Linux or Windows** with **CGO enabled** and **Jolt static libraries** available for the platform, MoonBASIC links to the **Jolt Physics** engine. This is the "gold standard" for the engine and supports:
- **`CharacterVirtual`**: Native kinematic character controller (KCC) with robust sweep tests, stair stepping, and slope management.
- **Rigid Body Dynamics**: Full multi-threaded solver for cubes, spheres, and complex meshes.
- **Shared Memory Sync**: Physics state is written to a shared buffer and synced back to entities each frame via `syncEntitiesFromPhysics`.

### Path B: Stub / non-desktop / no native Jolt

When **CGO is disabled**, or the OS is **not** Linux/Windows (e.g. some builds), or **Jolt libraries are not linked**, MoonBASIC uses **`physics3d` stub files** (no rigid-body solver hits) and **`runtime/player` stub files** that return **clear errors** for **`PLAYER.*` / `CHAR.*` / `CHARACTER.*`** KCC commands (see `errPlayerRequiresCGOJolt` in the player package). There is **no second Go “host” character solver**: desktop **Windows and Linux** with **CGO + Jolt** are the supported path for real KCC physics.

### Vendored Jolt Go API (`third_party/jolt-go`)

The module **`github.com/bbitechnologies/jolt-go`** is **replaced** in this repo by [`third_party/jolt-go`](../third_party/jolt-go) (see root `go.mod`). It intentionally exposes a **small** C wrapper surface (bodies, shapes, raycasts, `AddImpulse`, linear velocity, etc.). It does **not** mirror the full Jolt C++ API.

For contributors and script authors, the practical rules are:

| Area | Behavior on Linux or Windows + CGO (native Jolt) |
|------|-------------------------|
| **`JOINT.CREATEHINGE` / `JOINT.CREATEPOINT`** | Allocates a **placeholder** joint handle so scripts and stubs stay aligned. Real hinge/point constraints are **not** created until the wrapper grows constraint APIs. |
| **`BODY3D.SETGRAVITYFACTOR`**, **`SETDAMPING`**, **`LOCKAXIS`**, **`SETCCD`** | Parsed, validated, then return a deterministic **\"not implemented on native backend\"** error (no silent no-op). |
| **`BODY3D.GETLINEARVEL` / `SETLINEARVEL`** | Uses **`GetLinearVelocity` / `SetLinearVelocity`** on `BodyInterface` (supported). Returns a **3-element numeric array** handle (same pattern as `BODY3D.GETPOS`). |
| **`BODY3D.GETANGULARVEL` / `SETANGULARVEL`** | Getter returns **zeros**; setter returns **\"not implemented on native backend\"** (angular state setter not exposed). |
| **`BODY3D.APPLYFORCE`** | Implemented as **`AddImpulse(F × dt)`** using the module fixed timestep (`PHYSICS3D.SETTIMESTEP` / default 1/60 s), because the wrapper has **`AddImpulse`** but not **`AddForce`**. Treat script “force” as **impulse per physics step** in spirit. |
| **`BODY3D.APPLYTORQUE`** | Returns **\"not implemented on native backend\"** (torque API not exposed in the wrapper). |
| **`BODY3D.GETMASS`** | Returns **`1.0`** (mass not queried from Jolt in this binding). |
| **Ground probe hook** | **`mbphysics3d.SetRaycastHook`** is implemented on **both** [`pick_cgo.go`](../runtime/physics3d/pick_cgo.go) (Jolt path) and [`pick_stub.go`](../runtime/physics3d/pick_stub.go) (no Jolt) so `mbentity` can register a floor query with the same API on every build. |

gopls on a **Windows** machine may type-check Linux-only files when **`GOOS=linux`** is set in `gopls.build.env`; vendored **`jolt-go`** sources should still analyze cleanly.

---

## Build tag contract for physics3d

To avoid **duplicate symbols** (stub + native compiled together), use **mutually exclusive** tags:

| Role | Typical `//go:build` |
|------|----------------------|
| Native Jolt implementation (bodies, picks, collision bridge, matrix export) | ``(linux || windows) && cgo`` |
| Stub / no native Jolt (same Go signatures, no-ops or “no hit”) | ``(!linux && !windows) || !cgo`` |

Do **not** use ``!linux || !cgo`` for `physics3d` stubs: on **Windows + CGO** that expression is true and **overlaps** the native path.

**Windows:** place prebuilt **`libJolt.a`** and **`libjolt_wrapper.a`** under [`third_party/jolt-go/jolt/lib/windows_amd64/`](../third_party/jolt-go/jolt/lib/windows_amd64/README.md) or build them with [`build-libs-windows.ps1`](../third_party/jolt-go/scripts/build-libs-windows.ps1). Without them, **link** fails even when **compile** succeeds.

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
1. Implement the **Jolt** path in `*_cgo.go` (common for Linux and Windows).
2. Implement an **identical Go signature stub** in `*_stub.go` for builds without native Jolt (errors or no-ops as appropriate).
3. Ensure both paths expose the **same manifest keys** in `compiler/builtinmanifest/commands.json`.

---

## 4. Building Native Jolt on Windows

To enable Path A on Windows, you must link the native static libraries:
1. **Toolchain**: Install **MinGW-w64** (via MSYS2) and ensure `gcc` is on your `PATH`.
2. **Libraries**: Run `third_party/jolt-go/scripts/build-libs-windows.ps1` to compile `libJolt.a` and `libjolt_wrapper.a`.
3. **Build**: Run MoonBASIC with `CGO_ENABLED=1`:
   `go run -tags fullruntime . --run examples/physics_demo.mb`

---

## 5. Troubleshooting

- **Flickering Grounded State**: Check the `joltGroundRayStartLift` and `joltGroundPastFeetSkin` constants in `entity_phys_sync_cgo.go`. The ray must start slightly above the pivot to avoid self-hits.
- **Micro-Jitter on Slopes**: Ensure the `groundNormal` check in the probe logic properly filters out wall hits (normal Y < 0.28).

---

## 6. KCC contacts vs rigid-body collision queue

Jolt **`CharacterVirtual`** uses an internal **`CharacterContactListener`** path. Scripts can drain a small event queue via **`CHARACTERREF.DRAINCONTACTS`** (and related listener toggles). That pipeline is **separate** from **`PHYSICS3D.PROCESSCOLLISIONS`**, which feeds the rigid-body collision callback queue used for dynamic bodies.

- **Do not assume** the same events appear in both places; order and pairing differ.
- **One-way platforms** use object layer **`ONE_WAY` (4)** on static/kinematic bodies; the character listener weakens contact response when the hit comes from **below** (see `third_party/jolt-go/jolt/wrapper/character.cpp`).
- If you need a **single** script callback for both KCC and rigid bodies, fan in manually (e.g. poll `DRAINCONTACTS` in the same frame as `PROCESSCOLLISIONS`) or keep two handlers and document ordering for your game.
