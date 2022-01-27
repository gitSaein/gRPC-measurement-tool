package runner

import (
	"time"

	m "gRPC_measurement_tool/measure"
)

type Worker struct {
	Wid       uint64
	StartTime time.Time
	EndTime   time.Time
	Total     time.Duration
	ErrorList *[]m.ErrorStatus
}

type Job struct {
}
