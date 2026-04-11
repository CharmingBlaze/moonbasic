# Scene management (`SCENE.*`)

Registers **named scenes** backed by user **`FUNCTION`** loaders, optional **per-frame update/draw** hooks, and optional **transitions** when switching scenes. Implemented in `runtime/mbscene`.

> **Note:** This module is **not** a glTF level loader. For the roadmap toward Blender → scoped scene graph, deduplicated assets, and Jolt buffer integration, see [SCENE_ENGINE_BRIEF.md](SCENE_ENGINE_BRIEF.md).

The module must receive a **user-function invoker** from the host (same mechanism as tweens and behavior trees): without it, loads fail.

---

## Registration

### `Scene.Register(sceneId, loadFunctionName)`

Maps a string id (stored uppercase) to the **name of a parameterless user function** that runs when the scene is loaded. That function should create state, load assets, and wire subsystems.

### `Scene.SetHandlers(updateFunctionName, drawFunctionName)`

Sets global names for:

- **Update** — called as `update(dt)` with delta time in seconds.
- **Draw** — called with **no arguments** once per frame after transition polling.

Names are folded to uppercase.

---

## Loading

### `Scene.Load(sceneId)` / `Scene.LoadAsync(sceneId)`

- **Load** — runs the loader **immediately** and clears any pending async load.
- **LoadAsync** — queues the id; the load runs at the **start** of the next `SCENE.UPDATE` call.

### `Scene.LoadWithTransition(sceneId, kind, duration)`

Starts a transition **out**, then loads the scene when the transition completes, then optionally fades **in**:

- `kind` — `"fade"` or `"wipe"` (case-insensitive).
- `duration` — seconds, must be positive.

Uses [TRANSITION](TRANSITION.md) (`TRANSITION.FADEOUT`, `TRANSITION.WIPE`, `TRANSITION.FADEIN`, `TRANSITION.ISDONE`). Requires an active runtime registry so `TRANSITION.*` can be invoked.

---

## Frame hooks

### `Scene.Update(dt)`

1. Runs a pending **async** load if any.
2. Advances **transition** state (loads pending scene after fade-out when transition finishes).
3. Invokes the registered **update** function with `dt`.

### `Scene.Draw()`

Polls transitions again, then invokes the **draw** function.

### `Scene.Current()` → string

Returns the current scene id (uppercase), or empty if none.

---

## Typical loop

```basic
Scene.SetHandlers "MYUPDATE", "MYDRAW"
Scene.Register "LEVEL1", "LOAD_LEVEL1"

Scene.Load "LEVEL1"

WHILE NOT Window.ShouldClose()
    Scene.Update(Time.Delta())
    ; ... game rendering ...
    Scene.Draw()
WEND
```
