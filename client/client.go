package main

import (
	"log"
	"time"

	c "gRPC_measurement_tool/cmd"
	h "gRPC_measurement_tool/handler"
	m "gRPC_measurement_tool/measure"
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
	option = c.Basic()
}

func main() {
	// runtime.GOMAXPROCS(runtime.NumCPU()) // CPU 개수를 구한 뒤 사용할 최대 CPU 개수 설정
	// fmt.Println(runtime.GOMAXPROCS(0))   // 설정 값 출력
	log.Println("Measure start...")
	ch_worker := make(chan m.Worker)
	ch_result := make(chan []*m.Worker)
	ch_left_rps := make(chan int)
	ch_done := make(chan bool)
	workers := []*m.Worker{}

	// 1. total 값 + rps 로 주고,
	// 2. 1초마다 호출 rps 만큼 호출 하고, total에서 실제 실행된 수 만큼 뺴준다.
	// 3. total값 만큼 다 호출하면 끝낸다.
	// totalR := option.RT
	startInitAt = time.Now()
	workCnt := h.MakeWorkCount(option.RT, option.RPS)
	go h.Work(option.RPS, option.RT, ch_worker, workCnt)
	go h.Jobs(ch_worker, ch_result, workers, ch_left_rps, option)
	go h.CheckLeftRequest(ch_left_rps, ch_worker, workCnt, ch_done, option.RPS)
	h.Result(ch_worker, ch_result, workCnt, ch_done, startInitAt, option)

}
