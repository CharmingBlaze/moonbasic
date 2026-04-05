# Particles — 3D `PARTICLE.*` vs 2D `PARTICLE2D.*`

## 3D — `PARTICLE.*`

GPU-style particle systems (Raylib **`ParticleSystem`** wrappers) live under **`PARTICLE.*`** in **`runtime/mbparticles`**: create with **`PARTICLE.MAKE`**, configure emitters, then **`PARTICLE.UPDATE`** / **`PARTICLE.DRAW`** each frame. See **[PARTICLES.md](PARTICLES.md)** and **`docs/API_CONSISTENCY.md`** for the full list.

## 2D — `PARTICLE2D.*`

Simple CPU circles for prototyping are documented in **[SPRITE.md](SPRITE.md)** (`PARTICLE2D.*` alongside **`SPRITE.*`**).

---

## See also

- [SPRITE.md](SPRITE.md) — `PARTICLE2D.*`
- [MODEL.md](MODEL.md) — scene-scale effects with meshes
