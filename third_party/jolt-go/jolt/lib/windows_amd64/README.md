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

## Toolchain matching (LTO / GCC)

If linking fails with **LTO** or **undefined reference** errors that mention a **GCC version mismatch**, rebuild **`libJolt.a`** and **`libjolt_wrapper.a`** with the **same** MinGW **`gcc` / `g++`** you use for **`CGO`**. Vendored archives built with a different compiler may not link cleanly. Prefer running [`build-libs-windows.ps1`](../../../scripts/build-libs-windows.ps1) on the machine that will compile the Go binary, then link again.
