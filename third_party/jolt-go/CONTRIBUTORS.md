# Contributor Guide

This guide is for contributors who need to rebuild the pre-built binaries for jolt-go.

## Overview

jolt-go ships with pre-built static libraries for:
- `darwin/arm64` (macOS Apple Silicon)
- `linux/amd64` (Linux x86-64)
- `linux/arm64` (Linux ARM64)

These binaries are committed to the repository in the `jolt/lib/` directory, ensuring the package works immediately with `go get`.

**Current Jolt Physics Version:** v5.4.0

## Prerequisites

### Required Tools
- **Go 1.20+**
- **Git**
- **CMake 3.15+**
- **Docker** (for Linux builds)
- **macOS with Apple Silicon** (for darwin/arm64 builds)

## Setup

### 1. Clone/Update Jolt Physics

Clone Jolt Physics to a location of your choice:

```bash
git clone https://github.com/jrouwe/JoltPhysics.git ~/JoltPhysics
cd ~/JoltPhysics

# Optionally check out a specific version
git checkout v5.4.0
# Or use latest main
git pull origin main
```

### 2. Set JOLT_SRC Environment Variable

```bash
export JOLT_SRC=~/JoltPhysics  # or wherever you cloned it
```

## Rebuilding Binaries

### Build All Platforms

```bash
cd /path/to/jolt-go
export JOLT_SRC=/path/to/JoltPhysics
./scripts/build-libs.sh all
```

### Build Specific Platform

```bash
# Build darwin/arm64 only (requires macOS ARM64)
./scripts/build-libs.sh darwin_arm64

# Build linux/amd64 only (uses Docker)
./scripts/build-libs.sh linux_amd64

# Build linux/arm64 only (uses Docker)
./scripts/build-libs.sh linux_arm64
```

### Build Output

Binaries will be placed in:
- `jolt/lib/{platform}/libJolt.a`
- `jolt/lib/{platform}/libjolt_wrapper.a`

## Testing Binaries

After building, test that everything works:

```bash
go run example/main.go
```

## Updating Jolt Physics Version

When updating to a new Jolt Physics version:

1. **Update Jolt source:**
   ```bash
   cd $JOLT_SRC
   git pull origin main
   # Or checkout specific tag: git checkout v5.5.0
   ```

2. **Rebuild all binaries:**
   ```bash
   cd /path/to/jolt-go
   ./scripts/build-libs.sh all
   ```

3. **Test:**
   ```bash
   go run example/main.go
   go test ./...  # if you have tests
   ```

4. **Commit and tag the binaries:**
   ```bash
   git add jolt/lib/
   git commit -m "Update binaries to Jolt Physics v5.5.0"
   git tag v0.2.0
   git push origin main --tags
   ```

### Release Process

To create a new release after updating binaries:

   ```bash
   git tag v0.2.0
   git push origin main --tags
   ```
