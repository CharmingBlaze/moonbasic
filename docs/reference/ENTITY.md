# Entity commands

moonBASIC **entities** are lightweight **integer ids** (handles to rows in the host entity store—not heavy “OOP objects”). **`ENTITY.*`** is the canonical prefix; **PascalCase** **`Entity.*`** names in this doc are **facade aliases** where registered (see [`entity_blitz_facade_cgo.go`](../../runtime/mbentity/entity_blitz_facade_cgo.go)). **CGO** builds link the full implementation; see [`runtime/mbentity/`](../../runtime/mbentity/).

**Scene vs level:** **`SCENE.*`** is for **game scene** switching (**mbscene**). Loading glTF levels uses **`LEVEL.LOAD`** / **`LEVEL.FINDENTITY`** (or **`ENTITY.FIND`** by name)—not **`Scene.FindEntity`**.

## “Entity.*” API (engine-style names)

PascalCase names in other engines map to **`ENTITY.*`** as follows. **Rotation naming:** **`ENTITY.ROTATE`** / **`ENTITY.TURN`** add Euler deltas (degrees); **`ENTITY.SETROTATION`** / **`ENTITY.ROTATEENTITY`** set absolute pitch/yaw/roll.

| Concept | moonBASIC |
|--------|-----------|
| Load model from path | **`ENTITY.LOAD(path)`** or **`ENTITY.LOADMESH`** — formats supported by Raylib (e.g. glTF/OBJ/… per build). Optional 2nd arg: parent entity. |
| Cube / sphere primitives | **`ENTITY.CREATECUBE(size)`** / **`ENTITY.CREATEBOX(w, h, d)`** — box; **`ENTITY.CREATESPHERE(radius)`** (default 16 segments) or **`ENTITY.CREATESPHERE(radius, segments)`**. |
| World position | **`ENTITY.POSITION(id, x, y, z)`** or **`ENTITY.SETPOSITION`** — optional 5th arg local vs global (see runtime). |
| Move along local axes | **`ENTITY.MOVE(id, forward, right, up)`** — same as **`MOVEENTITY`**; uses yaw/pitch for forward/right/up. |
| Absolute rotation | **`ENTITY.SETROTATION(id, pitch, yaw, roll)`** — alias of **`ENTITY.ROTATEENTITY`**. |
| Turn (add rotation) | **`ENTITY.TURN`** / **`ENTITY.ROTATE`** / **`TURNENTITY`** — delta pitch/yaw/roll. |
| Scale | **`ENTITY.SCALE(id, sx, sy, sz)`**. |
| Parent / unparent | **`ENTITY.PARENT(child, parent)`**; **`ENTITY.UNPARENT(id)`** — alias of **`ENTITY.PARENTCLEAR`** (keeps world position). |
| Look at world point | **`ENTITY.LOOKAT(id, targetX, targetY, targetZ)`** — sets pitch/yaw toward the point. (**`ENTITY.POINTENTITY`** aims at another **entity** on XZ.) |
| Distance | **`ENTITY.DISTANCE(id1, id2)`** or **`ENTITY.GETDISTANCE`** (alias) — world-space distance between pivots. **`Entity.GetDistance`** / **`EntityDistance`** (Blitz). |
| Type / role (Blender extras) | **`ENTITY.ISTYPE(id, pattern)`** — **`true`** if **`ENTITY`** name, Blender **`tag`**, or metadata keys **`type`**, **`entity_type`**, **`kind`**, **`category`** (case-insensitive) match **`pattern`** (glob: **`*`**, **`?`**). |
| Find by custom property | **`ENTITY.FINDBYPROPERTY(key, value)`** — returns a **numeric array** of entity ids whose **`ENTITY.GETMETADATA`** row matches **`value`** for **`key`** (exact or glob on value). |
| In-world messages | **`ENTITY.SENDMESSAGE(targetId, msg)`** queues a string; **`ENTITY.POLLMESSAGE(id)`** pops one message (FIFO) or **`""`**. Use in your game loop (not automatic networking). |
| Visibility | **`ENTITY.SETVISIBLE(id, toggle)`** — alias of **`ENTITY.VISIBLE`**. |

