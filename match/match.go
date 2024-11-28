package match

import (
	"context"
	"database/sql"
	"nakama-game-server/api"
	"nakama-game-server/constain"

	"github.com/heroiclabs/nakama-common/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

type nakamaRpcFunc func(context.Context, runtime.Logger, *sql.DB, runtime.NakamaModule, string) (string, error)

func RpcFindMatch(marshaler *protojson.MarshalOptions, unmarshaler *protojson.UnmarshalOptions) nakamaRpcFunc {
	return func(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
		// No available matches found, create a new one.
		matchID, err := nk.MatchCreate(ctx, constain.ModuleName, map[string]interface{}{"fast": 1})
		if err != nil {
			logger.Error("error creating match: %v", err)
			return "", constain.ErrInternalError
		}

		response, err := marshaler.Marshal(&api.RpcFindMatchResponse{MatchId: matchID})
		if err != nil {
			logger.Error("error marshaling response payload: %v", err.Error())
			return "", constain.ErrMarshal
		}

		return string(response), nil
	}
}
