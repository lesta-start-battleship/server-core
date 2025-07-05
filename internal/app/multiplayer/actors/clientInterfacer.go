package actors

import (
	"lesta-battleship/server-core/pkg/packets"
)

type ClientInterfacer interface {
	Id() string
	ConnectTo(Actor)
	GetPacket(senderId string, packet packets.Packet)
	ReadPump()
	WritePump()
	Stop()
}
