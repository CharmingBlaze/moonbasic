package enet_test

import (
	"fmt"
	"github.com/codecat/go-enet"
	"log"
)

func Example_clientServer() {
	// This example demonstrates some basic usage of the enet library.
	// Here we set up a client & server, send a message between them, then
	// disconnect & cleanup.

	port := uint16(1234)

	fmt.Printf("enet version: %s\n", enet.LinkedVersion())

	// Initialize enet
	enet.Initialize()

	// Make our server.
	server, err := enet.NewHost(enet.NewListenAddress(port), 32, 1, 0, 0)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to create server: %w", err))
	}

	// For this example, we're going to wait until a disconnect event has been
	// properly handled. Set this up here.
	disconnected := make(chan bool, 0)

	// Setup our server handling running in a separate goroutine.
	go func() {
		for true {
			ev := server.Service(10)

			switch ev.GetType() {
			case enet.EventConnect:
				fmt.Printf("[SERVER] new connection from client\n")
			case enet.EventReceive:
				fmt.Printf("[SERVER] received packet from client: %s\n", ev.GetPacket().GetData())

				// We must destroy the packet when we're done with it
				ev.GetPacket().Destroy()

				// Send back a message to the client.
				err := ev.GetPeer().SendString("message received!", 0, enet.PacketFlagReliable)
				if err != nil {
					log.Fatal(err)
				}
			case enet.EventDisconnect:
				fmt.Printf("[SERVER] client disconnected")
				close(disconnected)
			}
		}
	}()

	// Make a client that will speak to the server.
	client, err := enet.NewHost(nil, 32, 1, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Connect to the server.
	peer, err := client.Connect(enet.NewAddress("localhost", port), 1, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Keep checking the client until we get a response from the server.
	done := false
	for !done {
		ev := client.Service(10)

		switch ev.GetType() {
		case enet.EventReceive:
			fmt.Printf("[CLIENT] received packet from server: %s\n", string(ev.GetPacket().GetData()))
			ev.GetPacket().Destroy()
			done = true
		case enet.EventNone:
			// If nothing else to do, send a packet.
			fmt.Printf("[CLIENT] sending packet to server\n")
			err = peer.SendString("hello world", 0, enet.PacketFlagReliable)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	// Immediately disconnect the client.
	peer.DisconnectNow(0)

	// Wait for the disconnection to be handled by the server.
	<-disconnected

	// Cleanup.
	client.Destroy()
	server.Destroy()
	enet.Deinitialize()

	// Output:
	// enet version: 1.3.17
	// [SERVER] new connection from client
	// [CLIENT] sending packet to server
	// [SERVER] received packet from client: hello world
	// [CLIENT] received packet from server: message received!
	// [SERVER] client disconnected
}
