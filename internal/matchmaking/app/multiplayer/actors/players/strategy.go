package players

import (
	"fmt"

	"github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors"
	"github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors/players/strategies"
	"github.com/lesta-battleship/server-core/pkg/matchmaking/packets"
)

type Strategy interface {
	HandlePacket(senderId string, packet packets.Packet)
	OnExit()

	fmt.Stringer
}

func SetInHub(player *Player, hub actors.Actor) {
	player.ChangeStrategy(&strategies.InHub{Player: player, Hub: hub})
}

func SetInRoom(player *Player, room actors.Actor) {
	player.ChangeStrategy(&strategies.InRoom{Player: player, Room: room})
}

func SetInSearch(player *Player, matchmaker actors.Matchmaker) {
	player.ChangeStrategy(&strategies.InSearch{Player: player, Matchmaker: matchmaker})
}
