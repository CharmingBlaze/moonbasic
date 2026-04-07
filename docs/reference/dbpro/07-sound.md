# DBPro — Sound

moonBASIC: **`AUDIO.*`**, **`SOUND.*`** (wave-backed) — [AUDIO.md](../AUDIO.md).

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **LOAD SOUND (file$, snd)** | ✓ **`AUDIO.LOADSOUND`** | Returns handle. |
| **DELETE SOUND** | ≈ stop + release patterns | **`SOUND.FREE`** / unload depending on path. |
| **PLAY SOUND** / **LOOP SOUND** | ✓ **`AUDIO.PLAY`** | |
| **STOP SOUND** | ✓ **`AUDIO.STOP`** | |
| **SET SOUND VOLUME** / **PAN** / **SPEED** | ✓ **`AUDIO.SETSOUNDVOLUME`**, **`SETSOUNDPAN`**, **`SETSOUNDPITCH`** | |
| **SOUND PLAYING** | ✓ **`AUDIO.ISSOUNDPLAYING`** | |
