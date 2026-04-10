# Files

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **ReadFile (file)** | **Flat:** **`READFILE`** → **`FILE.READALLTEXT`**; also **`READALLTEXT`**, **`FILE.READALLTEXT`** | Whole file as string — no path handle. |
| **WriteFile (path, text)** | **Flat:** **`WRITEFILE`** → **`FILE.WRITEALLTEXT`**; **`WRITEALLTEXT`**, **`FILE.WRITEALLTEXT`** | **Note:** flat **`WRITEFILE`** is whole-file write. Legacy stream **`WRITEFILE`** was **`FILE.WRITE`** — use **`FILE.WRITE`** for open-handle streaming. |
| **CloseFile (handle)** | **`CLOSEFILE`**, **`FILE.CLOSE`** | Close handles from **`OPENFILE`**. |
| **ReadLine / WriteLine** | **Flat:** **`READLINE`**, **`WRITELINE`** → **`FILE.READLINE`**, **`FILE.WRITELN`**; also **`READFILE`**, **`WRITEFILELN`** | |
| **FileExists** | **Flat:** **`FILEEXISTS`** → **`FILE.EXISTS`** | |
| **DeleteFile / CopyFile** | **Flat:** **`DELETEFILE`** → **`UTIL.DELETEFILE`**; **`COPYFILE`** → **`UTIL.COPYFILE`** | Same as existing globals (re-registered by facade). |
