package rigo

import (
	"fmt"
	"sync"
)

// Message represents a simple message structure
type Message struct {
	ID      int
	Content string
}

// MessageQueue represents a simple in-memory message queue
type MessageQueue struct {
	messages []Message
	lock     sync.Mutex
}

// Produce adds a message to the message queue
func (mq *MessageQueue) Produce(wg *sync.WaitGroup, message Message) {
	defer wg.Done()
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.messages = append(mq.messages, message)
	fmt.Printf("Produced: %+v\n", message)
}

// Consume removes and returns a message from the message queue
func (mq *MessageQueue) Consume(wg *sync.WaitGroup, result chan<- Message) {
	defer wg.Done()
	mq.lock.Lock()
	defer mq.lock.Unlock()

	if len(mq.messages) == 0 {
		fmt.Println("Queue is empty")
		return
	}

	message := mq.messages[0]
	mq.messages = mq.messages[1:]
	fmt.Printf("Consumed: %+v\n", message)
	result <- message
}
