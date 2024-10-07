package rigo

import "sync"

type Channel struct {
	mu    sync.Mutex
	Name  string
	queue IQueue
}

func NewChannel(name string) *Channel {
	return &Channel{
		Name:  name,
		mu:    sync.Mutex{},
		queue: newQueue(),
	}
}
