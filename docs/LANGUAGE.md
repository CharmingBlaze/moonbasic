# moonBASIC Language Reference

This document covers the core features of the moonBASIC language.

For **built-in APIs** (window, draw, time, files, …), how to structure a game loop, and platform notes, see the [Programming Guide](PROGRAMMING.md), [Command Index](COMMANDS.md), and the full registry [API_CONSISTENCY.md](API_CONSISTENCY.md).

---

## Variables and Types

## Variables and Types

Variables are created when you first assign a value to them. Their type is determined implicitly by the value assigned.

| Type      | Example                  |
|-----------|--------------------------|
| String    | `name = "Player 1"`      |
| Float     | `speed = 150.5`          |
| Boolean   | `alive = TRUE`           |
| Integer   | `score = 1000`           |

The language is dynamically typed; a variable can hold any value (implicit `Any` type).

**Case insensitivity:** Language keywords and built-in command names (for example `Namespace.Method` / `NAMESPACE.METHOD`) are matched **without regard to letter case**. Prefer a consistent style in new code; see [STYLE_GUIDE.md](../STYLE_GUIDE.md).

### Scope

Variables are global by default. You can declare variables with local scope inside functions using the `LOCAL` keyword.

```basic
FUNCTION MyFunc()
    LOCAL message = "This is a local string"
    PRINT message
ENDFUNCTION
```

- **`GLOBAL`**: Explicitly declare a variable as global (this is the default behavior).
- **`LOCAL`**: Declare a variable that only exists within the current `FUNCTION`.
- **`STATIC`**: Declare a variable inside a function that retains its value between function calls.

### Record types (`TYPE` … `ENDTYPE`)

You can define **named record types** at global scope (before use). Types are **value data only** (no methods).

```basic
TYPE Platform
    x, y, z
    w, h, d
    r, g, b
ENDTYPE
```

Allocate a **typed array** with **`DIM name AS TypeName(count)`**. Indices run from `0` to `count - 1`. Set an element with the **`TypeName(...)` constructor**, passing values in **field declaration order**:

```basic
CONST N = 4
DIM plat AS Platform(N)
plat(0) = Platform(0.0, 1.5, 6.0, 4.0, 0.4, 4.0, 255, 60, 200)
PRINT plat(0).x
plat(0).r = 200
```

Read and write fields with **dot notation** on an indexed element: `arr(i).field = expr`.

**`ERASE(name)`** frees a typed array the same way as other `DIM` arrays when you are done with it. See [Array commands](reference/ARRAY.md) for `DIM`, lengths, and heap behaviour.

**`ERASE ALL`** (the identifier **`ALL`**, case-insensitive) frees **every** VM heap object (arrays, cameras, textures, models, etc.), then sets **all** global and operand-stack values that held a handle to **null**. Equivalent callable form: **`FREE.ALL`**. Use at shutdown or scene resets — not in the middle of an expression. Does **not** replace **`ENTITY.CLEARSCENE`** / **`ENTITY.FREE`** for numeric entity IDs. Avoid naming a variable **`ALL`** if you need per-variable **`ERASE`**. Details: [MEMORY.md](MEMORY.md).

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
SELECT fruit
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

### Multiple results with a float array

There is no multi-value `RETURN` tuple yet. For landing on a box top from a sphere, **`BOXTOPLAND`** returns a **single float** (landing centre Y or `0.0`) — see [GAMEHELPERS.md](reference/GAMEHELPERS.md). For other cases where you need two or more numbers, pack them into a **small array** and return that handle. The caller reads `result(0)`, `result(1)`, then **`ERASE(result)`** when done.

```basic
FUNCTION PlatformSnap(px, py, pz, pvy, pr, bx, by, bz, bw, bh, bd)
    DIM r(2)
    r(0) = 0.0
    r(1) = py
    landY = BOXTOPLAND(px, py, pz, pvy, pr, bx, by, bz, bw, bh, bd)
    IF landY > 0.0 THEN
        r(0) = 1.0
        r(1) = landY
    ENDIF
    RETURN r
ENDFUNCTION

h = PlatformSnap(px, py, pz, pvy, pr, 0.0, 1.5, 6.0, 4.0, 0.4, 4.0)
IF h(0) THEN
    py = h(1)
ENDIF
ERASE h
```

Use **`LOCAL`** inside `FUNCTION` for temporaries so names do not collide with globals.

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

---

## INCLUDE (splitting programs across files)

Use **`INCLUDE "path.mb"`** to merge another source file **at compile time**. The path is resolved relative to the **file that contains the `INCLUDE`**. If the file is not found there, the compiler searches extra roots from **`MOONBASIC_PATH`** and installed package directories (see [PACKAGES.md](PACKAGES.md)).

- **Case**: The keyword is case-insensitive (`include` and `INCLUDE` are the same).
- **Order**: Types and functions from the included file are pulled in first; then its top-level statements are inserted where the `INCLUDE` line was.
- **Duplicates**: Including the same file more than once (even from different parents) is **ignored after the first occurrence** — the file is parsed once, so shared modules like `game.mb` or `menu.mb` do not duplicate code or work.
- **Cycles**: `A` includes `B` includes `A` is an error (clear message at compile time).

Example — main file pulls in a menu module:

```basic
INCLUDE "menu.mb"

Window.Open(800, 600, "Game")
; ... rest of main ...
Window.Close()
```

Where `menu.mb` might define shared functions or data used by the main game.

Compile-time only: there is **no** runtime cost per frame from `INCLUDE`; the merged program is what gets bytecode-generated. Transient parse memory uses the compiler arena during `moonbasic` / `CompileFile` and is released after compilation.

See also: [MEMORY.md](MEMORY.md) (brief note on compile-time merge).
