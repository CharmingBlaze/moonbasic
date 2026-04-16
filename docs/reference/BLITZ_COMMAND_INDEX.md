# Blitz-style command index (moonBASIC mapping)

Classic **Blitz3D / BlitzPlus** used short globals (`Graphics3D`, `CreateCube`, `Line`, …). moonBASIC is **not** a byte-compatible runtime: rendering uses **Raylib** (GPU), files and types differ, and APIs are usually **`NAMESPACE.NAME`** (identifiers are case-insensitive — see [LANGUAGE.md](../LANGUAGE.md)).

This page maps **familiar Blitz names** to **implemented** moonBASIC commands, with **parity notes**. For narrative “how to think in moonBASIC”, see [BLITZ3D.md](BLITZ3D.md) and [BLITZ2025.md](BLITZ2025.md). The full registry is [`commands.json`](../../compiler/builtinmanifest/commands.json); human list: [API_CONSISTENCY.md](../API_CONSISTENCY.md).

**Legend:** **✓** close match · **≈** same role, different args or workflow · **—** no direct equivalent (use the suggested alternative or script a helper).

---

## Registry aliases (same implementation, Blitz-friendly names)

| Command | Same handler as |
|---------|-----------------|
| **`DRAW.PLOT`** | **`DRAW.PIXEL`** — **6 args** `(x, y, r, g, b, a)` |
| **`DRAW.OVAL`** | **`DRAW.ELLIPSE`** — **8 args** `(cx, cy, rx, ry, r, g, b, a)` |
| **`DRAW.OVALLINES`** | **`DRAW.ELLIPSELINES`** |
| **`LOADTEXTURE`** | **`TEXTURE.LOAD`** |
| **`FREETEXTURE`** | **`TEXTURE.FREE`** |
| **`TEXTUREWIDTH`** | **`TEXTURE.WIDTH`** |
| **`TEXTUREHEIGHT`** | **`TEXTURE.HEIGHT`** |

Also registered: **`OPENFILE`**, **`CLOSEFILE`**, **`READFILE`**, **`SEEKFILE`**, **`FILEPOS`**, **`FILESIZE`**, **`KEYDOWN`**, **`KEYHIT`**, **`MOUSEX`**, … — see [FILE.md](FILE.md), [INPUT.md](INPUT.md), [GAMEHELPERS.md](GAMEHELPERS.md).

---

## 2D drawing (BlitzPlus-style buffers)

Blitz **2D** used a software framebuffer, **current pen color**, and **Origin/Viewport**. moonBASIC **2D** draws in **screen or Camera2D space** with **per-call RGBA** (no hidden pen stack). Off-screen drawing uses **`RENDERTARGET.*`** or **`IMAGE.*`** CPU surfaces — see [DRAW2D.md](DRAW2D.md), [RENDER.md](RENDER.md), [IMAGE.md](IMAGE.md).

| Blitz-style | Call in moonBASIC | Notes |
|-------------|-------------------|--------|
| **Plot** | **`DRAW.PLOT`** or **`DRAW.PIXEL`** | **6 args** — set colour every call. |
| **Line** | **`DRAW.LINE`** | **8 args** including **`r,g,b,a`**. |
| **Rect** | **`DRAW.RECTANGLE`** / **`DRAW.RECTLINES`** | Filled vs outline. |
| **Oval** | **`DRAW.OVAL`** / **`DRAW.OVALLINES`** | Same as ellipse; from a Blitz **w,h box**: **`cx=x+w/2`**, **`cy=y+h/2`**, **`rx=w/2`**, **`ry=h/2`**. |
| **Text** | **`DRAW.TEXT`** | [DRAW2D.md](DRAW2D.md) |
| **Color** | *(locals / per-call)* | No global pen. |
| **Origin** | **`CAMERA2D.*`** | Offset / target / zoom / rotation. |
| **Viewport** | **`RENDER.SETSCISSOR`** or **`RENDERTARGET.*`** | |
| **SetBuffer** | **`RENDERTARGET.BEGIN`** / **`END`** | |
| **CopyRect** / pixels | **`IMAGE.*`**, **`DRAW.TEXTURE*`** | [IMAGE.md](IMAGE.md) |

---

## World / scene

