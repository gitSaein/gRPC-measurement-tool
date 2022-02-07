package main

import (
	"log"
	"sync"
	"time"

	c "gRPC_measurement_tool/cmd"
	u "gRPC_measurement_tool/config"
	h "gRPC_measurement_tool/handler"
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

var kkk = 0

func job(wait *sync.WaitGroup, cmd m.Option, worker *m.Worker, job *m.Job) {
	startAt := time.Now()
	setting := u.SettingOptions(cmd)
	h.HandleReponse(setting.Error, worker, job, cmd, m.SetOption, startAt)

	defer func() {
		kkk += 1
		worker.Jobs = append(worker.Jobs, job)
		wait.Done()
		setting.CancelFunc()
		job.Duration = time.Since(startAt)
		job.TimeStamp = time.Now()

		// log.Printf("wid: %v jid: %v  %s\n", worker.WId, job.JId, job.TimeStamp)
	}()

	// Set up a connection to the server.
	startAt1 := time.Now()
	conn, err := grpc.DialContext(setting.Context, cmd.Target, setting.Options...)
	h.HandleReponse(err, worker, job, cmd, m.DialOpen, startAt1)

	if conn != nil {
		defer func() {
			startAt2 := time.Now()
			err = conn.Close()
			h.HandleReponse(err, worker, job, cmd, m.DialClose, startAt2)
		}()
	}

}
func NormalWorker(wait *sync.WaitGroup, report *m.Report, cmd m.Option) {
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
		// go job(wg, option, worker)
	}

	wg.Wait() //Go루틴 모두 끝날 때까지 대기

}

func WorkerWithTickerJob(report *m.Report, cmd m.Option, worker *m.Worker) {
	startAt := time.Now()

	// tick := time.Tick(1 * time.Second)
	// end := time.NewTimer(time.Duration(cmd.LoadMaxDuration) * time.Second)

	wg := new(sync.WaitGroup)
	j := &m.Job{JId: u.GetID()}
	wg.Add(worker.RPS)
	ccc := 0
	go func() {
		for i := 0; i < worker.RPS; i++ {
			job(wg, option, worker, j)
			ccc += 1
		}
	}()

	wg.Wait() //Go루틴 모두 끝날 때까지 대기

	defer func() {
		worker.Duration = time.Since(startAt)
		log.Printf("[%d] - [%d] : %d", worker.WId, worker.RPS, ccc)
	}()

}

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU()) // CPU 개수를 구한 뒤 사용할 최대 CPU 개수 설정
	// fmt.Println(runtime.GOMAXPROCS(0))   // 설정 값 출력
	log.Println("Measure start...")

	report := &m.Report{}
	var mutex = &sync.Mutex{}
	defer func() {

		report.Total = time.Since(startInitAt)
		log.Printf("Measure end... %d", kkk)
		kkk = 0
		m.PrintResult(report, option)

	}()

	totalR := option.RT

	for {
		time.Sleep(time.Duration(1) * time.Second)
		log.Println("tick")
		wg := new(sync.WaitGroup)
		wg.Add(option.WorkerCnt)

		go func() {
			log.Printf("left RT: %d", totalR)

			for i := 0; i < option.WorkerCnt; i++ {
				worker := &m.Worker{}
				worker.WId = uint64(i)
				h.ShareRpsPerWorer(option, i, worker, report)
				defer func() {
					report.Workers = append(report.Workers, worker)
					mutex.Lock()
					totalR -= len(worker.Jobs)
					mutex.Unlock()
					wg.Done()
				}()
				go WorkerWithTickerJob(report, option, worker)
			}
		}()
		wg.Wait() //Go루틴 모두 끝날 때까지 대기
	}

}
