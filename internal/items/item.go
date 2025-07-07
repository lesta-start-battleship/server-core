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
}

func UseItem(id ItemID, state *game.States, itemsList map[ItemID]*Item, params map[string]any) (string, error) {
	item, ok := itemsList[id]
	if !ok {
		return "", fmt.Errorf("item with id %d not found", id)
	}

	// if err := UseInventoryItem(id, userJWT); err != nil {
	// 	return "", fmt.Errorf("failed to use item in inventory: %w", err)
	// }

	return RunScript(item.Script, state, params)
}

// func GetAllItems() ([]Item, error) {
// 	resp, err := http.Get(BaseItemsAPIURL + "/items/")
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	var items []Item
// 	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
// 		return nil, err
// 	}
// 	return items, nil
// }

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

// func UseInventoryItem(itemID ItemID, userJWT string) error {
// 	url := BaseItemsAPIURL + "/inventory/use_item"
// 	payload := map[string]interface{}{"item_id": int(itemID), "amount": 1}

// 	body, _ := json.Marshal(payload)
// 	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
// 	if err != nil {
// 		return err
// 	}

// 	req.Header.Set("Authorization", "Bearer "+userJWT)
// 	req.Header.Set("Content-Type", "application/json")
// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return fmt.Errorf("use_item returned status %d", resp.StatusCode)
// 	}
// 	return nil
// }
