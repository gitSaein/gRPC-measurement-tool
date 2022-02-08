package handler

import (
	"fmt"
	c "gRPC_measurement_tool/config"
	m "gRPC_measurement_tool/measure"

	"log"
	"math"
	"sync"
	"time"

	"google.golang.org/grpc"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func ShareRpsPerWorer(rps int, rt int, wno int, w m.Worker) m.Worker {

	perRps := rps
	alreadyGet := wno * perRps
	leftTotalRps := rt - alreadyGet

	if leftTotalRps > 0 && leftTotalRps < perRps {
		w.RPS = leftTotalRps
	} else {
		w.RPS = perRps
	}
	return w

}

func job(option m.Option, wid uint64) *m.Job {
	startAt := time.Now()
	job := &m.Job{JId: c.GetID()}
	setting := c.SettingOptions(option)
	HandleReponse(setting.Error, wid, job, m.SetOption, startAt)
	startAt1 := time.Now()
	conn, err := grpc.DialContext(setting.Context, option.Target, setting.Options...)
	HandleReponse(err, wid, job, m.DialOpen, startAt1)
	if conn != nil {
		startAt2 := time.Now()
		err = conn.Close()
		HandleReponse(err, wid, job, m.DialClose, startAt2)
	}
	setting.CancelFunc()
	job.Duration = time.Since(startAt)
	job.TimeStamp = time.Now()
	return job

}

func MakeWorkCount(total int, rps int) int {
	left := total % rps
	workCnt := int(ToFixed(float64(total/rps), 0))
	if left > 0 && left < rps {
		workCnt = workCnt + 1
	}
	return workCnt
}

func Work(rps int, rt int, ch_worker chan m.Worker, workCnt int) {
	for wno := 0; wno < workCnt; wno++ {
		go func(wno int) {
			worker := m.Worker{}
			worker.WId = c.GetID()
			worker = ShareRpsPerWorer(rps, rt, wno, worker)
			ch_worker <- worker
		}(wno)
	}

}

func Jobs(ch_worker chan m.Worker, ch_result chan []*m.Worker, workers []*m.Worker, ch_left_rps chan int, option m.Option) {
	for {
		select {
		case worker := <-ch_worker:
			time.Sleep(1 * time.Second)
			startAt := time.Now()
			jobs := []*m.Job{}

			var wg sync.WaitGroup
			wg.Add(worker.RPS)
			for i := 0; i < worker.RPS; i++ {
				go func() {
					defer wg.Done()
					j := job(option, worker.WId)
					jobs = append(jobs, j)
				}()
			}
			wg.Wait()
			workers = append(workers, &m.Worker{Jobs: jobs, Duration: time.Since(startAt)})
			left := worker.RPS - len(jobs)
			ch_result <- workers
			ch_left_rps <- left
			log.Printf("  done: %v (left: %v)\n", len(jobs), left)
			fmt.Print("▶ ")

		}

	}
}

func CheckLeftRequest(ch_left_rps chan int, ch_worker chan m.Worker, workCnt int, ch_done chan bool, rps int) {
	cnt := 0
	total_left := 0
	for {
		select {
		case left_rps := <-ch_left_rps:
			cnt += 1
			total_left += left_rps
			if workCnt <= cnt {
				if total_left == 0 {
					ch_done <- true
					close(ch_left_rps)
				}

				workCnt := MakeWorkCount(total_left, rps)
				go Work(rps, total_left, ch_worker, workCnt)
				total_left = 0
			}
		}
	}

}

func Result(ch_worker chan m.Worker, ch_result chan []*m.Worker, workCnt int, ch_done chan bool, startInitAt time.Time, option m.Option) {
	report := &m.Report{}

	for {
		select {
		case workers := <-ch_result:
			report.Workers = workers
		case <-ch_done:
			report.Total = time.Since(startInitAt)
			log.Println("Measure end...")
			m.PrintResult(report, option)
			close(ch_worker)
			close(ch_result)
			close(ch_done)
			return
		}

	}
}
