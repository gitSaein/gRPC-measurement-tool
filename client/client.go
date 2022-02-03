package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	c "gRPC_measurement_tool/cmd"
	u "gRPC_measurement_tool/config"
	errorModel "gRPC_measurement_tool/handler"
	m "gRPC_measurement_tool/measure"

	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

var (
	name        string
	option      m.Option
	startInitAt time.Time
)

// 프로그램 실행시 호출
func init() {
	startInitAt = time.Now()
	option = c.Basic()
}

func job(wait *sync.WaitGroup, cmd m.Option, worker *m.Worker) {
	startAt := time.Now()
	job := &m.Job{JId: u.GetID()}
	setting := u.SettingOptions(cmd)
	errorModel.HandleReponse(setting.Error, worker, job, cmd, m.SetOption, startAt)
	defer func() {
		worker.Jobs = append(worker.Jobs, job)
		wait.Done()
		setting.CancelFunc()
	}()

	// Set up a connection to the server.
	startAt = time.Now()
	conn, err := grpc.DialContext(setting.Context, cmd.Target, setting.Options...)
	errorModel.HandleReponse(err, worker, job, cmd, m.DialOpen, startAt)

	if conn != nil {
		defer func() {
			startAt = time.Now()
			err = conn.Close()
			errorModel.HandleReponse(err, worker, job, cmd, m.DialClose, startAt)
		}()

	}
}

func Worker(wait *sync.WaitGroup, report *m.Report, cmd m.Option) {
	wg := new(sync.WaitGroup)
	wg.Add(option.RPS)

	worker := &m.Worker{}
	worker.WId = u.GetID()
	defer func() {
		report.Workers = append(report.Workers, worker)
		wait.Done()
	}()
	for i := 0; i < option.RPS; i++ {
		go job(wg, option, worker)
	}

	wg.Wait() //Go루틴 모두 끝날 때까지 대기

}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // CPU 개수를 구한 뒤 사용할 최대 CPU 개수 설정
	fmt.Println(runtime.GOMAXPROCS(0))   // 설정 값 출력
	log.Println("Measure start...")

	report := &m.Report{}
	defer func() {
		report.Total = time.Since(startInitAt)
		m.PrintResult(report, option)
		log.Println("Measure end...")

	}()

	tick := time.Tick(1 * time.Second)
	boom := time.After(10 * time.Second)

	select {
	case <-tick:
		log.Println("tick")

		wg := new(sync.WaitGroup)
		wg.Add(option.WorkerCnt)

		for i := 0; i < option.WorkerCnt; i++ {
			go Worker(wg, report, option)
		}

		wg.Wait() //Go루틴 모두 끝날 때까지 대기

	case <-boom:
		log.Println("BOOM!")
		return
	}

}
