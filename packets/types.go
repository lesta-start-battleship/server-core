package packets

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Ship struct {
	ID       int     `json:"id,omitempty"`
	Len      int     `json:"len"`
	Coords   Coord   `json:"coords"`
	Bearings bool    `json:"bearings"`
	Health   int     `json:"health,omitempty"`
	Decks    []Coord `json:"decks,omitempty"`
}
