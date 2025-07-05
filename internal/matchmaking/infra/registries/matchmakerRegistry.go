package registries

import "github.com/lesta-battleship/server-core/internal/matchmaking/app/multiplayer/actors/matchmakers"

type MatchmakerRegistry struct {
	players map[string]*matchmakers.Matchmaker
}

func NewMatchmakerRegistry() *MatchmakerRegistry {
	return &MatchmakerRegistry{players: make(map[string]*matchmakers.Matchmaker)}
}

func (r *MatchmakerRegistry) Track(id string, player *matchmakers.Matchmaker) {
	r.players[id] = player
}

func (r *MatchmakerRegistry) Find(id string) *matchmakers.Matchmaker {
	return r.players[id]
}

func (r *MatchmakerRegistry) Matchmakers() map[string]*matchmakers.Matchmaker {
	return r.players
}

func (r *MatchmakerRegistry) Delete(id string) {
	delete(r.players, id)
}
