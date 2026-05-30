# Char Commands

`CHAR.*` is the entity-integer alias namespace for the Jolt Kinematic Character Controller — every command mirrors one in [PLAYER.md](PLAYER.md) or [CHARACTERREF.md](CHARACTERREF.md). Prefer those namespaces in new code.

## Aliases

| `CHAR.*` command | Preferred equivalent |
|---|---|
| `CHAR.CREATE(entity)` | `PLAYER.CREATE(entity)` |
| `CHAR.CREATE(entity, r, h)` | `PLAYER.CREATE(entity, r, h)` |
| `CHAR.MOVE(entity, dx, dz, speed)` | `PLAYER.MOVE(entity, dx, dz, speed)` |
| `CHAR.JUMP(entity, impulse)` | `PLAYER.JUMP(entity, impulse)` |
| `CHAR.ISGROUNDED(entity)` | `PLAYER.ISGROUNDED(entity)` |
| `CHAR.ISGROUNDED(entity, coyote)` | `PLAYER.ISGROUNDED(entity, coyote)` |
| `CHAR.GETGROUNDSTATE(entity)` | `PLAYER.GETGROUNDSTATE(entity)` |
| `CHAR.ISONSTEEPSLOPE(entity)` | `PLAYER.ISONSTEEPSLOPE(entity)` |
| `CHAR.GETPOSITIONX/Y/Z(entity)` | `PLAYER.GETPOSITIONX/Y/Z(entity)` |
| `CHAR.GETVELOCITYX/Y/Z(entity)` | `PLAYER.GETVELOCITYX/Y/Z(entity)` |
| `CHAR.GETSPEED(entity)` | `PLAYER.GETSPEED(entity)` |
| `CHAR.GETONSLOPE(entity)` | `PLAYER.GETONSLOPE(entity)` |
| `CHAR.GETONWALL(entity)` | `PLAYER.GETONWALL(entity)` |
| `CHAR.GETSLOPEANGLE(entity)` | `PLAYER.GETSLOPEANGLE(entity)` |
| `CHAR.GETISJUMPING(entity)` | `PLAYER.GETISJUMPING(entity)` |
| `CHAR.GETISFALLING(entity)` | `PLAYER.GETISFALLING(entity)` |
| `CHAR.GETISSLIDING(entity)` | `PLAYER.GETISSLIDING(entity)` |
| `CHAR.GETCEILING(entity)` | `PLAYER.GETCEILING(entity)` |
| `CHAR.SETPADDING(entity, p)` | `PLAYER.SETPADDING(entity, p)` |
| `CHAR.SETSTEP(entity, h)` | `PLAYER.SETSTEPOFFSET(entity, h)` |
| `CHAR.SETSLOPE(entity, deg)` | `PLAYER.SETSLOPELIMIT(entity, deg)` |
| `CHAR.STICK(entity, dist)` | `PLAYER.SETSTICKFLOOR(entity, dist)` |
| `CHAR.MOVEWITHCAMERA(entity, cam, speed, mx, mz)` | `PLAYER.MOVEWITHCAMERA(entity, cam, speed, mx, mz)` |
| `CHAR.NAVTO(entity, x, y, z)` | `PLAYER.NAVTO(entity, x, y, z)` |
| `CHAR.NAVUPDATE(entity)` | `PLAYER.NAVUPDATE(entity)` |
| `CHAR.DIST(entityA, entityB)` | `ENTITY.DIST(entityA, entityB)` |
| `CHAR.UPDATE(dt)` | `PLAYER.UPDATE(dt)` |
| `CHAR.GETVX/VY/VZ(entity)` | `PLAYER.GETVELOCITYX/Y/Z(entity)` |

## Extended Command Reference (getter aliases)

All `CHAR.*` getters below are direct aliases of the matching `PLAYER.*` forms.

| Command | Equivalent |
|--------|------------|
| `CHAR.GETX(e)` / `CHAR.GETY(e)` / `CHAR.GETZ(e)` | `PLAYER.GETX/Y/Z` |
| `CHAR.GETYAW(e)` / `CHAR.GETPITCH(e)` / `CHAR.GETROLL(e)` | `PLAYER.GETYAW/PITCH/ROLL` |
| `CHAR.GETROTATIONYAW(e)` / `GETROTATIONPITCH` / `GETROTATIONROLL` | `PLAYER.GETROTATION*` |
| `CHAR.GETPOSITIONY(e)` / `CHAR.GETPOSITIONZ(e)` | `PLAYER.GETPOSITIONY/Z` |
| `CHAR.GETVELOCITYY(e)` / `CHAR.GETVELOCITYZ(e)` | `PLAYER.GETVELOCITYY/Z` |
| `CHAR.GETVY(e)` / `CHAR.GETVZ(e)` | `PLAYER.GETVY/VZ` |
| `CHAR.GETGROUNDVELOCITYX(e)` / `CHAR.GETGROUNDVELOCITYY(e)` / `CHAR.GETGROUNDVELOCITYZ(e)` | `PLAYER.GETGROUNDVELOCITY*` |
| `CHAR.MAKE(entity, r, h)` | Deprecated alias of `CHAR.CREATE`. |
| `CHAR.GETGROUNDED(e)` | `PLAYER.GETGROUNDED` |
| `CHAR.GETHEIGHT(e)` / `CHAR.GETRADIUS(e)` | `PLAYER.GETHEIGHT/RADIUS` |
| `CHAR.GETCAPSULEHEIGHT(e)` / `CHAR.GETCAPSULERADIUS(e)` | `PLAYER.GETCAPSULE*` |
| `CHAR.GETGRAVITY(e)` / `CHAR.GETGRAVITYSCALE(e)` | `PLAYER.GETGRAVITY/SCALE` |
| `CHAR.GETMAXSLOPE(e)` / `CHAR.GETSTEPHEIGHT(e)` / `CHAR.GETSNAPDISTANCE(e)` | `PLAYER.GETMAXSLOPE/STEPHEIGHT/SNAPDISTANCE` |
| `CHAR.GETFRICTION(e)` | `PLAYER.GETFRICTION` |
| `CHAR.GETSHAPETYPE(e)` | `PLAYER.GETSHAPETYPE` |
| `CHAR.GETLAYER(e)` / `CHAR.GETMASK(e)` | `PLAYER.GETLAYER/MASK` |
| `CHAR.GETCOLLISIONENABLED(e)` | `PLAYER.GETCOLLISIONENABLED` |

---

## See also

- [PLAYER.md](PLAYER.md) — canonical KCC API
- [CHARACTERREF.md](CHARACTERREF.md) — handle-based KCC
- [CHARCONTROLLER.md](CHARCONTROLLER.md) — capsule controller
