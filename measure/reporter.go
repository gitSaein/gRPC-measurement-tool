package measure

import (
	"fmt"
	"sort"
	"time"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc/connectivity"
)

type Option struct {
	RT              int
	WorkerCnt       int
	Timeout         time.Duration
	LoadMaxDuration time.Duration
	IsTls           bool
	Call            string
	Target          string
	RPS             int
}

type ErrorStatus struct {
	Wid uint64
	Jid uint64
	// The status code, which should be an enum value of [google.rpc.Code][google.rpc.Code].
	Code int32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// A developer-facing error message, which should be in English. Any
	// user-facing error message should be localized and sent in the
	// [google.rpc.Status.details][google.rpc.Status.details] field, or localized by the client.
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	// A list of messages that carry the error details.  There is a common set of
	// message types for APIs to use.
	Details   []*any.Any `protobuf:"bytes,3,rep,name=details,proto3" json:"details,omitempty"`
	Timestamp time.Time
}

type ConnectState struct {
	ConnectState connectivity.State
	Duration     time.Duration
	TimeStamp    time.Time
}

type Status string

const (
	OK    = Status("OK")
	ERROR = Status("ERROR")
)

type ProcessName string

const (
	SetOption  = ProcessName("SetOption")
	CallMethod = ProcessName("CallMethod")
	DialOpen   = ProcessName("DialOpen")
	DialClose  = ProcessName("DialClose")
)

type Worker struct {
	WId      uint64
	Jobs     []*Job
	Duration time.Duration
}
type Job struct {
	JId       uint64
	States    *ConnectState
	Errors    []*ErrorStatus
	Process   []*Process
	Duration  time.Duration
	TimeStamp time.Time
}

type Process struct {
	Name     ProcessName
	Status   string
	Duration time.Duration
}

type ErrorResult struct {
	Count  int
	Errors []*ErrorStatus
}

type HistogramData struct {
	WId       uint64
	JId       uint64
	Timestamp time.Time
}

type Report struct {
	StartTime   time.Time
	EndTime     time.Time
	Total       time.Duration
	Min         time.Duration
	Max         time.Duration
	Avg         time.Duration
	ErrorResult ErrorResult
	SuccessCnt  int
	Workers     []*Worker
	States      []*ConnectState
	Histogram   []*HistogramData
}

func PrintResult(report *Report, cmd Option) {

	minMaxAverage(report)
	CheckResultCnt(report)
	fmt.Println("Summary:")
	fmt.Printf("  Target: %v\n", cmd.Target)
	fmt.Printf("  RequestCount: %v\n", cmd.RT)
	fmt.Printf("  Total: %v\n", report.Total)
	fmt.Println(" Request latency:")
	fmt.Printf("   Min: %v\n", report.Min)
	fmt.Printf("   Max: %v\n", report.Max)
	fmt.Printf("   Avg: %v\n", report.Avg)

	// fmt.Println("  Options:")
	// fmt.Printf("     rt: %v\n     w: %v\n     timeout: %v\n     load-max-duration: %v\n     isTls: %v\n     call: %v\n     target: %v\n     rps: %v\n   ",
	// 	cmd.RT, cmd.WorkerCnt, cmd.Timeout, cmd.LoadMaxDuration, cmd.IsTls, cmd.Call, cmd.Target, cmd.RPS)
	fmt.Println()

	// if len(report.Workers) > 0 {
	// 	fmt.Println("Process Tracking:")
	// 	fmt.Println("  Worker    Job    State   Process            Duration")
	// 	for _, worker := range report.Workers {
	// 		fmt.Printf("  [%-5v] %-5v\n", worker.WId, worker.Duration)
	// 		for _, job := range worker.Jobs {
	// 			fmt.Printf("  [%-5v] [%-5v]  %-5v\n", worker.WId, job.JId, job.Duration)
	// 			for _, process := range job.Process {
	// 				fmt.Printf("  [%-5v] [%-5v] [%-5v] [%-15v]  %-5v\n", worker.WId, job.JId, process.Status, process.Name, process.Duration)
	// 			}
	// 		}
	// 	}

	// }
	// fmt.Println()

	// if len(report.States) > 0 {
	// 	fmt.Println("Dial State Trace:")
	// 	fmt.Println("  State       duration:")
	// 	for _, state := range report.States {
	// 		fmt.Printf("  [%v]       %v\n", state.ConnectState, state.Duration)
	// 	}
	// }
	// fmt.Println()

	if report.ErrorResult.Count > 0 {
		fmt.Println("Errors:")
		fmt.Println("  Code       message:")
		for _, state := range report.ErrorResult.Errors {
			fmt.Printf("  [%-5v]    %-5v\n", state.Code, state.Message)
		}
	}
	fmt.Println()

	makeHistogramData(report)

	// for _, h := range report.Histogram {

	// }

	okLats := make([]float64, 0)
	for _, d := range report.Histogram {
		okLats = append(okLats, float64(d.Timestamp.Second()))
	}

	sort.Float64s(okLats)

	if len(okLats) > 0 {
		var fastestNum, slowestNum float64
		fastestNum = okLats[0]
		slowestNum = okLats[len(okLats)-1]

		histogramRet := histogram(okLats, slowestNum, fastestNum)
		fmt.Printf("Response time histogram:\n")
		fmt.Printf("%s", histogramPrintString(histogramRet))
	}

}
