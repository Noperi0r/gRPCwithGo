package main

import (
	"fmt"
	pb "gRPCwithGo/prg04_serverstreaming/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedServerStreamingServer
}

func (s *server) GetServerResponse(request *pb.Number, stream pb.ServerStreaming_GetServerResponseServer) error {
	fmt.Printf("Server processing gRPC server-streaming {%d}. \n", request.Value)

	messages := []string{
		"message #1",
		"message #2",
		"message #3",
		"message #4",
		"message #5",
	}

	// python yield처럼 하나씩 스트리밍 전송
	for _, msg := range messages {
		response := &pb.Message{
			Message: msg,
		}

		if err := stream.Send(response); err != nil {
			return err
		}
	}

	return nil // 모든 메시지 정상 streaming

}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()                          // 서버 생성
	pb.RegisterServerStreamingServer(s, &server{}) // 서비스 등록

	fmt.Println("Starting server. Listening on port 50051.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
