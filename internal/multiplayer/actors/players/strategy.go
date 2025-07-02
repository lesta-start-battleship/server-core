package players

import (
	"fmt"
	"lesta-battleship/server-core/internal/multiplayer/actors"
	"lesta-battleship/server-core/internal/multiplayer/actors/players/states"
	"lesta-battleship/server-core/pkg/packets"
)

type Strategy interface {
	HandlePacket(senderId string, packet packets.Packet)
	OnExit()

	fmt.Stringer
}

func SetInRoom(player *Player, room actors.Actor) {
	player.ChangeStrategy(&states.InRoom{Room: room})
}

func SetInSearch(player *Player, matchmaker actors.Actor) {
	player.ChangeStrategy(&states.InSearch{Matchmaker: matchmaker})
}
