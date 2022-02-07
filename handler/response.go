package handler

import (
	"time"

	"google.golang.org/grpc/status"

	m "gRPC_measurement_tool/measure"
)

func HandleReponse(err error, worker m.Worker, job *m.Job, option m.Option, p m.ProcessName, startAt time.Time) {
	if err != nil {
		st, _ := status.FromError(err)
		job.Errors = append(job.Errors, &m.ErrorStatus{
			Wid:       worker.WId,
			Jid:       job.JId,
			Code:      st.Proto().Code,
			Message:   st.Proto().Message,
			Details:   st.Proto().Details,
			Timestamp: time.Now(),
		})

		job.Process = append(job.Process, &m.Process{
			Name:     p,
			Status:   string(m.ERROR),
			Duration: time.Since(startAt),
		})
	} else {
		job.Process = append(job.Process, &m.Process{
			Name:     p,
			Status:   string(m.OK),
			Duration: time.Since(startAt),
		})
	}
}