| Blitz-style | moonBASIC | Notes |
|-------------|-----------|--------|
| **CreateCamera(parent)** | **≈** **`CAMERA.CREATE`** (deprecated `CAMERA.MAKE`) | Returns **camera handle**; parenting differs from Blitz entity graph. |
| **CreateLight(type, parent)** | **≈** **`LIGHT.CREATE`** / deprecated **`LIGHT.MAKE`** + **`LIGHT.SET*`** | See [LIGHT.md](LIGHT.md). |
| **CreatePivot(parent)** | **—** | Use **empty entity** / **group** patterns or **`ENTITY.PARENT`**. |
| **CreateListener(parent)** | **—** | No audio listener entity; use **spatial** audio APIs if exposed. |
| **RenderWorld** | **—** *(intentional)* | Not a builtin — use **`CAMERA.Begin`/`End`** or **`RENDER.Begin3D`/`End3D`**, then **`ENTITY.DRAWALL`**, then **`RENDER.FRAME`**. See [BLITZ3D.md](BLITZ3D.md) § Raylib render pipeline. |
| **UpdateWorld** | **—** *(intentional)* | Not a builtin — use **`ENTITY.UPDATE(dt)`** with **`TIME.DELTA()`**. See [BLITZ3D.md](BLITZ3D.md). |
| **Flip** | **—** *(intentional)* | Not a builtin — use **`RENDER.FRAME`** to present. |
| **ClearWorld(…)** | **≈** `ENTITY.CLEARSCENE` / `SCENE.CLEARSCENE` | See [BLITZ2025.md](BLITZ2025.md). |
| **AmbientLight(r, g, b)** | **≈** `RENDER.SETAMBIENT` / `FOG.SETCOLOR` context | Plus **`RENDER.CLEAR`** for sky colour. |
| **FogColor** / **FogRange** | **✓** `FOG.SETCOLOR`, `FOG.SETRANGE` / `SETNEAR`+`SETFAR` | [BLITZ2025.md](BLITZ2025.md) |
| **WireFrame(enable)** | **≈** `RENDER.SETWIREFRAME` | |
| **Dither(enable)** | **—** | Use **`IMAGE.DITHER`** on CPU images if needed. |
| **AntiAlias(enable)** | **≈** `RENDER.SETMSAA` | Platform/window dependent. |
| **TrisRendered()** | **—** | No built-in triangle counter; use profiler / GPU tools. |

---

## Entity commands

moonBASIC uses **`ENTITY.*`** with **integer entity ids** (and optional **`ENTITYREF`** dot-syntax from **`CUBE()`** / **`SPHERE()`** — [BLITZ3D.md](BLITZ3D.md)). Many **Blitz command names** are registered as **`ENTITY.POSITIONENTITY`**, **`ENTITY.ROTATEENTITY`**, etc. — see [`entity_blitz_cgo.go`](../../runtime/mbentity/entity_blitz_cgo.go).

| Blitz-style | moonBASIC (typical) | Notes |
|-------------|---------------------|--------|
| **CreateCube(parent)** | **✓** `ENTITY.CREATECUBE` / **`CUBE()`** | Parent: **`ENTITY.PARENT`**. |
| **CreateSphere(segments, parent)** | **✓** `ENTITY.CREATESPHERE` / **`SPHERE(r)`** | Arity differs — see handler. |
| **CreateCylinder** / **CreateCone** | **✓** `ENTITY.CREATECYLINDER` | Cone: **mesh** path **`MESH.MAKECONE`** + entity attach pattern. |
| **CreatePlane(divisions, parent)** | **✓** `ENTITY.CREATEPLANE` | |
| **CreateMesh(parent)** | **✓** `ENTITY.CREATEMESH` | |
| **CreateBrush(r, g, b)** | **—** | Use **material** / **`ENTITY.COLOR`** workflows. |
| **CreateSurface(mesh)** | **—** | Raylib **mesh** is internal; use **`MESH.*`** procedural APIs. |
| **CreateSprite(parent)** | **≈** `SPRITE` / sprite system | See [SPRITE.md](SPRITE.md). |
| **CreateMirror(parent)** | **—** | Planar reflection not a single command. |
| **CopyEntity** / **FreeEntity** | **✓** `ENTITY.COPY` / `ENTITY.FREE` | |
| **PositionEntity** … **AlignToVector** | **✓** `ENTITY.POSITIONENTITY`, `ENTITY.ROTATEENTITY`, `ENTITY.TURNENTITY`, `ENTITY.MOVEENTITY`, `ENTITY.TRANSLATEENTITY`, `ENTITY.POINTENTITY`, `ENTITY.ALIGNTOVECTOR` | |
| **EntityX** … **EntityRoll** | **✓** `ENTITY.ENTITYX` … `ENTITY.ENTITYROLL` | Optional **global** flag where implemented. |
| **EntityScaleX/Y/Z** | **≈** `ENTITY.SCALE` | Often one **`SCALE`** with 3 components. |
| **EntityParent** / **GetParent** | **✓** `ENTITY.PARENT` / **`ENTITY.PARENTCLEAR`** | **GetParent**: use stored id or scene data. |
| **HideEntity** / **ShowEntity** | **✓** `ENTITY.HIDE` / `ENTITY.SHOW` | |
| **EntityVisible** / **EntityInView** | **≈** `ENTITY.*` + **`CAMERA.ISONSCREEN`** | |
| **EntityOrder** / **EntityAutoFade** / **EntityFX** / **EntityAlpha** / **EntityColor** / **EntityShininess** / **EntityTexture** / **EntityBlend** / **EntityType** / **EntityPickMode** / **EntityRadius** / **EntityBox** | **✓** `ENTITY.ORDER`, `ENTITY.ALPHA`, `ENTITY.COLOR`, `ENTITY.SHININESS`, `ENTITY.TEXTURE`, `ENTITY.BLEND`, `ENTITY.TYPE`, `ENTITY.PICKMODE`, `ENTITY.RADIUS`, `ENTITY.BOX` | Names are **`ENTITY.*`** in moonBASIC. |
| **Collisions(src, dst, method, response)** | **≈** `ENTITY.COLLIDE` / physics **`PHYSICS3D.*`** | Not a 1:1 Blitz collision stack; see [COLLISION.md](COLLISION.md), [PHYSICS3D.md](PHYSICS3D.md). |
| **EntityCollided** / **CountCollisions** / **CollisionX** … **CollisionTriangle** | **≈** `ENTITY.COLLIDED`, `ENTITY.COLLISIONX` … | See entity module; full Blitz collision index may differ. |

