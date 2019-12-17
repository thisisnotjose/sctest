package events

import (
	"github.com/thisisnotjose/sctest/internal/types"
	"github.com/thisisnotjose/sctest/internal/users"
)

// ProcessFollow Adds the relationship of the users in the follow registry
func ProcessFollow(ctx *types.Context, evt types.Event) {
	// Fetch the follow list of the followed user and if it doesn't have one create it
	if _, ok := ctx.FollowRegistry[evt.ReceiverUserID]; !ok {
		ctx.FollowRegistry[evt.ReceiverUserID] = make(map[int]bool)
	}

	// Add the sender to the followers of the user
	followers, _ := ctx.FollowRegistry[evt.ReceiverUserID]
	followers[evt.EmitterUserID] = true

	users.SendEventToUser(ctx, evt.ReceiverUserID, evt)
}
