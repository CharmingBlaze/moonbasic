# moonBASIC API consistency

This document is generated from `compiler/builtinmanifest/commands.json`.

**Contributor contract:** Treat this file as the authoritative checklist of **registered overloads** (name, arity, and manifest metadata). New builtins belong in **`compiler/builtinmanifest/commands.json`**; refresh this doc after manifest edits so tooling, reviews, and external contributors stay aligned.

Refresh: `go run ./tools/apidoc` (from the repository root).

## Related documentation

- **[ERROR_MESSAGES.md](../ERROR_MESSAGES.md)** — compile-time vs runtime errors, did-you-mean, heap handle hints.
- **[ROADMAP.md](../ROADMAP.md)** — phased engineering plan (polish → rendering → 2D → systems → …).
- **[COMMAND_AUDIT.md](../COMMAND_AUDIT.md)** — namespace → primary `docs/reference/*.md` file; run **`go run ./tools/cmdaudit`** to verify every manifest namespace maps to an existing reference page (exit code **2** if a namespace is unmapped or a referenced file is missing).
- **[reference/API_CONVENTIONS.md](../reference/API_CONVENTIONS.md)** — consistent verbs (`LOAD`, `SETPOS`, `SETSCALE`, …) across object types.

## Naming conventions

- **Registry / source form**: `NS.ACTION` in uppercase with a dot (e.g. `CAMERA.SETPOS`).
- **Handle methods** (on a handle value): `cam.SetPos` dispatches to `CAMERA.SETPOS`. **`SetPosition`** is an alias for **`SetPos`** where both are registered (same handler).
- **Spatial handles** (`Camera3D`, `Body3D`, `Model`, `Sprite`, `Light2D`): use **`SETPOS`** for position. Aliases **`SETPOSITION`** exist for **Camera**, **Model**, **Body3D**, **Sprite**, **Light2D** — same implementation as `SETPOS`.
- **3D lights** (`LIGHT.*`): use **`LIGHT.SETDIR`** for the directional sun (normalized). **`LIGHT.SETPOS`** stores point/spot position; **`LIGHT.SETTARGET`** moves the shadow frustum look-at; **`RENDER.SETAMBIENT`** sets PBR ambient tint.
- **`MODEL.SETPOS`**: sets the model root transform to a **translation matrix** (replaces prior rotation/scale on that matrix).
- **Creation verbs**: prefer **`*.CREATE`** for procedural handles; deprecated **`*.MAKE`** aliases point at the same handlers where registered. **`*.LOAD`** for assets (`SPRITE.LOAD`, `MODEL.LOAD`); materials use `MATERIAL.MAKEDEFAULT` / `MATERIAL.MAKEPBR`.
- **Cross-type patterns**: see **[API_CONVENTIONS.md](../reference/API_CONVENTIONS.md)**.

## Default values (common no-arg `CREATE` paths)

| Command | Defaults |
|----------|----------|
| `CAMERA.CREATE` (deprecated `CAMERA.MAKE`) | position (0, 2, 8), target (0, 0, 0), up (0, 1, 0), FOV 45°, perspective |
| `LIGHT.CREATE` (deprecated `LIGHT.MAKE`) | kind `directional`, white, intensity 1.0, direction toward normalized (-1,-2,-1) |
| `BODY3D.CREATE` (deprecated `BODY3D.MAKE`) | no args → **DYNAMIC** motion type |
| `MATERIAL.MAKEDEFAULT` / `MAKEPBR` | see `runtime/mbmodel3d` (material modules) |

## Debug watch overlay

`DEBUG.WATCH(label$, value)` stores rows; `DEBUG.WATCHCLEAR` clears them. With **CGO** and Raylib, the window pipeline may draw a **top-left overlay** each frame (`runtime/mbdebug/overlay_cgo.go`) when **`DEBUG.ENABLE`** was called or the host enabled **`Registry.DebugMode`** (e.g. **`--debug`**). **`DEBUG.DISABLE`** clears the user override. Without CGO, watches are stored but not drawn.

## Errors

- **Compile-time**: unknown `NS.METHOD` → did-you-mean within namespace + manifest listing (see `compiler/semantic/cmdhint.go`).
- **Runtime**: VM wraps native errors with **source file and line** when available (`vm/vm.go`). Unknown registry keys → `runtime.FormatUnknownRegistryCommand`.

## Commands by namespace

### ABS

- **`ABS`** - args: any

### ACOS

- **`ACOS`** - args: any

### ADDFORCE

- **`ADDFORCE`** - args: handle, float, float, float — Easy Mode: Body.AddForce(x, y, z)

### ADDIMPULSE

- **`ADDIMPULSE`** - args: handle, float, float, float — Easy Mode: Body.AddImpulse(x, y, z)

### AMBIENTLIGHT

- **`AMBIENTLIGHT`** - args: int, int, int, float -> returns void — Easy Mode: Set global ambient light (r, g, b, intensity)

### ANGLE

- **`ANGLE.DIFFERENCE`** - args: float, float -> returns float — Shortest signed angle from a to b in degrees (alias of MATH.ANGLEDIFF)

### ANGLEDIFF

- **`ANGLEDIFF`** - args: any, any

### ANGLEDIFFRAD

- **`ANGLEDIFFRAD`** - args: float, float -> returns float — Shortest signed angle difference b-a in radians

### ANGLETO

- **`ANGLETO`** - args: float, float, float, float -> returns float — Heading in degrees [0,360) on XZ from (x1,z1) to (x2,z2)

### ANIM

- **`ANIM.ADDTRANSITION`** - args: handle, string, string, string
- **`ANIM.DEFINE`** - args: handle, string, int, int, float, bool
- **`ANIM.SETPARAM`** - args: handle, string, any
- **`ANIM.UPDATE`** - args: handle, float

### ARGC

- **`ARGC`** - args: (none)

### ARRAYCONTAINS

- **`ARRAYCONTAINS`** - args: handle, any -> returns bool

### ARRAYCOPY

- **`ARRAYCOPY`** - args: handle, handle

### ARRAYFILL

- **`ARRAYFILL`** - args: handle, any

### ARRAYFIND

- **`ARRAYFIND`** - args: handle, any -> returns int

### ARRAYFREE

- **`ARRAYFREE`** - args: handle

### ARRAYJOINS

- **`ARRAYJOINS`** - args: handle, string -> returns string

### ARRAYJOINS$

- **`ARRAYJOINS$`** - args: handle, string -> returns string

### ARRAYLEN

- **`ARRAYLEN`** - args: handle -> returns int

### ARRAYPOP

- **`ARRAYPOP`** - args: handle -> returns any

### ARRAYPUSH

- **`ARRAYPUSH`** - args: handle, any

### ARRAYREVERSE

- **`ARRAYREVERSE`** - args: handle

### ARRAYSHIFT

- **`ARRAYSHIFT`** - args: handle -> returns any

### ARRAYSLICE

- **`ARRAYSLICE`** - args: handle, int, int -> returns handle

### ARRAYSORT

- **`ARRAYSORT`** - args: handle

### ARRAYSPLICE

- **`ARRAYSPLICE`** - args: handle, int, int

### ARRAYUNSHIFT

- **`ARRAYUNSHIFT`** - args: handle, any

### ASC

- **`ASC`** - args: string

### ASIN

- **`ASIN`** - args: any

### ASSERT

- **`ASSERT`** - args: any, string

### ATAN

- **`ATAN`** - args: any

### ATAN2

- **`ATAN2`** - args: any, any

### ATLAS

- **`ATLAS.FREE`** - args: handle
- **`ATLAS.GETSPRITE`** - args: handle, string -> returns handle
- **`ATLAS.LOAD`** - args: string, string -> returns handle

### ATN

- **`ATN`** - args: any

### AUDIO

- **`AUDIO.CLOSE`** - args: (none)
- **`AUDIO.GETMUSICLENGTH`** - args: handle -> returns float
- **`AUDIO.GETMUSICTIME`** - args: handle -> returns float
- **`AUDIO.INIT`** - args: (none)
- **`AUDIO.ISMUSICPLAYING`** - args: handle -> returns bool
- **`AUDIO.ISSOUNDPLAYING`** - args: handle -> returns bool
- **`AUDIO.LISTENERCAMERA`** - args: handle
- **`AUDIO.LOADMUSIC`** - args: string -> returns handle
- **`AUDIO.LOADSOUND`** - args: string -> returns handle
- **`AUDIO.PAUSE`** - args: handle
- **`AUDIO.PLAY`** - args: handle
- **`AUDIO.RESUME`** - args: handle
- **`AUDIO.SEEKMUSIC`** - args: handle, float
- **`AUDIO.SETMASTERVOLUME`** - args: float
- **`AUDIO.SETMUSICPITCH`** - args: handle, float
- **`AUDIO.SETMUSICVOLUME`** - args: handle, float
- **`AUDIO.SETSOUNDPAN`** - args: handle, float
- **`AUDIO.SETSOUNDPITCH`** - args: handle, float
- **`AUDIO.SETSOUNDVOLUME`** - args: handle, float
- **`AUDIO.STOP`** - args: handle
- **`AUDIO.UPDATEMUSIC`** - args: handle

### AUDIOSTREAM

- **`AUDIOSTREAM.CREATE`** - args: int, int, int -> returns handle
- **`AUDIOSTREAM.FREE`** - args: handle
- **`AUDIOSTREAM.ISPLAYING`** - args: handle -> returns bool
- **`AUDIOSTREAM.ISREADY`** - args: handle -> returns bool
- **`AUDIOSTREAM.MAKE`** - args: int, int, int -> returns handle — DEPRECATED alias of AUDIOSTREAM.CREATE. Use AUDIOSTREAM.CREATE.
- **`AUDIOSTREAM.PAUSE`** - args: handle
- **`AUDIOSTREAM.PLAY`** - args: handle
- **`AUDIOSTREAM.RESUME`** - args: handle
- **`AUDIOSTREAM.SETPAN`** - args: handle, float
- **`AUDIOSTREAM.SETPITCH`** - args: handle, float
- **`AUDIOSTREAM.SETVOLUME`** - args: handle, float
- **`AUDIOSTREAM.STOP`** - args: handle
- **`AUDIOSTREAM.UPDATE`** - args: handle, handle

### AXIS

- **`AXIS`** - args: any, any -> returns float — Easy Mode: INPUT.AXIS(INPUT(), k1, k2)
- **`AXIS`** - args: int, int -> returns float

### ActiveShader

- **`ActiveShader`** - args: handle — Alias of POST.ADDSHADER — full-screen post shader for the render pipeline

### AddTriangle

- **`AddTriangle`** - args: handle, int, int, int

### AddVertex

- **`AddVertex`** - args: handle, float, float, float -> returns int

### AddWheel

- **`AddWheel`** - args: any

### Animate

- **`Animate`** - args: int, any, any

### AppTitle

- **`AppTitle`** - args: string — Alias of WINDOW.SETTITLE

### ApplyEntityImpulse

- **`ApplyEntityImpulse`** - args: int, float, float, float

### BALL

- **`BALL`** - args: float, float, float, float, int, int, int, int — alias of DRAW3D.SPHERE — solid sphere

### BALLW

- **`BALLW`** - args: float, float, float, float, int, int, int, int, int, int — alias of DRAW3D.SPHEREWIRES — wire sphere

### BAND

- **`BAND`** - args: any, any

### BBOX

- **`BBOX.CHECK`** - args: handle, handle -> returns bool
- **`BBOX.CHECKSPHERE`** - args: handle, float, float, float, float -> returns bool
- **`BBOX.CREATE`** - args: float, float, float, float, float, float -> returns handle
- **`BBOX.FREE`** - args: handle
- **`BBOX.FROMMODEL`** - args: handle -> returns handle
- **`BBOX.MAKE`** - args: float, float, float, float, float, float -> returns handle — DEPRECATED alias of BBOX.CREATE. Use BBOX.CREATE.

### BCLEAR

- **`BCLEAR`** - args: any, int

### BCOUNT

- **`BCOUNT`** - args: any

### BIN

- **`BIN`** - args: int -> returns string

### BIN$

- **`BIN$`** - args: int -> returns string

### BIOME

- **`BIOME.CREATE`** - args: string -> returns handle — Create a named biome state handle (temperature/humidity via BIOME.SET*).
- **`BIOME.FREE`** - args: handle — Release biome handle.
- **`BIOME.MAKE`** - args: string -> returns handle — DEPRECATED alias of BIOME.CREATE. Use BIOME.CREATE.
- **`BIOME.SETHUMIDITY`** - args: handle, float — Set biome humidity (0..1).
- **`BIOME.SETTEMP`** - args: handle, float — Set biome temperature (celsius).

### BLSHIFT

- **`BLSHIFT`** - args: any, int

### BNOT

- **`BNOT`** - args: any

### BODY

- **`BODY3D.LOCKAXIS`** - args: handle, int — Lock motion/rotation axes (flags: 1=X, 2=Y, 4=Z, 8=RotX, 16=RotY, 32=RotZ).
- **`BODY3D.SETCCD`** - args: handle, bool — Enable/disable Continuous Collision Detection.
- **`BODY3D.SETDAMPING`** - args: handle, float, float — Set linear and angular damping.
- **`BODY3D.SETGRAVITYFACTOR`** - args: handle, float — Set gravity multiplier (0.0 = weightless).

### BODY2D

- **`BODY2D.ADDCIRCLE`** - args: handle, float
- **`BODY2D.ADDPOLYGON`** - args: handle, handle
- **`BODY2D.ADDRECT`** - args: handle, float, float
- **`BODY2D.APPLYFORCE`** - args: handle, float, float
- **`BODY2D.APPLYIMPULSE`** - args: handle, float, float
- **`BODY2D.COLLIDED`** - args: handle -> returns int
- **`BODY2D.COLLISIONNORMAL`** - args: handle -> returns handle
- **`BODY2D.COLLISIONOTHER`** - args: handle -> returns handle
- **`BODY2D.COLLISIONPOINT`** - args: handle -> returns handle
- **`BODY2D.COMMIT`** - args: handle, float, float -> returns handle
- **`BODY2D.CREATE`** - args: string -> returns handle
- **`BODY2D.FREE`** - args: handle
- **`BODY2D.GETPOS`** - args: handle -> returns handle
- **`BODY2D.GETROT`** - args: handle -> returns float
- **`BODY2D.MAKE`** - args: string -> returns handle — DEPRECATED alias of BODY2D.CREATE. Use BODY2D.CREATE.
- **`BODY2D.ROT`** - args: handle -> returns float
- **`BODY2D.SETANGULARVELOCITY`** - args: handle, float
- **`BODY2D.SETFRICTION`** - args: handle, float
- **`BODY2D.SETLINEARVELOCITY`** - args: handle, float, float
- **`BODY2D.SETMASS`** - args: handle, float
- **`BODY2D.SETPOS`** - args: handle, float, float
- **`BODY2D.SETPOSITION`** - args: handle, float, float — DEPRECATED alias of BODY2D.SETPOS. Use BODY2D.SETPOS.
- **`BODY2D.SETRESTITUTION`** - args: handle, float
- **`BODY2D.SETROT`** - args: handle, float
- **`BODY2D.X`** - args: handle -> returns float
- **`BODY2D.Y`** - args: handle -> returns float

### BODY3D

- **`BODY3D.ACTIVATE`** - args: handle
- **`BODY3D.ADDBOX`** - args: handle, float, float, float
- **`BODY3D.ADDCAPSULE`** - args: handle, float, float
- **`BODY3D.ADDMESH`** - args: handle, handle
- **`BODY3D.ADDSPHERE`** - args: handle, float
- **`BODY3D.APPLYFORCE`** - args: handle, float, float, float
- **`BODY3D.APPLYIMPULSE`** - args: handle, float, float, float
- **`BODY3D.BUFFERINDEX`** - args: handle -> returns int
- **`BODY3D.COLLIDED`** - args: handle -> returns int
- **`BODY3D.COLLISIONNORMAL`** - args: handle -> returns handle
- **`BODY3D.COLLISIONOTHER`** - args: handle -> returns handle
- **`BODY3D.COLLISIONPOINT`** - args: handle -> returns handle
- **`BODY3D.COMMIT`** - args: handle, float, float, float -> returns handle
- **`BODY3D.CREATE`** - args: (none) -> returns handle
- **`BODY3D.CREATE`** - args: string -> returns handle
- **`BODY3D.CREATE`** - args: string
- **`BODY3D.DEACTIVATE`** - args: handle
- **`BODY3D.FREE`** - args: handle
- **`BODY3D.GETPOS`** - args: handle -> returns handle
- **`BODY3D.GETROT`** - args: handle -> returns handle
- **`BODY3D.GETSCALE`** - args: handle -> returns handle — Returns [sx,sy,sz] scale factors for primitive bodies (box/sphere/capsule); mesh bodies report 1,1,1
- **`BODY3D.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of BODY3D.CREATE. Use BODY3D.CREATE.
- **`BODY3D.MAKE`** - args: string — DEPRECATED alias of BODY3D.CREATE. Use BODY3D.CREATE.
- **`BODY3D.MAKE`** - args: string -> returns handle — DEPRECATED alias of BODY3D.CREATE. Use BODY3D.CREATE.
- **`BODY3D.SETANGULARVEL`** - args: handle, float, float, float
- **`BODY3D.SETFRICTION`** - args: handle, float
- **`BODY3D.SETLINEARVEL`** - args: handle, float, float, float
- **`BODY3D.SETMASS`** - args: handle, float
- **`BODY3D.SETPOS`** - args: handle, float, float, float
- **`BODY3D.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of BODY3D.SETPOS. Use BODY3D.SETPOS.
- **`BODY3D.SETRESTITUTION`** - args: handle, float
- **`BODY3D.SETROT`** - args: handle, float, float, float
- **`BODY3D.SETSCALE`** - args: handle, float, float, float — Scales collision shape for primitive bodies built via ADDBOX/ADDSPHERE/ADDCAPSULE or SHAPE.CREATE*; not supported for mesh (ADDMESH)
- **`BODY3D.X`** - args: handle -> returns float
- **`BODY3D.Y`** - args: handle -> returns float
- **`BODY3D.Z`** - args: handle -> returns float

### BODYREF

- **`BODYREF.ENABLECOLLISION`** - args: handle, bool — Enables/Disables body participation in physics.
- **`BODYREF.SETLAYER`** - args: handle, int — Sets the Jolt collision layer.
- **`BODYREF.SETPOS`** - args: handle, float, float, float — Moves a Kinematic/Static/Trigger body.
- **`BODYREF.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of BODYREF.SETPOS. Use BODYREF.SETPOS.
- **`BODYREF.SETROTATION`** - args: handle, float, float, float — Sets body orientation (Euler degrees).

### BOOL

- **`BOOL`** - args: any -> returns bool

### BOR

- **`BOR`** - args: any, any

### BOX

- **`BOX`** - args: float, float, float, float, float, float, int, int, int, int — alias of DRAW3D.CUBE — solid axis-aligned box

### BOX2D

- **`BOX2D.BODYCREATE`** - args: float, float, int
- **`BOX2D.FIXTUREBOX`** - args: float, float, float, float
- **`BOX2D.FIXTURECIRCLE`** - args: float
- **`BOX2D.WORLDCREATE`** - args: float, float
- **`BOX2D.WORLDSTEP`** - args: float, int, int

### BOXTOPLAND

- **`BOXTOPLAND`** - args: float, float, float, float, float, float, float, float, float, float, float -> returns float — Sphere vs box top: landing centre Y or 0.0 if no landing

### BOXW

- **`BOXW`** - args: float, float, float, float, float, float, int, int, int, int — alias of DRAW3D.CUBEWIRES — wire box

### BRSHIFT

- **`BRSHIFT`** - args: any, int

### BSET

- **`BSET`** - args: any, int

### BSPHERE

- **`BSPHERE.CHECK`** - args: handle, handle -> returns bool
- **`BSPHERE.CHECKBOX`** - args: handle, handle -> returns bool
- **`BSPHERE.CREATE`** - args: float, float, float, float -> returns handle
- **`BSPHERE.FREE`** - args: handle
- **`BSPHERE.MAKE`** - args: float, float, float, float -> returns handle — DEPRECATED alias of BSPHERE.CREATE. Use BSPHERE.CREATE.

### BTEST

- **`BTEST`** - args: any, int

### BTOGGLE

- **`BTOGGLE`** - args: any, int

### BTREE

- **`BTREE.ADDACTION`** - args: handle, string
- **`BTREE.ADDCONDITION`** - args: handle, string
- **`BTREE.CREATE`** - args: (none) -> returns handle
- **`BTREE.FREE`** - args: handle
- **`BTREE.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of BTREE.CREATE. Use BTREE.CREATE.
- **`BTREE.RUN`** - args: handle, handle, float
- **`BTREE.SEQUENCE`** - args: handle -> returns handle

### BXOR

- **`BXOR`** - args: any, any

### BrushFX

- **`BrushFX`** - args: handle, any

### BrushShininess

- **`BrushShininess`** - args: handle, float

### BrushTexture

- **`BrushTexture`** - args: handle, handle, any

### CAM

- **`CAM`** - args: (none) -> returns handle — Short Blitz-style camera alias; same as CAMERA.CREATE (deprecated name: CAMERA.MAKE)

### CAMERA

- **`CAMERA.BEGIN`** - args: handle
- **`CAMERA.CAMERAFOLLOW`** - args: handle, int, float, float, float
- **`CAMERA.CLEARFPSMODE`** - args: handle
- **`CAMERA.CREATE`** - args: (none)
- **`CAMERA.CREATE`** - args: (none) -> returns handle — Returns a Camera3D heap handle (canonical; deprecated alias: CAMERA.MAKE)
- **`CAMERA.END`** - args: (none)
- **`CAMERA.END`** - args: handle
- **`CAMERA.FOLLOW`** - args: handle, handle, float, float — Spring math camera tracker.
- **`CAMERA.FOLLOW`** - args: handle, float, float, float, float, float, float, float
- **`CAMERA.FOLLOWENTITY`** - args: handle, int, float, float, float
- **`CAMERA.FREE`** - args: handle
- **`CAMERA.GETACTIVE`** - args: (none) -> returns handle
- **`CAMERA.GETMATRIX`** - args: handle -> returns handle
- **`CAMERA.GETPOS`** - args: handle -> returns handle
- **`CAMERA.GETRAY`** - args: handle, float, float
- **`CAMERA.GETROT`** - args: handle -> returns handle
- **`CAMERA.GETTARGET`** - args: handle -> returns handle
- **`CAMERA.GETVIEWRAY`** - args: float, float, handle, int, int
- **`CAMERA.GETYAW`** - args: handle — Alias of CAMERA.YAW.
- **`CAMERA.ISONSCREEN`** - args: handle, float, float, float -> returns bool
- **`CAMERA.ISONSCREEN`** - args: handle, float, float, float, float -> returns bool
- **`CAMERA.LERPTO`** - args: handle, int, float — Smoothly interpolate camera target toward an entity.
- **`CAMERA.LOOKAT`** - args: handle, float, float, float
- **`CAMERA.LOOKATENTITY`** - args: handle, int — Sets camera target to entity world position (same idea as Blitz PointAt)
- **`CAMERA.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of CAMERA.CREATE. Returns a Camera3D heap handle.
- **`CAMERA.MAKE`** - args: (none) — DEPRECATED alias of CAMERA.CREATE. Use CAMERA.CREATE.
- **`CAMERA.MOUSERAY`** - args: handle -> returns handle
- **`CAMERA.MOVE`** - args: handle, float, float, float
- **`CAMERA.ORBIT`** - args: handle, float, float, float, float, float, float
- **`CAMERA.ORBITAROUND`** - args: handle, float, float, float, float, float, float
- **`CAMERA.ORBITAROUNDEG`** - args: handle, float, float, float, float, float, float
- **`CAMERA.ORBITCAMERA`** - args: handle, float, float, float -> returns float
- **`CAMERA.ORBITENTITY`** - args: handle, int, float, float, float
- **`CAMERA.PICK`** - args: handle, float, float -> returns handle
- **`CAMERA.POINTATENTITY`** - args: handle, int — Alias of CAMERA.LOOKATENTITY
- **`CAMERA.PROJECT`** - args: handle, float, float, float -> returns handle — Alias of CAMERA.WORLDTOSCREEN — world point to screen [sx,sy]
- **`CAMERA.RAYCASTMOUSE`** - args: handle -> returns int — Raycast from mouse through camera; returns entity id or 0.
- **`CAMERA.ROTATE`** - args: handle, float, float, float
- **`CAMERA.SETACTIVE`** - args: handle
- **`CAMERA.SETFOV`** - args: handle, float
- **`CAMERA.SETFPSMODE`** - args: handle, float
- **`CAMERA.SETMODE`** - args: handle, any — 0/1 or perspective/orthographic — alias-friendly CAMERA.SETPROJECTION
- **`CAMERA.SETORBIT`** - args: handle, float, float, float, float, float, float
- **`CAMERA.SETORBITKEYS`** - args: handle, float, float — Raylib key codes for orbit yaw (0 disables that side).
- **`CAMERA.SETORBITKEYSPEED`** - args: handle, float — Keyboard orbit yaw rate in radians per second.
- **`CAMERA.SETORBITLIMITS`** - args: handle, float, float, float, float — Clamp pitch (radians) and orbit distance for CAMERA.ORBIT (entity).
- **`CAMERA.SETORBITSPEED`** - args: handle, float, float — Mouse drag sensitivity and mouse wheel zoom scale for orbit-follow.
- **`CAMERA.SETPOS`** - args: handle, float, float, float
- **`CAMERA.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of CAMERA.SETPOS. Use CAMERA.SETPOS.
- **`CAMERA.SETPROJECTION`** - args: handle, int
- **`CAMERA.SETRANGE`** - args: handle, float, float
- **`CAMERA.SETTARGET`** - args: handle, float, float, float
- **`CAMERA.SETTARGETENTITY`** - args: handle, int
- **`CAMERA.SETUP`** - args: handle, float, float, float
- **`CAMERA.SHAKE`** - args: handle, float, float
- **`CAMERA.SMOOTHEXP`** - args: float, float, float, float -> returns float — Exponential smoothing: current toward target using (1-exp(-smoothHz*dt)); for orbit angles
- **`CAMERA.TURN`** - args: handle, float, float, float
- **`CAMERA.TURNLEFT`** - args: handle, float -> returns float
- **`CAMERA.TURNRIGHT`** - args: handle, float -> returns float
- **`CAMERA.UNPROJECT`** - args: handle, float, float -> returns handle — Screen (x,y) to world ray — alias of CAMERA.GETRAY / PICK
- **`CAMERA.UPDATEFPS`** - args: handle
- **`CAMERA.USEMOUSEORBIT`** - args: handle, bool — Enable/disable mouse contribution to CAMERA.ORBIT (entity) orbit-follow.
- **`CAMERA.USEORBITRIGHTMOUSE`** - args: handle, bool — If true (default), mouse orbit only while right button is held; if false, mouse moves orbit without RMB.
- **`CAMERA.WORLDTOSCREEN`** - args: handle, float, float, float -> returns handle
- **`CAMERA.WORLDTOSCREEN2D`** - args: handle, float, float, float -> returns handle
- **`CAMERA.XZBASIS`** - args: handle -> returns array — Returns planar [fwdX, fwdZ, rightX, rightZ] vectors for camera-relative movement.
- **`CAMERA.YAW`** - args: handle — Orbit yaw in radians (internal state) for aligning entities with cam.Orbit(entity, dist).
- **`CAMERA.ZOOM`** - args: handle, float

### CAMERA2D

- **`CAMERA2D.BEGIN`** - args: (none)
- **`CAMERA2D.BEGIN`** - args: handle
- **`CAMERA2D.CREATE`** - args: (none) -> returns handle
- **`CAMERA2D.END`** - args: (none)
- **`CAMERA2D.FOLLOW`** - args: handle, handle, float, float
- **`CAMERA2D.FREE`** - args: handle
- **`CAMERA2D.GETMATRIX`** - args: handle -> returns handle
- **`CAMERA2D.GETPOS`** - args: handle -> returns array
- **`CAMERA2D.GETROTATION`** - args: handle -> returns float
- **`CAMERA2D.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of CAMERA2D.CREATE. Use CAMERA2D.CREATE.
- **`CAMERA2D.ROTATION`** - args: handle -> returns float
- **`CAMERA2D.SCREENTOWORLD`** - args: handle, float, float -> returns handle
- **`CAMERA2D.SETOFFSET`** - args: handle, float, float
- **`CAMERA2D.SETROTATION`** - args: handle, float
- **`CAMERA2D.SETTARGET`** - args: handle, float, float
- **`CAMERA2D.SETZOOM`** - args: handle, float
- **`CAMERA2D.TARGETX`** - args: handle -> returns float
- **`CAMERA2D.TARGETY`** - args: handle -> returns float
- **`CAMERA2D.WORLDTOSCREEN`** - args: handle, float, float -> returns handle
- **`CAMERA2D.ZOOMIN`** - args: handle, float
- **`CAMERA2D.ZOOMOUT`** - args: handle, float
- **`CAMERA2D.ZOOMTOMOUSE`** - args: handle, float

### CAMERA2DOFFSET

- **`CAMERA2DOFFSET`** - args: handle, float, float — Easy Mode: CAMERA2D.SETOFFSET(cam, x, y)

### CAMERA2DROTATION

- **`CAMERA2DROTATION`** - args: handle, float — Easy Mode: CAMERA2D.SETROTATION(cam, r)

### CAMERA2DTARGET

- **`CAMERA2DTARGET`** - args: handle, float, float — Easy Mode: CAMERA2D.SETTARGET(cam, x, y)

### CAMERA2DZOOM

- **`CAMERA2DZOOM`** - args: handle, float — Easy Mode: CAMERA2D.SETZOOM(cam, z)

### CAMERAFOLLOW

- **`CAMERAFOLLOW`** - args: handle, int, float, float, float — Easy Mode: CAMERA.FOLLOWENTITY(cam, ent, dist, height, smooth)

### CAMERAPICK

- **`CAMERAPICK`** - args: handle, float, float -> returns handle — Easy Mode: CAMERA.PICK(cam, x, y)

### CAMERAZOOM

- **`CAMERAZOOM`** - args: handle, float — Easy Mode: CAMERA.ZOOM(cam, z)

### CAP

- **`CAP`** - args: float, float, float, float, float, float, float, int, int, int, int, int, int — alias of DRAW3D.CAPSULE — solid capsule

### CAPW

- **`CAPW`** - args: float, float, float, float, float, float, float, int, int, int, int, int, int — alias of DRAW3D.CAPSULEWIRES — wire capsule

### CEIL

- **`CEIL`** - args: any

### CHAR

- **`CHAR.CREATE`** - args: int — Alias of PLAYER.CREATE: (entity) or (entity, radius#, height#); allocates Jolt CharacterVirtual and clears scripted gravity/velocity for stable KCC (Linux+CGO)
- **`CHAR.CREATE`** - args: int, float, float — Alias of PLAYER.CREATE(entity, radius#, height#)
- **`CHAR.DIST`** - args: int, int -> returns float — Alias of ENTITY.DIST — distance between two entities
- **`CHAR.GETCEILING`** - args: int -> returns bool — Alias of PLAYER.GETCEILING
- **`CHAR.GETCOLLISIONENABLED`** - args: int -> returns bool
- **`CHAR.GETFRICTION`** - args: int -> returns float
- **`CHAR.GETGRAVITYSCALE`** - args: int -> returns float
- **`CHAR.GETGROUNDSTATE`** - args: int -> returns int — Alias of PLAYER.GETGROUNDSTATE
- **`CHAR.GETGROUNDVELOCITYX`** - args: int -> returns float
- **`CHAR.GETGROUNDVELOCITYY`** - args: int -> returns float
- **`CHAR.GETGROUNDVELOCITYZ`** - args: int -> returns float
- **`CHAR.GETHEIGHT`** - args: int -> returns float
- **`CHAR.GETISFALLING`** - args: int -> returns bool
- **`CHAR.GETISJUMPING`** - args: int -> returns bool
- **`CHAR.GETISSLIDING`** - args: int -> returns bool — Alias of PLAYER.GETISSLIDING
- **`CHAR.GETLAYER`** - args: int -> returns int
- **`CHAR.GETMASK`** - args: int -> returns int
- **`CHAR.GETMAXSLOPE`** - args: int -> returns float
- **`CHAR.GETONSLOPE`** - args: int -> returns bool
- **`CHAR.GETONWALL`** - args: int -> returns bool
- **`CHAR.GETPOSITIONX`** - args: int -> returns float — Alias of PLAYER.GETPOSITIONX
- **`CHAR.GETPOSITIONY`** - args: int -> returns float
- **`CHAR.GETPOSITIONZ`** - args: int -> returns float
- **`CHAR.GETRADIUS`** - args: int -> returns float
- **`CHAR.GETROTATIONPITCH`** - args: int -> returns float
- **`CHAR.GETROTATIONROLL`** - args: int -> returns float
- **`CHAR.GETROTATIONYAW`** - args: int -> returns float
- **`CHAR.GETSLOPEANGLE`** - args: int -> returns float
- **`CHAR.GETSNAPDISTANCE`** - args: int -> returns float
- **`CHAR.GETSPEED`** - args: int -> returns float
- **`CHAR.GETSTEPHEIGHT`** - args: int -> returns float
- **`CHAR.GETVELOCITYX`** - args: int -> returns float
- **`CHAR.GETVELOCITYY`** - args: int -> returns float
- **`CHAR.GETVELOCITYZ`** - args: int -> returns float
- **`CHAR.GETVX`** - args: int — Returns horizontal velocity X for kinematic controller
- **`CHAR.GETVY`** - args: int — Returns vertical velocity Y for kinematic controller
- **`CHAR.GETVZ`** - args: int — Returns horizontal velocity Z for kinematic controller
- **`CHAR.ISGROUNDED`** - args: int -> returns bool — Alias of PLAYER.ISGROUNDED
- **`CHAR.ISGROUNDED`** - args: int, float -> returns bool — KCC ground test with optional coyote grace (seconds)
- **`CHAR.ISONSTEEPSLOPE`** - args: int -> returns bool — Alias of PLAYER.ISONSTEEPSLOPE
- **`CHAR.JUMP`** - args: int, float — Alias of PLAYER.JUMP
- **`CHAR.MAKE`** - args: int — DEPRECATED alias of CHAR.CREATE. Use CHAR.CREATE. Alias of PLAYER.CREATE: (entity) or (entity, radius#, height#); allocates Jolt CharacterVirtual and clears scripted gravity/velocity for stable KCC (Linux+CGO)
- **`CHAR.MAKE`** - args: int, float, float — DEPRECATED alias of CHAR.CREATE. Use CHAR.CREATE. Alias of PLAYER.CREATE(entity, radius#, height#)
- **`CHAR.MOVE`** - args: int, float, float, float — KCC world move: (entity, dirX#, dirZ#, speed#) → horizontal velocity = dir * speed; slides on walls (CharacterVirtual; Linux+CGO)
- **`CHAR.MOVEWITHCAM`** - args: int, handle, float, float, float — Alias of CHAR.MOVEWITHCAMERA / PLAYER.MOVEWITHCAMERA
- **`CHAR.MOVEWITHCAMERA`** - args: int, handle, float, float, float — Alias of PLAYER.MOVEWITHCAMERA
- **`CHAR.NAVTO`** - args: int, float, float, float — Alias of PLAYER.NAVTO
- **`CHAR.NAVTO`** - args: int, float, float, float, float — Alias of PLAYER.NAVTO (5-arg)
- **`CHAR.NAVTO`** - args: int, float, float, float, float, float — Alias of PLAYER.NAVTO (6-arg)
- **`CHAR.NAVUPDATE`** - args: int — Alias of PLAYER.NAVUPDATE
- **`CHAR.SETPADDING`** - args: int, float — Alias of PLAYER.SETPADDING (KCC skin padding)
- **`CHAR.SETSLOPE`** - args: int, float — Alias of PLAYER.SETSLOPELIMIT
- **`CHAR.SETSTEP`** - args: int, float — Alias of PLAYER.SETSTEPOFFSET / stair step-up height
- **`CHAR.STICK`** - args: int, float — Alias of PLAYER.SETSTICKFLOOR — glue to floor within max step down (world units)
- **`CHAR.UPDATE`** - args: float — Update kinematic character solver with delta time (legacy; prefer CHARACTERREF.UPDATE / UPDATEPHYSICS)

### CHARACTER

- **`CHARACTER.CREATE`** - args: any -> returns handle — Entity-bound capsule with default dimensions (0.4, 1.0). Linux/Windows fullruntime with CGO+Jolt.
- **`CHARACTER.CREATE`** - args: any, float, float -> returns handle — Entity-bound capsule: (visualEntity, radius, height).
- **`CHARACTER.CREATE`** - args: float, float, float, float, float -> returns handle — World-positioned capsule: (radius, height, x, y, z). Windows/Linux fullruntime with CGO+Jolt: Jolt CharacterVirtual.
- **`CHARACTER.MAKE`** - args: any -> returns handle — DEPRECATED alias of CHARACTER.CREATE. Use CHARACTER.CREATE. Entity-bound capsule with default dimensions (0.4, 1.0). Linux/Windows fullruntime with CGO+Jolt.
- **`CHARACTER.MAKE`** - args: any, float, float -> returns handle — DEPRECATED alias of CHARACTER.CREATE. Use CHARACTER.CREATE. Entity-bound capsule: (visualEntity, radius, height). Windows/Linux fullruntime with CGO+Jolt: Jolt CharacterVirtual.
- **`CHARACTER.MAKE`** - args: float, float, float, float, float -> returns handle — DEPRECATED alias of CHARACTER.CREATE. Use CHARACTER.CREATE. World-positioned capsule: (radius, height, x, y, z). Windows/Linux fullruntime with CGO+Jolt: Jolt CharacterVirtual.
- **`CHARACTERREF.SETFRICTION`** - args: handle, float — Kinematic sliding friction (0..1).
- **`CHARACTERREF.SETGRAVITY`** - args: handle, float — Gravity multiplier for this character (default 1.0).
- **`CHARACTERREF.SETPOS`** - args: handle, float, float, float — Teleports the character capsule to a new world position.
- **`CHARACTERREF.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of CHARACTERREF.SETPOS. Use CHARACTERREF.SETPOS.
- **`CHARACTERREF.SETSETTING`** - args: handle, string, float — Generic float setting for character physics (Gravity, Friction, StepHeight, MaxSlope, SnapDist).

