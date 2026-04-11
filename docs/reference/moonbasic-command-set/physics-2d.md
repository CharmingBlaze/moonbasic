# Physics 2D (Box2D)

| Designed | moonBASIC | Memory / notes |
|----------|------------|----------------|
| **Physics2D.Start([gx, gy])** | **`Physics2D.Start()`** | Initializes the 2D physics world. Default gravity is `(0, 500)`. |
| **Physics2D.Stop()** | **`Physics2D.Stop()`** | Shuts down the simulation and frees internal buffers. |
| **Physics2D.Step()** | **`Physics2D.Step()`** | Advances the simulation (call once per frame). |
| **Body2D.Create(type)** | **`Body2D.Make()`** | Creates a body definition (`"static"`, `"dynamic"`, `"kinematic"`). |
| **Body2D.AddRect(def, w, h)** | **`Body2D.AddRect()`** | Adds a rectangle collision shape to the definition. |
| **Body2D.AddCircle(def, r)** | **`Body2D.AddCircle()`** | Adds a circle collision shape to the definition. |
| **Body2D.Commit(def, x, y)** | **`Body2D.Commit()`** | Finalizes the body and adds it to the world. Returns a **body handle**. |
| **Body2D.SetPos(id, x, y)** | **`Body2D.SetPos()`** | Teleports a 2D body to a new position. |
| **Body2D.X(id)** / **Body2D.Y(id)** | **`Body2D.X()`** / **`Body2D.Y()`** | Returns the current X or Y coordinate of the body's center. |
| **Body2D.Rot(id)** | **`Body2D.Rot()`** | Returns the body rotation in radians. |
| **Body2D.Free(id)** | **`Body2D.Free()`** | Removes a body from the simulation and frees its memory. |

**Legacy:** **`BOX2D.*`** aliases remain (`WORLDCREATE`, `BODYCREATE`, `FIXTUREBOX`, `FIXTURECIRCLE`).

See also: [PHYSICS2D.md](../PHYSICS2D.md).
