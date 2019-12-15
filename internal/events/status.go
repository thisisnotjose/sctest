package events

import (
	"fmt"

	"github.com/thisisnotjose/sctest/internal/types"
)

// ProcessStatusUpdate Adds the relationship of the users in the follow registry
func ProcessStatusUpdate(ctx *types.Context, evt types.Event) {
	if followers, ok := ctx.FollowRegistry[evt.EmitterUserID]; ok {
		for follower := range followers {
			clientConn, ok := ctx.UsersPool[follower]
			if ok {
				fmt.Fprint(clientConn, evt.Payload)
			}
		}
	}
}
