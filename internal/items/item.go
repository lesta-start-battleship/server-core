package items

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lesta-battleship/server-core/internal/config"
	"github.com/lesta-battleship/server-core/internal/game"
)

type ItemID int

type Item struct {
	Name        string `json:"name"`
	Kind        string `json:"kind"`
	Description string `json:"description"`
	Script      string `json:"script"`
	ID          ItemID `json:"id"`
	UseLimit    int    `json:"use_limit"`
	Cooldown    int    `json:"cooldown"`
}

type ItemInput struct {
	X         int `json:"x"`
	Y         int `json:"y"`
	Direction int `json:"direction,omitempty"`
	ItemID    int `json:"item_id"`
}

type ItemUsageData struct {
	UsedTimes    int
	LastUsedTurn int
}

type ItemEffect struct {
	Type   string       `json:"type"` // "open", "heal", "shoot" и так далее карочи, заеб
	Coords []game.Coord `json:"coords"`
}

func GetAllItems() (map[ItemID]*Item, error) {
	resp, err := http.Get(config.GetAllItemsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var items []Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}

	result := make(map[ItemID]*Item)
	for i := range items {
		item := &items[i]
		result[item.ID] = item
	}

	return result, nil
}

func GetUserItems(userJWT string) (map[ItemID]int, error) {
	url := config.GetAllUserItemsURl
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+userJWT)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user_inventory returned status %d", resp.StatusCode)
	}
	var inv struct {
		UserID      int `json:"user_id"`
		LinkedItems []struct {
			ItemID int `json:"item_id"`
			Amount int `json:"amount"`
		} `json:"linked_items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&inv); err != nil {
		return nil, err
	}

	result := make(map[ItemID]int)
	for _, item := range inv.LinkedItems {
		result[ItemID(item.ItemID)] = item.Amount
	}

	return result, nil
}
