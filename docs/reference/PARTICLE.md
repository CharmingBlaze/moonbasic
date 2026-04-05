# Particles — 3D `PARTICLE.*` vs 2D `PARTICLE2D.*`

There is **no** separate **`PARTICLE3D.*`** namespace. **3D** billboard-style particles use **`PARTICLE.*`** (Raylib **`ParticleSystem`** wrappers in **`runtime/mbparticles`**). **2D** CPU circles use **`PARTICLE2D.*`** (see [SPRITE.md](SPRITE.md)).

---

## 3D — `PARTICLE.*`

Create with **`PARTICLE.MAKE`**, configure emitters (**`PARTICLE.SETEMITRATE`**, **`PARTICLE.SETLIFETIME`**, …), then each frame **`PARTICLE.UPDATE`** and **`PARTICLE.DRAW`** inside your **`CAMERA.BEGIN` / `CAMERA.END`** pass. Full API and examples: **[PARTICLES.md](PARTICLES.md)** and **`docs/API_CONSISTENCY.md`**.

Handle type name: **`Particle`**. Method dispatch uses **`PARTICLE.*`** registry keys.

---

## 2D — `PARTICLE2D.*`

Documented with sprites in **[SPRITE.md](SPRITE.md)** (`PARTICLE2D.MAKE`, draw helpers).

---

## See also

- [PARTICLES.md](PARTICLES.md) — full **`PARTICLE.*`** reference
- [SPRITE.md](SPRITE.md) — **`PARTICLE2D.*`**
- [MODEL.md](MODEL.md) — mesh-based scene effects
- [CAMERA.md](CAMERA.md) — 3D pass for **`PARTICLE.DRAW`**
