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

moonBASIC uses a dot-notation module system for its game engine commands. These all use the form `Module.Command(args)`. They can also be called on a handle variable: `cam = Camera.Make()` then `cam.SetPos(0, 5, 10)`.

> **Note:** Commands listed as `**[PARTIAL]**` or `**[MISSING]**` in this section are planned features that are not yet fully implemented.

---

### Window — [Reference](reference/WINDOW.md)

| Command | Description |
|---|---|
| `Window.Open(w, h, title)` | Opens the window. |
| `Window.Close()` | Closes the window and exits. |
| `Window.ShouldClose()` | Returns `TRUE` when user requests close. |
| `Window.SetFPS(fps)` | Sets the target frame rate. |
| `Window.SetTitle(title)` | Updates the window title. |

---

### Render — [Reference](reference/RENDER.md)

| Command | Description |
|---|---|
| `Render.Clear(r, g, b, [a])` | Clears the screen. |
| `Render.Frame()` | Presents the rendered frame. |
| `Render.DrawFPS(x, y)` | Draws the current FPS. |
| `Render.Width()` / `Render.Height()` | Returns framebuffer dimensions. |

---

### Camera — [Reference](reference/CAMERA.md)

| Command | Description |
|---|---|
| `Camera.Make()` | Creates a 3D camera handle. |
| `Camera.Begin(cam)` | Enters 3D mode. |
| `Camera.End()` | Exits 3D mode. |
| `Camera.SetPos(cam, x, y, z)` | Sets camera position. |
| `Camera.SetTarget(cam, x, y, z)` | Sets look-at point. |

---

### Entity — [Reference](reference/ENTITY.md)

| Command | Description |
|---|---|
| `Entity.Load(path)` | Loads a 3D model. |
| `Entity.CreateCube(size)` | Creates a cube entity. |
| `Entity.Position(id, x, y, z)` | Sets world position. |
| `Entity.Turn(id, p, y, r)` | Adds rotation delta. |
| `Entity.Free(id)` | Frees the entity. |

---

### Draw (2D) — [Reference](reference/DRAW2D.md)

| Command | Description |
|---|---|
| `Draw.Rectangle(x, y, w, h, r, g, b, a)` | Filled rectangle. |
| `Draw.RectangleRounded(...)` | Rounded rectangle. |
| `Draw.Circle` / `Draw.Ellipse` / `Draw.Ring` / `Draw.Triangle` / `Draw.Poly` | Filled primitives (see reference for wire variants, gradients, arcs). |
| `Draw.CircleSector` / `Draw.CircleGradient` | Sector and radial gradient. |
| `Draw.Line` / `Draw.LineEx` / `Draw.LineBezier*` | Lines and Bézier strokes. |
| `Draw.Spline*` | Splines from a point array + thickness + color. |
| `Draw.Texture*` / `Draw.TextureNPatch` | Textured quads, tiling, n-patch (tint is **required** on `Draw.Texture`: 7 args). |
| `Draw.Text` / `Draw.TextEx` / `Draw.TextFont` / `Draw.TextPro` | Text and measurement helpers. |

### Draw (3D) — [Reference](reference/DRAW3D.md)

| Command | Description |
|---|---|
| `Draw3D.Grid` / `Draw.Grid` | 3D reference grid (`Camera.Begin` / `End`). |
| `Draw3D.Line` / `Draw3D.Point` / `Draw3D.Sphere*` / `Draw3D.Cube*` / `Draw3D.Cylinder*` / `Draw3D.Capsule*` / `Draw3D.Plane` / `Draw3D.BBox` | Primitives (see reference for arities). |
| `BOX` / `BOXW` / `WIRECUBE` / `BALL` / `BALLW` / `GRID3` / `FLAT` / `CAP` / `CAPW` | **Short global names** — same handlers as `DRAW3D.CUBE` / `CUBEWIRES` / … (`WIRECUBE` = Blitz **WireCube** — see [DRAW3D.md](reference/DRAW3D.md), [BLITZ3D.md](reference/BLITZ3D.md)). |
| `Draw3D.Ray` | Debug-draw a ray from a 6-float array handle. |
| `Draw3D.Billboard` / `Draw3D.BillboardRec` | Textured billboards (require active 3D camera). |
| `Draw.Line3D` / `Draw.Sphere` / … | **`Draw.*`** aliases of the same `DRAW3D.*` handlers (see [DRAW3D.md](reference/DRAW3D.md)). |

