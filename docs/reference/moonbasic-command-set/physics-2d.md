# Physics 2D (Box2D)

| Designed | moonBASIC | Memory / notes |
|----------|------------|----------------|
| **Physics2D.Start([gx, gy])** | **`PHYSICS2D.START`** | Initializes the 2D physics world. Default gravity is `(0, 500)`. |
| **Physics2D.Stop()** | **`PHYSICS2D.STOP`** | Shuts down the simulation and frees internal buffers. |
| **Physics2D.Step()** | **`PHYSICS2D.STEP`** | Advances the simulation (call once per frame). |
| **Body2D.Create(type)** | **`BODY2D.CREATE`** (deprecated **`BODY2D.MAKE`**) | Creates a body definition (`"static"`, `"dynamic"`, `"kinematic"`). |
| **Body2D.AddRect(def, w, h)** | **`BODY2D.ADDRECT`** | Adds a rectangle collision shape to the definition. |
| **Body2D.AddCircle(def, r)** | **`BODY2D.ADDCIRCLE`** | Adds a circle collision shape to the definition. |
| **Body2D.Commit(def, x, y)** | **`BODY2D.COMMIT`** | Finalizes the body and adds it to the world. Returns a **body handle**. |
| **Body2D.SetPos(id, x, y)** | **`BODY2D.SETPOS`** (deprecated **`BODY2D.SETPOSITION`**) | Teleports a 2D body to a new position. |
| **Body2D.X(id)** / **Body2D.Y(id)** | **`BODY2D.X`** / **`BODY2D.Y`** | Returns the current X or Y coordinate of the body's center. |
| **Body2D.Rot(id)** | **`BODY2D.ROT`** | Returns the body rotation in radians. |
| **Body2D.Free(id)** | **`BODY2D.FREE`** | Removes a body from the simulation and frees its memory. |

**Legacy:** **`BOX2D.*`** aliases remain (`WORLDCREATE`, `BODYCREATE`, `FIXTUREBOX`, `FIXTURECIRCLE`).

See also: [PHYSICS2D.md](../PHYSICS2D.md).
