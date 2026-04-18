# Game Engine Patterns

Common high-level patterns for world interaction, material management, and lighting in the MoonBASIC engine.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. **Spatial Queries:** Use `RAY.HITMODEL` or `CAMERA.UNPROJECT` to interact with the 3D world.
2. **Materials:** Customize appearance with `MATERIAL.SETTEXTURE` and `TEXTURE.SETFILTER`.
3. **Lighting:** Setup the scene's mood with `LIGHT.CREATEPOINT` and `RENDER.SETAMBIENT`.
4. **2D Juice:** Drive animated UI and sprites with `SPRITE.PLAY` and `SPRITE.SETFRAME`.

---

## 1. Collision and Raycasting (3D)

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `CAMERA.UNPROJECT(cam, x, y)`| Handle, Float, Float| Array | Screen `[x, y]` to world `[x, y, z]`. |
| `RAY.HITMODEL(ray, model)` | Handle, Handle | Boolean | `TRUE` if ray hits the model. |
| `RAY.HITDISTANCE(ray, model)`| Handle, Handle | Float | Distance to the first hit point. |
| `ENTITY.DISTANCE(e1, e2)` | Integer, Integer | Float | 3D world distance. |

---

## 2. Texture and Material

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `MATERIAL.SETTEXTURE(m, t, type)`| Handle, Handle, Int| None | Binds texture to slot (diffuse, normal, etc). |
| `TEXTURE.SETFILTER(t, mode)`| Handle, Integer | None | Sets point vs. bilinear filtering. |
| `TEXTURE.SETWRAP(t, mode)` | Handle, Integer | None | Sets repeat vs. clamp tiling. |

---

## 3. Lights and Ambient

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `LIGHT.CREATEPOINT(x, y, z, r, g, b, e)`| Float... | Handle | Creates a point light source. |
| `LIGHT.SETTARGET(light, x, y, z)`| Handle, Float... | None | Aims directional/spot lights. |
| `RENDER.SETAMBIENT(r, g, b, a)`| Float... | None | Sets global scene base lighting. |

---

## Full Example

A simple scene demonstrating raycasting against a model and dynamic lighting.

```basic
WINDOW.OPEN(1280, 720, "Engine Patterns Demo")
cam = CAMERA.CREATE()
mdl = MODEL.LOAD("cube.glb").SETPOS(0, 0, 0)
lt = LIGHT.CREATEPOINT(5, 5, 5, 255, 200, 100, 10.0)

RENDER.SETAMBIENT(40, 40, 50, 255)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    
    ; 1. Raycasting: Check if mouse is over model
    ray = CAMERA.GETRAY(cam, INPUT.MOUSEX(), INPUT.MOUSEY())
    IF RAY.HITMODEL(ray, mdl)
        ENTITY.SETCOLOR(mdl, 255, 0, 0, 255)
    ELSE
        ENTITY.SETCOLOR(mdl, 255, 255, 255, 255)
    END IF

    ; 2. Rendering
    RENDER.CLEAR(10, 10, 20)
    RENDER.BEGIN3D(cam)
        MODEL.DRAW(mdl)
    RENDER.END3D()
    RENDER.FRAME()
WEND
```

## See also

- [ENTITY.md](ENTITY.md) — full entity API
- [MATERIAL.md](MATERIAL.md) — advanced shader and texture properties
- [CAMERA.md](CAMERA.md) — projection and viewport management
