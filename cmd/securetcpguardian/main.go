package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/sigsegv1989/securetcpguardian/internal/server"
)

func main() {
	// Define command-line flags for server IP address and port
	listenAddress := flag.String("listen-address", "localhost", "Server listen address (default: localhost)")
	listenPort := flag.Int("listen-port", 8080, "Server listen port (default: 8080)")
	flag.Parse()

	// Create a channel to listen for OS signals (e.g., Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the TCP server
	server := server.NewTCPServer(fmt.Sprintf("%s:%d", *listenAddress, *listenPort))
	go server.Start()

	fmt.Printf("TCP server started on %s:%d. Press Ctrl+C to exit.\n", *listenAddress, *listenPort)
	<-sigChan // Wait for a signal to gracefully shut down
	server.Stop()
}
