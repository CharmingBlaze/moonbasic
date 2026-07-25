# MoonBASIC API Standardization Directive

**TO:** Compiler Engineering Team  
**FROM:** Language Design Authority  
**RE:** Complete API Consistency Overhaul  
**PRIORITY:** HIGH - Foundation for v1.0

---

## Executive Decision: API Philosophy

After thorough analysis, we are **standardizing moonBASIC around the Namespace.Method pattern** with the following core principles:

### ✅ APPROVED Design Decisions

1. **Primary API Style:** `Namespace.Method()` (e.g., `CAMERA.CREATE()`, `MODEL.SETPOS()`)
2. **Creation Verb:** `CREATE` is now the **standard** (deprecate all `MAKE` variants)
3. **Position Method:** `SETPOS` is **canonical** (deprecate `SETPOSITION`)
4. **Easy Mode:** Remains as **convenience layer only** (not primary documentation)
5. **Universal Methods:** All spatial handles **must** implement `.pos()`, `.rot()`, `.scale()` where applicable
6. **Method Design:** Favor **chainable, simple methods** over multi-argument functions

---

## Part 1: Global Naming Standards

### Creation Pattern (REQUIRED)

```basic
' ✅ CORRECT - Use CREATE for all object instantiation
camera = CAMERA.CREATE()
model = MODEL.LOAD("player.glb")      ' File loading still uses LOAD
light = LIGHT.CREATEPOINT()           ' Specific types use CREATE<Type>
body = BODY3D.CREATE()

' ❌ WRONG - MAKE is deprecated
camera = CAMERA.MAKE()                ' DEPRECATED
light = LIGHT.MAKE()                  ' DEPRECATED
```

**Action Required:**
- Rename ALL `*.MAKE` → `*.CREATE` in manifest
- Rename ALL `*.MAKE<Type>` → `*.CREATE<Type>` (e.g., `MAKECUBE` → `CREATECUBE`)
- Keep `*.MAKE` as deprecated aliases temporarily (with warnings)
- Add migration guide for users

### Position/Transform Pattern (REQUIRED)

```basic
' ✅ CORRECT - Primary names
CAMERA.SETPOS(cam, x, y, z)
MODEL.SETROT(model, pitch, yaw, roll)
SPRITE.SETSCALE(sprite, sx, sy, sz)

' ⚠️ ACCEPTABLE - Deprecated aliases (keep for now)
CAMERA.SETPOSITION(cam, x, y, z)     ' Alias to SETPOS

' ❌ WRONG - Inconsistent naming
CAMERA.POSITION(cam, x, y, z)        ' REMOVE
```

**Action Required:**
- Keep `SETPOS` as canonical
- Add `SETPOSITION` as alias everywhere for consistency (then deprecate)
- Document `SETPOS` in all examples and docs

---

## Part 2: Universal Handle Methods

### MANDATORY: All Spatial Handles Must Implement

Every handle representing a spatial object (Camera, Model, Sprite, Light, Body3D, etc.) **MUST** expose these methods where applicable:

```basic
' Position (REQUIRED for all spatial objects)
handle.pos(x, y, z)              ' Equivalent to NAMESPACE.SETPOS(handle, x, y, z)
pos = handle.pos()               ' Get position (returns array or Vec3)

' Rotation (REQUIRED for rotatable objects)
handle.rot(pitch, yaw, roll)     ' Euler angles in degrees
handle.rot(yaw)                  ' Single-axis rotation for 2D objects
rot = handle.rot()               ' Get rotation

' Scale (REQUIRED for scalable objects)
handle.scale(sx, sy, sz)         ' Non-uniform scale
handle.scale(s)                  ' Uniform scale
scale = handle.scale()           ' Get scale

' Color (REQUIRED for renderable objects)
handle.col(r, g, b)              ' RGB 0-255
handle.col(r, g, b, a)           ' RGBA 0-255
color = handle.col()             ' Get color

' Alpha (REQUIRED for renderable objects)
handle.alpha(a)                  ' 0.0-1.0
alpha = handle.alpha()           ' Get alpha

' Cleanup (REQUIRED for ALL heap handles)
handle.free()                    ' Equivalent to NAMESPACE.FREE(handle)
```

