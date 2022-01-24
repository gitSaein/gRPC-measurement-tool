package main

import (
	"sync"
	"time"

	c "gRPC_measurement_tool/cmd"
	errorModel "gRPC_measurement_tool/error"
	m "gRPC_measurement_tool/measure"
	u "gRPC_measurement_tool/util"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

var (
	name    string
	option  m.Option
	startAt time.Time
)

// 프로그램 실행시 호출
func init() {
	startAt = time.Now()
	option = c.Basic()
}

func connectServer(wait *sync.WaitGroup, cmd m.Option, report *m.Report) {
	pid, opts, err, ctx, cancel := u.SetOption(cmd, startAt, report)
	report = errorModel.HandleError(err, pid, report, cmd, m.SetOption, startAt)
	defer wait.Done()
	if cancel != nil {
		defer cancel()
	}

	// Set up a connection to the server.
	startAt = time.Now()
	conn, err := grpc.DialContext(ctx, cmd.IP+":"+cmd.Port, opts...)
	if conn != nil {
		defer conn.Close() // 프로그램 종료시 conn.Close() 호출
	}
	report = errorModel.HandleError(err, pid, report, cmd, m.Dial, startAt)

	// go u.CheckDialConnection(conn, ctx, pid, startAt, report)

	// reply, err := u.CallMethod(cmd, conn, ctx)
	// report = errorModel.HandleError(err, pid, report, cmd, m.CallMethod)
	// log.Printf("message: %v", reply.GetMessage())

	report.Total = time.Since(startAt)
	m.PrintResult(report, cmd)
}

func main() {

	report := &m.Report{}

	wg := new(sync.WaitGroup)
	wg.Add(option.Tr)

	for i := 0; i < option.Tr; i++ {
		go connectServer(wg, option, report)
	}
	wg.Wait() //Go루틴 모두 끝날 때까지 대기

}
