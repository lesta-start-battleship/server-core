package game

type OpenCellCommand struct {
	Coords Coord // - получаем от юзера
	PrevState CellState
}

func NewOpenCellCommand(target Coord) *OpenCellCommand {
	return &OpenCellCommand{
		Coords: target,
	}
}

func (c *OpenCellCommand) Apply(gs *GameState) error {
	// проверка валидности координаты 
	if !gs.isInside(c.Coords) {
		return nil
	}
	c.PrevState = gs.Field[c.Coords.X][c.Coords.Y]
	gs.Field[c.Coords.X][c.Coords.Y].State = Open
	return nil
}

func (c *OpenCellCommand) Undo(gs *GameState) {
	if !gs.isInside(c.Coords) {
		return 
	}
	gs.Field[c.Coords.X][c.Coords.Y] = c.PrevState
}