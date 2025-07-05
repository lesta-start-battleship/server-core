package matchmakers

import (
	"fmt"
	"lesta-battleship/server-core/internal/app/multiplayer/actors/matchmakers/strategies"
	"lesta-battleship/server-core/pkg/packets"
)

type Strategy interface {
	HandlePacket(senderId string, packet packets.Packet)
	OnExit()

	fmt.Stringer
}

var strategiesMap = map[MatchType]func(*Matchmaker){
	RandomMatch: SetRandom,
	RankedMatch: SetRanked,
	CustomMatch: SetCustom,
}

func SetStrategy(matchmaker *Matchmaker, matchType MatchType) {
	strategiesMap[matchType](matchmaker)
}

func SetRandom(matchmaker *Matchmaker) {
	matchmaker.ChangeStrategy(&strategies.Random{Matchmaker: matchmaker, Hub: matchmaker.hub, Queue: matchmaker.queue})
}

func SetRanked(matchmaker *Matchmaker) {
	matchmaker.ChangeStrategy(&strategies.Ranked{Matchmaker: matchmaker, Hub: matchmaker.hub, Queue: matchmaker.queue})
}

func SetCustom(matchmaker *Matchmaker) {
	matchmaker.ChangeStrategy(&strategies.Custom{Matchmaker: matchmaker, Hub: matchmaker.hub})
}
