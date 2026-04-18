# ENet Commands

`ENET.*` commands are documented in [NET.md](NET.md) under the **ENET Commands (Low-level)** section.

## Quick Reference

| Command | Description |
|---|---|
| `ENET.INITIALIZE()` | Start raw ENet library |
| `ENET.DEINITIALIZE()` | Stop ENet |
| `ENET.CREATEHOST(addr, port, peers, channels, inBps, outBps)` | Create raw ENet host |
| `ENET.HOSTSERVICE(host, timeout)` | Dispatch events |
| `ENET.HOSTBROADCAST(host, channel, flags, packet)` | Broadcast to all peers |
| `ENET.PEERSEND(peer, channel, packet)` | Send to one peer |
| `ENET.PEERPING(peer)` | Measure latency |

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `ENET.MAKEHOST(port, maxPeers)` | Deprecated alias of `ENET.HOST`. |

## See also

- [NET.md](NET.md) — full network stack
- [PEER.md](PEER.md) — peer handles
