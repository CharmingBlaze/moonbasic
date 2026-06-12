# Randomness and procedural rolls — dice, loot, and seeds

> Fair rolls, loot tables, spawn scatter, and reproducible runs with `MATH.RAND`, `RNDF`, `CHANCE`, and seeds.

**Namespaces:** `MATH` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#data](../../COMMAND_REGISTRY.md#data) · **Vectors for scatter:** [VEC2-MATH.md](VEC2-MATH.md) · [MATH-3D-GAMEPLAY.md](MATH-3D-GAMEPLAY.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use which random API](#when-to-use-which-random-api)
- [Core workflow](#core-workflow)
- [Command guide](#command-guide)
- [Loot table pattern](#loot-table-pattern)
- [Procedural spawn scatter](#procedural-spawn-scatter)
- [Full example — dice and loot](#full-example--dice-and-loot)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | Int rand, float rand, probability checks, optional seed |
| **Typical games** | Damage variance, crit chance, loot, star fields |
| **Reproducible** | `MATH.RNDSEED` for daily seed / debug replay |

---

## When to use which random API

| Need | Use |
|------|-----|
| Integer inclusive range | `MATH.RAND(lo, hi)` or `MATH.RND(lo, hi)` |
| Float range | `MATH.RNDF(lo, hi)` or `MATH.RANGE` |
| 0..n-1 | `MATH.RND(n)` |
| “30% crit?” | `MATH.CHANCE(0.3)` |
| Same run tomorrow | `MATH.RNDSEED(dayNumber)` at boot |
| Pick array index | `MATH.RAND(1, count)` with your array |

---

## Core workflow

1. **Optional seed** at level load — `MATH.RNDSEED(seed)` for reproducible tests.
2. **Roll** when event fires (not every frame unless noise).
3. **Clamp** result to gameplay bounds — `MATH.CLAMP`.
4. **Weighted loot** — roll bucket then item (see below).

---

## Command guide

```basic
dice = MATH.RAND(1, 6)
offset = MATH.RNDF(-1, 1)
IF MATH.CHANCE(0.25) THEN crit = 1

; Debug reproducibility
MATH.RNDSEED(12345)
a = MATH.RAND(1, 100)
```

| Command | Why |
|---------|-----|
| `MATH.RAND(min, max)` | Int inclusive |
| `MATH.RND` | Variants — see registry |
| `MATH.RNDF(min, max)` | Float uniform |
| `MATH.CHANCE(p)` | Bool with probability p |
| `MATH.RNDSEED(n)` | Reset RNG sequence |

---

## Loot table pattern

```basic
roll = MATH.RAND(1, 100)
IF roll <= 50 THEN
    item = "coin"
ELSE IF roll <= 80 THEN
    item = "gem"
ELSE
    item = "rare_sword"
ENDIF
```

**Why tiers:** One roll selects bucket — easier than floating weights for beginners.

---

## Procedural spawn scatter

```basic
; Ring spawn on XZ
i = MATH.RAND(1, 8)
pt = MATH.CIRCLEPOINT(cx, cz, radius, i, 8)
; or random angle:
ang = MATH.RNDF(0, MATH.TAU())
sx = cx + MATH.COS(ang) * radius
sz = cz + MATH.SIN(ang) * radius
```

2D screen scatter: `MATH.RAND(0, APP.WIDTH())`, `MATH.RAND(0, APP.HEIGHT())`.

---

## Full example — dice and loot

```basic
APP.OPEN(400, 280, "Random")
APP.SETFPS(60)
MATH.RNDSEED(42)

lastRoll = 0
lastItem = "none"

WHILE NOT APP.SHOULDCLOSE()
    IF INPUT.KEYHIT(KEY_SPACE) THEN
        lastRoll = MATH.RAND(1, 6)
        r = MATH.RAND(1, 100)
        IF r <= 60 THEN lastItem = "herb"
        ELSE IF r <= 90 THEN lastItem = "ore"
        ELSE lastItem = "relic"
    ENDIF

    crit = 0
    IF MATH.CHANCE(0.1) THEN crit = 1

    DRAW.TEXT("Roll " + lastRoll + "  Loot " + lastItem, 10, 10, 16, 255, 255, 255)
    IF crit THEN DRAW.TEXT("CRIT!", 10, 30, 16, 255, 80, 80)
    RENDER.FRAME()
WEND
APP.CLOSE()
```

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| `RAND` every frame for one event | Roll on `KEYHIT` / timer |
| Wrong inclusive range | Verify min ≤ max in registry |
| Seed after rolls | Seed **before** sequence you want fixed |
| `CHANCE(30)` instead of `0.3` | Chance is 0..1 float |
| Predictable if unseeded | OK for arcade; seed for roguelike |

---

## See also

- [MATH-2D-GAMEPLAY.md](MATH-2D-GAMEPLAY.md) — `CIRCLEPOINT` placement
- [TERRAIN-OPEN-WORLD.md](../TERRAIN-OPEN-WORLD.md) — scatter on heightfield
- [SAVE-AND-PROGRESS.md](../SAVE-AND-PROGRESS.md) — store seed in save
