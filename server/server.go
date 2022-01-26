package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "gRPC_measurement_tool/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	interceptor "gRPC_measurement_tool/interceptors"
)

const (
	port = ":50051"
)

var listener *net.Listener

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

type Server struct {
	proto           string
	addr            string
	networkListener net.Listener
	mu              sync.Mutex
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {

	dd, ok := ctx.Deadline()
	if !ok {
		log.Printf("%v", dd)
	}
	// deadline
	for i := 0; i < 5; i++ {

		err := ctx.Err()
		if err != nil {
			log.Printf("Deadline Error: %v", err)

			if ctx.Err() == context.DeadlineExceeded {
				return nil, status.Errorf(codes.DeadlineExceeded, "HelloworldService.SayHello DeadlineExceeded")
			}
		}

	}

	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func ListenAndGrpcServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(interceptor.UnaryServerInterceptor))
	tls := false
	if tls {
		serverCrt := "../cert/server.crt"
		serverPem := "../cert/server.pem"
		creds, sslErr := credentials.NewServerTLSFromFile(serverCrt, serverPem)
		if sslErr != nil {
			log.Fatalf("failed to load creds: %v", sslErr)
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)

	pb.RegisterGreeterServer(s, &server{}) // helloworld_grpc.pb.go 에 있음
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { //grpc 서버 시작
		log.Fatalf("failed to serve: %v", err)
	}

}

func main() {
	ListenAndGrpcServer()
}
