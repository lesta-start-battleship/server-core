package game

import (
	"errors"
)

type PlaceSubmarineCommand struct {
	ship *Ship
}

func NewPlaceSubmarineCommand(coord Coord, bearings bool) *PlaceSubmarineCommand {
	return &PlaceSubmarineCommand{
		ship: &Ship{
			ID: 11,
			Len:      3,
			Coords:   coord,
			Bearings: bearings,
			Health:   3,
			Decks:    make(map[Coord]bool),
		},
	}
}
func (c *PlaceSubmarineCommand) Apply(states *States) error {
	gs := states.PlayerState

	if !gs.isInside(c.ship.Coords) {
		return errors.New("starting coordinate is out of bounds")
	}

	// Проверка выхода за пределы по длине корабля
	for i := 0; i < c.ship.Len; i++ {
		x, y := c.ship.Coords.X, c.ship.Coords.Y
		if c.ship.Bearings == Vertical {
			y += i
		} else {
			x += i
		}
		if !gs.IsInside(x, y) {
			return errors.New("ship placement goes out of bounds")
		}
	}

	// Проверка на пересечение или соседство с другими кораблями
	for i := -1; i <= c.ship.Len; i++ {
		for j := -1; j <= 1; j++ {
			var x, y int
			if c.ship.Bearings == Vertical {
				x = c.ship.Coords.X + j
				y = c.ship.Coords.Y + i
			} else {
				x = c.ship.Coords.X + i
				y = c.ship.Coords.Y + j
			}
			if gs.IsInside(x, y) && gs.Field[x][y].ShipID != Empty {
				return errors.New("ship overlaps or is adjacent to another")
			}
		}
	}

	// Выдача ID
	/*if err := issueID(gs, c.ship); err != nil {
		return err
	}*/

	// Установка палуб и обновление поля
	for i := 0; i < c.ship.Len; i++ {
		var x, y int
		if c.ship.Bearings == Vertical {
			x = c.ship.Coords.X
			y = c.ship.Coords.Y + i
		} else {
			x = c.ship.Coords.X + i
			y = c.ship.Coords.Y
		}
		c.ship.Decks[Coord{X: x, Y: y}] = Whole
		gs.Field[x][y].ShipID = c.ship.ID
	}

	// Добавление корабля
	gs.Ships[c.ship.ID] = c.ship
	gs.NumShips += 1
	return nil
}

func (c *PlaceSubmarineCommand) Undo(states *States) {
	gs := states.PlayerState

	for coord := range c.ship.Decks {
		gs.Field[coord.X][coord.Y].ShipID = Empty
	}

	gs.Ships[c.ship.ID] = nil
	gs.NumShips -= 1
}

func (c *PlaceSubmarineCommand) GetDeckCoords() []Coord {
	coords := make([]Coord, c.ship.Len)
	for i := 0; i < c.ship.Len; i++ {
		x := c.ship.Coords.X
		y := c.ship.Coords.Y
		if c.ship.Bearings == Vertical {
			y += i
		} else {
			x += i
		}
		coords[i] = Coord{X: x, Y: y}
	}
	return coords
}

func (gs *GameState) FindSubmarineByCoord(coord Coord) *Ship {
	cell := gs.Field[coord.X][coord.Y]
	if cell.ShipID == Empty {
		return nil
	}
	ship := gs.Ships[cell.ShipID]
	if ship != nil && ship.Contains(coord) {
		return ship
	}
	return nil
}

func (c *PlaceSubmarineCommand) Ship() *Ship {
	return c.ship
}
