# Check Commands

Entity visibility checks.

## Commands

### `CHECK.INVIEW(entityId)` 

Returns `TRUE` if `entityId` is within the current camera frustum. Same test as `ENTITY.INFRUSTUM`. Use to skip drawing off-screen entities.

```basic
IF CHECK.INVIEW(enemy) THEN
    ENTITY.DRAW(enemy)
END IF
```

---

## See also

- [ENTITY.md](ENTITY.md) — `ENTITY.INFRUSTUM`, `ENTITY.SETVISIBLE`
- [CAMERA.md](CAMERA.md) — camera setup
