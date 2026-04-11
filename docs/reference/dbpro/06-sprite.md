# DBPro — Sprite (2D)

moonBASIC **`Sprite.*`** is **handle-based** (load texture/atlas, then draw), not always the same as DBPro’s integer **id** + **Sprite(id, x, y, img)**. See [SPRITE.md](../SPRITE.md).

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **LOAD SPRITE** | ✓ **`Sprite.Load()`** | Returns **handle**. |
| **DELETE SPRITE** | ✓ **`Sprite.Free()`** | |
| **PASTE SPRITE** | ✓ **`Sprite.Draw()`** | |
| **POSITION SPRITE** | ✓ **`Sprite.SetPos()`** | |
| **MIRROR SPRITE** | ✓ **`Sprite.FlipH()`** | |
| **FLIP SPRITE** | ✓ **`Sprite.FlipV()`** | |
| **ROTATE SPRITE** | ✓ **`Sprite.Rotate()`** | |
| **SCALE SPRITE** | ✓ **`Sprite.Scale()`** | |
| **SET SPRITE FRAME** | ✓ **`Sprite.SetFrame()`** | |
| **SPRITE HIT** | ✓ **`Sprite.Hit()`** | |
| **SPRITE COLLISION** | ✓ **`Sprite.Hit()`** | |
| **SET SPRITE PRIORITY** | ≈ draw order / **Z** / **layer** | Engine uses **order** in your loop or **sprite batch** features if present. |
| **SET SPRITE IMAGE** | ≈ reload / swap texture handle | |
| **SPRITE HIT** | ✓ **`Sprite.Hit()`**, **`Sprite.PointHit()`** | |