## Spatial macros (`ENTITY.X`, `ENTITY.Y`, …) and bounds safety

Shorthand **`ENTITY.X(id)`**, **`ENTITY.Y(id)`**, **`ENTITY.Z(id)`**, **`ENTITY.P` / `W` / `YAW` / `R`** compile to fast bytecode that reads/writes the host **SoA spatial buffer** when the full entity runtime is linked.

- **Literal ids**: If **`id`** is a **numeric literal**, the **compiler** rejects negative values and indices **≥ 2²⁴** (`runtime.MaxEntitySpatialIndex`).
- **Dynamic ids**: Non-constant indices are checked **at run time** by the VM (same bounds). In-bounds SoA slots that are **not** active entities produce a clear **`ENTITY:`** error (no silent stale reads/writes).
- **Details**: [COMPILER_SPEC.md](../COMPILER_SPEC.md) · [ARCHITECTURE.md](../../ARCHITECTURE.md) §8.3.

## Quick links

- **3D skeletal clips & unified model API** — [ANIMATION_3D.md](ANIMATION_3D.md) (**`ENTITY.PLAY`** / **`PLAYNAME`**, **`ENTITY.LOADANIMATIONS`**, **`ENTITY.DRAW`**, **`GETBOUNDS`**, **`RAYHIT`**, …).
- **glTF level markers / layers** — [LEVEL.md](LEVEL.md) (**`LEVEL.LOAD`**, **`LEVEL.GETSPAWN`**, **`LEVEL.SHOWLAYER`** — distinct from **`SCENE.*`** game scenes).
- **Blitz-style names** (`PositionEntity`, `CreateSphere`, …) are mapped under **`ENTITY.POSITIONENTITY`**, **`ENTITY.CREATESPHERE`**, etc. — see [`entity_blitz_cgo.go`](../../runtime/mbentity/entity_blitz_cgo.go), the **[Blitz command index](BLITZ_COMMAND_INDEX.md)**, and the concise **[Blitz essential API](BLITZ_ESSENTIAL_API.md)** (Position vs Move, Rotate vs Turn, Parent, Distance, …).
- **Dot-syntax handles** (`cube.Pos`, `sphere.Turn`) use **`ENTITYREF`** from **`CUBE()`** / **`SPHERE()`** — [BLITZ3D.md](BLITZ3D.md).
- **Scene save/load / clear** — [BLITZ2025.md](BLITZ2025.md), **`ENTITY.SAVESCENE`**, **`ENTITY.LOADSCENE`**, **`ENTITY.CLEARSCENE`**.

## Modern Blitz-style shorthands

- **`UpdatePhysics()`** — same as **`UPDATEPHYSICS`** (blitzengine): one call per frame for **`ENTITY.UPDATE(Time.Delta)`** plus best-effort world / 2D / 3D simulation steps; then **`Render.Clear`** → **`RENDER.Begin3D`** / **`DrawEntities`** / **`RENDER.End3D`** → **`Render.Frame`**.
- **`DrawEntities()`** — same as **`ENTITY.DRAWALL`** (scene graph draw pass).
- **`CreatePivot()`** — empty transform node (invisible, for parenting).
- **`CreateCube(...)`** — `CreateCube()` / `CreateCube(w,h,d)` / `CreateCube(parent)` / `CreateCube(parent, w,h,d)`; see [`entity_blitz_cgo.go`](../../runtime/mbentity/entity_blitz_cgo.go).
- **Jolt (Linux+CGO):** **`ENTITY.LINKPHYSBUFFER(entity, bufferIndex)`** ties an entity to a **`BODY3D`** matrix slot (from **`BODY3D.BUFFERINDEX`**). After **`PHYSICS3D.STEP`**, translation from the shared buffer updates the entity pose. **`ENTITY.CLEARPHYSBUFFER(entity)`** removes the link.

