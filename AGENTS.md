# Agent / assistant notes

- **Contributor workflow and commands:** [CONTRIBUTING.md](CONTRIBUTING.md), [docs/DEVELOPER.md](docs/DEVELOPER.md) (includes **VS Code / gopls** `fullruntime` + **`gopls_stub`** for terrain stubs on Windows and **`scripts/check_builds.sh`**).
- **Default `go run .` only compiles to `.mbc`.** Graphical programs need **`-tags fullruntime`** and **`moonrun`** or **`moonbasic --run`** (built with fullruntime).
- **Manifest changes:** edit `compiler/builtinmanifest/commands.json`, then run `go run . --check` on a relevant sample; regenerate `docs/API_CONSISTENCY.md` with `go run ./tools/apidoc` when the public API surface changes.
