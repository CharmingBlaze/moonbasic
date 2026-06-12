# Game systems — complete command registry

> Every registered command for the **40 beginner systems** and their related namespaces.

Generated from `compiler/builtinmanifest/commands.json` (same source as [API_CONSISTENCY.md](../API_CONSISTENCY.md)).

**How to use this page:**

- Learn **why** and **when** in [00-START.md](00-START.md), [GUIDES.md](GUIDES.md) (entity, collision, UI, multiplayer), and `01-CORE` … `11-TOOLING`.
- Look up **arity** and return types here or in [API_CONSISTENCY.md](../API_CONSISTENCY.md).
- Deep behavior: [COMMAND_AUDIT.md](../COMMAND_AUDIT.md) → `docs/reference/*.md`.
- Validate a script: `moonbasic --check yourgame.mb`.
- In-game help: `HELP("ENTITY.SETPOS")`.

**Case:** Command names are **case-insensitive** in source.

---

## Table of contents

- [Core: window, time, render, scene, entity](#core-window-time)
- [Camera and light](#camera-light)
- [Meshes, models, materials, textures, asset packs](#assets)
- [Input and actions](#input-action)
- [Physics, bodies, collision, picking](#physics)
- [Audio (2D and 3D)](#audio)
- [2D sprites, tilemaps, terrain, particles, animation](#2d-world)
- [UI, fonts, and text](#ui-text)
- [Save, files, JSON, math, vectors](#data)
- [Debug, timers](#debug-timer)
- [Globals and language builtins](#globals)
- [All other engine namespaces](#all-other-namespaces)

---

## Core: window, time, render, scene, entity

Guide: [01-CORE.md](01-CORE.md)

### WINDOW

- **`WINDOW.CANOPEN`** — args: int, int, string → bool
- **`WINDOW.CHECKFLAG`** — args: int → bool
- **`WINDOW.CLEARFLAG`** — args: int
- **`WINDOW.CLOSE`** — args: (none)
- **`WINDOW.DPISCALE`** — args: (none) → float
- **`WINDOW.GETFPS`** — args: (none) → int
- **`WINDOW.GETMONITORCOUNT`** — args: (none) → int
- **`WINDOW.GETMONITORHEIGHT`** — args: int → int
- **`WINDOW.GETMONITORNAME`** — args: int → string
- **`WINDOW.GETMONITORREFRESHRATE`** — args: int → int
- **`WINDOW.GETMONITORWIDTH`** — args: int → int
- **`WINDOW.GETPOSITIONX`** — args: (none) → int
- **`WINDOW.GETPOSITIONY`** — args: (none) → int
- **`WINDOW.GETSCALEDPIX`** — args: (none) → float
- **`WINDOW.GETSCALEDPIY`** — args: (none) → float
- **`WINDOW.HEIGHT`** — args: (none) → int
- **`WINDOW.ISFULLSCREEN`** — args: (none) → bool
- **`WINDOW.ISRESIZED`** — args: (none) → bool
- **`WINDOW.LOADINGMODE`** — args: (none) → bool — Current loading-mode flag from WINDOW.SETLOADINGMODE
- **`WINDOW.MAXIMIZE`** — args: (none)
- **`WINDOW.MINIMIZE`** — args: (none)
- **`WINDOW.OPEN`** — args: int, int, string
- **`WINDOW.RESTORE`** — args: (none)
- **`WINDOW.SETFLAG`** — args: int
- **`WINDOW.SETFPS`** — args: int
- **`WINDOW.SETICON`** — args: string
- **`WINDOW.SETLOADINGMODE`** — args: bool — When true, TERRAIN.DRAW skips drawing so RENDER.FRAME still polls OS events during mesh builds
- **`WINDOW.SETMAXSIZE`** — args: int, int
- **`WINDOW.SETMINSIZE`** — args: int, int
- **`WINDOW.SETMONITOR`** — args: int
- **`WINDOW.SETMSAA`** — args: int — MSAA sample count hint before/during window use (2+ enables GPU MSAA); Easy Mode alias: SetMSAA
- **`WINDOW.SETOPACITY`** — args: float
- **`WINDOW.SETPOS`** — args: int, int
- **`WINDOW.SETPOSITION`** — args: int, int — DEPRECATED alias of WINDOW.SETPOS. Use WINDOW.SETPOS. Deprecated alias of WINDOW.SETPOS â€” set window client-area position in screen pixels
- **`WINDOW.SETSIZE`** — args: int, int
- **`WINDOW.SETSTATE`** — args: int
- **`WINDOW.SETTARGETFPS`** — args: int
- **`WINDOW.SETTITLE`** — args: string
- **`WINDOW.SHOULDCLOSE`** — args: (none)
- **`WINDOW.TOGGLEFULLSCREEN`** — args: (none)
- **`WINDOW.WIDTH`** — args: (none) → int
- **`WORLD.FLASH`** — args: handle, float — Tints the screen temporarily (damage effects, etc).

### TIME

- **`TIME`** — args: (none)
- **`TIME`** — args: (none) → string
- **`TIME.DELTA`** — args: (none)
- **`TIME.DELTA`** — args: (none) → float
- **`TIME.DELTA`** — args: float, float → float
- **`TIME.GET`** — args: (none)
- **`TIME.GET`** — args: (none) → float
- **`TIME.GETFPS`** — args: (none) → float
- **`TIME.MILLIS`** — args: (none) → int
- **`TIME.SETMAXDELTA`** — args: float
- **`TIME.UPDATE`** — args: (none)

### SYSTEM

- **`SYSTEM.CPUNAME`** — args: (none) → string
- **`SYSTEM.EXECUTE`** — args: string → int
- **`SYSTEM.EXIT`** — args: (none)
- **`SYSTEM.FREEMEMORY`** — args: (none) → int
- **`SYSTEM.GETCLIPBOARD`** — args: (none) → string
- **`SYSTEM.GETENV`** — args: string → string
- **`SYSTEM.GPUNAME`** — args: (none) → string
- **`SYSTEM.ISDEBUGBUILD`** — args: (none) → bool
- **`SYSTEM.LOCALE`** — args: (none) → string
- **`SYSTEM.OPENURL`** — args: string
- **`SYSTEM.SETCLIPBOARD`** — args: string
- **`SYSTEM.SETENV`** — args: string, string
- **`SYSTEM.TOTALMEMORY`** — args: (none) → int
- **`SYSTEM.USERNAME`** — args: (none) → string
- **`SYSTEM.VERSION`** — args: (none) → string — MoonBasic release string (e.g. 1.0.0-GOLD); informational only.

### RENDER

- **`RENDER.BEGIN`** — args: (none)
- **`RENDER.BEGIN`** — args: handle
- **`RENDER.BEGIN3D`** — args: handle — Alias for CAMERA.BEGIN: 3D camera heap handle from CAMERA.CREATE or CreateCamera (deprecated alias CAMERA.MAKE)
- **`RENDER.BEGINFRAME`** — args: (none)
- **`RENDER.BEGINMODE2D`** — args: (none)
- **`RENDER.BEGINMODE3D`** — args: (none)
- **`RENDER.BEGINSHADER`** — args: handle
- **`RENDER.CLEAR`** — args: (none)
- **`RENDER.CLEAR`** — args: handle
- **`RENDER.CLEAR`** — args: int, int, int
- **`RENDER.CLEAR`** — args: int, int, int, int
- **`RENDER.CLEARCACHE`** — args: (none)
- **`RENDER.CLEARSCISSOR`** — args: (none)
- **`RENDER.DRAWFPS`** — args: int, int
- **`RENDER.END`** — args: (none)
- **`RENDER.END3D`** — args: (none) — Alias for CAMERA.END (no arguments)
- **`RENDER.ENDFRAME`** — args: (none)
- **`RENDER.ENDMODE2D`** — args: (none)
- **`RENDER.ENDMODE3D`** — args: (none)
- **`RENDER.ENDSHADER`** — args: (none)
- **`RENDER.FRAME`** — args: (none)
- **`RENDER.HEIGHT`** — args: (none) → int
- **`RENDER.SCREENSHOT`** — args: string
- **`RENDER.SET2DAMBIENT`** — args: int, int, int, int
- **`RENDER.SET2DAmbIENT`** — args: int, int, int, int
- **`RENDER.SETAMBIENT`** — args: float, float, float
- **`RENDER.SETAMBIENT`** — args: float, float, float, float
- **`RENDER.SETBACKGROUND`** — args: int, int, int
- **`RENDER.SETBLEND`** — args: int
- **`RENDER.SETBLENDMODE`** — args: int
- **`RENDER.SETBLOOM`** — args: float — POST.BLOOM threshold; intensity defaults to 1
- **`RENDER.SETBLOOM`** — args: float, float — POST.BLOOM threshold and intensity
- **`RENDER.SETCULLFACE`** — args: int
- **`RENDER.SETDEPTHMASK`** — args: bool
- **`RENDER.SETDEPTHTEST`** — args: bool
- **`RENDER.SETDEPTHWRITE`** — args: bool
- **`RENDER.SETFOG`** — args: float, float, float, float, float, float — Fog RGB, near, far, density â€” FOG.* + WORLD.FOGDENSITY
- **`RENDER.SETFPS`** — args: int
- **`RENDER.SETIBLINTENSITY`** — args: float
- **`RENDER.SETIBLSPLIT`** — args: float, float
- **`RENDER.SETMODE`** — args: string
- **`RENDER.SETMSAA`** — args: bool
- **`RENDER.SETPOSTPROCESS`** — args: handle
- **`RENDER.SETSCISSOR`** — args: int, int, int, int
- **`RENDER.SETSHADOWMAPSIZE`** — args: int
- **`RENDER.SETSKYBOX`** — args: string
- **`RENDER.SETTONEMAPPING`** — args: int
- **`RENDER.SETWIREFRAME`** — args: bool
- **`RENDER.WIDTH`** — args: (none) → int

### SCENE

- **`SCENE.APPLYPHYSICS`** — args: handle — Automatically parses glTF Extras to generate Jolt colliders.
- **`SCENE.CLEARSCENE`** — args: (none)
- **`SCENE.CURRENT`** — args: (none) → string
- **`SCENE.DRAW`** — args: (none)
- **`SCENE.LOAD`** — args: string
- **`SCENE.LOADASYNC`** — args: string
- **`SCENE.LOADSCENE`** — args: any
- **`SCENE.LOADWITHTRANSITION`** — args: string, string, float
- **`SCENE.REGISTER`** — args: string, string
- **`SCENE.SAVESCENE`** — args: any
- **`SCENE.SETHANDLERS`** — args: string, string
- **`SCENE.SWITCH`** — args: handle, float — Smoothly transitions levels.
- **`SCENE.UPDATE`** — args: float

### ENTITY

- **`ENTITY.ADDFORCE`** — args: int, float, float, float
- **`ENTITY.ADDPHYSICS`** — args: int, string, string — One-line Jolt body: motion (static/dynamic), shape (box/capsule/sphere)
- **`ENTITY.ADDPHYSICS`** — args: int, string, string, float
- **`ENTITY.ADDTRAIL`** — args: handle, int
- **`ENTITY.ADDTRIANGLE`** — args: handle, int, int, int
- **`ENTITY.ADDVERTEX`** — args: handle, float, float, float → int
- **`ENTITY.ADDWOBBLE`** — args: handle, float, float
- **`ENTITY.ALIGNTOVECTOR`** — args: int, float, float, float, int
- **`ENTITY.ALPHA`** — args: int, float → handle — Easy Mode: Set entity transparency (0.0 to 1.0)
- **`ENTITY.ANIMATE`** — args: int, any, any
- **`ENTITY.ANIMATETOWARD`** — args: int, float, float, float, float — Linear world lerp over duration (seconds); advanced in ENTITY.UPDATE
- **`ENTITY.ANIMCOUNT`** — args: int → int
- **`ENTITY.ANIMINDEX`** — args: int → int
- **`ENTITY.ANIMLENGTH`** — args: int → float
- **`ENTITY.ANIMNAME`** — args: any, int → string
- **`ENTITY.ANIMTIME`** — args: int → float
- **`ENTITY.APPLYGRAVITY`** — args: int, float, float
- **`ENTITY.APPLYIMPULSE`** — args: int, float, float, float — Same as ENTITY.ADDFORCE / ApplyEntityImpulse (velocity change; not Jolt BodyInterface impulse until exposed)
- **`ENTITY.APPLYTORQUE`** — args: handle, float, float, float — Spins physics object.
- **`ENTITY.ATTACH`** — args: handle, handle, float, float, float — Welds entities together with offset.
- **`ENTITY.BLEND`** — args: int, int
- **`ENTITY.BOX`** — args: int, float, float, float
- **`ENTITY.CANSEE`** — args: int, int, float, float → bool — Vision cone (degrees) + max distance + unobstructed Jolt ray to target
- **`ENTITY.CHECKCOLLISION`** — args: int, int → bool — True if two entities had a Jolt contact last step (same as EntityCollided)
- **`ENTITY.CHECKRADIUS`** — args: handle, float, string → handle — Check sensor
- **`ENTITY.CLAMPTOTERRAIN`** — args: int, handle — Sets Y from terrain height at entity XZ (offset 0); alias of TERRAIN.SNAPY argument order swap
- **`ENTITY.CLEARPHYSBUFFER`** — args: int — Remove physics matrix buffer binding from entity
- **`ENTITY.CLEARSCENE`** — args: (none)
- **`ENTITY.COLLIDE`** — args: int, int
- **`ENTITY.COLLIDED`** — args: int → bool
- **`ENTITY.COLLISIONLAYER`** — args: int, int — Reserved 0..31 layer id for future Jolt bitmask filtering (stored on entity)
- **`ENTITY.COLLISIONNX`** — args: int → float
- **`ENTITY.COLLISIONNY`** — args: int → float
- **`ENTITY.COLLISIONNZ`** — args: int → float
- **`ENTITY.COLLISIONOTHER`** — args: int → int
- **`ENTITY.COLLISIONX`** — args: int → float
- **`ENTITY.COLLISIONY`** — args: int → float
- **`ENTITY.COLLISIONZ`** — args: int → float
- **`ENTITY.COLOR`** — args: int, handle → handle
- **`ENTITY.COLOR`** — args: int, int, int, int → handle
- **`ENTITY.COLOR`** — args: int, int, int, int, int → handle
- **`ENTITY.COLORPULSE`** — args: handle, handle, handle, float — Pulses color.
- **`ENTITY.COPY`** — args: int → int
- **`ENTITY.COUNTCHILDREN`** — args: int → int
- **`ENTITY.CREATE`** — args: (none) → int
- **`ENTITY.CREATEBOX`** — args: float → int — Uniform cube: size used for width, height, and depth (alias ENTITY.CREATECUBE)
- **`ENTITY.CREATEBOX`** — args: float, float, float → int
- **`ENTITY.CREATECONE`** — args: int, int, int, int → handle
- **`ENTITY.CREATECUBE`** — args: float, float, float → handle
- **`ENTITY.CREATECYLINDER`** — args: float, float, int → int
- **`ENTITY.CREATEENTITY`** — args: (none) → int
- **`ENTITY.CREATEMESH`** — args: any → int — Procedural mesh: optional parentEntity; use AddVertex/UpdateMesh
- **`ENTITY.CREATEPLANE`** — args: float → int
- **`ENTITY.CREATESPHERE`** — args: float → int — Radius only â€” default 16 segments
- **`ENTITY.CREATESPHERE`** — args: float, int → int
- **`ENTITY.CREATESPRITE`** — args: string → int
- **`ENTITY.CREATESPRITE`** — args: string, int → int
- **`ENTITY.CREATESPRITE`** — args: handle, float, float → int — Billboard from TEXTURE handle (atlas / TEXTURE.LOADANIM)
- **`ENTITY.CREATESPRITE`** — args: handle, float, float, int → int
- **`ENTITY.CREATESURFACE`** — args: int → handle
- **`ENTITY.CROSSFADE`** — args: int, any, float
- **`ENTITY.CURRENTANIM`** — args: any → string
- **`ENTITY.CUTJUMP`** — args: handle
- **`ENTITY.DAMAGE`** — args: handle, float
- **`ENTITY.DELTAX`** — args: int, int → float
- **`ENTITY.DELTAY`** — args: int, int → float
- **`ENTITY.DELTAZ`** — args: int, int → float
- **`ENTITY.DIST`** — args: int, int → float — 3D distance between two entities (alias of ENTITY.DISTANCE semantics)
- **`ENTITY.DISTANCE`** — args: int, int → float
- **`ENTITY.DISTANCETO`** — args: handle, handle → float — Returns distance.
- **`ENTITY.DRAW`** — args: int
- **`ENTITY.DRAWALL`** — args: (none)
- **`ENTITY.EMITPARTICLES`** — args: handle, handle — Attaches particles to entity.
- **`ENTITY.ENTITIESINBOX`** — args: float, float, float, float, float, float
- **`ENTITY.ENTITIESINGROUP`** — args: any
- **`ENTITY.ENTITIESINRADIUS`** — args: float, float, float, float
- **`ENTITY.ENTITYPITCH`** — args: int, any → float
- **`ENTITY.ENTITYROLL`** — args: int, any → float
- **`ENTITY.ENTITYX`** — args: int, any → float
- **`ENTITY.ENTITYY`** — args: int, any → float
- **`ENTITY.ENTITYYAW`** — args: int, any → float
- **`ENTITY.ENTITYZ`** — args: int, any → float
- **`ENTITY.EXPLODE`** — args: handle, int — Instantly explodes object.
- **`ENTITY.EXTRACTANIMSEQ`** — args: int, any, any
- **`ENTITY.FADE`** — args: handle, float, float, float — Lerps alpha.
- **`ENTITY.FIND`** — args: any → int
- **`ENTITY.FINDBONE`** — args: int, any → int
- **`ENTITY.FINDBYPROPERTY`** — args: string, string → handle
- **`ENTITY.FINDCHILD`** — args: int, string → int
- **`ENTITY.FLEE`** — args: handle, handle, float, float — Runs away.
- **`ENTITY.FLOOR`** — args: int → float
- **`ENTITY.FREE`** — args: int
- **`ENTITY.FREEENTITIES`** — args: handle
- **`ENTITY.FX`** — args: int, int
- **`ENTITY.GETALPHA`** — args: int → float
- **`ENTITY.GETBONEPOS`** — args: int, string → handle
- **`ENTITY.GETBONEROT`** — args: int, string → handle
- **`ENTITY.GETBOUNDS`** — args: int → handle
- **`ENTITY.GETBUOYANCY`** — args: int → float — Alias of PHYSICS.GETBUOYANCY
- **`ENTITY.GETCHILD`** — args: int, int → int
- **`ENTITY.GETCLOSESTWITHTAG`** — args: int, float, string → int — Nearest entity within radius matching name/tag glob (same rules as PLAYER.GETNEARBY)
- **`ENTITY.GETCOLOR`** — args: int → array
- **`ENTITY.GETDISTANCE`** — args: int, int → float
- **`ENTITY.GETGROUNDNORMAL`** — args: int → handle — World ground normal under entity (CharacterVirtual if PLAYER.CREATE; else short downward Jolt ray)
- **`ENTITY.GETMETADATA`** — args: int, string → string
- **`ENTITY.GETOVERLAPCOUNT`** — args: int, string → int — Counts tagged entities whose pivot lies in zone entity world AABB (sphere prefilter)
- **`ENTITY.GETPOS`** — args: int → handle
- **`ENTITY.GETPOSITION`** — args: int → handle
- **`ENTITY.GETROT`** — args: int → handle
- **`ENTITY.GETSCALE`** — args: int → handle
- **`ENTITY.GETSTATE`** — args: handle → int — Returns string AI state.
- **`ENTITY.GETXZ`** — args: int → handle
- **`ENTITY.GHOSTMODE`** — args: handle, float — Disables collisions temporarily.
- **`ENTITY.GRAVITY`** — args: int, float
- **`ENTITY.GROUNDED`** — args: int → bool
- **`ENTITY.GROUPADD`** — args: any, int
- **`ENTITY.GROUPCREATE`** — args: any
- **`ENTITY.GROUPREMOVE`** — args: any, int
- **`ENTITY.HASTAG`** — args: int, string → bool — Glob match on Blender tag or entity name only (stricter than ENTITY.ISTYPE)
- **`ENTITY.HIDE`** — args: int
- **`ENTITY.INFRUSTUM`** — args: int → bool — True if entity AABB intersects active CAMERA.BEGIN frustum (same as ENTITY.INVIEW without passing camera)
- **`ENTITY.INFRUSTUM`** — args: handle, handle → int — Boolean bounds.
- **`ENTITY.INSTANCE`** — args: int → int
- **`ENTITY.INSTANCEGRID`** — args: int, int, int, float → int
- **`ENTITY.INVIEW`** — args: int, handle → bool
- **`ENTITY.ISALIVE`** — args: handle → bool
- **`ENTITY.ISPLAYING`** — args: int → bool
- **`ENTITY.ISSUBMERGED`** — args: int → float — Fraction 0..1 of entity vertical extent below water surface (any overlapping WATER volume)
- **`ENTITY.ISTYPE`** — args: int, string → bool
- **`ENTITY.ISWALLSLIDING`** — args: handle → bool
- **`ENTITY.JUMP`** — args: int, float
- **`ENTITY.LINEOFSIGHT`** — args: int, int → bool — Unobstructed Jolt ray from observer eye to target (no FOV); sensors still occlude until filtered
- **`ENTITY.LINKPHYSBUFFER`** — args: int, int — Bind entity to Jolt shared matrix slot index (use BODY3D.BUFFERINDEX on the body)
- **`ENTITY.LOAD`** — args: any → int — Alias of ENTITY.LOADMESH â€” static model path (Raylib-supported formats), optional parentEntity
- **`ENTITY.LOADANIMATEDMESH`** — args: any → int
- **`ENTITY.LOADANIMATIONS`** — args: int, string
- **`ENTITY.LOADMESH`** — args: any → int
- **`ENTITY.LOADSCENE`** — args: any
- **`ENTITY.LOADSPRITE`** — args: string → int
- **`ENTITY.LOADSPRITE`** — args: string, int → int
- **`ENTITY.LOOKAT`** — args: handle, float, float — Instantly rotates an entity to face a point.
- **`ENTITY.LOOKAT`** — args: int, float, float, float — Face world point (entity, targetX, targetY, targetZ); sets pitch/yaw
- **`ENTITY.MAGNETTO`** — args: handle, float, float, float, float
- **`ENTITY.MAKE`** — args: (none) → int — DEPRECATED alias of ENTITY.CREATE. Use ENTITY.CREATE.
- **`ENTITY.MAKEBOX`** — args: float → int — DEPRECATED alias of ENTITY.CREATEBOX. Use ENTITY.CREATEBOX. Uniform cube: size used for width, height, and depth (alias ENTITY.CREATECUBE)
- **`ENTITY.MAKEBOX`** — args: float, float, float → int — DEPRECATED alias of ENTITY.CREATEBOX. Use ENTITY.CREATEBOX.
- **`ENTITY.MAKECONE`** — args: int, int, int, int → handle — DEPRECATED alias of ENTITY.CREATECONE. Use ENTITY.CREATECONE(...).
- **`ENTITY.MAKECUBE`** — args: float, float, float → int — DEPRECATED alias of ENTITY.CREATECUBE. Use ENTITY.CREATECUBE.
- **`ENTITY.MAKECYLINDER`** — args: float, float, int → int — DEPRECATED alias of ENTITY.CREATECYLINDER. Use ENTITY.CREATECYLINDER.
- **`ENTITY.MAKEENTITY`** — args: (none) → int — DEPRECATED alias of ENTITY.CREATEENTITY. Use ENTITY.CREATEENTITY.
- **`ENTITY.MAKEMESH`** — args: any → int — DEPRECATED alias of ENTITY.CREATEMESH. Use ENTITY.CREATEMESH. Procedural mesh: optional parentEntity; use AddVertex/UpdateMesh
- **`ENTITY.MAKEPLANE`** — args: float → int — DEPRECATED alias of ENTITY.CREATEPLANE. Use ENTITY.CREATEPLANE.
- **`ENTITY.MAKESPHERE`** — args: float → int — DEPRECATED alias of ENTITY.CREATESPHERE. Use ENTITY.CREATESPHERE. Radius only â€” default 16 segments
- **`ENTITY.MAKESPHERE`** — args: float, int → int — DEPRECATED alias of ENTITY.CREATESPHERE. Use ENTITY.CREATESPHERE.
- **`ENTITY.MAKESPRITE`** — args: string → int — DEPRECATED alias of ENTITY.CREATESPRITE. Use ENTITY.CREATESPRITE.
- **`ENTITY.MAKESPRITE`** — args: string, int → int — DEPRECATED alias of ENTITY.CREATESPRITE. Use ENTITY.CREATESPRITE.
- **`ENTITY.MAKESPRITE`** — args: handle, float, float → int — DEPRECATED alias of ENTITY.CREATESPRITE. Use ENTITY.CREATESPRITE. Billboard from TEXTURE handle (atlas / TEXTURE.LOADANIM)
- **`ENTITY.MAKESPRITE`** — args: handle, float, float, int → int — DEPRECATED alias of ENTITY.CREATESPRITE. Use ENTITY.CREATESPRITE.
- **`ENTITY.MAKESURFACE`** — args: int → handle — DEPRECATED alias of ENTITY.CREATESURFACE. Use ENTITY.CREATESURFACE.
- **`ENTITY.MATRIXELEMENT`** — args: int, int, int → float
- **`ENTITY.MOVE`** — args: int, float, float, float
- **`ENTITY.MOVECAMERARELATIVE`** — args: int, float, float, handle — World XZ step from camera yaw: forward/strafe are deltas (typically speed*dt*input); camera is a Camera3D handle.
- **`ENTITY.MOVEENTITY`** — args: int, float, float, float
- **`ENTITY.MOVERELATIVE`** — args: int, float, float, float, float
- **`ENTITY.MOVETOWARD`** — args: handle, handle, float — Moves an entity toward another entity at constant speed (XZ toward target, Y preserved).
- **`ENTITY.MOVETOWARD`** — args: handle, float, float, float — Moves an entity toward a coordinate.
- **`ENTITY.MOVEWITHCAMERA`** — args: int, handle, float, float, float — Horizontal walk velocity (units/s) from camera XZ strafe basis (eyeâ†’target on ground). forwardAxis/strafeAxis are typically Input.Axis âˆ’1..1; preserves vertical velocity. Dot: player.MoveWithCamera(cam, â€¦).
- **`ENTITY.NAVTO`** — args: handle, float, float, float
- **`ENTITY.ONDEATHDROP`** — args: handle, string
- **`ENTITY.ONHIT`** — args: handle, string — Fires MB callback on collision.
- **`ENTITY.ORDER`** — args: int, int
- **`ENTITY.OUTLINE`** — args: int, float, handle — Apply a highlighted outline effect to a model.
- **`ENTITY.P`** — args: int → float — Easy Mode: Get Pitch of entity
- **`ENTITY.P`** — args: int, float — Easy Mode: Set Pitch of entity
- **`ENTITY.PARENT`** — args: int, int — Attach child entity to parent (optional third arg: global preserve world position; default true)
- **`ENTITY.PARENT`** — args: int, int, any
- **`ENTITY.PARENTCLEAR`** — args: int
- **`ENTITY.PATROL`** — args: handle, handle, float — Loops an entity across a WAYPOINT array handle at speed.
- **`ENTITY.PHYSICS`** — args: int, string, float — Quickly setup a physics body for an entity (auto-sizes based on model/shape).
- **`ENTITY.PHYSICS`** — args: int, string, float, float, float — Quickly setup a physics body with mass, friction, and restitution.
- **`ENTITY.PHYSICS`** — args: int, string, float, float, float, bool — Quickly setup a physics body with mass, friction, restitution, and CCD enabled.
- **`ENTITY.PHYSICSMOTION`** — args: int, string — Toggle physics motion type (STATIC, DYNAMIC, KINEMATIC).
- **`ENTITY.PICK`** — args: int, float → bool
- **`ENTITY.PICKMODE`** — args: int, int
- **`ENTITY.PLAY`** — args: int, any
- **`ENTITY.PLAYNAME`** — args: int, string
- **`ENTITY.POINTAT`** — args: int, int
- **`ENTITY.POINTENTITY`** — args: int, int
- **`ENTITY.POLLMESSAGE`** — args: int → string
- **`ENTITY.POS`** — args: int, float, float, float → handle — Easy Mode shorthand for positioning an entity
- **`ENTITY.POSITION`** — args: int, float, float, float, any — Alias of ENTITY.SETPOS â€” set world or local position
- **`ENTITY.POSITIONENTITY`** — args: int, float, float, float, any
- **`ENTITY.PUSH`** — args: int, float, float, float — Apply Jolt impulse (requires ENTITY.ADDPHYSICS)
- **`ENTITY.PUSHOUTOFGEOMETRY`** — args: int — Best-effort depenetration: nudges entity world Y up slightly
- **`ENTITY.R`** — args: int → float — Easy Mode: Get Roll of entity
- **`ENTITY.R`** — args: int, float — Easy Mode: Set Roll of entity
- **`ENTITY.RADIUS`** — args: int, float
- **`ENTITY.RAYCAST`** — args: handle, float → handle — Raycast sensor
- **`ENTITY.RAYCAST`** — args: float, float, float, float, float, float, float → int — Jolt ray cast; returns first hit entity or 0 (same query path as PHYSICS3D/PICK)
- **`ENTITY.RAYHIT`** — args: int, float, float, float, float, float, float → bool
- **`ENTITY.RGB`** — args: int, int, int, int — Easy Mode: Set entity color (id, r, g, b)
- **`ENTITY.ROT`** — args: int, float, float, float → handle — Easy Mode shorthand for rotating an entity (absolute)
- **`ENTITY.ROTATE`** — args: int, float, float, float
- **`ENTITY.ROTATEENTITY`** — args: int, float, float, float, any
- **`ENTITY.SAVESCENE`** — args: any
- **`ENTITY.SCA`** — args: int, float, float, float — Easy Mode shorthand for scaling an entity (absolute)
- **`ENTITY.SCALE`** — args: int, float, float, float → handle
- **`ENTITY.SCROLLMATERIAL`** — args: int, float, float — Add (du,dv) to material 0 scroll (same as MODEL.SCROLLTEXTURE)
- **`ENTITY.SENDMESSAGE`** — args: int, string
- **`ENTITY.SETANIMATION`** — args: int, handle, float — Cycle IMAGE.LOADSEQUENCE/LOADGIF frames onto sprite texture at fps
- **`ENTITY.SETANIMATION`** — args: int, int, float — Second arg 0 clears image-sequence animation
- **`ENTITY.SETANIMATION`** — args: int, handle, float, bool
- **`ENTITY.SETANIMFRAME`** — args: int, float
- **`ENTITY.SETANIMINDEX`** — args: int, any
- **`ENTITY.SETANIMLOOP`** — args: int, any
- **`ENTITY.SETANIMSPEED`** — args: int, float
- **`ENTITY.SETANIMTIME`** — args: int, float
- **`ENTITY.SETBOUNCE`** — args: int, float
- **`ENTITY.SETBOUNCINESS`** — args: int, float — Sets restitution (bounciness) on an entity's Jolt body; 0 = no bounce. Alias of PHYSICS.BOUNCE.
- **`ENTITY.SETBUOYANCY`** — args: int, float — Alias of PHYSICS.SETBUOYANCY â€” per-entity density hint for buoyancy
- **`ENTITY.SETCOLLISIONGROUP`** — args: int, int — Alias for ENTITY.COLLISIONLAYER (collision group / layer 0..31)
- **`ENTITY.SETCULLMODE`** — args: handle, int
- **`ENTITY.SETDETAILTEXTURE`** — args: int, handle — Bind secondary map as MATERIAL_MAP_NORMAL for blending/detail
- **`ENTITY.SETFRICTION`** — args: int, float
- **`ENTITY.SETGRAVITY`** — args: int, float
- **`ENTITY.SETGRAVITYSCALE`** — args: handle, float
- **`ENTITY.SETHEALTH`** — args: handle, float
- **`ENTITY.SETMASS`** — args: int, float
- **`ENTITY.SETNAME`** — args: int, any
- **`ENTITY.SETPOS`** — args: int, float, float, float, any
- **`ENTITY.SETPOSITION`** — args: int, float, float, float, any — DEPRECATED alias of ENTITY.SETPOS. Use ENTITY.SETPOS. Deprecated alias of ENTITY.SETPOS â€” set world or local position
- **`ENTITY.SETROTATION`** — args: int, float, float, float, any — Absolute pitch/yaw/roll degrees â€” alias of ENTITY.ROTATEENTITY
- **`ENTITY.SETSHADER`** — args: int, handle
- **`ENTITY.SETSHADER`** — args: handle, int — Binds an active Shader Library component to the entity.
- **`ENTITY.SETSPRITEFRAME`** — args: int, int — Set atlas frame on billboard bound to a TEXTURE object
- **`ENTITY.SETSTATIC`** — args: int
- **`ENTITY.SETTAG`** — args: handle, string
- **`ENTITY.SETTEXTUREFLIP`** — args: handle, float, float — Modifies UV scaling for horizontal/vertical mirroring.
- **`ENTITY.SETTEXTUREMAP`** — args: int, any, handle
- **`ENTITY.SETTEXTURESCROLL`** — args: handle, float, float — Injects offsets into the shader for animated water/lava.
- **`ENTITY.SETTRIGGER`** — args: int
- **`ENTITY.SETVISIBLE`** — args: int, any — Alias of ENTITY.VISIBLE
- **`ENTITY.SETWEIGHT`** — args: handle, float — Changes entity mass.
- **`ENTITY.SHININESS`** — args: int, float
- **`ENTITY.SHOW`** — args: int
- **`ENTITY.SLIDE`** — args: int, any
- **`ENTITY.SNAPTO`** — args: int, int — Instantly align one entity to another's transform.
- **`ENTITY.SPRITEVIEWMODE`** — args: int, int
- **`ENTITY.SQUASH`** — args: int, float, float — Juice: squash scale Y then tween back
- **`ENTITY.STOPANIM`** — args: int
- **`ENTITY.TAG`** — args: handle, string — Sets spatial tag.
- **`ENTITY.TEXTURE`** — args: int, any
- **`ENTITY.TFORMPOINT`** — args: float, float, float, int, int → handle
- **`ENTITY.TFORMVECTOR`** — args: float, float, float, int, int → handle
- **`ENTITY.TRANSITION`** — args: int, string, float
- **`ENTITY.TRANSLATE`** — args: int, float, float, float
- **`ENTITY.TRANSLATEENTITY`** — args: int, float, float, float, any
- **`ENTITY.TURN`** — args: int, float, float, float — Add pitch/yaw/roll degrees â€” alias of ENTITY.ROTATE / TURNENTITY
- **`ENTITY.TURNENTITY`** — args: int, float, float, float, any
- **`ENTITY.TURNTOWARD`** — args: handle, float, float, float — Slowly rotates the entity to face a target over time.
- **`ENTITY.TWEEN`** — args: int, string, any, float, string — Animate properties (position, scale, rotation) using easing functions (bounce, elastic, etc).
- **`ENTITY.TYPE`** — args: int, int
- **`ENTITY.UNPARENT`** — args: int — Alias of ENTITY.PARENTCLEAR â€” detach and keep world position
- **`ENTITY.UPDATE`** — args: float
- **`ENTITY.UPDATEMESH`** — args: int
- **`ENTITY.VELOCITY`** — args: int, float, float, float
- **`ENTITY.VERTEXX`** — args: handle, int → float
- **`ENTITY.VERTEXY`** — args: handle, int → float
- **`ENTITY.VERTEXZ`** — args: handle, int → float
- **`ENTITY.VISIBLE`** — args: int, any
- **`ENTITY.W`** — args: int → float — Easy Mode: Get Yaw (W) of entity
- **`ENTITY.W`** — args: int, float — Easy Mode: Set Yaw (W) of entity
- **`ENTITY.WANDER`** — args: handle, float, float, float, float — Moves an NPC randomly within a zone.
- **`ENTITY.WASGROUNDED`** — args: handle → bool
- **`ENTITY.WITHINRADIUS`** — args: handle, handle, float → bool — True if 3D distance between entities is <= maxDistance (simple sphere check; not Jolt physics).
- **`ENTITY.X`** — args: int → float — Easy Mode: Get X position of entity
- **`ENTITY.X`** — args: int, float — Easy Mode: Set X position of entity
- **`ENTITY.Y`** — args: int → float
- **`ENTITY.Y`** — args: int, float
- **`ENTITY.Z`** — args: int → float
- **`ENTITY.Z`** — args: int, float
- **`PHYSICS.CCD`** — args: int, bool — Enable Continuous Collision Detection to prevent high-speed tunneling.

*412 overloads in this section.*

---

## Camera and light

Guide: [02-CAMERA-LIGHT.md](02-CAMERA-LIGHT.md)

### CAMERA

- **`CAMERA.BEGIN`** — args: handle → handle
- **`CAMERA.CAMERAFOLLOW`** — args: handle, int, float, float, float
- **`CAMERA.CLEARFPSMODE`** — args: handle
- **`CAMERA.CREATE`** — args: (none) → handle — Returns a Camera3D heap handle (canonical; deprecated alias: CAMERA.MAKE)
- **`CAMERA.CREATE`** — args: (none)
- **`CAMERA.END`** — args: (none)
- **`CAMERA.END`** — args: handle
- **`CAMERA.FOLLOW`** — args: handle, handle, float, float — Spring math camera tracker.
- **`CAMERA.FOLLOW`** — args: handle, float, float, float, float, float, float, float
- **`CAMERA.FOLLOWENTITY`** — args: handle, int, float, float, float
- **`CAMERA.FOV`** — args: handle → float — Property alias for CAMERA.GETFOV
- **`CAMERA.FREE`** — args: handle
- **`CAMERA.GETACTIVE`** — args: (none) → handle
- **`CAMERA.GETFOV`** — args: handle → float — Get camera field of view.
- **`CAMERA.GETMATRIX`** — args: handle → handle
- **`CAMERA.GETPOS`** — args: handle → handle
- **`CAMERA.GETPROJECTION`** — args: handle → int — Returns the camera projection mode (0=Persp, 1=Ortho)
- **`CAMERA.GETRAY`** — args: handle, float, float
- **`CAMERA.GETROT`** — args: handle → handle
- **`CAMERA.GETTARGET`** — args: handle → handle
- **`CAMERA.GETUP`** — args: handle → array — Returns the camera UP vector as a Vec3 handle
- **`CAMERA.GETVIEWRAY`** — args: float, float, handle, int, int
- **`CAMERA.GETYAW`** — args: handle — Alias of CAMERA.YAW.
- **`CAMERA.ISONSCREEN`** — args: handle, float, float, float → bool
- **`CAMERA.ISONSCREEN`** — args: handle, float, float, float, float → bool
- **`CAMERA.LERPTO`** — args: handle, int, float → handle — Smoothly interpolate camera target toward an entity.
- **`CAMERA.LOOKAT`** — args: handle, float, float, float → handle
- **`CAMERA.LOOKATENTITY`** — args: handle, int — Sets camera target to entity world position (same idea as Blitz PointAt)
- **`CAMERA.MAKE`** — args: (none) — DEPRECATED alias of CAMERA.CREATE. Use CAMERA.CREATE.
- **`CAMERA.MAKE`** — args: (none) → handle — DEPRECATED alias of CAMERA.CREATE. Returns a Camera3D heap handle.
- **`CAMERA.MOUSERAY`** — args: handle → handle
- **`CAMERA.MOVE`** — args: handle, float, float, float → handle
- **`CAMERA.ORBIT`** — args: handle, float, float, float, float, float, float
- **`CAMERA.ORBITAROUND`** — args: handle, float, float, float, float, float, float
- **`CAMERA.ORBITAROUNDEG`** — args: handle, float, float, float, float, float, float
- **`CAMERA.ORBITCAMERA`** — args: handle, float, float, float → float
- **`CAMERA.ORBITENTITY`** — args: handle, int, float, float, float
- **`CAMERA.PICK`** — args: handle, float, float → handle
- **`CAMERA.POINTATENTITY`** — args: handle, int — Alias of CAMERA.LOOKATENTITY
- **`CAMERA.POS`** — args: handle → array — Property alias for CAMERA.GETPOS
- **`CAMERA.PROJECT`** — args: handle, float, float, float → handle — Alias of CAMERA.WORLDTOSCREEN â€” world point to screen [sx,sy]
- **`CAMERA.PROJECTION`** — args: handle → int — Property alias for CAMERA.GETPROJECTION
- **`CAMERA.RAYCASTMOUSE`** — args: handle → int — Raycast from mouse through camera; returns entity id or 0.
- **`CAMERA.ROT`** — args: handle → array — Property alias for CAMERA.GETROT
- **`CAMERA.ROTATE`** — args: handle, float, float, float
- **`CAMERA.SETACTIVE`** — args: handle → handle
- **`CAMERA.SETFOV`** — args: handle, float → handle
- **`CAMERA.SETFPSMODE`** — args: handle, float → handle
- **`CAMERA.SETMODE`** — args: handle, any → handle — 0/1 or perspective/orthographic â€” alias-friendly CAMERA.SETPROJECTION
- **`CAMERA.SETORBIT`** — args: handle, float, float, float, float, float, float → handle
- **`CAMERA.SETORBITKEYS`** — args: handle, float, float → handle — Raylib key codes for orbit yaw (0 disables that side).
- **`CAMERA.SETORBITKEYSPEED`** — args: handle, float → handle — Keyboard orbit yaw rate in radians per second.
- **`CAMERA.SETORBITLIMITS`** — args: handle, float, float, float, float → handle — Clamp pitch (radians) and orbit distance for CAMERA.ORBIT (entity).
- **`CAMERA.SETORBITSPEED`** — args: handle, float, float → handle — Mouse drag sensitivity and mouse wheel zoom scale for orbit-follow.
- **`CAMERA.SETPOS`** — args: handle, float, float, float → handle
- **`CAMERA.SETPOSITION`** — args: handle, float, float, float → handle — DEPRECATED alias of CAMERA.SETPOS. Use CAMERA.SETPOS.
- **`CAMERA.SETPROJECTION`** — args: handle, int → handle
- **`CAMERA.SETRANGE`** — args: handle, float, float → handle
- **`CAMERA.SETTARGET`** — args: handle, float, float, float → handle
- **`CAMERA.SETTARGETENTITY`** — args: handle, int → handle
- **`CAMERA.SETUP`** — args: handle, float, float, float → handle
- **`CAMERA.SHAKE`** — args: handle, float, float → handle
- **`CAMERA.SMOOTHEXP`** — args: float, float, float, float → float — Exponential smoothing: current toward target using (1-exp(-smoothHz*dt)); for orbit angles
- **`CAMERA.TARGET`** — args: handle → array — Property alias for CAMERA.GETTARGET
- **`CAMERA.TURN`** — args: handle, float, float, float
- **`CAMERA.TURNLEFT`** — args: handle, float → float
- **`CAMERA.TURNRIGHT`** — args: handle, float → float
- **`CAMERA.UNPROJECT`** — args: handle, float, float → handle — Screen (x,y) to world ray â€” alias of CAMERA.GETRAY / PICK
- **`CAMERA.UP`** — args: handle → array — Property alias for CAMERA.GETUP
- **`CAMERA.UPDATEFPS`** — args: handle
- **`CAMERA.USEMOUSEORBIT`** — args: handle, bool — Enable/disable mouse contribution to CAMERA.ORBIT (entity) orbit-follow.
- **`CAMERA.USEORBITRIGHTMOUSE`** — args: handle, bool — If true (default), mouse orbit only while right button is held; if false, mouse moves orbit without RMB.
- **`CAMERA.WORLDTOSCREEN`** — args: handle, float, float, float → handle
- **`CAMERA.WORLDTOSCREEN2D`** — args: handle, float, float, float → handle
- **`CAMERA.XZBASIS`** — args: handle → array — Returns planar [fwdX, fwdZ, rightX, rightZ] vectors for camera-relative movement.
- **`CAMERA.YAW`** — args: handle — Orbit yaw in radians (internal state) for aligning entities with cam.Orbit(entity, dist).
- **`CAMERA.ZOOM`** — args: handle, float

### CAMERA2D

- **`CAMERA2D.BEGIN`** — args: (none)
- **`CAMERA2D.BEGIN`** — args: handle
- **`CAMERA2D.CREATE`** — args: (none) → handle
- **`CAMERA2D.END`** — args: (none)
- **`CAMERA2D.FOLLOW`** — args: handle, handle, float, float
- **`CAMERA2D.FREE`** — args: handle
- **`CAMERA2D.GETMATRIX`** — args: handle → handle
- **`CAMERA2D.GETOFFSET`** — args: handle → array
- **`CAMERA2D.GETPOS`** — args: handle → array
- **`CAMERA2D.GETROTATION`** — args: handle → float
- **`CAMERA2D.GETZOOM`** — args: handle → float
- **`CAMERA2D.MAKE`** — args: (none) → handle — DEPRECATED alias of CAMERA2D.CREATE. Use CAMERA2D.CREATE.
- **`CAMERA2D.ROTATION`** — args: handle → float
- **`CAMERA2D.SCREENTOWORLD`** — args: handle, float, float → handle
- **`CAMERA2D.SETOFFSET`** — args: handle, float, float
- **`CAMERA2D.SETROTATION`** — args: handle, float
- **`CAMERA2D.SETTARGET`** — args: handle, float, float
- **`CAMERA2D.SETZOOM`** — args: handle, float
- **`CAMERA2D.TARGETX`** — args: handle → float
- **`CAMERA2D.TARGETY`** — args: handle → float
- **`CAMERA2D.WORLDTOSCREEN`** — args: handle, float, float → handle
- **`CAMERA2D.ZOOMIN`** — args: handle, float
- **`CAMERA2D.ZOOMOUT`** — args: handle, float
- **`CAMERA2D.ZOOMTOMOUSE`** — args: handle, float

### LIGHT

- **`LIGHT.COLOR`** — args: handle → handle — Property alias for LIGHT.GETCOLOR
- **`LIGHT.CREATE`** — args: (none) → handle
- **`LIGHT.CREATE`** — args: string → handle
- **`LIGHT.CREATEDIRECTIONAL`** — args: float, float, float, float, float, float, float → handle — Directional light: direction vector (dx,dy,dz), RGB, energy â€” direction is normalized
- **`LIGHT.CREATEPOINT`** — args: float, float, float, float, float, float, float → handle — Point light at (x,y,z) with RGB (0-255 or 0-1) and intensity (energy)
- **`LIGHT.CREATESPOT`** — args: float, float, float, float, float, float, float, float, float, float, float → handle — Spot: position, target point, RGB, outer cone degrees, energy
- **`LIGHT.DIR`** — args: handle → array — Property alias for LIGHT.GETDIR
- **`LIGHT.ENABLE`** — args: handle, bool → handle
- **`LIGHT.ENABLED`** — args: handle → bool — Property alias for LIGHT.ISENABLED
- **`LIGHT.FREE`** — args: handle
- **`LIGHT.GETCOLOR`** — args: handle → handle — (Returns Color instance handle)
- **`LIGHT.GETCOLOR`** — args: handle → handle — Get light color as Color instance.
- **`LIGHT.GETCOLOR`** — args: handle → handle — Returns a Color heap handle with RGBA components (0-255).
- **`LIGHT.GETDIR`** — args: handle → handle — Get light direction as Vec3.
- **`LIGHT.GETDIR`** — args: handle → array
- **`LIGHT.GETENERGY`** — args: handle → float
- **`LIGHT.GETINNERCONE`** — args: handle → float
- **`LIGHT.GETINNERCONE`** — args: handle → float — Get spotlight inner cone angle.
- **`LIGHT.GETINTENSITY`** — args: handle → float
- **`LIGHT.GETINTENSITY`** — args: handle → float — Get light intensity.
- **`LIGHT.GETINTENSITY`** — args: handle → float
- **`LIGHT.GETOUTERCONE`** — args: handle → float — Get spotlight outer cone angle.
- **`LIGHT.GETOUTERCONE`** — args: handle → float
- **`LIGHT.GETPOS`** — args: handle → array
- **`LIGHT.GETPOS`** — args: handle → handle — Get light position as Vec3.
- **`LIGHT.GETRANGE`** — args: handle → float
- **`LIGHT.GETRANGE`** — args: handle → float — Get light range.
- **`LIGHT.GETROT`** — args: handle → array — Returns [p, y, r] Euler rotation of the light
- **`LIGHT.GETSHADOW`** — args: handle → bool — Check if light has shadows enabled.
- **`LIGHT.GETSHADOW`** — args: handle → bool
- **`LIGHT.INTENSITY`** — args: handle → float — Property alias for LIGHT.GETINTENSITY
- **`LIGHT.ISENABLED`** — args: handle → int
- **`LIGHT.MAKE`** — args: (none) → handle — DEPRECATED alias of LIGHT.CREATE. Use LIGHT.CREATE.
- **`LIGHT.MAKE`** — args: string → handle — DEPRECATED alias of LIGHT.CREATE. Use LIGHT.CREATE.
- **`LIGHT.MAKEDIRECTIONAL`** — args: float, float, float, float, float, float, float → handle — DEPRECATED alias of LIGHT.CREATEDIRECTIONAL. Use LIGHT.CREATEDIRECTIONAL. Directional light: direction vector (dx,dy,dz), RGB, energy â€” direction is normalized
- **`LIGHT.MAKEPOINT`** — args: float, float, float, float, float, float, float → handle — DEPRECATED alias of LIGHT.CREATEPOINT. Use LIGHT.CREATEPOINT. Point light at (x,y,z) with RGB (0â€“255 or 0â€“1) and intensity (energy)
- **`LIGHT.MAKESPOT`** — args: float, float, float, float, float, float, float, float, float, float, float → handle — DEPRECATED alias of LIGHT.CREATESPOT. Use LIGHT.CREATESPOT. Spot: position, target point, RGB, outer cone degrees, energy
- **`LIGHT.POS`** — args: handle → array — Property alias for LIGHT.GETPOS
- **`LIGHT.RANGE`** — args: handle → float — Property alias for LIGHT.GETRANGE
- **`LIGHT.SETCOLOR`** — args: handle, float, float, float → handle
- **`LIGHT.SETCOLOR`** — args: handle, float, float, float, float → handle
- **`LIGHT.SETDIR`** — args: handle, float, float, float → handle
- **`LIGHT.SETINNERCONE`** — args: handle, float → handle
- **`LIGHT.SETINTENSITY`** — args: handle, float → handle
- **`LIGHT.SETOUTERCONE`** — args: handle, float → handle
- **`LIGHT.SETPOS`** — args: handle, float, float, float → handle
- **`LIGHT.SETPOSITION`** — args: handle, float, float, float → handle — DEPRECATED alias of LIGHT.SETPOS. Use LIGHT.SETPOS.
- **`LIGHT.SETRANGE`** — args: handle, float → handle
- **`LIGHT.SETROT`** — args: handle, float, float, float → handle — Sets light orientation using Euler angles (pitch, yaw, roll)
- **`LIGHT.SETSHADOW`** — args: handle, bool → handle
- **`LIGHT.SETSHADOWBIAS`** — args: handle, float → handle
- **`LIGHT.SETSTATE`** — args: handle, bool → handle — Alias of LIGHT.ENABLE
- **`LIGHT.SETTARGET`** — args: handle, float, float, float → handle
- **`LIGHT.SHADOW`** — args: handle → bool — Property alias for LIGHT.GETSHADOW

### LIGHT2D

- **`LIGHT2D.CREATE`** — args: (none) → handle
- **`LIGHT2D.FREE`** — args: handle
- **`LIGHT2D.GETCOLOR`** — args: handle → handle — (Returns Color instance handle)
- **`LIGHT2D.GETINTENSITY`** — args: handle → float
- **`LIGHT2D.GETPOS`** — args: handle → array
- **`LIGHT2D.GETRADIUS`** — args: handle → float
- **`LIGHT2D.MAKE`** — args: (none) → handle — DEPRECATED alias of LIGHT2D.CREATE. Use LIGHT2D.CREATE.
- **`LIGHT2D.SETCOLOR`** — args: handle, int, int, int, int
- **`LIGHT2D.SETINTENSITY`** — args: handle, float
- **`LIGHT2D.SETPOS`** — args: handle, float, float
- **`LIGHT2D.SETPOSITION`** — args: handle, float, float — DEPRECATED alias of LIGHT2D.SETPOS. Use LIGHT2D.SETPOS.
- **`LIGHT2D.SETRADIUS`** — args: handle, float

*167 overloads in this section.*

---

## Meshes, models, materials, textures, asset packs

Guide: [03-ASSETS.md](03-ASSETS.md)

### MESH

- **`MESH.CREATECAPSULE`** — args: float, float, int, int
- **`MESH.CREATECONE`** — args: float, float, int
- **`MESH.CREATECUBE`** — args: float, float, float → handle — Alias of MESH.MAKECUBE
- **`MESH.CREATECUBE`** — args: float, float, float
- **`MESH.CREATECUBICMAP`** — args: handle, float, float, float
- **`MESH.CREATECUSTOM`** — args: handle, handle → handle
- **`MESH.CREATECYLINDER`** — args: float, float, int
- **`MESH.CREATEHEIGHTMAP`** — args: handle, float, float, float
- **`MESH.CREATEKNOT`** — args: float, float, int, int
- **`MESH.CREATEPLANE`** — args: float, float, int, int → handle — Alias of MESH.MAKEPLANE â€” procedural plane mesh handle
- **`MESH.CREATEPLANE`** — args: float, float, int, int
- **`MESH.CREATEPOLY`** — args: int, float
- **`MESH.CREATESPHERE`** — args: float, int, int
- **`MESH.CREATESPHERE`** — args: float, int, int → handle — Alias of MESH.MAKESPHERE
- **`MESH.CREATETORUS`** — args: float, float, int, int
- **`MESH.CUBE`** — args: float, float, float
- **`MESH.DRAW`** — args: handle, handle, handle
- **`MESH.DRAWAT`** — args: handle, handle, float, float, float
- **`MESH.DRAWINSTANCED`** — args: handle, handle, handle, int
- **`MESH.DRAWROTATED`** — args: handle, handle, float, float, float
- **`MESH.EXPORT`** — args: handle, string
- **`MESH.FREE`** — args: handle
- **`MESH.GENERATEBOUNDS`** — args: handle
- **`MESH.GENERATELOD`** — args: handle, float, float
- **`MESH.GENERATELODCHAIN`** — args: handle, any
- **`MESH.GENERATENORMALS`** — args: handle
- **`MESH.GENTANGENTS`** — args: handle
- **`MESH.GETBBOXMAXX`** — args: handle
- **`MESH.GETBBOXMAXY`** — args: handle
- **`MESH.GETBBOXMAXZ`** — args: handle
- **`MESH.GETBBOXMINX`** — args: handle
- **`MESH.GETBBOXMINY`** — args: handle
- **`MESH.GETBBOXMINZ`** — args: handle
- **`MESH.GETBOUNDS`** — args: handle → handle
- **`MESH.LOAD`** — args: string → handle
- **`MESH.MAKECAPSULE`** — args: float, float, int, int — DEPRECATED alias of MESH.CREATECAPSULE. Use MESH.CREATECAPSULE.
- **`MESH.MAKECONE`** — args: float, float, int — DEPRECATED alias of MESH.CREATECONE. Use MESH.CREATECONE.
- **`MESH.MAKECUBE`** — args: float, float, float → handle — DEPRECATED alias of MESH.CREATECUBE. Use MESH.CREATECUBE. Alias of MESH.MAKECUBE
- **`MESH.MAKECUBE`** — args: float, float, float — DEPRECATED alias of MESH.CREATECUBE. Use MESH.CREATECUBE.
- **`MESH.MAKECUBICMAP`** — args: handle, float, float, float — DEPRECATED alias of MESH.CREATECUBICMAP. Use MESH.CREATECUBICMAP.
- **`MESH.MAKECUSTOM`** — args: handle, handle → handle — DEPRECATED alias of MESH.CREATECUSTOM. Use MESH.CREATECUSTOM.
- **`MESH.MAKECYLINDER`** — args: float, float, int — DEPRECATED alias of MESH.CREATECYLINDER. Use MESH.CREATECYLINDER.
- **`MESH.MAKEHEIGHTMAP`** — args: handle, float, float, float — DEPRECATED alias of MESH.CREATEHEIGHTMAP. Use MESH.CREATEHEIGHTMAP.
- **`MESH.MAKEKNOT`** — args: float, float, int, int — DEPRECATED alias of MESH.CREATEKNOT. Use MESH.CREATEKNOT.
- **`MESH.MAKEPLANE`** — args: float, float, int, int → handle — DEPRECATED alias of MESH.CREATEPLANE. Use MESH.CREATEPLANE. Alias of MESH.MAKEPLANE â€” procedural plane mesh handle
- **`MESH.MAKEPLANE`** — args: float, float, int, int — DEPRECATED alias of MESH.CREATEPLANE. Use MESH.CREATEPLANE.
- **`MESH.MAKEPOLY`** — args: int, float — DEPRECATED alias of MESH.CREATEPOLY. Use MESH.CREATEPOLY.
- **`MESH.MAKESPHERE`** — args: float, int, int — DEPRECATED alias of MESH.CREATESPHERE. Use MESH.CREATESPHERE.
- **`MESH.MAKESPHERE`** — args: float, int, int → handle — DEPRECATED alias of MESH.CREATESPHERE. Use MESH.CREATESPHERE. Alias of MESH.MAKESPHERE
- **`MESH.MAKETORUS`** — args: float, float, int, int — DEPRECATED alias of MESH.CREATETORUS. Use MESH.CREATETORUS.
- **`MESH.OPTIMISEALL`** — args: handle
- **`MESH.OPTIMISEFETCH`** — args: handle
- **`MESH.OPTIMISEOVERDRAW`** — args: handle, float
- **`MESH.OPTIMISEVERTEXCACHE`** — args: handle
- **`MESH.OPTIMIZEALL`** — args: handle
- **`MESH.OPTIMIZEFETCH`** — args: handle
- **`MESH.OPTIMIZEOVERDRAW`** — args: handle, float
- **`MESH.OPTIMIZEVERTEXCACHE`** — args: handle
- **`MESH.PLANE`** — args: float, float, int, int
- **`MESH.SPHERE`** — args: float, int, int
- **`MESH.TRIANGLECOUNT`** — args: handle → int
- **`MESH.UPDATEVERTEX`** — args: handle, int, float, float, float, float, float, float, float, float
- **`MESH.UPDATEVERTICES`** — args: handle, handle
- **`MESH.UPLOAD`** — args: handle, bool
- **`MESH.VERTEXCOUNT`** — args: handle → int

### MODEL

- **`INSTANCE.GETALPHA`** — args: handle → float — Get instance 0 alpha.
- **`INSTANCE.GETCOLOR`** — args: handle → handle — Get instance 0 color handle. (Returns Color instance handle)
- **`MODEL.ADDCHILD`** — args: handle, handle
- **`MODEL.ALPHA`** — args: handle → float — Property alias for MODEL.GETALPHA
- **`MODEL.ANIMCOUNT`** — args: handle → int
- **`MODEL.ANIMDONE`** — args: handle → bool
- **`MODEL.ANIMNAME`** — args: handle, int → string
- **`MODEL.ATTACHTO`** — args: handle, handle
- **`MODEL.CHILDCOUNT`** — args: handle → int
- **`MODEL.CLONE`** — args: handle
- **`MODEL.COLOR`** — args: handle → handle — Property alias for MODEL.GETCOLOR
- **`MODEL.CREATE`** — args: handle → handle
- **`MODEL.CREATEBOX`** — args: float, float, float → handle
- **`MODEL.CREATEBOX`** — args: float, float, float, bool → handle
- **`MODEL.CREATECAPSULE`** — args: float, float → handle — EntityRef capsule primitive (radius, height); draw matches Jolt capsule when using ENTITY.ADDPHYSICS capsule
- **`MODEL.CREATEINSTANCED`** — args: string, int → handle
- **`MODEL.DETACH`** — args: handle
- **`MODEL.DRAW`** — args: handle → handle
- **`MODEL.DRAWAT`** — args: handle, float, float, float, float, float, float, float, float, float
- **`MODEL.DRAWEX`** — args: handle, float, float, float, float, float, float, float, float, float, float, int, int, int, int
- **`MODEL.DRAWWIRES`** — args: handle, int, int, int, int
- **`MODEL.EXISTS`** — args: handle
- **`MODEL.FREE`** — args: handle
- **`MODEL.GETALPHA`** — args: handle → float
- **`MODEL.GETALPHA`** — args: handle → float — Get model alpha (0..1).
- **`MODEL.GETCHILD`** — args: handle, int → handle
- **`MODEL.GETCOLOR`** — args: handle → handle
- **`MODEL.GETCOLOR`** — args: handle → handle — Get model color as a Color instance handle. (Returns Color instance handle)
- **`MODEL.GETFRAME`** — args: handle → int
- **`MODEL.GETMATERIALCOUNT`** — args: handle
- **`MODEL.GETPARENT`** — args: handle → handle
- **`MODEL.GETPOS`** — args: handle → handle
- **`MODEL.GETROT`** — args: handle → handle
- **`MODEL.GETSCALE`** — args: handle → handle
- **`MODEL.HIDE`** — args: handle
- **`MODEL.INSTANCE`** — args: handle
- **`MODEL.ISLOADED`** — args: handle → bool
- **`MODEL.ISPLAYING`** — args: handle → bool
- **`MODEL.ISVISIBLE`** — args: handle → bool
- **`MODEL.LIMBCOUNT`** — args: handle → int
- **`MODEL.LIMBX`** — args: handle, int → float
- **`MODEL.LOAD`** — args: string
- **`MODEL.LOADANIMATIONS`** — args: handle, string
- **`MODEL.LOADASYNC`** — args: string → handle
- **`MODEL.LOADLOD`** — args: string, string, string → handle
- **`MODEL.LOOP`** — args: handle, bool
- **`MODEL.MAKE`** — args: handle → handle — DEPRECATED alias of MODEL.CREATE. Use MODEL.CREATE.
- **`MODEL.MAKEBOX`** — args: float, float, float → handle — DEPRECATED alias of MODEL.CREATEBOX. Use MODEL.CREATEBOX.
- **`MODEL.MAKEBOX`** — args: float, float, float, bool → handle — DEPRECATED alias of MODEL.CREATEBOX. Use MODEL.CREATEBOX.
- **`MODEL.MAKECAPSULE`** — args: float, float → handle — DEPRECATED alias of MODEL.CREATECAPSULE. Use MODEL.CREATECAPSULE. EntityRef capsule primitive (radius, height); draw matches Jolt capsule when using ENTITY.ADDPHYSICS capsule
- **`MODEL.MAKEINSTANCED`** — args: string, int → handle — DEPRECATED alias of MODEL.CREATEINSTANCED. Use MODEL.CREATEINSTANCED.
- **`MODEL.MOVE`** — args: handle, float, float, float
- **`MODEL.PLAY`** — args: handle, string → handle
- **`MODEL.PLAYIDX`** — args: handle, int
- **`MODEL.POS`** — args: handle → handle — Property alias for MODEL.GETPOS
- **`MODEL.REMOVECHILD`** — args: handle, handle
- **`MODEL.ROT`** — args: handle → handle — Property alias for MODEL.GETROT
- **`MODEL.ROTATE`** — args: handle, float, float, float
- **`MODEL.ROTATETEXTURE`** — args: handle, float
- **`MODEL.SCALE`** — args: handle → handle — Property alias for MODEL.GETSCALE
- **`MODEL.SCALETEXTURE`** — args: handle, float, float
- **`MODEL.SCROLLTEXTURE`** — args: handle, float, float
- **`MODEL.SETALPHA`** — args: handle, int → handle
- **`MODEL.SETAMBIENTCOLOR`** — args: handle, int, int, int → handle
- **`MODEL.SETBLEND`** — args: handle, int → handle
- **`MODEL.SETCASTSHADOW`** — args: handle, bool → handle
- **`MODEL.SETCOLOR`** — args: handle, int, int, int, int → handle
- **`MODEL.SETCULL`** — args: handle, bool → handle
- **`MODEL.SETDEPTH`** — args: handle, int → handle
- **`MODEL.SETDIFFUSE`** — args: handle, int, int, int → handle
- **`MODEL.SETEMISSIVE`** — args: handle, int, int, int → handle
- **`MODEL.SETFOG`** — args: handle, bool → handle
- **`MODEL.SETGPUSKINNING`** — args: handle, bool → handle
- **`MODEL.SETINSTANCEPOS`** — args: handle, int, float, float, float → handle
- **`MODEL.SETINSTANCESCALE`** — args: handle, int, float, float, float → handle
- **`MODEL.SETLIGHTING`** — args: handle, bool → handle
- **`MODEL.SETLIMBPOS`** — args: handle, int, float, float, float → handle
- **`MODEL.SETLODDISTANCES`** — args: handle, float, float, float → handle
- **`MODEL.SETMATERIAL`** — args: handle, int, handle → handle
- **`MODEL.SETMATERIALSHADER`** — args: handle, int, handle → handle
- **`MODEL.SETMATERIALTEXTURE`** — args: handle, int, int, handle → handle
- **`MODEL.SETMATRIX`** — args: handle, handle → handle
- **`MODEL.SETMETAL`** — args: handle, float → handle
- **`MODEL.SETMODELMESHMATERIAL`** — args: handle, int, int → handle
- **`MODEL.SETPOS`** — args: handle, float, float, float → handle
- **`MODEL.SETPOSITION`** — args: handle, float, float, float → handle — DEPRECATED alias of MODEL.SETPOS. Use MODEL.SETPOS.
- **`MODEL.SETRECEIVESHADOW`** — args: handle, bool → handle
- **`MODEL.SETROT`** — args: handle, float, float, float → handle
- **`MODEL.SETROUGH`** — args: handle, float → handle
- **`MODEL.SETSCALE`** — args: handle, float, float, float → handle
- **`MODEL.SETSCALEUNIFORM`** — args: handle, float → handle
- **`MODEL.SETSPECULAR`** — args: handle, int, int, int → handle
- **`MODEL.SETSPECULARPOW`** — args: handle, float → handle
- **`MODEL.SETSPEED`** — args: handle, float → handle
- **`MODEL.SETSTAGEBLEND`** — args: handle, int, float → handle
- **`MODEL.SETSTAGEROTATE`** — args: handle, int, float → handle
- **`MODEL.SETSTAGESCALE`** — args: handle, int, float, float → handle
- **`MODEL.SETSTAGESCROLL`** — args: handle, int, float, float → handle
- **`MODEL.SETTEXTURESTAGE`** — args: handle, int, handle → handle
- **`MODEL.SETWIREFRAME`** — args: handle, bool → handle
- **`MODEL.SHOW`** — args: handle
- **`MODEL.STOP`** — args: handle → handle
- **`MODEL.TOTALFRAMES`** — args: handle → int
- **`MODEL.UPDATEANIM`** — args: handle, float → handle
- **`MODEL.UPDATEINSTANCES`** — args: handle
- **`MODEL.X`** — args: handle → float
- **`MODEL.Y`** — args: handle → float
- **`MODEL.Z`** — args: handle → float

### MATERIAL

- **`MATERIAL.AUTOFILTER`** — args: any
- **`MATERIAL.BULKASSIGN`** — args: string, handle → int
- **`MATERIAL.CREATE`** — args: (none) → handle
- **`MATERIAL.CREATEDEFAULT`** — args: (none)
- **`MATERIAL.CREATEPBR`** — args: (none) → handle
- **`MATERIAL.FREE`** — args: handle
- **`MATERIAL.MAKE`** — args: (none) → handle — DEPRECATED alias of MATERIAL.CREATE. Use MATERIAL.CREATE.
- **`MATERIAL.MAKEDEFAULT`** — args: (none) — DEPRECATED alias of MATERIAL.CREATEDEFAULT. Use MATERIAL.CREATEDEFAULT.
- **`MATERIAL.MAKEPBR`** — args: (none) → handle — DEPRECATED alias of MATERIAL.CREATEPBR. Use MATERIAL.CREATEPBR.
- **`MATERIAL.SETCOLOR`** — args: handle, int, int, int, int, int
- **`MATERIAL.SETEFFECT`** — args: handle, string
- **`MATERIAL.SETEFFECTPARAM`** — args: handle, string, float
- **`MATERIAL.SETFLOAT`** — args: handle, int, float
- **`MATERIAL.SETSECONDARYTEXTURE`** — args: int, handle — Alias of ENTITY.SETDETAILTEXTURE
- **`MATERIAL.SETSHADER`** — args: handle, handle
- **`MATERIAL.SETTEXTURE`** — args: handle, int, handle
- **`MATERIAL.SETUVSCROLL`** — args: int, float, float — Alias of ENTITY.SCROLLMATERIAL (mesh material 0)

### TEXTURE

- **`TEXTURE.FREE`** — args: handle
- **`TEXTURE.FROMIMAGE`** — args: handle
- **`TEXTURE.GENCHECKED`** — args: int, int, int, int, handle, handle → handle
- **`TEXTURE.GENCOLOR`** — args: int, int, int, int, int, int → handle
- **`TEXTURE.GENGRADIENTH`** — args: int, int, handle, handle → handle
- **`TEXTURE.GENGRADIENTV`** — args: int, int, handle, handle → handle
- **`TEXTURE.GENWHITENOISE`** — args: int, int → handle
- **`TEXTURE.GENWHITENOISE`** — args: int, int, float → handle
- **`TEXTURE.GETHEIGHT`** — args: handle → int — Same as TEXTURE.HEIGHT; handle-chain friendly name.
- **`TEXTURE.GETSIZE`** — args: handle → handle — Texture dimensions as Vec2 (width, height).
- **`TEXTURE.GETWIDTH`** — args: handle → int — Same as TEXTURE.WIDTH; handle-chain friendly name.
- **`TEXTURE.HEIGHT`** — args: handle → int
- **`TEXTURE.ISLOADED`** — args: handle → bool
- **`TEXTURE.LOAD`** — args: string
- **`TEXTURE.LOADANIM`** — args: string, int, int → handle — TEXTURE.LOAD + SETGRID in one call
- **`TEXTURE.LOADASYNC`** — args: string → handle
- **`TEXTURE.PLAY`** — args: handle, float, bool — Auto-advance atlas frames; call TEXTURE.TICKALL each frame
- **`TEXTURE.RELOAD`** — args: handle
- **`TEXTURE.SETDEFAULTFILTER`** — args: int
- **`TEXTURE.SETDISTORTION`** — args: handle, float — Shader-side distortion amount hint
- **`TEXTURE.SETFILTER`** — args: handle, int
- **`TEXTURE.SETFRAME`** — args: handle, int — Select atlas frame index (0-based)
- **`TEXTURE.SETGRID`** — args: handle, int, int — Spritesheet layout: columns x rows of equal frames
- **`TEXTURE.SETUVSCROLL`** — args: handle, float, float — Source-rectangle scroll speeds for sampled UVs
- **`TEXTURE.SETWRAP`** — args: handle, int
- **`TEXTURE.STOPANIM`** — args: handle
- **`TEXTURE.TICKALL`** — args: (none) — Advance all playing atlas animations (optional dt via overload)
- **`TEXTURE.TICKALL`** — args: float
- **`TEXTURE.UPDATE`** — args: handle, handle
- **`TEXTURE.WIDTH`** — args: handle → int

### ASSET

- **`ASSET.LOADPACK`** — args: string
- **`ASSET.MODEL`** — args: string → handle
- **`ASSET.SOUND`** — args: string → handle
- **`ASSET.TEXTURE`** — args: string → handle
- **`ASSET.UNLOAD`** — args: (none)

### BBOX

- **`BBOX.CHECK`** — args: handle, handle → bool
- **`BBOX.CHECKSPHERE`** — args: handle, float, float, float, float → bool
- **`BBOX.CREATE`** — args: float, float, float, float, float, float → handle
- **`BBOX.FREE`** — args: handle
- **`BBOX.FROMMODEL`** — args: handle → handle
- **`BBOX.MAKE`** — args: float, float, float, float, float, float → handle — DEPRECATED alias of BBOX.CREATE. Use BBOX.CREATE.

### BSPHERE

- **`BSPHERE.CHECK`** — args: handle, handle → bool
- **`BSPHERE.CHECKBOX`** — args: handle, handle → bool
- **`BSPHERE.CREATE`** — args: float, float, float, float → handle
- **`BSPHERE.FREE`** — args: handle
- **`BSPHERE.MAKE`** — args: float, float, float, float → handle — DEPRECATED alias of BSPHERE.CREATE. Use BSPHERE.CREATE.

*236 overloads in this section.*

---

## Input and actions

Guide: [04-INPUT.md](04-INPUT.md)

### INPUT

- **`INPUT`** — args: string → string
- **`INPUT.ACTIONAXIS`** — args: string → float
- **`INPUT.ACTIONDOWN`** — args: string → bool
- **`INPUT.ACTIONPRESSED`** — args: string → bool
- **`INPUT.ACTIONRELEASED`** — args: string → bool
- **`INPUT.AXIS`** — args: any, any → float — Two-key axis: -1, 0, or 1 (negKey vs posKey)
- **`INPUT.AXISDEG`** — args: any, any, float, float → float — Input.Axis(neg,pos)*DEGPERSEC(degPerSec,dt) â€” radians this frame
- **`INPUT.CHARPRESSED`** — args: (none) → int
- **`INPUT.GAMEPADAXIS`** — args: int, int → float
- **`INPUT.GAMEPADAXISCOUNT`** — args: int → int
- **`INPUT.GAMEPADBUTTONCOUNT`** — args: int → int
- **`INPUT.GAMEPADBUTTONDOWN`** — args: int, int → bool
- **`INPUT.GETGAMEPADAXISVALUE`** — args: int, int → float
- **`INPUT.GETINACTIVITY`** — args: (none) → float — Returns time in seconds since the last user interaction.
- **`INPUT.GETKEYNAME`** — args: int → string
- **`INPUT.GETMOUSEWORLDPOS`** — args: handle, int, int → handle
- **`INPUT.GETTOUCHPOINTID`** — args: int → int
- **`INPUT.ISGAMEPADAVAILABLE`** — args: int → bool
- **`INPUT.JOYBUTTON`** — args: int → bool
- **`INPUT.JOYDOWN`** — args: any, any → bool
- **`INPUT.JOYX`** — args: (none) → float
- **`INPUT.JOYY`** — args: (none) → float
- **`INPUT.KEYDOWN`** — args: any
- **`INPUT.KEYDOWN`** — args: int → bool
- **`INPUT.KEYHIT`** — args: any → bool
- **`INPUT.KEYPRESSED`** — args: any
- **`INPUT.KEYUP`** — args: any
- **`INPUT.KEYUP`** — args: int → bool
- **`INPUT.LOADMAPPINGS`** — args: string
- **`INPUT.LOCKMOUSE`** — args: bool
- **`INPUT.MAPGAMEPADAXIS`** — args: string, int, int
- **`INPUT.MAPGAMEPADBUTTON`** — args: string, int, int
- **`INPUT.MAPKEY`** — args: string, int
- **`INPUT.MOUSEDELTA`** — args: (none) → handle
- **`INPUT.MOUSEDELTAX`** — args: (none) → float
- **`INPUT.MOUSEDELTAY`** — args: (none) → float
- **`INPUT.MOUSEDELTA_X`** — args: (none) → float
- **`INPUT.MOUSEDELTA_Y`** — args: (none) → float
- **`INPUT.MOUSEDOWN`** — args: int
- **`INPUT.MOUSEDX`** — args: (none) → float — Alias of INPUT.MOUSEDELTAX
- **`INPUT.MOUSEDY`** — args: (none) → float — Alias of INPUT.MOUSEDELTAY
- **`INPUT.MOUSEHIT`** — args: int → bool
- **`INPUT.MOUSEPRESSED`** — args: int → bool
- **`INPUT.MOUSERELEASED`** — args: int → bool
- **`INPUT.MOUSEWHEEL`** — args: (none) → float — Alias of INPUT.MOUSEWHEELMOVE
- **`INPUT.MOUSEWHEELMOVE`** — args: (none) → float
- **`INPUT.MOUSEX`** — args: (none)
- **`INPUT.MOUSEXSPEED`** — args: (none) → float
- **`INPUT.MOUSEY`** — args: (none)
- **`INPUT.MOUSEYSPEED`** — args: (none) → float
- **`INPUT.MOVEDIR`** — args: float, float → handle
- **`INPUT.MOVEMENT2D`** — args: any, any, any, any → handle — 2-float array [forward, strafe] from two Axis pairs; ERASE when done
- **`INPUT.ORBIT`** — args: any, any, float, float → float — Alias of INPUT.AXISDEG â€” orbit / yaw delta this frame
- **`INPUT.SAVEMAPPINGS`** — args: string
- **`INPUT.SETGAMEPADMAPPINGS`** — args: string → int
- **`INPUT.SETMOUSEOFFSET`** — args: int, int
- **`INPUT.SETMOUSEPOS`** — args: int, int — Warp OS cursor to client pixel (x,y); pair with CURSOR.DISABLE for game-style recenter
- **`INPUT.SETMOUSESCALE`** — args: float, float
- **`INPUT.TOUCHCOUNT`** — args: (none) → int
- **`INPUT.TOUCHPRESSED`** — args: int → bool
- **`INPUT.TOUCHX`** — args: int → int
- **`INPUT.TOUCHY`** — args: int → int

### ACTION

- **`ACTION.BINDGAMEPAD`** — args: string, int, int
- **`ACTION.BINDKEY`** — args: string, int
- **`ACTION.DOWN`** — args: string → bool
- **`ACTION.HIT`** — args: string → bool
- **`ACTION.MAPAXIS`** — args: string, int, int
- **`ACTION.MAPJOY`** — args: string, int, int
- **`ACTION.MAPKEY`** — args: string, int
- **`ACTION.MAPMOUSE`** — args: string, int
- **`ACTION.PRESSED`** — args: string → bool
- **`ACTION.RELEASED`** — args: string → bool
- **`ACTION.RESET`** — args: (none)
- **`ACTION.VALUE`** — args: string → float

### GAMEPAD

- **`GAMEPAD`** — args: (none) → bool
- **`GAMEPAD.AXIS`** — args: int, int → float
- **`GAMEPAD.BUTTON`** — args: int, int → bool

### CURSOR

- **`CURSOR.DISABLE`** — args: (none)
- **`CURSOR.ENABLE`** — args: (none)
- **`CURSOR.HIDE`** — args: (none)
- **`CURSOR.ISENABLED`** — args: (none) → bool
- **`CURSOR.ISHIDDEN`** — args: (none)
- **`CURSOR.ISONSCREEN`** — args: (none)
- **`CURSOR.SET`** — args: int
- **`CURSOR.SHOW`** — args: (none)

### GESTURE

- **`GESTURE.ENABLE`** — args: int
- **`GESTURE.GETDETECTED`** — args: (none)
- **`GESTURE.GETDRAGANGLE`** — args: (none)
- **`GESTURE.GETDRAGVECTORX`** — args: (none)
- **`GESTURE.GETDRAGVECTORY`** — args: (none)
- **`GESTURE.GETHOLDDURATION`** — args: (none)
- **`GESTURE.GETPINCHANGLE`** — args: (none)
- **`GESTURE.GETPINCHVECTORX`** — args: (none)
- **`GESTURE.GETPINCHVECTORY`** — args: (none)
- **`GESTURE.ISDETECTED`** — args: int

*95 overloads in this section.*

---

## Physics, bodies, collision, picking

Guide: [05-PHYSICS.md](05-PHYSICS.md)

### PHYSICS

- **`PHYSICS.AUTOCREATE`** — args: int
- **`PHYSICS.BOXCAST`** — args: any
- **`PHYSICS.DISABLE`** — args: any
- **`PHYSICS.ENABLE`** — args: any
- **`PHYSICS.EXPLOSION`** — args: float, float, float, float, float — Applies physical impulse radially.
- **`PHYSICS.GETBUOYANCY`** — args: int → float — Reads stored buoyancy density (default 0)
- **`PHYSICS.RAYCAST`** — args: float, float, float, float, float, float, float → handle
- **`PHYSICS.SETBUOYANCY`** — args: int, float — Stores per-entity buoyancy density for future Jolt/WASM fluid coupling (gameplay hint today)
- **`PHYSICS.SETGRAVITY`** — args: float, float, float
- **`PHYSICS.SETSUBSTEPS`** — args: int
- **`PHYSICS.SPHERECAST`** — args: any
- **`PHYSICS.START`** — args: (none)
- **`PHYSICS.STEP`** — args: float
- **`PHYSICS.STOP`** — args: (none)
- **`PHYSICS.TORQUE`** — args: handle, float, float, float

### PHYSICS3D

- **`AERO.SETDRAG`** — args: handle, float — Apply air resistance coefficient.
- **`AERO.SETLIFT`** — args: handle, float — Set lift coefficient for a physics body.
- **`AERO.SETTHRUST`** — args: handle, float — Apply local Z-axis thrust power.
- **`BODY3D.APPLYTORQUE`** — args: handle, float, float, float
- **`BODY3D.GETANGULARVEL`** — args: handle → handle
- **`BODY3D.GETLINEARVEL`** — args: handle → handle — Get linear velocity as a 3-element array (alias of BODY3D.GETVELOCITY).
- **`BODY3D.GETMASS`** — args: handle → float
- **`BODY3D.GETVELOCITY`** — args: handle → handle — Get linear velocity as a 3-element numeric array.
- **`BODY3D.SETVELOCITY`** — args: handle, float, float, float — Set linear velocity.
- **`BODYREF.FREE`** — args: handle — Destroy a physics body (handle method).
- **`BODYREF.GETPOSITION`** — args: handle → handle
- **`BODYREF.GETROTATION`** — args: handle → handle
- **`BODYREF.GETVELOCITY`** — args: handle → handle
- **`BODYREF.SETVELOCITY`** — args: handle, float, float, float
- **`DEBUG.DRAWBODY`** — args: any — Debug draw a physics body wireframe (no-op on stub builds).
- **`DEBUG.DRAWCHARACTER`** — args: any — Debug draw a KCC capsule (no-op on stub builds).
- **`JOINT.CREATEHINGE`** — args: handle, handle, float, float, float, float, float, float → handle — Creates a hinge joint between two bodies at (px,py,pz) around axis (ax,ay,az).
- **`JOINT.CREATEPOINT`** — args: handle, handle, float, float, float → handle — Creates a point-to-point (ball and socket) joint between two bodies at (px,py,pz).
- **`JOINT.FREE`** — args: handle — Destroys a physics joint/constraint.
- **`JOINT.MAKEHINGE`** — args: handle, handle, float, float, float, float, float, float → handle — DEPRECATED alias of JOINT.CREATEHINGE. Use JOINT.CREATEHINGE. Creates a hinge joint between two bodies at (px,py,pz) around axis (ax,ay,az).
- **`JOINT.MAKEPOINT`** — args: handle, handle, float, float, float → handle — DEPRECATED alias of JOINT.CREATEPOINT. Use JOINT.CREATEPOINT. Creates a point-to-point (ball and socket) joint between two bodies at (px,py,pz).
- **`PHYSICS.AUTO`** — args: int, string, float — Alias for ENTITY.PHYSICS.
- **`PHYSICS.BOUNCE`** — args: int, float — Modular building: Sets bounciness (restitution) for a pending physics body.
- **`PHYSICS.BUILD`** — args: int, float — Modular building: Finalizes and commits the physics body with given mass.
- **`PHYSICS.FORCE`** — args: int, float, float, float — Entity-First: Applies a continuous force to an entity's physics body.
- **`PHYSICS.FRICTION`** — args: int, float — Modular building: Sets friction for a pending physics body.
- **`PHYSICS.GETGRAVITYX`** — args: (none) → float — Get gravity X (alias of PHYSICS3D.GETGRAVITYX).
- **`PHYSICS.GETGRAVITYY`** — args: (none) → float — Get gravity Y (alias of PHYSICS3D.GETGRAVITYY).
- **`PHYSICS.GETGRAVITYZ`** — args: (none) → float — Get gravity Z (alias of PHYSICS3D.GETGRAVITYZ).
- **`PHYSICS.GRAVITY`** — args: int, float — Entity-First: Scale the gravity factor for a specific entity (e.g. 0.0 for zero-g).
- **`PHYSICS.IMPULSE`** — args: int, float, float, float — Entity-First: Applies an instant impulse to an entity's physics body.
- **`PHYSICS.SETROT`** — args: int, float, float, float — Entity-First: Instantly sets the rotation of an entity's physics body (Euler radians).
- **`PHYSICS.SHAPE`** — args: int, string — Modular building: Sets the physics shape for a pending body.
- **`PHYSICS.SIZE`** — args: int, float, float, float — Modular building: Sets dimensions for a pending physics shape.
- **`PHYSICS.VELOCITY`** — args: int, float, float, float — Entity-First: Sets the linear velocity of an entity's physics body.
- **`PHYSICS.WAKE`** — args: int — Entity-First: Forces a sleeping physics body to wake up.
- **`PHYSICS3D.DEBUGDRAW`** — args: int
- **`PHYSICS3D.GETGRAVITYX`** — args: (none) → float — Get current gravity X component.
- **`PHYSICS3D.GETGRAVITYY`** — args: (none) → float — Get current gravity Y component.
- **`PHYSICS3D.GETGRAVITYZ`** — args: (none) → float — Get current gravity Z component.
- **`PHYSICS3D.GETMATRIXBUFFER`** — args: (none) → handle — Get the shared matrix buffer for render interpolation.
- **`PHYSICS3D.GETSCRATCHFLOAT`** — args: int → float — Read a scratch float from the physics scratch buffer.
- **`PHYSICS3D.MOUSEHIT`** — args: handle → handle — Raycast from mouse through camera; returns [x,y,z] array or nil.
- **`PHYSICS3D.ONCOLLISION`** — args: handle, handle, string
- **`PHYSICS3D.PROCESSCOLLISIONS`** — args: (none)
- **`PHYSICS3D.RAYCAST`** — args: float, float, float, float, float, float, float → handle
- **`PHYSICS3D.SETGRAVITY`** — args: float, float, float
- **`PHYSICS3D.SETSUBSTEPS`** — args: int
- **`PHYSICS3D.SETTIMESTEP`** — args: float — Set the fixed physics simulation timestep (e.g. 60.0, 90.0).
- **`PHYSICS3D.START`** — args: (none)
- **`PHYSICS3D.STEP`** — args: (none)
- **`PHYSICS3D.STOP`** — args: (none)
- **`PHYSICS3D.SYNCWASMTOPHYSREGS`** — args: int, int — Sync WASM physics view to VM registers (count, firstReg).
- **`PHYSICS3D.UPDATE`** — args: (none) — Advance the 3D physics simulation (same implementation as PHYSICS3D.STEP; optional frame dt like STEP)
- **`SHAPE.GETDEPTH`** — args: handle → float — Get shape dimension 3 (half-extent Z).
- **`SHAPE.GETHEIGHT`** — args: handle → float — Get shape dimension 2 (half-extent Y or height).
- **`SHAPE.GETRADIUS`** — args: handle → float — Get shape radius (same as SHAPE.GETWIDTH for spheres).
- **`SHAPE.GETSIZEX`** — args: handle → float — Get shape X dimension.
- **`SHAPE.GETSIZEY`** — args: handle → float — Get shape Y dimension.
- **`SHAPE.GETSIZEZ`** — args: handle → float — Get shape Z dimension.
- **`SHAPE.GETTYPE`** — args: handle → int — Get the shape type (1=Box, 2=Sphere, 3=Capsule, 4=Cylinder).
- **`SHAPE.GETWIDTH`** — args: handle → float — Get shape dimension 1 (half-extent X or radius).
- **`SHAPEREF.FREE`** — args: handle — Destroy a collision shape.
- **`VEHICLE.CONTROL`** — args: int, float, float, float — Update all vehicle inputs (vid, throttle, steer, brake).
- **`VEHICLE.CREATE`** — args: int, int → int — Create a vehicle controller for an entity.
- **`VEHICLE.MAKE`** — args: int, int → int — DEPRECATED alias of VEHICLE.CREATE. Use VEHICLE.CREATE. Create a vehicle controller for an entity.
- **`VEHICLE.SETSTEER`** — args: int, float — Set vehicle steering (-1 to 1).
- **`VEHICLE.SETTHROTTLE`** — args: int, float — Set vehicle throttle (-1 to 1).
- **`VEHICLE.SETTUNING`** — args: int, float, float, float, float — Tune suspension (vid, spring, damp, maxSpeed, steerSpeed).
- **`VEHICLE.SETWHEEL`** — args: int, int, float, float, float, float — Configure a wheel (vid, idx, ox, oy, oz, radius).
- **`VEHICLE.STEP`** — args: float — Step all vehicle simulations by dt.
- **`VEHICLE.WHEELX`** — args: int, int → float
- **`VEHICLE.WHEELY`** — args: int, int → float
- **`VEHICLE.WHEELZ`** — args: int, int → float
- **`WORLD.SETUP`** — args: (none) — Initialise physics world with default gravity (-9.81).
- **`WORLD.SETUP`** — args: float — Initialise physics world with custom Y gravity.

### BODY

- **`BODY.ADDCAPSULE`** — args: handle, float, float
- **`BODY.ADDDYNAMICBOX`** — args: handle, float, float, float
- **`BODY.ADDSPHERE`** — args: handle, float
- **`BODY.ADDSTATICBOX`** — args: handle, float, float, float
- **`BODY.APPLYFORCE`** — args: handle, float, float, float
- **`BODY.APPLYIMPULSE`** — args: handle, float, float, float
- **`BODY.SETBOUNCE`** — args: handle, float
- **`BODY.SETFRICTION`** — args: handle, float
- **`BODY.SETMASS`** — args: handle, float
- **`BODY3D.LOCKAXIS`** — args: handle, int — Lock motion/rotation axes (flags: 1=X, 2=Y, 4=Z, 8=RotX, 16=RotY, 32=RotZ).
- **`BODY3D.SETCCD`** — args: handle, bool → handle — Enable/disable Continuous Collision Detection.
- **`BODY3D.SETDAMPING`** — args: handle, float, float → handle — Set linear and angular damping.
- **`BODY3D.SETGRAVITYFACTOR`** — args: handle, float → handle — Set gravity multiplier (0.0 = weightless).

### BODY2D

- **`BODY2D.ADDCIRCLE`** — args: handle, float
- **`BODY2D.ADDPOLYGON`** — args: handle, handle
- **`BODY2D.ADDRECT`** — args: handle, float, float
- **`BODY2D.APPLYFORCE`** — args: handle, float, float
- **`BODY2D.APPLYIMPULSE`** — args: handle, float, float
- **`BODY2D.COLLIDED`** — args: handle → int
- **`BODY2D.COLLISIONNORMAL`** — args: handle → handle
- **`BODY2D.COLLISIONOTHER`** — args: handle → handle
- **`BODY2D.COLLISIONPOINT`** — args: handle → handle
- **`BODY2D.COMMIT`** — args: handle, float, float → handle
- **`BODY2D.CREATE`** — args: string → handle
- **`BODY2D.FREE`** — args: handle
- **`BODY2D.GETANGULARVELOCITY`** — args: handle → float
- **`BODY2D.GETFRICTION`** — args: handle → float
- **`BODY2D.GETLINEARVELOCITY`** — args: handle → handle
- **`BODY2D.GETMASS`** — args: handle → float
- **`BODY2D.GETPOS`** — args: handle → handle
- **`BODY2D.GETRESTITUTION`** — args: handle → float
- **`BODY2D.GETROT`** — args: handle → float
- **`BODY2D.MAKE`** — args: string → handle — DEPRECATED alias of BODY2D.CREATE. Use BODY2D.CREATE.
- **`BODY2D.ROT`** — args: handle → float
- **`BODY2D.SETANGULARVELOCITY`** — args: handle, float
- **`BODY2D.SETFRICTION`** — args: handle, float
- **`BODY2D.SETLINEARVELOCITY`** — args: handle, float, float
- **`BODY2D.SETMASS`** — args: handle, float
- **`BODY2D.SETPOS`** — args: handle, float, float
- **`BODY2D.SETPOSITION`** — args: handle, float, float — DEPRECATED alias of BODY2D.SETPOS. Use BODY2D.SETPOS.
- **`BODY2D.SETRESTITUTION`** — args: handle, float
- **`BODY2D.SETROT`** — args: handle, float
- **`BODY2D.X`** — args: handle → float
- **`BODY2D.Y`** — args: handle → float

### BODY3D

- **`BODY3D.ACTIVATE`** — args: handle
- **`BODY3D.ADDBOX`** — args: handle, float, float, float
- **`BODY3D.ADDCAPSULE`** — args: handle, float, float
- **`BODY3D.ADDMESH`** — args: handle, handle
- **`BODY3D.ADDSPHERE`** — args: handle, float
- **`BODY3D.ANGULARVEL`** — args: handle → handle
- **`BODY3D.APPLYFORCE`** — args: handle, float, float, float
- **`BODY3D.APPLYIMPULSE`** — args: handle, float, float, float
- **`BODY3D.BOUNCE`** — args: handle → float
- **`BODY3D.BUFFERINDEX`** — args: handle → int
- **`BODY3D.COLLIDED`** — args: handle → int
- **`BODY3D.COLLISIONNORMAL`** — args: handle → handle
- **`BODY3D.COLLISIONOTHER`** — args: handle → handle
- **`BODY3D.COLLISIONPOINT`** — args: handle → handle
- **`BODY3D.COMMIT`** — args: handle, float, float, float → handle
- **`BODY3D.CREATE`** — args: (none) → handle
- **`BODY3D.CREATE`** — args: string → handle
- **`BODY3D.CREATE`** — args: string
- **`BODY3D.DEACTIVATE`** — args: handle
- **`BODY3D.FREE`** — args: handle
- **`BODY3D.FRICTION`** — args: handle → float
- **`BODY3D.GETCCD`** — args: handle → bool
- **`BODY3D.GETDAMPING`** — args: handle → handle
- **`BODY3D.GETFRICTION`** — args: handle → float
- **`BODY3D.GETGRAVITYFACTOR`** — args: handle → float
- **`BODY3D.GETPOS`** — args: handle → handle
- **`BODY3D.GETRESTITUTION`** — args: handle → float
- **`BODY3D.GETROT`** — args: handle → handle
- **`BODY3D.GETSCALE`** — args: handle → handle — Returns [sx,sy,sz] scale factors for primitive bodies (box/sphere/capsule); mesh bodies report 1,1,1
- **`BODY3D.MAKE`** — args: (none) → handle — DEPRECATED alias of BODY3D.CREATE. Use BODY3D.CREATE.
- **`BODY3D.MAKE`** — args: string — DEPRECATED alias of BODY3D.CREATE. Use BODY3D.CREATE.
- **`BODY3D.MAKE`** — args: string → handle — DEPRECATED alias of BODY3D.CREATE. Use BODY3D.CREATE.
- **`BODY3D.MASS`** — args: handle → float
- **`BODY3D.POS`** — args: handle → handle
- **`BODY3D.RESTITUTION`** — args: handle → float
- **`BODY3D.ROT`** — args: handle → handle
- **`BODY3D.SCALE`** — args: handle → handle
- **`BODY3D.SETANGULARVEL`** — args: handle, float, float, float → handle
- **`BODY3D.SETFRICTION`** — args: handle, float → handle
- **`BODY3D.SETLINEARVEL`** — args: handle, float, float, float → handle
- **`BODY3D.SETMASS`** — args: handle, float → handle
- **`BODY3D.SETPOS`** — args: handle, float, float, float → handle
- **`BODY3D.SETPOSITION`** — args: handle, float, float, float — DEPRECATED alias of BODY3D.SETPOS. Use BODY3D.SETPOS.
- **`BODY3D.SETRESTITUTION`** — args: handle, float → handle
- **`BODY3D.SETROT`** — args: handle, float, float, float → handle
- **`BODY3D.SETSCALE`** — args: handle, float, float, float → handle — Scales collision shape for primitive bodies built via ADDBOX/ADDSPHERE/ADDCAPSULE or SHAPE.CREATE*; not supported for mesh (ADDMESH)
- **`BODY3D.VEL`** — args: handle → handle
- **`BODY3D.VELOCITY`** — args: handle → handle
- **`BODY3D.X`** — args: handle → float
- **`BODY3D.Y`** — args: handle → float
- **`BODY3D.Z`** — args: handle → float

### BODYREF

- **`BODYREF.ENABLECOLLISION`** — args: handle, bool — Enables/Disables body participation in physics.
- **`BODYREF.SETLAYER`** — args: handle, int — Sets the Jolt collision layer.
- **`BODYREF.SETPOS`** — args: handle, float, float, float — Moves a Kinematic/Static/Trigger body.
- **`BODYREF.SETPOSITION`** — args: handle, float, float, float — DEPRECATED alias of BODYREF.SETPOS. Use BODYREF.SETPOS.
- **`BODYREF.SETROTATION`** — args: handle, float, float, float — Sets body orientation (Euler degrees).

### COLLISION

- **`BBOX.GETMAX`** — args: handle → handle — Get bounding box max corner as Vec3 handle.
- **`BBOX.GETMIN`** — args: handle → handle — Get bounding box min corner as Vec3 handle.
- **`BBOX.MAX`** — args: handle → handle — Property alias for BBOX.GETMAX.
- **`BBOX.MIN`** — args: handle → handle — Property alias for BBOX.GETMIN.
- **`BBOX.SETMAX`** — args: handle, float, float, float → handle — Set bounding box max corner. Returns handle.
- **`BBOX.SETMIN`** — args: handle, float, float, float → handle — Set bounding box min corner. Returns handle.
- **`BSPHERE.GETPOS`** — args: handle → handle — Get bounding sphere center as Vec3 handle.
- **`BSPHERE.GETRADIUS`** — args: handle → float — Get bounding sphere radius.
- **`BSPHERE.POS`** — args: handle → handle — Property alias for BSPHERE.GETPOS.
- **`BSPHERE.RADIUS`** — args: handle → float — Property alias for BSPHERE.GETRADIUS.
- **`BSPHERE.SETPOS`** — args: handle, float, float, float → handle — Set bounding sphere center. Returns handle.
- **`BSPHERE.SETPOSITION`** — args: handle, float, float, float → handle — DEPRECATED alias of BSPHERE.SETPOS. Use BSPHERE.SETPOS.
- **`BSPHERE.SETRADIUS`** — args: handle, float → handle — Set bounding sphere radius. Returns handle.
- **`COLLISION.AABBOVERLAP3D`** — args: handle, handle, handle, handle → bool — 3D AABB overlap using min/max corners for each box (four VEC3 handles; same math as AABBCOLLIDE).
- **`COLLISION.BOXOVERLAP2D`** — args: handle, handle, handle, handle → bool — 2D AABB overlap using four VEC2 handles: position and size for each box (fewer scalars than BOXCOLLIDE).
- **`COLLISION.CIRCLEBOX2D`** — args: handle, float, handle, handle → bool — 2D circle vs axis-aligned box: center VEC2, radius, box position VEC2, box size VEC2.
- **`COLLISION.CIRCLEOVERLAP2D`** — args: handle, float, handle, float → bool — 2D circle-circle test: two VEC2 centers and two radii (four arguments total).
- **`COLLISION.LINESEGINTERSECT2D`** — args: handle, handle, handle, handle → bool — 2D segment intersection: endpoints of segment A and segment B as VEC2 handles.
- **`COLLISION.POINTINAABB3D`** — args: handle, handle, handle → bool — 3D point in axis-aligned box: point VEC3, box min corner VEC3, box size VEC3.
- **`COLLISION.POINTINBOX2D`** — args: handle, handle, handle → bool — 2D point-in-axis-aligned-box using VEC2 point, box min corner, and box size.
- **`COLLISION.POINTONSEG2D`** — args: handle, handle, handle, float → bool — 2D point-near-segment test: point, segment endpoints, distance threshold (matches POINTONLINE math).
- **`COLLISION.SPHEREBOX3D`** — args: handle, float, handle, handle → bool — 3D sphere vs axis-aligned box: sphere center VEC3, radius, box min corner VEC3, box size VEC3.
- **`COLLISION.SPHEREOVERLAP3D`** — args: handle, float, handle, float → bool — 3D sphere-sphere overlap: two VEC3 centers and two radii.
- **`RAY.DIR`** — args: handle → handle — Property alias for RAY.GETDIR.
- **`RAY.GETDIR`** — args: handle → handle — Get ray direction as Vec3 handle.
- **`RAY.GETPOS`** — args: handle → handle — Get ray origin as Vec3 handle.
- **`RAY.POS`** — args: handle → handle — Property alias for RAY.GETPOS.
- **`RAY.SETDIR`** — args: handle, float, float, float → handle — Set ray direction. Returns ray handle.
- **`RAY.SETPOS`** — args: handle, float, float, float → handle — Set ray origin. Returns ray handle.
- **`RAY.SETPOSITION`** — args: handle, float, float, float → handle — DEPRECATED alias of RAY.SETPOS. Use RAY.SETPOS.

### PICK

- **`PICK.CAST`** — args: (none) → int — Run Jolt raycast from staged params; returns entity or 0
- **`PICK.DIRECTION`** — args: float, float, float — Stage ray direction; length is max travel unless PICK.MAXDIST set
- **`PICK.DIST`** — args: (none) → float — Distance along ray to last hit
- **`PICK.DISTANCE`** — args: (none) → float
- **`PICK.ENTITY`** — args: (none) → int — Entity from last pick (linked BODY3D only)
- **`PICK.FROMCAMERA`** — args: handle, float, float — Stage ray from camera handle and screen pixels (sets default MAXDIST if unset)
- **`PICK.HIT`** — args: (none) → bool — Whether last PICK.CAST / SCREENCAST hit
- **`PICK.LAYERMASK`** — args: int — Bit i accepts ENTITY.COLLISIONLAYER i; 0 accepts all
- **`PICK.MAXDIST`** — args: float — Optional max ray length (normalize direction then scale)
- **`PICK.MOUSE`** — args: handle → bool
- **`PICK.NX`** — args: (none) → float — Last pick surface normal X
- **`PICK.NY`** — args: (none) → float — Last pick surface normal Y
- **`PICK.NZ`** — args: (none) → float — Last pick surface normal Z
- **`PICK.ORIGIN`** — args: float, float, float — Stage ray origin for PICK.CAST (Linux+CGO Jolt)
- **`PICK.RADIUS`** — args: float — Reserved; non-zero returns error until sphere pick exists
- **`PICK.RAY`** — args: float, float, float, float, float, float → bool
- **`PICK.SCREENCAST`** — args: handle, float, float → int — FROMCAMERA then CAST; returns entity or 0
- **`PICK.X`** — args: (none) → float — Last pick hit world X
- **`PICK.Y`** — args: (none) → float — Last pick hit world Y
- **`PICK.Z`** — args: (none) → float — Last pick hit world Z

### RAY

- **`RAY.CREATE`** — args: float, float, float, float, float, float → handle
- **`RAY.FREE`** — args: handle
- **`RAY.HITBOX`** — args: handle, float, float, float, float, float, float → bool
- **`RAY.HITBOX_DISTANCE`** — args: handle, float, float, float, float, float, float → float
- **`RAY.HITBOX_HIT`** — args: handle, float, float, float, float, float, float → bool
- **`RAY.HITBOX_NORMALX`** — args: handle, float, float, float, float, float, float → float
- **`RAY.HITBOX_NORMALY`** — args: handle, float, float, float, float, float, float → float
- **`RAY.HITBOX_NORMALZ`** — args: handle, float, float, float, float, float, float → float
- **`RAY.HITBOX_POINTX`** — args: handle, float, float, float, float, float, float → float
- **`RAY.HITBOX_POINTY`** — args: handle, float, float, float, float, float, float → float
- **`RAY.HITBOX_POINTZ`** — args: handle, float, float, float, float, float, float → float
- **`RAY.HITMESH`** — args: handle, handle, handle → bool
- **`RAY.HITMESH_DISTANCE`** — args: handle, handle, handle → float
- **`RAY.HITMESH_HIT`** — args: handle, handle, handle → bool
- **`RAY.HITMESH_NORMALX`** — args: handle, handle, handle → float
- **`RAY.HITMESH_NORMALY`** — args: handle, handle, handle → float
- **`RAY.HITMESH_NORMALZ`** — args: handle, handle, handle → float
- **`RAY.HITMESH_POINTX`** — args: handle, handle, handle → float
- **`RAY.HITMESH_POINTY`** — args: handle, handle, handle → float
- **`RAY.HITMESH_POINTZ`** — args: handle, handle, handle → float
- **`RAY.HITMODEL`** — args: handle, handle → bool
- **`RAY.HITMODEL_DISTANCE`** — args: handle, handle → float
- **`RAY.HITMODEL_HIT`** — args: handle, handle → bool
- **`RAY.HITMODEL_NORMALX`** — args: handle, handle → float
- **`RAY.HITMODEL_NORMALY`** — args: handle, handle → float
- **`RAY.HITMODEL_NORMALZ`** — args: handle, handle → float
- **`RAY.HITMODEL_POINTX`** — args: handle, handle → float
- **`RAY.HITMODEL_POINTY`** — args: handle, handle → float
- **`RAY.HITMODEL_POINTZ`** — args: handle, handle → float
- **`RAY.HITPLANE`** — args: handle, float, float, float, float → bool
- **`RAY.HITPLANE_DISTANCE`** — args: handle, float, float, float, float → float
- **`RAY.HITPLANE_HIT`** — args: handle, float, float, float, float → bool
- **`RAY.HITPLANE_NORMALX`** — args: handle, float, float, float, float → float
- **`RAY.HITPLANE_NORMALY`** — args: handle, float, float, float, float → float
- **`RAY.HITPLANE_NORMALZ`** — args: handle, float, float, float, float → float
- **`RAY.HITPLANE_POINTX`** — args: handle, float, float, float, float → float
- **`RAY.HITPLANE_POINTY`** — args: handle, float, float, float, float → float
- **`RAY.HITPLANE_POINTZ`** — args: handle, float, float, float, float → float
- **`RAY.HITSPHERE`** — args: handle, float, float, float, float → bool
- **`RAY.HITSPHERE_DISTANCE`** — args: handle, float, float, float, float → float
- **`RAY.HITSPHERE_HIT`** — args: handle, float, float, float, float → bool
- **`RAY.HITSPHERE_NORMALX`** — args: handle, float, float, float, float → float
- **`RAY.HITSPHERE_NORMALY`** — args: handle, float, float, float, float → float
- **`RAY.HITSPHERE_NORMALZ`** — args: handle, float, float, float, float → float
- **`RAY.HITSPHERE_POINTX`** — args: handle, float, float, float, float → float
- **`RAY.HITSPHERE_POINTY`** — args: handle, float, float, float, float → float
- **`RAY.HITSPHERE_POINTZ`** — args: handle, float, float, float, float → float
- **`RAY.HITTRIANGLE`** — args: handle, float, float, float, float, float, float, float, float, float → bool
- **`RAY.HITTRIANGLE_DISTANCE`** — args: handle, float, float, float, float, float, float, float, float, float → float
- **`RAY.HITTRIANGLE_HIT`** — args: handle, float, float, float, float, float, float, float, float, float → bool
- **`RAY.HITTRIANGLE_NORMALX`** — args: handle, float, float, float, float, float, float, float, float, float → float
- **`RAY.HITTRIANGLE_NORMALY`** — args: handle, float, float, float, float, float, float, float, float, float → float
- **`RAY.HITTRIANGLE_NORMALZ`** — args: handle, float, float, float, float, float, float, float, float, float → float
- **`RAY.HITTRIANGLE_POINTX`** — args: handle, float, float, float, float, float, float, float, float, float → float
- **`RAY.HITTRIANGLE_POINTY`** — args: handle, float, float, float, float, float, float, float, float, float → float
- **`RAY.HITTRIANGLE_POINTZ`** — args: handle, float, float, float, float, float, float, float, float, float → float
- **`RAY.INTERSECTSMODEL`** — args: handle, handle → bool — Alias of RAY.HITMODEL
- **`RAY.INTERSECTSMODEL_DISTANCE`** — args: handle, handle → float
- **`RAY.INTERSECTSMODEL_HIT`** — args: handle, handle → bool
- **`RAY.INTERSECTSMODEL_NORMALX`** — args: handle, handle → float
- **`RAY.INTERSECTSMODEL_NORMALY`** — args: handle, handle → float
- **`RAY.INTERSECTSMODEL_NORMALZ`** — args: handle, handle → float
- **`RAY.INTERSECTSMODEL_POINTX`** — args: handle, handle → float
- **`RAY.INTERSECTSMODEL_POINTY`** — args: handle, handle → float
- **`RAY.INTERSECTSMODEL_POINTZ`** — args: handle, handle → float
- **`RAY.MAKE`** — args: float, float, float, float, float, float → handle — DEPRECATED alias of RAY.CREATE. Use RAY.CREATE.

### RAY2D

- **`RAY2D.HITCIRCLE_DISTANCE`** — args: float, float, float, float, float, float, float → float — Distance along ray to hit (0 if miss)
- **`RAY2D.HITCIRCLE_HIT`** — args: float, float, float, float, float, float, float → bool — 2D ray vs circle â€” hit
- **`RAY2D.HITCIRCLE_POINTX`** — args: float, float, float, float, float, float, float → float
- **`RAY2D.HITCIRCLE_POINTY`** — args: float, float, float, float, float, float, float → float
- **`RAY2D.HITRECT_DISTANCE`** — args: float, float, float, float, float, float, float, float → float
- **`RAY2D.HITRECT_HIT`** — args: float, float, float, float, float, float, float, float → bool — 2D ray vs axis-aligned rect (minx,miny,maxx,maxy)
- **`RAY2D.HITRECT_POINTX`** — args: float, float, float, float, float, float, float, float → float
- **`RAY2D.HITRECT_POINTY`** — args: float, float, float, float, float, float, float, float → float
- **`RAY2D.HITSEGMENT_DISTANCE`** — args: float, float, float, float, float, float, float, float → float
- **`RAY2D.HITSEGMENT_HIT`** — args: float, float, float, float, float, float, float, float → bool — 2D ray vs segment (x1,y1)-(x2,y2)
- **`RAY2D.HITSEGMENT_POINTX`** — args: float, float, float, float, float, float, float, float → float
- **`RAY2D.HITSEGMENT_POINTY`** — args: float, float, float, float, float, float, float, float → float

*319 overloads in this section.*

---

## Audio (2D and 3D)

Guide: [06-AUDIO.md](06-AUDIO.md)

### AUDIO

- **`AUDIO.CLOSE`** — args: (none)
- **`AUDIO.GETMUSICLENGTH`** — args: handle → float
- **`AUDIO.GETMUSICPITCH`** — args: handle → float — Get music pitch (1.0 default; tracked after SETMUSICPITCH).
- **`AUDIO.GETMUSICTIME`** — args: handle → float
- **`AUDIO.GETMUSICVOLUME`** — args: handle → float — Get music volume (0..1).
- **`AUDIO.GETSOUNDPAN`** — args: handle → float — Get sound pan (0.5 default).
- **`AUDIO.GETSOUNDPITCH`** — args: handle → float — Get sound pitch (1.0 default).
- **`AUDIO.GETSOUNDVOLUME`** — args: handle → float — Get sound volume (0..1).
- **`AUDIO.INIT`** — args: (none)
- **`AUDIO.ISMUSICPLAYING`** — args: handle → bool
- **`AUDIO.ISSOUNDPLAYING`** — args: handle → bool
- **`AUDIO.LISTENERCAMERA`** — args: handle
- **`AUDIO.LOADMUSIC`** — args: string → handle
- **`AUDIO.LOADSOUND`** — args: string → handle
- **`AUDIO.PAUSE`** — args: handle → handle
- **`AUDIO.PLAY`** — args: handle → handle
- **`AUDIO.PLAYMUSIC`** — args: handle
- **`AUDIO.PLAYRNDSOUND`** — args: handle, int
- **`AUDIO.PLAYSOUND`** — args: handle
- **`AUDIO.PLAYVARYSOUND`** — args: handle, float, float
- **`AUDIO.RESUME`** — args: handle → handle
- **`AUDIO.SEEKMUSIC`** — args: handle, float
- **`AUDIO.SETMASTERVOLUME`** — args: float
- **`AUDIO.SETMUSICPITCH`** — args: handle, float
- **`AUDIO.SETMUSICVOLUME`** — args: handle, float
- **`AUDIO.SETSOUNDPAN`** — args: handle, float
- **`AUDIO.SETSOUNDPITCH`** — args: handle, float
- **`AUDIO.SETSOUNDVOLUME`** — args: handle, float
- **`AUDIO.SETVOLUME`** — args: handle, float
- **`AUDIO.STOP`** — args: handle → handle
- **`AUDIO.STOPMUSIC`** — args: handle
- **`AUDIO.STOPSOUND`** — args: handle
- **`AUDIO.UPDATEMUSIC`** — args: handle

### AUDIOSTREAM

- **`AUDIOSTREAM.CREATE`** — args: int, int, int → handle
- **`AUDIOSTREAM.FREE`** — args: handle
- **`AUDIOSTREAM.GETPAN`** — args: handle → float
- **`AUDIOSTREAM.GETPITCH`** — args: handle → float
- **`AUDIOSTREAM.GETVOLUME`** — args: handle → float
- **`AUDIOSTREAM.ISPLAYING`** — args: handle → bool
- **`AUDIOSTREAM.ISREADY`** — args: handle → bool
- **`AUDIOSTREAM.MAKE`** — args: int, int, int → handle — DEPRECATED alias of AUDIOSTREAM.CREATE. Use AUDIOSTREAM.CREATE.
- **`AUDIOSTREAM.PAUSE`** — args: handle
- **`AUDIOSTREAM.PLAY`** — args: handle
- **`AUDIOSTREAM.RESUME`** — args: handle
- **`AUDIOSTREAM.SETPAN`** — args: handle, float
- **`AUDIOSTREAM.SETPITCH`** — args: handle, float
- **`AUDIOSTREAM.SETVOLUME`** — args: handle, float
- **`AUDIOSTREAM.STOP`** — args: handle
- **`AUDIOSTREAM.UPDATE`** — args: handle, handle

### SOUND

- **`SOUND.ATTACH`** — args: handle, handle — Pins a sound to an entity.
- **`SOUND.FREE`** — args: handle
- **`SOUND.FROMWAVE`** — args: handle → handle
- **`SOUND.PLAY3D`** — args: handle, float, float, float, float — Plays 3D spatialized audio.

*53 overloads in this section.*

---

## 2D sprites, tilemaps, terrain, particles, animation

Guide: [07-2D-WORLD.md](07-2D-WORLD.md)

### SPRITE

- **`SPRITE`** — args: string → handle
- **`SPRITE.ALPHA`** — args: handle → float — Property alias for SPRITE.GETALPHA
- **`SPRITE.COLOR`** — args: handle → handle — Property alias for SPRITE.GETCOLOR
- **`SPRITE.DEFANIM`** — args: handle, string
- **`SPRITE.DRAW`** — args: handle, int, int → handle
- **`SPRITE.FREE`** — args: handle
- **`SPRITE.GETALPHA`** — args: handle → float
- **`SPRITE.GETCOLOR`** — args: handle → handle — RGBA as floats (A channel 0â€“255) (Returns Color instance handle)
- **`SPRITE.GETPOS`** — args: handle → array
- **`SPRITE.GETROT`** — args: handle → handle — Returns [0, 0, roll] radians (2D screen rotation)
- **`SPRITE.GETSCALE`** — args: handle → handle — Returns [sx, sy, 1] scale factors (2D draw uses DrawTexturePro)
- **`SPRITE.HIT`** — args: handle, handle — True if the two sprites' drawn quads overlap (same scale, origin, and rotation as SPRITE.DRAW / DrawTexturePro; SAT on quad corners).
- **`SPRITE.LOAD`** — args: string
- **`SPRITE.PLAY`** — args: handle, int, int, float, bool → handle — Animate frames start..end at speed (frames/sec); call SPRITE.UPDATEANIM with Time.Delta()
- **`SPRITE.PLAYANIM`** — args: handle, string
- **`SPRITE.POINTHIT`** — args: handle, float, float — True if (x,y) lies inside the sprite's drawn quad (same space as SPRITE.DRAW position plus SETPOS offsets).
- **`SPRITE.POS`** — args: handle → array — Property alias for SPRITE.GETPOS
- **`SPRITE.ROT`** — args: handle → array — Property alias for SPRITE.GETROT
- **`SPRITE.SCALE`** — args: handle → array — Property alias for SPRITE.GETSCALE
- **`SPRITE.SETALPHA`** — args: handle, float → handle
- **`SPRITE.SETCOLOR`** — args: handle, int, int, int → handle
- **`SPRITE.SETCOLOR`** — args: handle, int, int, int, float → handle
- **`SPRITE.SETFRAME`** — args: handle, int → handle — Manual frame index (strip / DEFANIM); stops SPRITE.PLAY range playback
- **`SPRITE.SETORIGIN`** — args: handle, float, float → handle — Pivot offset in pixels (subtracted from draw position)
- **`SPRITE.SETPOS`** — args: handle, float, float → handle
- **`SPRITE.SETPOSITION`** — args: handle, float, float → handle — DEPRECATED alias of SPRITE.SETPOS. Use SPRITE.SETPOS.
- **`SPRITE.SETROT`** — args: handle, float → handle — Sets rotation in radians (CCW)
- **`SPRITE.SETSCALE`** — args: handle, float, float → handle
- **`SPRITE.UPDATEANIM`** — args: handle, float → handle

### TILEMAP

- **`TILEMAP.COLLISIONAT`** — args: handle, int, int → int
- **`TILEMAP.DRAW`** — args: handle
- **`TILEMAP.DRAWLAYER`** — args: handle, int
- **`TILEMAP.FREE`** — args: handle
- **`TILEMAP.GETTILE`** — args: handle, int, int, int → int
- **`TILEMAP.HEIGHT`** — args: handle → int
- **`TILEMAP.ISSOLID`** — args: handle, int, int → bool
- **`TILEMAP.ISSOLIDCATEGORY`** — args: handle, int, int, int → bool
- **`TILEMAP.LAYERCOUNT`** — args: handle → int
- **`TILEMAP.LAYERNAME`** — args: handle, int → string
- **`TILEMAP.LOAD`** — args: string → handle
- **`TILEMAP.MERGECOLLISIONLAYER`** — args: handle, int, int
- **`TILEMAP.SETCOLLISION`** — args: handle, int, int, int
- **`TILEMAP.SETTILE`** — args: handle, int, int, int, int
- **`TILEMAP.SETTILESIZE`** — args: handle, int, int
- **`TILEMAP.WIDTH`** — args: handle → int

### TERRAIN

- **`TERRAIN.APPLYMAP`** — args: handle, handle → handle — Apply CPU image as terrain diffuse + splat sample; rebuilds loaded chunk meshes
- **`TERRAIN.APPLYTILES`** — args: handle, handle, int → int — Copy template entity to each non-empty tile on layer 0; returns count placed
- **`TERRAIN.APPLYTILES`** — args: handle, handle, int, int → int — Same as 3-arg form with explicit tile layer index
- **`TERRAIN.CREATE`** — args: int, int
- **`TERRAIN.CREATE`** — args: int, int, float → handle
- **`TERRAIN.DRAW`** — args: handle → handle
- **`TERRAIN.FILLFLAT`** — args: handle, float
- **`TERRAIN.FILLPERLIN`** — args: handle, float, float
- **`TERRAIN.FREE`** — args: handle
- **`TERRAIN.GETDETAIL`** — args: handle → float
- **`TERRAIN.GETHEIGHT`** — args: handle, float, float → float
- **`TERRAIN.GETNORMAL`** — args: handle, float, float → handle — Unit terrain normal (heap vec3) for slope tilt
- **`TERRAIN.GETPOS`** — args: handle → handle — Get terrain position as Vec3.
- **`TERRAIN.GETPOS`** — args: handle → array — Returns [x,y,z] position of terrain.
- **`TERRAIN.GETROT`** — args: handle → array — Returns [x,y,z] rotation of terrain.
- **`TERRAIN.GETSCALE`** — args: handle → handle
- **`TERRAIN.GETSLOPE`** — args: handle, float, float → float
- **`TERRAIN.GETSPLAT`** — args: handle, float, float → int — Diffuse/splat map red channel 0..255 (-1 if no map); use for footstep ids
- **`TERRAIN.LOAD`** — args: string, string → handle — Heightmap image path + optional diffuse/splat path; GPU mesh + CPU splat sample
- **`TERRAIN.LOWER`** — args: handle, float, float, float, float → handle
- **`TERRAIN.MAKE`** — args: int, int — DEPRECATED alias of TERRAIN.CREATE. Use TERRAIN.CREATE.
- **`TERRAIN.MAKE`** — args: int, int, float → handle — DEPRECATED alias of TERRAIN.CREATE. Use TERRAIN.CREATE.
- **`TERRAIN.PLACE`** — args: handle, int, float, float, float
- **`TERRAIN.RAISE`** — args: handle, float, float, float, float → handle
- **`TERRAIN.RAYCAST`** — args: handle, float, float, float, float, float, float → handle — Ray vs terrain only; float array [hit, x, y, z]; max ray length is large by default
- **`TERRAIN.SETASYNCMESHBUILD`** — args: handle, bool → handle — When true, CPU heightmap prep runs on a background goroutine; GenMeshHeightmap still runs on the main thread when jobs drain (use with WINDOW.SETLOADINGMODE / mesh budget).
- **`TERRAIN.SETCHUNKSIZE`** — args: handle, int → handle
- **`TERRAIN.SETDETAIL`** — args: handle, float → handle — LOD factor in (0,1]: lower = coarser chunk meshes
- **`TERRAIN.SETMESHBUILDBUDGET`** — args: handle, int → handle — Max chunk mesh GPU rebuilds per WORLD.UPDATE tick; 0 = unlimited (default). Use 1â€“4 to avoid UI thread stalls.
- **`TERRAIN.SETPOS`** — args: handle, float, float, float → handle
- **`TERRAIN.SETPOSITION`** — args: handle, float, float, float → handle — DEPRECATED alias of TERRAIN.SETPOS. Use TERRAIN.SETPOS.
- **`TERRAIN.SETROT`** — args: handle, float — Set Y-axis rotation of terrain.
- **`TERRAIN.SETROT`** — args: handle, float, float, float — Set full X,Y,Z rotation of terrain.
- **`TERRAIN.SETSCALE`** — args: handle, float, float, float → handle — Non-uniform scale: XZ stretch per cell, Y height multiplier (marks chunks dirty)
- **`TERRAIN.SNAPY`** — args: handle, int, float

### PARTICLE

- **`PARTICLE.COUNT`** — args: handle → int
- **`PARTICLE.CREATE`** — args: (none) → handle
- **`PARTICLE.DRAW`** — args: handle → handle
- **`PARTICLE.DRAW`** — args: handle, handle → handle
- **`PARTICLE.FREE`** — args: handle
- **`PARTICLE.GETALPHA`** — args: handle → float
- **`PARTICLE.GETCOLOR`** — args: handle → handle — (Returns Color instance handle)
- **`PARTICLE.GETPOS`** — args: handle → array
- **`PARTICLE.GETSIZE`** — args: handle → handle — Emitter start/end size as Vec2 (sizeStartMin, sizeEndMin); aligns with PARTICLE.SETSIZE.
- **`PARTICLE.GETVELOCITY`** — args: handle → array — Emitter base direction (vx, vy, vz) last set with PARTICLE.SETVELOCITY (VEC3-compatible handle).
- **`PARTICLE.ISALIVE`** — args: handle → int
- **`PARTICLE.MAKE`** — args: (none) → handle — DEPRECATED alias of PARTICLE.CREATE. Use PARTICLE.CREATE.
- **`PARTICLE.PLAY`** — args: handle → handle
- **`PARTICLE.SETBILLBOARD`** — args: handle, bool → handle
- **`PARTICLE.SETBURST`** — args: handle, int → handle
- **`PARTICLE.SETCOLOR`** — args: handle, int, int, int, int → handle
- **`PARTICLE.SETCOLOREND`** — args: handle, int, int, int, int → handle
- **`PARTICLE.SETDIRECTION`** — args: handle, float, float, float → handle
- **`PARTICLE.SETEMITRATE`** — args: handle, float → handle
- **`PARTICLE.SETENDCOLOR`** — args: handle, int, int, int, int → handle
- **`PARTICLE.SETENDSIZE`** — args: handle, float, float → handle
- **`PARTICLE.SETGRAVITY`** — args: handle, float → handle
- **`PARTICLE.SETGRAVITY`** — args: handle, float, float, float → handle
- **`PARTICLE.SETLIFETIME`** — args: handle, float, float → handle
- **`PARTICLE.SETPOS`** — args: handle, float, float, float → handle
- **`PARTICLE.SETPOSITION`** — args: handle, float, float, float → handle — DEPRECATED alias of PARTICLE.SETPOS. Use PARTICLE.SETPOS.
- **`PARTICLE.SETRATE`** — args: handle, float → handle
- **`PARTICLE.SETSIZE`** — args: handle, float, float → handle
- **`PARTICLE.SETSPEED`** — args: handle, float, float → handle
- **`PARTICLE.SETSPREAD`** — args: handle, float → handle
- **`PARTICLE.SETSTARTCOLOR`** — args: handle, int, int, int, int → handle
- **`PARTICLE.SETSTARTSIZE`** — args: handle, float, float → handle
- **`PARTICLE.SETTEXTURE`** — args: handle, handle → handle
- **`PARTICLE.SETVELOCITY`** — args: handle, float, float, float, float → handle
- **`PARTICLE.STOP`** — args: handle → handle
- **`PARTICLE.UPDATE`** — args: handle, float → handle

### ANIM

- **`ANIM.ADDTRANSITION`** — args: handle, string, string, string
- **`ANIM.DEFINE`** — args: handle, string, int, int, float, bool
- **`ANIM.SETPARAM`** — args: handle, string, any
- **`ANIM.UPDATE`** — args: handle, float

### WORLD

- **`WORLD.DAYNIGHTCYCLE`** — args: float — Rotates global sunlight over duration (seconds).
- **`WORLD.EXPLOSION`** — args: float, float, float, float, float — Alias of PHYSICS.EXPLOSION
- **`WORLD.FOGCOLOR`** — args: int, int, int
- **`WORLD.FOGDENSITY`** — args: float
- **`WORLD.FOGMODE`** — args: int
- **`WORLD.GETRAY`** — args: float, float, handle → handle — Returns Array [px,py,pz,dx,dy,dz]
- **`WORLD.GRAVITY`** — args: float, float, float — Alias: forwards to PHYSICS3D.SETGRAVITY (global Jolt gravity)
- **`WORLD.HITSTOP`** — args: float — Freeze gameplay delta for duration (wall-clock seconds) â€” impact frames
- **`WORLD.ISREADY`** — args: handle → bool
- **`WORLD.MOUSE2D`** — args: handle → handle — Mouse position through Camera2D; float array [wx,wy]
- **`WORLD.MOUSEFLOOR`** — args: handle, float → handle — Alias of WORLD.MOUSEFLOOR3D â€” mouse ray vs plane y=floorY â†’ [wx,wz] or NIL
- **`WORLD.MOUSEFLOOR3D`** — args: handle, float → handle — Mouse ray vs plane y=floorY; float array [wx,wz] or NIL
- **`WORLD.MOUSEPICK`** — args: handle → int — Alias of WORLD.MOUSETOENTITY â€” entity id under mouse cursor (physics ray; Linux+Jolt)
- **`WORLD.MOUSETOENTITY`** — args: handle → int — Jolt ray pick at cursor (Linux+CGO); entity or 0. Same as CAMERA.RAYCASTMOUSE
- **`WORLD.MOUSETOFLOOR`** — args: handle, float → handle — Alias of WORLD.MOUSEFLOOR3D
- **`WORLD.PRELOAD`** — args: handle, int
- **`WORLD.SCREENSHAKE`** — args: float, float — Shakes the primary camera.
- **`WORLD.SETAMBIENCE`** — args: handle, float — Plays a looping background track.
- **`WORLD.SETCENTER`** — args: float, float
- **`WORLD.SETCENTERENTITY`** — args: int
- **`WORLD.SETGRAVITY`** — args: float, float, float — Alias of PHYSICS3D.SETGRAVITY
- **`WORLD.SETREFLECTION`** — args: int
- **`WORLD.SETREVERB`** — args: int — Changes echo.
- **`WORLD.SETTIMESCALE`** — args: float — Alias of GAME.SETTIMESCALE
- **`WORLD.SETVEGETATION`** — args: handle, handle, float — Scatter helper: terrain + billboard entity reserved + density; uses internal SCATTER set
- **`WORLD.SHAKE`** — args: float, float — Alias of WORLD.SCREENSHAKE â€” screen impact via active camera
- **`WORLD.STATUS`** — args: (none) → string
- **`WORLD.STREAMENABLE`** — args: bool
- **`WORLD.TOSCREEN`** — args: int → handle — WORLD.TOSCREEN(entity) â€” screen [x,y] for entity world position via active 3D camera
- **`WORLD.TOSCREEN`** — args: float, float, float → handle — World to screen using active CAMERA.BEGIN 3D camera; returns float array [sx,sy]
- **`WORLD.TOSCREEN`** — args: float, float, float, handle → handle — Returns 2D Screen coords given 3D World coords and Camera.
- **`WORLD.TOWORLD`** — args: float, float, float → handle — Unproject screen x,y with depth along view ray (active 3D camera); returns [wx,wy,wz]
- **`WORLD.TOWORLD`** — args: float, float, float, handle → handle — Returns 3D World coords from 2D.
- **`WORLD.UPDATE`** — args: float

### CHUNK

- **`CHUNK.COUNT`** — args: handle → int
- **`CHUNK.GENERATE`** — args: handle, int, int → handle
- **`CHUNK.ISLOADED`** — args: handle, int, int → bool
- **`CHUNK.SETRANGE`** — args: handle, float, float → handle

*158 overloads in this section.*

---

## UI, fonts, and text

Guide: [08-UI-TEXT.md](08-UI-TEXT.md)

### GUI

- **`GUI.BUTTON`** — args: float, float, float, float, string → bool
- **`GUI.CHECKBOX`** — args: float, float, float, float, string, bool → bool
- **`GUI.COLORBARALPHA`** — args: float, float, float, float, string, float → float
- **`GUI.COLORBARHUE`** — args: float, float, float, float, string, float → float
- **`GUI.COLORPANEL`** — args: float, float, float, float, string, int, int, int, int → handle
- **`GUI.COLORPANELHSV`** — args: float, float, float, float, string, handle → int
- **`GUI.COLORPICKER`** — args: float, float, float, float, string, int, int, int, int → handle
- **`GUI.COLORPICKERHSV`** — args: float, float, float, float, string, handle → int
- **`GUI.COMBOBOX`** — args: float, float, float, float, string, int → int
- **`GUI.DISABLE`** — args: (none)
- **`GUI.DISABLETOOLTIP`** — args: (none)
- **`GUI.DRAWICON`** — args: int, int, int, int, int, int, int, int
- **`GUI.DRAWRECTANGLE`** — args: float, float, float, float, int, int, int, int, int, int, int, int, int
- **`GUI.DRAWTEXT`** — args: string, float, float, float, float, int, int, int, int, int
- **`GUI.DROPDOWNBOX`** — args: float, float, float, float, string, handle → bool
- **`GUI.DUMMYREC`** — args: float, float, float, float, string
- **`GUI.ENABLE`** — args: (none)
- **`GUI.ENABLETOOLTIP`** — args: (none)
- **`GUI.FADE`** — args: int, int, int, int, float → handle
- **`GUI.GETCOLOR`** — args: int, int → handle
- **`GUI.GETSTATE`** — args: (none) → int
- **`GUI.GETSTYLE`** — args: int, int → int
- **`GUI.GETTEXTBOUNDS`** — args: int, float, float, float, float → handle
- **`GUI.GETTEXTSIZE`** — args: (none) → int
- **`GUI.GETTEXTWIDTH`** — args: string → int
- **`GUI.GRID`** — args: float, float, float, float, string, float, int, handle → int
- **`GUI.GROUPBOX`** — args: float, float, float, float, string
- **`GUI.ICONTEXT`** — args: int, string → string
- **`GUI.ISLOCKED`** — args: (none) → bool
- **`GUI.LABEL`** — args: float, float, float, float, string
- **`GUI.LABELBUTTON`** — args: float, float, float, float, string → bool
- **`GUI.LINE`** — args: float, float, float, float, string
- **`GUI.LISTVIEW`** — args: float, float, float, float, string, handle → int
- **`GUI.LISTVIEWEX`** — args: float, float, float, float, string, handle → int
- **`GUI.LOADDEFAULTSTYLE`** — args: (none)
- **`GUI.LOADICONS`** — args: string, bool
- **`GUI.LOADICONSMEM`** — args: string, bool
- **`GUI.LOADSTYLE`** — args: string
- **`GUI.LOADSTYLEMEM`** — args: string
- **`GUI.LOCK`** — args: (none)
- **`GUI.MESSAGEBOX`** — args: float, float, float, float, string, string, string → int
- **`GUI.PANEL`** — args: float, float, float, float, string
- **`GUI.PROGRESSBAR`** — args: float, float, float, float, string, string, float, float, float → float
- **`GUI.SCROLLBAR`** — args: float, float, float, float, int, int, int → int
- **`GUI.SCROLLPANEL`** — args: float, float, float, float, string, float, float, float, float, handle
- **`GUI.SETALPHA`** — args: float
- **`GUI.SETCOLOR`** — args: int, int, int, int, int, int
- **`GUI.SETFONT`** — args: handle
- **`GUI.SETICONSCALE`** — args: int
- **`GUI.SETSTATE`** — args: int
- **`GUI.SETSTYLE`** — args: int, int, int
- **`GUI.SETTEXTALIGN`** — args: int
- **`GUI.SETTEXTALIGNVERT`** — args: int
- **`GUI.SETTEXTLINEHEIGHT`** — args: int
- **`GUI.SETTEXTSIZE`** — args: int
- **`GUI.SETTEXTSPACING`** — args: int
- **`GUI.SETTEXTWRAP`** — args: int
- **`GUI.SETTOOLTIP`** — args: string
- **`GUI.SLIDER`** — args: float, float, float, float, string, string, float, float, float → float
- **`GUI.SLIDERBAR`** — args: float, float, float, float, string, string, float, float, float → float
- **`GUI.SPINNER`** — args: float, float, float, float, string, int, int, int, bool → int
- **`GUI.STATUSBAR`** — args: float, float, float, float, string
- **`GUI.TABBAR`** — args: float, float, float, float, string, handle → int
- **`GUI.TEXTBOX`** — args: float, float, float, float, string, int, bool → string
- **`GUI.TEXTINPUTBOX`** — args: float, float, float, float, string, string, string, string, int, handle → int
- **`GUI.TEXTINPUTLAST`** — args: (none) → string
- **`GUI.THEMEAPPLY`** — args: string
- **`GUI.THEMENAMES`** — args: (none) → string
- **`GUI.TOGGLE`** — args: float, float, float, float, string, bool → bool
- **`GUI.TOGGLEGROUP`** — args: float, float, float, float, string → int
- **`GUI.TOGGLEGROUPAT`** — args: float, float, float, float, string, int → int
- **`GUI.TOGGLESLIDER`** — args: float, float, float, float, string, int → int
- **`GUI.UNLOCK`** — args: (none)
- **`GUI.VALUEBOX`** — args: float, float, float, float, string, int, int, int, bool → int
- **`GUI.VALUEBOXFLOAT`** — args: float, float, float, float, string, float, string, bool → float
- **`GUI.VALUEBOXFLOATTEXT`** — args: (none) → string
- **`GUI.WINDOWBOX`** — args: float, float, float, float, string → bool

### FONT

- **`FONT.DRAWDEFAULT`** — args: (none)
- **`FONT.FREE`** — args: handle
- **`FONT.LOAD`** — args: string
- **`FONT.LOADBDF`** — args: string, int
- **`FONT.SETDEFAULT`** — args: handle

### DRAW

- **`DRAW.ARC`** — args: float, float, float, float, float, float, int, int, int, int
- **`DRAW.BILLBOARD`** — args: handle, float, float, float, float, int, int, int, int
- **`DRAW.BILLBOARDREC`** — args: handle, float, float, float, float, float, float, float, float, int, int, int, int
- **`DRAW.BOUNDINGBOX`** — args: float, float, float, float, float, float, int, int, int, int
- **`DRAW.CAPSULE`** — args: float, float, float, float, float, float, float, float, float, int, int, int, int, int, int
- **`DRAW.CAPSULEWIRES`** — args: float, float, float, float, float, float, float, float, float, int, int, int, int, int, int
- **`DRAW.CENTERTEXT`** — args: string, int, int, int, int, int, int, int
- **`DRAW.CIRCLE`** — args: int, int, float, int, int, int, int
- **`DRAW.CIRCLEGRADIENT`** — args: int, int, int, int, int, int, int, int, int, int, int
- **`DRAW.CIRCLELINES`** — args: int, int, float, int, int, int, int
- **`DRAW.CIRCLESECTOR`** — args: int, int, int, int, int, int, int, int, int, int
- **`DRAW.CROSSHAIR`** — args: int, int, int, int, int, int
- **`DRAW.CUBE`** — args: float, float, float, float, float, float, float, int, int, int, int
- **`DRAW.CUBEWIRES`** — args: float, float, float, float, float, float, float, int, int, int, int
- **`DRAW.CYLINDER`** — args: float, float, float, float, float, float, float, int, int, int, int
- **`DRAW.CYLINDERWIRES`** — args: float, float, float, float, float, float, float, int, int, int, int
- **`DRAW.DOT`** — args: float, float, float, int, int, int, int
- **`DRAW.ELLIPSE`** — args: int, int, float, float, int, int, int, int
- **`DRAW.ELLIPSELINES`** — args: int, int, float, float, int, int, int, int
- **`DRAW.GETPIXELCOLOR`** — args: int, int → array
- **`DRAW.GRID`** — args: int, float
- **`DRAW.GRID2D`** — args: int, int, int, int, int
- **`DRAW.HEALTHBAR`** — args: int, int, int, int, int, float, int, int, int, int, int, int, int, int, int, int, int, int
- **`DRAW.LINE`** — args: int, int, int, int, int, int, int, int
- **`DRAW.LINE3D`** — args: float, float, float, float, float, float, int, int, int, int
- **`DRAW.LINEBEZIER`** — args: float, float, float, float, float, int, int, int, int
- **`DRAW.LINEBEZIERCUBIC`** — args: float, float, float, float, float, float, float, float, float, int, int, int, int
- **`DRAW.LINEBEZIERQUAD`** — args: float, float, float, float, float, float, float, int, int, int, int
- **`DRAW.LINEEX`** — args: float, float, float, float, float, int, int, int, int
- **`DRAW.OUTLINETEXT`** — args: string, int, int, int, int, int, int, int, int, int, int, int
- **`DRAW.OVAL`** — args: int, int, float, float, int, int, int, int
- **`DRAW.OVALLINES`** — args: int, int, float, float, int, int, int, int
- **`DRAW.PIXEL`** — args: int, int, int, int, int, int
- **`DRAW.PIXELV`** — args: float, float, int, int, int, int
- **`DRAW.PLANE`** — args: float, float, float, float, float, float, int, int, int, int
- **`DRAW.PLOT`** — args: int, int, int, int, int, int
- **`DRAW.POINT3D`** — args: float, float, float, int, int, int, int
- **`DRAW.POLY`** — args: float, float, int, float, float, int, int, int, int
- **`DRAW.POLYLINES`** — args: float, float, int, float, float, float, int, int, int, int
- **`DRAW.PROGRESSBAR`** — args: int, int, int, int, int, float, int, int, int, int, int, int, int, int
- **`DRAW.RAY`** — args: handle, int, int, int, int
- **`DRAW.RECTANGLE`** — args: int, int, int, int, int, int, int, int
- **`DRAW.RECTANGLE_ROUNDED`** — args: int, int, int, int, int, int, int, int, int
- **`DRAW.RECTGRAD`** — args: int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int, int
- **`DRAW.RECTGRADH`** — args: int, int, int, int, int, int, int, int, int, int, int, int
- **`DRAW.RECTGRADV`** — args: int, int, int, int, int, int, int, int, int, int, int, int
- **`DRAW.RECTGRID`** — args: int, int, int, int, int, int
- **`DRAW.RECTLINES`** — args: int, int, int, int, int, float, int, int, int, int
- **`DRAW.RECTPRO`** — args: int, int, int, int, float, float, float, int, int, int, int
- **`DRAW.RIGHTTEXT`** — args: string, int, int, int, int, int, int, int
- **`DRAW.RING`** — args: float, float, float, float, float, float, int, int, int, int, int
- **`DRAW.RINGLINES`** — args: float, float, float, float, float, float, int, int, int, int, int
- **`DRAW.SETPIXELCOLOR`** — args: int, int, int, int, int, int
- **`DRAW.SHADOWTEXT`** — args: string, int, int, int, int, int, int, int, int, int, int, int
- **`DRAW.SPHERE`** — args: float, float, float, float, int, int, int, int
- **`DRAW.SPHEREWIRES`** — args: float, float, float, float, int, int, int, int, int, int
- **`DRAW.SPLINEBASIS`** — args: handle, float, int, int, int, int
- **`DRAW.SPLINEBEZIERCUBIC`** — args: handle, float, int, int, int, int
- **`DRAW.SPLINEBEZIERQUAD`** — args: handle, float, int, int, int, int
- **`DRAW.SPLINECATMULLROM`** — args: handle, float, int, int, int, int
- **`DRAW.SPLINELINEAR`** — args: handle, float, int, int, int, int
- **`DRAW.TEXT`** — args: string, int, int, int, int, int, int, int
- **`DRAW.TEXTEX`** — args: handle, string, float, float, float, float, int, int, int, int
- **`DRAW.TEXTFONT`** — args: handle, string, float, float, float, float, int, int, int, int — Same handler as DRAW.TEXTEX â€” DrawTextEx with a loaded FONT.* handle
- **`DRAW.TEXTFONTWIDTH`** — args: handle, string, float, float → float
- **`DRAW.TEXTPRO`** — args: handle, string, float, float, float, float, float, float, float, int, int, int, int
- **`DRAW.TEXTURE`** — args: handle, int, int, int, int, int, int
- **`DRAW.TEXTUREEX`** — args: handle, float, float, float, float, int, int, int, int
- **`DRAW.TEXTUREFLIPPED`** — args: handle
- **`DRAW.TEXTUREFULL`** — args: handle
- **`DRAW.TEXTURENPATCH`** — args: handle, int, int, int, int, int, int, int, int, int, int, int, int
- **`DRAW.TEXTUREPRO`** — args: handle, float, float, float, float, float, float, float, float, float, float, int, int, int, int
- **`DRAW.TEXTUREREC`** — args: handle, float, float, float, float, float, float, int, int, int, int
- **`DRAW.TEXTURETILED`** — args: handle, float, float, float, float, float, float, float, float, float, float, float, float, int, int, int, int
- **`DRAW.TEXTUREV`** — args: handle, float, float, int, int, int, int
- **`DRAW.TEXTWIDTH`** — args: string, int → int
- **`DRAW.TRIANGLE`** — args: float, float, float, float, float, float, int, int, int, int
- **`DRAW.TRIANGLELINES`** — args: float, float, float, float, float, float, int, int, int, int

### TEXT

- **`TEXT`** — args: int, int, string
- **`TEXT.DRAW`** — args: string, int, int
- **`TEXT.DRAWFONT`** — args: handle, string, int, int
- **`TEXT.SIZE`** — args: string → int

### COLOR

- **`COLOR.A`** — args: handle → int
- **`COLOR.B`** — args: handle → int
- **`COLOR.BRIGHTNESS`** — args: handle, float → handle
- **`COLOR.CLAMP`** — args: float, float, float → handle
- **`COLOR.CONTRAST`** — args: handle, float → handle
- **`COLOR.FADE`** — args: handle, float → handle
- **`COLOR.FREE`** — args: handle
- **`COLOR.FROMHSV`** — args: float, float, float → handle
- **`COLOR.G`** — args: handle → int
- **`COLOR.HEX`** — args: string → handle
- **`COLOR.HSV`** — args: float, float → handle — COLOR.HSV(index, total) â€” evenly spaced hues on the wheel
- **`COLOR.HSV`** — args: float, float, float → handle
- **`COLOR.INVERT`** — args: handle → handle
- **`COLOR.LERP`** — args: handle, handle, float → handle
- **`COLOR.R`** — args: handle → int
- **`COLOR.RGB`** — args: int, int, int → handle
- **`COLOR.RGBA`** — args: int, int, int, int → handle
- **`COLOR.TOHEX`** — args: handle → string
- **`COLOR.TOHSV`** — args: handle → handle
- **`COLOR.TOHSVX`** — args: handle → float
- **`COLOR.TOHSVY`** — args: handle → float
- **`COLOR.TOHSVZ`** — args: handle → float

*186 overloads in this section.*

---

## Save, files, JSON, math, vectors

Guide: [09-DATA.md](09-DATA.md)

### SAVE

- **`SAVE.DATA`** — args: string, string — Writes JSON data.
- **`SAVE.GET`** — args: string → string — Reads JSON data.
- **`SAVE.READ`** — args: string
- **`SAVE.READFILE`** — args: string
- **`SAVE.SET`** — args: string, string
- **`SAVE.WRITE`** — args: string
- **`SAVE.WRITEFILE`** — args: string

### FILE

- **`FILE.CLOSE`** — args: handle
- **`FILE.DELETE`** — args: string → bool
- **`FILE.EOF`** — args: handle → bool
- **`FILE.EXISTS`** — args: any → bool
- **`FILE.GETEOF`** — args: handle → bool — Alias of FILE.EOF getter.
- **`FILE.GETPOS`** — args: handle → int — Alias of FILE.POS getter.
- **`FILE.GETSIZE`** — args: handle → int — Alias of FILE.SIZE getter.
- **`FILE.OPEN`** — args: string, string → handle
- **`FILE.OPENREAD`** — args: string → handle
- **`FILE.OPENWRITE`** — args: string → handle
- **`FILE.READALLTEXT`** — args: any → string
- **`FILE.READLINE`** — args: handle
- **`FILE.READTEXT`** — args: string → string
- **`FILE.SEEK`** — args: handle, int → handle
- **`FILE.SIZE`** — args: handle → int
- **`FILE.TELL`** — args: handle → int
- **`FILE.WRITE`** — args: handle, string — Write string to file without appending a newline.
- **`FILE.WRITEALLTEXT`** — args: any, any
- **`FILE.WRITELN`** — args: handle, string — Write string to file and append a newline.
- **`FILE.WRITETEXT`** — args: string, string

### JSON

- **`JSON.APPEND`** — args: handle, any
- **`JSON.CLEAR`** — args: handle
- **`JSON.CREATE`** — args: (none) → handle
- **`JSON.DELETE`** — args: handle, string
- **`JSON.FREE`** — args: handle
- **`JSON.GET`** — args: handle, string → any
- **`JSON.GETARRAY`** — args: handle, string → handle
- **`JSON.GETBOOL`** — args: handle, string → bool
- **`JSON.GETFLOAT`** — args: handle, string → float
- **`JSON.GETINT`** — args: handle, string → int
- **`JSON.GETOBJECT`** — args: handle, string → handle
- **`JSON.GETSTRING`** — args: handle, string → string
- **`JSON.HAS`** — args: handle, string → bool
- **`JSON.KEYS`** — args: handle → handle
- **`JSON.LEN`** — args: handle → int
- **`JSON.LOADFILE`** — args: any → handle
- **`JSON.MAKE`** — args: (none) → handle — DEPRECATED alias of JSON.CREATE. Use JSON.CREATE.
- **`JSON.MAKEARRAY`** — args: (none) → handle
- **`JSON.MINIFY`** — args: handle → string
- **`JSON.PARSE`** — args: string → handle
- **`JSON.PARSESTRING`** — args: string → handle
- **`JSON.PRETTY`** — args: handle → string
- **`JSON.QUERY`** — args: handle, string → any
- **`JSON.SAVEFILE`** — args: any, any
- **`JSON.SET`** — args: handle, string, string
- **`JSON.SETBOOL`** — args: handle, bool
- **`JSON.SETFLOAT`** — args: handle, string, float
- **`JSON.SETINT`** — args: handle, string, int
- **`JSON.SETNULL`** — args: handle, string
- **`JSON.SETSTRING`** — args: handle, string, string
- **`JSON.STRINGIFY`** — args: handle → string
- **`JSON.TOCSV`** — args: handle → string
- **`JSON.TOFILE`** — args: handle, string
- **`JSON.TOFILEPRETTY`** — args: handle, string
- **`JSON.TOSTRING`** — args: handle → string
- **`JSON.TYPE`** — args: handle, string → string

### MATH

- **`MATH.ABS`** — args: any
- **`MATH.ACOS`** — args: any
- **`MATH.ANGLEDIFF`** — args: any, any
- **`MATH.ANGLEDIFFRAD`** — args: float, float → float — Same as ANGLEDIFFRAD
- **`MATH.ANGLETO`** — args: float, float, float, float → float — Same as ANGLETO
- **`MATH.APPROACH`** — args: float, float, float → float
- **`MATH.ASIN`** — args: any
- **`MATH.ATAN`** — args: any
- **`MATH.ATAN2`** — args: any, any
- **`MATH.ATN`** — args: any
- **`MATH.CEIL`** — args: any
- **`MATH.CHANCE`** — args: float → bool
- **`MATH.CIRCLEPOINT`** — args: float, float, float, float, float → handle
- **`MATH.CLAMP`** — args: any, any, any
- **`MATH.COS`** — args: any
- **`MATH.COSD`** — args: any
- **`MATH.CURVE`** — args: float, float, float → float — Alias of CURVE â€” current + (target-current)/divisor (divisor clamped to >= 1)
- **`MATH.DEG2RAD`** — args: any
- **`MATH.DEGPERSEC`** — args: any, any
- **`MATH.DIST2D`** — args: float, float, float, float → float — Same as DIST2D
- **`MATH.DISTSQ2D`** — args: float, float, float, float → float — Same as DISTSQ2D
- **`MATH.E`** — args: (none)
- **`MATH.EXP`** — args: any
- **`MATH.FIX`** — args: any
- **`MATH.FLOOR`** — args: any
- **`MATH.HDIST`** — args: float, float, float, float → float — Same as HDIST
- **`MATH.HDISTSQ`** — args: float, float, float, float → float — Same as HDISTSQ
- **`MATH.INVERSE_LERP`** — args: float, float, float → float
- **`MATH.LERP`** — args: any, any, any
- **`MATH.LERPANGLE`** — args: float, float, float → float
- **`MATH.LOG`** — args: any
- **`MATH.LOG10`** — args: any
- **`MATH.LOG2`** — args: any
- **`MATH.MAX`** — args: any, any
- **`MATH.MIN`** — args: any, any
- **`MATH.NEWX`** — args: float, float, float → float — currentX + MOVEX(yaw,1,0)*dist â€” yaw in radians (XZ forward step)
- **`MATH.NEWZ`** — args: float, float, float → float — currentZ + MOVEZ(yaw,1,0)*dist â€” yaw in radians
- **`MATH.PI`** — args: (none)
- **`MATH.PINGPONG`** — args: any, any
- **`MATH.POW`** — args: any, any
- **`MATH.RAD2DEG`** — args: any
- **`MATH.RAND`** — args: any, any → int — Same as RAND
- **`MATH.RANGE`** — args: float, float → float
- **`MATH.REMAP`** — args: float, float, float, float, float → float
- **`MATH.RND`** — args: (none)
- **`MATH.RND`** — args: any
- **`MATH.RND`** — args: any, any → int — Inclusive int range â€” same as RND(lo, hi)
- **`MATH.RNDF`** — args: any, any
- **`MATH.RNDSEED`** — args: any
- **`MATH.ROUND`** — args: any
- **`MATH.ROUND`** — args: any, any
- **`MATH.SATURATE`** — args: float → float
- **`MATH.SGN`** — args: any
- **`MATH.SIGN`** — args: any
- **`MATH.SIN`** — args: any
- **`MATH.SIND`** — args: any
- **`MATH.SMOOTHERSTEP`** — args: any, any, any → float — Same as SMOOTHERSTEP
- **`MATH.SMOOTHSTEP`** — args: any, any, any
- **`MATH.SQR`** — args: any
- **`MATH.SQRT`** — args: any
- **`MATH.TAN`** — args: any
- **`MATH.TAND`** — args: any
- **`MATH.TAU`** — args: (none)
- **`MATH.WRAP`** — args: any, any, any
- **`MATH.WRAPANGLE`** — args: any
- **`MATH.WRAPANGLE180`** — args: any
- **`MATH.YAWFROMXZ`** — args: float, float → float — Same as YAWFROMXZ

### VEC3

- **`VEC3.ADD`** — args: handle, handle → handle
- **`VEC3.ANGLE`** — args: handle, handle → float
- **`VEC3.CREATE`** — args: float, float, float → handle
- **`VEC3.CROSS`** — args: handle, handle → handle
- **`VEC3.DIST`** — args: handle, handle → float
- **`VEC3.DIST`** — args: float, float, float, float, float, float → float
- **`VEC3.DISTANCE`** — args: handle, handle → float
- **`VEC3.DISTSQ`** — args: float, float, float, float, float, float → float
- **`VEC3.DIV`** — args: handle, float → handle
- **`VEC3.DOT`** — args: handle, handle → float
- **`VEC3.EQUALS`** — args: handle, handle → bool
- **`VEC3.FREE`** — args: handle
- **`VEC3.LENGTH`** — args: handle → float
- **`VEC3.LENGTH`** — args: float, float, float → float
- **`VEC3.LERP`** — args: handle, handle, float → handle
- **`VEC3.MAKE`** — args: float, float, float → handle — DEPRECATED alias of VEC3.CREATE. Use VEC3.CREATE.
- **`VEC3.MUL`** — args: handle, float → handle
- **`VEC3.NEGATE`** — args: handle → handle
- **`VEC3.NORMALIZE`** — args: handle → handle
- **`VEC3.NORMALIZE`** — args: float, float, float → handle
- **`VEC3.ORTHONORMALIZE`** — args: handle, handle
- **`VEC3.PROJECT`** — args: handle, handle → handle
- **`VEC3.REFLECT`** — args: handle, handle → handle
- **`VEC3.ROTATEBYQUAT`** — args: handle, handle → handle
- **`VEC3.SET`** — args: handle, float, float, float
- **`VEC3.SUB`** — args: handle, handle → handle
- **`VEC3.TRANSFORMMAT4`** — args: handle, handle → handle
- **`VEC3.VEC3`** — args: float, float, float → handle
- **`VEC3.VECADD`** — args: handle, handle → handle
- **`VEC3.VECCROSS`** — args: handle, handle → handle
- **`VEC3.VECDOT`** — args: handle, handle → float
- **`VEC3.VECLENGTH`** — args: handle → float
- **`VEC3.VECNORMALIZE`** — args: handle → handle
- **`VEC3.VECSCALE`** — args: handle, float → handle
- **`VEC3.VECSUB`** — args: handle, handle → handle
- **`VEC3.X`** — args: handle → float
- **`VEC3.Y`** — args: handle → float
- **`VEC3.Z`** — args: handle → float

### VEC2

- **`VEC2.ADD`** — args: handle, handle → handle
- **`VEC2.ANGLE`** — args: handle, handle → float
- **`VEC2.CREATE`** — args: float, float → handle
- **`VEC2.DIST`** — args: handle, handle → float
- **`VEC2.DIST`** — args: float, float, float, float → float
- **`VEC2.DISTANCE`** — args: handle, handle → float
- **`VEC2.DISTSQ`** — args: float, float, float, float → float
- **`VEC2.FREE`** — args: handle
- **`VEC2.LENGTH`** — args: handle → float
- **`VEC2.LENGTH`** — args: float, float → float
- **`VEC2.LERP`** — args: handle, handle, float → handle
- **`VEC2.MAKE`** — args: float, float → handle — DEPRECATED alias of VEC2.CREATE. Use VEC2.CREATE.
- **`VEC2.MOVE_TOWARD`** — args: float, float, float, float, float → handle
- **`VEC2.MUL`** — args: handle, float → handle
- **`VEC2.NORMALIZE`** — args: handle → handle
- **`VEC2.NORMALIZE`** — args: float, float → handle
- **`VEC2.PUSHOUT`** — args: float, float, float, float, float → handle
- **`VEC2.ROTATE`** — args: handle, float → handle
- **`VEC2.SET`** — args: handle, float, float
- **`VEC2.SUB`** — args: handle, handle → handle
- **`VEC2.TRANSFORMMAT4`** — args: handle, handle → handle
- **`VEC2.X`** — args: handle → float
- **`VEC2.Y`** — args: handle → float

### CONFIG

- **`CONFIG.DELETE`** — args: string
- **`CONFIG.GETBOOL`** — args: string → bool
- **`CONFIG.GETFLOAT`** — args: string → float
- **`CONFIG.GETINT`** — args: string → int
- **`CONFIG.GETSTRING`** — args: string → string
- **`CONFIG.HAS`** — args: string → bool
- **`CONFIG.LOAD`** — args: string
- **`CONFIG.SAVE`** — args: string
- **`CONFIG.SETBOOL`** — args: string, bool
- **`CONFIG.SETFLOAT`** — args: string, float
- **`CONFIG.SETINT`** — args: string, int
- **`CONFIG.SETSTRING`** — args: string, string

### CSV

- **`CSV.COLCOUNT`** — args: handle → int
- **`CSV.FREE`** — args: handle
- **`CSV.FROMSTRING`** — args: string → handle
- **`CSV.GET`** — args: handle, int, int → string
- **`CSV.LOAD`** — args: string → handle
- **`CSV.ROWCOUNT`** — args: handle → int
- **`CSV.SAVE`** — args: handle, string
- **`CSV.SET`** — args: handle, int, int, string
- **`CSV.TOJSON`** — args: handle → handle
- **`CSV.TOSTRING`** — args: handle → string

*213 overloads in this section.*

---

## Debug, timers

Guide: [10-DEBUG-TIMER.md](10-DEBUG-TIMER.md)

### DEBUG

- **`CONSOLE.LOG`** — args: string — Add a message to the scrolling on-screen debug console.
- **`CONSOLE.LOG`** — args: string, handle — Add a colored message to the scrolling on-screen debug console.
- **`DEBUG.ASSERT`** — args: any, string
- **`DEBUG.BREAKPOINT`** — args: (none)
- **`DEBUG.DISABLE`** — args: (none)
- **`DEBUG.DRAWBODY`** — args: handle — Renders body collision shape.
- **`DEBUG.DRAWBOX`** — args: float, float, float, float, float, float, int, int, int
- **`DEBUG.DRAWCHARACTER`** — args: handle — Renders capsule wireframe and ground probes for character Controller.
- **`DEBUG.DRAWLINE`** — args: float, float, float, float, float, float, int, int, int
- **`DEBUG.DRAWPHYSICS`** — args: bool — Toggle collision wireframe visualization.
- **`DEBUG.DUMPHEAP`** — args: (none) — Professional: Scan all active handles and print to diagnostics.
- **`DEBUG.ENABLE`** — args: (none)
- **`DEBUG.GCSTATS`** — args: (none)
- **`DEBUG.HEAPSTATS`** — args: (none)
- **`DEBUG.INSPECT`** — args: int — Display live transform/state info for an entity.
- **`DEBUG.ISENABLED`** — args: (none) → bool
- **`DEBUG.LISTCOMMANDS`** — args: (none) — Professional: List all registered built-in commands.
- **`DEBUG.LOG`** — args: string
- **`DEBUG.LOGFILE`** — args: string, string
- **`DEBUG.PRINT`** — args: string
- **`DEBUG.PRINT`** — args: any
- **`DEBUG.PRINT`** — args: string, any
- **`DEBUG.PRINT`** — args: string, any, any
- **`DEBUG.PRINT`** — args: string, any, any, any
- **`DEBUG.PRINT`** — args: string, any, any, any, any
- **`DEBUG.PRINT`** — args: string, any, any, any, any, any
- **`DEBUG.PRINT`** — args: string, any, any, any, any, any, any
- **`DEBUG.PRINT`** — args: string, any, any, any, any, any, any, any
- **`DEBUG.PRINT`** — args: string, any, any, any, any, any, any, any, any
- **`DEBUG.PRINT`** — args: string, any, any, any, any, any, any, any, any, any
- **`DEBUG.PRINT`** — args: string, any, any, any, any, any, any, any, any, any, any
- **`DEBUG.PRINTL`** — args: string, any
- **`DEBUG.PROFILEEND`** — args: string
- **`DEBUG.PROFILEREPORT`** — args: (none)
- **`DEBUG.PROFILESTART`** — args: string
- **`DEBUG.SHOWFPSGRAPH`** — args: bool — Show or hide the real-time FPS graph overlay.
- **`DEBUG.STACKTRACE`** — args: (none)
- **`DEBUG.WATCH`** — args: string, any
- **`DEBUG.WATCHCLEAR`** — args: (none)
- **`SYSTEM.MONITOR`** — args: (none) — Toggle the system performance monitor (FPS, RAM).
- **`SYSTEM.MONITOR`** — args: bool — Toggle the system performance monitor (FPS, RAM).

### TIMER

- **`TIMER`** — args: (none) → float
- **`TIMER`** — args: (none)
- **`TIMER.AFTER`** — args: float, string → int
- **`TIMER.CANCEL`** — args: int
- **`TIMER.CREATE`** — args: float → handle — Simulation timer: duration in seconds; use TIMER.START/UPDATE with delta time
- **`TIMER.DONE`** — args: handle → bool
- **`TIMER.EVERY`** — args: float, string → int
- **`TIMER.FINISHED`** — args: handle → bool
- **`TIMER.FRACTION`** — args: handle → float
- **`TIMER.FREE`** — args: handle
- **`TIMER.GETLOOP`** — args: handle → bool — Whether the sim timer (GameTimerSim) repeats after each cycle; last set with TIMER.SETLOOP.
- **`TIMER.MAKE`** — args: float → handle — DEPRECATED alias of TIMER.CREATE. Use TIMER.CREATE.
- **`TIMER.NEW`** — args: float → handle — Wall-clock deadline timer (time.Now-based)
- **`TIMER.REMAINING`** — args: handle → float
- **`TIMER.RESET`** — args: handle, float
- **`TIMER.REWIND`** — args: handle
- **`TIMER.SETLOOP`** — args: handle, any
- **`TIMER.START`** — args: handle
- **`TIMER.STOP`** — args: handle
- **`TIMER.UPDATE`** — args: handle, float

*61 overloads in this section.*

---

## Globals and language builtins

Guide: [11-TOOLING.md](11-TOOLING.md) · [LANGUAGE.md](../LANGUAGE.md)

### PRINT

- **`PRINT`** — args: any — Print values to stdout, space-separated, with newline.

### HELP

- **`HELP`** — args: string — Live Discovery: Show arguments and description for any command.

### ABS

- **`ABS`** — args: any

### SIN

- **`SIN`** — args: any

### COS

- **`COS`** — args: any

### RAND

- **`RAND`** — args: any, any → int — Same as RND(min, max) â€” inclusive integer range
- **`RAND`** — args: int, int → int — Easy Mode: Random int in range
- **`RAND.CREATE`** — args: int → handle
- **`RAND.FREE`** — args: handle
- **`RAND.MAKE`** — args: int → handle — DEPRECATED alias of RAND.CREATE. Use RAND.CREATE.
- **`RAND.NEXT`** — args: handle, int, int → int
- **`RAND.NEXTF`** — args: handle → float

### LEN

- **`LEN`** — args: string

### STR

- **`STR`** — args: any → string — Convert a value to string (canonical; same as legacy STR).

### VAL

- **`VAL`** — args: string → float

### CHR

- **`CHR`** — args: int → string

### ASC

- **`ASC`** — args: string

Language keywords (`IF`, `WHILE`, `FUNCTION`, `IMPORT "file.mb"`, …) are documented in [LANGUAGE.md](../LANGUAGE.md).

CLI (`moonbasic new`, `moonrun`, `moonbasic --check`, …): [11-TOOLING.md](11-TOOLING.md).

---

## All other engine namespaces

The full engine registers **930** dotted namespaces and thousands of overloads. Everything is listed in [API_CONSISTENCY.md](../API_CONSISTENCY.md).

Namespace → reference file map: [COMMAND_AUDIT.md](../COMMAND_AUDIT.md).

| Namespace | Overloads | Primary reference |
|-------------|----------:|-------------------|
| `AABBCOLLIDE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#aabbcollide) |
| `ACOS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#acos) |
| `ADDFORCE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#addforce) |
| `ADDIMPULSE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#addimpulse) |
| `ADDSURFACE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#addsurface) |
| `ALIGNTOVECTOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#aligntovector) |
| `AMBIENTLIGHT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#ambientlight) |
| `ANGLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#angle) |
| `ANGLEDIFF` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#anglediff) |
| `ANGLEDIFFRAD` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#anglediffrad) |
| `ANGLETO` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#angleto) |
| `ANIMLENGTH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#animlength) |
| `APP` | 10 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#app) |
| `APPLYENTITYFORCE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#applyentityforce) |
| `APPLYENTITYTORQUE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#applyentitytorque) |
| `APPROACH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#approach) |
| `ARGB` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#argb) |
| `ARGC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#argc) |
| `ARRAY` | 20 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#array) |
| `ARRAYFILL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#arrayfill) |
| `ARRAYFREE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#arrayfree) |
| `ARRAYJOINS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#arrayjoins) |
| `ARRAYLEN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#arraylen) |
| `ARRAYPUSH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#arraypush) |
| `ASIN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#asin) |
| `ASSERT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#assert) |
| `ATAN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#atan) |
| `ATAN2` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#atan2) |
| `ATLAS` | 3 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#atlas) |
| `ATN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#atn) |
| `AUDIO3D` | 5 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#audio3d) |
| `AVAILVIDMEM` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#availvidmem) |
| `AXIS` | 3 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#axis) |
| `ActiveShader` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#activeshader) |
| `AddTriangle` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#addtriangle) |
| `AddVertex` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#addvertex) |
| `AddWheel` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#addwheel) |
| `Animate` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#animate) |
| `AppTitle` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#apptitle) |
| `ApplyEntityImpulse` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#applyentityimpulse) |
| `BALL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#ball) |
| `BALLW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#ballw) |
| `BAND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#band) |
| `BANKSIZE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#banksize) |
| `BCLEAR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#bclear) |
| `BCOUNT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#bcount) |
| `BIN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#bin) |
| `BIOME` | 7 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#biome) |
| `BLSHIFT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#blshift) |
| `BNOT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#bnot) |
| `BOOL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#bool) |
| `BOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#bor) |
| `BOX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#box) |
| `BOX2D` | 5 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#box2d) |
| `BOXCOLLIDE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#boxcollide) |
| `BOXTOPLAND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#boxtopland) |
| `BOXW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#boxw) |
| `BRSHIFT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#brshift) |
| `BRUSHALPHA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#brushalpha) |
| `BRUSHBLEND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#brushblend) |
| `BRUSHCOLOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#brushcolor) |
| `BSET` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#bset) |
| `BTEST` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#btest) |
| `BTOGGLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#btoggle) |
| `BTREE` | 7 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#btree) |
| `BXOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#bxor) |
| `BrushFX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#brushfx) |
| `BrushShininess` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#brushshininess) |
| `BrushTexture` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#brushtexture) |
| `CAM` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cam) |
| `CAMERA2DOFFSET` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camera2doffset) |
| `CAMERA2DROTATION` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camera2drotation) |
| `CAMERA2DTARGET` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camera2dtarget) |
| `CAMERA2DZOOM` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camera2dzoom) |
| `CAMERAFOGCOLOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camerafogcolor) |
| `CAMERAFOGMODE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camerafogmode) |
| `CAMERAFOGRANGE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camerafogrange) |
| `CAMERAFOLLOW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camerafollow) |
| `CAMERAPICK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camerapick) |
| `CAMERAPROJECT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cameraproject) |
| `CAMERARANGE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camerarange) |
| `CAMERAVIEWPORT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cameraviewport) |
| `CAMERAZOOM` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camerazoom) |
| `CAP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cap) |
| `CAPSULE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#capsule) |
| `CAPW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#capw) |
| `CEIL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#ceil) |
| `CHAR` | 66 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#char) |
| `CHARACTER` | 11 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#character) |
| `CHARACTERREF` | 40 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#characterref) |
| `CHARCONTROLLER` | 27 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#charcontroller) |
| `CHECK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#check) |
| `CHOOSE` | 11 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#choose) |
| `CIRCLEBOXCOLLIDE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#circleboxcollide) |
| `CIRCLECOLLIDE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#circlecollide) |
| `CIRCLEPOINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#circlepoint) |
| `CLAMP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#clamp) |
| `CLAMPENTITY2D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#clampentity2d) |
| `CLEAR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#clear) |
| `CLEARWORLD` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#clearworld) |
| `CLIENT` | 6 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#client) |
| `CLIPBOARD` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#clipboard) |
| `CLOSEFILE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#closefile) |
| `CLOUD` | 7 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cloud) |
| `CLS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cls) |
| `COLLISIONENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisionentity) |
| `COLLISIONFORCE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisionforce) |
| `COLLISIONS` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisions) |
| `COLLISIONX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisionx) |
| `COLLISIONZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisionz) |
| `COLORPRINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#colorprint) |
| `COL_BLACK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_black) |
| `COL_BLUE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_blue) |
| `COL_CYAN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_cyan) |
| `COL_DARKGRAY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_darkgray) |
| `COL_GRAY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_gray) |
| `COL_GREEN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_green) |
| `COL_LIGHTGRAY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_lightgray) |
| `COL_MAGENTA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_magenta) |
| `COL_ORANGE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_orange) |
| `COL_RED` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_red) |
| `COL_TRANSPARENT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_transparent) |
| `COL_WHITE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_white) |
| `COL_YELLOW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#col_yellow) |
| `COMMAND` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#command) |
| `COMPUTESHADER` | 8 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#computeshader) |
| `CONNECT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#connect) |
| `CONTAINS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#contains) |
| `CONTROLLER` | 6 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#controller) |
| `COPYBANK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#copybank) |
| `COPYENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#copyentity) |
| `COPYFILE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#copyfile) |
| `COSD` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cosd) |
| `COUNT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#count) |
| `COUNTCHILDREN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#countchildren) |
| `COUNTCOLLISIONS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#countcollisions) |
| `COUNTTRIANGLES` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#counttriangles) |
| `COUNTVERTICES` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#countvertices) |
| `CREATEBANK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createbank) |
| `CREATEBODY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createbody) |
| `CREATEBODY2D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createbody2d) |
| `CREATECAMERA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createcamera) |
| `CREATECAMERA2D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createcamera2d) |
| `CREATECONE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createcone) |
| `CREATECUBE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createcube) |
| `CREATECYLINDER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createcylinder) |
| `CREATEEMITTER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createemitter) |
| `CREATELIGHT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createlight) |
| `CREATEMESH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createmesh) |
| `CREATEMIRROR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createmirror) |
| `CREATEPLANE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createplane) |
| `CREATESPHERE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createsphere) |
| `CREATESPRITE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createsprite) |
| `CREATESPRITE3D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createsprite3d) |
| `CREATETERRAIN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createterrain) |
| `CREATETEXTURE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createtexture) |
| `CREATEWORLD` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createworld) |
| `CUBE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cube) |
| `CULL` | 23 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cull) |
| `CURRENTDATE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#currentdate) |
| `CURRENTTIME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#currenttime) |
| `CURVE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#curve) |
| `CURVEANGLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#curveangle) |
| `CURVEVALUE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#curvevalue) |
| `CVDOUBLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cvdouble) |
| `CVFLOAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cvfloat) |
| `CVINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cvint) |
| `CVLONG` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cvlong) |
| `CVSHORT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cvshort) |
| `CameraFOV` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camerafov) |
| `CameraLookAt` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#cameralookat) |
| `CameraShake` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camerashake) |
| `CameraSmoothFollow` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#camerasmoothfollow) |
| `CollisionForce` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisionforce) |
| `CollisionNX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisionnx) |
| `CollisionNY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisionny) |
| `CollisionNZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisionnz) |
| `CollisionPX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisionpx) |
| `CollisionPY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisionpy) |
| `CollisionPZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisionpz) |
| `CollisionY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#collisiony) |
| `CountCollisions` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#countcollisions) |
| `CreateBrush` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createbrush) |
| `CreateCube` | 4 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createcube) |
| `CreateLight` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createlight) |
| `CreatePivot` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createpivot) |
| `CreatePointLight` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createpointlight) |
| `CreateSurface` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createsurface) |
| `CreateVehicle` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#createvehicle) |
| `DATA` | 10 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#data) |
| `DATE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#date) |
| `DATETIME` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#datetime) |
| `DAY` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#day) |
| `DB` | 14 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#db) |
| `DECAL` | 20 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#decal) |
| `DEG2RAD` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#deg2rad) |
| `DEGPERSEC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#degpersec) |
| `DELAY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#delay) |
| `DELETEDIR` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#deletedir) |
| `DELETEFILE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#deletefile) |
| `DELTATIME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#deltatime) |
| `DIREXISTS` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#direxists) |
| `DIST2D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#dist2d) |
| `DIST3D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#dist3d) |
| `DISTANCE2D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#distance2d) |
| `DISTANCE3D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#distance3d) |
| `DISTANCESQ2D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#distancesq2d) |
| `DISTANCESQ3D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#distancesq3d) |
| `DISTSQ2D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#distsq2d) |
| `DRAW3D` | 17 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#draw3d) |
| `DRAWBBOX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawbbox) |
| `DRAWBILLBOARD` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawbillboard) |
| `DRAWBILLBOARDREC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawbillboardrec) |
| `DRAWCAP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawcap) |
| `DRAWCAPW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawcapw) |
| `DRAWCIRCLE2` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawcircle2) |
| `DRAWCIRCLE2W` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawcircle2w) |
| `DRAWCUBE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawcube) |
| `DRAWCUBEWIRES` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawcubewires) |
| `DRAWCYLINDER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawcylinder) |
| `DRAWCYLINDERW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawcylinderw) |
| `DRAWELLIPSE2` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawellipse2) |
| `DRAWELLIPSE2W` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawellipse2w) |
| `DRAWEMITTER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawemitter) |
| `DRAWGRID3D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawgrid3d) |
| `DRAWLINE2` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawline2) |
| `DRAWLINE3D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawline3d) |
| `DRAWPLANE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawplane) |
| `DRAWPOINT3D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawpoint3d) |
| `DRAWPOLY2` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawpoly2) |
| `DRAWPOLY2W` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawpoly2w) |
| `DRAWPRIM2D` | 14 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawprim2d) |
| `DRAWPRIM3D` | 17 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawprim3d) |
| `DRAWRAY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawray) |
| `DRAWRECT2` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawrect2) |
| `DRAWRECT2W` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawrect2w) |
| `DRAWRING2` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawring2) |
| `DRAWRING2W` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawring2w) |
| `DRAWSPHERE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawsphere) |
| `DRAWSPHEREW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawspherew) |
| `DRAWTEX2` | 7 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawtex2) |
| `DRAWTEXPRO` | 10 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawtexpro) |
| `DRAWTEXREC` | 8 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawtexrec) |
| `DRAWTRI2` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawtri2) |
| `DRAWTRI2W` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawtri2w) |
| `DUMP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#dump) |
| `DrawEntities` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawentities) |
| `DrawEntity` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#drawentity) |
| `E` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#e) |
| `EASEIN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easein) |
| `EASEIN3` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easein3) |
| `EASEINBACK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeinback) |
| `EASEINBOUNCE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeinbounce) |
| `EASEINELASTIC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeinelastic) |
| `EASEINOUT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeinout) |
| `EASEINOUT3` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeinout3) |
| `EASEINOUTSINE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeinoutsine) |
| `EASEINSINE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeinsine) |
| `EASELERP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easelerp) |
| `EASEOUT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeout) |
| `EASEOUT3` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeout3) |
| `EASEOUTBACK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeoutback) |
| `EASEOUTBOUNCE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeoutbounce) |
| `EASEOUTELASTIC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeoutelastic) |
| `EASEOUTSINE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#easeoutsine) |
| `EFFECT` | 24 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#effect) |
| `ELAPSED` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#elapsed) |
| `EMITPARTICLE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#emitparticle) |
| `EMITTERALIVE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#emitteralive) |
| `EMITTERCOUNT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#emittercount) |
| `EMITTERPOS` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#emitterpos) |
| `ENDGAME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#endgame) |
| `ENDSWITH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#endswith) |
| `ENEMY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#enemy) |
| `ENET` | 8 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#enet) |
| `ENT` | 17 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#ent) |
| `ENTHIT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#enthit) |
| `ENTITYALPHA` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityalpha) |
| `ENTITYANIMATETOWARD` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityanimatetoward) |
| `ENTITYAUTOFADE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityautofade) |
| `ENTITYBLEND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityblend) |
| `ENTITYBOX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitybox) |
| `ENTITYCOLLIDED` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitycollided) |
| `ENTITYCOLOR` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitycolor) |
| `ENTITYDISTANCE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitydistance) |
| `ENTITYFLOOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityfloor) |
| `ENTITYFX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityfx) |
| `ENTITYINVIEW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityinview) |
| `ENTITYJUMP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityjump) |
| `ENTITYNAME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityname) |
| `ENTITYORDER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityorder) |
| `ENTITYPARENT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityparent) |
| `ENTITYPHYSICSTOUCH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityphysicstouch) |
| `ENTITYPICK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitypick) |
| `ENTITYPITCH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitypitch) |
| `ENTITYRADIUS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityradius) |
| `ENTITYREF` | 3 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityref) |
| `ENTITYROLL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityroll) |
| `ENTITYSCALEX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityscalex) |
| `ENTITYSCALEY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityscaley) |
| `ENTITYSCALEZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityscalez) |
| `ENTITYSHININESS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityshininess) |
| `ENTITYTEXTURE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitytexture) |
| `ENTITYTYPE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitytype) |
| `ENTITYVISIBLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityvisible) |
| `ENTITYX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityx) |
| `ENTITYY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityy) |
| `ENTITYYAW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityyaw) |
| `ENTITYZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityz) |
| `ENTP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entp) |
| `ENTPITCH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entpitch) |
| `ENTR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entr) |
| `ENTRAD` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entrad) |
| `ENTROLL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entroll) |
| `ENTTYPE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#enttype) |
| `ENTW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entw) |
| `ENTX` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entx) |
| `ENTY` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#enty) |
| `ENTYAW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entyaw) |
| `ENTZ` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entz) |
| `ENVIRON` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#environ) |
| `EOF` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#eof) |
| `ERASE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#erase) |
| `ERR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#err) |
| `ERRLINE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#errline) |
| `EVENT` | 28 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#event) |
| `EXP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#exp) |
| `EmitSound` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#emitsound) |
| `EntityAnimTime` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityanimtime) |
| `EntityAnimateToward` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityanimatetoward) |
| `EntityApplyImpulse` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityapplyimpulse) |
| `EntityCanSee` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitycansee) |
| `EntityCheckCollision` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitycheckcollision) |
| `EntityCollided` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitycollided) |
| `EntityCollisionLayer` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitycollisionlayer) |
| `EntityEmission` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityemission) |
| `EntityFriction` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityfriction) |
| `EntityGetClosestWithTag` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitygetclosestwithtag) |
| `EntityGetGroundNormal` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitygetgroundnormal) |
| `EntityGetOverlapCount` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitygetoverlapcount) |
| `EntityGrounded` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitygrounded) |
| `EntityHasTag` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityhastag) |
| `EntityHitsType` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityhitstype) |
| `EntityInFrustum` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityinfrustum) |
| `EntityLineOfSight` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitylineofsight) |
| `EntityMass` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitymass) |
| `EntityMoveCameraRelative` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitymovecamerarelative) |
| `EntityNormalMap` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitynormalmap) |
| `EntityPBR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitypbr) |
| `EntityPushOutOfGeometry` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitypushoutofgeometry) |
| `EntityRaycast` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityraycast) |
| `EntityRestitution` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityrestitution) |
| `EntitySetCollisionGroup` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entitysetcollisiongroup) |
| `EntityShadow` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#entityshadow) |
| `ExtractAnimSeq` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#extractanimseq) |
| `FBMNOISE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#fbmnoise) |
| `FILEEXISTS` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#fileexists) |
| `FILEPOS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#filepos) |
| `FILESIZE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#filesize) |
| `FINDCHILD` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#findchild) |
| `FINISH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#finish) |
| `FITMESH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#fitmesh) |
| `FIX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#fix) |
| `FLAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#flat) |
| `FLIPMESH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#flipmesh) |
| `FLOAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#float) |
| `FLOOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#floor) |
| `FLUSHKEYS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#flushkeys) |
| `FLUSHMOUSE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#flushmouse) |
| `FOG` | 5 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#fog) |
| `FOGCOLOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#fogcolor) |
| `FOGDENSITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#fogdensity) |
| `FOGMODE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#fogmode) |
| `FORMAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#format) |
| `FORMATINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#formatint) |
| `FORMATSCORE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#formatscore) |
| `FORMATTIME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#formattime) |
| `FORMATTIME2` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#formattime2) |
| `FPS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#fps) |
| `FRAMECOUNT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#framecount) |
| `FREE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#free) |
| `FREEBANK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#freebank) |
| `FREEBRUSH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#freebrush) |
| `FREEENTITIES` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#freeentities) |
| `FREEENTITY` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#freeentity) |
| `FREESOUND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#freesound) |
| `FREETEXTURE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#freetexture) |
| `FindBone` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#findbone) |
| `GAME` | 50 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#game) |
| `GAMEPAUSE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#gamepause) |
| `GETCHILD` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getchild) |
| `GETCOLLISIONENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getcollisionentity) |
| `GETDIR` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getdir) |
| `GETDIRS` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getdirs) |
| `GETDROPPEDFILES` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getdroppedfiles) |
| `GETENTITYBRUSH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getentitybrush) |
| `GETFILEEXT` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getfileext) |
| `GETFILEMODTIME` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getfilemodtime) |
| `GETFILENAME` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getfilename) |
| `GETFILENAMENOEXT` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getfilenamenoext) |
| `GETFILEPATH` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getfilepath) |
| `GETFILES` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getfiles) |
| `GETFILESIZE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getfilesize) |
| `GETJOY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getjoy) |
| `GETKEY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getkey) |
| `GETPARENT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getparent) |
| `GETSURFACEBRUSH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#getsurfacebrush) |
| `GETTEXTCODEPOINTCOUNT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#gettextcodepointcount) |
| `GPUNAME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#gpuname) |
| `GRAPHICS` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#graphics) |
| `GRAPHICSDEPTH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#graphicsdepth) |
| `GRAPHICSHEIGHT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#graphicsheight) |
| `GRAPHICSWIDTH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#graphicswidth) |
| `GRID` | 14 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#grid) |
| `GRID3` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#grid3) |
| `Graphics3D` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#graphics3d) |
| `HASHFLOAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#hashfloat) |
| `HASHINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#hashint) |
| `HDIST` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#hdist) |
| `HDISTSQ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#hdistsq) |
| `HEX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#hex) |
| `HIDEENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#hideentity) |
| `HIDEPOINTER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#hidepointer) |
| `HITCOUNT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#hitcount) |
| `HITENT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#hitent) |
| `HOUR` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#hour) |
| `IIF` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#iif) |
| `IMAGE` | 63 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#image) |
| `INSTANCE` | 24 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#instance) |
| `INSTR` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#instr) |
| `INT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#int) |
| `INTERP` | 10 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#interp) |
| `INVERSE_LERP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#inverse_lerp) |
| `ISALPHA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#isalpha) |
| `ISALPHANUM` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#isalphanum) |
| `ISFILEDROPPED` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#isfiledropped) |
| `ISHANDLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#ishandle) |
| `ISNULL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#isnull) |
| `ISNUMERIC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#isnumeric) |
| `ISTYPE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#istype) |
| `JOIN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#join) |
| `JOINT` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joint) |
| `JOINT2D` | 4 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joint2d) |
| `JOINT3D` | 5 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joint3d) |
| `JOLT` | 19 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#jolt) |
| `JOYDOWN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joydown) |
| `JOYHAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joyhat) |
| `JOYHIT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joyhit) |
| `JOYPITCH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joypitch) |
| `JOYROLL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joyroll) |
| `JOYU` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joyu) |
| `JOYV` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joyv) |
| `JOYX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joyx) |
| `JOYXDIR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joyxdir) |
| `JOYY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joyy) |
| `JOYYAW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joyyaw) |
| `JOYYDIR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joyydir) |
| `JOYZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#joyz) |
| `KEEPPLAYERINBOUNDS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#keepplayerinbounds) |
| `KEY` | 4 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#key) |
| `KEYDOWN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#keydown) |
| `KEYHIT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#keyhit) |
| `KEYUP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#keyup) |
| `KINEMATIC` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#kinematic) |
| `KINEMATICREF` | 3 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#kinematicref) |
| `KeyDown` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#keydown) |
| `LANDBOX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#landbox) |
| `LANDBOXES` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#landboxes) |
| `LEFT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#left) |
| `LERP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#lerp) |
| `LEVEL` | 16 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#level) |
| `LIGHTCONE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#lightcone) |
| `LIGHTPOINTAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#lightpointat) |
| `LIGHTPOSITION` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#lightposition) |
| `LINE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#line) |
| `LINE3D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#line3d) |
| `LINECOLLIDE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#linecollide) |
| `LINEPICK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#linepick) |
| `LISTEN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#listen) |
| `LOADANIMTEXTURE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loadanimtexture) |
| `LOADBRUSH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loadbrush) |
| `LOADFONT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loadfont) |
| `LOADIMAGE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loadimage) |
| `LOADMESH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loadmesh) |
| `LOADMUSIC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loadmusic) |
| `LOADSOUND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loadsound) |
| `LOADSPRITE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loadsprite) |
| `LOADTERRAIN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loadterrain) |
| `LOADTEXTURE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loadtexture) |
| `LOBBY` | 9 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#lobby) |
| `LOCATE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#locate) |
| `LOG` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#log) |
| `LOG10` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#log10) |
| `LOG2` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#log2) |
| `LOOPSOUND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loopsound) |
| `LOWER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#lower) |
| `LSET` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#lset) |
| `LTRIM` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#ltrim) |
| `LightColor` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#lightcolor) |
| `LightRange` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#lightrange) |
| `Listener` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#listener) |
| `Load3DSound` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#load3dsound) |
| `LoadAnimMesh` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#loadanimmesh) |
| `MAKEDIR` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#makedir) |
| `MAKEDIRS` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#makedirs) |
| `MAT4` | 17 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mat4) |
| `MATRIX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#matrix) |
| `MAX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#max) |
| `MEASURETEXT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#measuretext) |
| `MEASURETEXTEX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#measuretextex) |
| `MEM` | 19 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mem) |
| `MESHDEPTH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#meshdepth) |
| `MESHHEIGHT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#meshheight) |
| `MESHWIDTH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#meshwidth) |
| `MID` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mid) |
| `MILLISECOND` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#millisecond) |
| `MILLISECS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#millisecs) |
| `MIN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#min) |
| `MINUTE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#minute) |
| `MKDOUBLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mkdouble) |
| `MKFLOAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mkfloat) |
| `MKINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mkint) |
| `MKLONG` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mklong) |
| `MKSHORT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mkshort) |
| `MODIFYTERRAIN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#modifyterrain) |
| `MONTH` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#month) |
| `MOUSE` | 17 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mouse) |
| `MOUSEDOWN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mousedown) |
| `MOUSEDX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mousedx) |
| `MOUSEDY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mousedy) |
| `MOUSEHIT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mousehit) |
| `MOUSEWHEEL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mousewheel) |
| `MOUSEX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mousex) |
| `MOUSEY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mousey) |
| `MOUSEZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mousez) |
| `MOVE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#move) |
| `MOVECAMERA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#movecamera) |
| `MOVEENTITY` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#moveentity) |
| `MOVEENTITY2D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#moveentity2d) |
| `MOVEFILE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#movefile) |
| `MOVEMOUSE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#movemouse) |
| `MOVEPLAYER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#moveplayer) |
| `MOVER` | 7 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mover) |
| `MOVESPRITE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#movesprite) |
| `MOVESTEPX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#movestepx) |
| `MOVESTEPZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#movestepz) |
| `MOVEX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#movex) |
| `MOVEZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#movez) |
| `MUSIC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#music) |
| `MUSICVOLUME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#musicvolume) |
| `MilliSecs` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#millisecs) |
| `MouseWheel` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#mousewheel) |
| `MoveEntity` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#moveentity) |
| `NAMEENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#nameentity) |
| `NAV` | 23 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#nav) |
| `NAVAGENT` | 31 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#navagent) |
| `NET` | 22 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#net) |
| `NETREADFLOAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#netreadfloat) |
| `NETREADINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#netreadint) |
| `NETREADSTRING` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#netreadstring) |
| `NETSENDFLOAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#netsendfloat) |
| `NETSENDINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#netsendint) |
| `NETSENDSTRING` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#netsendstring) |
| `NET_RELIABLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#net_reliable) |
| `NET_UNRELIABLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#net_unreliable) |
| `NEWXVALUE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#newxvalue) |
| `NEWYVALUE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#newyvalue) |
| `NEWZVALUE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#newzvalue) |
| `NOISE` | 34 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#noise) |
| `OCT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#oct) |
| `OPENFILE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#openfile) |
| `ORBITDISTDELTA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#orbitdistdelta) |
| `ORBITPITCHDELTA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#orbitpitchdelta) |
| `ORBITYAWDELTA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#orbityawdelta) |
| `OSCILLATE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#oscillate) |
| `OVAL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#oval) |
| `PACKET` | 4 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#packet) |
| `PAINTSURFACE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#paintsurface) |
| `PARTICLE2D` | 6 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#particle2d) |
| `PARTICLE3D` | 34 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#particle3d) |
| `PARTICLECOLOR` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#particlecolor) |
| `PARTICLEEMITRATE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#particleemitrate) |
| `PARTICLELIFE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#particlelife) |
| `PARTICLES` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#particles) |
| `PARTICLESPEED` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#particlespeed) |
| `PARTICLEVELOCITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#particlevelocity) |
| `PATH` | 6 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#path) |
| `PAUSEGAME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pausegame) |
| `PEEKBYTE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#peekbyte) |
| `PEEKFLOAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#peekfloat) |
| `PEEKINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#peekint) |
| `PEEKSHORT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#peekshort) |
| `PEER` | 5 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#peer) |
| `PERLIN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#perlin) |
| `PHYSICS2D` | 11 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#physics2d) |
| `PHYSICSCOLLISIONFORCE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#physicscollisionforce) |
| `PHYSICSCOLLISIONNX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#physicscollisionnx) |
| `PHYSICSCOLLISIONNY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#physicscollisionny) |
| `PHYSICSCOLLISIONNZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#physicscollisionnz) |
| `PHYSICSCOLLISIONPX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#physicscollisionpx) |
| `PHYSICSCOLLISIONPY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#physicscollisionpy) |
| `PHYSICSCOLLISIONPZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#physicscollisionpz) |
| `PHYSICSCOLLISIONY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#physicscollisiony) |
| `PHYSICSCONTACTCOUNT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#physicscontactcount) |
| `PI` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pi) |
| `PICKEDDISTANCE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pickeddistance) |
| `PICKEDENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pickedentity) |
| `PICKEDNX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pickednx) |
| `PICKEDNY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pickedny) |
| `PICKEDNZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pickednz) |
| `PICKEDSURFACE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pickedsurface) |
| `PICKEDTRIANGLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pickedtriangle) |
| `PICKEDX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pickedx) |
| `PICKEDY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pickedy) |
| `PICKEDZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pickedz) |
| `PINGPONG` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pingpong) |
| `PLAYER` | 116 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#player) |
| `PLAYER2D` | 11 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#player2d) |
| `PLAYMUSIC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#playmusic) |
| `PLAYSOUND` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#playsound) |
| `PLOT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#plot) |
| `POINT3D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#point3d) |
| `POINTDIR2D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pointdir2d) |
| `POINTDIR3D` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pointdir3d) |
| `POINTENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pointentity) |
| `POINTINAABB` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pointinaabb) |
| `POINTINBOX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pointinbox) |
| `POINTINCIRCLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pointincircle) |
| `POINTONLINE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pointonline) |
| `POKEBYTE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pokebyte) |
| `POKEFLOAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pokefloat) |
| `POKEINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pokeint) |
| `POKESHORT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pokeshort) |
| `POOL` | 8 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pool) |
| `POSENT` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#posent) |
| `POSITIONCAMERA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#positioncamera) |
| `POSITIONENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#positionentity) |
| `POSITIONTEXTURE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#positiontexture) |
| `POST` | 8 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#post) |
| `POW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pow) |
| `PP_BLOOM` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pp_bloom) |
| `PP_CRT_SCANLINES` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pp_crt_scanlines) |
| `PP_PIXELATE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#pp_pixelate) |
| `PRINTAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#printat) |
| `PRINTCOLOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#printcolor) |
| `PRINTLN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#println) |
| `PROP` | 3 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#prop) |
| `PaintEntity` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#paintentity) |
| `QUAT` | 13 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#quat) |
| `QUIT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#quit) |
| `RAD2DEG` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rad2deg) |
| `RANDOMELEMENT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#randomelement) |
| `RANDOMIZE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#randomize) |
| `RAYLIB` | 33 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#raylib) |
| `READALLTEXT` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#readalltext) |
| `READBANK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#readbank) |
| `READBYTE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#readbyte) |
| `READFILE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#readfile) |
| `READFLOAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#readfloat) |
| `READINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#readint) |
| `READLINE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#readline) |
| `READSHORT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#readshort) |
| `READSTRING` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#readstring) |
| `RECT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rect) |
| `REMAP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#remap) |
| `RENAMEFILE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#renamefile) |
| `RENDERTARGET` | 6 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rendertarget) |
| `REPEAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#repeat) |
| `REPLACE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#replace) |
| `RES` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#res) |
| `RESETENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#resetentity) |
| `RESIZEBANK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#resizebank) |
| `RESUMEGAME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#resumegame) |
| `REVERSE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#reverse) |
| `RGB` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rgb) |
| `RGBA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rgba) |
| `RGBB` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rgbb) |
| `RGBBRIGHTEN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rgbbrighten) |
| `RGBDARKEN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rgbdarken) |
| `RGBFADE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rgbfade) |
| `RGBG` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rgbg) |
| `RGBMIX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rgbmix) |
| `RGBR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rgbr) |
| `RIGHT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#right) |
| `RND` | 4 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rnd) |
| `RNDF` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rndf) |
| `RNDRANGE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rndrange) |
| `RNDSEED` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rndseed) |
| `ROTATECAMERA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rotatecamera) |
| `ROTATEENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rotateentity) |
| `ROTATETEXTURE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rotatetexture) |
| `ROTENT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rotent) |
| `ROUND` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#round) |
| `ROWS` | 5 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rows) |
| `RPC` | 24 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rpc) |
| `RSET` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rset) |
| `RTRIM` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#rtrim) |
| `SATURATE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#saturate) |
| `SCALEENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#scaleentity) |
| `SCALENT` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#scalent) |
| `SCALESPRITE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#scalesprite) |
| `SCALETEXTURE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#scaletexture) |
| `SCATTER` | 5 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#scatter) |
| `SCREENHEIGHT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#screenheight) |
| `SCREENWIDTH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#screenwidth) |
| `SECOND` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#second) |
| `SEEDRND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#seedrnd) |
| `SEEKFILE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#seekfile) |
| `SERVER` | 8 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#server) |
| `SERVICENET` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#servicenet) |
| `SETALPHA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setalpha) |
| `SETAMBIENT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setambient) |
| `SETBLOOM` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setbloom) |
| `SETCLEARCOLOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setclearcolor) |
| `SETCOLOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setcolor) |
| `SETCUBEFACE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setcubeface) |
| `SETCUBEMODE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setcubemode) |
| `SETDIR` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setdir) |
| `SETFOG` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setfog) |
| `SETFPS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setfps) |
| `SETGRAVITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setgravity) |
| `SETORIGIN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setorigin) |
| `SETVIEWPORT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setviewport) |
| `SETVSYNC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setvsync) |
| `SETWIREFRAME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setwireframe) |
| `SGN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#sgn) |
| `SHADER` | 13 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#shader) |
| `SHADER_CEL_STYLED` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#shader_cel_styled) |
| `SHADER_PBR_LIT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#shader_pbr_lit) |
| `SHADER_PS1_RETRO` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#shader_ps1_retro) |
| `SHADER_WATER_PROCEDURAL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#shader_water_procedural) |
| `SHAKECAMERA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#shakecamera) |
| `SHAPE` | 8 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#shape) |
| `SHOWENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#showentity) |
| `SHOWPOINTER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#showpointer) |
| `SHUFFLE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#shuffle) |
| `SIGN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#sign) |
| `SIMPLEX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#simplex) |
| `SIND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#sind) |
| `SKY` | 9 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#sky) |
| `SKYCOLOR` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#skycolor) |
| `SLEEP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#sleep) |
| `SMOOTHERSTEP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#smootherstep) |
| `SMOOTHSTEP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#smoothstep) |
| `SOUNDPAN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#soundpan) |
| `SOUNDVOLUME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#soundvolume) |
| `SPACE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#space) |
| `SPAWNER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spawner) |
| `SPC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spc) |
| `SPHERE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#sphere) |
| `SPHEREBOXCOLLIDE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#sphereboxcollide) |
| `SPHERECOLLIDE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spherecollide) |
| `SPLIT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#split) |
| `SPRITEALPHA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spritealpha) |
| `SPRITEBATCH` | 6 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spritebatch) |
| `SPRITECOLLIDE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spritecollide) |
| `SPRITECOLOR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spritecolor) |
| `SPRITEGROUP` | 7 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spritegroup) |
| `SPRITEHIT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spritehit) |
| `SPRITEIMAGE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spriteimage) |
| `SPRITELAYER` | 7 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spritelayer) |
| `SPRITEMODE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spritemode) |
| `SPRITEUI` | 4 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spriteui) |
| `SPRITEVIEWMODE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#spriteviewmode) |
| `SQR` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#sqr) |
| `SQRT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#sqrt) |
| `STARTSWITH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#startswith) |
| `STATIC` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#static) |
| `STEER` | 10 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#steer) |
| `STOP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#stop) |
| `STOPMUSIC` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#stopmusic) |
| `STOPSOUND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#stopsound) |
| `STOPWATCH` | 4 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#stopwatch) |
| `STRING` | 11 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#string) |
| `SWITCH` | 8 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#switch) |
| `SYSTEMPROPERTY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#systemproperty) |
| `SetAnimTime` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setanimtime) |
| `SetMSAA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setmsaa) |
| `SetPostProcess` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setpostprocess) |
| `SetSSAO` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#setssao) |
| `SoundPitch` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#soundpitch) |
| `SoundVolume` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#soundvolume) |
| `TAB` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#tab) |
| `TABLE` | 12 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#table) |
| `TAN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#tan) |
| `TAND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#tand) |
| `TAU` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#tau) |
| `TERRAINDETAIL` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#terraindetail) |
| `TERRAINHEIGHT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#terrainheight) |
| `TERRAINSHADING` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#terrainshading) |
| `TERRAINSIZE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#terrainsize) |
| `TERRAINX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#terrainx) |
| `TERRAINZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#terrainz) |
| `TEXTDRAW` | 7 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#textdraw) |
| `TEXTEXOBJ` | 7 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#textexobj) |
| `TEXTOBJ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#textobj) |
| `TEXTOBJEX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#textobjex) |
| `TEXTURECOORDS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#texturecoords) |
| `TEXTUREHEIGHT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#textureheight) |
| `TEXTURENAME` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#texturename) |
| `TEXTUREWIDTH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#texturewidth) |
| `TFormVector` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#tformvector) |
| `THROW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#throw) |
| `TICKCOUNT` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#tickcount) |
| `TIMEMS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#timems) |
| `TIMESTAMP` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#timestamp) |
| `TOTALVIDMEM` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#totalvidmem) |
| `TRACE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#trace) |
| `TRANSFORM` | 16 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#transform) |
| `TRANSITION` | 5 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#transition) |
| `TRIGGER` | 6 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#trigger) |
| `TRIM` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#trim) |
| `TURNCAMERA` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#turncamera) |
| `TURNENTITY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#turnentity) |
| `TWEEN` | 16 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#tween) |
| `TYPEOF` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#typeof) |
| `TranslateEntity` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#translateentity) |
| `UI` | 4 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#ui) |
| `UPDATEEMITTER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#updateemitter) |
| `UPDATENORMALS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#updatenormals) |
| `UPDATEPHYSICS` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#updatephysics) |
| `UPDW` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#updw) |
| `UPPER` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#upper) |
| `UTIL` | 25 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#util) |
| `UpdateMesh` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#updatemesh) |
| `VERTEXNX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#vertexnx) |
| `VERTEXNY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#vertexny) |
| `VERTEXNZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#vertexnz) |
| `VERTEXU` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#vertexu) |
| `VERTEXV` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#vertexv) |
| `VORONOI` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#voronoi) |
| `VertexX` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#vertexx) |
| `VertexY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#vertexy) |
| `VertexZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#vertexz) |
| `WAIT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#wait) |
| `WAITKEY` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#waitkey) |
| `WAITMOUSE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#waitmouse) |
| `WATER` | 32 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#water) |
| `WAVE` | 6 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#wave) |
| `WEATHER` | 8 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#weather) |
| `WEIGHTEDRND` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#weightedrnd) |
| `WIND` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#wind) |
| `WINDOWHEIGHT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#windowheight) |
| `WINDOWWIDTH` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#windowwidth) |
| `WIRECUBE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#wirecube) |
| `WRAP` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#wrap) |
| `WRAPANGLE` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#wrapangle) |
| `WRAPANGLE180` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#wrapangle180) |
| `WRAPVALUE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#wrapvalue) |
| `WRITE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#write) |
| `WRITEALLTEXT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#writealltext) |
| `WRITEBANK` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#writebank) |
| `WRITEBYTE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#writebyte) |
| `WRITEFILE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#writefile) |
| `WRITEFILELN` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#writefileln) |
| `WRITEFLOAT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#writefloat) |
| `WRITEINT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#writeint) |
| `WRITELINE` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#writeline) |
| `WRITESHORT` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#writeshort) |
| `WRITESTRING` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#writestring) |
| `YAWFROMXZ` | 1 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#yawfromxz) |
| `YEAR` | 2 | [API_CONSISTENCY.md](../API_CONSISTENCY.md#year) |

*Beginner-system overloads above: 1917 · Other namespaces: 2238 · See [API_CONSISTENCY.md](../API_CONSISTENCY.md) for the complete machine-readable list.*
