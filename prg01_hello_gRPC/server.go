package main

import (
	"context"                              // context parameter
	"fmt"                                  // print
	pb "gRPCwithGo/prg01_hello_gRPC/proto" // protobuffer
	"log"                                  // 에러 로그
	"net"                                  // 네트워크 연결

	"google.golang.org/grpc" // grpc 모듈
)

var _ = context.TODO
var _ = fmt.Println
var _ = log.Fatal
var _ = net.Listen
var _ = pb.MyNumber{}
var _ = grpc.NewServer
