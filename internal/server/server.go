package server

import (
	"net"

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

	// You can implement your request handling logic here.
	// For simplicity, we just send a "Hello, World!" response.
	response := "HTTP/1.1 200 OK \r\n\r\nHello, World!\r\n"
	_, err := conn.Write([]byte(response))
	if err != nil {
		s.logger.Errorf("Error writing response: %v", err)
	}
}
