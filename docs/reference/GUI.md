# GUI Commands

Immediate-mode raygui widgets: buttons, sliders, text boxes, color pickers, themes, and more.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. `RENDER.CLEAR` (and your scene).
2. Call `GUI.*` widget functions each frame (immediate mode — no persistent widget objects).
3. `RENDER.FRAME`.

Apply themes with `GUI.THEMEAPPLY`, customize with `GUI.SETSTYLE` / `GUI.SETCOLOR`. Requires **CGO** for full raygui; a minimal subset works without CGO on Windows.

---

## Coordinates, colors, handles

- Most widgets: **`x, y, width, height`** as numeric rectangle (float-friendly).
- **Byte colors** `r,g,b,a` (0–255) where noted.
- **RGBA out** (`GUI.Getcolor`, `GUI.Colorpicker`, `GUI.Fade`, …): **4-float heap array** `[r,g,b,a]` — use `arr(0)`…`arr(3)`.
- **Rectangle out** (`GUI.GetTextbounds`): **4-float array** `[x, y, width, height]`.
- **Stateful widgets** use **`DIM`** numeric arrays passed **by handle**; sizes are documented below.

List separators in strings use **`;`**, matching raygui (tabs, combos, list rows).

---

## Customizing appearance (everything you can change)

1. **Whole look + embedded font** — `GUI.THEMEAPPLY(name)` for built-in / bundled [raygui styles](https://github.com/raysan5/raygui/tree/master/styles), or `GUI.LOADSTYLE(path)` / `GUI.LOADSTYLEMEM(path)` for a binary `.rgs` file. Names accepted by `GUI.THEMEAPPLY` are listed by **`GUI.THEMENAMES`** (semicolon-separated).
2. **Reset** — `GUI.LOADDEFAULTSTYLE` or `GUI.THEMEAPPLY("DEFAULT")` / `"RESET"`.
3. **Per-property colors** — `GUI.SETCOLOR(control, property, r,g,b,a)` with **`GCTL_*`** + **`GPROP_*`** globals (see tables below).
4. **Per-property integers** — `GUI.SETSTYLE(control, property, value)` / `GUI.GETSTYLE`; use the same **`GCTL_*`** / **`GPROP_*`** IDs. Numeric style values include border width, padding, sizes, alignment indices, etc.
5. **Global text metrics** — `GUI.SETTEXTSIZE`, `SETTEXTSPACING`, `SETTEXTLINEHEIGHT`, `SETTEXTWRAP`, `SETTEXTALIGN`, `SETTEXTALIGNVERT`, `GUI.GETTEXTSIZE`.
6. **Font you control** — `FONT.LOAD` … then `GUI.SETFONT(fontHandle)`. After a theme load, the theme’s font is active until you set another.
7. **Transparency & focus** — `GUI.SETALPHA`, `GUI.SETSTATE` / `GUI.GETSTATE` with **`GUI_STATE_*`**.
8. **Icons** — `GUI.LOADICONS` / `GUI.LOADICONSMEM`, `GUI.SETICONSCALE`, `GUI.ICONTEXT`, `GUI.DRAWICON`.
9. **Low-level drawing** — `GUI.DRAWRECTANGLE`, `GUI.DRAWTEXT`, `GUI.GETTEXTBOUNDS` (respect current style).

**Reading colors back:** `GUI.GETCOLOR(control, property)` → RGBA array handle.

---

## Style & theme globals (`GCTL_*`, `GPROP_*`, …)

VM globals match raygui **ControlID** / **PropertyID** integers. **Important:** many **`GPROP_*`** values reuse the same number for different controls; always pass the **`GCTL_*`** that matches the widget you are styling.

### Control IDs (`GCTL_*`) 

| Global | Typical use |
|--------|-------------|
| `GCTL_DEFAULT` | Global defaults, text size, line/background colors, shared padding |
| `GCTL_LABEL` … `GCTL_STATUSBAR` | Per–control-kind styling (`GCTL_BUTTON`, `GCTL_SLIDER`, `GCTL_LISTVIEW`, …) |

(Full set: `GCTL_DEFAULT`, `GCTL_LABEL`, `GCTL_BUTTON`, `GCTL_TOGGLE`, `GCTL_SLIDER`, `GCTL_PROGRESSBAR`, `GCTL_CHECKBOX`, `GCTL_COMBOBOX`, `GCTL_DROPDOWNBOX`, `GCTL_TEXTBOX`, `GCTL_VALUEBOX`, `GCTL_CONTROL11`, `GCTL_LISTVIEW`, `GCTL_COLORPICKER`, `GCTL_SCROLLBAR`, `GCTL_STATUSBAR`.)

---

### Base properties (colors, border, text box layout) 

| Global | Meaning |
|--------|---------|
| `GPROP_BORDER_COLOR_NORMAL` … `GPROP_TEXT_COLOR_DISABLED` | Border / fill / text per GUI state |
| `GPROP_BORDER_WIDTH` | Border thickness |
| `GPROP_TEXT_PADDING` | Inset for text |
| `GPROP_TEXT_ALIGNMENT` | Horizontal alignment index |

---

### DEFAULT-only (global) extended 

| Global | Meaning |
|--------|---------|
| `GPROP_TEXT_SIZE` | Font pixel height |
| `GPROP_TEXT_SPACING` | Glyph spacing |
| `GPROP_LINE_COLOR` | Line control color |
| `GPROP_BACKGROUND_COLOR` | Panel-style background |
| `GPROP_TEXT_LINE_SPACING` | Line spacing |
| `GPROP_TEXT_ALIGNMENT_VERTICAL` | Vertical alignment |
| `GPROP_TEXT_WRAP_MODE` | Wrap mode |

---

### Per-control extended (same integers, different meaning per `GCTL_*`) 

Examples: `GPROP_TOGGLE_GROUP_PADDING`, `GPROP_SLIDER_WIDTH`, `GPROP_SCROLLBAR_*`, `GPROP_LIST_*`, `GPROP_COLOR_SELECTOR_SIZE`, … — see `runtime/mbgui/seed.go` for the full list and comments.

---

### Other enums 

| Globals | Use |
|---------|-----|
| `GUI_STATE_NORMAL` … `GUI_STATE_DISABLED` | `GUI.SETSTATE` / `GUI.GETSTATE` |
| `GUI_TEXT_ALIGN_LEFT` / `CENTER` / `RIGHT` | Horizontal text |
| `GUI_TEXT_ALIGN_TOP` / `MIDDLE` / `BOTTOM` | Vertical text |
| `GUI_TEXT_WRAP_NONE` / `CHAR` / `WORD` | `GUI.SETTEXTWRAP` |
| `GUI_SCROLLBAR_LEFT` / `GUI_SCROLLBAR_RIGHT` | List scrollbar side |

**`GUI.GETTEXTBOUNDS(control, x,y,w,h)`** — `control` is a **`GCTL_*`** value (raygui `ControlID`).

---

## Stateful array layouts

| Command | Array (numeric `DIM`) | Notes |
|---------|----------------------|--------|
| `GUI.TABBAR` | **1** float: active tab index | Returns close-tab index or **-1** |
| `GUI.SCROLLPANEL` | **6** floats: `scrollX, scrollY, viewX, viewY, viewW, viewH` | Updated each frame |
| `GUI.DROPDOWNBOX` | **2** floats: `activeItem`, `editMode` (0/1) | Return = dropdown open? |
| `GUI.LISTVIEW` | **2** floats: `scrollIndex`, `activeLine` | Return = selected index |
| `GUI.LISTVIEWEX` | **3** floats: `focus`, `scrollIndex`, `active` | Return = selected index |
| `GUI.COLORPICKERHSV` / `COLORPANELHSV` | **3** floats: H, S, V | In/out |
| `GUI.TEXTINPUTBOX` | **1** float: secret view 0/1 | In/out |
| `GUI.GRID` | **2** floats: mouse cell | In/out |

---

## Command reference — all `Gui.*`

Notation: **float**, **string**, **optional bool** where shown. **`→`** return type.

### Global, lock, alpha, font, style, themes 

| Command | Arguments | Returns |
|---------|-----------|---------|
| `Gui.Enable()` / `Gui.Disable()` | — | — |
| `Gui.Lock()` / `Gui.Unlock()` | — | — |
| `Gui.IsLocked()` | — | bool |
| `Gui.SetAlpha(alpha)` | `alpha` | — |
| `Gui.SetState(state)` | `state` (use `GUI_STATE_*`) | — |
| `Gui.GetState()` | — | int |
| `Gui.SetFont(font)` | `fontHandle` | — |
| `Gui.SetStyle(control, property, value)` | `control, property, value` | — |
| `Gui.GetStyle(control, property)` | `control, property` | int |
| `Gui.GetColor(control, property)` | `control, property` | RGBA handle |
| `Gui.SetColor(control, property, r, g, b, a)` | `control, property, r, g, b, a` | — |
| `Gui.SetTextSize(n)` | `n` | — |
| `Gui.SetTextSpacing(n)` | `n` | — |
| `Gui.SetTextLineHeight(n)` | `n` | — |
| `Gui.SetTextWrap(mode)` | `mode` (`GUI_TEXT_WRAP_*`) | — |
| `Gui.SetTextAlign(mode)` | `mode` (`GUI_TEXT_ALIGN_LEFT`/`CENTER`/`RIGHT`) | — |
| `Gui.SetTextAlignVert(mode)` | `mode` (`GUI_TEXT_ALIGN_TOP`/`MIDDLE`/`BOTTOM`) | — |
| `Gui.GetTextSize()` | — | int |
| `Gui.ThemeApply(name)` | `name` | — |
| `Gui.ThemeNames()` | — | string (`;`-separated theme names) |
| `Gui.LoadStyle(path)` | `path` | — |
| `Gui.LoadDefaultStyle()` | — | — |
| `Gui.LoadStyleMem(path)` | `path` (binary `.rgs` read on host) | — |
| `Gui.LoadIcons(path, loadNames)` | `path, loadNames` | — |

---

### Layout 

| Command | Arguments | Returns |
|---------|-----------|---------|
| `Gui.WindowBox(x, y, w, h, title)` | `x, y, w, h, title` | bool (close pressed) |
| `Gui.GroupBox(x, y, w, h, text)` / `Gui.Line(x, y, w, h, text)` / `Gui.Panel(x, y, w, h, text)` | `x, y, w, h, text` | — |
| `Gui.TabBar(x, y, w, h, tabs, stateHandle)` | `x, y, w, h, tabs, stateHandle` | int (close tab or -1) |
| `Gui.ScrollPanel(px, py, pw, ph, title, cx, cy, cw, ch, stateHandle)` | `px, py, pw, ph, title, cx, cy, cw, ch, stateHandle` | — |

---

### Basic controls 

---

### `GUI.BUTTON(label, x, y, w, h)`
Draws a clickable button.

- **Arguments**:
    - `label`: (String) Text to display.
    - `x, y, w, h`: (Float) Rectangle dimensions.
- **Returns**: (Boolean) `TRUE` if clicked this frame.

---

### `GUI.LABEL(label, x, y, w, h)`
Draws static text.

- **Returns**: (None)

---

### `GUI.TOGGLE(label, x, y, w, h, active)`
Draws a toggle button.

- **Arguments**:
    - `active`: (Boolean) Current state.
- **Returns**: (Boolean) The new state.

---

### `GUI.SLIDER(label, x, y, value, min, max, w, h)`
Draws a horizontal slider.

- **Arguments**:
    - `value`: (Float) Current value.
    - `min, max`: (Float) Range.
- **Returns**: (Float) The updated value.

---

### `GUI.CHECKBOX(label, x, y, w, h, checked)`
Draws a checkbox.

- **Returns**: (Boolean) The new checked state.

---

### `GUI.TEXTBOX(label, x, y, w, h, text, maxLen, editMode)`
Draws a text input box.

- **Arguments**:
    - `text`: (String) Current text.
    - `maxLen`: (Integer) Max characters.
    - `editMode`: (Boolean) Whether it is currently being edited.
- **Returns**: (String) The updated text.

---

### `GUI.LISTVIEW(label, x, y, w, h, items, stateHandle)`
Draws a list view.

- **Arguments**:
    - `items`: (String) Semicolon-separated list.
    - `stateHandle`: (Handle) DIM array for scroll/active state.
- **Returns**: (Integer) The selected index.

---

### Color pickers & dialogs 

| Command | Arguments | Returns |
|---------|-----------|---------|
| `Gui.ColorPanel(x, y, w, h, text, r, g, b, a)` / `Gui.ColorPicker(x, y, w, h, text, r, g, b, a)` | `x, y, w, h, text, r, g, b, a` | RGBA handle |
| `Gui.ColorBarAlpha(x, y, w, h, text, value)` / `Gui.ColorBarHue(x, y, w, h, text, value)` | `x, y, w, h, text, value` | float |
| `Gui.ColorPickerHSV(x, y, w, h, text, hsvHandle)` / `Gui.ColorPanelHSV(x, y, w, h, text, hsvHandle)` | `x, y, w, h, text, hsvHandle` | int |
| `Gui.MessageBox(x, y, w, h, title, message, buttons)` | `x, y, w, h, title, message, buttons` | int |
| `Gui.TextInputBox(x, y, w, h, title, message, buttons, text, maxLen, secretHandle)` | `x, y, w, h, title, message, buttons, text, maxLen, secretHandle` | int |
| `Gui.TextInputLast()` | — (after `TextInputBox`) | string |
| `Gui.Grid(x, y, w, h, text, spacing, subdivs, cellHandle)` | `x, y, w, h, text, spacing, subdivs, cellHandle` | int |

---

### Tooltips & icons 

| Command | Arguments | Returns |
|---------|-----------|---------|
| `Gui.EnableTooltip()` / `Gui.DisableTooltip()` | — | — |
| `Gui.SetTooltip(text)` | `text` | — |
| `Gui.IconText(iconId, text)` | `iconId, text` | string |
| `Gui.DrawIcon(iconId, x, y, pixelSize, r, g, b, a)` | `iconId, x, y, pixelSize, r, g, b, a` | — |
| `Gui.SetIconScale(scale)` | `scale` | — |
| `Gui.GetTextWidth(text)` | `text` | int |

---

### Styled drawing helpers 

| Command | Arguments | Returns |
|---------|-----------|---------|
| `Gui.Fade(r, g, b, a, alpha)` | `r, g, b, a, alpha` | RGBA handle |
| `Gui.DrawRectangle(x, y, w, h, borderW, br, bg, bb, ba, fr, fg, fb, fa)` | `x, y, w, h, borderW, br, bg, bb, ba, fr, fg, fb, fa` | — |
| `Gui.DrawText(text, x, y, w, h, align, r, g, b, a)` | `text, x, y, w, h, align, r, g, b, a` | — |
| `Gui.GetTextBounds(control, x, y, w, h)` | `control, x, y, w, h` | rect handle |

---

## `Gui.ThemeApply` names

- **`DEFAULT`** / **`RESET`** — `Gui.LoadDefaultStyle`
- **`LIGHT`** — built-in light palette (no `.rgs`)
- **`BUILTIN_DARK`** — small built-in dark palette (no `.rgs`)
- **Bundled raygui `.rgs`:** `AMBER`, `ASHES`, `BLUISH`, `CANDY`, `CHERRY`, `CYBER`, `DARK`, `ENEFETE`, `GENESIS`, `JUNGLE`, `LAVANDA`, `RLTECH`, `SUNNY`, `TERMINAL` (case-insensitive)

Call **`GUI.THEMENAMES`** for an authoritative `;`-separated list at runtime.

---

## Examples in this repo

- `examples/gui_basics/main.mb` — window, label, button  
- `examples/gui_theme/main.mb` — `GUI.THEMEAPPLY`, text sizing  
- `examples/gui_form/main.mb` — fields, slider, checkbox, tabs  

---

## Full Example

```basic
WINDOW.OPEN(640, 480, "GUI Demo")
WINDOW.SETFPS(60)

sliderVal = 50
checked = FALSE

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(40, 40, 50)

    IF GUI.BUTTON("Click Me", 20, 20, 120, 30) THEN
        PRINT "Button pressed!"
    ENDIF

    GUI.LABEL("Volume:", 20, 70, 80, 30)
    sliderVal = GUI.SLIDER("", 110, 70, sliderVal, 0, 100, 200, 30)

    checked = GUI.CHECKBOX("Enable", 20, 120, 20, 20, checked)

    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## Design notes

- **No `GUI.GETFONT`** — see introduction; use `GUI.GETTEXTSIZE` / `GUI.GETTEXTWIDTH` / `GUI.GETSTYLE` for metrics.
- **Raygui C API extras** not wrapped by raylib-go (e.g. raw icon pointer access) are unavailable until added upstream.
