# Res Commands

Resolve resource paths relative to the running executable and check file existence.

## Commands

### `RES.PATH(localPath)` 

Returns the absolute path of `localPath` resolved relative to the executable directory. Use for portable asset loading that works regardless of working directory.

```basic
tex = TEXTURE.LOAD(RES.PATH("assets/player.png"))
```

---

### `RES.EXISTS(path)` 

Returns `TRUE` if `path` exists on disk. Alias of `UTIL.FILEEXISTS`.

```basic
IF RES.EXISTS(RES.PATH("save.json")) THEN
    ; load save data
END IF
```

---

## See also

- [FILE.md](FILE.md) — file I/O
- [SAVE.md](SAVE.md) — game save data
- [CONFIG.md](CONFIG.md) — config file loading
