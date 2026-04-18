# Model Commands

Load, position, and draw 3D models (glTF, GLB, OBJ, IQM, B3D) with materials and transforms.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Load with `MODEL.LOAD` or create from a mesh with `MODEL.CREATE`.
2. Position with `MODEL.SETPOS`, rotate with `MODEL.SETROT`, scale with `MODEL.SETSCALE`.
3. Draw between `RENDER.BEGIN3D` / `RENDER.END3D` with `MODEL.DRAW`.
4. Free with `MODEL.FREE`.

For raw meshes see [MESH.md](MESH.md). For skeletal animation see [ANIMATION_3D.md](ANIMATION_3D.md).

---

### `MODEL.LOAD(path)`
Loads a 3D model file (glTF, GLB, OBJ, IQM, B3D).

- **Arguments**:
    - `path`: (String) File path.
- **Returns**: (Handle) The new model handle.
- **Example**:
    ```basic
    hero = MODEL.LOAD("hero.glb")
    ```

---

### `MODEL.CREATE(mesh)`
Builds a model from an existing **`Mesh`** handle.

- **Returns**: (Handle) The new model handle.

---

### `MODEL.DRAW(handle)`
Draws the model using its current transform.

- **Returns**: (Handle) The model handle (for chaining).

---

### `MODEL.SETPOS(handle, x, y, z)` / `SETROT` / `SETSCALE`
Sets the model's world position, Euler rotation (radians), or scale.

- **Returns**: (Handle) The model handle (for chaining).

---

### `MODEL.SETMATERIAL(handle, index, material)`
Replaces a material slot in the model.

- **Returns**: (Handle) The model handle (for chaining).

---

### `MODEL.FREE(handle)`
Unloads the model and its resources.

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
