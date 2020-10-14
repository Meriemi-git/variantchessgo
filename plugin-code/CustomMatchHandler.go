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
	presences map[string]runtime.Presence
}

// VariantUser comment
type VariantUser struct {
	UserID   string
	Username string
	Color    string
}

// Match comment
type Match struct{}

func main() {

}

// MatchInit comment
func (m *Match) MatchInit(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, params map[string]interface{}) (interface{}, int, string) {
	state := &MatchState{
		presences: make(map[string]runtime.Presence),
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
	for _, presence := range presences {
		mState.presences[presence.GetUserId()] = presence
	}

	if len(mState.presences) >= 2 {
		var userList []VariantUser
		var recipients []runtime.Presence
		// Initialize random num generator
		rand.Seed(time.Now().UnixNano())
		whiteIdx := rand.Intn(2)
		loopIndex := 0
		for _, connected := range mState.presences {
			recipients = append(recipients, connected)
			if whiteIdx == loopIndex {
				userList = append(userList, VariantUser{
					UserID:   connected.GetUserId(),
					Color:    "white",
					Username: connected.GetUsername(),
				})
			} else {
				userList = append(userList, VariantUser{
					UserID:   connected.GetUserId(),
					Color:    "black",
					Username: connected.GetUsername(),
				})
			}
			loopIndex = loopIndex + 1
		}

		response, err := json.Marshal(userList)
		if err == nil {
			dispatcher.BroadcastMessage(2, []byte(string(response)), recipients, nil, true)
		}
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
		var sender runtime.Presence
		var recipients []runtime.Presence
		for _, presence := range mState.presences {
			if presence.GetUserId() != message.GetUserId() {
				recipients = append(recipients, presence)
			} else {
				sender = presence
			}
		}
		dispatcher.BroadcastMessage(1, message.GetData(), recipients, sender, true)
	}

	return mState
}

// MatchTerminate comment
func (m *Match) MatchTerminate(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, dispatcher runtime.MatchDispatcher, tick int64, state interface{}, graceSeconds int) interface{} {
	message := "Server shutting down in " + strconv.Itoa(graceSeconds) + " seconds."
	dispatcher.BroadcastMessage(666, []byte(message), nil, nil, true)
	return state
}
