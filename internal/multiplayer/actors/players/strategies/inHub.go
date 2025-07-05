package strategies

import (
	"lesta-battleship/server-core/internal/multiplayer/actors"
	"lesta-battleship/server-core/pkg/packets"
	"log"
)

type InHub struct {
	Player actors.Actor
	Hub    actors.Actor
}

func (s *InHub) HandlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case *packets.Disconnect:
		s.handleDisconnect(senderId, packet)
	default:
		log.Printf("Player %q: Received incorrect packet %T from %q", s.Player.Id(), packet, senderId)
	}
}

func (s *InHub) handleDisconnect(senderId string, packet *packets.Disconnect) {
	s.Hub.GetPacket(senderId, packets.NewDisconnect(senderId))
}

func (s *InHub) OnExit() {
	s.Hub = nil
}

func (s InHub) String() string {
	return "InHub"
}
