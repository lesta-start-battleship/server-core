package strategies

import (
	"lesta-battleship/server-core/internal/multiplayer/actors"
	"lesta-battleship/server-core/internal/multiplayer/actors/players"
	"lesta-battleship/server-core/pkg/packets"
	"log"
)

type Ranked struct {
	Matchmaker actors.Matchmaker

	Hub   actors.Actor
	Queue map[string]*players.Player
}

func (s *Ranked) HandlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case packets.JoinSearch:
		s.handleJoinSearch(senderId, packet)
	case packets.Disconnect:
		s.handleDisconnect(senderId, packet)
	default:
		log.Printf("Matchmaker %q: Got incorrect packet %t from %q", s.Matchmaker.Id(), packet, senderId)
	}
}

func (s *Ranked) handleJoinSearch(senderId string, packet packets.JoinSearch) {
	s.Matchmaker.AddToQueue(senderId)

	for secondId := range s.Queue {
		if senderId != secondId {
			room := s.Matchmaker.CreateRoom()

			room.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packets.ConnectPlayer{Id: senderId}})
			room.GetPacket(secondId, packets.Packet{SenderId: secondId, Body: packets.ConnectPlayer{Id: secondId}})

			s.Matchmaker.RemoveFromQueue(senderId)
			s.Matchmaker.RemoveFromQueue(secondId)

			return
		}
	}
}

func (s *Ranked) handleDisconnect(senderId string, packet packets.Disconnect) {
	s.Matchmaker.RemoveFromQueue(senderId)

	s.Hub.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packet})
}

func (s *Ranked) OnExit() {
	s.Hub = nil
	s.Matchmaker = nil
}

func (s *Ranked) String() string {
	return "Ranked"
}
