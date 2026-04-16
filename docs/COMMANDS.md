# moonBASIC Command Index

This page is a **topic-oriented** index of moonBASIC built-ins (globals and `NAMESPACE.NAME` APIs). It does not name every overload; the compiler manifest is the source of truth.

**How to use these commands in a real program:** see the [Programming Guide](PROGRAMMING.md) (game loop, types, CGO) and runnable projects under [`examples/`](../examples/README.md). Copy-paste snippets also live in [Examples](EXAMPLES.md).

## Complete registry (every command)

| Resource | Purpose |
|----------|---------|
| [`compiler/builtinmanifest/commands.json`](../compiler/builtinmanifest/commands.json) | Machine-readable manifest: one row per overload for the compiler and tools (count changes as APIs evolve). |
| [API_CONSISTENCY.md](API_CONSISTENCY.md) | Human-readable list: **every** registered name with argument kinds, grouped by namespace. Regenerate: `go run ./tools/apidoc`. Optional per-row `description` in the manifest appears here when set. |
| [COMMAND_AUDIT.md](COMMAND_AUDIT.md) | **Namespace → doc map:** overload counts per namespace, primary reference page, one-line blurb, and file-exists checks. Regenerate: `go run ./tools/cmdaudit`. |
| [API_CONVENTIONS.md](reference/API_CONVENTIONS.md) | **Cross-type naming:** `LOAD` vs `MAKE`, `SETPOS`, scale/rotate patterns, and aliases. |
| [reference/BLITZ_COMMAND_INDEX.md](reference/BLITZ_COMMAND_INDEX.md) | **Blitz3D / BlitzPlus → moonBASIC:** familiar names (`Plot`, `CreateCube`, `CameraZoom`, …) mapped to dotted APIs and notes on parity. |
| [reference/dbpro/README.md](reference/dbpro/README.md) | **DarkBASIC Professional (DBPro) → moonBASIC:** modular section files (`01-objects-3d.md`, …) mapping DBPro commands to namespaces. |
| [reference/moonbasic-command-set/README.md](reference/moonbasic-command-set/README.md) | **Designed MoonBASIC command set** (Blitz spirit / DBPro power / simple API): modular tables with **memory** column and real **`NAMESPACE.NAME`** mappings. |
| [MEMORY.md](MEMORY.md) | **`FREE`** / **`ERASE`**, VM heap tags (including physics joints and network packets where applicable). |

Use **API_CONSISTENCY.md** when you need to verify that a name exists or which arity the manifest allows. Use **COMMAND_AUDIT.md** for a high-level map of all dotted namespaces (`WINDOW`, `RENDER`, …) and where they are documented.

## Engine namespaces

Dotted commands are grouped by their first segment (e.g. `WINDOW.OPEN` → namespace `WINDOW`). The **authoritative table** (counts, links, and a short explanation per namespace) is **[COMMAND_AUDIT.md](COMMAND_AUDIT.md)** — it is generated from the manifest so it cannot drift. Topics below still give narrative detail, DONE/PARTIAL status, and examples.

---

**Legend:**
- **[DONE]**: Implemented, tested, and ready to use.
- **[PARTIAL]**: Partially implemented or missing tests. May not work as expected.
- **[MISSING]**: Not yet implemented.

---

## Output / Input

- `PRINT(args...)` **[DONE]**: Prints values to the console, separated by spaces, with a newline.
- `PRINTLN(args...)` **[DONE]**: Same as `PRINT()`.
- `WRITE(args...)` **[DONE]**: Prints values without a trailing newline.
- `INPUT(prompt)` **[DONE]**: Prompts the user for console input.
- `CLS()` **[DONE]**: Clears the console screen.
- `LOCATE(row, column)` **[DONE]**: Moves the console cursor to a specific row and column.
- `TAB` **[PARTIAL]**: Coming soon.
- `SPC` **[PARTIAL]**: Coming soon.

---

## Type Conversion

- `INT(value)` **[DONE]**: Converts a value to an integer.
- `FLOAT(value)` **[DONE]**: Converts a value to a float.
- `STR(value)` **[DONE]**: Converts a value to a string.
- `VAL(string)` **[DONE]**: Parses a string to a float number.
- `ASC(string)` **[DONE]**: Returns the ASCII code for the first character of a string.
- `CHR(code)` **[DONE]**: Returns a string from an ASCII code.
- `BOOL(value)` **[PARTIAL]**: Coming soon.
- `FIX(value)` **[DONE]**: Truncates a float toward zero (e.g. `FIX(-3.7)` → `-3`).
- `TYPEOF(variable)` **[DONE]**: Returns the type of a variable as a string.
- `ISNULL(value)` **[DONE]**: Checks if a value is null.
- `ISHANDLE(value)` **[DONE]**: Checks if a value is a handle.
- `ISTYPE(variable, type)` **[DONE]**: Checks if a variable is of a certain type.

---

## String Manipulation

