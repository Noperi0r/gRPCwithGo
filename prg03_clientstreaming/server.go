package main

import (
	"fmt"
	"io"
	"log"
	"net"

	pb "gRPCwithGo/prg03_clientstreaming/proto"

	"google.golang.org/grpc"
)

// 서버 구조체
type server struct {
	pb.UnimplementedClientStreamingServer
}

// stream으로 메시지를 여러 개 받고, 마지막에 응답 1개만 리턴
func (s *server) GetServerResponse(stream pb.ClientStreaming_GetServerResponseServer) error {
	fmt.Println("Server processing gRPC client-streaming.")

	count := int32(0)

	// 메시지 클라로부터 계속 받기
	for {
		message, err := stream.Recv()

		if err == io.EOF {
			// 최종 전송
			return stream.SendAndClose(&pb.Number{
				Value: count,
			})
		}

		if err != nil {
			return err
		}

		// 메시지 받았음 (실제로는 아무것도 안 하고 카운트만 증가)
		fmt.Printf("Received from client: %s\n", message.Message)
		count++
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterClientStreamingServer(s, &server{})

	fmt.Println("Starting server. Listening on port 50051.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
