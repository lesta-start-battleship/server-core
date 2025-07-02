package rooms

type RoomRegistry interface {
	Track(id string, room *Room)
	Find(id string) *Room
	Rooms() map[string]*Room
	Delete(id string)
}
