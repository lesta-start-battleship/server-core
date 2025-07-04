package strategies

import (
	"lesta-battleship/server-core/internal/multiplayer/actors"
	"lesta-battleship/server-core/pkg/packets"
	"log"
)

type InSearch struct {
	Player     actors.Actor
	Matchmaker actors.Matchmaker
}

func (s *InSearch) HandlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case *packets.CreateRoom:
		s.handleCreateRoom(senderId, packet)
	case *packets.JoinRoom:
		s.handleJoinRoom(senderId, packet)
	case *packets.Disconnect:
		s.handleLeaveSearch(senderId, packet)
	default:
		log.Printf("Player %q: Received incorrect packet %T from %q", s.Player.Id(), packet, senderId)
	}
}

func (s *InSearch) handleCreateRoom(senderId string, packet *packets.CreateRoom) {
	s.Matchmaker.GetPacket(senderId, packets.NewCreateRoom(senderId))
}

func (s *InSearch) handleJoinRoom(senderId string, packet *packets.JoinRoom) {
	s.Matchmaker.GetPacket(senderId, packets.NewJoinRoom(senderId, packet.Id))
}

func (s *InSearch) handleLeaveSearch(senderId string, packet *packets.Disconnect) {
	s.Matchmaker.GetPacket(senderId, packets.NewDisconnect(senderId))
}

func (s *InSearch) OnExit() {
	s.Matchmaker = nil
}

func (s *InSearch) String() string {
	return "InSearch"
}