## Jolt collision groups, queries, and AI helpers (Linux + **`PHYSICS3D.START`**)

| Command | Purpose |
|--------|---------|
| **`ENTITY.SETCOLLISIONGROUP(id, group)`** | Alias of **`ENTITY.COLLISIONLAYER`** — stores **0..31** for **`PICK.LAYERMASK`** / future simulation filtering. |
| **`ENTITY.CHECKCOLLISION(a, b)`** | Same as **`EntityCollided`** — **`true`** if the pair had a Jolt contact since the last **`PHYSICS3D.STEP`** (requires **`ENTITY.LINKPHYSBUFFER`** on both sides where applicable). |
| **`ENTITY.RAYCAST(ox, oy, oz, dx, dy, dz, maxDist)`** → **entity** | First hit entity along the ray segment (**Jolt** query path shared with **`PICK.*`** / **`PickCastEntityID`**). Returns **0** if none. |
| **`ENTITY.GETGROUNDNORMAL(id)`** → **vec3 handle** | With **`PLAYER.CREATE`**, uses **`CharacterVirtual.GetGroundNormal`**; otherwise a short downward Jolt ray. Fallback normal **`(0,1,0)`** if no hit. |
| **`ENTITY.APPLYIMPULSE(id, fx, fy, fz)`** | Same as **`ENTITY.ADDFORCE`** / **`ApplyEntityImpulse`** (host velocity integration). Not **`BodyInterface::AddImpulse`** until the Jolt C wrapper exposes it. |
| **`ENTITY.CANSEE(observer, target, fovDeg, maxDist)`** → **bool** | Vision cone + line-of-sight: forward from observer eye height (**~1.65**), aim at target at the same offset, **`PickCastEntityID`** along that segment must hit **target** (or no physics hit). |
| **`ENTITY.GETCLOSESTWITHTAG(id, radius, tag)`** → **entity** | Same tag rules as **`PLAYER.GETNEARBY`**, but only the **nearest** match (**0** if none). |
| **`ENTITY.PUSHOUTOFGEOMETRY(id)`** | Best-effort depenetration: nudges world **Y** up slightly; full recovery belongs in Jolt body / character settings when exposed. |
| **`ENTITY.HASTAG(id, pattern)`** → **bool** | **`path.Match`** on **Blender `tag`** or **entity name** only (stricter than **`ENTITY.ISTYPE`**, which also checks metadata **type** fields). **`EntityHasTag`** alias. |
| **`ENTITY.INFRUSTUM(id)`** → **bool** | Same frustum test as **`ENTITY.INVIEW`**, but uses the **active** **`CAMERA.BEGIN`** camera (no camera handle argument). Returns **false** outside a Begin/End 3D block. |
| **`ENTITY.LINEOFSIGHT(observer, target)`** → **bool** | Straight segment from observer eye (~**1.65** m) to target “eye” height — first **Jolt** hit must be **target** (or no physics hit). Does **not** skip trigger/sensor bodies until those use collision layers / masks. |
| **`ENTITY.GETOVERLAPCOUNT(zoneId, tag)`** → **int** | Counts entities matching **`tag`** (same glob rules as **`PLAYER.GETNEARBY`**) whose **pivot** lies inside **zoneId**’s world **AABB** (sphere prefilter + axis test). |
| **`ENTITY.ANIMATETOWARD(id, x, y, z, duration)`** | Linear **world** lerp of the entity root to **(x,y,z)** over **duration** seconds (driven inside **`ENTITY.UPDATE(dt)`**). |

Detailed normals along an arbitrary ray: **`PHYSICS3D.RAYCAST`** (returns a small result array including the surface normal).

