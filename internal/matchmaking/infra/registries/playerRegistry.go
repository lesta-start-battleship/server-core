package registries

import "github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors/players"

type PlayerRegistry struct {
	players map[string]*players.Player
}

func NewPlayerRegistry() *PlayerRegistry {
	return &PlayerRegistry{players: make(map[string]*players.Player)}
}

func (r *PlayerRegistry) Track(id string, player *players.Player) {
	r.players[id] = player
}

func (r *PlayerRegistry) Find(id string) *players.Player {
	return r.players[id]
}

func (r *PlayerRegistry) Players() map[string]*players.Player {
	return r.players
}

func (r *PlayerRegistry) Delete(id string) {
	delete(r.players, id)
}
