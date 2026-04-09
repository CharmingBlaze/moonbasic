# Entity commands

moonBASIC **entities** are lightweight **integer ids** (physics/visual game objects) with **`ENTITY.*`** builtins. **CGO** builds link the full implementation; see the registry in [`runtime/mbentity/`](../../runtime/mbentity/).

## Quick links

- **Blitz-style names** (`PositionEntity`, `CreateSphere`, …) are mapped under **`ENTITY.POSITIONENTITY`**, **`ENTITY.CREATESPHERE`**, etc. — see [`entity_blitz_cgo.go`](../../runtime/mbentity/entity_blitz_cgo.go) and the **[Blitz command index](BLITZ_COMMAND_INDEX.md)**.
- **Dot-syntax handles** (`cube.Pos`, `sphere.Turn`) use **`ENTITYREF`** from **`CUBE()`** / **`SPHERE()`** — [BLITZ3D.md](BLITZ3D.md).
- **Scene save/load / clear** — [BLITZ2025.md](BLITZ2025.md), **`ENTITY.SAVESCENE`**, **`ENTITY.LOADSCENE`**, **`ENTITY.CLEARSCENE`**.

## Modern Blitz-style shorthands

- **`UpdatePhysics()`** — same as **`UPDATEPHYSICS`** (blitzengine): one call per frame for **`ENTITY.UPDATE(Time.Delta)`** plus best-effort world / 2D / 3D simulation steps; then **`Render.Clear`** → **`RENDER.Begin3D`** / **`DrawEntities`** / **`RENDER.End3D`** → **`Render.Frame`**.
- **`DrawEntities()`** — same as **`ENTITY.DRAWALL`** (scene graph draw pass).
- **`CreatePivot()`** — empty transform node (invisible, for parenting).
- **`CreateCube(...)`** — `CreateCube()` / `CreateCube(w,h,d)` / `CreateCube(parent#)` / `CreateCube(parent#, w,h,d)`; see [`entity_blitz_cgo.go`](../../runtime/mbentity/entity_blitz_cgo.go).
- **Jolt (Linux+CGO):** **`ENTITY.LINKPHYSBUFFER(entity#, bufferIndex#)`** ties an entity to a **`BODY3D`** matrix slot (from **`BODY3D.BUFFERINDEX`**). After **`PHYSICS3D.STEP`**, translation from the shared buffer updates the entity pose. **`ENTITY.CLEARPHYSBUFFER(entity#)`** removes the link.

## Movement, rule collisions, and space transforms

These globals mirror Blitz-style names; canonical forms are **`ENTITY.*`** / **`MOVEENTITY`** where noted.

### `MoveEntity(entity#, forward#, right#, up#)`

- **Arguments:** **`entity#`** — entity id; **`forward#`**, **`right#`**, **`up#`** — distances to move along that entity’s **local** axes (from its **pitch** and **yaw**; roll is not used for the basis).
- **Behavior:** Same as **`MOVEENTITY`** and **`ENTITY.MOVE`**. The engine builds a forward vector from yaw/pitch, derives right from the world up cross forward, then up from right cross forward, and adds **`forward·fwd + right·right + up·up`** to the entity’s **world** position (parent-aware).
- **Use for:** Walking relative to facing (e.g. set **`RotateEntity(player, 0, camYaw, 0)`** then **`MoveEntity(player, speed*dt, 0, 0)`** for forward).
- **Not for:** A fixed world offset — use **`TranslateEntity`** instead.

### `TranslateEntity(entity#, dx#, dy#, dz#)`

- **Arguments:** **`entity#`**; **`dx#`**, **`dy#`**, **`dz#`** — delta in **world** space (applied to world position, then converted back to local if parented).
- **Behavior:** Same as **`ENTITY.TRANSLATE`** / **`ENTITY.TRANSLATEENTITY`**.
- **Use for:** Nudging lights, props, or anything that should move **`(1,0,0)`** in world axes regardless of rotation.

### `EntityHitsType(entity#, type#)` → **bool**

