# Mesh Commands

CPU mesh geometry: procedural creation, loading, GPU upload, and drawing.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create with `MESH.MAKECUBE`, `MESH.MAKESPHERE`, or `MESH.LOAD`.
2. Upload to GPU with `MESH.UPLOAD`.
3. Draw with `MESH.DRAW` using a material and transform handle.
4. Free with `MESH.FREE`.

Meshes are **not** full models — see [MODEL.md](MODEL.md) for assets with bundled materials.

---

### `MESH.MAKECUBE(w, h, d)` / `MAKESPHERE` / `MAKEPLANE`
Creates a procedural mesh.

- **Arguments**:
    - `w, h, d`: (Float) Dimensions.
- **Returns**: (Handle) The new mesh handle.
- **Example**:
    ```basic
    cube = MESH.MAKECUBE(2, 2, 2)
    ```

---

### `MESH.LOAD(path)`
Loads the first submesh from a model file.

- **Returns**: (Handle) The new mesh handle.

---

### `MESH.UPLOAD(mesh [, dynamic])`
Uploads mesh data to the GPU.

- **Arguments**:
    - `mesh`: (Handle) The mesh to upload.
    - `dynamic`: (Boolean, Optional) `TRUE` if the mesh will be updated frequently.
- **Returns**: (Handle) The mesh handle (for chaining).

---

### `MESH.DRAW(mesh, material, x, y, z)`
Draws a mesh with a specific material and position.

- **Arguments**:
    - `mesh`: (Handle) The mesh to draw.
    - `material`: (Handle) The material to apply.
    - `x, y, z`: (Float) World position.
- **Returns**: (Handle) The mesh handle (for chaining).

---

### `MESH.FREE(handle)`
Unloads the mesh from the GPU and frees the handle.

---

### `MESH.VERTEXCOUNT(handle)` / `TRIANGLECOUNT`
Returns the number of vertices or triangles in the mesh.

- **Returns**: (Integer)

---

### `MESH.GETBOUNDS(handle)`
Returns the bounding box of the mesh.

- **Returns**: (Handle) A 6-float array handle `[minX, minY, minZ, maxX, maxY, maxZ]`.

---

## Full Example

```basic
WINDOW.OPEN(800, 600, "Mesh Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 3, -5)
CAMERA.SETTARGET(cam, 0, 0, 0)

cube = MESH.MAKECUBE(2, 2, 2)
MESH.UPLOAD(cube)

mat = MATERIAL.DEFAULT()

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(30, 30, 50)
    RENDER.BEGIN3D(cam)
        MESH.DRAW(cube, mat, 0, 0, 0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

MESH.FREE(cube)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## Common mistakes

- **`Window.Open`** first — GPU init must be ready (see test).
- **Skipping `MESH.Upload`** — draw may fail or warn if the mesh is not uploaded.
- **`MESH.DrawRotated`** — **radians**, not degrees.
- **Confusing mesh and model** — **`MODEL.*`** loads whole assets; **`MESH.*`** is one mesh + your material + transform.

---

## See also

- [MODEL.md](MODEL.md) — full models and materials
- [IMAGE.md](IMAGE.md) — heightmap / cubicmap sources
- [CAMERA.md](CAMERA.md) — 3D camera
