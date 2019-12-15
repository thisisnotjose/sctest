package processors

import (
	"github.com/thisisnotjose/sctest/internal/events"
	"github.com/thisisnotjose/sctest/internal/types"
)

type processor struct {
	ctx *types.Context
}

// NewEventsProcessor returns a new events processor
func NewEventsProcessor(ctx *types.Context) types.Processor {
	return processor{
		ctx: ctx,
	}
}

// Start kicks off the listener on the events channel and for every message it
// adds it to the queue and then processes as many messages as it can
func (p processor) Start() {
	lastSeqNo := 0

	for newEvent := range p.ctx.EventChannel {
		p.ctx.EventQueue[newEvent.Sequence] = newEvent
		lastSeqNo = processEventQueue(p.ctx, lastSeqNo)
	}
}

// processEventQueue takes events from the queue and processes sequentially,
// it stops when the next item to be processed is not on the queue
func processEventQueue(ctx *types.Context, lastSeqNo int) int {
	for {
		nextEvent, ok := ctx.EventQueue[lastSeqNo+1]
		if !ok {
			break
		}

		delete(ctx.EventQueue, lastSeqNo+1)
		processEvent(ctx, nextEvent)
		lastSeqNo++
	}

	return lastSeqNo
}

func processEvent(ctx *types.Context, evt types.Event) {
	switch evt.EventType {
	case "F":
		events.ProcessFollow(ctx, evt)
	case "U":
		events.ProcessUnfollow(ctx, evt)
	case "P":
		events.ProcessPrivateMessage(ctx, evt)
	case "B":
		events.ProcessBroadcast(ctx, evt)
	case "S":
		events.ProcessStatusUpdate(ctx, evt)
	}
	//fmt.Printf("message processed: %+v\n", evt)
}
