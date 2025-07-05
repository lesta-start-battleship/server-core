package strategies

import (
	"log"

	"github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors"
	"github.com/lesta-battleship/server-core/pkg/matchmaking/packets"
)

type Custom struct {
	Matchmaker actors.Matchmaker

	Hub actors.Actor
}

func (s *Custom) HandlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case *packets.JoinSearch:
		s.handleJoinSearch(senderId, packet)
	case *packets.CreateRoom:
		s.handleCreateRoom(senderId, packet)
	case *packets.JoinRoom:
		s.handleJoinRoom(senderId, packet)
	case *packets.Disconnect:
		s.handleDisconnect(senderId, packet)
	default:
		log.Printf("Matchmaker %q: Got incorrect packet %t from %q", s.Matchmaker.Id(), packet, senderId)
	}
}

func (s *Custom) handleJoinSearch(senderId string, packet *packets.JoinSearch) {
	s.Matchmaker.AddToQueue(senderId)
}

func (s *Custom) handleCreateRoom(senderId string, packet *packets.CreateRoom) {
	s.Matchmaker.RemoveFromQueue(senderId)

	room := s.Matchmaker.CreateRoom()

	s.Matchmaker.ConnectToRoom(room.Id(), senderId)
}

func (s *Custom) handleJoinRoom(senderId string, packet *packets.JoinRoom) {
	s.Matchmaker.RemoveFromQueue(senderId)

	roomId := packet.Id

	s.Matchmaker.ConnectToRoom(roomId, senderId)
}

func (s *Custom) handleDisconnect(senderId string, packet *packets.Disconnect) {
	s.Matchmaker.RemoveFromQueue(senderId)

	s.Hub.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packet})
}

func (s *Custom) OnExit() {
	s.Hub = nil
	s.Matchmaker = nil
}

func (s *Custom) String() string {
	return "Custom"
}
