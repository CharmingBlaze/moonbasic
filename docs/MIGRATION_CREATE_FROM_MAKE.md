# Migration Guide: `MAKE` to `CREATE`

MoonBASIC now standardizes creation commands on `CREATE`.

## What Changed

- Canonical creation commands are now `NAMESPACE.CREATE` and `NAMESPACE.CREATE<Type>`.
- Legacy `NAMESPACE.MAKE` and `NAMESPACE.MAKE<Type>` names remain available as deprecated aliases during migration.
- Canonical position setter is `SETPOS`; `SETPOSITION` remains an alias for compatibility.

## Quick Rename Rules

- `*.MAKE` -> `*.CREATE`
- `*.MAKE<Type>` -> `*.CREATE<Type>`
- Keep all other arguments and behavior the same.

## Common Examples

```basic
' Before
cam = CAMERA.MAKE()
light = LIGHT.MAKEPOINT()
model = MODEL.MAKE(mesh)

' After
cam = CAMERA.CREATE()
light = LIGHT.CREATEPOINT()
model = MODEL.CREATE(mesh)
```

```basic
' Before
MODEL.SETPOSITION(m, x, y, z)

' After
MODEL.SETPOS(m, x, y, z)
```

## Recommended Migration Workflow

1. Replace `MAKE` names with `CREATE` names.
2. Replace `SETPOSITION` calls with `SETPOS`.
3. Run `go run . --check <your_script.mb>` on updated scripts.
4. Keep legacy aliases only for temporary compatibility branches.

## Compatibility Window

- `MAKE` aliases still resolve in the current transition period.
- New code and all documentation should use `CREATE`.
- Future major versions may remove deprecated `MAKE` aliases.

