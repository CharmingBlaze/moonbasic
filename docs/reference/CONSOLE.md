# Console & Text Output Commands

Commands for printing text to the console or terminal where `moonbasic` is running. These are primarily used for debugging, as graphical text is handled by `Draw.Text`.

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

### `INPUT(prompt, [default])`

Prompts the user for text input from the console.

- `prompt`: The message to display to the user.
- `default`: (Optional) The value to return if the user just presses Enter.

Returns the string entered by the user.

```basic
name = INPUT("What is your name? ", "Player1")
PRINT "Hello, " + name
```

---

### `CLS()`

Clears the console screen. This uses an ANSI escape code and may not work in all terminal emulators.

---

### `LOCATE(row, column)`

Moves the console cursor to the specified row and column. Like `CLS`, this uses ANSI escape codes.

```basic
CLS()
LOCATE 10, 5
PRINT "This text is on row 10, column 5."
```
