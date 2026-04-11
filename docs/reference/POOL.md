# Object pools (`POOL.*`)

Reuse **heap handles** produced by a **factory user function**, with optional **reset** callbacks when returning objects to the pool. Implemented in `runtime/mbpool`.

Typical use: bullets, particles, or UI rows where allocation churn should be avoided.

---

## Creation

### `Pool.Make(name, capacity)` → handle

`capacity` must be a positive integer — maximum **checked-out** objects (`GET` fails if `busy` count would exceed `max`).

### `Pool.SetFactory(pool, factoryFunctionName)` / `Pool.SetReset(pool, resetFunctionName)`

- **Factory** — user function invoked with **no arguments**; must return a **handle** to push into the pool or hand to the caller.
- **Reset** — optional; called as `reset(handle)` when returning an object to the free list.

Names are resolved at `GET` / `RETURN` time via the runtime’s user invoker.

### `Pool.Prewarm(pool)`

Allocates up to **capacity** objects by calling the factory repeatedly and storing them in the free list. Fails if the factory is not set.

---

## Use

### `Pool.Get(pool)` → handle

Checks out an object: pops from **free** if any, otherwise calls **factory**. Errors if at capacity or factory unset.

### `Pool.Return(pool, handle)`

Moves `handle` from **busy** to **free** after optional **reset**; errors if the handle was not checked out from this pool.

### `Pool.Free(pool)`

Frees the pool object itself and releases all **busy** and **free** child handles via the heap.

---

## Contract

- **Factory** must return a valid handle the pool may own.
- Returning a handle the pool does not track results in an error.
