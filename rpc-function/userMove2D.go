package rpcFunction

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
)

// Struct field tags should use lowercase JSON key names
type UserMove struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func RpcUserMove2D(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	type request struct {
		X int
		Y int
	}

	logger.Debug("User move RPC is called")
	var req request
	json.Unmarshal([]byte(payload), &req)
	fmt.Println("payload data: ", req.X)
	// Initialize the response with valid field names (capitalized)
	response := &UserMove{X: 2, Y: 2}

	// Marshal the response to JSON
	out, err := json.Marshal(response)
	if err != nil {
		logger.Error("Error marshalling response type to JSON: %v", err) // Corrected formatting
		return "", runtime.NewError("Cannot marshal type", 13)
	}

	// Debug log the JSON output for debugging
	logger.Debug("Marshalled JSON response: %s", string(out))
	fmt.Println("Response:", string(out)) // Print for local testing

	return string(out), nil
}
