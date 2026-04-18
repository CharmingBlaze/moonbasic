# BTree Commands

Behaviour tree builder for AI logic. Compose sequences of conditions and actions into a tree, then `RUN` it each frame against an agent handle.

## Core Workflow

1. `BTREE.CREATE()` — create a root tree handle.
2. `BTREE.SEQUENCE(tree)` — add a sequence node.
3. `BTREE.ADDCONDITION(node, conditionName)` and `BTREE.ADDACTION(node, actionName)` — populate it.
4. Each frame: `BTREE.RUN(tree, agentHandle, dt)` — execute the tree.
5. `BTREE.FREE(tree)` when done.

---

## Creation

### `BTREE.CREATE()` 

Creates a new root behaviour tree handle.

---

## Tree Building

### `BTREE.SEQUENCE(tree)` 

Adds a sequence node to the tree. A sequence runs its children left-to-right and succeeds only if all children succeed. Returns the sequence node handle.

---

### `BTREE.ADDCONDITION(node, conditionName)` 

Adds a named condition leaf to `node`. The condition is evaluated by the agent's registered condition callback (defined in game code).

---

### `BTREE.ADDACTION(node, actionName)` 

Adds a named action leaf to `node`. The action is executed when reached and all preceding conditions pass.

---

## Execution

### `BTREE.RUN(tree, agentHandle, dt)` 

Runs the behaviour tree for one frame. `agentHandle` is passed to condition/action callbacks so they can read agent state.

---

## Lifetime

### `BTREE.FREE(tree)` 

Destroys the tree handle.

---

## Full Example

Simple patrol-and-attack tree for an enemy.

```basic
WINDOW.OPEN(960, 540, "BTree Demo")
WINDOW.SETFPS(60)

; build the behaviour tree
tree = BTREE.CREATE()
seq  = BTREE.SEQUENCE(tree)
BTREE.ADDCONDITION(seq, "playerInRange")
BTREE.ADDACTION(seq,    "attackPlayer")

; fallback patrol
patrol = BTREE.SEQUENCE(tree)
BTREE.ADDACTION(patrol, "patrolWaypoints")

enemy = NAVAGENT.CREATE(0)
NAVAGENT.SETPOS(enemy, 5, 0, 0)
NAVAGENT.SETSPEED(enemy, 3.0)

px = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN px = px + 4 * dt
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN px = px - 4 * dt

    ; BTree callbacks would check dist(enemy, player) and move accordingly
    BTREE.RUN(tree, enemy, dt)
    NAVAGENT.UPDATE(enemy, dt)

    RENDER.CLEAR(20, 25, 35)
    DRAW.TEXT("Enemy X: " + STR(INT(NAVAGENT.X(enemy))), 10, 10, 18, 200, 200, 200, 255)
    DRAW.TEXT("Player X: " + STR(INT(px)), 10, 35, 18, 80, 200, 80, 255)
    RENDER.FRAME()
WEND

BTREE.FREE(tree)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `BTREE.MAKE(...)` | Deprecated alias of `BTREE.CREATE`. |

---

## See also

- [NAVAGENT.md](NAVAGENT.md) — navigation agents
- [STEER.md](STEER.md) — steering forces
- [ANIM.md](ANIM.md) — driving animations from AI states
