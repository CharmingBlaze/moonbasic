# Pool Commands

Object pool for reusing heap handles with factory and reset callbacks, reducing allocation churn.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create a pool with `POOL.MAKE`.
2. Set a factory function with `POOL.SETFACTORY` and optionally a reset with `POOL.SETRESET`.
3. Prewarm with `POOL.PREWARM` to pre-allocate objects.
4. Check out objects with `POOL.GET`, return them with `POOL.RETURN`.
5. Free the pool with `POOL.FREE`.

---

### `POOL.MAKE(name, capacity)` 

`capacity` must be a positive integer — maximum **checked-out** objects (`GET` fails if `busy` count would exceed `max`).

---

### `POOL.SETFACTORY(poolHandle, factoryFunctionName)` 

Sets the factory function (called with no arguments, must return a handle).

---

### `POOL.SETRESET(poolHandle, resetFunctionName)` 

Sets an optional reset function called as `reset(handle)` when returning an object.

---

### `POOL.PREWARM(poolHandle)` 

Pre-allocates up to capacity by calling the factory. Fails if no factory is set.

---

### `POOL.GET(poolHandle)` 

Checks out an object from the pool. Pops from the free list or calls the factory. Errors if at capacity.

---

### `POOL.RETURN(poolHandle, objectHandle)` 

Returns an object to the pool’s free list, calling reset if set.

---

### `POOL.FREE(poolHandle)` 

Frees the pool and all its managed objects.

---

## Full Example

This example pools bullet entities to avoid per-frame allocation.

```basic
FUNCTION MakeBullet()
    b = ENTITY.CREATE()
    ENTITY.SETSCALE(b, 0.1, 0.1, 0.1)
    RETURN b
END FUNCTION

FUNCTION ResetBullet(b)
    ENTITY.SETPOS(b, 0, -100, 0)
END FUNCTION

pool = POOL.MAKE("bullets", 50)
POOL.SETFACTORY(pool, "MakeBullet")
POOL.SETRESET(pool, "ResetBullet")
POOL.PREWARM(pool)

; Fire a bullet
bullet = POOL.GET(pool)
ENTITY.SETPOS(bullet, px, py, pz)

; Return when off-screen
POOL.RETURN(pool, bullet)

; Cleanup
POOL.FREE(pool)
```
