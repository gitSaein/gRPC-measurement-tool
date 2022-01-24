package error

import (
	"time"

	"google.golang.org/grpc/status"

	m "gRPC_measurement_tool/measure"
)

func HandleError(err error, pid uint64, report *m.Report, option m.Option, p m.Process, startAt time.Time) *m.Report {

	if err != nil {
		errorStatus := &m.ErrorStatus{}

		st, _ := status.FromError(err)

		errorStatus.Code = st.Proto().Code
		errorStatus.Message = st.Proto().Message
		errorStatus.Details = st.Proto().Details
		errorStatus.Timestamp = time.Now()
		report.Errors = append(report.Errors, errorStatus)

		response := &m.ResponseState{}
		response.Status = string(m.ERROR)
		response.Process = p
		response.Duration = time.Since(startAt)
		report.ResponseState = append(report.ResponseState, response)

		return report

	} else {
		response := &m.ResponseState{}
		response.Status = string(m.OK)
		response.Process = p
		response.Duration = time.Since(startAt)
		report.ResponseState = append(report.ResponseState, response)

	}

	return report

}
