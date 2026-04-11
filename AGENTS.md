# Agent / assistant notes

- **Contributor workflow and commands:** [CONTRIBUTING.md](CONTRIBUTING.md), [docs/DEVELOPER.md](docs/DEVELOPER.md) (includes **VS Code / gopls** `fullruntime` + **`gopls_stub`** for terrain stubs on Windows and **`scripts/check_builds.sh`**).
- **Default `go run .` only compiles to `.mbc`.** Graphical programs need **`-tags fullruntime`** and **`moonrun`** or **`moonbasic --run`** (built with fullruntime).
- **Manifest changes:** edit `compiler/builtinmanifest/commands.json`, then run `go run . --check` on a relevant sample; regenerate `docs/API_CONSISTENCY.md` with `go run ./tools/apidoc` when the public API surface changes.
- **Windows and Linux:** Register new builtins on every split build path (`*_cgo.go` / `*_stub.go` or equivalent) so both OSes expose the same manifest keys—use stubs or no-ops where native code is platform-specific (e.g. Jolt on Linux+CGO vs physics stubs on Windows). Run `bash scripts/check_builds.sh` on Linux/macOS; on Windows use the same `go build` tag matrix from [docs/DEVELOPER.md](docs/DEVELOPER.md).
