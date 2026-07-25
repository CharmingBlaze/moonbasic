# Networking without CGO (`!cgo`)

When `CGO_ENABLED=0`, moonBASIC uses stubs for the ENet-backed natives in `runtime/net/` (see `enet_cgo.go` vs `stub.go`). Stock **libenet** is not linked.

## Goals for a pure-Go layer

- **No C toolchain**: implement a small **UDP**-based session layer in Go using `net.UDPConn`.
- **Semantics**: games need **reliable ordered** channels for some messages and **unreliable** for others — mirror ENet’s *channels* concept even if the wire format differs.

## Wire compatibility

| Approach | Wire compatible with stock ENet? | Notes |
|----------|-----------------------------------|--------|
| **Go UDP + custom framing** | **No** (unless you implement ENet byte-for-byte) | Simplest; version the protocol (`PROTO` handshake / semver in moonBASIC terms). |
| **KCP** ([kcp-go](https://github.com/xtaci/kcp-go)) | **No** | Good for game traffic; tune for latency vs reliability. |
| **Full ENet in Go** | **Yes** | Large effort; only needed for interoperability with existing ENet binaries. |

## Recommended migration path

1. Ship a **new protocol version** when moving from C ENet to pure Go (document in `docs/reference/NETWORK.md` when implemented).
2. Optional **dual-stack** period: detect peer capability or compile-time flag — outside the scope of this stub-era doc.

Until a native implementation lands, `CGO_ENABLED=1` remains the path for ENet-compatible networked builds.
