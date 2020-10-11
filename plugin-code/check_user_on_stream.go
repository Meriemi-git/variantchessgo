package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/heroiclabs/nakama-common/runtime"
)

type SessionCheckResponse struct {
	AlreadyConnected bool `json:"already_connected"`
}

func CheckUserSessionExists(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, ok := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !ok {
		return "", runtime.NewError("no ID for user; must be authenticated", 3)
	}

	var input map[string]string
	err := json.Unmarshal([]byte(payload), &input)
	if err != nil {
		return "", err
	}
	var authToken = input["VARIANT_CHESS_TOKEN"]

	count, err := nk.StreamCount(0, userID, "", "")
	if err != nil {
		return "", fmt.Errorf("unable to count notification stream for user: %s", userID)
	}

	if count > 1 {
		content := map[string]interface{}{
			"USER_ID":             userID,
			"VARIANT_CHESS_TOKEN": authToken,
		}
		nk.NotificationSend(ctx, userID, "disconnect previous", content, 666, "", true)
	}
	response, err := json.Marshal(&SessionCheckResponse{AlreadyConnected: count > 1})
	if err != nil {
		logger.Error("unable to encode json: %v", err)
		return "", errors.New("failed to encode json")
	}
	return string(response), nil
}
