package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/midedickson/rigo"
)

func main() {
	// read the configuration
	godotenv.Load()

	// initiate the options
	opt := rigo.NewOptions()

	// Start the rigo server, abstract the initialization
	server, err := rigo.NewServer(opt)
	if err != nil {
		fmt.Println("Error occured while starting server:", err)
	}
	for {
		// Accept incoming connections
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		// Handle the connection in a new goroutine
		go rigo.HandleConnection(conn)
	}
}
