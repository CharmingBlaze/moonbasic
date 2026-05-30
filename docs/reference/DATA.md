# Data Commands

Compression (Zstandard), Base64 encoding, and cryptographic hash helpers.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Pass a string of raw bytes into any `DATA.*` function. Compression and encoding functions return transformed strings; hash functions return integers or hex strings.

---

### `DATA.COMPRESS(raw)` 

Compresses the input with Zstandard. Returns opaque binary data as a string.

---

### `DATA.DECOMPRESS(zstdData)` 

Decompresses a buffer previously produced by `DATA.COMPRESS`. Errors on invalid input.

---

### `DATA.ENCODEBASE64(raw)` 

Returns the standard Base64 encoding of the input string.

---

### `DATA.DECODEBASE64(encoded)` 

Decodes a Base64 string. Errors on invalid input.

---

### `DATA.CRC32(raw)` 

Returns the IEEE CRC-32 checksum as an integer.

---

### `DATA.COMPUTECRC32(raw)` 

Alias for `DATA.CRC32`.

---

### `DATA.MD5(raw)` 

Returns the MD5 digest as a 32-character lowercase hex string.

---

### `DATA.COMPUTEMD5(raw)` 

Alias for `DATA.MD5`.

---

### `DATA.SHA1(raw)` 

Returns the SHA-1 digest as a 40-character lowercase hex string.

---

### `DATA.COMPUTESHA1(raw)` 

Alias for `DATA.SHA1`.

---

## Full Example

This example compresses data, encodes it, and verifies with a hash.

```basic
original = "Hello, moonBASIC! This is test data for compression."
hash_before = DATA.SHA1(original)

; Compress and Base64-encode for transport
compressed = DATA.COMPRESS(original)
encoded = DATA.ENCODEBASE64(compressed)
PRINT "Encoded length: " + STR(LEN(encoded))

; Decode and decompress
decoded = DATA.DECODEBASE64(encoded)
restored = DATA.DECOMPRESS(decoded)
hash_after = DATA.SHA1(restored)

ASSERT(hash_before = hash_after, "Round-trip hash mismatch!")
PRINT "Round-trip OK: " + hash_after
```
