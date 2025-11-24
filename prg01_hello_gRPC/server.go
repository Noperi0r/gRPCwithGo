package main

import (
	"context"
	"fmt"
	pb "gRPCwithGo/prg01_hello_gRPC/proto"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct { // 서버 구조체(클래스 정의)
	pb.UnimplementedMyServiceServer // 상속을 중괄호에 정의
}

// s *server: python self와 동일.
// req *pb.MyNumber: request
func (s *server) MyFunction(ctx context.Context, req *pb.MyNumber) (*pb.MyNumber, error) { // *pb.MyNumber, error 리턴값 두개 > 응답, 에러
	result := MyFunc(req.Value)
	fmt.Println("RPC Service result: ", result)
	return &pb.MyNumber{Value: result}, nil // 리턴값 2개(&는 주소값이고 MyNumber 객체 생성, nil은 None, 에러 없음 의미)
}

func main() {
	// 네트워크 리스너
	lis, err := net.Listen("tcp", ":50051") // := > 변수 선언과 할당 동시에 진행
	if err != nil {                         // 에러 명시적으로 체크
		log.Fatalf("listen fail: %v", err) // log.Fatalf > 에러 출력 후 프로그램 종료
	}

	s := grpc.NewServer() // grpc 서버 생성

	pb.RegisterMyServiceServer(s, &server{}) // 서비스 등록

	fmt.Println("Starting server. Listening on port 50051.")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
