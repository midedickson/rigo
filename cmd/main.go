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
	fmt.Printf("Starting Rigo Server on port: %s\n", opt.Port)
	server, err := rigo.NewServer(opt)
	if err != nil {
		fmt.Println("Error occured while starting server:", err)
	}
	fmt.Printf("Rigo Server Started on port: %s 🚀\n", opt.Port)

	for {
		// Accept incoming connections
		conn, err := server.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Received connection, Listening for commands... 🤝")
		// Handle every new connection in a goroutine
		go rigo.HandleConnection(conn)
	}
}
