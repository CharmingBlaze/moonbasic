# Draw Commands

Commands for immediate-mode 2D and 3D drawing. These commands render shapes, text, textures, and 3D primitives directly to the screen without creating persistent objects. All draw calls must happen between `Render.Clear()` and `Render.Frame()`.

## Core Concepts

- **Immediate mode** — Draw commands render instantly and don't persist. You must redraw everything every frame.
- **2D drawing** — Screen-space coordinates. Origin (0,0) is top-left. X increases right, Y increases down.
- **3D drawing** — World-space coordinates. Must be inside a `Camera.Begin` / `Camera.End` or `Render.Begin3D` / `Render.End3D` block.
- **Color** — All colors are RGBA with values 0–255 per channel. Alpha 255 = fully opaque, 0 = invisible.

---

## 2D Text

### `Draw.Text(text, x, y, fontSize, r, g, b, a)`

Draws text at a screen position using the default font.

- `text` (string) — Text to draw.
- `x` (int) — X position in pixels.
- `y` (int) — Y position in pixels.
- `fontSize` (int) — Font size in pixels.
- `r`, `g`, `b`, `a` (int) — Color (0–255 each).

**How it works:** Uses Raylib's default bitmap font. For custom fonts, see the `FONT` namespace.

```basic
Draw.Text("Score: " + STR(score), 10, 10, 24, 255, 255, 255, 255)
Draw.Text("Game Over", 500, 300, 48, 255, 50, 50, 200)
```

---

## 2D Shapes

### `Draw.Rectangle(x, y, width, height, r, g, b, a)`

Draws a filled rectangle. **8 arguments.**

- `x`, `y` (int) — Top-left corner position.
- `width`, `height` (int) — Size in pixels.
- `r`, `g`, `b`, `a` (int) — Fill color.

**How it works:** Delegates to `rt.Driver.Video.DrawRectangle` which calls Raylib's `DrawRectangle`. Color components are clamped to `uint8`.

```basic
; Health bar background
Draw.Rectangle(10, 50, 200, 20, 60, 60, 60, 255)
; Health bar fill
Draw.Rectangle(10, 50, health * 2, 20, 50, 200, 50, 255)
```

---

### `Draw.RectLines(x, y, w, h, thickness, r, g, b, a)`

Draws a rectangle outline with configurable line thickness. **9 arguments.**

- `x`, `y` (int) — Top-left corner.
- `w`, `h` (int) — Width, height.
- `thickness` (float) — Line thickness in pixels.
- `r`, `g`, `b`, `a` (int) — Color.

```basic
; Selection highlight
Draw.RectLines(selX, selY, selW, selH, 2.0, 255, 255, 0, 255)
```

---

### `Draw.RectPro(x, y, w, h, originX, originY, rotation, r, g, b, a)`

Draws a filled rectangle with rotation around an origin point. **11 arguments.**

- `originX`, `originY` (float) — Rotation pivot.
- `rotation` (float) — Angle in degrees.

---

### `Draw.RectGradV(x, y, w, h, topR, topG, topB, topA, botR, botG, botB, botA)`

Draws a rectangle with a vertical gradient (top color to bottom color).

---

### `Draw.RectGradH(x, y, w, h, leftR, leftG, leftB, leftA, rightR, rightG, rightB, rightA)`

Draws a rectangle with a horizontal gradient.

---

### `Draw.RectGrad(x, y, w, h, topLeftR, topLeftG, topLeftB, topLeftA, ...)`

Draws a rectangle with four different corner colors (full gradient).

---

### `Draw.Circle(centerX, centerY, radius, r, g, b, a)`

Draws a filled circle. **7 arguments.**

- `centerX`, `centerY` (int) — Center position.
- `radius` (float) — Circle radius in pixels.
- `r`, `g`, `b`, `a` (int) — Fill color.

**How it works:** Delegates to `rt.Driver.Video.DrawCircle`.

```basic
Draw.Circle(640, 360, 50, 255, 200, 100, 255)
```

---

### `Draw.CircleLines(centerX, centerY, radius, r, g, b, a)`

Draws a circle outline. **7 arguments.** Same signature as `Draw.Circle`.