---

## Camera

| Blitz-style | moonBASIC | Notes |
|-------------|-----------|--------|
| **CameraRange** | **≈** `CAMERA.SETRANGE` | Near/far planes. |
| **CameraZoom** | **✓** `CAMERA.ZOOM` | FOV-based in moonBASIC — see [CAMERA.md](CAMERA.md). |
| **CameraFogMode/Color/Range** | **≈** `FOG.*` + camera clear | Global fog helpers. |
| **CameraClsColor** / **ClsMode** | **≈** `RENDER.CLEAR` | Clear colour + depth. |
| **CameraProjMode** | **≈** `CAMERA.SETPROJECTION` | |
| **CameraViewport** | **≈** scissor / render target | |
| **CameraPick** / **CameraProject** | **≈** `CAMERA.GETRAY` / `CAMERA.WORLDTOSCREEN` | |
| **ProjectedX/Y/Z** | **≈** `CAMERA.WORLDTOSCREEN` return values | Check arity in manifest. |

---

## Lights

| Blitz-style | moonBASIC | Notes |
|-------------|-----------|--------|
| **LightColor** | **✓** `LIGHT.SETCOLOR` | |
| **LightRange** | **≈** `LIGHT.SETRANGE` / intensity | |
| **LightConeAngles** | **≈** `LIGHT.SETINNERCONE` / `SETOUTERCONE` | See [LIGHT.md](LIGHT.md). |
| **LightMesh** | **—** | No light-volume mesh helper. |

---

## Meshes and surfaces

| Blitz-style | moonBASIC | Notes |
|-------------|-----------|--------|
| **LoadMesh** / **LoadAnimMesh** | **✓** `ENTITY.LOADMESH` / `ENTITY.LOADANIMATEDMESH` | Also **`MESH.LOAD`**. |
| **FitMesh** / **FlipMesh** / **PaintMesh** / **ScaleMesh** / **RotateMesh** / **PositionMesh** | **≈** `MESH.*` + transforms | See [MESH.md](MESH.md). |
| **MeshWidth** / **Height** / **Depth** | **≈** `MESH.GETBBOX*` | Bounding box extents. |
| **MeshesIntersect** | **—** | Use **physics** overlap or custom test. |
| **CountSurfaces** / **GetSurface** | **—** | Raylib **`Mesh`** exposes **`MESH.TRIANGLECOUNT`** / **`VERTEXCOUNT`**; no Blitz “surface” split. |
| **AddVertex** / **AddTriangle** / **Vertex*** / **TriangleVertex** | **≈** `MESH.MAKECUSTOM` / `UPDATEVERTEX` | Different model — see [MESH.md](MESH.md). |

---

## Textures

| Blitz-style | Call in moonBASIC | Notes |
|-------------|-------------------|--------|
| **LoadTexture** | **`LOADTEXTURE`** or **`TEXTURE.LOAD`** | Alias — same handler. [TEXTURE.md](TEXTURE.md) |
| **FreeTexture** | **`FREETEXTURE`** or **`TEXTURE.FREE`** | Alias. |
| **CreateTexture** | **`IMAGE.CREATE`** + **`TEXTURE.FROMIMAGE`** | Or load from file with **`LOADTEXTURE`**. |
| **TextureWidth** / **Height** | **`TEXTUREWIDTH`** / **`TEXTUREHEIGHT`** or **`TEXTURE.WIDTH`** / **`HEIGHT`** | Aliases. |
| **TextureBlend** / UV transforms | **`TEXTURE.SETFILTER`**, **`SETWRAP`**, **`DRAW.TEXTUREPRO`** | No Blitz-style matrix stack. |

