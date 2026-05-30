# Enemy Commands

High-level enemy movement shorthand. Moves an entity along a `PATH` handle each frame.

## Commands

### `ENEMY.FOLLOWPATH(entityId, pathHandle, speed)` 

Advances `entityId` along `pathHandle` toward the next waypoint at `speed` world units per second. When the last waypoint is reached the entity stops.

---

## See also

- [PATH.md](PATH.md) — path waypoint handles
- [NAV.md](NAV.md) — building paths
- [ENT.md](ENT.md) — health, teams, shoot
- [BTREE.md](BTREE.md) — AI behaviour trees
