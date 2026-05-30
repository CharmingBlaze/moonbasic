# moonBASIC Language Reference

This document covers the core features of the moonBASIC language.

For **built-in APIs** (window, draw, time, files, …), how to structure a game loop, and platform notes, see the [Programming Guide](PROGRAMMING.md), [Command Index](COMMANDS.md), and the full registry [API_CONSISTENCY.md](API_CONSISTENCY.md).

---

## Variables and Types

Variables are created when you first assign a value to them. Their type is determined implicitly by the value assigned.

| Type      | Example                  |
|-----------|--------------------------|
| String    | `name = "Player 1"`      |
| Float     | `speed = 150.5`          |
| Boolean   | `alive = TRUE`           |
| Integer   | `score = 1000`           |

The language is dynamically typed; a variable can hold any value (implicit `Any` type).

**No `#` / `$` / `?` / `%` suffixes:** moonBASIC does not use Blitz-style suffix characters on identifiers (`name#`, `name$`, etc.). Use normal names and typing rules above; see [STYLE_GUIDE.md](../STYLE_GUIDE.md).

**Case insensitivity:** Keywords, variables, functions, and built-in command names are matched **without regard to letter case** — `score`, `Score`, and `SCORE` are the same identifier; `Window.Open` and `WINDOW.open` are the same call. Prefer a consistent style in new code; see [STYLE_GUIDE.md](../STYLE_GUIDE.md). Implementation detail only: the lexer records a normalized spelling for names; commands and symbols are matched against the manifest and VM using the usual **uppercase** dotted registry keys — you do not need to type those keys yourself in source.

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

### Constants and enums

**`CONST`** declares a compile-time constant at module scope:

```basic
CONST MAX_HP = 100
```

**`ENUM`** groups related integer constants that auto-increment from `0`:

```basic
ENUM State
    IDLE
    WALK
    RUN
    JUMP
ENDENUM

IF playerState = State.IDLE THEN ...
```

Each member becomes a global constant named **`EnumName_MEMBER`** (for example **`STATE_IDLE`**, **`STATE_WALK`**). You can also read members as **`State.IDLE`**, **`State.WALK`**, and so on — the compiler resolves these to the same integer values.

### String interpolation

Use **`$"..."`** for inline expressions in string literals (similar to C# / Swift):

```basic
PRINT($"Score: {score}  Health: {hp}")
PRINT($"Angle: {yaw:.2f}")
```

- **`{expr}`** is converted with **`STR(expr)`** and concatenated.
- **`{expr:fmt}`** uses **`FORMAT(expr, "%fmt")`** (for example **`{hp:.1f}`**).

This is compile-time desugaring — no new runtime opcode. For positional placeholders, **`STRING.INTERP("Hi {0}", name)`** still works; see [STRING.md](reference/STRING.md).

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

### FOR EACH … IN

Iterate every element of a **`DIM`** array (1-based indices internally):

```basic
DIM enemies(10)
FOR EACH e IN enemies
    IF e.hp <= 0 THEN ENTITY.FREE(e.mesh)
NEXT
```

The loop variable **`e`** is assigned each element in order. This works with typed arrays (`DIM arr AS Type(n)`) and untyped numeric arrays.

### FOR … = EACH(Type)

Iterate every live **`NEW(Type)`** instance (tracked by the VM at runtime):

```basic
TYPE Enemy
    FIELD hp
ENDTYPE

FOR e = EACH(Enemy)
    IF e.hp <= 0 THEN DELETE e
NEXT
```

The loop variable is a handle to each instance. Instances are registered when created with **`NEW`** and removed when **`DELETE`**d.

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

### Multiple return values

Return several values with a comma-separated **`RETURN`**, then unpack them with comma-separated assignment:

```basic
FUNCTION GetPlayerPos()
    RETURN px, py, pz
ENDFUNCTION

x, y, z = GetPlayerPos()
```

The compiler packs multiple values into a temporary **1-based array** for the caller; **`a, b, c = expr`** unpacks that array into separate variables. You do **not** need to **`ERASE`** the temporary — only **`ERASE`** handles you explicitly stored in a **`DIM`** variable.

For helpers that still return a **single float** (for example **`BOXTOPLAND`**), assign to one variable — see [GAMEHELPERS.md](reference/GAMEHELPERS.md).

**Legacy pattern (still valid):** return a **`DIM`** array handle yourself when you need a long-lived tuple on the heap:

```basic
FUNCTION PlatformSnap(...)
    DIM r(2)
    r(0) = 1.0
    r(1) = landY
    RETURN r
ENDFUNCTION
h = PlatformSnap(...)
; read h(1), then ERASE h when done
```

Use **`LOCAL`** inside `FUNCTION` for temporaries so names do not collide with globals.

### Function references

Pass callbacks without fragile string names using **`@FunctionName`** or an anonymous function literal:

```basic
FUNCTION OnHit(a, b)
    PRINT "hit!"
ENDFUNCTION

PHYSICS3D.ONCOLLISION(bodyA, bodyB, @OnHit)
TWEEN.ONCOMPLETE(tw, @OnDone)

; Anonymous callback
onHit = FUNCTION(a, b)
    PRINT "hit!"
ENDFUNCTION

; Call a stored reference
cb = @OnHit
cb(bodyA, bodyB)
```

The compiler emits a first-class function-reference value (`KindFunc`). APIs that accept callbacks accept either a **string name** (legacy) or **`@func`** / a **`FUNCTION() … ENDFUNCTION`** expression assigned to a variable.

### Coroutines

Use block syntax or **`COROUTINE.START`**:

```basic
COROUTINE patrol
    WHILE TRUE
        PRINT "step"
        COROUTINE.WAIT(1.0)
        YIELD
    WEND
ENDCOROUTINE

; patrol handle is started automatically; resumes each frame
```

Or the explicit form:

```basic
FUNCTION Patrol()
    WHILE TRUE
        PRINT "step"
        COROUTINE.WAIT(1.0)
        YIELD
    WEND
ENDFUNCTION

co = COROUTINE.START(@Patrol)
```

### Typed function signatures (optional)

Add **`AS type`** on parameters and return values for compile-time checks:

```basic
FUNCTION Add(a AS FLOAT, b AS FLOAT) AS FLOAT
    RETURN a + b
ENDFUNCTION
```

Supported types: **`INTEGER`**, **`FLOAT`**, **`STRING`**, **`BOOL`**, or a user **`TYPE`** name. Multi-return: **`AS INT, INT`**.

Set asset roots before loading files:

```basic
ASSET.PATH("assets/")
model = MODEL.LOAD("player.glb")
```

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

## INCLUDE and IMPORT (splitting programs across files)

**`INCLUDE "path.mb"`** merges another **`.mb`** file at **compile time** (not **`.md`** — docs under `docs/` are reference only). Full behavior: path resolution, include-once, cycles, shared globals/functions, and **`MOONBASIC_PATH`**: see the WAVE reference **[INCLUDE.md](reference/INCLUDE.md)**.

**`IMPORT "package"`** loads a package entry file from configured package roots (for example `physics_helper/main.mb`). Search paths are set via **`MOONBASIC_PKG`** / **`SetPackageRoots`** — same infrastructure as shared libraries. Use **`INCLUDE`** for files next to your game; use **`IMPORT`** for reusable packages installed under a packages folder. See [ROADMAP.md](ROADMAP.md) for the full package manifest story.

Quick facts: case-insensitive keyword; paths resolve relative to the **including** file, then optional package roots; duplicate includes of the same file are skipped; circular includes error at compile time; **no** per-frame runtime cost (merge is compile-time only).