### Type-Specific Methods (Camera Example)

```basic
' Camera-specific methods (in addition to universal .pos, .rot, .scale)
cam.look(targetX, targetY, targetZ)   ' Look at point
cam.look(targetEntity)                 ' Look at entity
cam.turn(pitch, yaw, roll)            ' Relative rotation
cam.fov(degrees)                      ' Set field of view
cam.zoom(factor)                      ' Set zoom level
cam.shake(intensity, duration)        ' Camera shake effect
cam.orbit(targetX, targetY, targetZ, yaw, pitch, distance)  ' Orbit mode

' All cameras still have universal methods
cam.pos(0, 10, 20)
cam.rot(0, 45, 0)
cam.free()
```

### Type-Specific Methods (Body3D Example)

```basic
' Physics body methods (in addition to universal .pos, .rot)
body.vel(vx, vy, vz)              ' Linear velocity
body.angvel(wx, wy, wz)           ' Angular velocity
body.force(fx, fy, fz)            ' Apply force
body.impulse(ix, iy, iz)          ' Apply impulse
body.torque(tx, ty, tz)           ' Apply torque
body.mass(m)                      ' Set mass
body.friction(f)                  ' Set friction
body.restitution(r)               ' Set bounciness

' All bodies still have universal methods
body.pos(x, y, z)
body.rot(pitch, yaw, roll)
body.scale(sx, sy, sz)           ' Scale collision shape if supported
body.free()
```

### Complete Universal Method Matrix

| Handle Type | .pos() | .rot() | .scale() | .col() | .alpha() | .free() | Type-Specific |
|-------------|--------|--------|----------|--------|----------|---------|---------------|
| **CAMERA** | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | .look(), .turn(), .fov(), .zoom(), .shake(), .orbit() |
| **CAMERA2D** | ✅ | ✅* | ❌ | ❌ | ❌ | ✅ | .target(), .zoom() |
| **MODEL** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | .texture(), .shader() |
| **SPRITE** | ✅ | ✅* | ✅ | ✅ | ✅ | ✅ | .frame(), .anim() |
| **LIGHT** | ✅ | ✅ | ❌ | ✅ | ❌ | ✅ | .dir(), .intensity(), .shadows() |
| **BODY3D** | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | .vel(), .force(), .impulse(), .mass() |
| **BODY2D** | ✅ | ✅* | ✅ | ❌ | ❌ | ✅ | .vel(), .force(), .impulse(), .mass() |
| **PARTICLE** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | .emit(), .burst() |
| **DECAL** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | .project() |
| **TERRAIN** | ✅ | ✅ | ✅ | ❌ | ❌ | ✅ | .height(), .smooth() |
| **TEXTURE** | ❌ | ❌ | ❌ | ❌ | ❌ | ✅ | .width(), .height(), .filter() |
| **NAVAGENT** | ✅ | ✅ | ❌ | ❌ | ❌ | ✅ | .goto(), .speed() |
| **WATER** | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | .wave(), .flow() |

*\* 2D objects use single-value .rot(angle) for Z-axis rotation only*

**v0.9 engine note:** This matrix is the **design target**. Runtime handle dispatch (`vm/handlecall.go`) and `NAMESPACE.GET*` builtins cover the **primary spatial heap tags** used in examples (camera, entity, model, 3D/2D bodies, lights, particles, nav agent, etc.). Some rows (e.g. DECAL/TERRAIN/WATER) may still be **namespace-first** until a stable handle type is wired end-to-end. **SPRITE** stores **pos / rot / scale / color / alpha** for **`DrawTexturePro`**; **`SPRITE.HIT`** / **`SPRITE.POINTHIT`** use the same **scaled** quad, **origin**, and **rotation** as **`SPRITE.DRAW`** (see `runtime/sprite/sprite_hit_cgo.go`). See `docs/DIRECTIVE_IMPLEMENTATION_TRACEABILITY.md` and `docs/reference/UNIVERSAL_HANDLE_METHODS.md`.

