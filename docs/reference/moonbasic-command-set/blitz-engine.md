# Blitz-style flat commands (engine facade)

**Conventions:** [STYLE_GUIDE.md](../../../STYLE_GUIDE.md) — flat globals map to **`WINDOW.*`**, **`RENDER.*`**, **`ENTITY.*`**, …; prefer registry names in new non-Blitz code.

The **`blitzengine`** runtime module registers **uppercase flat globals** (`APPTITLE`, `GRAPHICS`, `PLOT`, …) that forward to the existing dotted API (`WINDOW.*`, `DRAW.*`, `ENTITY.*`, …). Implementation: [`runtime/blitzengine/`](../../../runtime/blitzengine/).

**Reserved words:** `END` is a keyword — use **`FINISH`** to terminate the program (same as **`ENDGAME`**).

**Breaking change:** **`WRITEFILE`** from this facade maps to **`FILE.WRITEALLTEXT`** `(path, text)` (whole-file write). Streaming write to an open handle remains **`FILE.WRITE`** / legacy patterns.

---

## Quick map (designed name → flat command)

| Area | Flat commands |
|------|----------------|
| Program | `APPTITLE()`, `SETFPS()`, `DELTATIME()`, `TIMEMS()`, `FINISH()` (`SLEEP` stays core `SLEEP`) |
| Display | `GRAPHICS()`, `GRAPHICS3D()`, `SETVSYNC()`, `SETCLEARCOLOR()`, `CLEAR()` — end each frame with **`RENDER.FRAME()`** (Blitz **Flip** is not a separate flat command). |
| 2D | `SETCOLOR()`, `SETALPHA()`, `SETORIGIN()`, `SETVIEWPORT()`, `PLOT()`, `LINE()`, `RECT()`, `OVAL()`, `TEXT()` |
| World | `CREATEWORLD()`, `CLEARWORLD()`, `SETAMBIENT()`, `SETFOG()`, `SETWIREFRAME()` — use **`ENTITY.UPDATE(dt)`** + **`RENDER.BEGIN3D(cam)`** / **`RENDER.END3D()`** (or **`CAMERA.BEGIN`/`CAMERA.END`**) + **`DrawEntities()`** instead of Blitz **UpdateWorld** / **RenderWorld** ([BLITZ3D.md](../BLITZ3D.md)). |
| Camera | `CREATECAMERA()`, `POSITIONCAMERA()`, `ROTATECAMERA()`, `MOVECAMERA()`, `CAMERARANGE()`, `CAMERAZOOM()`, `CAMERAVIEWPORT()`, `CAMERAPROJECT()`, `CAMERAPICK()` |
| Lights | `CREATELIGHT()`, `LIGHTCOLOR()`, `LIGHTRANGE()`, `LIGHTCONE()`, `LIGHTPOSITION()`, `LIGHTPOINTAT()` |
| Entities | `CREATECUBE()`, `CREATESPHERE()`, `CREATEPLANE()`, `CREATEMESH()`, `COPYENTITY()`, `FREEENTITY()`, transforms/getters/state — see [`entities.md`](entities.md) |
| Mesh | `LOADMESH()`, `LOADANIMMESH()`, `MESHWIDTH()`, `MESHHEIGHT()`, `MESHDEPTH()`; surface/vertex API **not** implemented (error) |
| Textures | `CREATETEXTURE()`, `TEXTURECOORDS()` → `Texture.SetWrap()`; `SCALETEXTURE()` / `ROTATETEXTURE()` error (use `Draw.TexturePro()` / filters) |
| Sprites | `CREATESPRITE()`, `SPRITE()` (pos+draw), `MOVESPRITE()`, `SPRITEHIT()`; `SPRITECOLOR()` / `SPRITEALPHA()` no-op |
| Sound | `LOADSOUND()`, `PLAYSOUND()`, `LOOPSOUND()`, `STOPSOUND()`, `SOUNDVOLUME()`, `SOUNDPAN()`, `SOUNDPITCH()` |
| Input | `MOUSEZ()`, `MOUSEDOWN()`, `MOUSEHIT()`, `MOVEMOUSE()`, `FLUSHKEYS()`, `FLUSHMOUSE()` (no-op flush); `KEYDOWN()` / `KEYHIT()` / `MOUSEX()` / `MOUSEY()` remain existing shortcuts |
| Files | `READFILE()`, `WRITEFILE()`, `CLOSEFILE()`, `READLINE()`, `WRITELINE()`, `FILEEXISTS()`, `DELETEFILE()`, `COPYFILE()` |

Details per subsystem stay in the section files linked from [README.md](README.md). The full dotted registry is still [API_CONSISTENCY.md](../../API_CONSISTENCY.md).
