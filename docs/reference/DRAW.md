# Draw Commands

Immediate-mode 2D and 3D primitive drawing. All calls must be inside the appropriate render scope (`RENDER.CLEAR` → draw → `RENDER.FRAME`; 3D calls inside `RENDER.BEGIN3D` / `RENDER.END3D`).

## Core Workflow

```basic
RENDER.CLEAR(r, g, b)
; 2D draws here (no camera transform needed)
DRAW.RECTANGLE(x, y, w, h, r, g, b, a)
DRAW.TEXT("hello", x, y, size, r, g, b, a)
; 3D draws between BEGIN3D/END3D
RENDER.BEGIN3D(cam)
    DRAW.CUBE(x, y, z, w, h, d, rot, r, g, b, a)
    DRAW.SPHERE(x, y, z, radius, r, g, b, a)
    DRAW.LINE3D(x1, y1, z1, x2, y2, z2, r, g, b, a)
RENDER.END3D()
RENDER.FRAME()
```

---

## 2D Shapes

### `DRAW.RECTANGLE(x, y, w, h, r, g, b, a)` 

Filled rectangle at `(x, y)` with width `w` and height `h`.

---

### `DRAW.RECTANGLE_ROUNDED(x, y, w, h, roundness, segments, r, g, b, a)` 

Filled rectangle with rounded corners. `roundness` is 0.0–1.0 (proportion of the shortest side used as corner radius). `segments` controls smoothness (e.g. 8).

---

### `DRAW.RECTLINES(x, y, w, h, lineWidth, thick, r, g, b, a)` 

Outlined rectangle with `thick` border.

---

### `DRAW.RECTPRO(x, y, w, h, originX, originY, rot, r, g, b, a)` 

Rectangle with rotation `rot` (degrees) around `(originX, originY)`.

---

### `DRAW.RECTGRAD(x, y, w, h, r1, g1, b1, a1, r2, g2, b2, a2, ...)` 

Rectangle with corner color gradients.

---

### `DRAW.RECTGRADH(x, y, w, h, r1, g1, b1, a1, r2, g2, b2, a2)` 

Horizontal gradient rectangle (left to right).

---

### `DRAW.RECTGRADV(x, y, w, h, r1, g1, b1, a1, r2, g2, b2, a2)` 

Vertical gradient rectangle (top to bottom).

---

### `DRAW.RECTGRID(x, y, w, h, cellsX, cellsY)` 

Draws a grid of lines inside a rectangle.

---

### `DRAW.CIRCLE(x, y, radius, r, g, b, a)` 

Filled circle.

---

### `DRAW.CIRCLELINES(x, y, radius, r, g, b, a)` 

Circle outline.

---

### `DRAW.CIRCLESECTOR(x, y, radius, startAngle, endAngle, segs, r, g, b, a)` 

Filled circle sector (pie slice).

---

### `DRAW.CIRCLEGRADIENT(x, y, radius, r1, g1, b1, a1, r2, g2, b2, a2)` 

Circle with radial gradient from centre to edge.

---

### `DRAW.ELLIPSE(x, y, rx, ry, r, g, b, a)` / `DRAW.OVAL(...)` 

Filled ellipse with radii `rx` and `ry`.

---

### `DRAW.ELLIPSELINES(x, y, rx, ry, r, g, b, a)` / `DRAW.OVALLINES(...)` 

Ellipse outline.

---

### `DRAW.RING(cx, cy, innerR, outerR, startAngle, endAngle, segs, r, g, b, a)` 

Filled ring / donut sector.

---

### `DRAW.RINGLINES(...)` 

Ring outline.

---

### `DRAW.TRIANGLE(x1, y1, x2, y2, x3, y3, r, g, b, a)` 

Filled triangle.

---

### `DRAW.TRIANGLELINES(...)` 

Triangle outline.

---

### `DRAW.POLY(cx, cy, sides, radius, rot, r, g, b, a)` 

Filled regular polygon.

---

### `DRAW.POLYLINES(cx, cy, sides, radius, rot, thick, r, g, b, a)` 

Polygon outline with thickness.

---

### `DRAW.LINE(x1, y1, x2, y2, r, g, b, a)` 

Simple line between two points.

---

### `DRAW.LINEEX(x1, y1, x2, y2, thick, r, g, b, a)` 