**Action Required:**
1. Audit ALL handle types against this matrix
2. Implement missing universal methods
3. Document type-specific methods clearly
4. Ensure getter/setter symmetry (`.pos()` gets, `.pos(x,y,z)` sets)

---

## Part 3: API Ergonomics - Reduce Argument Count

### ❌ BEFORE: Long Argument Lists

```basic
' Hard to read, hard to remember parameter order
CAMERA.SETUP(cam, 0, 10, 20, 0, 0, 0, 0, 1, 0, 45, TRUE)
LIGHT.CREATESPOT(255, 200, 150, 1.0, 0, 5, 0, 0, -1, 0, 30, 45)
SPRITE.CREATE(texHandle, 100, 200, 32, 32, 255, 255, 255, 255, 1.0, 1.0, 0.0)
```

### ✅ AFTER: Builder Pattern with Method Chaining

```basic
' Clear, readable, self-documenting
cam = CAMERA.CREATE()
    .pos(0, 10, 20)
    .look(0, 0, 0)
    .fov(45)
    
light = LIGHT.CREATESPOT()
    .pos(0, 5, 0)
    .dir(0, -1, 0)
    .col(255, 200, 150)
    .intensity(1.0)
    .cone(30, 45)
    
sprite = SPRITE.CREATE(texHandle)
    .pos(100, 200)
    .size(32, 32)
    .col(255, 255, 255)
    .alpha(1.0)
```

**Design Principles:**

1. **Required params in CREATE():** Only params that MUST be set (usually just asset path or type)
2. **Optional params via methods:** Everything else via chainable setters
3. **Sensible defaults:** Objects work immediately after CREATE() with good defaults
4. **Method chaining:** All setters return `self` for chaining

**Action Required:**
1. Identify all functions with 5+ parameters
2. Refactor to CREATE + chainable methods
3. Maintain backward compatibility with deprecated multi-arg versions
4. Document chaining pattern in STYLE_GUIDE.md

---

## Part 4: Namespace Organization

### Standard Namespace Structure

Every namespace should follow this pattern:

```basic
NAMESPACE.CREATE()          ' Create with defaults
NAMESPACE.CREATE<Type>()    ' Create specific type
NAMESPACE.LOAD(path)        ' Load from file
NAMESPACE.FREE(handle)      ' Free resource

NAMESPACE.SETPOS(h, x, y, z)     ' Canonical setter
NAMESPACE.SETROT(h, p, y, r)     ' Canonical setter  
NAMESPACE.SETSCALE(h, sx, sy, sz)' Canonical setter
NAMESPACE.SETCOL(h, r, g, b)     ' Canonical setter

NAMESPACE.GETPOS(h)         ' Canonical getter
NAMESPACE.GETROT(h)         ' Canonical getter
NAMESPACE.GETSCALE(h)       ' Canonical getter
```

### Namespace Examples

