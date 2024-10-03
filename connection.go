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

func (c *Connection) openChannel(name string) *Channel {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exists := c.channels[name]; exists {
		return channel // Return the existing channel if it already exists
	}

	channel := NewChannel(name)
	c.channels[name] = channel
	return channel
}

func (c *Connection) getChannel(name string) (*Channel, bool) {
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
	// example command: "CHANNEL <channel>"
	if len(parts) < 2 {
		return fmt.Errorf("channel name is required for CHANNEL command")
	}
	name := parts[1]
	c.openChannel(name)
	return nil
}

func (c *Connection) handleProduceCommand(wg *sync.WaitGroup, parts []string) error {
	// message parts format: "PRODUCE <channel> <message>"
	if len(parts) < 3 {
		return fmt.Errorf("channel name and mesage content are required for PRODUCE command")
	}
	channel, exists := c.getChannel(parts[1])
	if !exists {
		return fmt.Errorf("channel '%s' not found", parts[1])
	}
	message := Message{Content: strings.Join(parts[2:], " ")}
	channel.Produce(wg, &message)
	return nil
}

func (c *Connection) handleConsumeCommand(wg *sync.WaitGroup, parts []string) (*Message, error) {
	// message parts format: "CONSUME <channel>"
	if len(parts) < 2 {
		return nil, fmt.Errorf("channel name content are required for CONSUME command")
	}
	channel, exists := c.getChannel(parts[1])
	if !exists {
		return nil, fmt.Errorf("channel '%s' not found", parts[1])
	}

	m := channel.Consume(wg)
	return m, nil
}

func (c *Connection) handleQuitCommand() error {
	// close the connection
	c.conn.Close()
	return nil
}

func readCommandFromConnection(rawConn net.Conn) (string, error) {
	buf := make([]byte, 1024)
	n, err := rawConn.Read(buf)
	if err != nil {
		return "", err
	}
	command := strings.TrimSpace(string(buf[:n]))
	if len(command) == 0 {
		return "", fmt.Errorf("empty command received")
	}
	return command, nil
}

func HandleConnection(rawConn net.Conn) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	conn := NewConnection(rawConn)
	for {
		commandString, err := readCommandFromConnection(conn.conn)
		if err != nil {
			fmt.Printf("Error reading command: %v\n", err)
			break
		}
		err = conn.RunCommand(wg, commandString)
		if err != nil {
			fmt.Printf("Error handling command: %v\n", err)
			break
		}
	}
}