## Movement, rule collisions, and space transforms

These globals mirror Blitz-style names; canonical forms are **`ENTITY.*`** / **`MOVEENTITY`** where noted.

### `ENTITY.GETPOS(entity#)` → **handle**

- **Arguments:** `entity` (int entity id).
- **Returns:** a 3-float tuple-like array handle `[x, y, z]` for destructuring.
- **Use case:** convenient one-call position reads in game loops.

```basic
px#, py#, pz# = Entity.GetPos(player)
```

`ENTITY.GETPOSITION(entity)` remains available and returns a vec3 handle for the handle-based vector API.

### `MoveEntity(entity#, forward#, right#, up#)`

- **Arguments:** **`entity`** — entity id; **`forward`**, **`right`**, **`up`** — distances to move along that entity’s **local** axes (from its **pitch** and **yaw**; roll is not used for the basis).
- **Behavior:** Same as **`MOVEENTITY`** and **`ENTITY.MOVE`**. The engine builds a forward vector from yaw/pitch, derives right from the world up cross forward, then up from right cross forward, and adds **`forward·fwd + right·right + up·up`** to the entity’s **world** position (parent-aware).
- **Use for:** Walking relative to facing (e.g. set **`RotateEntity(player, 0, camYaw, 0)`** then **`MoveEntity(player, speed*dt, 0, 0)`** for forward).
- **Not for:** A fixed world offset — use **`TranslateEntity`** instead.

### `TranslateEntity(entity#, dx#, dy#, dz#)`

- **Arguments:** **`entity`**; **`dx`**, **`dy`**, **`dz`** — delta in **world** space (applied to world position, then converted back to local if parented).
- **Behavior:** Same as **`ENTITY.TRANSLATE`** / **`ENTITY.TRANSLATEENTITY`**.
- **Use for:** Nudging lights, props, or anything that should move **`(1,0,0)`** in world axes regardless of rotation.

### `EntityHitsType(entity#, type#)` → **bool**

- **Arguments:** **`entity`** — mover/query entity; **`type`** — integer **collision type** previously set with **`EntityType`** / **`ENTITY.TYPE`** on **other** entities (e.g. ground = **`2`**).
- **Returns:** **`TRUE`** if, **after the last `ENTITY.UPDATE(dt)`** (or **`UPDATEPHYSICS`**), **`entity`** has a rule-based hit whose other body’s **`EntityType`** equals **`type`**. Otherwise **`FALSE`**.
- **Relation to `ENTITYCOLLIDED`:** Same test as **`ENTITYCOLLIDED(entity, type) <> 0`**; **`ENTITYCOLLIDED`** returns the **other entity’s id** or **`0`** if you need the handle.
- **Prerequisites:** Register pairs with **`COLLISIONS(srcType, dstType, method, response)`** (e.g. sphere-vs-box **`method`** **`2`**) and run **`ENTITY.UPDATE`** each frame. **Not** the same as **`EntityCollided(a, b)`**, which is the **two-entity Jolt** contact query (Linux + linked buffers).

### `TFormVector(x#, y#, z#, srcEntity#, dstEntity#)` → **handle**

- **Arguments:** Direction or vector components **`x`**, **`y`**, **`z`** in **`srcEntity`**’s **local** space; **`srcEntity`** and **`dstEntity`** are entity ids.
- **Returns:** **Heap handle** to a **3-element float array** (same convention as **`ENTITY.GETPOSITION`**): read components via array access or helpers your script style supports.
- **Behavior:** Alias of **`ENTITY.TFORMVECTOR`**. Transforms the vector by the **linear** part of the world matrix chain (direction only, no translation).
- **Use for:** Camera-relative directions, wind in ship space, etc. **Note:** There is no **`entity = 0`** “world” shortcut; use an axis-aligned **pivot entity** at the origin if you need world as a space.

