# Real-Time Broadcasting Chatroom (Go Concurrency)

A real-time chatroom implementation using Go's concurrency features (goroutines, channels, and mutexes).

## Files

- `server.go` - TCP server with concurrent client handling and broadcasting
- `client.go` - TCP client with concurrent send/receive loops
- `README.md` - This file

## Features

- Real-time message broadcasting to all connected clients
- Join/leave notifications for all users
- No self-echo (clients don't receive their own messages from server)
- Thread-safe client management using Mutex
- Concurrent message handling using goroutines and channels
- Graceful disconnection handling
- Support for multiple simultaneous clients

## Architecture

### Server

- Maintains a thread-safe map of connected clients
- Each client has a dedicated outgoing message channel
- Uses goroutines for concurrent client handling
- Broadcasts messages to all clients except the sender

### Client

- Sends username on connection
- Runs two concurrent goroutines:
  - Receiver: listens for broadcasts from server
  - Sender: reads user input and sends to server
- Displays own messages locally without server echo

## Prerequisites

- Go 1.16 or higher

## Running the Chatroom

### Step 1: Start the Server

```bash
go run server.go
```

You should see:

```
Chat server started on :1234
```

### Step 2: Start Client(s)

In a new terminal window:

```bash
go run client.go
```

You'll be prompted to enter a username. After that, you can start typing messages.

**Open multiple clients** in different terminals to simulate real-time chat between users.

### Step 3: Send Messages

Type your message and press Enter. Your message appears locally, and all other clients receive it instantly.

### Step 4: Exit

Type `exit` or press `Ctrl+C` to quit the client.

## Example Usage

**Terminal 1 (Server):**

```
$ go run server.go
Chat server started on :1234
```

**Terminal 2 (Client - Alice):**

```
$ go run client.go
Enter your username: Alice
Welcome to the chatroom, Alice!
Type 'exit' to quit or press Ctrl+C.
------------------------------------
You: Hello everyone!
Bob: Hi Alice!
User Charlie joined
Charlie: Hey folks!
```

**Terminal 3 (Client - Bob):**

```
$ go run client.go
Enter your username: Bob
Welcome to the chatroom, Bob!
Type 'exit' to quit or press Ctrl+C.
------------------------------------
User Bob joined
Alice: Hello everyone!
You: Hi Alice!
User Charlie joined
Charlie: Hey folks!
```

## Technical Details

- **Concurrency:** Each client connection runs in its own goroutine
- **Broadcasting:** Server uses channels to send messages to clients concurrently
- **Thread Safety:** Mutex protects shared client list from race conditions
- **No Self-Echo:** Clients print their own messages locally; server excludes sender from broadcasts
