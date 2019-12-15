package handlers

import (
	"log"
	"strconv"
	"strings"

	"github.com/thisisnotjose/sctest/internal/types"
)

// NewEventsHandler Registers a new user into the system
func NewEventsHandler(ctx *types.Context) types.Handler {
	return func(conn types.Connection, message string) {
		eventsHandler(ctx, conn, message)
	}
}

func eventsHandler(ctx *types.Context, conn types.Connection, message string) {
	var err error
	eventParts := strings.Split(message, "|")

	sequenceNo, err := strconv.Atoi(eventParts[0])
	if err != nil {
		log.Printf("Couldn't process message %v", message)
		return
	}

	if len(eventParts) < 1 {
		log.Printf("Couldn't process message %v", message)
		return
	}

	eventType := eventParts[1]

	emitterUserID := 0
	if len(eventParts) > 2 {
		emitterUserID, err = strconv.Atoi(eventParts[2])
		if err != nil {
			log.Fatal(err)
		}
	}

	receiverUserID := 0
	if len(eventParts) > 3 {
		receiverUserID, err = strconv.Atoi(eventParts[3])
		if err != nil {
			log.Fatal(err)
		}
	}

	evt := types.Event{
		Sequence:       sequenceNo,
		ReceiverUserID: receiverUserID,
		EmitterUserID:  emitterUserID,
		EventType:      eventType,
		Payload:        message + "\n",
	}

	// Push the event to the channel it will get picked up
	// by the event processor to be added to the events queue
	ctx.EventChannel <- evt
}