---

### Texture — [Reference](reference/TEXTURE.md)

| Command | Description |
|---|---|
| `Texture.Load(path)` | Loads a texture handle from disk. |
| `Texture.Free(id)` | Unloads a texture from memory. |
| `Texture.FromImage(id)` | Creates a texture from an Image handle. |

---

### Image (CPU) — [Reference](reference/IMAGE.md)

| Command | Description |
|---|---|
| `Image.Load(path)` | Loads a CPU pixel buffer from disk. |
| `Image.Make(w, h)` | Creates a new blank Image handle. |
| `Image.Free(id)` | Frees Image memory. |
| `Image.Export(id, path)` | Saves an Image to a file. |

---

### Font — [Reference](reference/FONT.md)

| Command | Description |
|---|---|
| `Font.Load(path)` | Loads a `.ttf` or `.otf` font file. |
| `Font.Free(id)` | Unloads a font from memory. |

---

### GUI (raygui) — [Reference](reference/GUI.md)

| Command | Description |
|---|---|
| `Gui.Button(label, x, y, w, h)` | Draws a clickable button. |
| `Gui.Label(text, x, y)` | Draws a text label. |
| `Gui.Slider(label, x, y, val, min, max)` | Draws a slider. |
| `Gui.SetFont(id)` | Sets the active GUI font. |

Runnable demos: `examples/gui_basics/main.mb`, `examples/gui_form/main.mb`.

---

### Sprite & Atlas — [Sprite](reference/SPRITE.md) · [Atlas](reference/ATLAS.md)

| Command | Description |
|---|---|
| `Sprite.Load(path)` | Loads a sprite from disk. |
| `Sprite.Draw(id, x, y)` | Draws a sprite at pixel coordinates. |
| `Sprite.SetPos(id, x, y)` | Sets a float draw offset. |
| `Sprite.Free(id)` | Frees a sprite handle. |
| `Atlas.Load(path)` | Loads a texture atlas JSON. |
| `Atlas.GetSprite(id, name)` | Retrieves a sprite from an atlas. |

---

### JSON, CSV & DB — [JSON](reference/JSON.md) · [CSV](reference/CSV.md) · [DB](reference/DATABASE.md)

| Command | Description |
|---|---|
| `JSON.Parse(path)` | Loads a JSON file. |
| `JSON.GetString(id, path)` | Reads a value from JSON. |
| `CSV.Load(path)` | Loads a CSV file. |
| `CSV.Get(id, row, col)` | Reads a cell from a CSV. |
| `DB.Open(path)` | Opens a SQLite database. |
| `DB.Query(id, sql)` | Runs a SQL query. |

---

### Model, Mesh & Material — [Reference](reference/MODEL.md) · [3D animation](reference/ANIMATION_3D.md)

| Command | Description |
|---|---|
| `Model.Load(path)` | Loads a 3D model file. |
| `Model.Draw(handle)` | Draws a model using its root transform. |
| `Model.SetPos(id, x, y, z)` | Sets model position. |
| `Model.Free(handle)` | Unloads a model from memory. |
| `Mesh.MakeCube(w, h, d)` | Creates a procedural box mesh. |
| `Mesh.Upload(id, dynamic)` | Uploads mesh data to GPU. |
| `Mesh.Draw(id, mat, matrix)` | Draws a single mesh. |
| `Mesh.Free(handle)` | Unloads a mesh from memory. |
| `Material.MakeDefault()` | Creates a default PBR material. |
| `Material.SetTexture(id, slot, tex)` | Assigns a texture to a map slot. |
| `Material.Free(handle)` | Frees a material. |

---

### Physics 3D (Jolt) — [Reference](reference/PHYSICS3D.md)

| Command | Description |
|---|---|
| `Physics3D.Start()` | Initializes the 3D physics world. |
| `Physics3D.Stop()` | Shuts down the 3D physics world. |
| `Physics3D.Step()` | Advances simulation one step. |
| `Body3D.Make(type)` | Creates a body definition. |
| `Body3D.Commit(def, x, y, z)` | Finalizes body into the world. |

---

### Physics 2D (Box2D) — [Reference](reference/PHYSICS2D.md)

| Command | Description |
|---|---|
| `Physics2D.Start()` | Initializes the 2D physics world. |
| `Physics2D.Stop()` | Shuts down the 2D physics world. |
| `Physics2D.Step()` | Advances the simulation. |
| `Body2D.Make(type)` | Creates a 2D body definition. |
| `Body2D.Commit(def, x, y)` | Finalizes body into the world. |

