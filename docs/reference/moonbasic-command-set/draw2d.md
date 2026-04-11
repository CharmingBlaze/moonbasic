# 2D drawing

No **global pen** — colour is **per call** or stored in your own variables.

| Designed | moonBASIC | Memory / notes |
|----------|------------|----------------|
| **SetColor(r, g, b)** | **`Draw.SetColor()`** | Active 2D draw color. |
| **Line(x1, y1, x2, y2)** | **`Draw.Line()`** | |
| **Rect(x, y, w, h)** | **`Draw.Rectangle()`** | |
| **Oval(x, y, w, h)** | **`Draw.Ellipse()`** | |
| **Text(x, y, text)** | **`Draw.Text()`** | |
| **Cls()** | **`Render.Clear()`** | |
| **SetOrigin(x, y)** | **`Camera2D.SetOffset()`** | 2D camera state. |
| **SetViewport(x, y, w, h)** | **`Render.SetScissor()`** | Clipping rect. |

See [DRAW2D.md](../DRAW2D.md), [DRAW_WRAPPERS.md](../DRAW_WRAPPERS.md).
