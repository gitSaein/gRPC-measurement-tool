package measure

import "time"

func histogram(report *Report) {

}

func minMaxAverage(report *Report) {

	var min time.Duration
	var max time.Duration
	var avg time.Duration
	var total time.Duration
	for idx, worker := range report.Workers {
		for _, job := range worker.Jobs {

			if min > job.Duration || idx == 0 {
				min = job.Duration
			}

			if max < job.Duration || idx == 0 {
				max = job.Duration
			}
			total += job.Duration
		}
		avg = time.Duration(int64(total) / int64(len(worker.Jobs)*len(report.Workers)))
	}

	report.Min = min
	report.Max = max
	report.Avg = avg
}

func indexOf(code int32, data []*ErrorStatus) int {
	for k, v := range data {
		if code == v.Code {
			return k
		}
	}
	return -1
}

func CheckResultCnt(report *Report, worker *Worker, job *Job) {
	for _, err := range job.Errors {
		report.ErrorResult.Count += 1
		if indexOf(err.Code, report.ErrorResult.Errors) == -1 {
			report.ErrorResult.Errors = append(report.ErrorResult.Errors, err)
		}
	}
}