## Scene hierarchy & world utilities (Blitz-style)

- **`ENTITY.VISIBLE(entity, visible)`** / **`EntityVisible`** — sets the same flag as **`ENTITY.HIDE`** / **`ENTITY.SHOW`** (`visible` = false hides the entity).
- **`ENTITY.COUNTCHILDREN(parent)`** — number of **direct** children (stable order = reparent / create order).
- **`ENTITY.GETCHILD(parent, index)`** — direct child entity at `index` (0-based).
- **`ENTITY.FINDCHILD(rootEntity, name)`** — breadth-first search **under** `root` (not global; use **`ENTITY.FIND`** for global name lookup). Names come from **`ENTITY.SETNAME`**.
- **`ENTITY.TFORMPOINT(x, y, z, srcEntity, dstEntity)`** / **`ENTITY.TFORMVECTOR(...)`** — same semantics as **`TFormVector`** / **`ENTITY.TFORMVECTOR`** above; **`TFORMPOINT`** includes translation (full matrix); **`TFORMVECTOR`** is direction-only. Returns a **3-float numeric array handle** (same pattern as **`ENTITY.GETPOSITION`**).
- **`ENTITY.DELTAX` / `DELTAY` / `DELTAZ(entityA, entityB)`** — world-space axis delta **B − A** between origins.
- **`ENTITY.MATRIXELEMENT(entity, row, col)`** — one element of the **world** matrix; **row/col 0..3**, **column-major** (same as **`MAT4.GETELEMENT`** / Raylib `rl.Matrix`).
- **`ENTITY.INVIEW(entity, camera)`** — conservative frustum test for the entity bounds vs the given **`CAMERA.MAKE`** handle (aspect from current framebuffer). **`ENTITY.SETCULLMODE`** force visible/hidden still applies first.

## 3D sprites (billboards)

- **`LOADSPRITE(path)`** / **`LOADSPRITE(path, parent)`** — **`ENTITY.LOADSPRITE`** / **`ENTITY.CREATESPRITE`** are aliases; optional **parent** parents the new sprite like **`ENTITY.PARENT`** (child starts at local origin).
- **`SPRITEMODE`** / **`ENTITY.SPRITEVIEWMODE`** / **`SPRITEVIEWMODE`** — **`1`** = Y-axis billboard, **`2`** = full camera-facing billboard, **`3`** = static quad (see implementation in [`entity_cgo.go`](../../runtime/mbentity/entity_cgo.go)).

## Bulk free (`FREEENTITIES` / `ENTITY.FREEENTITIES`)

**`FREEENTITIES(arrayHandle)`** walks a **numeric** entity array (e.g. **`DIM badGuy AS HANDLE(n)`** / integer slots holding entity ids) and calls **`FreeEntity`** on each non-zero entry. Use at shutdown or level unload instead of hand-written **`FOR i = 1 TO n : FreeEntity(...) : NEXT`**.

## Terrain vs entity#

Heightmap **terrain** is a separate **`TERRAIN.*`** heap object ([TERRAIN.md](TERRAIN.md)), not an entity. Use **`Terrain.GetHeight`** for height queries, **`Terrain.Place`** to position an entity on the surface in one call, or **`Terrain.SnapY`** to adjust **Y** only. **`Terrain.Raise`** / **`Terrain.Lower`** edit heights. **Jolt** heightfield shapes for terrain are **not** exposed in the current `jolt-go` binding; physics for terrain remains mesh/other shapes until extended.

## Performance / roadmap notes

- **Stencil mirrors** (reflection planes) are not implemented yet.
- **Heavy billboard counts:** many **`LOADSPRITE`** instances may benefit from future batching in **`ENTITY.DRAWALL`**; profile hot paths first.

## Procedural meshes (`ENTITY.CREATEMESH`)

