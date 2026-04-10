# Networking (ENet)

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **Net.Start / Stop** | **`NET.START`**, **`NET.STOP`** | Call **`NET.START`** before creating hosts. |
| **Net.CreateServer / CreateClient** | **`NET.CREATESERVER`**, **`NET.CREATECLIENT`** | Returns **host** handle; **`NET.CLOSE`** frees it. |
| **Net.SetChannels** | **`NET.SETCHANNELS`** `(count)` | Sets channel limit **for the next** **`CREATESERVER`** / **`CREATECLIENT`** (1–32). Default **1**. |
| **Net.Connect** | **`NET.CONNECT`** `(clientHost, host, port)` | Returns **peer** handle. Uses the configured channel count. |
| **Net.Service** | **`NET.SERVICE`** `(host, timeout_ms)` | Pumps the host with a blocking wait up to **`timeout_ms`** on the first **`Service`** call, then drains with **0** timeout. |
| **Net.Update** | **`NET.UPDATE`** `(host)` | Same pump with **0** ms wait (non-blocking style). |
| **Net.Receive** | **`NET.RECEIVE`** `(host)` | Returns **event** handle or **0**. |
| **Net.EventType / Peer / Data** | **`EVENT.TYPE`**, **`EVENT.PEER`**, **`EVENT.DATA`**, **`EVENT.CHANNEL`** | **`EVENT.TYPE`**: `1` connect, `2` disconnect, `3` receive. **`EVENT.FREE`** required. |
| **Net.Send** | **`PEER.SEND`** `(peer, channel, data, reliable)` | **Channel index before payload string.** |
| **Net.Broadcast** | **`NET.BROADCAST`** `(server, channel, data, reliable)` | |
| **Net.Disconnect** | **`PEER.DISCONNECT`** `(peer)` | |
| **Net.Flush** | **`NET.FLUSH`** `(host)` | Returns an error: upstream **`go-enet`** does not expose **`enet_host_flush`** on the **`Host`** interface. Pump with **`NET.UPDATE`** / **`NET.SERVICE`**. |
| **Packet.Create / Data / Free** | **`PACKET.CREATE`**, **`PACKET.DATA`**, **`PACKET.FREE`** | **`PEER.SENDPACKET`** `(peer, packet, channel)` transfers ownership to ENet and frees the VM packet handle. |

See also: [NETWORK.md](../NETWORK.md).
