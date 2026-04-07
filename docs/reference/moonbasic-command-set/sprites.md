# Sprites (2D)

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **CreateSprite (img)** | **`SPRITE.LOAD`** (texture/path) | **Sprite handle** — **`SPRITE.FREE`**. |
| **Sprite (id, x, y)** | **`SPRITE.SETPOS`** + **`SPRITE.DRAW`** in loop | |
| **MoveSprite** | **`SPRITE.SETPOS`** | |
| **SpriteImage** | Reload / swap underlying texture | New texture may need **`TEXTURE.FREE`**. |
| **SpriteColor / Alpha** | Tint in draw path or sprite state | |
| **SpriteHit** | **`SPRITE.HIT`**, **`POINTHIT`** | |
