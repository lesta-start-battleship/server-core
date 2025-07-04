package strategies

import (
	"lesta-battleship/server-core/internal/multiplayer/actors"
	"lesta-battleship/server-core/pkg/packets"
	"log"
)

type InRoom struct {
	Player actors.Actor
	Room   actors.Actor
}

func (s *InRoom) HandlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case *packets.PlayerMessage:
		s.handleBroadcast(senderId, packet)
	case *packets.Disconnect:
		s.handleDisconnect(senderId, packet)
	default:
		log.Printf("Player %q: Received incorrect packet %T from %q", s.Player.Id(), packet, senderId)
	}
}

func (s *InRoom) handleBroadcast(senderId string, packet *packets.PlayerMessage) {
	s.Room.GetPacket(senderId, packets.NewPlayerMessage(senderId, packet.Msg))
}

func (s *InRoom) handleDisconnect(senderId string, packet *packets.Disconnect) {
	s.Room.GetPacket(senderId, packets.NewDisconnect(senderId))
}

func (s *InRoom) OnExit() {
	s.Room = nil
}

func (s InRoom) String() string {
	return "InRoom"
}
