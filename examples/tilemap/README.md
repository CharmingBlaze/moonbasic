# Tilemap example

Loads a small [Tiled](https://www.mapeditor.org/) map (`assets/level1.tmx`), draws it with `TILEMAP.DRAW`, and uses the `collision` layer for movement blocking via `TILEMAP.ISSOLID`.

```bash
CGO_ENABLED=1 go run . examples/tilemap/main.mb
```

Assets live beside the `.tmx` file (`tiles.png` is a 2-tile 16×16 strip). Edit the map in Tiled and re-run — keep CSV encoding and a layer named `collision` for solid tiles.
