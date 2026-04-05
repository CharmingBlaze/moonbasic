# Light Commands

Commands for creating and controlling lights in a 3D scene. The runtime registers **`LIGHT.*`** (see `runtime/mblight`). In source you can write **`Light.Make`** etc.; the compiler emits the same uppercase keys.

---

### `Light.Make` (`LIGHT.MAKE(type$)`)

Creates a new light.

- `type$`: The type of light. Can be `"directional"`, `"point"`, or `"spot"`.

Returns a handle to the light.

```basic
; Create a directional light to act as a sun
sun = Light.Make("directional")
```

---

### `Light.SetDir` (`LIGHT.SETDIR(lightHandle, x#, y#, z#)`)

Sets the direction for a `"directional"` or `"spot"` light.

- `lightHandle`: The handle of the light.
- `x#`, `y#`, `z#`: The direction vector.

---

### `Light.SetShadow` (`LIGHT.SETSHADOW(lightHandle, castShadows?)`)

Enables or disables shadow casting for this light.

- `lightHandle`: The handle of the light.
- `castShadows?`: `TRUE` to enable shadows, `FALSE` to disable.

---

### `Render.SetAmbient(r#, g#, b#)`

Sets the ambient light color for the scene.

- `r#`, `g#`, `b#`: The RGB color components, typically from 0.0 to 1.0.
