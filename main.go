package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

// func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
//     logger.Info("Hello World!")
//     return nil
// }

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
    initStart := time.Now()

	errHealthCheck := initializer.RegisterRpc("health-check",RpcHealthCheck)
	if errHealthCheck != nil {
		return errHealthCheck
	}

	errUserMove := initializer.RegisterRpc("user-move",rpcUserMove2D)
	if errUserMove != nil {
		return errUserMove
	}

	logger.Info("Module loaded in %dms", time.Since(initStart).Milliseconds())
    return nil
}