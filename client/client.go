package main

import (
	"fmt"
	"log"
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

func job(cmd m.Option, worker m.Worker) m.Worker {
	startAt := time.Now()
	job := &m.Job{JId: u.GetID()}
	setting := u.SettingOptions(cmd)
	h.HandleReponse(setting.Error, worker, job, cmd, m.SetOption, startAt)
	// Set up a connection to the server.
	startAt1 := time.Now()
	conn, err := grpc.DialContext(setting.Context, cmd.Target, setting.Options...)
	h.HandleReponse(err, worker, job, cmd, m.DialOpen, startAt1)

	kkk += 1
	if conn != nil {
		startAt2 := time.Now()
		err = conn.Close()
		h.HandleReponse(err, worker, job, cmd, m.DialClose, startAt2)
	}
	setting.CancelFunc()
	job.Duration = time.Since(startAt)
	job.TimeStamp = time.Now()
	worker.Jobs = append(worker.Jobs, job)
	// log.Printf("wid: %v jid: %v  %s\n", worker.WId, job.JId, job.TimeStamp)
	return worker

}

func WorkerWithTickerJob(cmd m.Option, worker m.Worker) m.Worker {
	startAt := time.Now()
	worker = job(option, worker)
	worker.Duration = time.Since(startAt)
	return worker
}

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU()) // CPU 개수를 구한 뒤 사용할 최대 CPU 개수 설정
	// fmt.Println(runtime.GOMAXPROCS(0))   // 설정 값 출력
	log.Println("Measure start...")

	report := &m.Report{}
	defer func() {

		report.Total = time.Since(startInitAt)
		log.Printf("Measure end... %d", kkk)
		kkk = 0
		m.PrintResult(report, option)

	}()

	ch := make(chan m.Worker)
	ch_result := make(chan m.Worker)
	req_cnt := make(chan int)
	ch_done := make(chan bool)

	// 1. total 값 + rps 로 주고,
	// 2. 1초마다 호출 rps 만큼 호출 하고, total에서 실제 실행된 수 만큼 뺴준다.
	// 3. total값 만큼 다 호출하면 끝낸다.
	// totalR := option.RT

	left := option.RT % option.RPS
	loopCnt := int(h.ToFixed(float64(option.RT/option.RPS), 0))
	if left > 0 && left < option.RPS {
		loopCnt = loopCnt + 1
	}

	for i := 0; i < loopCnt; i++ {
		go func(j int) {
			work(ch, j)
		}(i)
	}

	report = jobz(ch, ch_result, req_cnt, option.RT, ch_done)

	a := <-req_cnt
	fmt.Println(a)

}

func jobz(ch chan m.Worker, ch_result chan m.Worker, req_cnt chan int, totalCnt int, ch_done chan bool) *m.Report {

	for {
		select {
		case worker := <-ch:
			fmt.Println(worker)
			time.Sleep(time.Duration(1) * time.Second)
			log.Println("tick")
			jobs := []m.Job{}
			for i := 0; i < worker.RPS; i++ {
				go func() {
					ch_result <- WorkerWithTickerJob(option, worker)
				}()
				work_result := <-ch_result
				jobs = append(jobs, work_result)
				log.Print(work_result)
				log.Printf("%v - %v", len(report.Workers), len(work_result.Jobs))
			}
		case <-ch_done:
			return report
		}
	}

}

func work(ch chan m.Worker, wno int) {

	worker := m.Worker{}
	worker.WId = u.GetID()
	worker = h.ShareRpsPerWorer(option, wno, worker)
	ch <- worker

}
