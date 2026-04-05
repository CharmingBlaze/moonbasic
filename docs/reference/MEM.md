# Binary memory buffers (`MEM.*`)

Fixed-size **byte buffers** on the VM heap (`Mem` handles) for low-level packing, binary protocols, or interop-style layouts. Implemented in `runtime/mbmem`.

**Limits:** Each block is at most **256 MiB** (`maxMemBlock`). Integer offsets and sizes must be non-negative and in range.

**Endianness:** Multi-byte reads and writes use **little-endian** order (`binary.LittleEndian`).

---

## Allocation

### `Mem.Make(size)` → handle

Allocates `size` bytes, zero-initialized.

### `Mem.Free(mem)`

Frees the buffer.

### `Mem.Size(mem)` → int

Byte length of the buffer.

### `Mem.Clear(mem)`

Sets every byte to zero (`clear`).

### `Mem.Copy(src, dst, srcOff, dstOff, size)`

Copies `size` bytes from `src[srcOff:]` to `dst[dstOff:]`. Overlap rules follow Go `copy`; out-of-range access errors.

---

## Read / write primitives

Offsets are **byte indices**. `GET` forms return values; `SET` forms return nothing.

| Command | Width | Notes |
|---------|-------|------|
| `Mem.GetByte` / `SetByte` | 1 byte | `SetByte` value 0–255. |
| `Mem.GetWord` / `SetWord` | 2 bytes | `SetWord` value 0–65535. |
| `Mem.GetDword` / `SetDword` | 4 bytes | `SetDword` uses low 32 bits (two’s complement as `uint32`). |
| `Mem.GetFloat` / `SetFloat` | 4 bytes | IEEE **32-bit** float. |
| `Mem.GetString` / `SetString` | C string | **Get:** reads bytes until **NUL** or end of buffer. **Set:** writes bytes plus trailing **NUL**; needs `len(s)+1` bytes free at offset. |
