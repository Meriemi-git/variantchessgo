package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/runtime"
)

type CustomError struct {
	message    string
	statusCode int
}

func (this *CustomError) New(message string) {
	this.message = message
}
func (this *CustomError) Error() string {
	return this.message
}

func (this *CustomError) SetCode(code int) {
	this.statusCode = code
}

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
	logger.Info("On Before")
	vars := in.GetAccount().Vars
	signType := vars["signType"]
	email := vars["mail"]
	logger.Info("signType %s", signType)
	logger.Info("mail %s", email)
	if signType == "SIGNIN" {
		rows, err := db.QueryContext(ctx, "SELECT email, metadata FROM users WHERE email=$1 OR email IS NULL", email)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			var metadata []byte
			var userMail sql.NullString
			if err := rows.Scan(&userMail, &metadata); err != nil {
				return nil, err
			}
			if userMail.Valid {
				if userMail.String == email {
					var custErr CustomError
					custErr.New("do you want to link your account")
					custErr.SetCode(406)
					return nil, &custErr
				}
			} else {
				var input map[string]string
				err := json.Unmarshal([]byte(metadata), &input)
				if err != nil {
					return nil, err
				}
				if input["mail"] == email {
					return nil, errors.New("you have already signi with this account please signup instead")
				}
			}
		}
		rerr := rows.Close()
		if rerr != nil {
			return nil, err
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return in, nil
}
