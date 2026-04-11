# System Commands

Commands for interacting with the operating system and environment.

---

## Program Control

### `System.Exit()`

Terminates the program immediately with the specified exit code. This is a hard exit; `Window.Close()` is preferred for a clean shutdown.

---

## Host Environment

### `System.Version()`

Returns the MoonBasic release label (e.g., `"1.0.0-GOLD"`). This value is useful for logging and debugging purposes.

### `System.OS()`

Returns a string identifying the operating system (`"windows"`, `"linux"`, `"darwin"`, etc.). This can be used to implement platform-specific behavior or optimizations.

### `System.Arch()`

Returns the CPU architecture (`"amd64"`, `"arm64"`, etc.). This can be used to implement architecture-specific optimizations or workarounds.

### `System.Env(name)`

Returns the value of an environment variable as a string, or `""` if not set. This allows scripts to interact with the environment and access external configuration.

### `System.OpenURL(url)`

Opens a URL in the user's default web browser. This can be used to provide links to documentation, tutorials, or other online resources.

### `System.ClipboardSetText(text)`

Sets the OS clipboard text. This can be used to provide a simple way to copy text to the clipboard.

### `System.ClipboardGetText()`

Returns the current OS clipboard text. This can be used to retrieve text from the clipboard and process it in the script.

---

## Command Line Arguments

These commands allow you to read arguments passed to your script when it was run from the terminal.

### `System.Args()`

Returns a string array handle containing command-line arguments. This allows scripts to parse and process command-line arguments.

### `System.IsCGO()`

Returns `TRUE` if the current build has CGO enabled (required for most 3D/GPU features). If `FALSE`, most 3D and physics commands will be unavailable or return errors.

---

## Full Example: Command Line Parser

Save this script as `args_test.mb` and run it from your terminal like this:
`moonbasic args_test.mb hello world --version`

```basic
PRINT "moonBASIC Argument Parser"
PRINT "-------------------------"

arg_count = ARGC()
PRINT "Arguments received: " + STR(arg_count)

IF arg_count > 0 THEN
    FOR i = 0 TO arg_count - 1
        arg = COMMAND(i)
        PRINT "Arg " + STR(i) + ": " + arg

        IF arg = "--version" THEN
            PRINT "Version flag detected!"
        ENDIF
    NEXT
ELSE
    PRINT "No arguments were provided."
ENDIF
```
