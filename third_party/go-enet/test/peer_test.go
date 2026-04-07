package enet_test

import (
	"fmt"
	"github.com/codecat/go-enet"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

func TestPeerData(t *testing.T) {
	testData := []byte{0x1, 0x2, 0x3}

	// peer is connected to our server.
	// events will produce events as the server receives them.
	peer, events := createServerClient(t)

	// Wait for the server to respond with a connection.
	ev := <-events
	if data := ev.GetPeer().GetData(); data != nil {
		t.Fatalf("did not expect new peer to have data set, but has %x", data)
	}

	// Set some data against our peer and immediately check it's there.
	ev.GetPeer().SetData(testData)
	assertPeerData(t, ev.GetPeer(), testData, "immediate after set")

	// Send a message to the server.
	if err := peer.SendString("testmessage", 0, enet.PacketFlagReliable); err != nil {
		t.Fatal(err)
	}

	// Wait for the server to receive this message, then check the
	// server-side peer associated with this event has the data
	// we set previously.
	ev = <-events
	assertPeerData(t, ev.GetPeer(), testData, "on packet received")

	t.Run("clear-data", func(t *testing.T) {
		ev.GetPeer().SetData(nil)
		assertPeerData(t, ev.GetPeer(), nil, "nil set")
	})

	t.Run("empty-slice", func(t *testing.T) {
		ev.GetPeer().SetData([]byte{})
		assertPeerData(t, ev.GetPeer(), []byte{}, "empty set")
	})

	// Check that our data stored in C survives garbage collection.
	t.Run("survives-gc", func(t *testing.T) {
		ev.GetPeer().SetData([]byte{1, 2, 3})
		runtime.GC()
		assertPeerData(t, ev.GetPeer(), []byte{1, 2, 3}, "after GC")
	})

	// Sniffs for a potential memory leak in our set data implementation.
	// We expect SetData to clear whatever C memory was used previously.
	// This may end up being a flaky test, but will keep it in for now to
	// build confidence in our implementation.
	t.Run("memory-leaks", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.Skipf("skipping mem leak test as not running on a linux host")
		}

		noOfIncreases := 0
		last := currentMemory(t)

		// Assign a large string (10MB) to data over and over again, checking
		// for continuous increases in mem usage, with some threshold.
		for i := 0; i < 99; i++ {
			ev.GetPeer().SetData([]byte(strings.Repeat("x", 1024*1024*10)))

			// Detect a memory leak by checking if we're using more than 1MB last than the
			// previous for too many iterations.
			now := currentMemory(t)

			if now-last > 1024 {
				noOfIncreases++
			} else {
				// If it's not an increase, reset our counter.
				noOfIncreases = 0
			}

			// If we reach a threshold of 5 continuous increases, consider this a leak.
			if noOfIncreases > 5 {
				t.Fatal("potential memory leak detected")
			}

			last = now
		}
	})
}

func assertPeerData(t testing.TB, peer enet.Peer, expected []byte, msg string) {
	t.Helper()

	actual := peer.GetData()

	if (actual == nil) != (expected == nil) {
		t.Fatalf("%s: expected peer data to be present? %t vs actual: %t", msg, expected != nil, actual != nil)
	}

	if len(actual) != len(expected) {
		t.Fatalf("%s: expected peer data to have len %d vs actual %d", msg, len(expected), len(actual))
	}

	if string(actual) != string(expected) {
		t.Fatalf("%s: expected peer data to be %v, but it was %v", msg, expected, actual)
	}
}

// createServerClient creates a dummy enet server and client. The returned
// peer can be used to send messages to the server, and the blocking events
// channel returned will be given each event as the server picks it up.
func createServerClient(t *testing.T) (clientConn enet.Peer, serverEvents <-chan enet.Event) {
	port := getFreePort()

	done := make(chan bool, 0)
	events := make(chan enet.Event)

	t.Cleanup(func() {
		// Kill our background service routines for client & server.
		close(done)
	})

	// Create a server and continuously service it, exposing any captured events.
	server, err := enet.NewHost(enet.NewListenAddress(port), 10, 1, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for true {
			select {
			case <-done:
				return
			default:
				ev := server.Service(0)

				// Pass any event out to our channel. This will block
				// until a test consumes it.
				if ev.GetType() != enet.EventNone {
					events <- ev
				}
			}
		}
	}()

	// Create a client and continuously service it in the background.
	client, err := enet.NewHost(nil, 1, 1, 0, 0)
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for true {
			select {
			case <-done:
				return
			default:
				client.Service(0)
			}
		}
	}()

	// Connect to our server.
	peer, err := client.Connect(enet.NewAddress("localhost", port), 1, 0)
	if err != nil {
		t.Fatal(err)
	}

	return peer, events
}

var port uint16 = 49152

// getFreePort returns a unique private port. Note this doesn't guarantee
// it's free, but should be good enough from within docker tests.
func getFreePort() uint16 {
	port++
	return port
}

// currentMemory returns the memory usage of the current process according to
// the OS. This uses linux's proc FS to give a rough estimate based on VmSize.
// Note we don't want to use runtime.MemStats here as we're looking for the
// total memory (including C allocations).
func currentMemory(t testing.TB) int {
	// GC to give a stable measure.
	runtime.GC()

	pmapCmd := fmt.Sprintf("pmap %d | grep total | grep -Eo '[0-9]+'", os.Getpid())
	cmd := exec.Command("bash", "-c", pmapCmd)

	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("failed to run memory check command: %s", err)
	}

	kb, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		t.Fatalf("failed converting memory output provided by pmap to int: %s", err)
	}

	return kb
}