- `LEN(string)` **[DONE]**: Returns the length of a string.
- `LEFT(string, count)` **[DONE]**: Returns characters from the left side.
- `RIGHT(string, count)` **[DONE]**: Returns characters from the right side.
- `MID(string, start, [count])` **[DONE]**: Extracts a substring.
- `UPPER(string)` **[DONE]**: Converts to uppercase.
- `LOWER(string)` **[DONE]**: Converts to lowercase.
- `INSTR(string, sub, [start])` **[DONE]**: Finds substring position.
- `TRIM(string)` **[DONE]**: Removes whitespace.
- `REPLACE(string, old, new)` **[DONE]**: Replaces occurrences.
- `CONTAINS(string, sub)` **[DONE]**: Checks for substring.
- `STARTSWITH(string, sub)` **[DONE]**: Checks prefix.
- `ENDSWITH(string, sub)` **[DONE]**: Checks suffix.
- `SPLIT(string, sep)` **[DONE]**: Splits into array.
- `JOIN(array, sep)` **[DONE]**: Joins array into string.
- `BIN(value)` **[DONE]**: Integer to binary string.
- `HEX(value)` **[DONE]**: Integer to hex string.

---

## Math

All math functions are also available as `MATH.name(...)` (e.g. `MATH.SIN`, `MATH.CLAMP`). See [Math Reference](reference/MATH.md) for full details.

- `SIN(angle)` **[DONE]**: Returns the sine of an angle (radians).
- `COS(angle)` **[DONE]**: Returns the cosine of an angle (radians).
- `TAN(angle)` **[DONE]**: Returns the tangent of an angle (radians).
- `ATN(angle)` **[DONE]**: Returns the arctangent of a value (radians). Also available as `ATAN`.
- `ASIN(value)` **[DONE]**: Returns the arcsine (radians).
- `ACOS(value)` **[DONE]**: Returns the arccosine (radians).
- `ATAN(value)` **[DONE]**: Returns the arctangent (radians). Alias for `ATN`.
- `ATAN2(y, x)` **[DONE]**: Returns the two-argument arctangent (radians).
- `SQR(value)` **[DONE]**: Returns the square root. Alias for `SQRT`.
- `SQRT(value)` **[DONE]**: Returns the square root.
- `EXP(value)` **[DONE]**: Returns e raised to the power of `value`.
- `LOG(value)` **[DONE]**: Returns the natural logarithm.
- `LOG2(value)` **[DONE]**: Returns the base-2 logarithm.
- `LOG10(value)` **[DONE]**: Returns the base-10 logarithm.
- `POW(base, exp)` **[DONE]**: Returns `base` raised to the power of `exp`.
- `FLOOR(value)` **[DONE]**: Rounds down to the nearest integer.
- `CEIL(value)` **[DONE]**: Rounds up to the nearest integer.
- `ROUND(value, [decimals])` **[DONE]**: Rounds to nearest integer or to `decimals` decimal places.
- `ABS(value)` **[DONE]**: Returns the absolute value.
- `SGN(value)` **[DONE]**: Returns -1, 0, or 1 depending on the sign.
- `FIX(value)` **[DONE]**: Truncates toward zero (like `INT` but for floats).
- `MIN(a, b)` **[DONE]**: Returns the smaller of two values.
- `MAX(a, b)` **[DONE]**: Returns the larger of two values.
- `CLAMP(value, min, max)` **[DONE]**: Clamps a value between min and max.
- `LERP(a, b, t)` **[DONE]**: Linearly interpolates between `a` and `b` by `t` (0.0-1.0).
- `SMOOTHSTEP(lo, hi, x)` **[DONE]**: Smooth interpolation between 0 and 1.
- `PINGPONG(t, length)` **[DONE]**: Bounces `t` back and forth between 0 and `length`.
- `WRAP(value, min, max)` **[DONE]**: Wraps `value` within the range [min, max].
- `DEG2RAD(degrees)` **[DONE]**: Converts degrees to radians.
- `RAD2DEG(radians)` **[DONE]**: Converts radians to degrees.
- `WRAPANGLE(angle)` **[DONE]**: Wraps an angle to the range [0, 360).
- `WRAPANGLE180(angle)` **[DONE]**: Wraps an angle to the range [-180, 180).
- `ANGLEDIFF(from, to)` **[DONE]**: Returns the shortest signed angle difference in degrees.
- `RND([limit])` **[DONE]**: Returns a random integer from 0 to `limit-1`, or a random float [0,1) if no argument.
- `RNDF(min, max)` **[DONE]**: Returns a random float between `min` and `max`.
- `RNDSEED(seed)` **[DONE]**: Seeds the random number generator.
- `RANDOMIZE([seed])` **[DONE]**: Seeds the RNG from a value, or from the clock if omitted.
- `PI()` **[DONE]**: Returns π (3.14159...).
- `TAU()` **[DONE]**: Returns τ = 2π.
- `E()` **[DONE]**: Returns Euler's number (2.71828...).
- `MOVEX(yaw, forward, strafe)` **[DONE]**: Camera-relative world X on the XZ plane — see [Math Reference](reference/MATH.md).
- `MOVEZ(yaw, forward, strafe)` **[DONE]**: Camera-relative world Z on the XZ plane — see [Math Reference](reference/MATH.md).
- `MOVESTEPX` / `MOVESTEPZ` **[DONE]**: **`MOVEX`/`MOVEZ` × speed × `dt`** — see [Math Reference](reference/MATH.md).
- `LANDBOXES(...)` **[DONE]**: Best **`BOXTOPLAND`** snap Y over parallel box arrays — see [Game helpers](reference/GAMEHELPERS.md).
- `PLAYER.MOVERELATIVE(...)` **[DONE]**: 2-float **`[dx,dz]`** handle — same as **`MOVESTEPX`/`MOVESTEPZ`** — see [Game helpers](reference/GAMEHELPERS.md).
- **Less math** (`INPUT.MOVEDIR`, `INPUT.MOUSEDELTA`, `MATH.CIRCLEPOINT`, `VEC2.DIST` / `DISTSQ`, `VEC2.PUSHOUT`, `TERRAIN.SNAPY`, `WORLD.SETCENTERENTITY`, `ENTITY.GETXZ`, …) **[DONE]**: Shortcuts for spawn rings, distance checks, camera-relative WASD, terrain snap — see [Less math](reference/LESS_MATH.md). Additional **`MATH.*`** table helpers (`HDIST`, `YAWFROMXZ`, `SMOOTHERSTEP`, …) — see [Game math helpers](reference/GAME_MATH_HELPERS.md). Blitz-style engine helpers (`CAMERA.UNPROJECT`, `RAY.INTERSECTSMODEL_*`, `LIGHT.CREATEPOINT`, `SPRITE.PLAY`, `RES.PATH`, …) — see [Game engine patterns](reference/GAME_ENGINE_PATTERNS.md).
- `REMAP` / `MATH.REMAP` / `INVERSE_LERP` / `MATH.INVERSE_LERP` / `SATURATE` / `MATH.SATURATE` **[DONE]**: Range mapping and **[0,1]** clamp — see [Math Reference](reference/MATH.md).
- `VEC3.DIST` / `VEC3.DISTSQ` **[DONE]**: Scalar 3D distance / squared distance — see [Vec2/Vec3/Quat](reference/VEC_QUAT.md).
- `INTERP` / `STRING.INTERP` **[DONE]**: `"{0}"`…`"{9}"` template fill — see [STRING.md](reference/STRING.md); hot-path notes — [STRING_HEAP.md](reference/STRING_HEAP.md).
- `COLOR.TOHSV(handle)` **[DONE]**: `(h,s,v)` tuple from a color — see [COLOR.md](reference/COLOR.md).
- `VEC2.LENGTH(x, y)` / `VEC3.LENGTH(x, y, z)` **[DONE]**: Scalar vector lengths (no vector handle allocation).
- `VEC2.NORMALIZE(x, y)` / `VEC3.NORMALIZE(x, y, z)` **[DONE]**: Scalar normalize helpers returning tuple-like arrays for destructuring.
- `VEC2.MOVE_TOWARD(fromX, fromY, toX, toY, maxDist)` **[DONE]**: Move toward target by max distance; returns `(x, y)` tuple-like array.
- `ENTITY.GETPOS(entity)` **[DONE]**: Returns `(x, y, z)` tuple-like array for destructuring assignment.
- `COLOR.FROMHSV(h, s, v)` / `COLOR.CLAMP(r, g, b)` **[DONE]**: Procedural color helpers (handle color and tuple clamp).

