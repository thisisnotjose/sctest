package servers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/thisisnotjose/sctest/internal/types"
)

type eventsServer struct {
	handler types.Handler
}

// NewEventsServer returns an instance of a server
func NewEventsServer(handler types.Handler) types.Server {
	return eventsServer{
		handler: handler,
	}
}

// Listen will start the server in the given port, handler will be called after each control character
func (s eventsServer) Listen(port int) {
	eventListener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal(err)
	}
	defer eventListener.Close()

	fmt.Printf("Listening for operators on %d\n", port)

outer:
	for {
		conn, err := eventListener.Accept()

		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		reader := bufio.NewReader(conn)

		for {
			payloadRaw, err := reader.ReadString('\n')

			if err == io.EOF {
				conn.Close()
				continue outer

			} else if err != nil {
				log.Fatal(err)
			}

			payload := strings.TrimSpace(payloadRaw)

			fmt.Printf("Message received: %s\n", payload)

			s.handler(conn, payload)
		}
	}
}