### CHARACTERREF

- **`CHARACTERREF.ADDVELOCITY`** - args: handle, float, float, float — Accumulates world-space velocity (m/s).
- **`CHARACTERREF.GETCEILING`** - args: handle -> returns bool — True if ceiling/head contact detected on last move
- **`CHARACTERREF.GETGROUNDVELOCITY`** - args: handle -> returns handle — Returns [vx,vy,vz] array for ground/platform velocity (same as CHARCONTROLLER.GETGROUNDVELOCITY)
- **`CHARACTERREF.GETISSLIDING`** - args: handle -> returns bool — True on steep-ground sliding (Jolt ground state)
- **`CHARACTERREF.GETPOSITION`** - args: handle -> returns array — Returns [x, y, z] position of the character capsule.
- **`CHARACTERREF.GETROT`** - args: handle -> returns array — Approximate [pitch,yaw,roll] (radians) from linear velocity direction (Jolt CharacterVirtual has no orientation API).
- **`CHARACTERREF.GETSLOPEANGLE`** - args: handle -> returns float — Current surface normal slope angle in degrees.
- **`CHARACTERREF.GETSPEED`** - args: handle -> returns float — Current scalar speed (m/s).
- **`CHARACTERREF.ISGROUNDED`** - args: handle -> returns bool — True if supported by floor geometry.
- **`CHARACTERREF.JUMP`** - args: handle, float — Applies an upward vertical impulse.
- **`CHARACTERREF.MOVE`** - args: handle, float, float, float — Translate character position by (dx, dy, dz); bypasses velocity integration.
- **`CHARACTERREF.ONSLOPE`** - args: handle -> returns bool — True if standing on a floor steeper than epsilon.
- **`CHARACTERREF.ONWALL`** - args: handle -> returns bool — True if colliding with vertical or steep geometry.
- **`CHARACTERREF.SETAIRCONTROL`** - args: handle, float — Scales horizontal move input while airborne
- **`CHARACTERREF.SETGROUNDCONTROL`** - args: handle, float — Scales horizontal move input while on ground
- **`CHARACTERREF.SETJUMPBUFFER`** - args: handle, float — Jump buffer window (seconds) for CHARACTERREF.JUMP while airborne
- **`CHARACTERREF.SETMAXSLOPE`** - args: handle, float — Maximum walkable slope angle in degrees.
- **`CHARACTERREF.SETSNAPDISTANCE`** - args: handle, float — Stick-to-ground snap distance (m).
- **`CHARACTERREF.SETSTEPHEIGHT`** - args: handle, float — Maximum height character can step up onto (m).
- **`CHARACTERREF.SETVELOCITY`** - args: handle, float, float, float — Sets the integrated world-space velocity (m/s).
- **`CHARACTERREF.UPDATE`** - args: handle — Updates character simulation for one physics frame. Handles sweeps and snapping.

### CHARCONTROLLER

- **`CHARACTERREF.DRAINCONTACTS`** - args: handle -> returns handle — Drain pending contact events from a KCC; returns array of contact info.
- **`CHARACTERREF.FREE`** - args: handle — Destroy a Kinematic Character Controller.
- **`CHARACTERREF.SETCONTACTLISTENER`** - args: handle, string — Set a callback function name to be called on KCC contact events.
- **`CHARCONTROLLER.CREATE`** - args: float, float, float, float, float -> returns handle
- **`CHARCONTROLLER.FREE`** - args: handle
- **`CHARCONTROLLER.GETGROUNDNORMAL`** - args: handle -> returns handle — Ground contact normal [nx, ny, nz]; stub returns up or zero when airborne
- **`CHARCONTROLLER.GETGROUNDVELOCITY`** - args: handle -> returns handle — Velocity projected to ground plane (Jolt CharacterVirtual.GetGroundVelocity)
- **`CHARCONTROLLER.GETLINEARVEL`** - args: handle -> returns handle — World linear velocity [vx, vy, vz] (Jolt); stub uses internal velocity
- **`CHARCONTROLLER.GETPOS`** - args: handle -> returns handle
- **`CHARCONTROLLER.GROUNDSTATE`** - args: handle -> returns int — Jolt EGroundState: 0 OnGround, 1 OnSteepGround, 2 NotSupported, 3 InAir (stub: 0 or 3)
- **`CHARCONTROLLER.ISGROUNDED`** - args: handle -> returns bool
- **`CHARCONTROLLER.MAKE`** - args: float, float, float, float, float -> returns handle — DEPRECATED alias of CHARCONTROLLER.CREATE. Use CHARCONTROLLER.CREATE.
- **`CHARCONTROLLER.MOVE`** - args: handle, float, float, float
- **`CHARCONTROLLER.SETPOS`** - args: handle, float, float, float
- **`CHARCONTROLLER.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of CHARCONTROLLER.SETPOS. Use CHARCONTROLLER.SETPOS.
- **`CHARCONTROLLER.TELEPORT`** - args: handle, float, float, float — Snap capsule to (x,y,z) and clear linear velocity
- **`CHARCONTROLLER.X`** - args: handle -> returns float
- **`CHARCONTROLLER.Y`** - args: handle -> returns float
- **`CHARCONTROLLER.Z`** - args: handle -> returns float

### CHECK

- **`CHECK.INVIEW`** - args: int -> returns bool — Same frustum test as ENTITY.INFRUSTUM (active CAMERA.BEGIN)

### CHOOSE

- **`CHOOSE`** - args: any, any
- **`CHOOSE`** - args: any, any, any
- **`CHOOSE`** - args: any, any, any, any
- **`CHOOSE`** - args: any, any, any, any, any
- **`CHOOSE`** - args: any, any, any, any, any, any
- **`CHOOSE`** - args: any, any, any, any, any, any, any
- **`CHOOSE`** - args: any, any, any, any, any, any, any, any
- **`CHOOSE`** - args: any, any, any, any, any, any, any, any, any
- **`CHOOSE`** - args: any, any, any, any, any, any, any, any, any, any
- **`CHOOSE`** - args: any, any, any, any, any, any, any, any, any, any, any
- **`CHOOSE`** - args: any, any, any, any, any, any, any, any, any, any, any, any

### CHR

- **`CHR`** - args: int -> returns string

### CHR$

- **`CHR$`** - args: int -> returns string

### CHUNK

- **`CHUNK.COUNT`** - args: handle -> returns int
- **`CHUNK.GENERATE`** - args: handle, int, int
- **`CHUNK.ISLOADED`** - args: handle, int, int -> returns bool
- **`CHUNK.SETRANGE`** - args: handle, float, float

### CIRCLEPOINT

- **`CIRCLEPOINT`** - args: float, float, float, float, float -> returns handle

### CLAMP

- **`CLAMP`** - args: any, any, any

### CLAMPENTITY2D

- **`CLAMPENTITY2D`** - args: handle, float, float, float, float

### CLIENT

- **`CLIENT.CONNECT`** - args: string, int
- **`CLIENT.ONCONNECT`** - args: string
- **`CLIENT.ONMESSAGE`** - args: string
- **`CLIENT.ONSYNC`** - args: string
- **`CLIENT.STOP`** - args: (none)
- **`CLIENT.TICK`** - args: float

### CLIPBOARD

- **`CLIPBOARD.GETIMAGE`** - args: (none) -> returns handle

### CLOSEFILE

- **`CLOSEFILE`** - args: handle

### CLOUD

- **`CLOUD.CREATE`** - args: (none) -> returns handle — Procedural cloud volume handle (CPU-side; coverage via CLOUD.SETCOVERAGE).
- **`CLOUD.DRAW`** - args: handle — Draw cloud volume (stub path).
- **`CLOUD.FREE`** - args: handle — Release cloud handle.
- **`CLOUD.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of CLOUD.CREATE. Use CLOUD.CREATE.
- **`CLOUD.SETCOVERAGE`** - args: handle, float — Set cloud coverage amount (0..1).
- **`CLOUD.UPDATE`** - args: handle, float — Advance cloud state (dt in seconds).

### CLS

- **`CLS`** - args: (none)

### COLLISIONS

- **`COLLISIONS`** - args: int, int, int, int -> returns void — Easy Mode: Set global collision response rule (srcType, dstType, method, response)
- **`COLLISIONS`** - args: int, int, int, int — Easy Mode: Define collision rules between types

### COLOR

- **`COLOR.A`** - args: handle -> returns int
- **`COLOR.B`** - args: handle -> returns int
- **`COLOR.BRIGHTNESS`** - args: handle, float -> returns handle
- **`COLOR.CLAMP`** - args: float, float, float -> returns handle
- **`COLOR.CONTRAST`** - args: handle, float -> returns handle
- **`COLOR.FADE`** - args: handle, float -> returns handle
- **`COLOR.FREE`** - args: handle
- **`COLOR.FROMHSV`** - args: float, float, float -> returns handle
- **`COLOR.G`** - args: handle -> returns int
- **`COLOR.HEX`** - args: string -> returns handle
- **`COLOR.HSV`** - args: float, float -> returns handle — COLOR.HSV(index, total) — evenly spaced hues on the wheel
- **`COLOR.HSV`** - args: float, float, float -> returns handle
- **`COLOR.INVERT`** - args: handle -> returns handle
- **`COLOR.LERP`** - args: handle, handle, float -> returns handle
- **`COLOR.R`** - args: handle -> returns int
- **`COLOR.RGB`** - args: int, int, int -> returns handle
- **`COLOR.RGBA`** - args: int, int, int, int -> returns handle
- **`COLOR.TOHEX`** - args: handle -> returns string
- **`COLOR.TOHSV`** - args: handle -> returns handle
- **`COLOR.TOHSVX`** - args: handle -> returns float
- **`COLOR.TOHSVY`** - args: handle -> returns float
- **`COLOR.TOHSVZ`** - args: handle -> returns float

### COLORPRINT

- **`COLORPRINT`** - args: int, int, int, string — Print colored text to console

### COMMAND

- **`COMMAND`** - args: (none)
- **`COMMAND`** - args: int

### COMMAND$

- **`COMMAND$`** - args: (none)
- **`COMMAND$`** - args: int

### COMPUTESHADER

- **`COMPUTESHADER.BUFFERFREE`** - args: handle
- **`COMPUTESHADER.BUFFERMAKE`** - args: int -> returns handle
- **`COMPUTESHADER.DISPATCH`** - args: handle, int, int, int
- **`COMPUTESHADER.FREE`** - args: handle
- **`COMPUTESHADER.LOAD`** - args: string -> returns handle
- **`COMPUTESHADER.SETBUFFER`** - args: handle, int, handle
- **`COMPUTESHADER.SETFLOAT`** - args: handle, string, float
- **`COMPUTESHADER.SETINT`** - args: handle, string, int

### CONNECT

- **`CONNECT`** - args: string, int -> returns handle — Easy Mode: NET.CONNECT(host, port)

### CONTAINS

- **`CONTAINS`** - args: string, string

### CONTROLLER

- **`CONTROLLER.CREATE`** - args: float, float, float, float, float -> returns handle
- **`CONTROLLER.FREE`** - args: handle
- **`CONTROLLER.GROUNDED`** - args: handle -> returns bool
- **`CONTROLLER.JUMP`** - args: handle, float
- **`CONTROLLER.MAKE`** - args: float, float, float, float, float -> returns handle — DEPRECATED alias of CONTROLLER.CREATE. Use CONTROLLER.CREATE.
- **`CONTROLLER.MOVE`** - args: handle, float, float, float

### COPYFILE

- **`COPYFILE`** - args: string, string
- **`COPYFILE`** - args: string, string -> returns bool

### COS

- **`COS`** - args: any

### COSD

- **`COSD`** - args: any

### COUNT

- **`COUNT`** - args: string, string -> returns int

### COUNT$

- **`COUNT$`** - args: string, string -> returns int

### COUNTCOLLISIONS

- **`COUNTCOLLISIONS`** - args: handle -> returns int — Easy Mode: Get number of active collisions for entity

### CREATEBODY

- **`CREATEBODY`** - args: int, int -> returns handle — Easy Mode: PHYSICS3D.CREATEBODY(type, shape)

### CREATEBODY2D

- **`CREATEBODY2D`** - args: int, int -> returns handle — Easy Mode: PHYSICS2D.CREATEBODY(type, shape)

### CREATECAMERA

- **`CREATECAMERA`** - args: (none) -> returns handle — Blitz-style: forwards to CAMERA.CREATE (3D heap camera)

### CREATECAMERA2D

- **`CREATECAMERA2D`** - args: (none) -> returns handle — Easy Mode: forwards to CAMERA2D.CREATE (same as deprecated CAMERA2D.MAKE)

### CREATECUBE

- **`CREATECUBE`** - args: (none) -> returns int — Easy Mode: ENTITY.CREATECUBE(1, 1, 1)

### CREATEEMITTER

- **`CREATEEMITTER`** - args: (none) -> returns handle — Easy Mode: Create a 3D particle emitter

### CREATELIGHT

- **`CREATELIGHT`** - args: (none) -> returns handle — Blitz-style: forwards to LIGHT.CREATE (same as deprecated LIGHT.MAKE)

### CUBE

- **`CUBE`** - args: (none) -> returns handle — Blitz-style static box entity (1×1×1); use CUBE(w,h,d) for dimensions — ENTITYREF handle
- **`CUBE`** - args: float, float, float -> returns handle — Blitz-style static box entity — ENTITYREF handle

### CURSOR

- **`CURSOR.DISABLE`** - args: (none)
- **`CURSOR.ENABLE`** - args: (none)
- **`CURSOR.HIDE`** - args: (none)
- **`CURSOR.ISHIDDEN`** - args: (none)
- **`CURSOR.ISONSCREEN`** - args: (none)
- **`CURSOR.SET`** - args: int
- **`CURSOR.SHOW`** - args: (none)

### CURVE

- **`CURVE`** - args: float, float, float -> returns float — Easy Mode: Blitz-style smooth follower (value, target, divisor)

### CURVEANGLE

- **`CURVEANGLE`** - args: float, float, float -> returns float — Like CURVEVALUE for degrees (360 wrap)

### CURVEVALUE

- **`CURVEVALUE`** - args: float, float, float -> returns float — DBPro-style: move current toward target by (target-current)/speed per call

### CVDOUBLE

- **`CVDOUBLE`** - args: string

### CVFLOAT

- **`CVFLOAT`** - args: string

### CVINT

- **`CVINT`** - args: string

### CVLONG

- **`CVLONG`** - args: string

### CVSHORT

- **`CVSHORT`** - args: string

### CameraFOV

- **`CameraFOV`** - args: handle, float

### CameraLookAt

- **`CameraLookAt`** - args: handle, float, float, float

### CameraShake

- **`CameraShake`** - args: handle, float, float

### CameraSmoothFollow

- **`CameraSmoothFollow`** - args: handle, int, float

### CollisionForce

- **`CollisionForce`** - args: (none) -> returns float — Penetration-depth proxy for impact strength (not true Jolt impulse on this path)

### CollisionNX

- **`CollisionNX`** - args: (none) -> returns float — World normal X from last successful EntityCollided query this frame

### CollisionNY

- **`CollisionNY`** - args: (none) -> returns float — World normal Y from last successful EntityCollided query

### CollisionNZ

- **`CollisionNZ`** - args: (none) -> returns float — World normal Z from last successful EntityCollided query

### CollisionPX

- **`CollisionPX`** - args: (none) -> returns float — Contact point X (shape query) after last EntityCollided

### CollisionPY

- **`CollisionPY`** - args: (none) -> returns float — Contact point Y after last EntityCollided

### CollisionPZ

- **`CollisionPZ`** - args: (none) -> returns float — Contact point Z after last EntityCollided

### CollisionY

- **`CollisionY`** - args: (none) -> returns float — Alias for CollisionPY (contact Y)

### CountCollisions

- **`CountCollisions`** - args: int -> returns int — Count Jolt contact pairs involving entity# this frame (distinct from COUNTCOLLISIONS legacy hits)

### CreateBrush

- **`CreateBrush`** - args: float, float, float -> returns handle

### CreateCube

- **`CreateCube`** - args: (none) -> returns int — Default 1x1x1 axis-aligned box; returns entity#
- **`CreateCube`** - args: int -> returns int — 1x1x1 box parented to entity# (parent entity id)
- **`CreateCube`** - args: float, float, float -> returns int — Box with width, height, depth (no parent)
- **`CreateCube`** - args: int, float, float, float -> returns int — Box (w,h,d) parented to entity#

### CreateLight

- **`CreateLight`** - args: any, any -> returns handle — Blitz-style: type 1=directional, 2=point, 3=spot; optional parent entity# stored for future attachment

### CreatePivot

