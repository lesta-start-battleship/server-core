package matchmakers

import (
	"fmt"
	"lesta-battleship/server-core/internal/multiplayer/actors/matchmakers/strategies"
	"lesta-battleship/server-core/pkg/packets"
)

type Strategy interface {
	HandlePacket(senderId string, packet packets.Packet)
	OnExit()

	fmt.Stringer
}

func SetRandom(matchmaker *Matchmaker) {
	matchmaker.ChangeStrategy(&strategies.Random{Matchmaker: matchmaker, Hub: matchmaker.hub, Queue: matchmaker.queue})
}

func SetRanked(matchmaker *Matchmaker) {
	matchmaker.ChangeStrategy(&strategies.Ranked{Matchmaker: matchmaker, Hub: matchmaker.hub, Queue: matchmaker.queue})
}
