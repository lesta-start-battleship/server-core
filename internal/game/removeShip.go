package game

import "errors"

type RemoveShipCommand struct {
	Coords Coord // - получаем от юзера
	ship *Ship // - инициализирует код и использует для бекапа
}

func (c *RemoveShipCommand) Apply(gs *GameState) error {
	// проверка валидности координаты
	if !gs.isInside(c.Coords) {
		return errors.New("out of bounds")
	}
	if cellState := gs.Field[c.Coords.X][c.Coords.Y]; cellState.ShipID != Empty  {
		c.ship = gs.Ships[cellState.ShipID]
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
	} else {
		return errors.New("empty cell")
	}
	return nil
}

func (c *RemoveShipCommand) Undo(gs *GameState) {
	// добавление корабля на карту
	gs.Ships[c.ship.ID] = c.ship

	// изменение карты
	for x := c.ship.Coords.X - 1; x <= c.ship.Coords.X + c.ship.Len; x++ {
		mx, my := x, c.ship.Coords.Y
		if c.ship.Bearings == Vertical {
			mx, my = my, mx
		}
		// c.ship.Decks[Coord{mx, my}] = Whole
		gs.Field[mx][my].ShipID = c.ship.ID
	}

	// изменения юзера
	gs.NumShips += 1
}
