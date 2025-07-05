package players

type PlayerRegistry interface {
	Track(id string, room *Player)
	Find(id string) *Player
	Delete(id string)
}
