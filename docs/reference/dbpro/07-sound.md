# DBPro — Sound

moonBASIC: **`AUDIO.*`**, **`SOUND.*`** (wave-backed) — [AUDIO.md](../AUDIO.md).

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **LOAD SOUND (file, snd)** | ✓ **`Audio.LoadSound()`** | Returns **handle**. |
| **DELETE SOUND (snd)** | ✓ **`Sound.Free()`** | |
| **PLAY SOUND (snd)** | ✓ **`Audio.Play()`** | |
| **LOOP SOUND (snd)** | ≈ **`Audio.Play()`** (check loop flag) | |
| **STOP SOUND (snd)** | ✓ **`Audio.Stop()`** | |
| **PAUSE SOUND (snd)** | ✓ **`Audio.Pause()`** | |
| **RESUME SOUND (snd)** | ✓ **`Audio.Resume()`** | |
| **SET SOUND VOLUME (snd, vol)** | ✓ **`Audio.SetVolume()`** | |
| **LOAD MUSIC (file, mus)** | ✓ **`Audio.LoadMusic()`** | Returns **handle**. |
| **DELETE MUSIC (mus)** | ✓ **`Music.Free()`** | |
| **PLAY MUSIC (mus)** | ✓ **`Audio.Play()`** | |
| **STOP MUSIC (mus)** | ✓ **`Audio.Stop()`** | |
| **PAUSE MUSIC (mus)** | ✓ **`Audio.Pause()`** | |
| **RESUME MUSIC (mus)** | ✓ **`Audio.Resume()`** | |
| **SET MUSIC VOLUME (mus, vol)** | ✓ **`Audio.SetVolume()`** | |
| **LOAD 3D SOUND** | ✓ **`Audio.LoadSound()`** | Use with **`Audio.Listener()`**. |
| **PLAY 3D SOUND** | ✓ **`Entity.EmitSound()`** | Spatial positioning. |
