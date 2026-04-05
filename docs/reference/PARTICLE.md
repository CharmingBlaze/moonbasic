# Particles — `PARTICLE.*`, `PARTICLE3D.*`, and `PARTICLE2D.*`

**3D** billboard particles are implemented in **`runtime/mbparticles`** (CPU simulation + **`DrawBillboard`**). The canonical registry keys are **`PARTICLE.*`**. **`PARTICLE3D.*`** is a **full alias set** — same handlers and handle type (**`Particle`**) — for documentation and DarkBASIC-style naming.

**2D** CPU circles use **`PARTICLE2D.*`** (see [SPRITE.md](SPRITE.md)).

---

## `PARTICLE.*` / `PARTICLE3D.*` (3D)

| Command | Notes |
|---------|--------|
| **`PARTICLE.MAKE`** / **`PARTICLE3D.MAKE`** | No args → emitter **handle**. |
| **`PARTICLE.FREE`** / **`PARTICLE3D.FREE`** | Free emitter. |
| **`PARTICLE.SETTEXTURE`** | `(emitter, textureHandle)` |
| **`PARTICLE.SETEMITRATE`** / **`PARTICLE.SETRATE`** | `(emitter, per_sec)` — **`SETRATE`** is an alias. **`PARTICLE3D.SETRATE`** same. |
| **`PARTICLE.SETPOS`** | `(emitter, x#, y#, z#)` |
| **`PARTICLE.SETLIFETIME`** | `(emitter, min#, max#)` seconds |
| **`PARTICLE.SETVELOCITY`** | `(emitter, vx#, vy#, vz#, spread)` — **spread** is random component noise (same units as velocity components). |
| **`PARTICLE.SETDIRECTION`** | `(emitter, vx#, vy#, vz#)` — base direction; combine with **`SETSPREAD`**. |
| **`PARTICLE.SETSPREAD`** | `(emitter, angle)` — random jitter added to each velocity component at spawn. |
| **`PARTICLE.SETSPEED`** | `(emitter, min#, max#)` — per-particle scalar applied after direction+spread. Default `1..1`. |
| **`PARTICLE.SETSTARTSIZE`** | `(emitter, min#, max#)` — spawn size range. |
| **`PARTICLE.SETENDSIZE`** | `(emitter, min#, max#)` — end-of-life size range. |
| **`PARTICLE.SETSIZE`** | `(emitter, start#, end#)` — sets **both** start and end to **single** values (legacy shorthand). |
| **`PARTICLE.SETCOLOR`** / **`PARTICLE.SETSTARTCOLOR`** | `(emitter, r, g, b, a)` 0–255 |
| **`PARTICLE.SETCOLOREND`** / **`PARTICLE.SETENDCOLOR`** | End color |
| **`PARTICLE.SETGRAVITY`** | `(emitter, g)` **legacy** → `(0, g, 0)`, or `(emitter, gx#, gy#, gz#)` |
| **`PARTICLE.SETBURST`** | `(emitter, count)` — spawn **count** particles immediately (capped). |
| **`PARTICLE.SETBILLBOARD`** | `(emitter, TRUE/FALSE)` — **`TRUE`**: **`DrawBillboard`**. **`FALSE`**: draw **cubes** at particle positions (debug / non-camera-facing). |
| **`PARTICLE.PLAY`** / **`PARTICLE.STOP`** | Start/stop continuous emission. |
| **`PARTICLE.UPDATE`** | `(emitter, dt#)` |
| **`PARTICLE.DRAW`** | `(emitter)` uses **`CAMERA.BEGIN` … `CAMERA.END`** active camera, or **`(emitter, cameraHandle)`** for an explicit **`Camera3D`** handle. |
| **`PARTICLE.ISALIVE`** | `→ int` (`1` = still playing **or** live particles remain). |
| **`PARTICLE.COUNT`** | `→ int` live particles |

Every row above exists under **`PARTICLE3D.*`** as well (e.g. **`PARTICLE3D.SETTEXTURE`**, **`PARTICLE3D.DRAW`**, …).

---

## Handle methods

On a **`Particle`** handle, method calls map to **`PARTICLE.*`** keys (see **`vm/handlecall.go`**), e.g. **`emitter.SetPos`**, **`emitter.Play`**, **`emitter.SetStartColor`**.

---

## See also

- [PARTICLES.md](PARTICLES.md) — longer examples and workflow
- [CAMERA.md](CAMERA.md) — 3D pass for **`PARTICLE.DRAW`**
- [SPRITE.md](SPRITE.md) — **`PARTICLE2D.*`**
