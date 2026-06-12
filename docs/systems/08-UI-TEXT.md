# UI and text systems: UI, FONT, TEXT

> Immediate-mode menus (`GUI.*`), font loading, and HUD text drawing.

**All commands:** [COMMAND_REGISTRY.md#ui-text](COMMAND_REGISTRY.md#ui-text)

**Deep guide:** [guides/UI-MENUS.md](guides/UI-MENUS.md)

**See also:** [reference/GUI.md](../reference/GUI.md) · [reference/FONT.md](../reference/FONT.md) · [examples/gui_basics](../examples/gui_basics/main.mb)

---

## Table of contents

- [UI system](#ui-system)
- [FONT system](#font-system)
- [TEXT system](#text-system)
- [Full example](#full-example)
- [See also](#see-also)

---

## UI system

Two layers: full **`GUI.*`** (raygui-style widgets) and lightweight **`UI.*`** helpers.

### GUI workflow (menus and tools)

1. Each frame after world draw (or in a dedicated UI pass): call widget functions.
2. `GUI.THEMEAPPLY("dark")` once at startup for built-in themes.
3. Use `GUI.BUTTON`, `GUI.LABEL`, `GUI.SLIDER`, etc.

**Example:**

```basic
IF GUI.BUTTON(20, 20, 120, 30, "Start") THEN
    startGame = true
ENDIF
```

Runnable demos: `examples/gui_basics`, `examples/gui_theme`, `examples/gui_form`.

---

### GUI essentials

| Command | Description |
|---------|-------------|
| `GUI.BUTTON(x, y, w, h, label)` | Clickable button — returns true when pressed |
| `GUI.LABEL(x, y, text)` | Static text |
| `GUI.SLIDER(x, y, w, label, value, min, max)` | Float slider |
| `GUI.CHECKBOX(x, y, label, checked)` | Toggle |
| `GUI.TEXTBOX(x, y, w, h, text)` | Text entry |
| `GUI.THEMEAPPLY(name)` | Apply theme |
| `GUI.SETCOLOR(control, r, g, b, a)` | Style color |

Full catalog: [reference/GUI.md](../reference/GUI.md).

---

### UI helpers (3D / simple widgets)

| Command | Description |
|---------|-------------|
| `UI.BUTTON(...)` | Simple UI button |
| `UI.LABEL3D(x, y, z, text)` | World-space label |

Checklist `UI.BEGIN` / `UI.END` immediate-mode blocks map to **`GUI.*`** patterns — see [FINAL_POLISH_SYSTEMS.md](../FINAL_POLISH_SYSTEMS.md).

---

## FONT system

TrueType fonts for crisp HUD text.

### `FONT.LOAD(path [, size])`

Loads a `.ttf` file.

| Argument | Type | Description |
|----------|------|-------------|
| path | string | Font file |
| size | int | Optional pixel size |

**Returns:** `handle`

**Example:**

```basic
font = FONT.LOAD("assets/fonts/arial.ttf", 32)
```

---

### `FONT.FREE(font)`

Releases the font handle.

---

## TEXT system

Draw strings without full GUI layout.

### Checklist aliases → canonical

| Alias | Canonical |
|-------|-----------|
| `TEXT.DRAW(msg, x, y)` | `DRAW.TEXT` |
| `TEXT.DRAWFONT(font, msg, x, y)` | `DRAW.TEXTEX` / font-aware draw |
| `TEXT.SIZE(msg)` | `DRAW.TEXTWIDTH` / measure helpers |

---

### `DRAW.TEXT(text, x, y, size, r, g, b [, a])`

Default-font HUD text — no `.ttf` required.

**Example:**

```basic
DRAW.TEXT("Score: " + score, 10, 10, 20, 255, 255, 255)
```

---

### `DRAW.TEXTEX(font, text, x, y, size, spacing, r, g, b, a)`

Draw with a loaded font handle.

**Example:**

```basic
DRAW.TEXTEX(font, "Hello", 20, 40, 24, 2, 255, 200, 100, 255)
```

---

### `TEXT.DRAW` / `TEXT.DRAWFONT` / `TEXT.SIZE`

Alias names registered for checklist compatibility — same behavior as `DRAW.*` equivalents.

**Example:**

```basic
TEXT.DRAW("Paused", 300, 280)
w = TEXT.SIZE("Hello")
```

---

## Full example

```basic
APP.OPEN(640, 480, "UI + Text")
APP.SETFPS(60)

score = 0
cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)

WHILE NOT APP.SHOULDCLOSE()
    IF GUI.BUTTON(20, 20, 100, 28, "Point") THEN score = score + 1

    RENDER.CLEAR(30, 32, 40)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()

    DRAW.TEXT("Score: " + score, 20, 60, 18, 255, 255, 255)
    RENDER.FRAME()
WEND

APP.CLOSE()
```

---

## See also

- [04-INPUT](04-INPUT.md) — mouse for GUI
- [PROGRAMMING.md](../PROGRAMMING.md) — string interpolation `$"Score: {score}"`
