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
		EventChannel:     make(chan types.Event, 1),
	}

	eventsServer := servers.NewEventsServer(handlers.NewOperatorHandler(&ctx))
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

// func main() {
// 	clientPool := make(map[int]net.Conn)
// 	followRegistry := map[int]map[int]bool{}

// 	seqNoToMessage := make(map[int][]string)

// 	go func() {
// 		lastSeqNo := 0

// 		eventListener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", eventPort))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer eventListener.Close()

// 		fmt.Printf("Listening for events on %d\n", eventPort)

// 	outer:
// 		for {
// 			conn, err := eventListener.Accept()

// 			if err != nil {
// 				log.Fatal(err)
// 			}
// 			defer conn.Close()

// 			reader := bufio.NewReader(conn)

// 			for {
// 				payloadRaw, err := reader.ReadString('\n')

// 				if err == io.EOF {
// 					conn.Close()
// 					continue outer

// 				} else if err != nil {
// 					log.Fatal(err)
// 				}

// 				payload := strings.TrimSpace(payloadRaw)

// 				fmt.Printf("Message received: %s\n", payload)

// 				payloadParts := strings.Split(payload, "|")

// 				incomingSeqNo, err := strconv.Atoi(payloadParts[0])
// 				if err != nil {
// 					log.Fatal(err)
// 				}

// 				seqNoToMessage[incomingSeqNo] = payloadParts

// 				for {

// 					nextMessage, ok := seqNoToMessage[lastSeqNo+1]
// 					delete(seqNoToMessage, lastSeqNo+1)

// 					if !ok {
// 						break
// 					}

// 					nextPayload := strings.Join(nextMessage, "|") + "\n"
// 					kind := strings.TrimSpace(nextMessage[1])

// 					switch kind {
// 					case "F":
// 						fromUserID, err := strconv.Atoi(nextMessage[2])
// 						if err != nil {
// 							log.Fatal(err)
// 						}
// 						toUserID, err := strconv.Atoi(nextMessage[3])
// 						if err != nil {
// 							log.Fatal(err)
// 						}

// 						if _, ok := followRegistry[toUserID]; !ok {
// 							followRegistry[toUserID] = make(map[int]bool)
// 						}

// 						followers, _ := followRegistry[toUserID]
// 						followers[fromUserID] = true

// 						clientConn, ok := clientPool[toUserID]
// 						if ok {
// 							fmt.Fprint(clientConn, nextPayload)
// 						}

// 					case "U":
// 						fromUserID, err := strconv.Atoi(nextMessage[2])
// 						if err != nil {
// 							log.Fatal(err)
// 						}
// 						toUserID, err := strconv.Atoi(nextMessage[3])
// 						if err != nil {
// 							log.Fatal(err)
// 						}

// 						if followers, ok := followRegistry[toUserID]; ok {
// 							delete(followers, fromUserID)
// 						}

// 					case "P":
// 						toUserID, err := strconv.Atoi(nextMessage[3])
// 						if err != nil {
// 							log.Fatal(err)
// 						}

// 						if clientConn, ok := clientPool[toUserID]; ok {
// 							fmt.Fprint(clientConn, nextPayload)
// 						}

// 					case "B":
// 						for _, clientConn := range clientPool {
// 							fmt.Fprint(clientConn, nextPayload)
// 						}

// 					case "S":
// 						fromUserID, err := strconv.Atoi(nextMessage[2])
// 						if err != nil {
// 							log.Fatal(err)
// 						}

// 						if followers, ok := followRegistry[fromUserID]; ok {
// 							for follower := range followers {
// 								clientConn, ok := clientPool[follower]
// 								if ok {
// 									fmt.Fprint(clientConn, nextPayload)
// 								}
// 							}
// 						}
// 					}

// 					lastSeqNo = lastSeqNo + 1
// 				}
// 			}
// 		}
// 	}()

// 	eventListener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", clientPort))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer eventListener.Close()

// 	fmt.Printf("Listening for client requests on %d\n", clientPort)

// 	for {
// 		conn, err := eventListener.Accept()
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		reader := bufio.NewReader(conn)

// 		userIDRaw, err := reader.ReadString('\n')

// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		userIDStr := strings.TrimSpace(userIDRaw)

// 		userID, err := strconv.Atoi(userIDStr)
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		clientPool[userID] = conn

// 		fmt.Printf("User connected: %d (%d total)\n", userID, len(clientPool))
// 	}
// }
