package cmdflag

import (
	"flag"
	"fmt"
)

type Command struct {
	Tr      int
	Timeout int
	IsTls   bool
	Call    string
}

func Basic() Command {
	tr := flag.Int("tr", 1, "total request")          // 명령줄 옵션을 받은 뒤 문자열로 저장
	timeout := flag.Int("timeout", 0, "timeout (ms)") // 명령줄 옵션을 받은 뒤 정수로 저장
	isTls := flag.Bool("isTls", false, "tls 인증여부")    // 명령줄 옵션을 받은 뒤 실수로 저장
	call := flag.String("call", "", "call method")    // 명령줄 옵션을 받은 뒤 불로 저장

	flag.Parse() // 명령줄 옵션의 내용을 각 자료형별로 분석

	if flag.NFlag() == 0 { // 명령줄 옵션의 개수가 0개이면
		flag.Usage() // 명령줄 옵션 기본 사용법 출력
	}

	cmd := &Command{*tr, *timeout, *isTls, *call}

	fmt.Println("rc:", cmd.Tr)
	fmt.Println("timeout:", cmd.Timeout)
	fmt.Println("isTls:", cmd.IsTls)
	fmt.Println("call:", cmd.Call)

	return *cmd
}
