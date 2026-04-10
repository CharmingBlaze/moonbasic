# Input Commands

Commands for handling keyboard, mouse, and gamepad input.

---

## Direct Input

These commands directly query the state of a specific key or button.

### `Input.KeyDown(key)` / flat `KEYDOWN(key)`

Returns `TRUE` if the specified key is currently being held down.

- `key`: The key code to check (e.g., `KEY_SPACE`, `KEY_W`).
- **`KEYDOWN`** is the same behaviour as **`Input.KeyDown`** (shorter name for game loops).

### `Input.KeyPressed(key)` / flat `KEYPRESSED(key)`

Returns `TRUE` only on the frame the specified key was first pressed. Use this for **jump**, **shoot**, **menu confirm**, or any action that should fire **once per press**.

### `Input.KeyUp(key)`

Returns `TRUE` only on the frame the specified key was released.

### `Input.Axis(negKey, posKey)` → **float**

Returns a **discrete axis** in **`{-1.0, 0.0, 1.0}`** from two key codes. Use for **WASD-style** movement without separate accumulators:

- **`1.0`** — `posKey` is held and `negKey` is not.
- **`-1.0`** — `negKey` is held and `posKey` is not.
- **`0.0`** — both held, or neither held.

```basic
f# = Input.Axis(KEY_S, KEY_W)   ; forward: W positive, S negative (if your camera treats +f as forward, swap keys)
s# = Input.Axis(KEY_A, KEY_D)   ; strafe: D positive, A negative
```

The registry name is **`INPUT.AXIS`** (see [API_CONSISTENCY.md](../API_CONSISTENCY.md)).

### `Input.AxisDeg(negKey, posKey, degreesPerSec#, dt#)` → **float**

Bundled **rotation delta** for this frame: **same result** as **`Input.Axis(negKey, posKey) * DEGPERSEC(degreesPerSec, dt)`** — radians to add to a yaw (e.g. **`camYaw`**). Registry: **`INPUT.AXISDEG`**.

```basic
camYaw# = camYaw# + Input.AxisDeg(KEY_Q, KEY_E, 77.0, dt#)
```

### `Input.Orbit(negKey, posKey, degreesPerSec#, dt#)` → **float**

**Alias** of **`Input.AxisDeg`** (registry **`INPUT.ORBIT`**) — same arguments and result; use whichever reads clearer for camera orbit.

---

### `Input.Movement2D(keyBack, keyForward, keyLeft, keyRight)` → **handle**

Returns a **2-float array** **`[forwardAxis, strafeAxis]`** using two **`Input.Axis`** pairs — same as:

```basic
f# = Input.Axis(KEY_S, KEY_W)
s# = Input.Axis(KEY_A, KEY_D)
```

Registry: **`INPUT.MOVEMENT2D`**. **ERASE** the handle when you no longer need it (each frame if you allocate every frame). For zero extra allocations, call **`Input.Axis`** twice instead.

---

Typical use: **Q/E orbit** with **`Camera.OrbitAround`** — see **`examples/mario64/main_v2.mb`** and [CAMERA.md](CAMERA.md).

---

### When to use which (keyboard)

| API | Fires when | Good for |
|-----|------------|----------|
| **`KEYDOWN` / `Input.KeyDown`** | Every frame while held | Movement, strafe, holding a key to rotate |
| **`KEYPRESSED` / `Input.KeyPressed`** | First frame only | Jump, fire once, UI toggle |
| **`KEYRELEASED` / `Input.KeyReleased`** | The frame the key goes up | Charge attacks, hold-to-aim release |

Example: move with **`KEYDOWN`**, jump with **`KEYPRESSED`** so holding space does not retrigger mid-air.

```basic
IF KEYDOWN(KEY_W) THEN player_z# = player_z# - speed# * DT()
IF on_ground? AND KEYPRESSED(KEY_SPACE) THEN vel_y# = jump_strength#
```

---

### `Input.MouseDown(button)`

Returns `TRUE` if the specified mouse button is currently being held down.

- `button`: The mouse button to check (e.g., `MOUSE_LEFT_BUTTON`).

### `Input.MouseX()` / `Input.MouseY()`

Returns the current X or Y coordinate of the mouse cursor.

### `Input.MouseWheel()` / `INPUT.MOUSEWHEELMOVE()` / `MouseWheel()`

Returns the **mouse wheel** delta for this frame (implementation follows Raylib). **`Input.MouseWheel`** is a shorter alias for the same behavior as **`INPUT.MOUSEWHEELMOVE`**.

