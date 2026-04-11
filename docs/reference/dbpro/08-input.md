# DBPro — Input

moonBASIC: **`Input.*()`**, shortcuts **`KeyDown()`**, **`MouseX()`**, … — [INPUT.md](../INPUT.md).

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **UPKEY()** / **DOWNKEY()** | ✓ **`Input.KeyDown()`** | Use `KEY_UP`, `KEY_DOWN`. |
| **LEFTKEY()** / **RIGHTKEY()** | ✓ **`Input.KeyDown()`** | |
| **MOUSEX()** / **MOUSEY()** | ✓ **`Input.MouseX()`** / **`Input.MouseY()`** | |
| **MOUSECLICK()** | ✓ **`Input.MouseButtonDown()`** | |
| **MOUSEMOVE()** | ✓ **`Input.MouseDelta()`** | |
| **SCANCODE()** | ✓ **`Input.GetKeyPressed()`** | |
| **KEYSTATE()** | ✓ **`Input.KeyDown()`** | |
| **JOYSTICK X/Y/Z** / **FIRE A** | ✓ **`Input.JoyX()`**, **`Input.JoyY()`**, **`Input.JoyButton()`** / **`Input.JoyDown()`** | See **`Input.Joy*()`** in manifest. |
