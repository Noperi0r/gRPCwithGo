package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "gRPCwithGo/prg03_clientstreaming/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewClientStreamingClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// rpc 호출
	stream, err := client.GetServerResponse(ctx)
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}

	messages := []string{
		"message #1",
		"message #2",
		"message #3",
		"message #4",
		"message #5",
	}

	// 메시지를 하나씩 스트리밍 전송
	for _, msg := range messages {
		message := &pb.Message{
			Message: msg,
		}

		if err := stream.Send(message); err != nil {
			log.Fatalf("error while sending: %v", err)
		}

		fmt.Printf("[client to server] %s\n", message.Message)
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving: %v", err)
	}

	fmt.Printf("[server to client] %d\n", response.Value)
}
