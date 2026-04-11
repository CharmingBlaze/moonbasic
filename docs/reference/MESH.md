# Mesh — `MESH.*`

**CPU mesh geometry** (Raylib `Mesh`) stored in the VM heap (`TypeName` **`Mesh`**, `TagMesh`). Use **`MESH.Upload`** to send vertex data to the GPU, then **`MESH.Draw`** or **`MESH.DrawRotated`** with a **`MATERIAL.*`** handle. Meshes are **not** full **`MODEL.*`** objects (no bundled materials or hierarchy).

**Requires CGO** (same as **`MODEL.*`**, **`Camera.*`**, **`MATERIAL.*`**).

Registry keys use **dots and uppercase** (e.g. `MESH.MakeCube`).

---

### `Mesh.MakeCube(w, h, d)`
Creates a procedural box mesh. Returns a **mesh handle**.

### `Mesh.MakeSphere(r, rings, slices)`
Creates a procedural sphere mesh.

### `Mesh.Load(path)`
Loads the first submesh from a model file. Returns a **mesh handle**.

---

### `Mesh.Upload(mesh, dynamic)`
Uploads mesh data to the GPU. Set `dynamic` to `TRUE` if the mesh will be updated frequently.

### `Mesh.Draw(mesh, material, matrix)`
Draws a mesh using a `Material` handle and a `Transform` matrix handle.

### `Mesh.Free(handle)`
Unloads the mesh from the GPU and frees the handle from memory.

---

### `Mesh.VertexCount(handle)` / `Mesh.TriangleCount(handle)`
Returns the number of vertices or triangles in the mesh.

### `Mesh.GetBounds(handle)`
Returns a 6-element float array handle `[minX, minY, minZ, maxX, maxY, maxZ]`.

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
