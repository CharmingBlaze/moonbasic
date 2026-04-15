# 3D Hop (`mario64`)

Small **third-person** demos in a **Blitz3D-style** spirit: walk on a plane, jump, **orbit the camera**, physics-backed primitives. Several variants:

| File | What to notice |
|------|----------------|
| **`modern_blitz_hop.mb`** | **Minimal loop:** **`cam.Orbit(player, distance)`** (engine-owned yaw/pitch/zoom), **`cam.Yaw()`** for facing, **`player.Move`** in units/sec. **Start here:** **[`MODERN_BLITZ_HOP_BEGINNER.md`](MODERN_BLITZ_HOP_BEGINNER.md)** (line-by-line tutorial). Optional **orbit configuration** (see below). |
| **`modern_blitz_hop_kcc.mb`** | Same orbit camera, but the hero uses **Jolt KCC** (**`CHAR.CREATE`**, **`CHAR.MOVEWITHCAMERA`**, **`CHAR.JUMP`**) instead of **`ENTITY.PHYSICS`** on the player. **Linux + CGO + fullruntime.** See **[`docs/reference/KCC.md`](../../docs/reference/KCC.md)**. |
| **`main_orbit_simple.mb`** | **Easiest read:** commented “map” at the top, **`CONST`** palette + world bounds, one floor + one box — **`ORBITYAWDELTA` / `ORBITPITCHDELTA` / `ORBITDISTDELTA`**, **`MOVESTEPX`/`Z`**, **`LANDBOXES`**, **`Camera.SetOrbit`**, **`ERASE ALL`**. |
| **`main.mb`** | **Default pick:** same hop as before, **implicit types** (no `#` / `$` / `?` suffixes), **Draw3D** + **Camera** only — no skybox or entity graph. |
| **`main_entities.mb`** | **Engine-style:** **CreateCube** / **CreateSphere**, **COLLISIONS**, **EntityGrounded** (coyote), **EntityMoveCameraRelative**, **Camera.OrbitEntity**, **CopyEntity** platforms, **ENTITY.UPDATE**, **DrawEntities**, child **hat** on **player**. |
| **`main_v2.mb`** | **Recommended teaching path:** parallel arrays for platforms, but **`Input.Axis`**, **`MOVEX`/`MOVEZ`**, **`BOXTOPLAND`** float return, **`IIF$`**, and **one line** for orbit yaw (`Input.Axis(KEY_Q, KEY_E) * DEGPERSEC(...)`). Heavily commented. |
| **`main_v3.mb`** | Same logic with **`TYPE` / `DIM AS`** — one `Platform` array instead of nine arrays. Uses **`Input.Orbit`** and **`MOVESTEPX`/`MOVESTEPZ`**. |

---

## `modern_blitz_hop.mb` — orbit-follow API

**New to this sample?** Read **[`MODERN_BLITZ_HOP_BEGINNER.md`](MODERN_BLITZ_HOP_BEGINNER.md)** for prerequisites, how to run, a frame-by-frame mental model, and a **line-by-line** walkthrough of the source.

The engine keeps **yaw**, **pitch**, and **orbit distance** inside the camera. Each frame you call **`cam.Orbit(player, 12.0)`** and read **`cam.Yaw()`** to rotate the player so **WASD** matches the view.

**Default controls:** **WASD** move · **right-drag** orbit · **Q/E** yaw orbit · **wheel** zoom · **Space** jump · **Esc** quit.

Full sample (as in the repo):

```moonbasic
; ==========================================
; moonBASIC 64 — Clean orbit hop (engine-owned yaw/pitch/dist)
; ==========================================
; cam.Orbit(entity, distance) — R-drag + Q/E + wheel; cam.Yaw() for player facing.
; Move(forward,right,up) is units per second (dt applied internally on scripted builds).
;
;   go run . --check examples/mario64/modern_blitz_hop.mb
;   CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/mario64/modern_blitz_hop.mb

Window.Open(1280, 720, "moonBASIC 64")
Window.SetFPS(60)

PHYSICS3D.START()
WORLD.Gravity(0, -40, 0)

cam = CreateCamera()
cam.SetFOV(60)

player = Model.CreateCapsule(0.4, 1.0)
player.Pos(0, 5, 0)
player.Color(255, 60, 60)
ENTITY.PHYSICS(player, "CAPSULE", 1.0, 0.9, 0.0)

floor = Model.CreateBox(100, 2, 100)
floor.Pos(0, -1, 0)
floor.Color(60, 200, 90)
ENTITY.PHYSICS(floor, "BOX", 0.0, 0.9, 0.0)

WHILE NOT (KEYDOWN(KEY_ESCAPE) OR Window.ShouldClose())

    cam.Orbit(player, 12.0)

    player.SetRot(0, cam.Yaw(), 0)

    fwd = Input.Axis(KEY_S, KEY_W)
    side = Input.Axis(KEY_A, KEY_D)
    player.Move(10.0 * fwd, 10.0 * side, 0)

    IF KEYPRESSED(KEY_SPACE) AND player.IsGrounded() THEN
        player.Jump(12.0)
        player.Squash(0.5, 0.3)
    ENDIF

    UPDATEPHYSICS()

    RENDER.Clear(100, 150, 250)

    cam.Begin()
        ENTITY.DRAWALL()
    cam.End()

    DRAW.TEXT("WASD move · R-drag orbit · Q/E yaw · wheel zoom · Space jump · ESC", 20, 20, 14, 255, 255, 255, 255)

    RENDER.FRAME()
WEND

Window.Close()
```

