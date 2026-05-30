# System Commands

OS environment, program control, command-line arguments, and system info.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Query system properties with `SYSTEM.VERSION`, `SystemProperty`, `SYSTEM.CPUNAME`, etc. Read command-line args with `ARGC` / `COMMAND`. Exit with `SYSTEM.EXIT`.

---

## Program control

### `SYSTEM.EXIT([code])`
Terminates the program immediately.

- **Arguments**:
    - `code`: (Integer, Optional) Process exit code.
- **Returns**: (None)

---

### `SYSTEM.VERSION()`
Returns the MoonBASIC release label.

- **Returns**: (String)

---

### `SYSTEM.CPUNAME()` / `GPUNAME()`
Returns the hardware model strings.

- **Returns**: (String)

---

### `SYSTEM.GETENV(name)` / `SETENV(name, value)`
Reads or writes environment variables.

- **Returns**: (String) For `GETENV`.

---

### `SYSTEM.OPENURL(url)`
Opens a URL in the system's default browser.

---

### `SYSTEM.GETCLIPBOARD()` / `SETCLIPBOARD(text)`
Accesses the OS text clipboard.

---

### `ARGC()`
Returns the number of command-line arguments.

- **Returns**: (Integer)

---

### `COMMAND([index])`
Returns the full command line or a specific argument.

- **Arguments**:
    - `index`: (Integer, Optional) 0-based argument index.
- **Returns**: (String)

---

## Full Example: command-line parser

Save as `args_test.mb` and run, e.g. `moonbasic args_test.mb hello world --version`:

```basic
PRINT("moonBASIC argument parser")
PRINT("-------------------------")

arg_count = ARGC()
PRINT("Arguments received: " + STR(arg_count))

IF arg_count > 0 THEN
    FOR i = 0 TO arg_count - 1
        arg = COMMAND(i)
        PRINT("Arg " + STR(i) + ": " + arg)

        IF arg = "--version" THEN
            PRINT("Version flag detected!")
        ENDIF
    NEXT
ELSE
    PRINT("No arguments were provided.")
ENDIF

PRINT("OS: " + SystemProperty("os") + "  arch: " + SystemProperty("arch"))
```

---

## Extended Command Reference

### Process & environment

| Command | Description |
|--------|-------------|
| `SYSTEM.EXECUTE(cmd)` | Run a shell command string; returns exit code. |
| `SYSTEM.USERNAME()` | Returns the OS user name string. |
| `SYSTEM.LOCALE()` | Returns the OS locale string (e.g. `"en-AU"`). |
| `SYSTEM.MONITOR(id)` | Returns monitor info array `[w, h, refreshRate]` for monitor `id`. |
| `SYSTEM.ISDEBUGBUILD()` | Returns `TRUE` if running a debug build. |

### Memory

| Command | Description |
|--------|-------------|
| `SYSTEM.TOTALMEMORY()` | Returns total system RAM in bytes. |
| `SYSTEM.FREEMEMORY()` | Returns available system RAM in bytes. |

## See also

- [DEBUG.md](DEBUG.md) — `DEBUG.HEAPSTATS`, `DEBUG.GCSTATS`
- [UTIL.md](UTIL.md) — file system helpers
