package main

import (
	"context"
	"fmt"
	pb "gRPCwithGo/prg01_hello_gRPC/proto"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// 서버 연결
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close() // defer > 함수 끝날 때 자동 실행. 따라서 main 함수 종료 시 연결 자동 종료

	client := pb.NewMyServiceClient(conn) // 클라이언트 생성

	ctx, cancel := context.WithTimeout(context.Background(), time.Second) // Context 생성(타임아웃도 설정)
	defer cancel()

	// 요청 메시지 생성, RPC 호출
	response, err := client.MyFunction(ctx, &pb.MyNumber{Value: 4})
	if err != nil {
		log.Fatalf("could not call: %v", err)
	}

	// 결과 출력
	fmt.Printf("gRPC result: %d\n", response.Value)

}
