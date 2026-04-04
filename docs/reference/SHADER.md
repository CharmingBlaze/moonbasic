# Shader Commands

Commands for loading and using custom shaders.

---

### `Shader.Load(vertexPath$, fragmentPath$)`

Loads a GLSL vertex and fragment shader from files. Returns a handle to the shader program.

- `vertexPath$`: Path to the vertex shader file (`.vs`).
- `fragmentPath$`: Path to the fragment shader file (`.fs`).

---

**Note:** The runtime does **not** currently expose global `Render.BeginShader` / `Render.EndShader`. Apply shaders via **`Model.SetMaterialShader`** / material APIs and **`Shader.Set*`** uniforms (below), or other module-specific paths.

---

### `Model.SetMaterialShader(model, materialIndex, shaderHandle)`

Assigns a shader to a **material slot** on a model (`MODEL.SETMATERIALSHADER`).

---

## Uniforms and lifecycle

Uniform names are **string** arguments (bytecode string table).

| Command | Purpose |
|--------|---------|
| `Shader.Free(shader)` | Unloads shader GPU program (heap handle). |
| `Shader.GetLoc(shader, name$)` | Raylib location index (`int`). |
| `Shader.SetFloat(shader, name$, value#)` | |
| `Shader.SetVec2` / `SetVec3` / `SetVec4` | Name + float components. |
| `Shader.SetInt(shader, name$, value)` | |
| `Shader.SetTexture(shader, name$, textureHandle)` | Binds a `TEXTURE` heap handle. |

```basic
sh = Shader.Load("custom.vs", "custom.fs")
t# = TIME.GET()
Shader.SetFloat(sh, "uTime", t#)
Shader.Free(sh)
```
