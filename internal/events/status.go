package events

import (
	"fmt"

	"github.com/thisisnotjose/sctest/internal/types"
)

// ProcessStatusUpdate Adds the relationship of the users in the follow registry
func ProcessStatusUpdate(ctx *types.Context, evt types.Event) {
	if followers, ok := ctx.FollowRegistry[evt.EmitterUserID]; ok {
		fmt.Printf("FOLLOW REGISTRY OF USER %v FOUND %v\n", evt.EmitterUserID, followers)
		for follower := range followers {
			clientConn, ok := ctx.UsersPool[follower]
			if ok {
				fmt.Fprint(clientConn, evt.Payload)
				_, err := fmt.Printf("NOTIFIED USER %v\n", follower)
				if err != nil {
					fmt.Printf("FAILED NOTIFIED USER %v because %v\n", evt.ReceiverUserID, err)
				}
			}
		}
	}
}
