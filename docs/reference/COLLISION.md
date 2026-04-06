# 2D/3D collision math (`BOXCOLLIDE`, …)

Pure collision tests registered at **top-level** keys (no `GAME.` prefix) from **`runtime/mbgame/register_collision.go`**:

- **`BOXCOLLIDE`**, **`CIRCLECOLLIDE`**, **`POINTINBOX`**, **`POINTINCIRCLE`**, **`CIRCLEBOXCOLLIDE`**, **`LINECOLLIDE`**, **`POINTONLINE`**, **`SPHERECOLLIDE`**, **`AABBCOLLIDE`**, **`SPHEREBOXCOLLIDE`**, **`POINTINAABB`**
- **Sphere vs AABB top landing:** **`BOXTOPLAND`** returns a **float** landing centre Y (or **`0.0`**). **`LANDBOXES`** scans the same test over **parallel box arrays** — see **[GAMEHELPERS.md](GAMEHELPERS.md)**.
- Distance helpers: **`DISTANCE2D`**, **`DISTANCE3D`**, **`DISTANCESQ2D`**, **`DISTANCESQ3D`**
- **Raycasts:** **`RAY.*`** (3D, Raylib / CGO) and **`RAY2D.*`** (circle, rect, segment) — see **[RAYCAST.md](RAYCAST.md)**.

**Sprites:** axis-aligned tests are **`Sprite.Hit`**, alias **`SPRITECOLLIDE`**, and **`Sprite.PointHit`** — see **[SPRITE.md](SPRITE.md)**.
