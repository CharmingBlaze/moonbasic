# Audio Stream Commands

Commands for working with raw audio streams. This is an advanced feature for generating or processing audio in real-time.

## Core Workflow

1.  **Create Stream**: Use `AudioStream.Make()` to create a stream with a specific sample rate, bit depth, and channel count.
2.  **Generate or Load PCM Data**: Create an array holding raw audio data (Pulse-Code Modulation).
3.  **Update Stream**: In a loop, use `AudioStream.Update()` to push your PCM data to the stream's buffer.
4.  **Control Playback**: Use `AudioStream.Play()`, `AudioStream.Pause()`, etc.
5.  **Cleanup**: Free the stream with `AudioStream.Free()`.

---

### `AudioStream.Make(sampleRate, bitDepth, channels)`

Creates a new audio stream for custom audio playback.

- `sampleRate`: The number of samples per second (e.g., 44100).
- `bitDepth`: The number of bits per sample (e.g., 16).
- `channels`: The number of channels (1 for mono, 2 for stereo).

Returns a handle to the audio stream.

---

### `AudioStream.Update(streamHandle, pcmDataArray)`

Sends a chunk of raw PCM data to the audio stream's buffer for playback.

- `streamHandle`: The handle of the audio stream.
- `pcmDataArray`: A handle to a 1D array containing the raw audio samples.

---

### `AudioStream.Play(streamHandle)`

Starts playback of the audio stream.

---

### `AudioStream.IsPlaying(streamHandle)`

Returns `TRUE` if the stream is currently playing.

---

### `AudioStream.Free(streamHandle)`

Frees the audio stream resource.
