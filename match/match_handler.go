package match

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand/v2"

	"github.com/heroiclabs/nakama-common/runtime"
)

var _ runtime.Match = &UserMove2DMatch{}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// MatchState holds the state of the match, including all players and their positions.
type MatchState struct {
	Players map[string]Position `json:"players"`
}

const (
	OpCodeMove = 1
)

type UserMove2DMatch struct {
	State  *MatchState
	Logger runtime.Logger
	Nk     runtime.NakamaModule
}

// MatchInit implements runtime.Match.
func (u *UserMove2DMatch) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	fmt.Println("Game is init with param: ", params)

	u.Logger = logger
	u.Nk = nk
	u.State = &MatchState{
		Players: make(map[string]Position),
	}
	tickRate := 1 // Match updates 10 times per second.
	label := "UserMove2D Match"
	logger.Info("Match initialized with label: %s", label)
	return u.State, tickRate, label
}

// MatchJoin implements runtime.Match.
func (u *UserMove2DMatch) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	fmt.Println("Game join is init with")
	gameState := state.(*MatchState)

	for _, presence := range presences {
		gameState.Players[presence.GetUserId()] = Position{
			X: rand.Float64() * 100,
			Y: rand.Float64() * 100,
		}
		u.Logger.Info("Player %s joined", presence.GetUserId())
	}

	return gameState
}

// MatchJoinAttempt implements runtime.Match.
func (u *UserMove2DMatch) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	return state, true, ""
}

// MatchLeave implements runtime.Match.
func (u *UserMove2DMatch) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	gameState := state.(*MatchState)
	for _, presence := range presences {
		delete(gameState.Players, presence.GetUserId())
		dispatcher.BroadcastMessage(2, []byte(fmt.Sprintf("%s left the match.", presence.GetUsername())), nil, nil, true)
		u.Logger.Info("Player %s left the match.", presence.GetUsername())
	}
	return gameState
}

// MatchLoop implements runtime.Match.
func (u *UserMove2DMatch) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	fmt.Println("Game loop is init with param: ", messages)

	// return gameState
	gameState := state.(*MatchState)

	for _, message := range messages {
		switch message.GetOpCode() {
		case OpCodeMove:
			var move Position
			err := json.Unmarshal(message.GetData(), &move)
			if err != nil {
				u.Logger.Error("Invalid move data: %v", err)
				continue
			}

			userID := message.GetUserId()
			if _, ok := gameState.Players[userID]; ok {

				u.Logger.Info("Player %s moved to X: %f, Y: %f", userID, move.X, move.Y)
				data, _ := json.Marshal(map[string]interface{}{
					"user_id": userID,
					"x":       move.X,
					"y":       move.Y,
				})
				dispatcher.BroadcastMessage(OpCodeMove, data, nil, nil, true)
			}
		}
	}

	return gameState
}

// MatchSignal implements runtime.Match.
func (u *UserMove2DMatch) MatchSignal(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, data string) (interface{}, string) {
	return state, ""
}

// MatchTerminate implements runtime.Match.
func (u *UserMove2DMatch) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	return state
}
