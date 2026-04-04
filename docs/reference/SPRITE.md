# Sprite Commands

Commands for working with animated sprites from Aseprite files.

## Aseprite Workflow

moonBASIC has first-class support for `.ase` and `.aseprite` files. The workflow is designed around using **tags** in Aseprite to define your animations.

1.  **Create Animations in Aseprite**: Create your sprite and define animations by creating frame tags (e.g., a tag named "run" that spans frames 1-8).
2.  **Load Sprite**: Use `Sprite.Load()` to load the `.aseprite` file.
3.  **Define Animations**: Use `Sprite.DefAnim()` for each tag you want to use in your game.
4.  **Play Animation**: Call `Sprite.PlayAnim()` to start an animation.
5.  **Update and Draw**: In your main loop, call `Sprite.UpdateAnim()` to advance the animation timer and `Sprite.Draw()` to show it on screen.

---

### `Sprite.Load(filePath$)`

Loads an Aseprite file and its associated data (layers, frames, tags). Returns a handle to the sprite.

- `filePath$`: The path to the Aseprite file.

---

### `Sprite.DefAnim(spriteHandle, animName$)`

Informs moonBASIC about an animation cycle you want to use. This must match the name of a tag in your Aseprite file exactly.

- `spriteHandle`: The handle of the sprite.
- `animName$`: The name of the animation tag (e.g., "walk", "idle").

---

### `Sprite.PlayAnim(spriteHandle, animName$)`

Sets the specified animation as the current one to be played.

- `spriteHandle`: The handle of the sprite.
- `animName$`: The name of the animation to play.

---

### `Sprite.UpdateAnim(spriteHandle, deltaTime#)`

Updates the sprite's animation playback based on the elapsed time. This should be called every frame for smooth animation.

- `spriteHandle`: The handle of the sprite.
- `deltaTime#`: The time elapsed since the last frame, usually from `Time.Delta()`.

---

### `Sprite.Draw(spriteHandle, x, y)`

Draws the current frame of the sprite's active animation at the specified coordinates.

- `spriteHandle`: The handle of the sprite to draw.
- `x`, `y`: The top-left position to draw the sprite.

---

## Animation State Machine

For more complex characters, you can use the `ANIM` commands to build a state machine.

### `ANIM.Define(spriteHandle, animName$, startFrame, endFrame, fps#, looping?)`

Defines a single animation clip for a sprite.

### `ANIM.AddTransition(spriteHandle, fromAnim$, toAnim$, condition$)`

Creates a transition between two animations that triggers when a parameter (set via `ANIM.SETPARAM`) becomes true.

### `ANIM.SETPARAM(spriteHandle, paramName$, value)`

Sets a parameter in the animation state machine, which can trigger transitions.

### `ANIM.UPDATE(spriteHandle, deltaTime#)`

Updates the animation state machine. Use this instead of `Sprite.UpdateAnim` when using the state machine.

---

## Full Example

This example assumes you have an Aseprite file named `player.aseprite` with two tags: "idle" and "run".

```basic
Window.Open(800, 600, "Sprite Animation Example")
Window.SetFPS(60)

; 1. Load sprite
player = Sprite.Load("player.aseprite")

; 2. Define animations from Aseprite tags
Sprite.DefAnim(player, "idle")
Sprite.DefAnim(player, "run")

; 3. Play the initial animation
Sprite.PlayAnim(player, "idle")

player_x = 350
player_y = 250

WHILE NOT Window.ShouldClose()
    ; --- LOGIC ---
    ; Change animation based on input
    IF Input.KeyDown(KEY_RIGHT) OR Input.KeyDown(KEY_LEFT) THEN
        Sprite.PlayAnim(player, "run")
    ELSE
        Sprite.PlayAnim(player, "idle")
    ENDIF

    ; 4. Update the animation every frame
    Sprite.UpdateAnim(player, Time.Delta())

    ; --- DRAWING ---
    Render.Clear(30, 40, 50)
    Render.BeginMode2D()
        ; 5. Draw the sprite
        Sprite.Draw(player, player_x, player_y)
    Render.EndMode2D()
    Render.Frame()
WEND

Window.Close()
```
