package items

import (
	"testing"

	"github.com/lesta-battleship/server-core/internal/game"
)

func makeTestState() *game.States {
	return &game.States{
		PlayerState: &game.GameState{
			Field: [10][10]game.CellState{},
			Ships: make([]*game.Ship, 12),
		},
		EnemyState: &game.GameState{
			Field: [10][10]game.CellState{},
			Ships: make([]*game.Ship, 12),
		},
	}
}

func TestRunScript_Nakhimov(t *testing.T) {
	script := `{"input": "Координаты выбранной клетки (x, y)", "actions": [{ "OPEN_CELL": {"x": "x", "y": "y"} }, { "OPEN_CELL": {"x": "x", "y": "y+1"} }, { "OPEN_CELL": {"x": "x+1", "y": "y"} }, { "OPEN_CELL": {"x": "x", "y": "y-1"} }, { "OPEN_CELL":{"x": "x-1", "y": "y"} }, { "SET_NAHIMOV_STATUS": {"status": "1"} }, { "END_PLAYER_ACTION": "None" }]}`
	input := ItemInput{X: 4, Y: 4}
	state := makeTestState()
	effects, err := RunScript(script, state, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantCoords := map[[2]int]bool{
		{4, 4}: true,
		{4, 5}: true,
		{5, 4}: true,
		{4, 3}: true,
		{3, 4}: true,
	}
	gotCoords := map[[2]int]bool{}
	for _, eff := range effects {
		if eff.Type != "open" {
			continue
		}
		for _, c := range eff.Coords {
			gotCoords[[2]int{c.X, c.Y}] = true
		}
	}
	for coord := range wantCoords {
		if !gotCoords[coord] {
			t.Errorf("missing coord: %v", coord)
		}
	}
}

func TestRunScript_Kon(t *testing.T) {
	script := `{"input": "Координаты выбранной клетки (x, y) и направление 'direction'", "actions": [{ "SWITCH_CASE": {"1": [ {"Name": "OPEN_CELL", "Args": {"x": "$x", "y": "$y"}}, {"Name": "OPEN_CELL", "Args": {"x": "$x", "y": "$y+1"}}, {"Name": "OPEN_CELL", "Args": {"x": "$x", "y": "$y+2"}}, {"Name": "OPEN_CELL", "Args": {"x": "$x-1", "y": "$y+2"}} ], "2": [ {"Name": "OPEN_CELL", "Args": {"x": "$x", "y": "$y"}}, {"Name": "OPEN_CELL", "Args": {"x": "$x", "y": "$y+1"}}, {"Name": "OPEN_CELL", "Args": {"x": "$x", "y": "$y+2"}}, {"Name": "OPEN_CELL", "Args": {"x": "$x+1", "y": "$y+2"}} ]} }, { "Name": "END_PLAYER_ACTION", "Args": "None" }]}`
	input := ItemInput{X: 2, Y: 2, Direction: 1}
	state := makeTestState()
	effects, err := RunScript(script, state, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	wantCoords := map[[2]int]bool{
		{2, 2}: true,
		{2, 3}: true,
		{2, 4}: true,
		{1, 4}: true,
	}
	gotCoords := map[[2]int]bool{}
	for _, eff := range effects {
		if eff.Type != "open" {
			continue
		}
		for _, c := range eff.Coords {
			gotCoords[[2]int{c.X, c.Y}] = true
		}
	}
	for coord := range wantCoords {
		if !gotCoords[coord] {
			t.Errorf("missing coord: %v", coord)
		}
	}
}

func TestRunScript_Ferz(t *testing.T) {
	script := `{"input": "Координаты выбранной клетки (x, y)", "actions": [{"Name": "OPEN_CELL", "Args": {"x": "$x", "y": "$y"}}, {"Name": "OPEN_CELL", "Args": {"x": "{'Name': 'RAND', 'Args': 'None'} - $FIELD_SIZE + $x", "y": "{'Name': 'PREV_RAND', 'Args': 'None'} - $FIELD_SIZE + $y"}}, {"Name": "OPEN_CELL", "Args": {"x": "{'Name': 'RAND', 'Args': 'None'} - $FIELD_SIZE + $x", "y": "$y - {'Name': 'PREV_RAND', 'Args': 'None'} + $FIELD_SIZE"}}, {"Name": "OPEN_CELL", "Args": {"x": { "Name": "RAND", "Args": "None" }, "y": "$y"}}, {"Name": "OPEN_CELL", "Args": {"x": "$x", "y": { "Name": "RAND", "Args": "None" }}}]}`
	input := ItemInput{X: 5, Y: 5}
	state := makeTestState()
	effects, err := RunScript(script, state, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(effects) == 0 {
		t.Fatalf("no effects returned")
	}
	// Проверяем, что хотя бы одна клетка совпадает с (5,5)
	found := false
	for _, eff := range effects {
		if eff.Type != "open" {
			continue
		}
		for _, c := range eff.Coords {
			if c.X == 5 && c.Y == 5 {
				found = true
			}
		}
	}
	if !found {
		t.Errorf("expected coord (5,5) in effects")
	}
}

func TestRunScript_Submarine(t *testing.T) {
	script := `{"input": "Координаты выбранной клетки (x, y)", "actions": [{ "Name": "PLACE_SUBMARINE", "Args": {"x": "x", "y": "y"} }, { "Name":"END_PLAYER_ACTION", "Args": "None" }]}`
	input := ItemInput{X: 1, Y: 1}
	state := makeTestState()
	// Добавим корабль длиной 3 на (1,1) горизонтально
	ship := &game.Ship{
		ID:       1,
		Len:      3,
		Coords:   game.Coord{X: 1, Y: 1},
		Bearings: game.Horizontal,
		Health:   3,
		Decks: map[game.Coord]bool{
			{X: 1, Y: 1}: true,
			{X: 2, Y: 1}: true,
			{X: 3, Y: 1}: true,
		},
	}
	state.PlayerState.Ships[1] = ship
	for c := range ship.Decks {
		state.PlayerState.Field[c.X][c.Y].ShipID = 1
	}
	effects, err := RunScript(script, state, input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, eff := range effects {
		if eff.Type == "place-submarine" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected place-submarine effect")
	}
}
