# Mesh — `MESH.*`

CPU mesh geometry for **`MESH.Draw`** / **`MESH.DRAWROTATED`** with a **material** and **matrix** handle. Meshes are **Raylib `Mesh`** values stored in the VM heap (`TagMesh`).

**Requires CGO** (same as **`MODEL.*`**, **`Camera.*`**).

---

## Procedural generators

| Command | Arguments | Notes |
|--------|-----------|--------|
| `MESH.MAKEPOLY` | `sides, radius#` | Regular polygon prism |
| `MESH.MAKEPLANE` | `width#, length#, resX, resZ` | Subdivided plane |
| `MESH.MAKECUBE` | `width#, height#, length#` | Axis-aligned box |
| `MESH.MAKESPHERE` | `radius#, rings, slices` | UV sphere |
| `MESH.MAKECYLINDER` | `radius#, height#, slices` | |
| `MESH.MAKECONE` | `radius#, height#, slices` | |
| `MESH.MAKETORUS` | `radius#, size#, radSeg, sides` | |
| `MESH.MAKEKNOT` | `radius#, size#, radSeg, sides` | Trefoil-style knot |
| `MESH.MAKEHEIGHTMAP` | `image, sizeX#, sizeY#, sizeZ#` | **Image** heightfield → mesh |
| `MESH.MAKECUBICMAP` | `image, cubeX#, cubeY#, cubeZ#` | Cubic map from **image** |

### Legacy aliases (same Raylib generators)

| Command | Same as |
|--------|---------|
| `MESH.CUBE` | `MESH.MAKECUBE` |
| `MESH.SPHERE` | `MESH.MAKESPHERE` |
| `MESH.PLANE` | `MESH.MAKEPLANE` |

---

## GPU upload, draw, edit

| Command | Arguments | Notes |
|--------|-----------|--------|
| `MESH.UPLOAD` | `(mesh)` | Uploads mesh to GPU (Raylib `UploadMesh`) |
| `MESH.DRAW` | `(mesh, material, matrix)` | Standard draw |
| `MESH.DRAWROTATED` | `(mesh, material, matrix, rotationAxisX, rotationAxisY, rotationAxisZ, angleDeg#)` | Axis-angle variant |
| `MESH.UPDATEVERTEX` | `(mesh, vertexIndex, x#, y#, z#)` | CPU vertex write |
| `MESH.GENTANGENTS` | `(mesh)` | Regenerate tangents/binormals |

---

## Bounding box (local mesh bounds)

| Command | Returns |
|--------|---------|
| `MESH.GETBBOXMINX` … `MESH.GETBBOXMAXZ` | `(mesh)` → float components |

---

## Lifecycle

| Command | Notes |
|--------|--------|
| `MESH.FREE` | Unloads mesh data (`UnloadMesh`); call when done |

---

## Common mistakes

- **Drawing without upload** — Call **`MESH.UPLOAD`** after procedural generation before **`MESH.DRAW`** in the usual Raylib workflow.
- **Confusing mesh and model** — A **`MODEL.*`** bundles meshes + materials; **`MESH.*`** is a single mesh + your **`MATERIAL.*`** + **`MAT4` / transform** matrix.

---

## See also

- [MODEL.md](MODEL.md) — full models, hierarchy, materials
- [IMAGE.md](IMAGE.md) — heightmap / cubicmap source images
