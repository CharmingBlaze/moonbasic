# String Commands

Commands for manipulating and querying strings.

## Core Concepts

-   **Type Inference**: Variables do not require suffixes. The compiler infers the type from the context or the first assignment.
-   **Concatenation**: Use the `+` operator to join strings: `"Hello " + "World"`.
-   **Conversion**: Use `STR(number)` to convert a number to a string, and `VAL(string)` to convert a string to a number.

---

## Slicing & Substrings

### `LEN(text)`

Returns the number of characters in a string.

### `LEFT(text, count)` / `RIGHT(text, count)`

Returns the leftmost or rightmost `count` characters of a string.

### `MID(text, start, [count])`

Extracts a substring. `start` is a 1-based index. If `count` is omitted, it returns the rest of the string.

---

## Formatting & Case

### `UPPER(text)` / `LOWER(text)`

Converts a string to uppercase or lowercase.

### `TRIM(text)`

Removes leading and trailing whitespace.

### `REPLACE(text, find, replace)`

Replaces all occurrences of `find` with `replace`.

---

## Templated interpolation (`INTERP` / `STRING.INTERP`)

`INTERP(template, value0, [value1, ...])` substitutes placeholders **`{0}`** … **`{9}`** in `template`, using up to **10** values. **`STRING.INTERP`** is the same function (namespace-style name).

Values are turned into text the same way as typical `PRINT` / `STR` display rules (numbers, booleans, handles show a readable form). This **builds a new string** — ideal for **HUD and debug**, not for **tight physics inner loops**; see [STRING_HEAP.md](STRING_HEAP.md).

```basic
msg = INTERP("Score {0} / {1}  hp {2}", score, maxScore, INT(hp))
label = STRING.INTERP("pos {0}, {1}", px, pz)
```

---

## Searching & Querying

### `INSTR(text, find, [start])`

Returns the 1-based index of the first occurrence of `find`. Returns `0` if not found.

### `CONTAINS(text, find)`

Returns `TRUE` if `text` contains `find`.

### `STARTSWITH(text, find)` / `ENDSWITH(text, find)`

Returns `TRUE` if `text` starts or ends with `find`.

---

## Splitting & Joining

### `SPLIT(text, separator)`

Splits a string into an array of strings using the `separator`. Returns a handle to the new array.

### `JOIN(arrayHandle, separator)`

Joins an array of strings into a single string, separated by `separator`.

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
