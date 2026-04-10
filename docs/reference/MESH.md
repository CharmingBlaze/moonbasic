# Mesh — `MESH.*`

**CPU mesh geometry** (Raylib `Mesh`) stored in the VM heap (`TypeName` **`Mesh`**, `TagMesh`). Use **`MESH.UPLOAD`** to send vertex data to the GPU, then **`MESH.DRAW`** or **`MESH.DRAWROTATED`** with a **`MATERIAL.*`** handle. Meshes are **not** full **`MODEL.*`** objects (no bundled materials or hierarchy).

**Requires CGO** (same as **`MODEL.*`**, **`Camera.*`**, **`MATERIAL.*`**).

Registry keys use **dots and uppercase** (e.g. `MESH.MAKECUBE`).

---

### Mesh.MakeCube / procedural generators

```basic
m = MESH.MAKECUBE(width#, height#, length#)
m = MESH.MAKESPHERE(radius#, rings, slices)
m = MESH.MAKEPLANE(width#, length#, resX, resZ)
; … see table below
```

**Parameters** — Raylib **`GenMesh*`** rules apply (units are world units; planes use **resX/resZ** subdivisions).

**Returns** — mesh handle.

**Notes** — After **`Window.Open`**, call **`MESH.UPLOAD(m, dynamic?)`** before **`MESH.DRAW`** (Raylib uploads VBOs). **`dynamic?`** is **`TRUE`** if you will **`MESH.UPDATEVERTEX`** every frame.

| Command | Arguments |
|---|---|
| `MESH.MAKEPOLY` | `sides, radius` |
| `MESH.MAKEPLANE` | `width, length, resX, resZ` |
| `MESH.MAKECUBE` | `width, height, length` |
| `MESH.MAKESPHERE` | `radius, rings, slices` |
| `MESH.MAKECYLINDER` | `radius, height, slices` |
| `MESH.MAKECONE` | `radius, height, slices` |
| `MESH.MAKETORUS` | `radius, size, radSeg, sides` |
| `MESH.MAKEKNOT` | `radius, size, radSeg, sides` |
| `MESH.MAKEHEIGHTMAP` | `image, sizeX, sizeY, sizeZ` |
| `MESH.MAKECUBICMAP` | `image, cubeX, cubeY, cubeZ` |

### Legacy aliases

| Command | Same generator |
|---|---|
| `MESH.CUBE` | `MESH.MAKECUBE` |
| `MESH.SPHERE` | `MESH.MAKESPHERE` |
| `MESH.PLANE` | `MESH.MAKEPLANE` |
| `MESH.CREATECUBE` | `MESH.MAKECUBE` |
| `MESH.CREATESPHERE` | `MESH.MAKESPHERE` |
| `MESH.CREATEPLANE` | `MESH.MAKEPLANE` |

---

### Mesh.Upload

```basic
MESH.UPLOAD(m, dynamic?)
```

Uploads mesh data to the GPU (`UploadMesh`). **`dynamic?`** must be **boolean** — use **`TRUE`** for meshes you update per frame with **`MESH.UPDATEVERTEX`**.

> **Common mistake:** Passing **`0`/`1`** instead of **`FALSE`/`TRUE`** — the semantic checker expects a **bool** for the second argument.

---

### Mesh.Draw

```basic
MESH.DRAW(m, material, matrix)
```

**`matrix`** is a **`MAT4`** handle, or **`0`** for identity. Uses `DrawMesh` with the given **material** and **transform**.

**Phase** — Call inside **`CAMERA.BEGIN`** / **`CAMERA.END`** (3D mode) for correct projection.

---

### Mesh.DrawRotated

```basic
MESH.DRAWROTATED(m, material, rx#, ry#, rz#)
```

Builds **`MatrixRotateXYZ(Vector3{rx, ry, rz})`** — angles are **radians** (not degrees). **No** separate axis-angle overload in this binding.

---

### Mesh.UpdateVertex

```basic
MESH.UPDATEVERTEX(m, idx, x#, y#, z#, nx#, ny#, nz#, u#, v#)
```

Writes one vertex’s **position**, **normal**, and **UV**. Requires **`MESH.UPLOAD(m, TRUE)`** for dynamic meshes; after edits, Raylib buffer updates run when **`VaoID != 0`**.

---

### Mesh.VertexCount / Mesh.TriangleCount

```basic
n = MESH.VERTEXCOUNT(m)
t = MESH.TRIANGLECOUNT(m)
```

**Returns** — integers from the Raylib **`Mesh`** (after generation).

---

### Bounding box (axis-aligned, local)

```basic
mnx = MESH.GETBBOXMINX(m)
; … MESH.GETBBOXMAXZ(m)
```

Six accessors, each **`(m)`** → float component.

---

### Mesh.GenTangents (`MESH.GENTANGENTS`)

```basic
MESH.GENTANGENTS(m)
```

**Current builds:** With **CGO**, this command **errors at runtime** — the **raylib-go** tangent generator is not wired in the CGO `rmodels` path. Do not rely on it until the runtime exposes **`GenMeshTangents`** for CGO builds.

---

### Mesh.Free

```basic
MESH.FREE(m)
```

Unloads GPU and CPU mesh data (`UnloadMesh`). Idempotent object **`Free`** via heap; a second **`MESH.FREE`** on the same handle fails (stale handle).

**`MESH.LOAD(path)`** — loads the **first submesh** of a file via Raylib **`LoadModel`**; the runtime keeps the parent model alive until **`MESH.FREE`** (do not mix with manual unload of the same file elsewhere).

**`MESH.EXPORT(m, path)`** — writes **`ExportMesh`** (e.g. OBJ).

**`MESH.GETBOUNDS(m)`** — returns a **6-element float array** `[minX,minY,minZ,maxX,maxY,maxZ]` (caller should **`ERASE`** when done).

**`MESH.DRAWAT(m, material, x, y, z)`** — draws with **`MatrixTranslate`** (identity rotation).

**`MESH.DRAWINSTANCED(m, material, transforms_array, count)`** — **`transforms_array`** is a flat float array of **`count×16`** values (column-major **`MAT4`** per instance).

**`MESH.MAKECUSTOM(verts, indices)`** — **`verts`** length **`8×N`** (`x,y,z,nx,ny,nz,u,v` per vertex); **`indices`** length **`3×T`** (triangle list). **`MESH.MAKECAPSULE`** and **meshoptimizer**-style **`MESH.OPTIMISE*`** are not available in this binding (runtime returns a clear error).

---

## Example (upload, camera, draw)

See **`testdata/mesh_complete_test.mb`**: **`MESH.MAKECUBE`** → **`MESH.UPLOAD`** → **`MATERIAL.MAKEDEFAULT`** → **`CAMERA.BEGIN`** / **`MESH.DRAWROTATED`** / **`CAMERA.END`**.

---

## Common mistakes

- **`Window.Open`** first — GPU init must be ready (see test).
- **Skipping `MESH.UPLOAD`** — draw may fail or warn if the mesh is not uploaded.
- **`MESH.DRAWROTATED`** — **radians**, not degrees.
- **Confusing mesh and model** — **`MODEL.*`** loads whole assets; **`MESH.*`** is one mesh + your material + transform.

---

## See also

- [MODEL.md](MODEL.md) — full models and materials
- [IMAGE.md](IMAGE.md) — heightmap / cubicmap sources
- [CAMERA.md](CAMERA.md) — 3D camera
