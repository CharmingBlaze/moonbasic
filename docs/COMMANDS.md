# moonBASIC Command Index

This page lists every command available in moonBASIC, grouped by category.

**How to use these commands in a real program:** see the [Programming Guide](PROGRAMMING.md) (game loop, types, CGO) and runnable projects under [`examples/`](../examples/README.md). Copy-paste snippets also live in [Examples](EXAMPLES.md).

**Legend:**
- **[DONE]**: Implemented, tested, and ready to use.
- **[PARTIAL]**: Partially implemented or missing tests. May not work as expected.
- **[MISSING]**: Not yet implemented.

---

## Output / Input

- `PRINT` **[DONE]**: Print values to stdout, space-separated, with a newline.
- `PRINTLN` **[DONE]**: Same as `PRINT`.
- `WRITE` **[DONE]**: Prints values without a trailing newline.
- `INPUT` **[DONE]**: Prompts the user for console input.
- `CLS` **[DONE]**: Clears the console screen.
- `LOCATE` **[DONE]**: Moves the console cursor to a specific row and column.
- `TAB` **[PARTIAL]**: Coming soon.
- `SPC` **[PARTIAL]**: Coming soon.

---

## Type Conversion

- `INT` **[DONE]**: Converts a value to an integer.
- `FLOAT` **[DONE]**: Converts a value to a float.
- `STR$` **[DONE]**: Converts a value to a string.
- `VAL` **[DONE]**: Parses a string to a float number.
- `ASC` **[DONE]**: Returns the ASCII code for the first character of a string.
- `CHR$` **[DONE]**: Returns a string from an ASCII code.
- `BOOL` **[PARTIAL]**: Coming soon.
- `FIX` **[DONE]**: Truncates a float toward zero (e.g. `FIX(-3.7)` → `-3`).
- `TYPEOF` **[DONE]**: Returns the type of a variable as a string.
- `ISNULL` **[DONE]**: Checks if a value is null.
- `ISHANDLE` **[DONE]**: Checks if a value is a handle.
- `ISTYPE` **[DONE]**: Checks if a variable is of a certain type.

---

## String Manipulation

- `LEN` **[DONE]**: Returns the length of a string.
- `LEFT$` **[DONE]**: Returns a specified number of characters from the left side of a string.
- `RIGHT$` **[DONE]**: Returns a specified number of characters from the right side of a string.
- `MID$` **[DONE]**: Extracts a substring from a string.
- `UPPER$` **[DONE]**: Converts a string to uppercase.
- `LOWER$` **[DONE]**: Converts a string to lowercase.
- `INSTR` **[DONE]**: Finds the position of the first occurrence of a substring in a string.
- `TRIM$` **[DONE]**: Removes whitespace from both ends of a string.
- `LTRIM$` **[DONE]**: Removes whitespace from the left side of a string.
- `RTRIM$` **[DONE]**: Removes whitespace from the right side of a string.
- `REPLACE$` **[DONE]**: Replaces all occurrences of a substring with another substring.
- `CONTAINS` **[DONE]**: Checks if a string contains a specific substring.
- `STARTSWITH` **[DONE]**: Checks if a string starts with a specific substring.
- `ENDSWITH` **[DONE]**: Checks if a string ends with a specific substring.
- `SPLIT$` **[DONE]**: Splits a string into an array of substrings.
- `JOIN$` **[DONE]**: Joins the elements of a string array into a single string.
- `SPACE$` **[DONE]**: Returns a string consisting of a specified number of spaces.
- `STRING$` **[DONE]**: Returns a string of a specified length, consisting of a repeated character.
- `REVERSE$` **[DONE]**: Reverses a string.
- `REPEAT$` **[DONE]**: Repeats a string a specified number of times.
- `COUNT$` **[DONE]**: Counts the number of occurrences of a substring.
- `LSET$` **[DONE]**: Left-aligns a string within a specified length.
- `RSET$` **[DONE]**: Right-aligns a string within a specified length.
- `ISALPHA` **[DONE]**: Checks if a string contains only alphabetic characters.
- `ISNUMERIC` **[PARTIAL]**: Coming soon.
- `ISALPHANUM` **[DONE]**: Checks if a string contains only alphanumeric characters.
- `FORMAT$` **[DONE]**: Formats a number into a string.
- `BIN$` **[DONE]**: Converts an integer to a binary string.
- `HEX$` **[DONE]**: Converts an integer to a hexadecimal string.
- `OCT$` **[DONE]**: Converts an integer to an octal string.
- `MKSHORT$` **[DONE]**: Converts a short integer to a 2-byte string.
- `MKINT$` **[DONE]**: Converts an integer to a 4-byte string.
- `MKLONG$` **[DONE]**: Converts a long integer to an 8-byte string.
- `MKFLOAT$` **[DONE]**: Converts a float to a 4-byte string.
- `MKDOUBLE$` **[DONE]**: Converts a double to an 8-byte string.
- `CVSHORT` **[DONE]**: Converts a 2-byte string to a short integer.
- `CVINT` **[DONE]**: Converts a 4-byte string to an integer.
- `CVLONG` **[DONE]**: Converts an 8-byte string to a long integer.
- `CVFLOAT` **[DONE]**: Converts a 4-byte string to a float.
- `CVDOUBLE` **[DONE]**: Converts an 8-byte string to a double.

