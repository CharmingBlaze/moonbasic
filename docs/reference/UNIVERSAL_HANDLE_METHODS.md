# Universal handle methods

**Policy:** [API Standardization Directive](../API_STANDARDIZATION_DIRECTIVE.md) (Part 2).

moonBASIC maps **handle** method calls (`myHandle.pos(...)`, `myHandle.free()`, …) to **`NAMESPACE.COMMAND`** builtins. Normalization and dispatch live in [`vm/handlecall.go`](../../vm/handlecall.go) (`normalizeHandleMethod`, `handleCallDispatch`, `handleCallBuiltin`, `HandleCallSuggestions`). With **zero arguments**, universal pose names (e.g. `pos()`, `rot()`, `scale()`) resolve to **`GET*`** builtins where implemented (e.g. `MODEL.GETPOS`); with arguments they resolve to **`SET*`**.

## Core Workflow

- **Zero-arg getter:** `handle.pos()` → `NAMESPACE.GETPOS(handle)`.
- **With-args setter:** `handle.pos(x, y, z)` → `NAMESPACE.SETPOS(handle, x, y, z)`.
- **Free:** `handle.free()` → `NAMESPACE.FREE(handle)` for any heap type.
- **Add new handle type:** Update `handleCallBuiltin` in `vm/handlecall.go`, register in `runtime/*/register*.go`, extend `commands.json`.

## Mapping rules

- Method names are **case-insensitive** in source; dispatch uses uppercase **`TYPE.ACTION`** registry keys. Common method aliases are folded (for example `pos`, `position`, `setposition` → `SETPOS` for spatial types where registered). **`SPRITE.HIT`** and **`SPRITE.POINTHIT`** use the same **scaled** quad, **origin**, and **rotation** as **`SPRITE.DRAW`** (raylib **`DrawTexturePro`**), not a separate axis-aligned box. Handle calls **`spr.hit(other)`** / **`spr.collide(other)`** dispatch to **`SPRITE.HIT`**; **`spr.pointHit(x, y)`** → **`SPRITE.POINTHIT`**.
- Each heap **tag** (camera, model, sprite, light, body, …) defines which methods are valid and which registry key is invoked. Examples: zero-arg **`cam.rot()`** → **`CAMERA.GETROT`** (euler radians, same order as **`MODEL.GETROT`**); **`entity.rot()`** / **`entity.scale()`** → **`ENTITY.GETROT`** / **`ENTITY.GETSCALE`** (pitch/yaw/roll and XYZ scale from the entity store); **`entity.col()`** / **`entity.alpha()`** → **`ENTITY.GETCOLOR`** (RGBA 0–255 in the array’s fourth slot) / **`ENTITY.GETALPHA`** (0.0–1.0); **`hero.rot()`** (character handle) → **`CHARACTERREF.GETROT`** (approximate euler from **linear velocity**; no capsule orientation in Jolt’s API); **`agent.rot()`** (nav agent) → **`NAVAGENT.GETROT`** (euler from **waypoint tangent** or **steering velocity**); 3D Jolt **`body.scale()`** → **`BODY3D.GETSCALE`** (**[sx,sy,sz]** factors for **box/sphere/capsule** primitives; **mesh** bodies report **1,1,1**; **`body.scale(x,y,z)`** → **`BODY3D.SETSCALE`**); 2D **`spr.rot()`** / **`spr.scale()`** / **`spr.col()`** / **`spr.alpha()`** → **`SPRITE.GETROT`** (**[0,0,roll]** radians) / **`SPRITE.GETSCALE`** (**[sx,sy,1]**) / **`SPRITE.GETCOLOR`** / **`SPRITE.GETALPHA`** (tint used by **`DrawTexturePro`**); 2D physics **`body.rot()`** → **`BODY2D.GETROT`** (scalar angle; **`body.rot(rad)`** → **`BODY2D.SETROT`**); particle **`emitter.col()`** / **`emitter.alpha()`** → **`PARTICLE.GETCOLOR`** (start RGBA 0–255) / **`PARTICLE.GETALPHA`** (start A as 0.0–1.0); light **`light.col()`** → **`LIGHT.GETCOLOR`** ([r,g,b,a] 0–255 floats); 2D light **`light2d.col()`** → **`LIGHT2D.GETCOLOR`**; instanced model **`inst.pos()`** / **`rot()`** / **`scale()`** → **`INSTANCE.GETPOS`** / **`GETROT`** / **`GETSCALE`** (transform of **instance index 0**; use **`MODEL.SETINSTANCEPOS`** / **`INSTANCE.SETROT`** with an index for other slots).

## Directive alignment

The API standardization directive requires **`.pos` / `.rot` / `.scale`** (where applicable), **`.col` / `.alpha`** for renderables, and **`.free()`** for heap objects. Coverage depends on:

1. **Manifest** entries for `NAMESPACE.SET*` / `GET*` / `FREE`.
2. **`handleCallBuiltin`** entries for the handle tag so script `.method` resolves.
3. **Runtime** implementations backing those registry keys.

When adding a new handle type, update **`handleCallBuiltin`** and **`HandleCallSuggestions`**, register the namespace commands in the appropriate **`runtime/*/register*.go`**, and extend **`compiler/builtinmanifest/commands.json`** so the compiler and LSP agree on arity and types.

## Full Example

Handle method chaining on entity, camera, and sprite.

```basic
WINDOW.OPEN(960, 540, "Handle Methods Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
cam.pos(0, 6, -10)        ; CAMERA.SETPOS
cam.target(0, 0, 0)

e = ENTITY.CREATECUBE(1.0)
e.pos(0, 1, 0)            ; ENTITY.SETPOS
e.col(80, 160, 255, 255)  ; ENTITY.SETCOLOR

t = 0.0
WHILE NOT WINDOW.SHOULDCLOSE()
    t = t + TIME.DELTA()
    e.rot(0, t * 45, 0)   ; ENTITY.SETROT

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW.GRID(10, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

e.free()                  ; ENTITY.FREE
cam.free()                ; CAMERA.FREE
WINDOW.CLOSE()
```

---

## See also

- [API_CONVENTIONS.md](API_CONVENTIONS.md)
- [STYLE_GUIDE.md](../../STYLE_GUIDE.md)
- [API_CONSISTENCY.md](../API_CONSISTENCY.md) (regenerate with `go run ./tools/apidoc`)
