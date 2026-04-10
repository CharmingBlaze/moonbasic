# Array Commands

Commands for creating and manipulating arrays.

## Core Concepts

- **Declaration**: Arrays are declared with `DIM` or typed declaration syntax `name AS Type(...)`.
- **Indexing**: Arrays are **1-based**. `DIM a(10)` uses indices `1..10`.
- **Dimensions**: Arrays support any number of comma-separated dimensions (`a(10,10,5,2,...)`).
- **Storage**: Arrays are flat heap-backed storage internally (row-major), with runtime bounds checking.
- **Type hints**: `AS INTEGER` / `AS FLOAT` / `AS STRING` are stored as hints today; runtime remains dynamic for array slots.

---

## Declaration Syntax

### Classic `DIM`

```basic
DIM scores(10)
DIM grid(20, 15)
DIM names$(5)
```

### Typed declaration (preferred)

```basic
enemies AS INTEGER(100)
grid AS INTEGER(10, 10)
map AS INTEGER(10, 10, 5)
```

### Arrays of user `TYPE`

```basic
TYPE Enemy
    x AS FLOAT
    y AS FLOAT
    health AS INTEGER
ENDTYPE

enemies AS Enemy(100)
enemies(1).x = 32.0
enemies(1).health = 100
```

`DIM name AS TypeName(...)` also works for compatibility.

---

## Access and Safety

- Access uses one parenthesized index list: `arr(i)`, `grid(x, y)`, `map(x, y, z)`.
- Runtime enforces:
  - bounds checks per dimension,
  - dimension-count checks,
  - allocation-size limits,
  - stale-handle protection after free.

Out-of-bounds errors include array name, dimension, index, and valid range.

---

## Length

Use `.length` to get the first dimension size:

```basic
FOR i = 1 TO enemies.length
    PRINT enemies(i).health
NEXT i
```

For multidimensional arrays, `.length` returns dimension 1 size.

---

## Memory Management

- `ERASE(name)` frees a `DIM`/typed array and clears the variable.
- `ARRAYFREE(handle)` frees a heap array handle directly.
- `ERASE ALL` / `FREE.ALL` frees all heap objects and nulls handle globals/stack values.

See [MEMORY.md](../MEMORY.md).
