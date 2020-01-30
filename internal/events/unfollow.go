package events

import (
	"github.com/thisisnot/sctest/internal/types"
)

// ProcessUnfollow Adds the relationship of the users in the follow registry
func ProcessUnfollow(ctx *types.Context, evt types.Event) {
	if followers, ok := ctx.FollowRegistry[evt.ReceiverUserID]; ok {
		delete(followers, evt.EmitterUserID)
	}
}