---

## Logic

- `IIF(condition, trueVal, falseVal)` **[DONE]**: Returns `trueVal` if condition is true, otherwise `falseVal`. Inline if-then-else. **Both branches are evaluated.**
- `IIF(condition, trueVal, falseVal)` **[DONE]**: String variant of `IIF` — see [Math Reference](reference/MATH.md).
- `CHOOSE(index, val1, val2, ...)` **[DONE]**: Returns the value at position `index` (1-based) from the argument list.
- `SWITCH(expr, case1, val1, ..., caseN, valN, default)` **[DONE]**: Returns the value paired with the first matching case, or `default` if none match.

---

## Array Operations

- `DIM` **[DONE]**: Declares an array (including **`DIM name AS TypeName(n)`** for record types — see [LANGUAGE.md](LANGUAGE.md)).
- `REDIM` **[PARTIAL]**: Coming soon.
- `ERASE` **[DONE]**: Frees a `DIM` or typed array — see [ARRAY.md](reference/ARRAY.md). **`ERASE ALL`** frees the entire VM heap; **`FREE.ALL`** is the same as a callable — [MEMORY.md](MEMORY.md).
- `ARRAYLEN` **[DONE]**: Returns the length of an array.
- `ARRAYFILL` **[PARTIAL]**: Coming soon.
- `ARRAYCOPY` **[PARTIAL]**: Coming soon.
- `ARRAYSORT` **[PARTIAL]**: Coming soon.
- `ARRAYREVERSE` **[PARTIAL]**: Coming soon.
- `ARRAYFIND` **[PARTIAL]**: Coming soon.
- `ARRAYCONTAINS` **[PARTIAL]**: Coming soon.
- `ARRAYPUSH` **[PARTIAL]**: Coming soon.
- `ARRAYPOP` **[PARTIAL]**: Coming soon.
- `ARRAYSHIFT` **[PARTIAL]**: Coming soon.
- `ARRAYUNSHIFT` **[PARTIAL]**: Coming soon.
- `ARRAYSPLICE` **[PARTIAL]**: Coming soon.
- `ARRAYSLICE` **[PARTIAL]**: Coming soon.
- `ARRAYJOINS` **[PARTIAL]**: Coming soon.

---

## File I/O

