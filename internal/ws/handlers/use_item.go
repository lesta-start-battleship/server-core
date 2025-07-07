package handlers

import (
	"errors"
	"log"

	"github.com/gorilla/websocket"
	"github.com/lesta-battleship/server-core/internal/event"
	"github.com/lesta-battleship/server-core/internal/items"
	"github.com/lesta-battleship/server-core/internal/match"
)

func HandleUseItem(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn, input WSInput, dispatcher *event.MatchEventDispatcher) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	if room.Status != "playing" {
		err := errors.New("game not started")
		SendError(conn, err.Error())
		return err
	}

	if room.Turn != player.ID {
		err := errors.New("not your turn")
		SendError(conn, err.Error())
		return err
	}

	// проверим доступность предмета
	if player.Items[items.ItemID(input.ItemID)] <= 0 {
		err := errors.New("item not available or already used")
		SendError(conn, err.Error())
		return err
	}

	log.Println("item dostupen")

	// найдём сам предмет из общей коллекции
	itemData, ok := room.Items[items.ItemID(input.ItemID)]
	if !ok {
		err := errors.New("item metadata not found")
		SendError(conn, err.Error())
		return err
	}
	log.Println("nashli iz obshey kollekcii")

	result, err := items.UseItem(items.ItemID(input.ItemID), player.States, room.Items, input.Params)
	if err != nil {
		log.Printf("[WS] Use item error: %v", err)
		SendError(conn, err.Error())
		return err
	}
	log.Println("usenuli item")

	player.Items[items.ItemID(input.ItemID)]--

	usedItem := event.Item{
		PlayerID: player.ID,
		ItemID:   input.ItemID,
	}
	if err := dispatcher.DispatchUsedItem(usedItem); err != nil {
		log.Printf("[KAFKA] Failed to dispatch used item: %v", err)
	}

	Broadcast(room, EventItemUsed, ItemUsedResponse{
		ItemID: items.ItemID(input.ItemID),
		Name:   itemData.Name,
		By:     player.ID,
		Result: result,
	})

	return nil
}
