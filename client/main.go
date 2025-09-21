package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "grpc-chat/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)
	stream, err := client.ChatStream(context.Background())
	if err != nil {
		log.Fatalf("Failed to open stream: %v", err)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	user, _ := reader.ReadString('\n')
	user = user[:len(user)-1]

	// Receiving goroutine
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Receive error: %v", err)
				break
			}
			fmt.Printf("\n[%s]: %s\n", msg.User, msg.Message)
		}
	}()

	// Sending loop
	for {
		fmt.Print("You: ")
		text, _ := reader.ReadString('\n')
		stream.Send(&pb.ChatMessage{
			User:      user,
			Message:   text,
			Timestamp: time.Now().Unix(),
		})
	}
}

