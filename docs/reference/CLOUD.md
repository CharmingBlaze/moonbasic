# Clouds (`CLOUD.*`)

**Coverage** and timing state for volumetric-style clouds. The current implementation stores parameters and reserves a draw hook; visual detail may be minimal until shaders are extended. **CGO** required.

---

## `Cloud.Create()` → handle (canonical; deprecated `Cloud.Make()` / `CLOUD.MAKE`)

Creates a cloud layer object (registry **`CLOUD.CREATE`**).

---

## `Cloud.Free(cloud)`

Frees the cloud handle.

---

## `Cloud.Update(cloud, dt)` / `Cloud.Draw(cloud)`

Advance simulation time and draw (draw may be a no-op depending on build).

---

## `Cloud.SetCoverage(cloud, coverage)`

**`coverage`** in **0–1** range (clamped), affecting density/opacity where implemented.

---

## See also

- [SKY.md](SKY.md) — draw order
- [WEATHER.md](WEATHER.md) — precipitation coverage
