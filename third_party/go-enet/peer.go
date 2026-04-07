package enet

// #include <enet/enet.h>
import "C"
import (
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"
)

// Peer is a peer which data packets may be sent or received from
type Peer interface {
	GetAddress() Address
	GetConnectId() uint
	GetIncomingPeerId() uint16
	GetIncomingSessionId() uint8

	Disconnect(data uint32)
	DisconnectNow(data uint32)
	DisconnectLater(data uint32)

	SendBytes(data []byte, channel uint8, flags PacketFlags) error
	SendString(str string, channel uint8, flags PacketFlags) error
	SendPacket(packet Packet, channel uint8) error

	// SetData sets an arbitrary value against a peer. This is useful to attach some
	// application-specific data for future use, such as an identifier.
	//
	// http://enet.bespin.org/structENetPeer.html#a1873959810db7ac7a02da90469ee384e
	//
	// Note that due to the way the enet library works, if using this you are
	// responsible for clearing this data when the peer is finished with.
	// SetData(nil) will free underlying memory and avoid any leaks.
	//
	// See http://enet.bespin.org/Tutorial.html#ManageHost for an example of this
	// in the underlying library.
	SetData(data []byte)

	// GetData returns an application-specific value that's been set
	// against this peer. This returns nil if no data has been set.
	//
	// http://enet.bespin.org/structENetPeer.html#a1873959810db7ac7a02da90469ee384e
	GetData() []byte

	//rtt
	GetRoundTripTime() uint32
	GetLastRoundTripTime() uint32
	GetLowestRoundTripTime() uint32
	GetRoundTripTimeVariance() uint32
	GetLastRoundTripTimeVariance() uint32
	GetHighestRoundTripTimeVariance() uint32
}

type enetPeer struct {
	cPeer *C.struct__ENetPeer
}

func (peer enetPeer) GetAddress() Address {
	return &enetAddress{
		cAddr: peer.cPeer.address,
	}
}

func (peer enetPeer) GetConnectId() uint {
	return uint(peer.cPeer.connectID)
}

func (peer enetPeer) GetIncomingPeerId() uint16 {
	return uint16(peer.cPeer.incomingPeerID)
}

func (peer enetPeer) GetIncomingSessionId() uint8 {
	return uint8(peer.cPeer.incomingSessionID)
}

func (peer enetPeer) Disconnect(data uint32) {
	C.enet_peer_disconnect(
		peer.cPeer,
		(C.enet_uint32)(data),
	)
}

func (peer enetPeer) DisconnectNow(data uint32) {
	C.enet_peer_disconnect_now(
		peer.cPeer,
		(C.enet_uint32)(data),
	)
}

func (peer enetPeer) DisconnectLater(data uint32) {
	C.enet_peer_disconnect_later(
		peer.cPeer,
		(C.enet_uint32)(data),
	)
}

func (peer enetPeer) SendBytes(data []byte, channel uint8, flags PacketFlags) error {
	packet, err := NewPacket(data, flags)
	if err != nil {
		return err
	}
	return peer.SendPacket(packet, channel)
}

func (peer enetPeer) SendString(str string, channel uint8, flags PacketFlags) error {
	packet, err := NewPacket([]byte(str), flags)
	if err != nil {
		return err
	}
	return peer.SendPacket(packet, channel)
}

func (peer enetPeer) SendPacket(packet Packet, channel uint8) error {
	C.enet_peer_send(
		peer.cPeer,
		(C.enet_uint8)(channel),
		packet.(enetPacket).cPacket,
	)
	return nil
}

func (peer enetPeer) SetData(data []byte) {
	if len(data) > math.MaxUint32 {
		panic(fmt.Sprintf("maximum peer data length is uint32 (%d)", math.MaxUint32))
	}

	// Free any data that was previously stored against this peer.
	existing := unsafe.Pointer(peer.cPeer.data)
	if existing != nil {
		C.free(existing)
	}

	// If nil, set this explicitly.
	if data == nil {
		peer.cPeer.data = nil
		return
	}

	// First 4 bytes stores how many bytes we have. This is so we can C.GoBytes when
	// retrieving which requires a byte length to read.
	b := make([]byte, len(data)+4)
	binary.LittleEndian.PutUint32(b, uint32(len(data)))
	// Join this header + data in to a contiguous slice
	copy(b[4:], data)
	// And write it out to C memory, storing our pointer.
	peer.cPeer.data = unsafe.Pointer(C.CBytes(b))
}

func (peer enetPeer) GetData() []byte {
	ptr := unsafe.Pointer(peer.cPeer.data)

	if ptr == nil {
		return nil
	}

	// First 4 bytes are the bytes length.
	header := []byte{
		*(*byte)(unsafe.Add(ptr, 0)),
		*(*byte)(unsafe.Add(ptr, 1)),
		*(*byte)(unsafe.Add(ptr, 2)),
		*(*byte)(unsafe.Add(ptr, 3)),
	}

	return []byte(C.GoBytes(
		// Take from the start of the data.
		unsafe.Add(ptr, 4),
		// As many bytes as were indicated in the header.
		C.int(binary.LittleEndian.Uint32(header)),
	))
}

// rtt impl
func (peer enetPeer) GetRoundTripTime() uint32 {
	return uint32(peer.cPeer.roundTripTime)
}

func (peer enetPeer) GetLastRoundTripTime() uint32 {
	return uint32(peer.cPeer.lastRoundTripTime)
}

func (peer enetPeer) GetLowestRoundTripTime() uint32 {
	return uint32(peer.cPeer.lowestRoundTripTime)
}

func (peer enetPeer) GetRoundTripTimeVariance() uint32 {
	return uint32(peer.cPeer.roundTripTimeVariance)
}

func (peer enetPeer) GetLastRoundTripTimeVariance() uint32 {
	return uint32(peer.cPeer.lastRoundTripTimeVariance)
}

func (peer enetPeer) GetHighestRoundTripTimeVariance() uint32 {
	return uint32(peer.cPeer.highestRoundTripTimeVariance)
}
