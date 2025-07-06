package items

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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

var BaseItemsAPIURL string = "http://37.9.53.107"

func init() {
	if apiURL := os.Getenv("ITEMS_API_URL"); apiURL != "" {
		BaseItemsAPIURL = apiURL
	}
}

func UseItem(id ItemID, state *game.States, itemsList map[ItemID]*Item, params map[string]interface{}, userJWT string) (string, error) {
	item, ok := itemsList[id]
	if !ok {
		return "", fmt.Errorf("item with id %d not found", id)
	}

	if err := UseInventoryItem(id, userJWT); err != nil {
		return "", fmt.Errorf("failed to use item in inventory: %w", err)
	}

	return RunScript(item.Script, state, params)
}

func GetAllItems() ([]Item, error) {
	resp, err := http.Get(BaseItemsAPIURL + "/items/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var items []Item
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, err
	}
	return items, nil
}

func GetNumberItems(userJWT string) (map[ItemID]int, error) {
	url := BaseItemsAPIURL + "/inventory/user_inventory"
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

func UseInventoryItem(itemID ItemID, userJWT string) error {
	url := BaseItemsAPIURL + "/inventory/use_item"
	payload := map[string]interface{}{"item_id": int(itemID), "amount": 1}

	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+userJWT)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("use_item returned status %d", resp.StatusCode)
	}
	return nil
}
