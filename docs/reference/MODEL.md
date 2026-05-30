# Model Commands

Load, position, and draw 3D models (glTF, GLB, OBJ, IQM, B3D) with materials and transforms.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Load with `MODEL.LOAD` or create from a mesh with `MODEL.CREATE`.
2. Position with `MODEL.SETPOS`, rotate with `MODEL.SETROT`, scale with `MODEL.SETSCALE`.
3. Draw between `RENDER.BEGIN3D` / `RENDER.END3D` with `MODEL.DRAW`.
4. Free with `MODEL.FREE`.

For raw meshes see [MESH.md](MESH.md). For skeletal animation see [ANIMATION_3D.md](ANIMATION_3D.md).

---

### `MODEL.LOAD(path)`
Loads a 3D model file (glTF, GLB, OBJ, IQM, B3D).

- **Arguments**:
    - `path`: (String) File path.
- **Returns**: (Handle) The new model handle.
- **Example**:
    ```basic
    hero = MODEL.LOAD("hero.glb")
    ```

---

### `MODEL.CREATE(mesh)`
Builds a model from an existing **`Mesh`** handle.

- **Returns**: (Handle) The new model handle.

---

### `MODEL.DRAW(handle)`
Draws the model using its current transform.

- **Returns**: (Handle) The model handle (for chaining).

---

### `MODEL.SETPOS(handle, x, y, z)` / `SETROT` / `SETSCALE`
Sets the model's world position, Euler rotation (radians), or scale.

- **Returns**: (Handle) The model handle (for chaining).

---

### `MODEL.SETMATERIAL(handle, index, material)`
Replaces a material slot in the model.

- **Returns**: (Handle) The model handle (for chaining).

---

### `MODEL.FREE(handle)`
Unloads the model and its resources.

---

## Full Example (load and draw)

```basic
WINDOW.OPEN(1280, 720, "Model Example")
WINDOW.SETFPS(60)
cam = CAMERA.CREATE()

mdl = MODEL.LOAD("assets/character.glb")
MODEL.SETPOS(mdl, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(20, 20, 20)
    RENDER.BEGIN3D(cam)
        MODEL.DRAW(mdl)
    RENDER.END3D()
    RENDER.FRAME()
WEND

MODEL.FREE(mdl)
WINDOW.CLOSE()
```

---

## Common mistakes

- **`MODEL.DRAW(mdl, matrix)`** — not supported; use **`MODEL.SETPOS`** (canonical) or deprecated **`MODEL.SETPOSITION`**, **`SETMATRIX`**, **`DRAWAT`**.
- **`mod` as a variable name** — **`MOD`** is reserved in moonBASIC; use **`mdl`** or **`modelHandle`**.
- **Double-free after `MODEL.CREATE` (mesh → model)** — **`MODEL.MAKE`** is deprecated with the same arity; follow **`MODEL.FREE`** then **`MESH.FREE`** (mesh slot only) as in the test, or read **`consumedByModel`** behaviour above.

---

## Extended Command Reference

### Creation & loading

| Command | Description |
|--------|-------------|
| `MODEL.CREATEBOX(w,h,d)` / `MODEL.MAKEBOX(w,h,d)` | Box mesh model. |
| `MODEL.CREATECAPSULE(r,h)` / `MODEL.MAKECAPSULE(r,h)` | Capsule mesh model. |
| `MODEL.CREATEINSTANCED(mesh, count)` / `MODEL.MAKEINSTANCED(mesh, count)` | GPU-instanced model with `count` slots. |
| `MODEL.LOADASYNC(path, callback)` | Non-blocking load; fires `callback(handle)` on completion. |
| `MODEL.LOADLOD(pathLow, pathMid, pathHigh)` | Load LOD chain. |
| `MODEL.LOADANIMATIONS(mdl, path)` | Load extra animation clips from a separate file onto `mdl`. |
| `MODEL.CLONE(mdl)` | Shallow-clone a model handle (shares mesh/material data). |
| `MODEL.EXISTS(mdl)` | Returns `TRUE` if handle is valid. |
| `MODEL.ISLOADED(mdl)` | Returns `TRUE` if async load has completed. |

