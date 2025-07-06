package game

const (
	Empty      int  = 0
	Open       int  = 1
	Close      int  = 2
	Vertical   bool = true
	Horizontal bool = false
	Hit        bool = true
	Whole      bool = false
)

type Coord struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Ship struct {
	ID int `json:"id"`
	// Type   ShipType `json:"type"`
	Len      int            `json:"len"`
	Coords   Coord          `json:"coords"`
	Bearings bool           `json:"bearings"` // ориентация
	Health   int            `json:"health"`
	Decks    map[Coord]bool `json:"decks"`
	// Coords []Coord  `json:"coords"`
}

type GameState struct {
	Field    [10][10]CellState
	Ships    []*Ship
	NumShips int
}

type States struct {
	PlayerState *GameState
	EnemyState  *GameState
}

type CellState struct {
	State  int `json:"state"`
	ShipID int `json:"shipid"`
}

const (
	Battleship int = 1
	Cruiser    int = 2
	Destroyer  int = 3
	Submarine  int = 4
)

func NewGameState() *GameState {
	return &GameState{
		Ships: make([]*Ship, 11),
	}
}

func (gs *GameState) isInside(c Coord) bool {
	return c.X >= 0 && c.X < 10 && c.Y >= 0 && c.Y < 10
}

func (gs *GameState) IsInside(x, y int) bool {
	return x >= 0 && x < 10 && y >= 0 && y < 10
}
