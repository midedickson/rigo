package rigo

import (
	"sync"
)

// Queue represents a simple in-memory message queue
type Queue struct {
	messages []Message
	lock     sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{
		lock: sync.Mutex{},
	}
}
