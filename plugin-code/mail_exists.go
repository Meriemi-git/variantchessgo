package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func MailExists(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var input map[string]string
	err := json.Unmarshal([]byte(payload), &input)
	if err != nil {
		return "", err
	}

	var email sql.NullString
	err2 := db.QueryRowContext(ctx, "SELECT email FROM users WHERE email=$1 ", input["email"]).Scan(&email)
	switch {
	case err2 == sql.ErrNoRows:
		return "false", nil
	case err2 != nil:
		return "false", err2
	default:
		return "true", err2
	}
}
