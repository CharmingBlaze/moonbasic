# Textures

| Designed | moonBASIC | Notes |
|----------|------------|-------|
| **LoadTexture(file)** | **`Texture.Load()`** | Returns a **texture handle**. |
| **CreateTexture(w, h)** | **`Texture.Load()`** | Or use **`Image.Make()`** + **`Texture.FromImage()`**. |
| **FreeTexture(id)** | **`Texture.Free()`** | Unloads GPU data. |
| **TextureWidth(id)** | **`Texture.Width()`** | |
| **TextureHeight(id)** | **`Texture.Height()`** | |
| **TextureFilter(id, mode)** | **`Texture.SetFilter()`** | |
| **EntityTexture(entity, id)** | **`Entity.Texture()`** | |
