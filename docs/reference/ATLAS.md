# Texture Atlas Commands

Commands for working with texture atlases, which are large images containing many smaller sub-images or sprites. This is an efficient way to manage 2D assets.

## Core Workflow

1.  **Create an Atlas**: Use a tool like TexturePacker to pack your individual sprite images into a single sheet and export a `.json` data file (in the "JSON (Array)" format).
2.  **Load Atlas**: Use `Atlas.Load()` with both the image sheet and the JSON data file.
3.  **Get Sprites**: Use `Atlas.GetSprite()` to retrieve handles to individual sprites from the atlas by their original filename.
4.  **Use Sprites**: Use the sprite handles with `Sprite.Draw()` and other sprite commands.
5.  **Free Atlas**: Call `Atlas.Free()` when you are done to unload the texture and all associated sprite data.

---

### `Atlas.Load(imagePath$, jsonPath$)`

Loads a texture atlas from an image and a JSON data file. Returns a handle to the atlas.

- `imagePath$`: The path to the atlas texture sheet (e.g., `.png`).
- `jsonPath$`: The path to the JSON data file.

---

### `Atlas.GetSprite(atlasHandle, spriteName$)`

Retrieves a handle to a single sprite within the atlas.

- `atlasHandle`: The handle of the loaded atlas.
- `spriteName$`: The original filename of the sprite as it was packed into the atlas (e.g., `"player.png"`).

---

### `Atlas.Free(atlasHandle)`

Frees the atlas texture and all associated sprite data.

---

## See also

- [SPRITE.md](SPRITE.md) — **`Sprite.Draw`**, strip frames, **`Anim.*`**
- [IMAGE.md](IMAGE.md) — CPU images before GPU upload

---

## Full Example

Assume you have `game_atlas.png` and `game_atlas.json`, and that the atlas contains `player.png` and `enemy.png`.

```basic
Window.Open(800, 600, "Texture Atlas Example")
Window.SetFPS(60)

; 1. Load the atlas
my_atlas = Atlas.Load("game_atlas.png", "game_atlas.json")

; 2. Get individual sprite handles
player_sprite = Atlas.GetSprite(my_atlas, "player.png")
enemy_sprite = Atlas.GetSprite(my_atlas, "enemy.png")

WHILE NOT Window.ShouldClose()
    Render.Clear(20, 20, 20)
    Camera2D.Begin()
        ; 3. Use the sprite handles to draw
        Sprite.Draw(player_sprite, 100, 250)
        Sprite.Draw(enemy_sprite, 400, 250)
    Camera2D.End()
    Render.Frame()
WEND

; 4. Free the entire atlas
Atlas.Free(my_atlas)
Window.Close()
```
