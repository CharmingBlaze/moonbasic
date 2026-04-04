# Bitwise Commands

Commands for low-level bitwise manipulation of integer values. These are useful for working with flags, compact data formats, or certain kinds of algorithms.

## Core Concepts

Bitwise operations treat integers as a sequence of binary digits (bits). A common use case is to store multiple boolean `TRUE`/`FALSE` states (flags) in a single integer variable.

For example, you can define flags for player abilities:

```basic
GLOBAL FLAG_JUMP = 1  ; Binary 0001
GLOBAL FLAG_SPRINT = 2 ; Binary 0010
GLOBAL FLAG_CROUCH = 4 ; Binary 0100
GLOBAL FLAG_INVIS = 8  ; Binary 1000
```

---

### `BAND(a, b)` / `BOR(a, b)` / `BXOR(a, b)`

-   `BAND` (AND): Returns bits that are set in *both* `a` and `b`.
-   `BOR` (OR): Returns bits that are set in *either* `a` or `b`.
-   `BXOR` (XOR): Returns bits that are set in one but not both.

---

### `BTEST(value, bit)`

Returns `TRUE` if the specified `bit` (0-indexed) is set (is 1) in `value`.

### `BSET(value, bit)`

Returns a new value with the specified `bit` set to 1.

### `BCLEAR(value, bit)`

Returns a new value with the specified `bit` cleared to 0.

### `BTOGGLE(value, bit)`

Returns a new value with the specified `bit` flipped (0 to 1, or 1 to 0).

---

## Full Example: Using Bitwise Flags

This example uses bitwise commands to manage a set of player ability flags.

```basic
; Define flags as powers of 2
CONST FLAG_JUMP = 1
CONST FLAG_SPRINT = 2
CONST FLAG_STEALTH = 4

; Start the player with jump and sprint abilities
player_flags = BOR(FLAG_JUMP, FLAG_SPRINT)

PRINT "Initial Flags: " + BIN$(player_flags)

; Check if the player has the stealth ability
IF BAND(player_flags, FLAG_STEALTH) THEN
    PRINT "Player has stealth."
ELSE
    PRINT "Player does NOT have stealth."
ENDIF

PRINT "\nGranting stealth..."
player_flags = BOR(player_flags, FLAG_STEALTH)
PRINT "New Flags: " + BIN$(player_flags)

; Check for stealth again
IF BAND(player_flags, FLAG_STEALTH) THEN
    PRINT "Player now has stealth!"
ENDIF

PRINT "\nRemoving sprint..."
player_flags = BXOR(player_flags, FLAG_SPRINT) ; Use XOR to toggle it off
PRINT "Final Flags: " + BIN$(player_flags)
```
