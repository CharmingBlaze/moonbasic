# Input systems: INPUT and ACTION

> Raw keyboard, mouse, and gamepad polling plus named action maps for game-friendly controls.

**All commands:** [COMMAND_REGISTRY.md#input-action](COMMAND_REGISTRY.md#input-action)

**Deep guide:** [guides/CAMERA-AND-INPUT.md](guides/CAMERA-AND-INPUT.md)

**See also:** [01-CORE](01-CORE.md) · [reference/INPUT.md](../reference/INPUT.md) · [reference/ACTION.md](../reference/ACTION.md)

**Case:** `KEY_W`, `key_w`, and `Key_W` are equivalent ([LANGUAGE.md](../LANGUAGE.md)).

---

## Table of contents

- [INPUT system](#input-system)
- [ACTION system](#action-system)
- [Full example](#full-example)
- [See also](#see-also)

---

## INPUT system

Low-level device state — poll every frame inside your main loop.

### Core workflow

1. Each frame, read `INPUT.KEYDOWN` / `KEYHIT` for movement and menus.
2. Read `INPUT.MOUSEX` / `MOUSEY` and button helpers for UI and aiming.
3. Optional: `INPUT.GAMEPADCONNECTED` and axis/button queries.

---

### Keyboard

| Command | Description |
|---------|-------------|
| `INPUT.KEYDOWN(key)` | Key held this frame |
| `INPUT.KEYHIT(key)` | Key pressed this frame (edge) |
| `INPUT.KEYUP(key)` | Key released this frame |
| `INPUT.KEYPRESSED(key)` | Alias for edge-pressed |

**Example:**

```basic
IF INPUT.KEYDOWN(KEY_W) THEN
    player.move(0, 0, speed * APP.DELTA())
ENDIF
IF INPUT.KEYHIT(KEY_SPACE) THEN
    jump()
ENDIF
```

Common keys: `KEY_ESCAPE`, `KEY_W`, `KEY_A`, `KEY_S`, `KEY_D`, `KEY_SPACE`.

---

### Mouse

| Command | Description |
|---------|-------------|
| `INPUT.MOUSEDOWN(button)` | Button held |
| `INPUT.MOUSEHIT(button)` | Button pressed this frame |
| `INPUT.MOUSEX()` / `INPUT.MOUSEY()` | Cursor in window pixels |
| `INPUT.MOUSEDELTA_X()` / `MOUSEDELTA_Y()` | Motion since last frame |
| `INPUT.MOUSEWHEEL()` | Scroll delta |

**Example:**

```basic
mx = INPUT.MOUSEX()
my = INPUT.MOUSEY()
IF INPUT.MOUSEHIT(MOUSE_LEFT) THEN
    ; click at mx, my
ENDIF
```

---

### Gamepad

| Command | Description |
|---------|-------------|
| `INPUT.GAMEPADCONNECTED(index)` | Pad present |
| `INPUT.GAMEPADBUTTONDOWN(pad, button)` | Button held |
| `INPUT.GAMEPADAXIS(pad, axis)` | Stick value −1…1 |

**Aliases:** checklist `GAMEPAD*` names map to these — see [reference/INPUT.md](../reference/INPUT.md).

**Example:**

```basic
IF INPUT.GAMEPADCONNECTED(0) THEN
    lx = INPUT.GAMEPADAXIS(0, PAD_LEFT_X)
ENDIF
```

---

## ACTION system

Named actions sit above raw keys — map once, read `ACTION.DOWN("Forward")` in gameplay code.

### Core workflow

1. At startup: `ACTION.MAPKEY` / `BINDKEY` / `BINDGAMEPAD`.
2. Each frame: `ACTION.DOWN`, `ACTION.PRESSED` (or `HIT`), `ACTION.VALUE` for axes.
3. `ACTION.RESET` to clear bindings when changing control schemes.

---

### `ACTION.MAPKEY(action, key)` / `ACTION.BINDKEY(action, key)`

Binds a keyboard key to a string action name.

| Argument | Type | Description |
|----------|------|-------------|
| action | string | Name (e.g. `"Jump"`) |
| key | int | Key constant |

**Returns:** nothing

**Aliases:** `ACTION.BINDKEY`

**Example:**

```basic
ACTION.MAPKEY("Jump", KEY_SPACE)
ACTION.MAPKEY("Forward", KEY_W)
```

---

### `ACTION.BINDGAMEPAD(action, button)`

Maps a gamepad button to an action.

**Example:**

```basic
ACTION.BINDGAMEPAD("Jump", PAD_A)
```

---

### `ACTION.DOWN(action)` / `ACTION.PRESSED(action)` / `ACTION.HIT(action)`

| Command | Description |
|---------|-------------|
| `ACTION.DOWN(name)` | Action active this frame |
| `ACTION.PRESSED(name)` | Edge pressed |
| `ACTION.HIT(name)` | Alias for pressed |
| `ACTION.RELEASED(name)` | Edge released |
| `ACTION.VALUE(name)` | Axis-style action (−1…1) |

**Example:**

```basic
IF ACTION.DOWN("Forward") THEN
    player.move(0, 0, speed * APP.DELTA())
ENDIF
IF ACTION.HIT("Jump") THEN
    doJump()
ENDIF
```

---

### `ACTION.MAPAXIS(action, axis)` / `ACTION.MAPMOUSE(action, button)`

Bind stick axes or mouse buttons to named actions.

See [reference/ACTION.md](../reference/ACTION.md) for full mapping tables.

---

## Full example

```basic
APP.OPEN(640, 480, "Input")
APP.SETFPS(60)

ACTION.MAPKEY("Forward", KEY_W)
ACTION.MAPKEY("Back", KEY_S)
ACTION.BINDKEY("Jump", KEY_SPACE)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 2, -6)

cube = ENTITY.CREATECUBE(1, 1, 1)

WHILE NOT APP.SHOULDCLOSE()
    IF ACTION.DOWN("Forward") THEN cube.move(0, 0, 3 * APP.DELTA())
    IF ACTION.DOWN("Back") THEN cube.move(0, 0, -3 * APP.DELTA())
    IF ACTION.HIT("Jump") THEN cube.move(0, 2, 0)

    RENDER.CLEAR(20, 22, 30)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

APP.CLOSE()
```

---

## See also

- [05-PHYSICS](05-PHYSICS.md) — movement + collision
- [08-UI-TEXT](08-UI-TEXT.md) — GUI mouse interaction
- [examples/gamepad](../examples/gamepad/README.md) — controller sample
