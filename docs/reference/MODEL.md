# Model — `MODEL.*`

A **model** is a Raylib **`Model`**: meshes, materials, optional animation data, and a root transform. Registry keys use **dots and uppercase** (e.g. `MODEL.LOAD`). This page matches the **current** runtime, not legacy PascalCase-only docs.

**Requires CGO.** See **[MESH.md](MESH.md)** for raw **`MESH.*`** geometry.

---

### Model.Load

```basic
mdl = MODEL.LOAD(path$)
```

Loads **GLTF/GLB/OBJ/IQM/B3D** (Raylib loaders). Returns a model handle.

**Returns** — handle.

---

### Model.Make

```basic
mdl = MODEL.MAKE(mesh)
```

Builds a model from an existing **`MESH`** handle via **`LoadModelFromMesh`**. The engine marks the mesh as **consumed by the model**: GPU data is owned by the new model; **`MESH.FREE(mesh)`** does **not** call `UnloadMesh` again (avoids double-free). Always **`MODEL.FREE(mdl)`** to release the model (which unloads the mesh data).

> **Common mistake:** Calling **`MESH.FREE`** before **`MODEL.FREE`** expecting to unload VRAM first — with **`MODEL.MAKE`**, free the **model** first, then release the **mesh** handle for the heap slot only.

---

### Model.Draw

```basic
MODEL.DRAW(mdl)
```

**One argument** — the model, **LOD model**, or **instanced model** handle. There is **no** separate matrix argument: use **`MODEL.SETPOS`** / internal transforms on the Raylib model (see below).

**Phase:** Call between **`CAMERA.BEGIN`** and **`CAMERA.END`** for 3D.

---

### Model.SetPos

```basic
MODEL.SETPOS(mdl, x#, y#, z#)
MODEL.SETPOSITION(mdl, x#, y#, z#)
```

Sets the model’s root transform to **translation** (replaces rotation/scale from the previous matrix). Works for **`modelObj`** and **`lodModelObj`**.

---

### Model.Free

```basic
MODEL.FREE(mdl)
```

Unloads the model and meshes/materials (`UnloadModel`). Pair with **`MODEL.LOAD`** / **`MODEL.MAKE`**.

---

### Model.GetMaterialCount

```basic
n = MODEL.GETMATERIALCOUNT(mdl)
```

**Returns** — integer.

---

## Materials and appearance (selected)

| Command | Notes |
|---|---|
| `MODEL.SETDIFFUSE` | `(mdl, r, g, b)` — albedo tint, 0–255 |
| `MODEL.SETSPECULAR` / `SETSPECULARPOW` | Specular colour / shininess |
| `MODEL.SETEMISSIVE` | Emissive colour |
| `MODEL.SETAMBIENTCOLOR` | Ambient tint |
| `MODEL.SETALPHA` | Alpha channel on albedo maps |
| `MODEL.SETMATERIAL` | Replace material by index |
| `MODEL.SETMATERIALTEXTURE` / `SETMATERIALSHADER` | Per-slot texture/shader |
| `MODEL.SETMODELMESHMATERIAL` | Assign material index per mesh |

Texture **stage** helpers (`SETTEXTURESTAGE`, `SETSTAGEBLEND`, `SCROLLTEXTURE`, …) are for multi-layer scrolling effects — see runtime registration in **`model_texture_stages_cgo.go`**.

---

## Render state toggles

| Command | Purpose |
|---|---|
| `MODEL.SETWIREFRAME` | Wireframe overlay hint |
| `MODEL.SETCULL` | Face culling |
| `MODEL.SETLIGHTING` / `SETFOG` | Lighting / fog hints |
| `MODEL.SETBLEND` | Blend mode |
| `MODEL.SETDEPTH` | Depth test/write hints |
| `MODEL.SETGPUSKINNING` | GPU skinning when available |

---

## Scene graph (handles)

| Command | Notes |
|---|---|
| `MODEL.CLONE` | Duplicate model data |
| `MODEL.INSTANCE` | Shared-mesh instance |
| `MODEL.ATTACHTO` / `DETACH` | Parent/child |
| `MODEL.EXISTS` | Valid handle check |

---

## LOD

| Command | Notes |
|---|---|
| `MODEL.LOADLOD` | Multiple paths + LOD distances |
| `MODEL.SETLODDISTANCES` | Per-LOD ranges |

---

## Instancing

| Command | Notes |
|---|---|
| `MODEL.MAKEINSTANCED` | `path$`, instance count |
| `MODEL.SETINSTANCEPOS` / `SETINSTANCESCALE` | Per-instance |
| `MODEL.UPDATEINSTANCES` | Rebuild instance matrices |
| `MODEL.DRAW` | Same as regular draw (handles instanced object) |

---

## Example (procedural mesh → model)

See **`testdata/model_complete_test.mb`**: **`MESH.MAKECUBE`** → **`MODEL.MAKE`** → **`MODEL.SETPOS`** / **`MODEL.SETDIFFUSE`** → **`CAMERA.BEGIN`** / **`MODEL.DRAW`** / **`CAMERA.END`**.

---

## Common mistakes

- **`MODEL.DRAW(mdl, matrix)`** — not supported; use **`MODEL.SETPOS`** (and future rotation APIs if added).
- **`mod` as a variable name** — **`MOD`** is reserved in moonBASIC; use **`mdl`** or **`modelHandle`**.
- **Double-free after `MODEL.MAKE`** — follow **`MODEL.FREE`** then **`MESH.FREE`** (mesh slot only) as in the test, or read **`consumedByModel`** behaviour above.

---

## See also

- [MESH.md](MESH.md) — procedural meshes, **`MESH.UPLOAD`**, **`MESH.DRAW`**
- [CAMERA.md](CAMERA.md) — 3D camera
- [LIGHT.md](LIGHT.md) — PBR lighting
- [SHADER.md](SHADER.md) — custom materials via shaders