- **Arguments:** **`entity#`** — mover/query entity; **`type#`** — integer **collision type** previously set with **`EntityType`** / **`ENTITY.TYPE`** on **other** entities (e.g. ground = **`2`**).
- **Returns:** **`TRUE`** if, **after the last `ENTITY.UPDATE(dt)`** (or **`UPDATEPHYSICS`**), **`entity#`** has a rule-based hit whose other body’s **`EntityType`** equals **`type#`**. Otherwise **`FALSE`**.
- **Relation to `ENTITYCOLLIDED`:** Same test as **`ENTITYCOLLIDED(entity#, type#) <> 0`**; **`ENTITYCOLLIDED`** returns the **other entity’s id** or **`0`** if you need the handle.
- **Prerequisites:** Register pairs with **`COLLISIONS(srcType, dstType, method, response)`** (e.g. sphere-vs-box **`method`** **`2`**) and run **`ENTITY.UPDATE`** each frame. **Not** the same as **`EntityCollided(a, b)`**, which is the **two-entity Jolt** contact query (Linux + linked buffers).

### `TFormVector(x#, y#, z#, srcEntity#, dstEntity#)` → **handle**

- **Arguments:** Direction or vector components **`x#`**, **`y#`**, **`z#`** in **`srcEntity#`**’s **local** space; **`srcEntity#`** and **`dstEntity#`** are entity ids.
- **Returns:** **Heap handle** to a **3-element float array** (same convention as **`ENTITY.GETPOSITION`**): read components via array access or helpers your script style supports.
- **Behavior:** Alias of **`ENTITY.TFORMVECTOR`**. Transforms the vector by the **linear** part of the world matrix chain (direction only, no translation).
- **Use for:** Camera-relative directions, wind in ship space, etc. **Note:** There is no **`entity# = 0`** “world” shortcut; use an axis-aligned **pivot entity** at the origin if you need world as a space.

## Scene hierarchy & world utilities (Blitz-style)

- **`ENTITY.VISIBLE(entity#, visible)`** / **`EntityVisible`** — sets the same flag as **`ENTITY.HIDE`** / **`ENTITY.SHOW`** (`visible` = false hides the entity).
- **`ENTITY.COUNTCHILDREN(parent#)`** — number of **direct** children (stable order = reparent / create order).
- **`ENTITY.GETCHILD(parent#, index#)`** — direct child entity# at `index#` (0-based).
- **`ENTITY.FINDCHILD(rootEntity#, name$)`** — breadth-first search **under** `root` (not global; use **`ENTITY.FIND`** for global name lookup). Names come from **`ENTITY.SETNAME`**.
- **`ENTITY.TFORMPOINT(x#, y#, z#, srcEntity#, dstEntity#)`** / **`ENTITY.TFORMVECTOR(...)`** — same semantics as **`TFormVector`** / **`ENTITY.TFORMVECTOR`** above; **`TFORMPOINT`** includes translation (full matrix); **`TFORMVECTOR`** is direction-only. Returns a **3-float numeric array handle** (same pattern as **`ENTITY.GETPOSITION`**).
- **`ENTITY.DELTAX` / `DELTAY` / `DELTAZ(entityA#, entityB#)`** — world-space axis delta **B − A** between origins.
- **`ENTITY.MATRIXELEMENT(entity#, row, col)`** — one element of the **world** matrix; **row/col 0..3**, **column-major** (same as **`MAT4.GETELEMENT`** / Raylib `rl.Matrix`).
- **`ENTITY.INVIEW(entity#, camera)`** — conservative frustum test for the entity bounds vs the given **`CAMERA.MAKE`** handle (aspect from current framebuffer). **`ENTITY.SETCULLMODE`** force visible/hidden still applies first.

## 3D sprites (billboards)

- **`LOADSPRITE(path$)`** / **`LOADSPRITE(path$, parent#)`** — **`ENTITY.LOADSPRITE`** / **`ENTITY.CREATESPRITE`** are aliases; optional **parent#** parents the new sprite like **`ENTITY.PARENT`** (child starts at local origin).
- **`SPRITEMODE`** / **`ENTITY.SPRITEVIEWMODE`** / **`SPRITEVIEWMODE`** — **`1`** = Y-axis billboard, **`2`** = full camera-facing billboard, **`3`** = static quad (see implementation in [`entity_cgo.go`](../../runtime/mbentity/entity_cgo.go)).

## Terrain vs entity#

