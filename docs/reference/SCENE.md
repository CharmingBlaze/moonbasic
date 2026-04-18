# Scene Commands

Named scene registration, loading, and per-frame update/draw hooks with optional transitions.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Register scenes with `SCENE.REGISTER`, mapping an ID to a loader function.
2. Set per-frame hooks with `SCENE.SETHANDLERS`.
3. Load a scene with `SCENE.LOAD` or `SCENE.LOADASYNC`.
4. Each frame, call `SCENE.UPDATE(dt)` and `SCENE.DRAW()`.
5. Free with `SCENE.FREE` when done.

For file-based level loading see [LEVEL.md](LEVEL.md). For transitions see [TRANSITION.md](TRANSITION.md).

---

### `SCENE.REGISTER(id, loaderName)`
Maps a scene ID to a user function for initialization.

- **Arguments**:
    - `id`: (String) Unique scene identifier.
    - `loaderName`: (String) Name of the BASIC function to call.
- **Returns**: (None)

---

### `SCENE.LOAD(id)` / `LOADASYNC`
Loads a scene and runs its registration hook.

- **Returns**: (Handle) The scene handle.

---

### `SCENE.SETHANDLERS(updateName, drawName)`
Sets global hooks for the per-frame loop.

- **Returns**: (None)

---

### `SCENE.UPDATE(dt)` / `DRAW()`
Advances the scene state or renders the active scene.

---

### `SCENE.CURRENT()`
Returns the ID string of the active scene.

- **Returns**: (String)

---

## Full Example

```basic
SCENE.SETHANDLERS("MYUPDATE", "MYDRAW")
SCENE.REGISTER("LEVEL1", "LOAD_LEVEL1")
SCENE.LOAD("LEVEL1")

WHILE NOT WINDOW.SHOULDCLOSE()
    SCENE.UPDATE(TIME.DELTA())
    ; ... game rendering ...
    SCENE.DRAW()
WEND
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `SCENE.LOADSCENE(path)` | Load a scene file and replace the current scene. |
| `SCENE.SAVESCENE(path)` | Serialize and save the current scene to file. |
| `SCENE.SWITCH(name)` | Switch to a named scene (registered with `SCENE.REGISTER`). |
| `SCENE.LOADWITHTRANSITION(name, transition)` | Switch scene with a named transition effect. |
| `SCENE.CLEARSCENE()` | Remove all entities and reset the scene. |
| `SCENE.APPLYPHYSICS(bool)` | Enable/disable automatic physics step during `SCENE.UPDATE`. |

## See also

- [ENTITY.md](ENTITY.md) — entity creation and management
- [WORLD.md](WORLD.md) — `WORLD.SETCENTER`, `WORLD.UPDATE`
