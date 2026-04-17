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

### `SCENE.LOAD(id)` 
Loads a scene from a file and runs its loader function immediately. Returns a **scene handle**.

---

### `SCENE.LOADASYNC(id)` 
Queues a scene to load at the start of the next update cycle.

---

### `SCENE.FREE(sceneHandle)` 
Unloads the scene and frees all associated entities and resources.

---

### `SCENE.REGISTER(id, funcName)` 
Maps a scene string ID to a parameterless user function that acts as the loader.

---

### `SCENE.SETHANDLERS(updateFunc, drawFunc)` 
Sets global names for the per-frame update and draw functions.

---

### `SCENE.UPDATE(dt)` 
Advances transitions, runs async loads, and invokes the registered update function.

---

### `SCENE.DRAW()` 
Polls active transitions and invokes the registered draw function.

---

### `SCENE.CURRENT()` 
Returns the ID string of the currently active scene.

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
