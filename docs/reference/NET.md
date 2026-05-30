# Net / Server / Client / Peer / Packet / ENet / RPC Commands

Beginner-oriented overview and learning order: **[MULTIPLAYER.md](MULTIPLAYER.md)**.

moonBASIC's network stack. Multiple layers are available depending on control level needed:

| Layer | Use when |
|---|---|
| **`SERVER.*` / `CLIENT.*`** | High-level lobby game server with entity sync |
| **`NET.*`** | Mid-level ENet wrapper with peer handles |
| **`ENET.*`** | Legacy names for the same stack — see **[ENET.md](ENET.md)** |
| **`RPC.*`** | Remote procedure calls over an active connection |
| **`PEER.*` / `PACKET.*`** | Per-peer send and raw packet building |

---

## Core Workflow (High-level)

**Server side:**
```basic
SERVER.START(port, maxClients)
SERVER.ONCONNECT("OnConnect")
SERVER.ONMESSAGE("OnMessage")
WHILE running
    SERVER.TICK(dt)
WEND
SERVER.STOP()
```

**Client side:**
```basic
CLIENT.CONNECT(host, port)
CLIENT.ONMESSAGE("OnMessage")
WHILE running
    CLIENT.TICK(dt)
WEND
CLIENT.STOP()
```

---

## SERVER Commands

### `SERVER.START(port, maxClients)` 

Starts a game server on `port` accepting up to `maxClients` connections.

---

### `SERVER.STOP()` 

Shuts down the server and disconnects all clients.

---

### `SERVER.ONCONNECT(handlerName)` 

Registers a function called when a client connects. The handler receives a peer handle.

---

### `SERVER.ONDISCONNECT(handlerName)` 

Registers a function called when a client disconnects.

---

### `SERVER.ONMESSAGE(handlerName)` 

Registers a function called when a message arrives. The handler receives the message string and peer handle.

---

### `SERVER.SYNCENTITY(entityHandle, rate)` 

Marks an entity for automatic replication to all clients at `rate` times per second.

---

### `SERVER.SETTICKRATE(rate)` 

Sets the server tick rate (updates per second).

---

### `SERVER.TICK(dt)` 

Processes pending network events. Call every frame.

---

## CLIENT Commands

### `CLIENT.CONNECT(host, port)` 

Connects to a server at `host:port`.

---

### `CLIENT.STOP()` 

Disconnects from the server.

---

### `CLIENT.ONCONNECT(handlerName)` 

Registers a function called when the connection is established.

---

### `CLIENT.ONMESSAGE(handlerName)` 

Registers a function called on incoming messages.

---

### `CLIENT.ONSYNC(handlerName)` 

Registers a function called when synced entity data arrives.

---

### `CLIENT.TICK(dt)` 

Processes pending network events. Call every frame.

---

## NET Commands (Mid-level)

### `NET.START()` / `NET.STOP()` 

Initialise / shut down the NET subsystem.

---

### `NET.CREATESERVER(port, maxPeers)` / `NET.MAKESERVER(port, maxPeers)` 

Creates a server host handle on `port`. `NET.MAKESERVER` is an alias.

---

### `NET.CREATECLIENT()` / `NET.MAKECLIENT()` 

Creates a client host handle (not yet connected). `NET.MAKECLIENT` is an alias.

---

### `NET.CONNECT(hostHandle, address, port)` / `NET.CONNECT(address, port)` 

Connects a client host to `address:port`. Returns a peer handle, or simplified single-command form.

---

### `NET.UPDATE(hostHandle)` / `NET.SERVICE(hostHandle, timeout)` 

Dispatches queued events for the host.

---

### `NET.RECEIVE(hostHandle)` 

Returns the next queued event handle (or `0` if none). Read via `EVENT.*`.

---

### `NET.BROADCAST(hostHandle, channel, message, reliable)` 

Broadcasts a string message to all peers on `channel`. `reliable` = `TRUE` for guaranteed delivery.

---

### `NET.PEERCOUNT(hostHandle)` 

Returns the number of connected peers.

---

### `NET.GETPING(peerHandle)` 

Returns round-trip time in milliseconds.

---

### `NET.SETTIMEOUT(hostHandle, ms)` / `NET.SETBANDWIDTH(hostHandle, inBps, outBps)` / `NET.SETCHANNELS(count)` 

Tuning options for timeout, bandwidth, and channel count.

---

### `NET.FLUSH(hostHandle)` 

Forces immediate send of queued packets.

---

### `NET.CLOSE(handle)` 

Closes a host or peer handle.

---

### `NET.HOST(port)` / `NET.SEND(channel, message)` / `NET.SYNC(entityId)` 

Simplified one-liners for quick hosting, sending, and entity sync.

---

## PEER Commands

### `PEER.SEND(peer, channel, message, reliable)` 

Sends a string message to a single peer.

---

### `PEER.SENDPACKET(peer, packet, channel)` 

Sends a raw `PACKET.*` handle to a peer.

---

### `PEER.DISCONNECT(peer)` 

Gracefully disconnects a peer.

---

### `PEER.IP(peer)` 

Returns the peer's IP address string.

---

### `PEER.PING(peer)` 

Returns the peer's round-trip time in milliseconds.

---

## PACKET Commands

### `PACKET.CREATE(data)` 

Creates a raw packet handle with `data` string payload. Returns a **packet handle**.

---

### `PACKET.DATA(packet)` 

Returns the payload string from a packet.

---

### `PACKET.FREE(packet)` 

Frees the packet handle.

---

## ENET Commands (Low-level)

Legacy **`ENET.*`** names are documented on **[ENET.md](ENET.md)** (same implementation as **`NET.*`** / **`PEER.*`** / **`PACKET.*`**).

---

## RPC Commands

### `RPC.CALL(funcName [, arg1...arg7])` 

Calls a named function on all connected peers.

---

### `RPC.CALLTO(peerHandle, funcName [, arg1...arg7])` 

Calls a named function on a specific peer.

---

### `RPC.CALLSERVER(funcName [, arg1...arg7])` 

Calls a named function on the server (client → server direction).

---

## Full Example

Minimal two-player chat using the high-level SERVER/CLIENT layer.

```basic
; === SERVER (run this in one instance) ===
SERVER.START(7777, 8)
SERVER.SETTICKRATE(30)

FUNCTION OnMessage(msg, peer)
    PRINT "Client says: " + msg
    PEER.SEND(peer, 0, "Echo: " + msg, TRUE)
END FUNCTION
SERVER.ONMESSAGE("OnMessage")

WHILE NOT WINDOW.SHOULDCLOSE()
    SERVER.TICK(TIME.DELTA())
    RENDER.FRAME()
WEND
SERVER.STOP()

; === CLIENT (run in another instance) ===
CLIENT.CONNECT("127.0.0.1", 7777)

FUNCTION OnMsg(msg)
    PRINT "Server says: " + msg
END FUNCTION
CLIENT.ONMESSAGE("OnMsg")

WHILE NOT WINDOW.SHOULDCLOSE()
    CLIENT.TICK(TIME.DELTA())
    IF INPUT.KEYPRESSED(KEY_SPACE) THEN
        NET.SEND(0, "Hello!")
    END IF
    RENDER.FRAME()
WEND
CLIENT.STOP()
```

---

## See also

- [NETWORK.md](NETWORK.md) — lobby/matchmaking overview
- [EVENT.md](EVENT.md) — event system used by `NET.RECEIVE`
- [LOBBY.md](LOBBY.md) — lobby discovery and joining
