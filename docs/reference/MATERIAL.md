# Material Commands

Material handles control how meshes are shaded: texture slots, shader assignment, color tint, UV scroll, and PBR parameters.

## Core Workflow

1. `MATERIAL.CREATE()` or `MATERIAL.CREATEDEFAULT()` — get a material handle.
2. `MATERIAL.SETTEXTURE(mat, slot, texHandle)` — bind a texture to a slot.
3. `MATERIAL.SETSHADER(mat, shaderHandle)` — optional custom shader.
4. Pass the handle to `MESH.DRAW(mesh, mat, transform)` each frame.
5. `MATERIAL.FREE(mat)` when done.

---

## Creation

### `MATERIAL.CREATE()` 

Creates a blank material with no textures or shader assigned. Returns a **material handle**.

---

### `MATERIAL.CREATEDEFAULT()` 

Creates a material pre-loaded with the engine default PBR shader. Returns a **material handle**. Use this as a starting point for most opaque meshes.

---

### `MATERIAL.CREATEPBR()` 

Creates a material configured for the PBR (physically-based rendering) path. Returns a **material handle**.

---

### `MATERIAL.FREE(mat)` 

Releases the material handle and frees GPU-side references (does not unload bound textures).

- *Handle shortcut*: `mat.free()`

---

## Textures & Shader

### `MATERIAL.SETTEXTURE(mat, slot, texHandle)` 

Binds a texture to a numbered slot on the material. Slot `0` = albedo/diffuse, `1` = normal map, `2` = metalness/roughness (PBR). `texHandle` from `TEXTURE.LOAD`.

- *Handle shortcut*: `mat.setTexture(slot, texHandle)`

---

### `MATERIAL.SETSECONDARYTEXTURE(entityId, texHandle)` 

Sets the secondary (detail) texture on an entity's material 0. Alias of `ENTITY.SETDETAILTEXTURE`.

---

### `MATERIAL.SETSHADER(mat, shaderHandle)` 

Assigns a custom shader to the material. `shaderHandle` from `SHADER.LOAD`.

- *Handle shortcut*: `mat.setShader(shaderHandle)`

---

## Color & Effects

### `MATERIAL.SETCOLOR(mat, mapIndex, r, g, b, a)` 

Sets the color tint for a specific material map. `mapIndex` `0` = diffuse, `1` = specular, `2` = normal, etc. Values 0–255.

- *Handle shortcut*: `mat.setColor(mapIndex, r, g, b, a)`

---

### `MATERIAL.SETFLOAT(mat, locIndex, value)` 

Sets a float uniform on the material shader at `locIndex`.

---

### `MATERIAL.SETEFFECT(mat, effectName)` 

Applies a named post-process or render effect to this material (e.g. `"cel"`, `"pbr_lit"`).

---

### `MATERIAL.SETEFFECTPARAM(mat, paramName, value)` 

Sets a named float parameter for the material's effect.

---

### `MATERIAL.SETUVSCROLL(entityId, u, v)` 

Scrolls the UV coordinates of mesh material 0 on the given entity by `(u, v)` per frame. Alias of `ENTITY.SCROLLMATERIAL`.

---

## Bulk & Auto

### `MATERIAL.BULKASSIGN(pattern, matHandle)` 

Assigns a material to all loaded entities whose name matches `pattern` (glob). Returns the count of entities affected.

---

### `MATERIAL.AUTOFILTER(value)` 

Sets the global texture filter mode used when auto-assigning materials. Accepts a filter constant.

---

## Full Example

A mesh with a textured PBR material and color tint.

```basic
WINDOW.OPEN(960, 540, "Material Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 2, -6)
CAMERA.SETTARGET(cam, 0, 0, 0)

sun = LIGHT.CREATE("directional")
LIGHT.SETDIR(sun, -0.5, -1, -0.3)
LIGHT.SETCOLOR(sun, 255, 240, 200, 255)

mesh = MESH.CREATECUBE(2, 2, 2)
tex  = TEXTURE.LOAD("assets/crate.png")
mat  = MATERIAL.CREATEPBR()
MATERIAL.SETTEXTURE(mat, 0, tex)
MATERIAL.SETCOLOR(mat, 0, 255, 200, 160, 255)    ; warm tint

t = TRANSFORM.IDENTITY()

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(20, 20, 30)
    RENDER.BEGIN3D(cam)
        MESH.DRAW(mesh, mat, t)
        DRAW3D.GRID(10, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

MATERIAL.FREE(mat)
TEXTURE.UNLOAD(tex)
MESH.UNLOAD(mesh)
LIGHT.FREE(sun)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `MATERIAL.MAKE(...)` | Deprecated alias of `MATERIAL.CREATE`. |
| `MATERIAL.MAKEDEFAULT()` | Create a default (white, unlit) material. Alias of `MATERIAL.DEFAULT`. |
| `MATERIAL.MAKEPBR(albedoTex, normalTex, roughness, metalness)` | Create a PBR material in one call. |

---

## See also

- [MESH.md](MESH.md) — mesh creation and drawing
- [TEXTURE.md](TEXTURE.md) — texture loading
- [SHADER.md](SHADER.md) — custom shader programs
- [LIGHT.md](LIGHT.md) — lighting for PBR materials