---

## Sound

| Blitz-style | moonBASIC | Notes |
|-------------|-----------|--------|
| **LoadSound** / **FreeSound** | **✓** `AUDIO.LOADSOUND` / `AUDIO.STOP`+unload patterns | See [AUDIO.md](AUDIO.md). |
| **Load3DSound** | **—** | Use **spatial** APIs if listed in manifest. |
| **PlaySound** / **LoopSound** | **✓** `AUDIO.PLAY` | |
| **StopChannel** / **PauseChannel** / **ResumeChannel** | **≈** `AUDIO.STOP` / `PAUSE` / `RESUME` | |
| **ChannelPitch** / **Volume** / **Pan** | **≈** `AUDIO.SETSOUNDPITCH` / `SETSOUNDVOLUME` / `SETSOUNDPAN` | |
| **ChannelPlaying** | **≈** `AUDIO.ISSOUNDPLAYING` | |
| **EmitSound(sound, entity)** | **—** | Position sound by **3D coords** or future spatial helper. |

---

## Input

| Blitz-style | Call in moonBASIC | Notes |
|-------------|-------------------|--------|
| **KeyDown** / **KeyHit** | **`KEYDOWN`** / **`KEYHIT`** (shortcuts) or **`INPUT.KEYDOWN`** / **`INPUT.KEYHIT`** or **`KEY.DOWN`** / **`KEY.HIT`** | [INPUT.md](INPUT.md) |
| **GetKey** | Poll **`KEYDOWN`** / **`KEYHIT`** | No blocking **`WaitKey`**. |
| **MouseX/Y/Z** | **`INPUT.MOUSEX`**, **`MOUSEY`**, wheel | Flat **`MOUSEX`** may exist via shortcuts — check manifest. |
| **MouseHit** / **MouseDown** | **`INPUT.MOUSEHIT`**, **`INPUT.MOUSEDOWN`** | |
| **MoveMouse** | **`INPUT.SETMOUSEPOS`** | |
| **FlushKeys** | *(n/a)* | Input is polled each frame. |

---

## File I/O

Use **`FILE.OPEN`** / **`OPENFILE`**, **`FILE.CLOSE`** / **`CLOSEFILE`**, **`FILE.READLINE`** / **`READFILE`**, **`FILE.WRITE`**, **`FILE.SEEK`** / **`SEEKFILE`**, **`FILE.TELL`** / **`FILEPOS`**, **`FILE.SIZE`** / **`FILESIZE`**, **`READALLTEXT`**, **`WRITEALLTEXT`**, … — full list in [FILE.md](FILE.md) and [API_CONSISTENCY.md](../API_CONSISTENCY.md). Binary typed reads/writes: **`FILE.READBYTE`** etc. if present in manifest.

---

## Math

| Blitz-style | moonBASIC | Notes |
|-------------|-----------|--------|
| **Sin** … **ATan2** | **✓** `SIN`, `COS`, `TAN`, `ASIN`, `ACOS`, `ATAN`, `ATAN2` | Also **`MATH.SIN`** etc. |
| **Sqr** / **Floor** / **Ceil** / **Abs** | **✓** `SQR`, `FLOOR`, `CEIL`, `ABS` | |
| **Rnd** / **SeedRnd** | **✓** `RND`, `RANDOMIZE` / `RNDSEED` | |
| **Rand(low, high)** | **≈** `RND` with range | See [MATH.md](MATH.md). |

---

## Wrappers and aliases

- **Object-style draw helpers** (short **`.Draw`**, **`.Pos`**) are optional — [DRAW_WRAPPERS.md](DRAW_WRAPPERS.md).
- **Entity dot-syntax** on **`CUBE()`** / **`SPHERE()`** maps to **`ENTITY.*`** — [BLITZ3D.md](BLITZ3D.md).
- New **flat** Blitz globals (e.g. a second `LINE` that shadows user code) are intentionally avoided; prefer **`DRAW.*`** and **`ENTITY.*`** as documented.

---

## See also

- [API_CONSISTENCY.md](../API_CONSISTENCY.md) — every registered name  
- [COMMAND_AUDIT.md](../COMMAND_AUDIT.md) — namespace → doc coverage  
- [MISSING_COMMANDS.md](../MISSING_COMMANDS.md) — engine gaps  
- [DarkBASIC Professional map (dbpro/)](dbpro/README.md) — DBPro-style command list in modular files  
