# Free Commands

Global handle cleanup.

## Commands

### `FREE.ALL()` 

Frees **all** active handles of every type (entities, textures, bodies, sounds, etc.) in one call. Use during scene transitions or at shutdown when you want a clean slate without tracking individual handles.

> **Warning:** This frees everything. Do not call mid-scene unless you intend to rebuild all resources.

---

## See also

- [ENTITY.md](ENTITY.md) — `ENTITY.FREEALL`
- [TEXTURE.md](TEXTURE.md) — `TEXTURE.UNLOAD`
- [AUDIO.md](AUDIO.md) — audio cleanup
