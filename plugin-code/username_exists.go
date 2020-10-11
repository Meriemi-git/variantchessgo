package main

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/heroiclabs/nakama-common/runtime"
)

func UsernameExists(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, payload string) (string, error) {
	var input map[string]string
	err := json.Unmarshal([]byte(payload), &input)
	if err != nil {
		return "", err
	}

	if input["searchGoogleAccount"] == "false" {
		var email sql.NullString
		err := db.QueryRowContext(ctx, "SELECT email FROM users WHERE email=$1 ", input["email"]).Scan(&email)
		switch {
		case err == sql.ErrNoRows:
			return "false", nil
		case err != nil:
			return "false", err
		default:
			return "true", err
		}
	} else {
		rows, err := db.QueryContext(ctx, "SELECT metadata FROM users WHERE email IS NULL ")
		if err != nil {
			return "false", err
		}
		defer rows.Close()
		for rows.Next() {
			var metadata []byte
			if err := rows.Scan(&metadata); err != nil {
				return "false", err
			}
			var data map[string]string
			err := json.Unmarshal([]byte(metadata), &data)
			if err != nil {
				return "false", nil
			}

			if data["mail"] == input["email"] {
				return "true", nil
			}
		}
		return "false", nil
	}
}
