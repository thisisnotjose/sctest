package users

import (
	"bufio"
	"fmt"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thisisnotjose/sctest/internal/types"
)

type TestNotificationTuple struct {
	ID  int
	evt types.Event
}

func resetContext(conn types.Connection) types.Context {
	return types.Context{
		EventQueue:     make(map[int]types.Event),
		FollowRegistry: make(map[int]map[int]bool),
		UsersPool: map[int]types.Connection{
			1: conn,
			2: conn,
			3: conn,
		},
		EventsPort:        9090,
		SubscriptionPort:  9099,
		EventChannel:      make(chan types.Event, 100),
		DeadLetterChannel: make(chan string, 100),
	}
}

func TestSendEventToUser(t *testing.T) {
	server, conn := net.Pipe()
	closeServer, closedConn := net.Pipe()
	closedConn.Close()
	closeServer.Close()

	go func() {
		reader := bufio.NewReader(server)

		for {
			payloadRaw, err := reader.ReadString('\n')

			if err != nil {
				break
			}

			fmt.Printf("Server message: %s\n", payloadRaw)
		}

		server.Close()
	}()

	ctx := resetContext(conn)

	tests := []struct {
		name                string
		messages            []TestNotificationTuple
		wantDLQ             bool
		passCloseConnection bool
	}{
		{
			"SuccessfulNotification",
			[]TestNotificationTuple{
				TestNotificationTuple{
					1,
					types.Event{
						EmitterUserID:  1,
						ReceiverUserID: 2,
						Sequence:       1,
						Payload:        "1|F|1|2\n",
					},
				},
			},
			false,
			false,
		},
		{
			"UserNotFoundInPool",
			[]TestNotificationTuple{
				TestNotificationTuple{
					7, // only users 1,2 and 3 are in the pool
					types.Event{
						EmitterUserID:  1,
						ReceiverUserID: 7,
						Sequence:       1,
						Payload:        "1|F|1|7\n",
					},
				},
			},
			true,
			false,
		},
		{
			"WritingToDeadConnection",
			[]TestNotificationTuple{
				TestNotificationTuple{
					1, // only users 1,2 and 3 are in the pool
					types.Event{
						EmitterUserID:  1,
						ReceiverUserID: 7,
						Sequence:       1,
						Payload:        "1|F|1|7\n",
					},
				},
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.messages {
				if tt.passCloseConnection {
					ctx.UsersPool[1] = closedConn
				}
				SendEventToUser(&ctx, msg.ID, msg.evt)
				if tt.wantDLQ {
					deadletter := <-ctx.DeadLetterChannel
					assert.Equal(t, msg.evt.Payload, deadletter)
				}
				if tt.passCloseConnection {
					// if we changed the context reset it for the next test
					ctx = resetContext(conn)
				}
			}
		})
	}

	server.Close()
	conn.Close()
}

func TestSendEventToAllUsers(t *testing.T) {
	server, conn := net.Pipe()
	closeServer, closedConn := net.Pipe()
	closedConn.Close()
	closeServer.Close()

	go func() {
		reader := bufio.NewReader(server)

		for {
			payloadRaw, err := reader.ReadString('\n')

			if err != nil {
				break
			}

			fmt.Printf("Server message: %s\n", payloadRaw)
		}

		server.Close()
	}()

	ctx := resetContext(conn)

	tests := []struct {
		name                string
		messages            []TestNotificationTuple
		wantDLQ             bool
		passCloseConnection bool
	}{
		{
			"SuccessfulNotification",
			[]TestNotificationTuple{
				TestNotificationTuple{
					1,
					types.Event{
						EmitterUserID:  1,
						ReceiverUserID: 2,
						Sequence:       1,
						Payload:        "1|F|1|2\n",
					},
				},
			},
			false,
			false,
		},
		{
			"WritingToCloseConnection",
			[]TestNotificationTuple{
				TestNotificationTuple{
					7, // only users 1,2 and 3 are in the pool
					types.Event{
						EmitterUserID:  1,
						ReceiverUserID: 7,
						Sequence:       1,
						Payload:        "1|F|1|7\n",
					},
				},
			},
			true,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, msg := range tt.messages {
				if tt.passCloseConnection {
					ctx.UsersPool[1] = closedConn
				}
				SendEventToAllUsers(&ctx, msg.evt)
				if tt.wantDLQ {
					deadletter := <-ctx.DeadLetterChannel
					assert.Equal(t, msg.evt.Payload, deadletter)
				}
				if tt.passCloseConnection {
					// if we changed the context reset it for the next test
					ctx = resetContext(conn)
				}
			}

		})
	}

	server.Close()
	conn.Close()
}
