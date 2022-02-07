package cmd

import (
	"flag"
	m "gRPC_measurement_tool/measure"
)

func Basic() m.Option {

	rt := flag.Int("rt", 900, "request Total count")
	rps := flag.Int("rps", 900, "rps")
	timeout := flag.Int("timeout", 1000, "timeout(ms)")
	loadMaxDuration := flag.Int("lmd", 5, "load max duration(s)")
	isTls := flag.Bool("isTls", false, "tls 인증여부")
	call := flag.String("call", "", "call method")
	target := flag.String("target", "localhost:50051", "target")
	w := flag.Int("w", 3, "total worker")

	flag.Parse() // 명령줄 옵션의 내용을 각 자료형별로 분석

	// if flag.NFlag() == 0 { // 명령줄 옵션의 개수가 0개이면
	// 	flag.Usage() // 명령줄 옵션 기본 사용법 출력
	// }

	cmd := &m.Option{
		RT:              *rt,
		RPS:             *rps,
		WorkerCnt:       *w,
		Timeout:         *timeout,
		LoadMaxDuration: *loadMaxDuration,
		IsTls:           *isTls,
		Call:            *call,
		Target:          *target,
	}

	return *cmd
}
