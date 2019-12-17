package events

import (
	"github.com/thisisnotjose/sctest/internal/types"
	"github.com/thisisnotjose/sctest/internal/users"
)

// ProcessBroadcast Adds the relationship of the users in the follow registry
func ProcessBroadcast(ctx *types.Context, evt types.Event) {
	users.SendEventToAllUsers(ctx, evt)
}
