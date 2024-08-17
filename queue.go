package rigo

import (
	"sync"
)

// MessageQueue represents a simple in-memory message queue
type MessageQueue struct {
	Name     string
	messages []Message
	lock     sync.Mutex
}

var queueTable map[string]*MessageQueue

func Queue(queueName string) *MessageQueue {
	messageQueue := &MessageQueue{Name: queueName}
	if _, ok := queueTable[queueName]; ok {
		return queueTable[queueName]
	} else {
		queueTable[queueName] = messageQueue
	}
	return messageQueue
}