---

## Math

All math functions are also available as `MATH.name(...)` (e.g. `MATH.SIN`, `MATH.CLAMP`). See [Math Reference](reference/MATH.md) for full details.

- `SIN(angle#)` **[DONE]**: Returns the sine of an angle (radians).
- `COS(angle#)` **[DONE]**: Returns the cosine of an angle (radians).
- `TAN(angle#)` **[DONE]**: Returns the tangent of an angle (radians).
- `ATN(angle#)` **[DONE]**: Returns the arctangent of a value (radians). Also available as `ATAN`.
- `ASIN(value#)` **[DONE]**: Returns the arcsine (radians).
- `ACOS(value#)` **[DONE]**: Returns the arccosine (radians).
- `ATAN(value#)` **[DONE]**: Returns the arctangent (radians). Alias for `ATN`.
- `ATAN2(y#, x#)` **[DONE]**: Returns the two-argument arctangent (radians).
- `SQR(value#)` **[DONE]**: Returns the square root. Alias for `SQRT`.
- `SQRT(value#)` **[DONE]**: Returns the square root.
- `EXP(value#)` **[DONE]**: Returns e raised to the power of `value`.
- `LOG(value#)` **[DONE]**: Returns the natural logarithm.
- `LOG2(value#)` **[DONE]**: Returns the base-2 logarithm.
- `LOG10(value#)` **[DONE]**: Returns the base-10 logarithm.
- `POW(base#, exp#)` **[DONE]**: Returns `base` raised to the power of `exp`.
- `FLOOR(value#)` **[DONE]**: Rounds down to the nearest integer.
- `CEIL(value#)` **[DONE]**: Rounds up to the nearest integer.
- `ROUND(value#, [decimals])` **[DONE]**: Rounds to nearest integer or to `decimals` decimal places.
- `ABS(value#)` **[DONE]**: Returns the absolute value.
- `SGN(value#)` **[DONE]**: Returns -1, 0, or 1 depending on the sign.
- `FIX(value#)` **[DONE]**: Truncates toward zero (like `INT` but for floats).
- `MIN(a#, b#)` **[DONE]**: Returns the smaller of two values.
- `MAX(a#, b#)` **[DONE]**: Returns the larger of two values.
- `CLAMP(value#, min#, max#)` **[DONE]**: Clamps a value between min and max.
- `LERP(a#, b#, t#)` **[DONE]**: Linearly interpolates between `a` and `b` by `t` (0.0-1.0).
- `SMOOTHSTEP(lo#, hi#, x#)` **[DONE]**: Smooth interpolation between 0 and 1.
- `PINGPONG(t#, length#)` **[DONE]**: Bounces `t` back and forth between 0 and `length`.
- `WRAP(value#, min#, max#)` **[DONE]**: Wraps `value` within the range [min, max].
- `DEG2RAD(degrees#)` **[DONE]**: Converts degrees to radians.
- `RAD2DEG(radians#)` **[DONE]**: Converts radians to degrees.
- `WRAPANGLE(angle#)` **[DONE]**: Wraps an angle to the range [0, 360).
- `WRAPANGLE180(angle#)` **[DONE]**: Wraps an angle to the range [-180, 180).
- `ANGLEDIFF(from#, to#)` **[DONE]**: Returns the shortest signed angle difference in degrees.
- `RND([limit])` **[DONE]**: Returns a random integer from 0 to `limit-1`, or a random float [0,1) if no argument.
- `RNDF(min#, max#)` **[DONE]**: Returns a random float between `min` and `max`.
- `RNDSEED(seed)` **[DONE]**: Seeds the random number generator.
- `RANDOMIZE([seed])` **[DONE]**: Seeds the RNG from a value, or from the clock if omitted.
- `PI()` **[DONE]**: Returns π (3.14159...).
- `TAU()` **[DONE]**: Returns τ = 2π.
- `E()` **[DONE]**: Returns Euler's number (2.71828...).

