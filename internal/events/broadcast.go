package events

import (
	"fmt"

	"github.com/thisisnotjose/sctest/internal/types"
)

// ProcessBroadcast Adds the relationship of the users in the follow registry
func ProcessBroadcast(ctx *types.Context, evt types.Event) {
	for _, clientConn := range ctx.UsersPool {
		fmt.Fprint(clientConn, evt.Payload)
	}
}