### How the orbit loop fits together

1. **`CreateCamera()`** (Easy Mode; forwards to **`CAMERA.CREATE`**) creates the camera. Optionally call **orbit configuration** commands right after (see below).
2. **Each frame** in your **`WHILE`**: call **`cam.Orbit(player, distance)`** so the engine updates hidden yaw/pitch/zoom and moves the eye around the player.
3. **`cam.Yaw()`** reads the **horizontal** angle of that orbit (radians). Use **`player.SetRot(0, cam.Yaw(), 0)`** so the character faces the way the camera looks, and **`player.Move(...)`** stays aligned with **WASD**.

You only configure orbit **once** at startup. The game loop stays one line: **`cam.Orbit(player, 12.0)`**.

---

### Orbit configuration (optional) — step by step

**When to call:** after **`cam = CreateCamera()`** (and **`cam.SetFOV`**, etc.), **before** the **`WHILE`** loop.

**Two ways to spell the same call:**

- **Dot style (short):** **`cam.UseMouseOrbit(FALSE)`**
- **Full name:** **`Camera.UseMouseOrbit(cam, FALSE)`**

The tables below use **dot style** on **`cam`**.

---

#### Defaults (you can skip everything)

If you never call the settings below, you get:

| Setting | Default |
|--------|---------|
| Mouse moves orbit | **On** |
| Mouse only while **right button** held | **Yes** (hold RMB to drag-view) |
| Keys to yaw the orbit | **Q** (left), **E** (right) |
| Mouse sensitivity | **0.005** |
| Wheel zoom scale | **1.0** |
| Keyboard yaw speed | **1.5** radians per second |
| Pitch limits | **−1.5** to **1.5** radians |
| Distance (zoom) limits | **2** to **50** world units |

---

#### `cam.UseMouseOrbit(useMouse)` — turn mouse orbit on or off

| | |
|--|--|
| **Arguments** | **`useMouse`** — **`TRUE`** / **`FALSE`** (or **`1`** / **`0`**) |
| **Use when** | You need the mouse for something else (crosshair, UI) and only want **keys + wheel** to move the camera. |

**Example — keyboard + wheel only (mouse does not orbit):**

```moonbasic
cam = CreateCamera()
cam.UseMouseOrbit(FALSE)
; Q/E and mouse wheel still work (unless you also change keys / disable them)
```

---

#### `cam.UseOrbitRightMouse(requireRightMouse)` — RMB drag vs “always” mouse

| | |
|--|--|
| **Arguments** | **`requireRightMouse`** — **`TRUE`** = orbit with mouse **only while right button is down** (default). **`FALSE`** = moving the mouse orbits **without** holding RMB (closer to a free-fly inspector). |

**Example — orbit by moving the mouse without holding RMB:**

```moonbasic
cam.UseOrbitRightMouse(FALSE)
```

**Example — go back to “MMO style” (only drag while RMB held):**

```moonbasic
cam.UseOrbitRightMouse(TRUE)
```

---

#### `cam.SetOrbitKeys(leftKey, rightKey)` — which keys spin the orbit left/right

| | |
|--|--|
| **Arguments** | Raylib key codes, e.g. **`KEY_Q`**, **`KEY_E`**, **`KEY_LEFT`**, **`KEY_RIGHT`**. Use **`0`** for one side to turn off that direction. Use **`0, 0`** to disable **keyboard** orbit entirely (mouse/wheel unchanged). |

**Example — arrow keys instead of Q/E:**

```moonbasic
cam.SetOrbitKeys(KEY_LEFT, KEY_RIGHT)
```

**Example — no keyboard orbit (mouse + wheel only):**

```moonbasic
cam.SetOrbitKeys(0, 0)
```

---

#### `cam.SetOrbitLimits(minPitch, maxPitch, minDist, maxDist)` — stop flips and extreme zoom

| | |
|--|--|
| **Arguments** | **`minPitch`**, **`maxPitch`** — radians (down/up tilt). **`minDist`**, **`maxDist`** — how close/far the camera can be in **world units**. |

**Rough pitch guide:** **0** ≈ level; negative looks down; positive looks up. **`~±1.2`** rad is a bit tighter than the default **`±1.5`**.