---

## Logic

- `IIF(condition, trueVal, falseVal)` **[DONE]**: Returns `trueVal` if condition is true, otherwise `falseVal`. Inline if-then-else.
- `CHOOSE(index, val1, val2, ...)` **[DONE]**: Returns the value at position `index` (1-based) from the argument list.
- `SWITCH(expr, case1, val1, ..., caseN, valN, default)` **[DONE]**: Returns the value paired with the first matching case, or `default` if none match.

---

## Array Operations

- `DIM` **[DONE]**: Declares an array.
- `REDIM` **[PARTIAL]**: Coming soon.
- `ERASE` **[PARTIAL]**: Coming soon.
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
- `ARRAYJOINS$` **[PARTIAL]**: Coming soon.

---

## File I/O

- `OPENFILE` **[PARTIAL]**: Coming soon.
- `CLOSEFILE` **[PARTIAL]**: Coming soon.
- `READFILE$` **[PARTIAL]**: Coming soon.
- `FILE.WRITE` **[DONE]**: Writes raw bytes to a file without a newline.
- `FILE.WRITELN` **[DONE]**: Writes a string to a file, followed by a newline.
- `WRITEFILE` **[DONE]**: Alias for `FILE.WRITE`.
- `WRITEFILELN` **[DONE]**: Alias for `FILE.WRITELN`.
- `READALLTEXT$` **[DONE]**: Reads the entire content of a file into a string.
- `WRITEALLTEXT` **[DONE]**: Writes a string to a file, overwriting existing content.
- `READBYTE` **[PARTIAL]**: Coming soon.
- `WRITEBYTE` **[PARTIAL]**: Coming soon.
- `READSHORT` **[PARTIAL]**: Coming soon.
- `WRITESHORT` **[PARTIAL]**: Coming soon.
- `READINT` **[PARTIAL]**: Coming soon.
- `WRITEINT` **[PARTIAL]**: Coming soon.
- `READFLOAT` **[PARTIAL]**: Coming soon.
- `WRITEFLOAT` **[PARTIAL]**: Coming soon.
- `READSTRING$` **[PARTIAL]**: Coming soon.
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
- `GETDIR$` **[DONE]**: Gets the current working directory.
- `SETDIR` **[DONE]**: Sets the current working directory.
- `GETFILES$` **[PARTIAL]**: Coming soon.
- `GETDIRS$` **[DONE]**: Gets a list of subdirectories in a path.
- `GETFILEEXT$` **[DONE]**: Returns the extension of a file path.
- `GETFILENAME$` **[DONE]**: Returns the file name from a path.
- `GETFILENAMENOEXT$` **[DONE]**: Returns the file name without the extension.
- `GETFILEPATH$` **[DONE]**: Returns the directory path from a file path.
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
- `DATE$` **[DONE]**: Returns the current date as a string.
- `TIME$` **[DONE]**: Returns the current time as a string.
- `DATETIME$` **[DONE]**: Returns the current date and time as a string.
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
- `ENVIRON$` **[DONE]**: Gets an environment variable.
- `COMMAND$` **[DONE]**: Gets a command-line argument by index.
- `ARGC` **[DONE]**: Gets the number of command-line arguments.

