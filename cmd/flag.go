package cmd

import (
	"flag"
	m "gRPC_measurement_tool/measure"
)

func Basic() m.Option {
	tr := flag.Int("tr", 1, "total request")
	timeout := flag.Duration("timeout", 5000, "timeout (ms)")
	isTls := flag.Bool("isTls", false, "tls 인증여부")
	call := flag.String("call", "", "call method")
	target := flag.String("target", "localhost:50051", "target")

	flag.Parse() // 명령줄 옵션의 내용을 각 자료형별로 분석

	// if flag.NFlag() == 0 { // 명령줄 옵션의 개수가 0개이면
	// 	flag.Usage() // 명령줄 옵션 기본 사용법 출력
	// }

	cmd := &m.Option{*tr, *timeout, *isTls, *call, *target}

	return *cmd
}