- `OPENFILE` **[PARTIAL]**: Coming soon.
- `CLOSEFILE` **[PARTIAL]**: Coming soon.
- `READFILE` **[PARTIAL]**: Coming soon.
- `FILE.WRITE` **[DONE]**: Writes raw bytes to a file without a newline.
- `FILE.WRITELN` **[DONE]**: Writes a string to a file, followed by a newline.
- `WRITEFILE` **[DONE]**: Alias for `FILE.WRITE`.
- `WRITEFILELN` **[DONE]**: Alias for `FILE.WRITELN`.
- `READALLTEXT` **[DONE]**: Reads the entire content of a file into a string.
- `WRITEALLTEXT` **[DONE]**: Writes a string to a file, overwriting existing content.
- `READBYTE` **[PARTIAL]**: Coming soon.
- `WRITEBYTE` **[PARTIAL]**: Coming soon.
- `READSHORT` **[PARTIAL]**: Coming soon.
- `WRITESHORT` **[PARTIAL]**: Coming soon.
- `READINT` **[PARTIAL]**: Coming soon.
- `WRITEINT` **[PARTIAL]**: Coming soon.
- `READFLOAT` **[PARTIAL]**: Coming soon.
- `WRITEFLOAT` **[PARTIAL]**: Coming soon.
- `READSTRING` **[PARTIAL]**: Coming soon.
- `WRITESTRING` **[PARTIAL]**: Coming soon.
- `FILEPOS` **[PARTIAL]**: Coming soon.
- `SEEKFILE` **[PARTIAL]**: Coming soon.
- `EOF` **[PARTIAL]**: Coming soon.
- `FILESIZE` **[PARTIAL]**: Coming soon.
- `FILEEXISTS` **[DONE]**: Checks if a file exists.
- `DIREXISTS` **[DONE]**: Checks if a directory exists.
- `DELETEFILE` **[DONE]**: Deletes a file.
- `COPYFILE` **[DONE]**: Copies a file.
- `MOVEFILE` **[DONE]**: Moves or renames a file.
- `RENAMEFILE` **[DONE]**: Renames a file.
- `MAKEDIR` **[DONE]**: Creates a directory.
- `MAKEDIRS` **[DONE]**: Creates a directory and all parent directories.
- `DELETEDIR` **[DONE]**: Deletes a directory.
- `GETDIR` **[DONE]**: Gets the current working directory.
- `SETDIR` **[DONE]**: Sets the current working directory.
- `GETFILES` **[PARTIAL]**: Coming soon.
- `GETDIRS` **[DONE]**: Gets a list of subdirectories in a path.
- `GETFILEEXT` **[DONE]**: Returns the extension of a file path.
- `GETFILENAME` **[DONE]**: Returns the file name from a path.
- `GETFILENAMENOEXT` **[DONE]**: Returns the file name without the extension.
- `GETFILEPATH` **[DONE]**: Returns the directory path from a file path.
- `GETFILESIZE` **[DONE]**: Returns the size of a file in bytes.
- `GETFILEMODTIME` **[DONE]**: Returns the last modification time of a file.

---

## Date & Time

- `YEAR` **[DONE]**: Returns the current year.
- `MONTH` **[DONE]**: Returns the current month.
- `DAY` **[DONE]**: Returns the current day.
- `HOUR` **[DONE]**: Returns the current hour.
- `MINUTE` **[DONE]**: Returns the current minute.
- `SECOND` **[DONE]**: Returns the current second.
- `MILLISECOND` **[DONE]**: Returns the current millisecond.
- `TIMESTAMP` **[DONE]**: Returns the number of seconds since the Unix epoch.
- `DATE` **[DONE]**: Returns the current date as a string.
- `TIME` **[DONE]**: Returns the current time as a string.
- `DATETIME` **[DONE]**: Returns the current date and time as a string.
- `TICKCOUNT` **[DONE]**: Returns the number of milliseconds since the program started.
- `TIMER` **[DONE]**: Returns the elapsed time in seconds since the program started.

---

## Bitwise Operations

- `BAND` **[DONE]**: Bitwise AND.
- `BOR` **[DONE]**: Bitwise OR.
- `BXOR` **[DONE]**: Bitwise XOR.
- `BNOT` **[DONE]**: Bitwise NOT.
- `BLSHIFT` **[DONE]**: Bitwise left shift.
- `BRSHIFT` **[DONE]**: Bitwise right shift.
- `BTEST` **[DONE]**: Tests a specific bit.
- `BSET` **[DONE]**: Sets a specific bit to 1.
- `BCLEAR` **[DONE]**: Clears a specific bit to 0.
- `BTOGGLE` **[DONE]**: Toggles a specific bit.
- `BCOUNT` **[DONE]**: Counts the number of set bits (1s).

---

## Audio

- `Audio.Init` **[DONE]**: Initializes the audio device.
- `Audio.Close` **[DONE]**: Closes the audio device.
- `Audio.LoadSound` **[DONE]**: Loads a sound effect.
- `Audio.LoadMusic` **[DONE]**: Loads streaming music.
- `Audio.Play` **[DONE]**: Plays a sound or music.
- `AudioStream.Make` **[DONE]**: Creates a raw audio stream.
- `AudioStream.Update` **[DONE]**: Updates a stream with PCM data.
- `Wave.Load` **[DONE]**: Loads a `.wav` file into memory.
- `Wave.Export` **[DONE]**: Saves a wave handle to a `.wav` file.
- `Sound.FromWave` **[DONE]**: Creates a playable sound from wave data.

