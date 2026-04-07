# DBPro — Mesh / limb

DBPro **limbs** map only partially to moonBASIC: you get **model** / **mesh** handles and some **hierarchy** queries, not a full DBPro limb editor.

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **MAKE MESH FROM OBJECT** | ≈ **`MESH.*`**, **`MODEL` export/extract patterns** | See [MESH.md](../MESH.md). |
| **DELETE MESH** | ✓ **`MESH.FREE`** | |
| **SET MESH (obj, mesh_id)** | ≈ **`MODEL` attach / reload** | Depends on asset pipeline. |
| **ADD LIMB** / **DELETE LIMB** | — | No DBPro-style limb list; use **`MODEL.CHILDCOUNT`**, **`MODEL.LIMBCOUNT`** where exposed. |
| **SET LIMB TEXTURE** / **ALPHA** / **COLOR** / **LIGHT** / **WIREFRAME** / **SHADING** | ≈ **`MODEL.SETCOLOR`**, **`MATERIAL.*`**, **`TEXTURE.*`** | Per-submesh control varies. |

## Limb transform

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **POSITION LIMB** / **ROTATE LIMB** / **SCALE LIMB** | ≈ **`MODEL.SETMATRIX`**, child transforms | Often done via **animation** or **scene** tools, not one call per limb index. |
| **LIMB POSITION X/Y/Z** | ≈ **`MODEL.X/Y/Z`** (whole model) or custom | No universal **limb_id** accessor in DBPro form. |
