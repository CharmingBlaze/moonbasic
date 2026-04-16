# Model Commands

Commands for loading, creating, transforming, rendering, and animating 3D models. Models are the primary way to get complex 3D geometry into your scene. They support glTF/GLB files, procedural mesh primitives, PBR materials, skeletal animation, LOD, instancing, and hierarchical attachment.

## Core Concepts

- **Model** — A renderable 3D object with one or more meshes, materials, and an optional skeleton. Created from file (`MODEL.LOAD`) or from a mesh handle (`MODEL.CREATE`).
- **Mesh** — The raw vertex/triangle data. Created procedurally or extracted from a model.
- **Material** — Surface properties (color, texture, metalness, roughness, shader). Each mesh slot in a model has its own material.
- **Instance** — A lightweight clone that shares the parent model's mesh data but has its own transform and material overrides. Very cheap to render.
- **LOD** — Level of Detail. Multiple model files at different polygon counts, swapped based on camera distance.
- All models are **heap handles** that must be freed.

---

## Loading

### `Model.Load(filePath)`

Loads a 3D model from a file. Supports `.glb`, `.gltf`, `.obj`, `.iqm`, `.vox`.

- `filePath` (string) — Path to model file.

**Returns:** `handle`

**How it works:** Reads the file, uploads mesh data to the GPU, creates default materials, and optionally loads embedded textures. The model is positioned at the origin with identity transform.

```basic
playerModel = Model.Load("assets/player.glb")
```

---

### `Model.LoadAsync(filePath)`

Begins loading a model in the background for streaming.

- `filePath` (string) — Path to model file.

**Returns:** `handle` — Pending model handle.

---

### `Model.Create(meshHandle)` / `Model.Make(meshHandle)`

Creates a model from a mesh handle. Use this when you generate meshes procedurally.

- `meshHandle` (handle) — Mesh created with `Mesh.*` commands.

**Returns:** `handle`

```basic
cubeMesh = Mesh.CreateCube(2, 2, 2)
cubeModel = Model.Create(cubeMesh)
```

---

### `Model.Free(modelHandle)` / `Instance.Free(instanceHandle)`

Frees a model or instance from GPU memory.

```basic
Model.Free(playerModel)
```

---

## Position & Transform

### `Model.SetPos(modelHandle, x, y, z)` / `Model.SetPosition(modelHandle, x, y, z)` / `model.pos(x, y, z)`

Sets the model's world position by writing a translation into its root transform matrix.

- `x`, `y`, `z` (float) — World coordinates.

```basic
Model.SetPos(playerModel, 0, 0, 0)
```

---

### `Model.GetPos(modelHandle)` / `Model.X(modelHandle)` / `Model.Y(modelHandle)` / `Model.Z(modelHandle)`

Returns the model's position. `GetPos` returns a Vec3 handle; `X`/`Y`/`Z` return individual float components.

---

### `Model.SetRot(modelHandle, pitch, yaw, roll)` / `model.rot(p, y, r)`

Sets the model's rotation in Euler degrees.

---

### `Model.Rotate(modelHandle, pitch, yaw, roll)`

Adds relative rotation to the model's current orientation.

---

### `Model.GetRot(modelHandle)`

Returns the model's rotation as a Vec3 handle.

---

### `Model.SetScale(modelHandle, sx, sy, sz)` / `model.scale(sx, sy, sz)`

Sets non-uniform scale.

---

### `Model.SetScaleUniform(modelHandle, scale)`

Sets uniform scale on all axes.

---

### `Model.GetScale(modelHandle)`

Returns the scale as a Vec3 handle.

---

### `Model.Move(modelHandle, forward, right, up)`

Moves the model relative to its orientation (like `Entity.Move`).

---

### `Model.SetMatrix(modelHandle, matrixHandle)`

Sets the model's root transform to a custom Mat4 matrix for full control.

---

## Rendering

### `Model.DrawAt(modelHandle, x, y, z, scale)` / `Model.DrawAt(modelHandle, x, y, z, sx, sy, sz)`

