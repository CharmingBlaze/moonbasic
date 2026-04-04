# moonBASIC API consistency

This document is generated from `compiler/builtinmanifest/commands.json`.

Refresh: `go run ./tools/apidoc > docs/API_CONSISTENCY.md` (from the repository root).

## Naming conventions

- **Registry / source form**: `NS.ACTION` in uppercase with a dot (e.g. `CAMERA.SETPOS`).
- **Handle methods** (on a handle value): `cam.SetPos` dispatches to `CAMERA.SETPOS`. **`SetPosition`** is an alias for **`SetPos`** on spatial types.
- **Spatial handles** (`Camera3D`, `Body3D`, `CharController`, `Model`, `Sprite`): use **`SetPos`** / `SETPOS` (sprite: x,y; others: x,y,z).
- **`MODEL.SETPOS`**: sets the model root transform to a **translation matrix** (replaces prior rotation/scale on that matrix).
- **`LIGHT.MAKE`**: zero arguments -> directional, white, intensity 1.0 (config handle for future lighting).
- **`BODY3D.MAKE`**: zero arguments -> **DYNAMIC** motion builder.
- **Errors**: include **file and line** when available; unknown commands use **did-you-mean** against the live registry.

## Commands by namespace

### ABS

- **`ABS`** - args: any

### ACOS

- **`ACOS`** - args: any

### ANGLEDIFF

- **`ANGLEDIFF`** - args: any, any

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

### ATN

- **`ATN`** - args: any

### AUDIO

- **`AUDIO.CLOSE`** - args: (none)
- **`AUDIO.INIT`** - args: (none)
- **`AUDIO.LOADMUSIC`** - args: string
- **`AUDIO.LOADSOUND`** - args: string
- **`AUDIO.PAUSE`** - args: (none)
- **`AUDIO.PLAY`** - args: any
- **`AUDIO.RESUME`** - args: (none)
- **`AUDIO.STOP`** - args: (none)

### AUDIOSTREAM

- **`AUDIOSTREAM.FREE`** - args: handle
- **`AUDIOSTREAM.ISPLAYING`** - args: handle -> returns bool
- **`AUDIOSTREAM.ISREADY`** - args: handle -> returns bool
- **`AUDIOSTREAM.MAKE`** - args: int, int, int -> returns handle
- **`AUDIOSTREAM.PAUSE`** - args: handle
- **`AUDIOSTREAM.PLAY`** - args: handle
- **`AUDIOSTREAM.RESUME`** - args: handle
- **`AUDIOSTREAM.SETPAN`** - args: handle, float
- **`AUDIOSTREAM.SETPITCH`** - args: handle, float
- **`AUDIOSTREAM.SETVOLUME`** - args: handle, float
- **`AUDIOSTREAM.STOP`** - args: handle
- **`AUDIOSTREAM.UPDATE`** - args: handle, handle

### BAND

- **`BAND`** - args: any, any

### BBOX

- **`BBOX.CHECK`** - args: handle, handle -> returns bool
- **`BBOX.CHECKSPHERE`** - args: handle, float, float, float, float -> returns bool
- **`BBOX.FREE`** - args: handle
- **`BBOX.FROMMODEL`** - args: handle -> returns handle
- **`BBOX.MAKE`** - args: float, float, float, float, float, float -> returns handle

### BCLEAR

- **`BCLEAR`** - args: any, int

### BCOUNT

- **`BCOUNT`** - args: any

### BIN$

- **`BIN$`** - args: int

### BLSHIFT

- **`BLSHIFT`** - args: any, int

### BNOT

- **`BNOT`** - args: any

### BODY2D

- **`BODY2D.ADDCIRCLE`** - args: handle, float
- **`BODY2D.ADDRECT`** - args: handle, float, float
- **`BODY2D.APPLYFORCE`** - args: handle, float, float
- **`BODY2D.APPLYIMPULSE`** - args: handle, float, float
- **`BODY2D.COMMIT`** - args: handle, float, float -> returns handle
- **`BODY2D.FREE`** - args: handle
- **`BODY2D.GETPOS`** - args: handle -> returns handle
- **`BODY2D.GETROT`** - args: handle -> returns float
- **`BODY2D.MAKE`** - args: string -> returns handle
- **`BODY2D.ROT`** - args: handle -> returns float
- **`BODY2D.SETFRICTION`** - args: handle, float
- **`BODY2D.SETMASS`** - args: handle, float
- **`BODY2D.SETPOS`** - args: handle, float, float
- **`BODY2D.SETRESTITUTION`** - args: handle, float
- **`BODY2D.SETROT`** - args: handle, float
- **`BODY2D.X`** - args: handle -> returns float
- **`BODY2D.Y`** - args: handle -> returns float

### BODY3D

- **`BODY3D.ADDBOX`** - args: handle, float, float, float
- **`BODY3D.ADDCAPSULE`** - args: handle, float, float
- **`BODY3D.ADDMESH`** - args: handle, handle
- **`BODY3D.ADDSPHERE`** - args: handle, float
- **`BODY3D.APPLYFORCE`** - args: handle, float, float, float
- **`BODY3D.APPLYIMPULSE`** - args: handle, float, float, float
- **`BODY3D.COMMIT`** - args: handle, float, float, float -> returns handle
- **`BODY3D.FREE`** - args: handle
- **`BODY3D.GETPOS`** - args: handle -> returns handle
- **`BODY3D.GETROT`** - args: handle -> returns handle
- **`BODY3D.MAKE`** - args: (none) -> returns handle
- **`BODY3D.MAKE`** - args: (none) -> returns handle
- **`BODY3D.MAKE`** - args: string -> returns handle
- **`BODY3D.MAKE`** - args: string
- **`BODY3D.SETANGULARVEL`** - args: handle, float, float, float
- **`BODY3D.SETFRICTION`** - args: handle, float
- **`BODY3D.SETLINEARVEL`** - args: handle, float, float, float
- **`BODY3D.SETMASS`** - args: handle, float
- **`BODY3D.SETPOS`** - args: handle, float, float, float
- **`BODY3D.SETPOSITION`** - args: handle, float, float, float
- **`BODY3D.SETRESTITUTION`** - args: handle, float
- **`BODY3D.SETROT`** - args: handle, float, float, float
- **`BODY3D.X`** - args: handle -> returns float
- **`BODY3D.Y`** - args: handle -> returns float
- **`BODY3D.Z`** - args: handle -> returns float