---

### Character Controller — [Reference](reference/CHARCONTROLLER.md)

| Command | Description |
|---|---|
| `CharController.Make(r, h, x, y, z)` | Creates a capsule controller. |
| `CharController.Move(id, dx, dy, dz)` | Moves with collision detection. |
| `CharController.IsGrounded(id)` | Returns `TRUE` if on a surface. |
| `CharController.Free(id)` | Frees the controller. |

---

### Audio — [Reference](reference/AUDIO.md)

| Command | Description |
|---|---|
| `Audio.Init()` | Initializes the audio device. |
| `Audio.LoadSound(path)` | Loads a sound effect. |
| `Audio.LoadMusic(path)` | Loads streaming music. |
| `Audio.Play(handle)` | Plays a sound or music track. |
| `Audio.UpdateMusic(handle)` | Updates music buffer (call per frame). |
| `AudioStream.Update(handle, pcmArray)` | Pushes PCM data to the stream. |
| `AudioStream.Play(handle)` | Starts the audio stream. |
| `AudioStream.IsPlaying(handle)` | Returns `TRUE` if the stream is playing. |
| `AudioStream.Free(handle)` | Frees the audio stream. |

---

### Network (ENet) — [Reference](reference/NETWORK.md) · [Command set (host/events)](reference/moonbasic-command-set/network-enet.md) · [Helpers](reference/moonbasic-command-set/network-helpers.md)

| Command | Description |
|---|---|
| `Net.Start()` | Initializes the networking system. |
| `Net.Stop()` | Shuts down the networking system. |
| `Net.CreateServer(port, maxClients)` | Creates a server host. Returns a handle. |
| `Net.CreateClient()` | Creates a client host. Returns a handle. |
| `Net.Connect(clientHandle, address, port)` | Connects a client to a server. Returns a peer handle. |
| `Net.Update(hostHandle)` | Processes network packets. **Call every frame.** |
| `Net.Receive(hostHandle)` | Returns the next event handle, or `0` if none queued. |
| `Net.Broadcast(serverHandle, channel, data, reliable)` | Sends a message to all connected clients. |
| `Peer.Send(peerHandle, channel, data, reliable)` | Sends a message to a specific peer. |
| `Event.Type(handle)` | Returns the event type (`EVENT_CONNECT`, `EVENT_DISCONNECT`, `EVENT_RECEIVE`). |
| `Event.Peer(handle)` | Returns the peer associated with the event. |
| `Event.Data(handle)` | Returns the string data of a `RECEIVE` event. |
| `Event.Free(handle)` | Frees the event. **Must be called for every event received.** |

---

### Time — [Reference](reference/TIME.md)

| Command | Description |
|---|---|
| `Time.Delta()` | Returns seconds elapsed since last frame. |
| `Time.Get()` | Returns total seconds elapsed since start. |

---

### Input — [Reference](reference/INPUT.md)

| Command | Description |
|---|---|
| `Input.KeyDown(key)` | True while a key is held. |
| `Input.KeyPressed(key)` | True on the first frame of a press. |
| `Input.MouseX()` / `Input.MouseY()` | Current mouse pixel coordinates. |
| `Input.Axis(neg, pos)` | Returns -1, 0, or 1 based on keys. |

---

### Transform — [Reference](reference/TRANSFORM.md)

| Command | Description |
|---|---|
| `Transform.Identity()` | New identity matrix handle. |
| `Transform.Translate(x, y, z)` | Translation matrix. |
| `Transform.Rotate(p, y, r)` | Euler rotation matrix (radians). |
| `Transform.Multiply(a, b)` | Combines two matrices. |
| `Transform.Free(id)` | Frees matrix handle. |

---

### Shader — [Reference](reference/SHADER.md)

| Command | Description |
|---|---|
| `Shader.Load(vs, fs)` | Loads GLSL vertex/fragment shaders. |
| `Shader.Free(id)` | Unloads shader from GPU. |

---

### Light — [Reference](reference/LIGHT.md)

| Command | Description |
|---|---|
| `Light.Make(type)` | Creates a light handle. |
| `Light.SetDir(id, x, y, z)` | Sets directional light vector. |
| `Light.Free(id)` | Frees light resource. |

---

