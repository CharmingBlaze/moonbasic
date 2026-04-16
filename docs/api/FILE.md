# File Commands

Commands for reading, writing, and managing files on disk. moonBASIC provides high-level file I/O for text and binary data, directory listing, file existence checks, and path utilities.

## Core Concepts

- **File paths** — Relative to the program's working directory, or absolute paths. Use `/` as the separator (works on all platforms).
- **Text I/O** — Read and write entire files as strings.
- **Binary I/O** — Read and write raw byte data.
- **No open/close** — moonBASIC uses single-call read/write functions. There is no file handle to manage.

---

## Reading

### `File.ReadText(filePath)`

Reads an entire text file into a string.

- `filePath` (string) — Path to the file.

**Returns:** `string`

```basic
contents = File.ReadText("config.txt")
PRINT contents
```

---

### `File.ReadLines(filePath)`

Reads a text file and returns an array of lines.

- `filePath` (string) — Path to the file.

**Returns:** Array of strings.

---

### `File.ReadBytes(filePath)`

Reads a file as raw binary data.

- `filePath` (string) — Path to the file.

**Returns:** `handle` — Binary data handle.

---

## Writing

### `File.WriteText(filePath, content)`

Writes a string to a file, overwriting any existing content.

- `filePath` (string) — Output path.
- `content` (string) — Text to write.

```basic
File.WriteText("save.txt", "Level: 5\nScore: 12345")
```

---

### `File.AppendText(filePath, content)`

Appends text to the end of an existing file (creates if missing).

- `filePath` (string) — Output path.
- `content` (string) — Text to append.

```basic
File.AppendText("log.txt", "Event: Player died at " + STR(Time.Millisecs()) + "\n")
```

---

### `File.WriteBytes(filePath, dataHandle)`

Writes binary data to a file.

- `filePath` (string) — Output path.
- `dataHandle` (handle) — Binary data handle.

---

## File System

### `File.Exists(filePath)`

Returns `TRUE` if the file exists.

- `filePath` (string) — Path to check.

**Returns:** `bool`

```basic
IF File.Exists("save.json") THEN
    LoadGame("save.json")
ELSE
    PRINT "No save file found"
ENDIF
```

---

### `File.Delete(filePath)`

Deletes a file.

- `filePath` (string) — File to delete.

---

### `File.Copy(sourcePath, destPath)`

Copies a file.

---

### `File.Rename(oldPath, newPath)`

Renames or moves a file.

---

### `File.Size(filePath)`

Returns the file size in bytes.

**Returns:** `int`

---

### `File.ListDir(directoryPath)`

Returns a list of files and subdirectories in a directory.

**Returns:** Array of file names.

---

### `File.IsDir(path)`

Returns `TRUE` if the path is a directory.

**Returns:** `bool`

---

### `File.MakeDir(path)`

Creates a directory (and parent directories if needed).

---

### `File.GetCwd()`

Returns the current working directory.

**Returns:** `string`

---

## Full Example

A simple logging system and save file manager.

```basic
; === Logging ===
FUNCTION LogEvent(msg)
    timestamp = STR(Time.Millisecs())
    File.AppendText("game.log", "[" + timestamp + "] " + msg + "\n")
END FUNCTION

; === Save System ===
FUNCTION SaveProgress(level, score)
    data = JSON.Create()
    JSON.SetInt(data, "level", level)
    JSON.SetInt(data, "score", score)
    JSON.SetFloat(data, "timestamp", Time.Millisecs())
    JSON.ToFilePretty(data, "progress.json")
    JSON.Free(data)
    LogEvent("Progress saved: level " + STR(level))
END FUNCTION

FUNCTION LoadProgress()
    IF NOT File.Exists("progress.json") THEN
        PRINT "No save file"
        RETURN
    ENDIF

    data = JSON.Parse("progress.json")
    level = JSON.GetInt(data, "level")
    score = JSON.GetInt(data, "score")
    JSON.Free(data)

    PRINT "Loaded: Level " + STR(level) + " Score " + STR(score)
    LogEvent("Progress loaded")
END FUNCTION

; Usage
SaveProgress(5, 12345)
LoadProgress()
```

---

## See Also

- [JSON](JSON.md) — Structured data I/O
- [NET](NET.md) — Network file transfers
