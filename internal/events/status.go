package events

import (
	"github.com/thisisnot/sctest/internal/types"
	"github.com/thisisnot/sctest/internal/users"
)

// ProcessStatusUpdate Adds the relationship of the users in the follow registry
func ProcessStatusUpdate(ctx *types.Context, evt types.Event) {
	users.SendNotificationToAllFollowers(ctx, evt.EmitterUserID, evt)
}