```basic
Draw.CircleLines(640, 360, 50, 255, 255, 255, 200)
```

---

### `Draw.CircleSector(cx, cy, radius, startAngle, endAngle, segments, r, g, b, a)`

Draws a filled sector (pie slice) of a circle. **10 arguments.**

- `cx`, `cy` (int) — Center.
- `radius` (float) — Radius.
- `startAngle`, `endAngle` (float) — Arc range in degrees.
- `segments` (int) — Number of segments (higher = smoother).
- `r`, `g`, `b`, `a` (int) — Color.

**How it works:** Calls Raylib's `DrawCircleSector`. Useful for pie charts, radial menus, and cooldown indicators.

```basic
; Cooldown indicator (quarter circle fills as cooldown progresses)
Draw.CircleSector(100, 100, 40, 0, cooldownPct * 360, 32, 50, 200, 50, 255)
```

---

### `Draw.CircleGradient(cx, cy, radius, innerR, innerG, innerB, innerA, outerR, outerG, outerB, outerA)`

Draws a circle with a radial gradient from inner color (center) to outer color (edge). **11 arguments.**

```basic
; Glowing aura effect
Draw.CircleGradient(playerX, playerY, 60, 255, 255, 100, 200, 255, 255, 100, 0)
```

---

### `Draw.Ellipse(centerX, centerY, radiusH, radiusV, r, g, b, a)`

Draws a filled ellipse. Also registered as `Draw.Oval`. **8 arguments.**

- `centerX`, `centerY` (int) — Center position.
- `radiusH` (float) — Horizontal radius.
- `radiusV` (float) — Vertical radius.
- `r`, `g`, `b`, `a` (int) — Fill color.

---

### `Draw.EllipseLines(centerX, centerY, radiusH, radiusV, r, g, b, a)`

Draws an ellipse outline. Also registered as `Draw.OvalLines`. **8 arguments.**

---

### `Draw.Ring(cx, cy, innerRadius, outerRadius, startAngle, endAngle, segments, r, g, b, a)`

Draws a filled ring (annulus) or arc. Useful for radial health bars.

---

### `Draw.RingLines(cx, cy, innerRadius, outerRadius, startAngle, endAngle, segments, r, g, b, a)`

Draws a ring outline.

---

### `Draw.Line(x1, y1, x2, y2, r, g, b, a)`

Draws a 1-pixel line between two points. **8 arguments.**

- `x1`, `y1` (int) — Start point.
- `x2`, `y2` (int) — End point.
- `r`, `g`, `b`, `a` (int) — Line color.

**How it works:** Calls Raylib's `DrawLine`.

```basic
Draw.Line(0, 0, 1280, 720, 255, 255, 255, 128)
```

---

### `Draw.LineEx(x1, y1, x2, y2, thickness, r, g, b, a)`

Draws a thick line. **9 arguments.**

- `x1`, `y1`, `x2`, `y2` (float) — Start and end (float precision).
- `thickness` (float) — Line width in pixels.
- `r`, `g`, `b`, `a` (int) — Color.

**How it works:** Uses `DrawLineEx` for anti-aliased thick lines.

```basic
Draw.LineEx(100, 100, 500, 300, 3.0, 255, 100, 100, 255)
```

---

### `Draw.LineBezier(x1, y1, x2, y2, thickness, r, g, b, a)`

Draws a cubic Bézier curve between two points. **9 arguments.**

- `x1`, `y1` (float) — Start point.
- `x2`, `y2` (float) — End point.
- `thickness` (float) — Line width.
- `r`, `g`, `b`, `a` (int) — Color.

```basic
; Smooth connection line (like node editor wires)
Draw.LineBezier(100, 200, 500, 400, 2.0, 100, 200, 255, 255)
```

---

### `Draw.LineBezierQuad(x1, y1, cx, cy, x2, y2, thickness, r, g, b, a)`

Draws a quadratic Bézier curve with one control point. **11 arguments.**

- `cx`, `cy` (float) — Control point.

---

### `Draw.LineBezierCubic(x1, y1, c1x, c1y, c2x, c2y, x2, y2, thickness, r, g, b, a)`

