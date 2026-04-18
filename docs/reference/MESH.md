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

## Extended Command Reference

### Creation aliases

| Command | Description |
|--------|-------------|
| `MESH.CREATECUBE(size)` | Cube mesh. Alias: `MESH.CREATECUBE`. |
| `MESH.CREATESPHERE(r, rings, slices)` | Sphere mesh. |
| `MESH.CREATECYLINDER(r, h, slices)` | Cylinder mesh. |
| `MESH.CREATECONE(r, h, slices)` | Cone mesh. |
| `MESH.CREATECAPSULE(r, h, rings, slices)` | Capsule mesh. |
| `MESH.CREATEPLANE(w, d, resX, resZ)` | Plane with resolution. |
| `MESH.CREATETORUS(r, size, radSeg, sides)` | Torus mesh. |
| `MESH.CREATEKNOT(r, size, radSeg, sides)` | Knot/trefoil mesh. |
| `MESH.CREATEPOLY(sides, r)` | Regular polygon mesh. |
| `MESH.CREATECUBICMAP(img, scale)` | Voxel terrain mesh from cubicmap image. |
| `MESH.CREATEHEIGHTMAP(img, sx, sy, sz)` | Terrain mesh from heightmap image. |
| `MESH.CREATECUSTOM(vertCount, triCount)` | Empty custom mesh for manual vertex fill. |
| `MESH.MAKECAPSULE(r,h,rings,slices)` | Alias of `MESH.CREATECAPSULE`. |
| `MESH.MAKECONE(r,h,slices)` | Alias of `MESH.CREATECONE`. |
| `MESH.MAKECYLINDER(r,h,slices)` | Alias of `MESH.CREATECYLINDER`. |
| `MESH.MAKECUBICMAP(img, scale)` | Alias of `MESH.CREATECUBICMAP`. |
| `MESH.MAKEHEIGHTMAP(img, sx,sy,sz)` | Alias of `MESH.CREATEHEIGHTMAP`. |
| `MESH.MAKECUSTOM(verts, tris)` | Alias of `MESH.CREATECUSTOM`. |
| `MESH.MAKETORUS(r,size,radSeg,sides)` | Alias of `MESH.CREATETORUS`. |
| `MESH.MAKEKNOT(r,size,radSeg,sides)` | Alias of `MESH.CREATEKNOT`. |
| `MESH.MAKEPOLY(sides, r)` | Alias of `MESH.CREATEPOLY`. |

### Bounds & LOD

| Command | Description |
|--------|-------------|
| `MESH.GENERATEBOUNDS(mesh)` | Compute and store AABB/sphere bounds. |
| `MESH.GETBBOXMINX(mesh)` / `MESH.GETBBOXMINY(mesh)` / `MESH.GETBBOXMINZ(mesh)` | Bounding box minimum per axis. |
| `MESH.GETBBOXMAXX(mesh)` / `MESH.GETBBOXMAXY(mesh)` / `MESH.GETBBOXMAXZ(mesh)` | Bounding box maximum per axis. |
| `MESH.GENERATELOD(mesh, ratio)` | Generate one LOD level at `ratio` poly reduction. |
| `MESH.GENERATELODCHAIN(mesh, levels)` | Generate a multi-level LOD chain. |

### Normals & tangents

| Command | Description |
|--------|-------------|
| `MESH.GENERATENORMALS(mesh)` | Recompute smooth vertex normals. |
| `MESH.GENTANGENTS(mesh)` | Compute tangent vectors for normal mapping. |

### Vertex updates

| Command | Description |
|--------|-------------|
| `MESH.UPDATEVERTEX(mesh, index, x, y, z)` | Set position of a single vertex. |
| `MESH.UPDATEVERTICES(mesh, floatArray)` | Bulk-update all positions from a flat float array. |

### Optimisation

| Command | Description |
|--------|-------------|
| `MESH.OPTIMISEVERTEXCACHE(mesh)` / `MESH.OPTIMIZEVERTEXCACHE` | Reorder triangles for post-transform cache. |
| `MESH.OPTIMISEOVERDRAW(mesh)` / `MESH.OPTIMIZEOVERDRAW` | Reorder to reduce overdraw. |
| `MESH.OPTIMISEFETCH(mesh)` / `MESH.OPTIMIZEFETCH` | Reorder vertices for pre-transform cache. |
| `MESH.OPTIMISEALL(mesh)` / `MESH.OPTIMIZEALL` | Apply all optimisations. |

### Export & draw variants

| Command | Description |
|--------|-------------|
| `MESH.EXPORT(mesh, path)` | Export mesh to file (OBJ/GLB). |
| `MESH.DRAWAT(mesh, matHandle, x, y, z, sx, sy, sz)` | Draw mesh with explicit transform and material. |
| `MESH.DRAWINSTANCED(mesh, matHandle, transforms)` | Draw mesh with a transform array for instancing. |

---

## See also

- [MODEL.md](MODEL.md) — full models and materials
- [IMAGE.md](IMAGE.md) — heightmap / cubicmap sources
- [CAMERA.md](CAMERA.md) — 3D camera