---

## Module Commands

moonBASIC uses a dot-notation module system for its game engine commands. These all use the form `Module.Command(args)`. They can also be called on a handle variable: `cam = Camera.Make()` then `cam.SetPos(0, 5, 10)`.

> **Note:** Commands listed as `**[PARTIAL]**` or `**[MISSING]**` in this section are planned features that are not yet fully implemented.

---

### Window — [Reference](reference/WINDOW.md)

| Command | Description |
|---|---|
| `Window.Open(w, h, title$)` | Opens the window. Must be the first command called. |
| `Window.Close()` | Closes the window and exits. |
| `Window.ShouldClose()` | Returns `TRUE` when the user requests close (use in main loop). |
| `Window.SetFPS(fps)` | Sets the target frame rate. |
| `Window.SetTitle(title$)` | Updates the window title at runtime. |
| `Window.SetPosition(x, y)` | Moves the window on screen. |
| `Window.SetSize(w, h)` | Resizes the window. |
| `Window.SetMinSize(w, h)` | Sets the minimum resizable size. |
| `Window.SetMaxSize(w, h)` | Sets the maximum resizable size. |
| `Window.SetIcon(path$)` | Sets the window icon from an image file. |
| `Window.SetOpacity(alpha#)` | Sets window transparency (0.0–1.0). |
| `Window.SetMonitor(idx)` | Moves the window to a specific monitor. |
| `Window.GetMonitorCount()` | Returns the number of connected monitors. |
| `Window.GetMonitorWidth(idx)` | Returns a monitor's width in pixels. |
| `Window.GetMonitorHeight(idx)` | Returns a monitor's height in pixels. |
| `Window.GetMonitorRefreshRate(idx)` | Returns a monitor's refresh rate in Hz. |
| `Window.GetMonitorName(idx)` | Returns a monitor's name string. |
| `Window.GetPositionX()` | Returns the window's X position on screen. |
| `Window.GetPositionY()` | Returns the window's Y position on screen. |
| `Window.GetScaleDPIX()` | Returns the DPI X scale factor. |
| `Window.GetScaleDPIY()` | Returns the DPI Y scale factor. |
| `Window.SetFlag(flag)` | Sets a window state flag (e.g., `FLAG_RESIZABLE`). |
| `Window.ClearFlag(flag)` | Clears a window state flag. |
| `Window.CheckFlag(flag)` | Returns `TRUE` if the window flag is set. |

---

### Render — [Reference](reference/RENDER.md)

| Command | Description |
|---|---|
| `Render.Clear(r, g, b, [a])` | Clears the screen to a color. Call at the start of each frame. |
| `Render.Frame()` | Presents the rendered frame. Call at the end of each frame. |
| `Render.BeginShader(shaderHandle)` | Applies a custom shader to all subsequent drawing. |
| `Render.EndShader()` | Ends custom shader mode. |
| `Render.DrawFPS(x, y)` | Draws the current FPS on screen. |
| `Render.Screenshot(path$)` | Saves a screenshot to a PNG file. |
| `Render.SetBlend(mode)` | Sets the blend mode (e.g., `BLEND_ALPHA`, `BLEND_ADDITIVE`). |
| `Render.SetAmbient(r#, g#, b#)` | Sets the 3D scene's ambient light color (0.0–1.0). |

---

### Camera — [Reference](reference/CAMERA.md)

| Command | Description |
|---|---|
| `Camera.Make()` | Creates a new 3D camera, returns a handle. |
| `cam.SetPos(x#, y#, z#)` | Sets the camera's world position. |
| `cam.SetTarget(x#, y#, z#)` | Sets the point the camera looks at. |
| `cam.SetFOV(fov#)` | Sets the vertical field of view in degrees. |
| `cam.Begin()` | Enters 3D rendering mode with this camera. |
| `cam.End()` | Exits 3D rendering mode. |
| `cam.GetMatrix()` | Returns the camera's view matrix handle. |
| `Camera2D.Make()` | Creates a 2D scrolling camera, returns a handle. |
| `cam2d.Begin()` | Enters 2D camera rendering mode. |
| `cam2d.End()` | Exits 2D camera rendering mode. |
| `cam2d.SetTarget(x#, y#)` | Sets the 2D camera's world target. |
| `cam2d.SetOffset(x#, y#)` | Sets the screen-space offset (for centering). |
| `cam2d.SetZoom(zoom#)` | Sets the zoom level (1.0 = normal). |
| `cam2d.SetRotation(angle#)` | Sets the camera rotation in degrees. |

