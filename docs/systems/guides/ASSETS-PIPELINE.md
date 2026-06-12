# Assets pipeline — manifests, paths, and memory

> Load textures, models, and sounds by **id** from one JSON pack — ship a folder players can run.

**Namespaces:** `ASSET` · `TEXTURE` · `MODEL` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#assets](../COMMAND_REGISTRY.md#assets) · [03-ASSETS.md](../03-ASSETS.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [Why use ASSET.LOADPACK](#why-use-assetloadpack)
- [Core workflow](#core-workflow)
- [Manifest format](#manifest-format)
- [Using handles in gameplay](#using-handles-in-gameplay)
- [Direct load vs pack](#direct-load-vs-pack)
- [Memory and reload](#memory-and-reload)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Approach | When |
|----------|------|
| **`ASSET.LOADPACK` + ids** | Real projects, multiple file types |
| **`TEXTURE.LOAD(path)`** | Single quick test |
| **`ENTITY.LOAD(path)`** | One hero model without manifest |

**Working directory:** Paths in manifests are relative to the **manifest file** location. Run `moonrun` from the project root (or set paths accordingly).

---

## Why use ASSET.LOADPACK

- **One place** to list every file — artists add rows, code stays `ASSET.TEXTURE("crate")`.
- **Case-insensitive ids** — `Crate` and `crate` match.
- **Cached handles** — second `ASSET.TEXTURE("crate")` returns same GPU memory.
- **Ship with game** — zip `assets/` + `main.mb` + `moonrun`.

---

## Core workflow

1. Create `assets/assets.json` listing files.
2. At startup: `ASSET.LOADPACK("assets/assets.json")`.
3. Fetch: `tex = ASSET.TEXTURE("player")`, `ASSET.MODEL("hero")`, `ASSET.SOUND("jump")`.
4. Assign to entities/materials/audio (see other guides).
5. On full reload: `ASSET.UNLOAD()` then load new pack.

```basic
ASSET.LOADPACK("assets/assets.json")
heroModel = ASSET.MODEL("hero")
hero = ENTITY.LOAD("")   ; or SETMODEL — see ENTITY guide
```

---

## Manifest format

```json
{
  "textures": {
    "player": "textures/player.png",
    "crate": "textures/crate.png"
  },
  "models": {
    "hero": "models/hero.glb"
  },
  "sounds": {
    "jump": "audio/jump.wav",
    "land": "audio/land.wav"
  }
}
```

Paths are relative to `assets/` folder containing the JSON.

---

## Using handles in gameplay

| Asset | Typical use |
|-------|-------------|
| Texture | `MATERIAL.SETTEXTURE(mat, tex)` |
| Model | `ENTITY.SETMODEL(ent, model)` or `ENTITY.LOAD` |
| Sound | `AUDIO.PLAYSOUND(ASSET.SOUND("jump"))` |

---

## Direct load vs pack

| `TEXTURE.LOAD("path.png")` | `ASSET.TEXTURE("id")` |
|----------------------------|------------------------|
| Fast prototype | Production layout |
| Path in every file | Path once in JSON |
| Manual `TEXTURE.FREE` | `ASSET.UNLOAD` frees pack |

---

## Memory and reload

- **`TEXTURE.FREE` / `MODEL.FREE`** — single asset.
- **`ASSET.UNLOAD`** — entire pack cache (level transition).
- **`ERASE ALL` / shutdown** — see [MEMORY.md](../../MEMORY.md).

**Why unload:** GPU memory on texture/model reload (e.g. season swap).

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Wrong working directory | Run from folder containing `assets/` |
| Id not in JSON | Error at load — check manifest |
| Reload pack without unload | Orphaned GPU handles |
| Mix release versions | Ship same moonBASIC tag as dev |

---

## See also

- [ENTITY-SYSTEM.md](ENTITY-SYSTEM.md)
- [AUDIO-FEEDBACK.md](AUDIO-FEEDBACK.md)
- `moonbasic new` creates `assets/` folder
