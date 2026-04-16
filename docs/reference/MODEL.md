# Model — `MODEL.*`

A **model** is a Raylib **`Model`**: meshes, materials, optional animation data, and a root transform. Registry keys use **dots and uppercase** (e.g. `MODEL.LOAD`). This page matches the **current** runtime, not legacy PascalCase-only docs.

**Conventions:** [STYLE_GUIDE.md](../../STYLE_GUIDE.md), [API_CONVENTIONS.md](API_CONVENTIONS.md) — reference pages use uppercase **`NAMESPACE.ACTION`**; Easy Mode (`Model.Load`, …) is [compatibility only](../../STYLE_GUIDE.md#easy-mode-compatibility-layer).

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — see [WAVE.md](WAVE.md) (registry-first headings, **Full Example** at the end).

**Requires CGO.** See **[MESH.md](MESH.md)** for raw **`MESH.*`** geometry.

---

### `MODEL.LOAD(path)`
Loads a 3D model file (glTF, GLB, OBJ, IQM, B3D). Returns a **model handle**.

### `MODEL.MAKE(mesh)` / `MODEL.CREATE(mesh)`
Builds a model from an existing **`Mesh`** handle. The model takes ownership of the mesh GPU data. Prefer registry **`MODEL.CREATE`** (canonical); **`MODEL.MAKE`** is a deprecated alias.

---

### `MODEL.DRAW(handle)`
Draws the model using its root transform. Call between **`RENDER.BEGIN3D(cam)`** and **`RENDER.END3D()`** (active 3D camera **`cam`**) for 3D rendering.

### `MODEL.SETPOS(handle, x, y, z)` (canonical; deprecated `MODEL.SETPOSITION`)
Sets the model's root transform to a specific world position.

### `MODEL.SETROT(handle, pitch, yaw, roll)`
Sets the model's absolute Euler rotation in **radians**.

### `MODEL.SETSCALE(handle, sx, sy, sz)`
Sets the non-uniform scale of the model.

---

### `MODEL.SETMATERIAL(handle, index, mat)`
Replaces a specific material slot in the model with a **`Material`** handle.

### `MODEL.FREE(handle)`
Unloads the model and its associated meshes/materials from memory and frees the heap slot.

---

## Full Example (load and draw)

```basic
WINDOW.OPEN(1280, 720, "Model Example")
WINDOW.SETFPS(60)
cam = CAMERA.CREATE()

mdl = MODEL.LOAD("assets/character.glb")
MODEL.SETPOS(mdl, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(20, 20, 20)
    RENDER.BEGIN3D(cam)
        MODEL.DRAW(mdl)
    RENDER.END3D()
    RENDER.FRAME()
WEND

MODEL.FREE(mdl)
WINDOW.CLOSE()
```

---

## Common mistakes

- **`MODEL.DRAW(mdl, matrix)`** — not supported; use **`MODEL.SETPOS`** (canonical) or deprecated **`MODEL.SETPOSITION`**, **`SETMATRIX`**, **`DRAWAT`**.
- **`mod` as a variable name** — **`MOD`** is reserved in moonBASIC; use **`mdl`** or **`modelHandle`**.
- **Double-free after `MODEL.CREATE` (mesh → model)** — **`MODEL.MAKE`** is deprecated with the same arity; follow **`MODEL.FREE`** then **`MESH.FREE`** (mesh slot only) as in the test, or read **`consumedByModel`** behaviour above.

---

## See also

- [ANIMATION_3D.md](ANIMATION_3D.md) — skeletal clips: **`MODEL.*`** vs **`ENTITY.*`**
- [MESH.md](MESH.md) — procedural meshes, **`MESH.UPLOAD`**, **`MESH.DRAW`**
- [CAMERA.md](CAMERA.md) — 3D camera
- [LIGHT.md](LIGHT.md) — PBR lighting
- [SHADER.md](SHADER.md) — custom materials via shaders
