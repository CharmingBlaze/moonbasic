# Clipboard Commands

Access the system clipboard.

## Commands

### `CLIPBOARD.GETIMAGE()` 

Returns a texture handle containing the current clipboard image (if any), or `0` if the clipboard does not contain an image.

```basic
tex = CLIPBOARD.GETIMAGE()
IF tex THEN
    DRAW.TEXTURE(tex, 0, 0, 400, 300, 255, 255, 255, 255)
END IF
```

---

## See also

- [TEXTURE.md](TEXTURE.md) — texture handles
- [FILE.md](FILE.md) — file-based image loading
