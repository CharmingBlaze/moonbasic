# Multiplayer — two programs over the network

> Run a **host** and **client** (two `moonrun` processes) that talk over **UDP** with reliable channels — LAN, same PC, or known IP.

**Namespaces:** `SERVER` · `CLIENT` · `RPC` · `NET` · `PEER` · `LOBBY` · **Status:** Shipped · **Platform:** full runtime (Windows / Linux)

**Commands:** [COMMAND_REGISTRY.md](../COMMAND_REGISTRY.md) (search NET, SERVER) · [reference/MULTIPLAYER.md](../../reference/MULTIPLAYER.md)

**Walkthrough:** [../../tutorials/FIRST_MULTIPLAYER_GAME.md](../../tutorials/FIRST_MULTIPLAYER_GAME.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use multiplayer APIs](#when-to-use-multiplayer-apis)
- [Choose the right layer](#choose-the-right-layer)
- [Core workflow (high level)](#core-workflow-high-level)
- [Host program](#host-program)
- [Client program](#client-program)
- [RPC — call functions on the other side](#rpc--call-functions-on-the-other-side)
- [What moonBASIC does not ship](#what-moonbasic-does-not-ship)
- [Full local test](#full-local-test)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **Transport** | UDP via ENet (reliable + unreliable channels) |
| **Processes** | Usually **one server** `.mb` + **one client** `.mb` per player (or N clients) |
| **Tick** | You call `SERVER.TICK` / `CLIENT.TICK` every frame — **no hidden game thread** |
| **Matchmaking** | You share IP/port (Discord, LAN) — not Steam-style global browser |

**Why tick in your loop:** moonBASIC does not run your gameplay on a background network thread. Network I/O is polled when you tick — same pattern as input and physics.

---

## When to use multiplayer APIs

**Use when:**

- Co-op or competitive games on LAN or dedicated server.
- Small player counts (dozens, not MMO scale on one process).
- You design **server authority** (server validates actions).

**Skip when:**

- Single-player only.
- You need global matchmaking / voice — integrate external services (see below).

---

## Choose the right layer

| Layer | APIs | Best for |
|-------|------|----------|
| **High (start here)** | `SERVER.*`, `CLIENT.*`, `RPC.*` | Game loops, function RPCs |
| **Mid** | `NET.*`, `PEER.*`, `NET.RECEIVE` | Custom packets, channels |
| **Legacy names** | `ENET.*` | Old samples — prefer `NET.*` |
| **Lobby handles** | `LOBBY.*` | **In-process only** — not Internet discovery |

**Why high level first:** `SERVER.START` + `CLIENT.CONNECT` + `RPC.CALLSERVER` covers “host runs logic, client asks server” without manual packet parsing.

---

## Core workflow (high level)

### Server

1. `SERVER.START(port, maxClients)` — **Why:** Opens UDP listen socket.
2. Loop: `SERVER.TICK(dt)` — **Why:** Polls connections, delivers messages, runs RPC handlers.
3. `SERVER.STOP()` — **Why:** Clean shutdown.

### Client

1. `CLIENT.CONNECT(host, port)` — **Why:** Outgoing connection to server.
2. `CLIENT.ONCONNECT("MyHandler")` — **Why:** Run setup when connected (e.g. send first RPC).
3. Loop: `CLIENT.TICK(dt)` — **Why:** Same polling model as server.
4. `CLIENT.STOP()`.

**dt** is often `0.016` (~60 Hz) or `APP.DELTA()` if you also render a window.

---

## Host program

Minimal host from [`testdata/mp_host.mb`](../../../testdata/mp_host.mb):

```basic
; Host: waits for RPC "PING" from client
SERVER.START(27777, 8)
done = 0
WHILE done = 0
    SERVER.TICK(0.016)
WEND
SERVER.STOP()

FUNCTION PING(msg, peerH)
    ; msg from client; peerH identifies connection
    done = 1
ENDFUNCTION
```

**Why `FUNCTION PING`:** `RPC.CALLSERVER("PING", …)` on the client invokes a **user function** on the server by name.

---

## Client program

From [`testdata/mp_client.mb`](../../../testdata/mp_client.mb):

```basic
CLIENT.CONNECT("127.0.0.1", 27777)
CLIENT.ONCONNECT("ONCONNECTED")

FUNCTION ONCONNECTED()
    RPC.CALLSERVER("PING", "hello")
ENDFUNCTION

i = 0
WHILE i < 2000
    CLIENT.TICK(0.016)
    i = i + 1
WEND
CLIENT.STOP()
```

**Why `ONCONNECT`:** RPCs sent before connect fail — wait for the handshake.

---

## RPC — call functions on the other side

| Command | Direction |
|---------|-----------|
| `RPC.CALLSERVER(name, …)` | Client → server function |
| `RPC.CALL(name, …)` | Server → all clients (pattern varies — see NET doc) |
| `RPC.CALLTO(peer, name, …)` | Target one peer |

**Why RPC:** Structured calls instead of parsing raw strings in `ONMESSAGE`. Uses dedicated channel `chRPC` internally.

Handlers are normal moonBASIC functions:

```basic
FUNCTION SpawnEnemy(type, peerH)
    ; server-only spawn logic
ENDFUNCTION
```

---

## What moonBASIC does not ship

| Feature | What to do |
|---------|------------|
| Global server browser | Your web API or Steam — distribute IP/port |
| Voice chat | Discord, Steam Voice, etc. |
| Anti-cheat | Validate on server you control |
| `LOBBY.FIND("game", "name")` | Only searches **in-process** lobbies — not WAN |

---

## Full local test

**Terminal A (host):**

```bash
moonrun testdata/mp_host.mb
```

**Terminal B (client):**

```bash
moonrun testdata/mp_client.mb
```

**Check without running:**

```bash
moonbasic --check testdata/mp_host.mb
moonbasic --check testdata/mp_client.mb
```

**Firewall:** Windows may prompt on first UDP listen — allow for LAN tests.

**Port:** Pick free UDP (samples use `27777`).

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| No `TICK` in loop | Nothing receives — call every frame |
| RPC before connect | Use `ONCONNECT` handler |
| Same port twice | One server per port |
| Expect `LOBBY.FIND` online | Use known host/IP |
| Compiler-only install for play | Need full runtime for live net |

---

## See also

- [NETWORKING-LOW-LEVEL.md](NETWORKING-LOW-LEVEL.md) — `NET.*` / `PEER.*` when RPC is not enough
- [FIRST_MULTIPLAYER_GAME.md](../../tutorials/FIRST_MULTIPLAYER_GAME.md)
- [reference/NETWORK.md](../../reference/NETWORK.md) — mid-level patterns
- [reference/LOBBY.md](../../reference/LOBBY.md) — lobby limits
