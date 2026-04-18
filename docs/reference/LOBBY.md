# Lobby Commands

Game lobby discovery and session management: create named lobbies, advertise them, find others, and join.

## Core Workflow

1. **Host**: `LOBBY.CREATE(name, maxPlayers)` → configure → `LOBBY.START(lobby)`.
2. **Client**: `LOBBY.FIND(gameName, filter)` → pick a result → `LOBBY.JOIN(lobby)`.
3. Use `LOBBY.SETPROPERTY` to attach metadata (map name, game mode, etc.).
4. `LOBBY.FREE(lobby)` when done.

---

## Creation

### `LOBBY.CREATE(name, maxPlayers)` 

Creates a new lobby descriptor with the given display `name` and `maxPlayers` limit. Returns a **lobby handle**.

---

## Configuration

### `LOBBY.SETPROPERTY(lobby, key, value)` 

Sets a string metadata property on the lobby. Use to advertise map, mode, version, etc.

---

### `LOBBY.SETHOST(lobby, address, port)` 

Sets the connection host and port players will connect to when joining.

---

## Start

### `LOBBY.START(lobby)` 

Advertises the lobby so other clients can discover it via `LOBBY.FIND`.

---

## Discovery

### `LOBBY.FIND(gameName, filter)` 

Searches for lobbies matching `gameName`. `filter` is an optional property filter string. Returns a **lobby list handle** whose entries are accessed by iterating results.

---

### `LOBBY.GETNAME(lobby)` 

Returns the display name of a lobby (from a find result).

---

## Join

### `LOBBY.JOIN(lobby)` 

Connects to the lobby's advertised host. Triggers `CLIENT.ONCONNECT` on success.

---

## Lifetime

### `LOBBY.FREE(lobby)` 

Frees the lobby handle.

---

## Full Example

Host advertising a lobby and a client finding and joining it.

```basic
; === HOST ===
lobby = LOBBY.CREATE("My Game", 8)
LOBBY.SETPROPERTY(lobby, "map", "level1")
LOBBY.SETPROPERTY(lobby, "mode", "deathmatch")
LOBBY.SETHOST(lobby, "192.168.1.100", 7777)
LOBBY.START(lobby)
SERVER.START(7777, 8)

WHILE NOT WINDOW.SHOULDCLOSE()
    SERVER.TICK(TIME.DELTA())
    RENDER.FRAME()
WEND

LOBBY.FREE(lobby)
SERVER.STOP()

; === CLIENT ===
results = LOBBY.FIND("My Game", "")
IF results THEN
    LOBBY.JOIN(results)
END IF
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `LOBBY.MAKE(...)` | Deprecated alias of `LOBBY.CREATE`. |

---

## See also

- [NET.md](NET.md) — server/client connection
- [NETWORK.md](NETWORK.md) — network overview and lobby pattern