### BOOL

- **`BOOL`** - args: any -> returns bool

### BOR

- **`BOR`** - args: any, any

### BOX2D

- **`BOX2D.BODYCREATE`** - args: float, float, int
- **`BOX2D.FIXTUREBOX`** - args: float, float, float, float
- **`BOX2D.FIXTURECIRCLE`** - args: float
- **`BOX2D.WORLDCREATE`** - args: float, float
- **`BOX2D.WORLDSTEP`** - args: float, int, int

### BRSHIFT

- **`BRSHIFT`** - args: any, int

### BSET

- **`BSET`** - args: any, int

### BSPHERE

- **`BSPHERE.CHECK`** - args: handle, handle -> returns bool
- **`BSPHERE.CHECKBOX`** - args: handle, handle -> returns bool
- **`BSPHERE.FREE`** - args: handle
- **`BSPHERE.MAKE`** - args: float, float, float, float -> returns handle

### BTEST

- **`BTEST`** - args: any, int

### BTOGGLE

- **`BTOGGLE`** - args: any, int

### BXOR

- **`BXOR`** - args: any, any

### CAMERA

- **`CAMERA.BEGIN`** - args: handle
- **`CAMERA.END`** - args: (none)
- **`CAMERA.GETMATRIX`** - args: handle -> returns handle
- **`CAMERA.GETRAY`** - args: handle, float, float
- **`CAMERA.GETVIEWRAY`** - args: float, float, handle, int, int
- **`CAMERA.MAKE`** - args: (none)
- **`CAMERA.MOVE`** - args: handle, float, float, float
- **`CAMERA.SETFOV`** - args: handle, float
- **`CAMERA.SETPOS`** - args: handle, float, float, float
- **`CAMERA.SETPOSITION`** - args: handle, float, float, float
- **`CAMERA.SETTARGET`** - args: handle, float, float, float

### CAMERA2D

- **`CAMERA2D.BEGIN`** - args: (none)
- **`CAMERA2D.END`** - args: (none)

### CEIL

- **`CEIL`** - args: any

### CHARCONTROLLER

- **`CHARCONTROLLER.FREE`** - args: handle
- **`CHARCONTROLLER.GETPOS`** - args: handle -> returns handle
- **`CHARCONTROLLER.ISGROUNDED`** - args: handle -> returns bool
- **`CHARCONTROLLER.MAKE`** - args: float, float, float, float, float -> returns handle
- **`CHARCONTROLLER.MOVE`** - args: handle, float, float, float
- **`CHARCONTROLLER.SETPOS`** - args: handle, float, float, float
- **`CHARCONTROLLER.SETPOSITION`** - args: handle, float, float, float
- **`CHARCONTROLLER.X`** - args: handle -> returns float
- **`CHARCONTROLLER.Y`** - args: handle -> returns float
- **`CHARCONTROLLER.Z`** - args: handle -> returns float

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

### CHR$

- **`CHR$`** - args: int

### CLAMP

- **`CLAMP`** - args: any, any, any

### CLIPBOARD

- **`CLIPBOARD.GETIMAGE`** - args: (none) -> returns handle

### CLOSEFILE

- **`CLOSEFILE`** - args: handle

### CLS

- **`CLS`** - args: (none)

### COLOR

- **`COLOR.A`** - args: handle -> returns int
- **`COLOR.B`** - args: handle -> returns int
- **`COLOR.BRIGHTNESS`** - args: handle, float -> returns handle
- **`COLOR.CONTRAST`** - args: handle, float -> returns handle
- **`COLOR.FADE`** - args: handle, float -> returns handle
- **`COLOR.FREE`** - args: handle
- **`COLOR.G`** - args: handle -> returns int
- **`COLOR.HEX`** - args: string -> returns handle
- **`COLOR.HSV`** - args: float, float, float -> returns handle
- **`COLOR.INVERT`** - args: handle -> returns handle
- **`COLOR.LERP`** - args: handle, handle, float -> returns handle
- **`COLOR.R`** - args: handle -> returns int
- **`COLOR.RGB`** - args: int, int, int -> returns handle
- **`COLOR.RGBA`** - args: int, int, int, int -> returns handle
- **`COLOR.TOHEX`** - args: handle -> returns string
- **`COLOR.TOHSVX`** - args: handle -> returns float
- **`COLOR.TOHSVY`** - args: handle -> returns float
- **`COLOR.TOHSVZ`** - args: handle -> returns float

### COMMAND$

- **`COMMAND$`** - args: (none)
- **`COMMAND$`** - args: int

### CONTAINS

- **`CONTAINS`** - args: string, string

### COPYFILE

- **`COPYFILE`** - args: string, string

### COS

- **`COS`** - args: any

### COUNT$

- **`COUNT$`** - args: string, string -> returns int

### CURSOR

- **`CURSOR.DISABLE`** - args: (none)
- **`CURSOR.ENABLE`** - args: (none)
- **`CURSOR.HIDE`** - args: (none)
- **`CURSOR.ISHIDDEN`** - args: (none)
- **`CURSOR.ISONSCREEN`** - args: (none)
- **`CURSOR.SET`** - args: int
- **`CURSOR.SHOW`** - args: (none)

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

