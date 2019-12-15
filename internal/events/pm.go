package events

import (
	"fmt"

	"github.com/thisisnotjose/sctest/internal/types"
)

// ProcessPrivateMessage Adds the relationship of the users in the follow registry
func ProcessPrivateMessage(ctx *types.Context, evt types.Event) {
	if clientConn, ok := ctx.UsersPool[evt.ReceiverUserID]; ok {
		fmt.Fprint(clientConn, evt.Payload)
		fmt.Printf("NOTIFIED USER %v\n", evt.ReceiverUserID)
	}
}
