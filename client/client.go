package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	cmdflag "gRPC_measurement_tool/cmdflag"
	errorModel "gRPC_measurement_tool/error"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

var (
	name    string
	command cmdflag.Command
	startAt time.Time
)

// 프로그램 실행시 호출
func init() {
	startAt = time.Now()
	command = cmdflag.Basic()

}

func checkcheck(conn *grpc.ClientConn, ctx context.Context, pid uint64) {
	for {
		is_changed_status := conn.WaitForStateChange(ctx, conn.GetState())
		if is_changed_status {
			currentState := conn.GetState()

			elapsed := time.Since(startAt)
			log.Printf("[client-pid: %v] server-status: '%s', take-time: %s, arrive-time: %v", pid, currentState, elapsed, time.Now())

			if currentState == connectivity.Ready {
				break
			}
			if currentState == connectivity.Shutdown || currentState == connectivity.TransientFailure {
				break
			}
		}

	}
}

func connectServer(wait sync.WaitGroup, done chan bool, cmd cmdflag.Command) {
	pid, opts, err := cmdflag.GetInitSetting(cmd, startAt)
	errorModel.CheckErrorState(err, pid)
	defer wait.Done()

	ctx, cancel := cmdflag.GetInitTimeout(cmd)
	str := fmt.Sprintf("%v", pid)
	md := metadata.Pairs("pid", str)
	ctx = metadata.NewOutgoingContext(ctx, md)

	if cancel != nil {
		defer cancel()
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, opts...)
	defer conn.Close() // 프로그램 종료시 conn.Close() 호출
	errorModel.CheckErrorState(err, pid)

	go checkcheck(conn, ctx, pid)

	reply, err := cmdflag.GetInitCall(cmd, conn, ctx)
	errorModel.CheckErrorState(err, pid)

	elapsed := time.Since(startAt)

	log.Printf("Result: [client-pid: %v] [take-time: %v] [server-status: '%s'] [message: %s] ", pid, elapsed, conn.GetState(), reply.GetMessage())

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