### DATE$

- **`DATE$`** - args: (none)

### DATETIME$

- **`DATETIME$`** - args: (none)

### DAY

- **`DAY`** - args: (none)

### DEBUG

- **`DEBUG.ASSERT`** - args: any, string
- **`DEBUG.BREAKPOINT`** - args: (none)
- **`DEBUG.GCSTATS`** - args: (none)
- **`DEBUG.HEAPSTATS`** - args: (none)
- **`DEBUG.LOG`** - args: string
- **`DEBUG.LOGFILE`** - args: string, string
- **`DEBUG.PRINT`** - args: any
- **`DEBUG.PRINTL`** - args: string, any
- **`DEBUG.PROFILEEND`** - args: string
- **`DEBUG.PROFILEREPORT`** - args: (none)
- **`DEBUG.PROFILESTART`** - args: string
- **`DEBUG.STACKTRACE`** - args: (none)
- **`DEBUG.WATCH`** - args: string, any
- **`DEBUG.WATCHCLEAR`** - args: (none)

### DEG2RAD

- **`DEG2RAD`** - args: any

### DELETEDIR

- **`DELETEDIR`** - args: string

### DELETEFILE

- **`DELETEFILE`** - args: string

### DIREXISTS

- **`DIREXISTS`** - args: string

### DRAW

- **`DRAW.RECTANGLE`** - args: int, int, int, int, int, int, int, int
- **`DRAW.RECTANGLE_ROUNDED`** - args: int, int, int, int, int, int, int, int, int
- **`DRAW.TEXTURE`** - args: handle, int, int, int, int, int, int
- **`DRAW.TEXTURENPATCH`** - args: handle, int, int, int, int, int, int, int, int, int, int, int, int

### DUMP

- **`DUMP`** - args: any

### E

- **`E`** - args: (none)

### ENDSWITH

- **`ENDSWITH`** - args: string, string -> returns bool

### ENET

- **`ENET.CREATEHOST`** - args: string, int, int, int, int
- **`ENET.DEINITIALIZE`** - args: (none)
- **`ENET.HOSTBROADCAST`** - args: handle, int, int, handle
- **`ENET.HOSTSERVICE`** - args: handle, int
- **`ENET.INITIALIZE`** - args: (none)
- **`ENET.PEERPING`** - args: handle
- **`ENET.PEERSEND`** - args: handle, int, handle

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

- **`EVENT.DATA`** - args: handle -> returns string
- **`EVENT.FREE`** - args: handle
- **`EVENT.ISPLAYING`** - args: (none) -> returns bool
- **`EVENT.LISTCLEAR`** - args: handle
- **`EVENT.LISTCOUNT`** - args: handle -> returns int
- **`EVENT.LISTEXPORT`** - args: handle, string
- **`EVENT.LISTFREE`** - args: handle
- **`EVENT.LISTLOAD`** - args: string -> returns handle
- **`EVENT.LISTMAKE`** - args: string -> returns handle
- **`EVENT.PEER`** - args: handle -> returns handle
- **`EVENT.RECPLAYING`** - args: (none) -> returns bool
- **`EVENT.RECSTART`** - args: (none)
- **`EVENT.RECSTOP`** - args: (none)
- **`EVENT.REPLAY`** - args: handle
- **`EVENT.SETACTIVELIST`** - args: handle
- **`EVENT.TYPE`** - args: handle -> returns int

### EXP

- **`EXP`** - args: any

### FILE

- **`FILE.CLOSE`** - args: handle
- **`FILE.EOF`** - args: handle
- **`FILE.OPEN`** - args: string, string
- **`FILE.READLINE`** - args: handle
- **`FILE.SEEK`** - args: handle, int
- **`FILE.SIZE`** - args: handle -> returns int
- **`FILE.TELL`** - args: handle -> returns int
- **`FILE.WRITE`** - args: handle, string
- **`FILE.WRITELN`** - args: handle, string

### FILEEXISTS

- **`FILEEXISTS`** - args: string

### FILEPOS

- **`FILEPOS`** - args: handle

### FILESIZE

- **`FILESIZE`** - args: handle

### FIX

- **`FIX`** - args: any

### FLOAT

- **`FLOAT`** - args: any

### FLOOR

- **`FLOOR`** - args: any

### FONT

- **`FONT.DRAWDEFAULT`** - args: (none)
- **`FONT.FREE`** - args: handle
- **`FONT.LOAD`** - args: string
- **`FONT.LOADBDF`** - args: string, int

### FORMAT$

- **`FORMAT$`** - args: any, string -> returns string

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

### GETDIR$

- **`GETDIR$`** - args: (none)

### GETDIRS$

- **`GETDIRS$`** - args: string

### GETDROPPEDFILES

- **`GETDROPPEDFILES`** - args: (none)

### GETFILEEXT$

- **`GETFILEEXT$`** - args: string

### GETFILEMODTIME

- **`GETFILEMODTIME`** - args: string

### GETFILENAME$

- **`GETFILENAME$`** - args: string

### GETFILENAMENOEXT$

- **`GETFILENAMENOEXT$`** - args: string

### GETFILEPATH$

- **`GETFILEPATH$`** - args: string

### GETFILES$

- **`GETFILES$`** - args: string

### GETFILESIZE

- **`GETFILESIZE`** - args: string

### HEX$

- **`HEX$`** - args: int

### HOUR

- **`HOUR`** - args: (none)

### IIF

- **`IIF`** - args: any, any, any

### IMAGE

