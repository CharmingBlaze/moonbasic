# Server Commands

`SERVER.*` commands are documented in [NET.md](NET.md) under the **SERVER Commands** section.

## Quick Reference

| Command | Description |
|---|---|
| `SERVER.START(port, maxClients)` | Start the game server |
| `SERVER.STOP()` | Shut down the server |
| `SERVER.ONCONNECT(handler)` | Callback on client connect |
| `SERVER.ONDISCONNECT(handler)` | Callback on client disconnect |
| `SERVER.ONMESSAGE(handler)` | Callback on message |
| `SERVER.SYNCENTITY(entity, rate)` | Auto-replicate an entity |
| `SERVER.SETTICKRATE(rate)` | Set server tick rate |
| `SERVER.TICK(dt)` | Process network events (call each frame) |

## See also

- [NET.md](NET.md) — full network stack documentation
- [CLIENT.md](CLIENT.md) — client side
- [LOBBY.md](LOBBY.md) — lobby discovery
