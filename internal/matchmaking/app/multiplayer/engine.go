package multiplayer

import (
	"log"

	"github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors"
	"github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors/hubs"
	"github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors/matchmakers"
	"github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors/players"
	"github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors/rooms"
	"github.com/lesta-battleship/server-core/internal/matchmaking/infra"
	"github.com/lesta-battleship/server-core/pkg/matchmaking/packets"
)

type Engine struct {
	matchmakerRegistry matchmakers.MatchmakerRegistry
	roomRegistry       rooms.RoomRegistry
	playerRegistry     players.PlayerRegistry

	hub actors.Actor
}

func NewEngine(
	matchmakerRegistry matchmakers.MatchmakerRegistry,
	roomRegistry rooms.RoomRegistry,
	playerRegistry players.PlayerRegistry,
) *Engine {
	return &Engine{
		matchmakerRegistry: matchmakerRegistry,
		roomRegistry:       roomRegistry,
		playerRegistry:     playerRegistry,

		hub: nil,
	}
}

func (e *Engine) CreateHub() *hubs.Hub {
	hub := hubs.NewHub(e.matchmakerRegistry, e.roomRegistry, e.playerRegistry)
	go hub.Start()
	e.hub = hub

	return hub
}

func (e *Engine) CreatePlayer(interfacer actors.ClientInterfacer) *players.Player {
	player := players.NewPlayer(interfacer.Id(), interfacer)
	players.SetInHub(player, e.hub)
	go player.Start()
	e.playerRegistry.Track(player.Id(), player)

	log.Printf("Engine: Created player %q", player.Id())

	return player
}

func (e *Engine) CreateMatchmaker(matchType matchmakers.MatchType) *matchmakers.Matchmaker {
	matchmaker := matchmakers.NewMatchmaker(
		infra.GenerateId(),
		e.playerRegistry,
		e.roomRegistry,
		e.hub)
	matchmakers.SetStrategy(matchmaker, matchType)
	go matchmaker.Start()
	e.matchmakerRegistry.Track(matchType.String(), matchmaker)

	log.Printf("Engine: Created matchmaker %q", matchmaker.Id())

	return matchmaker
}

func (e *Engine) SendToMatchmaking(player *players.Player, matchType matchmakers.MatchType) {
	e.hub.GetPacket(player.Id(), packets.NewJoinSearch(player.Id(), matchType.String()))
}
