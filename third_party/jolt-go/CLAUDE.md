# jolt-go: LLM Assistant Context

## What This Is

Go bindings for Jolt Physics (C++ engine). Pre-built binaries included in repository. No user compilation required.

**Architecture:** Go → CGO → C Wrapper → Jolt Physics C++

**Platforms:** macOS ARM64, Linux x86-64

## File Structure

- `cgo.go` - Package documentation and CGO compiler flags
- `cgo_{platform}_{arch}.go` - Platform-specific linker flags
- `jolt/*.go` - Feature-specific Go APIs (body, character, physics system, vectors)
- `jolt/wrapper/*.{cpp,h}` - C wrapper around Jolt C++ API (opaque pointers)
- `jolt/lib/{platform}/` - Pre-built static libraries (libJolt.a, libjolt_wrapper.a), committed to git
- `scripts/build-libs.sh` - Builds binaries for all platforms
- `scripts/docker/` - Docker build environment for Linux
- `example/main.go` - Falling sphere demo
- `.github/workflows/build-binaries.yml` - CI/CD for building and testing

## Core Values

1. **Never break zero-compilation guarantee** - Users must never need to compile C++
2. **ABI compatibility is critical** - Wrapper and Jolt must use identical compiler flags
3. **Test both platforms** - darwin/arm64 AND linux/amd64 before any release
4. **Opaque pointers only** - Never expose C++ types directly to Go
5. **Memory management** - Every `Create*()` needs `Destroy()`

## Development Instructions

### When Modifying Go API (`jolt/*.go` files)
1. Add C wrapper function in `jolt/wrapper/*.{cpp,h}` first
2. Use `extern "C"` and opaque pointers in wrapper
3. Add Go functions in appropriate `jolt/*.go` file (body, character, physics system, etc.)
4. Rebuild binaries: `./scripts/build-libs.sh all`
5. Test: `go run example/main.go`
6. Update README.md API section if adding public functions

### When Modifying C Wrapper (`jolt/wrapper/*.{cpp,h}` files)
1. Keep flags matching Jolt build: `-DJPH_DISABLE_CUSTOM_ALLOCATOR -DJPH_PROFILE_ENABLED -DJPH_DEBUG_RENDERER -DJPH_OBJECT_STREAM`
2. Return error codes (0=success, -1=fail), never throw exceptions
3. Rebuild binaries for BOTH platforms
4. Update CONTRIBUTORS.md if changing build process

### When Updating Jolt Physics Version
1. Update `$JOLT_SRC` checkout to new version
2. Rebuild: `./scripts/build-libs.sh all`
3. Test: `go run example/main.go`
4. Update version number in CONTRIBUTORS.md (line 13)
5. Update README.md if new features exposed (line 10)
6. Commit binaries: `git add jolt/lib/ && git commit -m "Update to Jolt vX.X.X"`
7. Tag release: `git tag vX.X.X && git push origin main --tags`

### When Changing Build System (`scripts/`)
1. Test darwin build: `./scripts/build-libs.sh darwin_arm64`
2. Test linux build: `./scripts/build-libs.sh linux_amd64`
3. Verify output in `jolt/lib/{platform}/`
4. Update CONTRIBUTORS.md if process changes

### When Adding Features
1. Check if feature exists in upstream Jolt Physics
2. Add wrapper functions (keep minimal surface area)
3. Add Go API functions (use Go idioms)
4. Update `example/main.go` if demonstrating new capability
5. Update README.md Quick Start or API Overview sections

### When Fixing Bugs
1. Identify layer: Go, CGO, C wrapper, or Jolt
2. Fix in appropriate file
3. If wrapper change: rebuild binaries
4. Test with example
5. Update README.md Troubleshooting if user-facing issue

### Documentation Updates

**README.md** - Update when:
- Adding/changing public API functions
- Changing installation process
- Adding platform support
- Updating Jolt Physics version (line 10)
- Adding examples or usage patterns

**CONTRIBUTORS.md** - Update when:
- Changing build process or scripts
- Updating Jolt Physics version (line 13)
- Modifying compiler flags
- Changing CI/CD workflow
- Adding platform support

**CLAUDE.md (this file)** - Update when:
- Changing project architecture
- Adding new directories or major files
- Modifying core principles or workflows
- Updating critical rules or constraints

## Guidelines

**Be concise.** No fluff, no unnecessary explanations.

**Think from first principles.** Understand why before changing what.

**Follow industry best practices:**
- Go: Effective Go conventions, clear error handling
- C++: RAII, no raw pointers, const correctness
- CGO: Minimize crossings, opaque pointers, explicit memory management
- Cross-platform: Test both, avoid platform assumptions

**Verify before suggesting.** Read actual code, don't guess.

**Prioritize maintainability.** Clear code > clever code.

**When uncertain:** Check README.md (user view) or CONTRIBUTORS.md (build details).