---

## Program Control & Debugging

- `END` **[DONE]**: Terminates the program immediately.
- `QUIT` **[DONE]**: Terminates the program immediately.
- `STOP` **[DONE]**: Pauses the program and enters debug mode (if available).
- `WAIT` **[DONE]**: Pauses program execution for a specified number of milliseconds.
- `SLEEP` **[DONE]**: Alias for `WAIT`.
- `ASSERT` **[DONE]**: Asserts that a condition is true; if not, it halts with an error.
- `DUMP` **[PARTIAL]**: Coming soon.
- `TRACE` **[PARTIAL]**: Coming soon.
- `PRINTAT` **[PARTIAL]**: Coming soon.
- `PRINTCOLOR` **[PARTIAL]**: Coming soon.

---

## System & Host

- `ISFILEDROPPED` **[PARTIAL]**: Coming soon.
- `GETDROPPEDFILES` **[PARTIAL]**: Coming soon.
- `ENVIRON` **[DONE]**: Gets an environment variable.
- `COMMAND` **[DONE]**: Gets a command-line argument by index.
- `ARGC` **[DONE]**: Gets the number of command-line arguments.

---

## Module Commands

moonBASIC uses a dot-notation module system for its game engine commands. Tables below list **registry keys** (`WINDOW.*`, `RENDER.*`, …) first; **Easy Mode** (`Window.Open`, `Render.Clear`, …) compiles to the same keys — see [STYLE_GUIDE.md](../STYLE_GUIDE.md) and [EASY_MODE.md](EASY_MODE.md). They can also be called on a handle variable: `cam = CreateCamera()` (Easy Mode → `CAMERA.CREATE`) then `cam.SetPos(0, 5, 10)` — deprecated `Camera.Make()` / `CAMERA.MAKE` still compile with a warning.

> **Note:** Commands listed as `**[PARTIAL]**` or `**[MISSING]**` in this section are planned features that are not yet fully implemented.

---

### Window — [Reference](reference/WINDOW.md)

| Command | Description |
|---|---|
| `WINDOW.OPEN(w, h, title)` | Opens the window. |
| `WINDOW.CLOSE()` | Closes the window and exits. |
| `WINDOW.SHOULDCLOSE()` | Returns `TRUE` when user requests close. |
| `WINDOW.SETFPS(fps)` | Sets the target frame rate. |
| `WINDOW.SETTITLE(title)` | Updates the window title. |

---

### Render — [Reference](reference/RENDER.md)

| Command | Description |
|---|---|
| `RENDER.CLEAR(r, g, b, [a])` | Clears the screen. |
| `RENDER.FRAME()` | Presents the rendered frame. |
| `RENDER.DRAWFPS(x, y)` | Draws the current FPS. |
| `RENDER.WIDTH()` / `RENDER.HEIGHT()` | Returns framebuffer dimensions. |

---

### Camera — [Reference](reference/CAMERA.md)

| Command | Description |
|---|---|
| `CreateCamera()` / `CAMERA.CREATE()` | Creates a 3D camera handle (deprecated: `Camera.Make()` / `CAMERA.MAKE`). |
| `CAMERA.BEGIN(cam)` | Enters 3D mode. |
| `CAMERA.END()` | Exits 3D mode. |
| `CAMERA.SETPOS(cam, x, y, z)` | Sets camera position. |
| `CAMERA.SETTARGET(cam, x, y, z)` | Sets look-at point. |

---

### Entity — [Reference](reference/ENTITY.md)

| Command | Description |
|---|---|
| `ENTITY.LOAD(path)` | Loads a 3D model. |
| `ENTITY.CREATECUBE(size)` | Creates a cube entity. |
| `ENTITY.SETPOS(id, x, y, z)` | Sets world position. |
| `ENTITY.TURN(id, p, y, r)` | Adds rotation delta. |
| `ENTITY.FREE(id)` | Frees the entity. |

---

### Draw (2D) — [Reference](reference/DRAW2D.md)

| Command | Description |
|---|---|
| `DRAW.RECTANGLE(x, y, w, h, r, g, b, a)` | Filled rectangle. |
| `DRAW.RECTANGLE_ROUNDED(...)` | Rounded rectangle. |
| `DRAW.CIRCLE` / `DRAW.ELLIPSE` / `DRAW.RING` / `DRAW.TRIANGLE` / `DRAW.POLY` / `DRAW.OVAL` | Filled primitives (see reference for `*LINES` wire variants, `DRAW.ARC`, etc.). |
| `DRAW.LINE` / `DRAW.LINEEX` / `DRAW.LINEBEZIER` | Lines and Bézier strokes. |
| `DRAW.TEXTURE` / `DRAW.TEXTURENPATCH` / … | Textured quads, tiling, n-patch (tint is **required** on **`DRAW.TEXTURE`**: 7 args). |
| `DRAW.TEXT` / `DRAW.TEXTEX` / `DRAW.TEXTFONT` / … | Text and measurement helpers. |

### Draw (3D) — [Reference](reference/DRAW3D.md)

