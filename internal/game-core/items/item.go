package items

import (
	"encoding/json"
	// "fmt"
	"lesta-battleship/server-core/internal/game-core/game"
	"net/http"
)

type ItemID int

type Item struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	Script      string `json:"script"`
	ID          ItemID    `json:"id"`
}

type ItemsInfo struct {
	Items map[ItemID]*Item
	ItemsPlayer1 map[ItemID]int
	ItemsPlayer2 map[ItemID]int
}

func GetItemsInfo(player1ID, player2ID string) (*ItemsInfo, error) {
	items, err := GetItemsFunc()
	if err != nil {
		return nil, err
	}
	items_1, err := GetNumberItems(player1ID)
	if err != nil {
		return nil, err
	}

	items_2, err := GetNumberItems(player1ID)
	if err != nil {
		return nil, err
	}	
	return &ItemsInfo{
		Items: items,
		ItemsPlayer1: items_1,
		ItemsPlayer2: items_2,
	}, nil

}

func GetItemsFunc() (map[ItemID]*Item, error) {
	r, err := http.Get("http://37.9.53.107/items/")
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var items []Item
	if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
		return nil, err
	}

	// TODO
	// превратить слайс в map
	// return items, nil
	return nil, nil
}

func GetNumberItems(ID string) (map[ItemID]int, error) {
	// TODO: достать количество предметов каждого типа у игрока
	return nil, nil
}

// TODO: сделай обертку, по вызову которой мы будем сообщать сервису предметов, об использовании предмета

func UseItem(id ItemID, state *game.States, itemsList map[ItemID]*Item, params map[string]interface{}) (string, error) {
	var item *Item
	/* TODO: сделай использование через map
	for i := range itemsList {
		if itemsList[i].ID == id {
			item = &itemsList[i]
			break
		}
	}
	if item == nil {
		return "", fmt.Errorf("item with id %d not found", id)
	}
	*/
	return RunScript(item.Script, state, params)
	// TODO: добавить логику сообщения сервису об использовании предмета
}