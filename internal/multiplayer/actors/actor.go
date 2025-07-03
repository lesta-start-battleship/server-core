package actors

import "lesta-battleship/server-core/pkg/packets"

type Actor interface {
	Id() string
	GetPacket(senderId string, packet packets.Packet)
	Start()
	Stop()
}

type Matchmaker interface {
	Actor
	CreateRoom() Actor
	AddToQueue(playerId string)
	RemoveFromQueue(playerId string)
}
