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

func UseItem(id ItemID, state *game.States, itemsList map[ItemID]*Item, params map[string]any) (string, error) {
	item, ok := itemsList[id]
	if !ok {
		return "", fmt.Errorf("item with id %d not found", id)
	}

	if item.UseLimit > 0 {
		used := 0
		if v, ok := params["used_count"].(int); ok {
			used = v
		}
		if used >= item.UseLimit {
			return "", fmt.Errorf("use limit reached for item %d", id)
		}
	}

	if item.Cooldown > 0 {
		if lastTurn, ok := params["last_used_turn"].(int); ok {
			currentTurn := 0
			if v, ok := params["turn"].(int); ok {
				currentTurn = v
			}
			if currentTurn > 0 && lastTurn > 0 && currentTurn-lastTurn < item.Cooldown {
				return "", fmt.Errorf("cooldown not expired for item %d", id)
			}
		}
	}

	return RunScript(item.Script, state, params)
}

// нам нужен в мааап
func GetAllItems() (map[ItemID]*Item, error) {
	resp, err := http.Get(config.InventoryServiceURL + "/items/")
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
	url := config.InventoryServiceURL + "/inventory/user_inventory"
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
