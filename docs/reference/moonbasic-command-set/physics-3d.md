# Physics 3D (Jolt)

| Designed | moonBASIC | Memory / notes |
|----------|------------|----------------|
| **Physics3D.Start()** | **`Physics3D.Start()`** | Initializes the 3D physics world. **Linux + CGO + jolt-go** only. |
| **Physics3D.Stop()** | **`Physics3D.Stop()`** | Shuts down the simulation and frees resources. |
| **Physics3D.Step()** | **`Physics3D.Step()`** | Advances simulation (call once per frame). |
| **Physics3D.SetGravity(x, y, z)** | **`Physics3D.SetGravity()`** | Sets the global gravity vector. |
| **Body3D.Create(type)** | **`Body3D.Make()`** | Creates a body definition (`"static"`, `"dynamic"`, `"kinematic"`). |
| **Body3D.AddBox(def, w, h, d)** | **`Body3D.AddBox()`** | Adds a box collision shape to the definition. |
| **Body3D.Commit(def, x, y, z)** | **`Body3D.Commit()`** | Finalizes the body and adds it to the world. Returns a **body handle**. |
| **Body3D.SetPos(id, x, y, z)** | **`Body3D.SetPos()`** | Teleports a body to a new position. |
| **Body3D.SetLinearVel(id, vx, vy, vz)** | **`Body3D.SetLinearVel()`** | Sets linear velocity directly. |
| **Body3D.GetMatrix(id)** | **`Body3D.GetMatrix()`** | Returns transform matrix handle for visual sync. |
| **Body3D.Free(id)** | **`Body3D.Free()`** | Removes a body and frees its memory. |
| **Body3D.Collided(a, b)** | **`Body3D.Collided()`** | Returns TRUE if two bodies are in contact. |

**Aliases:** **`PHYSICS.*`** mirrors several **`PHYSICS3D.*`** names.

See also: [PHYSICS3D.md](../PHYSICS3D.md).
