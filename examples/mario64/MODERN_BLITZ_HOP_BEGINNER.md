# Beginner’s guide: `modern_blitz_hop.mb`

This document explains **[`modern_blitz_hop.mb`](modern_blitz_hop.mb)** from the ground up. You do **not** need to know 3D math first: the engine stores camera **yaw**, **pitch**, and **zoom distance** for you. Your job is to create objects, call **`cam.Orbit`**, move the player, and draw.

---

## What you will learn

- How a minimal **moonBASIC** 3D program is structured (window → physics → objects → loop → draw).
- What **`cam.Orbit(player, distance)`** does and why **`cam.Yaw()`** pairs with **`player.SetRot`**.
- How **WASD** movement stays aligned with the camera.
- How **jump**, **physics**, and the **main loop** fit together.

---

## What this program does

You control a **red capsule** on a **green floor**. The **camera** orbits around the player: hold **right mouse** and drag to look around, **Q** / **E** spin the view, the **mouse wheel** zooms. **WASD** walks relative to where the camera looks. **Space** jumps. **Esc** exits.

---

## What you need installed

- **Go** (to build and run moonBASIC).
- A **C compiler** and **CGO** enabled (**`CGO_ENABLED=1`**) when using the full 3D runtime (Raylib + physics). On Windows, a typical setup uses **MSYS2** or **Visual Studio** build tools so CGO can link Raylib.

The repo’s **[`CONTRIBUTING.md`](../../CONTRIBUTING.md)** and **[`docs/DEVELOPER.md`](../../docs/DEVELOPER.md)** describe editor setup and build tags.

---

## How to run it

From the **repository root**:

**Check that the script compiles** (bytecode only, no game window):

```bash
go run . --check examples/mario64/modern_blitz_hop.mb
```

**Run the game** (needs **`fullruntime`** so the window, camera, and physics work):

```bash
go run -tags fullruntime . --run examples/mario64/modern_blitz_hop.mb
```

If your project uses a separate runner:

```bash
go run -tags fullruntime ./cmd/moonrun examples/mario64/modern_blitz_hop.mb
```

If something fails, confirm **CGO** is on and you built with **`fullruntime`**. See **[`AGENTS.md`](../../AGENTS.md)** in the repo root for a short note on **`go run`** vs graphical programs.

---

## Big picture: one frame of the game

Each time through the **`WHILE`** loop, the program roughly does this **in order**:

1. **Update the camera** so it sits on a sphere around the player (**`cam.Orbit`**).
2. **Rotate the player** to face the same horizontal direction as the camera (**`cam.Yaw`** → **`SetRot`**).
3. **Read input** and **move** the player (**`Input.Axis`** → **`Move`**).
4. **Jump** if requested and allowed (**`KEYPRESSED`**, **`IsGrounded`**).
5. **Step the physics world** (**`UPDATEPHYSICS`**).
6. **Clear the screen**, **draw the 3D scene** (**`cam.Begin`** … **`ENTITY.DRAWALL`** … **`cam.End`**), draw **HUD text**, **present** the frame (**`RENDER.FRAME`**).

Order matters: the camera uses the player’s **position** from last frame’s physics; then movement and **`UPDATEPHYSICS`** advance the world for **next** frame.

---

## Concepts (short glossary)