---

## Cursor

Registry names are prefixed with **`CURSOR.`** (e.g. **`CURSOR.HIDE()`**). Use these for games that should not show the system pointer over the window.

### `CURSOR.HIDE()` / `CURSOR.SHOW()`

- **`CURSOR.HIDE()`** — Hides the **OS mouse cursor** while it is over the window. **Mouse position still updates** (`INPUT.MOUSEX` / `INPUT.MOUSEY`); only the drawn pointer disappears.
- **`CURSOR.SHOW()`** — Shows the cursor again.

Call **`CURSOR.HIDE()`** after **`WINDOW.OPEN`** (or once the player is “in game”), and call **`CURSOR.SHOW()`** before **`WINDOW.CLOSE`** so the shell / desktop gets a visible cursor back after exit.

**3D orbit / drag:** Hiding the cursor does **not** enable “FPS mode” or raw mouse deltas. For **game-style** play, use **`CURSOR.DISABLE()`** so the virtual cursor stays **centered** and **`INPUT.MOUSEDELTAX`** / **`INPUT.MOUSEDELTAY`** report movement (see **`examples/mario64/main_entities.mb`**). Optionally call **`INPUT.SETMOUSEPOS(cx, cy)`** once after **`WINDOW.OPEN`** (e.g. half of width/height) to park the hardware cursor before relative mode.

### `INPUT.SETMOUSEPOS(x, y)`

Warps the mouse to **client-area pixel** coordinates. Often used with **`CURSOR.DISABLE()`** to re-seat the cursor after **`WINDOW.OPEN`** or when regaining focus.

### `CURSOR.DISABLE()` / `CURSOR.ENABLE()` (raw / relative mouse)

- **`CURSOR.DISABLE()`** — Raylib **disables** the cursor and switches to **relative** mouse mode (movement reported as deltas; pointer behaves like a centered “virtual” game cursor). Pair with **`CURSOR.ENABLE()`** before showing the OS cursor again (e.g. before **`WINDOW.CLOSE`**).
- Prefer **`CURSOR.HIDE()`** only if you need the pointer invisible but still want **absolute** **`INPUT.MOUSEX`** / **`INPUT.MOUSEY`** (no relative deltas).

### Other

- **`CURSOR.ISHIDDEN()`** — Whether the cursor is currently hidden.
- **`CURSOR.SET(id)`** — Sets the system cursor shape (Raylib / platform cursor id).

---

## Action Mapping

The action mapping system is a powerful way to handle input. Instead of checking for specific keys, you define abstract "actions" and then check the state of those actions. This makes it easy to support multiple input devices and allow for user-configurable controls.

### 1. Define Mappings

First, map physical inputs to action names. This is typically done once at the start of your program.

- `Input.MapKey(action, key)`: Maps a keyboard key to an action.
- `Input.MapGamepadButton(action, gamepad, button)`: Maps a gamepad button.
- `Input.MapGamepadAxis(action, gamepad, axis, direction)`: Maps a gamepad axis direction (e.g., left stick up) to an action.

```basic
; Define a "move_right" action for multiple inputs
Input.MapKey("move_right", KEY_D)
Input.MapKey("move_right", KEY_RIGHT) ; Also map the right arrow key
Input.MapGamepadButton("move_right", 0, GAMEPAD_BUTTON_RIGHT_FACE_RIGHT) ; D-pad right
```

### 2. Check Actions in the Game Loop

In your main loop, check the state of the action by its name, not the key.

- `Input.ActionDown(action)`: Returns `TRUE` if any mapped input is currently active (held down).
- `Input.ActionPressed(action)`: Returns `TRUE` only on the frame the action was first triggered.
- `Input.ActionReleased(action)`: Returns `TRUE` on the frame the action was released.

```basic
; In the main loop, check the abstract action
IF Input.ActionDown("move_right") THEN
    player_x# = player_x# + speed# * Time.Delta()
ENDIF
```

### `Input.ActionAxis(action$)`

For analog controls like a joystick, this returns the axis value from -1.0 to 1.0.

```basic
; Map joystick axes
Input.MapGamepadAxis("move_horizontal", 0, GAMEPAD_AXIS_LEFT_X)

; Get the analog value
move_x# = Input.ActionAxis("move_horizontal")
player_x# = player_x# + move_x# * speed# * Time.Delta()
```
