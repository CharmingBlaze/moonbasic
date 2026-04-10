# Essential Blitz-style API (implemented in moonBASIC)

Quick map from **familiar Blitz3D / DBPro names** to **moonBASIC** commands. Identifiers are case-insensitive; registry uses **`NAMESPACE.NAME`**. Full parity notes: [BLITZ_COMMAND_INDEX.md](BLITZ_COMMAND_INDEX.md), [BLITZ3D.md](BLITZ3D.md).

---

## 1. Entity system (the “Blitz feel”)

Entities use **integer ids** (`entity`) or **`CUBE()`** / **`SPHERE()`** handles with dot methods — see [ENTITY.md](ENTITY.md), [BLITZ3D.md](BLITZ3D.md).

| Concept | moonBASIC |
|--------|-----------|
| **Position** (absolute world) | **`ENTITY.POSITIONENTITY`** / **`PositionEntity`** — also **`ENTITY.SETPOSITION`**. Dot: **`obj.Pos(x,y,z)`**. |
| **Move** (along entity **local** forward/right/up from pitch/yaw) | **`ENTITY.MOVEENTITY`** / **`MoveEntity`** / **`ENTITY.MOVE`** — **not** a world delta; use **`ENTITY.TRANSLATEENTITY`** / **`TranslateEntity`** for world **`dx,dy,dz`**. |
| **Rotate** (absolute euler **radians**) | **`ENTITY.ROTATEENTITY`** / **`RotateEntity`**. |
| **Turn** (relative euler) | **`ENTITY.TURNENTITY`** / **`TurnEntity`**. |
| **Scale** | **`ENTITY.SCALE`** / **`ScaleEntity`**. |
| **Parent** | **`ENTITY.PARENT`** / **`EntityParent`** — child inherits transforms. **`ENTITY.PARENTCLEAR`**, **`GetParent`**. |
| **Color / tint** | **`ENTITY.COLOR`** / **`EntityColor`** — RGB **0–255**; alpha via **`ENTITY.ALPHA`** / **`EntityAlpha`**. |
| **Distance** | **`ENTITY.DISTANCE`** / **`EntityDistance`** (two entity ids). |

---

## 2. Meshes & primitives (persistent handles)

**CPU mesh** handles (**`MESH.*`**) — upload with **`MESH.UPLOAD`**, draw with **`MESH.DRAW`** inside **`CAMERA.Begin`/`End`**. See [MESH.md](MESH.md).

| Desired name | moonBASIC |
|-------------|-----------|
| **Mesh.CreateCube** | **`MESH.CREATECUBE`** — alias of **`MESH.MAKECUBE`**. |
| **Mesh.CreateSphere** | **`MESH.CREATESPHERE`** — alias of **`MESH.MAKESPHERE`**. |
| **Mesh.CreatePlane** | **`MESH.CREATEPLANE`** — alias of **`MESH.MAKEPLANE`**. |
| **Mesh.Load** | **`MESH.LOAD(path)`** — format support follows **Raylib** (e.g. **`.obj`**, **`.gltf`** / **`.glb`**, etc.; see Raylib docs for your build). |

For **scene entities** (persistent like Blitz **CreateCube**), use **`ENTITY.CREATECUBE`** / **`CreateCube`** or **`CUBE()`** — [ENTITY.md](ENTITY.md).

---

## 3. Camera & picking

| Desired | moonBASIC |
|--------|-----------|
| **Pick / ray from screen `(x,y)`** | **`CAMERA.GETRAY(cam, sx, sy)`** or **`CAMERA.PICK`** (alias of **`GETRAY`** in CGO builds) — returns **6-float** origin+dir. Use **`RAY.HIT*_*`**, **`ENTITY.PICK`**, or physics queries — [CAMERA.md](CAMERA.md), [RAYCAST.md](RAYCAST.md). There is **no** single **`PickEntity(x,y)`** that returns an entity id; combine **ray + hit test**. |
| **Project 3D → 2D** | **`CAMERA.WORLDTOSCREEN`** or **`CAMERA.PROJECT`** (alias) — returns **`[sx, sy]`** handle. Blitz-style: **`CameraProject`**. |
| **Camera points at entity** | **`CAMERA.LOOKATENTITY`** / **`CAMERA.POINTATENTITY`** — sets target to entity position. Or **`Camera.LookAt(x,y,z)`** / **`CAMERA.SETTARGET`**. |

---

## 4. 2D / screen space (`CAMERA2D.Begin`/`End` or raw screen)

| Desired | moonBASIC |
|--------|-----------|
| **Sprite.Load** | **`SPRITE.LOAD(path)`** — [SPRITE.md](SPRITE.md). |
| **Sprite.Draw** | **`SPRITE.DRAW`** — frame from **`SPRITE.DEFANIM`** / atlas. |
| **Sprite collision** | **`SPRITE.HIT`** / **`SPRITE.POINTHIT`** — bounding tests; not pixel-perfect by default — [SPRITE.md](SPRITE.md). |
| **Viewport / clip rect** | **`RENDER.SETSCISSOR(x, y, w, h)`** — restricts drawing (minimaps, split screen) — [RENDER.md](RENDER.md). |

---

## 5. Logic & “game juice”

| Desired | moonBASIC |
|--------|-----------|
| **CurveValue** | **`CURVEVALUE(dest, current, speed)`** — DBPro-style approach each call — [GAMEHELPERS.md](GAMEHELPERS.md), [mbgame `register_math.go`](../../runtime/mbgame/register_math.go). |
| **CurveAngle** | **`CURVEANGLE(destDeg, srcDeg, speed)`** — degrees, shortest path — same doc. |
| **Rand(min, max)** | **`RAND(min, max)`** or **`RND(min, max)`** — **inclusive integers**. **`RND(n)`** is still **0..n−1**; **`RNDF(lo, hi)`** is float range — [MATH.md](MATH.md). |

---

## See also

- [PROGRAMMING.md](../PROGRAMMING.md) — main loop, **`TIME.DELTA`**, **`ENTITY.UPDATE`**, **`RENDER.FRAME`**
- [API_CONSISTENCY.md](../API_CONSISTENCY.md) — every registered name
