package rigo

import "sync"

type Channel struct {
	mu    sync.Mutex
	Name  string
	queue *Queue
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name:  name,
		mu:    sync.Mutex{},
		queue: NewQueue(),
	}
}
