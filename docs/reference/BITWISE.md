# Bitwise Commands

Commands for low-level bitwise manipulation of integer values. These are useful for working with flags, compact data formats, or certain kinds of algorithms.

## Core Concepts

Bitwise operations treat integers as a sequence of binary digits (bits). A common use case is to store multiple boolean `TRUE`/`FALSE` states (flags) in a single integer variable.

For example, you can define flags for player abilities:

```basic
GLOBAL FLAG_JUMP = 1  ; Binary 0001
GLOBAL FLAG_SPRINT = 2 ; Binary 0010
```

---

### `Bit.And(a, b)`
Returns the bitwise **AND** of two 32-bit integers.

### `Bit.Or(a, b)`
Returns the bitwise **OR** of two 32-bit integers.

### `Bit.Xor(a, b)`
Returns the bitwise **XOR** of two 32-bit integers.

### `Bit.Not(a)`
Returns the bitwise **NOT** (one's complement) of an integer.

---

### `Bit.Shl(v, n)`
Returns `v` shifted left by `n` bits.

### `Bit.Shr(v, n)`
Returns `v` shifted right by `n` bits.

---

### `Bit.Get(v, bitIndex)`
Returns `TRUE` if the bit at `bitIndex` (0–31) is set, `FALSE` otherwise.

### `Bit.Set(v, bitIndex)`
Returns a new integer with the bit at `bitIndex` set to 1.

### `Bit.Clear(v, bitIndex)`
Returns a new integer with the bit at `bitIndex` set to 0.

### `Bit.Toggle(v, bitIndex)`
Returns a new value with the specified bit flipped (0 to 1, or 1 to 0).

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
