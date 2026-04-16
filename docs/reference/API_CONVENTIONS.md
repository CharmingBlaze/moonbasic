# API conventions — consistent object commands

**Directive:** [API Standardization Directive](../API_STANDARDIZATION_DIRECTIVE.md) (CREATE vs MAKE, SETPOS, universal handle methods, phased rollout).

moonBASIC registers commands as **`NAMESPACE.ACTION`** (uppercase, dot). In source you can write **`Object.Method`** style; names are matched case-insensitively and resolve to those registry keys (see [PROGRAMMING.md](../PROGRAMMING.md)). This page describes **recommended verbs** so similar concepts look alike across types (`Load`, `SetPos`, `Scale`, `Rotate`, …).

---

## 1. Standard verbs

| Intent | Typical registry name | Notes |
|--------|-------------------------|--------|
| Create from a **file** | `*.LOAD` | `MODEL.LOAD`, `SPRITE.LOAD`, `TEXTURE.LOAD`, `FONT.LOAD`, … |
| Create **procedurally** | `*.CREATE` | `CAMERA.CREATE`, `LIGHT.CREATE`, `BODY3D.CREATE`, `MODEL.CREATE` (from mesh), … |
| Release a handle | `*.FREE` | Pair with the matching `LOAD` / `CREATE` / `CREATE*` |
| Set **position** | `*.SETPOS` | Canonical name; **`SETPOSITION`** is registered as an **alias** where listed in the manifest |
| Set **uniform scale** | `*.SETSCALE` | Use when the type exposes scaling (e.g. instances); not every handle has this yet |
| Set **rotation** | `*.SETROT`, `*.LOOKAT`, or `TRANSFORM.*` | Cameras often use **`LOOKAT`** / **`SETTARGET`**; models may use matrix / texture-stage rotates — see per-type docs |
| Draw / step | `*.DRAW`, `*.UPDATE`, … | Namespace-specific |
| Advance **3D physics** (Jolt) | `PHYSICS3D.UPDATE`, `PHYSICS3D.STEP` | Same implementation; optional frame **`dt`** — see [PHYSICS3D.md](PHYSICS3D.md) |

**Aliases:** The manifest may register migration aliases (for example `*.MAKE*` and `*.SETPOSITION`). Prefer canonical **`CREATE`** and **`SETPOS`** in all new code.

---

## 2. What exists today (quick map)

| Family | Load / make | Position | Scale / rotate (summary) |
|--------|-------------|----------|---------------------------|
| **3D model** | `MODEL.LOAD`, `MODEL.CREATE` | `MODEL.SETPOS` | Root transform is **translation only** in `SETPOS`; use **`TRANSFORM.*`** / materials for full TRS (see [MODEL.md](MODEL.md)) |
| **2D sprite** | `SPRITE.LOAD` | `SPRITE.SETPOS` | Draw position + internal offset; no separate scale/rotate on the sprite handle today |
| **Camera** | `CAMERA.CREATE` | `CAMERA.SETPOS` | Orientation via **`SETTARGET`**, **`MOVE`**, **`LOOKAT`**, **`SETUP`** — not `SETROT` |
| **3D body** | `BODY3D` builder + `COMMIT` | `BODY3D.SETPOS` | Physics orientation APIs vary; see [PHYSICS3D.md](PHYSICS3D.md) |
| **Light** | `LIGHT.CREATE` | `LIGHT.SETPOS` / `LIGHT.SETDIR` | Directional vs point/spot differ; see [LIGHT.md](LIGHT.md) |
| **Texture** | `TEXTURE.LOAD` | N/A | N/A |

This table is **descriptive** (current engine), not a promise that every row will gain every verb.

---

## 3. Naming style in examples

| Style | Example |
|-------|---------|
| Registry (docs, errors) | `MODEL.SETPOS(mdl, x, y, z)` |
| Mixed-case in scripts | `Model.SetPos(mdl, x, y, z)` or `model.SetPos(...)` on a handle |

---

## 4. Adding new commands

When you add a new handle type to the manifest and runtime:

1. Prefer **`LOAD`** (asset) vs **`CREATE`** (procedural) for canonical naming.
2. Use **`SETPOS`** for position; add **`SETPOSITION`** only if you need backward compatibility with an older name.
3. If you add scale/rotation, prefer **`SETSCALE`** / **`SETROT`** (or document **`LOOKAT`**-style if it is camera-like).
4. Register aliases by pointing at the **same** Go handler (see existing `SETPOS` / `SETPOSITION` pairs).

Full arity and types: **[API_CONSISTENCY.md](../API_CONSISTENCY.md)** (`go run ./tools/apidoc`).
