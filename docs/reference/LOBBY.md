# Lobby Commands

Game lobby discovery and session management: create named lobbies, attach metadata, and (within a single runtime) scan for matching sessions.

## Scope (read this first)

- **`LOBBY.FIND` is in-process only.** It scans lobbies registered in the **same moonBASIC runtime / process** (`lobbyHandles` in the implementation). It does **not** perform Internet-wide or Steam-style discovery, and there is no built-in dedicated lobby server.
- **Arguments are property key/value pairs**, not a fuzzy “game name” search: `LOBBY.FIND(key$, value$)` lowercases the key and returns handles whose `LOBBY.SETPROPERTY` map contains that exact pair.
- **Cross-machine play** still needs a real address: set `LOBBY.SETHOST` with the host players can reach, or skip lobbies entirely and call `CLIENT.CONNECT` / `NET.CONNECT` with an IP and port you obtained out-of-band (chat, web API, your own matchmaking service).

## Core Workflow

1. **Host**: `LOBBY.CREATE(name, maxPlayers)` → `LOBBY.SETPROPERTY` / `LOBBY.SETHOST` → `LOBBY.START(lobby)`.
2. **Same runtime (optional)**: `LOBBY.FIND(key$, value$)` → inspect the returned **array handle** (possibly empty) → pick a lobby handle → `LOBBY.JOIN(lobbyHandle)` (requires a real `LOBBY.SETHOST`).
3. **Remote client without shared memory**: use **`CLIENT.CONNECT(host, port)`** or **`NET.CONNECT`** with the host/port you already know—**do not** expect `LOBBY.FIND` on machine B to discover a lobby created only on machine A.
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

### `LOBBY.FIND(key$, value$)` 

Returns a **heap-backed numeric array** of matching **lobby handles**. The implementation compares `key` (lower-cased) against `LOBBY.SETPROPERTY` keys and requires an exact string `value` match. When nothing matches, you still get an array (possibly length `1` with a `0` sentinel—treat **length / contents** as the source of truth in your build). This is **not** a cloud lobby browser.

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

Host creates a lobby, advertises connection info, and starts the high-level server. A **second machine** cannot magically discover this via `LOBBY.FIND`; it must already know the host/port (here `192.168.1.100:7777`) and call `CLIENT.CONNECT` / `NET.CONNECT`. The `LOBBY.FIND` block shows **in-process** filtering only (for example diagnostics or tooling running inside the same executable).

```basic
; === HOST ===
lobby = LOBBY.CREATE("My Game", 8)
LOBBY.SETPROPERTY(lobby, "map", "level1")
LOBBY.SETPROPERTY(lobby, "mode", "deathmatch")
LOBBY.SETHOST(lobby, "192.168.1.100", 7777)
LOBBY.START(lobby)
SERVER.START(7777, 8)

; Same-process sanity check: FIND matches the property bag above.
matches = LOBBY.FIND("map", "level1")
; Use ARRAY.LEN(matches) / iterate per ARRAY.md if you need to inspect results.

WHILE NOT WINDOW.SHOULDCLOSE()
    SERVER.TICK(TIME.DELTA())
    RENDER.FRAME()
WEND

LOBBY.FREE(lobby)
SERVER.STOP()

; === REMOTE CLIENT (typical) ===
CLIENT.CONNECT("192.168.1.100", 7777)

FUNCTION OnNet(msg)
    ; Incoming user-channel strings (see NET.md).
ENDFUNCTION
CLIENT.ONMESSAGE("OnNet")

WHILE NOT WINDOW.SHOULDCLOSE()
    CLIENT.TICK(TIME.DELTA())
    RENDER.FRAME()
WEND
CLIENT.STOP()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `LOBBY.MAKE(...)` | Deprecated alias of `LOBBY.CREATE`. |

---

## See also

- [FIRST_MULTIPLAYER_GAME.md](../tutorials/FIRST_MULTIPLAYER_GAME.md) — two-process `moonrun` walkthrough
- [MULTIPLAYER.md](MULTIPLAYER.md) — honest multiplayer scope + learning path
- [NET.md](NET.md) — server/client connection
- [NETWORK.md](NETWORK.md) — network overview and lobby pattern
