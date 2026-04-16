# Texture atlas (`ATLAS.*`, `SPRITE.*`)

Texture **atlases** are large images containing many smaller sub-images or sprites—efficient for 2D asset management.

**Conventions:** [STYLE_GUIDE.md](../../STYLE_GUIDE.md), [API_CONVENTIONS.md](API_CONVENTIONS.md) — reference pages use uppercase **`NAMESPACE.ACTION`**; Easy Mode (`Atlas.Load`, …) is [compatibility only](../../STYLE_GUIDE.md#easy-mode-compatibility-layer).

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — see [WAVE.md](WAVE.md) (registry-first headings, **Full Example** at the end).

## Core Workflow

1. **Pack:** Use a tool like TexturePacker to pack sprites into one sheet and export a `.json` data file (e.g. “JSON (Array)” format).
2. **Load:** **`ATLAS.LOAD(imagePath, jsonPath)`** — returns an atlas handle.
3. **Sprites:** **`ATLAS.GETSPRITE(atlasHandle, spriteName)`** — handles for each packed filename.
4. **Draw:** **`SPRITE.DRAW(spriteHandle, x, y)`** (and other **`SPRITE.*`** commands as needed).
5. **Free:** **`ATLAS.FREE(atlasHandle)`** when done.

---

### `ATLAS.LOAD(imagePath, jsonPath)`
Loads a texture atlas from an image and a JSON data file. Returns a handle to the atlas.

- `imagePath`: Path to the atlas texture sheet (e.g. `.png`).
- `jsonPath`: Path to the JSON data file.

---

### `ATLAS.GETSPRITE(atlasHandle, spriteName)`
Retrieves a handle to a single sprite within the atlas.

- `atlasHandle`: Handle of the loaded atlas.
- `spriteName`: Original filename of the sprite as packed (e.g. `"player.png"`).

---

### `ATLAS.FREE(atlasHandle)`
Frees the atlas texture and all associated sprite data.

---

## See also

- [SPRITE.md](SPRITE.md) — **`SPRITE.DRAW`**, strip frames, **`ANIM.*`**
- [IMAGE.md](IMAGE.md) — CPU images before GPU upload

---

## Full Example

Assume `game_atlas.png` and `game_atlas.json`, with sprites named `player.png` and `enemy.png`.

```basic
WINDOW.OPEN(800, 600, "Texture Atlas Example")
WINDOW.SETFPS(60)

my_atlas = ATLAS.LOAD("game_atlas.png", "game_atlas.json")

player_sprite = ATLAS.GETSPRITE(my_atlas, "player.png")
enemy_sprite = ATLAS.GETSPRITE(my_atlas, "enemy.png")

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(20, 20, 20)
    CAMERA2D.BEGIN()
        SPRITE.DRAW(player_sprite, 100, 250)
        SPRITE.DRAW(enemy_sprite, 400, 250)
    CAMERA2D.END()
    RENDER.FRAME()
WEND

ATLAS.FREE(my_atlas)
WINDOW.CLOSE()
```
