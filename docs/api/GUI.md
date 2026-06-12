# GUI Commands

Commands for creating immediate-mode GUI widgets using the raygui library. moonBASIC's GUI system provides buttons, labels, sliders, checkboxes, text inputs, dropdowns, spinners, color pickers, progress bars, and more. Widgets are drawn every frame and return user interaction results.

## Core Concepts

- **Immediate mode** — GUI widgets are drawn and polled every frame. There is no persistent widget state — you call the function, it draws the widget and returns the result.
- **Rectangle layout** — Every widget takes `x, y, width, height` to define its screen rectangle.
- **Return values** — Interactive widgets return whether they were clicked, their current value, or the selected index.
- **Themes** — GUI appearance can be customized with `.rgs` style files.
- **raygui backend** — All widgets delegate to the raygui C library embedded in the Raylib build.

---

## Buttons

### `GUI.Button(x, y, w, h, text)`

Draws a button and returns `TRUE` if it was clicked this frame.

- `x`, `y` (int) — Position.
- `w`, `h` (int) — Size.
- `text` (string) — Button label.

**Returns:** `bool`

```basic
IF GUI.Button(10, 10, 120, 30, "Start Game") THEN
    gameStarted = TRUE
ENDIF
```

---

### `GUI.ButtonLabel(x, y, w, h, text)`

