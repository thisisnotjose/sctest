package events

import (
	"github.com/thisisnot/sctest/internal/types"
)

// ProcessDeadLetter adds the event to the dead letter queue
func ProcessDeadLetter(ctx *types.Context, payload string) {
	ctx.DeadLetterChannel <- payload
}
