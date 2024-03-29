package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/midedickson/rigo/pkg/rigo"
)

func main() {
	// plan here is to remove main function and convert to mqueue package
	// may or may not run on a server.
	var wg sync.WaitGroup
	messageQueue := rigo.MessageQueue{}
	result := make(chan rigo.Message, 5)

	// Start the producer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			message := rigo.Message{ID: i, Content: fmt.Sprintf("Message %d", i)}
			wg.Add(1)
			go messageQueue.Produce(&wg, message)
			time.Sleep(time.Second)
		}
	}()

	// Start the consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			wg.Add(1)
			go messageQueue.Consume(&wg, result)
			time.Sleep(time.Second)
		}
	}()

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the result channel after all consumers are done
	close(result)
}
