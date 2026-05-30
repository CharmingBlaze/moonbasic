# UI Commands

Immediate-mode UI widgets: buttons, progress bars, inventory icons, and 3D labels. Draw each frame inside the render loop.

For sprite-based UI panels see [SPRITEBATCH.md](SPRITEBATCH.md). For full retained UI see `GUI.*` commands.

## Core Workflow

Call UI commands each frame between `RENDER.CLEAR` and `RENDER.FRAME`. They draw and optionally return interaction state.

---

## Commands

### `UI.BUTTON(label, x, y, w, h)` 

Draws a clickable button at screen position `(x, y)` with size `(w, h)`. Returns `1` on the frame it is clicked, `0` otherwise.

---

### `UI.PROGRESSBAR(x, y, w, h, value, r, g, b, a)` 

Draws a horizontal progress bar. `value` is 0.0–1.0 fill level. Color `(r, g, b, a)` for the fill.

---

### `UI.INVENTORYICON(texHandle, x, y)` 

Draws a texture as an inventory icon at screen position `(x, y)`.

---

### `UI.LABEL3D(text, entityHandle, camHandle)` 

Draws a 2D text label above a 3D entity, projected into screen space using `camHandle`.

---

## Full Example

A HUD with a health bar, button, and 3D label.

```basic
WINDOW.OPEN(800, 600, "UI Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

target = ENTITY.CREATECUBE(1.0)
ENTITY.SETPOS(target, 0, 1, 0)

hp     = 1.0
score  = 0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    ENTITY.UPDATE(dt)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW.GRID(10, 1.0)
    RENDER.END3D()

    ; HUD
    UI.PROGRESSBAR(10, 10, 200, 20, hp, 80, 200, 80, 255)
    IF UI.BUTTON("Take Damage", 10, 40, 120, 28) THEN
        hp = MAX(0, hp - 0.1)
    END IF
    IF UI.BUTTON("Heal", 140, 40, 80, 28) THEN
        hp = MIN(1, hp + 0.1)
    END IF

    UI.LABEL3D("Target", target, cam)

    RENDER.FRAME()
WEND

ENTITY.FREE(target)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [DRAW.md](DRAW.md) — `DRAW.PROGRESSBAR`, `DRAW.HEALTHBAR`, `DRAW.TEXT`
- [SPRITEBATCH.md](SPRITEBATCH.md) — sprite-based UI
- [GUI.md](GUI.md) — retained-mode GUI panels
