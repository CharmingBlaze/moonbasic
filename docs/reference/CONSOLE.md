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

### `Console.Input(prompt, [default])`

Prompts the user for text input from the console.

- `prompt`: The message to display to the user.
- `default`: (Optional) The value to return if the user just presses Enter.

Returns the string entered by the user.

```basic
name = Console.Input("What is your name? ", "Player1")
Console.Print("Hello, " + name)
```

---

### `Console.Log(message)`

Writes a message to the internal console.

---

### `Console.Clear()`

Clears the console buffer.

---

### `Console.Show()` / `Console.Hide()`

Toggles the console overlay.

---

### `Console.SetColor(r, g, b, a)`

Sets console text color.

### `Console.SetBackground(r, g, b, a)`

Sets console background color.

---

### `Console.Locate(row, column)`

Moves the console cursor to the specified row and column. Like `Console.Clear`, this uses ANSI escape codes.

```basic
Console.Clear()
Console.Locate(10, 5)
Console.Print("This text is on row 10, column 5.")