| Command | Description |
|---|---|
| `DRAW3D.GRID` / `Draw.Grid` | 3D reference grid (**`CAMERA.BEGIN`** / **`CAMERA.END`** or **`RENDER.BEGIN3D`** / **`RENDER.END3D`**). |
| `DRAW3D.LINE` / `DRAW3D.POINT` / `DRAW3D.SPHERE*` / `DRAW3D.CUBE*` / `DRAW3D.CYLINDER*` / `DRAW3D.CAPSULE*` / `DRAW3D.PLANE` / `DRAW3D.BBOX` | Primitives (see reference for arities). |
| `BOX` / `BOXW` / `WIRECUBE` / `BALL` / `BALLW` / `GRID3` / `FLAT` / `CAP` / `CAPW` | **Short global names** — same handlers as `DRAW3D.CUBE` / `CUBEWIRES` / … (`WIRECUBE` = Blitz **WireCube** — see [DRAW3D.md](reference/DRAW3D.md), [BLITZ3D.md](reference/BLITZ3D.md)). |
| `DRAW3D.RAY` | Debug-draw a ray from a 6-float array handle. |
| `DRAW3D.BILLBOARD` / `DRAW3D.BILLBOARDREC` | Textured billboards (require active 3D camera). |
| `DRAW3D.LINE3D` / `DRAW3D.SPHERE` / … | **`Draw.*`** aliases of the same `DRAW3D.*` handlers (see [DRAW3D.md](reference/DRAW3D.md)). |

---

### Texture — [Reference](reference/TEXTURE.md)

| Command | Description |
|---|---|
| `TEXTURE.LOAD(path)` | Loads a texture handle from disk. |
| `TEXTURE.FREE(id)` | Unloads a texture from memory. |
| `TEXTURE.FROMIMAGE(id)` | Creates a texture from an Image handle. |

---

### Image (CPU) — [Reference](reference/IMAGE.md)

| Command | Description |
|---|---|
| `IMAGE.LOAD(path)` | Loads a CPU pixel buffer from disk. |
| `IMAGE.MAKE(w, h)` | Creates a new blank Image handle (see also **`IMAGE.MAKEBLANK`**). |
| `IMAGE.FREE(id)` | Frees Image memory. |
| `IMAGE.EXPORT(id, path)` | Saves an Image to a file. |

---

### Font — [Reference](reference/FONT.md)

| Command | Description |
|---|---|
| `FONT.LOAD(path)` | Loads a `.ttf` or `.otf` font file. |
| `FONT.FREE(id)` | Unloads a font from memory. |

---

### GUI (raygui) — [Reference](reference/GUI.md)

| Command | Description |
|---|---|
| `GUI.BUTTON(...)` | Draws a clickable button (see [GUI.md](reference/GUI.md) for arity). |
| `GUI.LABEL(...)` | Draws a text label. |
| `GUI.SLIDER(...)` / `GUI.SLIDERBAR(...)` | Draws a slider control. |
| `GUI.SETFONT(id)` | Sets the active GUI font. |

Runnable demos: `examples/gui_basics/main.mb`, `examples/gui_form/main.mb`.

---

### Sprite & Atlas — [Sprite](reference/SPRITE.md) · [Atlas](reference/ATLAS.md)

| Command | Description |
|---|---|
| `SPRITE.LOAD(path)` | Loads a sprite from disk. |
| `SPRITE.DRAW(id, x, y)` | Draws a sprite at pixel coordinates. |
| `SPRITE.SETPOS(id, x, y)` | Sets a float draw offset. |
| `SPRITE.FREE(id)` | Frees a sprite handle. |
| `ATLAS.LOAD(path)` | Loads a texture atlas JSON. |
| `ATLAS.GETSPRITE(id, name)` | Retrieves a sprite from an atlas. |

---

### JSON, CSV & DB — [JSON](reference/JSON.md) · [CSV](reference/CSV.md) · [DB](reference/DATABASE.md)

| Command | Description |
|---|---|
| `JSON.LOADFILE(path)` / `JSON.PARSE(...)` | Load or parse JSON (see reference). |
| `JSON.GETSTRING(id, path)` | Reads a value from JSON. |
| `CSV.LOAD(path)` | Loads a CSV file (see [CSV.md](reference/CSV.md)). |
| `CSV.GET(id, row, col)` | Reads a cell from a CSV. |
| `DB.OPEN(path)` | Opens a SQLite database. |
| `DB.QUERY(db, sql [, ...params])` | Runs a SQL query (returns rows handle). |

---

### Model, Mesh & Material — [Reference](reference/MODEL.md) · [3D animation](reference/ANIMATION_3D.md)

| Command | Description |
|---|---|
| `MODEL.LOAD(path)` | Loads a 3D model file. |
| `MODEL.DRAW(handle)` | Draws a model using its root transform. |
| `MODEL.SETPOS(id, x, y, z)` | Sets model position. |
| `MODEL.FREE(handle)` | Unloads a model from memory. |
| `MESH.MAKECUBE(w, h, d)` | Creates a procedural box mesh. |
| `MESH.UPLOAD(id, dynamic)` | Uploads mesh data to GPU. |
| `MESH.DRAW(id, mat, matrix)` | Draws a single mesh. |
| `MESH.FREE(handle)` | Unloads a mesh from memory. |
| `MATERIAL.MAKEDEFAULT()` | Creates a default PBR material. |
| `MATERIAL.SETTEXTURE(id, slot, tex)` | Assigns a texture to a map slot. |
| `MATERIAL.FREE(handle)` | Frees a material. |

