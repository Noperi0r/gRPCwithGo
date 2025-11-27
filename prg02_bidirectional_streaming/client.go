package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	pb "gRPCwithGo/prg02_bidirectional_streaming/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 서버 연결
	conn, err := grpc.NewClient("localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewBidirectionalClient(conn) // 클라이언트 생성

	// context 생성
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stream, err := client.GetServerResponse(ctx)
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}

	// 5. 전송할 메시지 리스트
	messages := []string{
		"message #1",
		"message #2",
		"message #3",
		"message #4",
		"message #5",
	}

	// 6. 송신과 수신을 동시에 처리하기 위해 고루틴 사용
	// 고루틴 = Go의 경량 스레드 (Python의 asyncio 개념)

	// 6-1. 수신용 고루틴 (별도 스레드에서 실행)
	waitc := make(chan struct{}) // 동기화용 채널
	go func() {
		for {
			// 서버로부터 메시지 받기
			response, err := stream.Recv()
			if err == io.EOF {
				// 서버가 스트림 종료
				close(waitc) // 채널 닫기 (메인 고루틴에 알림)
				return
			}
			if err != nil {
				log.Fatalf("error while receiving: %v", err)
			}
			fmt.Printf("[server to client] %s\n", response.Message)
		}
	}()

	// 6-2. 송신 (메인 고루틴에서 실행)
	for _, msg := range messages {
		message := &pb.Message{
			Message: msg,
		}

		// 서버에 메시지 전송
		if err := stream.Send(message); err != nil {
			log.Fatalf("error while sending: %v", err)
		}

		fmt.Printf("[client to server] %s\n", message.Message)

		// 약간의 지연 (선택사항)
		time.Sleep(500 * time.Millisecond)
	}

	// 7. 모든 메시지 전송 완료, 스트림 닫기
	stream.CloseSend()

	// 8. 수신 고루틴이 종료될 때까지 대기
	<-waitc // 채널에서 값 받기 (채널이 닫힐 때까지 대기)
}
