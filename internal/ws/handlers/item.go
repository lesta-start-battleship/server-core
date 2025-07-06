package handlers

import (
	"errors"
	"log"

	"github.com/lesta-battleship/server-core/internal/event"
	"github.com/lesta-battleship/server-core/internal/items"
	"github.com/lesta-battleship/server-core/internal/match"

	"github.com/gorilla/websocket"
)

func HandleItem(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn, input WSInput, dispatcher *event.MatchEventDispatcher) error {
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

	if player.Items[items.ItemID(input.ItemID)] == 0 {
		err := errors.New("the user does not have this item")
		SendError(conn, err.Error())
		return err
	}
	// TODO: можем ли использовать данный предмет(ограничения на применения в 1 игре/ограничения на применения несколько ходов подряд)

	// TODO: пока нет идей как это делать в общем случае, предлагайте варианты
	params := map[string]any{"x": 5, "y": 5}

	_, err := items.UseItem(items.ItemID(input.ItemID), player.States, room.Items, params)
	if err != nil {
		SendError(conn, err.Error())
		return err
	}

	usedItem := event.Item{
		PlayerID: player.ID,
		ItemID:   input.ItemID,
	}
	if err := dispatcher.DispatchUsedItem(usedItem); err != nil {
		log.Printf("[KAFKA] Failed to dispatch used item: %v", err)
	}

	// TODO: придумать как в общем виде выдавать результат исполнения команд
	Broadcast(room, EventItemUsed, ItemUsedResponse{
		Coords: affectedCoords, // координаты
		ItemID: input.ItemID,
		By:     player.ID,
	})

	return nil
}
