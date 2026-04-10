# Networking helpers (typed send / read)

These builtins make binary **`NETSEND*`** / **`NETREAD*`** easier than raw strings.

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **NetSendString** | **`NETSENDSTRING`** `(peer, text)` | Sends **4-byte little-endian length** + **raw UTF-8** bytes (reliable, channel **0**). Pair with **`NETREADSTRING`** on the receiver after **`EVENT.DATA`**. |
| **NetSendInt** | **`NETSENDINT`** `(peer, value)` | **4** bytes, little-endian **`int32`** (reliable, channel **0**). |
| **NetSendFloat** | **`NETSENDFLOAT`** `(peer, value)` | **8** bytes, **`float64`** LE (reliable, channel **0**). |
| **NetReadString / Int / Float** | **`NETREADSTRING`**, **`NETREADINT`**, **`NETREADFLOAT`** | Read cursor is filled when **`EVENT.DATA`** runs on a **`RECEIVE`** event (full payload copied). Call **`EVENT.DATA`** first, then **`NETREAD*`** in order. **`NETREADSTRING`** expects the **`NETSENDSTRING`** length prefix. |

**Event order:** dequeue with **`NET.RECEIVE`**, branch on **`EVENT.TYPE`**, for receive events call **`EVENT.DATA`** (updates read buffer), then **`NETREAD*`** as needed, then **`EVENT.FREE`**.

See also: [network-enet.md](network-enet.md), [NETWORK.md](../NETWORK.md).
