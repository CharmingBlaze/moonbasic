# Weather, fog, and wind (`WEATHER.*`, `FOG.*`, `WIND.*`)

Combined **weather state** (type, coverage), **distance fog** parameters stored for the module, and a **global wind** strength. **CGO** required for full behavior; fog may not call every Raylib entry point if a symbol is unavailable — state is still tracked.

---

## `WEATHER.*`

| Command | Role |
|--------|------|
| `Weather.Make()` | Creates a weather controller handle. |
| `Weather.Free(weather)` | Frees it. |
| `Weather.Update(weather, dt#)` | Per-frame update. |
| `Weather.Draw(weather)` | Draw pass (particles/effects as implemented). |
| `Weather.SetType(weather, name$)` | Sets named preset (e.g. `"clear"`, `"rain"` — see runtime). |
| `Weather.GetCoverage(weather)` | Normalized coverage **0–1**. |
| `Weather.GetType(weather)` | Current type string. |

---

## `FOG.*` (global/module state)

| Command | Role |
|--------|------|
| `Fog.Enable(enabled?)` | Turns fog application on/off where supported. |
| `Fog.SetNear(near#)` / `Fog.SetFar(far#)` | Distance fog start/end. |
| `Fog.SetColor(r, g, b, a)` | Fog color components **0–255**. |

**Common mistake:** Expecting **`FOG`** to duplicate **`Render`** fog APIs — this namespace is **weather-scoped** state; combine with your render pipeline as documented in runtime.

---

## `WIND.*`

| Command | Role |
|--------|------|
| `Wind.Set(strength#, dx#, dz#)` | Sets wind **strength** and a horizontal direction on the XZ plane (components need not be normalized). |
| `Wind.GetStrength()` | Reads current strength. |

---

## See also

- [PARTICLES.md](PARTICLES.md)
- [SKY.md](SKY.md)
