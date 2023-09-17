package server

import (
	"context"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// TCPServer represents a simple TCP server.
type TCPServer struct {
	address string
	running bool
	logger  *log.Logger
}

// NewTCPServer creates a new TCPServer instance.
func NewTCPServer(address string) *TCPServer {
	// Create a new Logger instance
	logger := log.New()

	// Set the logger format to include timestamps
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	return &TCPServer{
		address: address,
		running: false,
		logger:  logger,
	}
}

// Start starts the TCP server.
func (s *TCPServer) Start() {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		s.logger.Errorf("Error: %v", err)
		return
	}
	defer listener.Close()

	s.running = true
	s.logger.Infof("Listening on %s...", s.address)

	for s.running {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Errorf("Error accepting connection: %v", err)
			continue
		}

		go s.handleRequest(conn)
	}
}

// Stop stops the TCP server.
func (s *TCPServer) Stop() {
	s.running = false
}

func (s *TCPServer) handleRequest(conn net.Conn) {
	defer conn.Close()

	// Handle the incoming request here
	s.logger.Infof("Accepted connection from %s", conn.RemoteAddr())

	// Create a context with a 120-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Second)
	defer cancel()

	// Create a channel to signal when the operation is complete
	done := make(chan struct{})

	func() {
		defer close(done)

		// Create a buffer for reading messages
		buf := make([]byte, 1024)

		for {
			select {
			case <-ctx.Done():
				err := ctx.Err()
				if err == context.DeadlineExceeded {
					s.logger.Infof("Connection from %s timed out", conn.RemoteAddr())
				} else {
					s.logger.Errorf("Context error: %v", err)
				}
				return
			default:
				// Set a read deadline for each message read
				conn.SetReadDeadline(time.Now().Add(120 * time.Second))

				n, err := conn.Read(buf)
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						s.logger.Infof("Read timeout from %s", conn.RemoteAddr())
					} else {
						s.logger.Errorf("Error reading from client: %v", err)
					}
					return
				}

				clientMessage := string(buf[:(n - 1)])
				s.logger.Infof("Received message from client: %s", clientMessage)

				//response := "HTTP/1.1 200 OK \r\n\r\n" + clientMessage + "!\r\n"
				_, err = conn.Write([]byte(clientMessage + "!\n"))
				if err != nil {
					s.logger.Errorf("Error writing response: %v", err)
				}

				// reset the read deadline for the next message
				conn.SetReadDeadline(time.Time{})
			}
		}
	}()

	<-done
}
