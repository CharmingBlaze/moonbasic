# Sprite Commands

Commands for creating and controlling 2D sprites. Sprites are textured quads with position, rotation, scale, color, animation frames, and layer-based draw ordering. They are the primary building block for 2D games.

## Core Concepts

- **Sprite** — A 2D textured rectangle with position, size, rotation, color, alpha, and optional animation frames.
- **Texture source** — Each sprite is backed by a texture handle. Multiple sprites can share the same texture.
- **Animation** — Sprites support frame-based animation from a spritesheet grid.
- **Layers** — Sprites render in layer order (lower layers draw first = further back).
- All sprites are **heap handles** and must be freed.

---

## Creation

### `Sprite.Create(textureHandle)` / `Sprite.Make(textureHandle)`

Creates a new sprite from a texture.

- `textureHandle` (handle) — Texture loaded with `Texture.Load`.

**Returns:** `handle`

```basic
playerTex = Texture.Load("assets/player.png")
player = Sprite.Create(playerTex)
```

---

### `Sprite.Free(spriteHandle)`

Frees a sprite handle.

---

## Position & Transform

### `Sprite.SetPos(spriteHandle, x, y)` / `sprite.pos(x, y)`

Sets the sprite's screen position.

- `x`, `y` (float) — Position in pixels.

```basic
Sprite.SetPos(player, 100, 200)
```

---

### `Sprite.SetRot(spriteHandle, angle)` / `sprite.rot(angle)`

Sets the sprite's rotation in degrees.

- `angle` (float) — Rotation angle.

---

### `Sprite.SetScale(spriteHandle, sx, sy)` / `sprite.scale(sx, sy)`

Sets the sprite's scale.

- `sx`, `sy` (float) — Scale factors (1.0 = original size).

---

### `Sprite.SetOrigin(spriteHandle, ox, oy)`

Sets the sprite's rotation/scale origin point (pivot).

- `ox`, `oy` (float) — Origin in pixels relative to sprite top-left.

```basic
; Center pivot
Sprite.SetOrigin(player, 16, 16)
```

---

## Appearance

### `Sprite.SetColor(spriteHandle, r, g, b)` / `sprite.col(r, g, b)`

Sets the sprite tint color.

---

### `Sprite.SetAlpha(spriteHandle, alpha)` / `sprite.alpha(a)`

Sets sprite transparency (0.0–1.0).

---

### `Sprite.SetFlip(spriteHandle, flipX, flipY)`

Flips the sprite horizontally and/or vertically.

- `flipX`, `flipY` (bool) — Flip flags.

```basic
; Flip when facing left
IF facingLeft THEN
    Sprite.SetFlip(player, TRUE, FALSE)
ELSE
    Sprite.SetFlip(player, FALSE, FALSE)
ENDIF
```

---

### `Sprite.SetLayer(spriteHandle, layer)`

Sets the sprite's draw layer. Lower layers render behind higher layers.

- `layer` (int) — Layer number.

```basic
Sprite.SetLayer(background, 0)
Sprite.SetLayer(player, 5)
Sprite.SetLayer(foreground, 10)
```

---

### `Sprite.Show(spriteHandle)` / `Sprite.Hide(spriteHandle)`

Toggle sprite visibility.

---

## Animation

### `Sprite.SetFrameSize(spriteHandle, frameW, frameH)`

Defines the frame dimensions for spritesheet animation. The texture is divided into a grid of frames.

- `frameW`, `frameH` (int) — Frame size in pixels.

```basic
; 32x32 pixel frames in a spritesheet
Sprite.SetFrameSize(player, 32, 32)
```

---

### `Sprite.SetFrame(spriteHandle, frameIndex)` / `sprite.frame(index)`

Sets the current animation frame (0-based).

- `frameIndex` (int) — Frame index.

```basic
; Animate by cycling frames
animTimer = animTimer + dt * 10
Sprite.SetFrame(player, INT(animTimer) MOD 4)
```

---

### `Sprite.SetAnim(spriteHandle, startFrame, endFrame, fps, loop)`

Sets up automatic frame animation.

- `startFrame` (int) — First frame.
- `endFrame` (int) — Last frame.
- `fps` (float) — Animation speed in frames per second.
- `loop` (bool) — `TRUE` to loop.

```basic
Sprite.SetAnim(player, 0, 3, 10, TRUE)   ; Walk cycle, 4 frames at 10fps
```

---

## Rendering

### `Sprite.Draw(spriteHandle)`

Draws a single sprite.

---

### `Sprite.DrawAll()`

Draws all visible sprites in layer order.

```basic
Render.Clear(40, 40, 60)
Sprite.DrawAll()
Render.Frame()
```

---

## Collision

### `Sprite.GetRect(spriteHandle)`

Returns the sprite's bounding rectangle for collision detection.

**Returns:** x, y, width, height

---

### `Sprite.Overlaps(spriteA, spriteB)`

Returns `TRUE` if two sprites overlap (AABB collision).

**Returns:** `bool`

```basic
IF Sprite.Overlaps(player, enemy) THEN
    PRINT "Hit!"
ENDIF
```

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `CreateSprite(tex)` | `Sprite.Create(tex)` |
| `PositionSprite(s, x, y)` | `Sprite.SetPos(s, x, y)` |
| `RotateSprite(s, a)` | `Sprite.SetRot(s, a)` |
| `ScaleSprite(s, sx, sy)` | `Sprite.SetScale(s, sx, sy)` |
| `FreeSprite(s)` | `Sprite.Free(s)` |

---

## Full Example

A 2D game with animated sprites and layer-based rendering.

```basic
Window.Open(800, 600, "Sprite Demo")
Window.SetFPS(60)

; Load textures
bgTex = Texture.Load("assets/background.png")
playerTex = Texture.Load("assets/player_sheet.png")

; Background sprite
bg = Sprite.Create(bgTex)
Sprite.SetPos(bg, 0, 0)
Sprite.SetLayer(bg, 0)

; Player sprite with animation
player = Sprite.Create(playerTex)
Sprite.SetPos(player, 400, 300)
Sprite.SetFrameSize(player, 32, 32)
Sprite.SetOrigin(player, 16, 16)
Sprite.SetLayer(player, 5)
Sprite.SetAnim(player, 0, 3, 10, TRUE)

playerX = 400
playerY = 300
speed = 200

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()

    ; Movement
    IF Input.KeyDown(KEY_A) THEN
        playerX = playerX - speed * dt
        Sprite.SetFlip(player, TRUE, FALSE)
    ENDIF
    IF Input.KeyDown(KEY_D) THEN
        playerX = playerX + speed * dt
        Sprite.SetFlip(player, FALSE, FALSE)
    ENDIF
    IF Input.KeyDown(KEY_W) THEN playerY = playerY - speed * dt
    IF Input.KeyDown(KEY_S) THEN playerY = playerY + speed * dt

    Sprite.SetPos(player, playerX, playerY)

    ; Render
    Render.Clear(0, 0, 0)
    Sprite.DrawAll()

    Draw.Text("WASD = Move", 10, 10, 16, 255, 255, 255, 255)
    Render.Frame()
WEND

Sprite.Free(player)
Sprite.Free(bg)
Texture.Free(playerTex)
Texture.Free(bgTex)
Window.Close()
```

---

## See Also

- [TEXTURE](TEXTURE.md) — Loading textures for sprites
- [CAMERA](CAMERA.md) — Camera2D for scrolling 2D worlds
- [INPUT](INPUT.md) — Input for sprite movement
- [DRAW](DRAW.md) — Immediate-mode 2D drawing alternative
