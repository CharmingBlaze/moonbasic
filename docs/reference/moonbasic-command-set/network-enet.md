# Networking (ENet)

### `Net.Start()`
Initializes the networking system.

### `Net.Stop()`
Shuts down the networking system.

---

### `Net.CreateServer(port, maxClients)`
Creates a server host. Returns a handle.

### `Net.CreateClient()`
Creates a client host. Returns a handle.

---

### `Net.Connect(client, address, port)`
Connects to a server. Returns a peer handle.

### `Net.Update(host)`
Processes packets (non-blocking).

---

### `Net.Receive(host)`
Returns an event handle or 0.

### `Event.Type(id)`
Returns 1:Connect, 2:Disconnect, 3:Receive.

### `Event.Data(id)`
Returns received string data.

### `Event.Free(id)`
Frees the event resource.

Registry-first **`NET.*`** / **`PEER.*`** / **`EVENT.*`** map to the same implementation as legacy **`ENET.*`** — see **[ENET.md](../ENET.md)** for `ENET.CREATEHOST`, `ENET.HOSTSERVICE`, etc.

See also: [NETWORK.md](../NETWORK.md).