```basic
' ===== CAMERA Namespace =====
CAMERA.CREATE()                        ' Default 3D camera
CAMERA.CREATE2D()                      ' 2D camera
CAMERA.CREATEORTHOGRAPHIC()            ' Orthographic camera
CAMERA.FREE(cam)

CAMERA.SETPOS(cam, x, y, z)
CAMERA.SETROT(cam, pitch, yaw, roll)   ' NEW - was missing
CAMERA.SETFOV(cam, degrees)
CAMERA.SETTARGET(cam, x, y, z)
CAMERA.SETUP(cam, upX, upY, upZ)

CAMERA.GETPOS(cam)                     ' Returns Vec3 or array
CAMERA.GETROT(cam)                     ' NEW - returns rotation
CAMERA.GETFOV(cam)

' ===== MODEL Namespace =====
MODEL.LOAD(path)                       ' From file
MODEL.CREATECUBE(w, h, d)             ' Was MAKECUBE
MODEL.CREATECAPSULE(r, h)             ' Was MAKECAPSULE
MODEL.CREATESPHERE(r, rings, slices)  ' Was MAKESPHERE
MODEL.FREE(model)

MODEL.SETPOS(model, x, y, z)
MODEL.SETROT(model, pitch, yaw, roll)  ' NEW - implement this
MODEL.SETSCALE(model, sx, sy, sz)      ' NEW - implement this
MODEL.SETCOL(model, r, g, b)
MODEL.SETALPHA(model, a)

MODEL.GETPOS(model)
MODEL.GETROT(model)                    ' NEW
MODEL.GETSCALE(model)                  ' NEW

' ===== LIGHT Namespace =====
LIGHT.CREATE()                         ' Default directional
LIGHT.CREATEPOINT()                    ' Was MAKEPOINT
LIGHT.CREATEDIRECTIONAL()              ' Was MAKEDIRECTIONAL  
LIGHT.CREATESPOT()                     ' Was MAKESPOT
LIGHT.FREE(light)

LIGHT.SETPOS(light, x, y, z)           ' For point/spot
LIGHT.SETDIR(light, x, y, z)           ' For directional
LIGHT.SETCOL(light, r, g, b)
LIGHT.SETINTENSITY(light, value)
LIGHT.SETSHADOWS(light, enabled)       ' NEW

LIGHT.GETPOS(light)
LIGHT.GETDIR(light)
LIGHT.GETINTENSITY(light)

' ===== BODY3D Namespace =====
BODY3D.CREATE()                        ' Default dynamic
BODY3D.CREATESTATIC()                  ' Static body
BODY3D.CREATEKINEMATIC()               ' Kinematic body
BODY3D.FREE(body)

BODY3D.SETPOS(body, x, y, z)
BODY3D.SETROT(body, pitch, yaw, roll)  ' NEW - add quaternion support
BODY3D.SETSCALE(body, sx, sy, sz)      ' NEW - scales collision shape
BODY3D.SETVEL(body, vx, vy, vz)
BODY3D.SETMASS(body, mass)

BODY3D.GETPOS(body)
BODY3D.GETROT(body)                    ' NEW
BODY3D.GETVEL(body)
BODY3D.GETMASS(body)
```

**Action Required:**
1. Audit every namespace for missing methods
2. Add GET* methods for all SET* methods
3. Ensure CREATE/FREE symmetry
4. Document all methods in reference docs

---

## Part 5: Easy Mode Design

### Easy Mode as Convenience Layer

Easy Mode provides **global shortcuts** to Namespace.Method calls. These are **convenience only** and should:

1. Map 1:1 to Namespace.Method equivalents
2. Be documented as "shortcuts" not primary API
3. Use clear, Blitz3D-compatible names
4. NOT introduce new behavior (just syntax sugar)

### Easy Mode Standard Mappings

```basic
' ===== Creation Shortcuts =====
CreateCamera()          → CAMERA.CREATE()
CreateCamera2D()        → CAMERA.CREATE2D()
CreateLight()           → LIGHT.CREATE()
CreateModel(path)       → MODEL.LOAD(path)
CreateCube(w,h,d)       → MODEL.CREATECUBE(w,h,d)
CreateSprite(tex)       → SPRITE.CREATE(tex)
CreateBody()            → BODY3D.CREATE()

' ===== Spatial Shortcuts =====
PositionEntity(e,x,y,z) → ENTITY.SETPOS(e,x,y,z)
RotateEntity(e,p,y,r)   → ENTITY.SETROT(e,p,y,r)
ScaleEntity(e,sx,sy,sz) → ENTITY.SETSCALE(e,sx,sy,sz)
MoveEntity(e,f,r,u)     → ENTITY.MOVE(e,f,r,u)
TurnEntity(e,p,y,r)     → ENTITY.TURN(e,p,y,r)

EntityX(e)              → ENTITY.GETPOS(e)[0]
EntityY(e)              → ENTITY.GETPOS(e)[1]
EntityZ(e)              → ENTITY.GETPOS(e)[2]

' ===== Rendering Shortcuts =====
EntityColor(e,r,g,b)    → ENTITY.SETCOL(e,r,g,b)
EntityAlpha(e,a)        → ENTITY.SETALPHA(e,a)
ShowEntity(e)           → ENTITY.SHOW(e)
HideEntity(e)           → ENTITY.HIDE(e)

' ===== Cleanup Shortcuts =====
FreeEntity(e)           → ENTITY.FREE(e)
FreeCamera(c)           → CAMERA.FREE(c)
FreeModel(m)            → MODEL.FREE(m)
```

