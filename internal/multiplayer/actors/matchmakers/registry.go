package matchmakers

type MatchmakerRegistry interface {
	Track(id string, room *Matchmaker)
	Find(id string) *Matchmaker
	Delete(id string)
}
