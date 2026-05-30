# Physics 3D (Jolt)

| Designed | moonBASIC | Memory / notes |
|----------|------------|----------------|
| **Physics3D.Start()** | **`PHYSICS3D.START`** | Initializes the 3D physics world. **Linux + CGO + jolt-go** only. |
| **Physics3D.Stop()** | **`PHYSICS3D.STOP`** | Shuts down the simulation and frees resources. |
| **Physics3D.Step()** / **Physics3D.Update()** | **`PHYSICS3D.STEP`** / **`PHYSICS3D.UPDATE`** | Same implementation — advance simulation once per frame (optional **`dt`**). Prefer **`UPDATE`** in new scripts. |
| **Physics3D.SetGravity(x, y, z)** | **`PHYSICS3D.SETGRAVITY`** | Sets the global gravity vector. |
| **Body3D.Create(type)** | **`BODY3D.CREATE`** (deprecated **`BODY3D.MAKE`**) | Creates a body definition (`"static"`, `"dynamic"`, `"kinematic"`). |
| **Body3D.AddBox(def, w, h, d)** | **`BODY3D.ADDBOX`** | Adds a box collision shape to the definition. |
| **Body3D.Commit(def, x, y, z)** | **`BODY3D.COMMIT`** | Finalizes the body and adds it to the world. Returns a **body handle**. |
| **Body3D.SetPos(id, x, y, z)** | **`BODY3D.SETPOS`** (deprecated **`BODY3D.SETPOSITION`**) | Teleports a body to a new position. |
| **Body3D.SetLinearVel(id, vx, vy, vz)** | **`BODY3D.SETLINEARVEL`** | Sets linear velocity directly. |
| **Transform / render sync** | **`BODY3D.BUFFERINDEX`**, **`PHYSICS3D.GETMATRIXBUFFER`** | Shared matrix pool; pair visuals with **`ENTITY.LINKPHYSBUFFER`** (see [PHYSICS3D.md](../PHYSICS3D.md)). |
| **Body3D.Free(id)** | **`BODY3D.FREE`** | Removes a body and frees its memory. |
| **Body3D.Collided(a, b)** | **`BODY3D.COLLIDED`** | Returns TRUE if two bodies are in contact. |

**Aliases:** **`PHYSICS.*`** mirrors several **`PHYSICS3D.*`** names.

See also: [PHYSICS3D.md](../PHYSICS3D.md).
