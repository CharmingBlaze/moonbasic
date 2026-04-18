# AudioStream Commands

Low-level streaming audio buffers for real-time audio generation or processing.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create a stream with `AUDIOSTREAM.CREATE` (or `AUDIOSTREAM.MAKE`), specifying sample rate, bit depth, and channels.
2. Start playback with `AUDIOSTREAM.PLAY`.
3. Each frame, check `AUDIOSTREAM.ISREADY` and push PCM data with `AUDIOSTREAM.UPDATE`.
4. Control with `AUDIOSTREAM.PAUSE`, `AUDIOSTREAM.RESUME`, `AUDIOSTREAM.STOP`.
5. Adjust volume, pitch, and pan as needed.
6. Free with `AUDIOSTREAM.FREE`.

For file-based audio see [AUDIO.md](AUDIO.md). For in-memory wave manipulation see [WAVE.md](../reference/WAVE.md).

---

### `AUDIOSTREAM.CREATE(sampleRate, bitDepth, channels)` 

Creates a new audio stream. Returns a handle.

- `sampleRate`: Samples per second (e.g. 44100).
- `bitDepth`: Bits per sample (e.g. 16).
- `channels`: 1 for mono, 2 for stereo.

---

### `AUDIOSTREAM.MAKE(sampleRate, bitDepth, channels)` 

Alias for `AUDIOSTREAM.CREATE`.

---

### `AUDIOSTREAM.UPDATE(streamHandle, pcmDataHandle)` 

Pushes a chunk of raw PCM data to the stream's buffer.

- `pcmDataHandle`: Handle to a memory buffer or array of audio samples.

---

### `AUDIOSTREAM.ISREADY(streamHandle)` 

Returns `TRUE` if the stream's buffer is ready to accept more data.

---

### `AUDIOSTREAM.ISPLAYING(streamHandle)` 

Returns `TRUE` if the stream is currently playing.

---

### `AUDIOSTREAM.PLAY(streamHandle)` 

Starts playback of the audio stream.

---

### `AUDIOSTREAM.PAUSE(streamHandle)` 

Pauses playback. Resume with `AUDIOSTREAM.RESUME`.

---

### `AUDIOSTREAM.RESUME(streamHandle)` 

Resumes a paused stream.

---

### `AUDIOSTREAM.STOP(streamHandle)` 

Stops playback and resets the stream position.

---

### `AUDIOSTREAM.SETVOLUME(streamHandle, volume)` 

Sets the playback volume (0.0–1.0).

---

### `AUDIOSTREAM.SETPITCH(streamHandle, pitch)` 

Sets the pitch multiplier (1.0 = normal).

---

### `AUDIOSTREAM.SETPAN(streamHandle, pan)` 

Sets the stereo pan (−1.0 = left, 0.0 = centre, 1.0 = right).

---

### `AUDIOSTREAM.FREE(streamHandle)` 

Frees the audio stream resource.

---

## Full Example

This example creates a stream and fills it with a simple sine tone.

```basic
AUDIO.INIT()
sampleRate = 44100
stream = AUDIOSTREAM.CREATE(sampleRate, 16, 1)
AUDIOSTREAM.PLAY(stream)
AUDIOSTREAM.SETVOLUME(stream, 0.5)

; Generate and push a 440Hz sine tone
buf = MEM.CREATE(sampleRate * 2)
FOR i = 0 TO sampleRate - 1
    sample = INT(SIN(2.0 * PI * 440.0 * FLOAT(i) / FLOAT(sampleRate)) * 32000)
    MEM.WRITESHORT(buf, i * 2, sample)
NEXT

WHILE NOT WINDOW.SHOULDCLOSE()
    IF AUDIOSTREAM.ISREADY(stream)
        AUDIOSTREAM.UPDATE(stream, buf)
    END IF
    RENDER.BEGINFRAME()
    RENDER.ENDFRAME()
WEND

AUDIOSTREAM.FREE(stream)
MEM.FREE(buf)
AUDIO.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `AUDIOSTREAM.GETVOLUME(stream)` | Returns current volume 0.0–1.0. |
| `AUDIOSTREAM.GETPITCH(stream)` | Returns current pitch multiplier. |
| `AUDIOSTREAM.GETPAN(stream)` | Returns current stereo pan -1.0–1.0. |

## See also

- [AUDIO.md](AUDIO.md) — `AUDIO.SEEKMUSIC`, `AUDIO.SETMASTERVOLUME`
- [MEM.md](MEM.md) — raw buffer for PCM data
