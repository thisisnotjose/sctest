package users

import (
	"fmt"

	"github.com/thisisnotjose/sctest/internal/types"
)

// SendEventToUser sends the payload of an event to a user as long as the connection is in the usersPool
// if not it sends the payload to the dead letter queue implementation
func SendEventToUser(ctx *types.Context, userID int, evt types.Event) {
	clientConn, ok := ctx.UsersPool[userID]
	if ok {
		fmt.Fprint(clientConn, evt.Payload)
	} else {
		// If the user is not connected we send the payload to the dead letter queue
		ctx.DeadLetterChannel <- evt.Payload
	}
}

// SendEventToAllUsers sends an event payload to all users in the usersPool
func SendEventToAllUsers(ctx *types.Context, evt types.Event) {
	// If the event is sent to all users don't register the missing events to the dead queue.
	for _, clientConn := range ctx.UsersPool {
		fmt.Fprint(clientConn, evt.Payload)
	}
}

// SendNotificationToAllFollowers sends an event payload to all followers of a user
func SendNotificationToAllFollowers(ctx *types.Context, userID int, evt types.Event) {
	if followers, ok := ctx.FollowRegistry[userID]; ok {
		for follower := range followers {
			SendEventToUser(ctx, follower, evt)
		}
	}
}
