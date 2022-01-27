package measure

import (
	"fmt"
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

type Process string

const (
	SetOption  = Process("SetOption")
	CallMethod = Process("CallMethod")
	DialOpen   = Process("DialOpen")
	DialClose  = Process("DialClose")
)

type ResponseState struct {
	Wid      uint64
	Status   string
	Process  Process
	Duration time.Duration
}

type Report struct {
	Wid           uint64
	StartTime     time.Time
	EndTime       time.Time
	Total         time.Duration
	States        []*ConnectState
	ResponseState []*ResponseState
	Errors        []*ErrorStatus
}

func PrintResult(report *Report, cmd Option) {
	fmt.Println()
	fmt.Println("Summary:")
	fmt.Printf("  Target: %v\n", cmd.Target)
	fmt.Printf("  Total: %v\n", report.Total)
	fmt.Println("  Options:")
	fmt.Printf("     %v\n", cmd)

	fmt.Println()

	if len(report.ResponseState) > 0 {
		fmt.Println("Process Tracking:")
		fmt.Println("  Pid      State   Process            Duration")
		for _, state := range report.ResponseState {
			fmt.Printf("  [%-5v]  [%-5v] [%-15v]  %-5v\n", state.Wid, state.Status, state.Process, state.Duration)
		}
	}
	fmt.Println()

	if len(report.States) > 0 {
		fmt.Println("Dial State Trace:")
		fmt.Println("  State       duration:")
		for _, state := range report.States {
			fmt.Printf("  [%v]       %v\n", state.ConnectState, state.Duration)
		}
	}
	fmt.Println()

	if len(report.Errors) > 0 {
		fmt.Println("Error Description:")
		fmt.Println("  Pid      code         message:")
		for _, state := range report.Errors {
			fmt.Printf("  %-5v   [%-5v]       %-5v\n", state.Wid, state.Code, state.Message)
		}
	}
	fmt.Println()

}
