package main

import (
	"context"
	"database/sql"
	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	/*	if err := initializer.RegisterBeforeAuthenticateGoogle(OnUserAuthentBefore); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}*/

	if err := initializer.RegisterAfterAuthenticateGoogle(OnUserAuthentAfter); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}

	if err := initializer.RegisterRpc("user_exists", UserExists); err != nil {
		logger.Error("Unable to register: %v", err)
		return err
	}
	return nil
}