---

### Draw — [Reference](reference/DRAW2D.md)

| Command | Description |
|---|---|
| `Draw.Rectangle(x, y, w, h, r, g, b, a)` | Draws a filled rectangle. |
| `Draw.RectangleRounded(x, y, w, h, radius, r, g, b, a)` | Draws a filled rectangle with rounded corners. |
| `Draw.Circle(cx, cy, radius, r, g, b, a)` | Draws a filled circle. |
| `Draw.Line(x1, y1, x2, y2, r, g, b, a)` | Draws a line segment. |
| `Draw.Text(text$, x, y, size, r, g, b, a)` | Draws text with the default font. |
| `Draw.TextFont(fontHandle, text$, x, y, size, spacing, r, g, b, a)` | Draws text with a custom font. |
| `Draw.Texture(texHandle, x, y, [r, g, b, a])` | Draws a texture at a position. |
| `Draw.TextureNPatch(texHandle, l, t, r, b, x, y, w, h, r, g, b, a)` | Draws a 9-patch texture for UI elements. |
| `Draw.Grid(slices, spacing#)` | Draws a 3D grid (use inside `cam.Begin()`/`cam.End()`). |

---

### Texture — [Reference](reference/TEXTURE.md)

| Command | Description |
|---|---|
| `Texture.Load(path$)` | Loads a PNG/JPG from disk into GPU memory. Returns a handle. |
| `Texture.Free(handle)` | Unloads a texture from GPU memory. |
| `Texture.FromImage(imgHandle)` | Creates a GPU texture from an in-memory image. |
| `Texture.GenWhiteNoise(w, h)` | Generates a white noise texture procedurally. |

---

### Font — [Reference](reference/FONT.md)

| Command | Description |
|---|---|
| `Font.Load(path$)` | Loads a `.ttf` or `.otf` font file. Returns a handle. |
| `Font.Free(handle)` | Unloads a font from memory. |

---

### GUI (raygui) — [Reference](reference/GUI.md)

Immediate-mode widgets from raygui (`GUI.*`). Requires a **CGO** build like `DRAW.*` / `WINDOW.*`.

| Command | Description |
|---|---|
| `GUI.ENABLE` / `GUI.DISABLE` | Toggle global GUI input. |
| `GUI.BUTTON` / `GUI.LABEL` / `GUI.SLIDER` / … | See [GUI.md](reference/GUI.md) for the full list (layout, lists, color pickers, dialogs, tooltips). |

Runnable demos: `examples/gui_basics/main.mb`, `examples/gui_form/main.mb`.

---

### Sprite & Atlas — [Reference: Sprite](reference/SPRITE.md) · [Atlas](reference/ATLAS.md)

| Command | Description |
|---|---|
| `Sprite.Load(path$)` | Loads an `.aseprite` file. Returns a handle. |
| `Sprite.DefAnim(handle, name$)` | Registers an animation tag for use. |
| `Sprite.PlayAnim(handle, name$)` | Sets the active animation. |
| `Sprite.UpdateAnim(handle, dt#)` | Advances the animation timer. Call every frame. |
| `Sprite.Draw(handle, x, y)` | Draws the current animation frame. |
| `Atlas.Load(imgPath$, jsonPath$)` | Loads a texture atlas (image + JSON data). |
| `Atlas.GetSprite(handle, name$)` | Gets a sub-sprite handle from the atlas. |
| `Atlas.Free(handle)` | Frees the atlas and all its sprites. |

---

### Model, Mesh & Material — [Reference](reference/MODEL.md)

