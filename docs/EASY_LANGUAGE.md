# Easy power: helpers without hiding math

moonBASIC aims for a **low floor and a high ceiling**: you can ship games with **named patterns** (movement, snapping, easing, spawning) while **`MATH.*`, `VEC2.*`, `VEC3.*`, and friends stay available** for anyone who wants full control.

This page is the **design stance** for contributors and doc authors. It does not remove or deprecate low-level APIs.

---

## 1. Two tracks, one language

| Track | Role | Audience |
|--------|------|----------|
| **Expressive helpers** | One line per *intent* (“move toward player”, “snap feet to terrain”, “orbit camera”) | Beginners and fast iteration |
| **Full math / vectors** | Primitives and composition (`SIN`, `SQRT`, `LERP`, matrices, …) | Custom algorithms, tools, and learning |

Helpers should **read like gameplay**, not like a textbook. When a pattern appears in **multiple samples** or **three+ lines of the same structure**, it is a candidate for a builtin or a documented idiom.

---

## 2. Where helpers already live

- **[reference/LESS_MATH.md](reference/LESS_MATH.md)** — Distance, spawn rings, camera-relative WASD, terrain snap, streaming center; sample: `examples/terrain_chase/main.mb`.
- **[reference/GAME_MATH_HELPERS.md](reference/GAME_MATH_HELPERS.md)** — `MATH.*` helpers for 2D/3D games: horizontal (`XZ`) distance, `YAWFROMXZ`, radian angle diff, `SMOOTHERSTEP`, `DIST2D` / `DISTSQ2D`.
- **[reference/GAMEHELPERS.md](reference/GAMEHELPERS.md)** — Landing, platform lists, third-person orbit deltas.
- **[reference/QOL.md](reference/QOL.md)** — `mbgame` shortcuts (`DT`, `SCREENW`, …) and links to movement, easing, noise.
- **[reference/EASING.md](reference/EASING.md)** — Smooth motion without hand-tuned lerp everywhere.
- **[reference/MOVEMENT.md](reference/MOVEMENT.md)** — Wrapping and step values on axes.
- **[COMMANDS.md](COMMANDS.md)** — Topic index; use it to discover “is there already a command for this?”

Prefer **reusing** these namespaces (`INPUT.*`, `Terrain.*`, `VEC2.*`, `MATH.*`, `ENTITY.*`, `World.*`) before inventing parallel names.

---

## 3. Rules of thumb for new helpers

1. **Name the task, not the formula** — e.g. `MOVE_TOWARD`, `PUSHOUT`, `SnapY`, not `SCALAR_PROJ_ON_RAY`.
2. **Prefer multi-return tuples** — `x, y = VEC2.NORMALIZE(dx, dy)` keeps user code short and avoids mystery temporaries.
3. **Document units** — radians vs degrees, world vs screen, and what “forward” means for camera-relative APIs.
4. **Safe edge cases** — zero-length vectors, degenerate circles, clamped delta; helpers should not surprise with NaNs or divide-by-zero when avoidable.
5. **One manifest row + implementation + `go run ./tools/apidoc`** when the public surface changes ([CONTRIBUTING.md](../CONTRIBUTING.md)).

---

## 4. Pattern categories (roadmap-style)

These are **directional** ideas—not a commitment list. They illustrate the kind of APIs that reduce “math literacy” as a prerequisite:

- **Space / camera** — Already: `INPUT.MOVEDIR`, orbit samples; possible extensions: more packaged camera rigs *as optional* wrappers, always alongside raw `Camera.*`.
- **2D / XZ gameplay** — Already: `VEC2.DIST`, `DISTSQ`, `PUSHOUT`, `MOVE_TOWARD`; extend when the same 5–10 lines repeat in examples.
- **Angles** — Already: `MATH.LERPANGLE`; avoid asking users to fix wrap by hand when a helper can mean “shortest path”.
- **Scalars over time** — Already: `MATH.APPROACH`, easing tables; good for zoom, health bars, UI without picking arbitrary lerp factors.
- **World / streaming** — Already: `World.SetCenterEntity`; keep terrain and world commands discoverable next to entity position.
- **Input** — Combine related reads (`MOUSEDELTA`) so one frame’s input is one obvious call.

When in doubt, add a **short subsection** to [LESS_MATH.md](reference/LESS_MATH.md) or the relevant `reference/*.md` page *before* adding a new global synonym—avoid namespace sprawl.

---

## 5. What we do not do

- **We do not remove** `MATH.SQRT`, trig, or vector primitives to “force” helpers.
- **We do not** hide performance tradeoffs when they matter (`DISTSQ` vs `DIST`)—we document them in one line.
- **We do not** replace the visible `WHILE` game loop with a black-box engine loop; helpers shorten *bodies*, not control flow ([PROGRAMMING.md](PROGRAMMING.md) §3).

---

## See also

- [PROGRAMMING.md](PROGRAMMING.md) — loop structure and destructuring.
- [reference/API_CONVENTIONS.md](reference/API_CONVENTIONS.md) — naming across modules.
- [GETTING_STARTED.md](GETTING_STARTED.md) — first projects and examples.
