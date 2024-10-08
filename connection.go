package rigo

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

type Connection struct {
	conn     net.Conn
	mu       sync.Mutex
	channels map[string]IChannel
}

func newConnection(conn net.Conn) *Connection {
	return &Connection{
		conn:     conn,
		mu:       sync.Mutex{},
		channels: make(map[string]IChannel),
	}
}

func (c *Connection) openChannel(name string) IChannel {
	c.mu.Lock()
	defer c.mu.Unlock()

	if channel, exists := c.channels[name]; exists {
		return channel // Return the existing channel if it already exists
	}

	channel := NewChannel(name)
	c.channels[name] = channel
	return channel
}

func (c *Connection) getChannel(name string) (IChannel, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	channel, exists := c.channels[name]
	return channel, exists
}

func (c *Connection) runCommand(wg *sync.WaitGroup, commandString string) error {
	parts := SplitCommand(commandString)
	command := parts[0]

	switch command {
	case CHANNEL:
		err := c.handleChannelCommand(parts)
		if err == nil {
			c.conn.Write([]byte("OK"))
		}
		return err
	case PRODUCE:
		err := c.handleProduceCommand(wg, parts)
		if err == nil {
			c.conn.Write([]byte("OK"))
		}
		return err
	case CONSUME:
		m, err := c.handleConsumeCommand(parts)
		if err != nil {
			return err
		}
		if m != nil {
			c.conn.Write([]byte(m.Content))
		} else {
			c.conn.Write([]byte("EMPTY"))
		}
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

func (c *Connection) handleConsumeCommand(parts []string) (*Message, error) {
	// message parts format: "CONSUME <channel>"
	if len(parts) < 2 {
		return nil, fmt.Errorf("channel name content are required for CONSUME command")
	}
	channel, exists := c.getChannel(parts[1])
	if !exists {
		return nil, fmt.Errorf("channel '%s' not found", parts[1])
	}

	m := channel.Consume()
	return m, nil
}

func (c *Connection) handleQuitCommand() error {
	// close the connection
	c.conn.Close()
	return nil
}

func readCommandString(bufR *bufio.Reader) (string, error) {
	// Read until a newline delimiter is encountered
	command, err := bufR.ReadString('\n')
	if err != nil {
		return "", err
	}
	command = strings.TrimSpace(command)
	if len(command) == 0 {
		return "", errEmptyCommand
	}
	return command, nil
}

func HandleConnection(rawConn net.Conn) {
	wg := &sync.WaitGroup{}
	defer wg.Wait()

	conn := newConnection(rawConn)
	// Create a buffered reader to read the incoming stream
	reader := bufio.NewReader(conn.conn)
	for {
		commandString, err := readCommandString(reader)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed.")
				break
			}
			if err == errEmptyCommand {
				fmt.Println("Empty command received.")
				continue
			}
			rawConn.Write([]byte("err: " + err.Error()))
			fmt.Printf("Error reading command: %v\n", err)
			continue
		}
		fmt.Println("final command string: " + commandString)
		err = conn.runCommand(wg, commandString)
		if err != nil {
			rawConn.Write([]byte("err: " + err.Error()))
			fmt.Printf("Error handling command: %v\n", err)
			continue
		}
	}
}
