# Model — `MODEL.*`

A **model** is a Raylib **`Model`**: meshes, materials, optional animation data, and a root transform. Registry keys use **dots and uppercase** (e.g. `MODEL.LOAD`). This page matches the **current** runtime, not legacy PascalCase-only docs.

**Requires CGO.** See **[MESH.md](MESH.md)** for raw **`MESH.*`** geometry.

---

### `Model.Load(path)`
Loads a 3D model file (glTF, GLB, OBJ, IQM, B3D). Returns a **model handle**.

### `Model.Make(mesh)`
Builds a model from an existing **`Mesh`** handle. The model takes ownership of the mesh GPU data.

---

### `Model.Draw(handle)`
Draws the model using its root transform. Call between **`Camera.Begin()`** and **`Camera.End()`** for 3D rendering.

### `Model.SetPos(handle, x, y, z)`
Sets the model's root transform to a specific world position.

### `Model.SetRot(handle, pitch, yaw, roll)`
Sets the model's absolute Euler rotation in **radians**.

### `Model.SetScale(handle, sx, sy, sz)`
Sets the non-uniform scale of the model.

---

### `Model.SetMaterial(handle, index, mat)`
Replaces a specific material slot in the model with a **`Material`** handle.

### `Model.Free(handle)`
Unloads the model and its associated meshes/materials from memory and frees the heap slot.

---

## Example (Load and Draw)

```basic
Window.Open(1280, 720, "Model Example")
mdl = Model.Load("assets/character.glb")
Model.SetPos(mdl, 0, 0, 0)

WHILE NOT Window.ShouldClose()
    Render.Clear(20, 20, 20)
    Camera.Begin(cam)
        Model.Draw(mdl)
    Camera.End()
    Render.Frame()
WEND

Model.Free(mdl)
Window.Close()
```

---

## Common mistakes

- **`MODEL.DRAW(mdl, matrix)`** — not supported; use **`MODEL.SETPOSITION`** / **`SETMATRIX`** / **`DRAWAT`**.
- **`mod` as a variable name** — **`MOD`** is reserved in moonBASIC; use **`mdl`** or **`modelHandle`**.
- **Double-free after `MODEL.MAKE`** — follow **`MODEL.FREE`** then **`MESH.FREE`** (mesh slot only) as in the test, or read **`consumedByModel`** behaviour above.

---

## See also

- [ANIMATION_3D.md](ANIMATION_3D.md) — skeletal clips: **`MODEL.*`** vs **`ENTITY.*`**
- [MESH.md](MESH.md) — procedural meshes, **`MESH.UPLOAD`**, **`MESH.DRAW`**
- [CAMERA.md](CAMERA.md) — 3D camera
- [LIGHT.md](LIGHT.md) — PBR lighting
- [SHADER.md](SHADER.md) — custom materials via shaders