**Action Required:**
1. Implement Easy Mode as thin wrappers (no logic)
2. Update EASY_MODE.md to clearly state it's a convenience layer
3. Update all examples to prefer Namespace.Method
4. Keep Easy Mode for Blitz3D migration path

---

## Part 6: Documentation Standards

### STYLE_GUIDE.md (CREATE THIS FILE)

```markdown
# MoonBASIC Style Guide

## Official API Style

moonBASIC uses **Namespace.Method** as the canonical API style.

### ✅ Recommended Style

```basic
' Creation
camera = CAMERA.CREATE()
model = MODEL.LOAD("assets/player.glb")
light = LIGHT.CREATEPOINT()

' Configuration (prefer method chaining)
camera.pos(0, 10, 20)
      .look(0, 0, 0)
      .fov(60)

' Alternative: Individual calls
CAMERA.SETPOS(camera, 0, 10, 20)
CAMERA.SETTARGET(camera, 0, 0, 0)
CAMERA.SETFOV(camera, 60)

' Cleanup
CAMERA.FREE(camera)
MODEL.FREE(model)
```

### ⚠️ Acceptable for Migration

```basic
' Easy Mode (Blitz3D compatibility)
camera = CreateCamera()
PositionCamera(camera, 0, 10, 20)
CameraFOV(camera, 60)
FreeCamera(camera)
```

### ❌ Avoid

```basic
' Inconsistent mixing
camera = CreateCamera()         ' Easy Mode
CAMERA.SETPOS(camera, 0, 10, 20)  ' Namespace.Method
FreeCamera(camera)              ' Easy Mode again
```

## Naming Conventions

- **Variables:** `camelCase` (e.g., `playerModel`, `mainCamera`)
- **Constants:** `SCREAMING_SNAKE_CASE` (e.g., `MAX_SPEED`, `GRAVITY`)
- **Types:** `PascalCase` (e.g., `PlayerData`, `EnemyStats`)
- **Commands:** `Namespace.Method` or `NAMESPACE.METHOD` (case-insensitive)

## Function Design

### Prefer Chainable Methods

```basic
' ✅ Good - Clear, readable
light = LIGHT.CREATESPOT()
    .pos(10, 20, 5)
    .dir(0, -1, 0)
    .col(255, 200, 150)
    .intensity(2.0)
    .cone(25, 45)

' ❌ Avoid - Hard to read
light = LIGHT.CREATESPOT(10, 20, 5, 0, -1, 0, 255, 200, 150, 2.0, 25, 45)
```

### Keep Functions Focused

```basic
' ✅ Good - Single responsibility
CAMERA.SETPOS(cam, x, y, z)
CAMERA.SETTARGET(cam, tx, ty, tz)
CAMERA.SETFOV(cam, degrees)

' ❌ Avoid - Does too much
CAMERA.SETUP(cam, x, y, z, tx, ty, tz, ux, uy, uz, fov, perspective)
```
```

### Reference Documentation Template

Every namespace needs a reference doc following this template:

