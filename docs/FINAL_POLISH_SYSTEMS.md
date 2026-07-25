# Final polish systems checklist

This document maps the **40 beginner systems** from [`audit/archive/final-polish-and-docs.md`](audit/archive/final-polish-and-docs.md) to moonBASIC’s shipped commands, notes on memory, and gaps.

**Beginner docs:** Full system guides (style guide format) live in **[systems/README.md](systems/README.md)** — start with [01-CORE](systems/01-CORE.md), then follow the build order in that index.

## Quick reference

| System | Status | Canonical commands |
|--------|--------|-------------------|
| APP | Aliases | `APP.*` → `WINDOW.*` / `TIME.*` / `SYSTEM.VERSION` |
| RENDER | Core + aliases | `RENDER.CLEAR`, `RENDER.FRAME`, `RENDER.BEGIN`/`END` (→ `BEGIN3D`/`END3D`) |
| SCENE | Shipped | `SCENE.REGISTER`, `SCENE.SWITCH`, `SCENE.DRAW`, `SCENE.SAVESCENE`, `SCENE.LOADSCENE` |
| ENTITY | Shipped | `ENTITY.CREATE`, `CREATECUBE`, `SETPOSITION`, `TURN`, hierarchy, tags |
| CAMERA | Shipped | `CAMERA.CREATE`, `SETACTIVE`, `LOOKAT`, `BEGIN`/`END` |
| LIGHT | Shipped | `LIGHT.CREATEPOINT`, `CREATEDIRECTIONAL`, `CREATESPOT` |
| MESH | Shipped | `MESH.CUBE`, `SPHERE`, `PLANE`, `CREATECYLINDER`, `UPLOAD` |
| MODEL | Shipped | `MODEL.LOAD`, `MODEL.FREE`, `MODEL.ANIMCOUNT` |
| MATERIAL | Shipped | `MATERIAL.CREATE`, `SETCOLOR`, `SETTEXTURE` |
| TEXTURE | Shipped | `TEXTURE.LOAD`, `WIDTH`, `HEIGHT`, `FREE` |
| ANIM | Dual API | Entity: `ENTITY.PLAY` / `STOPANIM`; FSM: `ANIM.DEFINE`, `UPDATE` |
| INPUT | Shipped | `INPUT.KEYDOWN`, `KEYHIT`, `MOUSEX`, `GAMEPAD*` aliases |
| ACTION | Shipped | `ACTION.MAPKEY`, `DOWN`, `PRESSED` (`HIT` alias) |
| PHYSICS | Shipped | `PHYSICS.START`, `SETGRAVITY`, `STEP` |
| BODY | Aliases | `BODY.*` → `ENTITY.ADDPHYSICS` + property setters |
| COLLISION | Geometry API | Overlap helpers; beginner slide API → use `ENTITY` + physics |
| PICK | Shipped | `PICK.SCREENCAST`, `CAST`, `HIT`, `X/Y/Z`, `DIST` |
| AUDIO | Shipped | `AUDIO.LOADSOUND`, `PLAY`, `STOP`, `LOADMUSIC`, `UPDATEMUSIC` |
| AUDIO3D | Aliases | `AUDIO3D.*` → spatial `SOUND.PLAY3D`, `AUDIO.LISTENERCAMERA` |
| UI | Shipped | `GUI.*` for menus; `UI.BUTTON`, `UI.LABEL3D` for widgets |
| FONT/TEXT | Shipped | `FONT.LOAD`, `DRAW.TEXT`, `TEXT.*` aliases |
| SPRITE | Shipped | `SPRITE.LOAD`, `SETPOSITION`, `DRAW` |
| TILEMAP | Shipped | `TILEMAP.LOAD`, `DRAW`, `GETTILE`, `SETTILE` |
| TERRAIN | Shipped | `TERRAIN.CREATE`, `GETHEIGHT`, `LOAD`, `DRAW` |
| PARTICLE | Shipped | `PARTICLE.CREATE`, `SETRATE`, `PLAY`, `STOP` |
| TIMER | Shipped | Handle timers + `TIMER.AFTER`/`EVERY`/`CANCEL` callbacks |
| SAVE | Shipped | `SAVE.DATA`/`GET`, `SAVE.WRITE`/`READ` (file paths) |
| ASSET | Shipped | `ASSET.LOADPACK`, `TEXTURE`, `MODEL`, `SOUND`, `UNLOAD` |
| FILE | Shipped | `FILE.EXISTS`, `READTEXT`, `WRITETEXT`, `DELETE` |
| JSON | Shipped | `JSON.PARSE`, typed getters, `TOSTRING`/`STRINGIFY` |
| MATH | Shipped | `MATH.RAND`, `CLAMP`, `LERP`, `VEC3.DISTANCE` |
| VEC3 | Shipped | `VEC3.CREATE`, `ADD`, `NORMALIZE`, `LENGTH` |
| DEBUG | Shipped | `DEBUG.LOG`, `DRAWLINE`, `DRAWBOX`, `SHOWFPSGRAPH` |
| ERROR | Compiler | File/line diagnostics, “did you mean” suggestions |
| PROJECT | CLI | `moonbasic new`, `run`, `build` |
| PACKAGE | CLI | `moonbasic pack`, `moonbasic package windows|linux` |
| MODULE | Language | `IMPORT "file.mb"` |
| HELP | Shipped | `HELP("ENTITY.SETPOSITION")` |
| TEST | CLI | `moonbasic test` |
| TEMPLATE | CLI | `moonbasic new --template 3d|platformer|ui` |

## Foundation example

Run (full runtime):

```bash
moonrun examples/foundation/main.mb
```

Check only (no window):

```bash
moonbasic --check examples/foundation/main.mb
```

## Memory notes

- **ASSET.LOADPACK** caches loaded handles; call **ASSET.UNLOAD** or reload pack to free GPU/audio memory.
- **TEXTURE.FREE** / **MODEL.FREE** / **FREESOUND** free individual assets.
- **ENTITY.FREE** / **TIMER.CANCEL** / **TIMER.FREE** release handles.
- **TIMER.AFTER/EVERY** callbacks are removed when fired (AFTER) or via **TIMER.CANCEL**.

## Canonical vs checklist names

Prefer checklist names in new tutorials (`APP.OPEN`, `RENDER.BEGIN`). They are aliases; LSP lists both where registered.

## Case insensitivity

moonBASIC treats **commands**, **keywords**, **variables**, **function names**, **KEY_*** constants, **SAVE** keys, and **ASSET** pack ids as **case-insensitive** in source. Examples that compile identically:

- `APP.OPEN` / `app.open` / `App.Open`
- `WHILE` / `while` / `While`
- `KEY_ESCAPE` / `key_escape`
- `FUNCTION Tick` … `tick()` / `TICK()`

File paths (`FILE.*`, `IMPORT`, asset paths on disk) still follow the OS filesystem rules.

See [LANGUAGE.md](LANGUAGE.md) and `testdata/case_insensitive.mb`.

See also: [COMMANDS.md](COMMANDS.md), [GETTING_STARTED.md](GETTING_STARTED.md).
