package main

import (
	"sync"

	"github.com/thisisnotjose/sctest/internal/handlers"
	"github.com/thisisnotjose/sctest/internal/processors"
	"github.com/thisisnotjose/sctest/internal/servers"
	"github.com/thisisnotjose/sctest/internal/types"
)

const eventPort = 9090
const clientPort = 9099

func main() {
	var wg sync.WaitGroup

	ctx := types.Context{
		EventQueue:       make(map[int]types.Event),
		FollowRegistry:   make(map[int]map[int]bool),
		UsersPool:        make(map[int]types.Connection),
		EventsPort:       9090,
		SubscriptionPort: 9099,
		EventChannel:     make(chan types.Event, 100),
	}

	eventsServer := servers.NewEventsServer(handlers.NewEventsHandler(&ctx))
	subscriptionServer := servers.NewSubscriptionServer(handlers.NewSubscriptionHandler(&ctx))

	// you can also add these one at
	// a time if you need to
	wg.Add(2)

	go func() {
		defer wg.Done()
		eventsServer.Listen(ctx.EventsPort)
	}()

	go func() {
		defer wg.Done()
		subscriptionServer.Listen(ctx.SubscriptionPort)
	}()

	go func() {
		defer wg.Done()
		processors.NewEventsProcessor(&ctx).Start()
	}()

	wg.Wait()
}
