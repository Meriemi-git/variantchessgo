package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	"github.com/heroiclabs/nakama-common/runtime"
)

// MatchState comment
type MatchState struct {
	presences     map[string]runtime.Presence
	whitePlayerID string
	text          string
}

// Match comment
type Match struct{}

func main() {

}

// MatchInit comment
func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	presences := make(map[string]runtime.Presence)
	state := &MatchState{
		presences:     presences,
		text:          "Match state de test",
		whitePlayerID: "",
	}
	// Initialize random num generator
	rand.Seed(time.Now().UnixNano())
	whiteIdx := rand.Intn(2)
	logger.Warn("Random white idx : %d", whiteIdx)
	loopIdx := 0
	for _, presence := range presences {
		if whiteIdx == loopIdx {
			state.whitePlayerID = presence.GetUserId()
		}
		loopIdx++
	}
	tickRate := 1
	label := "TestMatchLabel"

	return state, tickRate, label
}

// MatchJoinAttempt comment
func (m *Match) MatchJoinAttempt(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presence runtime.Presence, metadata map[string]string) (interface{}, bool, string) {
	acceptUser := true
	return state, acceptUser, ""
}

// MatchJoin comment
func (m *Match) MatchJoin(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	mState, _ := state.(*MatchState)
	for _, p := range presences {
		mState.presences[p.GetUserId()] = p
	}
	message, err := json.Marshal(mState)
	if err == nil {
		logger.Warn("Send message on MatchJoin")
		dispatcher.BroadcastMessage(2, []byte(message), presences, nil, true)
	}
	return mState
}

// MatchLeave comment
func (m *Match) MatchLeave(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, presences []runtime.Presence) interface{} {
	mState, _ := state.(*MatchState)
	for _, p := range presences {
		delete(mState.presences, p.GetUserId())
	}

	return mState
}

// MatchLoop comment
func (m *Match) MatchLoop(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, messages []runtime.MatchData) interface{} {
	mState, _ := state.(*MatchState)

	for _, message := range messages {
		const nbPresence = 3
		var presences []runtime.Presence
		for _, presence := range mState.presences {
			presences = append(presences, presence)
		}
		dispatcher.BroadcastMessage(1, message.GetData(), presences, nil, true)
	}

	return mState
}

// MatchTerminate comment
func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	message := "Server shutting down in " + strconv.Itoa(graceSeconds) + " seconds."
	dispatcher.BroadcastMessage(666, []byte(message), nil, nil, true)
	return state
}
