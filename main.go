package main

import (
	"github.com/midedickson/rigo/pkg/rigo"
)

func main() {
	server := rigo.NewServer()
	server.ListenAndServe()
}
