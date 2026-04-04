# Mat4 (legacy)

**Prefer [`Transform.*`](TRANSFORM.md)** for new code—the names describe what you are doing (moving and rotating objects in 3D).

`Mat4.*` remains available for compatibility. Mapping:

| Legacy | Use instead |
|--------|-------------|
| `Mat4.Identity()` | `Transform.Identity()` |
| `Mat4.FromTranslation(x,y,z)` | `Transform.Translation(x,y,z)` |
| `Mat4.FromRotation` / `Mat4.Rotation` | `Transform.Rotation` |
| `Mat4.FromScale` | `Transform.Scale` |
| `Mat4.SetRotation` | `Transform.SetRotation` |
| `Mat4.Multiply`, `Inverse`, `Transpose` | `Transform.Multiply`, … |
| `Mat4.LookAt`, `Perspective`, `Ortho` | `Transform.LookAt`, … |
| `Mat4.GetElement` | `Transform.GetElement` |
| `Mat4.TransformX/Y/Z(m, x, y, z)` | `Transform.ApplyX/Y/Z(m, x, y, z)` |
| `Mat4.Free` | `Transform.Free` |

Full tutorials and behavior: **[TRANSFORM.md](TRANSFORM.md)**.
