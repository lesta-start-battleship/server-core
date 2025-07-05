package rooms

type RoomRegistry interface {
	Track(id string, room *Room)
	Find(id string) *Room
	Delete(id string)
}
