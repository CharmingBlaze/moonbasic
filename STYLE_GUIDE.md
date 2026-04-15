# MoonBASIC Style Guide

## Official API Style

MoonBASIC uses `Namespace.Method` as the canonical API style.

### Recommended Style

```basic
' Creation
camera = CAMERA.CREATE()
model = MODEL.LOAD("assets/player.glb")
light = LIGHT.CREATEPOINT()

' Configuration (prefer method chaining)
camera.pos(0, 10, 20)
      .look(0, 0, 0)
      .fov(60)

' Alternative: namespace calls
CAMERA.SETPOS(camera, 0, 10, 20)
CAMERA.SETTARGET(camera, 0, 0, 0)
CAMERA.SETFOV(camera, 60)

' Cleanup
CAMERA.FREE(camera)
MODEL.FREE(model)
```

### Migration Style (Easy Mode)

```basic
camera = CreateCamera()
PositionEntity(camera, 0, 10, 20)
FreeCamera(camera)
```

Easy Mode is a compatibility and convenience layer only. New docs and examples should use canonical namespace commands first.

## Naming Conventions

- Variables: `camelCase` (for example: `playerModel`, `mainCamera`)
- Constants: `SCREAMING_SNAKE_CASE` (for example: `MAX_SPEED`)
- Types: `PascalCase` (for example: `PlayerData`)
- Commands: `Namespace.Method` in docs; command lookup remains case-insensitive

## API Design Rules

- Use `CREATE` for object construction (`MAKE` remains temporary deprecated alias).
- Prefer `SETPOS` as canonical position setter (`SETPOSITION` remains temporary alias).
- Keep method signatures short and focused.
- Prefer chainable setters over long multi-argument constructor calls.
- Keep symmetry: every `SET*` should have a matching `GET*`; every `CREATE` should have a matching `FREE`.

## Universal Spatial Methods

Spatial handles should expose these methods where applicable:

- `.pos()` / `.pos(x, y, z)`
- `.rot()` / `.rot(pitch, yaw, roll)` (or `.rot(angle)` for 2D)
- `.scale()` / `.scale(sx, sy, sz)`
- `.col()` / `.col(r, g, b[, a])` for renderable handles
- `.alpha()` / `.alpha(a)` for renderable handles
- `.free()`

## Easy Mode Mapping Rule

Easy Mode wrappers must be thin aliases to canonical namespace methods and must not introduce behavior differences.

## Long-argument APIs and builders

Prefer **`NAMESPACE.CREATE()`** with defaults, then **handle or namespace setters** (ideally chainable) instead of single calls with many positional parameters. Legacy multi-argument overloads may remain during migration; the compiler can flag deprecated aliases (see `--strict-deprecated` on `moonbasic --check`). New APIs should cap rarely-needed optional configuration at **four or fewer positional parameters** and move extras to named/chainable methods.