---

### Visibility & hierarchy

| Command | Description |
|--------|-------------|
| `MODEL.SHOW(mdl)` | Make model visible. |
| `MODEL.HIDE(mdl)` | Make model invisible. |
| `MODEL.ISVISIBLE(mdl)` | Returns `TRUE` if visible. |
| `MODEL.ATTACHTO(mdl, parentMdl)` | Parent `mdl` to `parentMdl`. |
| `MODEL.DETACH(mdl)` | Clear parent; return to world space. |
| `MODEL.ADDCHILD(mdl, childMdl)` | Add child to model hierarchy. |
| `MODEL.REMOVECHILD(mdl, childMdl)` | Remove a child from hierarchy. |
| `MODEL.GETPARENT(mdl)` | Returns parent model handle (or 0). |
| `MODEL.GETCHILD(mdl, index)` | Returns child handle at `index`. |
| `MODEL.CHILDCOUNT(mdl)` | Returns number of direct children. |

---

### Transform queries

| Command | Description |
|--------|-------------|
| `MODEL.GETPOS(mdl)` | Returns `[x,y,z]` world position. |
| `MODEL.GETROT(mdl)` | Returns `[pitch,yaw,roll]` in radians. |
| `MODEL.GETSCALE(mdl)` | Returns `[sx,sy,sz]`. |
| `MODEL.SETSCALEUNIFORM(mdl, s)` | Set uniform XYZ scale. |
| `MODEL.MOVE(mdl, dx, dy, dz)` | Translate by delta (world space). |

---

### Color & alpha

| Command | Description |
|--------|-------------|
| `MODEL.COLOR(mdl, r, g, b, a)` | Set RGBA tint. |
| `MODEL.SETCOLOR(mdl, r, g, b, a)` | Alias of `MODEL.COLOR`. |
| `MODEL.GETCOLOR(mdl)` | Returns `[r,g,b,a]` tint. |
| `MODEL.ALPHA(mdl, a)` | Set alpha 0.0–1.0. |
| `MODEL.SETALPHA(mdl, a)` | Alias of `MODEL.ALPHA`. |
| `MODEL.GETALPHA(mdl)` | Returns current alpha. |

---

### Material & shader

| Command | Description |
|--------|-------------|
| `MODEL.SETDIFFUSE(mdl, r, g, b, a)` | Set PBR diffuse/albedo color. |
| `MODEL.SETEMISSIVE(mdl, r, g, b, a)` | Set emissive color. |
| `MODEL.SETMETAL(mdl, v)` | Set metalness 0.0–1.0. |
| `MODEL.SETROUGH(mdl, v)` | Set roughness 0.0–1.0. |
| `MODEL.SETSPECULAR(mdl, r, g, b)` | Set specular color. |
| `MODEL.SETSPECULARPOW(mdl, p)` | Set specular exponent. |
| `MODEL.SETAMBIENTCOLOR(mdl, r, g, b)` | Set ambient color override. |
| `MODEL.SETMATERIALSHADER(mdl, slot, shaderHandle)` | Assign shader to material slot. |
| `MODEL.SETMATERIALTEXTURE(mdl, slot, texHandle)` | Set texture on a material slot. |
| `MODEL.SETTEXTURESTAGE(mdl, stage, texHandle)` | Set texture on a render stage. |
| `MODEL.SETSTAGEBLEND(mdl, stage, mode)` | Blend mode for a texture stage. |
| `MODEL.SETSTAGEROTATE(mdl, stage, angle)` | UV rotation for a texture stage. |
| `MODEL.SETSTAGESCALE(mdl, stage, sx, sy)` | UV scale for a texture stage. |
| `MODEL.SETSTAGESCROLL(mdl, stage, ux, uy)` | UV scroll for a texture stage. |
| `MODEL.ROTATETEXTURE(mdl, angle)` | Rotate base texture UVs. |
| `MODEL.SCALETEXTURE(mdl, sx, sy)` | Scale base texture UVs. |
| `MODEL.SCROLLTEXTURE(mdl, ux, uy)` | Scroll base texture UVs per frame. |
| `MODEL.SETBLEND(mdl, mode)` | Alpha blend mode (`NONE`, `ALPHA`, `ADD`). |
| `MODEL.SETCULL(mdl, mode)` | Face culling mode. |
| `MODEL.SETDEPTH(mdl, enabled)` | Enable/disable depth test. |
| `MODEL.SETWIREFRAME(mdl, bool)` | Wireframe rendering. |
| `MODEL.SETLIGHTING(mdl, bool)` | Enable/disable lighting. |
| `MODEL.SETFOG(mdl, bool)` | Enable/disable fog. |
| `MODEL.SETGPUSKINNING(mdl, bool)` | Use GPU vs CPU skinning. |
| `MODEL.SETMODELMESHMATERIAL(mdl, meshIdx, matHandle)` | Override material on a sub-mesh. |
| `MODEL.GETMATERIALCOUNT(mdl)` | Returns number of materials. |
| `MODEL.SETCASTSHADOW(mdl, bool)` | Enable/disable shadow casting. |
| `MODEL.SETRECEIVESHADOW(mdl, bool)` | Enable/disable shadow receiving. |

