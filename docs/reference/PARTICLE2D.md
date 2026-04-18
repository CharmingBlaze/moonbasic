# Particle2D Commands

2D screen-space particle emitters for effects like sparks, smoke trails, and explosions in 2D games.

## Core Workflow

1. `PARTICLE2D.CREATE(r, g, b, a, size)` — create an emitter with a base color and particle size.
2. Each frame: `PARTICLE2D.EMIT(emitter, x, y, vx, vy, lifetime)` — spawn particles at a position with velocity.
3. `PARTICLE2D.UPDATE(emitter, dt)` — advance particle lifetimes.
4. `PARTICLE2D.DRAW(emitter)` — render all live particles.
5. `PARTICLE2D.FREE(emitter)` when done.

---

## Creation

### `PARTICLE2D.CREATE(r, g, b, a, size)` 

Creates a 2D particle emitter. `r, g, b, a` set the base particle color (0–255). `size` is the particle radius in pixels. Returns an **emitter handle**.

---

## Emission

### `PARTICLE2D.EMIT(emitter, x, y, vx, vy, lifetime)` 

Spawns a particle at world position `(x, y)` with velocity `(vx, vy)` pixels/second. `lifetime` is how long the particle lives in seconds before fading out.

---

## Update & Draw

### `PARTICLE2D.UPDATE(emitter, dt)` 

Advances all live particles by `dt` seconds. Removes expired particles.

---

### `PARTICLE2D.DRAW(emitter)` 

Draws all live particles. Call between `CAMERA2D.BEGIN` / `CAMERA2D.END` if using a scrolling camera.

---

## Lifetime

### `PARTICLE2D.FREE(emitter)` 

Destroys the emitter and all its particles.

---

## Full Example

Sparks emitting from the mouse cursor when the left button is held.

```basic
WINDOW.OPEN(800, 600, "Particle2D Demo")
WINDOW.SETFPS(60)

sparks = PARTICLE2D.CREATE(255, 200, 60, 255, 3)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    IF MOUSE.DOWN(0) THEN
        FOR i = 1 TO 5
            angle = RND(0, 360) * 3.14159 / 180.0
            speed = RNDF(80, 200)
            vx    = COS(angle) * speed
            vy    = SIN(angle) * speed
            PARTICLE2D.EMIT(sparks, MOUSE.X(), MOUSE.Y(), vx, vy, RNDF(0.3, 0.8))
        NEXT i
    END IF

    PARTICLE2D.UPDATE(sparks, dt)

    RENDER.CLEAR(10, 10, 20)
    PARTICLE2D.DRAW(sparks)
    RENDER.FRAME()
WEND

PARTICLE2D.FREE(sparks)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `PARTICLE2D.MAKE(...)` | Deprecated alias of `PARTICLE2D.CREATE`. |

---

## See also

- [PARTICLE3D.md](PARTICLE3D.md) — 3D world-space particles
- [PARTICLE.md](PARTICLE.md) — legacy particle emitters
- [DRAW2D.md](DRAW2D.md) — primitive 2D drawing
