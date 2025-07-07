package items

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/lesta-battleship/server-core/internal/game"
)

func printBoard(field [10][10]game.CellState) {
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			fmt.Printf("%d ", field[x][y].ShipID)
		}
		fmt.Println()
	}
	fmt.Println()
}

func TestRunScript_REMOVE_SHIP(t *testing.T) {
	gs := &game.GameState{}
	gs.Field = [10][10]game.CellState{}
	gs.Ships = make([]*game.Ship, 11)
	// Ставим "корабль" в (2,3)
	ship := &game.Ship{
		ID:       1,
		Len:      1,
		Coords:   game.Coord{X: 2, Y: 3},
		Bearings: game.Vertical,
		Health:   1,
		Decks:    map[game.Coord]bool{game.Coord{X: 2, Y: 3}: game.Whole},
	}
	gs.Ships[1] = ship
	gs.Field[2][3] = game.CellState{State: 2, ShipID: 1}
	gs.NumShips = 1
	state := &game.States{PlayerState: gs}

	action := `[ { "Name": "REMOVE_SHIP", "Args": { "x": 2, "y": 3 } } ]`
	_, err := RunScript(action, state, nil)
	if err != nil {
		t.Errorf("REMOVE_SHIP error: %v", err)
	}
	fmt.Println("REMOVE_SHIP result:")
	printBoard(gs.Field)
}

func TestRunScript_PLACE_SHIP(t *testing.T) {
	gs := &game.GameState{}
	gs.Field = [10][10]game.CellState{}
	gs.Ships = make([]*game.Ship, 11)
	state := &game.States{PlayerState: gs}

	action := `[ { "Name": "PLACE_SHIP", "Args": { "length": 3, "x": 1, "y": 1, "bearings": true } } ]`
	_, err := RunScript(action, state, nil)
	if err != nil {
		t.Errorf("PLACE_SHIP error: %v", err)
	}
	fmt.Println("PLACE_SHIP result:")
	printBoard(gs.Field)
}

func TestRunScript_HEAL_SHIP(t *testing.T) {
	gs := &game.GameState{}
	gs.Field = [10][10]game.CellState{}
	gs.Ships = make([]*game.Ship, 11)
	// Корабль с повреждённой палубой в (5,5)
	ship := &game.Ship{
		ID:       2,
		Len:      1,
		Coords:   game.Coord{X: 5, Y: 5},
		Bearings: game.Vertical,
		Health:   1,
		Decks:    map[game.Coord]bool{game.Coord{X: 5, Y: 5}: game.Hit},
	}
	gs.Ships[2] = ship
	gs.Field[5][5] = game.CellState{State: 1, ShipID: 2}
	gs.NumShips = 1
	state := &game.States{PlayerState: gs}

	action := `[ { "Name": "HEAL_SHIP", "Args": { "x": 5, "y": 5 } } ]`
	_, err := RunScript(action, state, nil)
	if err != nil {
		t.Errorf("HEAL_SHIP error: %v", err)
	}
	fmt.Println("HEAL_SHIP result:")
	printBoard(gs.Field)
}

func createInventoryForTest(userID int, role, jwt string) error {
	url := "http://37.9.53.107/inventory/"
	payload := map[string]interface{}{"user_id": userID, "role": role}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK && resp.StatusCode != 409 {
		return fmt.Errorf("inventory create returned status %d", resp.StatusCode)
	}
	return nil
}

func addItemToInventoryForTest(itemID, amount int, jwt string) error {
	url := "http://37.9.53.107/inventory/add_item"
	payload := map[string]interface{}{"item_id": itemID, "amount": amount}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("add_item returned status %d", resp.StatusCode)
	}
	return nil
}

func TestUseItem_Integration(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoiYWRtaW4ifQ.JzhxV6sJhyTWgr4F-_EeDHg3-urRQiZUWYU9EvMZNHU"
	userID := 1
	role := "admin"

	// 1. Получить список всех предметов
	items, err := GetAllItems()
	if err != nil {
		t.Fatalf("GetAllItems error: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("Нет предметов для теста")
	}
	t.Logf("Получено %d предметов: %+v", len(items), items)
	itemID := ItemID(items[0].ID)

	// 2. Создать инвентарь игрока
	if err := createInventoryForTest(userID, role, token); err != nil {
		t.Fatalf("createInventoryForTest error: %v", err)
	}
	t.Log("Инвентарь создан/уже существует")

	// 3. Добавить предметы в инвентарь
	if err := addItemToInventoryForTest(int(itemID), 2, token); err != nil {
		t.Fatalf("addItemToInventoryForTest error: %v", err)
	}
	t.Logf("Добавлено 2 предмета с itemID=%d", itemID)

	// 4. Получить инвентарь игрока
	inv, err := GetNumberItems(token)
	if err != nil {
		t.Fatalf("GetNumberItems error: %v", err)
	}
	t.Logf("Инвентарь до использования: %+v", inv)
	if inv[itemID] < 1 {
		t.Fatalf("В инвентаре нет нужного предмета")
	}

	// 5. Проверить применение UseItem
	state := &game.States{}
	item := &Item{ID: itemID, Script: "[]"}
	itemsList := map[ItemID]*Item{itemID: item}
	_, err = UseItem(itemID, state, itemsList, map[string]interface{}{})
	if err != nil {
		t.Fatalf("UseItem error: %v", err)
	}
	t.Logf("UseItem успешно применён к itemID=%d", itemID)

	// 6. Проверить, что количество предмета уменьшилось
	inv2, err := GetNumberItems(token)
	if err != nil {
		t.Fatalf("GetNumberItems error: %v", err)
	}
	t.Logf("Инвентарь после использования: %+v", inv2)
	if inv2[itemID] != inv[itemID]-1 {
		t.Fatalf("Количество предмета не уменьшилось после UseItem: было %d, стало %d", inv[itemID], inv2[itemID])
	}
	t.Logf("Тест завершён успешно: количество предмета уменьшилось корректно")
}

func TestUseItem_NoItemInInventory(t *testing.T) {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJyb2xlIjoiYWRtaW4ifQ.JzhxV6sJhyTWgr4F-_EeDHg3-urRQiZUWYU9EvMZNHU"

	// 1. Получить список всех предметов
	items, err := GetAllItems()
	if err != nil {
		t.Fatalf("GetAllItems error: %v", err)
	}
	if len(items) < 2 {
		t.Fatal("Нужно минимум два предмета для теста")
	}
	t.Logf("Получено %d предметов: %+v", len(items), items)
	itemID := ItemID(items[1].ID) // используем второй предмет, который не добавлялся

	// 2. Получить инвентарь игрока
	inv, err := GetNumberItems(token)
	if err != nil {
		t.Fatalf("GetNumberItems error: %v", err)
	}
	if inv[itemID] > 0 {
		t.Fatalf("В инвентаре неожиданно есть предмет")
	}

	// 3. Попробовать использовать предмет
	state := &game.States{}
	item := &Item{ID: itemID, Script: "[]"}
	itemsList := map[ItemID]*Item{itemID: item}
	_, err = UseItem(itemID, state, itemsList, map[string]interface{}{})
	if err == nil {
		t.Fatalf("UseItem должен вернуть ошибку, если предмета нет в инвентаре")
	}
	t.Logf("Ожидаемая ошибка: %v", err)
}
