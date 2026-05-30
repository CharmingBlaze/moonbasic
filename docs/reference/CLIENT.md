# Client Commands

`CLIENT.*` commands are documented in [NET.md](NET.md) under the **CLIENT Commands** section.

## Quick Reference

| Command | Description |
|---|---|
| `CLIENT.CONNECT(host, port)` | Connect to a server |
| `CLIENT.STOP()` | Disconnect |
| `CLIENT.ONCONNECT(handler)` | Callback when connected |
| `CLIENT.ONMESSAGE(handler)` | Callback on message received |
| `CLIENT.ONSYNC(handler)` | Callback on entity sync |
| `CLIENT.TICK(dt)` | Process network events (call each frame) |

## See also

- [NET.md](NET.md) — full network stack documentation
- [SERVER.md](SERVER.md) — server side
- [RPC.md](RPC.md) — remote procedure calls
