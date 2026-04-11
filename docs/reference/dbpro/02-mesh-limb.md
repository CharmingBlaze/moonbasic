# DBPro — Mesh / limb

DBPro **limbs** map only partially to moonBASIC: you get **model** / **mesh** handles and some **hierarchy** queries, not a full DBPro limb editor.

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **MAKE MESH FROM OBJECT (mesh, obj)** | ✓ **`Model.Make()`** | Create model from mesh handle. |
| **ADD MESH TO VERTEXDATA** | ≈ **`Mesh.Upload()`** | |
| **DELETE MESH (mesh)** | ✓ **`Mesh.Free()`** | |
| **SET MESH (obj, mesh_id)** | ≈ **`Model.Attach()`** / **`Model.Reload()`** | Depends on asset pipeline. |
| **ADD LIMB (obj, limb, mesh)** | ≈ **`Entity.Parent()`** / **`Entity.FindBone()`** | Limb logic maps to hierarchy/bones. |
| **REMOVE LIMB (obj, limb)** | ✓ **`Entity.Unparent()`** | |
| **LIMB EXIST (obj, limb)** | ✓ **`Entity.FindChild()`** | |
| **LIMB COUNT (obj)** | ✓ **`Entity.CountChildren()`** | |
| **SET LIMB TEXTURE** / **ALPHA** / **COLOR** / **LIGHT** / **WIREFRAME** / **SHADING** | ≈ **`Model.SetColor()`**, **`Material.*()`**, **`Texture.*()`** | Per-submesh control varies. |

---

## Limb transform

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **POSITION LIMB** / **ROTATE LIMB** / **SCALE LIMB** | ≈ **`Model.SetMatrix()`**, child transforms | Often done via **animation** or **scene** tools, not one call per limb index. |
| **LIMB POSITION X/Y/Z** | ≈ **`Model.X()`** / **`Model.Y()`** / **`Model.Z()`** (whole model) or custom | No universal **limb_id** accessor in DBPro form. |
