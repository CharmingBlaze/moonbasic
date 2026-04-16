# String Commands

Commands for manipulating and querying strings.

## Core Concepts

-   **Type Inference**: Variables do not require suffixes. moonBASIC does **not** use Blitz-style **`#` / `$` / `?` / `%`** on names; infer types from assignment or use `DIM` / `AS` ([STYLE_GUIDE.md](../../STYLE_GUIDE.md)).
-   **Concatenation**: Use the `+` operator to join strings: `"Hello " + "World"`.
-   **Conversion**: Use **`STR(value)`** to convert a value to a string, and **`FORMAT(value, pattern)`** for printf-style formatting (manifest canonical). Legacy **`STR$`** / **`FORMAT$`** are deprecated aliases (same runtime).
-   **Slice / search / binary helpers**: Prefer **`LEFT`**, **`RIGHT`**, **`MID`**, **`TRIM`**, **`SPLIT`**, **`JOIN`**, **`HEX`**, **`BIN`**, **`OCT`**, **`MKINT`**, … — each has a legacy **`…$`** alias in the manifest. See [API_CONSISTENCY.md](../API_CONSISTENCY.md).

---

## Slicing & Substrings

### `String.Len(s)`
Returns the number of characters in a string.

### `String.Upper(s)` / `String.Lower(s)`
Converts a string to all uppercase or all lowercase letters.

### `String.Trim(s)`
Removes whitespace from both ends of a string.

---

### `String.Mid(s, start [, count])`
Extracts a substring from `s`, starting at the 1-based index `start`. If `count` is omitted, returns the rest of the string.

### `String.Replace(s, old, new)`
Returns a new string with all occurrences of the `old` substring replaced by `new`.

### `String.Split(s, separator)`
Splits `s` into an array of substrings based on `separator`. Returns a handle to a string list.

---

### `String.Contains(s, find)`
Returns `TRUE` if the substring `find` is found within `s`.

### `String.Join(handle, separator)`
Joins the elements of a string list array into a single string, separated by `separator`.

---

## Full Example: Parsing Data

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
