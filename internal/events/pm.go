package events

import (
	"github.com/thisisnotjose/sctest/internal/types"
	"github.com/thisisnotjose/sctest/internal/users"
)

// ProcessPrivateMessage Adds the relationship of the users in the follow registry
func ProcessPrivateMessage(ctx *types.Context, evt types.Event) {
	users.SendEventToUser(ctx, evt.ReceiverUserID, evt)
}
