# Spawner Commands

Timed entity spawner: create a spawner that emits entities at a rate, at a world position.

## Commands

### `SPAWNER.MAKE(templateEntity, x, y, z)` 

Creates a spawner at world position `(x, y, z)` that clones `templateEntity` on each spawn event. Returns a **spawner handle**.

---

## See also

- [PROP.md](PROP.md) — static prop placement
- [ENTITY.md](ENTITY.md) — entity system
- [EVENT.md](EVENT.md) — trigger spawners from events
