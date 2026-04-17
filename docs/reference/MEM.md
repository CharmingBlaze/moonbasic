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

### `MEM.MAKE(size)` 

Allocates `size` bytes, zero-initialized.

---

### `MEM.FREE(memHandle)` 

Frees the buffer.

---

### `MEM.SIZE(memHandle)` 

Returns the byte length of the buffer.

---

### `MEM.CLEAR(memHandle)` 

Sets every byte to zero.

---

### `MEM.COPY(src, dst, srcOff, dstOff, size)` 

Copies `size` bytes from `src[srcOff:]` to `dst[dstOff:]`.

---

### `MEM.GETBYTE(memHandle, offset)` / `MEM.SETBYTE(memHandle, offset, value)` 

Read/write a single byte (0–255).

---

### `MEM.GETWORD(memHandle, offset)` / `MEM.SETWORD(memHandle, offset, value)` 

Read/write a 16-bit unsigned integer (0–65535), little-endian.

---

### `MEM.GETDWORD(memHandle, offset)` / `MEM.SETDWORD(memHandle, offset, value)` 

Read/write a 32-bit integer, little-endian.

---

### `MEM.GETFLOAT(memHandle, offset)` / `MEM.SETFLOAT(memHandle, offset, value)` 

Read/write an IEEE 32-bit float, little-endian.

---

### `MEM.GETSTRING(memHandle, offset)` / `MEM.SETSTRING(memHandle, offset, value)` 

Read/write a NUL-terminated C string. `SETSTRING` needs `LEN(value)+1` bytes free at offset.

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