- **`ENTITY.CREATEMESH`** / **`CreateMesh([parent])`** — allocates a **blank** procedural mesh (no default cube). The entity stays **hidden** until **`UpdateMesh`**. Optional **parent** works like **`CreateCube(parent)`**.
- **`CreateSurface` / `ENTITY.CREATESURFACE`** — returns the **surface handle** (same as the internal **`MeshBuilder`** heap object) used by **`AddVertex`** / **`AddTriangle`**.
- **`AddVertex`**, **`AddTriangle`**, **`UpdateMesh`**, **`VertexX` / `Y` / `Z`** — CPU-side builder + GPU upload via **`LoadModelFromMesh`** (smooth normals from triangle fans). **`ENTITY.FREE`** unloads the model and frees the builder.
- **`EmitSound(sound, entity)`** — plays a one-shot sound with **distance attenuation** and **stereo pan** vs the last **`Listener(cam)`** / **`AUDIO.LISTENERCAMERA`** (see [AUDIO.md](AUDIO.md) spatial notes).

## Skeletal animation, bone sockets, materials (Raylib)

Full **command matrix** (entities vs **`MODEL.*`**), time scaling, and limitations: **[ANIMATION_3D.md](ANIMATION_3D.md)**.

- **`ENTITY.LOADANIMATEDMESH` / `LoadAnimMesh`** — loads a model and **`LoadModelAnimations`**; first pose is applied with **`UpdateModelAnimation`** + **`UpdateModelAnimationBones`**.
- **`ENTITY.ANIMATE` / `Animate`** — **`ENTITY.ANIMATE(entity [, mode, speed])`**: mode **`0`–`1`** = loop, **`2`** = ping-pong, **`3`+** = clamp at end of clip. (Older scripts used mode **`1`** for clamp; use **`3`** now.) **100 ms skeletal cross-fade** between clips is not implemented (single active pose from Raylib).
- **`ENTITY.EXTRACTANIMSEQ` / `ExtractAnimSeq`** — **`(entity, startFrame, endFrame)`** inclusive clip range for the **current** animation; use **`ENTITY.SETANIMINDEX`** to pick which **`ModelAnimation`** is active.
- **`ENTITY.SETANIMINDEX`** — select animation clip index (resets time to 0). **`ENTITY.ANIMINDEX`** / **`ENTITY.ANIMCOUNT`** — read active index and loaded clip count.
- **`ENTITY.FINDBONE` / `FindBone`** — returns a **hidden** entity whose **world matrix** tracks a named bone on the host model each frame (parent props with **`ENTITY.LOADMESH(path, parent)`** or **`ENTITY.PARENT`**). If the host model is freed, sockets are invalidated.
- **`ENTITY.SETANIMTIME` / `SetAnimTime`**, **`ENTITY.ANIMTIME` / `EntityAnimTime`** — continuous animation time (not always an integer frame index).
- **Brushes:** **`CreateBrush`**, **`BrushTexture`**, **`BrushFX`**, **`BrushShininess`**, **`PaintEntity`** — heap **`Brush`** handle; **`PaintEntity`** copies color/texture/FX/shininess onto the entity. **`BrushFX`**: bit **`1`** = full-bright tint boost in **`ENTITY.DRAWALL`**, bit **`16`** = additive blending (Raylib **`BlendAdditive`**). Full PBR/shader swaps are future work.
- **`EntityShadow`** — stores **`shadowCast`** on the entity (**`0`** default, **`1`** / **`2`** reserved); hooking into the deferred shadow pass for **`ENTITY.DRAWALL`** models is not wired yet.

## Reference tables

- **[API_CONSISTENCY.md](../API_CONSISTENCY.md)** — search for **`ENTITY.`** for every overload and arity.
- **[ANIMATION_3D.md](ANIMATION_3D.md)** — skeletal clips, **`ENTITY.UPDATE`**, bone sockets.
- **[GAMEHELPERS.md](GAMEHELPERS.md)** — movement, landing, camera follow.