- **`IMAGE.CLEARBACKGROUND`** - args: handle, int, int, int, int
- **`IMAGE.COLORBRIGHTNESS`** - args: handle, int
- **`IMAGE.COLORCONTRAST`** - args: handle, float
- **`IMAGE.COLORGRAYSCALE`** - args: handle
- **`IMAGE.COLORINVERT`** - args: handle
- **`IMAGE.COLORREPLACE`** - args: handle, int, int, int, int, int, int, int, int
- **`IMAGE.COLORTINT`** - args: handle, int, int, int, int
- **`IMAGE.CROP`** - args: handle, int, int, int, int
- **`IMAGE.DRAWCIRCLE`** - args: handle, int, int, int, int, int, int, int
- **`IMAGE.DRAWLINE`** - args: handle, int, int, int, int, int, int, int, int
- **`IMAGE.DRAWPIXEL`** - args: handle, int, int, int, int, int, int
- **`IMAGE.DRAWRECT`** - args: handle, int, int, int, int, int, int, int, int
- **`IMAGE.DRAWTEXT`** - args: handle, int, int, string, int, int, int, int, int
- **`IMAGE.EXPORT`** - args: handle, string
- **`IMAGE.FLIPH`** - args: handle
- **`IMAGE.FLIPV`** - args: handle
- **`IMAGE.FREE`** - args: handle
- **`IMAGE.GETBBOXH`** - args: handle, float
- **`IMAGE.GETBBOXW`** - args: handle, float
- **`IMAGE.GETBBOXX`** - args: handle, float
- **`IMAGE.GETBBOXY`** - args: handle, float
- **`IMAGE.GETCOLORA`** - args: handle, int, int
- **`IMAGE.GETCOLORB`** - args: handle, int, int
- **`IMAGE.GETCOLORG`** - args: handle, int, int
- **`IMAGE.GETCOLORR`** - args: handle, int, int
- **`IMAGE.HEIGHT`** - args: handle
- **`IMAGE.LOAD`** - args: string
- **`IMAGE.LOADRAW`** - args: string, int, int, int, int
- **`IMAGE.MAKEBLANK`** - args: int, int, int, int, int, int
- **`IMAGE.MAKECOPY`** - args: handle
- **`IMAGE.MAKETEXT`** - args: string, int, int, int, int, int
- **`IMAGE.RESIZE`** - args: handle, int, int
- **`IMAGE.RESIZENN`** - args: handle, int, int
- **`IMAGE.ROTATE`** - args: handle, int
- **`IMAGE.ROTATECCW`** - args: handle
- **`IMAGE.ROTATECW`** - args: handle
- **`IMAGE.WIDTH`** - args: handle

### INPUT

- **`INPUT`** - args: string -> returns string
- **`INPUT.GAMEPADAXISCOUNT`** - args: int -> returns int
- **`INPUT.GAMEPADBUTTONCOUNT`** - args: int -> returns int
- **`INPUT.GETKEYNAME`** - args: int -> returns string
- **`INPUT.GETMOUSEWORLDPOS`** - args: handle, int, int -> returns handle
- **`INPUT.GETTOUCHPOINTID`** - args: int -> returns int
- **`INPUT.KEYDOWN`** - args: any
- **`INPUT.KEYPRESSED`** - args: any
- **`INPUT.KEYUP`** - args: any
- **`INPUT.MOUSEDOWN`** - args: int
- **`INPUT.MOUSEX`** - args: (none)
- **`INPUT.MOUSEY`** - args: (none)
- **`INPUT.SETGAMEPADMAPPINGS`** - args: string -> returns int
- **`INPUT.SETMOUSEOFFSET`** - args: int, int
- **`INPUT.SETMOUSESCALE`** - args: float, float
- **`INPUT.TOUCHCOUNT`** - args: (none) -> returns int
- **`INPUT.TOUCHPRESSED`** - args: int -> returns bool
- **`INPUT.TOUCHX`** - args: int -> returns int
- **`INPUT.TOUCHY`** - args: int -> returns int

### INSTR

- **`INSTR`** - args: string, string
- **`INSTR`** - args: string, string, int

### INT

- **`INT`** - args: any

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

### JOIN$

- **`JOIN$`** - args: handle, string -> returns string

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

- **`JSON.FREE`** - args: handle
- **`JSON.GETBOOL`** - args: handle, string -> returns bool
- **`JSON.GETFLOAT`** - args: handle, string -> returns float
- **`JSON.GETINT`** - args: handle, string -> returns int
- **`JSON.GETSTRING`** - args: handle, string -> returns string
- **`JSON.MAKE`** - args: (none) -> returns handle
- **`JSON.PARSE`** - args: string -> returns handle
- **`JSON.SETBOOL`** - args: handle, string, bool
- **`JSON.SETFLOAT`** - args: handle, string, float
- **`JSON.SETINT`** - args: handle, string, int
- **`JSON.SETSTRING`** - args: handle, string, string
- **`JSON.TOSTRING`** - args: handle -> returns string

### LEFT$

- **`LEFT$`** - args: string, int

### LEN

- **`LEN`** - args: string

### LERP

- **`LERP`** - args: any, any, any

### LIGHT

- **`LIGHT.MAKE`** - args: (none) -> returns handle
- **`LIGHT.MAKE`** - args: string -> returns handle

### LOCATE

- **`LOCATE`** - args: int, int

### LOG

- **`LOG`** - args: any

### LOG10

- **`LOG10`** - args: any

### LOG2

- **`LOG2`** - args: any

### LOWER$

- **`LOWER$`** - args: string

### LSET$

- **`LSET$`** - args: string, int -> returns string

### LTRIM$

- **`LTRIM$`** - args: string -> returns string

### MAKEDIR

- **`MAKEDIR`** - args: string

### MAKEDIRS

- **`MAKEDIRS`** - args: string

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

