# Textures

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **LoadTexture (file$)** | **`LOADTEXTURE`**, **`TEXTURE.LOAD`** | **Heap handle** — **`FREETEXTURE`** / **`TEXTURE.FREE`**. |
| **CreateTexture (w, h)** | **`IMAGE.MAKE`** + **`TEXTURE.FROMIMAGE`** | **Image** + **texture** handles — free both when done. |
| **FreeTexture** | **`FREETEXTURE`**, **`TEXTURE.FREE`** | |
| **TextureWidth / Height** | **`TEXTUREWIDTH`**, **`TEXTUREHEIGHT`**, **`TEXTURE.WIDTH`**, **`HEIGHT`** | Read-only — no free. |
| **ScaleTexture / RotateTexture / TextureCoords** | **`TEXTURE.SETWRAP`**, **`SETFILTER`**, **`DRAW.TEXTUREPRO`** | |
| **EntityTexture** | **`ENTITY.TEXTURE`** | |