Draws a model at a specific position with a given scale, without modifying the model's stored transform.

- `x`, `y`, `z` (float) — Draw position.
- `scale` or `sx, sy, sz` (float) — Draw scale.

```basic
Camera.Begin(cam)
    Model.DrawAt(treeModel, 10, 0, 5, 1.0)
    Model.DrawAt(treeModel, -8, 0, 3, 0.8)   ; Reuse same model
Camera.End(cam)
```

---

### `Model.DrawEx(modelHandle, x, y, z, rotAxisX, rotAxisY, rotAxisZ, rotAngle, scale)`

Draws a model with position, rotation axis+angle, and scale in one call.

---

### `Model.DrawWires(modelHandle, x, y, z, scale)`

Draws a model in wireframe mode.

---

### `Model.Show(modelHandle)` / `Model.Hide(modelHandle)` / `Model.IsVisible(modelHandle)`

Toggle model visibility. Hidden models are skipped by `Entity.DrawAll()`.

---

## Materials & Appearance

### `Model.SetColor(modelHandle, r, g, b, a)` / `model.col(r, g, b)`

Sets the tint color applied to all materials on the model.

- `r`, `g`, `b`, `a` (int) — Color (0–255).

```basic
Model.SetColor(enemyModel, 255, 100, 100, 255)
```

---

### `Model.SetAlpha(modelHandle, alpha)` / `model.alpha(a)`

Sets the model's alpha transparency (0.0 = invisible, 1.0 = opaque).

---

### `Model.SetMetal(modelHandle, metalness)`

Sets the PBR metalness value (0.0 = dielectric, 1.0 = metal).

---

### `Model.SetRough(modelHandle, roughness)`

Sets the PBR roughness value (0.0 = mirror, 1.0 = matte).

---

### `Model.SetMaterial(modelHandle, materialHandle)`

Assigns a material handle to the model (replaces the default material).

---

### `Model.SetMaterialTexture(modelHandle, mapIndex, textureHandle)`

Sets a texture on a specific material map (diffuse, normal, metallic, etc.).

- `mapIndex` (int) — Material map index.
- `textureHandle` (handle) — Texture handle.

---

### `Model.SetMaterialShader(modelHandle, shaderHandle)`

Assigns a custom shader to the model's material.

---

### `Model.SetDiffuse(modelHandle, r, g, b)` / `Model.SetSpecular(modelHandle, r, g, b)` / `Model.SetSpecularPow(modelHandle, power)` / `Model.SetEmissive(modelHandle, r, g, b)` / `Model.SetAmbientColor(modelHandle, r, g, b)`

Set legacy lighting material properties.

---

### `Model.SetWireframe(modelHandle, enabled)` / `Model.SetCull(modelHandle, mode)` / `Model.SetLighting(modelHandle, enabled)` / `Model.SetFog(modelHandle, enabled)` / `Model.SetBlend(modelHandle, mode)` / `Model.SetDepth(modelHandle, enabled)`

Render state overrides per model.

---

## Skeletal Animation

### `Model.LoadAnimations(modelHandle, filePath)`

Loads animation data from a file (`.glb`, `.iqm`) and binds it to the model.

- `filePath` (string) — Path to animation file.

```basic
Model.LoadAnimations(playerModel, "assets/player_anims.glb")
```

---

### `Model.UpdateAnim(modelHandle, animIndex, frame)`

Updates the model's skeleton to a specific animation frame.

- `animIndex` (int) — Animation clip index (0-based).
- `frame` (int) — Frame number.

```basic
animFrame = animFrame + 30 * dt
Model.UpdateAnim(playerModel, 0, INT(animFrame))
```

---

### `Model.PlayIdx(modelHandle, animIndex)`

Starts playing an animation by index.

---

### `Model.Stop(modelHandle)`

Stops the current animation.

---

### `Model.Loop(modelHandle, enabled)`

Enables/disables animation looping.

---

### `Model.SetSpeed(modelHandle, fps)`

