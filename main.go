package main

import (
	"github.com/midedickson/rigo/pkg/rigo"
)

func main() {
	// Start the rigo server, abstract the initialization
	server := rigo.NewServer()
	server.ListenAndServe()
}
