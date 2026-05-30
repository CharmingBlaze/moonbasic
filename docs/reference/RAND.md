# Rand Commands

Seeded pseudo-random number generator handles. Use when you need reproducible random sequences independent of global `RND` / `RNDF`.

For global random functions see [MATH.md](MATH.md) (`RND`, `RNDF`, `RNDSEED`).

## Core Workflow

1. `RAND.CREATE(seed)` — create a seeded RNG handle.
2. `RAND.NEXT(rng, min, max)` — get a random integer in `[min, max]`.
3. `RAND.NEXTF(rng)` — get a random float in `[0.0, 1.0)`.
4. `RAND.FREE(rng)` when done.

---

## Creation

### `RAND.CREATE(seed)` 

Creates a seeded pseudo-random generator. Same seed always produces the same sequence. Returns a **rng handle**.

---

## Generation

### `RAND.NEXT(rng, min, max)` 

Returns a random integer in the closed range `[min, max]`.

---

### `RAND.NEXTF(rng)` 

Returns a random float in `[0.0, 1.0)`.

---

## Lifetime

### `RAND.FREE(rng)` 

Frees the RNG handle.

---

## Full Example

Procedural dungeon room layout with a fixed seed for repeatability.

```basic
WINDOW.OPEN(800, 600, "Rand Demo")
WINDOW.SETFPS(60)

SEED = 42
rng = RAND.CREATE(SEED)

; generate 20 room positions reproducibly
rooms = ARRAY.MAKE(20)
FOR i = 0 TO 19
    rx = RAND.NEXT(rng, 0, 24) * 30
    ry = RAND.NEXT(rng, 0, 18) * 30
    ARRAY.SET(rooms, i, rx * 10000 + ry)   ; pack x,y
NEXT i

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(10, 10, 20)
    FOR i = 0 TO 19
        packed = ARRAY.GET(rooms, i)
        rx = INT(packed / 10000)
        ry = packed MOD 10000
        DRAW.RECT(rx, ry, 24, 18, 60, 100, 160, 255)
    NEXT i
    DRAW.TEXT("Seed: " + STR(SEED), 10, 10, 18, 200, 200, 200, 255)
    RENDER.FRAME()
WEND

RAND.FREE(rng)
WINDOW.CLOSE()
```

---

## See also

- [MATH.md](MATH.md) — global `RND`, `RNDF`, `RNDSEED`, `RNDRANGE`
- [NOISE.md](NOISE.md) — Perlin/Simplex procedural noise
- [SCATTER.md](SCATTER.md) — random placement of props
