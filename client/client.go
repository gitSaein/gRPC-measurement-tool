package main

import (
	"bytes"
	"context"
	"flag"
	"log"
	"runtime"
	"strconv"
	"sync"
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
var count int32

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
func getGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func initialized() (uint64, []grpc.DialOption, time.Time) {
	start := time.Now()
	pid := getGID()
	log.Printf("[client-pid: %v] arrive-time: %v", pid, start)

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(interceptor.Identity{ID: pid, StartAt: start}.UnaryClient),
	}

	return pid, opts, start
}

func connectServer(wait sync.WaitGroup, done chan bool) {
	pid, opts, start := initialized()
	defer wait.Done()

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
	elapsed := time.Since(start)

	log.Printf("[client-pid: %v] message: %s, server-status: '%s', take-time: %v arrive-time: %v", pid, r.GetMessage(), conn.GetState(), elapsed, time.Now())

}

func main() {
	connect_count := 5

	var wait sync.WaitGroup
	wait.Add(connect_count)

	done := make(chan bool)
	for i := 0; i < connect_count; i++ {
		go connectServer(wait, done)
	}

	close(done)
	wait.Wait() //Go루틴 모두 끝날 때까지 대기

}
