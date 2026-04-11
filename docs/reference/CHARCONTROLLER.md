# Character Controller Commands

Commands for creating and managing a kinematic character controller for 3D worlds. This provides a way to handle player movement that is driven by input rather than physics forces, while still respecting the collision geometry of the world.

For an **entity-based** wrapper (**`PLAYER.CREATE`**, **`PLAYER.MOVE`**, look targets, tag queries), see [PLAYER.md](PLAYER.md).

## Core Workflow

1.  **Start Physics**: The character controller relies on the 3D physics world. Start it with `Physics3D.Start()`.
2.  **Create Controller**: Use `CharController.Make()` to create the controller, defining its shape and initial position.
3.  **Update**: In the main loop, get user input and use `CharController.Move()` to update the controller's position.
4.  **Synchronize**: Use `CharController.GetPos()` or the `X/Y/Z` commands to sync your visual model with the controller.

---

### `CharController.Make(radius, height, x, y, z)`
Creates a new virtual character controller with a capsule shape at the specified world position. Returns a **controller handle**.

### `CharController.Free(handle)`
Frees the character controller resource and releases its heap slot.

---

### `CharController.Move(handle, dx, dy, dz)`
Updates the character's position based on a desired displacement vector. The controller automatically handles collisions with the physics world.

### `CharController.IsGrounded(handle)`
Returns `TRUE` if the character controller is currently standing on a surface (floor).

---

### `CharController.GetPos(handle)`
Returns a 3-float array handle `[x, y, z]` representing the controller's current world position.

### `CharController.X(handle)` / `CharController.Y(handle)` / `CharController.Z(handle)`
Returns the individual world coordinate component of the controller's position.

---

## Full Example

```basic
Window.Open(960, 540, "Character Controller")
Window.SetFPS(60)

; 1. Start Physics
Physics3D.Start()
Physics3D.SetGravity(0, -10, 0)

; Setup camera and floor
cam = Camera.Make()
cam.SetTarget(0, 5, 0)
floor_def = Body3D.Make("static")
Body3D.AddBox(floor_def, 100, 1, 100)
floor_body = Body3D.Commit(floor_def, 0, 0, 0)
floor_mesh = Mesh.MakeCube(100, 1, 100)
mat = Material.MakeDefault()

; 2. Create Controller
player = CharController.Make(0.5, 2.0, 0, 5, 0)
player_mesh = Mesh.MakeCapsule(0.5, 2.0, 16, 16)

WHILE NOT Window.ShouldClose()
    Physics3D.Step()

    ; 3. Update controller from input
    speed = 5.0 * Time.Delta()
    dx = 0
    dz = 0
    IF Input.KeyDown(KEY_W) THEN dz = -speed
    IF Input.KeyDown(KEY_S) THEN dz = speed
    IF Input.KeyDown(KEY_A) THEN dx = -speed
    IF Input.KeyDown(KEY_D) THEN dx = speed
    CharController.Move(player, dx, 0, dz)

    ; 4. Synchronize visuals
    player_x = CharController.X(player)
    player_y = CharController.Y(player)
    player_z = CharController.Z(player)
    cam.SetPos(player_x, player_y + 10, player_z + 15)
    cam.SetTarget(player_x, player_y, player_z)

    player_transform = Transform.Translation(player_x, player_y, player_z)

    Render.Clear(20, 30, 40)
    cam.Begin()
        Mesh.Draw(floor_mesh, mat, Body3D.GetMatrix(floor_body))
        Mesh.Draw(player_mesh, mat, player_transform)
        Draw.Grid(100, 1.0)
    cam.End()
    Render.Frame()
WEND

CharController.Free(player)
Physics3D.Stop()
Window.Close()
```
