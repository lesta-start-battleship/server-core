package game

type OpenCellCommand struct {
	Coords    Coord // - получаем от юзера
	PrevState CellState
	ShipFound bool
}

func NewOpenCellCommand(target Coord) *OpenCellCommand {
	return &OpenCellCommand{
		Coords: target,
	}
}

func (c *OpenCellCommand) Apply(states *States) error {
	gs := states.EnemyState
	// проверка валидности координаты
	if !gs.isInside(c.Coords) {
		return nil
	}
	c.PrevState = gs.Field[c.Coords.X][c.Coords.Y]
	gs.Field[c.Coords.X][c.Coords.Y].State = Open

	if c.PrevState.ShipID != Empty {
		c.ShipFound = true
	}
	return nil
}

func (c *OpenCellCommand) Undo(states *States) {
	gs := states.EnemyState
	if !gs.isInside(c.Coords) {
		return
	}
	gs.Field[c.Coords.X][c.Coords.Y] = c.PrevState
}