- **`MATERIAL.FREE`** - args: handle
- **`MATERIAL.MAKEDEFAULT`** - args: (none)
- **`MATERIAL.SETCOLOR`** - args: handle, int, int, int, int, int
- **`MATERIAL.SETFLOAT`** - args: handle, int, float
- **`MATERIAL.SETSHADER`** - args: handle, handle
- **`MATERIAL.SETTEXTURE`** - args: handle, int, handle

### MATH

- **`MATH.ABS`** - args: any
- **`MATH.ACOS`** - args: any
- **`MATH.ANGLEDIFF`** - args: any, any
- **`MATH.ASIN`** - args: any
- **`MATH.ATAN`** - args: any
- **`MATH.ATAN2`** - args: any, any
- **`MATH.ATN`** - args: any
- **`MATH.CEIL`** - args: any
- **`MATH.CLAMP`** - args: any, any, any
- **`MATH.COS`** - args: any
- **`MATH.DEG2RAD`** - args: any
- **`MATH.E`** - args: (none)
- **`MATH.EXP`** - args: any
- **`MATH.FIX`** - args: any
- **`MATH.FLOOR`** - args: any
- **`MATH.LERP`** - args: any, any, any
- **`MATH.LOG`** - args: any
- **`MATH.LOG10`** - args: any
- **`MATH.LOG2`** - args: any
- **`MATH.MAX`** - args: any, any
- **`MATH.MIN`** - args: any, any
- **`MATH.PI`** - args: (none)
- **`MATH.PINGPONG`** - args: any, any
- **`MATH.POW`** - args: any, any
- **`MATH.RAD2DEG`** - args: any
- **`MATH.RND`** - args: (none)
- **`MATH.RND`** - args: any
- **`MATH.RNDF`** - args: any, any
- **`MATH.RNDSEED`** - args: any
- **`MATH.ROUND`** - args: any
- **`MATH.ROUND`** - args: any, any
- **`MATH.SGN`** - args: any
- **`MATH.SIN`** - args: any
- **`MATH.SMOOTHSTEP`** - args: any, any, any
- **`MATH.SQR`** - args: any
- **`MATH.SQRT`** - args: any
- **`MATH.TAN`** - args: any
- **`MATH.TAU`** - args: (none)
- **`MATH.WRAP`** - args: any, any, any
- **`MATH.WRAPANGLE`** - args: any
- **`MATH.WRAPANGLE180`** - args: any

### MATRIX

- **`MATRIX.FREE`** - args: handle

### MAX

- **`MAX`** - args: any, any

### MEM

- **`MEM.CLEAR`** - args: handle
- **`MEM.COPY`** - args: handle, handle, int, int, int
- **`MEM.FREE`** - args: handle
- **`MEM.GETBYTE`** - args: handle, int -> returns int
- **`MEM.GETDWORD`** - args: handle, int -> returns int
- **`MEM.GETFLOAT`** - args: handle, int -> returns float
- **`MEM.GETSTRING`** - args: handle, int -> returns string
- **`MEM.GETWORD`** - args: handle, int -> returns int
- **`MEM.MAKE`** - args: int -> returns handle
- **`MEM.SETBYTE`** - args: handle, int, int
- **`MEM.SETDWORD`** - args: handle, int, int
- **`MEM.SETFLOAT`** - args: handle, int, float
- **`MEM.SETSTRING`** - args: handle, int, string
- **`MEM.SETWORD`** - args: handle, int, int
- **`MEM.SIZE`** - args: handle -> returns int

### MESH

- **`MESH.CUBE`** - args: float, float, float
- **`MESH.DRAW`** - args: handle, handle, handle
- **`MESH.DRAWROTATED`** - args: handle, handle, float, float, float
- **`MESH.FREE`** - args: handle
- **`MESH.GENTANGENTS`** - args: handle
- **`MESH.GETBBOXMAXX`** - args: handle
- **`MESH.GETBBOXMAXY`** - args: handle
- **`MESH.GETBBOXMAXZ`** - args: handle
- **`MESH.GETBBOXMINX`** - args: handle
- **`MESH.GETBBOXMINY`** - args: handle
- **`MESH.GETBBOXMINZ`** - args: handle
- **`MESH.MAKECONE`** - args: float, float, int
- **`MESH.MAKECUBE`** - args: float, float, float
- **`MESH.MAKECUBICMAP`** - args: handle, float, float, float
- **`MESH.MAKECYLINDER`** - args: float, float, int
- **`MESH.MAKEHEIGHTMAP`** - args: handle, float, float, float
- **`MESH.MAKEKNOT`** - args: float, float, int, int
- **`MESH.MAKEPLANE`** - args: float, float, int, int
- **`MESH.MAKEPOLY`** - args: int, float
- **`MESH.MAKESPHERE`** - args: float, int, int
- **`MESH.MAKETORUS`** - args: float, float, int, int
- **`MESH.PLANE`** - args: float, float, int, int
- **`MESH.SPHERE`** - args: float, int, int
- **`MESH.UPDATEVERTEX`** - args: handle, int, float, float, float, float, float, float, float, float
- **`MESH.UPLOAD`** - args: handle, bool

### MID$

- **`MID$`** - args: string, int
- **`MID$`** - args: string, int, int

### MILLISECOND

- **`MILLISECOND`** - args: (none)

### MIN

- **`MIN`** - args: any, any

### MINUTE

- **`MINUTE`** - args: (none)

### MKDOUBLE$

- **`MKDOUBLE$`** - args: any

### MKFLOAT$

- **`MKFLOAT$`** - args: any

### MKINT$

- **`MKINT$`** - args: any

### MKLONG$

- **`MKLONG$`** - args: any

### MKSHORT$

- **`MKSHORT$`** - args: any

### MODEL

