# Bitwise Commands

Commands for low-level bitwise manipulation of integer values. These are useful for working with flags, compact data formats, or certain kinds of algorithms.

## Core Workflow

Bitwise operations treat integers as a sequence of binary digits (bits). A common use case is to store multiple boolean `TRUE`/`FALSE` states (flags) in a single integer variable.

For example, you can define flags for player abilities:

```basic
GLOBAL FLAG_JUMP = 1  ; Binary 0001
GLOBAL FLAG_SPRINT = 2 ; Binary 0010
```

---

### `BAND(a, b)` / `BOR` / `BXOR`
Performs bitwise logic on two integers.

- **Returns**: (Integer)

---

### `BNOT(v)`
Performs bitwise negation (ones' complement).

- **Returns**: (Integer)

---

### `BSHL(v, count)` / `BSHR`
Shifts bits left or right.

- **Returns**: (Integer)

---

### `BSET(v, index)` / `BCLEAR` / `BTOGGLE`
Returns a new integer with a specific bit (0–31) modified.

- **Returns**: (Integer)

---

### `BGET(v, index)`
Returns `TRUE` if the specified bit is set.

- **Returns**: (Boolean)

---

## Full Example

This example uses bitwise commands to manage a set of player ability flags.

```basic
; Define flags as powers of 2
CONST FLAG_JUMP = 1
CONST FLAG_SPRINT = 2
CONST FLAG_STEALTH = 4

; Start the player with jump and sprint abilities
player_flags = BOR(FLAG_JUMP, FLAG_SPRINT)

PRINT "Initial Flags: " + BIN(player_flags)

; Check if the player has the stealth ability
IF BAND(player_flags, FLAG_STEALTH) THEN
    PRINT "Player has stealth."
ELSE
    PRINT "Player does NOT have stealth."
ENDIF

PRINT "\nGranting stealth..."
player_flags = BOR(player_flags, FLAG_STEALTH)
PRINT "New Flags: " + BIN(player_flags)

; Check for stealth again
IF BAND(player_flags, FLAG_STEALTH) THEN
    PRINT "Player now has stealth!"
ENDIF

PRINT "\nRemoving sprint..."
player_flags = BXOR(player_flags, FLAG_SPRINT) ; Use XOR to toggle it off
PRINT "Final Flags: " + BIN(player_flags)
```
