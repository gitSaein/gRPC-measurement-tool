package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strconv"
	"time"

	pb "gRPC_measurement_tool/protos"
	rpc "gRPC_measurement_tool/rpc"

	interceptor "gRPC_measurement_tool/interceptors"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

var (
	name string
)

// 프로그램 실행시 호출
func init() {
	flag.StringVar(&name, "name", defaultName, "input name") // 커맨드 라인 명령: cmd> *.exe -name [value] : https://gobyexample.com/command-line-flags
	flag.Parse()                                             //  // 커맨드 라인 명령 시작
}

func CheckHttpHeader(ctx context.Context) {

}

func CheckServerStatus(conn *grpc.ClientConn) {
	client := rpc.NewGrpcHealthClient(conn)

	for {
		ok, err := client.Check(context.Background())

		if !ok || err != nil {
			log.Panicf("can't connect grpc server: %v, code: %v\n", err, grpc.Code(err))
		}
	}

}

func initialized() []grpc.DialOption {
	start := time.Now()
	pid := strconv.Itoa(os.Getpid())

	log.Printf("[client-pid: %v] start at. %s", pid, start)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.Identity{ID: pid, StartAt: start}.UnaryClient),
	}

	return opts
}

func main() {

	opts := initialized()

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close() // 프로그램 종료시 conn.Close() 호출

	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Second)
	defer cancel()

	// 서버의 rpc 호출
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})

	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())

}