- **`MODEL.ATTACHTO`** - args: handle, handle
- **`MODEL.CLONE`** - args: handle
- **`MODEL.DETACH`** - args: handle
- **`MODEL.EXISTS`** - args: handle
- **`MODEL.FREE`** - args: handle
- **`MODEL.GETMATERIALCOUNT`** - args: handle
- **`MODEL.INSTANCE`** - args: handle
- **`MODEL.LOAD`** - args: string
- **`MODEL.ROTATETEXTURE`** - args: handle, float
- **`MODEL.SCALETEXTURE`** - args: handle, float, float
- **`MODEL.SCROLLTEXTURE`** - args: handle, float, float
- **`MODEL.SETALPHA`** - args: handle, int
- **`MODEL.SETAMBIENTCOLOR`** - args: handle, int, int, int
- **`MODEL.SETBLEND`** - args: handle, int
- **`MODEL.SETCULL`** - args: handle, bool
- **`MODEL.SETDEPTH`** - args: handle, int
- **`MODEL.SETDIFFUSE`** - args: handle, int, int, int
- **`MODEL.SETEMISSIVE`** - args: handle, int, int, int
- **`MODEL.SETFOG`** - args: handle, bool
- **`MODEL.SETGPUSKINNING`** - args: handle, bool
- **`MODEL.SETLIGHTING`** - args: handle, bool
- **`MODEL.SETMATERIAL`** - args: handle, int, handle
- **`MODEL.SETMATERIALSHADER`** - args: handle, int, handle
- **`MODEL.SETMATERIALTEXTURE`** - args: handle, int, int, handle
- **`MODEL.SETMODELMESHMATERIAL`** - args: handle, int, int
- **`MODEL.SETPOS`** - args: handle, float, float, float
- **`MODEL.SETPOSITION`** - args: handle, float, float, float
- **`MODEL.SETSPECULAR`** - args: handle, int, int, int
- **`MODEL.SETSPECULARPOW`** - args: handle, float
- **`MODEL.SETSTAGEBLEND`** - args: handle, int, float
- **`MODEL.SETSTAGEROTATE`** - args: handle, int, float
- **`MODEL.SETSTAGESCALE`** - args: handle, int, float, float
- **`MODEL.SETSTAGESCROLL`** - args: handle, int, float, float
- **`MODEL.SETTEXTURESTAGE`** - args: handle, int, handle
- **`MODEL.SETWIREFRAME`** - args: handle, bool

### MONTH

- **`MONTH`** - args: (none)

### MOVEFILE

- **`MOVEFILE`** - args: string, string

### NET

- **`NET.BROADCAST`** - args: handle, int, string, bool
- **`NET.CLOSE`** - args: handle
- **`NET.CONNECT`** - args: handle, string, int -> returns handle
- **`NET.CREATECLIENT`** - args: (none) -> returns handle
- **`NET.CREATESERVER`** - args: int, int -> returns handle
- **`NET.GETPING`** - args: handle -> returns int
- **`NET.PEERCOUNT`** - args: handle -> returns int
- **`NET.RECEIVE`** - args: handle -> returns handle
- **`NET.SETBANDWIDTH`** - args: handle, int, int
- **`NET.SETTIMEOUT`** - args: handle, int
- **`NET.START`** - args: (none)
- **`NET.STOP`** - args: (none)
- **`NET.UPDATE`** - args: handle

### OCT$

- **`OCT$`** - args: int -> returns string

### OPENFILE

- **`OPENFILE`** - args: string, string

### PEER

- **`PEER.DISCONNECT`** - args: handle
- **`PEER.IP`** - args: handle -> returns string
- **`PEER.PING`** - args: handle -> returns int
- **`PEER.SEND`** - args: handle, int, string, bool

### PHYSICS2D

- **`PHYSICS2D.ONCOLLISION`** - args: handle, handle, string
- **`PHYSICS2D.PROCESSCOLLISIONS`** - args: (none)
- **`PHYSICS2D.SETGRAVITY`** - args: float, float
- **`PHYSICS2D.START`** - args: (none)
- **`PHYSICS2D.STEP`** - args: (none)
- **`PHYSICS2D.STOP`** - args: (none)

### PHYSICS3D

- **`PHYSICS3D.ONCOLLISION`** - args: handle, handle, string
- **`PHYSICS3D.ONCOLLISION`** - args: handle, handle, string
- **`PHYSICS3D.PROCESSCOLLISIONS`** - args: (none)
- **`PHYSICS3D.RAYCAST`** - args: float, float, float, float, float, float, float -> returns handle
- **`PHYSICS3D.SETGRAVITY`** - args: float, float, float
- **`PHYSICS3D.SETGRAVITY`** - args: float, float, float
- **`PHYSICS3D.SETSUBSTEPS`** - args: int
- **`PHYSICS3D.START`** - args: (none)
- **`PHYSICS3D.START`** - args: (none)
- **`PHYSICS3D.STEP`** - args: (none)
- **`PHYSICS3D.STEP`** - args: (none)
- **`PHYSICS3D.STOP`** - args: (none)
- **`PHYSICS3D.STOP`** - args: (none)

### PI

- **`PI`** - args: (none)

### PINGPONG

- **`PINGPONG`** - args: any, any

### POW

- **`POW`** - args: any, any

### PRINT

- **`PRINT`** - args: any

### PRINTAT

- **`PRINTAT`** - args: int, int, any

### PRINTCOLOR

- **`PRINTCOLOR`** - args: int, int, int, any

### PRINTLN

- **`PRINTLN`** - args: any

### QUIT

- **`QUIT`** - args: (none)

### RAD2DEG

- **`RAD2DEG`** - args: any

### RAND

- **`RAND.FREE`** - args: handle
- **`RAND.MAKE`** - args: int -> returns handle
- **`RAND.NEXT`** - args: handle, int, int -> returns int
- **`RAND.NEXTF`** - args: handle -> returns float

