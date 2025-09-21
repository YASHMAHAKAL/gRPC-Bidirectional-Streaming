# gRPC Chat Demo

A real-time chat application demonstrating bidirectional streaming with gRPC, featuring both Go and Python client implementations.

## Overview

This project showcases a simple chat system built with gRPC that allows multiple clients to connect and exchange messages in real-time. The server is implemented in Go, and clients are available in both Go and Python.

## Features

- **Real-time messaging**: Bidirectional streaming allows instant message delivery
- **Multi-client support**: Multiple users can join the chat simultaneously
- **Cross-language compatibility**: Go server with Go and Python clients
- **Simple protocol**: Clean protobuf definition for easy extension

## Architecture

```
┌─────────────────┐    gRPC Stream    ┌──────────────────┐
│   Go Client     │◄─────────────────►│                  │
└─────────────────┘                   │    Go Server     │
┌─────────────────┐    gRPC Stream    │   (Port 50053)   │
│  Python Client  │◄─────────────────►│                  │
└─────────────────┘                   └──────────────────┘
┌─────────────────┐    gRPC Stream
│   Go Client     │◄─────────────────►
└─────────────────┘
```

## Protocol Definition

The chat service uses a simple protobuf definition:

```protobuf
service ChatService {
  rpc ChatStream(stream ChatMessage) returns (stream ChatMessage);
}

message ChatMessage {
  string user = 1;
  string message = 2;
  int64 timestamp = 3;
}
```

## Prerequisites

### For Go components:
- Go 1.24.4 or higher
- Protocol Buffer compiler (`protoc`)

### For Python client:
- Python 3.x
- grpcio and grpcio-tools packages

## Installation & Setup

### 1. Clone the repository
```bash
git clone <repository-url>
cd gRPC-demo
```

### 2. Install Go dependencies
```bash
go mod download
```

### 3. Install Python dependencies (for Python client)
```bash
pip install grpcio grpcio-tools
```

### 4. Generate protobuf files (if needed)

For Go:
```bash
protoc --go_out=. --go_grpc_out=. proto/chat.proto
```

For Python:
```bash
python -m grpc_tools.protoc -I proto --python_out=pyclient --grpc_python_out=pyclient proto/chat.proto
```

## Running the Application

### 1. Start the Server
```bash
# Option 1: Run directly
go run server/main.go

# Option 2: Build and run
go build -o bin/server server/main.go
./bin/server
```

The server will start listening on `localhost:50053`.

### 2. Start Client(s)

#### Go Client
```bash
# Option 1: Run directly  
go run client/main.go

# Option 2: Build and run
go build -o bin/client client/main.go
./bin/client
```

#### Python Client
```bash
cd pyclient
python client.py
```

### 3. Start Chatting
1. When prompted, enter your username
2. Type messages and press Enter to send
3. Messages from other connected users will appear automatically
4. Use Ctrl+C to exit

## Project Structure

```
gRPC-demo/
├── README.md
├── go.mod                 # Go module definition
├── go.sum                 # Go dependency checksums
├── bin/                   # Compiled binaries
│   ├── client
│   └── server
├── proto/                 # Protocol buffer definitions
│   ├── chat.proto         # Service definition
│   ├── chat.pb.go         # Generated Go code
│   ├── chat_grpc.pb.go    # Generated Go gRPC code
│   ├── chat_pb2.py        # Generated Python code
│   └── chat_pb2_grpc.py   # Generated Python gRPC code
├── server/                # Go server implementation
│   └── main.go
├── client/                # Go client implementation
│   └── main.go
└── pyclient/              # Python client implementation
    ├── client.py
    ├── chat_pb2.py        # Generated Python code
    └── chat_pb2_grpc.py   # Generated Python gRPC code
```

## How It Works

### Server Implementation
- **Concurrent Streams**: Maintains a map of active client streams
- **Message Broadcasting**: Forwards messages from one client to all other connected clients
- **Connection Management**: Automatically handles client connections and disconnections
- **Thread Safety**: Uses mutexes to safely manage shared state

### Client Implementation
- **Bidirectional Streaming**: Maintains persistent connection for real-time communication
- **Concurrent Operations**: Separate goroutines/threads for sending and receiving messages
- **User Interface**: Simple command-line interface for message input
- **Error Handling**: Graceful handling of connection issues

## Example Usage

```bash
# Terminal 1 - Start server
$ go run server/main.go
2024/09/21 Server listening on :50053

# Terminal 2 - Go client
$ go run client/main.go
Enter your name: Alice
You: Hello everyone!

# Terminal 3 - Python client  
$ cd pyclient && python client.py
Enter your username: Bob
You: Hi Alice!

# Back to Terminal 2
[Bob]: Hi Alice!
You: Nice to meet you Bob!

# Back to Terminal 3  
[Alice]: Nice to meet you Bob!
You: 
```

## Development Notes

### Extending the Protocol
To add new features, modify `proto/chat.proto` and regenerate the language-specific code:

```bash
# Regenerate Go code
protoc --go_out=. --go_grpc_out=. proto/chat.proto

# Regenerate Python code
python -m grpc_tools.protoc -I proto --python_out=pyclient --grpc_python_out=pyclient proto/chat.proto
```

### Adding New Clients
The server accepts connections from any gRPC client implementing the `ChatService`. You can easily add clients in other languages by:
1. Generating protobuf code for your target language
2. Implementing the bidirectional streaming logic
3. Connecting to `localhost:50053`

## Dependencies

### Go Dependencies
- `google.golang.org/grpc` - gRPC Go implementation
- `google.golang.org/protobuf` - Protocol buffer support

### Python Dependencies  
- `grpcio` - gRPC Python implementation
- `grpcio-tools` - Protocol buffer compiler for Python

## Troubleshooting

### Common Issues

1. **Port already in use**: Change the port in both server and client code
2. **Connection refused**: Ensure the server is running before starting clients
3. **Import errors in Python**: Make sure you're running from the correct directory and have generated protobuf files

### Server Logs
The server logs all received messages and connection events to help with debugging.

## License

This project is a demonstration/educational example. Feel free to use and modify as needed.
