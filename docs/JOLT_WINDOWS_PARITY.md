# Jolt parity: Windows vs Linux [RESOLVED]

As of v1.3.1, MoonBASIC supports full native **Jolt Physics** parity on Windows using CGO and prebuilt static libraries.

## Current State

- **Cross-Platform CGO**: Architecture has been migrated from `*_linux.go` to platform-agnostic `*_cgo.go` files.
- **Windows Support**: Native Jolt is supported on Windows x86-64 when `CGO_ENABLED=1` is set.
- **Unified API**: The "Easy Mode" API and "Host KCC" (Path B) are synchronized with the native Jolt behavior (Path A), ensuring scripts are portable even in non-CGO environments.

## How to Build on Windows

To link Jolt natively on Windows, you must follow these steps:

### 1. Requirements
- **Go 1.25.3+**
- **MinGW-w64** (via MSYS2). Ensure `gcc` is on your `PATH`.
- **CMake** (required for building the Jolt static libraries).
- A local checkout of **[JoltPhysics](https://github.com/jrouwe/JoltPhysics)** (e.g. **v5.4.0** to match vendored **jolt-go**). The build script does **not** clone it; you point **`JPH_SRC`** at that directory.

### 2. Build Jolt Static Libraries
MoonBASIC requires `libJolt.a` and `libjolt_wrapper.a`. Use [third_party/jolt-go/scripts/build-libs-windows.ps1](../third_party/jolt-go/scripts/build-libs-windows.ps1) from the repository root with **`JPH_SRC`** set:

```powershell
# Example: JoltPhysics cloned beside moonbasic, or set to your path
$env:JPH_SRC = "C:\path\to\JoltPhysics"
powershell -File third_party/jolt-go/scripts/build-libs-windows.ps1
```

The script will:
- Configure and compile the Jolt core with CMake from **`JPH_SRC`** (under `Build/windows_amd64_release`).
- Compile the C++ wrapper sources in **jolt-go** and archive **`libjolt_wrapper.a`**.
- Place both **`libJolt.a`** and **`libjolt_wrapper.a`** in **`third_party/jolt-go/jolt/lib/windows_amd64/`**.

See also [third_party/jolt-go/jolt/lib/windows_amd64/README.md](../third_party/jolt-go/jolt/lib/windows_amd64/README.md).

### 3. Compile MoonBASIC with Physics
Use the `fullruntime` tag and enable CGO:
```powershell
$env:CGO_ENABLED="1"
go run -tags fullruntime ./cmd/moonrun examples/physics_demo.mb
```

## Technical Implementation Details

For other programmers and maintainers:

- **LDFLAGS**: Windows-specific linking is handled in `third_party/jolt-go/jolt/cgo_windows_amd64.go`. It links `-lJolt`, `-ljolt_wrapper`, and the required C++ standard libraries (`-lstdc++`).
- **Build Tags**: Promotion of modules was achieved by switching from `linux && cgo` to `(linux || windows) && cgo`.
- **Stub Sync**: `runtime/charcontroller/stub.go` implements a lightweight AABB capsule solver that mirrors the native `CHARACTERREF.*` API surface.

## Summary

The "multi-day effort" to bring Jolt to Windows is complete. Windows is no longer a second-class citizen for physics development in MoonBASIC.
