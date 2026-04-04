# Camera Commands

Commands for creating and controlling 2D and 3D cameras.

---

## 3D Camera

### `Camera.Make()`

Creates a new 3D camera and returns a handle to it.

```basic
cam = Camera.Make()
```

---

### `Camera.SetPos(cameraHandle, x#, y#, z#)`

Sets the position of the camera in 3D space.

- `cameraHandle`: The handle of the camera.
- `x#`, `y#`, `z#`: The world coordinates.

```basic
cam.SetPos(0, 10, 20)
```

---

### `Camera.SetTarget(cameraHandle, x#, y#, z#)`

Sets the point in 3D space that the camera will look at.

- `cameraHandle`: The handle of the camera.
- `x#`, `y#`, `z#`: The world coordinates of the target.

```basic
cam.SetTarget(0, 0, 0)
```

---

### `Camera.SetFOV(cameraHandle, fov#)`

Sets the vertical field of view of the camera.

- `cameraHandle`: The handle of the camera.
- `fov#`: The field of view in degrees.

---

### `Camera.Begin(cameraHandle)`

Enters 3D rendering mode using the specified camera's settings.

- `cameraHandle`: The handle of the camera to use.

---

### `Camera.End()`

Ends 3D camera mode.

---

### `Camera.GetMatrix(cameraHandle)`

Returns a handle to the camera's view matrix. Useful for custom shaders or advanced rendering.

---

## 2D Camera

### `Camera2D.Make()`

Creates a new 2D camera for scrolling, zooming, and rotating the 2D view. Returns a handle.

---

### `Camera2D.Begin(cameraHandle)`

Enters 2D rendering mode with the specified 2D camera's transformations.

---

### `Camera2D.End()`

Exits 2D camera mode.

---

### `Camera2D.SetTarget(cameraHandle, x#, y#)`

Sets the camera's target position (the point it will center on).

---

### `Camera2D.SetOffset(cameraHandle, x#, y#)`

Sets the camera's screen offset. Useful for centering the camera on the screen.

---

### `Camera2D.SetZoom(cameraHandle, zoom#)`

Sets the camera's zoom level. `1.0` is normal zoom.

---

### `Camera2D.SetRotation(cameraHandle, angle#)`

Sets the camera's rotation in degrees.
