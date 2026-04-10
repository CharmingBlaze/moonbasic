# System Commands

Commands for interacting with the operating system and environment.

---

## Program Control

### `System.Exit()`

Terminates the program immediately. This is a hard exit; `Window.Close()` is preferred for a clean shutdown.

---

## Host Environment

### `System.GetEnv(varName$)`

Gets the value of an environment variable.

- `varName`: The name of the environment variable (e.g., "PATH").

### `System.OpenURL(url$)`

Opens a URL in the user's default web browser.

### `System.GetClipboard()` / `System.SetClipboard(text$)`

Gets or sets the content of the system clipboard.

---

## Command Line Arguments

These commands allow you to read arguments passed to your script when it was run from the terminal.

### `ARGC()`

Returns the number of command-line arguments. The script path itself is not counted.

### `COMMAND$(index)`

Returns the command-line argument at the specified index.

- `index`: The 0-based index of the argument.

---

## Full Example: Command Line Parser

Save this script as `args_test.mb` and run it from your terminal like this:
`moonbasic args_test.mb hello world --version`

```basic
PRINT "moonBASIC Argument Parser"
PRINT "-------------------------"

arg_count = ARGC()
PRINT "Arguments received: " + STR$(arg_count)

IF arg_count > 0 THEN
    FOR i = 0 TO arg_count - 1
        arg$ = COMMAND$(i)
        PRINT "Arg " + STR$(i) + ": " + arg$

        IF arg$ = "--version" THEN
            PRINT "Version flag detected!"
        ENDIF
    NEXT
ELSE
    PRINT "No arguments were provided."
ENDIF
```
