# Packet Commands

`PACKET.*` commands are documented in [NET.md](NET.md) under the **PACKET Commands** section.

## Quick Reference

| Command | Description |
|---|---|
| `PACKET.CREATE(data)` | Create a raw packet handle |
| `PACKET.DATA(packet)` | Read payload string |
| `PACKET.FREE(packet)` | Free the handle |

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `PACKET.MAKE(...)` | Deprecated alias of `PACKET.CREATE`. |

## See also

- [NET.md](NET.md) — full network stack
- [PEER.md](PEER.md) — sending packets to peers