Sets the animation playback speed in frames per second.

---

### `Model.SetGPUSkinning(modelHandle, enabled)`

Enables GPU-based skeletal skinning for better performance with many animated models.

---

### `Model.GetFrame(modelHandle)`

Returns the current animation frame index.

**Returns:** `int`

---

### `Model.TotalFrames(modelHandle)`

Returns the total frame count of the current animation clip.

**Returns:** `int`

```basic
; Loop animation manually
frame = frame + 1
IF frame >= Model.TotalFrames(player) THEN frame = 0
Model.UpdateAnim(player, 0, frame)
```

---

### `Model.IsPlaying(modelHandle)`

Returns `TRUE` if the model is currently playing an animation.

**Returns:** `bool`

---

### `Model.AnimDone(modelHandle)`

Returns `TRUE` if the current animation has finished (non-looping animations only).

**Returns:** `bool`

**How it works:** Returns true when `animFrame >= frameCount - 1` and the animation is not playing. For looping animations, this is never true while playing.

```basic
; Wait for death animation to finish, then remove
IF Model.AnimDone(enemy) THEN
    Model.Free(enemy)
ENDIF
```

---

### `Model.AnimCount(modelHandle)`

Returns the number of animation clips loaded for this model.

**Returns:** `int`

---

### `Model.AnimName(modelHandle, animIndex)` / `Model.AnimName$(modelHandle, animIndex)`

Returns the name of an animation clip by index.

**Returns:** `string`

```basic
; List all animations
FOR i = 0 TO Model.AnimCount(player) - 1
    PRINT Model.AnimName(player, i)
NEXT
```

---

### `Model.LimbCount(modelHandle)`

Returns the number of bones in the model's skeleton.

**Returns:** `int`

---

## Hierarchy & Instancing

### `Model.Clone(modelHandle)`

Creates a deep copy of a model (own meshes, materials, transform).

**Returns:** `handle`

---

### `Model.Instance(modelHandle)`

Creates a lightweight instance that shares mesh data with the parent.

**Returns:** `handle`

```basic
; 100 trees from one model
tree = Model.Load("assets/tree.glb")
FOR i = 0 TO 99
    inst = Model.Instance(tree)
    Model.SetPos(inst, RND(-50, 50), 0, RND(-50, 50))
    Model.SetScaleUniform(inst, RND(0.8, 1.2))
NEXT
```

---

### `Model.AttachTo(childHandle, parentHandle)` / `Model.Detach(childHandle)`

Attaches a model as a child of another (inherits parent transform). Detach removes the relationship.

---

### `Model.GetParent(modelHandle)` / `Model.ChildCount(modelHandle)`

Query hierarchy relationships.

---

### `Model.Exists(modelHandle)`

Returns `TRUE` if the model handle is still valid.

**Returns:** `bool`

---

## LOD (Level of Detail)

### `Model.LoadLOD(filePath, lodLevel)`

Loads a model as a specific LOD level.

- `lodLevel` (int) — LOD index (0 = highest detail).

**Returns:** `handle`

---

### `Model.SetLODDistances(modelHandle, d1, d2, d3)`

Sets the distances at which LOD levels switch.

- `d1`, `d2`, `d3` (float) — Distance thresholds.

```basic
Model.SetLODDistances(treeModel, 20, 50, 100)
```

---

## Mesh Commands

### Procedural Mesh Generation

All mesh creation functions return a heap `handle`.

