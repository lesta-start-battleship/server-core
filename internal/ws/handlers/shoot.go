package handlers

import (
	"errors"
	"lesta-battleship/server-core/internal/game"
	"lesta-battleship/server-core/internal/match"
	"lesta-battleship/server-core/internal/transaction"
	"log"

	"github.com/gorilla/websocket"
)

func HandleFire(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn, input EventInput) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	log.Printf("[SHOOT] %s firing at (%d,%d)", player.ID, input.X, input.Y)

	if room.Status != "playing" {
		err := errors.New("game not started")
		Send(conn, "shoot_error", err.Error())
		return err
	}
	if room.Turn != player.ID {
		err := errors.New("not your turn")
		Send(conn, "shoot_error", err.Error())
		return err
	}

	var target *match.PlayerConn
	if room.Player1.ID == player.ID {
		target = room.Player2
	} else {
		target = room.Player1
	}

	targetCell := target.States.PlayerState.Field[input.X][input.Y]
	if targetCell.State == game.Open {
		err := errors.New("you already shot at this cell")
		Send(conn, "shoot_error", err.Error())
		return err
	}

	cmd := game.NewShootCommand(game.Coord{X: input.X, Y: input.Y})
	tx := transaction.NewTransaction()
	tx.Add(cmd)

	if err := tx.Execute(player.States); err != nil {
		log.Println("[SHOOT] Error:", err)
		Send(conn, "shoot_error", err.Error())
		return err
	}

	gameOver := true
	for _, ship := range target.States.PlayerState.Ships {
		if ship == nil {
			continue
		}
		for _, deck := range ship.Decks {
			if deck == game.Whole {
				gameOver = false
				break
			}
		}
		if !gameOver {
			break
		}
	}

	Broadcast(room, "shoot_result", map[string]any{
		"x":         input.X,
		"y":         input.Y,
		"by":        player.ID,
		"hit":       cmd.Success,
		"next_turn": target.ID,
		"game_over": gameOver,
	})

	if gameOver {
		room.Status = "ended"
		room.WinnerID = player.ID

		Broadcast(room, "game_end", map[string]any{"winner": player.ID})
		if room.Player1.Conn != nil {
			_ = room.Player1.Conn.Close()
		}
		if room.Player2.Conn != nil {
			_ = room.Player2.Conn.Close()
		}
		match.Rooms.Delete(room.RoomID)
	} else {
		room.Turn = target.ID
	}

	return nil
}
