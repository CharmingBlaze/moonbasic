# Scene management (`SCENE.*`)

Registers **named scenes** backed by user **`FUNCTION`** loaders, optional **per-frame update/draw** hooks, and optional **transitions** when switching scenes. Implemented in `runtime/mbscene`.

> **Note:** This module is **not** a glTF level loader. For the roadmap toward Blender → scoped scene graph, deduplicated assets, and Jolt buffer integration, see [SCENE_ENGINE_BRIEF.md](SCENE_ENGINE_BRIEF.md).

The module must receive a **user-function invoker** from the host (same mechanism as tweens and behavior trees): without it, loads fail.

---

## Registration

### `Scene.Load(id)`
Loads a scene from a file and runs its loader function immediately. Returns a **scene handle**.

### `Scene.LoadAsync(id)`
Queues a scene to load at the start of the next update cycle.

### `Scene.Free(handle)`
Unloads the scene and frees all associated entities and native resources.

---

### `Scene.Register(id, funcName)`
Maps a scene string ID to a parameterless user function name that acts as the loader.

### `Scene.SetHandlers(updateFunc, drawFunc)`
Sets global names for the per-frame update and draw functions.

---

### `Scene.Update(dt)`
Advances transitions, runs async loads, and invokes the registered update function.

### `Scene.Draw()`
Polls active transitions and invokes the registered draw function.

### `Scene.Current()`
Returns the ID string of the currently active scene.

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