| Command | Arguments | Description |
|---------|-----------|-------------|
| `Mesh.CreateCube(w, h, d)` | width, height, depth | Axis-aligned box |
| `Mesh.MakeCube(w, h, d)` | (alias) | Same as CreateCube |
| `Mesh.CreateSphere(r, rings, slices)` | radius, rings, slices | UV sphere |
| `Mesh.MakeSphere(r, rings, slices)` | (alias) | Same |
| `Mesh.CreatePlane(w, d, resX, resZ)` | width, depth, X subdivisions, Z subdivisions | Flat ground plane |
| `Mesh.MakePlane(w, d, resX, resZ)` | (alias) | Same |
| `Mesh.MakeCylinder(radius, height, slices)` | radius, height, slices | Cylinder |
| `Mesh.MakeCone(radius, height, slices)` | radius, height, slices | Cone |
| `Mesh.MakeTorus(radius, size, radSegs, sides)` | radius, tube size, radial segments, sides | Torus/donut |
| `Mesh.MakeKnot(radius, size, radSegs, sides)` | radius, tube size, radial segments, sides | Trefoil knot |
| `Mesh.MakePoly(sides, radius)` | sides, radius | 2D polygon extruded |
| `Mesh.MakeHeightmap(imageHandle, size)` | image, scale | Terrain from heightmap image |
| `Mesh.MakeCubicmap(imageHandle, size)` | image, scale | Voxel-style map from image |

**How it works:** Each calls the corresponding Raylib `GenMesh*` function, then allocates the resulting mesh on the heap.

```basic
; Create a procedural torus
mesh = Mesh.MakeTorus(2, 0.5, 32, 16)
model = Model.Create(mesh)
```

---

### Mesh Operations

### `Mesh.Upload(meshHandle)`

Uploads mesh vertex data to the GPU. Called automatically by `Model.Create`, but useful if you modify vertices manually.

---

### `Mesh.Draw(meshHandle, materialHandle, x, y, z, sx, sy, sz)`

Draws a mesh directly with a material at a position and scale, without creating a model.

---

### `Mesh.DrawRotated(meshHandle, materialHandle, x, y, z, rx, ry, rz, sx, sy, sz)`

Draws a mesh with position, rotation (Euler), and scale.

---

### `Mesh.UpdateVertex(meshHandle, vertexIndex, x, y, z, nx, ny, nz, u, v)`

Modifies a single vertex's position, normal, and UV. Use for mesh deformation.

- `vertexIndex` (int) — 0-based vertex index.
- `x`, `y`, `z` (float) — New position.
- `nx`, `ny`, `nz` (float) — New normal.
- `u`, `v` (float) — New texture coordinates.

---

### `Mesh.GenTangents(meshHandle)`

Generates tangent vectors for the mesh (needed for normal mapping). Only available in purego builds.

---

### Mesh Queries

| Command | Returns |
|---------|---------|
| `Mesh.VertexCount(meshHandle)` | Number of vertices |
| `Mesh.TriangleCount(meshHandle)` | Number of triangles |
| `Mesh.GetBBoxMinX(meshHandle)` | Bounding box minimum X |
| `Mesh.GetBBoxMinY(meshHandle)` | Bounding box minimum Y |
| `Mesh.GetBBoxMinZ(meshHandle)` | Bounding box minimum Z |
| `Mesh.GetBBoxMaxX(meshHandle)` | Bounding box maximum X |
| `Mesh.GetBBoxMaxY(meshHandle)` | Bounding box maximum Y |
| `Mesh.GetBBoxMaxZ(meshHandle)` | Bounding box maximum Z |

```basic
; Check mesh size
mesh = Mesh.CreateCube(2, 2, 2)
PRINT "Vertices: " + STR(Mesh.VertexCount(mesh))
PRINT "Triangles: " + STR(Mesh.TriangleCount(mesh))
PRINT "Width: " + STR(Mesh.GetBBoxMaxX(mesh) - Mesh.GetBBoxMinX(mesh))
```

---

### `Mesh.Free(meshHandle)`

Frees a mesh from GPU memory.

---

## Shader Commands

### `Shader.Load(vertPath, fragPath)`

Loads a custom shader from vertex and fragment shader files.

- `vertPath` (string) — Vertex shader path (or `""` for default).
- `fragPath` (string) — Fragment shader path.

**Returns:** `handle`

---

### `Shader.GetLoc(shaderHandle, uniformName)`

Gets the location of a uniform variable in the shader.

**Returns:** `int`

---