| Command | Description |
|---|---|
| `Model.Load(path$)` | Loads a 3D model file (`.gltf`, `.glb`, `.obj`). |
| `Model.Draw(handle, matrixHandle)` | Draws a model with a transform matrix. |
| `Model.SetMaterial(handle, idx, matHandle)` | Replaces one of a model's materials. |
| `Model.Free(handle)` | Unloads a model from memory. |
| `Mesh.MakeCube(w#, h#, d#)` | Creates a procedural box mesh. |
| `Mesh.MakeSphere(radius#, rings, slices)` | Creates a procedural sphere mesh. |
| `Mesh.MakePlane(w#, len#, resX, resZ)` | Creates a procedural flat plane mesh. |
| `Mesh.MakeHeightmap(imgHandle, w#, h#, len#)` | Creates a terrain mesh from a heightmap image. |
| `Mesh.Draw(meshHandle, matHandle, matrixHandle)` | Draws a single mesh with a material and transform. |
| `Mesh.Free(handle)` | Unloads a mesh from memory. |
| `Material.MakeDefault()` | Creates a default PBR material. |
| `Material.MakePBR()` | Creates a full PBR material with shadow support. |
| `Material.SetTexture(handle, slot, texHandle)` | Assigns a texture to a material map slot. |
| `Material.SetColor(handle, slot, r, g, b, a)` | Sets a material map slot's tint color. |
| `Material.SetFloat(handle, slot, value#)` | Sets a material map slot's float value. |
| `Material.SetShader(handle, shaderHandle)` | Applies a custom shader to a material. |
| `Material.Free(handle)` | Frees a material. |

---

### Physics 3D (Jolt) — [Reference](reference/PHYSICS3D.md)

| Command | Description |
|---|---|
| `Physics3D.Start()` | Initializes the 3D physics world. |
| `Physics3D.Stop()` | Shuts down the 3D physics world. |
| `Physics3D.Step()` | Advances the simulation one step. Call once per frame. |
| `Physics3D.SetGravity(x#, y#, z#)` | Sets global gravity. |
| `Body3D.Make(type$)` | Creates a body definition (`"static"`, `"dynamic"`, `"kinematic"`). |
| `Body3D.AddBox(defHandle, w#, h#, d#)` | Adds a box collision shape. |
| `Body3D.AddSphere(defHandle, radius#)` | Adds a sphere collision shape. |
| `Body3D.AddCapsule(defHandle, height#, radius#)` | Adds a capsule collision shape. |
| `Body3D.Commit(defHandle, x#, y#, z#)` | Finalizes the body and adds it to the world. |
| `Body3D.SetPos(handle, x#, y#, z#)` | Teleports a body to a new position. |
| `Body3D.SetMass(handle, mass#)` | Sets the mass of a dynamic body. |
| `Body3D.ApplyForce(handle, x#, y#, z#)` | Applies a continuous force. |
| `Body3D.ApplyImpulse(handle, x#, y#, z#)` | Applies an instant impulse. |
| `Body3D.SetLinearVel(handle, vx#, vy#, vz#)` | Sets linear velocity directly. |
| `Body3D.GetMatrix(handle)` | Returns the body's transform matrix (for rendering). |
| `Body3D.Free(handle)` | Removes a body from the simulation. |

---

### Physics 2D (Box2D) — [Reference](reference/PHYSICS2D.md)

| Command | Description |
|---|---|
| `Physics2D.Start()` | Initializes the 2D physics world. |
| `Physics2D.Stop()` | Shuts down the 2D physics world. |
| `Physics2D.Step()` | Advances the simulation. Call once per frame. |
| `Physics2D.SetGravity(x#, y#)` | Sets global 2D gravity (positive Y = down). |
| `Body2D.Make(type$)` | Creates a 2D body definition. |
| `Body2D.AddRect(defHandle, w#, h#)` | Adds a rectangle collision shape. |
| `Body2D.AddCircle(defHandle, radius#)` | Adds a circle collision shape. |
| `Body2D.Commit(defHandle, x#, y#)` | Adds the body to the world at a position. |
| `Body2D.SetPos(handle, x#, y#)` | Teleports a 2D body. |
| `Body2D.X(handle)` / `Body2D.Y(handle)` | Returns the body's X or Y position. |
| `Body2D.Rot(handle)` | Returns the body's rotation in degrees. |
| `Body2D.ApplyForce(handle, x#, y#)` | Applies a continuous 2D force. |
| `Body2D.ApplyImpulse(handle, x#, y#)` | Applies an instant 2D impulse. |
| `Body2D.Free(handle)` | Removes a body from the simulation. |

