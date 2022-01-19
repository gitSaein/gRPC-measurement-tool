package main

import (
	"context"
	"log"
	"sync"
	"time"

	cmdflag "gRPC_measurement_tool/cmdflag"
	rpc "gRPC_measurement_tool/rpc"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

var (
	name    string
	command cmdflag.Command
)

// 프로그램 실행시 호출
func init() {
	command = cmdflag.Basic()

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

func connectServer(wait sync.WaitGroup, done chan bool, cmd cmdflag.Command) {
	pid, opts, start := cmdflag.GetInitSetting(cmd)
	defer wait.Done()

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close() // 프로그램 종료시 conn.Close() 호출

	ctx, cancel := cmdflag.GetInitTimeout(cmd)
	if cancel != nil {
		defer cancel()
	}
	reply, err := cmdflag.GetInitCall(cmd, conn, ctx)
	if reply == nil {
		log.Fatalf("could not greet: %v", err)
	}
	elapsed := time.Since(start)

	log.Printf("[client-pid: %v] message: %s, server-status: '%s', take-time: %v arrive-time: %v", pid, reply.GetMessage(), conn.GetState(), elapsed, time.Now())

}

func main() {

	if command.Tr < 0 {
		log.Fatalf("invalid total request value")
	}

	var wait sync.WaitGroup
	wait.Add(command.Tr)

	done := make(chan bool)
	for i := 0; i < command.Tr; i++ {
		go connectServer(wait, done, command)
	}

	close(done)
	wait.Wait() //Go루틴 모두 끝날 때까지 대기

}
