package players

import (
	"fmt"
	"lesta-battleship/server-core/internal/multiplayer/actors"
	"lesta-battleship/server-core/internal/multiplayer/actors/players/strategies"
	"lesta-battleship/server-core/pkg/packets"
)

type Strategy interface {
	HandlePacket(senderId string, packet packets.Packet)
	OnExit()

	fmt.Stringer
}

func SetInRoom(player *Player, room actors.Actor) {
	player.ChangeStrategy(&strategies.InRoom{Player: player, Room: room})
}

func SetInSearch(player *Player, matchmaker actors.Matchmaker) {
	player.ChangeStrategy(&strategies.InSearch{Player: player, Matchmaker: matchmaker})
}
