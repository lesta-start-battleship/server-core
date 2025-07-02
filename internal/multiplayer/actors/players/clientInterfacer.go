package players

import (
	"lesta-battleship/server-core/pkg/packets"
)

type ClientInterfacer interface {
	ConnectTo(*Player)
	GetPacket(senderId string, packet packets.Packet)
	ReadPump()
	WritePump()
	Stop()
}
