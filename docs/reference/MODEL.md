# 3D Model & Mesh Commands

Commands for creating, loading, and drawing 3D models and meshes.

---

## Models

Models are complex objects that contain one or more meshes, along with their own materials and transformation data.

### `Model.Load(filePath$)`

Loads a 3D model from a file (e.g., `.gltf`, `.glb`, `.obj`). This command loads all associated meshes and materials. It returns a handle to the model.

- `filePath$`: The path to the model file.

```basic
; Load a spaceship model
ship_model = Model.Load("assets/ship.glb")
```

---

### `Model.Draw(modelHandle, matrixHandle)`

Draws a complete model, including all of its meshes, using a given transformation matrix.

- `modelHandle`: The handle of the model to draw.
- `matrixHandle`: A handle to a `Mat4` transformation matrix that defines the model's position, rotation, and scale.

```basic
; Create a transformation matrix for the ship
ship_transform = Mat4.Identity()
Mat4.SetTranslation(ship_transform, 0, 5, 0) ; Move it up

; In the main loop
Render.Clear(0,0,0)
cam.Begin()
    Model.Draw(ship_model, ship_transform)
cam.End()
Render.Frame()
```

---

### `Model.SetMaterial(modelHandle, materialIndex, materialHandle)`

Replaces one of a model's existing materials with a new one. This is useful for customizing the appearance of a loaded model.

- `modelHandle`: The handle of the model.
- `materialIndex`: The zero-based index of the material to replace.
- `materialHandle`: The handle of the new material.

---

### `Model.Free(modelHandle)`

Unloads a model and all its associated meshes and materials from memory. This is important for preventing memory leaks.

---

## Meshes

Meshes are the building blocks of 3D objects, containing vertex data like positions and texture coordinates. You can create them procedurally or load them as part of a model.

### `Mesh.MakeCube(width#, height#, depth#)`

Creates a procedural cube mesh and returns a handle to it.

```basic
; Create a 2x2x2 cube mesh
cube_mesh = Mesh.MakeCube(2, 2, 2)
```

---

### `Mesh.Draw(meshHandle, materialHandle, matrixHandle)`

Draws a single mesh with a specific material and transformation. This is more low-level than `Model.Draw`.

- `meshHandle`: The handle of the mesh.
- `materialHandle`: The handle of the material to apply.
- `matrixHandle`: The transformation matrix for the mesh.

```basic
; Create a cube, a material, and a transform
cube_mesh = Mesh.MakeCube(2, 2, 2)
my_material = Material.MakeDefault()
cube_transform = Mat4.Identity()

; In the main loop, draw the single mesh
Mesh.Draw(cube_mesh, my_material, cube_transform)
```

---

### Other Procedural Meshes

- `Mesh.MakeSphere(radius#, rings, slices)`: Creates a sphere.
- `Mesh.MakePlane(width#, length#, resX, resZ)`: Creates a flat plane.

---

### `Mesh.Free(meshHandle)`

Unloads a mesh from memory.

---

## Materials

Materials define the surface appearance of a mesh (color, texture, shininess).

### `Material.MakeDefault()`

Creates a new material with default PBR (Physically-Based Rendering) settings. Returns a handle.

---

### `Material.MakePBR()`

Creates a full PBR material with support for shadow sampling and additional
uniform slots. Use this when working with custom PBR shaders.

---

### `Material.SetTexture(materialHandle, mapType, textureHandle)`

Assigns a texture to a specific map channel of a material.

- `materialHandle`: The handle of the material.
- `mapType`: The material map type. Common values:
  - `MATERIAL_MAP_ALBEDO` (also `MAP_DIFFUSE`) — the base color texture.
  - `MATERIAL_MAP_METALNESS` — metalness map.
  - `MATERIAL_MAP_ROUGHNESS` — roughness map.
  - `MATERIAL_MAP_NORMAL` — normal map.
- `textureHandle`: The handle of the texture, loaded with `Texture.Load()`.

```basic
; Create a material and apply a texture to it
stone_mat = Material.MakeDefault()
stone_tex = Texture.Load("assets/stone.png")
Material.SetTexture(stone_mat, MATERIAL_MAP_ALBEDO, stone_tex)
```

---

### `Material.SetColor(materialHandle, mapType, r, g, b, a)`

Sets the tint color of a material map slot. This is a fast way to color a mesh
without needing a texture — the default material renders with this color.

- `mapType`: Typically `MATERIAL_MAP_ALBEDO`.
- `r`, `g`, `b`, `a`: Color components (0–255).

```basic
; Create a solid red cube without any texture file
red_mat = Material.MakeDefault()
Material.SetColor(red_mat, MATERIAL_MAP_ALBEDO, 220, 50, 50, 255)
```

---

### `Material.SetFloat(materialHandle, mapType, value#)`

Sets the float value of a material map slot. Used for metalness and roughness
values when not using a texture map.

```basic
; Make a shiny metal surface
chrome_mat = Material.MakeDefault()
Material.SetFloat(chrome_mat, MATERIAL_MAP_METALNESS, 1.0)
Material.SetFloat(chrome_mat, MATERIAL_MAP_ROUGHNESS, 0.1)
```

---

### `Material.SetShader(materialHandle, shaderHandle)`

Assigns a custom GLSL shader to a material. See [Shader Reference](SHADER.md).

---

### `Material.Free(materialHandle)`

Frees a material resource. This does not free textures assigned to the material;
call `Texture.Free()` for those separately.
