# Windows x64 Jolt static libraries

Place **`libJolt.a`** and **`libjolt_wrapper.a`** in this directory so CGO can link native Jolt on Windows (`third_party/jolt-go/jolt/cgo_windows_amd64.go`).

## Build locally

1. Clone [JoltPhysics](https://github.com/jrouwe/JoltPhysics) and set **`JPH_SRC`** to that directory.
2. Install **CMake** and **MinGW-w64** (`g++`, `ar` on `PATH`).
3. From the repository root:

```powershell
powershell -File third_party/jolt-go/scripts/build-libs-windows.ps1
```

The script writes both archives here. Re-run after upgrading Jolt or changing the C++ wrapper.

## Policy

These binaries are large; they may be built locally or committed per release policy. Without them, **`go build -tags fullruntime`** with **`CGO_ENABLED=1`** will fail at link time on Windows when importing Jolt.

## Toolchain matching (GCC)

[`build-libs-windows.ps1`](../../../scripts/build-libs-windows.ps1) compiles Jolt **without LTO** (`-fno-lto`, **`CMAKE_INTERPROCEDURAL_OPTIMIZATION=OFF`**) so **`libJolt.a`** contains plain objects—see **“Windows full-runtime PE link model”** in [`docs/BUILDING.md`](../../../../../docs/BUILDING.md). If you still see **undefined reference** errors, rebuild both archives with the **same** MinGW **`g++`** you use for **`go build`**.
