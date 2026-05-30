# Wave Commands

Commands for loading and manipulating raw wave audio data in memory.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**): **Core Workflow** → per-command sections with **`---`** separators → **Full Example**.

## Core Workflow

`WAVE` commands are used for offline audio processing. You can load a wave file, perform operations like cropping or formatting, and then either save it back to a file or convert it into a playable `SOUND`.

---

### `WAVE.LOAD(filePath)` 

Loads a wave file (`.wav`) into memory. Returns a handle to the wave data.

- `filePath`: The path to the `.wav` file.

---

### `WAVE.COPY(waveHandle)` 

Creates a new, independent copy of a wave resource.

---

### `WAVE.CROP(waveHandle, startFrame, endFrame)` 

Crops the wave data to a new length, from `startFrame` to `endFrame`.

---

### `WAVE.EXPORT(waveHandle, filePath)` 

Saves the wave data to a new `.wav` file.

---

### `WAVE.FREE(waveHandle)` 

Frees the wave data from memory.

---

### `SOUND.FROMWAVE(waveHandle)` 

Creates a playable `SOUND` handle from raw `WAVE` data. The resulting sound can be used with **`AUDIO.PLAY()`**.

---

## Full Example

This example loads a wave file, creates a shorter version by cropping it, and saves the result as a new file.

```basic
AUDIO.INIT()

; Load the original wave
original_wave = WAVE.LOAD("my_sound.wav")
ASSERT(original_wave <> 0, "Failed to load my_sound.wav")

; Create a copy to modify
cropped_wave = WAVE.COPY(original_wave)

; Crop the copy to the first 22050 frames (0.5 seconds at 44100Hz)
WAVE.CROP(cropped_wave, 0, 22050)

; Export the cropped version
WAVE.EXPORT(cropped_wave, "my_sound_short.wav")
PRINT "Created my_sound_short.wav"

; You can also create a playable sound from it
playable_sound = SOUND.FROMWAVE(cropped_wave)
AUDIO.PLAY(playable_sound)
SLEEP 1000

; Cleanup
WAVE.FREE(original_wave)
WAVE.FREE(cropped_wave)
SOUND.FREE(playable_sound)
AUDIO.CLOSE()
```
