# moonBASIC package format (v1)

Packages ship as a **directory** or **zip** containing compiled bytecode and metadata for the `moonbasic install` command.

## Layout

```
mylib/
  manifest.json    # required
  mylib.mbc        # compiled library (name from manifest.entry_mbc)
  *.mb             # optional: sources for INCLUDE (if you ship sources)
```

## manifest.json

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | yes | Package id (lowercase, `[a-z0-9_-]+`). |
| `version` | string | yes | Semantic version, e.g. `1.0.0`. |
| `description` | string | no | Short human-readable summary. |
| `moonbasic` | string | no | Engine version constraint (npm-style range), e.g. `>=1.2.0`. |
| `entry_mbc` | string | yes | Filename of the `.mbc` file inside the bundle. |
| `deps` | object | no | Map of package name → version range for future resolver. |
| `sha256_mbc` | string | no | Hex SHA-256 of `entry_mbc` for verification. |

## Version pinning

- Installed copies live under the user cache: `packages/<name>/<version>/`.
- Multiple versions can coexist; games pin by path or future lockfile.
- `moonbasic list` shows name, version, and install path.

## Registry index (remote `install <name>`)

When `MOONBASIC_REGISTRY` points to a JSON URL (default unset), `moonbasic install <name>` downloads that index and resolves `packages[name].url` to a zip or tarball.

Index shape:

```json
{
  "packages": {
    "math_extra": {
      "version": "1.0.0",
      "url": "https://github.com/org/moonbasic-pkgs/releases/download/v1/math_extra.zip",
      "sha256": "optional-hex-of-entire-zip"
    }
  }
}
```

## Publishing

`moonbasic publish <dir>` validates `manifest.json`, builds a zip (manifest + mbc + optional sources), and with `GITHUB_TOKEN` + `MOONBASIC_PUBLISH_REPO=owner/repo` uploads the zip as a release asset (requires an existing GitHub release tag or creates draft flow — see `moonbasic publish -help`).
