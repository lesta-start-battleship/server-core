package actors

import "lesta-battleship/server-core/pkg/packets"

type Actor interface {
	GetPacket(senderId string, packet packets.Packet)
	Start()
	Stop()
}