Heightmap **terrain** is a separate **`TERRAIN.*`** heap object ([TERRAIN.md](TERRAIN.md)), not an entity#. Use **`TERRAIN.GETHEIGHT`** for height queries and **`TERRAIN.RAISE`** / **`TERRAIN.LOWER`** for edits. **Jolt** heightfield shapes for terrain are **not** exposed in the current `jolt-go` binding; physics for terrain remains mesh/other shapes until extended.

## Performance / roadmap notes

- **Stencil mirrors** (reflection planes) are not implemented yet.
- **Heavy billboard counts:** many **`LOADSPRITE`** instances may benefit from future batching in **`ENTITY.DRAWALL`**; profile hot paths first.

## Procedural meshes (`ENTITY.CREATEMESH`)

- **`ENTITY.CREATEMESH`** / **`CreateMesh([parent#])`** — allocates a **blank** procedural mesh (no default cube). The entity stays **hidden** until **`UpdateMesh`**. Optional **parent#** works like **`CreateCube(parent#)`**.
- **`CreateSurface` / `ENTITY.CREATESURFACE`** — returns the **surface handle** (same as the internal **`MeshBuilder`** heap object) used by **`AddVertex`** / **`AddTriangle`**.
- **`AddVertex`**, **`AddTriangle`**, **`UpdateMesh`**, **`VertexX` / `Y` / `Z`** — CPU-side builder + GPU upload via **`LoadModelFromMesh`** (smooth normals from triangle fans). **`ENTITY.FREE`** unloads the model and frees the builder.
- **`EmitSound(sound, entity#)`** — plays a one-shot sound with **distance attenuation** and **stereo pan** vs the last **`Listener(cam)`** / **`AUDIO.LISTENERCAMERA`** (see [AUDIO.md](AUDIO.md) spatial notes).

## Skeletal animation, bone sockets, materials (Raylib)

- **`ENTITY.LOADANIMATEDMESH` / `LoadAnimMesh`** — loads a model and **`LoadModelAnimations`**; first pose is applied with **`UpdateModelAnimation`** + **`UpdateModelAnimationBones`**.
- **`ENTITY.ANIMATE` / `Animate`** — **`ENTITY.ANIMATE(entity [, mode, speed])`**: mode **`0`–`1`** = loop, **`2`** = ping-pong, **`3`+** = clamp at end of clip. (Older scripts used mode **`1`** for clamp; use **`3`** now.) **100 ms skeletal cross-fade** between clips is not implemented (single active pose from Raylib).
- **`ENTITY.EXTRACTANIMSEQ` / `ExtractAnimSeq`** — **`(entity#, startFrame#, endFrame#)`** inclusive clip range for the **current** animation; use **`ENTITY.SETANIMINDEX`** to pick which **`ModelAnimation`** is active.
- **`ENTITY.SETANIMINDEX`** — select animation clip index (resets time to 0).
- **`ENTITY.FINDBONE` / `FindBone`** — returns a **hidden** entity# whose **world matrix** tracks a named bone on the host model each frame (parent props with **`ENTITY.LOADMESH(path$, parent#)`** or **`ENTITY.PARENT`**). If the host model is freed, sockets are invalidated.
- **`ENTITY.SETANIMTIME` / `SetAnimTime`**, **`ENTITY.ANIMTIME` / `EntityAnimTime`** — continuous animation time (not always an integer frame index).
- **Brushes:** **`CreateBrush`**, **`BrushTexture`**, **`BrushFX`**, **`BrushShininess`**, **`PaintEntity`** — heap **`Brush`** handle; **`PaintEntity`** copies color/texture/FX/shininess onto the entity. **`BrushFX`**: bit **`1`** = full-bright tint boost in **`ENTITY.DRAWALL`**, bit **`16`** = additive blending (Raylib **`BlendAdditive`**). Full PBR/shader swaps are future work.
- **`EntityShadow`** — stores **`shadowCast`** on the entity (**`0`** default, **`1`** / **`2`** reserved); hooking into the deferred shadow pass for **`ENTITY.DRAWALL`** models is not wired yet.

## Reference tables

- **[API_CONSISTENCY.md](../API_CONSISTENCY.md)** — search for **`ENTITY.`** for every overload and arity.
- **[GAMEHELPERS.md](GAMEHELPERS.md)** — movement, landing, camera follow.
