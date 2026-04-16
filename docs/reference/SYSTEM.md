# System (`SYSTEM.*`, `ARGC`, `COMMAND`, `SystemProperty`)

Commands for interacting with the operating system and environment.

**Conventions:** [STYLE_GUIDE.md](../../STYLE_GUIDE.md), [API_CONVENTIONS.md](API_CONVENTIONS.md) — reference pages use uppercase **`NAMESPACE.ACTION`** where applicable; Easy Mode (`System.Exit`, …) is [compatibility only](../../STYLE_GUIDE.md#easy-mode-compatibility-layer).

**Page shape:** [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) — see [WAVE.md](WAVE.md) (registry-first headings, **Full Example** at the end).

---

## Program control

### `SYSTEM.EXIT([code])`

Terminates the program immediately with the given exit code (when supported). Prefer **`WINDOW.CLOSE()`** for a normal graphics shutdown when a window is open.

---

## Host environment

### `SYSTEM.VERSION()`

Returns the MoonBasic release label (e.g. `"1.0.0-GOLD"`). Useful for logging and debugging.

### `SystemProperty(key)`

Returns a small set of OS/runtime facts keyed by **`key`** (string). Examples: **`"os"`** / **`"os_name"`** → OS id (`"windows"`, `"linux"`, …); **`"arch"`** → CPU architecture (`"amd64"`, `"arm64"`, …); **`"cpu_cores"`**, **`"compiler"`**, etc. Unknown keys return **`""`**.

### `SYSTEM.CPUNAME()` / `SYSTEM.GPUNAME()`

CPU model string and primary GPU name (best-effort; may be **`"(unavailable)"`** on some setups).

### `SYSTEM.TOTALMEMORY()` / `SYSTEM.FREEMEMORY()`

Host RAM totals (bytes), via the same probes as the runtime host module.

### `SYSTEM.GETENV(name)` / `SYSTEM.SETENV(name, value)`

Read or set an environment variable. Aliases: **`ENVIRON`**, **`ENVIRON$`** (same arity as **`SYSTEM.GETENV`**).

### `SYSTEM.OPENURL(url)`

Opens a URL in the system default browser.

### `SYSTEM.GETCLIPBOARD()` / `SYSTEM.SETCLIPBOARD(text)`

Read or write the OS text clipboard.

### `SYSTEM.LOCALE()` / `SYSTEM.USERNAME()`

Current locale hint and current username (best-effort).

### `SYSTEM.ISDEBUGBUILD()`

**`TRUE`** when the build or environment indicates a development/debug configuration (see runtime for exact rules).

### `SYSTEM.EXECUTE(cmdline)`

Runs a shell command (`cmd /C` on Windows, `sh -c` elsewhere). Returns an exit code on success; errors may surface as runtime errors.

---

## Command-line arguments

### `ARGC()`

Returns the number of host arguments available to the script (see **`COMMAND`**).

### `COMMAND([index])`

With **no** arguments, returns the full command line as one string. With **`index`**, returns the **`index`**-th argument (0-based), or **`""`** if out of range.

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
