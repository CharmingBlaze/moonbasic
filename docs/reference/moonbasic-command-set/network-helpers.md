# Networking helpers (typed send / read)
 
These builtins make binary **`NETSEND*`** / **`NETREAD*`** easier than raw strings.
 
### `NetStart()`
Initializes the networking system.
 
### `NetStop()`
Shuts down the networking system.
 
---
 
### `NetUpdate(host)`
Processes packets (non-blocking).
 
### `NetReceive(host)`
Returns an event handle or 0.
 
---
 
### `NetEventType(ev)`
Returns event type.
 
### `NetEventData(ev)`
Returns received string data.
 
### `NetEventFree(ev)`
Frees the event resource.
 
**Event order:** dequeue with **`NetReceive()`**, branch on **`EventType()`**, for receive events call **`EventData()`** (updates read buffer), then **`NetRead*()`** as needed, then **`EventFree()`**.
 
| Designed | moonBASIC | Memory / notes |
|----------|------------|----------------|
| **NetStart()** | **`Net.Start()`** | Initializes networking. |
| **NetStop()** | **`Net.Stop()`** |  |
| **NetUpdate(host)** | **`Net.Update()`** | Processes packets. |
| **NetReceive(host)** | **`Net.Receive()`** | Returns event or 0. |
| **NetSend(peer, s)** | **`Peer.Send()`** |  |
| **NetBroadcast(host, s)** | **`Net.Broadcast()`** |  |
| **NetEventType(ev)** | **`Event.Type()`** |  |
| **NetEventPeer(ev)** | **`Event.Peer()`** |  |
| **NetEventData(ev)** | **`Event.Data()`** |  |
| **NetEventFree(ev)** | **`Event.Free()`** |  |
 
See also: [network-enet.md](network-enet.md), [NETWORK.md](../NETWORK.md).
# Networking helpers (typed send / read)
 
These builtins make binary **`NETSEND*`** / **`NETREAD*`** easier than raw strings.
 
### `NetStart()`
Initializes the networking system. Must be called before any other network commands.
 
### `NetStop()`
Shuts down the networking system and releases all host resources.
 
---
 
### `NetUpdate(host)`
Processes incoming packets and sends outgoing ones for the specified host (non-blocking).
 
### `NetReceive(host)`
Returns a handle to the next available network event, or `0` if the queue is empty.
 
---
 
### `NetEventType(ev)`
Returns the type of the event (1: Connect, 2: Disconnect, 3: Receive).
 
### `NetEventData(ev)`
Returns the payload of a `RECEIVE` event as a string.
 
### `NetEventFree(ev)`
Releases the event handle. **Must be called for every event received** to avoid leaks.
 
**Event order:** dequeue with **`NetReceive()`**, branch on **`EventType()`**, for receive events call **`EventData()`** (updates read buffer), then **`NetRead*()`** as needed, then **`EventFree()`**.
 
| Designed | moonBASIC | Memory / notes |
|----------|------------|----------------|
| **NetStart()** | **`Net.Start()`** | Initializes networking. |
| **NetStop()** | **`Net.Stop()`** |  |
| **NetUpdate(host)** | **`Net.Update()`** | Processes packets. |
| **NetReceive(host)** | **`Net.Receive()`** | Returns event or 0. |
| **NetSend(peer, s)** | **`Peer.Send()`** |  |
| **NetBroadcast(host, s)** | **`Net.Broadcast()`** |  |
| **NetEventType(ev)** | **`Event.Type()`** |  |
| **NetEventPeer(ev)** | **`Event.Peer()`** |  |
| **NetEventData(ev)** | **`Event.Data()`** |  |
| **NetEventFree(ev)** | **`Event.Free()`** |  |
 
See also: [network-enet.md](network-enet.md), [NETWORK.md](../NETWORK.md).