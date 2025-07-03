package strategies

import (
	"lesta-battleship/server-core/internal/multiplayer/actors"
	"lesta-battleship/server-core/pkg/packets"
)

type InSearch struct {
	Matchmaker   actors.Actor
}

func (s *InSearch) HandlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case packets.Disconnect:
		s.handleLeaveSearch(senderId, packet)
	}
}

func (s *InSearch) handleLeaveSearch(senderId string, packet packets.Disconnect) {
	s.Matchmaker.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packet})
}

func (s *InSearch) OnExit() {
	s.Matchmaker = nil
}

func (s InSearch) String() string {
	return "InSearch"
}
