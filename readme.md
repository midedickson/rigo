# Rigo Queue Server

**Rigo** started as "Redis In GO," but it has evolved into a robust and standalone queueing system designed for modern applications. This system can be used in two ways:

- As a **standalone server** accessed via TCP clients.
- Embedded within applications, giving developers complete control over queue management.

## Feature Goals

- **Standalone Queue Server**: Can be hosted as a separate service and managed via a TCP client.
- **Embedded Queue**: Can be installed locally into an application and fully managed by developers.
- **Message Limit per Channel**: Supports buffering with message limits for each channel.
- **Connection Tracking**: Tracks and records the number of active connections to the queue server.

## Getting Started

To test the current queue system, you can run the example provided.

1. Clone the repository:

   ```bash
   git clone https://github.com/midedickson/rigo.git
   cd rigo-server
   ```

2. Run the example:

   ```bash
   go run cmd/main.go
   ```

This will run the basic example of the Rigo queue system in action.

## Usage

### Standalone Server Mode

To use Rigo as a standalone server, you can launch it and connect via a TCP client (like the [Rigo Client CLI](https://github.com/midedickson/rigo-client)).

1. Start the server:

   ```bash
   go run cmd/main.go
   ```

2. Connect to the server using a TCP client. For example, you can use the **Rigo Client CLI** to send and receive messages from the queue.

### Embedded Queue Mode

Rigo can also be integrated directly into your Go application as a local queue system:

1. Import the Rigo queue package in your Go application:

   ```go
   import "github.com/midedickson/rigo"
   ```

2. Use the queue in your application:

   ```go
   channel := rigo.NewChannel()
   channel.Produce("channel1", "message")
   message := queue.Consume("channel1")
   fmt.Println("Dequeued message:", message)
   ```

This allows you to fully manage the queue locally without needing an external server.

## Key Features

- **Message Limit per Channel**: You can set a buffer size to limit the number of messages in each queue channel.
- **Connection Count Recording**: The system records and tracks the number of active TCP connections to the server, providing insights into load and usage.

## Example Workflow

1. Start the server:

   ```bash
   go run server/main.go
   ```

2. Use the Rigo Client CLI to create a channel and enqueue a message:

   ```bash
   ./rigo-client 127.0.0.1 8080
   127.0.0.1:8080> CHANNEL job123
   Response: OK
   127.0.0.1:8080> PRODUCE job123 message
   Response: OK


   ```

3. Dequeue the message:

   ```bash
   127.0.0.1:8080> CONSUME job123
   Response: message
   ```

## Future Ideas

- **Standalone Server**: Continue developing the server as a standalone, distributed system.
- **Embedded Queue**: Extend the embedded queue capabilities to provide additional control and flexibility for developers.
- **Enhanced Buffering**: Add more control over buffering, including the ability to dynamically resize message limits per channel.
- **Advanced Metrics**: Track and monitor server metrics, including connection counts and message throughput.

## Related Projects

- **Rigo Client CLI**: Interact with the Rigo server via a command-line interface. [Rigo Client CLI](https://github.com/midedickson/rigo-client)
