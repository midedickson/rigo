package rigo

import (
	"sync"
)

// Queue represents a simple in-memory message queue
type Queue struct {
	messages []Message
	lock     sync.Mutex
}

func newQueue() *Queue {
	return &Queue{
		lock: sync.Mutex{},
	}
}
