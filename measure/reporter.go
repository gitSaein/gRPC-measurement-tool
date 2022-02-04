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
	Timeout         int
	LoadMaxDuration int
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
	RPS      int
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

type JobResult struct {
	ErrCnt   int
	OkCnt    int
	TotalCnt int
	Errors   []*ErrorStatus
}

type HistogramData struct {
	WId       uint64
	JId       uint64
	Timestamp time.Time
}

type Report struct {
	StartTime  time.Time
	EndTime    time.Time
	Total      time.Duration
	Min        time.Duration
	Max        time.Duration
	Avg        time.Duration
	JobResult  JobResult
	SuccessCnt int
	Workers    []*Worker
	States     []*ConnectState
	Histogram  []*HistogramData
}

func PrintResult(report *Report, cmd Option) {

	minMaxAverage(report)
	CheckResultCnt(report)
	fmt.Println("Summary:")
	fmt.Println(" Options:")
	fmt.Printf("   Woker: %v\n   Rps: %v\n   Dial Connection timeout: %v\n   Load-max-duration: %v\n   Tls: %v\n   Target: %v\n",
		cmd.WorkerCnt, cmd.RPS, time.Duration(cmd.Timeout)*time.Millisecond, time.Duration(cmd.LoadMaxDuration)*time.Second, cmd.IsTls, cmd.Target)
	fmt.Println(" Request latency:")
	fmt.Printf("   Total(s): %v\n", report.Total)
	fmt.Printf("   Min: %v\n", report.Min)
	fmt.Printf("   Max: %v\n", report.Max)
	fmt.Printf("   Avg: %v\n", report.Avg)
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

	fmt.Println(" Response Result:")
	fmt.Printf("  Total: [%-5v] OK:[%-5v] Failed:[%-5v]\n", report.JobResult.TotalCnt, report.JobResult.OkCnt, report.JobResult.ErrCnt)

	if report.JobResult.ErrCnt > 0 {
		fmt.Println(" Errors:")
		fmt.Println("  Code       message:")
		for _, state := range report.JobResult.Errors {
			fmt.Printf("  [%-5v]    %-5v\n", state.Code, state.Message)
		}
	}
	makeHistogramData(report)

	okLats := make([]float64, 0)
	for _, d := range report.Histogram {
		okLats = append(okLats, float64(d.Timestamp.Second()))
	}

	sort.Float64s(okLats)
	fmt.Println()

	if len(okLats) > 0 {
		var fastestNum, slowestNum float64
		fastestNum = okLats[0]
		slowestNum = okLats[len(okLats)-1]

		histogramRet := histogram(okLats, slowestNum, fastestNum)
		fmt.Printf(" Response time histogram(RPS):\n")
		fmt.Printf("%s", histogramPrintString(histogramRet))
	}

}
