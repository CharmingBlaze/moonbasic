# Data transforms (`DATA.*`)

String-in / string-out helpers for **compression**, **Base64**, and **checksums**. Implemented in `runtime/mbdata`.

---

## Compression (Zstandard)

### `Data.Compress(raw)` → string

Compresses the input bytes with **Zstandard** (`zstd`). Returns binary data as a Go string payload (opaque bytes).

### `Data.Decompress(zstdData)` → string

Decompresses a buffer previously produced by `Data.Compress`. Errors if the payload is invalid.

---

## Base64

### `Data.EncodeBase64(raw)` → string

Standard Base64 encoding.

### `Data.DecodeBase64(encoded)` → string

Standard Base64 decoding. Errors on invalid input.

---

## Hashes and digests

All of these take a **string** of raw bytes and return either an integer or a **lowercase hex** string.

| Command | Returns | Implementation |
|---------|---------|------------------|
| `Data.CRC32` / `Data.ComputeCRC32` | int64 (uint32 value) | IEEE CRC32 |
| `Data.MD5` / `Data.ComputeMD5` | string (32 hex chars) | MD5 |
| `Data.SHA1` / `Data.ComputeSHA1` | string (40 hex chars) | SHA-1 |

`Compute*` names are aliases of the same functions.
