# Array Commands

Commands for creating and manipulating arrays.

## Core Concepts

-   **Declaration**: Arrays are created with the `DIM` keyword. The number in parentheses is the *size*, not the upper bound. `DIM my_array(10)` creates an array with indices 0 through 9.
-   **Handles**: When you pass an array to a command (like `ARRAYLEN`), you are passing its handle, not its contents.
-   **Types**: An array holds values of the same type, determined by the variable's suffix (e.g., `my_strings$(10)`).

---

### `DIM`

Declares a new array. This is a language keyword, not a command.

```basic
; Declare an array of 10 integers (indices 0-9)
DIM scores(10)

; Declare a 2D array for a tilemap
DIM map(20, 15)

; Declare an array of strings
DIM names$(5)
```

---

### `ARRAYLEN(arrayHandle)`

Returns the total number of elements in an array.

```basic
DIM my_array(20)
PRINT ARRAYLEN(my_array) ; Outputs 20
```

---

### `ARRAYFREE(arrayHandle)`

Frees the memory used by an array. This is especially important for arrays returned by commands like `SPLIT$`.

---

## Full Example: Populating and Reading an Array

```basic
; Create an array to hold 5 high scores
DIM high_scores(5)

; Populate the array using a loop
FOR i = 0 TO ARRAYLEN(high_scores) - 1
    high_scores(i) = (5 - i) * 1000
NEXT

; Print the contents of the array
PRINT "High Scores:"
FOR i = 0 TO ARRAYLEN(high_scores) - 1
    PRINT STR$(i+1) + ". " + STR$(high_scores(i))
NEXT
```

---

## Other Commands

- `REDIM`: **[PARTIAL]** Coming soon.
- `ERASE`: **[PARTIAL]** Coming soon.
- `ARRAYFILL`: **[PARTIAL]** Coming soon.
- `ARRAYCOPY`: **[PARTIAL]** Coming soon.
- `ARRAYSORT`: **[PARTIAL]** Coming soon.
- `ARRAYPUSH`: **[PARTIAL]** Coming soon.
- `ARRAYPOP`: **[PARTIAL]** Coming soon.
