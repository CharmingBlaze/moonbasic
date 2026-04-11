# String Commands

Commands for manipulating and querying strings.

## Core Concepts

-   **Type Inference**: Variables do not require suffixes. The compiler infers the type from the context or the first assignment.
-   **Concatenation**: Use the `+` operator to join strings: `"Hello " + "World"`.
-   **Conversion**: Use `STR(number)` to convert a number to a string, and `VAL(string)` to convert a string to a number.

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
