# go-enet
Enet bindings for Go using cgo.

## Installation
First, you might need to install enet as a dependency:

* **Windows**: Nothing to do - should work out of the box with the supplied headers and library.
* **Linux**: Install the enet development package with your package manager.
	* On Debian-based systems: `apt install libenet-dev`
* **MacOS**: Install the enet package with brew: `brew install enet`

Since this module uses cgo, you must also make sure you have a C compiler in your `PATH`. For Linux and Mac, this is usually pretty straight forward. On Windows, you may have to use something like [MSYS2](https://www.msys2.org/).

When ready, get the module like this:

```
$ go get github.com/codecat/go-enet
```

## Usage
```go
import "github.com/codecat/go-enet"
```

The API is mostly the same as the C API, except it's more object-oriented.

## Server example
This is a basic server example that responds to packets `"ping"` and `"bye"`.

```go
package main

import (
	"github.com/codecat/go-enet"
	"github.com/codecat/go-libs/log"
)

func main() {
	// Initialize enet
	enet.Initialize()

	// Create a host listening on 0.0.0.0:8095
	host, err := enet.NewHost(enet.NewListenAddress(8095), 32, 1, 0, 0)
	if err != nil {
		log.Error("Couldn't create host: %s", err.Error())
		return
	}

	// The event loop
	for true {
		// Wait until the next event
		ev := host.Service(1000)

		// Do nothing if we didn't get any event
		if ev.GetType() == enet.EventNone {
			continue
		}

		switch ev.GetType() {
		case enet.EventConnect: // A new peer has connected
			log.Info("New peer connected: %s", ev.GetPeer().GetAddress())

		case enet.EventDisconnect: // A connected peer has disconnected
			log.Info("Peer disconnected: %s", ev.GetPeer().GetAddress())

		case enet.EventReceive: // A peer sent us some data
			// Get the packet
			packet := ev.GetPacket()

			// We must destroy the packet when we're done with it
			defer packet.Destroy()

			// Get the bytes in the packet
			packetBytes := packet.GetData()

			// Respond "pong" to "ping"
			if string(packetBytes) == "ping" {
				ev.GetPeer().SendString("pong", ev.GetChannelID(), enet.PacketFlagReliable)
				continue
			}

			// Disconnect the peer if they say "bye"
			if string(packetBytes) == "bye" {
				log.Info("Bye!")
				ev.GetPeer().Disconnect(0)
				continue
			}
		}
	}

	// Destroy the host when we're done with it
	host.Destroy()

	// Uninitialize enet
	enet.Deinitialize()
}
```

## Client example
This is a basic client example that sends a ping to the server every second that there is no event.

```go
package main

import (
	"github.com/codecat/go-enet"
	"github.com/codecat/go-libs/log"
)

func main() {
	// Initialize enet
	enet.Initialize()

	// Create a client host
	client, err := enet.NewHost(nil, 1, 1, 0, 0)
	if err != nil {
		log.Error("Couldn't create host: %s", err.Error())
		return
	}

	// Connect the client host to the server
	peer, err := client.Connect(enet.NewAddress("127.0.0.1", 8095), 1, 0)
	if err != nil {
		log.Error("Couldn't connect: %s", err.Error())
		return
	}

	// The event loop
	for true {
		// Wait until the next event
		ev := client.Service(1000)

		// Send a ping if we didn't get any event
		if ev.GetType() == enet.EventNone {
			peer.SendString("ping", 0, enet.PacketFlagReliable)
			continue
		}

		switch ev.GetType() {
		case enet.EventConnect: // We connected to the server
			log.Info("Connected to the server!")

		case enet.EventDisconnect: // We disconnected from the server
			log.Info("Lost connection to the server!")

		case enet.EventReceive: // The server sent us data
			packet := ev.GetPacket()
			log.Info("Received %d bytes from server", len(packet.GetData()))
			packet.Destroy()
		}
	}

	// Destroy the host when we're done with it
	client.Destroy()

	// Uninitialize enet
	enet.Deinitialize()
}
```
