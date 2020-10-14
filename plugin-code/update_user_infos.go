package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/heroiclabs/nakama-common/runtime"
)

// UpdateUserInfos comment
func UpdateUserInfos(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	userID, userIDOk := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !userIDOk {
		logger.Info("invalid context vars")
		return "", errors.New("invalid context")
	}

	var input map[string]string
	err := json.Unmarshal([]byte(payload), &input)
	if err != nil {
		return "", err
	}
	logger.Info("Getting infos")
	var displayName = input["displayName"]
	logger.Info("displayName " + displayName)
	var timeZone = input["timeZone"]
	logger.Info("timeZone " + timeZone)
	var langTag = input["langTag"]
	logger.Info("langTag " + langTag)
	var location = input["location"]
	logger.Info("location " + location)
	var avatarURL = input["avatarUrl"]
	logger.Info("avatarUrl " + avatarURL)

	var existingID sql.NullString
	sqlErr := db.QueryRowContext(ctx, "SELECT display_name FROM users WHERE display_name=$1", displayName).Scan(&existingID)
	logger.Info("Query proceed")
	switch {
	case sqlErr == sql.ErrNoRows:
		if err := nk.AccountUpdateId(ctx, userID, "", nil, displayName, timeZone, "location", langTag, "avatar_url"); err != nil {
			logger.Error("Account update error: %s", err.Error())
			return "", err
		}
		logger.Info("Update account proceed")
		return "", nil
	case sqlErr != nil:
		logger.Info("Erro sql")
		return "", sqlErr
	default:
		logger.Info("Already in use")
		customError := runtime.Error{Code: 6, Message: "displayName already in use"}
		return "", &customError
	}
}