| Idea | Plain meaning |
|------|----------------|
| **Handle** | A number that refers to an object the engine created (camera, entity). You pass it into commands instead of raw pointers. |
| **`CreateCamera()`** | Creates a 3D camera (Easy Mode → **`CAMERA.CREATE`**). **`cam`** stores that handle. |
| **Orbit** | The camera sits at some **distance** from the player and looks at them. **Yaw** spins left/right; **pitch** looks up/down. |
| **`cam.Orbit(player, 12.0)`** | “Each frame, recompute my orbit around **`player`** using base distance **12** (world units), and apply mouse / keys / wheel according to defaults or your [orbit settings](README.md#orbit-configuration-optional--step-by-step).” |
| **`cam.Yaw()`** | Horizontal angle of that orbit, in **radians**, so you can rotate the character to **face** the view. |
| **`player.Move(...)`** | Moves the entity; arguments are **units per second** on the scripted path (the runtime applies **delta time** where needed). |
| **`UPDATEPHYSICS`** | Advances Jolt (or the bundled physics step) so collisions and gravity apply. |
| **`ENTITY.DRAWALL()`** | Draws every entity registered with the entity system. |

---

## Walkthrough: the source file

Below, line numbers match the current **[`modern_blitz_hop.mb`](modern_blitz_hop.mb)**. Comments starting with **`;`** are ignored by the compiler.

### Lines 1–8 — header comments

These lines document **what** the sample is and **how** to compile or run it. They are for **you**, not the engine. Keep similar notes in your own projects so “future you” remembers the command line.

---

### Lines 10–11 — window and frame rate

```moonbasic
Window.Open(1280, 720, "moonBASIC 64")
Window.SetFPS(60)
```

- **`Window.Open`** creates the OS window: width, height, title.
- **`Window.SetFPS(60)`** asks for about 60 updates per second (smooth motion).

---

### Lines 15–16 — physics world

```moonbasic
PHYSICS3D.START()
WORLD.Gravity(0, -40, 0)
```

- **`PHYSICS3D.START`** turns on the 3D physics engine used by **`ENTITY.PHYSICS`** / **`AddPhysics`** (see [PHYSICS_ERGONOMICS.md](../../docs/PHYSICS_ERGONOMICS.md)).
- **`WORLD.Gravity(x, y, z)`** sets gravity. **−40** on **Y** means “down” is the negative Y direction, so things fall toward the floor.

---

### Lines 18–19 — camera

```moonbasic
cam = CreateCamera()
cam.SetFOV(60)
```

- **`CreateCamera()`** returns a **camera handle** stored in **`cam`**.
- **`SetFOV(60)`** sets vertical field of view to **60°** (wider = see more at once; narrower = more “zoomed in” feel).

---

### Lines 21–26 — player (dynamic capsule)

```moonbasic
player = Model.CreateCapsule(0.4, 1.0)
player.Pos(0, 5, 0)
player.Color(255, 60, 60)
ENTITY.PHYSICS(player, "CAPSULE", 1.0, 0.9, 0.0)
```

- **`CreateCapsule(radius, height)`** builds a **capsule** entity: the mesh is a rounded capsule (not a plain cylinder), aligned with **Jolt** when you add physics.
- **`Pos`** places it in world space: here **5** units up so it starts above the floor.
- **`Color`** is RGB **0–255** (red-ish capsule).
- **`ENTITY.PHYSICS(player, "CAPSULE", 1.0, 0.9, 0.0)`** — entity-first setup: **mass 1** = dynamic, **friction 0.9**, **restitution 0** (no bounce). Equivalent to older **`AddPhysics("dynamic", "capsule")`** plus **`SetBounciness`**, but **one call**; full options in [PHYSICS_ERGONOMICS.md](../../docs/PHYSICS_ERGONOMICS.md).

---

### Lines 28–32 — floor (static box)

```moonbasic
floor = Model.CreateBox(100, 2, 100)
floor.Pos(0, -1, 0)
floor.Color(60, 200, 90)
ENTITY.PHYSICS(floor, "BOX", 0.0, 0.9, 0.0)
```

- A large thin **box** acts as the ground. **Mass `0.0`** = **static**; friction/restitution match the player so the floor is not a “trampoline.”
- **`Pos(0, -1, 0)`** centers it so the top surface is near **Y = 0** (with the given height **2**, the center at **−1** puts the top at **0**).

---

### Lines 33–35 — prime orbit (before the loop)

```moonbasic
cam.Orbit(player, 12.0)
player.SetRot(0, cam.Yaw(), 0)
```

Places the camera **once** before the first frame is drawn so you do not start with the default camera pose. Inside the loop, **`UPDATEPHYSICS`** runs **before** **`cam.Orbit`** so the camera follows the same-frame player position.

---

### Lines 37–62 — main loop

```moonbasic
WHILE NOT (KEYDOWN(KEY_ESCAPE) OR Window.ShouldClose())
```

Repeat until the user presses **Escape** or closes the window.

---

#### Lines 39–41 — camera-relative walk

```moonbasic
    fwd = Input.Axis(KEY_S, KEY_W)
    side = Input.Axis(KEY_A, KEY_D)
    ENTITY.MOVEWITHCAMERA(player, cam, fwd, side, 10.0)
```

- **`Input.Axis(neg, pos)`** returns roughly **−1..1** depending on keys (**S** vs **W**, **A** vs **D**).
- **`ENTITY.MOVEWITHCAMERA`** sets horizontal **walk velocity** in the camera’s ground plane (orbit yaw), so **W** matches what you see after **right-drag** orbit. **10.0** is speed in world units per second.

---

#### Lines 43–46 — jump

```moonbasic
    IF KEYPRESSED(KEY_SPACE) AND player.IsGrounded() THEN
        player.Jump(12.0)
        player.Squash(0.5, 0.3)
    ENDIF
```

- Jump only if **Space** was pressed **this frame** and the player is **grounded**.
- **`Jump`** / **`Squash`** — see engine docs for units.

---

#### Lines 48–51 — physics, then orbit + facing

```moonbasic
    UPDATEPHYSICS()
    cam.Orbit(player, 12.0)
    player.SetRot(0, cam.Yaw(), 0)
```

- **`UPDATEPHYSICS`** advances **ENTITY.UPDATE** and **PHYSICS3D.STEP** (Jolt on Linux+CGO).
- **`cam.Orbit`** after physics keeps the third-person rig from lagging a frame.
- **`SetRot(0, cam.Yaw(), 0)`** makes the capsule face orbit **yaw** (radians).

---

#### Lines 53–61 — draw 3D, HUD, present

```moonbasic
    RENDER.Clear(100, 150, 250)

    cam.Begin()
        ENTITY.DRAWALL()
    cam.End()

    DRAW.TEXT("WASD = camera-relative · ...", 20, 20, 14, 255, 255, 255, 255)

    RENDER.FRAME()
```

- **`RENDER.Clear`** — background RGB.
- **`cam.Begin` / `cam.End`** — 3D pass with this camera.
- **`ENTITY.DRAWALL`** — player, floor, etc.
- **`DRAW.TEXT`** / **`RENDER.FRAME`** — HUD and swap buffers.

**Orbit controls:** right-drag, **Q/E**, wheel (see **[README](README.md)**).

---

### Lines 64–65 — shutdown

```moonbasic
Window.Close()
```

Releases the window when the loop exits.

---

## Customizing the camera (optional)

You do **not** need to change anything for the sample to work. If you want **keyboard-only** orbit, **no mouse orbit** (for aiming), different keys, or zoom limits, see the step-by-step tables and recipes in **[`README.md` — Orbit configuration](README.md#orbit-configuration-optional--step-by-step)** and the reference in **[`docs/reference/CAMERA.md`](../../docs/reference/CAMERA.md)**.

---

## Where to go next

| Goal | Document / sample |
|------|---------------------|
| More orbit recipes (shooter, free mouse, etc.) | **[`README.md`](README.md)** |
| Full **`CAMERA.*`** reference | **[`docs/reference/CAMERA.md`](../../docs/reference/CAMERA.md)** |
| Manual yaw/pitch with **`Camera.SetOrbit`** | **`main_orbit_simple.mb`**, **[`GAMEHELPERS.md`](../../docs/reference/GAMEHELPERS.md)** |
| Blitz-style naming map | **[`docs/reference/BLITZ3D.md`](../../docs/reference/BLITZ3D.md)** |
| Input axes | **[`docs/reference/INPUT.md`](../../docs/reference/INPUT.md)** |

Welcome to moonBASIC 3D — start by changing **speed**, **jump**, **colors**, and **orbit distance**, then peek at **`README.md`** when you want tighter control over the camera.
