# Benchmarks

Run from the repo root:

```bash
moonbasic --benchmark benchmarks/alu.mb
```

Print a line containing **`MOONBENCH`** (e.g. `PRINT("MOONBENCH ops=" + STR$(n))`) so the harness can attach script-reported metrics. The driver also prints wall-clock time and `runtime.MemStats` delta on stderr.

| Script | Intent |
|--------|--------|
| `alu.mb` | Integer / float loop workload |
| `strings.mb` | String concatenation |
| `drawcalls.mb` | 2D rectangles per frame (needs display / Xvfb) |
| `models.mb` | Placeholder until 3D model draw is benchmarked |
| `physics.mb` | Math-heavy stand-in until physics is isolated in bench |
