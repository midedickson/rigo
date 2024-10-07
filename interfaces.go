package rigo

import "sync"

type IChannel interface {
	Consume() *Message
	Produce(wg *sync.WaitGroup, message *Message)
}

type IQueue interface {
	consume() *Message
	produce(wg *sync.WaitGroup, message *Message)
}