---

### Character Controller — [Reference](reference/CHARCONTROLLER.md)

| Command | Description |
|---|---|
| `CharController.Make(radius#, height#, x#, y#, z#)` | Creates a capsule character controller. |
| `CharController.Move(handle, dx#, dy#, dz#)` | Moves the controller with collision detection. |
| `CharController.IsGrounded(handle)` | Returns `TRUE` if standing on a surface. |
| `CharController.X(handle)` / `.Y()` / `.Z()` | Returns the controller's position component. |
| `CharController.Free(handle)` | Frees the controller. |

---

### Audio — [Reference](reference/AUDIO.md)

| Command | Description |
|---|---|
| `Audio.Init()` | Initializes the audio device. |
| `Audio.Close()` | Closes the audio device. |
| `Audio.LoadSound(path$)` | Loads a sound effect into memory. Returns a handle. |
| `Audio.LoadMusic(path$)` | Loads a music file for streaming. Returns a handle. |
| `Audio.Play(handle)` | Plays a sound or music track. |
| `Audio.Stop(handle)` | Stops playback. |
| `Audio.Pause(handle)` | Pauses playback. |
| `Audio.Resume(handle)` | Resumes paused playback. |
| `Audio.UpdateMusic(handle)` | Updates the music stream buffer. **Call every frame** for music. |
| `Sound.Free(handle)` | Unloads a sound effect. |
| `Music.Free(handle)` | Unloads a music stream. |
| `AudioStream.Make(sampleRate, bitDepth, channels)` | Creates a raw PCM audio stream. |
| `AudioStream.Update(handle, pcmArray)` | Pushes PCM data to the stream. |
| `AudioStream.Play(handle)` | Starts the audio stream. |
| `AudioStream.IsPlaying(handle)` | Returns `TRUE` if the stream is playing. |
| `AudioStream.Free(handle)` | Frees the audio stream. |

---

### Network (ENet) — [Reference](reference/NETWORK.md)

| Command | Description |
|---|---|
| `Net.Start()` | Initializes the networking system. |
| `Net.Stop()` | Shuts down the networking system. |
| `Net.CreateServer(port, maxClients)` | Creates a server host. Returns a handle. |
| `Net.CreateClient()` | Creates a client host. Returns a handle. |
| `Net.Connect(clientHandle, address$, port)` | Connects a client to a server. Returns a peer handle. |
| `Net.Update(hostHandle)` | Processes network packets. **Call every frame.** |
| `Net.Receive(hostHandle)` | Returns the next event handle, or `0` if none queued. |
| `Net.Broadcast(serverHandle, channel, data$, reliable?)` | Sends a message to all connected clients. |
| `Peer.Send(peerHandle, channel, data$, reliable?)` | Sends a message to a specific peer. |
| `Event.Type(handle)` | Returns the event type (`EVENT_CONNECT`, `EVENT_DISCONNECT`, `EVENT_RECEIVE`). |
| `Event.Peer(handle)` | Returns the peer associated with the event. |
| `Event.Data(handle)` | Returns the string data of a `RECEIVE` event. |
| `Event.Free(handle)` | Frees the event. **Must be called for every event received.** |

---

### Time — [Reference](reference/TIME.md)

| Command | Description |
|---|---|
| `Time.Delta()` | Returns seconds elapsed since last frame. Use this for frame-rate-independent movement. |
| `Time.Get()` | Returns total seconds elapsed since the program started. |

---

### Mat4 / Matrix Math — [Reference](reference/MAT4.md)

