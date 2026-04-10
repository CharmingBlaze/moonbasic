# Game engine patterns (Blitz-style helpers)

moonBASIC keeps the render loop in **`RENDER.FRAME`**, **`CAMERA.Begin`/`End`**, **`ENTITY.UPDATE`**, etc. These commands focus on **state**, **spatial math**, **assets**, and **common gameplay** so you rarely need to hand-roll the same Raylib math.

**Entity CRUD / transform API** (`ENTITY.LOAD`, **`ENTITY.POSITION`**, **`ENTITY.TURN`** vs **`ENTITY.SETROTATION`**, **`ENTITY.LOOKAT`**, …) — see **[ENTITY.md](ENTITY.md)** (“Entity.* API” table).

**Camera, lights, fog, bloom, screenshots** — **[CAMERA_LIGHT_RENDER.md](CAMERA_LIGHT_RENDER.md)**.

**3D skeletal clips** (entities + **`ENTITY.UPDATE`**, or **`MODEL.LOADANIMATIONS`** / **`MODEL.UPDATEANIM`**) — **[ANIMATION_3D.md](ANIMATION_3D.md)**.

## 1. Collision and raycasting (3D)

| Idea | moonBASIC |
|------|-----------|
| Screen → world ray (“unproject”) | **`CAMERA.UNPROJECT(cam, screenX, screenY)`** — same implementation as **`CAMERA.GETRAY`** / **`CAMERA.PICK`**. |
| Ray vs loaded model | **`RAY.HITMODEL_*`** or **`RAY.INTERSECTSMODEL_*`** (aliases: same ray, same model handle). Use **`_HIT`** for boolean, **`_DISTANCE`** for hit distance along the ray. |
| Per-mesh + matrix (custom pose) | **`RAY.HITMESH_*`** with mesh handle + matrix handle. |
| Distance between entities | **`ENTITY.DISTANCE(e1, e2)`** — also **`MATH.HDIST`** / **`HDIST`** for horizontal (XZ) distance between two points. |

## 2. Texture and material

| Idea | moonBASIC |
|------|-----------|
| Swap maps on a material | **`MATERIAL.SETTEXTURE`** (diffuse, normal, specular, etc. — see manifest map-type constants). |
| Pixelated vs smooth | **`TEXTURE.SETFILTER`** (e.g. point vs bilinear). |
| Tiling | **`TEXTURE.SETWRAP`**. |

## 3. Lights and world ambient

| Idea | moonBASIC |
|------|-----------|
| Point light at position with color + energy | **`LIGHT.CREATEPOINT(x, y, z, r, g, b, energy)`** — returns a light handle. RGB accepts **0–255** or **0.0–1.0** (same heuristic as **`LIGHT.SETCOLOR`**). |
| Generic light | **`LIGHT.MAKE("point")`** / **`"directional"`** / **`"spot"`**, then **`LIGHT.SETPOS`**, **`LIGHT.SETCOLOR`**, **`LIGHT.SETINTENSITY`**. |
| Spotlight aim | **`LIGHT.SETTARGET`** (shadow frustum look-at; spot/directional semantics depend on kind). |
| Scene base level | **`RENDER.SETAMBIENT`** — so unlit areas are not pure black. |

## 4. Sprite animation (2D)

| Idea | moonBASIC |
|------|-----------|
| Manual frame | **`SPRITE.SETFRAME(sprite, frameIndex)`** — clamps to the strip from **`SPRITE.DEFANIM`**. |
| Range playback | **`SPRITE.PLAY(sprite, start, end, speed, loop)`** — **speed** is frames per second; each frame call **`SPRITE.UPDATEANIM(sprite, Time.Delta())`**. |
| Pivot | **`SPRITE.SETORIGIN(sprite, ox, oy)`** — offset in pixels (applied when drawing). |
| Named states / atlas FSM | **`ANIM.DEFINE`**, **`ANIM.UPDATE`**, etc. (see **`SPRITE`** / animation docs). |

## 5. “Magic” math

| Idea | moonBASIC |
|------|-----------|
| Move toward without overshooting | **`MOVE.TOWARD`** — alias of **`MATH.APPROACH`**. |
| Linear interpolation | **`MOVE.LERP`** — alias of **`MATH.LERP`**. |
| Shortest angle delta (degrees) | **`ANGLE.DIFFERENCE`** — alias of **`MATH.ANGLEDIFF`** / **`ANGLEDIFF`**. For radians, **`MATH.ANGLEDIFFRAD`**. |

## 6. Paths and file checks

| Idea | moonBASIC |
|------|-----------|
| Path next to the executable | **`RES.PATH(localPath)`** — absolute paths are left as-is after cleaning. |
| Exists before loading | **`RES.EXISTS(path)`** — same practical use as **`UTIL.FILEEXISTS`** / **`FILEEXISTS`**. |

---

**Refresh generated API tables:** `go run ./tools/apidoc` after manifest changes.
