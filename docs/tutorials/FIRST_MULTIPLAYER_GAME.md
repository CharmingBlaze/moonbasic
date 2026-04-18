# Your first multiplayer run (two processes)

This walkthrough gets you from zero to **two moonBASIC programs talking over UDP on one machine** (`127.0.0.1`). For **scope** (what the engine includes vs. what you integrate yourself), **layers** (`SERVER.*` vs `NET.*`), and **lobby limits**, read **[reference/MULTIPLAYER.md](../reference/MULTIPLAYER.md)** first or right after this page.

---

## What you need

- **Windows** or **Linux** with a **full runtime** build: **CGO on**, **`-tags fullruntime`**, and the **`moonrun`** entrypoint (or a `moonbasic` binary built with full runtime). See [BUILDING.md](../BUILDING.md) and [DEVELOPER.md](../DEVELOPER.md).
- **Two terminal windows** in the same clone of this repository.

---

## Steps

1. **Pick a UDP port** that nothing else is using (this repo’s samples use **`27777`**).
2. **Firewall** — The first time your build listens on UDP, **Windows Defender Firewall** may prompt you; allow access for local development if you want LAN tests later.
3. **Host** — In terminal A, from the repo root, run the high-level server sample.

   **Windows (PowerShell):**

   ```powershell
   $env:CGO_ENABLED="1"; go run -tags fullruntime ./cmd/moonrun testdata/mp_host.mb
   ```

   **Linux / macOS / Git Bash:**

   ```bash
   CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun testdata/mp_host.mb
   ```

4. **Client** — In terminal B, run the matching client (same pattern: set **`CGO_ENABLED=1`** then `go run … testdata/mp_client.mb`).

   The client connects to **`127.0.0.1:27777`**, then sends one **`RPC.CALLSERVER("PING", …)`** after **`CLIENT.ONCONNECT`** runs; the host’s **`FUNCTION PING`** receives it. See the source in **`testdata/mp_host.mb`** and **`testdata/mp_client.mb`**.

5. **Stop** — Use **Ctrl+C** in each terminal when you are done (or close the windows).

---

## Without running the game (compile check only)

From the repo root, the compiler can validate the same files without ENet linking at run time:

```bash
go run . --check testdata/mp_host.mb
go run . --check testdata/mp_client.mb
```

Mid-level JSON ping-pong samples:

```bash
go run . --check testdata/net_server.mb
go run . --check testdata/net_client.mb
```

---

## Next

- **[MULTIPLAYER.md](../reference/MULTIPLAYER.md)** — learning path, `LOBBY.*` semantics, and “not in the engine” (voice, global matchmaking).
- **[NET.md](../reference/NET.md)** — full command reference for `SERVER.*`, `CLIENT.*`, `RPC.*`, `NET.*`.
