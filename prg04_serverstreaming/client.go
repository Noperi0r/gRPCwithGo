package main

import (
	"context"
	"fmt"
	pb "gRPCwithGo/prg04_serverstreaming/proto"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewServerStreamingClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // context 생성, streaming은 시간 더 필요해서 10초 설정
	defer cancel()

	request := &pb.Number{Value: 5}

	// rpc 호출
	stream, err := client.GetServerResponse(ctx, request)
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}

	// 무한루프 stream 메시지 받기
	for {
		// 메시지 하나 받기
		response, err := stream.Recv()

		// 스트림 종료
		if err == io.EOF {
			fmt.Println("Stream ended.")
			break
		}

		if err != nil {
			log.Fatalf("error while receiving: %v", err)
		}

		fmt.Printf("[server to client] %s\n", response.Message)
	}

}
