package main

import (
	"log"
	"net"
	"sync"

	pb "grpc-chat/proto"
	"google.golang.org/grpc"
)

type chatServer struct {
	pb.UnimplementedChatServiceServer
	mu      sync.Mutex
	streams map[string]pb.ChatService_ChatStreamServer
}

func newChatServer() *chatServer {
	return &chatServer{
		streams: make(map[string]pb.ChatService_ChatStreamServer),
	}
}

func (s *chatServer) ChatStream(stream pb.ChatService_ChatStreamServer) error {
	var user string

	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Printf("Client disconnected or error: %v", err)
			s.mu.Lock()
			delete(s.streams, user)
			s.mu.Unlock()
			return err
		}

		user = msg.User
		log.Printf("[%s]: %s", user, msg.Message)

		s.mu.Lock()
		s.streams[user] = stream
		for u, st := range s.streams {
			if u != user {
				_ = st.Send(msg) // ignore error for simplicity
			}
		}
		s.mu.Unlock()
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, newChatServer())

	log.Println("Chat server running on port 50053...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

