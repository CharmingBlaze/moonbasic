# Wind Commands

Set a global wind vector that affects cloth, rope, particles, and foliage shaders.

## Core Workflow

1. `WIND.SET(x, y, z)` — set the global wind direction and strength vector.
2. `WIND.GETSTRENGTH()` — read back the wind magnitude.
3. Systems that react to wind (cloth, particle3D with gravity enabled, foliage) pick this up automatically.

---

## Commands

### `WIND.SET(x, y, z)` 

Sets the global wind vector. The vector encodes both direction and speed — e.g. `WIND.SET(3, 0, 0)` = 3 m/s blowing in the +X direction.

---

### `WIND.GETSTRENGTH()` 

Returns the magnitude of the current wind vector as a float.

---

## Full Example

Gusty wind that oscillates with a sine wave.

```basic
WINDOW.OPEN(960, 540, "Wind Demo")
WINDOW.SETFPS(60)

t = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    t  = t + dt

    ; oscillating gust
    strength = SIN(t * 0.8) * 4.0
    WIND.SET(strength, 0, 0)

    RENDER.CLEAR(80, 120, 160)
    DRAW.TEXT("Wind: " + STR(WIND.GETSTRENGTH()), 10, 10, 20, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

WIND.SET(0, 0, 0)
WINDOW.CLOSE()
```

---

## See also

- [CLOTH_ROPE_LIGHTING.md](CLOTH_ROPE_LIGHTING.md) — cloth simulation (reacts to wind)
- [PARTICLE3D.md](PARTICLE3D.md) — particles with gravity/wind
- [WEATHER.md](WEATHER.md) — rain, snow, storm effects
