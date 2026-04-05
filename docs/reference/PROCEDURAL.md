# Noise and random (`PERLIN`, `FBMNOISE`, `RNDRANGE`, …)

Procedural noise and RNG helpers ship in **`runtime/mbgame`** (`register_ease_noise_rand.go` and related). **`RNDRANGE`** / **`RND`** use the module’s **`rand.Rand`** instance.

For **deterministic** runs, seed RNG through the documented **`RANDOMIZE`** / **`SEED`** surface if your game requires reproducibility (see manifest).
