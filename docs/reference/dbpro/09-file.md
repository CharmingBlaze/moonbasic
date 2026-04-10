# DBPro — File I/O

moonBASIC: **`FILE.*`** and flat aliases **`OPENFILE`**, **`CLOSEFILE`**, … — [FILE.md](../FILE.md).

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **OPEN TO READ (file, handle)** | ✓ **`FILE.OPEN`**, **`OPENFILE`** with **`"r"`** | Returns handle. |
| **OPEN TO WRITE** | ✓ **`FILE.OPEN`** with **`"w"`** / **`"a"`** | |
| **CLOSE FILE** | ✓ **`FILE.CLOSE`**, **`CLOSEFILE`** | |
| **READ BYTE** / **WRITE BYTE** | ≈ **`MEM.*`** / raw read patterns | moonBASIC file API is primarily **text/stream**; binary workflows use **[MEM.md](../MEM.md)** or **[DATA.md](../DATA.md)** as needed — check [API_CONSISTENCY.md](../../API_CONSISTENCY.md). |
| **READ STRING** / **WRITE STRING** | ✓ **`FILE.READLINE`**, **`FILE.WRITE`**, **`WRITELN`** | |
| **FILE EXIST** | ✓ **`FILEEXISTS`**, **`FILE.EXISTS`** | |
| **DELETE FILE** | ✓ **`DELETEFILE`** | |
| **COPY FILE** | ✓ **`COPYFILE`** | |
