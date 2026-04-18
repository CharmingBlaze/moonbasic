# Mem Commands

Fixed-size byte buffers for low-level binary packing, protocols, and interop layouts.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Allocate a buffer with `MEM.MAKE`.
2. Write data with `MEM.SETBYTE`, `MEM.SETWORD`, `MEM.SETFLOAT`, `MEM.SETSTRING`, etc.
3. Read back with the matching `MEM.GET*` commands.
4. Copy between buffers with `MEM.COPY`.
5. Free with `MEM.FREE`.

**Limits:** Max 256 MiB per block. Little-endian byte order for multi-byte types.

---

### `MEM.MAKE(size)` / `FREE`
Allocates or releases a fixed-size byte buffer.

- **Arguments**:
    - `size`: (Integer) Buffer size in bytes.
- **Returns**: (Handle) The new memory handle.

---

### `MEM.SETBYTE(handle, offset, value)` / `SETWORD` / `SETFLOAT`
Writes binary data into the buffer at a specific offset.

- **Arguments**:
    - `handle`: (Handle) The buffer.
    - `offset`: (Integer) 0-based byte offset.
    - `value`: (Integer/Float) Data to write.
- **Returns**: (Handle) The memory handle (for chaining).

---

### `MEM.GETBYTE(handle, offset)` / `GETWORD` / `GETFLOAT`
Reads binary data from the buffer.

- **Returns**: (Integer / Float)

---

### `MEM.SETSTRING(handle, offset, text)` / `GETSTRING`
Reads or writes NUL-terminated strings.

---

### `MEM.COPY(src, dst, srcOff, dstOff, count)`
Copies raw bytes between two buffers.

---

## Full Example

This example packs a position into a binary buffer and reads it back.

```basic
buf = MEM.MAKE(12)
MEM.SETFLOAT(buf, 0, 1.5)
MEM.SETFLOAT(buf, 4, 2.0)
MEM.SETFLOAT(buf, 8, 3.5)

x = MEM.GETFLOAT(buf, 0)
y = MEM.GETFLOAT(buf, 4)
z = MEM.GETFLOAT(buf, 8)
PRINT "Position: " + STR(x) + ", " + STR(y) + ", " + STR(z)

MEM.FREE(buf)
```
