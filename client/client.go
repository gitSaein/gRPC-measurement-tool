package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "grpc/protos"
	rpc "grpc/rpc"

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

func main() {

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close() // 프로그램 종료시 conn.Close() 호출

	//server health check
	client := rpc.NewGrpcHealthClient(conn)

	ok, err := client.Check(context.Background())

	if !ok || err != nil {
		log.Printf("can't connect grpc server: %v, code: %v\n", err, grpc.Code(err))
	} else {
		log.Println("connect the grpc server successfully")
	}

	c := pb.NewGreeterClient(conn)
	log.Printf("connected status: %v", conn.GetState())
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()

	// 서버의 rpc 호출
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.GetMessage())

}