Draws a cubic Bézier curve with two control points. **13 arguments.**

---

### Spline Drawing

Splines draw smooth curves through arrays of points. All take `(pointsArray, thickness, r, g, b, a)` — **6 arguments**.

The `pointsArray` is a flat array of `[x1, y1, x2, y2, ...]` coordinate pairs.

| Command | Curve Type |
|---------|------------|
| `Draw.SplineLinear(pts, thick, r, g, b, a)` | Straight line segments |
| `Draw.SplineBasis(pts, thick, r, g, b, a)` | B-spline (smooth, doesn't pass through all points) |
| `Draw.SplineCatmullRom(pts, thick, r, g, b, a)` | Catmull-Rom (smooth, passes through all points) |
| `Draw.SplineBezierQuad(pts, thick, r, g, b, a)` | Quadratic Bézier segments |
| `Draw.SplineBezierCubic(pts, thick, r, g, b, a)` | Cubic Bézier segments |

**How it works:** Extracts x,y pairs from the array handle using `ArrayGetFloat`, builds a `[]rl.Vector2` slice (reused via `sync.Pool` to reduce GC pressure), and passes to the Raylib spline draw function.

```basic
; Draw a smooth path through waypoints
DIM pts(8)
pts(0) = 100 : pts(1) = 300
pts(2) = 300 : pts(3) = 100
pts(4) = 500 : pts(5) = 400
pts(6) = 700 : pts(7) = 200
Draw.SplineCatmullRom(pts, 2.0, 255, 200, 100, 255)
```

---

### `Draw.Triangle(x1, y1, x2, y2, x3, y3, r, g, b, a)`

Draws a filled triangle from three vertices. **10 arguments.**

- `x1`, `y1`, `x2`, `y2`, `x3`, `y3` (float) — Vertex positions.
- `r`, `g`, `b`, `a` (int) — Fill color.

**How it works:** Delegates to `rt.Driver.Video.DrawTriangle`. Vertices are wound counter-clockwise by Raylib.

---

### `Draw.TriangleLines(x1, y1, x2, y2, x3, y3, r, g, b, a)`

Draws a triangle outline. **10 arguments.**

---

### `Draw.Poly(cx, cy, sides, radius, rotation, r, g, b, a)`

Draws a filled regular polygon. **9 arguments.**

- `cx`, `cy` (float) — Center.
- `sides` (int) — Number of sides (3 = triangle, 6 = hexagon, etc.).
- `radius` (float) — Circumradius.
- `rotation` (float) — Rotation in degrees.
- `r`, `g`, `b`, `a` (int) — Color.

**How it works:** Calls Raylib's `DrawPoly`.

```basic
; Hexagonal tile
Draw.Poly(400, 300, 6, 50, 0, 100, 150, 200, 255)
```

---

### `Draw.PolyLines(cx, cy, sides, radius, rotation, r, g, b, a)`

Draws a polygon outline. Same args as `Draw.Poly`.

---

## Pixel-Level Drawing

### `Draw.Pixel(x, y, r, g, b, a)` / `Draw.Plot(x, y, r, g, b, a)`

Draws a single pixel. **6 arguments.** `Draw.Plot` is a Blitz-style alias.

- `x`, `y` (int) — Pixel position.

```basic
; Starfield
FOR i = 0 TO 199
    Draw.Pixel(starX(i), starY(i), 255, 255, 255, starBright(i))
NEXT
```

---

### `Draw.PixelV(x, y, r, g, b, a)`

Same as `Draw.Pixel` but accepts float coordinates (sub-pixel positioning). **6 arguments.**

---

### `Draw.Dot(x, y, size, r, g, b, a)`

Draws a filled dot (small circle). **7 arguments.** Useful for debug visualization.

- `x`, `y` (float) — Position.
- `size` (float) — Dot radius.

**How it works:** Internally calls `DrawCircleV`.

---

### `Draw.Arc(cx, cy, radius, startAngle, endAngle, thickness, r, g, b, a)`

Draws an arc (ring segment). **10 arguments.**

- `cx`, `cy` (float) — Center.
- `radius` (float) — Arc radius.
- `startAngle`, `endAngle` (float) — Arc range in degrees.
- `thickness` (float) — Ring thickness (clamped to radius).

**How it works:** Uses `DrawRing` with `innerRadius = radius - thickness`.

```basic
; Partial ring for loading indicator
Draw.Arc(640, 360, 40, 0, loadProgress * 360, 5, 100, 200, 255, 255)
```

---

### `Draw.GetPixelColor(x, y)`

Reads the color of a pixel on screen. **2 arguments.** Returns an array `[r, g, b, a]`.

**Returns:** `handle` — Array of 4 floats.

**How it works:** Takes a screenshot via `LoadImageFromScreen`, reads the pixel color, unloads the image. **Expensive** — do not call per-frame in hot loops.

---

### `Draw.Grid2D(spacing, r, g, b, a)`

Draws a full-screen 2D grid overlay. **5 arguments.** Useful for level editors and alignment.

- `spacing` (int) — Grid cell size in pixels.

**How it works:** Draws horizontal and vertical lines across the entire screen at `spacing` intervals.

```basic
; Level editor grid
Draw.Grid2D(32, 50, 50, 50, 100)
```

---

## Convenience Drawing

### `Draw.ProgressBar(x, y, w, h, value, max, r, g, b, a)`

Draws a progress/health bar with background and filled portion.

### `Draw.HealthBar(x, y, w, h, value, max, r, g, b, a)`

Alias for `Draw.ProgressBar`.

### `Draw.CenterText(text, y, fontSize, r, g, b, a)`

Draws text centered horizontally on screen.

### `Draw.RightText(text, x, y, fontSize, r, g, b, a)`

Draws text right-aligned to the given x position.

### `Draw.ShadowText(text, x, y, fontSize, r, g, b, a)`

Draws text with a dark drop shadow for readability on any background.

### `Draw.OutlineText(text, x, y, fontSize, r, g, b, a)`

Draws text with a dark outline.

### `Draw.Crosshair(x, y, size, r, g, b, a)`

Draws a crosshair (+) at the given position. Useful for FPS games.

### `Draw.RectGrid(x, y, cols, rows, cellW, cellH, r, g, b, a)`

Draws a rectangular grid of cells.

---

## 2D Text (Extended)

### `Draw.TextEx(fontHandle, text, x, y, fontSize, spacing, r, g, b, a)` / `Draw.TextFont(...)`

Draws text using a custom loaded font. Alias: `Draw.TextFont`.

- `fontHandle` (handle) — Font loaded with `Font.Load`.
- `spacing` (float) — Character spacing.

---

### `Draw.TextPro(fontHandle, text, posX, posY, originX, originY, rotation, fontSize, spacing, r, g, b, a)`

Draws text with rotation and origin control.

---

### `Draw.TextWidth(text, fontSize)`

Returns the pixel width of a text string at the given font size.

**Returns:** `int`

```basic
w = Draw.TextWidth("Hello", 24)
; Center text manually
Draw.Text("Hello", (screenW - w) / 2, 10, 24, 255, 255, 255, 255)
```

---

### `Draw.TextFontWidth(fontHandle, text, fontSize, spacing)`

Returns the pixel width of text using a custom font.

**Returns:** `int`

---

### `MeasureText(text, fontSize)` / `MeasureTextEx(fontHandle, text, fontSize, spacing)`

Easy Mode aliases for `Draw.TextWidth` / `Draw.TextFontWidth`.

---

### `GetTextCodepointCount(text)`

Returns the number of Unicode codepoints in a string.

**Returns:** `int`

---

### `Render.DrawFPS(x, y)`

Draws the current FPS counter at the given position using the default font. Registered in the draw module.

```basic
Render.DrawFPS(10, 10)
```

---

## 2D Texture Drawing (Immediate)

### `Draw.Texture(textureHandle, x, y, r, g, b, a)`

Draws a texture at a screen position with a tint color.

---

### `Draw.TextureV(textureHandle, posX, posY, r, g, b, a)`

Same as `Draw.Texture` but with float position.

---

### `Draw.TextureEx(textureHandle, posX, posY, rotation, scale, r, g, b, a)`

Draws a texture with rotation and scale.

---

### `Draw.TextureRec(textureHandle, srcX, srcY, srcW, srcH, dstX, dstY, r, g, b, a)`

Draws a sub-rectangle of a texture. Essential for **spritesheets/texture atlases**.

- `srcX`, `srcY`, `srcW`, `srcH` (float) — Source rectangle in the texture.
- `dstX`, `dstY` (float) — Screen destination.

```basic
; Draw frame 3 from a 32x32 spritesheet (4 columns)
frameX = (3 MOD 4) * 32
frameY = (3 / 4) * 32
Draw.TextureRec(sheet, frameX, frameY, 32, 32, 100, 100, 255, 255, 255, 255)
```

---

### `Draw.TexturePro(textureHandle, srcX, srcY, srcW, srcH, dstX, dstY, dstW, dstH, originX, originY, rotation, r, g, b, a)`

The most flexible texture draw — source rect, destination rect, origin, rotation, and tint. Used for scaled/rotated sprite rendering.

---

### `Draw.TextureFull(textureHandle, r, g, b, a)`

Draws a texture stretched to fill the entire screen.

---

### `Draw.TextureFlipped(textureHandle, x, y, r, g, b, a)`

Draws a texture flipped vertically.

---

### `Draw.TextureTiled(textureHandle, srcX, srcY, srcW, srcH, dstX, dstY, dstW, dstH, originX, originY, rotation, scale, r, g, b, a)`

Draws a texture tiled/repeated within a destination rectangle.

---

### `Draw.TextureNPatch(textureHandle, ...)`

Draws a 9-patch (nine-slice) texture for scalable UI panels, borders, and buttons.

---

## 2D Prim Objects

For more control over 2D shapes, create persistent 2D primitive handles with modifiable properties.

### `DrawPrim2D.Create(kind)` — Creates a 2D primitive handle.

Available kinds: `Circle`, `CircleLines`, `Ellipse`, `EllipseLines`, `Rect`, `RectLines`, `Line`, `Triangle`, `TriangleLines`, `Ring`, `RingLines`, `Poly`, `PolyLines`.

**Returns:** `handle`

Methods on prim objects: `.pos(x, y)`, `.size(w, h)`, `.color(r, g, b, a)`, `.draw()`, `.free()`.

```basic
; Create a reusable selection box
sel = DrawPrim2D.Create("RectLines")
sel.pos(100, 100)
sel.size(200, 150)
sel.color(255, 255, 0, 200)

; Each frame
sel.draw()

; Cleanup
sel.free()
```

---

## 2D Text Objects

### `TextObj(text, x, y, fontSize, r, g, b, a)`

Creates a persistent text draw object. Faster than `Draw.Text` when the same text is drawn every frame — avoids repeated string marshalling.

**Returns:** `handle`

Methods: `TextDraw.Pos(h, x, y)`, `TextDraw.Size(h, fontSize)`, `TextDraw.Color(h, r, g, b, a)`, `TextDraw.SetText(h, newText)`, `TextDraw.Draw(h)`, `TextDraw.Free(h)`.

```basic
score = TextObj("Score: 0", 10, 10, 24, 255, 255, 255, 255)

; In game loop
TextDraw.SetText(score, "Score: " + STR(points))
TextDraw.Draw(score)

; Cleanup
TextDraw.Free(score)
```

---

## 2D Texture Draw Objects

### `DrawTex2(textureHandle, x, y, r, g, b, a)`

Creates a persistent texture draw object for efficient repeated rendering.

**Returns:** `handle`

Methods: `DrawTex2.Pos(h, x, y)`, `DrawTex2.Color(h, r, g, b, a)`, `DrawTex2.SetTexture(h, texHandle)`, `DrawTex2.Draw(h)`, `DrawTex2.Free(h)`.

---

### `DrawTexRec(textureHandle, srcX, srcY, srcW, srcH, dstX, dstY, r, g, b, a)`

Creates a persistent texture-rectangle draw object (spritesheet slice).

**Returns:** `handle`

Methods: `DrawTexRec.Src(h, sx, sy, sw, sh)`, `DrawTexRec.Pos(h, x, y)`, `DrawTexRec.Color(h, r, g, b, a)` / `.Col(...)`, `DrawTexRec.SetTexture(h, tex)`, `DrawTexRec.Draw(h)`, `DrawTexRec.Free(h)`.

---

### `DrawTexPro(textureHandle, srcX, srcY, srcW, srcH, dstX, dstY, dstW, dstH, originX, originY, rotation, r, g, b, a)`

Creates a persistent pro-texture draw object with full transform control.

**Returns:** `handle`

Methods: `DrawTexPro.Src(h, ...)`, `DrawTexPro.Dst(h, ...)`, `DrawTexPro.Origin(h, x, y)`, `DrawTexPro.Rot(h, angle)`, `DrawTexPro.Color(h, r, g, b, a)` / `.Col(...)`, `DrawTexPro.SetTexture(h, tex)`, `DrawTexPro.Draw(h)`, `DrawTexPro.Free(h)`.

---

## 3D Primitives

All 3D draw commands must be inside a `Camera.Begin` / `Camera.End` block.

### `Draw.Grid(slices, spacing)`

Draws a ground grid for reference. Very useful during development.

- `slices` (int) — Number of grid divisions.
- `spacing` (float) — Distance between grid lines in world units.

```basic
Camera.Begin(cam)
    Draw.Grid(20, 1.0)
Camera.End(cam)
```

---

### `Draw.Cube(x, y, z, width, height, depth, r, g, b, a)`

Draws a solid colored cube at a world position.

- `x`, `y`, `z` (float) — Center position.
- `width`, `height`, `depth` (float) — Size.
- `r`, `g`, `b`, `a` (int) — Color.

```basic
Draw.Cube(0, 1, 0, 2, 2, 2, 200, 50, 50, 255)
```

---

### `Draw.CubeWires(x, y, z, width, height, depth, r, g, b, a)`

Draws a wireframe cube.

---

### `Draw.Sphere(x, y, z, radius, r, g, b, a)`

Draws a solid colored sphere.

- `x`, `y`, `z` (float) — Center position.
- `radius` (float) — Sphere radius.
- `r`, `g`, `b`, `a` (int) — Color.

```basic
Draw.Sphere(5, 1, 0, 1.0, 50, 200, 50, 255)
```

---

### `Draw.SphereWires(x, y, z, radius, r, g, b, a)`

Draws a wireframe sphere.

---

### `Draw.Cylinder(x, y, z, radiusTop, radiusBottom, height, slices, r, g, b, a)`

Draws a solid cylinder or cone.

- `radiusTop`, `radiusBottom` (float) — Top and bottom radii. Set one to 0 for a cone.
- `height` (float) — Cylinder height.
- `slices` (int) — Number of side segments.

---

### `Draw.CylinderWires(x, y, z, radiusTop, radiusBottom, height, slices, r, g, b, a)`

Draws a wireframe cylinder.

---

### `Draw.Capsule(startX, startY, startZ, endX, endY, endZ, radius, slices, rings, r, g, b, a)`

Draws a solid capsule (cylinder with hemispherical caps).

---

### `Draw.CapsuleWires(startX, startY, startZ, endX, endY, endZ, radius, slices, rings, r, g, b, a)`

Draws a wireframe capsule.

---

### `Draw.Plane(x, y, z, width, depth, r, g, b, a)`

Draws a flat plane (quad) at a world position.

---

### `Draw.Line3D(x1, y1, z1, x2, y2, z2, r, g, b, a)`

Draws a line in 3D space.

- `x1`, `y1`, `z1` (float) — Start point.
- `x2`, `y2`, `z2` (float) — End point.
- `r`, `g`, `b`, `a` (int) — Color.

```basic
; Draw an axis indicator
Draw.Line3D(0, 0, 0, 5, 0, 0, 255, 0, 0, 255)   ; X = red
Draw.Line3D(0, 0, 0, 0, 5, 0, 0, 255, 0, 255)   ; Y = green
Draw.Line3D(0, 0, 0, 0, 0, 5, 0, 0, 255, 255)   ; Z = blue
```

---

### `Draw.Point3D(x, y, z, r, g, b, a)`

Draws a point in 3D space.

---

### `Draw.BoundingBox(minX, minY, minZ, maxX, maxY, maxZ, r, g, b, a)`

Draws a wireframe axis-aligned bounding box.

---

### `Draw.Ray(originX, originY, originZ, dirX, dirY, dirZ, r, g, b, a)`

Draws a ray as a line from origin in the given direction.

---

### `Draw.Billboard(cameraHandle, textureHandle, x, y, z, size, r, g, b, a)`

Draws a textured quad that always faces the camera (billboard). Useful for particles, labels, and vegetation.

- `cameraHandle` (handle) — Camera for orientation.
- `textureHandle` (handle) — Texture to draw.
- `x`, `y`, `z` (float) — World position.
- `size` (float) — Billboard size.
- `r`, `g`, `b`, `a` (int) — Tint color.

---

### `Draw.BillboardRec(cameraHandle, textureHandle, sourceRect, x, y, z, sizeX, sizeY, r, g, b, a)`

Draws a billboard using a sub-rectangle of a texture (texture atlas support).

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `LINE3D(x1,y1,z1,x2,y2,z2,r,g,b,a)` | `Draw.Line3D(...)` |
| `DRAW3D.GRID(s, sp)` | `Draw.Grid(s, sp)` |
| `DRAW3D.CUBE(...)` | `Draw.Cube(...)` |
| `DRAW3D.SPHERE(...)` | `Draw.Sphere(...)` |

---

## Full Example

A complete demo showing 2D UI overlays and 3D primitive rendering.

```basic
Window.Open(1280, 720, "Draw Demo")
Window.SetFPS(60)

cam = Camera.Create()
cam.pos(0, 8, 15)
cam.look(0, 0, 0)
cam.fov(60)

angle = 0
health = 75

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()
    angle = angle + 45 * dt

    Render.Clear(25, 25, 40)

    ; === 3D Pass ===
    Camera.Begin(cam)
        Draw.Grid(20, 1.0)

        ; Rotating cube
        x = SIN(angle * 0.0174) * 4
        z = COS(angle * 0.0174) * 4
        Draw.Cube(x, 1, z, 1.5, 1.5, 1.5, 220, 80, 80, 255)

        ; Static sphere
        Draw.Sphere(-3, 1, 0, 1.0, 80, 220, 80, 255)

        ; Wireframe cylinder
        Draw.CylinderWires(3, 0, 0, 0.5, 1.0, 3, 16, 80, 80, 220, 255)

        ; Axis lines
        Draw.Line3D(0, 0, 0, 5, 0, 0, 255, 0, 0, 255)
        Draw.Line3D(0, 0, 0, 0, 5, 0, 0, 255, 0, 255)
        Draw.Line3D(0, 0, 0, 0, 0, 5, 0, 0, 255, 255)
    Camera.End(cam)

    ; === 2D Overlay ===
    Draw.Text("Draw Demo", 10, 10, 24, 255, 255, 255, 255)
    Draw.Text("FPS: " + STR(Window.GetFPS()), 10, 40, 18, 100, 255, 100, 255)

    ; Health bar
    Draw.Rectangle(10, 70, 200, 16, 60, 60, 60, 255)
    Draw.Rectangle(10, 70, health * 2, 16, 50, 200, 50, 255)
    Draw.Text("HP: " + STR(health), 220, 70, 16, 200, 200, 200, 255)

    ; Crosshair
    cx = Window.Width() / 2
    cy = Window.Height() / 2
    Draw.Line(cx - 10, cy, cx + 10, cy, 255, 255, 255, 200)
    Draw.Line(cx, cy - 10, cx, cy + 10, 255, 255, 255, 200)

    Render.Frame()
WEND

Camera.Free(cam)
Window.Close()
```

---

## See Also

- [RENDER](RENDER.md) — Frame lifecycle, blend modes, wireframe toggle
- [CAMERA](CAMERA.md) — Required for 3D drawing passes
- [TEXTURE](TEXTURE.md) — Textures for billboard and sprite drawing
- [FONT](FONT.md) — Custom fonts for text rendering
