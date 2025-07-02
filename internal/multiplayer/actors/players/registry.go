package players

type PlayerRegistry interface {
	Track(id string, room *Player)
	Find(id string) *Player
	Players() map[string]*Player
	Delete(id string)
}
