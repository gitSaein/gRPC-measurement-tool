package measure

import (
	"context"
	"fmt"
	"os"
	"sort"
	"time"

	h "github.com/aybabtme/uniplot/histogram"
	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Option struct {
	RT int
	// WorkerCnt       int
	Timeout int
	// LoadMaxDuration int
	IsTls  bool
	Call   string
	Target string
	RPS    int
}

type Setting struct {
	Options    []grpc.DialOption
	Error      error
	Context    context.Context
	CancelFunc context.CancelFunc
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
	RPSSet   float64
	Duration time.Duration
	Setting  *Setting
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
	WId      uint64
	JId      uint64
	Duration time.Duration
}

type Report struct {
	StartTime  time.Time
	EndTime    time.Time
	Total      time.Duration
	Min        time.Duration
	Max        time.Duration
	RPS        float64
	Avg        time.Duration
	JobResult  JobResult
	SuccessCnt int
	Workers    []*Worker
	States     []*ConnectState
	Histogram  []*HistogramData
}

func PrintResult(report *Report, option Option) {
	fmt.Println()
	CheckResultCnt(report)
	finalCalc(report)
	fmt.Println("Summary:")
	fmt.Println(" Options:")
	fmt.Printf("   Request Total: %v\n   Rps: %v\n   Dial Connection timeout: %v\n   Tls: %v\n   Target: %v\n",
		option.RT, option.RPS, time.Duration(option.Timeout)*time.Millisecond, option.IsTls, option.Target)
	fmt.Println(" Request latency:")
	fmt.Printf("   Total(s): %v\n", report.Total)
	fmt.Printf("   Min: %v\n", report.Min)
	fmt.Printf("   Max: %v\n", report.Max)
	fmt.Printf("   Avg: %v\n", report.Avg)
	fmt.Printf("   Request(s): %v\n", report.RPS)

	fmt.Printf(" Total: [%-5v] OK:[%-5v] Failed:[%-5v]\n", report.JobResult.TotalCnt, report.JobResult.OkCnt, report.JobResult.ErrCnt)

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
		okLats = append(okLats, float64(d.Duration))
	}

	sort.Float64s(okLats)
	fmt.Println()

	if len(okLats) > 0 {

		hist := h.Hist(10, okLats)
		fmt.Printf(" Response time histogram:\n")

		if err := h.Fprintf(os.Stdout, hist, h.Linear(5), func(v float64) string {
			return time.Duration(v).String()
		}); err != nil {
			panic(err)
		}

	}

}