---

### Animation

| Command | Description |
|--------|-------------|
| `MODEL.PLAY(mdl, name)` | Play animation by name. |
| `MODEL.PLAYIDX(mdl, index)` | Play animation by index. |
| `MODEL.STOP(mdl)` | Stop current animation. |
| `MODEL.LOOP(mdl, bool)` | Enable/disable loop. |
| `MODEL.ISPLAYING(mdl)` | Returns `TRUE` if animation is playing. |
| `MODEL.ANIMDONE(mdl)` | Returns `TRUE` if one-shot animation finished. |
| `MODEL.ANIMCOUNT(mdl)` | Returns total number of animations. |
| `MODEL.ANIMNAME(mdl, idx)` / `MODEL.ANIMNAME$(mdl, idx)` | Returns animation name by index. |
| `MODEL.GETFRAME(mdl)` | Returns current animation frame number. |
| `MODEL.TOTALFRAMES(mdl)` | Returns total frames of current animation. |
| `MODEL.UPDATEANIM(mdl, dt)` | Advance animation manually (if not auto-updated). |
| `MODEL.SETSPEED(mdl, speed)` | Set animation playback speed. |

---

### Limbs (bone accessors)

| Command | Description |
|--------|-------------|
| `MODEL.LIMBCOUNT(mdl)` | Returns number of bones/limbs. |
| `MODEL.LIMBX(mdl, bone, x)` / `LIMBY` / `LIMBZ` | Set local bone position offset. |
| `MODEL.SETLIMBPOS(mdl, bone, x, y, z)` | Set bone world offset. |
| `MODEL.SETLODDISTANCES(mdl, d0, d1, d2)` | Set LOD switch distances. |

---

### Draw variants

| Command | Description |
|--------|-------------|
| `MODEL.DRAWEX(mdl, x,y,z, p,y2,r, sx,sy,sz, r,g,b,a)` | Draw with explicit transform and tint. |
| `MODEL.DRAWWIRES(mdl)` | Draw model as wireframe overlay. |
| `MODEL.INSTANCE(mdl, index)` | Draw a specific instance slot. |
| `MODEL.SETINSTANCEPOS(mdl, index, x,y,z)` | Set world position of instance `index`. |
| `MODEL.SETINSTANCESCALE(mdl, index, sx,sy,sz)` | Set scale of instance `index`. |
| `MODEL.UPDATEINSTANCES(mdl)` | Upload instance transform buffer to GPU. |

---

## See also

- [ANIMATION_3D.md](ANIMATION_3D.md) — skeletal clips: **`MODEL.*`** vs **`ENTITY.*`**
- [MESH.md](MESH.md) — procedural meshes, **`MESH.UPLOAD`**, **`MESH.DRAW`**
- [CAMERA.md](CAMERA.md) — 3D camera
- [LIGHT.md](LIGHT.md) — PBR lighting
- [SHADER.md](SHADER.md) — custom materials via shaders
