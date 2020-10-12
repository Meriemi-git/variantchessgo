package main

import (
	"context"
	"database/sql"

	"github.com/heroiclabs/nakama-common/runtime"
)

// OnRegisterMatchMakerMatched comment
func OnRegisterMatchMakerMatched(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, entries []runtime.MatchmakerEntry) (string, error) {
	for _, e := range entries {
		logger.Info("Matched user '%s' named '%s'", e.GetPresence().GetUserId(), e.GetPresence().GetUsername())
		for k, v := range e.GetProperties() {
			logger.Info("Matched on '%s' value '%v'", k, v)
		}
	}
	params := map[string]interface{}{
		"some": "data",
	}
	matchID, err := nk.MatchCreate(ctx, "variantchess", params)
	if err != nil {
		return "", err
	}
	//matchId, err := nk.MatchCreate(ctx, "variantchess", map[string]interface{}{"invited": entries, "white_player": entries[0].GetPresence().GetUserId()})
	//if err != nil {
	//	return "", err
	//}

	return matchID, nil

}
