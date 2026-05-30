# Game jam example

Runnable demo using **`SPRITE.BUILTIN`** and **`SOUND.BUILTIN`** — no PNG or WAV files needed.

```bash
moonrun examples/gamejam/main.mb
```

Collect coins with arrow keys; space plays the jump sound.

Built-in sprite names: `player`, `enemy`, `bullet`, `tile`, `coin`, `heart`, `star`, `block`.

Built-in sound names: `jump`, `hit`, `coin`, `shoot`, `powerup`, `explode`, `select`, `error`.

Use **`FONT.BUILTIN()`** for the default Raylib font handle when drawing custom text.