**Example — shallower tilt + closer zoom range (good for brawlers / isometric-ish feel):**

```moonbasic
cam.SetOrbitLimits(-1.0, 1.0, 4.0, 25.0)
```

**Example — almost top-down (narrow pitch band):**

```moonbasic
cam.SetOrbitLimits(0.5, 1.1, 10.0, 40.0)
```

---

#### `cam.SetOrbitSpeed(mouseSens, wheelSens)` — how fast mouse drag and wheel feel

| | |
|--|--|
| **Arguments** | **`mouseSens`** — multiplier on mouse movement (yaw + pitch while orbiting). **`wheelSens`** — how strong **scroll zoom** is. Larger = faster. |

**Example — snappier mouse, stronger zoom:**

```moonbasic
cam.SetOrbitSpeed(0.008, 2.5)
```

**Example — gentler control:**

```moonbasic
cam.SetOrbitSpeed(0.003, 0.7)
```

---

#### `cam.SetOrbitKeySpeed(keyRadPerSec)` — how fast Q/E (or your keys) spin the orbit

| | |
|--|--|
| **Arguments** | **`keyRadPerSec`** — radians **per second** (not per frame). Higher = faster keyboard orbit. |

**Example — slower, precise keyboard orbit:**

```moonbasic
cam.SetOrbitKeySpeed(0.8)
```

**Example — snappier keyboard orbit:**

```moonbasic
cam.SetOrbitKeySpeed(2.5)
```

---

### Recipe: combine settings for common goals

**A — Shooter-style:** mouse aims (no orbit from mouse), orbit only with keys:

```moonbasic
cam.UseMouseOrbit(FALSE)
cam.SetOrbitKeys(KEY_Q, KEY_E)
cam.SetOrbitKeySpeed(1.5)
```

**B — “Always look” third person:** mouse orbits without holding RMB:

```moonbasic
cam.UseOrbitRightMouse(FALSE)
```

**C — Full keyboard + wheel, no mouse camera at all:**

```moonbasic
cam.UseMouseOrbit(FALSE)
cam.SetOrbitKeys(KEY_LEFT, KEY_RIGHT)
cam.SetOrbitLimits(-1.2, 1.2, 5.0, 30.0)
cam.SetOrbitSpeed(0.008, 2.0)
cam.SetOrbitKeySpeed(2.0)
```

**D — Defaults with tighter zoom/pitch only:**

```moonbasic
cam.SetOrbitLimits(-1.2, 1.2, 5.0, 30.0)
```

Full registry reference: **[CAMERA.md](../../docs/reference/CAMERA.md)**.

---

## Run

**`main.mb`** (Draw3D path, full runtime + Raylib):

```bash
go run -tags fullruntime . --run examples/mario64/main.mb
```

**`modern_blitz_hop.mb`** (entity + physics + orbit-follow):

```bash
go run -tags fullruntime . --run examples/mario64/modern_blitz_hop.mb
```

**Blitz-style variants** (often need **CGO** and the same **fullruntime** build if you use **`--run`** from the repo root):

```bash
go run -tags fullruntime . --run examples/mario64/main_orbit_simple.mb
```

**Controls (`main.mb` / orbit samples):** **Q/E** yaw, **right-drag** yaw/pitch, **wheel** zoom, **WASD** move, **Space** jump, **Esc** quit.

## Docs to read

- **[BLITZ3D.md](../../docs/reference/BLITZ3D.md)** — BlitzBasic3D → moonBASIC map (**`KEYHIT`**, **`WIRECUBE`**, **`Camera.Orbit`**, entities, …).  
- **[CAMERA.md](../../docs/reference/CAMERA.md)** — **`Camera.SetOrbit`**, **entity `Orbit`**, **`Camera.Yaw`**, orbit **configuration** builtins, **`Camera.OrbitAround`**, **`GetRay`**.  
- **`main_orbit_simple.mb`** — orbit deltas **`ORBITYAWDELTA` / `ORBITPITCHDELTA` / `ORBITDISTDELTA`** in [GAMEHELPERS.md](../../docs/reference/GAMEHELPERS.md); teardown **`ERASE ALL`** in [MEMORY.md](../../docs/MEMORY.md).  
- **Orbit camera** — `Camera.OrbitAround` in [CAMERA.md](../../docs/reference/CAMERA.md) (third-person on XZ + fixed eye height).  
- **Walk + orbit input** — `Input.Axis` in [INPUT.md](../../docs/reference/INPUT.md); pair **Q/E** with **`DEGPERSEC`** for degrees-per-second yaw.  
- **Movement** — `MOVEX` / `MOVEZ` in [MATH.md](../../docs/reference/MATH.md).  
- **Landing** — `BOXTOPLAND` / `LANDBOXES` in [GAMEHELPERS.md](../../docs/reference/GAMEHELPERS.md).