### Tilemap — [Reference](reference/TILEMAP.md)

| Command | Description |
|---|---|
| `Tilemap.Load(path)` | Loads a Tiled map handle. |
| `Tilemap.Draw(id, x, y)` | Draws all map layers. |
| `Tilemap.IsSolid(id, x, y)` | Collision check for tile. |
| `Tilemap.Free(id)` | Frees tilemap. |

---

### Particles — [Reference](reference/PARTICLES.md)

| Command | Description |
|---|---|
| `Particles.Make(max)` | Creates particle emitter. |
| `Particles.Emit(id, x, y, z, n)` | Emits a burst of particles. |
| `Particles.Free(id)` | Frees emitter. |
| `Particle.SetEmitRate(handle, rate)` | Sets particles emitted per second. |
| `Particle.SetLifetime(handle, min, max)` | Sets the particle lifetime range in seconds. |
| `Particle.SetVelocity(handle, vx, vy, vz, spread)` | Sets the emission velocity and spread. |
| `Particle.SetColor(handle, r, g, b, a)` | Sets the starting particle color. |
| `Particle.SetColorEnd(handle, r, g, b, a)` | Sets the ending particle color (fades to this). |
| `Particle.SetSize(handle, startSize, endSize)` | Sets the particle size over its lifetime. |
| `Particle.SetGravity(handle, gx, gy, gz)` | Sets per-emitter gravity. |
| `Particle.SetPos(handle, x, y, z)` | Sets the emitter position in world space. |
| `Particle.Play(handle)` | Starts the emitter. |
| `Particle.Update(handle, dt)` | Advances particle simulation. Call every frame. |
| `Particle.Draw(handle)` | Draws all active particles. |
| `Particle.Free(handle)` | Frees the particle emitter. |

---

### Open world — terrain, streaming, water, sky, weather ([TERRAIN](reference/TERRAIN.md), [WORLD](reference/WORLD.md), [WATER](reference/WATER.md), [SKY](reference/SKY.md), [CLOUD](reference/CLOUD.md), [WEATHER](reference/WEATHER.md), [SCATTER](reference/SCATTER.md), [BIOME](reference/BIOME.md), [NAVMESH](reference/NAVMESH.md))

| Command | Description |
|---|---|
| `Terrain.Make` / `Terrain.Free` | Create or free heightfield terrain (`TERRAIN.MAKE` overloads: 2 or 3 args). |
| `Terrain.SetPos` / `SetChunkSize` | World origin and chunk sample size. |
| `Terrain.FillPerlin` / `FillFlat` | Procedural or flat height fill. |
| `Terrain.GetHeight` / `GetSlope` | Sample height and slope at XZ. |
| `Terrain.Raise` / `Lower` | Brush sculpting. |
| `Terrain.Draw` | Draw loaded chunk meshes. |
| `Chunk.Generate` / `Count` / `SetRange` / `IsLoaded` | Chunk mesh build and streaming distances. |
| `World.SetCenter` / `Update` / `StreamEnable` | Streaming focal point and per-frame tick. |
| `World.Preload` / `Status` / `IsReady` | Preload radius and debug strings. |
| `Water.Make` / `Free` / `SetPos` | Water plane mesh. |
| `Water.Draw` / `Update` / `SetWaveHeight` | Render and animate waves. |
| `Water.GetWaveY` / `GetDepth` / `IsUnder` | Surface and depth queries. |
| `Water.SetShallowColor` / `SetDeepColor` | Water color tuning. |
| `Sky.Make` / `Free` / `Update` / `Draw` | Day/night sky dome. |
| `Sky.SetTime` / `SetDayLength` / `GetTimeHours` / `IsNight` | Time-of-day. |
| `Cloud.Make` / `Free` / `Update` / `Draw` / `SetCoverage` | Cloud layer state. |
| `Weather.Make` / `Free` / `Update` / `Draw` / `SetType` / `GetCoverage` / `GetType` | Weather controller. |
| `Fog.Enable` / `SetNear` / `SetFar` / `SetColor` | Fog distances and color. |
| `Wind.Set` / `GetStrength` | Wind vector and strength. |
| `Scatter.Create` / `Free` / `Apply` / `DrawAll` | Scatter instances on terrain. |
| `Prop.Place` / `Free` / `DrawAll` | Placed prop markers. |
| `Biome.Make` / `Free` / `SetTemp` / `SetHumidity` | Biome parameters. |
