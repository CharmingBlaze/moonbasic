# Input

| Designed | moonBASIC | Memory / notes |
|----------|------------|----------------|
| **KeyDown(key)** | **`Input.KeyDown()`** | Returns TRUE if held. |
| **KeyHit(key)** | **`Input.KeyPressed()`** | TRUE on first frame. |
| **MouseDown(btn)** | **`Input.MouseButtonDown()`** | |
| **MouseHit(btn)** | **`Input.MouseButtonPressed()`** | |
| **MouseX()** | **`Input.MouseX()`** | |
| **MouseY()** | **`Input.MouseY()`** | |
| **MouseZ()** | **`Input.MouseWheelMove()`** | |
| **MoveMouse(x, y)** | **`Input.SetMousePos()`** | |
| **Input.Axis(n, p)** | **`Input.Axis()`** | Returns -1, 0, or 1. |
| **Input.MouseDelta()** | **`Input.MouseDelta()`** | Returns [dx, dy] handle. |

See [INPUT.md](../INPUT.md), [INPUT helpers](../INPUT.md) for **`Input.Axis`**.