Thick line.

---

### `DRAW.LINEBEZIER(x1, y1, x2, y2, thick, r, g, b, a)` 

Bezier-curved line segment.

---

### `DRAW.LINEBEZIERCUBIC(x1, y1, c1x, c1y, c2x, c2y, x2, y2, thick, r, g, b, a)` 

Cubic bezier line.

---

### `DRAW.LINEBEZIERQUAD(x1, y1, cx, cy, x2, y2, thick, r, g, b, a)` 

Quadratic bezier line.

---

### `DRAW.ARC(cx, cy, radius, startAngle, endAngle, thick, r, g, b, a)` 

Arc outline.

---

### `DRAW.PIXEL(x, y, r, g, b, a)` / `DRAW.PLOT(...)` / `DRAW.DOT(...)` / `DRAW.PIXELV(...)` 

Single pixel / dot draw.

---

### `DRAW.CROSSHAIR(x, y, size, gap, thick, r, g, b, a)` 

Draws a crosshair at `(x, y)`.

---

### `DRAW.PROGRESSBAR(x, y, w, h, border, progress, r1, g1, b1, a1, r2, g2, b2, a2)` 

Progress bar with background and fill colors.

---

### `DRAW.HEALTHBAR(x, y, w, h, border, progress, ...)` 

Health bar (same as progressbar with extra color params).

---

### `DRAW.GRID2D(cells, offsetX, offsetY, r, g, b, a)` 

2D grid overlay.

---

## 2D Text

### `DRAW.TEXT(text, x, y, size, r, g, b, a)` 

Draws text using the default font.

---

### `DRAW.TEXTEX(font, text, x, y, size, spacing, r, g, b, a)` / `DRAW.TEXTFONT(...)` 

Draws text with a custom `font` handle (from `FONT.LOAD`).

---

### `DRAW.TEXTPRO(font, text, x, y, originX, originY, rot, size, spacing, r, g, b, a)` 

Text with rotation and custom origin.

---

### `DRAW.CENTERTEXT(text, y, size, r, g, b, a)` 

Draws text horizontally centred on screen.

---

### `DRAW.RIGHTTEXT(text, x, y, size, r, g, b, a)` 

Right-aligned text.

---

### `DRAW.SHADOWTEXT(text, x, y, size, r, g, b, a, sr, sg, sb, sa)` 

Text with a drop shadow.

---

### `DRAW.OUTLINETEXT(text, x, y, size, r, g, b, a, or, og, ob, oa)` 

Text with an outline color.

---

### `DRAW.TEXTWIDTH(text, size)` 

Returns the pixel width of `text` at `size`.

---

### `DRAW.TEXTFONTWIDTH(font, text, size, spacing)` 

Returns the pixel width using a specific font.

---

## 2D Textures

### `DRAW.TEXTURE(tex, x, y, w, h, r, g, b, a)` 

Draws a texture scaled to `(w, h)` at `(x, y)` with color tint.

---

### `DRAW.TEXTUREEX(tex, x, y, rot, scale, r, g, b, a)` 

Texture with rotation and uniform scale.

---

### `DRAW.TEXTUREV(tex, x, y, r, g, b, a)` 

Texture at position using vector args.

---

### `DRAW.TEXTUREREC(tex, srcX, srcY, srcW, srcH, x, y, r, g, b, a)` 

Texture with source rect crop.

---

### `DRAW.TEXTUREPRO(tex, srcX, srcY, srcW, srcH, dstX, dstY, dstW, dstH, originX, originY, rot, r, g, b, a)` 

Full pro texture draw with src/dst rects, origin, and rotation.

---

### `DRAW.TEXTURETILED(...)` 

Tiled texture fill across a destination rectangle.

---

### `DRAW.TEXTUREFLIPPED(tex)` 

Draw full texture flipped vertically.

---

### `DRAW.TEXTUREFULL(tex)` 

Draw full texture at `(0,0)`.

---

### `DRAW.TEXTURENPATCH(tex, ...)` 

9-patch / n-patch texture draw for UI panels.

---

## 2D Splines

### `DRAW.SPLINELINEAR(pointsHandle, thick, r, g, b, a)` 

Linear spline through points array.

---

