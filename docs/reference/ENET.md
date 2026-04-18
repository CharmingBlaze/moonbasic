# ENET Commands

Legacy **`ENET.*`** names map to the **same ENet implementation** as **`NET.*`**, **`PEER.*`**, **`EVENT.*`**, and **`PACKET.*`** ([`runtime/net/enet_legacy_cgo.go`](../../runtime/net/enet_legacy_cgo.go)). Prefer **`NET.START`** / **`NET.CREATESERVER`** for new code; use **`ENET.*`** when following older samples or tutorials.

**Build:** Requires **CGO**, linked **libenet**, and a **full runtime** game host (`-tags fullruntime`, `CGO_ENABLED=1`). Without CGO, builtins return a clear error (see [NET_ZERO_CGO.md](../architecture/NET_ZERO_CGO.md)).

---

## Core Workflow

1. **`ENET.INITIALIZE()`** — same as **`NET.START()`**; initializes the ENet library once per process.
2. **`ENET.CREATEHOST`** (or **`NET.CREATESERVER`** / **`NET.CREATECLIENT`**) — allocate a host handle.
3. **`ENET.HOSTSERVICE`** or **`NET.UPDATE`** / **`NET.SERVICE`** — pump the host each frame; **`NET.RECEIVE`** returns **`EVENT.*`** handles.
4. Send with **`PEER.SEND`**, **`PEER.SENDPACKET`**, **`ENET.PEERSEND`**, or **`NET.BROADCAST`** / **`ENET.HOSTBROADCAST`** (packet-based).
5. **`ENET.DEINITIALIZE()`** — same as **`NET.STOP()`**; tears down hosts and deinitializes ENet.

---

## Library lifecycle

### `ENET.INITIALIZE()` 

Initializes the ENet C library. Idempotent after the first successful call in a process.

- **Arguments**: None.
- **Returns**: None.

---

### `ENET.DEINITIALIZE()` 

Closes all hosts tracked by the runtime and calls **`enet_deinitialize`**.

- **Arguments**: None.
- **Returns**: None.

---

## Host creation

### `ENET.CREATEHOST(address$, port, maxPeers, channels, bandwidth)` 

Creates a listening or bound host. **`address$`** empty means listen on all interfaces (same as **`NET.CREATESERVER`** with a port-only listen address). If **`address$`** is non-empty, it is passed to ENet as the bind host (hostname or IP). **`bandwidth`** is applied to **both** incoming and outgoing host bandwidth limits (bytes per second; **`0`** = unlimited in ENet). **`channels`** sets the ENet channel count (1–32) for this host, like **`NET.SETCHANNELS`** before **`NET.CREATESERVER`**.

- **Arguments**:
  - `address$` (string): Bind address; empty = any.
  - `port` (number): UDP port (0–65535).
  - `maxPeers` (number): Maximum peers (≥ 1).
  - `channels` (number): Channel count (1–32).
  - `bandwidth` (number): Incoming/outgoing bandwidth cap (≥ 0); large values are clamped to 32-bit.
- **Returns**: (handle) Host handle (`NetHost`).

---

### `ENET.MAKEHOST(address$, port, maxPeers, channels, bandwidth)` 

Alias: **`ENET.CREATEHOST`**.

---

## Host I/O

### `ENET.HOSTSERVICE(host, timeout_ms)` 

Same as **`NET.SERVICE(host, timeout_ms)`**: runs **`enet_host_service`** up to **`timeout_ms`** milliseconds (use **`0`** for a non-blocking pump).

- **Arguments**:
  - `host` (handle): Host from **`ENET.CREATEHOST`** / **`NET.CREATESERVER`** / **`NET.CREATECLIENT`**.
  - `timeout_ms` (number): Wait budget in milliseconds.
- **Returns**: None.

---

### `ENET.HOSTBROADCAST(host, channel, flags, packet)` 

Broadcasts a **`PACKET.*`** payload to all connected peers. **`flags`** is reserved (packet reliability is determined when the packet was created with **`PACKET.CREATE`** / **`enet.NewPacket`**).

- **Arguments**:
  - `host` (handle): Server host.
  - `channel` (number): ENet channel index (0–255).
  - `flags` (number): Reserved for API compatibility.
  - `packet` (handle): Packet handle; ownership transfers to ENet (handle is freed on success).
- **Returns**: None.

---

## Peer helpers

### `ENET.PEERSEND(peer, channel, packet)` 

Sends a packet to one peer. Argument order is **`(peer, channel, packet)`**; behavior matches **`PEER.SENDPACKET(peer, packet, channel)`** (see [PEER.md](PEER.md)).

- **Arguments**:
  - `peer` (handle): **`EVENT.PEER`** / **`PEER.*`** target.
  - `channel` (number): Channel index (0–255).
  - `packet` (handle): Packet handle; ownership transfers to ENet.
- **Returns**: None.

---

### `ENET.PEERPING(peer)` 

Returns measured round-trip time in milliseconds (same as **`PEER.PING`** / **`NET.GETPING`**).

- **Arguments**:
  - `peer` (handle): Connected peer.
- **Returns**: (integer) RTT in ms.

---

## Full Example

Minimal compile-time check (no window); run with **`go run . --check testdata/enet_smoke.mb`**.

```basic
; Initialize ENet and create a host on UDP 27777 (8 peers, 2 channels, unlimited bandwidth).
ENET.INITIALIZE()
h = ENET.CREATEHOST("", 27777, 8, 2, 0)
ENET.HOSTSERVICE(h, 0)
ENET.DEINITIALIZE()
```

For a full client/server loop with JSON messages, use **`testdata/net_server.mb`** / **`testdata/net_client.mb`** and the **`NET.*`** names in [NETWORK.md](NETWORK.md).

---

## See also

- [NETWORK.md](NETWORK.md) — recommended mid-level workflow (`NET.UPDATE`, `NET.RECEIVE`).
- [NET.md](NET.md) — **`NET.*`**, **`SERVER.*`**, **`RPC.*`**, **`PACKET.*`**, **`PEER.*`**, **`EVENT.*`**.
- [MULTIPLAYER.md](MULTIPLAYER.md) — multiplayer scope and learning path.
- [PEER.md](PEER.md) — peer send/disconnect helpers.
