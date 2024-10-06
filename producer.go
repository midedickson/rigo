package rigo

import (
	"fmt"
	"sync"
)

// Produce adds a message to the message queue
func (mq *Queue) produce(wg *sync.WaitGroup, message *Message) {
	defer wg.Done()
	mq.lock.Lock()
	defer mq.lock.Unlock()
	mq.messages = append(mq.messages, *message)
	fmt.Printf("Produced: %+v\n", message)
}

func (channel *Channel) Produce(wg *sync.WaitGroup, message *Message) {
	wg.Add(1)
	go channel.queue.produce(wg, message)
}
