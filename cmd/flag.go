package cmd

import (
	"flag"
	m "gRPC_measurement_tool/measure"
)

func Basic() m.Option {
	tr := flag.Int("tr", 1, "total request")
	timeout := flag.Duration("timeout", 5000, "timeout (ms)")
	isTls := flag.Bool("isTls", true, "tls 인증여부")
	call := flag.String("call", "", "call method")
	ip := flag.String("ip", "localhost", "ip")
	port := flag.String("port", "50051", "port")

	flag.Parse() // 명령줄 옵션의 내용을 각 자료형별로 분석

	// if flag.NFlag() == 0 { // 명령줄 옵션의 개수가 0개이면
	// 	flag.Usage() // 명령줄 옵션 기본 사용법 출력
	// }

	cmd := &m.Option{*tr, *timeout, *isTls, *call, *ip, *port}

	return *cmd
}
