# Command backlog — what is done vs still open

The original long checklist lived here; most **Draw / Text / Texture / Window / Input / Camera / Camera2D / Audio / Shader uniforms / Image process / Vec–Quat math** items are now implemented under `runtime/` (CGO builds).

Use **`compiler/builtinmanifest/commands.json`** as the authoritative list of registered names and arities.

---

## Still open (engine or binding gaps)

| Area | Gap | Reason |
|------|-----|--------|
| **Physics3D / Body3D** | `GetLinearVel`, `GetAngularVel`, `IsActive`, `SetRotation`, forces, mass, friction, `GetRot` | Current **jolt-go** `BodyInterface` exposes position, activate/deactivate, create/set shape, etc. **Velocity, rotation, and impulses** need additional C wrapper exports in `jolt-go` before we can bind them safely. |
| **Body3D** | `Activate` / `Deactivate` | **Implemented** on **Linux + CGO** (`BODY3D.ACTIVATE`, `BODY3D.DEACTIVATE`). Stubs on other platforms. |
| **Physics3D** | `RAYCAST` | **Implemented** (Linux + CGO). Returns a 6-float array (hit, normal×3, fraction, reserved). See [reference/PHYSICS3D.md](reference/PHYSICS3D.md). |
| **Draw** | Hot-path heap allocations | Some helpers still allocate; see engine memory-safety notes (`ReleaseOnce`, pools). |

---

## Recently completed (reference docs)

| Doc | Topics |
|-----|--------|
| [reference/IMAGE.md](reference/IMAGE.md) | `IMAGE.DRAWIMAGE`, `DITHER`, `MIPMAPS`, `FORMAT`, `DRAWRECTLINES`, `ALPHACROP`, `ALPHACLEAR` |
| [reference/VEC_QUAT.md](reference/VEC_QUAT.md) | `VEC3.TRANSFORMMAT4`, `ANGLE`, `PROJECT`, `ORTHONORMALIZE`, `ROTATEBYQUAT`, `VEC2.TRANSFORMMAT4`, `QUAT.*` (incl. `TOEULER`, `FROMVEC3TOVEC3`, `FROMMAT4`, `TRANSFORM`) |
| [reference/SHADER.md](reference/SHADER.md) | `SHADER.FREE`, `GETLOC`, `SETFLOAT`/`VEC*`, `SETINT`, `SETTEXTURE` |
| [reference/PHYSICS3D.md](reference/PHYSICS3D.md) | `Body3D.Activate` / `Deactivate` |

---

## How to add a new command

1. Implement under the correct `runtime/<module>/` file (one concern per file when practical).
2. Register in the module’s `Register()`; add **`!cgo` stub** names to the matching `stub.go`.
3. Append **`commands.json`** with correct `args` / `returns`.
4. Add a **working example** to the relevant `docs/reference/*.md`.