- **`CreatePivot`** - args: (none) -> returns int — Create empty transform node (entity#) for parenting; invisible, no mesh

### CreatePointLight

- **`CreatePointLight`** - args: int, float, float, float -> returns handle

### CreateSurface

- **`CreateSurface`** - args: int -> returns handle

### CreateVehicle

- **`CreateVehicle`** - args: int

### DATA

- **`DATA.COMPRESS`** - args: string -> returns string
- **`DATA.COMPUTECRC32`** - args: string -> returns int
- **`DATA.COMPUTEMD5`** - args: string -> returns string
- **`DATA.COMPUTESHA1`** - args: string -> returns string
- **`DATA.CRC32`** - args: string -> returns int
- **`DATA.DECODEBASE64`** - args: string -> returns string
- **`DATA.DECOMPRESS`** - args: string -> returns string
- **`DATA.ENCODEBASE64`** - args: string -> returns string
- **`DATA.MD5`** - args: string -> returns string
- **`DATA.SHA1`** - args: string -> returns string

### DATE

- **`DATE`** - args: (none)
- **`DATE`** - args: (none) -> returns string

### DATE$

- **`DATE$`** - args: (none)
- **`DATE$`** - args: (none) -> returns string

### DATETIME

- **`DATETIME`** - args: (none)
- **`DATETIME`** - args: (none) -> returns string

### DATETIME$

- **`DATETIME$`** - args: (none)
- **`DATETIME$`** - args: (none) -> returns string

### DAY

- **`DAY`** - args: (none)
- **`DAY`** - args: (none) -> returns int

### DEBUG

- **`CONSOLE.LOG`** - args: string — Add a message to the scrolling on-screen debug console.
- **`CONSOLE.LOG`** - args: string, handle — Add a colored message to the scrolling on-screen debug console.
- **`DEBUG.ASSERT`** - args: any, string
- **`DEBUG.BREAKPOINT`** - args: (none)
- **`DEBUG.DRAWBODY`** - args: handle — Renders body collision shape.
- **`DEBUG.DRAWBOX`** - args: float, float, float, float, float, float, int, int, int
- **`DEBUG.DRAWCHARACTER`** - args: handle — Renders capsule wireframe and ground probes for character Controller.
- **`DEBUG.DRAWLINE`** - args: float, float, float, float, float, float, int, int, int
- **`DEBUG.DRAWPHYSICS`** - args: bool — Toggle collision wireframe visualization.
- **`DEBUG.DUMPHEAP`** - args: (none) — Professional: Scan all active handles and print to diagnostics.
- **`DEBUG.GCSTATS`** - args: (none)
- **`DEBUG.HEAPSTATS`** - args: (none)
- **`DEBUG.INSPECT`** - args: int — Display live transform/state info for an entity.
- **`DEBUG.LISTCOMMANDS`** - args: (none) — Professional: List all registered built-in commands.
- **`DEBUG.LOG`** - args: string
- **`DEBUG.LOGFILE`** - args: string, string
- **`DEBUG.PRINT`** - args: any
- **`DEBUG.PRINT`** - args: string
- **`DEBUG.PRINT`** - args: string, any
- **`DEBUG.PRINT`** - args: string, any, any
- **`DEBUG.PRINT`** - args: string, any, any, any
- **`DEBUG.PRINT`** - args: string, any, any, any, any
- **`DEBUG.PRINT`** - args: string, any, any, any, any, any
- **`DEBUG.PRINT`** - args: string, any, any, any, any, any, any
- **`DEBUG.PRINT`** - args: string, any, any, any, any, any, any, any
- **`DEBUG.PRINT`** - args: string, any, any, any, any, any, any, any, any
- **`DEBUG.PRINT`** - args: string, any, any, any, any, any, any, any, any, any
- **`DEBUG.PRINT`** - args: string, any, any, any, any, any, any, any, any, any, any
- **`DEBUG.PRINTL`** - args: string, any
- **`DEBUG.PROFILEEND`** - args: string
- **`DEBUG.PROFILEREPORT`** - args: (none)
- **`DEBUG.PROFILESTART`** - args: string
- **`DEBUG.SHOWFPSGRAPH`** - args: bool — Show or hide the real-time FPS graph overlay.
- **`DEBUG.STACKTRACE`** - args: (none)
- **`DEBUG.WATCH`** - args: string, any
- **`DEBUG.WATCHCLEAR`** - args: (none)
- **`SYSTEM.MONITOR`** - args: (none) — Toggle the system performance monitor (FPS, RAM).
- **`SYSTEM.MONITOR`** - args: bool — Toggle the system performance monitor (FPS, RAM).

### DECAL

- **`DECAL.CREATE`** - args: handle -> returns handle
- **`DECAL.DRAW`** - args: handle
- **`DECAL.FREE`** - args: handle
- **`DECAL.GETPOS`** - args: handle -> returns array — Returns [x, y, z] position of decal
- **`DECAL.GETROT`** - args: handle -> returns array — Returns [p, y, r] rotation of decal
- **`DECAL.GETSIZE`** - args: handle -> returns array — Returns [w, h] size of decal
- **`DECAL.MAKE`** - args: handle -> returns handle — DEPRECATED alias of DECAL.CREATE. Use DECAL.CREATE.
- **`DECAL.SETLIFETIME`** - args: handle, float
- **`DECAL.SETPOS`** - args: handle, float, float, float
- **`DECAL.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of DECAL.SETPOS.
- **`DECAL.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of DECAL.SETPOS. Use DECAL.SETPOS.
- **`DECAL.SETSIZE`** - args: handle, float, float

### DEG2RAD

- **`DEG2RAD`** - args: any

### DEGPERSEC

- **`DEGPERSEC`** - args: any, any

### DELAY

- **`DELAY`** - args: int -> returns void — Easy Mode: Blocking wait (ms)

### DELETEDIR

- **`DELETEDIR`** - args: string
- **`DELETEDIR`** - args: string -> returns bool

### DELETEFILE

- **`DELETEFILE`** - args: string
- **`DELETEFILE`** - args: string -> returns bool

### DIREXISTS

- **`DIREXISTS`** - args: string
- **`DIREXISTS`** - args: string -> returns bool

### DIST2D

- **`DIST2D`** - args: float, float, float, float -> returns float — 2D Euclidean distance; alias of DISTANCE2D under MATH

### DIST3D

- **`DIST3D`** - args: float, float, float, float, float, float -> returns float — Easy Mode: Distance between two points in 3D space

### DISTSQ2D

- **`DISTSQ2D`** - args: float, float, float, float -> returns float — Squared 2D distance

### DRAW

- **`DRAW.ARC`** - args: float, float, float, float, float, float, int, int, int, int
- **`DRAW.CIRCLE`** - args: int, int, float, int, int, int, int
- **`DRAW.CIRCLELINES`** - args: int, int, float, int, int, int, int
- **`DRAW.DOT`** - args: float, float, float, int, int, int, int
- **`DRAW.ELLIPSE`** - args: int, int, float, float, int, int, int, int
- **`DRAW.ELLIPSELINES`** - args: int, int, float, float, int, int, int, int
- **`DRAW.GETPIXELCOLOR`** - args: int, int -> returns array
- **`DRAW.GRID`** - args: int, float
- **`DRAW.GRID2D`** - args: int, int, int, int, int
- **`DRAW.LINE`** - args: int, int, int, int, int, int, int, int
- **`DRAW.LINE3D`** - args: float, float, float, float, float, float, int, int, int, int
- **`DRAW.LINEBEZIER`** - args: float, float, float, float, float, int, int, int, int
- **`DRAW.LINEEX`** - args: float, float, float, float, float, int, int, int, int
- **`DRAW.OVAL`** - args: int, int, float, float, int, int, int, int
- **`DRAW.OVALLINES`** - args: int, int, float, float, int, int, int, int
- **`DRAW.PIXEL`** - args: int, int, int, int, int, int
- **`DRAW.PLOT`** - args: int, int, int, int, int, int
- **`DRAW.POLY`** - args: float, float, int, float, float, int, int, int, int
- **`DRAW.POLYLINES`** - args: float, float, int, float, float, float, int, int, int, int
- **`DRAW.RECTANGLE`** - args: int, int, int, int, int, int, int, int
- **`DRAW.RECTANGLE_ROUNDED`** - args: int, int, int, int, int, int, int, int, int
- **`DRAW.RING`** - args: float, float, float, float, float, float, int, int, int, int, int
- **`DRAW.RINGLINES`** - args: float, float, float, float, float, float, int, int, int, int, int
- **`DRAW.TEXT`** - args: string, int, int, int, int, int, int, int
- **`DRAW.TEXTEX`** - args: handle, string, float, float, float, float, int, int, int, int
- **`DRAW.TEXTFONT`** - args: handle, string, float, float, float, float, int, int, int, int — Same handler as DRAW.TEXTEX — DrawTextEx with a loaded FONT.* handle
- **`DRAW.TEXTURE`** - args: handle, int, int, int, int, int, int
- **`DRAW.TEXTURENPATCH`** - args: handle, int, int, int, int, int, int, int, int, int, int, int, int
- **`DRAW.TRIANGLE`** - args: float, float, float, float, float, float, int, int, int, int
- **`DRAW.TRIANGLELINES`** - args: float, float, float, float, float, float, int, int, int, int

### DRAW3D

- **`DRAW3D.BBOX`** - args: float, float, float, float, float, float, int, int, int, int
- **`DRAW3D.BILLBOARD`** - args: handle, float, float, float, float, int, int, int, int
- **`DRAW3D.BILLBOARDREC`** - args: handle, float, float, float, float, float, float, float, float, float, int, int, int, int
- **`DRAW3D.CAPSULE`** - args: float, float, float, float, float, float, float, int, int, int, int, int, int
- **`DRAW3D.CAPSULEWIRES`** - args: float, float, float, float, float, float, float, int, int, int, int, int, int
- **`DRAW3D.CUBE`** - args: float, float, float, float, float, float, int, int, int, int
- **`DRAW3D.CUBEWIRES`** - args: float, float, float, float, float, float, int, int, int, int
- **`DRAW3D.CYLINDER`** - args: float, float, float, float, float, float, int, int, int, int, int
- **`DRAW3D.CYLINDERWIRES`** - args: float, float, float, float, float, float, int, int, int, int, int
- **`DRAW3D.GRID`** - args: int, float
- **`DRAW3D.GRID`** - args: int, float, float — XZ grid with optional Y offset (avoids z-fight with floor at Y=0)
- **`DRAW3D.LINE`** - args: float, float, float, float, float, float, int, int, int, int
- **`DRAW3D.PLANE`** - args: float, float, float, float, float, int, int, int, int
- **`DRAW3D.POINT`** - args: float, float, float, int, int, int, int
- **`DRAW3D.RAY`** - args: handle, int, int, int, int
- **`DRAW3D.SPHERE`** - args: float, float, float, float, int, int, int, int
- **`DRAW3D.SPHEREWIRES`** - args: float, float, float, float, int, int, int, int, int, int

### DRAWCUBE

- **`DRAWCUBE`** - args: (none) — Immediate-mode 3D box wrapper; use .Pos/.Size/.Color/.Draw (see DRAW_WRAPPERS.md)
- **`DRAWCUBE`** - args: float, float, float — DRAWCUBE(w,h,d) initial size

### DRAWEMITTER

- **`DRAWEMITTER`** - args: handle -> returns void — Easy Mode: Render particles from an emitter

### DRAWPOLY2

- **`DRAWPOLY2`** - args: int -> returns handle

### DRAWPRIM2D

- **`DRAWPRIM2D.DRAW`** - args: handle

### DRAWPRIM3D

- **`DRAWPRIM3D.DRAW`** - args: handle

### DRAWRING2

- **`DRAWRING2`** - args: (none) -> returns handle

### DRAWSPHERE

- **`DRAWSPHERE`** - args: float — DRAWCUBE-style sphere; radius

### DRAWTEXPRO

- **`DRAWTEXPRO`** - args: handle -> returns handle

### DRAWTEXREC

- **`DRAWTEXREC`** - args: handle -> returns handle

### DUMP

- **`DUMP`** - args: any

### DrawEntities

- **`DrawEntities`** - args: (none) — Alias for ENTITY.DRAWALL: draw all entities in the scene graph (no arguments)

### DrawEntity

- **`DrawEntity`** - args: int — Draw one entity (same as ENTITY.DRAW)

### E

- **`E`** - args: (none)

### EFFECT

- **`EFFECT.BLOOM`** - args: bool
- **`EFFECT.BLOOM`** - args: bool, float
- **`EFFECT.BLOOM`** - args: bool, float, float
- **`EFFECT.CHROMATICABERRATION`** - args: bool
- **`EFFECT.CHROMATICABERRATION`** - args: bool, float
- **`EFFECT.DEPTHOFFIELD`** - args: bool
- **`EFFECT.DEPTHOFFIELD`** - args: bool, float
- **`EFFECT.DEPTHOFFIELD`** - args: bool, float, float
- **`EFFECT.GRAIN`** - args: bool
- **`EFFECT.GRAIN`** - args: bool, float
- **`EFFECT.MOTIONBLUR`** - args: bool
- **`EFFECT.MOTIONBLUR`** - args: bool, float
- **`EFFECT.SHARPEN`** - args: bool
- **`EFFECT.SHARPEN`** - args: bool, float
- **`EFFECT.SSAO`** - args: bool
- **`EFFECT.SSAO`** - args: bool, float
- **`EFFECT.SSAO`** - args: bool, float, float
- **`EFFECT.SSR`** - args: bool
- **`EFFECT.SSR`** - args: bool, float
- **`EFFECT.SSR`** - args: bool, float, float
- **`EFFECT.TONEMAPPING`** - args: string
- **`EFFECT.VIGNETTE`** - args: bool
- **`EFFECT.VIGNETTE`** - args: bool, float

### EMITPARTICLE

- **`EMITPARTICLE`** - args: handle, int -> returns void — Easy Mode: Burst particles from an emitter
- **`EMITPARTICLE`** - args: handle, int

### EMITTERALIVE

- **`EMITTERALIVE`** - args: handle -> returns int — Easy Mode: Check if emitter is playing or has active particles

### EMITTERCOUNT

- **`EMITTERCOUNT`** - args: handle -> returns int — Easy Mode: Get number of active particles in emitter

### EMITTERPOS

- **`EMITTERPOS`** - args: handle, float, float, float -> returns void — Easy Mode: Reposition an emitter
- **`EMITTERPOS`** - args: handle, float, float, float

### ENDSWITH

- **`ENDSWITH`** - args: string, string -> returns bool

### ENEMY

- **`ENEMY.FOLLOWPATH`** - args: int, handle, float — Moves an entity along a PATH handle toward waypoints at speed (world units/sec)

### ENET

- **`ENET.CREATEHOST`** - args: string, int, int, int, int
- **`ENET.DEINITIALIZE`** - args: (none)
- **`ENET.HOSTBROADCAST`** - args: handle, int, int, handle
- **`ENET.HOSTSERVICE`** - args: handle, int
- **`ENET.INITIALIZE`** - args: (none)
- **`ENET.MAKEHOST`** - args: string, int, int, int, int — DEPRECATED alias of ENET.CREATEHOST. Use ENET.CREATEHOST.
- **`ENET.PEERPING`** - args: handle
- **`ENET.PEERSEND`** - args: handle, int, handle

### ENT

- **`ENT.DAMAGE`** - args: int, float — Reduce entity HP by amount#; triggers damage effects/logic
- **`ENT.DIST`** - args: int, int -> returns float — Alias of ENTITY.DIST
- **`ENT.FADE`** - args: int, float, float — Fade to target alpha over duration — convenience over ENTITY.FADE
- **`ENT.GETNEAREST`** - args: int, float, string -> returns handle — Alias of ENT.GET_NEAREST / PLAYER.GETNEARBY
- **`ENT.GET_NEAREST`** - args: int, float, string -> returns handle — Alias of PLAYER.GETNEARBY — entities with matching tag within radius (float array of ids)
- **`ENT.ONDEATH`** - args: int, int — Death-drop prefab with 100% chance — alias of ENTITY.ONDEATHDROP(entity, prefab, 100)
- **`ENT.ONDEATH`** - args: int, string — Prefab by ENTITY.SETNAME / registry name (same as int overload)
- **`ENT.SETHP`** - args: int, float, float — Alias of ENT.SET_HP / ENTITY.SETHEALTH
- **`ENT.SETTEAM`** - args: int, int — Alias of ENT.SET_TEAM
- **`ENT.SET_TEAM`** - args: int, int — Stores team id on entity (gameplay / friendly-fire bookkeeping)
- **`ENT.SHOOT`** - args: int, int, float -> returns int — Spawn ENTITY.COPY of prefab at shooter forward; sets host velocity (scripted projectile)
- **`ENT.SHOOT`** - args: int, string, float -> returns int — Prefab by registered name string
- **`ENT.TWEEN`** - args: int, float, float, float, float — Smooth move to world (x,y,z) over duration — alias of ENTITY.ANIMATETOWARD
- **`ENT.WOBBLE`** - args: int, float, float — Alias of ENTITY.ADDWOBBLE — bob amplitude and speed

### ENTHIT

- **`ENTHIT`** - args: handle, int -> returns handle — Shorthand: ENTITYCOLLIDED(ent, type)

### ENTITY

- **`ENTITY.ADDFORCE`** - args: int, float, float, float
- **`ENTITY.ADDPHYSICS`** - args: int, string, string — One-line Jolt body: motion (static/dynamic), shape (box/capsule/sphere)
- **`ENTITY.ADDPHYSICS`** - args: int, string, string, float
- **`ENTITY.ADDTRIANGLE`** - args: handle, int, int, int
- **`ENTITY.ADDVERTEX`** - args: handle, float, float, float -> returns int
- **`ENTITY.ALIGNTOVECTOR`** - args: int, float, float, float, int
- **`ENTITY.ALPHA`** - args: int, float — Easy Mode: Set entity transparency (0.0 to 1.0)
- **`ENTITY.ANIMATE`** - args: int, any, any
- **`ENTITY.ANIMATETOWARD`** - args: int, float, float, float, float — Linear world lerp over duration (seconds); advanced in ENTITY.UPDATE
- **`ENTITY.ANIMCOUNT`** - args: int -> returns int
- **`ENTITY.ANIMINDEX`** - args: int -> returns int
- **`ENTITY.ANIMLENGTH`** - args: int -> returns float
- **`ENTITY.ANIMNAME`** - args: any, int -> returns string
- **`ENTITY.ANIMNAME$`** - args: any, int -> returns string
- **`ENTITY.ANIMTIME`** - args: int -> returns float
- **`ENTITY.APPLYGRAVITY`** - args: int, float, float
- **`ENTITY.APPLYIMPULSE`** - args: int, float, float, float — Same as ENTITY.ADDFORCE / ApplyEntityImpulse (velocity change; not Jolt BodyInterface impulse until exposed)
- **`ENTITY.APPLYTORQUE`** - args: handle, float, float, float — Spins physics object.
- **`ENTITY.ATTACH`** - args: handle, handle, float, float, float — Welds entities together with offset.
- **`ENTITY.BLEND`** - args: int, int
- **`ENTITY.BOX`** - args: int, float, float, float
- **`ENTITY.CANSEE`** - args: int, int, float, float -> returns bool — Vision cone (degrees) + max distance + unobstructed Jolt ray to target
- **`ENTITY.CHECKCOLLISION`** - args: int, int -> returns bool — True if two entities had a Jolt contact last step (same as EntityCollided)
- **`ENTITY.CHECKRADIUS`** - args: handle, float, string -> returns handle — Check sensor
- **`ENTITY.CLAMPTOTERRAIN`** - args: int, handle — Sets Y from terrain height at entity XZ (offset 0); alias of TERRAIN.SNAPY argument order swap
- **`ENTITY.CLEARPHYSBUFFER`** - args: int — Remove physics matrix buffer binding from entity#
- **`ENTITY.CLEARSCENE`** - args: (none)
- **`ENTITY.COLLIDE`** - args: int, int
- **`ENTITY.COLLIDED`** - args: int -> returns bool
- **`ENTITY.COLLISIONLAYER`** - args: int, int — Reserved 0..31 layer id for future Jolt bitmask filtering (stored on entity)
- **`ENTITY.COLLISIONNX`** - args: int -> returns float
- **`ENTITY.COLLISIONNY`** - args: int -> returns float
- **`ENTITY.COLLISIONNZ`** - args: int -> returns float
- **`ENTITY.COLLISIONOTHER`** - args: int -> returns int
- **`ENTITY.COLLISIONX`** - args: int -> returns float
- **`ENTITY.COLLISIONY`** - args: int -> returns float
- **`ENTITY.COLLISIONZ`** - args: int -> returns float
- **`ENTITY.COLOR`** - args: int, handle
- **`ENTITY.COLOR`** - args: int, int, int, int
- **`ENTITY.COLOR`** - args: int, int, int, int, int
- **`ENTITY.COLORPULSE`** - args: handle, handle, handle, float — Pulses color.
- **`ENTITY.COPY`** - args: int -> returns int
- **`ENTITY.COUNTCHILDREN`** - args: int -> returns int
- **`ENTITY.CREATE`** - args: (none) -> returns int
- **`ENTITY.CREATEBOX`** - args: float -> returns int — Uniform cube: size# used for width, height, and depth (alias ENTITY.CREATECUBE)
- **`ENTITY.CREATEBOX`** - args: float, float, float -> returns int
- **`ENTITY.CREATECUBE`** - args: float, float, float -> returns int
- **`ENTITY.CREATECYLINDER`** - args: float, float, int -> returns int
- **`ENTITY.CREATEENTITY`** - args: (none) -> returns int
- **`ENTITY.CREATEMESH`** - args: any -> returns int — Procedural mesh: optional parentEntity#; use AddVertex/UpdateMesh
- **`ENTITY.CREATEPLANE`** - args: float -> returns int
- **`ENTITY.CREATESPHERE`** - args: float -> returns int — Radius only — default 16 segments
- **`ENTITY.CREATESPHERE`** - args: float, int -> returns int
- **`ENTITY.CREATESPRITE`** - args: string -> returns int
- **`ENTITY.CREATESPRITE`** - args: string, int -> returns int
- **`ENTITY.CREATESPRITE`** - args: handle, float, float -> returns int — Billboard from TEXTURE handle (atlas / TEXTURE.LOADANIM)
- **`ENTITY.CREATESPRITE`** - args: handle, float, float, int -> returns int
- **`ENTITY.CREATESURFACE`** - args: int -> returns handle
- **`ENTITY.CROSSFADE`** - args: int, any, float
- **`ENTITY.CURRENTANIM`** - args: any -> returns string
- **`ENTITY.CURRENTANIM$`** - args: any -> returns string
- **`ENTITY.DELTAX`** - args: int, int -> returns float
- **`ENTITY.DELTAY`** - args: int, int -> returns float
- **`ENTITY.DELTAZ`** - args: int, int -> returns float
- **`ENTITY.DIST`** - args: int, int -> returns float — 3D distance between two entities (alias of ENTITY.DISTANCE semantics)
- **`ENTITY.DISTANCE`** - args: int, int -> returns float
- **`ENTITY.DISTANCETO`** - args: handle, handle -> returns float — Returns distance.
- **`ENTITY.DRAW`** - args: int
- **`ENTITY.DRAWALL`** - args: (none)
- **`ENTITY.EMITPARTICLES`** - args: handle, handle — Attaches particles to entity.
- **`ENTITY.ENTITIESINBOX`** - args: float, float, float, float, float, float
- **`ENTITY.ENTITIESINGROUP`** - args: any
- **`ENTITY.ENTITIESINRADIUS`** - args: float, float, float, float
- **`ENTITY.ENTITYPITCH`** - args: int, any -> returns float
- **`ENTITY.ENTITYROLL`** - args: int, any -> returns float
- **`ENTITY.ENTITYX`** - args: int, any -> returns float
- **`ENTITY.ENTITYY`** - args: int, any -> returns float
- **`ENTITY.ENTITYYAW`** - args: int, any -> returns float
- **`ENTITY.ENTITYZ`** - args: int, any -> returns float
- **`ENTITY.EXPLODE`** - args: handle, int — Instantly explodes object.
- **`ENTITY.EXTRACTANIMSEQ`** - args: int, any, any
- **`ENTITY.FADE`** - args: handle, float, float, float — Lerps alpha.
- **`ENTITY.FIND`** - args: any -> returns int
- **`ENTITY.FINDBONE`** - args: int, any -> returns int
- **`ENTITY.FINDBYPROPERTY`** - args: string, string -> returns handle
- **`ENTITY.FINDCHILD`** - args: int, string -> returns int
- **`ENTITY.FLEE`** - args: handle, handle, float, float — Runs away.
- **`ENTITY.FLOOR`** - args: int -> returns float
- **`ENTITY.FREE`** - args: int
- **`ENTITY.FREEENTITIES`** - args: handle
- **`ENTITY.FX`** - args: int, int
- **`ENTITY.GETALPHA`** - args: int -> returns float
- **`ENTITY.GETBONEPOS`** - args: int, string -> returns handle
- **`ENTITY.GETBONEROT`** - args: int, string -> returns handle
- **`ENTITY.GETBOUNDS`** - args: int -> returns handle
- **`ENTITY.GETBUOYANCY`** - args: int -> returns float — Alias of PHYSICS.GETBUOYANCY
- **`ENTITY.GETCHILD`** - args: int, int -> returns int
- **`ENTITY.GETCLOSESTWITHTAG`** - args: int, float, string -> returns int — Nearest entity within radius matching name/tag glob (same rules as PLAYER.GETNEARBY)
- **`ENTITY.GETCOLOR`** - args: int -> returns array
- **`ENTITY.GETDISTANCE`** - args: int, int -> returns float
- **`ENTITY.GETGROUNDNORMAL`** - args: int -> returns handle — World ground normal under entity (CharacterVirtual if PLAYER.CREATE; else short downward Jolt ray)
- **`ENTITY.GETMETADATA`** - args: int, string -> returns string
- **`ENTITY.GETOVERLAPCOUNT`** - args: int, string -> returns int — Counts tagged entities whose pivot lies in zone entity world AABB (sphere prefilter)
- **`ENTITY.GETPOS`** - args: int -> returns handle
- **`ENTITY.GETPOSITION`** - args: int -> returns handle
- **`ENTITY.GETROT`** - args: int -> returns handle
- **`ENTITY.GETSCALE`** - args: int -> returns handle
- **`ENTITY.GETSTATE`** - args: handle -> returns int — Returns string AI state.
- **`ENTITY.GETXZ`** - args: int -> returns handle
- **`ENTITY.GHOSTMODE`** - args: handle, float — Disables collisions temporarily.
- **`ENTITY.GRAVITY`** - args: int, float
- **`ENTITY.GROUNDED`** - args: int -> returns bool
- **`ENTITY.GROUPADD`** - args: any, int
- **`ENTITY.GROUPCREATE`** - args: any
- **`ENTITY.GROUPREMOVE`** - args: any, int
- **`ENTITY.HASTAG`** - args: int, string -> returns bool — Glob match on Blender tag or entity name only (stricter than ENTITY.ISTYPE)
- **`ENTITY.HIDE`** - args: int
- **`ENTITY.INFRUSTUM`** - args: int -> returns bool — True if entity AABB intersects active CAMERA.BEGIN frustum (same as ENTITY.INVIEW without passing camera)
- **`ENTITY.INFRUSTUM`** - args: handle, handle -> returns int — Boolean bounds.
- **`ENTITY.INSTANCE`** - args: int -> returns int
- **`ENTITY.INSTANCEGRID`** - args: int, int, int, float -> returns int
- **`ENTITY.INVIEW`** - args: int, handle -> returns bool
- **`ENTITY.ISPLAYING`** - args: int -> returns bool
- **`ENTITY.ISSUBMERGED`** - args: int -> returns float — Fraction 0..1 of entity vertical extent below water surface (any overlapping WATER volume)
- **`ENTITY.ISTYPE`** - args: int, string -> returns bool
- **`ENTITY.JUMP`** - args: int, float
- **`ENTITY.LINEOFSIGHT`** - args: int, int -> returns bool — Unobstructed Jolt ray from observer eye to target (no FOV); sensors still occlude until filtered
- **`ENTITY.LINKPHYSBUFFER`** - args: int, int — Bind entity# to Jolt shared matrix slot index (use BODY3D.BUFFERINDEX on the body)
- **`ENTITY.LOAD`** - args: any -> returns int — Alias of ENTITY.LOADMESH — static model path (Raylib-supported formats), optional parentEntity#
- **`ENTITY.LOADANIMATEDMESH`** - args: any -> returns int
- **`ENTITY.LOADANIMATIONS`** - args: int, string
- **`ENTITY.LOADMESH`** - args: any -> returns int
- **`ENTITY.LOADSCENE`** - args: any
- **`ENTITY.LOADSPRITE`** - args: string -> returns int
- **`ENTITY.LOADSPRITE`** - args: string, int -> returns int
- **`ENTITY.LOOKAT`** - args: handle, float, float — Instantly rotates an entity to face a point.
- **`ENTITY.LOOKAT`** - args: int, float, float, float — Face world point (entity#, targetX#, targetY#, targetZ#); sets pitch/yaw
- **`ENTITY.MAKE`** - args: (none) -> returns int — DEPRECATED alias of ENTITY.CREATE. Use ENTITY.CREATE.
- **`ENTITY.MAKEBOX`** - args: float -> returns int — DEPRECATED alias of ENTITY.CREATEBOX. Use ENTITY.CREATEBOX. Uniform cube: size# used for width, height, and depth (alias ENTITY.CREATECUBE)
- **`ENTITY.MAKEBOX`** - args: float, float, float -> returns int — DEPRECATED alias of ENTITY.CREATEBOX. Use ENTITY.CREATEBOX.
- **`ENTITY.MAKECUBE`** - args: float, float, float -> returns int — DEPRECATED alias of ENTITY.CREATECUBE. Use ENTITY.CREATECUBE.
- **`ENTITY.MAKECYLINDER`** - args: float, float, int -> returns int — DEPRECATED alias of ENTITY.CREATECYLINDER. Use ENTITY.CREATECYLINDER.
- **`ENTITY.MAKEENTITY`** - args: (none) -> returns int — DEPRECATED alias of ENTITY.CREATEENTITY. Use ENTITY.CREATEENTITY.
- **`ENTITY.MAKEMESH`** - args: any -> returns int — DEPRECATED alias of ENTITY.CREATEMESH. Use ENTITY.CREATEMESH. Procedural mesh: optional parentEntity#; use AddVertex/UpdateMesh
- **`ENTITY.MAKEPLANE`** - args: float -> returns int — DEPRECATED alias of ENTITY.CREATEPLANE. Use ENTITY.CREATEPLANE.
- **`ENTITY.MAKESPHERE`** - args: float -> returns int — DEPRECATED alias of ENTITY.CREATESPHERE. Use ENTITY.CREATESPHERE. Radius only — default 16 segments
- **`ENTITY.MAKESPHERE`** - args: float, int -> returns int — DEPRECATED alias of ENTITY.CREATESPHERE. Use ENTITY.CREATESPHERE.
- **`ENTITY.MAKESPRITE`** - args: string -> returns int — DEPRECATED alias of ENTITY.CREATESPRITE. Use ENTITY.CREATESPRITE.
- **`ENTITY.MAKESPRITE`** - args: string, int -> returns int — DEPRECATED alias of ENTITY.CREATESPRITE. Use ENTITY.CREATESPRITE.
- **`ENTITY.MAKESPRITE`** - args: handle, float, float -> returns int — DEPRECATED alias of ENTITY.CREATESPRITE. Use ENTITY.CREATESPRITE. Billboard from TEXTURE handle (atlas / TEXTURE.LOADANIM)
- **`ENTITY.MAKESPRITE`** - args: handle, float, float, int -> returns int — DEPRECATED alias of ENTITY.CREATESPRITE. Use ENTITY.CREATESPRITE.
- **`ENTITY.MAKESURFACE`** - args: int -> returns handle — DEPRECATED alias of ENTITY.CREATESURFACE. Use ENTITY.CREATESURFACE.
- **`ENTITY.MATRIXELEMENT`** - args: int, int, int -> returns float
- **`ENTITY.MOVE`** - args: int, float, float, float
- **`ENTITY.MOVECAMERARELATIVE`** - args: int, float, float, handle — World XZ step from camera yaw: forward#/strafe# are deltas (typically speed*dt*input); camera is a Camera3D handle.
- **`ENTITY.MOVEENTITY`** - args: int, float, float, float
- **`ENTITY.MOVERELATIVE`** - args: int, float, float, float, float
- **`ENTITY.MOVETOWARD`** - args: handle, handle, float — Moves an entity toward another entity at constant speed (XZ toward target, Y preserved).
- **`ENTITY.MOVETOWARD`** - args: handle, float, float, float — Moves an entity toward a coordinate.
- **`ENTITY.MOVEWITHCAMERA`** - args: int, handle, float, float, float — Horizontal walk velocity (units/s) from camera XZ strafe basis (eye→target on ground). forwardAxis/strafeAxis are typically Input.Axis −1..1; preserves vertical velocity. Dot: player.MoveWithCamera(cam, …).
- **`ENTITY.ONHIT`** - args: handle, string — Fires MB callback on collision.
- **`ENTITY.ORDER`** - args: int, int
- **`ENTITY.OUTLINE`** - args: int, float, handle — Apply a highlighted outline effect to a model.
- **`ENTITY.P`** - args: int -> returns float — Easy Mode: Get Pitch of entity
- **`ENTITY.P`** - args: int, float — Easy Mode: Set Pitch of entity
- **`ENTITY.PARENT`** - args: int, int — Attach child entity to parent (optional third arg: global preserve world position; default true)
- **`ENTITY.PARENT`** - args: int, int, any
- **`ENTITY.PARENTCLEAR`** - args: int
- **`ENTITY.PATROL`** - args: handle, handle, float — Loops an entity across a WAYPOINT array handle at speed.
- **`ENTITY.PHYSICS`** - args: int, string, float — Quickly setup a physics body for an entity (auto-sizes based on model/shape).
- **`ENTITY.PHYSICS`** - args: int, string, float, float, float — Quickly setup a physics body with mass, friction, and restitution.
- **`ENTITY.PHYSICS`** - args: int, string, float, float, float, bool — Quickly setup a physics body with mass, friction, restitution, and CCD enabled.
- **`ENTITY.PHYSICSMOTION`** - args: int, string — Toggle physics motion type (STATIC, DYNAMIC, KINEMATIC).
- **`ENTITY.PICK`** - args: int, float -> returns bool
- **`ENTITY.PICKMODE`** - args: int, int
- **`ENTITY.PLAY`** - args: int, any
- **`ENTITY.PLAYNAME`** - args: int, string
- **`ENTITY.POINTAT`** - args: int, int
- **`ENTITY.POINTENTITY`** - args: int, int
- **`ENTITY.POLLMESSAGE`** - args: int -> returns string
- **`ENTITY.POS`** - args: int, float, float, float — Easy Mode shorthand for positioning an entity
- **`ENTITY.POSITION`** - args: int, float, float, float, any — Alias of ENTITY.SETPOS — set world or local position
- **`ENTITY.POSITIONENTITY`** - args: int, float, float, float, any
- **`ENTITY.PUSH`** - args: int, float, float, float — Apply Jolt impulse (requires ENTITY.ADDPHYSICS)
- **`ENTITY.PUSHOUTOFGEOMETRY`** - args: int — Best-effort depenetration: nudges entity world Y up slightly
- **`ENTITY.R`** - args: int -> returns float — Easy Mode: Get Roll of entity
- **`ENTITY.R`** - args: int, float — Easy Mode: Set Roll of entity
- **`ENTITY.RADIUS`** - args: int, float
- **`ENTITY.RAYCAST`** - args: handle, float -> returns handle — Raycast sensor
- **`ENTITY.RAYCAST`** - args: float, float, float, float, float, float, float -> returns int — Jolt ray cast; returns first hit entity# or 0 (same query path as PHYSICS3D/PICK)
- **`ENTITY.RAYHIT`** - args: int, float, float, float, float, float, float -> returns bool
- **`ENTITY.RGB`** - args: int, int, int, int — Easy Mode: Set entity color (id, r, g, b)
- **`ENTITY.ROT`** - args: int, float, float, float — Easy Mode shorthand for rotating an entity (absolute)
- **`ENTITY.ROTATE`** - args: int, float, float, float
- **`ENTITY.ROTATEENTITY`** - args: int, float, float, float, any
- **`ENTITY.SAVESCENE`** - args: any
- **`ENTITY.SCA`** - args: int, float, float, float — Easy Mode shorthand for scaling an entity (absolute)
- **`ENTITY.SCALE`** - args: int, float, float, float
- **`ENTITY.SCROLLMATERIAL`** - args: int, float, float — Add (du,dv) to material 0 scroll (same as MODEL.SCROLLTEXTURE)
- **`ENTITY.SENDMESSAGE`** - args: int, string
- **`ENTITY.SETANIMATION`** - args: int, int, float — Second arg 0 clears image-sequence animation
- **`ENTITY.SETANIMATION`** - args: int, handle, float — Cycle IMAGE.LOADSEQUENCE/LOADGIF frames onto sprite texture at fps
- **`ENTITY.SETANIMATION`** - args: int, handle, float, bool
- **`ENTITY.SETANIMFRAME`** - args: int, float
- **`ENTITY.SETANIMINDEX`** - args: int, any
- **`ENTITY.SETANIMLOOP`** - args: int, any
- **`ENTITY.SETANIMSPEED`** - args: int, float
- **`ENTITY.SETANIMTIME`** - args: int, float
- **`ENTITY.SETBOUNCE`** - args: int, float
- **`ENTITY.SETBOUNCINESS`** - args: int, float — Sets restitution (bounciness) on an entity's Jolt body; 0 = no bounce. Alias of PHYSICS.BOUNCE.
- **`ENTITY.SETBUOYANCY`** - args: int, float — Alias of PHYSICS.SETBUOYANCY — per-entity density hint for buoyancy
- **`ENTITY.SETCOLLISIONGROUP`** - args: int, int — Alias for ENTITY.COLLISIONLAYER (collision group / layer 0..31)
- **`ENTITY.SETDETAILTEXTURE`** - args: int, handle — Bind secondary map as MATERIAL_MAP_NORMAL for blending/detail
- **`ENTITY.SETFRICTION`** - args: int, float
- **`ENTITY.SETGRAVITY`** - args: int, float
- **`ENTITY.SETMASS`** - args: int, float
- **`ENTITY.SETNAME`** - args: int, any
- **`ENTITY.SETPOS`** - args: int, float, float, float, any
- **`ENTITY.SETPOSITION`** - args: int, float, float, float, any — DEPRECATED alias of ENTITY.SETPOS. Use ENTITY.SETPOS. Deprecated alias of ENTITY.SETPOS — set world or local position
- **`ENTITY.SETROTATION`** - args: int, float, float, float, any — Absolute pitch/yaw/roll degrees — alias of ENTITY.ROTATEENTITY
- **`ENTITY.SETSHADER`** - args: handle, int — Binds an active Shader Library component to the entity.
- **`ENTITY.SETSHADER`** - args: int, handle
- **`ENTITY.SETSPRITEFRAME`** - args: int, int — Set atlas frame on billboard bound to a TEXTURE object
- **`ENTITY.SETSTATIC`** - args: int
- **`ENTITY.SETTEXTUREFLIP`** - args: handle, float, float — Modifies UV scaling for horizontal/vertical mirroring.
- **`ENTITY.SETTEXTUREMAP`** - args: int, any, handle
- **`ENTITY.SETTEXTURESCROLL`** - args: handle, float, float — Injects offsets into the shader for animated water/lava.
- **`ENTITY.SETTRIGGER`** - args: int
- **`ENTITY.SETVISIBLE`** - args: int, any — Alias of ENTITY.VISIBLE
- **`ENTITY.SETWEIGHT`** - args: handle, float — Changes entity mass.
- **`ENTITY.SHININESS`** - args: int, float
- **`ENTITY.SHOW`** - args: int
- **`ENTITY.SLIDE`** - args: int, any
- **`ENTITY.SNAPTO`** - args: int, int — Instantly align one entity to another's transform.
- **`ENTITY.SPRITEVIEWMODE`** - args: int, int
- **`ENTITY.SQUASH`** - args: int, float, float — Juice: squash scale Y then tween back
- **`ENTITY.STOPANIM`** - args: int
- **`ENTITY.TAG`** - args: handle, string — Sets spatial tag.
- **`ENTITY.TEXTURE`** - args: int, any
- **`ENTITY.TFORMPOINT`** - args: float, float, float, int, int -> returns handle
- **`ENTITY.TFORMVECTOR`** - args: float, float, float, int, int -> returns handle
- **`ENTITY.TRANSITION`** - args: int, string, float
- **`ENTITY.TRANSLATE`** - args: int, float, float, float
- **`ENTITY.TRANSLATEENTITY`** - args: int, float, float, float, any
- **`ENTITY.TURN`** - args: int, float, float, float — Add pitch/yaw/roll degrees — alias of ENTITY.ROTATE / TURNENTITY
- **`ENTITY.TURNENTITY`** - args: int, float, float, float, any
- **`ENTITY.TURNTOWARD`** - args: handle, float, float, float — Slowly rotates the entity to face a target over time.
- **`ENTITY.TWEEN`** - args: int, string, any, float, string — Animate properties (position, scale, rotation) using easing functions (bounce, elastic, etc).
- **`ENTITY.TYPE`** - args: int, int
- **`ENTITY.UNPARENT`** - args: int — Alias of ENTITY.PARENTCLEAR — detach and keep world position
- **`ENTITY.UPDATE`** - args: float
- **`ENTITY.UPDATEMESH`** - args: int
- **`ENTITY.VELOCITY`** - args: int, float, float, float
- **`ENTITY.VERTEXX`** - args: handle, int -> returns float
- **`ENTITY.VERTEXY`** - args: handle, int -> returns float
- **`ENTITY.VERTEXZ`** - args: handle, int -> returns float
- **`ENTITY.VISIBLE`** - args: int, any
- **`ENTITY.W`** - args: int -> returns float — Easy Mode: Get Yaw (W) of entity
- **`ENTITY.W`** - args: int, float — Easy Mode: Set Yaw (W) of entity
- **`ENTITY.WANDER`** - args: handle, float, float, float, float — Moves an NPC randomly within a zone.
- **`ENTITY.WITHINRADIUS`** - args: handle, handle, float -> returns bool — True if 3D distance between entities is <= maxDistance (simple sphere check; not Jolt physics).
- **`ENTITY.X`** - args: int -> returns float — Easy Mode: Get X position of entity
- **`ENTITY.X`** - args: int, float — Easy Mode: Set X position of entity
- **`ENTITY.Y`** - args: int -> returns float
- **`ENTITY.Y`** - args: int, float
- **`ENTITY.Z`** - args: int -> returns float
- **`ENTITY.Z`** - args: int, float
- **`PHYSICS.CCD`** - args: int, bool — Enable Continuous Collision Detection to prevent high-speed tunneling.

### ENTITYALPHA

- **`ENTITYALPHA`** - args: int, float — Blitz-style: ENTITY.ALPHA(obj, alpha)
- **`ENTITYALPHA`** - args: handle, float — Professional: Set entity transparency (0.0=Invisible, 1.0=Solid).

### ENTITYBLEND

- **`ENTITYBLEND`** - args: handle, int — Professional: Set entity blend mode (0=Alpha, 1=Additive, 2=Multiply).

### ENTITYCOLLIDED

- **`ENTITYCOLLIDED`** - args: handle, int -> returns int — Easy Mode: Check if entity hit a specific type; returns handle of hit entity or 0

### ENTITYCOLOR

- **`ENTITYCOLOR`** - args: int, int, int, int — Blitz-style: ENTITY.COLOR(obj, r, g, b)
- **`ENTITYCOLOR`** - args: int, int, int, int, int — Easy Mode: ENTITY.COLOR(ent, r, g, b, a)

### ENTITYFLOOR

- **`ENTITYFLOOR`** - args: int -> returns bool — Easy Mode: Check if entity is on the floor

### ENTITYJUMP

- **`ENTITYJUMP`** - args: int, float — Easy Mode: Apply jump force to entity

### ENTITYPHYSICSTOUCH

- **`ENTITYPHYSICSTOUCH`** - args: int, int -> returns bool — Alias for EntityCollided

### ENTITYPITCH

- **`ENTITYPITCH`** - args: handle -> returns float — Easy Mode: Get entity Pitch orientation

### ENTITYRADIUS

- **`ENTITYRADIUS`** - args: handle, float -> returns void — Easy Mode: Set sphere collision radius for an entity

### ENTITYROLL

- **`ENTITYROLL`** - args: handle -> returns float — Easy Mode: Get entity Roll orientation

### ENTITYSHININESS

- **`ENTITYSHININESS`** - args: handle, float — Professional: Set entity specular highlight intensity.

### ENTITYTEXTURE

- **`ENTITYTEXTURE`** - args: handle, handle -> returns void — Easy Mode: Apply a texture handle to an entity

### ENTITYTYPE

- **`ENTITYTYPE`** - args: handle, int -> returns void — Easy Mode: Set collision group (1-32) for an entity

### ENTITYX

- **`ENTITYX`** - args: handle -> returns float — Easy Mode: Get entity X position

### ENTITYY

- **`ENTITYY`** - args: handle -> returns float — Easy Mode: Get entity Y position

### ENTITYYAW

- **`ENTITYYAW`** - args: handle -> returns float — Easy Mode: Get entity Yaw orientation

### ENTITYZ

- **`ENTITYZ`** - args: handle -> returns float — Easy Mode: Get entity Z position

### ENTPITCH

- **`ENTPITCH`** - args: handle -> returns float — Shorthand: ENTITYPITCH(ent)

### ENTRAD

- **`ENTRAD`** - args: int, float — Easy Mode: Set entity collision radius
- **`ENTRAD`** - args: handle, float — Shorthand: ENTITYRADIUS(ent, r)

### ENTROLL

- **`ENTROLL`** - args: handle -> returns float — Shorthand: ENTITYROLL(ent)

### ENTTYPE

- **`ENTTYPE`** - args: int, int — Easy Mode: Set entity collision type
- **`ENTTYPE`** - args: handle, int — Shorthand: ENTITYTYPE(ent, type)

### ENTX

- **`ENTX`** - args: int -> returns float
- **`ENTX`** - args: handle -> returns float — Shorthand: ENTITYX(ent)

### ENTY

- **`ENTY`** - args: int -> returns float
- **`ENTY`** - args: handle -> returns float — Shorthand: ENTITYY(ent)

### ENTYAW

- **`ENTYAW`** - args: handle -> returns float — Shorthand: ENTITYYAW(ent)

### ENTZ

- **`ENTZ`** - args: int -> returns float
- **`ENTZ`** - args: handle -> returns float — Shorthand: ENTITYZ(ent)

### ENVIRON

- **`ENVIRON`** - args: string

### ENVIRON$

- **`ENVIRON$`** - args: string

### EOF

- **`EOF`** - args: handle

### ERASE

- **`ERASE`** - args: handle

### ERR

- **`ERR`** - args: (none)

### ERRFILE$

- **`ERRFILE$`** - args: (none)

### ERRLINE

- **`ERRLINE`** - args: (none)

### ERRMSG$

- **`ERRMSG$`** - args: (none)

### EVENT

- **`EVENT.CHANNEL`** - args: handle -> returns int
- **`EVENT.DATA`** - args: handle -> returns string
- **`EVENT.FIRE`** - args: string
- **`EVENT.FIRE`** - args: string, any
- **`EVENT.FIRE`** - args: string, any, any
- **`EVENT.FIRE`** - args: string, any, any, any
- **`EVENT.FIRE`** - args: string, any, any, any, any
- **`EVENT.FIRE`** - args: string, any, any, any, any, any
- **`EVENT.FIRE`** - args: string, any, any, any, any, any, any
- **`EVENT.FIRE`** - args: string, any, any, any, any, any, any, any
- **`EVENT.FREE`** - args: handle
- **`EVENT.ISPLAYING`** - args: (none) -> returns bool
- **`EVENT.LISTCLEAR`** - args: handle
- **`EVENT.LISTCOUNT`** - args: handle -> returns int
- **`EVENT.LISTEXPORT`** - args: handle, string
- **`EVENT.LISTFREE`** - args: handle
- **`EVENT.LISTLOAD`** - args: string -> returns handle
- **`EVENT.LISTMAKE`** - args: string -> returns handle
- **`EVENT.OFF`** - args: string, string
- **`EVENT.ON`** - args: string, string
- **`EVENT.ONCE`** - args: string, string
- **`EVENT.PEER`** - args: handle -> returns handle
- **`EVENT.RECPLAYING`** - args: (none) -> returns bool
- **`EVENT.RECSTART`** - args: (none)
- **`EVENT.RECSTOP`** - args: (none)
- **`EVENT.REPLAY`** - args: handle
- **`EVENT.SETACTIVELIST`** - args: handle
- **`EVENT.TYPE`** - args: handle -> returns int

### EXP

- **`EXP`** - args: any

### EmitSound

- **`EmitSound`** - args: handle, int

### EntityAnimTime

- **`EntityAnimTime`** - args: int -> returns float

### EntityAnimateToward

- **`EntityAnimateToward`** - args: int, float, float, float, float — Alias for ENTITY.ANIMATETOWARD

### EntityApplyImpulse

- **`EntityApplyImpulse`** - args: int, float, float, float — Alias for ENTITY.APPLYIMPULSE

### EntityCanSee

- **`EntityCanSee`** - args: int, int, float, float -> returns bool — Alias for ENTITY.CANSEE

### EntityCheckCollision

- **`EntityCheckCollision`** - args: int, int -> returns bool — Alias for ENTITY.CHECKCOLLISION

### EntityCollided

- **`EntityCollided`** - args: int, int -> returns bool — True if two entities had a Jolt contact since last PHYSICS3D.STEP (Linux+CGO; link via ENTITY.LINKPHYSBUFFER)

### EntityCollisionLayer

- **`EntityCollisionLayer`** - args: int, int — Alias for ENTITY.COLLISIONLAYER

### EntityEmission

- **`EntityEmission`** - args: int, handle, float

### EntityFriction

- **`EntityFriction`** - args: int, float

### EntityGetClosestWithTag

- **`EntityGetClosestWithTag`** - args: int, float, string -> returns int — Alias for ENTITY.GETCLOSESTWITHTAG

### EntityGetGroundNormal

- **`EntityGetGroundNormal`** - args: int -> returns handle — Alias for ENTITY.GETGROUNDNORMAL

### EntityGetOverlapCount

- **`EntityGetOverlapCount`** - args: int, string -> returns int — Alias for ENTITY.GETOVERLAPCOUNT

### EntityGrounded

- **`EntityGrounded`** - args: int -> returns bool — True when the entity has floor support or is within coyote frames after leaving the ground (same as ENTITY.GROUNDED).

### EntityHasTag

- **`EntityHasTag`** - args: int, string -> returns bool — Alias for ENTITY.HASTAG

### EntityHitsType

- **`EntityHitsType`** - args: int, int -> returns bool — Args: (entity#, type#). True after ENTITY.UPDATE if entity# has a rule-based hit against another body whose EntityType equals type# (requires COLLISIONS + EntityType). For other entity id use ENTITYCOLLIDED; for Jolt pair test use EntityCollided(a#, b#).

### EntityInFrustum

- **`EntityInFrustum`** - args: int -> returns bool — Alias for ENTITY.INFRUSTUM

### EntityLineOfSight

- **`EntityLineOfSight`** - args: int, int -> returns bool — Alias for ENTITY.LINEOFSIGHT

### EntityMass

- **`EntityMass`** - args: int, float

### EntityMoveCameraRelative

- **`EntityMoveCameraRelative`** - args: int, float, float, handle — Same as ENTITY.MOVECAMERARELATIVE.

### EntityNormalMap

- **`EntityNormalMap`** - args: int, handle

### EntityPBR

- **`EntityPBR`** - args: int, float, float

### EntityPushOutOfGeometry

- **`EntityPushOutOfGeometry`** - args: int — Alias for ENTITY.PUSHOUTOFGEOMETRY

### EntityRaycast

- **`EntityRaycast`** - args: float, float, float, float, float, float, float -> returns int — Alias for ENTITY.RAYCAST

### EntityRestitution

- **`EntityRestitution`** - args: int, float

### EntitySetCollisionGroup

- **`EntitySetCollisionGroup`** - args: int, int — Alias for ENTITY.SETCOLLISIONGROUP

### EntityShadow

- **`EntityShadow`** - args: int, any

### ExtractAnimSeq

- **`ExtractAnimSeq`** - args: int, any, any

### FILE

- **`FILE.CLOSE`** - args: handle
- **`FILE.EOF`** - args: handle
- **`FILE.EXISTS`** - args: any -> returns bool
- **`FILE.OPEN`** - args: string, string
- **`FILE.READALLTEXT`** - args: any -> returns string
- **`FILE.READLINE`** - args: handle
- **`FILE.SEEK`** - args: handle, int
- **`FILE.SIZE`** - args: handle -> returns int
- **`FILE.TELL`** - args: handle -> returns int
- **`FILE.WRITE`** - args: handle, string — Write string to file without appending a newline.
- **`FILE.WRITEALLTEXT`** - args: any, any
- **`FILE.WRITELN`** - args: handle, string — Write string to file and append a newline.

### FILEEXISTS

- **`FILEEXISTS`** - args: string
- **`FILEEXISTS`** - args: string -> returns bool

### FILEPOS

- **`FILEPOS`** - args: handle

### FILESIZE

- **`FILESIZE`** - args: handle

### FIX

- **`FIX`** - args: any

### FLAT

- **`FLAT`** - args: float, float, float, float, float, int, int, int, int — alias of DRAW3D.PLANE — horizontal plane patch

### FLOAT

- **`FLOAT`** - args: any

### FLOOR

- **`FLOOR`** - args: any

### FOG

- **`FOG.ENABLE`** - args: bool
- **`FOG.SETCOLOR`** - args: int, int, int, int
- **`FOG.SETFAR`** - args: float
- **`FOG.SETNEAR`** - args: float
- **`FOG.SETRANGE`** - args: float, float

### FOGCOLOR

- **`FOGCOLOR`** - args: int, int, int — Environmental: Set global atmospheric haze color.

### FOGDENSITY

- **`FOGDENSITY`** - args: float — Environmental: Set thickness for exponential fog modes.

### FOGMODE

- **`FOGMODE`** - args: int — Environmental: Enable fog (0=Off, 1=Linear, 2=Exp, 3=Exp2).

### FONT

- **`FONT.DRAWDEFAULT`** - args: (none)
- **`FONT.FREE`** - args: handle
- **`FONT.LOAD`** - args: string
- **`FONT.LOADBDF`** - args: string, int

### FORMAT

- **`FORMAT`** - args: any, string -> returns string — Format a value with a printf-style pattern (canonical; same as legacy FORMAT$).

### FORMAT$

- **`FORMAT$`** - args: any, string -> returns string — DEPRECATED alias of FORMAT. Use FORMAT(value, pattern).

### FPS

- **`FPS`** - args: (none) -> returns int — Easy Mode: Get current frames per second

### FREE

- **`FREE.ALL`** - args: (none)

### FREEENTITIES

- **`FREEENTITIES`** - args: handle — Free every entity# stored in a numeric array (DIM badGuy AS HANDLE(n))

### FREEENTITY

- **`FREEENTITY`** - args: int — Blitz-style: ENTITY.FREE(obj)
- **`FREEENTITY`** - args: handle -> returns void — Easy Mode: Destroy an entity and free memory

### FREESOUND

- **`FREESOUND`** - args: handle -> returns void — Easy Mode: Free a sound asset

### FREETEXTURE

- **`FREETEXTURE`** - args: handle

### FindBone

- **`FindBone`** - args: int, any -> returns int

### GAME

- **`GAME.ANYKEY`** - args: (none) -> returns bool
- **`GAME.DT`** - args: (none) -> returns float
- **`GAME.ENDGAME`** - args: (none)
- **`GAME.FPS`** - args: (none) -> returns int
- **`GAME.GETTIMESCALE`** - args: (none) -> returns float — Current time scale (0 stored reads as 1 for delta)
- **`GAME.JOYBUTTON`** - args: int -> returns bool
- **`GAME.JOYX`** - args: (none) -> returns float
- **`GAME.JOYY`** - args: (none) -> returns float
- **`GAME.KEYCHAR`** - args: (none) -> returns int
- **`GAME.KEYDOWN`** - args: int -> returns bool
- **`GAME.KEYHIT`** - args: any -> returns bool
- **`GAME.KEYPRESSED`** - args: int -> returns bool
- **`GAME.KEYRELEASED`** - args: int -> returns bool
- **`GAME.MDX`** - args: (none) -> returns float
- **`GAME.MDY`** - args: (none) -> returns float
- **`GAME.MLEFT`** - args: (none) -> returns bool
- **`GAME.MLEFTPRESSED`** - args: (none) -> returns bool
- **`GAME.MMIDDLE`** - args: (none) -> returns bool
- **`GAME.MOUSEX`** - args: (none) -> returns int
- **`GAME.MOUSEXSPEED`** - args: (none) -> returns float
- **`GAME.MOUSEY`** - args: (none) -> returns int
- **`GAME.MOUSEYSPEED`** - args: (none) -> returns float
- **`GAME.MRIGHT`** - args: (none) -> returns bool
- **`GAME.MRIGHTPRESSED`** - args: (none) -> returns bool
- **`GAME.MWHEEL`** - args: (none) -> returns float
- **`GAME.MX`** - args: (none) -> returns int
- **`GAME.MY`** - args: (none) -> returns int
- **`GAME.ORBITDISTDELTA`** - args: float -> returns float
- **`GAME.ORBITPITCHDELTA`** - args: float -> returns float
- **`GAME.ORBITYAWDELTA`** - args: float, float, int, int, float -> returns float
- **`GAME.SCREENCX`** - args: (none) -> returns float
- **`GAME.SCREENCY`** - args: (none) -> returns float
- **`GAME.SCREENH`** - args: (none) -> returns int
- **`GAME.SCREENW`** - args: (none) -> returns int
- **`GAME.SETAESTHETIC`** - args: int — Selects global visual profile.
- **`GAME.SETPAUSE`** - args: int — Freezes physics/animation timers.
- **`GAME.SETTIMESCALE`** - args: float — Scales frame delta (0 = treated as 1); use for slow-mo / fast-forward with GAME.DT and TIME.DELTA
- **`GAME.SLOWMOTION`** - args: float, float — Slows time for cinematic.

### GESTURE

- **`GESTURE.ENABLE`** - args: int
- **`GESTURE.GETDETECTED`** - args: (none)
- **`GESTURE.GETDRAGANGLE`** - args: (none)
- **`GESTURE.GETDRAGVECTORX`** - args: (none)
- **`GESTURE.GETDRAGVECTORY`** - args: (none)
- **`GESTURE.GETHOLDDURATION`** - args: (none)
- **`GESTURE.GETPINCHANGLE`** - args: (none)
- **`GESTURE.GETPINCHVECTORX`** - args: (none)
- **`GESTURE.GETPINCHVECTORY`** - args: (none)
- **`GESTURE.ISDETECTED`** - args: int

### GETCOLLISIONENTITY

- **`GETCOLLISIONENTITY`** - args: handle, int -> returns handle — Easy Mode: Get handle of Nth collision

### GETDIR

- **`GETDIR`** - args: (none)
- **`GETDIR`** - args: (none) -> returns string

### GETDIR$

- **`GETDIR$`** - args: (none)
- **`GETDIR$`** - args: (none) -> returns string

### GETDIRS

- **`GETDIRS`** - args: string
- **`GETDIRS`** - args: string -> returns string

### GETDIRS$

- **`GETDIRS$`** - args: string
- **`GETDIRS$`** - args: string -> returns string

### GETDROPPEDFILES

- **`GETDROPPEDFILES`** - args: (none)

### GETFILEEXT

- **`GETFILEEXT`** - args: string
- **`GETFILEEXT`** - args: string -> returns string

### GETFILEEXT$

- **`GETFILEEXT$`** - args: string
- **`GETFILEEXT$`** - args: string -> returns string

### GETFILEMODTIME

- **`GETFILEMODTIME`** - args: string
- **`GETFILEMODTIME`** - args: string -> returns int

### GETFILENAME

- **`GETFILENAME`** - args: string
- **`GETFILENAME`** - args: string -> returns string

### GETFILENAME$

- **`GETFILENAME$`** - args: string
- **`GETFILENAME$`** - args: string -> returns string

### GETFILENAMENOEXT

- **`GETFILENAMENOEXT`** - args: string
- **`GETFILENAMENOEXT`** - args: string -> returns string

### GETFILENAMENOEXT$

- **`GETFILENAMENOEXT$`** - args: string
- **`GETFILENAMENOEXT$`** - args: string -> returns string

### GETFILEPATH

- **`GETFILEPATH`** - args: string
- **`GETFILEPATH`** - args: string -> returns string

### GETFILEPATH$

- **`GETFILEPATH$`** - args: string
- **`GETFILEPATH$`** - args: string -> returns string

### GETFILES

- **`GETFILES`** - args: string
- **`GETFILES`** - args: string -> returns string

### GETFILES$

- **`GETFILES$`** - args: string
- **`GETFILES$`** - args: string -> returns string

### GETFILESIZE

- **`GETFILESIZE`** - args: string
- **`GETFILESIZE`** - args: string -> returns int

### GRAPHICS

- **`GRAPHICS`** - args: int, int — Blitz-style: WINDOW.OPEN(w, h, 'moonBASIC')
- **`GRAPHICS`** - args: int, int, string — Blitz-style: WINDOW.OPEN(w, h, title$)

### GRID

- **`GRID.CREATE`** - args: int, int, float -> returns handle — Logical XZ tactical grid (width x depth cells, cell size)
- **`GRID.DRAW`** - args: handle, int, int, int
- **`GRID.DRAW`** - args: handle, int, int, int, int
- **`GRID.FOLLOWTERRAIN`** - args: handle, handle — Bake per-cell Y from terrain height
- **`GRID.FREE`** - args: handle
- **`GRID.GETCELL`** - args: handle, int, int -> returns int
- **`GRID.GETNEIGHBORS`** - args: handle, int, int, int -> returns handle — Entity IDs occupying cells in Chebyshev radius
- **`GRID.GETPATH`** - args: handle, float, float, float, float -> returns handle — Packed path [ix0,iz0, ix1,iz1, ...] or empty
- **`GRID.MAKE`** - args: int, int, float -> returns handle — DEPRECATED alias of GRID.CREATE. Use GRID.CREATE. Logical XZ tactical grid (width x depth cells, cell size)
- **`GRID.PLACEENTITY`** - args: handle, int, int, int
- **`GRID.RAYCAST`** - args: handle, float, float -> returns handle — Cell under mouse ray on XZ plane; [-1,-1] if miss
- **`GRID.SETCELL`** - args: handle, int, int, int
- **`GRID.SNAP`** - args: handle, int, int, int — Move entity to cell center (optional Y from GRID.FOLLOWTERRAIN)
- **`GRID.WORLDTOCELL`** - args: handle, float, float -> returns handle — Array handle [ix, iz]

### GRID3

- **`GRID3`** - args: int, float — alias of DRAW3D.GRID — XZ reference grid
- **`GRID3`** - args: int, float, float — alias of DRAW3D.GRID with Y offset

### GUI

- **`GUI.BUTTON`** - args: float, float, float, float, string -> returns bool
- **`GUI.CHECKBOX`** - args: float, float, float, float, string, bool -> returns bool
- **`GUI.COLORBARALPHA`** - args: float, float, float, float, string, float -> returns float
- **`GUI.COLORBARHUE`** - args: float, float, float, float, string, float -> returns float
- **`GUI.COLORPANEL`** - args: float, float, float, float, string, int, int, int, int -> returns handle
- **`GUI.COLORPANELHSV`** - args: float, float, float, float, string, handle -> returns int
- **`GUI.COLORPICKER`** - args: float, float, float, float, string, int, int, int, int -> returns handle
- **`GUI.COLORPICKERHSV`** - args: float, float, float, float, string, handle -> returns int
- **`GUI.COMBOBOX`** - args: float, float, float, float, string, int -> returns int
- **`GUI.DISABLE`** - args: (none)
- **`GUI.DISABLETOOLTIP`** - args: (none)
- **`GUI.DRAWICON`** - args: int, int, int, int, int, int, int, int
- **`GUI.DRAWRECTANGLE`** - args: float, float, float, float, int, int, int, int, int, int, int, int, int
- **`GUI.DRAWTEXT`** - args: string, float, float, float, float, int, int, int, int, int
- **`GUI.DROPDOWNBOX`** - args: float, float, float, float, string, handle -> returns bool
- **`GUI.DUMMYREC`** - args: float, float, float, float, string
- **`GUI.ENABLE`** - args: (none)
- **`GUI.ENABLETOOLTIP`** - args: (none)
- **`GUI.FADE`** - args: int, int, int, int, float -> returns handle
- **`GUI.GETCOLOR`** - args: int, int -> returns handle
- **`GUI.GETSTATE`** - args: (none) -> returns int
- **`GUI.GETSTYLE`** - args: int, int -> returns int
- **`GUI.GETTEXTBOUNDS`** - args: int, float, float, float, float -> returns handle
- **`GUI.GETTEXTSIZE`** - args: (none) -> returns int
- **`GUI.GETTEXTWIDTH`** - args: string -> returns int
- **`GUI.GRID`** - args: float, float, float, float, string, float, int, handle -> returns int
- **`GUI.GROUPBOX`** - args: float, float, float, float, string
- **`GUI.ICONTEXT`** - args: int, string -> returns string
- **`GUI.ISLOCKED`** - args: (none) -> returns bool
- **`GUI.LABEL`** - args: float, float, float, float, string
- **`GUI.LABELBUTTON`** - args: float, float, float, float, string -> returns bool
- **`GUI.LINE`** - args: float, float, float, float, string
- **`GUI.LISTVIEW`** - args: float, float, float, float, string, handle -> returns int
- **`GUI.LISTVIEWEX`** - args: float, float, float, float, string, handle -> returns int
- **`GUI.LOADDEFAULTSTYLE`** - args: (none)
- **`GUI.LOADICONS`** - args: string, bool
- **`GUI.LOADICONSMEM`** - args: string, bool
- **`GUI.LOADSTYLE`** - args: string
- **`GUI.LOADSTYLEMEM`** - args: string
- **`GUI.LOCK`** - args: (none)
- **`GUI.MESSAGEBOX`** - args: float, float, float, float, string, string, string -> returns int
- **`GUI.PANEL`** - args: float, float, float, float, string
- **`GUI.PROGRESSBAR`** - args: float, float, float, float, string, string, float, float, float -> returns float
- **`GUI.SCROLLBAR`** - args: float, float, float, float, int, int, int -> returns int
- **`GUI.SCROLLPANEL`** - args: float, float, float, float, string, float, float, float, float, handle
- **`GUI.SETALPHA`** - args: float
- **`GUI.SETCOLOR`** - args: int, int, int, int, int, int
- **`GUI.SETFONT`** - args: handle
- **`GUI.SETICONSCALE`** - args: int
- **`GUI.SETSTATE`** - args: int
- **`GUI.SETSTYLE`** - args: int, int, int
- **`GUI.SETTEXTALIGN`** - args: int
- **`GUI.SETTEXTALIGNVERT`** - args: int
- **`GUI.SETTEXTLINEHEIGHT`** - args: int
- **`GUI.SETTEXTSIZE`** - args: int
- **`GUI.SETTEXTSPACING`** - args: int
- **`GUI.SETTEXTWRAP`** - args: int
- **`GUI.SETTOOLTIP`** - args: string
- **`GUI.SLIDER`** - args: float, float, float, float, string, string, float, float, float -> returns float
- **`GUI.SLIDERBAR`** - args: float, float, float, float, string, string, float, float, float -> returns float
- **`GUI.SPINNER`** - args: float, float, float, float, string, int, int, int, bool -> returns int
- **`GUI.STATUSBAR`** - args: float, float, float, float, string
- **`GUI.TABBAR`** - args: float, float, float, float, string, handle -> returns int
- **`GUI.TEXTBOX`** - args: float, float, float, float, string, int, bool -> returns string
- **`GUI.TEXTINPUTBOX`** - args: float, float, float, float, string, string, string, string, int, handle -> returns int
- **`GUI.TEXTINPUTLAST`** - args: (none) -> returns string
- **`GUI.TEXTINPUTLAST$`** - args: (none) -> returns string
- **`GUI.THEMEAPPLY`** - args: string
- **`GUI.THEMENAMES`** - args: (none) -> returns string
- **`GUI.THEMENAMES$`** - args: (none) -> returns string
- **`GUI.TOGGLE`** - args: float, float, float, float, string, bool -> returns bool
- **`GUI.TOGGLEGROUP`** - args: float, float, float, float, string -> returns int
- **`GUI.TOGGLEGROUPAT`** - args: float, float, float, float, string, int -> returns int
- **`GUI.TOGGLESLIDER`** - args: float, float, float, float, string, int -> returns int
- **`GUI.UNLOCK`** - args: (none)
- **`GUI.VALUEBOX`** - args: float, float, float, float, string, int, int, int, bool -> returns int
- **`GUI.VALUEBOXFLOAT`** - args: float, float, float, float, string, float, string, bool -> returns float
- **`GUI.VALUEBOXFLOATTEXT`** - args: (none) -> returns string
- **`GUI.VALUEBOXFLOATTEXT$`** - args: (none) -> returns string
- **`GUI.WINDOWBOX`** - args: float, float, float, float, string -> returns bool

### Graphics3D

- **`Graphics3D`** - args: int, int — Resize window (w,h) with defaults: reserved depth and high-DPI mode
- **`Graphics3D`** - args: int, int, int, int — Resize window (w,h); depth reserved; mode bit0 = high-DPI flag

### HDIST

- **`HDIST`** - args: float, float, float, float -> returns float — Horizontal distance on XZ: hypot(x2-x1, z2-z1); ignores Y

### HDISTSQ

- **`HDISTSQ`** - args: float, float, float, float -> returns float — Squared HDIST for comparisons without sqrt

### HELP

- **`HELP`** - args: string — Live Discovery: Show arguments and description for any command.

### HEX

- **`HEX`** - args: int -> returns string

### HEX$

- **`HEX$`** - args: int -> returns string

### HIDEENTITY

- **`HIDEENTITY`** - args: handle -> returns void — Easy Mode: Hide an entity

### HITCOUNT

- **`HITCOUNT`** - args: handle -> returns int — Shorthand: COUNTCOLLISIONS(ent)

### HITENT

- **`HITENT`** - args: handle, int -> returns handle — Shorthand: GETCOLLISIONENTITY(ent, index)

### HOUR

- **`HOUR`** - args: (none)
- **`HOUR`** - args: (none) -> returns int

### IIF

- **`IIF`** - args: any, any, any

### IIF$

- **`IIF$`** - args: any, any, any — Inline conditional returning a string; both branches evaluated

### IMAGE

- **`IMAGE.ALPHACLEAR`** - args: handle, int, int, int, int, float
- **`IMAGE.ALPHACROP`** - args: handle, float
- **`IMAGE.CLEAR`** - args: handle, int, int, int, int
- **`IMAGE.CLEARBACKGROUND`** - args: handle, int, int, int, int
- **`IMAGE.COLORBRIGHTNESS`** - args: handle, int
- **`IMAGE.COLORCONTRAST`** - args: handle, float
- **`IMAGE.COLORGRAYSCALE`** - args: handle
- **`IMAGE.COLORINVERT`** - args: handle
- **`IMAGE.COLORREPLACE`** - args: handle, int, int, int, int, int, int, int, int
- **`IMAGE.COLORTINT`** - args: handle, int, int, int, int
- **`IMAGE.COPY`** - args: handle
- **`IMAGE.CREATE`** - args: int, int
- **`IMAGE.CREATE`** - args: int, int, int, int, int, int
- **`IMAGE.CREATEBLANK`** - args: int, int
- **`IMAGE.CREATEBLANK`** - args: int, int, int, int, int, int
- **`IMAGE.CREATECOPY`** - args: handle
- **`IMAGE.CREATETEXT`** - args: string, int, int, int, int, int
- **`IMAGE.CROP`** - args: handle, int, int, int, int
- **`IMAGE.DITHER`** - args: handle, int, int, int, int
- **`IMAGE.DRAWCIRCLE`** - args: handle, int, int, int, int, int, int, int
- **`IMAGE.DRAWIMAGE`** - args: handle, handle, float, float, float, float, float, float, float, float, int, int, int, int
- **`IMAGE.DRAWLINE`** - args: handle, int, int, int, int, int, int, int, int
- **`IMAGE.DRAWPIXEL`** - args: handle, int, int, int, int, int, int
- **`IMAGE.DRAWRECT`** - args: handle, int, int, int, int, int, int, int, int
- **`IMAGE.DRAWRECTLINES`** - args: handle, float, float, float, float, int, int, int, int, int
- **`IMAGE.DRAWTEXT`** - args: handle, int, int, string, int, int, int, int, int
- **`IMAGE.EXPORT`** - args: handle, string
- **`IMAGE.FLIPH`** - args: handle
- **`IMAGE.FLIPV`** - args: handle
- **`IMAGE.FORMAT`** - args: handle, int
- **`IMAGE.FREE`** - args: handle
- **`IMAGE.GETBBOXH`** - args: handle, float
- **`IMAGE.GETBBOXW`** - args: handle, float
- **`IMAGE.GETBBOXX`** - args: handle, float
- **`IMAGE.GETBBOXY`** - args: handle, float
- **`IMAGE.GETCOLORA`** - args: handle, int, int
- **`IMAGE.GETCOLORB`** - args: handle, int, int
- **`IMAGE.GETCOLORG`** - args: handle, int, int
- **`IMAGE.GETCOLORR`** - args: handle, int, int
- **`IMAGE.GETPIXEL`** - args: handle, int, int -> returns int — Packed pixel color (host byte order; typically ARGB-style int)
- **`IMAGE.HEIGHT`** - args: handle
- **`IMAGE.LOAD`** - args: string
- **`IMAGE.LOADGIF`** - args: string -> returns handle — Animated GIF to ImageSequence (cumulative frames)
- **`IMAGE.LOADRAW`** - args: string, int, int, int, int
- **`IMAGE.LOADSEQUENCE`** - args: string -> returns handle — Glob files matching prefix (e.g. assets/water_*.png) sorted by name
- **`IMAGE.MAKE`** - args: int, int — DEPRECATED alias of IMAGE.CREATE. Use IMAGE.CREATE.
- **`IMAGE.MAKE`** - args: int, int, int, int, int, int — DEPRECATED alias of IMAGE.CREATE. Use IMAGE.CREATE.
- **`IMAGE.MAKEBLANK`** - args: int, int — DEPRECATED alias of IMAGE.CREATEBLANK. Use IMAGE.CREATEBLANK.
- **`IMAGE.MAKEBLANK`** - args: int, int, int, int, int, int — DEPRECATED alias of IMAGE.CREATEBLANK. Use IMAGE.CREATEBLANK.
- **`IMAGE.MAKECOPY`** - args: handle — DEPRECATED alias of IMAGE.CREATECOPY. Use IMAGE.CREATECOPY.
- **`IMAGE.MAKETEXT`** - args: string, int, int, int, int, int — DEPRECATED alias of IMAGE.CREATETEXT. Use IMAGE.CREATETEXT.
- **`IMAGE.MIPMAPS`** - args: handle
- **`IMAGE.RESIZE`** - args: handle, int, int
- **`IMAGE.RESIZENN`** - args: handle, int, int
- **`IMAGE.ROTATE`** - args: handle, int
- **`IMAGE.ROTATECCW`** - args: handle
- **`IMAGE.ROTATECW`** - args: handle
- **`IMAGE.SETFILTER`** - args: handle, int — Raylib texture filter applied on IMAGE.TOTEXTURE / TEXTURE.FROMIMAGE
- **`IMAGE.TOTEXTURE`** - args: handle -> returns handle — Alias of TEXTURE.FROMIMAGE; respects IMAGE.SETFILTER when set
- **`IMAGE.WIDTH`** - args: handle

### INPUT

- **`INPUT`** - args: string -> returns string
- **`INPUT.ACTIONAXIS`** - args: string -> returns float
- **`INPUT.ACTIONDOWN`** - args: string -> returns bool
- **`INPUT.ACTIONPRESSED`** - args: string -> returns bool
- **`INPUT.ACTIONRELEASED`** - args: string -> returns bool
- **`INPUT.AXIS`** - args: any, any -> returns float — Two-key axis: -1, 0, or 1 (negKey vs posKey)
- **`INPUT.AXISDEG`** - args: any, any, float, float -> returns float — Input.Axis(neg,pos)*DEGPERSEC(degPerSec,dt) — radians this frame
- **`INPUT.GAMEPADAXISCOUNT`** - args: int -> returns int
- **`INPUT.GAMEPADBUTTONCOUNT`** - args: int -> returns int
- **`INPUT.GETINACTIVITY`** - args: (none) -> returns float — Returns time in seconds since the last user interaction.
- **`INPUT.GETKEYNAME`** - args: int -> returns string
- **`INPUT.GETMOUSEWORLDPOS`** - args: handle, int, int -> returns handle
- **`INPUT.GETTOUCHPOINTID`** - args: int -> returns int
- **`INPUT.JOYBUTTON`** - args: int -> returns bool
- **`INPUT.JOYDOWN`** - args: any, any -> returns bool
- **`INPUT.JOYX`** - args: (none) -> returns float
- **`INPUT.JOYY`** - args: (none) -> returns float
- **`INPUT.KEYDOWN`** - args: int -> returns bool
- **`INPUT.KEYDOWN`** - args: any
- **`INPUT.KEYHIT`** - args: any -> returns bool
- **`INPUT.KEYPRESSED`** - args: any
- **`INPUT.KEYUP`** - args: any
- **`INPUT.KEYUP`** - args: int -> returns bool
- **`INPUT.LOADMAPPINGS`** - args: string
- **`INPUT.MAPGAMEPADAXIS`** - args: string, int, int
- **`INPUT.MAPGAMEPADBUTTON`** - args: string, int, int
- **`INPUT.MAPKEY`** - args: string, int
- **`INPUT.MOUSEDELTA`** - args: (none) -> returns handle
- **`INPUT.MOUSEDELTAX`** - args: (none) -> returns float
- **`INPUT.MOUSEDELTAY`** - args: (none) -> returns float
- **`INPUT.MOUSEDOWN`** - args: int
- **`INPUT.MOUSEDX`** - args: (none) -> returns float — Alias of INPUT.MOUSEDELTAX
- **`INPUT.MOUSEDY`** - args: (none) -> returns float — Alias of INPUT.MOUSEDELTAY
- **`INPUT.MOUSEHIT`** - args: int -> returns bool
- **`INPUT.MOUSEWHEEL`** - args: (none) -> returns float — Alias of INPUT.MOUSEWHEELMOVE
- **`INPUT.MOUSEWHEELMOVE`** - args: (none) -> returns float
- **`INPUT.MOUSEX`** - args: (none)
- **`INPUT.MOUSEXSPEED`** - args: (none) -> returns float
- **`INPUT.MOUSEY`** - args: (none)
- **`INPUT.MOUSEYSPEED`** - args: (none) -> returns float
- **`INPUT.MOVEDIR`** - args: float, float -> returns handle
- **`INPUT.MOVEMENT2D`** - args: any, any, any, any -> returns handle — 2-float array [forward, strafe] from two Axis pairs; ERASE when done
- **`INPUT.ORBIT`** - args: any, any, float, float -> returns float — Alias of INPUT.AXISDEG — orbit / yaw delta this frame
- **`INPUT.SAVEMAPPINGS`** - args: string
- **`INPUT.SETGAMEPADMAPPINGS`** - args: string -> returns int
- **`INPUT.SETMOUSEOFFSET`** - args: int, int
- **`INPUT.SETMOUSEPOS`** - args: int, int — Warp OS cursor to client pixel (x,y); pair with CURSOR.DISABLE for game-style recenter
- **`INPUT.SETMOUSESCALE`** - args: float, float
- **`INPUT.TOUCHCOUNT`** - args: (none) -> returns int
- **`INPUT.TOUCHPRESSED`** - args: int -> returns bool
- **`INPUT.TOUCHX`** - args: int -> returns int
- **`INPUT.TOUCHY`** - args: int -> returns int

### INSTANCE

- **`INSTANCE.COUNT`** - args: handle -> returns int
- **`INSTANCE.CREATE`** - args: handle, int -> returns handle
- **`INSTANCE.CREATEINSTANCED`** - args: string, int -> returns handle
- **`INSTANCE.DRAW`** - args: handle
- **`INSTANCE.DRAWLOD`** - args: handle, handle, float
- **`INSTANCE.FREE`** - args: handle
- **`INSTANCE.GETPOS`** - args: handle -> returns array
- **`INSTANCE.GETROT`** - args: handle -> returns array
- **`INSTANCE.GETSCALE`** - args: handle -> returns array
- **`INSTANCE.MAKE`** - args: handle, int -> returns handle — DEPRECATED alias of INSTANCE.CREATE. Use INSTANCE.CREATE.
- **`INSTANCE.MAKEINSTANCED`** - args: string, int -> returns handle — DEPRECATED alias of INSTANCE.CREATEINSTANCED. Use INSTANCE.CREATEINSTANCED.
- **`INSTANCE.SETCOLOR`** - args: handle, int, float, float, float, float
- **`INSTANCE.SETCULLDISTANCE`** - args: handle, float
- **`INSTANCE.SETINSTANCEPOS`** - args: handle, int, float, float, float
- **`INSTANCE.SETINSTANCESCALE`** - args: handle, int, float, float, float
- **`INSTANCE.SETMATRIX`** - args: handle, int, handle
- **`INSTANCE.SETPOS`** - args: handle, int, float, float, float
- **`INSTANCE.SETPOSITION`** - args: handle, int, float, float, float — DEPRECATED alias of INSTANCE.SETPOS. Use INSTANCE.SETPOS.
- **`INSTANCE.SETROT`** - args: handle, int, float, float, float
- **`INSTANCE.SETSCALE`** - args: handle, int, float, float, float
- **`INSTANCE.UPDATEBUFFER`** - args: handle
- **`INSTANCE.UPDATEINSTANCES`** - args: handle

### INSTR

- **`INSTR`** - args: string, string
- **`INSTR`** - args: string, string, int

### INT

- **`INT`** - args: any

### INTERP

- **`INTERP`** - args: string, any -> returns string
- **`INTERP`** - args: string, any, any -> returns string
- **`INTERP`** - args: string, any, any, any -> returns string
- **`INTERP`** - args: string, any, any, any, any -> returns string
- **`INTERP`** - args: string, any, any, any, any, any -> returns string
- **`INTERP`** - args: string, any, any, any, any, any, any -> returns string
- **`INTERP`** - args: string, any, any, any, any, any, any, any -> returns string
- **`INTERP`** - args: string, any, any, any, any, any, any, any, any -> returns string
- **`INTERP`** - args: string, any, any, any, any, any, any, any, any, any -> returns string
- **`INTERP`** - args: string, any, any, any, any, any, any, any, any, any, any -> returns string

### INTERP$

- **`INTERP$`** - args: string, any -> returns string
- **`INTERP$`** - args: string, any, any -> returns string
- **`INTERP$`** - args: string, any, any, any -> returns string
- **`INTERP$`** - args: string, any, any, any, any -> returns string
- **`INTERP$`** - args: string, any, any, any, any, any -> returns string
- **`INTERP$`** - args: string, any, any, any, any, any, any -> returns string
- **`INTERP$`** - args: string, any, any, any, any, any, any, any -> returns string
- **`INTERP$`** - args: string, any, any, any, any, any, any, any, any -> returns string
- **`INTERP$`** - args: string, any, any, any, any, any, any, any, any, any -> returns string
- **`INTERP$`** - args: string, any, any, any, any, any, any, any, any, any, any -> returns string

### INVERSE_LERP

- **`INVERSE_LERP`** - args: float, float, float -> returns float

### ISALPHA

- **`ISALPHA`** - args: string -> returns bool

### ISALPHANUM

- **`ISALPHANUM`** - args: string -> returns bool

### ISFILEDROPPED

- **`ISFILEDROPPED`** - args: (none)

### ISHANDLE

- **`ISHANDLE`** - args: any

### ISNULL

- **`ISNULL`** - args: any

### ISNUMERIC

- **`ISNUMERIC`** - args: string -> returns bool

### ISTYPE

- **`ISTYPE`** - args: any, string

### JOIN

- **`JOIN`** - args: handle, string -> returns string

### JOIN$

- **`JOIN$`** - args: handle, string -> returns string

### JOINT

- **`JOINT.CREATEHINGE`** - args: handle, handle, float, float, float, float, float, float — Create a hinge joint (b1, b2, px, py, pz, ax, ay, az).
- **`JOINT.MAKEHINGE`** - args: handle, handle, float, float, float, float, float, float — DEPRECATED alias of JOINT.CREATEHINGE. Use JOINT.CREATEHINGE. Create a hinge joint (b1, b2, px, py, pz, ax, ay, az).

### JOINT2D

- **`JOINT2D.DISTANCE`** - args: handle, handle, float, float, float, float -> returns handle
- **`JOINT2D.FREE`** - args: handle
- **`JOINT2D.PRISMATIC`** - args: handle, handle, float, float, float, float -> returns handle
- **`JOINT2D.REVOLUTE`** - args: handle, handle, float, float -> returns handle

### JOINT3D

- **`JOINT3D.CONE`** - args: handle, handle, float, float, float, float, float -> returns handle
- **`JOINT3D.DELETE`** - args: handle
- **`JOINT3D.FIXED`** - args: handle, handle -> returns handle
- **`JOINT3D.HINGE`** - args: handle, handle, float, float, float, float, float, float -> returns handle
- **`JOINT3D.SLIDER`** - args: handle, handle, float, float, float, float, float, float -> returns handle

### JOLT

- **`JOLT.BODYCREATEDYNAMIC`** - args: (none)
- **`JOLT.BODYCREATEKINEMATIC`** - args: (none)
- **`JOLT.BODYCREATESTATIC`** - args: (none)
- **`JOLT.COLLISIONQUERY`** - args: handle
- **`JOLT.CONSTRAINTDISTANCE`** - args: handle, handle
- **`JOLT.CONSTRAINTFIXED`** - args: handle, handle
- **`JOLT.CONSTRAINTHINGE`** - args: handle, handle
- **`JOLT.CONSTRAINTPOINT`** - args: handle, handle
- **`JOLT.CONSTRAINTSLIDER`** - args: handle, handle
- **`JOLT.INIT`** - args: (none)
- **`JOLT.RAYCAST`** - args: float, float, float, float, float, float
- **`JOLT.SETGRAVITY`** - args: float, float, float
- **`JOLT.SHAPEBOX`** - args: float, float, float
- **`JOLT.SHAPECAPSULE`** - args: float, float
- **`JOLT.SHAPECYLINDER`** - args: float, float
- **`JOLT.SHAPEMESH`** - args: handle
- **`JOLT.SHAPESPHERE`** - args: float
- **`JOLT.SHUTDOWN`** - args: (none)
- **`JOLT.STEP`** - args: float

### JSON

- **`JSON.CREATE`** - args: (none) -> returns handle
- **`JSON.FREE`** - args: handle
- **`JSON.GETBOOL`** - args: handle, string -> returns bool
- **`JSON.GETFLOAT`** - args: handle, string -> returns float
- **`JSON.GETINT`** - args: handle, string -> returns int
- **`JSON.GETSTRING`** - args: handle, string -> returns string
- **`JSON.LOADFILE`** - args: any -> returns handle
- **`JSON.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of JSON.CREATE. Use JSON.CREATE.
- **`JSON.PARSE`** - args: string -> returns handle
- **`JSON.PARSESTRING`** - args: string -> returns handle
- **`JSON.SAVEFILE`** - args: any, any
- **`JSON.SETBOOL`** - args: handle, string, bool
- **`JSON.SETFLOAT`** - args: handle, string, float
- **`JSON.SETINT`** - args: handle, string, int
- **`JSON.SETSTRING`** - args: handle, string, string
- **`JSON.TOSTRING`** - args: handle -> returns string

### KEEPPLAYERINBOUNDS

- **`KEEPPLAYERINBOUNDS`** - args: handle

### KEY

- **`KEY`** - args: (none) -> returns handle
- **`KEY.DOWN`** - args: handle, any -> returns bool

### KEYDOWN

- **`KEYDOWN`** - args: any -> returns bool — Easy Mode: KEY.DOWN(KEY(), code)

### KEYHIT

- **`KEYHIT`** - args: any -> returns bool — Easy Mode: KEY.HIT(KEY(), code)

### KEYUP

- **`KEYUP`** - args: any -> returns bool — Easy Mode: KEY.UP(KEY(), code)

### KINEMATIC

- **`KINEMATIC.CREATE`** - args: handle -> returns handle — Creates a Kinematic Body from a shape handle.
- **`KINEMATIC.MAKE`** - args: handle -> returns handle — DEPRECATED alias of KINEMATIC.CREATE. Use KINEMATIC.CREATE. Creates a Kinematic Body from a shape handle.

### KINEMATICREF

- **`KINEMATICREF.SETVELOCITY`** - args: handle, float, float, float — Sets velocity for moving kinematic bodies.
- **`KINEMATICREF.UPDATE`** - args: handle — Resolves kinematic movement and collisions.

### KeyDown

- **`KeyDown`** - args: any -> returns bool — Alias for KEYDOWN / INPUT.KEYDOWN

### LANDBOX

- **`LANDBOX`** - args: float, float, float, float, float, any, any, any, any, any, any, any -> returns float — Alias of LANDBOXES

### LANDBOXES

- **`LANDBOXES`** - args: float, float, float, float, float, any, any, any, any, any, any, any -> returns float — Best BOXTOPLAND snap Y over count boxes (parallel float arrays)

### LEFT

- **`LEFT`** - args: string, int -> returns string

### LEFT$

- **`LEFT$`** - args: string, int -> returns string

### LEN

- **`LEN`** - args: string

### LERP

- **`LERP`** - args: any, any, any

### LEVEL

- **`LEVEL.APPLYPHYSICS`** - args: int
- **`LEVEL.AUTOCOLLIDE`** - args: (none) — Scans all active entities and automatically creates static mesh collisions for those marked as static with models.
- **`LEVEL.BINDSCRIPT`** - args: string, string
- **`LEVEL.FINDENTITY`** - args: string -> returns int
- **`LEVEL.GETMARKER`** - args: string -> returns handle
- **`LEVEL.GETSPAWN`** - args: string -> returns handle
- **`LEVEL.LOAD`** - args: string -> returns int
- **`LEVEL.LOADSKYBOX`** - args: string -> returns handle
- **`LEVEL.MATCHSCRIPTBIND`** - args: string -> returns string
- **`LEVEL.OPTIMIZE`** - args: int
- **`LEVEL.PRELOAD`** - args: string -> returns int
- **`LEVEL.SETROOT`** - args: string
- **`LEVEL.SHOWLAYER`** - args: string, any
- **`LEVEL.STATIC`** - args: any — Creates a static mesh collision body from an entity's current model. Use for optimized level geometry.
- **`LEVEL.SYNCLIGHTS`** - args: any

### LIGHT

- **`LIGHT.CREATE`** - args: (none) -> returns handle
- **`LIGHT.CREATE`** - args: string -> returns handle
- **`LIGHT.CREATEDIRECTIONAL`** - args: float, float, float, float, float, float, float -> returns handle — Directional light: direction vector (dx,dy,dz), RGB, energy — direction is normalized
- **`LIGHT.CREATEPOINT`** - args: float, float, float, float, float, float, float -> returns handle — Point light at (x,y,z) with RGB (0–255 or 0–1) and intensity (energy)
- **`LIGHT.CREATESPOT`** - args: float, float, float, float, float, float, float, float, float, float, float -> returns handle — Spot: position, target point, RGB, outer cone degrees, energy
- **`LIGHT.ENABLE`** - args: handle, bool
- **`LIGHT.FREE`** - args: handle
- **`LIGHT.GETCOLOR`** - args: handle -> returns array
- **`LIGHT.GETDIR`** - args: handle -> returns array
- **`LIGHT.GETINTENSITY`** - args: handle -> returns float
- **`LIGHT.GETPOS`** - args: handle -> returns array
- **`LIGHT.GETROT`** - args: handle -> returns array — Returns [p, y, r] Euler rotation of the light
- **`LIGHT.ISENABLED`** - args: handle -> returns int
- **`LIGHT.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of LIGHT.CREATE. Use LIGHT.CREATE.
- **`LIGHT.MAKE`** - args: string -> returns handle — DEPRECATED alias of LIGHT.CREATE. Use LIGHT.CREATE.
- **`LIGHT.MAKEDIRECTIONAL`** - args: float, float, float, float, float, float, float -> returns handle — DEPRECATED alias of LIGHT.CREATEDIRECTIONAL. Use LIGHT.CREATEDIRECTIONAL. Directional light: direction vector (dx,dy,dz), RGB, energy — direction is normalized
- **`LIGHT.MAKEPOINT`** - args: float, float, float, float, float, float, float -> returns handle — DEPRECATED alias of LIGHT.CREATEPOINT. Use LIGHT.CREATEPOINT. Point light at (x,y,z) with RGB (0–255 or 0–1) and intensity (energy)
- **`LIGHT.MAKESPOT`** - args: float, float, float, float, float, float, float, float, float, float, float -> returns handle — DEPRECATED alias of LIGHT.CREATESPOT. Use LIGHT.CREATESPOT. Spot: position, target point, RGB, outer cone degrees, energy
- **`LIGHT.SETCOLOR`** - args: handle, float, float, float
- **`LIGHT.SETCOLOR`** - args: handle, float, float, float, float
- **`LIGHT.SETDIR`** - args: handle, float, float, float
- **`LIGHT.SETINNERCONE`** - args: handle, float
- **`LIGHT.SETINTENSITY`** - args: handle, float
- **`LIGHT.SETOUTERCONE`** - args: handle, float
- **`LIGHT.SETPOS`** - args: handle, float, float, float
- **`LIGHT.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of LIGHT.SETPOS. Use LIGHT.SETPOS.
- **`LIGHT.SETRANGE`** - args: handle, float
- **`LIGHT.SETROT`** - args: handle, float, float, float — Sets light orientation using Euler angles (pitch, yaw, roll)
- **`LIGHT.SETSHADOW`** - args: handle, bool
- **`LIGHT.SETSHADOWBIAS`** - args: handle, float
- **`LIGHT.SETSTATE`** - args: handle, bool — Alias of LIGHT.ENABLE
- **`LIGHT.SETTARGET`** - args: handle, float, float, float

### LIGHT2D

- **`LIGHT2D.CREATE`** - args: (none) -> returns handle
- **`LIGHT2D.FREE`** - args: handle
- **`LIGHT2D.GETCOLOR`** - args: handle -> returns array
- **`LIGHT2D.GETPOS`** - args: handle -> returns array
- **`LIGHT2D.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of LIGHT2D.CREATE. Use LIGHT2D.CREATE.
- **`LIGHT2D.SETCOLOR`** - args: handle, int, int, int, int
- **`LIGHT2D.SETINTENSITY`** - args: handle, float
- **`LIGHT2D.SETPOS`** - args: handle, float, float
- **`LIGHT2D.SETPOSITION`** - args: handle, float, float — DEPRECATED alias of LIGHT2D.SETPOS. Use LIGHT2D.SETPOS.
- **`LIGHT2D.SETRADIUS`** - args: handle, float

### LINE3D

- **`LINE3D`** - args: float, float, float, float, float, float, int, int, int, int — Shorthand: DRAW3D.LINE(x1, y1, z1, x2, y2, z2, r, g, b, a)

### LISTEN

- **`LISTEN`** - args: int -> returns handle — Easy Mode: NET.HOST(port)

### LOADFONT

- **`LOADFONT`** - args: string, int -> returns handle — Easy Mode: FONT.LOAD(path, size)

### LOADIMAGE

- **`LOADIMAGE`** - args: string -> returns handle — Easy Mode: IMAGE.LOAD(path)

### LOADMESH

- **`LOADMESH`** - args: string -> returns handle — Easy Mode: MESH.LOAD(path)

### LOADMUSIC

- **`LOADMUSIC`** - args: string -> returns handle — Easy Mode: AUDIO.LOADMUSIC(path)

### LOADSOUND

- **`LOADSOUND`** - args: string -> returns handle — Easy Mode: Load a sound file

### LOADSPRITE

- **`LOADSPRITE`** - args: string -> returns int — Easy Mode: Load a 3D billboard sprite (entity#); optional parent entity#
- **`LOADSPRITE`** - args: string, int -> returns int — Load billboard sprite as child of parent entity#

### LOADTEXTURE

- **`LOADTEXTURE`** - args: string -> returns handle

### LOBBY

- **`LOBBY.CREATE`** - args: string, int -> returns handle
- **`LOBBY.FIND`** - args: string, string -> returns handle
- **`LOBBY.FREE`** - args: handle
- **`LOBBY.GETNAME`** - args: handle -> returns string
- **`LOBBY.JOIN`** - args: handle
- **`LOBBY.MAKE`** - args: string, int -> returns handle — DEPRECATED alias of LOBBY.CREATE. Use LOBBY.CREATE.
- **`LOBBY.SETHOST`** - args: handle, string, int
- **`LOBBY.SETPROPERTY`** - args: handle, string, string
- **`LOBBY.START`** - args: handle

### LOCATE

- **`LOCATE`** - args: int, int

### LOG

- **`LOG`** - args: any

### LOG10

- **`LOG10`** - args: any

### LOG2

- **`LOG2`** - args: any

### LOWER

- **`LOWER`** - args: string -> returns string

### LOWER$

- **`LOWER$`** - args: string -> returns string

### LSET

- **`LSET`** - args: string, int -> returns string

### LSET$

- **`LSET$`** - args: string, int -> returns string

### LTRIM

- **`LTRIM`** - args: string -> returns string

### LTRIM$

- **`LTRIM$`** - args: string -> returns string

### LightColor

- **`LightColor`** - args: handle, float, float, float — Alias of LIGHT.SETCOLOR (RGB)

### LightRange

- **`LightRange`** - args: handle, float — Alias of LIGHT.SETRANGE

### Listener

- **`Listener`** - args: handle — Sets spatial audio listener from a Camera3D handle (CAMERA.CREATE / CreateCamera; deprecated alias CAMERA.MAKE); call each frame before EmitSound

### Load3DSound

- **`Load3DSound`** - args: string -> returns handle — Loads WAV/OGG like AUDIO.LOADSOUND; use with Listener + EmitSound for 3D pan/attenuation

### LoadAnimMesh

- **`LoadAnimMesh`** - args: any -> returns int

### MAKEDIR

- **`MAKEDIR`** - args: string
- **`MAKEDIR`** - args: string -> returns bool

### MAKEDIRS

- **`MAKEDIRS`** - args: string
- **`MAKEDIRS`** - args: string -> returns bool

### MAT4

- **`MAT4.FREE`** - args: handle
- **`MAT4.FROMROTATION`** - args: float, float, float -> returns handle
- **`MAT4.FROMSCALE`** - args: float, float, float -> returns handle
- **`MAT4.FROMTRANSLATION`** - args: float, float, float -> returns handle
- **`MAT4.GETELEMENT`** - args: handle, int, int -> returns float
- **`MAT4.IDENTITY`** - args: (none) -> returns handle
- **`MAT4.INVERSE`** - args: handle -> returns handle
- **`MAT4.LOOKAT`** - args: float, float, float, float, float, float, float, float, float -> returns handle
- **`MAT4.MULTIPLY`** - args: handle, handle -> returns handle
- **`MAT4.ORTHO`** - args: float, float, float, float, float, float -> returns handle
- **`MAT4.PERSPECTIVE`** - args: float, float, float, float -> returns handle
- **`MAT4.ROTATION`** - args: float, float, float -> returns handle
- **`MAT4.SETROTATION`** - args: handle, float, float, float
- **`MAT4.TRANSFORMX`** - args: handle, float, float, float -> returns float
- **`MAT4.TRANSFORMY`** - args: handle, float, float, float -> returns float
- **`MAT4.TRANSFORMZ`** - args: handle, float, float, float -> returns float
- **`MAT4.TRANSPOSE`** - args: handle -> returns handle

### MATERIAL

- **`MATERIAL.AUTOFILTER`** - args: any
- **`MATERIAL.BULKASSIGN`** - args: string, handle -> returns int
- **`MATERIAL.CREATE`** - args: (none) -> returns handle
- **`MATERIAL.CREATEDEFAULT`** - args: (none)
- **`MATERIAL.CREATEPBR`** - args: (none) -> returns handle
- **`MATERIAL.FREE`** - args: handle
- **`MATERIAL.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of MATERIAL.CREATE. Use MATERIAL.CREATE.
- **`MATERIAL.MAKEDEFAULT`** - args: (none) — DEPRECATED alias of MATERIAL.CREATEDEFAULT. Use MATERIAL.CREATEDEFAULT.
- **`MATERIAL.MAKEPBR`** - args: (none) -> returns handle — DEPRECATED alias of MATERIAL.CREATEPBR. Use MATERIAL.CREATEPBR.
- **`MATERIAL.SETCOLOR`** - args: handle, int, int, int, int, int
- **`MATERIAL.SETEFFECT`** - args: handle, string
- **`MATERIAL.SETEFFECTPARAM`** - args: handle, string, float
- **`MATERIAL.SETFLOAT`** - args: handle, int, float
- **`MATERIAL.SETSECONDARYTEXTURE`** - args: int, handle — Alias of ENTITY.SETDETAILTEXTURE
- **`MATERIAL.SETSHADER`** - args: handle, handle
- **`MATERIAL.SETTEXTURE`** - args: handle, int, handle
- **`MATERIAL.SETUVSCROLL`** - args: int, float, float — Alias of ENTITY.SCROLLMATERIAL (mesh material 0)

### MATH

- **`MATH.ABS`** - args: any
- **`MATH.ACOS`** - args: any
- **`MATH.ANGLEDIFF`** - args: any, any
- **`MATH.ANGLEDIFFRAD`** - args: float, float -> returns float — Same as ANGLEDIFFRAD
- **`MATH.ANGLETO`** - args: float, float, float, float -> returns float — Same as ANGLETO
- **`MATH.APPROACH`** - args: float, float, float -> returns float
- **`MATH.ASIN`** - args: any
- **`MATH.ATAN`** - args: any
- **`MATH.ATAN2`** - args: any, any
- **`MATH.ATN`** - args: any
- **`MATH.CEIL`** - args: any
- **`MATH.CIRCLEPOINT`** - args: float, float, float, float, float -> returns handle
- **`MATH.CLAMP`** - args: any, any, any
- **`MATH.COS`** - args: any
- **`MATH.COSD`** - args: any
- **`MATH.CURVE`** - args: float, float, float -> returns float — Alias of CURVE — current + (target-current)/divisor (divisor clamped to >= 1)
- **`MATH.DEG2RAD`** - args: any
- **`MATH.DEGPERSEC`** - args: any, any
- **`MATH.DIST2D`** - args: float, float, float, float -> returns float — Same as DIST2D
- **`MATH.DISTSQ2D`** - args: float, float, float, float -> returns float — Same as DISTSQ2D
- **`MATH.E`** - args: (none)
- **`MATH.EXP`** - args: any
- **`MATH.FIX`** - args: any
- **`MATH.FLOOR`** - args: any
- **`MATH.HDIST`** - args: float, float, float, float -> returns float — Same as HDIST
- **`MATH.HDISTSQ`** - args: float, float, float, float -> returns float — Same as HDISTSQ
- **`MATH.INVERSE_LERP`** - args: float, float, float -> returns float
- **`MATH.LERP`** - args: any, any, any
- **`MATH.LERPANGLE`** - args: float, float, float -> returns float
- **`MATH.LOG`** - args: any
- **`MATH.LOG10`** - args: any
- **`MATH.LOG2`** - args: any
- **`MATH.MAX`** - args: any, any
- **`MATH.MIN`** - args: any, any
- **`MATH.NEWX`** - args: float, float, float -> returns float — currentX + MOVEX(yaw,1,0)*dist — yaw in radians (XZ forward step)
- **`MATH.NEWZ`** - args: float, float, float -> returns float — currentZ + MOVEZ(yaw,1,0)*dist — yaw in radians
- **`MATH.PI`** - args: (none)
- **`MATH.PINGPONG`** - args: any, any
- **`MATH.POW`** - args: any, any
- **`MATH.RAD2DEG`** - args: any
- **`MATH.RAND`** - args: any, any -> returns int — Same as RAND
- **`MATH.REMAP`** - args: float, float, float, float, float -> returns float
- **`MATH.RND`** - args: (none)
- **`MATH.RND`** - args: any
- **`MATH.RND`** - args: any, any -> returns int — Inclusive int range — same as RND(lo, hi)
- **`MATH.RNDF`** - args: any, any
- **`MATH.RNDSEED`** - args: any
- **`MATH.ROUND`** - args: any
- **`MATH.ROUND`** - args: any, any
- **`MATH.SATURATE`** - args: float -> returns float
- **`MATH.SGN`** - args: any
- **`MATH.SIGN`** - args: any
- **`MATH.SIN`** - args: any
- **`MATH.SIND`** - args: any
- **`MATH.SMOOTHERSTEP`** - args: any, any, any -> returns float — Same as SMOOTHERSTEP
- **`MATH.SMOOTHSTEP`** - args: any, any, any
- **`MATH.SQR`** - args: any
- **`MATH.SQRT`** - args: any
- **`MATH.TAN`** - args: any
- **`MATH.TAND`** - args: any
- **`MATH.TAU`** - args: (none)
- **`MATH.WRAP`** - args: any, any, any
- **`MATH.WRAPANGLE`** - args: any
- **`MATH.WRAPANGLE180`** - args: any
- **`MATH.YAWFROMXZ`** - args: float, float -> returns float — Same as YAWFROMXZ

### MATRIX

- **`MATRIX.FREE`** - args: handle

### MAX

- **`MAX`** - args: any, any

### MEM

- **`MEM.CLEAR`** - args: handle
- **`MEM.COPY`** - args: handle, handle, int, int, int
- **`MEM.CREATE`** - args: int -> returns handle
- **`MEM.FREE`** - args: handle
- **`MEM.GETBYTE`** - args: handle, int -> returns int
- **`MEM.GETDWORD`** - args: handle, int -> returns int
- **`MEM.GETFLOAT`** - args: handle, int -> returns float
- **`MEM.GETSTRING`** - args: handle, int -> returns string
- **`MEM.GETWORD`** - args: handle, int -> returns int
- **`MEM.MAKE`** - args: int -> returns handle — DEPRECATED alias of MEM.CREATE. Use MEM.CREATE.
- **`MEM.SETBYTE`** - args: handle, int, int
- **`MEM.SETDWORD`** - args: handle, int, int
- **`MEM.SETFLOAT`** - args: handle, int, float
- **`MEM.SETSTRING`** - args: handle, int, string
- **`MEM.SETWORD`** - args: handle, int, int
- **`MEM.SIZE`** - args: handle -> returns int

### MESH

- **`MESH.CREATECAPSULE`** - args: float, float, int, int
- **`MESH.CREATECONE`** - args: float, float, int
- **`MESH.CREATECUBE`** - args: float, float, float -> returns handle — Alias of MESH.MAKECUBE
- **`MESH.CREATECUBE`** - args: float, float, float
- **`MESH.CREATECUBICMAP`** - args: handle, float, float, float
- **`MESH.CREATECUSTOM`** - args: handle, handle -> returns handle
- **`MESH.CREATECYLINDER`** - args: float, float, int
- **`MESH.CREATEHEIGHTMAP`** - args: handle, float, float, float
- **`MESH.CREATEKNOT`** - args: float, float, int, int
- **`MESH.CREATEPLANE`** - args: float, float, int, int -> returns handle — Alias of MESH.MAKEPLANE — procedural plane mesh handle
- **`MESH.CREATEPLANE`** - args: float, float, int, int
- **`MESH.CREATEPOLY`** - args: int, float
- **`MESH.CREATESPHERE`** - args: float, int, int
- **`MESH.CREATESPHERE`** - args: float, int, int -> returns handle — Alias of MESH.MAKESPHERE
- **`MESH.CREATETORUS`** - args: float, float, int, int
- **`MESH.CUBE`** - args: float, float, float
- **`MESH.DRAW`** - args: handle, handle, handle
- **`MESH.DRAWAT`** - args: handle, handle, float, float, float
- **`MESH.DRAWINSTANCED`** - args: handle, handle, handle, int
- **`MESH.DRAWROTATED`** - args: handle, handle, float, float, float
- **`MESH.EXPORT`** - args: handle, string
- **`MESH.FREE`** - args: handle
- **`MESH.GENERATEBOUNDS`** - args: handle
- **`MESH.GENERATELOD`** - args: handle, float, float
- **`MESH.GENERATELODCHAIN`** - args: handle, any
- **`MESH.GENERATENORMALS`** - args: handle
- **`MESH.GENTANGENTS`** - args: handle
- **`MESH.GETBBOXMAXX`** - args: handle
- **`MESH.GETBBOXMAXY`** - args: handle
- **`MESH.GETBBOXMAXZ`** - args: handle
- **`MESH.GETBBOXMINX`** - args: handle
- **`MESH.GETBBOXMINY`** - args: handle
- **`MESH.GETBBOXMINZ`** - args: handle
- **`MESH.GETBOUNDS`** - args: handle -> returns handle
- **`MESH.LOAD`** - args: string -> returns handle
- **`MESH.MAKECAPSULE`** - args: float, float, int, int — DEPRECATED alias of MESH.CREATECAPSULE. Use MESH.CREATECAPSULE.
- **`MESH.MAKECONE`** - args: float, float, int — DEPRECATED alias of MESH.CREATECONE. Use MESH.CREATECONE.
- **`MESH.MAKECUBE`** - args: float, float, float -> returns handle — DEPRECATED alias of MESH.CREATECUBE. Use MESH.CREATECUBE. Alias of MESH.MAKECUBE
- **`MESH.MAKECUBE`** - args: float, float, float — DEPRECATED alias of MESH.CREATECUBE. Use MESH.CREATECUBE.
- **`MESH.MAKECUBICMAP`** - args: handle, float, float, float — DEPRECATED alias of MESH.CREATECUBICMAP. Use MESH.CREATECUBICMAP.
- **`MESH.MAKECUSTOM`** - args: handle, handle -> returns handle — DEPRECATED alias of MESH.CREATECUSTOM. Use MESH.CREATECUSTOM.
- **`MESH.MAKECYLINDER`** - args: float, float, int — DEPRECATED alias of MESH.CREATECYLINDER. Use MESH.CREATECYLINDER.
- **`MESH.MAKEHEIGHTMAP`** - args: handle, float, float, float — DEPRECATED alias of MESH.CREATEHEIGHTMAP. Use MESH.CREATEHEIGHTMAP.
- **`MESH.MAKEKNOT`** - args: float, float, int, int — DEPRECATED alias of MESH.CREATEKNOT. Use MESH.CREATEKNOT.
- **`MESH.MAKEPLANE`** - args: float, float, int, int -> returns handle — DEPRECATED alias of MESH.CREATEPLANE. Use MESH.CREATEPLANE. Alias of MESH.MAKEPLANE — procedural plane mesh handle
- **`MESH.MAKEPLANE`** - args: float, float, int, int — DEPRECATED alias of MESH.CREATEPLANE. Use MESH.CREATEPLANE.
- **`MESH.MAKEPOLY`** - args: int, float — DEPRECATED alias of MESH.CREATEPOLY. Use MESH.CREATEPOLY.
- **`MESH.MAKESPHERE`** - args: float, int, int — DEPRECATED alias of MESH.CREATESPHERE. Use MESH.CREATESPHERE.
- **`MESH.MAKESPHERE`** - args: float, int, int -> returns handle — DEPRECATED alias of MESH.CREATESPHERE. Use MESH.CREATESPHERE. Alias of MESH.MAKESPHERE
- **`MESH.MAKETORUS`** - args: float, float, int, int — DEPRECATED alias of MESH.CREATETORUS. Use MESH.CREATETORUS.
- **`MESH.OPTIMISEALL`** - args: handle
- **`MESH.OPTIMISEFETCH`** - args: handle
- **`MESH.OPTIMISEOVERDRAW`** - args: handle, float
- **`MESH.OPTIMISEVERTEXCACHE`** - args: handle
- **`MESH.OPTIMIZEALL`** - args: handle
- **`MESH.OPTIMIZEFETCH`** - args: handle
- **`MESH.OPTIMIZEOVERDRAW`** - args: handle, float
- **`MESH.OPTIMIZEVERTEXCACHE`** - args: handle
- **`MESH.PLANE`** - args: float, float, int, int
- **`MESH.SPHERE`** - args: float, int, int
- **`MESH.TRIANGLECOUNT`** - args: handle -> returns int
- **`MESH.UPDATEVERTEX`** - args: handle, int, float, float, float, float, float, float, float, float
- **`MESH.UPDATEVERTICES`** - args: handle, handle
- **`MESH.UPLOAD`** - args: handle, bool
- **`MESH.VERTEXCOUNT`** - args: handle -> returns int

### MID

- **`MID`** - args: string, int -> returns string
- **`MID`** - args: string, int, int -> returns string

### MID$

- **`MID$`** - args: string, int -> returns string
- **`MID$`** - args: string, int, int -> returns string

### MILLISECOND

- **`MILLISECOND`** - args: (none)
- **`MILLISECOND`** - args: (none) -> returns int

### MILLISECS

- **`MILLISECS`** - args: (none) -> returns int — Blitz-style: TIME.MILLIS()

### MIN

- **`MIN`** - args: any, any

### MINUTE

- **`MINUTE`** - args: (none)
- **`MINUTE`** - args: (none) -> returns int

### MKDOUBLE

- **`MKDOUBLE`** - args: any -> returns handle

### MKDOUBLE$

- **`MKDOUBLE$`** - args: any -> returns handle

### MKFLOAT

- **`MKFLOAT`** - args: any -> returns handle

### MKFLOAT$

- **`MKFLOAT$`** - args: any -> returns handle

### MKINT

- **`MKINT`** - args: any -> returns handle

### MKINT$

- **`MKINT$`** - args: any -> returns handle

### MKLONG

- **`MKLONG`** - args: any -> returns handle

### MKLONG$

- **`MKLONG$`** - args: any -> returns handle

### MKSHORT

- **`MKSHORT`** - args: any -> returns handle

### MKSHORT$

- **`MKSHORT$`** - args: any -> returns handle

### MODEL

- **`MODEL.ADDCHILD`** - args: handle, handle
- **`MODEL.ANIMCOUNT`** - args: handle -> returns int
- **`MODEL.ANIMDONE`** - args: handle -> returns bool
- **`MODEL.ANIMNAME`** - args: handle, int -> returns string
- **`MODEL.ANIMNAME$`** - args: handle, int -> returns string
- **`MODEL.ATTACHTO`** - args: handle, handle
- **`MODEL.CHILDCOUNT`** - args: handle -> returns int
- **`MODEL.CLONE`** - args: handle
- **`MODEL.CREATE`** - args: handle -> returns handle
- **`MODEL.CREATEBOX`** - args: float, float, float -> returns handle
- **`MODEL.CREATEBOX`** - args: float, float, float, bool -> returns handle
- **`MODEL.CREATECAPSULE`** - args: float, float -> returns handle — EntityRef capsule primitive (radius#, height#); draw matches Jolt capsule when using ENTITY.ADDPHYSICS capsule
- **`MODEL.CREATEINSTANCED`** - args: string, int -> returns handle
- **`MODEL.DETACH`** - args: handle
- **`MODEL.DRAW`** - args: handle
- **`MODEL.DRAWAT`** - args: handle, float, float, float, float, float, float, float, float, float
- **`MODEL.DRAWEX`** - args: handle, float, float, float, float, float, float, float, float, float, float, int, int, int, int
- **`MODEL.DRAWWIRES`** - args: handle, int, int, int, int
- **`MODEL.EXISTS`** - args: handle
- **`MODEL.FREE`** - args: handle
- **`MODEL.GETCHILD`** - args: handle, int -> returns handle
- **`MODEL.GETFRAME`** - args: handle -> returns int
- **`MODEL.GETMATERIALCOUNT`** - args: handle
- **`MODEL.GETPARENT`** - args: handle -> returns handle
- **`MODEL.GETPOS`** - args: handle -> returns handle
- **`MODEL.GETROT`** - args: handle -> returns handle
- **`MODEL.GETSCALE`** - args: handle -> returns handle
- **`MODEL.HIDE`** - args: handle
- **`MODEL.INSTANCE`** - args: handle
- **`MODEL.ISPLAYING`** - args: handle -> returns bool
- **`MODEL.ISVISIBLE`** - args: handle -> returns bool
- **`MODEL.LIMBCOUNT`** - args: handle -> returns int
- **`MODEL.LIMBX`** - args: handle, int -> returns float
- **`MODEL.LOAD`** - args: string
- **`MODEL.LOADANIMATIONS`** - args: handle, string
- **`MODEL.LOADLOD`** - args: string, string, string -> returns handle
- **`MODEL.LOOP`** - args: handle, bool
- **`MODEL.MAKE`** - args: handle -> returns handle — DEPRECATED alias of MODEL.CREATE. Use MODEL.CREATE.
- **`MODEL.MAKEBOX`** - args: float, float, float -> returns handle — DEPRECATED alias of MODEL.CREATEBOX. Use MODEL.CREATEBOX.
- **`MODEL.MAKEBOX`** - args: float, float, float, bool -> returns handle — DEPRECATED alias of MODEL.CREATEBOX. Use MODEL.CREATEBOX.
- **`MODEL.MAKECAPSULE`** - args: float, float -> returns handle — DEPRECATED alias of MODEL.CREATECAPSULE. Use MODEL.CREATECAPSULE. EntityRef capsule primitive (radius#, height#); draw matches Jolt capsule when using ENTITY.ADDPHYSICS capsule
- **`MODEL.MAKEINSTANCED`** - args: string, int -> returns handle — DEPRECATED alias of MODEL.CREATEINSTANCED. Use MODEL.CREATEINSTANCED.
- **`MODEL.MOVE`** - args: handle, float, float, float
- **`MODEL.PLAY`** - args: handle, string
- **`MODEL.PLAYIDX`** - args: handle, int
- **`MODEL.REMOVECHILD`** - args: handle, handle
- **`MODEL.ROTATE`** - args: handle, float, float, float
- **`MODEL.ROTATETEXTURE`** - args: handle, float
- **`MODEL.SCALETEXTURE`** - args: handle, float, float
- **`MODEL.SCROLLTEXTURE`** - args: handle, float, float
- **`MODEL.SETALPHA`** - args: handle, int
- **`MODEL.SETAMBIENTCOLOR`** - args: handle, int, int, int
- **`MODEL.SETBLEND`** - args: handle, int
- **`MODEL.SETCASTSHADOW`** - args: handle, bool
- **`MODEL.SETCOLOR`** - args: handle, int, int, int, int
- **`MODEL.SETCULL`** - args: handle, bool
- **`MODEL.SETDEPTH`** - args: handle, int
- **`MODEL.SETDIFFUSE`** - args: handle, int, int, int
- **`MODEL.SETEMISSIVE`** - args: handle, int, int, int
- **`MODEL.SETFOG`** - args: handle, bool
- **`MODEL.SETGPUSKINNING`** - args: handle, bool
- **`MODEL.SETINSTANCEPOS`** - args: handle, int, float, float, float
- **`MODEL.SETINSTANCESCALE`** - args: handle, int, float, float, float
- **`MODEL.SETLIGHTING`** - args: handle, bool
- **`MODEL.SETLIMBPOS`** - args: handle, int, float, float, float
- **`MODEL.SETLODDISTANCES`** - args: handle, float, float, float
- **`MODEL.SETMATERIAL`** - args: handle, int, handle
- **`MODEL.SETMATERIALSHADER`** - args: handle, int, handle
- **`MODEL.SETMATERIALTEXTURE`** - args: handle, int, int, handle
- **`MODEL.SETMATRIX`** - args: handle, handle
- **`MODEL.SETMETAL`** - args: handle, float
- **`MODEL.SETMODELMESHMATERIAL`** - args: handle, int, int
- **`MODEL.SETPOS`** - args: handle, float, float, float
- **`MODEL.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of MODEL.SETPOS. Use MODEL.SETPOS.
- **`MODEL.SETRECEIVESHADOW`** - args: handle, bool
- **`MODEL.SETROT`** - args: handle, float, float, float
- **`MODEL.SETROUGH`** - args: handle, float
- **`MODEL.SETSCALE`** - args: handle, float, float, float
- **`MODEL.SETSCALEUNIFORM`** - args: handle, float
- **`MODEL.SETSPECULAR`** - args: handle, int, int, int
- **`MODEL.SETSPECULARPOW`** - args: handle, float
- **`MODEL.SETSPEED`** - args: handle, float
- **`MODEL.SETSTAGEBLEND`** - args: handle, int, float
- **`MODEL.SETSTAGEROTATE`** - args: handle, int, float
- **`MODEL.SETSTAGESCALE`** - args: handle, int, float, float
- **`MODEL.SETSTAGESCROLL`** - args: handle, int, float, float
- **`MODEL.SETTEXTURESTAGE`** - args: handle, int, handle
- **`MODEL.SETWIREFRAME`** - args: handle, bool
- **`MODEL.SHOW`** - args: handle
- **`MODEL.STOP`** - args: handle
- **`MODEL.TOTALFRAMES`** - args: handle -> returns int
- **`MODEL.UPDATEANIM`** - args: handle, float
- **`MODEL.UPDATEINSTANCES`** - args: handle
- **`MODEL.X`** - args: handle -> returns float
- **`MODEL.Y`** - args: handle -> returns float
- **`MODEL.Z`** - args: handle -> returns float

### MONTH

- **`MONTH`** - args: (none)
- **`MONTH`** - args: (none) -> returns int

### MOUSE

- **`MOUSE`** - args: (none) -> returns handle — Singleton mouse input facade handle
- **`MOUSE.DX`** - args: handle -> returns float

### MOUSEDX

- **`MOUSEDX`** - args: (none) -> returns float — Easy Mode: MOUSE.DX(MOUSE())

### MOUSEDY

- **`MOUSEDY`** - args: (none) -> returns float — Easy Mode: MOUSE.DY(MOUSE())

### MOUSEHIT

- **`MOUSEHIT`** - args: int -> returns int — Easy Mode: Returns 1 if mouse button was pressed this frame

### MOUSEWHEEL

- **`MOUSEWHEEL`** - args: (none) -> returns float — Easy Mode: MOUSE.WHEEL(MOUSE())

### MOUSEX

- **`MOUSEX`** - args: (none) -> returns int — Easy Mode: Get absolute mouse X coordinate

### MOUSEY

- **`MOUSEY`** - args: (none) -> returns int — Easy Mode: Get absolute mouse Y coordinate

### MOUSEZ

- **`MOUSEZ`** - args: (none) -> returns int — Easy Mode: Get mouse wheel movement

### MOVE

- **`MOVE.LERP`** - args: float, float, float -> returns float — Alias of MATH.LERP
- **`MOVE.TOWARD`** - args: float, float, float -> returns float — Alias of MATH.APPROACH — move current toward target by at most maxDelta without overshooting

### MOVEENTITY

- **`MOVEENTITY`** - args: int, float, float, float — Blitz-style: ENTITY.MOVEENTITY(obj, x, y, z)
- **`MOVEENTITY`** - args: handle, float, float, float -> returns void — Easy Mode: Move entity relative to orientation

### MOVEENTITY2D

- **`MOVEENTITY2D`** - args: handle, float, float, float, float, float

### MOVEFILE

- **`MOVEFILE`** - args: string, string
- **`MOVEFILE`** - args: string, string -> returns bool

### MOVEPLAYER

- **`MOVEPLAYER`** - args: handle, float, float, float, float, float

### MOVER

- **`MOVER`** - args: (none) -> returns handle
- **`MOVER.MOVESTEPX`** - args: handle, float, float, float, float, float -> returns float
- **`MOVER.MOVESTEPZ`** - args: handle, float, float, float, float, float -> returns float
- **`MOVER.MOVEXZ`** - args: handle, float, float, float, float, float -> returns handle

### MOVESTEPX

- **`MOVESTEPX`** - args: float, float, float, float, float -> returns float — Same as MOVEX(yaw,f,s)*speed*dt — world X delta this frame

### MOVESTEPZ

- **`MOVESTEPZ`** - args: float, float, float, float, float -> returns float — Same as MOVEZ(yaw,f,s)*speed*dt — world Z delta this frame

### MOVEX

- **`MOVEX`** - args: float, float, float -> returns float — Camera-relative world X on XZ plane: yaw#, forward#, strafe#

### MOVEZ

- **`MOVEZ`** - args: float, float, float -> returns float — Camera-relative world Z on XZ plane: yaw#, forward#, strafe#

### MUSIC

- **`MUSIC.FREE`** - args: handle

### MUSICVOLUME

- **`MUSICVOLUME`** - args: handle, float — Easy Mode: AUDIO.SETMUSICVOLUME(music, vol)

### MilliSecs

- **`MilliSecs`** - args: (none) -> returns float — Milliseconds since Raylib init (CGO); monotonic wall ms on stub builds

### MouseWheel

- **`MouseWheel`** - args: (none) -> returns float — Alias of INPUT.MOUSEWHEELMOVE (Input.MouseWheel() style)

### MoveEntity

- **`MoveEntity`** - args: int, float, float, float — Args: (entity#, forward#, right#, up#). Move along entity local axes from pitch/yaw (same as MOVEENTITY / ENTITY.MOVE). For world-space offset use TranslateEntity(entity#, dx#, dy#, dz#).

### NAV

- **`NAV.ADDOBSTACLE`** - args: handle, handle
- **`NAV.ADDTERRAIN`** - args: handle, handle
- **`NAV.BAKE`** - args: handle, float, float -> returns handle — Builds a coarse walkability grid from a terrain heightmap (slope limit); returns nav handle and caches per terrain for NAV.GETPATH
- **`NAV.BUILD`** - args: int — Automatically scan the world for static geometry and bake the navigation grid.
- **`NAV.BUILD`** - args: handle
- **`NAV.CHASE`** - args: int, int, float, float — KCC follow: move toward target entity until within standoff gap (world units)
- **`NAV.CREATE`** - args: (none) -> returns handle
- **`NAV.CREATE`** - args: (none) -> returns int — Create a new navigation grid handle.
- **`NAV.DEBUGDRAW`** - args: int — Render a debug overlay of the navigation grid (Green=Walkable, Red=Blocked).
- **`NAV.DEBUGDRAW`** - args: handle
- **`NAV.FINDPATH`** - args: handle, float, float, float, float, float, float -> returns handle
- **`NAV.FREE`** - args: handle
- **`NAV.GETPATH`** - args: handle, float, float, float, float -> returns handle — A* path on last NAV.BAKE for this terrain (start/end XZ; Y sampled from terrain)
- **`NAV.GOTO`** - args: int, float, float, float — Alias of PLAYER.NAVTO — click-to-move for KCC (default arrival ~0.2 world units)
- **`NAV.GOTO`** - args: int, float, float, float, float — NAV.GOTO with arrival distance (alias of PLAYER.NAVTO)
- **`NAV.GOTO`** - args: int, float, float, float, float, float — NAV.GOTO with arrival and brake distance
- **`NAV.ISREACHABLE`** - args: handle, float, float, float, float -> returns bool — True if NAV.GETPATH would return a valid path
- **`NAV.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of NAV.CREATE. Use NAV.CREATE.
- **`NAV.MAKE`** - args: (none) -> returns int — DEPRECATED alias of NAV.CREATE. Use NAV.CREATE. Create a new navigation grid handle.
- **`NAV.PATROL`** - args: int, float, float, float, float, float — KCC ping-pong patrol between world XZ points A and B
- **`NAV.SETGRID`** - args: int, int, int, float, float, float — Initialize navigation grid dimensions: (handle, width, height, cellSize#, offsetX#, offsetY#)
- **`NAV.SETGRID`** - args: handle, int, int, float, float, float
- **`NAV.UPDATE`** - args: int — Alias of PLAYER.NAVUPDATE

### NAVAGENT

- **`NAVAGENT.APPLYFORCE`** - args: handle, float, float, float
- **`NAVAGENT.CREATE`** - args: handle -> returns handle
- **`NAVAGENT.CREATE`** - args: int -> returns int — Create a navigation agent for the specified grid handle.
- **`NAVAGENT.FREE`** - args: handle
- **`NAVAGENT.GETPOS`** - args: handle -> returns array
- **`NAVAGENT.GETROT`** - args: handle -> returns array — Approximate [pitch,yaw,roll] (radians) from path waypoint tangent or steering velocity
- **`NAVAGENT.ISATDESTINATION`** - args: int -> returns bool — Check if the agent has reached its destination.
- **`NAVAGENT.ISATDESTINATION`** - args: handle -> returns bool
- **`NAVAGENT.MAKE`** - args: handle -> returns handle — DEPRECATED alias of NAVAGENT.CREATE. Use NAVAGENT.CREATE.
- **`NAVAGENT.MAKE`** - args: int -> returns int — DEPRECATED alias of NAVAGENT.CREATE. Use NAVAGENT.CREATE. Create a navigation agent for the specified grid handle.
- **`NAVAGENT.MOVETO`** - args: handle, float, float, float
- **`NAVAGENT.MOVETO`** - args: int, float, float, float — Set the agent's target destination: (handle, x#, y#, z#)
- **`NAVAGENT.SETMAXFORCE`** - args: handle, float
- **`NAVAGENT.SETPOS`** - args: handle, float, float, float
- **`NAVAGENT.SETPOS`** - args: int, float, float, float — Set the agent's world-space position: (handle, x#, y#, z#)
- **`NAVAGENT.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of NAVAGENT.SETPOS. Use NAVAGENT.SETPOS.
- **`NAVAGENT.SETPOSITION`** - args: int, float, float, float — DEPRECATED alias of NAVAGENT.SETPOS. Use NAVAGENT.SETPOS. Set the agent's world-space position: (handle, x#, y#, z#)
- **`NAVAGENT.SETSPEED`** - args: int, float — Set the agent's movement speed: (handle, speed#)
- **`NAVAGENT.SETSPEED`** - args: handle, float
- **`NAVAGENT.UPDATE`** - args: int, float — Update the agent's movement: (handle, dt#)
- **`NAVAGENT.UPDATE`** - args: handle, float
- **`NAVAGENT.X`** - args: int -> returns float — Get the agent's current X position.
- **`NAVAGENT.X`** - args: handle -> returns float
- **`NAVAGENT.Y`** - args: int -> returns float — Get the agent's current Y position.
- **`NAVAGENT.Y`** - args: handle -> returns float
- **`NAVAGENT.Z`** - args: int -> returns float — Get the agent's current Z position.
- **`NAVAGENT.Z`** - args: handle -> returns float

### NET

- **`NET.BROADCAST`** - args: handle, int, string, bool
- **`NET.CLOSE`** - args: handle
- **`NET.CONNECT`** - args: string, int — Simplified command to join a server.
- **`NET.CONNECT`** - args: handle, string, int -> returns handle
- **`NET.CREATECLIENT`** - args: (none) -> returns handle
- **`NET.CREATESERVER`** - args: int, int -> returns handle
- **`NET.FLUSH`** - args: handle
- **`NET.GETPING`** - args: handle -> returns int
- **`NET.HOST`** - args: int — Simplified command to start a server.
- **`NET.MAKECLIENT`** - args: (none) -> returns handle — DEPRECATED alias of NET.CREATECLIENT. Use NET.CREATECLIENT.
- **`NET.MAKESERVER`** - args: int, int -> returns handle — DEPRECATED alias of NET.CREATESERVER. Use NET.CREATESERVER.
- **`NET.PEERCOUNT`** - args: handle -> returns int
- **`NET.RECEIVE`** - args: handle -> returns handle
- **`NET.SEND`** - args: int, string — Broadcast or send network data.
- **`NET.SERVICE`** - args: handle, int
- **`NET.SETBANDWIDTH`** - args: handle, int, int
- **`NET.SETCHANNELS`** - args: int
- **`NET.SETTIMEOUT`** - args: handle, int
- **`NET.START`** - args: (none)
- **`NET.STOP`** - args: (none)
- **`NET.SYNC`** - args: int — Mark an entity for network replication.
- **`NET.UPDATE`** - args: handle

### NETMSG$

- **`NETMSG$`** - args: (none) -> returns string — Easy Mode: NET.GETEVENTPACKET$()

### NETREADFLOAT

- **`NETREADFLOAT`** - args: (none) -> returns float

### NETREADINT

- **`NETREADINT`** - args: (none) -> returns int

### NETREADSTRING

- **`NETREADSTRING`** - args: (none) -> returns string

### NETSENDFLOAT

- **`NETSENDFLOAT`** - args: handle, float

### NETSENDINT

- **`NETSENDINT`** - args: handle, int

### NETSENDSTRING

- **`NETSENDSTRING`** - args: handle, string

### NOISE

- **`NOISE.CREATE`** - args: (none) -> returns handle
- **`NOISE.CREATECELLULAR`** - args: int, float, string -> returns handle
- **`NOISE.CREATEDOMAINWARP`** - args: int, float, float -> returns handle
- **`NOISE.CREATEFRACTAL`** - args: int, float, int, string -> returns handle
- **`NOISE.CREATEPERLIN`** - args: int, float -> returns handle
- **`NOISE.CREATESIMPLEX`** - args: int, float -> returns handle
- **`NOISE.FILLARRAY`** - args: handle, handle, int, int, float, float
- **`NOISE.FILLARRAYNORM`** - args: handle, handle, int, int, float, float
- **`NOISE.FILLIMAGE`** - args: handle, handle, float, float
- **`NOISE.FREE`** - args: handle
- **`NOISE.GET`** - args: handle, float, float -> returns float
- **`NOISE.GET3D`** - args: handle, float, float, float -> returns float
- **`NOISE.GETDOMAINWARPED`** - args: handle, float, float -> returns float
- **`NOISE.GETNORM`** - args: handle, float, float -> returns float
- **`NOISE.GETTILEABLE`** - args: handle, float, float, float, float -> returns float
- **`NOISE.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of NOISE.CREATE. Use NOISE.CREATE.
- **`NOISE.MAKECELLULAR`** - args: int, float, string -> returns handle — DEPRECATED alias of NOISE.CREATECELLULAR. Use NOISE.CREATECELLULAR.
- **`NOISE.MAKEDOMAINWARP`** - args: int, float, float -> returns handle — DEPRECATED alias of NOISE.CREATEDOMAINWARP. Use NOISE.CREATEDOMAINWARP.
- **`NOISE.MAKEFRACTAL`** - args: int, float, int, string -> returns handle — DEPRECATED alias of NOISE.CREATEFRACTAL. Use NOISE.CREATEFRACTAL.
- **`NOISE.MAKEPERLIN`** - args: int, float -> returns handle — DEPRECATED alias of NOISE.CREATEPERLIN. Use NOISE.CREATEPERLIN.
- **`NOISE.MAKESIMPLEX`** - args: int, float -> returns handle — DEPRECATED alias of NOISE.CREATESIMPLEX. Use NOISE.CREATESIMPLEX.
- **`NOISE.SETCELLULARDISTANCE`** - args: handle, string
- **`NOISE.SETCELLULARJITTER`** - args: handle, float
- **`NOISE.SETCELLULARTYPE`** - args: handle, string
- **`NOISE.SETDOMAINWARPAMPLITUDE`** - args: handle, float
- **`NOISE.SETDOMAINWARPTYPE`** - args: handle, string
- **`NOISE.SETFREQUENCY`** - args: handle, float
- **`NOISE.SETGAIN`** - args: handle, float
- **`NOISE.SETLACUNARITY`** - args: handle, float
- **`NOISE.SETOCTAVES`** - args: handle, int
- **`NOISE.SETPINGPONGSTRENGTH`** - args: handle, float
- **`NOISE.SETSEED`** - args: handle, int
- **`NOISE.SETTYPE`** - args: handle, string
- **`NOISE.SETWEIGHTEDSTRENGTH`** - args: handle, float

### OCT

- **`OCT`** - args: int -> returns string

### OCT$

- **`OCT$`** - args: int -> returns string

### OPENFILE

- **`OPENFILE`** - args: string, string

### ORBITDISTDELTA

- **`ORBITDISTDELTA`** - args: float -> returns float

### ORBITPITCHDELTA

- **`ORBITPITCHDELTA`** - args: float -> returns float

### ORBITYAWDELTA

- **`ORBITYAWDELTA`** - args: float, float, int, int, float -> returns float

### PACKET

- **`PACKET.CREATE`** - args: string -> returns handle
- **`PACKET.DATA`** - args: handle -> returns string
- **`PACKET.FREE`** - args: handle
- **`PACKET.MAKE`** - args: string -> returns handle — DEPRECATED alias of PACKET.CREATE. Use PACKET.CREATE.

### PARTICLE

- **`PARTICLE.COUNT`** - args: handle -> returns int
- **`PARTICLE.CREATE`** - args: (none) -> returns handle
- **`PARTICLE.DRAW`** - args: handle
- **`PARTICLE.DRAW`** - args: handle, handle
- **`PARTICLE.FREE`** - args: handle
- **`PARTICLE.GETALPHA`** - args: handle -> returns float
- **`PARTICLE.GETCOLOR`** - args: handle -> returns array
- **`PARTICLE.GETPOS`** - args: handle -> returns array
- **`PARTICLE.ISALIVE`** - args: handle -> returns int
- **`PARTICLE.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of PARTICLE.CREATE. Use PARTICLE.CREATE.
- **`PARTICLE.PLAY`** - args: handle
- **`PARTICLE.SETBILLBOARD`** - args: handle, bool
- **`PARTICLE.SETBURST`** - args: handle, int
- **`PARTICLE.SETCOLOR`** - args: handle, int, int, int, int
- **`PARTICLE.SETCOLOREND`** - args: handle, int, int, int, int
- **`PARTICLE.SETDIRECTION`** - args: handle, float, float, float
- **`PARTICLE.SETEMITRATE`** - args: handle, float
- **`PARTICLE.SETENDCOLOR`** - args: handle, int, int, int, int
- **`PARTICLE.SETENDSIZE`** - args: handle, float, float
- **`PARTICLE.SETGRAVITY`** - args: handle, float
- **`PARTICLE.SETGRAVITY`** - args: handle, float, float, float
- **`PARTICLE.SETLIFETIME`** - args: handle, float, float
- **`PARTICLE.SETPOS`** - args: handle, float, float, float
- **`PARTICLE.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of PARTICLE.SETPOS. Use PARTICLE.SETPOS.
- **`PARTICLE.SETRATE`** - args: handle, float
- **`PARTICLE.SETSIZE`** - args: handle, float, float
- **`PARTICLE.SETSPEED`** - args: handle, float, float
- **`PARTICLE.SETSPREAD`** - args: handle, float
- **`PARTICLE.SETSTARTCOLOR`** - args: handle, int, int, int, int
- **`PARTICLE.SETSTARTSIZE`** - args: handle, float, float
- **`PARTICLE.SETTEXTURE`** - args: handle, handle
- **`PARTICLE.SETVELOCITY`** - args: handle, float, float, float, float
- **`PARTICLE.STOP`** - args: handle
- **`PARTICLE.UPDATE`** - args: handle, float

### PARTICLE2D

- **`PARTICLE2D.CREATE`** - args: int, int, int, int, int -> returns handle
- **`PARTICLE2D.DRAW`** - args: handle
- **`PARTICLE2D.EMIT`** - args: handle, float, float, float, float, float
- **`PARTICLE2D.FREE`** - args: handle
- **`PARTICLE2D.MAKE`** - args: int, int, int, int, int -> returns handle — DEPRECATED alias of PARTICLE2D.CREATE. Use PARTICLE2D.CREATE.
- **`PARTICLE2D.UPDATE`** - args: handle, float

### PARTICLE3D

- **`PARTICLE3D.COUNT`** - args: handle -> returns int
- **`PARTICLE3D.CREATE`** - args: (none) -> returns handle
- **`PARTICLE3D.DRAW`** - args: handle
- **`PARTICLE3D.DRAW`** - args: handle, handle
- **`PARTICLE3D.FREE`** - args: handle
- **`PARTICLE3D.GETALPHA`** - args: handle -> returns float
- **`PARTICLE3D.GETCOLOR`** - args: handle -> returns array
- **`PARTICLE3D.GETPOS`** - args: handle -> returns array
- **`PARTICLE3D.ISALIVE`** - args: handle -> returns int
- **`PARTICLE3D.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of PARTICLE3D.CREATE. Use PARTICLE3D.CREATE.
- **`PARTICLE3D.PLAY`** - args: handle
- **`PARTICLE3D.SETBILLBOARD`** - args: handle, bool
- **`PARTICLE3D.SETBURST`** - args: handle, int
- **`PARTICLE3D.SETCOLOR`** - args: handle, int, int, int, int
- **`PARTICLE3D.SETCOLOREND`** - args: handle, int, int, int, int
- **`PARTICLE3D.SETDIRECTION`** - args: handle, float, float, float
- **`PARTICLE3D.SETEMITRATE`** - args: handle, float
- **`PARTICLE3D.SETENDCOLOR`** - args: handle, int, int, int, int
- **`PARTICLE3D.SETENDSIZE`** - args: handle, float, float
- **`PARTICLE3D.SETGRAVITY`** - args: handle, float
- **`PARTICLE3D.SETGRAVITY`** - args: handle, float, float, float
- **`PARTICLE3D.SETLIFETIME`** - args: handle, float, float
- **`PARTICLE3D.SETPOS`** - args: handle, float, float, float
- **`PARTICLE3D.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of PARTICLE3D.SETPOS. Use PARTICLE3D.SETPOS.
- **`PARTICLE3D.SETRATE`** - args: handle, float
- **`PARTICLE3D.SETSIZE`** - args: handle, float, float
- **`PARTICLE3D.SETSPEED`** - args: handle, float, float
- **`PARTICLE3D.SETSPREAD`** - args: handle, float
- **`PARTICLE3D.SETSTARTCOLOR`** - args: handle, int, int, int, int
- **`PARTICLE3D.SETSTARTSIZE`** - args: handle, float, float
- **`PARTICLE3D.SETTEXTURE`** - args: handle, handle
- **`PARTICLE3D.SETVELOCITY`** - args: handle, float, float, float, float
- **`PARTICLE3D.STOP`** - args: handle
- **`PARTICLE3D.UPDATE`** - args: handle, float

### PARTICLECOLOR

- **`PARTICLECOLOR`** - args: handle, int, int, int, int -> returns void — Easy Mode: Set emitter start color
- **`PARTICLECOLOR`** - args: handle, int, int, int, int

### PARTICLELIFE

- **`PARTICLELIFE`** - args: handle, float, float -> returns void — Easy Mode: Set emitter lifetime range
- **`PARTICLELIFE`** - args: handle, float, float

### PARTICLES

- **`PARTICLES.DRAWEMITTER`** - args: handle

### PARTICLESPEED

- **`PARTICLESPEED`** - args: handle, float, float -> returns void — Easy Mode: Set emitter speed range
- **`PARTICLESPEED`** - args: handle, float, float

### PATH

- **`PATH.FREE`** - args: handle
- **`PATH.ISVALID`** - args: handle -> returns bool
- **`PATH.NODECOUNT`** - args: handle -> returns int
- **`PATH.NODEX`** - args: handle, int -> returns float
- **`PATH.NODEY`** - args: handle, int -> returns float
- **`PATH.NODEZ`** - args: handle, int -> returns float

### PEER

- **`PEER.DISCONNECT`** - args: handle
- **`PEER.IP`** - args: handle -> returns string
- **`PEER.PING`** - args: handle -> returns int
- **`PEER.SEND`** - args: handle, int, string, bool
- **`PEER.SENDPACKET`** - args: handle, handle, int

### PHYSICS

- **`PHYSICS.AUTOCREATE`** - args: int
- **`PHYSICS.BOXCAST`** - args: any
- **`PHYSICS.DISABLE`** - args: any
- **`PHYSICS.ENABLE`** - args: any
- **`PHYSICS.EXPLOSION`** - args: float, float, float, float, float — Applies physical impulse radially.
- **`PHYSICS.GETBUOYANCY`** - args: int -> returns float — Reads stored buoyancy density (default 0)
- **`PHYSICS.RAYCAST`** - args: float, float, float, float, float, float, float -> returns handle
- **`PHYSICS.SETBUOYANCY`** - args: int, float — Stores per-entity buoyancy density for future Jolt/WASM fluid coupling (gameplay hint today)
- **`PHYSICS.SETGRAVITY`** - args: float, float, float
- **`PHYSICS.SETSUBSTEPS`** - args: int
- **`PHYSICS.SPHERECAST`** - args: any
- **`PHYSICS.START`** - args: (none)
- **`PHYSICS.STEP`** - args: float
- **`PHYSICS.STOP`** - args: (none)

### PHYSICS2D

- **`PHYSICS2D.DEBUGDRAW`** - args: int
- **`PHYSICS2D.GETDEBUGSEGMENTS`** - args: (none) -> returns handle
- **`PHYSICS2D.ONCOLLISION`** - args: handle, handle, string
- **`PHYSICS2D.PROCESSCOLLISIONS`** - args: (none)
- **`PHYSICS2D.SETGRAVITY`** - args: float, float
- **`PHYSICS2D.SETITERATIONS`** - args: int, int
- **`PHYSICS2D.SETSTEP`** - args: float
- **`PHYSICS2D.START`** - args: (none)
- **`PHYSICS2D.START`** - args: float, float — Optional initial gravity (gx#, gy#); same effect as SETGRAVITY after start.
- **`PHYSICS2D.STEP`** - args: (none)
- **`PHYSICS2D.STOP`** - args: (none)

### PHYSICS3D

- **`AERO.SETDRAG`** - args: handle, float — Apply air resistance coefficient.
- **`AERO.SETLIFT`** - args: handle, float — Set lift coefficient for a physics body.
- **`AERO.SETTHRUST`** - args: handle, float — Apply local Z-axis thrust power.
- **`BODY3D.APPLYTORQUE`** - args: handle, float, float, float
- **`BODY3D.GETANGULARVEL`** - args: handle -> returns handle
- **`BODY3D.GETLINEARVEL`** - args: handle -> returns handle — Get linear velocity as a 3-element array (alias of BODY3D.GETVELOCITY).
- **`BODY3D.GETMASS`** - args: handle -> returns float
- **`BODY3D.GETVELOCITY`** - args: handle -> returns handle — Get linear velocity as a 3-element numeric array.
- **`BODY3D.SETVELOCITY`** - args: handle, float, float, float — Set linear velocity.
- **`BODYREF.FREE`** - args: handle — Destroy a physics body (handle method).
- **`BODYREF.GETPOSITION`** - args: handle -> returns handle
- **`BODYREF.GETROTATION`** - args: handle -> returns handle
- **`BODYREF.GETVELOCITY`** - args: handle -> returns handle
- **`BODYREF.SETVELOCITY`** - args: handle, float, float, float
- **`DEBUG.DRAWBODY`** - args: any — Debug draw a physics body wireframe (no-op on stub builds).
- **`DEBUG.DRAWCHARACTER`** - args: any — Debug draw a KCC capsule (no-op on stub builds).
- **`JOINT.CREATEHINGE`** - args: handle, handle, float, float, float, float, float, float -> returns handle — Creates a hinge joint between two bodies at (px,py,pz) around axis (ax,ay,az).
- **`JOINT.CREATEPOINT`** - args: handle, handle, float, float, float -> returns handle — Creates a point-to-point (ball and socket) joint between two bodies at (px,py,pz).
- **`JOINT.FREE`** - args: handle — Destroys a physics joint/constraint.
- **`JOINT.MAKEHINGE`** - args: handle, handle, float, float, float, float, float, float -> returns handle — DEPRECATED alias of JOINT.CREATEHINGE. Use JOINT.CREATEHINGE. Creates a hinge joint between two bodies at (px,py,pz) around axis (ax,ay,az).
- **`JOINT.MAKEPOINT`** - args: handle, handle, float, float, float -> returns handle — DEPRECATED alias of JOINT.CREATEPOINT. Use JOINT.CREATEPOINT. Creates a point-to-point (ball and socket) joint between two bodies at (px,py,pz).
- **`PHYSICS.AUTO`** - args: int, string, float — Alias for ENTITY.PHYSICS.
- **`PHYSICS.BOUNCE`** - args: int, float — Modular building: Sets bounciness (restitution) for a pending physics body.
- **`PHYSICS.BUILD`** - args: int, float — Modular building: Finalizes and commits the physics body with given mass.
- **`PHYSICS.FORCE`** - args: int, float, float, float — Entity-First: Applies a continuous force to an entity's physics body.
- **`PHYSICS.FRICTION`** - args: int, float — Modular building: Sets friction for a pending physics body.
- **`PHYSICS.GETGRAVITYX`** - args: (none) -> returns float — Get gravity X (alias of PHYSICS3D.GETGRAVITYX).
- **`PHYSICS.GETGRAVITYY`** - args: (none) -> returns float — Get gravity Y (alias of PHYSICS3D.GETGRAVITYY).
- **`PHYSICS.GETGRAVITYZ`** - args: (none) -> returns float — Get gravity Z (alias of PHYSICS3D.GETGRAVITYZ).
- **`PHYSICS.GRAVITY`** - args: int, float — Entity-First: Scale the gravity factor for a specific entity (e.g. 0.0 for zero-g).
- **`PHYSICS.IMPULSE`** - args: int, float, float, float — Entity-First: Applies an instant impulse to an entity's physics body.
- **`PHYSICS.SETROT`** - args: int, float, float, float — Entity-First: Instantly sets the rotation of an entity's physics body (Euler radians).
- **`PHYSICS.SHAPE`** - args: int, string — Modular building: Sets the physics shape for a pending body.
- **`PHYSICS.SIZE`** - args: int, float, float, float — Modular building: Sets dimensions for a pending physics shape.
- **`PHYSICS.VELOCITY`** - args: int, float, float, float — Entity-First: Sets the linear velocity of an entity's physics body.
- **`PHYSICS.WAKE`** - args: int — Entity-First: Forces a sleeping physics body to wake up.
- **`PHYSICS3D.DEBUGDRAW`** - args: int
- **`PHYSICS3D.GETGRAVITYX`** - args: (none) -> returns float — Get current gravity X component.
- **`PHYSICS3D.GETGRAVITYY`** - args: (none) -> returns float — Get current gravity Y component.
- **`PHYSICS3D.GETGRAVITYZ`** - args: (none) -> returns float — Get current gravity Z component.
- **`PHYSICS3D.GETMATRIXBUFFER`** - args: (none) -> returns handle — Get the shared matrix buffer for render interpolation.
- **`PHYSICS3D.GETSCRATCHFLOAT`** - args: int -> returns float — Read a scratch float from the physics scratch buffer.
- **`PHYSICS3D.MOUSEHIT`** - args: handle -> returns handle — Raycast from mouse through camera; returns [x,y,z] array or nil.
- **`PHYSICS3D.ONCOLLISION`** - args: handle, handle, string
- **`PHYSICS3D.PROCESSCOLLISIONS`** - args: (none)
- **`PHYSICS3D.RAYCAST`** - args: float, float, float, float, float, float, float -> returns handle
- **`PHYSICS3D.SETGRAVITY`** - args: float, float, float
- **`PHYSICS3D.SETSUBSTEPS`** - args: int
- **`PHYSICS3D.SETTIMESTEP`** - args: float — Set the fixed physics simulation timestep (e.g. 60.0, 90.0).
- **`PHYSICS3D.START`** - args: (none)
- **`PHYSICS3D.STEP`** - args: (none)
- **`PHYSICS3D.STOP`** - args: (none)
- **`PHYSICS3D.SYNCWASMTOPHYSREGS`** - args: int, int — Sync WASM physics view to VM registers (count, firstReg).
- **`PHYSICS3D.UPDATE`** - args: (none) — Advance the 3D physics simulation (same implementation as PHYSICS3D.STEP; optional frame dt# like STEP)
- **`SHAPE.GETDEPTH`** - args: handle -> returns float — Get shape dimension 3 (half-extent Z).
- **`SHAPE.GETHEIGHT`** - args: handle -> returns float — Get shape dimension 2 (half-extent Y or height).
- **`SHAPE.GETRADIUS`** - args: handle -> returns float — Get shape radius (same as SHAPE.GETWIDTH for spheres).
- **`SHAPE.GETSIZEX`** - args: handle -> returns float — Get shape X dimension.
- **`SHAPE.GETSIZEY`** - args: handle -> returns float — Get shape Y dimension.
- **`SHAPE.GETSIZEZ`** - args: handle -> returns float — Get shape Z dimension.
- **`SHAPE.GETTYPE`** - args: handle -> returns int — Get the shape type (1=Box, 2=Sphere, 3=Capsule, 4=Cylinder).
- **`SHAPE.GETWIDTH`** - args: handle -> returns float — Get shape dimension 1 (half-extent X or radius).
- **`SHAPEREF.FREE`** - args: handle — Destroy a collision shape.
- **`VEHICLE.CONTROL`** - args: int, float, float, float — Update all vehicle inputs (vid, throttle, steer, brake).
- **`VEHICLE.CREATE`** - args: int, int -> returns int — Create a vehicle controller for an entity.
- **`VEHICLE.MAKE`** - args: int, int -> returns int — DEPRECATED alias of VEHICLE.CREATE. Use VEHICLE.CREATE. Create a vehicle controller for an entity.
- **`VEHICLE.SETSTEER`** - args: int, float — Set vehicle steering (-1 to 1).
- **`VEHICLE.SETTHROTTLE`** - args: int, float — Set vehicle throttle (-1 to 1).
- **`VEHICLE.SETTUNING`** - args: int, float, float, float, float — Tune suspension (vid, spring, damp, maxSpeed, steerSpeed).
- **`VEHICLE.SETWHEEL`** - args: int, int, float, float, float, float — Configure a wheel (vid, idx, ox, oy, oz, radius).
- **`VEHICLE.STEP`** - args: float — Step all vehicle simulations by dt.
- **`VEHICLE.WHEELX`** - args: int, int -> returns float
- **`VEHICLE.WHEELY`** - args: int, int -> returns float
- **`VEHICLE.WHEELZ`** - args: int, int -> returns float
- **`WORLD.SETUP`** - args: (none) — Initialise physics world with default gravity (-9.81).
- **`WORLD.SETUP`** - args: float — Initialise physics world with custom Y gravity.

### PI

- **`PI`** - args: (none)

### PICK

- **`PICK.CAST`** - args: (none) -> returns int — Run Jolt raycast from staged params; returns entity# or 0
- **`PICK.DIRECTION`** - args: float, float, float — Stage ray direction; length is max travel unless PICK.MAXDIST set
- **`PICK.DIST`** - args: (none) -> returns float — Distance along ray to last hit
- **`PICK.ENTITY`** - args: (none) -> returns int — Entity# from last pick (linked BODY3D only)
- **`PICK.FROMCAMERA`** - args: handle, float, float — Stage ray from camera handle and screen pixels (sets default MAXDIST if unset)
- **`PICK.HIT`** - args: (none) -> returns bool — Whether last PICK.CAST / SCREENCAST hit
- **`PICK.LAYERMASK`** - args: int — Bit i accepts ENTITY.COLLISIONLAYER i; 0 accepts all
- **`PICK.MAXDIST`** - args: float — Optional max ray length (normalize direction then scale)
- **`PICK.NX`** - args: (none) -> returns float — Last pick surface normal X
- **`PICK.NY`** - args: (none) -> returns float — Last pick surface normal Y
- **`PICK.NZ`** - args: (none) -> returns float — Last pick surface normal Z
- **`PICK.ORIGIN`** - args: float, float, float — Stage ray origin for PICK.CAST (Linux+CGO Jolt)
- **`PICK.RADIUS`** - args: float — Reserved; non-zero returns error until sphere pick exists
- **`PICK.SCREENCAST`** - args: handle, float, float -> returns int — FROMCAMERA then CAST; returns entity# or 0
- **`PICK.X`** - args: (none) -> returns float — Last pick hit world X
- **`PICK.Y`** - args: (none) -> returns float — Last pick hit world Y
- **`PICK.Z`** - args: (none) -> returns float — Last pick hit world Z

### PINGPONG

- **`PINGPONG`** - args: any, any

### PLAYER

- **`PLAYER.ADDIMPULSE`** - args: int, float, float, float — Adds to KCC linear velocity (velocity delta; no separate mass impulse)
- **`PLAYER.CREATE`** - args: int
- **`PLAYER.CREATE`** - args: handle — Initializes a Kinematic Character Controller in the Jolt buffer.
- **`PLAYER.CREATE`** - args: int, float, float — KCC with explicit capsule radius and height (world units)
- **`PLAYER.GETCAPSULEHEIGHT`** - args: (none) -> returns float — Alias of PLAYER.GETHEIGHT ()
- **`PLAYER.GETCAPSULEHEIGHT`** - args: int -> returns float
- **`PLAYER.GETCAPSULERADIUS`** - args: (none) -> returns float — Alias of PLAYER.GETRADIUS ()
- **`PLAYER.GETCAPSULERADIUS`** - args: int -> returns float
- **`PLAYER.GETCEILING`** - args: int -> returns bool — True if last integration saw a strong ceiling contact (head bump)
- **`PLAYER.GETCOLLISIONENABLED`** - args: int -> returns bool — Reserved; returns true
- **`PLAYER.GETCROUCH`** - args: int -> returns bool — Stored crouch flag (capsule resize not in Jolt wrapper yet)
- **`PLAYER.GETFOVKICK`** - args: int -> returns float — Reads stored FOV kick offset (degrees)
- **`PLAYER.GETFRICTION`** - args: int -> returns float — Gameplay friction (CHARACTERREF.SETFRICTION)
- **`PLAYER.GETGRAVITY`** - args: (none) -> returns float — Alias of PLAYER.GETGRAVITYSCALE (per-entity gravity scale, not world gravity)
- **`PLAYER.GETGRAVITY`** - args: int -> returns float
- **`PLAYER.GETGRAVITYSCALE`** - args: int -> returns float
- **`PLAYER.GETGROUNDED`** - args: (none) -> returns bool — Alias of PLAYER.ISGROUNDED ()
- **`PLAYER.GETGROUNDED`** - args: int -> returns bool
- **`PLAYER.GETGROUNDED`** - args: int, float -> returns bool — Coyote time variant (Linux+Jolt)
- **`PLAYER.GETGROUNDSTATE`** - args: (none) -> returns int — Implicit KCC subject (see PLAYER.ISGROUNDED ())
- **`PLAYER.GETGROUNDSTATE`** - args: int -> returns int — Jolt EGroundState: 0 OnGround, 1 OnSteepGround, 2 NotSupported, 3 InAir
- **`PLAYER.GETGROUNDVELOCITYX`** - args: int -> returns float — Ground/platform velocity X (Jolt GetGroundVelocity)
- **`PLAYER.GETGROUNDVELOCITYY`** - args: int -> returns float — Ground/platform velocity Y (Jolt GetGroundVelocity)
- **`PLAYER.GETGROUNDVELOCITYZ`** - args: int -> returns float — Ground/platform velocity Z (Jolt GetGroundVelocity)
- **`PLAYER.GETHEIGHT`** - args: int -> returns float — Capsule total height
- **`PLAYER.GETISFALLING`** - args: int -> returns bool
- **`PLAYER.GETISJUMPING`** - args: int -> returns bool
- **`PLAYER.GETISSLIDING`** - args: int -> returns bool — True when Jolt CharacterVirtual ground state is steep/slide
- **`PLAYER.GETLAYER`** - args: int -> returns int — KCC object layer id (Jolt CHARACTER layer = 2)
- **`PLAYER.GETLOOKTARGET`** - args: int, float -> returns int
- **`PLAYER.GETMASK`** - args: int -> returns int — Bitmask of object layers (0..4); typical full mask 31
- **`PLAYER.GETMAXSLOPE`** - args: int -> returns float — Configured max walk slope (degrees)
- **`PLAYER.GETNEARBY`** - args: int, float, string -> returns handle
- **`PLAYER.GETONSLOPE`** - args: (none) -> returns bool — Implicit KCC subject — steep slope (Jolt OnSteepGround)
- **`PLAYER.GETONSLOPE`** - args: int -> returns bool — Alias of PLAYER.ISONSTEEPSLOPE
- **`PLAYER.GETONWALL`** - args: int -> returns bool — Jolt ground state NotSupported (vertical contact)
- **`PLAYER.GETPITCH`** - args: (none) -> returns float — Alias of PLAYER.GETROTATIONPITCH ()
- **`PLAYER.GETPITCH`** - args: int -> returns float
- **`PLAYER.GETPOSITIONX`** - args: (none) -> returns float — Implicit KCC subject — same as Player.GetPositionX() with no args
- **`PLAYER.GETPOSITIONX`** - args: int -> returns float — KCC / entity world X; (entity) or () using implicit subject after PLAYER.CREATE / Character.Create
- **`PLAYER.GETPOSITIONY`** - args: (none) -> returns float — Implicit KCC subject
- **`PLAYER.GETPOSITIONY`** - args: int -> returns float — (entity) or () implicit KCC subject
- **`PLAYER.GETPOSITIONZ`** - args: (none) -> returns float — Implicit KCC subject
- **`PLAYER.GETPOSITIONZ`** - args: int -> returns float — (entity) or () implicit KCC subject
- **`PLAYER.GETRADIUS`** - args: int -> returns float — Capsule radius
- **`PLAYER.GETROLL`** - args: (none) -> returns float — Alias of PLAYER.GETROTATIONROLL ()
- **`PLAYER.GETROLL`** - args: int -> returns float
- **`PLAYER.GETROTATIONPITCH`** - args: int -> returns float — Entity world pitch (degrees), same basis as ENTITY.ENTITYPITCH
- **`PLAYER.GETROTATIONROLL`** - args: int -> returns float
- **`PLAYER.GETROTATIONYAW`** - args: int -> returns float
- **`PLAYER.GETSHAPETYPE`** - args: (none) -> returns string — Returns "capsule" for CharacterVirtual
- **`PLAYER.GETSHAPETYPE`** - args: int -> returns string
- **`PLAYER.GETSLOPEANGLE`** - args: int -> returns float — Ground tilt from vertical (degrees) on walkable floor
- **`PLAYER.GETSNAPDISTANCE`** - args: int -> returns float — Stick-to-floor max step down (positive); host uses CHAR.STICK or stepH+0.2
- **`PLAYER.GETSPEED`** - args: (none) -> returns float — Implicit KCC subject — same as Player.GetSpeed(entity)
- **`PLAYER.GETSPEED`** - args: int -> returns float — Scalar speed (m/s)
- **`PLAYER.GETSTANDNORMAL`** - args: int -> returns handle — Vec3 ground/floor normal under the player (CharacterVirtual or downward ray)
- **`PLAYER.GETSTEPHEIGHT`** - args: int -> returns float — Stair step-up height (WalkStairsStepUp Y on Jolt)
- **`PLAYER.GETSUBMERGEDFACTOR`** - args: int -> returns float — 0..1 fraction of capsule height below WATER surface (requires WATER.* volumes)
- **`PLAYER.GETSURFACETYPE`** - args: int -> returns string — Footstep label from downward ray hit entity metadata / Blender tag (else Default)
- **`PLAYER.GETVELOCITY`** - args: int -> returns handle — Heap vec3 of linear velocity (CharacterVirtual); requires PLAYER.CREATE
- **`PLAYER.GETVELOCITYX`** - args: int -> returns float — Linear velocity X (m/s); requires PLAYER.CREATE / CHAR.CREATE (deprecated CHAR.MAKE)
- **`PLAYER.GETVELOCITYY`** - args: int -> returns float
- **`PLAYER.GETVELOCITYZ`** - args: int -> returns float
- **`PLAYER.GETVX`** - args: int — Returns horizontal velocity X for kinematic controller
- **`PLAYER.GETVY`** - args: int — Returns vertical velocity Y for kinematic controller
- **`PLAYER.GETVZ`** - args: int — Returns horizontal velocity Z for kinematic controller
- **`PLAYER.GETX`** - args: (none) -> returns float — Alias of PLAYER.GETPOSITIONX () — implicit KCC subject
- **`PLAYER.GETX`** - args: int -> returns float — Alias of PLAYER.GETPOSITIONX (entity)
- **`PLAYER.GETY`** - args: (none) -> returns float — Alias of PLAYER.GETPOSITIONY ()
- **`PLAYER.GETY`** - args: int -> returns float
- **`PLAYER.GETYAW`** - args: (none) -> returns float — Alias of PLAYER.GETROTATIONYAW ()
- **`PLAYER.GETYAW`** - args: int -> returns float
- **`PLAYER.GETZ`** - args: (none) -> returns float — Alias of PLAYER.GETPOSITIONZ ()
- **`PLAYER.GETZ`** - args: int -> returns float
- **`PLAYER.GRAB`** - args: int, int — Welds target to player front each frame (target 0 releases); not a Jolt fixed constraint yet
- **`PLAYER.ISGROUNDED`** - args: (none) -> returns bool — Implicit KCC subject: the last PLAYER.CREATE / Character.Create capsule (entity-bound Jolt KCC)
- **`PLAYER.ISGROUNDED`** - args: int -> returns bool
- **`PLAYER.ISGROUNDED`** - args: int, float -> returns bool — Optional coyote time (seconds): true shortly after leaving ground
- **`PLAYER.ISMOVING`** - args: int -> returns bool — True if horizontal linear speed > ~0.05 (CharacterVirtual)
- **`PLAYER.ISONSTEEPSLOPE`** - args: (none) -> returns bool — Implicit KCC subject (see PLAYER.ISGROUNDED ())
- **`PLAYER.ISONSTEEPSLOPE`** - args: int -> returns bool — True if GetGroundState is OnSteepGround (Jolt CharacterVirtual)
- **`PLAYER.ISSUBMERGED`** - args: int -> returns bool — True when GETSUBMERGEDFACTOR exceeds ~45%
- **`PLAYER.ISSWIMMING`** - args: int -> returns bool — True when entity origin is inside a WATER volume column (bed..surface)
- **`PLAYER.JUMP`** - args: int, float
- **`PLAYER.MAKE`** - args: int — DEPRECATED alias of PLAYER.CREATE. Use PLAYER.CREATE.
- **`PLAYER.MAKE`** - args: handle — DEPRECATED alias of PLAYER.CREATE. Use PLAYER.CREATE. Initializes a Kinematic Character Controller in the Jolt buffer.
- **`PLAYER.MAKE`** - args: int, float, float — DEPRECATED alias of PLAYER.CREATE. Use PLAYER.CREATE. KCC with explicit capsule radius and height (world units)
- **`PLAYER.MOVE`** - args: int, float, float
- **`PLAYER.MOVERELATIVE`** - args: float, float, float, float, float -> returns handle — MOVESTEPX/Z combined — 2-float array [dx,dz]; ERASE when done
- **`PLAYER.MOVEWITHCAMERA`** - args: int, handle, float, float, float — WASD-style: (entity, camera, forwardAxis#, strafeAxis#, speed#) movement on XZ relative to camera view (Linux+CGO KCC)
- **`PLAYER.NAVTO`** - args: int, float, float, float — Click-to-move target: (entity, targetX#, targetZ#, speed# [, arrivalXZ# [, brakeDist#]]); use with PLAYER.NAVUPDATE each frame; soft brake near target (Linux+CGO KCC)
- **`PLAYER.NAVTO`** - args: int, float, float, float, float — NAVTO with arrival distance
- **`PLAYER.NAVTO`** - args: int, float, float, float, float, float — NAVTO with arrival and brake distance (soft stop)
- **`PLAYER.NAVUPDATE`** - args: int — Advances PLAYER.NAVTO / CHAR.NAVTO toward target with soft deceleration (Linux+CGO)
- **`PLAYER.ONTRIGGER`** - args: int, string
- **`PLAYER.PUSH`** - args: int, int, float — Applies forward horizontal force to target entity (host ENTITY.ADDFORCE path; scaled by player mass)
- **`PLAYER.SETAIRCONTROL`** - args: int, float — Scales horizontal PLAYER.MOVE input while airborne (1 = default)
- **`PLAYER.SETCROUCH`** - args: int, any — Sets crouch flag (gameplay; capsule height unchanged until wrapper supports it)
- **`PLAYER.SETFOVKICK`** - args: int, float — Stores extra FOV degrees; add Camera.SetFOV(base + Player.GetFovKick(id)) each frame
- **`PLAYER.SETGRAVITYSCALE`** - args: int, float — Scales CharacterVirtual gravity on Y (1=default; moon-jump / low-G zones)
- **`PLAYER.SETGROUNDCONTROL`** - args: int, float — Scales horizontal PLAYER.MOVE input while on ground (1 = default)
- **`PLAYER.SETJUMPBUFFER`** - args: int, float — Jump buffer window (seconds) when jump pressed in air; consumed on landing
- **`PLAYER.SETMASS`** - args: int, float — Stores gameplay mass (PLAYER.Push scaling); Jolt capsule mass is fixed at create
- **`PLAYER.SETPADDING`** - args: int, float — Character capsule skin padding (world units, >0); rebuilds CharacterVirtual (Linux+CGO)
- **`PLAYER.SETSLOPELIMIT`** - args: int, float — Rebuilds CharacterVirtual with MaxSlopeAngle = angle (degrees); requires PLAYER.CREATE (Linux+Jolt)
- **`PLAYER.SETSTATE`** - args: int, int
- **`PLAYER.SETSTEPHEIGHT`** - args: int, float — Stores max stair/curb step height for the player entity (reserved; Jolt runtime step tuning not exposed yet)
- **`PLAYER.SETSTEPOFFSET`** - args: int, float — Alias of PLAYER.SETSTEPHEIGHT; maps to Jolt ExtendedUpdate WalkStairsStepUp (Linux+CGO)
- **`PLAYER.SETSTICKFLOOR`** - args: int, float — Stick-to-floor max step down (world units); Jolt CharacterVirtual ExtendedUpdateSettings (Linux+CGO)
- **`PLAYER.SETVELOCITY`** - args: int, float, float, float — Sets KCC linear velocity (m/s); alias for CHARACTERREF.SETLINEARVELOCITY
- **`PLAYER.SNAPTOGROUND`** - args: int, handle, float — Sets entity Y from terrain height at entity XZ + offset (PLAYER.CREATE syncs capsule on Linux+Jolt)
- **`PLAYER.SWIM`** - args: int, float, float — Swim mode: buoyancy reduces downward gravity; drag damps horizontal motion; (0,0) disables
- **`PLAYER.SYNCANIM`** - args: int, any
- **`PLAYER.TELEPORT`** - args: int, float, float, float — Snaps capsule and entity to (x,y,z), clears linear velocity (no smoothing)
- **`PLAYER.UPDATE`** - args: float — Update kinematic character solver with delta time (legacy; prefer CHARACTERREF.UPDATE / UPDATEPHYSICS)

### PLAYER2D

- **`PLAYER2D.CLAMP`** - args: handle, float, float, float, float
- **`PLAYER2D.CREATE`** - args: (none) -> returns handle
- **`PLAYER2D.FREE`** - args: handle
- **`PLAYER2D.GETX`** - args: handle -> returns float
- **`PLAYER2D.GETZ`** - args: handle -> returns float
- **`PLAYER2D.KEEPINBOUNDS`** - args: handle
- **`PLAYER2D.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of PLAYER2D.CREATE. Use PLAYER2D.CREATE.
- **`PLAYER2D.MOVE`** - args: handle, float, float, float, float, float
- **`PLAYER2D.SETPOS`** - args: handle, float, float
- **`PLAYER2D.SETPOSITION`** - args: handle, float, float — DEPRECATED alias of PLAYER2D.SETPOS. Use PLAYER2D.SETPOS.

### PLAYMUSIC

- **`PLAYMUSIC`** - args: handle — Easy Mode: AUDIO.PLAY(music)

### PLAYSOUND

- **`PLAYSOUND`** - args: handle -> returns void — Easy Mode: Play a sound
- **`PLAYSOUND`** - args: handle — Easy Mode: AUDIO.PLAY(sound)

### POINT3D

- **`POINT3D`** - args: float, float, float, int, int, int, int — Shorthand: DRAW3D.POINT(x, y, z, r, g, b, a)

### POINTENTITY

- **`POINTENTITY`** - args: handle, handle -> returns void — Easy Mode: Point one entity at another

### POOL

- **`POOL.CREATE`** - args: string, int -> returns handle
- **`POOL.FREE`** - args: handle
- **`POOL.GET`** - args: handle -> returns handle
- **`POOL.MAKE`** - args: string, int -> returns handle — DEPRECATED alias of POOL.CREATE. Use POOL.CREATE.
- **`POOL.PREWARM`** - args: handle
- **`POOL.RETURN`** - args: handle, handle
- **`POOL.SETFACTORY`** - args: handle, string
- **`POOL.SETRESET`** - args: handle, string

### POSENT

- **`POSENT`** - args: int, float, float, float — Easy Mode: ENTITY.SETPOS(ent, x, y, z) (Blitz alias ENTITY.POSITIONENTITY)
- **`POSENT`** - args: handle, float, float, float — Shorthand: ENTITY.SETPOS(ent, x, y, z) (same as POSITIONENTITY; Blitz ENTITY.POSITIONENTITY)

### POSITIONCAMERA

- **`POSITIONCAMERA`** - args: handle, float, float, float — Easy Mode: CAMERA.SETPOS(cam, x, y, z)

### POSITIONENTITY

- **`POSITIONENTITY`** - args: int, float, float, float — Blitz-style global: ENTITY.SETPOS(ent, x, y, z); Blitz namespace alias ENTITY.POSITIONENTITY; deprecated ENTITY.SETPOSITION

### POST

- **`POST.ADD`** - args: string
- **`POST.ADDSHADER`** - args: handle
- **`POST.REMOVE`** - args: string
- **`POST.SETPARAM`** - args: string, string, float
- **`POST.SETTONEMAP`** - args: int

### POW

- **`POW`** - args: any, any

### PP_BLOOM

- **`PP_BLOOM`** - args: (none) -> returns int

### PP_CRT_SCANLINES

- **`PP_CRT_SCANLINES`** - args: (none) -> returns int

### PP_PIXELATE

- **`PP_PIXELATE`** - args: (none) -> returns int

### PRINT

- **`PRINT`** - args: any — Print values to stdout, space-separated, with newline.

### PRINTAT

- **`PRINTAT`** - args: int, int, any

### PRINTCOLOR

- **`PRINTCOLOR`** - args: int, int, int, any

### PRINTLN

- **`PRINTLN`** - args: any — Same as PRINT (newline after output).

### PaintEntity

- **`PaintEntity`** - args: int, handle

### QUAT

- **`QUAT.FREE`** - args: handle
- **`QUAT.FROMAXISANGLE`** - args: float, float, float, float -> returns handle
- **`QUAT.FROMEULER`** - args: float, float, float -> returns handle
- **`QUAT.FROMMAT4`** - args: handle -> returns handle
- **`QUAT.FROMVEC3TOVEC3`** - args: handle, handle -> returns handle
- **`QUAT.IDENTITY`** - args: (none) -> returns handle
- **`QUAT.INVERT`** - args: handle -> returns handle
- **`QUAT.MULTIPLY`** - args: handle, handle -> returns handle
- **`QUAT.NORMALIZE`** - args: handle -> returns handle
- **`QUAT.SLERP`** - args: handle, handle, float -> returns handle
- **`QUAT.TOEULER`** - args: handle -> returns handle
- **`QUAT.TOMAT4`** - args: handle -> returns handle
- **`QUAT.TRANSFORM`** - args: handle, handle -> returns handle

### QUIT

- **`QUIT`** - args: (none)

### RAD2DEG

- **`RAD2DEG`** - args: any

### RAND

- **`RAND`** - args: any, any -> returns int — Same as RND(min, max) — inclusive integer range
- **`RAND`** - args: int, int -> returns int — Easy Mode: Random int in range
- **`RAND.CREATE`** - args: int -> returns handle
- **`RAND.FREE`** - args: handle
- **`RAND.MAKE`** - args: int -> returns handle — DEPRECATED alias of RAND.CREATE. Use RAND.CREATE.
- **`RAND.NEXT`** - args: handle, int, int -> returns int
- **`RAND.NEXTF`** - args: handle -> returns float

### RANDOMIZE

- **`RANDOMIZE`** - args: (none)
- **`RANDOMIZE`** - args: any

### RAY

- **`RAY.CREATE`** - args: float, float, float, float, float, float -> returns handle
- **`RAY.FREE`** - args: handle
- **`RAY.HITBOX_DISTANCE`** - args: handle, float, float, float, float, float, float -> returns float
- **`RAY.HITBOX_HIT`** - args: handle, float, float, float, float, float, float -> returns bool
- **`RAY.HITBOX_NORMALX`** - args: handle, float, float, float, float, float, float -> returns float
- **`RAY.HITBOX_NORMALY`** - args: handle, float, float, float, float, float, float -> returns float
- **`RAY.HITBOX_NORMALZ`** - args: handle, float, float, float, float, float, float -> returns float
- **`RAY.HITBOX_POINTX`** - args: handle, float, float, float, float, float, float -> returns float
- **`RAY.HITBOX_POINTY`** - args: handle, float, float, float, float, float, float -> returns float
- **`RAY.HITBOX_POINTZ`** - args: handle, float, float, float, float, float, float -> returns float
- **`RAY.HITMESH_DISTANCE`** - args: handle, handle, handle -> returns float
- **`RAY.HITMESH_HIT`** - args: handle, handle, handle -> returns bool
- **`RAY.HITMESH_NORMALX`** - args: handle, handle, handle -> returns float
- **`RAY.HITMESH_NORMALY`** - args: handle, handle, handle -> returns float
- **`RAY.HITMESH_NORMALZ`** - args: handle, handle, handle -> returns float
- **`RAY.HITMESH_POINTX`** - args: handle, handle, handle -> returns float
- **`RAY.HITMESH_POINTY`** - args: handle, handle, handle -> returns float
- **`RAY.HITMESH_POINTZ`** - args: handle, handle, handle -> returns float
- **`RAY.HITMODEL_DISTANCE`** - args: handle, handle -> returns float
- **`RAY.HITMODEL_HIT`** - args: handle, handle -> returns bool
- **`RAY.HITMODEL_NORMALX`** - args: handle, handle -> returns float
- **`RAY.HITMODEL_NORMALY`** - args: handle, handle -> returns float
- **`RAY.HITMODEL_NORMALZ`** - args: handle, handle -> returns float
- **`RAY.HITMODEL_POINTX`** - args: handle, handle -> returns float
- **`RAY.HITMODEL_POINTY`** - args: handle, handle -> returns float
- **`RAY.HITMODEL_POINTZ`** - args: handle, handle -> returns float
- **`RAY.HITPLANE_DISTANCE`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITPLANE_HIT`** - args: handle, float, float, float, float -> returns bool
- **`RAY.HITPLANE_NORMALX`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITPLANE_NORMALY`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITPLANE_NORMALZ`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITPLANE_POINTX`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITPLANE_POINTY`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITPLANE_POINTZ`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITSPHERE_DISTANCE`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITSPHERE_HIT`** - args: handle, float, float, float, float -> returns bool
- **`RAY.HITSPHERE_NORMALX`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITSPHERE_NORMALY`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITSPHERE_NORMALZ`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITSPHERE_POINTX`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITSPHERE_POINTY`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITSPHERE_POINTZ`** - args: handle, float, float, float, float -> returns float
- **`RAY.HITTRIANGLE_DISTANCE`** - args: handle, float, float, float, float, float, float, float, float, float -> returns float
- **`RAY.HITTRIANGLE_HIT`** - args: handle, float, float, float, float, float, float, float, float, float -> returns bool
- **`RAY.HITTRIANGLE_NORMALX`** - args: handle, float, float, float, float, float, float, float, float, float -> returns float
- **`RAY.HITTRIANGLE_NORMALY`** - args: handle, float, float, float, float, float, float, float, float, float -> returns float
- **`RAY.HITTRIANGLE_NORMALZ`** - args: handle, float, float, float, float, float, float, float, float, float -> returns float
- **`RAY.HITTRIANGLE_POINTX`** - args: handle, float, float, float, float, float, float, float, float, float -> returns float
- **`RAY.HITTRIANGLE_POINTY`** - args: handle, float, float, float, float, float, float, float, float, float -> returns float
- **`RAY.HITTRIANGLE_POINTZ`** - args: handle, float, float, float, float, float, float, float, float, float -> returns float
- **`RAY.INTERSECTSMODEL_DISTANCE`** - args: handle, handle -> returns float
- **`RAY.INTERSECTSMODEL_HIT`** - args: handle, handle -> returns bool — Alias of RAY.HITMODEL_HIT — ray vs loaded MODEL mesh union
- **`RAY.INTERSECTSMODEL_NORMALX`** - args: handle, handle -> returns float
- **`RAY.INTERSECTSMODEL_NORMALY`** - args: handle, handle -> returns float
- **`RAY.INTERSECTSMODEL_NORMALZ`** - args: handle, handle -> returns float
- **`RAY.INTERSECTSMODEL_POINTX`** - args: handle, handle -> returns float
- **`RAY.INTERSECTSMODEL_POINTY`** - args: handle, handle -> returns float
- **`RAY.INTERSECTSMODEL_POINTZ`** - args: handle, handle -> returns float
- **`RAY.MAKE`** - args: float, float, float, float, float, float -> returns handle — DEPRECATED alias of RAY.CREATE. Use RAY.CREATE.

### RAY2D

- **`RAY2D.HITCIRCLE_DISTANCE`** - args: float, float, float, float, float, float, float -> returns float — Distance along ray to hit (0 if miss)
- **`RAY2D.HITCIRCLE_HIT`** - args: float, float, float, float, float, float, float -> returns bool — 2D ray vs circle — hit?
- **`RAY2D.HITCIRCLE_POINTX`** - args: float, float, float, float, float, float, float -> returns float
- **`RAY2D.HITCIRCLE_POINTY`** - args: float, float, float, float, float, float, float -> returns float
- **`RAY2D.HITRECT_DISTANCE`** - args: float, float, float, float, float, float, float, float -> returns float
- **`RAY2D.HITRECT_HIT`** - args: float, float, float, float, float, float, float, float -> returns bool — 2D ray vs axis-aligned rect (minx,miny,maxx,maxy)
- **`RAY2D.HITRECT_POINTX`** - args: float, float, float, float, float, float, float, float -> returns float
- **`RAY2D.HITRECT_POINTY`** - args: float, float, float, float, float, float, float, float -> returns float
- **`RAY2D.HITSEGMENT_DISTANCE`** - args: float, float, float, float, float, float, float, float -> returns float
- **`RAY2D.HITSEGMENT_HIT`** - args: float, float, float, float, float, float, float, float -> returns bool — 2D ray vs segment (x1,y1)-(x2,y2)
- **`RAY2D.HITSEGMENT_POINTX`** - args: float, float, float, float, float, float, float, float -> returns float
- **`RAY2D.HITSEGMENT_POINTY`** - args: float, float, float, float, float, float, float, float -> returns float

### RAYLIB

- **`RAYLIB.BEGINFRAME`** - args: (none)
- **`RAYLIB.BEGINSHADERMODE`** - args: handle
- **`RAYLIB.CLEARBACKGROUND`** - args: int, int, int
- **`RAYLIB.CLOSEWINDOW`** - args: (none)
- **`RAYLIB.DRAWCIRCLE`** - args: int, int, float, int, int, int, int
- **`RAYLIB.DRAWCUBE`** - args: float, float, float, float, float, float, int, int, int, int
- **`RAYLIB.DRAWFPS`** - args: int, int
- **`RAYLIB.DRAWLINE3D`** - args: float, float, float, float, float, float, int, int, int, int
- **`RAYLIB.DRAWMODEL`** - args: handle, float, float, float, float, float, float, float
- **`RAYLIB.DRAWRECTANGLE`** - args: int, int, int, int, int, int, int, int
- **`RAYLIB.DRAWSPHERE`** - args: float, float, float, float, int, int, int, int
- **`RAYLIB.DRAWTEXTURE`** - args: handle, int, int, int, int, int, int, int, int
- **`RAYLIB.ENDFRAME`** - args: (none)
- **`RAYLIB.ENDSHADERMODE`** - args: (none)
- **`RAYLIB.GETFPS`** - args: (none)
- **`RAYLIB.GETFRAMEBUFFERHEIGHT`** - args: (none)
- **`RAYLIB.GETFRAMEBUFFERWIDTH`** - args: (none)
- **`RAYLIB.GETMOUSEX`** - args: (none)
- **`RAYLIB.GETMOUSEY`** - args: (none)
- **`RAYLIB.GETTIME`** - args: (none)
- **`RAYLIB.INITWINDOW`** - args: int, int, string
- **`RAYLIB.ISKEYDOWN`** - args: int
- **`RAYLIB.ISKEYPRESSED`** - args: int
- **`RAYLIB.ISKEYRELEASED`** - args: int
- **`RAYLIB.ISMOUSEBUTTONDOWN`** - args: int
- **`RAYLIB.LOADMODEL`** - args: string
- **`RAYLIB.LOADSHADER`** - args: string, string
- **`RAYLIB.LOADTEXTURE`** - args: string
- **`RAYLIB.SETCAMERAMODE`** - args: handle, int
- **`RAYLIB.SETTARGETFPS`** - args: int
- **`RAYLIB.UNLOADTEXTURE`** - args: handle
- **`RAYLIB.UPDATECAMERA`** - args: handle, int
- **`RAYLIB.WINDOWSHOULDCLOSE`** - args: (none)

### READALLTEXT

- **`READALLTEXT`** - args: string
- **`READALLTEXT`** - args: string -> returns string

### READALLTEXT$

- **`READALLTEXT$`** - args: string
- **`READALLTEXT$`** - args: string -> returns string

### READBYTE

- **`READBYTE`** - args: handle

### READFILE

- **`READFILE`** - args: any

### READFILE$

- **`READFILE$`** - args: any

### READFLOAT

- **`READFLOAT`** - args: handle

### READINT

- **`READINT`** - args: handle

### READSHORT

- **`READSHORT`** - args: handle

### READSTRING

- **`READSTRING`** - args: handle, int

### READSTRING$

- **`READSTRING$`** - args: handle, int

### REMAP

- **`REMAP`** - args: float, float, float, float, float -> returns float

### RENAMEFILE

- **`RENAMEFILE`** - args: string, string
- **`RENAMEFILE`** - args: string, string -> returns bool

### RENDER

- **`RENDER.BEGIN3D`** - args: handle — Alias for CAMERA.BEGIN: 3D camera heap handle from CAMERA.CREATE or CreateCamera (deprecated alias CAMERA.MAKE)
- **`RENDER.BEGINFRAME`** - args: (none)
- **`RENDER.BEGINMODE2D`** - args: (none)
- **`RENDER.BEGINMODE3D`** - args: (none)
- **`RENDER.BEGINSHADER`** - args: handle
- **`RENDER.CLEAR`** - args: (none)
- **`RENDER.CLEAR`** - args: handle
- **`RENDER.CLEAR`** - args: int, int, int
- **`RENDER.CLEAR`** - args: int, int, int, int
- **`RENDER.CLEARCACHE`** - args: (none)
- **`RENDER.CLEARSCISSOR`** - args: (none)
- **`RENDER.DRAWFPS`** - args: int, int
- **`RENDER.END3D`** - args: (none) — Alias for CAMERA.END (no arguments)
- **`RENDER.ENDFRAME`** - args: (none)
- **`RENDER.ENDMODE2D`** - args: (none)
- **`RENDER.ENDMODE3D`** - args: (none)
- **`RENDER.ENDSHADER`** - args: (none)
- **`RENDER.FRAME`** - args: (none)
- **`RENDER.SCREENSHOT`** - args: string
- **`RENDER.SET2DAMBIENT`** - args: int, int, int, int
- **`RENDER.SET2DAmbIENT`** - args: int, int, int, int
- **`RENDER.SETAMBIENT`** - args: float, float, float
- **`RENDER.SETAMBIENT`** - args: float, float, float, float
- **`RENDER.SETBLEND`** - args: int
- **`RENDER.SETBLENDMODE`** - args: int
- **`RENDER.SETBLOOM`** - args: float — POST.BLOOM threshold; intensity defaults to 1
- **`RENDER.SETBLOOM`** - args: float, float — POST.BLOOM threshold and intensity
- **`RENDER.SETCULLFACE`** - args: int
- **`RENDER.SETDEPTHMASK`** - args: bool
- **`RENDER.SETDEPTHTEST`** - args: bool
- **`RENDER.SETDEPTHWRITE`** - args: bool
- **`RENDER.SETFOG`** - args: float, float, float, float, float, float — Fog RGB, near, far, density — FOG.* + WORLD.FOGDENSITY
- **`RENDER.SETFPS`** - args: int
- **`RENDER.SETIBLINTENSITY`** - args: float
- **`RENDER.SETIBLSPLIT`** - args: float, float
- **`RENDER.SETMODE`** - args: string
- **`RENDER.SETMSAA`** - args: bool
- **`RENDER.SETSCISSOR`** - args: int, int, int, int
- **`RENDER.SETSHADOWMAPSIZE`** - args: int
- **`RENDER.SETSKYBOX`** - args: string
- **`RENDER.SETTONEMAPPING`** - args: int
- **`RENDER.SETWIREFRAME`** - args: bool

### RENDERTARGET

- **`RENDERTARGET.BEGIN`** - args: handle
- **`RENDERTARGET.CREATE`** - args: int, int -> returns handle
- **`RENDERTARGET.END`** - args: (none)
- **`RENDERTARGET.FREE`** - args: handle
- **`RENDERTARGET.MAKE`** - args: int, int -> returns handle — DEPRECATED alias of RENDERTARGET.CREATE. Use RENDERTARGET.CREATE.
- **`RENDERTARGET.TEXTURE`** - args: handle -> returns handle

### REPEAT

- **`REPEAT`** - args: string, int -> returns string

### REPEAT$

- **`REPEAT$`** - args: string, int -> returns string

### REPLACE

- **`REPLACE`** - args: string, string, string -> returns string

### REPLACE$

- **`REPLACE$`** - args: string, string, string -> returns string

### RES

- **`RES.EXISTS`** - args: string -> returns bool — True if path exists on disk (same idea as UTIL.FILEEXISTS)
- **`RES.PATH`** - args: string -> returns string — Resolve localPath relative to the running executable directory (absolute paths unchanged)

### RESETENTITY

- **`RESETENTITY`** - args: handle -> returns void — Easy Mode: Reset entity velocity and collision state

### REVERSE

- **`REVERSE`** - args: string -> returns string

### REVERSE$

- **`REVERSE$`** - args: string -> returns string

### RIGHT

- **`RIGHT`** - args: string, int -> returns string

### RIGHT$

- **`RIGHT$`** - args: string, int -> returns string

### RND

- **`RND`** - args: (none) — RND() float in [0,1); RND(n) int in [0,n-1]; RND(lo,hi) inclusive int range.
- **`RND`** - args: any
- **`RND`** - args: any, any -> returns int — Inclusive random integer in [min, max]
- **`RND`** - args: float, float -> returns float — Easy Mode: Random float in range

### RNDF

- **`RNDF`** - args: any, any

### RNDSEED

- **`RNDSEED`** - args: any

### ROTATECAMERA

- **`ROTATECAMERA`** - args: handle, float, float, float — Easy Mode: CAMERA.ROTATE(cam, p, y, r)

### ROTATEENTITY

- **`ROTATEENTITY`** - args: int, float, float, float — Blitz-style: ENTITY.ROTATEENTITY(obj, p, y, r)

### ROTENT

- **`ROTENT`** - args: handle, float, float, float — Shorthand: ROTATEENTITY(ent, p, y, r)

### ROUND

- **`ROUND`** - args: any
- **`ROUND`** - args: any, any

### RPC

- **`RPC.CALL`** - args: string
- **`RPC.CALL`** - args: string, any
- **`RPC.CALL`** - args: string, any, any
- **`RPC.CALL`** - args: string, any, any, any
- **`RPC.CALL`** - args: string, any, any, any, any
- **`RPC.CALL`** - args: string, any, any, any, any, any
- **`RPC.CALL`** - args: string, any, any, any, any, any, any
- **`RPC.CALL`** - args: string, any, any, any, any, any, any, any
- **`RPC.CALLSERVER`** - args: string
- **`RPC.CALLSERVER`** - args: string, any
- **`RPC.CALLSERVER`** - args: string, any, any
- **`RPC.CALLSERVER`** - args: string, any, any, any
- **`RPC.CALLSERVER`** - args: string, any, any, any, any
- **`RPC.CALLSERVER`** - args: string, any, any, any, any, any
- **`RPC.CALLSERVER`** - args: string, any, any, any, any, any, any
- **`RPC.CALLSERVER`** - args: string, any, any, any, any, any, any, any
- **`RPC.CALLTO`** - args: handle, string
- **`RPC.CALLTO`** - args: handle, string, any
- **`RPC.CALLTO`** - args: handle, string, any, any
- **`RPC.CALLTO`** - args: handle, string, any, any, any
- **`RPC.CALLTO`** - args: handle, string, any, any, any, any
- **`RPC.CALLTO`** - args: handle, string, any, any, any, any, any
- **`RPC.CALLTO`** - args: handle, string, any, any, any, any, any, any
- **`RPC.CALLTO`** - args: handle, string, any, any, any, any, any, any, any

### RSET

- **`RSET`** - args: string, int -> returns string

### RSET$

- **`RSET$`** - args: string, int -> returns string

### RTRIM

- **`RTRIM`** - args: string -> returns string

### RTRIM$

- **`RTRIM$`** - args: string -> returns string

### SATURATE

- **`SATURATE`** - args: float -> returns float

### SAVE

- **`SAVE.DATA`** - args: string, string — Writes JSON data.
- **`SAVE.GET`** - args: string -> returns string — Reads JSON data.

### SCALENT

- **`SCALENT`** - args: int, float, float, float — Easy Mode: ENTITY.SCALEENTITY(ent, x, y, z)
- **`SCALENT`** - args: handle, float, float, float — Shorthand: SCALEENTITY(ent, x, y, z)

### SCALESPRITE

- **`SCALESPRITE`** - args: handle, float, float -> returns void — Easy Mode: Set sprite X/Y scale

### SCENE

- **`SCENE.APPLYPHYSICS`** - args: handle — Automatically parses glTF Extras to generate Jolt colliders.
- **`SCENE.CLEARSCENE`** - args: (none)
- **`SCENE.CURRENT`** - args: (none) -> returns string
- **`SCENE.DRAW`** - args: (none)
- **`SCENE.LOAD`** - args: string
- **`SCENE.LOADASYNC`** - args: string
- **`SCENE.LOADSCENE`** - args: any
- **`SCENE.LOADWITHTRANSITION`** - args: string, string, float
- **`SCENE.REGISTER`** - args: string, string
- **`SCENE.SAVESCENE`** - args: any
- **`SCENE.SETHANDLERS`** - args: string, string
- **`SCENE.SWITCH`** - args: handle, float — Smoothly transitions levels.
- **`SCENE.UPDATE`** - args: float

### SCREENHEIGHT

- **`SCREENHEIGHT`** - args: (none) -> returns int — Easy Mode: Get window height

### SCREENWIDTH

- **`SCREENWIDTH`** - args: (none) -> returns int — Easy Mode: Get window width

### SECOND

- **`SECOND`** - args: (none)
- **`SECOND`** - args: (none) -> returns int

### SEEKFILE

- **`SEEKFILE`** - args: handle, string, int

### SERVER

- **`SERVER.ONCONNECT`** - args: string
- **`SERVER.ONDISCONNECT`** - args: string
- **`SERVER.ONMESSAGE`** - args: string
- **`SERVER.SETTICKRATE`** - args: float
- **`SERVER.START`** - args: int, int
- **`SERVER.STOP`** - args: (none)
- **`SERVER.SYNCENTITY`** - args: handle, float
- **`SERVER.TICK`** - args: float

### SERVICENET

- **`SERVICENET`** - args: handle, int -> returns int — Easy Mode: NET.SERVICE(host, timeout)

### SETDIR

- **`SETDIR`** - args: string
- **`SETDIR`** - args: string -> returns bool

### SETGRAVITY

- **`SETGRAVITY`** - args: float, float, float — Easy Mode: PHYSICS3D.SETGRAVITY(x, y, z)

### SGN

- **`SGN`** - args: any

### SHADER

- **`SHADER.FREE`** - args: handle
- **`SHADER.FREE`** - args: int — Unloads shader.
- **`SHADER.GETLOC`** - args: handle, string -> returns int
- **`SHADER.LOAD`** - args: string, string
- **`SHADER.SETFLOAT`** - args: int, string, float — Injects a constant float to the custom shader uniformly.
- **`SHADER.SETFLOAT`** - args: handle, string, float
- **`SHADER.SETINT`** - args: handle, string, int
- **`SHADER.SETTEXTURE`** - args: handle, string, handle
- **`SHADER.SETTEXTURE`** - args: int, string, handle — Binds a texture resource to a sampler array element.
- **`SHADER.SETVEC2`** - args: handle, string, float, float
- **`SHADER.SETVEC3`** - args: handle, string, float, float, float
- **`SHADER.SETVEC4`** - args: handle, string, float, float, float, float
- **`SHADER.SETVECTOR`** - args: int, string, float, float, float — Injects a constant vec3.

### SHADER_CEL_STYLED

- **`SHADER_CEL_STYLED`** - args: (none) -> returns int

### SHADER_PBR_LIT

- **`SHADER_PBR_LIT`** - args: (none) -> returns int

### SHADER_PS1_RETRO

- **`SHADER_PS1_RETRO`** - args: (none) -> returns int

### SHADER_WATER_PROCEDURAL

- **`SHADER_WATER_PROCEDURAL`** - args: (none) -> returns int

### SHAKECAMERA

- **`SHAKECAMERA`** - args: handle, float, float — Easy Mode: CAMERA.SHAKE(cam, intensity, duration)

### SHAPE

- **`SHAPE.CREATEBOX`** - args: float, float, float -> returns handle — Creates a Jolt Box shape for collision geometry.
- **`SHAPE.CREATECAPSULE`** - args: float, float -> returns handle — Creates a Jolt Capsule shape: (radius, height).
- **`SHAPE.CREATECYLINDER`** - args: float, float -> returns handle — Creates a Jolt Cylinder shape: (radius, height).
- **`SHAPE.CREATESPHERE`** - args: float -> returns handle — Creates a Jolt Sphere shape.
- **`SHAPE.MAKEBOX`** - args: float, float, float -> returns handle — DEPRECATED alias of SHAPE.CREATEBOX. Use SHAPE.CREATEBOX. Creates a Jolt Box shape for collision geometry.
- **`SHAPE.MAKECAPSULE`** - args: float, float -> returns handle — DEPRECATED alias of SHAPE.CREATECAPSULE. Use SHAPE.CREATECAPSULE. Creates a Jolt Capsule shape: (radius, height).
- **`SHAPE.MAKECYLINDER`** - args: float, float -> returns handle — DEPRECATED alias of SHAPE.CREATECYLINDER. Use SHAPE.CREATECYLINDER. Creates a Jolt Cylinder shape: (radius, height).
- **`SHAPE.MAKESPHERE`** - args: float -> returns handle — DEPRECATED alias of SHAPE.CREATESPHERE. Use SHAPE.CREATESPHERE. Creates a Jolt Sphere shape.

### SHOWENTITY

- **`SHOWENTITY`** - args: handle -> returns void — Easy Mode: Show an entity

### SIGN

- **`SIGN`** - args: any

### SIN

- **`SIN`** - args: any

### SIND

- **`SIND`** - args: any

### SKY

- **`SKY.CREATE`** - args: (none) -> returns handle
- **`SKY.DRAW`** - args: handle
- **`SKY.FREE`** - args: handle
- **`SKY.GETTIMEHOURS`** - args: handle -> returns float
- **`SKY.ISNIGHT`** - args: handle -> returns bool
- **`SKY.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of SKY.CREATE. Use SKY.CREATE.
- **`SKY.SETDAYLENGTH`** - args: handle, float
- **`SKY.SETTIME`** - args: handle, float
- **`SKY.UPDATE`** - args: handle, float

### SKYCOLOR

- **`SKYCOLOR`** - args: int, int, int -> returns void — Easy Mode: Alias for Render.Clear(r, g, b)
- **`SKYCOLOR`** - args: int, int, int

### SLEEP

- **`SLEEP`** - args: any

### SMOOTHERSTEP

- **`SMOOTHERSTEP`** - args: any, any, any -> returns float — Ken Perlin smootherstep(edge0, edge1, x); clamps then 6t^5-15t^4+10t^3

### SMOOTHSTEP

- **`SMOOTHSTEP`** - args: any, any, any

### SOUND

- **`SOUND.ATTACH`** - args: handle, handle — Pins a sound to an entity.
- **`SOUND.FREE`** - args: handle
- **`SOUND.FROMWAVE`** - args: handle -> returns handle
- **`SOUND.PLAY3D`** - args: handle, float, float, float, float — Plays 3D spatialized audio.

### SOUNDVOLUME

- **`SOUNDVOLUME`** - args: handle, float -> returns void — Easy Mode: Set sound volume (0-1)

### SPACE

- **`SPACE`** - args: int -> returns string

### SPACE$

- **`SPACE$`** - args: int -> returns string

### SPC

- **`SPC`** - args: int

### SPHERE

- **`SPHERE`** - args: float -> returns handle — Blitz-style static sphere entity — ENTITYREF; optional 2nd arg segments (see ENTITY.CREATESPHERE)
- **`SPHERE`** - args: float, int -> returns handle — Blitz-style static sphere entity — ENTITYREF handle

### SPHERECOLLIDE

- **`SPHERECOLLIDE`** - args: handle, float -> returns void — Easy Mode: Set entity to use sphere collision with given radius

### SPLIT

- **`SPLIT`** - args: string, string -> returns handle

### SPLIT$

- **`SPLIT$`** - args: string, string -> returns handle

### SPRITE

- **`SPRITE.DEFANIM`** - args: handle, string
- **`SPRITE.DRAW`** - args: handle, int, int
- **`SPRITE.FREE`** - args: handle
- **`SPRITE.GETALPHA`** - args: handle -> returns float
- **`SPRITE.GETCOLOR`** - args: handle -> returns array — RGBA as floats (A channel 0–255)
- **`SPRITE.GETPOS`** - args: handle -> returns array
- **`SPRITE.GETROT`** - args: handle -> returns handle — Returns [0, 0, roll] radians (2D screen rotation)
- **`SPRITE.GETSCALE`** - args: handle -> returns handle — Returns [sx, sy, 1] scale factors (2D draw uses DrawTexturePro)
- **`SPRITE.HIT`** - args: handle, handle — True if the two sprites' drawn quads overlap (same scale, origin, and rotation as SPRITE.DRAW / DrawTexturePro; SAT on quad corners).
- **`SPRITE.LOAD`** - args: string
- **`SPRITE.PLAY`** - args: handle, int, int, float, bool — Animate frames start..end at speed (frames/sec); call SPRITE.UPDATEANIM with Time.Delta()
- **`SPRITE.PLAYANIM`** - args: handle, string
- **`SPRITE.POINTHIT`** - args: handle, float, float — True if (x,y) lies inside the sprite's drawn quad (same space as SPRITE.DRAW position plus SETPOS offsets).
- **`SPRITE.SETALPHA`** - args: handle, float
- **`SPRITE.SETCOLOR`** - args: handle, int, int, int
- **`SPRITE.SETCOLOR`** - args: handle, int, int, int, float
- **`SPRITE.SETFRAME`** - args: handle, int — Manual frame index (strip / DEFANIM); stops SPRITE.PLAY range playback
- **`SPRITE.SETORIGIN`** - args: handle, float, float — Pivot offset in pixels (subtracted from draw position)
- **`SPRITE.SETPOS`** - args: handle, float, float
- **`SPRITE.SETPOSITION`** - args: handle, float, float — DEPRECATED alias of SPRITE.SETPOS. Use SPRITE.SETPOS.
- **`SPRITE.SETROT`** - args: handle, float — Sets rotation in radians (CCW)
- **`SPRITE.SETSCALE`** - args: handle, float, float
- **`SPRITE.UPDATEANIM`** - args: handle, float

### SPRITEBATCH

- **`SPRITEBATCH.ADD`** - args: handle, handle, int, int
- **`SPRITEBATCH.CLEAR`** - args: handle
- **`SPRITEBATCH.CREATE`** - args: (none) -> returns handle
- **`SPRITEBATCH.DRAW`** - args: handle
- **`SPRITEBATCH.FREE`** - args: handle
- **`SPRITEBATCH.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of SPRITEBATCH.CREATE. Use SPRITEBATCH.CREATE.

### SPRITECOLLIDE

- **`SPRITECOLLIDE`** - args: handle, handle — Alias of SPRITE.HIT.

### SPRITEGROUP

- **`SPRITEGROUP.ADD`** - args: handle, handle
- **`SPRITEGROUP.CLEAR`** - args: handle
- **`SPRITEGROUP.CREATE`** - args: (none) -> returns handle
- **`SPRITEGROUP.DRAW`** - args: handle, int, int
- **`SPRITEGROUP.FREE`** - args: handle
- **`SPRITEGROUP.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of SPRITEGROUP.CREATE. Use SPRITEGROUP.CREATE.
- **`SPRITEGROUP.REMOVE`** - args: handle, handle

### SPRITELAYER

- **`SPRITELAYER.ADD`** - args: handle, handle
- **`SPRITELAYER.CLEAR`** - args: handle
- **`SPRITELAYER.CREATE`** - args: float -> returns handle
- **`SPRITELAYER.DRAW`** - args: handle, int, int
- **`SPRITELAYER.FREE`** - args: handle
- **`SPRITELAYER.MAKE`** - args: float -> returns handle — DEPRECATED alias of SPRITELAYER.CREATE. Use SPRITELAYER.CREATE.
- **`SPRITELAYER.SETZ`** - args: handle, float

### SPRITEMODE

- **`SPRITEMODE`** - args: handle, int -> returns void — Easy Mode: Set sprite billboard/blend mode

### SPRITEUI

- **`SPRITEUI.CREATE`** - args: handle, float, float -> returns handle
- **`SPRITEUI.DRAW`** - args: handle, int, int
- **`SPRITEUI.FREE`** - args: handle
- **`SPRITEUI.MAKE`** - args: handle, float, float -> returns handle — DEPRECATED alias of SPRITEUI.CREATE. Use SPRITEUI.CREATE.

### SPRITEVIEWMODE

- **`SPRITEVIEWMODE`** - args: handle, int -> returns void — Alias of SPRITEMODE: 1=Y billboard, 2=full billboard, 3=static quad

### SQR

- **`SQR`** - args: any

### SQRT

- **`SQRT`** - args: any

### STARTSWITH

- **`STARTSWITH`** - args: string, string -> returns bool

### STATIC

- **`STATIC.CREATE`** - args: handle -> returns handle — Creates a Static Body (environment) from a shape handle.
- **`STATIC.MAKE`** - args: handle -> returns handle — DEPRECATED alias of STATIC.CREATE. Use STATIC.CREATE. Creates a Static Body (environment) from a shape handle.

### STEER

- **`STEER.ARRIVE`** - args: handle, float, float, float, float -> returns handle
- **`STEER.AVOIDOBSTACLES`** - args: handle, float -> returns handle
- **`STEER.FLEE`** - args: handle, float, float, float -> returns handle
- **`STEER.FLOCK`** - args: handle, handle, float, float, float -> returns handle
- **`STEER.FOLLOWPATH`** - args: handle, handle -> returns handle
- **`STEER.GROUPADD`** - args: handle, handle
- **`STEER.GROUPCLEAR`** - args: handle
- **`STEER.GROUPMAKE`** - args: (none) -> returns handle
- **`STEER.SEEK`** - args: handle, float, float, float -> returns handle
- **`STEER.WANDER`** - args: handle, float, float -> returns handle

### STOP

- **`STOP`** - args: (none)

### STOPMUSIC

- **`STOPMUSIC`** - args: handle — Easy Mode: AUDIO.STOP(music)

### STOPWATCH

- **`STOPWATCH.ELAPSED`** - args: handle -> returns float
- **`STOPWATCH.FREE`** - args: handle
- **`STOPWATCH.NEW`** - args: (none) -> returns handle
- **`STOPWATCH.RESET`** - args: handle

### STR

- **`STR`** - args: any -> returns string — Convert a value to string (canonical; same as legacy STR$).

### STR$

- **`STR$`** - args: any -> returns string — DEPRECATED alias of STR. Use STR(...).

### STRING

- **`STRING`** - args: int, string -> returns string
- **`STRING.INTERP`** - args: string, any -> returns string
- **`STRING.INTERP`** - args: string, any, any -> returns string
- **`STRING.INTERP`** - args: string, any, any, any -> returns string
- **`STRING.INTERP`** - args: string, any, any, any, any -> returns string
- **`STRING.INTERP`** - args: string, any, any, any, any, any -> returns string
- **`STRING.INTERP`** - args: string, any, any, any, any, any, any -> returns string
- **`STRING.INTERP`** - args: string, any, any, any, any, any, any, any -> returns string
- **`STRING.INTERP`** - args: string, any, any, any, any, any, any, any, any -> returns string
- **`STRING.INTERP`** - args: string, any, any, any, any, any, any, any, any, any -> returns string
- **`STRING.INTERP`** - args: string, any, any, any, any, any, any, any, any, any, any -> returns string
- **`STRING.INTERP$`** - args: string, any -> returns string
- **`STRING.INTERP$`** - args: string, any, any -> returns string
- **`STRING.INTERP$`** - args: string, any, any, any -> returns string
- **`STRING.INTERP$`** - args: string, any, any, any, any -> returns string
- **`STRING.INTERP$`** - args: string, any, any, any, any, any -> returns string
- **`STRING.INTERP$`** - args: string, any, any, any, any, any, any -> returns string
- **`STRING.INTERP$`** - args: string, any, any, any, any, any, any, any -> returns string
- **`STRING.INTERP$`** - args: string, any, any, any, any, any, any, any, any -> returns string
- **`STRING.INTERP$`** - args: string, any, any, any, any, any, any, any, any, any -> returns string
- **`STRING.INTERP$`** - args: string, any, any, any, any, any, any, any, any, any, any -> returns string

### STRING$

- **`STRING$`** - args: int, string -> returns string

### SWITCH

- **`SWITCH`** - args: any, any, any, any
- **`SWITCH`** - args: any, any, any, any, any, any
- **`SWITCH`** - args: any, any, any, any, any, any, any, any
- **`SWITCH`** - args: any, any, any, any, any, any, any, any, any, any
- **`SWITCH`** - args: any, any, any, any, any, any, any, any, any, any, any, any
- **`SWITCH`** - args: any, any, any, any, any, any, any, any, any, any, any, any, any, any
- **`SWITCH`** - args: any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any
- **`SWITCH`** - args: any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any, any

### SYSTEM

- **`SYSTEM.CPUNAME`** - args: (none) -> returns string
- **`SYSTEM.EXECUTE`** - args: string -> returns int
- **`SYSTEM.EXIT`** - args: (none)
- **`SYSTEM.FREEMEMORY`** - args: (none) -> returns int
- **`SYSTEM.GETCLIPBOARD`** - args: (none) -> returns string
- **`SYSTEM.GETENV`** - args: string -> returns string
- **`SYSTEM.GPUNAME`** - args: (none) -> returns string
- **`SYSTEM.ISDEBUGBUILD`** - args: (none) -> returns bool
- **`SYSTEM.LOCALE`** - args: (none) -> returns string
- **`SYSTEM.OPENURL`** - args: string
- **`SYSTEM.SETCLIPBOARD`** - args: string
- **`SYSTEM.SETENV`** - args: string, string
- **`SYSTEM.TOTALMEMORY`** - args: (none) -> returns int
- **`SYSTEM.USERNAME`** - args: (none) -> returns string
- **`SYSTEM.VERSION`** - args: (none) -> returns string — MoonBasic release string (e.g. 1.0.0-GOLD); informational only.

### SetAnimTime

- **`SetAnimTime`** - args: int, float

### SetMSAA

- **`SetMSAA`** - args: int — Alias of WINDOW.SETMSAA — MSAA sample hint; 2+ enables GPU MSAA hint

### SetPostProcess

- **`SetPostProcess`** - args: handle — Full-screen post shader (alias of POST.ADDSHADER)

### SetSSAO

- **`SetSSAO`** - args: bool — Screen-space ambient occlusion (alias of EFFECT.SSAO enable)

### SoundPitch

- **`SoundPitch`** - args: handle, float — Alias of AUDIO.SETSOUNDPITCH

### SoundVolume

- **`SoundVolume`** - args: handle, float — Alias of AUDIO.SETSOUNDVOLUME

### TAB

- **`TAB`** - args: int

### TABLE

- **`TABLE.ADDROW`** - args: handle
- **`TABLE.COLCOUNT`** - args: handle -> returns int
- **`TABLE.CREATE`** - args: string -> returns handle
- **`TABLE.FREE`** - args: handle
- **`TABLE.FROMCSV`** - args: handle -> returns handle
- **`TABLE.FROMJSON`** - args: handle -> returns handle
- **`TABLE.GET`** - args: handle, int, string
- **`TABLE.MAKE`** - args: string -> returns handle — DEPRECATED alias of TABLE.CREATE. Use TABLE.CREATE.
- **`TABLE.ROWCOUNT`** - args: handle -> returns int
- **`TABLE.SET`** - args: handle, int, string, any
- **`TABLE.TOCSV`** - args: handle -> returns handle
- **`TABLE.TOJSON`** - args: handle -> returns handle

### TAN

- **`TAN`** - args: any

### TAND

- **`TAND`** - args: any

### TAU

- **`TAU`** - args: (none)

### TERRAIN

- **`TERRAIN.APPLYMAP`** - args: handle, handle — Apply CPU image as terrain diffuse + splat sample; rebuilds loaded chunk meshes
- **`TERRAIN.APPLYTILES`** - args: handle, handle, int -> returns int — Copy template entity to each non-empty tile on layer 0; returns count placed
- **`TERRAIN.APPLYTILES`** - args: handle, handle, int, int -> returns int — Same as 3-arg form with explicit tile layer index
- **`TERRAIN.CREATE`** - args: int, int
- **`TERRAIN.CREATE`** - args: int, int, float -> returns handle
- **`TERRAIN.DRAW`** - args: handle
- **`TERRAIN.FILLFLAT`** - args: handle, float
- **`TERRAIN.FILLPERLIN`** - args: handle, float, float
- **`TERRAIN.FREE`** - args: handle
- **`TERRAIN.GETHEIGHT`** - args: handle, float, float -> returns float
- **`TERRAIN.GETNORMAL`** - args: handle, float, float -> returns handle — Unit terrain normal (heap vec3) for slope tilt
- **`TERRAIN.GETSLOPE`** - args: handle, float, float -> returns float
- **`TERRAIN.GETSPLAT`** - args: handle, float, float -> returns int — Diffuse/splat map red channel 0..255 (-1 if no map); use for footstep ids
- **`TERRAIN.LOAD`** - args: string, string -> returns handle — Heightmap image path + optional diffuse/splat path; GPU mesh + CPU splat sample
- **`TERRAIN.LOWER`** - args: handle, float, float, float, float
- **`TERRAIN.MAKE`** - args: int, int — DEPRECATED alias of TERRAIN.CREATE. Use TERRAIN.CREATE.
- **`TERRAIN.MAKE`** - args: int, int, float -> returns handle — DEPRECATED alias of TERRAIN.CREATE. Use TERRAIN.CREATE.
- **`TERRAIN.PLACE`** - args: handle, int, float, float, float
- **`TERRAIN.RAISE`** - args: handle, float, float, float, float
- **`TERRAIN.RAYCAST`** - args: handle, float, float, float, float, float, float -> returns handle — Ray vs terrain only; float array [hit, x, y, z]; max ray length is large by default
- **`TERRAIN.SETASYNCMESHBUILD`** - args: handle, bool — When true, CPU heightmap prep runs on a background goroutine; GenMeshHeightmap still runs on the main thread when jobs drain (use with WINDOW.SETLOADINGMODE / mesh budget).
- **`TERRAIN.SETCHUNKSIZE`** - args: handle, int
- **`TERRAIN.SETDETAIL`** - args: handle, float — LOD factor in (0,1]: lower = coarser chunk meshes
- **`TERRAIN.SETMESHBUILDBUDGET`** - args: handle, int — Max chunk mesh GPU rebuilds per WORLD.UPDATE tick; 0 = unlimited (default). Use 1–4 to avoid UI thread stalls.
- **`TERRAIN.SETPOS`** - args: handle, float, float, float
- **`TERRAIN.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of TERRAIN.SETPOS. Use TERRAIN.SETPOS.
- **`TERRAIN.SETSCALE`** - args: handle, float, float, float — Non-uniform scale: XZ stretch per cell, Y height multiplier (marks chunks dirty)
- **`TERRAIN.SNAPY`** - args: handle, int, float

### TEXTDRAW

- **`TEXTDRAW.DRAW`** - args: handle

### TEXTOBJ

- **`TEXTOBJ`** - args: string -> returns handle

### TEXTURE

- **`TEXTURE.FREE`** - args: handle
- **`TEXTURE.FROMIMAGE`** - args: handle
- **`TEXTURE.GENCHECKED`** - args: int, int, int, int, handle, handle -> returns handle
- **`TEXTURE.GENCOLOR`** - args: int, int, int, int, int, int -> returns handle
- **`TEXTURE.GENGRADIENTH`** - args: int, int, handle, handle -> returns handle
- **`TEXTURE.GENGRADIENTV`** - args: int, int, handle, handle -> returns handle
- **`TEXTURE.GENWHITENOISE`** - args: int, int -> returns handle
- **`TEXTURE.GENWHITENOISE`** - args: int, int, float -> returns handle
- **`TEXTURE.HEIGHT`** - args: handle -> returns int
- **`TEXTURE.LOAD`** - args: string
- **`TEXTURE.LOADANIM`** - args: string, int, int -> returns handle — TEXTURE.LOAD + SETGRID in one call
- **`TEXTURE.PLAY`** - args: handle, float, bool — Auto-advance atlas frames; call TEXTURE.TICKALL each frame
- **`TEXTURE.RELOAD`** - args: handle
- **`TEXTURE.SETDEFAULTFILTER`** - args: int
- **`TEXTURE.SETDISTORTION`** - args: handle, float — Shader-side distortion amount hint
- **`TEXTURE.SETFILTER`** - args: handle, int
- **`TEXTURE.SETFRAME`** - args: handle, int — Select atlas frame index (0-based)
- **`TEXTURE.SETGRID`** - args: handle, int, int — Spritesheet layout: columns x rows of equal frames
- **`TEXTURE.SETUVSCROLL`** - args: handle, float, float — Source-rectangle scroll speeds for sampled UVs
- **`TEXTURE.SETWRAP`** - args: handle, int
- **`TEXTURE.STOPANIM`** - args: handle
- **`TEXTURE.TICKALL`** - args: (none) — Advance all playing atlas animations (optional dt via overload)
- **`TEXTURE.TICKALL`** - args: float
- **`TEXTURE.UPDATE`** - args: handle, handle
- **`TEXTURE.WIDTH`** - args: handle -> returns int

### TEXTUREHEIGHT

- **`TEXTUREHEIGHT`** - args: handle -> returns int

### TEXTUREWIDTH

- **`TEXTUREWIDTH`** - args: handle -> returns int

### TFormVector

- **`TFormVector`** - args: float, float, float, int, int -> returns handle — Args: (x#, y#, z#, srcEntity#, dstEntity#). Blitz alias of ENTITY.TFORMVECTOR: direction in src local space as linear transform into dst local space; returns heap handle to 3 float components (no world-entity shortcut).

### THROW

- **`THROW`** - args: int, string

### TICKCOUNT

- **`TICKCOUNT`** - args: (none)
- **`TICKCOUNT`** - args: (none) -> returns int

### TILEMAP

- **`TILEMAP.COLLISIONAT`** - args: handle, int, int -> returns int
- **`TILEMAP.DRAW`** - args: handle
- **`TILEMAP.DRAWLAYER`** - args: handle, int
- **`TILEMAP.FREE`** - args: handle
- **`TILEMAP.GETTILE`** - args: handle, int, int, int -> returns int
- **`TILEMAP.HEIGHT`** - args: handle -> returns int
- **`TILEMAP.ISSOLID`** - args: handle, int, int -> returns bool
- **`TILEMAP.ISSOLIDCATEGORY`** - args: handle, int, int, int -> returns bool
- **`TILEMAP.LAYERCOUNT`** - args: handle -> returns int
- **`TILEMAP.LAYERNAME`** - args: handle, int -> returns string
- **`TILEMAP.LOAD`** - args: string -> returns handle
- **`TILEMAP.MERGECOLLISIONLAYER`** - args: handle, int, int
- **`TILEMAP.SETCOLLISION`** - args: handle, int, int, int
- **`TILEMAP.SETTILE`** - args: handle, int, int, int, int
- **`TILEMAP.SETTILESIZE`** - args: handle, int, int
- **`TILEMAP.WIDTH`** - args: handle -> returns int

### TIME

- **`TIME`** - args: (none)
- **`TIME`** - args: (none) -> returns string
- **`TIME.DELTA`** - args: (none)
- **`TIME.DELTA`** - args: (none) -> returns float
- **`TIME.DELTA`** - args: float, float -> returns float
- **`TIME.GET`** - args: (none)
- **`TIME.GET`** - args: (none) -> returns float
- **`TIME.GETFPS`** - args: (none) -> returns float
- **`TIME.SETMAXDELTA`** - args: float

### TIME$

- **`TIME$`** - args: (none)
- **`TIME$`** - args: (none) -> returns string

### TIMER

- **`TIMER`** - args: (none) -> returns float
- **`TIMER`** - args: (none)
- **`TIMER.CREATE`** - args: float -> returns handle — Simulation timer: duration in seconds; use TIMER.START/UPDATE with delta time
- **`TIMER.DONE`** - args: handle -> returns bool
- **`TIMER.FINISHED`** - args: handle -> returns bool
- **`TIMER.FRACTION`** - args: handle -> returns float
- **`TIMER.FREE`** - args: handle
- **`TIMER.MAKE`** - args: float -> returns handle — DEPRECATED alias of TIMER.CREATE. Use TIMER.CREATE.
- **`TIMER.NEW`** - args: float -> returns handle — Wall-clock deadline timer (time.Now-based)
- **`TIMER.REMAINING`** - args: handle -> returns float
- **`TIMER.RESET`** - args: handle, float
- **`TIMER.REWIND`** - args: handle
- **`TIMER.SETLOOP`** - args: handle, any
- **`TIMER.START`** - args: handle
- **`TIMER.STOP`** - args: handle
- **`TIMER.UPDATE`** - args: handle, float

### TIMESTAMP

- **`TIMESTAMP`** - args: (none)
- **`TIMESTAMP`** - args: (none) -> returns float

### TRACE

- **`TRACE`** - args: any

### TRANSFORM

- **`TRANSFORM.APPLYX`** - args: handle, float, float, float -> returns float
- **`TRANSFORM.APPLYY`** - args: handle, float, float, float -> returns float
- **`TRANSFORM.APPLYZ`** - args: handle, float, float, float -> returns float
- **`TRANSFORM.FREE`** - args: handle
- **`TRANSFORM.GETELEMENT`** - args: handle, int, int -> returns float
- **`TRANSFORM.IDENTITY`** - args: (none) -> returns handle
- **`TRANSFORM.INVERSE`** - args: handle -> returns handle
- **`TRANSFORM.LOOKAT`** - args: float, float, float, float, float, float, float, float, float -> returns handle
- **`TRANSFORM.MULTIPLY`** - args: handle, handle -> returns handle
- **`TRANSFORM.ORTHO`** - args: float, float, float, float, float, float -> returns handle
- **`TRANSFORM.PERSPECTIVE`** - args: float, float, float, float -> returns handle
- **`TRANSFORM.ROTATION`** - args: float, float, float -> returns handle
- **`TRANSFORM.SCALE`** - args: float, float, float -> returns handle
- **`TRANSFORM.SETROTATION`** - args: handle, float, float, float
- **`TRANSFORM.TRANSLATION`** - args: float, float, float -> returns handle
- **`TRANSFORM.TRANSPOSE`** - args: handle -> returns handle

### TRANSITION

- **`TRANSITION.FADEIN`** - args: float
- **`TRANSITION.FADEOUT`** - args: float
- **`TRANSITION.ISDONE`** - args: (none) -> returns bool
- **`TRANSITION.SETCOLOR`** - args: int, int, int, int
- **`TRANSITION.WIPE`** - args: string, float

### TRIGGER

- **`TRIGGER.CREATE`** - args: handle -> returns handle — Creates a Trigger Body (non-solid sensor) from a shape handle.
- **`TRIGGER.CREATEFROMENTITY`** - args: int
- **`TRIGGER.CREATEZONE`** - args: float, float, float, float, float, float, string -> returns handle — Creates non-blocking zone firing hit tags.
- **`TRIGGER.MAKE`** - args: handle -> returns handle — DEPRECATED alias of TRIGGER.CREATE. Use TRIGGER.CREATE. Creates a Trigger Body (non-solid sensor) from a shape handle.
- **`TRIGGER.MAKEFROMENTITY`** - args: int — DEPRECATED alias of TRIGGER.CREATEFROMENTITY. Use TRIGGER.CREATEFROMENTITY.
- **`TRIGGER.MAKEZONE`** - args: float, float, float, float, float, float, string -> returns handle — DEPRECATED alias of TRIGGER.CREATEZONE. Use TRIGGER.CREATEZONE. Creates non-blocking zone firing hit tags.

### TRIM

- **`TRIM`** - args: string -> returns string

### TRIM$

- **`TRIM$`** - args: string -> returns string

### TURNCAMERA

- **`TURNCAMERA`** - args: handle, float, float, float — Easy Mode: CAMERA.TURN(cam, p, y, r)

### TURNENTITY

- **`TURNENTITY`** - args: handle, float, float, float -> returns void — Easy Mode: Incremental rotation

### TWEEN

- **`TWEEN.CREATE`** - args: (none) -> returns handle
- **`TWEEN.LOOP`** - args: handle, int
- **`TWEEN.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of TWEEN.CREATE. Use TWEEN.CREATE.
- **`TWEEN.ONCOMPLETE`** - args: handle, string
- **`TWEEN.START`** - args: handle
- **`TWEEN.STOP`** - args: handle
- **`TWEEN.THEN`** - args: handle, string, float, float, string
- **`TWEEN.TO`** - args: handle, string, float, float, string
- **`TWEEN.UPDATE`** - args: handle, float
- **`TWEEN.YOYO`** - args: handle

### TYPEOF

- **`TYPEOF`** - args: any

### TranslateEntity

- **`TranslateEntity`** - args: int, float, float, float — Args: (entity#, dx#, dy#, dz#). World-space position delta; same as ENTITY.TRANSLATE / ENTITY.TRANSLATEENTITY.

### UI

- **`UI.BUTTON`** - args: string, float, float, float, float -> returns int — Draws interactive button.
- **`UI.INVENTORYICON`** - args: handle, float, float — Draws a 3D model icon.
- **`UI.LABEL3D`** - args: string, handle, handle — Projects text into 3D world above object.
- **`UI.PROGRESSBAR`** - args: float, float, float, float, float, int — Draws progress bar.

### UPDATEEMITTER

- **`UPDATEEMITTER`** - args: handle, float -> returns void — Easy Mode: Update emitter simulation

### UPDATEPHYSICS

- **`UPDATEPHYSICS`** - args: (none) -> returns void — Easy Mode: ENTITY.UPDATE(Time.Delta); best-effort WORLD.UPDATE, PHYSICS2D.STEP, PHYSICS3D.UPDATE (same as PHYSICS3D.STEP)

### UPDW

- **`UPDW`** - args: float — Shorthand: ENTITY.UPDATE(dt) — use ENTITY.UPDATE in scripts

### UPPER

- **`UPPER`** - args: string -> returns string

### UPPER$

- **`UPPER$`** - args: string -> returns string

### UTIL

- **`UTIL.CHANGEDIR`** - args: string -> returns bool
- **`UTIL.CLEARDROPPEDFILES`** - args: (none)
- **`UTIL.COPYFILE`** - args: string, string
- **`UTIL.CREATEDIRECTORY`** - args: string -> returns bool
- **`UTIL.DELETEDIR`** - args: string
- **`UTIL.DELETEFILE`** - args: string -> returns bool
- **`UTIL.FILEEXISTS`** - args: string -> returns bool
- **`UTIL.GETDIR`** - args: (none) -> returns string
- **`UTIL.GETDIRFILES`** - args: string -> returns string
- **`UTIL.GETDIRS`** - args: string -> returns string
- **`UTIL.GETDROPPEDFILES`** - args: (none) -> returns string
- **`UTIL.GETFILEEXT`** - args: string -> returns string
- **`UTIL.GETFILEMODTIME`** - args: string -> returns int
- **`UTIL.GETFILENAME`** - args: string -> returns string
- **`UTIL.GETFILENAMENOEXT`** - args: string -> returns string
- **`UTIL.GETFILEPATH`** - args: string -> returns string
- **`UTIL.GETFILESIZE`** - args: string -> returns int
- **`UTIL.ISDIR`** - args: string -> returns bool
- **`UTIL.ISFILEDROPPED`** - args: (none) -> returns bool
- **`UTIL.ISFILENAMEVALID`** - args: string -> returns bool
- **`UTIL.LOADTEXT`** - args: string -> returns string
- **`UTIL.MAKEDIRECTORY`** - args: string -> returns bool — DEPRECATED alias of UTIL.CREATEDIRECTORY. Use UTIL.CREATEDIRECTORY.
- **`UTIL.MOVEFILE`** - args: string, string
- **`UTIL.RENAMEFILE`** - args: string, string
- **`UTIL.SAVETEXT`** - args: string, string

### UpdateMesh

- **`UpdateMesh`** - args: int

### VAL

- **`VAL`** - args: string -> returns float

### VEC2

- **`VEC2.ADD`** - args: handle, handle -> returns handle
- **`VEC2.ANGLE`** - args: handle, handle -> returns float
- **`VEC2.CREATE`** - args: float, float -> returns handle
- **`VEC2.DIST`** - args: handle, handle -> returns float
- **`VEC2.DIST`** - args: float, float, float, float -> returns float
- **`VEC2.DISTANCE`** - args: handle, handle -> returns float
- **`VEC2.DISTSQ`** - args: float, float, float, float -> returns float
- **`VEC2.FREE`** - args: handle
- **`VEC2.LENGTH`** - args: handle -> returns float
- **`VEC2.LENGTH`** - args: float, float -> returns float
- **`VEC2.LERP`** - args: handle, handle, float -> returns handle
- **`VEC2.MAKE`** - args: float, float -> returns handle — DEPRECATED alias of VEC2.CREATE. Use VEC2.CREATE.
- **`VEC2.MOVE_TOWARD`** - args: float, float, float, float, float -> returns handle
- **`VEC2.MUL`** - args: handle, float -> returns handle
- **`VEC2.NORMALIZE`** - args: handle -> returns handle
- **`VEC2.NORMALIZE`** - args: float, float -> returns handle
- **`VEC2.PUSHOUT`** - args: float, float, float, float, float -> returns handle
- **`VEC2.ROTATE`** - args: handle, float -> returns handle
- **`VEC2.SET`** - args: handle, float, float
- **`VEC2.SUB`** - args: handle, handle -> returns handle
- **`VEC2.TRANSFORMMAT4`** - args: handle, handle -> returns handle
- **`VEC2.X`** - args: handle -> returns float
- **`VEC2.Y`** - args: handle -> returns float

### VEC3

- **`VEC3.ADD`** - args: handle, handle -> returns handle
- **`VEC3.ANGLE`** - args: handle, handle -> returns float
- **`VEC3.CREATE`** - args: float, float, float -> returns handle
- **`VEC3.CROSS`** - args: handle, handle -> returns handle
- **`VEC3.DIST`** - args: handle, handle -> returns float
- **`VEC3.DIST`** - args: float, float, float, float, float, float -> returns float
- **`VEC3.DISTANCE`** - args: handle, handle -> returns float
- **`VEC3.DISTSQ`** - args: float, float, float, float, float, float -> returns float
- **`VEC3.DIV`** - args: handle, float -> returns handle
- **`VEC3.DOT`** - args: handle, handle -> returns float
- **`VEC3.EQUALS`** - args: handle, handle -> returns bool
- **`VEC3.FREE`** - args: handle
- **`VEC3.LENGTH`** - args: handle -> returns float
- **`VEC3.LENGTH`** - args: float, float, float -> returns float
- **`VEC3.LERP`** - args: handle, handle, float -> returns handle
- **`VEC3.MAKE`** - args: float, float, float -> returns handle — DEPRECATED alias of VEC3.CREATE. Use VEC3.CREATE.
- **`VEC3.MUL`** - args: handle, float -> returns handle
- **`VEC3.NEGATE`** - args: handle -> returns handle
- **`VEC3.NORMALIZE`** - args: handle -> returns handle
- **`VEC3.NORMALIZE`** - args: float, float, float -> returns handle
- **`VEC3.ORTHONORMALIZE`** - args: handle, handle
- **`VEC3.PROJECT`** - args: handle, handle -> returns handle
- **`VEC3.REFLECT`** - args: handle, handle -> returns handle
- **`VEC3.ROTATEBYQUAT`** - args: handle, handle -> returns handle
- **`VEC3.SET`** - args: handle, float, float, float
- **`VEC3.SUB`** - args: handle, handle -> returns handle
- **`VEC3.TRANSFORMMAT4`** - args: handle, handle -> returns handle
- **`VEC3.VEC3`** - args: float, float, float -> returns handle
- **`VEC3.VECADD`** - args: handle, handle -> returns handle
- **`VEC3.VECCROSS`** - args: handle, handle -> returns handle
- **`VEC3.VECDOT`** - args: handle, handle -> returns float
- **`VEC3.VECLENGTH`** - args: handle -> returns float
- **`VEC3.VECNORMALIZE`** - args: handle -> returns handle
- **`VEC3.VECSCALE`** - args: handle, float -> returns handle
- **`VEC3.VECSUB`** - args: handle, handle -> returns handle
- **`VEC3.X`** - args: handle -> returns float
- **`VEC3.Y`** - args: handle -> returns float
- **`VEC3.Z`** - args: handle -> returns float

### VertexX

- **`VertexX`** - args: handle, int -> returns float

### VertexY

- **`VertexY`** - args: handle, int -> returns float

### VertexZ

- **`VertexZ`** - args: handle, int -> returns float

### WAIT

- **`WAIT`** - args: any

### WATER

- **`WATER.CREATE`** - args: float, float, float, float, float -> returns handle — x, z, width, depth, water level (Y); same plane as WATER.MAKE
- **`WATER.CREATE`** - args: float, int, int, int, int -> returns handle
- **`WATER.DRAW`** - args: handle
- **`WATER.FREE`** - args: handle
- **`WATER.GETDEPTH`** - args: handle, float, float -> returns float
- **`WATER.GETWAVEY`** - args: handle, float, float -> returns float
- **`WATER.ISUNDER`** - args: handle, float, float, float -> returns bool
- **`WATER.MAKE`** - args: float, int, int, int, int -> returns handle — DEPRECATED alias of WATER.CREATE. Use WATER.CREATE.
- **`WATER.MAKE`** - args: float, float, float, float, float -> returns handle — DEPRECATED alias of WATER.CREATE. Use WATER.CREATE. x, z, width, depth, water level (Y); same plane as WATER.MAKE
- **`WATER.SETCOLOR`** - args: handle, int, float — Packed RGB diffuse (0xRRGGBB) and clarity (0..1 alpha, or 0..255); updates shallow/deep tint
- **`WATER.SETDEEPCOLOR`** - args: handle, int, int, int, int
- **`WATER.SETHEIGHT`** - args: handle, float
- **`WATER.SETPOS`** - args: handle, float, float, float
- **`WATER.SETPOSITION`** - args: handle, float, float, float — DEPRECATED alias of WATER.SETPOS. Use WATER.SETPOS.
- **`WATER.SETSHALLOWCOLOR`** - args: handle, int, int, int, int
- **`WATER.SETWAVE`** - args: handle, float, float — Sets wave frequency (speed) and amplitude (height)
- **`WATER.SETWAVEHEIGHT`** - args: handle, float
- **`WATER.SHOW`** - args: handle, bool
- **`WATER.UPDATE`** - args: handle, float

### WAVE

- **`WAVE.COPY`** - args: handle -> returns handle
- **`WAVE.CROP`** - args: handle, int, int
- **`WAVE.EXPORT`** - args: handle, string
- **`WAVE.FORMAT`** - args: handle, int, int, int
- **`WAVE.FREE`** - args: handle
- **`WAVE.LOAD`** - args: string -> returns handle

### WEATHER

- **`WEATHER.CREATE`** - args: (none) -> returns handle
- **`WEATHER.DRAW`** - args: handle
- **`WEATHER.FREE`** - args: handle
- **`WEATHER.GETCOVERAGE`** - args: handle -> returns float
- **`WEATHER.GETTYPE`** - args: handle -> returns string
- **`WEATHER.MAKE`** - args: (none) -> returns handle — DEPRECATED alias of WEATHER.CREATE. Use WEATHER.CREATE.
- **`WEATHER.SETTYPE`** - args: handle, string
- **`WEATHER.UPDATE`** - args: handle, float

### WIND

- **`WIND.GETSTRENGTH`** - args: (none) -> returns float
- **`WIND.SET`** - args: float, float, float

### WINDOW

- **`WINDOW.CANOPEN`** - args: int, int, string -> returns bool
- **`WINDOW.CHECKFLAG`** - args: int -> returns bool
- **`WINDOW.CLEARFLAG`** - args: int
- **`WINDOW.CLOSE`** - args: (none)
- **`WINDOW.GETFPS`** - args: (none) -> returns int
- **`WINDOW.GETMONITORCOUNT`** - args: (none) -> returns int
- **`WINDOW.GETMONITORHEIGHT`** - args: int -> returns int
- **`WINDOW.GETMONITORNAME`** - args: int -> returns string
- **`WINDOW.GETMONITORREFRESHRATE`** - args: int -> returns int
- **`WINDOW.GETMONITORWIDTH`** - args: int -> returns int
- **`WINDOW.GETPOSITIONX`** - args: (none) -> returns int
- **`WINDOW.GETPOSITIONY`** - args: (none) -> returns int
- **`WINDOW.GETSCALEDPIX`** - args: (none) -> returns float
- **`WINDOW.GETSCALEDPIY`** - args: (none) -> returns float
- **`WINDOW.HEIGHT`** - args: (none) -> returns int
- **`WINDOW.ISFULLSCREEN`** - args: (none) -> returns bool
- **`WINDOW.ISRESIZED`** - args: (none) -> returns bool
- **`WINDOW.LOADINGMODE`** - args: (none) -> returns bool — Current loading-mode flag from WINDOW.SETLOADINGMODE
- **`WINDOW.MAXIMIZE`** - args: (none)
- **`WINDOW.MINIMIZE`** - args: (none)
- **`WINDOW.OPEN`** - args: int, int, string
- **`WINDOW.RESTORE`** - args: (none)
- **`WINDOW.SETFLAG`** - args: int
- **`WINDOW.SETFPS`** - args: int
- **`WINDOW.SETICON`** - args: string
- **`WINDOW.SETLOADINGMODE`** - args: bool — When true, TERRAIN.DRAW skips drawing so RENDER.FRAME still polls OS events during mesh builds
- **`WINDOW.SETMAXSIZE`** - args: int, int
- **`WINDOW.SETMINSIZE`** - args: int, int
- **`WINDOW.SETMONITOR`** - args: int
- **`WINDOW.SETMSAA`** - args: int — MSAA sample count hint before/during window use (2+ enables GPU MSAA); Easy Mode alias: SetMSAA
- **`WINDOW.SETOPACITY`** - args: float
- **`WINDOW.SETPOS`** - args: int, int
- **`WINDOW.SETPOSITION`** - args: int, int — DEPRECATED alias of WINDOW.SETPOS. Use WINDOW.SETPOS. Deprecated alias of WINDOW.SETPOS — set window client-area position in screen pixels
- **`WINDOW.SETSIZE`** - args: int, int
- **`WINDOW.SETSTATE`** - args: int
- **`WINDOW.SETTARGETFPS`** - args: int
- **`WINDOW.SETTITLE`** - args: string
- **`WINDOW.SHOULDCLOSE`** - args: (none)
- **`WINDOW.TOGGLEFULLSCREEN`** - args: (none)
- **`WINDOW.WIDTH`** - args: (none) -> returns int
- **`WORLD.FLASH`** - args: handle, float — Tints the screen temporarily (damage effects, etc).

### WIRECUBE

- **`WIRECUBE`** - args: float, float, float, float, float, float, int, int, int, int — alias of DRAW3D.CUBEWIRES (Blitz3D WireCube spelling)

### WORLD

- **`WORLD.DAYNIGHTCYCLE`** - args: float — Rotates global sunlight over duration (seconds).
- **`WORLD.EXPLOSION`** - args: float, float, float, float, float — Alias of PHYSICS.EXPLOSION
- **`WORLD.FOGCOLOR`** - args: int, int, int
- **`WORLD.FOGDENSITY`** - args: float
- **`WORLD.FOGMODE`** - args: int
- **`WORLD.GETRAY`** - args: float, float, handle -> returns handle — Returns Array [px,py,pz,dx,dy,dz]
- **`WORLD.GRAVITY`** - args: float, float, float — Alias: forwards to PHYSICS3D.SETGRAVITY (global Jolt gravity)
- **`WORLD.HITSTOP`** - args: float — Freeze gameplay delta for duration (wall-clock seconds) — impact frames
- **`WORLD.ISREADY`** - args: handle -> returns bool
- **`WORLD.MOUSE2D`** - args: handle -> returns handle — Mouse position through Camera2D; float array [wx,wy]
- **`WORLD.MOUSEFLOOR`** - args: handle, float -> returns handle — Alias of WORLD.MOUSEFLOOR3D — mouse ray vs plane y=floorY → [wx,wz] or NIL
- **`WORLD.MOUSEFLOOR3D`** - args: handle, float -> returns handle — Mouse ray vs plane y=floorY; float array [wx,wz] or NIL
- **`WORLD.MOUSEPICK`** - args: handle -> returns int — Alias of WORLD.MOUSETOENTITY — entity id under mouse cursor (physics ray; Linux+Jolt)
- **`WORLD.MOUSETOENTITY`** - args: handle -> returns int — Jolt ray pick at cursor (Linux+CGO); entity# or 0. Same as CAMERA.RAYCASTMOUSE
- **`WORLD.MOUSETOFLOOR`** - args: handle, float -> returns handle — Alias of WORLD.MOUSEFLOOR3D
- **`WORLD.PRELOAD`** - args: handle, int
- **`WORLD.SCREENSHAKE`** - args: float, float — Shakes the primary camera.
- **`WORLD.SETAMBIENCE`** - args: handle, float — Plays a looping background track.
- **`WORLD.SETCENTER`** - args: float, float
- **`WORLD.SETCENTERENTITY`** - args: int
- **`WORLD.SETGRAVITY`** - args: float, float, float — Alias of PHYSICS3D.SETGRAVITY
- **`WORLD.SETREFLECTION`** - args: int
- **`WORLD.SETREVERB`** - args: int — Changes echo.
- **`WORLD.SETTIMESCALE`** - args: float — Alias of GAME.SETTIMESCALE
- **`WORLD.SETVEGETATION`** - args: handle, handle, float — Scatter helper: terrain + billboard entity reserved + density; uses internal SCATTER set
- **`WORLD.SHAKE`** - args: float, float — Alias of WORLD.SCREENSHAKE — screen impact via active camera
- **`WORLD.STATUS`** - args: (none) -> returns string
- **`WORLD.STREAMENABLE`** - args: bool
- **`WORLD.TOSCREEN`** - args: int -> returns handle — WORLD.TOSCREEN(entity#) — screen [x,y] for entity world position via active 3D camera
- **`WORLD.TOSCREEN`** - args: float, float, float -> returns handle — World to screen using active CAMERA.BEGIN 3D camera; returns float array [sx,sy]
- **`WORLD.TOSCREEN`** - args: float, float, float, handle -> returns handle — Returns 2D Screen coords given 3D World coords and Camera.
- **`WORLD.TOWORLD`** - args: float, float, float -> returns handle — Unproject screen x,y with depth along view ray (active 3D camera); returns [wx,wy,wz]
- **`WORLD.TOWORLD`** - args: float, float, float, handle -> returns handle — Returns 3D World coords from 2D.
- **`WORLD.UPDATE`** - args: float

### WRAP

- **`WRAP`** - args: any, any, any

### WRAPANGLE

- **`WRAPANGLE`** - args: any
- **`WRAPANGLE`** - args: float -> returns float — Easy Mode: Wrap angle to 0..360 or -PI..PI range

### WRAPANGLE180

- **`WRAPANGLE180`** - args: any

### WRITE

- **`WRITE`** - args: any — Print values space-separated without newline.

### WRITEALLTEXT

- **`WRITEALLTEXT`** - args: string, string

### WRITEBYTE

- **`WRITEBYTE`** - args: handle, int

### WRITEFILE

- **`WRITEFILE`** - args: handle, string

### WRITEFILELN

- **`WRITEFILELN`** - args: handle, string

### WRITEFLOAT

- **`WRITEFLOAT`** - args: handle, any

### WRITEINT

- **`WRITEINT`** - args: handle, int

### WRITESHORT

- **`WRITESHORT`** - args: handle, int

### WRITESTRING

- **`WRITESTRING`** - args: handle, string

### YAWFROMXZ

- **`YAWFROMXZ`** - args: float, float -> returns float — Yaw radians from flat direction (dx, dz); matches MOVEX/MOVEZ convention

### YEAR

- **`YEAR`** - args: (none)
- **`YEAR`** - args: (none) -> returns int