### `DRAW.SPLINEBEZIERQUAD(...)` / `DRAW.SPLINEBEZIERCUBIC(...)` / `DRAW.SPLINECATMULLROM(...)` / `DRAW.SPLINEBASIS(...)` 

Various spline types through a points handle.

---

## Pixel Access

### `DRAW.GETPIXELCOLOR(x, y)` 

Returns `[r, g, b, a]` of the pixel at screen position.

---

### `DRAW.SETPIXELCOLOR(x, y, r, g, b, a)` 

Sets a pixel directly.

---

## 3D Primitives

### `DRAW.CUBE(x, y, z, w, h, d, rot, r, g, b, a)` 

Filled cube at world position.

---

### `DRAW.CUBEWIRES(x, y, z, w, h, d, rot, r, g, b, a)` 

Cube wireframe.

---

### `DRAW.SPHERE(x, y, z, radius, r, g, b, a)` 

Filled sphere.

---

### `DRAW.SPHEREWIRES(x, y, z, radius, rings, slices, r, g, b, a)` 

Sphere wireframe.

---

### `DRAW.CYLINDER(x, y, z, topR, botR, h, slices, r, g, b, a)` 

Cylinder / frustum.

---

### `DRAW.CYLINDERWIRES(...)` 

Cylinder wireframe.

---

### `DRAW.CAPSULE(startX, startY, startZ, endX, endY, endZ, radius, slices, rings, r, g, b, a)` 

Filled capsule between two points.

---

### `DRAW.CAPSULEWIRES(...)` 

Capsule wireframe.

---

### `DRAW.PLANE(x, y, z, w, h, r, g, b, a)` 

Horizontal plane quad.

---

### `DRAW.BOUNDINGBOX(minX, minY, minZ, maxX, maxY, maxZ, r, g, b, a)` 

Axis-aligned bounding box wireframe.

---

### `DRAW.LINE3D(x1, y1, z1, x2, y2, z2, r, g, b, a)` 

3D line segment.

---

### `DRAW.POINT3D(x, y, z, r, g, b, a)` 

3D point.

---

### `DRAW.RAY(rayHandle, r, g, b, a)` 

Visualises a `RAY.*` handle as a line.

---

### `DRAW.BILLBOARD(tex, x, y, z, size, r, g, b, a)` 

Camera-facing billboard sprite at world position.

---

### `DRAW.BILLBOARDREC(tex, srcX, srcY, srcW, srcH, x, y, z, w, h, r, g, b, a)` 

Billboard with source rect.

---

### `DRAW.GRID(cells, cellSize)` 

XZ ground grid (3D, call inside `RENDER.BEGIN3D`).

---

## Full Example

2D HUD over a 3D scene.

```basic
WINDOW.OPEN(960, 540, "Draw Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

t = 0.0
WHILE NOT WINDOW.SHOULDCLOSE()
    t = t + TIME.DELTA()

    RENDER.CLEAR(20, 25, 35)

    RENDER.BEGIN3D(cam)
        DRAW.CUBE(0, 1, 0, 2, 2, 2, t * 30, 80, 160, 255, 255)
        DRAW.SPHERE(3, 1, 0, 1, 255, 120, 60, 255)
        DRAW.LINE3D(-5, 0, 0, 5, 0, 0, 200, 200, 200, 255)
        DRAW.GRID(20, 1.0)
    RENDER.END3D()

    ; HUD
    DRAW.RECTANGLE(10, 10, 200, 30, 0, 0, 0, 160)
    DRAW.TEXT("FPS: " + STR(GAME.FPS()), 14, 14, 18, 255, 255, 255, 255)
    DRAW.PROGRESSBAR(10, 50, 200, 16, 2, SIN(t) * 0.5 + 0.5, 40, 40, 40, 255, 80, 200, 80, 255)

    RENDER.FRAME()
WEND

CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [DRAW2D.md](DRAW2D.md) — 2D-specific draw shortcuts
- [DRAW3D.md](DRAW3D.md) — 3D-specific draw shortcuts
- [DRAWPRIM2D.md](DRAWPRIM2D.md) — retained 2D primitives
- [DRAWPRIM3D.md](DRAWPRIM3D.md) — retained 3D primitives
- [RENDER.md](RENDER.md) — frame begin/end, clear, 3D mode
