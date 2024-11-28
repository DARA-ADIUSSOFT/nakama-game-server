package main

import (
	"context"
	"database/sql"
	"time"

	constain "nakama-game-server/constain"
	match "nakama-game-server/match"
	rpcFunction "nakama-game-server/rpc-function"

	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	initStart := time.Now()

	marshaler := &protojson.MarshalOptions{
		UseEnumNumbers: true,
	}
	unmarshaler := &protojson.UnmarshalOptions{
		DiscardUnknown: false,
	}

	if err := initializer.RegisterRpc("Health-Check", RpcHealthCheck); err != nil {
		return err
	}

	if err := initializer.RegisterRpc("find_match", match.RpcFindMatch(marshaler, unmarshaler)); err != nil {
		return err
	}

	errUserMove := initializer.RegisterRpc("user-move", rpcFunction.RpcUserMove2D)
	if errUserMove != nil {
		return errUserMove
	}

	if err := initializer.RegisterMatch(constain.ModuleName, CreateMatch); err != nil {
		logger.Error("Failed to register match: %s", err.Error())
		return err
	}

	logger.Info("Module loaded in %dms", time.Since(initStart).Milliseconds())
	return nil
}

func CreateMatch(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule) (runtime.Match, error) {
	state := &match.MatchState{
		Players: make(map[string]match.Position),
	}
	return &match.UserMove2DMatch{State: state, Logger: logger, Nk: nk}, nil
}
