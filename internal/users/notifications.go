package users

import (
	"fmt"

	"github.com/thisisnotjose/sctest/internal/types"
)

// SendEventToUser sends the payload of an event to a user as long as the connection is in the usersPool
// if not it sends the payload to the dead letter queue implementation
func SendEventToUser(ctx *types.Context, userID int, evt types.Event) {
	var err error = nil
	clientConn, ok := ctx.UsersPool[userID]
	if ok {
		_, err = fmt.Fprint(clientConn, evt.Payload)
	}

	// If the user is not connected send it to the dead letter queue instead
	if !ok {
		ctx.DeadLetterChannel <- evt.Payload
	}

	// If we couldn't send the payload, send it to the dead letter queue instead
	if err != nil {
		fmt.Printf("Could not send message to user: %d error: %v", userID, err)
		ctx.DeadLetterChannel <- evt.Payload
	}
}

// SendEventToAllUsers sends an event payload to all users in the usersPool
func SendEventToAllUsers(ctx *types.Context, evt types.Event) {
	// If the event is sent to all users don't register the missing events to the dead queue.
	for userID, clientConn := range ctx.UsersPool {
		_, err := fmt.Fprint(clientConn, evt.Payload)

		// If we couldn't send the payload, send it to the dead letter queue instead
		if err != nil {
			fmt.Printf("Could not send message to user: %d error: %v", userID, err)
			ctx.DeadLetterChannel <- evt.Payload
		}
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
