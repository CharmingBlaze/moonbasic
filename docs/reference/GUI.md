# GUI — raygui (Raylib) widgets

moonBASIC exposes [raygui](https://github.com/raysan5/raygui) as the `GUI.*` namespace. Controls are **immediate-mode**: each frame you call the same functions; raygui tracks focus and state from rectangles and IDs.

Bindings follow [raylib-go raygui](https://github.com/gen2brain/raylib-go/tree/master/raygui) (raygui 4.x). **Every** exported control, style loader, tooltip, icon, and draw helper from that wrapper is available as a `GUI.*` command **except** font retrieval: there is no `GUI.GETFONT` because the active `rl.Font` is owned by raygui (often embedded in `.rgs` themes); double-freeing it through `FONT.FREE` would be unsafe. Use **`FONT.LOAD` + `GUI.SETFONT`** for fonts you allocate yourself.

**Requirements (full raygui):** **CGO** enabled (same as `DRAW.*`, `WINDOW.*`, `FONT.*`).

**Windows without CGO** (`CGO_ENABLED=0`, Raylib purego / `raylib.dll`): a **minimal** `GUI.*` implementation is built in—labels, buttons, basic toggles/checkboxes/combos, simple sliders, window/panel/group/line, theme defaults, and draw helpers. It is **not** pixel-identical to raygui. Commands that still need the C raygui stack (**`GUI.TEXTBOX`**, color-picker suite, **`GUI.MESSAGEBOX`**, **`GUI.LISTVIEW`**, **`GUI.TABBAR`**, `.rgs` **`GUI.LOADSTYLE`**, **`GUI.SETFONT`** with custom fonts, etc.) return an error asking you to rebuild with **CGO**. See [BUILDING.md](../BUILDING.md).

**Performance (purego GUI):** Internal widget state keys use **rounded pixel rects**. For stateful controls such as **`GUI.TOGGLEGROUP`**, use **stable** `x, y, width, height` (avoid coordinates that drift every frame). If an unusual number of distinct rectangles are used, the runtime may reset cached widget state to bound memory.

---

## Typical frame order

1. `Render.Clear(...)` (and your scene)
2. Optional: `Camera2D.Begin()` / `Camera.Begin()` if you draw in screen or 2D / 3D space
3. `GUI.*` for this frame
4. `Render.Frame()` / end camera

Call `GUI.Enable()` after `Window.Open` if you used `GUI.Disable()`.

---

## Coordinates, colors, handles

- Most widgets: **`x, y, width, height`** as numeric rectangle (float-friendly).
- **Byte colors** `r,g,b,a` (0–255) where noted.
- **RGBA out** (`GUI.GETCOLOR`, `GUI.COLORPICKER`, `GUI.FADE`, …): **4-float heap array** `[r,g,b,a]` — use `arr(0)`…`arr(3)`.
- **Rectangle out** (`GUI.GETTEXTBOUNDS`): **4-float array** `[x, y, width, height]`.
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

### Base properties (colors, border, text box layout)

| Global | Meaning |
|--------|---------|
| `GPROP_BORDER_COLOR_NORMAL` … `GPROP_TEXT_COLOR_DISABLED` | Border / fill / text per GUI state |
| `GPROP_BORDER_WIDTH` | Border thickness |
| `GPROP_TEXT_PADDING` | Inset for text |
| `GPROP_TEXT_ALIGNMENT` | Horizontal alignment index |

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

### Per-control extended (same integers, different meaning per `GCTL_*`)

Examples: `GPROP_TOGGLE_GROUP_PADDING`, `GPROP_SLIDER_WIDTH`, `GPROP_SCROLLBAR_*`, `GPROP_LIST_*`, `GPROP_COLOR_SELECTOR_SIZE`, … — see `runtime/mbgui/seed.go` for the full list and comments.

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

## Command reference — all `GUI.*`

Notation: **float**, **string**, **optional bool** where shown. **`→`** return type.

### Global, lock, alpha, font, style, themes

| Command | Arguments | Returns |
|---------|-----------|---------|
| `GUI.ENABLE` / `GUI.DISABLE` | — | — |
| `GUI.LOCK` / `GUI.UNLOCK` | — | — |
| `GUI.ISLOCKED` | — | bool |
| `GUI.SETALPHA` | `alpha` | — |
| `GUI.SETSTATE` | `state` (use `GUI_STATE_*`) | — |
| `GUI.GETSTATE` | — | int |
| `GUI.SETFONT` | `fontHandle` | — |
| `GUI.SETSTYLE` | `control, property, value` | — |
| `GUI.GETSTYLE` | `control, property` | int |
| `GUI.GETCOLOR` | `control, property` | RGBA handle |
| `GUI.SETCOLOR` | `control, property, r,g,b,a` | — |
| `GUI.SETTEXTSIZE` | `n` | — |
| `GUI.SETTEXTSPACING` | `n` | — |
| `GUI.SETTEXTLINEHEIGHT` | `n` | — |
| `GUI.SETTEXTWRAP` | `mode` (`GUI_TEXT_WRAP_*`) | — |
| `GUI.SETTEXTALIGN` | `mode` (`GUI_TEXT_ALIGN_LEFT`/`CENTER`/`RIGHT`) | — |
| `GUI.SETTEXTALIGNVERT` | `mode` (`GUI_TEXT_ALIGN_TOP`/`MIDDLE`/`BOTTOM`) | — |
| `GUI.GETTEXTSIZE` | — | int |
| `GUI.THEMEAPPLY` | `name` | — |
| `GUI.THEMENAMES` | — | string (`;`-separated theme names) |
| `GUI.LOADSTYLE` | `path` | — |
| `GUI.LOADDEFAULTSTYLE` | — | — |
| `GUI.LOADSTYLEMEM` | `path` (binary `.rgs` read on host) | — |
| `GUI.LOADICONS` | `path, loadNames` | — |
| `GUI.LOADICONSMEM` | `path, loadNames` | — |

### Layout

| Command | Arguments | Returns |
|---------|-----------|---------|
| `GUI.WINDOWBOX` | `x,y,w,h, title` | bool (close pressed) |
| `GUI.GROUPBOX` / `GUI.LINE` / `GUI.PANEL` | `x,y,w,h, text` | — |
| `GUI.TABBAR` | `x,y,w,h, tabs, stateHandle` | int (close tab or -1) |
| `GUI.SCROLLPANEL` | `px,py,pw,ph, title, cx,cy,cw,ch, stateHandle` | — |

### Basic controls

| Command | Arguments | Returns |
|---------|-----------|---------|
| `GUI.LABEL` | `x,y,w,h, text` | — |
| `GUI.BUTTON` / `GUI.LABELBUTTON` | `x,y,w,h, text` | bool |
| `GUI.TOGGLE` | `x,y,w,h, text, active` | bool |
| `GUI.TOGGLEGROUP` | `x,y,w,h, items` (`;` separated) | int (active index; starts at 0) |
| `GUI.TOGGLEGROUPAT` | `x,y,w,h, items, active` | int |
| `GUI.TOGGLESLIDER` | `x,y,w,h, text, active` | int |
| `GUI.CHECKBOX` | `x,y,w,h, text, checked` | bool |
| `GUI.COMBOBOX` | `x,y,w,h, items, active` | int |
| `GUI.DROPDOWNBOX` | `x,y,w,h, items, stateHandle` | bool |

### Text & numbers

| Command | Arguments | Returns |
|---------|-----------|---------|
| `GUI.TEXTBOX` | `x,y,w,h, text, maxLen, editMode` | string |
| `GUI.SPINNER` / `GUI.VALUEBOX` | `x,y,w,h, text, value, min, max, editMode` | int |
| `GUI.VALUEBOXFLOAT` | `x,y,w,h, label, value, textBuf, editMode` | float |
| `GUI.VALUEBOXFLOATTEXT` | — (after `VALUEBOXFLOAT` same frame) | string |

### Sliders, scroll, lists

| Command | Arguments | Returns |
|---------|-----------|---------|
| `GUI.SLIDER` / `GUI.SLIDERBAR` / `GUI.PROGRESSBAR` | `x,y,w,h, left, right, value, min, max` | float |
| `GUI.SCROLLBAR` | `x,y,w,h, value, min, max` | int |
| `GUI.STATUSBAR` / `GUI.DUMMYREC` | `x,y,w,h, text` | — |
| `GUI.LISTVIEW` | `x,y,w,h, items, stateHandle` | int |
| `GUI.LISTVIEWEX` | `x,y,w,h, items, stateHandle` | int |

### Color pickers & dialogs

| Command | Arguments | Returns |
|---------|-----------|---------|
| `GUI.COLORPANEL` / `GUI.COLORPICKER` | `x,y,w,h, text, r,g,b,a` | RGBA handle |
| `GUI.COLORBARALPHA` / `GUI.COLORBARHUE` | `x,y,w,h, text, value` | float |
| `GUI.COLORPICKERHSV` / `GUI.COLORPANELHSV` | `x,y,w,h, text, hsvHandle` | int |
| `GUI.MESSAGEBOX` | `x,y,w,h, title, message, buttons` | int |
| `GUI.TEXTINPUTBOX` | `x,y,w,h, title, message, buttons, text, maxLen, secretHandle` | int |
| `GUI.TEXTINPUTLAST` | — (after `TEXTINPUTBOX`) | string |
| `GUI.GRID` | `x,y,w,h, text, spacing, subdivs, cellHandle` | int |

### Tooltips & icons

| Command | Arguments | Returns |
|---------|-----------|---------|
| `GUI.ENABLETOOLTIP` / `GUI.DISABLETOOLTIP` | — | — |
| `GUI.SETTOOLTIP` | `text` | — |
| `GUI.ICONTEXT` | `iconId, text` | string |
| `GUI.DRAWICON` | `iconId, x, y, pixelSize, r,g,b,a` | — |
| `GUI.SETICONSCALE` | `scale` | — |
| `GUI.GETTEXTWIDTH` | `text` | int |

### Styled drawing helpers

| Command | Arguments | Returns |
|---------|-----------|---------|
| `GUI.FADE` | `r,g,b,a, alpha` | RGBA handle |
| `GUI.DRAWRECTANGLE` | `x,y,w,h, borderW, br,bg,bb,ba, fr,fg,fb,fa` | — |
| `GUI.DRAWTEXT` | `text, x,y,w,h, align, r,g,b,a` | — |
| `GUI.GETTEXTBOUNDS` | `control, x,y,w,h` | rect handle |

---

## `GUI.THEMEAPPLY` names

- **`DEFAULT`** / **`RESET`** — `GUI.LOADDEFAULTSTYLE`
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

## Design notes

- **No `GUI.GETFONT`** — see introduction; use `GUI.GETTEXTSIZE` / `GUI.GETTEXTWIDTH` / `GUI.GETSTYLE` for metrics.
- **Raygui C API extras** not wrapped by raylib-go (e.g. raw icon pointer access) are unavailable until added upstream.
