package handlers

import (
	"log"

	"github.com/lesta-battleship/server-core/internal/event"
	"github.com/lesta-battleship/server-core/internal/game"
	"github.com/lesta-battleship/server-core/internal/match"
	"github.com/lesta-battleship/server-core/internal/transaction"
	"github.com/lesta-battleship/server-core/internal/wsiface"
)

type ShootHandler struct{}

func (h *ShootHandler) EventName() string {
	return "shoot"
}

func (h *ShootHandler) Handle(input any, ctx *wsiface.Context) error {
	ctx.Room.Mutex.Lock()
	defer ctx.Room.Mutex.Unlock()

	wsInput, ok := input.(wsiface.WSInput)
	if !ok {
		return SendError(ctx.Conn, "invalid input format for shoot")
	}

	if ctx.Room.Status != "playing" {
		return SendError(ctx.Conn, "game not started")
	}
	if ctx.Room.Turn != ctx.Player.ID {
		return SendError(ctx.Conn, "not your turn")
	}
	if wsInput.X < 0 || wsInput.X >= 10 || wsInput.Y < 0 || wsInput.Y >= 10 {
		return SendError(ctx.Conn, "coordinates out of bounds")
	}

	var target *match.PlayerConn
	if ctx.Room.Player1.ID == ctx.Player.ID {
		target = ctx.Room.Player2
	} else {
		target = ctx.Room.Player1
	}

	// не работает с опен селл, т.е. не дает выстрелить еще раз
	// targetCell := target.States.PlayerState.Field[wsInput.X][wsInput.Y]
	// if targetCell.State == game.Open {
	// 	return SendError(ctx.Conn, "you already shot at this cell")
	// }

	cmd := game.NewShootCommand(game.Coord{X: wsInput.X, Y: wsInput.Y})
	tx := transaction.NewTransaction()
	tx.Add(cmd)

	if err := tx.Execute(ctx.Player.States); err != nil {
		log.Println("[SHOOT] Error:", err)
		return SendError(ctx.Conn, err.Error())
	}

	log.Printf("[SHOOT] %s firing at (%d,%d)", ctx.Player.ID, wsInput.X, wsInput.Y)

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

	if err := Broadcast(ctx.Room, wsiface.EventShootResult, wsiface.ShootResultResponse{
		X:        wsInput.X,
		Y:        wsInput.Y,
		By:       ctx.Player.ID,
		Hit:      cmd.Success,
		NextTurn: target.ID,
		GameOver: gameOver,
	}); err != nil {
		log.Printf("[SHOOT] Broadcast error: %v", err)
	}

	if gameOver {
		ctx.Room.Status = "ended"
		ctx.Room.WinnerID = ctx.Player.ID

		var matchResult event.MatchResult

		if ctx.Room.Mode == "ranked" {
			winGain, loserGain := match.GetRatingGain(ctx.Player.Rating, target.Rating)
			matchResult = event.MatchResult{
				WinnerID:  ctx.Player.ID,
				LoserID:   target.ID,
				MatchID:   ctx.Room.RoomID,
				MatchDate: ctx.Room.CreatedAt,
				MatchType: ctx.Room.Mode,
				Experience: &event.Experience{
					WinnerGain: 30,
					LoserGain:  15,
				},
				Rating: &event.Rating{
					WinnerGain: winGain,
					LoserGain:  loserGain,
				},
			}
		}

		if ctx.Room.Mode == "random" {
			matchResult = event.MatchResult{
				WinnerID:  ctx.Player.ID,
				LoserID:   target.ID,
				MatchID:   ctx.Room.RoomID,
				MatchDate: ctx.Room.CreatedAt,
				MatchType: ctx.Room.Mode,
				Experience: &event.Experience{
					WinnerGain: 30,
					LoserGain:  15,
				},
			}
		}

		if ctx.Room.Mode == "custom" {
			matchResult = event.MatchResult{
				WinnerID:  ctx.Player.ID,
				LoserID:   target.ID,
				MatchID:   ctx.Room.RoomID,
				MatchDate: ctx.Room.CreatedAt,
				MatchType: ctx.Room.Mode,
			}
		}

		if ctx.Room.Mode == "guild_war_match" {
			matchResult = event.MatchResult{
				WinnerID:  ctx.Player.ID,
				LoserID:   target.ID,
				MatchID:   ctx.Room.RoomID,
				WarID:     ctx.Room.GuildWarID,
				MatchDate: ctx.Room.CreatedAt,
				MatchType: ctx.Room.Mode,
				Experience: &event.Experience{
					WinnerGain: 30,
					LoserGain:  15,
				},
			}
		}

		if err := ctx.Dispatcher.DispatchMatchResult(matchResult); err != nil {
			log.Printf("[KAFKA] Failed to dispatch match result: %v", err)
		}

		_ = Broadcast(ctx.Room, wsiface.EventGameEnd, wsiface.GameEndResponse{
			Winner: ctx.Player.ID,
		})

		if ctx.Room.Player1.Conn != nil {
			_ = ctx.Room.Player1.Conn.Close()
		}
		if ctx.Room.Player2.Conn != nil {
			_ = ctx.Room.Player2.Conn.Close()
		}
		match.Rooms.Delete(ctx.Room.RoomID)
	} else {
		ctx.Room.Turn = target.ID
		ctx.Player.MoveCount++
	}

	return nil
}
