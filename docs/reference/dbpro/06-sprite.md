# DBPro — Sprite (2D)

moonBASIC **`SPRITE.*`** is **handle-based** (load texture/atlas, then draw), not always the same as DBPro’s integer **id** + **SPRITE (id, x, y, img)**. See [SPRITE.md](../SPRITE.md).

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **SPRITE (id, x, y, img)** | ≈ **`SPRITE.LOAD`** + **`SPRITE.SETPOS`** + **`SPRITE.DRAW`** | Typical frame: update position, draw. |
| **DELETE SPRITE** | ✓ **`SPRITE.FREE`** | |
| **MOVE SPRITE** | ✓ **`SPRITE.SETPOS`** / **`SETPOSITION`** | |
| **SIZE SPRITE** | ≈ scale via draw / texture | Check **`SPRITE.*`** for size helpers in manifest. |
| **SET SPRITE ALPHA** / **COLOR** | ≈ **`SPRITE`** + tint in draw pipeline | Often **`DRAW`** + texture; see examples. |
| **SET SPRITE PRIORITY** | ≈ draw order / **Z** / **layer** | Engine uses **order** in your loop or **sprite batch** features if present. |
| **SET SPRITE IMAGE** | ≈ reload / swap texture handle | |
| **SPRITE HIT** | ✓ **`SPRITE.HIT`**, **`SPRITE.POINTHIT`** | |
