# Tweens (`TWEEN.*`)

Keyframe-style **tween chains** that read and write **global variables** by name (float-backed). Implemented in `runtime/mbtween`.

**Requirements:** The host must configure **get/set global** accessors on the tween module; otherwise updates error. The tween holds a handle in the heap (`Tween` type).

---

## Building a tween

### `Tween.Make()` → handle

Creates an empty tween (default: **one** loop, no yoyo).

### `Tween.To(tween, varName$, target#, seconds#, easing$)`

Appends a segment that animates the **global** named `varName` (folded uppercase) from its **current value at segment start** toward `target` over `seconds` (must be > 0).

**Easing** names (case-insensitive; unknown names fall back to linear):

| Name | Effect |
|------|--------|
| `linear` or `""` | Linear |
| `easein` | Quadratic ease-in |
| `easeout` | Quadratic ease-out |
| `easeinout` | Quadratic ease-in-out |
| `bounce` | Ease-out bounce |
| `elastic` | Ease-out elastic |
| `back` | Ease-out overshoot |
| `circ` | Ease-out circular |
| `expo` | Ease-out exponential |
| `sine` | Ease-in-out sine |

### `Tween.Then(...)`

Alias of `Tween.To` — append another segment after the previous.

### `Tween.OnComplete(tween, functionName$)`

Registers a **parameterless user function** to run when the tween finishes all loops (not when stopped early). Cannot change while running.

### `Tween.Loop(tween, count)`

Sets loop count: `count <= 0` means **infinite** loops. Cannot change while running.

### `Tween.Yoyo(tween)`

When enabled, after the forward pass the tween runs **backward** through the same steps before counting a loop. Cannot change while running.

---

## Running

### `Tween.Start(tween)`

Begins playback from the first step (requires at least one `TO`).

### `Tween.Update(tween, dt#)`

Advances time; `dt` is clamped if negative. May complete steps, invoke **OnComplete**, or loop.

### `Tween.Stop(tween)`

Stops playback; does not call OnComplete.

---

## Notes

- You cannot append `TO` / change loop / yoyo / on-complete while `running`.
- Globals must hold **numeric** values for variables being tweened.
