# Camera, lights, and rendering (engine-style map)

moonBASIC uses **`CAMERA.*`** (heap **handles**, not integer IDs), **`LIGHT.*`** handles, and **`RENDER.*`** / **`FOG.*`** / **`POST.*`** for atmosphere and post-processing. **CGO** + Raylib required for full behavior.

---

## 1. Camera module

| Concept | moonBASIC |
|--------|-----------|
| Create camera | **`CAMERA.CREATE`** or **`CAMERA.MAKE`** / **`CAM`** — returns a **handle** (`Camera3D`). |
| Set eye position | **`CAMERA.SETPOS`** / **`CAMERA.SETPOSITION`** |
| Set look-at target | **`CAMERA.SETTARGET`** / **`CAMERA.LOOKAT`** |
| Field of view | **`CAMERA.SETFOV`** (vertical FOV in degrees, Raylib-style) |
| Perspective vs orthographic | **`CAMERA.SETMODE`** — **`0`** / **`"perspective"`** or **`1`** / **`"orthographic"`**; or **`CAMERA.SETPROJECTION`** `(cam, 0\|1)`. |
| Follow behind an entity | **`CAMERA.FOLLOWENTITY`** `(camera, entity, dist, height, smooth)` — third-person each frame. (**`CAMERA.FOLLOW`** is a lower-level variant with explicit world target + yaw — 8 args.) |
| Screen shake | **`CAMERA.SHAKE`** `(camera, amount, duration)` |
| World → screen (HUD) | **`CAMERA.PROJECT`** / **`CAMERA.WORLDTOSCREEN`** `(camera, wx, wy, wz)` → 2-float array handle **(x, y)** in screen pixels. |
| Screen → world ray | **`CAMERA.UNPROJECT`** / **`CAMERA.GETRAY`** / **`CAMERA.PICK`** `(camera, screenX, screenY)` → **ray handle**. |

See **[CAMERA.md](CAMERA.md)** for orbit, FPS mode, culling, and 2D cameras.

---

## 2. Light module

| Concept | moonBASIC |
|--------|-----------|
| Point (omni) | **`LIGHT.CREATEPOINT`** `(x, y, z, r, g, b, energy)` |
| Directional (sun) | **`LIGHT.CREATEDIRECTIONAL`** `(dx, dy, dz, r, g, b, energy)` — direction is normalized. |
| Spotlight | **`LIGHT.CREATESPOT`** `(x,y,z, tx,ty,tz, r,g,b, outerConeDeg, energy)` — aim is **position → target**; inner cone is derived. |
| Generic | **`LIGHT.MAKE`** `("point" \| "directional" \| "spot")` then **`LIGHT.SETPOS`**, **`LIGHT.SETDIR`**, **`LIGHT.SETTARGET`**, cones, **`LIGHT.SETRANGE`**, etc. |
| Color | **`LIGHT.SETCOLOR`** — channels **0–1** or **0–255** (same heuristic as elsewhere). |
| Position | **`LIGHT.SETPOSITION`** / **`LIGHT.SETPOS`** |
| On / off | **`LIGHT.ENABLE`** or **`LIGHT.SETSTATE`** `(light, enabled)` |
| Range | **`LIGHT.SETRANGE`** |
| Free | **`LIGHT.FREE`** |

PBR + directional shadows are described in **[LIGHT.md](LIGHT.md)** and the API consistency doc.

---

## 3. Atmosphere, fog, post, screenshots

| Concept | moonBASIC |
|--------|-----------|
| Scene ambient (PBR) | **`RENDER.SETAMBIENT`** `(r, g, b)` or **`(r, g, b, intensityScale)`** — fourth component scales RGB. |
| Fog + density in one call | **`RENDER.SETFOG`** `(r, g, b, start, end, density)` — enables **`FOG.*`**, sets color, range, and **`WORLD.FOGDENSITY`**. You can also use **`FOG.ENABLE`**, **`FOG.SETCOLOR`**, **`FOG.SETRANGE`**, **`WORLD.FOGDENSITY`** separately. |
| Bloom | **`RENDER.SETBLOOM`** `(threshold)` or **`(threshold, intensity)`** — forwards to **`POST.BLOOM`**. |
| Screenshot | **`RENDER.SCREENSHOT`** `(path)` — **`TakeScreenshot`** (Raylib). |
| Forward / deferred pipeline | **`RENDER.SETMODE`** `("forward" \| "deferred")` — internal pipeline flag (see **[RENDER.md](RENDER.md)**). |

**Sky / environment**

- **Procedural sky (time of day):** **`SKY.MAKE`**, **`SKY.DRAW`**, **`SKY.UPDATE`**, … — see **[SKY.md](SKY.md)**.
- **`RENDER.SETSKYBOX`** appears in the manifest for future / asset pipeline use; there is **no** dedicated HDR cubemap loader wired to that name yet. Prefer **`SKY.*`** or load environment data through your content path.

---

## See also

- **[GAME_ENGINE_PATTERNS.md](GAME_ENGINE_PATTERNS.md)** — rays, ambient, shortcuts.
- **[docs/API_CONSISTENCY.md](../API_CONSISTENCY.md)** — generated manifest index.

Refresh generated table: `go run ./tools/apidoc`.
