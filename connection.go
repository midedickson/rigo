package rigo

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

type Connection struct {
	conn     net.Conn
	mu       sync.Mutex
	channels map[string]*Channel
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn:     conn,
		mu:       sync.Mutex{},
		channels: make(map[string]*Channel),
	}
}

func (c *Connection) OpenChannel(name string) *Channel {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exists := c.channels[name]; exists {
		return channel // Return the existing channel if it already exists
	}

	channel := NewChannel(name)
	c.channels[name] = channel
	return channel
}

func (c *Connection) GetChannel(name string) (*Channel, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	channel, exists := c.channels[name]
	return channel, exists
}

func (c *Connection) RunCommand(wg *sync.WaitGroup, commandString string) error {
	parts := SplitCommand(commandString)
	command := parts[0]

	switch command {
	case CHANNEL:
		return c.handleChannelCommand(parts)
	case PRODUCE:
		return c.handleProduceCommand(wg, parts)
	case CONSUME:
		m, err := c.handleConsumeCommand(wg, parts)
		if err != nil {
			return err
		}
		c.conn.Write([]byte(m.Content))
		return nil
	case QUIT:
		return c.handleQuitCommand()
	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

func (c *Connection) handleChannelCommand(parts []string) error {
	if len(parts) < 2 {
		return fmt.Errorf("channel name is required for CHANNEL command")
	}
	name := parts[1]
	c.OpenChannel(name)
	return nil
}

func (c *Connection) handleProduceCommand(wg *sync.WaitGroup, parts []string) error {
	// message parts format: "PRODUCE <channel> <message>"
	if len(parts) < 3 {
		return fmt.Errorf("channel name and mesage content are required for PRODUCE command")
	}
	channel, exists := c.GetChannel(parts[1])
	if !exists {
		return fmt.Errorf("channel '%s' not found", parts[1])
	}
	message := Message{Content: strings.Join(parts[2:], " ")}
	channel.queue.Produce(wg, &message)
	return nil
}

func (c *Connection) handleConsumeCommand(wg *sync.WaitGroup, parts []string) (*Message, error) {
	// message parts format: "CONSUME <channel>"
	if len(parts) < 2 {
		return nil, fmt.Errorf("channel name content are required for CONSUME command")
	}
	channel, exists := c.GetChannel(parts[1])
	if !exists {
		return nil, fmt.Errorf("channel '%s' not found", parts[1])
	}

	m := channel.queue.Consume(wg)
	return m, nil
}

func (c *Connection) handleQuitCommand() error {
	// close the connection
	c.conn.Close()
	return nil
}