### `Shader.SetFloat(shaderHandle, locOrName, value)` / `Shader.SetVec2(...)` / `Shader.SetVec3(...)` / `Shader.SetVec4(...)` / `Shader.SetInt(...)` / `Shader.SetTexture(...)`

Set shader uniform values.

---

### `Shader.Free(shaderHandle)`

Frees a shader.

---

## Material Commands

### `Material.Create()` / `Material.MakeDefault()`

Creates a new material with default PBR properties.

**Returns:** `handle`

---

### `Material.SetColor(materialHandle, r, g, b, a)`

Sets the material's base color.

---

### `Material.SetTexture(materialHandle, mapIndex, textureHandle)`

Assigns a texture to a material map slot.

---

### `Material.Free(materialHandle)`

Frees a material.

---

### `Model.SetModelMeshMaterial(modelHandle, meshIndex, materialIndex)`

Assigns a material slot to a specific mesh within a multi-mesh model.

---

## Shader Constants

Built-in constants for selecting preset shaders:

| Constant | Value | Description |
|----------|-------|-------------|
| `SHADER_PBR_LIT` | 1 | Standard PBR lit shader |
| `SHADER_PS1_RETRO` | 2 | PS1-style retro shader (vertex jitter, affine textures) |
| `SHADER_CEL_STYLED` | 3 | Cel/toon shading |
| `SHADER_WATER_PROCEDURAL` | 4 | Animated procedural water |
| `PP_BLOOM` | 101 | Bloom post-process shader |
| `PP_CRT_SCANLINES` | 102 | CRT scanline effect |
| `PP_PIXELATE` | 103 | Pixelation post-process |

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `LoadMesh(path)` | `Model.Load(path)` |
| `CreateCube(w, h, d)` | `Mesh.CreateCube(w, h, d)` |
| `CreateSphere(r, r, s)` | `Mesh.CreateSphere(r, r, s)` |

---

## Full Example

Loading a model, animating it, and rendering with PBR materials.

```basic
Window.Open(1280, 720, "Model Demo")
Window.SetFPS(60)

cam = Camera.Create()
cam.pos(0, 3, 8)
cam.look(0, 1, 0)
cam.fov(60)

Render.SetAmbient(80, 80, 100)

; Load animated character
player = Model.Load("assets/character.glb")
Model.LoadAnimations(player, "assets/character.glb")
Model.SetPos(player, 0, 0, 0)
Model.SetMetal(player, 0.0)
Model.SetRough(player, 0.7)

; Create procedural floor
floorMesh = Mesh.CreatePlane(20, 20, 1, 1)
floorModel = Model.Create(floorMesh)
Model.SetColor(floorModel, 100, 130, 100, 255)

animFrame = 0

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()

    ; Animate
    animFrame = animFrame + 30 * dt
    Model.UpdateAnim(player, 0, INT(animFrame))

    ; Rotate model with arrows
    IF Input.KeyDown(KEY_LEFT) THEN Model.Rotate(player, 0, 90 * dt, 0)
    IF Input.KeyDown(KEY_RIGHT) THEN Model.Rotate(player, 0, -90 * dt, 0)

    Render.Clear(40, 50, 70)
    Camera.Begin(cam)
        Draw.Grid(20, 1.0)
        Model.DrawAt(floorModel, 0, 0, 0, 1.0)
        Model.DrawAt(player, 0, 0, 0, 1.0)
    Camera.End(cam)

    Draw.Text("Arrow keys = Rotate", 10, 10, 18, 255, 255, 255, 255)
    Render.Frame()
WEND

Model.Free(player)
Model.Free(floorModel)
Mesh.Free(floorMesh)
Camera.Free(cam)
Window.Close()
```

---

## See Also

- [TEXTURE](TEXTURE.md) — Textures for model materials
- [ENTITY](ENTITY.md) — Entities wrap models with position/collision
- [CAMERA](CAMERA.md) — Camera for rendering models
- [RENDER](RENDER.md) — Render pipeline and PBR ambient
