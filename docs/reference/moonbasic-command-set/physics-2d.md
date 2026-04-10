# Physics 2D (Box2D)

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **Physics2D.Start** (optional gravity) | **`PHYSICS2D.START`** with **0** or **2** args `(gx, gy)`; or **`PHYSICS2D.START`** then **`PHYSICS2D.SETGRAVITY`** | World is global; **`PHYSICS2D.STOP`** clears auxiliary state (contacts, debug segments). |
| **Physics2D.Stop** | **`PHYSICS2D.STOP`** | |
| **Physics2D.Step** | **`PHYSICS2D.STEP`** | Call after optional **`PHYSICS2D.SETSTEP`**, **`PHYSICS2D.SETITERATIONS`**. After each step, contact queries and debug segments update. |
| **Physics2D.DebugDraw** | **`PHYSICS2D.DEBUGDRAW`** `(mode)` — `0` off, non-zero collects segments | **`PHYSICS2D.GETDEBUGSEGMENTS`** returns a **float array** handle: `[x1,y1,x2,y2,…]` line pairs in world space; **`FREE`** the array when done. |
| **Body2D.Create** | **`BODY2D.MAKE`** `(kind)` → **`BODY2D.ADDRECT`** / **`ADDCIRCLE`** / **`ADDPOLYGON`** → **`BODY2D.COMMIT`** `(template, x, y)` | Template freed on commit; committed body returns a new handle. **`ADDRECT`**: `(h, w, h)` or `(h, w, h, density, friction, restitution)`; **`ADDCIRCLE`**: `(h, r)` or `(h, r, density, friction, restitution)`; **`ADDPOLYGON`**: `(h, points[])` flat float array `[x0,y0,x1,y1,…]` convex CCW. |
| **Body2D.Delete** | **`BODY2D.FREE`** | |
| **SetPosition / SetAngle / velocity** | **`BODY2D.SETPOS`**, **`SETROT`**, **`SETLINEARVELOCITY`**, **`SETANGULARVELOCITY`** | |
| **ApplyForce / Impulse** | **`BODY2D.APPLYFORCE`**, **`APPLYIMPULSE`** | |
| **GetX / Y / angle** | **`BODY2D.X`**, **`Y`**, **`ROT`** / **`GETROT`** | |
| **AddBox / Circle / Polygon** | **`BODY2D.ADDRECT`**, **`ADDCIRCLE`**, **`ADDPOLYGON`** | See **`COMMIT`** row for material extras. |
| **Joint2D.Distance / Revolute / Prismatic** | **`JOINT2D.DISTANCE`** `(bodyA, bodyB, ax, ay, bx, by)`; **`JOINT2D.REVOLUTE`** `(bodyA, bodyB, x, y)`; **`JOINT2D.PRISMATIC`** `(bodyA, bodyB, x, y, ax, ay)` | Joint returns a handle; **`JOINT2D.FREE`** destroys the joint. |
| **Joint2D.Delete** | **`JOINT2D.FREE`** | |
| **Collision queries** | **`BODY2D.COLLIDED`**, **`COLLISIONOTHER`**, **`COLLISIONNORMAL`**, **`COLLISIONPOINT`** | Updated after each **`PHYSICS2D.STEP`**. **`COLLISIONNORMAL`** / **`COLLISIONPOINT`** return **`Point2D`** instance handles. Bodies must be created after **`PHYSICS2D.START`** so **`UserData`** maps contacts. |

**Legacy:** **`BOX2D.*`** aliases remain (`WORLDCREATE`, `BODYCREATE`, `FIXTUREBOX`, `FIXTURECIRCLE`).

See also: [PHYSICS2D.md](../PHYSICS2D.md).
