package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "gRPCwithGo/prg02_bidirectional_streaming/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBidirectionalServer
}

// stream 파라미터 하나만 있음 송수신 모두 이 stream 사용
func (s *server) GetServerResponse(stream pb.Bidirectional_GetServerResponseServer) error {
	fmt.Println("Server processing gRPC bidirectional streaming.")

	for { // 메시지 cli로부터 받고 받은 메시지 돌려보내기

		message, err := stream.Recv()

		// 스트림을 종료 확인
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		fmt.Printf("Received from client: %s\n", message.Message)

		// 받은 메시지를 그대로 클라이언트에게 다시 보냄
		if err := stream.Send(message); err != nil {
			return err
		}

		fmt.Printf("Sent to client: %s\n", message.Message)
	}
}

// 4. main 함수
func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterBidirectionalServer(s, &server{})

	fmt.Println("Starting server. Listening on port 50051.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