Draws a label-style button (no background, just text that's clickable).

---

## Labels & Text

### `GUI.Label(x, y, w, h, text)`

Draws a text label.

- `text` (string) — Label text.

```basic
GUI.Label(10, 50, 200, 20, "Player Name:")
```

---

### `GUI.TextBox(x, y, w, h, text, maxChars, editMode)`

Draws an editable text box. Returns `TRUE` when the user presses Enter or clicks away.

- `text` (string) — Current text content.
- `maxChars` (int) — Maximum character count.
- `editMode` (bool) — `TRUE` if the text box is active for editing.

**Returns:** `bool` (toggled edit state)

```basic
playerName = "Hero"
editing = FALSE

; In game loop:
IF GUI.TextBox(10, 70, 200, 30, playerName, 20, editing) THEN
    editing = NOT editing
ENDIF
```

---

### `GUI.TextInput(x, y, w, h, text, maxLen)`

Simplified text input that manages edit state internally.

---

## Sliders

### `GUI.Slider(x, y, w, h, textLeft, textRight, value, minVal, maxVal)`

Draws a horizontal slider.

- `textLeft` (string) — Label on the left.
- `textRight` (string) — Label on the right.
- `value` (float) — Current value.
- `minVal`, `maxVal` (float) — Value range.

**Returns:** `float` — New value after user interaction.

```basic
volume = GUI.Slider(10, 100, 200, 20, "Vol", "100%", volume, 0.0, 1.0)
```

---

### `GUI.SliderBar(x, y, w, h, textLeft, textRight, value, minVal, maxVal)`

Same as `Slider` but with a filled bar visual.

---

## Checkboxes & Toggles

### `GUI.CheckBox(x, y, w, h, text, checked)`

Draws a checkbox.

- `checked` (bool) — Current state.

**Returns:** `bool` — New state.

```basic
showFPS = GUI.CheckBox(10, 130, 20, 20, "Show FPS", showFPS)
```

---

### `GUI.Toggle(x, y, w, h, text, active)`

Draws a toggle button.

**Returns:** `bool`

---

### `GUI.ToggleGroup(x, y, w, h, text, activeIndex)`

Draws a group of toggle buttons (radio-style).

- `text` (string) — Semicolon-separated labels (e.g., `"Easy;Normal;Hard"`).
- `activeIndex` (int) — Currently selected index.

**Returns:** `int` — New selected index.

```basic
difficulty = GUI.ToggleGroup(10, 160, 300, 30, "Easy;Normal;Hard", difficulty)
```

---

## Progress & Spinners

### `GUI.ProgressBar(x, y, w, h, textLeft, textRight, value, minVal, maxVal)`

Draws a progress bar.

- `value` (float) — Current progress.
- `minVal`, `maxVal` (float) — Range.

```basic
GUI.ProgressBar(10, 200, 200, 20, "HP", "100", health, 0, 100)
```

---

### `GUI.Spinner(x, y, w, h, text, value, minVal, maxVal, editMode)`

Draws a number spinner with +/- buttons.

**Returns:** `int` — New value.

---

### `GUI.ValueBox(x, y, w, h, text, value, minVal, maxVal, editMode)`

Draws an editable numeric value box.

**Returns:** `int`

---

## Dropdowns & Lists

### `GUI.DropdownBox(x, y, w, h, text, activeIndex, editMode)`

Draws a dropdown selection box.

- `text` (string) — Semicolon-separated options.
- `activeIndex` (int) — Currently selected index.
- `editMode` (bool) — `TRUE` when the dropdown is expanded.

**Returns:** `int` — New selection or toggle state.

```basic
weapons = "Sword;Bow;Staff;Axe"
selectedWeapon = GUI.DropdownBox(10, 240, 150, 30, weapons, selectedWeapon, dropdownOpen)
```

---

### `GUI.ListView(x, y, w, h, text, scrollIndex, activeIndex)`

Draws a scrollable list view.

- `text` (string) — Semicolon-separated items.
- `scrollIndex` (int) — Current scroll position.
- `activeIndex` (int) — Selected item index.

**Returns:** `int` — New selected index.

---

## Color Picker

### `GUI.ColorPicker(x, y, w, h, color)`

Draws a full color picker with hue bar.

- `color` (handle) — Current color handle.

**Returns:** `handle` — Selected color.

---

### `GUI.ColorPanel(x, y, w, h, color)`

Draws just the saturation/value panel without the hue bar.

---

### `GUI.ColorBarHue(x, y, w, h, hue)`

Draws a standalone hue selection bar.

---

## Panels & Groups

### `GUI.Panel(x, y, w, h)`

Draws a panel background (for grouping widgets).

---

### `GUI.GroupBox(x, y, w, h, text)`

Draws a labeled group box outline.

```basic
GUI.GroupBox(5, 5, 220, 180, "Settings")
; Draw widgets inside...
```

---

### `GUI.Line(x, y, w, h, text)`

Draws a horizontal separator line with optional text.

---

### `GUI.ScrollPanel(x, y, w, h, contentRect, scrollX, scrollY)`

Creates a scrollable panel area.

---

## Themes

### `GUI.ThemeApply(styleName)`

Applies a built-in raygui style theme.

- `styleName` (string) — Theme name (e.g., `"dark"`, `"jungle"`, `"candy"`, `"cherry"`, `"cyber"`, `"lavanda"`, `"terminal"`).

**How it works:** Loads the corresponding `.rgs` style file from the bundled `raygui_styles/` directory and applies it globally to all widgets.

```basic
GUI.ThemeApply("dark")
```

---

### `GUI.SetStyle(control, property, value)`

Sets a specific style property for fine-grained control.

---

### `GUI.GetStyle(control, property)`

Gets a style property value.

---

### `GUI.SetFont(fontHandle)` / `GUI.SetFontSize(size)`

Sets the GUI font and size.

---

## Full Example

A settings menu with various GUI widgets.

```basic
Window.Open(800, 600, "GUI Demo")
Window.SetFPS(60)

GUI.ThemeApply("dark")

volume = 0.7
brightness = 0.8
showFPS = TRUE
difficulty = 1
resolution = 0
playerName = "Hero"
nameEditing = FALSE

WHILE NOT Window.ShouldClose()
    Render.Clear(30, 30, 40)

    ; Title
    GUI.Label(20, 15, 200, 30, "GAME SETTINGS")
    GUI.Line(20, 45, 360, 10, "")

    ; Player name
    GUI.Label(20, 60, 100, 20, "Name:")
    IF GUI.TextBox(130, 60, 200, 25, playerName, 16, nameEditing) THEN
        nameEditing = NOT nameEditing
    ENDIF

    ; Difficulty
    GUI.Label(20, 100, 100, 20, "Difficulty:")
    difficulty = GUI.ToggleGroup(130, 100, 200, 25, "Easy;Normal;Hard", difficulty)

    ; Resolution
    GUI.Label(20, 140, 100, 20, "Resolution:")
    resolution = GUI.DropdownBox(130, 140, 200, 25, "1280x720;1920x1080;2560x1440", resolution, FALSE)

    ; Volume slider
    GUI.Label(20, 180, 100, 20, "Volume:")
    volume = GUI.Slider(130, 180, 200, 20, "0", "100", volume, 0.0, 1.0)

    ; Brightness slider
    GUI.Label(20, 210, 100, 20, "Brightness:")
    brightness = GUI.SliderBar(130, 210, 200, 20, "", "", brightness, 0.0, 1.0)

    ; Show FPS checkbox
    showFPS = GUI.CheckBox(20, 250, 20, 20, "Show FPS Counter", showFPS)

    ; Apply button
    IF GUI.Button(130, 300, 120, 35, "Apply") THEN
        PRINT "Settings applied!"
        PRINT "Name: " + playerName
        PRINT "Volume: " + STR(INT(volume * 100)) + "%"
    ENDIF

    ; Show FPS if enabled
    IF showFPS THEN
        Draw.Text("FPS: " + STR(Window.GetFPS()), 700, 10, 16, 100, 255, 100, 255)
    ENDIF

    Render.Frame()
WEND

Window.Close()
```

---

## See Also

- [WINDOW](WINDOW.md) — Window creation
- [INPUT](INPUT.md) — Input for custom widget interactions
- [DRAW](DRAW.md) — Custom drawing for non-raygui UI
- [FONT](FONT.md) — Custom fonts for GUI text
