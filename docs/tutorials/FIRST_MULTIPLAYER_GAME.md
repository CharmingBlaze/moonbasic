# Your first multiplayer run (two processes)

This walkthrough gets you from zero to **two moonBASIC programs talking over UDP on one machine** (`127.0.0.1`). For **scope** (what the engine includes vs. what you integrate yourself), **layers** (`SERVER.*` vs `NET.*`), and **lobby limits**, read **[reference/MULTIPLAYER.md](../reference/MULTIPLAYER.md)** first or right after this page.

---

## What you need

- **Full runtime** from [GitHub Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest) on **Windows** or **Linux** — you need **`moonrun`** (compiler-only archives do not run networked games with a window).
- **Two terminal windows** in the same folder that contains the sample scripts (`testdata/` in this repository).

---

## Steps

1. **Pick a UDP port** that nothing else is using (this repo’s samples use **`27777`**).
2. **Firewall** — The first time your game listens on UDP, **Windows Defender Firewall** may prompt you; allow access for local development if you want LAN tests later.
3. **Host** — In terminal A, from the repo root, run the high-level server sample:

   ```bash
   moonrun testdata/mp_host.mb
   ```

   On Windows:

   ```bat
   moonrun.exe testdata\mp_host.mb
   ```

4. **Client** — In terminal B, run the matching client:

   ```bash
   moonrun testdata/mp_client.mb
   ```

   The client connects to **`127.0.0.1:27777`**, then sends one **`RPC.CALLSERVER("PING", …)`** after **`CLIENT.ONCONNECT`** runs; the host’s **`FUNCTION PING`** receives it. See **`testdata/mp_host.mb`** and **`testdata/mp_client.mb`**.

5. **Stop** — Use **Ctrl+C** in each terminal when you are done (or close the windows).

---

## Without running the game (check only)

Validate the same files without opening a window or linking the network stack at run time:

```bash
moonbasic --check testdata/mp_host.mb
moonbasic --check testdata/mp_client.mb
```

Mid-level JSON ping-pong samples:

```bash
moonbasic --check testdata/net_server.mb
moonbasic --check testdata/net_client.mb
```

---

## Next

- **[MULTIPLAYER.md](../reference/MULTIPLAYER.md)** — learning path, `LOBBY.*` semantics, and “not in the engine” (voice, global matchmaking).
