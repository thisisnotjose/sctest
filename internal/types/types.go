package types

import (
	"net"
)

// Connection is a server TCP connection
type Connection net.Conn

// Server will listen to a port and execute Process() after receiving a LR character
type Server interface {
	Listen(port int)
}

// ProcessorStart will start a method that runs not dependent of a TCP connection
type ProcessorStart func(int) int

// Processor will start processing events sequentially as they are added into the event queue
type Processor interface {
	Start()
}

// Handler is a function called when processing a server message
type Handler func(conn Connection, message string)

// Context is the current context of the application
type Context struct {
	SubscriptionPort int
	EventsPort       int
	UsersPool        map[int]Connection
	FollowRegistry   map[int]map[int]bool // The follower-subscriber data table
	EventQueue       map[int]Event
	EventChannel     chan Event // A ordered array of events
}

// Event holds the basic format of a follower maze event
type Event struct {
	Sequence       int
	EventType      string
	EmitterUserID  int
	ReceiverUserID int
	Payload        string
}
