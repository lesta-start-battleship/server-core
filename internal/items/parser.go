package items

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lesta-battleship/server-core/internal/game"
	"github.com/lesta-battleship/server-core/internal/transaction"
)

type Action struct {
	Name string
	Args map[string]interface{}
}

func ParseScript(script string) ([]Action, error) {
	script = replaceSingleQuotes(script)
	var actions []Action

	err := json.Unmarshal([]byte(script), &actions)
	if err != nil {
		return nil, err
	}

	return actions, nil
}

func RunScript(script string, state *game.States, params map[string]interface{}) (string, error) {
	actions, err := ParseScript(script)
	if err != nil {
		return "", err
	}

	var lastResult string
	var prevRand float64

	tx := transaction.NewTransaction()
	for _, action := range actions {
		for k, v := range action.Args {
			if s, ok := v.(string); ok {
				val, err := evalExpr(s, params, prevRand)

				if err == nil {
					action.Args[k] = val
					if strings.Contains(strings.ToUpper(s), "RAND") {
						if f, ok := val.(float64); ok {
							prevRand = f
						}
					}
				}
			}
		}
		name := strings.ToUpper(action.Name)
		switch name {
		case "OPEN_CELL":
			x, okX := toFloat(action.Args["x"])
			y, okY := toFloat(action.Args["y"])
			if !okX || !okY {
				return "", fmt.Errorf("invalid args for open_cell")
			}

			tx.Add(
				game.NewOpenCellCommand(
					game.Coord{
						X: int(x),
						Y: int(y),
					}),
			)
		case "MAKE_SHOT":
			x, okX := toFloat(action.Args["x"])
			y, okY := toFloat(action.Args["y"])
			if !okX || !okY {
				return "", fmt.Errorf("invalid args for MAKE_SHOT")
			}

			tx.Add(
				game.NewShootCommand(
					game.Coord{
						X: int(x),
						Y: int(y),
					}),
			)
		//case "SET_CELL_STATUS":
		// TODO: удалить
		//case "SET_SHIP_COORDINATES":
		// REMOVE_SHIP + PLACE_SHIP
		case "REMOVE_SHIP":
			x, okX := toFloat(action.Args["x"])
			y, okY := toFloat(action.Args["y"])
			if !okX || !okY {
				return "", fmt.Errorf("invalid args for REMOVE_SHIP: x or y missing")
			}

			tx.Add(
				game.NewRemoveShipCommand(
					game.Coord{
						X: int(x),
						Y: int(y),
					}),
			)
		case "PLACE_SHIP":
			length, okLen := toFloat(action.Args["length"])
			x, okX := toFloat(action.Args["x"])
			y, okY := toFloat(action.Args["y"])
			bearings, okBearings := action.Args["bearings"].(bool)
			if !okLen || !okX || !okY || !okBearings {
				return "", fmt.Errorf("invalid args for PLACE_SHIP: length, x, y, or bearings missing")
			}

			tx.Add(
				game.NewPlaceShipCommand(int(length), game.Coord{X: int(x), Y: int(y)}, bearings),
			)
		case "HEAL_SHIP":
			x, okX := toFloat(action.Args["x"])
			y, okY := toFloat(action.Args["y"])
			if !okX || !okY {
				return "", fmt.Errorf("invalid args for HEAL_SHIP: x or y missing")
			}

			tx.Add(
				game.NewHealShipCommand(game.Coord{X: int(x), Y: int(y)}),
			)

		case "END_PLAYER_ACTION":
			lastResult = "end_action"
			// TODO: надо обсудить хотим ли давать предмету такой функционал и как его реализовывать
		default:
			return "", fmt.Errorf("unknown action: %s", action.Name)
		}
	}

	if err := tx.Execute(state); err != nil {
		return "", err
	}

	return lastResult, nil
}
