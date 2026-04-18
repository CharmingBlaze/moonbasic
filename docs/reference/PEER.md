# Peer Commands

`PEER.*` commands are documented in [NET.md](NET.md) under the **PEER Commands** section.

## Quick Reference

| Command | Description |
|---|---|
| `PEER.SEND(peer, channel, message, reliable)` | Send string to a peer |
| `PEER.SENDPACKET(peer, packet, channel)` | Send raw packet handle |
| `PEER.DISCONNECT(peer)` | Gracefully disconnect |
| `PEER.IP(peer)` | Peer IP address string |
| `PEER.PING(peer)` | Round-trip time in ms |

## See also

- [NET.md](NET.md) — full network stack
- [PACKET.md](PACKET.md) — raw packet handles
