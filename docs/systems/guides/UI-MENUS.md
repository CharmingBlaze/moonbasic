# UI and menus — HUD, buttons, and settings screens

> Draw **text**, **buttons**, **sliders**, and **themes** on top of your game — immediate-mode UI each frame.

**Namespaces:** `GUI` · `DRAW` · `TEXT` · `FONT` · `UI` · **Status:** Shipped · **Platform:** full runtime

**Commands:** [COMMAND_REGISTRY.md#ui-text](../COMMAND_REGISTRY.md#ui-text) · [reference/GUI.md](../../reference/GUI.md)

**Examples:** [`examples/gui_basics`](../../../examples/gui_basics/main.mb) · [`examples/gui_theme`](../../../examples/gui_theme/main.mb) · [`examples/gui_form`](../../../examples/gui_form/main.mb)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use which UI API](#when-to-use-which-ui-api)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [HUD text without fonts](#hud-text-without-fonts)
- [GUI widgets](#gui-widgets)
- [Themes and styling](#themes-and-styling)
- [Stateful widgets](#stateful-widgets)
- [Full example — menu + HUD](#full-example--menu--hud)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Layer | API | Mental model |
|-------|-----|--------------|
| **HUD numbers** | `DRAW.TEXT` | Draw string each frame at x,y |
| **Custom font** | `FONT.LOAD` + `DRAW.TEXTEX` | TTF sizing and spacing |
| **Menus / forms** | `GUI.BUTTON`, `SLIDER`, … | Immediate mode — **no widget objects** |
| **Aliases** | `TEXT.DRAW` | Same as `DRAW.TEXT` |

**Why immediate mode:** Widgets are recreated every frame from your code. State lives in **your variables** (`volume = 0.8`), not hidden UI objects.

---

## When to use which UI API

| Need | Use |
|------|-----|
| Score, FPS, debug line | `DRAW.TEXT` |
| Title screen buttons | `GUI.BUTTON` |
| Options slider | `GUI.SLIDER` |
| Text entry name | `GUI.TEXTBOX` |
| 3D label in world | `UI.LABEL3D` (see reference) |

**Draw order:** Usually **3D scene first** (`RENDER.BEGIN` … `END`), then **GUI/HUD**, then `RENDER.FRAME()`. Or GUI-only apps: clear → GUI → frame.

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| “Score: 42” overlay | `DRAW.TEXT` | `GUI.LABEL` for simple text |
| Clickable menu | `GUI.BUTTON` | Manual mouse rect + `INPUT` (unless learning) |
| Styled dark theme | `GUI.THEMEAPPLY("dark")` | Hand-draw every rect |
| Layout engine / retained UI | External or custom | moonBASIC GUI is immediate-mode |

---

## Core workflow

1. **Game draw** (optional 3D/2D world).
2. **GUI calls** — each returns `true` when activated (buttons).
3. **Update game state** from GUI results (`IF GUI.BUTTON(...) THEN start = true`).
4. **`RENDER.FRAME()`** — presents everything.

```basic
RENDER.CLEAR(30, 32, 40)
; ... world draw ...
IF GUI.BUTTON(20, 20, 120, 28, "Play") THEN playing = 1
DRAW.TEXT("HP: " + hp, 10, 10, 18, 255, 255, 255)
RENDER.FRAME()
```

**Why every frame:** Immediate-mode UI does not persist — skipping a frame skips input handling for that widget.

---

## HUD text without fonts

`DRAW.TEXT(text, x, y, size, r, g, b [, a])`

**Why:** Default bitmap font — no `assets/font.ttf` required. Good for prototypes.

```basic
DRAW.TEXT("Press SPACE", 300, 280, 20, 255, 255, 255, 255)
```

**Aliases:** `TEXT.DRAW` — same behavior.

**Measure width:** `DRAW.TEXTWIDTH` / `TEXT.SIZE` for centering.

---

## GUI widgets

| Widget | Returns | Typical use |
|--------|---------|-------------|
| `GUI.BUTTON(x,y,w,h, label)` | `true` if pressed | Menus |
| `GUI.LABEL(x,y, text)` | — | Static text |
| `GUI.SLIDER(x,y,w, label, value, min, max)` | new value | Volume, difficulty |
| `GUI.CHECKBOX(x,y, label, checked)` | new checked | Options |
| `GUI.TEXTBOX(x,y,w,h, text)` | edited string | Name entry |
| `GUI.DROPDOWNBOX` / `LISTVIEW` | selection | Lists — need `DIM` state arrays |

**Example button:**

```basic
IF GUI.BUTTON(100, 200, 160, 32, "Start Game") THEN
    mode = 1
ENDIF
```

**Why rectangle args:** `x, y, width, height` in **screen pixels** (top-left origin).

---

## Themes and styling

**Why themes:** One call sets colors, padding, and font for all widgets.

```basic
GUI.THEMEAPPLY("dark")    ; built-in name — see GUI.THEMENAMES
```

| Command | Why |
|---------|-----|
| `GUI.THEMEAPPLY(name)` | Load bundled raygui style |
| `GUI.LOADSTYLE(path)` | Custom `.rgs` file |
| `GUI.SETCOLOR(control, prop, r,g,b,a)` | Per-control tweak |
| `GUI.SETFONT(fontHandle)` | After `FONT.LOAD` |

Globals `GCTL_*` (control kind) and `GPROP_*` (property id) match raygui — see [reference/GUI.md](../../reference/GUI.md) tables.

---

## Stateful widgets

**Why `DIM` arrays:** List views, scroll panels, and dropdowns need persistent index state across frames.

```basic
DIM listItems(8)
; populate listItems — see examples/gui_form
```

Copy patterns from **`examples/gui_form/main.mb`** — do not guess array sizes.

---

## Full example — menu + HUD

```basic
APP.OPEN(640, 480, "UI demo")
APP.SETFPS(60)

GUI.THEMEAPPLY("dark")
playing = 0
score = 0
volume = 0.7

WHILE NOT APP.SHOULDCLOSE()
    IF playing = 0 THEN
        RENDER.CLEAR(25, 28, 35)
        IF GUI.BUTTON(220, 200, 200, 36, "Play") THEN playing = 1
        volume = GUI.SLIDER(220, 260, 200, "Volume", volume, 0, 1)
    ELSE
        RENDER.CLEAR(15, 18, 28)
        IF INPUT.KEYHIT(KEY_SPACE) THEN score = score + 1
        DRAW.TEXT("Score: " + score, 20, 20, 20, 255, 220, 100)
        IF GUI.BUTTON(20, 400, 100, 28, "Menu") THEN playing = 0
    ENDIF
    RENDER.FRAME()
WEND

APP.CLOSE()
```

`moonrun` required.

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| GUI before clear | `RENDER.CLEAR` first |
| Forget `RENDER.FRAME` | UI never appears |
| Store button state only in GUI | Use your `playing` / `volume` variables |
| Huge `TEXTBOX` without `DIM` buffer | Follow gui_form example |
| Mix GUI and 3D without `RENDER.END` | End 3D pass before GUI overlay |

---

## See also

- [08-UI-TEXT.md](../08-UI-TEXT.md) — overview
- [04-INPUT.md](../04-INPUT.md) — mouse for custom hit tests
- `moonbasic new --template ui` — starter project
