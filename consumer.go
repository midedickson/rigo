package rigo

import (
	"fmt"
)

// Consume removes and returns a message from the message queue
func (mq *Queue) consume() *Message {
	mq.lock.Lock()
	defer mq.lock.Unlock()

	if len(mq.messages) == 0 {
		fmt.Println("Queue is empty")
		return nil
	}

	message := mq.messages[0]
	mq.messages = mq.messages[1:]
	fmt.Printf("Consumed: %+v\n", message)
	return &message
}

func (channel *Channel) Consume() *Message {
	return channel.queue.consume()
}
