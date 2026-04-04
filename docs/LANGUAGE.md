# moonBASIC Language Reference

This document covers the core features of the moonBASIC language.

For **built-in APIs** (window, draw, time, files, …), how to structure a game loop, and platform notes, see the [Programming Guide](PROGRAMMING.md) and [Command Index](COMMANDS.md).

---

## Variables and Types

Variables are created when you first assign a value to them. Their type is determined by a suffix character at the end of the variable name.

| Suffix | Type      | Example                  |
|--------|-----------|--------------------------|
| `$`    | String    | `name$ = "Player 1"`     |
| `#`    | Float     | `speed# = 150.5`         |
| `?`    | Boolean   | `alive? = TRUE`          |
| (none) | Integer   | `score = 1000`           |

If no suffix is provided, the variable is an integer by default.

### Scope

Variables are global by default. You can declare variables with local scope inside functions using the `LOCAL` keyword.

```basic
FUNCTION MyFunc()
    LOCAL message$ = "This is a local string"
    PRINT message$
ENDFUNCTION
```

- **`GLOBAL`**: Explicitly declare a variable as global (this is the default behavior).
- **`LOCAL`**: Declare a variable that only exists within the current `FUNCTION`.
- **`STATIC`**: Declare a variable inside a function that retains its value between function calls.

---

## Control Flow

moonBASIC supports standard control flow structures.

### IF / THEN / ELSE

For conditional logic. `ELSEIF` and `ELSE` are optional.

```basic
IF score > 1000 THEN
    PRINT "High score!"
ELSEIF score > 500 THEN
    PRINT "Good job!"
ELSE
    PRINT "Try again!"
ENDIF
```

### SELECT / CASE

For choosing between multiple conditions. `DEFAULT` is optional.

```basic
SELECT fruit$
    CASE "apple"
        PRINT "An apple a day..."
    CASE "banana"
        PRINT "Potassium!"
    DEFAULT
        PRINT "That's not an apple or a banana."
ENDSELECT
```

### FOR / NEXT

Loops a specific number of times. `STEP` is optional and defaults to 1.

```basic
; Count from 1 to 10
FOR i = 1 TO 10
    PRINT i
NEXT

; Count down from 10 to 1 by -2
FOR i = 10 TO 1 STEP -2
    PRINT i
NEXT
```

### WHILE / WEND

Loops as long as a condition is true.

```basic
x = 0
WHILE x < 5
    PRINT x
    x = x + 1
WEND
```

### REPEAT / UNTIL

Loops until a condition becomes true. The loop body is always executed at least once.

```basic
x = 10
REPEAT
    PRINT x
    x = x - 1
UNTIL x = 0
```

### DO / LOOP

The `DO...LOOP` structure is flexible and can be combined with `WHILE` or `UNTIL` at the start or end of the loop.

```basic
; Loop while condition is true
DO WHILE x < 10
    x = x + 1
LOOP

; Loop until condition is true
DO UNTIL x = 10
    x = x + 1
LOOP

; Post-condition checks
DO
    x = x - 1
LOOP WHILE x > 0

DO
    x = x - 1
LOOP UNTIL x = 0
```

### Exiting Loops

You can exit a loop early using the `EXIT` command followed by the loop type.

- `EXIT FOR`
- `EXIT WHILE`
- `EXIT REPEAT`
- `EXIT DO`

---

## Functions

Create reusable blocks of code with `FUNCTION`. Functions can accept parameters and return a value.

```basic
; A function that takes two numbers and returns their sum
FUNCTION Add(a, b)
    RETURN a + b
ENDFUNCTION

result = Add(5, 10)
PRINT result ; Outputs 15
```

### Returning Values

Use the `RETURN` keyword to send a value back from a function. The type of the returned value is determined by the value itself (e.g., `RETURN 5` returns an integer, `RETURN "hello"` returns a string).

### Exiting Functions

You can exit a function at any point using `EXIT FUNCTION`.

```basic
FUNCTION CheckValue(val)
    IF val < 0 THEN
        PRINT "Value cannot be negative."
        EXIT FUNCTION
    ENDIF
    PRINT "Value is valid."
ENDFUNCTION
```
