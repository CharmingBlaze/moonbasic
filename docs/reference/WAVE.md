# Wave Commands

Commands for loading and manipulating raw wave audio data in memory.

## Core Workflow

`WAVE` commands are used for offline audio processing. You can load a wave file, perform operations like cropping or formatting, and then either save it back to a file or convert it into a playable `SOUND`.

---

### `Wave.Load(filePath$)`

Loads a wave file (`.wav`) into memory. Returns a handle to the wave data.

- `filePath$`: The path to the `.wav` file.

---

### `Wave.Copy(waveHandle)`

Creates a new, independent copy of a wave resource.

---

### `Wave.Crop(waveHandle, startFrame, endFrame)`

Crops the wave data to a new length, from `startFrame` to `endFrame`.

---

### `Wave.Export(waveHandle, filePath$)`

Saves the wave data to a new `.wav` file.

---

### `Wave.Free(waveHandle)`

Frees the wave data from memory.

---

### `Sound.FromWave(waveHandle)`

Creates a playable `SOUND` handle from raw `WAVE` data. The resulting sound can be used with `Audio.Play()`.

---

## Full Example

This example loads a wave file, creates a shorter version by cropping it, and saves the result as a new file.

```basic
Audio.Init()

; Load the original wave
original_wave = Wave.Load("my_sound.wav")
ASSERT(original_wave <> 0, "Failed to load my_sound.wav")

; Create a copy to modify
cropped_wave = Wave.Copy(original_wave)

; Crop the copy to the first 22050 frames (0.5 seconds at 44100Hz)
Wave.Crop(cropped_wave, 0, 22050)

; Export the cropped version
Wave.Export(cropped_wave, "my_sound_short.wav")
PRINT "Created my_sound_short.wav"

; You can also create a playable sound from it
playable_sound = Sound.FromWave(cropped_wave)
Audio.Play(playable_sound)
SLEEP 1000

; Cleanup
Wave.Free(original_wave)
Wave.Free(cropped_wave)
Sound.Free(playable_sound)
Audio.Close()
```