| Command | Description |
|---|---|
| `Mat4.Identity()` | Creates a new identity matrix. Returns a handle. |
| `Mat4.SetRotation(handle, rx#, ry#, rz#)` | Sets the rotation of an existing matrix (radians). |
| `Mat4.FromRotation(rx#, ry#, rz#)` | Creates a rotation matrix. Returns a handle. |
| `Mat4.FromScale(sx#, sy#, sz#)` | Creates a scale matrix. Returns a handle. |
| `Mat4.FromTranslation(x#, y#, z#)` | Creates a translation matrix. Returns a handle. |
| `Mat4.Multiply(a, b)` | Multiplies two matrices. Returns a new handle. |
| `Mat4.Inverse(handle)` | Returns the inverse of a matrix. |
| `Mat4.Transpose(handle)` | Returns the transpose of a matrix. |
| `Mat4.LookAt(eyeX#, eyeY#, eyeZ#, atX#, atY#, atZ#, upX#, upY#, upZ#)` | Creates a look-at view matrix. |
| `Mat4.Free(handle)` | Frees a matrix handle. |

---

### Shader — [Reference](reference/SHADER.md)

| Command | Description |
|---|---|
| `Shader.Load(vsPath$, fsPath$)` | Loads GLSL vertex and fragment shaders. Returns a handle. |

---

### Light — [Reference](reference/LIGHT.md)

| Command | Description |
|---|---|
| `Light.Make(type$)` | Creates a light (`"directional"`, `"point"`, `"spot"`). Returns a handle. |
| `Light.SetDir(handle, x#, y#, z#)` | Sets the direction for directional/spot lights. |
| `Light.SetShadow(handle, enabled?)` | Enables or disables shadow casting. |

---

### Tilemap — [Reference](reference/TILEMAP.md)

| Command | Description |
|---|---|
| `Tilemap.Load(path$)` | Loads a Tiled `.tmx` file. Returns a handle. |
| `Tilemap.Draw(handle, offsetX, offsetY)` | Draws all layers of the tilemap. |
| `Tilemap.DrawLayer(handle, layer, offsetX, offsetY)` | Draws a single named or indexed layer. |
| `Tilemap.GetTile(handle, layer, x, y)` | Returns the tile GID at a grid position. |
| `Tilemap.SetTile(handle, layer, x, y, gid)` | Sets a tile at a grid position at runtime. |
| `Tilemap.IsSolid(handle, x, y)` | Returns `TRUE` if the tile at (x,y) has collision. |
| `Tilemap.SetTileSize(handle, w, h)` | Overrides the display size of each tile. |
| `Tilemap.Width(handle)` | Returns the tilemap width in tiles. |
| `Tilemap.Height(handle)` | Returns the tilemap height in tiles. |
| `Tilemap.LayerCount(handle)` | Returns the number of layers. |
| `Tilemap.LayerName(handle, idx)` | Returns the name of a layer by index. |
| `Tilemap.Free(handle)` | Frees the tilemap from memory. |

---

### Particles — [Reference](reference/PARTICLES.md)

| Command | Description |
|---|---|
| `Particle.Make(maxCount)` | Creates a particle emitter. Returns a handle. |
| `Particle.SetTexture(handle, texHandle)` | Sets the particle texture. |
| `Particle.SetEmitRate(handle, rate#)` | Sets particles emitted per second. |
| `Particle.SetLifetime(handle, min#, max#)` | Sets the particle lifetime range in seconds. |
| `Particle.SetVelocity(handle, vx#, vy#, vz#, spread#)` | Sets the emission velocity and spread. |
| `Particle.SetColor(handle, r, g, b, a)` | Sets the starting particle color. |
| `Particle.SetColorEnd(handle, r, g, b, a)` | Sets the ending particle color (fades to this). |
| `Particle.SetSize(handle, startSize#, endSize#)` | Sets the particle size over its lifetime. |
| `Particle.SetGravity(handle, gx#, gy#, gz#)` | Sets per-emitter gravity. |
| `Particle.SetPos(handle, x#, y#, z#)` | Sets the emitter position in world space. |
| `Particle.Play(handle)` | Starts the emitter. |
| `Particle.Update(handle, dt#)` | Advances particle simulation. Call every frame. |
| `Particle.Draw(handle)` | Draws all active particles. |
| `Particle.Free(handle)` | Frees the particle emitter. |
