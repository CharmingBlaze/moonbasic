# String Commands

Commands for manipulating and querying strings.

## Core Concepts

-   **Suffix**: String variables must end with the `$` suffix (e.g., `my_string$`).
-   **Concatenation**: Use the `+` operator to join strings: `"Hello " + "World"`.
-   **Conversion**: Use `STR$(number)` to convert a number to a string, and `VAL(string$)` to convert a string to a number.

---

## Slicing & Substrings

### `LEN(text$)`

Returns the number of characters in a string.

### `LEFT$(text$, count)` / `RIGHT$(text$, count)`

Returns the leftmost or rightmost `count` characters of a string.

### `MID$(text$, start, [count])`

Extracts a substring. `start` is a 1-based index. If `count` is omitted, it returns the rest of the string.

---

## Formatting & Case

### `UPPER$(text$)` / `LOWER$(text$)`

Converts a string to uppercase or lowercase.

### `TRIM$(text$)`

Removes leading and trailing whitespace.

### `REPLACE$(text$, find$, replace$)`

Replaces all occurrences of `find$` with `replace$`.

---

## Searching & Querying

### `INSTR(text$, find$, [start])`

Returns the 1-based index of the first occurrence of `find$`. Returns `0` if not found.

### `CONTAINS(text$, find$)`

Returns `TRUE` if `text$` contains `find$`.

### `STARTSWITH(text$, find$)` / `ENDSWITH(text$, find$)`

Returns `TRUE` if `text$` starts or ends with `find$`.

---

## Splitting & Joining

### `SPLIT$(text$, separator$)`

Splits a string into an array of strings using the `separator$`. Returns a handle to the new array.

### `JOIN$(arrayHandle, separator$)`

Joins an array of strings into a single string, separated by `separator$`.

---

## Full Example: Parsing Data

This example demonstrates how to parse a comma-separated string, process the parts, and display them.

```basic
; A string containing player data
data$ = "player_1,100,55.5"

; Split the string into an array
parts = SPLIT$(data$, ",")

; Extract and convert the data
name$ = parts(0)
score = VAL(parts(1))
health# = VAL(parts(2))

; Modify and display the data
PRINT "Player: " + UPPER$(name$)
PRINT "Score: " + STR$(score)
PRINT "Health: " + STR$(health#) + "%"

; Remember to free the array handle created by SPLIT$
ARRAYFREE(parts)
```