---

### Physics 3D (Jolt) — [Reference](reference/PHYSICS3D.md)

| Command | Description |
|---|---|
| `PHYSICS3D.START()` | Initializes the 3D physics world. |
| `PHYSICS3D.STOP()` | Shuts down the 3D physics world. |
| `PHYSICS3D.STEP()` | Advances simulation one step. |
| `BODY3D.CREATE(...)` | Creates a body definition (see reference). |
| `BODY3D.COMMIT(def, x, y, z)` | Finalizes body into the world. |

---

### Physics 2D (Box2D) — [Reference](reference/PHYSICS2D.md)

| Command | Description |
|---|---|
| `PHYSICS2D.START()` | Initializes the 2D physics world. |
| `PHYSICS2D.STOP()` | Shuts down the 2D physics world. |
| `PHYSICS2D.STEP()` | Advances the simulation. |
| `BODY2D.CREATE(...)` | Creates a 2D body definition (see reference). |
| `BODY2D.COMMIT(def, x, y)` | Finalizes body into the world. |

---

### Character Controller — [Reference](reference/CHARCONTROLLER.md)

| Command | Description |
|---|---|
| `CHARCONTROLLER.CREATE(r, h, x, y, z)` | Creates a capsule controller. |
| `CHARCONTROLLER.MOVE(id, dx, dy, dz)` | Moves with collision detection. |
| `CHARCONTROLLER.ISGROUNDED(id)` | Returns `TRUE` if on a surface. |
| `CHARCONTROLLER.FREE(id)` | Frees the controller. |

---

### Audio — [Reference](reference/AUDIO.md)

| Command | Description |
|---|---|
| `AUDIO.INIT()` | Initializes the audio device. |
| `AUDIO.LOADSOUND(path)` | Loads a sound effect. |
| `AUDIO.LOADMUSIC(path)` | Loads streaming music. |
| `AUDIO.PLAY(handle)` | Plays a sound or music track. |
| `AUDIO.UPDATEMUSIC(handle)` | Updates music buffer (call per frame). |
| `AUDIOSTREAM.UPDATE(handle, pcmArray)` | Pushes PCM data to the stream. |
| `AUDIOSTREAM.PLAY(handle)` | Starts the audio stream. |
| `AUDIOSTREAM.ISPLAYING(handle)` | Returns `TRUE` if the stream is playing. |
| `AUDIOSTREAM.FREE(handle)` | Frees the audio stream. |

---

### Network (ENet) — [Reference](reference/NETWORK.md) · [Command set (host/events)](reference/moonbasic-command-set/network-enet.md) · [Helpers](reference/moonbasic-command-set/network-helpers.md)

| Command | Description |
|---|---|
| `NET.START()` | Initializes the networking system. |
| `NET.STOP()` | Shuts down the networking system. |
| `NET.CREATESERVER(port, maxClients)` | Creates a server host. Returns a handle. |
| `NET.CREATECLIENT()` | Creates a client host. Returns a handle. |
| `NET.CONNECT(clientHandle, address, port)` | Connects a client to a server. Returns a peer handle. |
| `NET.UPDATE(hostHandle)` | Processes network packets. **Call every frame.** |
| `NET.RECEIVE(hostHandle)` | Returns the next event handle, or `0` if none queued. |
| `NET.BROADCAST(serverHandle, channel, data, reliable)` | Sends a message to all connected clients. |
| `PEER.SEND(peerHandle, channel, data, reliable)` | Sends a message to a specific peer. |
| `EVENT.TYPE(handle)` | Returns the event type (`EVENT_CONNECT`, `EVENT_DISCONNECT`, `EVENT_RECEIVE`). |
| `EVENT.PEER(handle)` | Returns the peer associated with the event. |
| `EVENT.DATA(handle)` | Returns the string data of a `RECEIVE` event. |
| `EVENT.FREE(handle)` | Frees the event. **Must be called for every event received.** |

---

### Time — [Reference](reference/TIME.md)

| Command | Description |
|---|---|
| `TIME.DELTA()` | Returns seconds elapsed since last frame. |
| `TIME.GET()` | Returns total seconds elapsed since start. |

---

### Input — [Reference](reference/INPUT.md)

| Command | Description |
|---|---|
| `INPUT.KEYDOWN(key)` | True while a key is held. |
| `INPUT.KEYPRESSED(key)` | True on the first frame of a press. |
| `INPUT.MOUSEX()` / `INPUT.MOUSEY()` | Current mouse pixel coordinates. |
| `INPUT.AXIS(neg, pos)` | Returns -1, 0, or 1 based on keys. |

---

### Transform — [Reference](reference/TRANSFORM.md)

| Command | Description |
|---|---|
| `TRANSFORM.IDENTITY()` | New identity matrix handle. |
| `TRANSFORM.TRANSLATION(x, y, z)` | Translation matrix. |
| `TRANSFORM.ROTATION(p, y, r)` | Euler rotation matrix (radians). |
| `TRANSFORM.MULTIPLY(a, b)` | Combines two matrices. |
| `TRANSFORM.FREE(id)` | Frees matrix handle. |