```markdown
# NAMESPACE Reference

## Overview

Brief description of what this namespace handles.

## Creation

### NAMESPACE.CREATE()
Create a default [object type].

**Syntax:**
```basic
handle = NAMESPACE.CREATE()
```

**Returns:** Handle to the new object

**Example:**
```basic
obj = NAMESPACE.CREATE()
```

### NAMESPACE.CREATE<Type>()
Create a specific type of [object].

**Syntax:**
```basic
handle = NAMESPACE.CREATE<Type>(params...)
```

**Parameters:**
- `param1` (type) - Description

**Returns:** Handle to the new object

**Example:**
```basic
obj = NAMESPACE.CREATETYPE(value1, value2)
```

### NAMESPACE.LOAD(path)
Load [object] from file.

**Syntax:**
```basic
handle = NAMESPACE.LOAD(filepath$)
```

**Parameters:**
- `filepath` (string) - Path to asset file

**Returns:** Handle to the loaded object

**Example:**
```basic
obj = NAMESPACE.LOAD("assets/model.glb")
```

## Universal Methods

### .pos(x, y, z)
Set or get position.

**Syntax:**
```basic
obj.pos(x#, y#, z#)        ' Set position
position = obj.pos()       ' Get position
```

**Parameters (setter):**
- `x` (float) - X coordinate
- `y` (float) - Y coordinate  
- `z` (float) - Z coordinate

**Returns (getter):** Array [x, y, z] or Vec3

**Example:**
```basic
obj.pos(0, 10, 5)
pos = obj.pos()
Print "X: " + pos[0]
```

### .rot(pitch, yaw, roll)
Set or get rotation.

**Syntax:**
```basic
obj.rot(pitch#, yaw#, roll#)  ' Set rotation
rotation = obj.rot()          ' Get rotation
```

**Parameters (setter):**
- `pitch` (float) - X-axis rotation in degrees
- `yaw` (float) - Y-axis rotation in degrees
- `roll` (float) - Z-axis rotation in degrees

**Returns (getter):** Array [pitch, yaw, roll]

**Example:**
```basic
obj.rot(0, 45, 0)
rot = obj.rot()
```

[Continue for .scale(), .col(), .alpha(), .free()]

## Type-Specific Methods

[Document unique methods for this namespace]

## Easy Mode Shortcuts

- `CreateObject()` → `NAMESPACE.CREATE()`
- `PositionObject(obj, x, y, z)` → `NAMESPACE.SETPOS(obj, x, y, z)`
- `FreeObject(obj)` → `NAMESPACE.FREE(obj)`

[List all Easy Mode mappings]

## See Also

- [Related namespace 1]
- [Related namespace 2]
```

**Action Required:**
1. Create STYLE_GUIDE.md immediately
2. Update ALL reference docs to follow template
3. Ensure every documented method has example code
4. Cross-reference related namespaces

---

## Part 7: Implementation Checklist

**v0.9 status:** This checklist is **implemented in-tree** (manifest policy tests, runtime handle dispatch, compiler/LSP deprecations, `STYLE_GUIDE.md`, reference docs and migration guides, LSP/snippet ordering, migrated `examples/`). Evidence: `docs/DIRECTIVE_IMPLEMENTATION_TRACEABILITY.md`, `docs/DIRECTIVE_CONFORMANCE_REPORT.md`. Ongoing work: new namespaces and long-arg APIs follow the same patterns; v1.0 removes deprecated aliases (Part 8).

### Phase 1: Manifest Changes (Week 1)

- [x] Rename all `*.MAKE` → `*.CREATE` in `compiler/builtinmanifest/commands.json`
- [x] Add `*.MAKE` as deprecated aliases (emit warnings)
- [x] Add missing `SETPOSITION` aliases for consistency (DECAL, PARTICLE, etc.)
- [x] Add missing `GET*` methods for all `SET*` methods
- [x] Remove duplicate AUDIO registrations
- [x] Run `go run ./tools/apidoc` to regenerate API_CONSISTENCY.md
- [x] Verify all changes compile

### Phase 2: Runtime Implementation (Week 2-3)

- [x] Implement missing universal methods (`.rot()`, `.scale()` for all spatial objects)
- [x] Add method chaining support (return `self` from setters)
- [x] Implement getter methods for all spatial properties
- [x] Refactor long-argument functions to builder pattern
- [x] Add deprecation warnings for MAKE commands
- [x] Add deprecation warnings for multi-arg legacy functions
- [x] Test all changes with existing examples

### Phase 3: Documentation (Week 4)

- [x] Create STYLE_GUIDE.md
- [x] Update API_CONVENTIONS.md with new standards
- [x] Update EASY_MODE.md to clearly mark as convenience layer
- [x] Rewrite all namespace reference docs using template
- [x] Update LANGUAGE.md to clarify case insensitivity
- [x] Update all code examples to use new patterns
- [x] Create migration guide for existing users

### Phase 4: Tooling (Week 5)

- [x] Update VSCode extension autocomplete to prefer CREATE over MAKE
- [x] Update VSCode extension to suggest universal methods
- [x] Add linter rules for style guide compliance
- [x] Update syntax highlighting for new patterns
- [x] Add snippets for common patterns (chainable creation)

### Phase 5: Example Migration (Week 6)

- [x] Update all examples/ to use Namespace.Method style
- [x] Update all examples/ to use CREATE instead of MAKE
- [x] Update all examples/ to use method chaining where appropriate
- [x] Ensure examples are consistent with style guide
- [x] Add comments explaining patterns for learners

---

## Part 8: Breaking Change Timeline

### Immediate (v0.9)
- ✅ Add all missing methods
- ✅ Add deprecation warnings
- ✅ Update documentation
- ⚠️ Both MAKE and CREATE work (with warnings)

### Minor Version (v0.10)
- ⚠️ MAKE commands emit loud deprecation warnings
- ⚠️ Multi-arg legacy functions emit warnings
- ✅ All examples use new style
- ✅ Style guide is normative

### Major Version (v1.0)
- ❌ Remove MAKE entirely (CREATE only)
- ❌ Remove legacy multi-arg functions
- ❌ Remove SETPOSITION (SETPOS only)
- ✅ Clean, consistent API

---

## Part 9: Success Criteria

This standardization effort is complete when:

**v0.9 verification:** See `docs/DIRECTIVE_CONFORMANCE_REPORT.md` and `docs/DIRECTIVE_IMPLEMENTATION_TRACEABILITY.md` for the signed-off snapshot (`go test ./...`, manifest policy tests, apidoc, cmdaudit).

1. ✅ **Naming Consistency:** All creation uses CREATE, all position uses SETPOS
2. ✅ **Universal Methods:** Every spatial handle has .pos(), .rot(), .scale() where applicable
3. ✅ **Documentation:** Every namespace has complete reference docs following template
4. ✅ **Examples:** All examples demonstrate best practices from style guide
5. ✅ **Tooling:** VSCode extension guides users toward canonical API
6. ✅ **Ergonomics:** No function has more than 4 parameters (use chaining)
7. ✅ **Symmetry:** Every SET has a GET, every CREATE has a FREE
8. ✅ **Easy Mode:** Clearly documented as convenience layer, not primary API

---

## Questions & Escalation

**Questions?** Contact language design team.

**Blockers?** Escalate immediately - this is foundation for v1.0.

**Timeline Issues?** Adjust phases but maintain quality standards.

---

## Appendix A: Quick Reference Card

```basic
' ============================================
' MoonBASIC API Quick Reference
' ============================================

' ----- CREATION -----
object = NAMESPACE.CREATE()           ' Default create
object = NAMESPACE.CREATE<Type>()     ' Specific type
object = NAMESPACE.LOAD(path)         ' From file

' ----- UNIVERSAL SPATIAL METHODS -----
object.pos(x, y, z)                   ' Set position
object.rot(pitch, yaw, roll)          ' Set rotation
object.scale(sx, sy, sz)              ' Set scale
object.col(r, g, b, a)                ' Set color
object.alpha(a)                       ' Set alpha
object.free()                         ' Free resource

' ----- GETTERS -----
pos = object.pos()                    ' Get position [x,y,z]
rot = object.rot()                    ' Get rotation [p,y,r]
scale = object.scale()                ' Get scale [sx,sy,sz]

' ----- METHOD CHAINING -----
cam = CAMERA.CREATE()
    .pos(0, 10, 20)
    .look(0, 0, 0)
    .fov(60)

' ----- NAMESPACE FORM -----
NAMESPACE.SETPOS(obj, x, y, z)        ' Also available
NAMESPACE.SETROT(obj, p, y, r)
NAMESPACE.SETSCALE(obj, sx, sy, sz)

' ----- EASY MODE (LEGACY) -----
obj = CreateObject()                  ' Maps to NAMESPACE.CREATE()
PositionObject(obj, x, y, z)          ' Maps to NAMESPACE.SETPOS()
FreeObject(obj)                       ' Maps to NAMESPACE.FREE()
```

---

**END OF DIRECTIVE**

*This directive represents the official design decisions for moonBASIC API standardization. Implementation should follow the phases outlined, maintaining backward compatibility during the transition period.*
