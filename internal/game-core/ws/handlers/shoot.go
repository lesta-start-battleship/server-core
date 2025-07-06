package handlers

import (
	"errors"
	"lesta-battleship/server-core/internal/game-core/event"
	"lesta-battleship/server-core/internal/game-core/game"
	"lesta-battleship/server-core/internal/game-core/match"
	"lesta-battleship/server-core/internal/game-core/transaction"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func HandleFire(room *match.GameRoom, player *match.PlayerConn, conn *websocket.Conn, input EventInput, dispatcher *event.MatchEventDispatcher) error {
	room.Mutex.Lock()
	defer room.Mutex.Unlock()

	log.Printf("[SHOOT] %s firing at (%d,%d)", player.ID, input.X, input.Y)

	if room.Status != "playing" {
		err := errors.New("game not started")
		Send(conn, EventShootError, err.Error())
		return err
	}
	if room.Turn != player.ID {
		err := errors.New("not your turn")
		Send(conn, EventShootError, err.Error())
		return err
	}

	if input.X < 0 || input.X >= 10 || input.Y < 0 || input.Y >= 10 {
		err := errors.New("coordinates out of bounds")
		Send(conn, EventShootError, err.Error())
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
		Send(conn, EventShootError, err.Error())
		return err
	}

	cmd := game.NewShootCommand(game.Coord{X: input.X, Y: input.Y})
	tx := transaction.NewTransaction()
	tx.Add(cmd)

	if err := tx.Execute(player.States); err != nil {
		log.Println("[SHOOT] Error:", err)
		Send(conn, EventShootError, err.Error())
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

	Broadcast(room, EventShootResult, ShootResultResponse{
		X:        input.X,
		Y:        input.Y,
		By:       player.ID,
		Hit:      cmd.Success,
		NextTurn: target.ID,
		GameOver: gameOver,
	})

	if gameOver {
		room.Status = "ended"
		room.WinnerID = player.ID

		var winnerID, loserID string
		if room.Player1.ID == player.ID {
			winnerID = room.Player1.ID
			loserID = room.Player2.ID
		} else {
			winnerID = room.Player2.ID
			loserID = room.Player1.ID
		}

		matchResult := event.MatchResult{
			WinnerID:  winnerID,
			LoserID:   loserID,
			MatchID:   room.RoomID,
			MatchDate: time.Now(),
			MatchType: room.Mode,
			Experience: &event.Experience{
				WinnerGain: 30,
				LoserGain:  -15,
			},
		}

		if err := dispatcher.DispatchMatchResult(matchResult); err != nil {
			log.Printf("[KAFKA] Failed to dispatch match result: %v", err)
		}

		Broadcast(room, EventGameEnd, GameEndResponse{
			Winner: player.ID,
		})

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
