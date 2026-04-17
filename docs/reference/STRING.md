# String Commands

Commands for manipulating and querying strings.

## Core Workflow

-   **Type Inference**: Variables do not require suffixes. moonBASIC does **not** use Blitz-style **`#` / `$` / `?` / `%`** on names; infer types from assignment or use `DIM` / `AS` ([STYLE_GUIDE.md](../../STYLE_GUIDE.md)).
-   **Concatenation**: Use the `+` operator to join strings: `"Hello " + "World"`.
-   **Conversion**: Use **`STR(value)`** to convert a value to a string, and **`FORMAT(value, pattern)`** for printf-style formatting (manifest canonical). Legacy **`STR$`** / **`FORMAT$`** are deprecated aliases (same runtime).
-   **Slice / search / binary helpers**: Prefer **`LEFT`**, **`RIGHT`**, **`MID`**, **`TRIM`**, **`SPLIT`**, **`JOIN`**, **`HEX`**, **`BIN`**, **`OCT`**, **`MKINT`**, … — each has a legacy **`…$`** alias in the manifest. See [API_CONSISTENCY.md](../API_CONSISTENCY.md).

---

## Slicing & Substrings

### `LEN(s)` 
Returns the number of characters in a string.

---

### `UPPER(s)` / `LOWER(s)` 
Converts a string to all uppercase or all lowercase letters.

---

### `TRIM(s)` 
Removes whitespace from both ends of a string.

---

### `MID(s, start, count)` 
Extracts a substring from `s`, starting at the 1-based index `start`. If `count` is omitted, returns the rest of the string.

---

### `REPLACE(s, old, new)` 
Returns a new string with all occurrences of `old` replaced by `new`.

---

### `SPLIT(s, separator)` 
Splits `s` into an array of substrings. Returns a handle to a string list.

---

### `INSTR(s, find)` 
Returns the 1-based position of `find` within `s`, or `0` if not found.

---

### `JOIN(handle, separator)` 
Joins elements of a string list into a single string, separated by `separator`.

---

## Full Example

This example demonstrates how to parse a comma-separated string, process the parts, and display them.

```basic
; A string containing player data
data = "player_1,100,55.5"

; Split the string into an array
parts = SPLIT(data, ",")

; Extract and convert the data
name = parts(0)
score = VAL(parts(1))
health = VAL(parts(2))

; Modify and display the data
PRINT "Player: " + UPPER(name)
PRINT "Score: " + STR(score)
PRINT "Health: " + STR(health) + "%"

; Remember to free the array handle created by SPLIT
ARRAYFREE(parts)
```