---

### Shader — [Reference](reference/SHADER.md)

| Command | Description |
|---|---|
| `SHADER.LOAD(vs, fs)` | Loads GLSL vertex/fragment shaders. |
| `SHADER.FREE(id)` | Unloads shader from GPU. |

---

### Light — [Reference](reference/LIGHT.md)

| Command | Description |
|---|---|
| `LIGHT.CREATE*` | Creates a light handle (point / directional / spot — see reference). |
| `LIGHT.SETDIR(id, x, y, z)` | Sets directional light vector. |
| `LIGHT.FREE(id)` | Frees light resource. |

---

### Tilemap — [Reference](reference/TILEMAP.md)

| Command | Description |
|---|---|
| `TILEMAP.LOAD(path)` | Loads a Tiled map handle. |
| `TILEMAP.DRAW(id, x, y)` | Draws map layers (see reference for overloads). |
| `TILEMAP.ISSOLID(id, x, y)` / `TILEMAP.COLLISIONAT(...)` | Collision checks (see reference). |
| `TILEMAP.FREE(id)` | Frees tilemap. |

---

### Particles — [Reference](reference/PARTICLES.md)

| Command | Description |
|---|---|
| `PARTICLE.CREATE()` | Creates particle emitter. |
| `PARTICLE.SETEMITRATE` / `PARTICLE.SETLIFETIME` / `PARTICLE.SETVELOCITY` / … | Configure emitter (see reference). |
| `PARTICLE.SETPOS(handle, x, y, z)` | Emitter position in world space. |
| `PARTICLE.PLAY(handle)` / `PARTICLE.STOP(handle)` | Start or stop simulation. |
| `PARTICLE.UPDATE(handle, dt)` | Advances simulation — call every frame. |
| `PARTICLE.DRAW(handle)` | Draws active particles. |
| `PARTICLE.FREE(handle)` | Frees the emitter. |

---

### Open world — terrain, streaming, water, sky, weather ([TERRAIN](reference/TERRAIN.md), [WORLD](reference/WORLD.md), [WATER](reference/WATER.md), [SKY](reference/SKY.md), [CLOUD](reference/CLOUD.md), [WEATHER](reference/WEATHER.md), [BIOME](reference/BIOME.md), [SCATTER](reference/SCATTER.md), [NAVMESH](reference/NAVMESH.md))

| Command | Description |
|---|---|
| `TERRAIN.CREATE` / `TERRAIN.FREE` | Create or free heightfield terrain (**`TERRAIN.MAKE`** is a deprecated alias). |
| `TERRAIN.SETPOS` / `TERRAIN.SETCHUNKSIZE` | World origin and chunk sample size. |
| `TERRAIN.FILLPERLIN` / `TERRAIN.FILLFLAT` | Procedural or flat height fill. |
| `TERRAIN.GETHEIGHT` / `TERRAIN.GETSLOPE` | Sample height and slope at XZ. |
| `TERRAIN.RAISE` / `TERRAIN.LOWER` | Brush sculpting. |
| `TERRAIN.DRAW` | Draw loaded chunk meshes. |
| `CHUNK.GENERATE` / `CHUNK.COUNT` / `CHUNK.SETRANGE` / `CHUNK.ISLOADED` | Chunk mesh build and streaming distances. |
| `WORLD.SETCENTER` / `WORLD.UPDATE` / `WORLD.STREAMENABLE` | Streaming focal point and per-frame tick. |
| `WORLD.PRELOAD` / `WORLD.STATUS` / `WORLD.ISREADY` | Preload radius and readiness. |
| `WATER.CREATE` / `WATER.FREE` / `WATER.SETPOS` | Water plane mesh. |
| `WATER.DRAW` / `WATER.UPDATE` / `WATER.SETWAVEHEIGHT` | Render and animate waves. |
| `WATER.GETWAVEY` / `WATER.GETDEPTH` / `WATER.ISUNDER` | Surface and depth queries. |
| `WATER.SETSHALLOWCOLOR` / `WATER.SETDEEPCOLOR` | Water color tuning. |
| `SKY.CREATE` / `SKY.FREE` / `SKY.UPDATE` / `SKY.DRAW` | Day/night sky dome. |
| `SKY.SETTIME` / `SKY.SETDAYLENGTH` / `SKY.GETTIMEHOURS` / `SKY.ISNIGHT` | Time-of-day. |
| `CLOUD.CREATE` / `CLOUD.FREE` / `CLOUD.UPDATE` / `CLOUD.DRAW` / `CLOUD.SETCOVERAGE` | Cloud layer state. |
| `WEATHER.CREATE` / `WEATHER.FREE` / `WEATHER.UPDATE` / `WEATHER.DRAW` / `WEATHER.SETTYPE` / `WEATHER.GETCOVERAGE` / `WEATHER.GETTYPE` | Weather controller. |
| `FOG.ENABLE` / `FOG.SETNEAR` / `FOG.SETFAR` / `FOG.SETCOLOR` | Fog distances and color. |
| `WIND.SET` / `WIND.GETSTRENGTH` | Wind vector and strength. |
| `BIOME.CREATE` / `BIOME.FREE` / `BIOME.SETTEMP` / `BIOME.SETHUMIDITY` | Biome parameters. |
