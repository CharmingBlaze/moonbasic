# MoonBASIC Easy Mode Guide

MoonBASIC Easy Mode provides shorthands and property-style methods to make game development as fast and ergonomic as BlitzBasic.

## 1. Global Shorthands

These commands are available globally and act as shortcuts to the standard MoonBASIC modules.

| Easy Mode | Canonical MoonBASIC | Description |
|-----------|----------------------|-------------|
| `Graphics(w, h)` | `WINDOW.OPEN(w, h, "moonBASIC")` | Opens a game window. |
| `Graphics(w, h, title)` | `WINDOW.OPEN(w, h, title)` | Opens a game window with a title. |
| `PositionEntity(ent, x, y, z)` | `ENTITY.POSITIONENTITY(ent, x, y, z)` | Set an entity's absolute position (optional global flag on canonical API). |
| `RotateEntity(ent, p, y, r)` | `ENTITY.ROTATEENTITY(ent, p, y, r)` | Set an entity's absolute rotation (pitch, yaw, roll). |
| `MoveEntity(ent, f, r, u)` | `ENTITY.MOVE(ent, f, r, u)`, `MOVEENTITY` | Move along **local** forward, right, and up by **`f`**, **`r`**, **`u`** (from entity yaw/pitch). |
| `TranslateEntity(ent, dx, dy, dz)` | `ENTITY.TRANSLATE`, `ENTITY.TRANSLATEENTITY` | **World-space** delta **`(dx, dy, dz)`**; use for offsets that ignore entity facing. |
| `TFormVector(x, y, z, srcEnt, dstEnt)` | `ENTITY.TFORMVECTOR` | Transform direction **`(x,y,z)`** from **`srcEnt`** local space to **`dstEnt`** local space; returns **3-float array handle**. |
| `EntityHitsType(ent, type)` | (wrapper over `ENTITYCOLLIDED`) | **`TRUE`** if **`ent`** hit any other entity whose **`EntityType`** is **`type`** this frame (after **`ENTITY.UPDATE`** / **`UPDATEPHYSICS`**). |
| `EntityColor(obj, r, g, b)` | `ENTITY.COLOR(obj, r, g, b)` | Set an entity's color. |
| `EntityAlpha(obj, a)` | `ENTITY.ALPHA(obj, a)` | Set an entity's alpha transparency (0-1). |
| `FreeEntity(obj)` | `ENTITY.FREE(obj)` | Free an entity's memory. |
| `CreateCamera()` | `CAMERA.CREATE3D` | Create a standard 3D camera. |
| `TurnCamera(cam, p, y, r)` | `CAMERA.TURN(cam, p, y, r)` | Rotate camera relative to orientation. |
| `ShakeCamera(cam, i, d)` | `CAMERA.SHAKE(cam, intensity, duration)` | Shake the camera. |
| `CreateCamera2D()` | `CAMERA2D.MAKE()` | Create a 2D camera. |
| `Camera2DZoom(cam, z)` | `CAMERA2D.SETZOOM(cam, zoom)` | Set 2D camera zoom level. |
| `KeyHit(k)` | `INPUT.KEYPRESSED(k)` | Check if a key was pressed this frame. |
| `KeyDown(k)` | `INPUT.KEYDOWN(k)` | Check if a key is held down. |
| `MouseX()` | `INPUT.MOUSEX()` | Get mouse X position. |
| `MouseY()` | `INPUT.MOUSEY()` | Get mouse Y position. |
| `MouseHit(b)` | `INPUT.MOUSEPRESSED(b)` | Check if a mouse button was clicked. |
| `Millisecs()` | `TIME.MILLIS()` | Get milliseconds since the engine started. |
| `UpdatePhysics()` | `UPDATEPHYSICS` | One frame tick: `ENTITY.UPDATE(Time.Delta)` + best-effort world / 2D / 3D physics steps. |

## 2. Property-Style Handle Methods

Most engine handles now support unified property shorthands for easier manipulation within the game loop.

### 3D Cameras
```basic
cam = CreateCamera()
cam.pos(10, 10, 10)
cam.look(0, 0, 0)
cam.turn(0, 1, 0)     ' Turn camera
cam.zoom(1.5)         ' Set FOV/Zoom
cam.shake(1.0, 0.5)   ' Shake camera
```

### 2D Cameras
```basic
cam2d = CreateCamera2D()
cam2d.target(100, 100) ' Set target/position
cam2d.zoom(2.0)        ' Set zoom
cam2d.rot(45)          ' Set rotation
```

### Universal Methods:
- `.pos(x, y, z)` - Set position
- `.rot(p, y, r)` or `.rot(a)` - Set rotation
- `.scale(sx, sy, sz)` - Set scale
- `.size(w, h, d)` - Set dimensions
- `.col(r, g, b)` - Set color (0-255)
- `.alpha(a)` - Set alpha (0-1)
- `.free()` - Free handle memory


## 3. Physics & Networking (Extended)

### Physics (Jolt/Box2D)
```basic
body = CreateBody(TYPE_DYNAMIC, SHAPE_BOX)
body.pos(0, 10, 0)
body.force(0, -10, 0)  ' Apply force
body.vel(0, 1, 0)      ' Set linear velocity
```

### Networking (ENet)
```basic
' Server
server = Listen(1234)
If ServiceNet(server, 10) Then
    msg$ = NetMsg$()
    Print "Received: " + msg$
End If

' Client
client = Connect("127.0.0.1", 1234)
client.send("Hello World")
```

### Audio
```basic
snd = LoadSound("boom.wav")
snd.play()
snd.volume(0.5)

mus = LoadMusic("theme.ogg")
mus.play()
mus.pitch(1.2)
```
