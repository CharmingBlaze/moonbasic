# Network Commands

Commands for creating multiplayer games using ENet.

## Core Workflow

1.  **Initialize**: Call `Net.Start()` once.
2.  **Create Host**: Create a `Net.CreateServer()` or a `Net.CreateClient()`.
3.  **Connect**: If a client, use `Net.Connect()`.
4.  **Main Loop**: Inside the loop, call `Net.Update()` and then `Net.Receive()` repeatedly to process all incoming events for that frame.
5.  **Handle Events**: Use `Event.Type()` to check for connections, disconnections, and received data.
6.  **Send Data**: Use `Net.Broadcast()` (server) or `Peer.Send()` (client/server) to send messages.
7.  **Cleanup**: Call `Net.Stop()` before exiting.

---

## Host Management

### `Net.Start()` / `Net.Stop()`

Initializes and shuts down the entire networking system.

### `Net.CreateServer(port, maxClients)`

Creates a server host that listens for incoming connections.

### `Net.CreateClient()`

Creates a client host.

### `Net.Connect(clientHandle, address$, port)`

Connects a client to a server. Returns a handle to the server peer.

---

## Communication

### `Net.Update(hostHandle)`

This must be called every frame to process network packets.

### `Net.Receive(hostHandle)`

Retrieves the next available network event. Returns an event handle, or `0` if no events are waiting. You should call this in a loop until it returns `0`.

### `Net.Broadcast(serverHandle, channel, data$, reliable?)`

(Server-only) Sends a message to every connected client.

### `Peer.Send(peerHandle, channel, data$, reliable?)`

Sends a message to a specific peer.

- `reliable?`: `TRUE` guarantees delivery and order. `FALSE` is faster but packets can be lost or arrive out of order.

---

## Event Handling

When `Net.Receive()` returns a valid event handle, you must inspect it and then free it.

### `Event.Type(eventHandle)`

Returns the type of event:
- `EVENT_CONNECT`: A client connected (server-side) or the connection was successful (client-side).
- `EVENT_DISCONNECT`: A peer disconnected.
- `EVENT_RECEIVE`: Data was received.

### `Event.Peer(eventHandle)`

Returns the handle of the peer associated with the event.

### `Event.Data(eventHandle)`

For a `RECEIVE` event, this returns the string data that was sent.

### `Event.Free(eventHandle)`

Frees the event handle. **You must call this for every event you receive.**

---

## Server Example

```basic
; server.mb
Net.Start()
server = Net.CreateServer(1234, 32)

PRINT "Server started on port 1234..."

WHILE TRUE
    Net.Update(server)
    event = Net.Receive(server)
    WHILE event
        SELECT Event.Type(event)
            CASE EVENT_CONNECT
                PRINT "A client connected!"
            CASE EVENT_DISCONNECT
                PRINT "A client disconnected."
            CASE EVENT_RECEIVE
                PRINT "Got message: " + Event.Data(event)
                Net.Broadcast(server, 0, "Message received!", TRUE)
        ENDSELECT
        Event.Free(event)
        event = Net.Receive(server)
    WEND
WEND

Net.Stop()
```

## Client Example

```basic
; client.mb
Window.Open(400, 200, "Net Client")
Net.Start()
client = Net.CreateClient()
server_peer = Net.Connect(client, "127.0.0.1", 1234)

WHILE NOT Window.ShouldClose()
    Net.Update(client)
    event = Net.Receive(client)
    WHILE event
        IF Event.Type(event) = EVENT_RECEIVE THEN
            PRINT "Message from server: " + Event.Data(event)
        ENDIF
        Event.Free(event)
        event = Net.Receive(client)
    WEND

    IF Input.KeyPressed(KEY_SPACE) THEN
        Peer.Send(server_peer, 0, "Hello from the client!", TRUE)
    ENDIF

    Render.Clear(20,20,20)
    Render.Frame()
WEND

Net.Stop()
Window.Close()
```

---

## High-level server, client, RPC, and lobby (CGO)

The same ENet stack also exposes **opinionated** helpers in `runtime/net/mp_high_cgo.go`: a hosted **server tick**, **client tick**, **JSON RPC** over a dedicated channel, and **in-memory lobby** objects. These require **CGO** and a bound heap where noted.

### `SERVER.*` (dedicated host)

| Command | Purpose |
|---------|---------|
| `Server.Start(port, maxClients)` | Calls `Net.Start`, creates the server host, stores it globally. Fails if a server is already running. |
| `Server.Stop()` | Closes the server host and clears sync state. |
| `Server.OnConnect(functionName$)` / `Server.OnDisconnect(functionName$)` | Register user functions for peer connect/disconnect (names folded uppercase). |
| `Server.OnMessage(functionName$)` | Handler for non-RPC user traffic (see runtime for wire format). |
| `Server.SyncEntity(entityHandle, flags)` | Registers a **model** handle for periodic **transform sync** to clients (`flags` bitmask: transform bit = `1`). |
| `Server.SetTickRate(hz)` | Sets broadcast tick rate for sync flush inside `Server.Tick`. |
| `Server.Tick(dt#)` | Runs `Net.Update`, drains events, accumulates time, and periodically **broadcasts** transform sync packets. |

### `CLIENT.*`

| Command | Purpose |
|---------|---------|
| `Client.Connect(host$, port)` | Starts networking, creates a client host, connects, stores peer globally. |
| `Client.Stop()` | Closes client and clears peer. |
| `Client.OnConnect` / `Client.OnMessage` / `Client.OnSync` | Set user function names for connection, generic messages, and **sync** payloads (`Client.OnSync` receives decoded transform updates). |
| `Client.Tick(dt#)` | `Net.Update` + event drain for the client host. |

### `RPC.*` — JSON remote calls

Wire format is `MBRPC1:` + JSON `{"f":"FUNCNAME","a":[...]}` sent on **channel 2** (`chRPC`), **reliable**.

| Command | Role |
|---------|------|
| `RPC.Call(functionName$, ...)` | **Server → all clients** — broadcasts an RPC (server must be running). |
| `RPC.CallTo(peer, functionName$, ...)` | Send RPC to one **peer** handle. |
| `RPC.CallServer(functionName$, ...)` | **Client → server** — uses the connected server peer. |

Arguments are encoded as JSON numbers, strings, bools, or **handle ids as floats** (see `valueToJSONArg`). On the **server**, the callee receives **extra trailing peer handle** when your handler is invoked from `handleRPCPacket`.

### `LOBBY.*` — local lobby descriptors

Lightweight **heap** objects for matchmaking metadata (not network discovery by itself):

| Command | Purpose |
|---------|---------|
| `Lobby.Create(name$, maxPlayers)` → handle | Allocates a lobby; tracks it in a global list. |
| `Lobby.Free(lobby)` | Removes and frees the handle. |
| `Lobby.SetProperty(lobby, key$, value$)` | String properties (keys lowercased). |
| `Lobby.SetHost(lobby, host$, port)` | Address used by `Lobby.Join`. |
| `Lobby.Start(lobby)` | Marks lobby started (`started` flag). |
| `Lobby.Find(key$, value$)` → **heap array handle** | Lobbies whose property matches; handles stored as **floats** in the array (or a one-element array with `0` if none). |
| `Lobby.GetName(lobby)` → string | |
| `Lobby.Join(lobby)` | Calls `Client.Connect` with `Lobby`’s host/port (`SETHOST` required). |

For the lower-level **`Net.*` / `Peer.*` / `Event.*`** workflow, see the sections above.
