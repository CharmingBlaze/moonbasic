package enet

// #include <enet/enet.h>
import "C"
import (
	"errors"
	"unsafe"
)

// PacketFlags are bit constants
type PacketFlags uint32

const (
	// PacketFlagReliable packets must be received by the target peer and resend attempts
	// should be made until the packet is delivered
	PacketFlagReliable PacketFlags = C.ENET_PACKET_FLAG_RELIABLE

	// PacketFlagUnsequenced packets will not be sequenced with other packets not supported
	// for reliable packets
	PacketFlagUnsequenced PacketFlags = C.ENET_PACKET_FLAG_UNSEQUENCED

	// PacketFlagNoAllocate packets will not allocate data, and user must supply it instead
	PacketFlagNoAllocate PacketFlags = C.ENET_PACKET_FLAG_NO_ALLOCATE

	// PacketFlagUnreliableFragment packets will be fragmented using unreliable (instead of
	// reliable) sends if it exceeds the MTU
	PacketFlagUnreliableFragment PacketFlags = C.ENET_PACKET_FLAG_UNRELIABLE_FRAGMENT

	// PacketFlagSent specifies whether the packet has been sent from all queues it has been
	// entered into
	PacketFlagSent PacketFlags = C.ENET_PACKET_FLAG_SENT
)

// Packet may be sent to or received from a peer
type Packet interface {
	Destroy()
	GetData() []byte
	GetFlags() PacketFlags
}

type enetPacket struct {
	cPacket *C.struct__ENetPacket
}

func (packet enetPacket) Destroy() {
	C.enet_packet_destroy(packet.cPacket)
}

func (packet enetPacket) GetData() []byte {
	return C.GoBytes(
		unsafe.Pointer(packet.cPacket.data),
		(C.int)(packet.cPacket.dataLength),
	)
}

func (packet enetPacket) GetFlags() PacketFlags {
	return (PacketFlags)(packet.cPacket.flags)
}

// NewPacket creates a new packet to send to peers
func NewPacket(data []byte, flags PacketFlags) (Packet, error) {
	buffer := C.CBytes(data)
	packet := C.enet_packet_create(
		buffer,
		(C.size_t)(len(data)),
		(C.enet_uint32)(flags),
	)
	C.free(buffer)

	if packet == nil {
		return nil, errors.New("unable to create packet")
	}

	return enetPacket{
		cPacket: packet,
	}, nil
}
