# Character Physics (KCC)

Kinematic Character Controllers (KCC) provide gameplay-focused movement that is smoother and more controllable than pure rigid-body physics. In MoonBASIC, this is powered by Jolt's `CharacterVirtual` system on desktop platforms.

> [!NOTE]
> KCC requires **Windows or Linux + CGO + Jolt**. On other platforms, these commands return an error.

## Core Workflow

1. **Initialize Physics**: Call `PHYSICS3D.START()` once.
2. **Create Entity**: Create a model (e.g., `MODEL.CREATECAPSULE`) and an entity.
3. **Setup Controller**: Call `CHAR.CREATE(entity)` to attach a physical capsule to the entity.
4. **Game Loop**: Use `CHAR.MOVE` or `CHAR.MOVEWITHCAM` each frame.
5. **Detection**: Check `CHAR.ISGROUNDED()` for jumping logic.

---

## Setup & Configuration

### `CHAR.CREATE(entity, radius, height)` 

Attaches a Kinematic Character Controller to an entity. This clears any scripted gravity and allows the physics capsule to drive the entity's position.

- *Handle shortcut*: `entity.charCreate(radius, height)`
- **Arguments**:
  - `entity` (handle): The entity to control.
  - `radius` (float): Horizontal radius of the capsule.
  - `height` (float): Total vertical height of the capsule.
- **Returns**: (handle) Returns the entity handle for chaining.
- **Alias**: `PLAYER.CREATE`, `CHARACTER.CREATE`

- **Example**:
  ```basic
  ; Create a standard human-sized character
  hero = ENTITY.CREATE(MODEL.LOAD("hero.iqm"))
  hero.charCreate(0.4, 1.8)
  ```

---

### `CHAR.SETSTEP(entity, height)` 

Sets the maximum height that the character can automatically "step up" (e.g., stairs or curbs).

- *Handle shortcut*: `entity.setStep(height)`
- **Arguments**:
  - `entity` (handle): The character entity.
  - `height` (float): Maximum step height in world units.
- **Returns**: (handle) Returns the entity handle for chaining.
- **Alias**: `PLAYER.SETSTEPOFFSET`

---

### `CHAR.SETSLOPE(entity, degrees)` 

Sets the maximum slope angle (in degrees) that the character can walk up. Steeper slopes will be treated as walls.

- *Handle shortcut*: `entity.setSlope(degrees)`
- **Arguments**:
  - `entity` (handle): The character entity.
  - `degrees` (float): Max slope angle (e.g., 45.0).
- **Returns**: (handle) Returns the entity handle for chaining.
- **Alias**: `PLAYER.SETSLOPELIMIT`

---

## Movement

### `CHAR.MOVE(entity, dirX, dirZ, speed)` 

Moves the character in a world-space direction.

- *Handle shortcut*: `entity.move(dirX, dirZ, speed)`
- **Arguments**:
  - `entity` (handle): The character entity.
  - `dirX, dirZ` (float): Horizontal direction vector (normalized).
  - `speed` (float): Movement speed in units per second.
- **Returns**: (none)

---

### `CHAR.MOVEWITHCAM(entity, cam, fwd, side, speed)` 

Moves the character relative to a camera's orientation (standard WASD control).

- *Handle shortcut*: `entity.moveWithCam(cam, fwd, side, speed)`
- **Arguments**:
  - `entity` (handle): The character entity.
  - `cam` (handle): The camera used for orientation.
  - `fwd` (float): Forward/Back input (-1 to 1).
  - `side` (float): Left/Right input (-1 to 1).
  - `speed` (float): Movement speed.
- **Returns**: (none)
- **Alias**: `CHAR.MOVEWITHCAMERA`

- **Example**:
  ```basic
  fwd  = INPUT.AXIS(KEY_S, KEY_W)
  side = INPUT.AXIS(KEY_A, KEY_D)
  hero.moveWithCam(cam, fwd, side, 10.0)
  ```

---

### `CHAR.JUMP(entity, impulse)` 

Applies an upward vertical impulse to make the character jump.

- *Handle shortcut*: `entity.jump(impulse)`
- **Arguments**:
  - `entity` (handle): The character entity.
  - `impulse` (float): Upward strength of the jump.
- **Returns**: (none)

---

## World Interaction

### `CHAR.ISGROUNDED(entity, coyoteTime)` 

Checks if the character is currently standing on a supported surface.

- *Handle shortcut*: `entity.isGrounded(coyoteTime)`
- **Arguments**:
  - `entity` (handle): The character entity.
  - `coyoteTime` (float, optional): Extra grace period in seconds after leaving a ledge.
- **Returns**: (bool) `TRUE` if grounded or within coyote time.

---

### `CHAR.TELEPORT(entity, x, y, z)` 

Instantly moves the character to a new position and resets all velocities. Use this for spawning or warp points.

- *Handle shortcut*: `entity.teleport(x, y, z)`
- **Arguments**:
  - `entity` (handle): The character entity.
  - `x, y, z` (float): Target world position.
- **Returns**: (none)
- **Alias**: `PLAYER.TELEPORT`

---

## Full Example: Player Controller

```basic
WINDOW.OPEN(1280, 720, "Character Controller")
PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -15, 0)

; Setup Camera
cam = CAMERA.CREATE().pos(0, 10, 20).look(0, 0, 0)

; Create Ground
floorDef = BODY3D.CREATE("STATIC")
BODY3D.ADDBOX(floorDef, 50, 0.5, 50)
BODY3D.COMMIT(floorDef, 0, -0.5, 0)

; Create Player
heroModel = MODEL.CREATECAPSULE(0.5, 2.0)
hero      = ENTITY.CREATE(heroModel)
hero.charCreate(0.5, 2.0).setStep(0.4)

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS3D.UPDATE()

    ; Basic WASD movement
    fwd  = INPUT.AXIS(KEY_S, KEY_W)
    side = INPUT.AXIS(KEY_A, KEY_D)
    hero.moveWithCam(cam, fwd, side, 10.0)

    ; Jump logic with coyote time
    IF hero.isGrounded(0.1) AND INPUT.KEYPRESSED(KEY_SPACE) THEN
        hero.jump(8.0)
    ENDIF

    ; Update camera to follow player
    px, py, pz = hero.pos()
    cam.setPos(px, py + 5, pz + 10).look(px, py, pz)

    RENDER.CLEAR(30, 30, 30)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW.GRID(100, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [PHYSICS3D.md](PHYSICS3D.md) — World setup and rigid bodies.
- [ENTITY.md](ENTITY.md) — Entity system and models.
- [PLAYER.md](PLAYER.md) — Advanced player stats and navigation.
