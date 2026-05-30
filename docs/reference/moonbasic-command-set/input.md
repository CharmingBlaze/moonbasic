# Input

**Conventions:** [STYLE_GUIDE.md](../../../STYLE_GUIDE.md), [API_CONVENTIONS.md](../API_CONVENTIONS.md), full reference [INPUT.md](../INPUT.md). Tables use **registry keys** (`INPUT.*`); dotted **`Input.*`** is Easy Mode / compatibility.

| Designed | Registry (use in new examples) | Memory / notes |
|----------|-------------------------------|----------------|
| **KeyDown(key)** | **`INPUT.KEYDOWN(key)`** | TRUE if held. |
| **KeyHit(key)** | **`INPUT.KEYPRESSED(key)`** | TRUE on first frame. |
| **MouseDown(btn)** | **`INPUT.MOUSEDOWN(btn)`** | |
| **MouseHit(btn)** | **`INPUT.MOUSEHIT(btn)`** | If registered for your build. |
| **MouseX()** / **MouseY()** | **`INPUT.MOUSEX()`** / **`INPUT.MOUSEY()`** | |
| **MouseZ()** | **`INPUT.MOUSEWHEELMOVE()`** | Wheel delta. |
| **MoveMouse(x, y)** | **`INPUT.SETMOUSEPOS(x, y)`** | |
| **Axis(neg, pos)** | **`INPUT.AXIS(neg, pos)`** | −1, 0, or 1. |
| **MouseDelta** | **`INPUT.MOUSEDELTA()`** | `[dx, dy]` handle. |

See [INPUT.md](../INPUT.md) for **`INPUT.AXISDEG`**, **`INPUT.MOVEMENT2D`**, **`CURSOR.*`**, and action mapping.
