package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

func OnUserAuthentAfter(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, out *api.Session, in *api.AuthenticateGoogleRequest) error {
	logger.Info("OnUserAuthentAfter")
	vars, varsOk := ctx.Value(runtime.RUNTIME_CTX_VARS).(map[string]string)
	if !varsOk {
		logger.Info("invalid context vars")
		return errors.New("invalid context")
	}
	userId, userIdOk := ctx.Value(runtime.RUNTIME_CTX_USER_ID).(string)
	if !userIdOk {
		logger.Info("invalid context userId")
		return errors.New("invalid context")
	}

	signType := vars["signType"]
	mail := vars["mail"]
	logger.Info("signType : %s", signType)
	logger.Info("mail : %s", mail)
	logger.Info("Create : %b", out.Created)
	logger.Info("userId : %s", userId)
	if signType == "SIGNIN" {

	}
	return errors.New("try to create existing google account")
}

func OnUserAuthentBefore(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, in *api.AuthenticateGoogleRequest) (*api.AuthenticateGoogleRequest, error) {

	vars := in.GetAccount().Vars

	signType := vars["signType"]
	email := vars["mail"]
	//var id int
	logger.Info("mail : %s", email)
	logger.Info("signType : %s", signType)
	if signType == "SIGNIN" {
		logger.Info("In signin")
		rows, err := db.QueryContext(ctx, "SELECT metadata FROM users")
		logger.Info("after query")
		if err != nil {
			logger.Info("Error 1")
			return nil, err
		}
		defer rows.Close()
		logger.Info("after defer")
		for rows.Next() {
			logger.Info("In For")

			logger.Info("Mail is empty")
			var metadata []byte
			if err := rows.Scan(&metadata); err != nil {
				// Check for a scan error.
				// Query rows will be closed with defer.
				logger.Info("Error 3")
				return nil, err
			}
			var input map[string]string
			err := json.Unmarshal([]byte(metadata), &input)
			if err != nil {
				logger.Info("Error 4")
				return nil, err
			}
			if input["mail"] == email {
				logger.Info("error found google account with this name")
				return nil, errors.New("found google account with this name")
			}
		}
		rerr := rows.Close()
		if rerr != nil {
			logger.Info("Error 5")
			return nil, err
		}

		// Rows.Err will report the last error encountered by Rows.Scan.
		if err := rows.Err(); err != nil {
			logger.Info("Error 6")
			return nil, err
		}
	}
	return in, nil
}
