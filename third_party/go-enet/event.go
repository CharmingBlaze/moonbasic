package enet

// #include <enet/enet.h>
import "C"

// EventType is a type of event
type EventType int

const (
	// EventNone means that no event occurred within the specified time limit
	EventNone EventType = iota

	// EventConnect means that a connection request initiated by Host.Connect has completed
	// The peer field contains the peer which successfully connected
	EventConnect

	// EventDisconnect means that a peer has disconnected. This event is generated on a
	// successful completion of a disconnect initiated by Peer.Disconnect, if a peer has
	// timed out, or if a connection request intialized by Host.Connect has timed out. The
	// peer field contains the peer which disconnected. The data field contains user supplied
	// data describing the disconnection, or 0, if none is available.
	EventDisconnect

	// EventReceive means that a packet has been received from a peer. The peer field
	// specifies the peer which sent the packet. The channelID field specifies the channel
	// number upon which the packet was received. The packet field contains the packet that
	// was received; this packet must be destroyed with Packet.Destroy after use.
	EventReceive
)

// Event as returned by Host.Service()
type Event interface {
	GetType() EventType
	GetPeer() Peer
	GetChannelID() uint8
	GetData() uint32
	GetPacket() Packet
}

type enetEvent struct {
	cEvent C.struct__ENetEvent
}

func (event *enetEvent) GetType() EventType {
	return (EventType)(event.cEvent._type)
}

func (event *enetEvent) GetPeer() Peer {
	return enetPeer{
		cPeer: event.cEvent.peer,
	}
}

func (event *enetEvent) GetChannelID() uint8 {
	return (uint8)(event.cEvent.channelID)
}

func (event *enetEvent) GetData() uint32 {
	return (uint32)(event.cEvent.data)
}

func (event *enetEvent) GetPacket() Packet {
	return enetPacket{
		cPacket: event.cEvent.packet,
	}
}