### RANDOMIZE

- **`RANDOMIZE`** - args: (none)
- **`RANDOMIZE`** - args: any

### RAY

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
- **`RAY.MAKE`** - args: float, float, float, float, float, float -> returns handle

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

### READALLTEXT$

- **`READALLTEXT$`** - args: string

### READBYTE

- **`READBYTE`** - args: handle

### READFILE$

- **`READFILE$`** - args: handle

### READFLOAT

- **`READFLOAT`** - args: handle

### READINT

- **`READINT`** - args: handle

### READSHORT

- **`READSHORT`** - args: handle

### READSTRING$

- **`READSTRING$`** - args: handle, int

### RENAMEFILE

- **`RENAMEFILE`** - args: string, string

### RENDER

- **`RENDER.BEGINFRAME`** - args: (none)
- **`RENDER.BEGINMODE2D`** - args: (none)
- **`RENDER.BEGINMODE3D`** - args: (none)
- **`RENDER.BEGINSHADER`** - args: handle
- **`RENDER.CLEAR`** - args: (none)
- **`RENDER.CLEAR`** - args: handle
- **`RENDER.CLEAR`** - args: int, int, int
- **`RENDER.CLEAR`** - args: int, int, int, int
- **`RENDER.CLEARSCISSOR`** - args: (none)
- **`RENDER.DRAWFPS`** - args: int, int
- **`RENDER.ENDFRAME`** - args: (none)
- **`RENDER.ENDMODE2D`** - args: (none)
- **`RENDER.ENDMODE3D`** - args: (none)
- **`RENDER.ENDSHADER`** - args: (none)
- **`RENDER.FRAME`** - args: (none)
- **`RENDER.SCREENSHOT`** - args: string
- **`RENDER.SETAMBIENT`** - args: float, float, float
- **`RENDER.SETBLEND`** - args: int
- **`RENDER.SETBLENDMODE`** - args: int
- **`RENDER.SETCULLFACE`** - args: int
- **`RENDER.SETDEPTHMASK`** - args: bool
- **`RENDER.SETDEPTHTEST`** - args: bool
- **`RENDER.SETDEPTHWRITE`** - args: bool
- **`RENDER.SETFPS`** - args: int
- **`RENDER.SETIBLINTENSITY`** - args: float
- **`RENDER.SETIBLSPLIT`** - args: float, float
- **`RENDER.SETMSAA`** - args: bool
- **`RENDER.SETSCISSOR`** - args: int, int, int, int
- **`RENDER.SETSHADOWMAPSIZE`** - args: int
- **`RENDER.SETSKYBOX`** - args: string
- **`RENDER.SETWIREFRAME`** - args: bool

### REPEAT$

- **`REPEAT$`** - args: string, int -> returns string

### REPLACE$

- **`REPLACE$`** - args: string, string, string

### REVERSE$

- **`REVERSE$`** - args: string -> returns string

### RIGHT$

- **`RIGHT$`** - args: string, int

### RND

- **`RND`** - args: (none)
- **`RND`** - args: any

### RNDF

- **`RNDF`** - args: any, any

### RNDSEED

- **`RNDSEED`** - args: any

### ROUND

- **`ROUND`** - args: any
- **`ROUND`** - args: any, any

### RSET$

- **`RSET$`** - args: string, int -> returns string

### RTRIM$

- **`RTRIM$`** - args: string -> returns string

### SECOND

- **`SECOND`** - args: (none)

### SEEKFILE

- **`SEEKFILE`** - args: handle, string, int

### SETDIR

- **`SETDIR`** - args: string

### SGN

- **`SGN`** - args: any

### SHADER

- **`SHADER.LOAD`** - args: string, string

### SIN

- **`SIN`** - args: any

### SLEEP

- **`SLEEP`** - args: any

### SMOOTHSTEP

- **`SMOOTHSTEP`** - args: any, any, any

### SOUND

- **`SOUND.FREE`** - args: handle
- **`SOUND.FROMWAVE`** - args: handle -> returns handle

### SPACE$

- **`SPACE$`** - args: int -> returns string

### SPC

- **`SPC`** - args: int

### SPLIT$

- **`SPLIT$`** - args: string, string -> returns handle

### SPRITE

- **`SPRITE.DEFANIM`** - args: handle, string
- **`SPRITE.DRAW`** - args: handle, int, int
- **`SPRITE.HIT`** - args: handle, handle
- **`SPRITE.LOAD`** - args: string
- **`SPRITE.PLAYANIM`** - args: handle, string
- **`SPRITE.SETPOS`** - args: handle, float, float
- **`SPRITE.SETPOSITION`** - args: handle, float, float
- **`SPRITE.UPDATEANIM`** - args: handle, float

### SQR

- **`SQR`** - args: any

### SQRT

- **`SQRT`** - args: any

### STARTSWITH

- **`STARTSWITH`** - args: string, string -> returns bool

### STOP

- **`STOP`** - args: (none)

### STR$

- **`STR$`** - args: any

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

### TAB

- **`TAB`** - args: int

### TAN

- **`TAN`** - args: any

### TAU

- **`TAU`** - args: (none)

### TERRAIN

- **`TERRAIN.MAKE`** - args: int, int

### TEXTURE

- **`TEXTURE.FREE`** - args: handle
- **`TEXTURE.FROMIMAGE`** - args: handle
- **`TEXTURE.GENWHITENOISE`** - args: int, int
- **`TEXTURE.LOAD`** - args: string

### THROW

- **`THROW`** - args: int, string

### TICKCOUNT

- **`TICKCOUNT`** - args: (none)

### TIME

- **`TIME.DELTA`** - args: (none)
- **`TIME.GET`** - args: (none)

### TIME$

- **`TIME$`** - args: (none)

### TIMER

