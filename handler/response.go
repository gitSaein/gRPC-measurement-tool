package handler

import (
	"time"

	"google.golang.org/grpc/status"

	m "gRPC_measurement_tool/measure"
)

func HandleReponse(err error, wid uint64, report *m.Report, option m.Option, p m.Process, startAt time.Time) {
	if err != nil {
		handleError(err, wid, report, option, p, startAt)
	} else {
		handleOk(err, wid, report, option, p, startAt)
	}
}

func handleError(err error, wid uint64, report *m.Report, option m.Option, p m.Process, startAt time.Time) *m.Report {

	errorStatus := &m.ErrorStatus{}

	st, _ := status.FromError(err)

	errorStatus.Wid = wid
	errorStatus.Code = st.Proto().Code
	errorStatus.Message = st.Proto().Message
	errorStatus.Details = st.Proto().Details
	errorStatus.Timestamp = time.Now()
	report.Errors = append(report.Errors, errorStatus)

	response := &m.ResponseState{}
	response.Wid = wid
	response.Status = string(m.ERROR)
	response.Process = p
	response.Duration = time.Since(startAt)
	report.ResponseState = append(report.ResponseState, response)

	return report

}

func handleOk(err error, wid uint64, report *m.Report, option m.Option, p m.Process, startAt time.Time) *m.Report {
	response := &m.ResponseState{}

	response.Wid = wid
	response.Status = string(m.OK)
	response.Process = p
	response.Duration = time.Since(startAt)
	report.ResponseState = append(report.ResponseState, response)

	return report
}
