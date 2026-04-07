# Object-style draw and input wrappers

This reference describes the **optional** heap-backed wrappers around long-argument **`DRAW3D.*`**, **`DRAW.*`**, **`DRAW.TEXT`**, texture draws, **`INPUT.*`**, and **`LANDBOXES`** / **`PLAYER.MOVERELATIVE`**-style helpers. The original commands remain the canonical low-level API; wrappers call the same Raylib paths (or forward to existing natives).

**CGO:** All drawing wrappers require **Raylib** (same as **`DRAW3D.*`**). **`MOVER`** is pure Go. Input facades (**`MOUSE`**, **`KEY`**, **`GAMEPAD`**) require CGO builds that link Raylib input.

---

## Naming: `CUBE()` vs `DRAWCUBE()`

Identifiers are uppercased at parse time, so **`Cube()`** and **`CUBE()`** are the same.

| Global | Meaning |
|--------|---------|
| **`CUBE()`** / **`CUBE(w,h,d)`** | Blitz-style **entity** constructor → **`ENTITYREF`** (see [BLITZ3D.md](BLITZ3D.md)). |
| **`DRAWCUBE()`** / **`DRAWCUBE(w,h,d)`** | **Immediate-mode** box: state object with **`.Pos`**, **`.Draw`**, etc., mapped to **`DRAW3D.CUBE`** / wires. |

Use **`DRAW*`** prefixes for immediate-mode wrappers so entity and draw APIs do not collide.

---

## 3D immediate primitives

Constructors allocate a handle; methods dispatch to registry keys **`DRAWPRIM3D.*`** (handle is always the first argument at runtime).

| Constructor | Notes |
|-------------|--------|
| **`DRAWCUBE`**, **`DRAWCUBEWIRES`** | Optional **`(w#, h#, d#)`** or defaults **1,1,1**. |
| **`DRAWSPHERE(radius#)`** | Solid sphere. |
| **`DRAWSPHEREW(radius#, rings, slices)`** | Wire sphere. |
| **`DRAWCYLINDER`**, **`DRAWCYLINDERW`** | Defaults; use **`.Cyl`**, **`.Slices`**, **`.Pos`**, **`.Color`**. |
| **`DRAWCAP`**, **`DRAWCAPW`** | Capsule; **`.EndPoint`** for end position, **`.Radius`**, slices/rings. |
| **`DRAWPLANE`** | Optional **`(width#, depth#)`**. |
| **`DRAWBBOX`** | Use **`.BBox`** to set min/max corners. |
| **`DRAWRAY`** | **`.SetRay(rayArrayHandle)`** (6 floats: origin + direction). |
| **`DRAWLINE3D`**, **`DRAWPOINT3D`** | Line: **`.P2`**-style via **`.EndPoint`**; set endpoints with **`.Pos`** / **`.EndPoint`**. |
| **`DRAWGRID3D`** | Optional **`(slices, spacing#)`**. |
| **`DRAWBILLBOARD(tex)`**, **`DRAWBILLBOARDREC(tex)`** | **`.SetTexture`**, **`.SrcTex`** for rec variant; **inside `Camera.Begin`/`End`**. |

### Shared methods (3D)

| Method | Role |
|--------|------|
| **`Pos(x#, y#, z#)`** | Position (and start point for line/capsule). |
| **`Size(w#, h#, d#)`** | Box size; for sphere **`.Size(r)`** (one arg) sets radius where applicable. |
| **`Color(r,g,b,a)`** / **`Col(...)`** | 0–255. |
| **`Wire(on)`** | Toggle wireframe where supported (solid cube types). |
| **`Radius`**, **`EndPoint`**, **`Cyl`**, **`BBox`**, **`Slices`**, **`Rings`**, **`Grid`**, **`SetRay`**, **`SetTexture`**, **`SrcTex`** | Kind-specific. |
| **`Draw()`** | Issues the Raylib draw. |
| **`Free()`** | Releases the wrapper handle (idempotent; no GPU resource for pure immediate shapes). |

Internal registry prefix: **`DRAWPRIM3D`** (e.g. **`DRAWPRIM3D.POS`**). See [vm/handlecall.go](../../vm/handlecall.go) for the full method map.

---

## 2D immediate primitives

Constructors use a **`2`** suffix to avoid clashing with **`DRAW.CIRCLE`** and existing globals.

