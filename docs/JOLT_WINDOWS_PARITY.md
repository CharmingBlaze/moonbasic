# Jolt parity: Windows and Linux

Native **Jolt Physics** is supported on **Windows** and **Linux** when `moonrun` is built with **`CGO_ENABLED=1`**, the **`fullruntime`** tag, and the Jolt static libraries present.

If `moonrun -version` prints `Jolt backend: stub`, this binary was built **without CGO** (or on an unsupported OS). Soft stubs allow `BODY3D.CREATE` / `COMMIT` and aero/vehicle host helpers for API checks, but they are **not** a full physics engine. Official IDE / fullruntime release zips ship a CGO-linked `moonrun`.

## Build tags

| Path | Build constraint |
|------|------------------|
| Native Jolt (`*_cgo.go`) | `(linux \|\| windows) && cgo` |
| Soft stubs (`*_stub.go`) | `(!linux && !windows) \|\| !cgo` |

Do **not** use `!linux \|\| !cgo` for stubs — that wrongly includes stubs on **Windows+CGO**. See [PHYSICS.md](PHYSICS.md#build-tag-contract-for-physics3d) and [AGENTS.md](../AGENTS.md).

## Windows: link native Jolt

1. **Go** with a **MinGW-w64** toolchain (`gcc` / `g++` on `PATH`, typically MSYS2).
2. Build libs into `third_party/jolt-go/jolt/lib/windows_amd64/`:

```powershell
# Set JPH_SRC to a JoltPhysics checkout if needed
powershell -File third_party/jolt-go/scripts/build-libs-windows.ps1
```

3. Build / run:

```powershell
$env:CGO_ENABLED = "1"
go build -tags fullruntime -o moonrun.exe ./cmd/moonrun
.\moonrun.exe -version   # expect: Jolt backend: native (Windows/Linux CGO + Jolt)
```

Release CI (`release.yml` / `windows_fullruntime`) rebuilds the libs and links `moonrun` the same way. Details: [DEVELOPER.md](DEVELOPER.md#physics3d--jolt-on-windows-cgo), [BUILDING.md](BUILDING.md).

## Soft stub behavior

Without CGO, `BODY3D.COMMIT` still returns a **Body3D** handle (Euler integration + shared aero/vehicle Go helpers). `ENTITY.PHYSICS` / `ENTITY.ADDPHYSICS` remain no-ops. Prefer the official fullruntime binary for games.
