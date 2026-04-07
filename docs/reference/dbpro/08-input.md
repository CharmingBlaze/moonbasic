# DBPro — Input

moonBASIC: **`INPUT.*`**, shortcuts **`KEYDOWN`**, **`MOUSEX`**, … — [INPUT.md](../INPUT.md).

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **KEY STATE** | ≈ **`INPUT.KEYDOWN`** | |
| **KEY DOWN** | ✓ **`INPUT.KEYDOWN`**, **`KEYDOWN`** (shortcut) | |
| **KEY UP** | ✓ **`INPUT.KEYUP`** | |
| **MOUSE X** / **Y** / **Z** | ✓ **`INPUT.MOUSEX`**, **`MOUSEY`**, wheel | Flat **`MOUSEX`** may exist via **GAME** shortcuts. |
| **MOUSE CLICK** | ≈ **`INPUT.MOUSEHIT`**, **`MOUSEDOWN`** | |
| **MOUSE MOVE** | ✓ **`INPUT.SETMOUSEPOS`** | |
| **JOYSTICK X/Y/Z** / **FIRE A** | ✓ **`INPUT.JOYX`**, **`JOYY`**, **`JOYBUTTON`** / **`JOYDOWN`** | See **`INPUT.JOY*`** in manifest. |