| Constructor | Notes |
|-------------|--------|
| **`DRAWCIRCLE2(r#)`**, **`DRAWCIRCLE2W(r#)`** | Circle / circle lines. |
| **`DRAWELLIPSE2(rx#, ry#)`**, **`DRAWELLIPSE2W`** | Ellipse. |
| **`DRAWRECT2(w, h)`**, **`DRAWRECT2W`** | Axis-aligned rectangle. |
| **`DRAWLINE2`** | **`.Pos`** and **`.P2`** for endpoints. |
| **`DRAWTRI2`**, **`DRAWTRI2W`** | **`.Pos`**, **`.P2`**, **`.P3`** for vertices. |
| **`DRAWRING2`**, **`DRAWRING2W`** | Ring / ring lines: **`Size(inner#, outer#)`** or **`Ring(inner#, outer#, start#, end#)`**, **`Segs(n)`**, then **`Draw`**. |
| **`DRAWPOLY2(sides#)`**, **`DRAWPOLY2W(sides#)`** | Regular polygon; **`Size(radius#)`**, **`Sides`**, **`Rot`**, and for wire **`Thick`**. |

Registry: **`DRAWPRIM2D.*`**. Methods: **`Pos`**, **`Size`** (two floats, or one for circle radius), **`Color`**, **`Outline`**, **`P2`**, **`P3`**, **`Ring`**, **`Segs`**, **`Sides`**, **`Rot`**, **`Thick`**, **`Draw`**, **`Free`**.

---

## Text

| Constructor / API | Maps to |
|-------------------|---------|
| **`TEXTOBJ(text$)`** | Default font **`DRAW.TEXT`**-style state: **`TEXTDRAW.POS`**, **`SIZE`**, **`COLOR`**, **`SETTEXT`**, **`DRAW`**. |
| **`TEXTOBJEX(font, text$)`** | **`DRAW.TEXTEX`** via **`TEXTEXOBJ.*`**. |
| **`DRAWTEX2(tex)`** | **`DRAW.TEXTURE`**-style blit with position + tint (**`DRAWTEX2.*`**). |
| **`DRAWTEXREC(tex)`** | Sub-rectangle draw (**`DRAW.TEXTUREREC`**): **`Src`**, **`Pos`** (float), **`Color`**, **`Draw`**. |
| **`DRAWTEXPRO(tex)`** | **`DRAW.TEXTUREPRO`**: **`Src`**, **`Dst`**, **`Origin`**, **`Rot`**, **`Color`**, **`Draw`**. |

---

## Input facades (singletons)

| Global | Methods |
|--------|---------|
| **`MOUSE()`** | **`DX`**, **`DY`**, **`WHEEL`**, **`DOWN(button)`**, **`PRESSED`**, **`RELEASED`** |
| **`KEY()`** | **`DOWN(key)`**, **`HIT(key)`**, **`UP(key)`** |
| **`GAMEPAD()`** | **`AXIS(pad#, axis#)`**, **`BUTTON(pad#, button#)`** |

First **`MOUSE()`** / **`KEY()`** / **`GAMEPAD()`** call allocates a handle; later calls return the same handle (per module). Registry: **`MOUSE.DX`**, **`KEY.DOWN`**, etc.

---

## Movement helper

| Global | Methods |
|--------|---------|
| **`MOVER()`** | **`MoveXZ(handle, …)`** → **`[dx, dz]`** array. **`MoveStepX` / `MoveStepZ`** → forward to **`MOVESTEPX` / `MOVESTEPZ`** (same 5 numeric args after the mover handle). **`MoveRel`** → **`PLAYER.MOVERELATIVE`**. **`Land`** → **`LANDBOXES`**. **`Free`**. |

---

## Parity checklist (tooling)

- **`compiler/builtinmanifest/commands.json`** — includes representative wrapper keys; extend with more overload rows as needed.
- **`go run ./tools/cmdaudit`** — namespace coverage audit.
- **`go run ./tools/apidoc`** — refreshes **[API_CONSISTENCY.md](../API_CONSISTENCY.md)**.

---

## See also

- [DRAW3D.md](DRAW3D.md) — long-form **`DRAW3D.*`** and short names **`BOX`**, **`BALL`**, …
- [BLITZ3D.md](BLITZ3D.md) — **`CUBE()`** entities vs **`DRAWCUBE()`** immediate mode
- [DRAW2D.md](DRAW2D.md) — 2D drawing reference
- [INPUT.md](INPUT.md) — **`INPUT.*`** reference
- [GAMEHELPERS.md](GAMEHELPERS.md) — **`LANDBOXES`**, movement helpers
