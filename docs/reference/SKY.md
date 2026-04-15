# Sky (`SKY.*`)

Day/night **tinted sky dome** (drawn as a large sphere) with time-of-day in **hours** and configurable **day length**. **CGO + Raylib** required.

---

## `Sky.Create()` / `Sky.Make()` → handle

**Canonical:** **`SKY.CREATE`**. Deprecated alias: **`SKY.MAKE`**. Creates a sky object with default time and day length.

---

## `Sky.Free(sky)`

Frees the sky handle.

---

## `Sky.Update(sky, dt)`

Advances internal time using **`dt`** and day length.

---

## `Sky.Draw(sky)`

Draws the sky **before** terrain for typical frames (call order is user-defined).

---

## `Sky.SetTime(sky, hours)` / `Sky.SetDayLength(sky, seconds)`

**`SetTime`**: 0–24 style hours. **`SetDayLength`**: real-time seconds for a full cycle.

---

## `Sky.GetTimeHours(sky)` → float

Current simulated hour.

---

## `Sky.IsNight(sky)` → bool

**True** when the sun is below the horizon (implementation threshold).

---

## Common mistake

Calling **`Sky.Draw`** after opaque terrain — the sky should usually be **first** inside the camera block so depth works as expected.

---

## See also

- [CLOUD.md](CLOUD.md)
- [WEATHER.md](WEATHER.md)
