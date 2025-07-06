package game

import "errors"

type RemoveShipCommand struct {
	Coords Coord // - получаем от юзера
	ship   *Ship // - инициализирует код и использует для бекапа
}

func NewRemoveShipCommand(target Coord) *RemoveShipCommand {
	return &RemoveShipCommand{
		Coords: target,
	}
}

func (c *RemoveShipCommand) Apply(states *States) error {
	gs := states.PlayerState

	if !gs.isInside(c.Coords) {
		return errors.New("out of bounds")
	}

	cell := gs.Field[c.Coords.X][c.Coords.Y]
	if cell.ShipID == Empty {
		return errors.New("empty cell")
	}

	ship := gs.Ships[cell.ShipID]
	if ship == nil {
		return errors.New("ship not found")
	}

	c.ship = ship // backup

	// Удаляем палубы
	for i := 0; i < ship.Len; i++ {
		x := ship.Coords.X
		y := ship.Coords.Y
		if ship.Bearings == Vertical {
			y += i
		} else {
			x += i
		}
		gs.Field[x][y].ShipID = Empty
	}

	// Удаляем корабль
	gs.Ships[ship.ID] = nil
	gs.NumShips--

	return nil
}

func (c *RemoveShipCommand) Undo(states *States) {
	gs := states.PlayerState

	gs.Ships[c.ship.ID] = c.ship
	for i := 0; i < c.ship.Len; i++ {
		x := c.ship.Coords.X
		y := c.ship.Coords.Y
		if c.ship.Bearings == Vertical {
			y += i
		} else {
			x += i
		}
		gs.Field[x][y].ShipID = c.ship.ID
	}
	gs.NumShips++
}

func (c *RemoveShipCommand) GetDeckCoords() []Coord {
	if c.ship == nil {
		return nil
	}
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
