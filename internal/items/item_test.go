package items

// import (
// 	"testing"

// 	"os"

// 	"strings"

// 	"github.com/lesta-battleship/server-core/internal/game"
// )

// func TestUseItem_OpenCell(t *testing.T) {
// 	itemsList := map[ItemID]*Item{
// 		1: {
// 			Name:     "TestItem",
// 			Kind:     "test",
// 			Script:   `{"actions": [{"OPEN_CELL": {"x": "x", "y": "y"}}]}`,
// 			ID:       1,
// 			UseLimit: 1,
// 		},
// 	}
// 	state := &game.States{
// 		PlayerState: game.NewGameState(),
// 		EnemyState:  game.NewGameState(),
// 	}
// 	params := map[string]any{"x": 3, "y": 4}

// 	// До применения предмета клетка закрыта
// 	if state.EnemyState.Field[3][4].State == game.Open {
// 		t.Fatal("Cell should be closed before using item")
// 	}

// 	res, err := UseItem(1, state, itemsList, params)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	if res != "ok" {
// 		t.Errorf("unexpected result: %v", res)
// 	}

// 	if state.EnemyState.Field[3][4].State != game.Open {
// 		t.Errorf("Cell (3,4) should be open after using item")
// 	}
// }

// func TestUseItem_RealItemFromServer(t *testing.T) {
// 	if os.Getenv("INVENTORY_SERVICE_URL") == "" {
// 		t.Skip("INVENTORY_SERVICE_URL not set, skipping integration test")
// 	}
// 	itemsList, err := GetAllItems()
// 	if err != nil {
// 		t.Skipf("GetAllItems failed: %v", err)
// 	}
// 	var itemID ItemID
// 	for id, item := range itemsList {
// 		if item.Name == "Крест Нахимова" {
// 			itemID = id
// 			break
// 		}
// 	}
// 	if itemID == 0 {
// 		t.Skip("No 'Крест Нахимова' in loaded items")
// 	}
// 	state := &game.States{
// 		PlayerState: game.NewGameState(),
// 		EnemyState:  game.NewGameState(),
// 	}
// 	params := map[string]any{"x": 5, "y": 5}
// 	coords := [][2]int{{5, 5}, {5, 6}, {6, 5}, {5, 4}, {4, 5}}
// 	for _, c := range coords {
// 		if state.EnemyState.Field[c[0]][c[1]].State == game.Open {
// 			t.Fatalf("Cell (%d,%d) should be closed before using item", c[0], c[1])
// 		}
// 	}
// 	_, err = UseItem(itemID, state, itemsList, params)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}
// 	for _, c := range coords {
// 		if state.EnemyState.Field[c[0]][c[1]].State != game.Open {
// 			t.Errorf("Cell (%d,%d) should be open after using item", c[0], c[1])
// 		}
// 	}
// }

// func TestUseItem_AllRealItemsWithOpenCell(t *testing.T) {
// 	if os.Getenv("INVENTORY_SERVICE_URL") == "" {
// 		t.Skip("INVENTORY_SERVICE_URL not set, skipping integration test")
// 	}
// 	itemsList, err := GetAllItems()
// 	if err != nil {
// 		t.Skipf("GetAllItems failed: %v", err)
// 	}
// 	found := false
// 	for id, item := range itemsList {
// 		if item.Script == "" || !strings.Contains(item.Script, "OPEN_CELL") {
// 			continue
// 		}
// 		isSwitch := strings.Contains(item.Script, "SWITCH_CASE") || strings.Contains(item.Script, "SWICH_CASE")
// 		openedAny := false
// 		for direction := 1; direction <= 8; direction++ {
// 			params := map[string]any{"x": 5, "y": 5}
// 			if isSwitch {
// 				params["direction"] = direction
// 			}
// 			state := &game.States{
// 				PlayerState: game.NewGameState(),
// 				EnemyState:  game.NewGameState(),
// 			}
// 			_, err := UseItem(id, state, itemsList, params)
// 			if err != nil {
// 				t.Logf("UseItem failed for item %v (%v) direction=%d: %v", item.Name, id, direction, err)
// 				continue
// 			}
// 			opened := false
// 			openedCoords := [][2]int{}
// 			for i := 0; i < 10; i++ {
// 				for j := 0; j < 10; j++ {
// 					if state.EnemyState.Field[i][j].State == game.Open {
// 						opened = true
// 						openedCoords = append(openedCoords, [2]int{i, j})
// 					}
// 				}
// 			}
// 			if opened {
// 				t.Logf("Item %v (%v) direction=%d opened cells: %v", item.Name, id, direction, openedCoords)
// 				openedAny = true
// 				break
// 			}
// 		}
// 		if !openedAny {
// 			t.Errorf("No cell opened for item %v (%v) for any direction", item.Name, id)
// 		} else {
// 			found = true
// 		}
// 	}
// 	if !found {
// 		t.Skip("No items with OPEN_CELL found in loaded items")
// 	}
// }

// func printBoard(board [10][10]game.CellState) string {
// 	var sb strings.Builder
// 	for y := 0; y < 10; y++ {
// 		for x := 0; x < 10; x++ {
// 			if board[x][y].State == game.Open {
// 				sb.WriteByte('O')
// 			} else {
// 				sb.WriteByte('.')
// 			}
// 		}
// 		sb.WriteByte('\n')
// 	}
// 	return sb.String()
// }

// func TestUseItem_VisualizeEnemyBoards(t *testing.T) {
// 	if os.Getenv("INVENTORY_SERVICE_URL") == "" {
// 		t.Skip("INVENTORY_SERVICE_URL not set, skipping integration test")
// 	}
// 	itemsList, err := GetAllItems()
// 	if err != nil {
// 		t.Skipf("GetAllItems failed: %v", err)
// 	}
// 	for id, item := range itemsList {
// 		if item.Script == "" || !strings.Contains(item.Script, "OPEN_CELL") {
// 			continue
// 		}
// 		isSwitch := strings.Contains(item.Script, "SWITCH_CASE") || strings.Contains(item.Script, "SWICH_CASE")
// 		for direction := 1; direction <= 8; direction++ {
// 			params := map[string]any{"x": 5, "y": 5}
// 			if isSwitch {
// 				params["direction"] = direction
// 			}
// 			state := &game.States{
// 				PlayerState: game.NewGameState(),
// 				EnemyState:  game.NewGameState(),
// 			}
// 			_, err := UseItem(id, state, itemsList, params)
// 			if err != nil {
// 				t.Logf("UseItem failed for item %v (%v) direction=%d: %v", item.Name, id, direction, err)
// 				continue
// 			}
// 			boardStr := printBoard(state.EnemyState.Field)
// 			t.Logf("Item: %v (%v) direction=%d\n%s", item.Name, id, direction, boardStr)
// 			if !isSwitch {
// 				break // только один раз для не-switch предметов
// 			}
// 		}
// 	}
// }
