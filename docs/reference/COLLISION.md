# 2D/3D collision math (`BOXCOLLIDE`, …)

Pure collision tests registered at **top-level** keys (no `GAME.` prefix) from **`runtime/mbgame/register_collision.go`**:

- **`BOXCOLLIDE`**, **`CIRCLECOLLIDE`**, **`POINTINBOX`**, **`POINTINCIRCLE`**, **`CIRCLEBOXCOLLIDE`**, **`LINECOLLIDE`**, **`POINTONLINE`**, **`SPHERECOLLIDE`**, **`AABBCOLLIDE`**, **`SPHEREBOXCOLLIDE`**, **`POINTINAABB`**
- Distance helpers: **`DISTANCE2D`**, **`DISTANCE3D`**, **`DISTANCESQ2D`**, **`DISTANCESQ3D`**

**Sprites:** axis-aligned tests are **`Sprite.Hit`**, alias **`SPRITECOLLIDE`**, and **`Sprite.PointHit`** — see **[SPRITE.md](SPRITE.md)**.
