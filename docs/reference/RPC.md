# RPC Commands

`RPC.*` commands are documented in [NET.md](NET.md) under the **RPC Commands** section.

## Quick Reference

| Command | Description |
|---|---|
| `RPC.CALL(fn [, args...])` | Call a function on all peers |
| `RPC.CALLTO(peer, fn [, args...])` | Call a function on a specific peer |
| `RPC.CALLSERVER(fn [, args...])` | Call a function on the server |

Each overload supports up to 7 arguments.

## See also

- [NET.md](NET.md) — full network stack documentation
- [EVENT.md](EVENT.md) — event system
