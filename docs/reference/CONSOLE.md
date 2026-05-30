# Console Commands

Text output, input prompts, and console overlay control for debugging and terminal I/O.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Use `PRINT` / `WRITE` for quick debug output. For an in-game console overlay, use `CONSOLE.SHOW` / `CONSOLE.LOG` / `CONSOLE.CLEAR`. For graphical text rendering see [DRAW2D.md](DRAW2D.md).

---

### `PRINT(args...)` 

Prints one or more values to the console, separated by spaces, followed by a newline character.

```basic
name = "moonBASIC"
version = 0.1
PRINT "Welcome to", name, "version", version
; Output: Welcome to moonBASIC version 0.1
```

---

### `WRITE(args...)` 

Same as `PRINT`, but does *not* add a newline character at the end.

```basic
WRITE "Loading..."
; ... do some work ...
PRINT "Done!"
; Output: Loading...Done!
```

---

### `CONSOLE.INPUT(prompt, default)` 

Prompts the user for text input from the console. Returns the string entered, or `default` if the user presses Enter without typing.

- `prompt`: The message to display.
- `default`: Fallback value (optional).

---

### `CONSOLE.LOG(message)` 

Writes a message to the internal console overlay.

---

### `CONSOLE.CLEAR()` 

Clears the console buffer.

---

### `CONSOLE.SHOW()` / `CONSOLE.HIDE()` 

Shows or hides the console overlay.

---

### `CONSOLE.SETCOLOR(r, g, b, a)` 

Sets console text color (0–255 per channel).

---

### `CONSOLE.SETBACKGROUND(r, g, b, a)` 

Sets console background color (0–255 per channel).

---

### `CONSOLE.LOCATE(row, column)` 

Moves the console cursor to the specified row and column (ANSI escape codes).

---

## Full Example

This example uses the console overlay for debug logging.

```basic
CONSOLE.SHOW()
CONSOLE.SETCOLOR(0, 255, 0, 255)
CONSOLE.LOG("Game started")

WHILE NOT WINDOW.SHOULDCLOSE()
    CONSOLE.LOG("Frame: " + STR(TIME.GETFRAME()))
    RENDER.BEGINFRAME()
    RENDER.ENDFRAME()
WEND

CONSOLE.HIDE()