- **`TIMER`** - args: (none)

### TIMESTAMP

- **`TIMESTAMP`** - args: (none)

### TRACE

- **`TRACE`** - args: any

### TRIM$

- **`TRIM$`** - args: string

### TYPEOF

- **`TYPEOF`** - args: any

### UPPER$

- **`UPPER$`** - args: string

### UTIL

- **`UTIL.CHANGEDIR`** - args: string -> returns bool
- **`UTIL.CLEARDROPPEDFILES`** - args: (none)
- **`UTIL.FILEEXISTS`** - args: string -> returns bool
- **`UTIL.GETDIRFILES`** - args: string -> returns string
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
- **`UTIL.MAKEDIRECTORY`** - args: string -> returns bool
- **`UTIL.SAVETEXT`** - args: string, string

### VAL

- **`VAL`** - args: string -> returns float

### VEC2

- **`VEC2.ADD`** - args: handle, handle -> returns handle
- **`VEC2.ANGLE`** - args: handle, handle -> returns float
- **`VEC2.DISTANCE`** - args: handle, handle -> returns float
- **`VEC2.FREE`** - args: handle
- **`VEC2.LENGTH`** - args: handle -> returns float
- **`VEC2.LERP`** - args: handle, handle, float -> returns handle
- **`VEC2.MAKE`** - args: float, float -> returns handle
- **`VEC2.MUL`** - args: handle, float -> returns handle
- **`VEC2.NORMALIZE`** - args: handle -> returns handle
- **`VEC2.ROTATE`** - args: handle, float -> returns handle
- **`VEC2.SET`** - args: handle, float, float
- **`VEC2.SUB`** - args: handle, handle -> returns handle
- **`VEC2.X`** - args: handle -> returns float
- **`VEC2.Y`** - args: handle -> returns float

### VEC3

- **`VEC3.ADD`** - args: handle, handle -> returns handle
- **`VEC3.CROSS`** - args: handle, handle -> returns handle
- **`VEC3.DISTANCE`** - args: handle, handle -> returns float
- **`VEC3.DIV`** - args: handle, float -> returns handle
- **`VEC3.DOT`** - args: handle, handle -> returns float
- **`VEC3.EQUALS`** - args: handle, handle -> returns bool
- **`VEC3.FREE`** - args: handle
- **`VEC3.LENGTH`** - args: handle -> returns float
- **`VEC3.LERP`** - args: handle, handle, float -> returns handle
- **`VEC3.MAKE`** - args: float, float, float -> returns handle
- **`VEC3.MUL`** - args: handle, float -> returns handle
- **`VEC3.NEGATE`** - args: handle -> returns handle
- **`VEC3.NORMALIZE`** - args: handle -> returns handle
- **`VEC3.REFLECT`** - args: handle, handle -> returns handle
- **`VEC3.SET`** - args: handle, float, float, float
- **`VEC3.SUB`** - args: handle, handle -> returns handle
- **`VEC3.X`** - args: handle -> returns float
- **`VEC3.Y`** - args: handle -> returns float
- **`VEC3.Z`** - args: handle -> returns float

### WAIT

- **`WAIT`** - args: any

### WAVE

- **`WAVE.COPY`** - args: handle -> returns handle
- **`WAVE.CROP`** - args: handle, int, int
- **`WAVE.EXPORT`** - args: handle, string
- **`WAVE.FORMAT`** - args: handle, int, int, int
- **`WAVE.FREE`** - args: handle
- **`WAVE.LOAD`** - args: string -> returns handle

### WINDOW

- **`WINDOW.CHECKFLAG`** - args: int -> returns bool
- **`WINDOW.CLEARFLAG`** - args: int
- **`WINDOW.CLOSE`** - args: (none)
- **`WINDOW.GETMONITORCOUNT`** - args: (none) -> returns int
- **`WINDOW.GETMONITORHEIGHT`** - args: int -> returns int
- **`WINDOW.GETMONITORNAME`** - args: int -> returns string
- **`WINDOW.GETMONITORREFRESHRATE`** - args: int -> returns int
- **`WINDOW.GETMONITORWIDTH`** - args: int -> returns int
- **`WINDOW.GETPOSITIONX`** - args: (none) -> returns int
- **`WINDOW.GETPOSITIONY`** - args: (none) -> returns int
- **`WINDOW.GETSCALEDPIX`** - args: (none) -> returns float
- **`WINDOW.GETSCALEDPIY`** - args: (none) -> returns float
- **`WINDOW.MAXIMIZE`** - args: (none)
- **`WINDOW.MINIMIZE`** - args: (none)
- **`WINDOW.OPEN`** - args: int, int, string -> returns bool
- **`WINDOW.RESTORE`** - args: (none)
- **`WINDOW.SETFLAG`** - args: int
- **`WINDOW.SETFPS`** - args: int
- **`WINDOW.SETICON`** - args: string
- **`WINDOW.SETMAXSIZE`** - args: int, int
- **`WINDOW.SETMINSIZE`** - args: int, int
- **`WINDOW.SETMONITOR`** - args: int
- **`WINDOW.SETOPACITY`** - args: float
- **`WINDOW.SETPOSITION`** - args: int, int
- **`WINDOW.SETSIZE`** - args: int, int
- **`WINDOW.SETSTATE`** - args: int
- **`WINDOW.SETTARGETFPS`** - args: int
- **`WINDOW.SETTITLE`** - args: string
- **`WINDOW.SHOULDCLOSE`** - args: (none)

### WRAP

- **`WRAP`** - args: any, any, any

### WRAPANGLE

- **`WRAPANGLE`** - args: any

### WRAPANGLE180

- **`WRAPANGLE180`** - args: any

### WRITE

- **`WRITE`** - args: any

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

### YEAR

- **`YEAR`** - args: (none)

