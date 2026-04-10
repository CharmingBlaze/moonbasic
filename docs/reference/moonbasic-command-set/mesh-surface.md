# Mesh / surface

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **LoadMesh (file, parent)** | **`ENTITY.LOADMESH`**, **`MESH.LOAD`** | **Mesh/model** — **`MESH.FREE`**, **`ENTITY.FREE`** as appropriate. |
| **LoadAnimMesh** | **`ENTITY.LOADANIMATEDMESH`** | **Animations first** unload order — [MEMORY.md](../../MEMORY.md). |
| **MeshWidth / Height / Depth** | **`MESH.GETBBOX*`** (extents) | |
| **AddSurface / AddVertex / AddTriangle** | **`MESH.MAKECUSTOM`**, **`UPDATEVERTEX`**, procedural **`MESH.MAKE*`** | No DBPro-style surface index API — different model — [MESH.md](../MESH.md). |
| **VertexX/Y/Z** | **`MESH.UPDATEVERTEX`**, queries via bbox / mesh ops | |
