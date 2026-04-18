# Event Commands

Pub/sub event system: fire named events with optional data, subscribe handlers, record and replay event sequences.

## Core Workflow

1. `EVENT.ON(eventName, handlerName)` — subscribe a handler function.
2. `EVENT.FIRE(eventName [, arg1, ...])` — broadcast the event.
3. `EVENT.OFF(eventName, handlerName)` — unsubscribe.
4. Use `EVENT.ONCE` for one-shot subscriptions.

---

## Subscribe & Unsubscribe

### `EVENT.ON(eventName, handlerName)` 

Subscribes `handlerName` (a function name string) to `eventName`. The handler is called each time the event fires.

---

### `EVENT.ONCE(eventName, handlerName)` 

Same as `EVENT.ON` but automatically unsubscribes after the first call.

---

### `EVENT.OFF(eventName, handlerName)` 

Removes a previously subscribed handler.

---

## Fire

### `EVENT.FIRE(eventName [, arg1, arg2, ...arg7])` 

Broadcasts `eventName` to all subscribers. Up to 7 optional arguments are passed to each handler.

---

## Event Handle Data

### `EVENT.TYPE(evHandle)` 

Returns the event type id from a received event handle.

---

### `EVENT.PEER(evHandle)` 

Returns the peer handle that sent this event (network events).

---

### `EVENT.DATA(evHandle)` 

Returns the data payload string from the event handle.

---

### `EVENT.CHANNEL(evHandle)` 

Returns the channel id of a network event.

---

### `EVENT.FREE(evHandle)` 

Frees an event handle.

---

## Recording & Replay

### `EVENT.RECSTART()` 

Begins recording all fired events into the active list.

---

### `EVENT.RECSTOP()` 

Stops recording.

---

### `EVENT.RECPLAYING()` / `EVENT.ISPLAYING()` 

Returns `TRUE` if a replay is currently active.

---

### `EVENT.REPLAY(listHandle)` 

Replays all events from `listHandle` in order.

---

## Event Lists

### `EVENT.LISTMAKE(name)` 

Creates a named event list handle.

---

### `EVENT.LISTLOAD(filename)` 

Loads an event list from a file.

---

### `EVENT.LISTEXPORT(listHandle, filename)` 

Saves an event list to a file.

---

### `EVENT.SETACTIVELIST(listHandle)` 

Sets which list recording goes into.

---

### `EVENT.LISTCLEAR(listHandle)` 

Clears all events from the list.

---

### `EVENT.LISTCOUNT(listHandle)` 

Returns the number of events in the list.

---

### `EVENT.LISTFREE(listHandle)` 

Frees the list handle.

---

## Full Example

Score system using events.

```basic
WINDOW.OPEN(800, 450, "Event Demo")
WINDOW.SETFPS(60)

score = 0

FUNCTION OnScore(amount)
    score = score + amount
    PRINT "Score: " + STR(score)
END FUNCTION

EVENT.ON("score", "OnScore")

WHILE NOT WINDOW.SHOULDCLOSE()
    IF INPUT.KEYPRESSED(KEY_SPACE) THEN
        EVENT.FIRE("score", 10)
    END IF
    IF INPUT.KEYPRESSED(KEY_X) THEN
        EVENT.FIRE("score", 100)
    END IF

    RENDER.CLEAR(20, 20, 40)
    DRAW.TEXT("Score: " + STR(score), 10, 10, 24, 255, 255, 255, 255)
    DRAW.TEXT("SPACE = +10  X = +100", 10, 45, 18, 180, 180, 180, 255)
    RENDER.FRAME()
WEND

EVENT.OFF("score", "OnScore")
WINDOW.CLOSE()
```

---

## See also

- [TRIGGER.md](TRIGGER.md) — zone triggers that fire events
- [NETWORK.md](NETWORK.md) — network events via `EVENT.PEER`
- [TIMER.md](TIMER.md) — timed event dispatch
