package game

import (
	"errors"
	//"fmt"
)

type PlaceShipCommand struct {
	ship *Ship
}

func NewPlaceShipCommand(len int, coord Coord, bearings bool) PlaceShipCommand {
	return PlaceShipCommand{
		ship: &Ship{
			Len: len,
			Coords: coord,
			Bearings: bearings,
			Health: len,
			Decks: make(map[Coord]bool),
		},
	}
}
func (c *PlaceShipCommand) Apply(states *States) error {
	gs := states.PlayerState
	// проверка валидности координаты
	if !gs.isInside(c.ship.Coords) {
		return errors.New("out of bounds")
	}
	// проверка валидности размера корабля
	if err := issueID(gs, c.ship); err != nil {
		return err
	}
	// проверка валидности места
	for x := c.ship.Coords.X - 1; x <= c.ship.Coords.X + c.ship.Len + 1; x++ {
		for y := c.ship.Coords.Y - 1; y <= c.ship.Coords.X + 1; y++ {
			mx, my := x, y
			if c.ship.Bearings == Vertical {
				mx, my = my, mx
			}
			
			if !gs.IsInside(mx, my) || gs.Field[mx][my].ShipID == Empty  {
				return errors.New("bad place")
			}
		}
	}

	// добавление корабля на карту
	gs.Ships[c.ship.ID] = c.ship

	// выдача координат кораблю, изменение карты
	for x := c.ship.Coords.X - 1; x <= c.ship.Coords.X + c.ship.Len; x++ {
		mx, my := x, c.ship.Coords.Y
		if c.ship.Bearings == Vertical {
			mx, my = my, mx
		}
		c.ship.Decks[Coord{mx, my}] = Whole
		gs.Field[mx][my].ShipID = c.ship.ID
	}

	// изменения юзера
	gs.NumShips += 1
	return nil
}

func (c *PlaceShipCommand) Undo(states *States) error {
	gs := states.PlayerState
	// изменение карты
	for x := c.ship.Coords.X - 1; x <= c.ship.Coords.X + c.ship.Len; x++ {
		mx, my := x, c.ship.Coords.Y
		if c.ship.Bearings == Vertical {
			mx, my = my, mx
		}
		gs.Field[mx][my].ShipID = Empty
	}
	// удаление корабля
	gs.Ships[c.ship.ID] = nil

	// изменения юзера
	gs.NumShips -= 1

	return nil
}

func issueID(gs *GameState, ship *Ship) error {
	switch ship.Len {
	case Battleship:
		for i := 1; i < 5; i++ {
			if gs.Ships[i] == nil {
				//gs.Ships[i] = ship
				ship.ID = i
			}
		}  
	case Cruiser:
		for i := 5; i < 8; i++ {
			if gs.Ships[i] == nil {
				//gs.Ships[i] = ship
				ship.ID = i
			}
		}
	case Destroyer:
		for i := 8; i < 10; i++ {
			if gs.Ships[i] == nil {
				//gs.Ships[i] = ship
				ship.ID = i
			}
		}
	case Submarine:
		for i := 10; i < 11; i++ {
			if gs.Ships[i] == nil {
				//gs.Ships[i] = ship
				ship.ID = i
			}
		}
	default:
		return errors.New("incorrect length of ship")
	}
	return nil
}