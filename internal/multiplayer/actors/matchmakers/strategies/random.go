package strategies

import (
	"lesta-battleship/server-core/internal/multiplayer/actors"
	"lesta-battleship/server-core/internal/multiplayer/actors/players"
	"lesta-battleship/server-core/pkg/packets"
)

type Random struct {
	Matchmaker actors.Matchmaker

	Hub        actors.Actor
	Queue map[string]*players.Player
}

func (s *Random) HandlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case packets.JoinSearch:
		s.handleJoinSearch(senderId, packet)
	case packets.Disconnect:
		s.handleDisconnect(senderId, packet)
	}
}

func (s *Random) handleJoinSearch(senderId string, packet packets.JoinSearch) {
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

func (s *Random) handleDisconnect(senderId string, packet packets.Disconnect) {
	s.Matchmaker.RemoveFromQueue(senderId)

	s.Hub.GetPacket(senderId, packets.Packet{SenderId: senderId, Body: packet})
}

func (s *Random) OnExit() {
	s.Hub = nil
	s.Matchmaker = nil
}

func (s *Random) String() string {
	return "Random"
}
