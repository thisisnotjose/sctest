package events

import (
	"github.com/thisisnot/sctest/internal/types"
	"github.com/thisisnot/sctest/internal/users"
)

// ProcessBroadcast Adds the relationship of the users in the follow registry
func ProcessBroadcast(ctx *types.Context, evt types.Event) {
	users.SendEventToAllUsers(ctx, evt)
}
