package main

import (
	"log"
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
		job.Duration = time.Since(startAt)
	}()

	// Set up a connection to the server.
	startAt1 := time.Now()
	conn, err := grpc.DialContext(setting.Context, cmd.Target, setting.Options...)
	errorModel.HandleReponse(err, worker, job, cmd, m.DialOpen, startAt1)

	if conn != nil {
		defer func() {
			startAt2 := time.Now()
			err = conn.Close()
			errorModel.HandleReponse(err, worker, job, cmd, m.DialClose, startAt2)
		}()
	}

}

func TickWorkers(report *m.Report, cmd m.Option) {

	tick := time.Tick(1 * time.Second)
	boom := time.After(option.LoadMaxDuration)

	requestCnt := option.RT
	for requestCnt > 0 {
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
}

func NormalWorkers(report *m.Report, cmd m.Option) {

	wg := new(sync.WaitGroup)
	wg.Add(option.WorkerCnt)

	for i := 0; i < option.WorkerCnt; i++ {
		go Worker(wg, report, option)
	}

	wg.Wait() //Go루틴 모두 끝날 때까지 대기

}

func Worker(wait *sync.WaitGroup, report *m.Report, cmd m.Option) {
	startAt := time.Now()

	wg := new(sync.WaitGroup)
	wg.Add(option.RT)

	worker := &m.Worker{}
	worker.WId = u.GetID()
	defer func() {
		worker.Duration = time.Since(startAt)
		report.Workers = append(report.Workers, worker)
		wait.Done()
	}()
	for i := 0; i < option.RT; i++ {
		go job(wg, option, worker)
	}

	wg.Wait() //Go루틴 모두 끝날 때까지 대기

}

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU()) // CPU 개수를 구한 뒤 사용할 최대 CPU 개수 설정
	// fmt.Println(runtime.GOMAXPROCS(0))   // 설정 값 출력
	log.Println("Measure start...")

	report := &m.Report{}
	defer func() {
		report.Total = time.Since(startInitAt)
		log.Println("Measure end...")
		m.PrintResult(report, option)

	}()

	if option.LoadMaxDuration > 0 {
		TickWorkers(report, option)
	} else {
		NormalWorkers(report, option)
	}
}
