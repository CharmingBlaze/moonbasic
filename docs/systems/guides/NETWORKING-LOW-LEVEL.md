# Networking — packets and peers (mid level)

> After `SERVER.*` / `CLIENT.*`, use **`NET.*`** and **`PEER.*`** for custom channels, binary payloads, and manual poll loops.

**Namespaces:** `NET` · `PEER` · `EVENT` · `PACKET` · **Status:** Shipped

**Start here first:** [MULTIPLAYER.md](MULTIPLAYER.md) (high-level RPC)

**Commands:** [COMMAND_REGISTRY.md](../COMMAND_REGISTRY.md) · [reference/NETWORK.md](../../reference/NETWORK.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [Choose your networking layer](#choose-your-networking-layer)
- [Mid-level workflow](#mid-level-workflow)
- [Key commands](#key-commands)
- [Channels and reliability](#channels-and-reliability)
- [Sanity check samples](#sanity-check-samples)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Layer | You write | Good for |
|-------|-----------|----------|
| **SERVER / CLIENT** | Handler functions + RPC | Most games |
| **NET / PEER** | `UPDATE` + `RECEIVE` loop | Custom protocols |
| **ENET.*** | Legacy names | Old samples |

**Why mid-level:** Full control over **which channel**, **reliable vs unreliable**, and **payload format** (string vs bytes).

---

## Choose your networking layer

| Need | Use |
|------|-----|
| “Call function on server” | `RPC.CALLSERVER` — [MULTIPLAYER.md](MULTIPLAYER.md) |
| Position sync at 20 Hz | `NET.BROADCAST` unreliable channel |
| Custom binary struct | `PACKET.*` + `PEER.SEND` |
| Chat string | High-level `ONMESSAGE` or ch 0 |

---

## Mid-level workflow

**Server:**

```basic
NET.START()
host = NET.CREATESERVER(27777, 32)
WHILE running
    NET.UPDATE(host)
    WHILE NET.RECEIVE(host)
        ; handle PEER.* / EVENT.*
    WEND
    ; game + NET.BROADCAST(...)
WEND
NET.STOP()
```

**Client:** `NET.CREATECLIENT` → `NET.CONNECT` → same `UPDATE` / `RECEIVE` drain.

**Why drain `RECEIVE` in a while loop:** Multiple packets may arrive per frame.

---

## Key commands

| Command | Role |
|---------|------|
| `NET.START` / `STOP` | Stack lifetime |
| `NET.CREATESERVER(port, max)` | Listen host handle |
| `NET.CREATECLIENT()` | Outgoing host |
| `NET.CONNECT(host, ip, port)` | Client connect |
| `NET.UPDATE(host)` | Poll sockets |
| `NET.RECEIVE(host)` | Dequeue one message (loop) |
| `PEER.SEND(peer, data, channel, reliable)` | Target one peer |
| `NET.BROADCAST(host, data, ch, reliable)` | All clients |

---

## Channels and reliability

High-level stack uses fixed channels internally:

| Channel | Typical use |
|---------|-------------|
| 0 | User strings / chat |
| 1 | Entity sync |
| 2 | RPC |

Mid-level: you pass **channel index** and **reliable** flag per send — match game design (position often unreliable, inventory reliable).

---

## Sanity check samples

```bash
moonbasic --check testdata/net_server.mb
moonbasic --check testdata/net_client.mb
moonrun testdata/net_server.mb    ; terminal A
moonrun testdata/net_client.mb    ; terminal B
```

High-level pair: `mp_host.mb` / `mp_client.mb` — [FIRST_MULTIPLAYER_GAME.md](../../tutorials/FIRST_MULTIPLAYER_GAME.md).

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| No `NET.UPDATE` | No packets received |
| Single `RECEIVE` per frame | Drain with `WHILE` |
| Reliable flood on positions | Use unreliable for high-rate sync |
| Skip high-level RPC | Use RPC when shape fits |

---

## See also

- [MULTIPLAYER.md](MULTIPLAYER.md)
- [reference/ENET.md](../../reference/ENET.md)
