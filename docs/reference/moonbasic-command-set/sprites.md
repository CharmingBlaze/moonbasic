# Sprites (2D)

| Designed | moonBASIC | Notes |
|----------|------------|-------|
| **LoadSprite(file)** | **`Sprite.Load()`** | Returns a **sprite handle**. |
| **DrawSprite(id, x, y)** | **`Sprite.Draw()`** | Renders at pixel coordinates. |
| **MoveSprite(id, x, y)** | **`Sprite.SetPos()`** | Sets float draw offset. |
| **SpriteHit(a, b)** | **`Sprite.Hit()`** | Overlap of the **drawn** quads (**`DrawTexturePro`**: scale, origin, rotation), not a plain axis-aligned frame rect. |
| **PointHit(id, x, y)** | **`Sprite.PointHit()`** | Point vs same quad as **`Sprite.Draw`** + **`SetPos`**. |
| **ScaleSprite(id, s)** | **`Sprite.Scale()`** | |
| **RotateSprite(id, a)** | **`Sprite.Rotate()`** | |
