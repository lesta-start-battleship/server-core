package strategies

import (
	"github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors"
	"github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors/players"
	"github.com/lesta-battleship/server-core/pkg/matchmaking/packets"
)

type Random struct {
	Matchmaker actors.Matchmaker

	Hub   actors.Actor
	Queue map[string]*players.Player
}

func (s *Random) HandlePacket(senderId string, packet packets.Packet) {
	switch packet := packet.Body.(type) {
	case *packets.JoinSearch:
		s.handleJoinSearch(senderId, packet)
	case *packets.Disconnect:
		s.handleDisconnect(senderId, packet)
	}
}

func (s *Random) handleJoinSearch(senderId string, packet *packets.JoinSearch) {
	s.Matchmaker.AddToQueue(senderId)

	for secondId := range s.Queue {
		if senderId != secondId {
			room := s.Matchmaker.CreateRoom()

			s.Matchmaker.ConnectToRoom(room.Id(), senderId)
			s.Matchmaker.ConnectToRoom(room.Id(), secondId)

			s.Matchmaker.RemoveFromQueue(senderId)
			s.Matchmaker.RemoveFromQueue(secondId)

			return
		}
	}
}

func (s *Random) handleDisconnect(senderId string, packet *packets.Disconnect) {
	s.Matchmaker.RemoveFromQueue(senderId)

	s.Hub.GetPacket(senderId, packets.NewDisconnect(senderId))
}

func (s *Random) OnExit() {
	s.Hub = nil
	s.Matchmaker = nil
}

func (s *Random) String() string {
	return "Random"
}
