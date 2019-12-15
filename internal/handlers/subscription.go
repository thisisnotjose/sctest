package handlers

import (
	"fmt"
	"log"
	"strconv"

	"github.com/thisisnotjose/sctest/internal/types"
)

// NewSubscriptionHandler Registers a new user into the system
func NewSubscriptionHandler(ctx *types.Context) types.Handler {
	return func(conn types.Connection, message string) {
		subscriptionHandler(ctx, conn, message)
	}
}

func subscriptionHandler(ctx *types.Context, conn types.Connection, message string) {
	userID, err := strconv.Atoi(message)
	if err != nil {
		log.Fatal(err)
	}

	ctx.UsersPool[userID] = conn

	fmt.Printf("User connected: %d (%d total)\n", userID, len(ctx.UsersPool))
}
