package measure

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

func makeHistogramData(report *Report) {
	for _, worker := range report.Workers {
		for _, job := range worker.Jobs {
			report.Histogram =
				append(report.Histogram, &HistogramData{worker.WId, job.JId, job.TimeStamp})
		}
	}
}

func minMaxAverage(report *Report) {

	var min time.Duration
	var max time.Duration
	var avg time.Duration
	var total time.Duration
	if len(report.Workers) > 0 {
		if len(report.Workers[0].Jobs) > 0 {
			min = report.Workers[0].Jobs[0].Duration
			max = report.Workers[0].Jobs[0].Duration
		}
	}

	for _, worker := range report.Workers {

		for _, job := range worker.Jobs {

			if min > job.Duration {
				min = job.Duration
			}

			if max < job.Duration {
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

func CheckResultCnt(report *Report) {
	if len(report.Workers) > 0 {
		for i, worker := range report.Workers {
			report.JobResult.TotalCnt += len(worker.Jobs)
			report.Workers[i].RPSSet = float64(len(worker.Jobs)) / worker.Duration.Seconds()
			for _, job := range worker.Jobs {
				if len(job.Errors) == 0 {
					report.JobResult.OkCnt += 1
				} else {
					for _, err := range job.Errors {
						report.JobResult.ErrCnt += 1
						if indexOf(err.Code, report.JobResult.Errors) == -1 {
							report.JobResult.Errors = append(report.JobResult.Errors, err)
						}
					}
				}

			}
		}
	}

}

// from ghz
type Bucket struct {
	// The Mark for histogram bucket in seconds
	Mark float64 `json:"mark"`

	// The count in the bucket
	Count int `json:"count"`

	// The frequency of results in the bucket as a decimal percentage
	Frequency float64 `json:"frequency"`
}

func histogram(latencies []float64, slowest, fastest float64) []Bucket {
	bc := 10
	buckets := make([]float64, bc+1)
	counts := make([]int, bc+1)
	bs := (slowest - fastest) / float64(bc)
	for i := 0; i < bc; i++ {
		buckets[i] = fastest + bs*float64(i)
	}
	buckets[bc] = slowest
	var bi int
	var max int
	for i := 0; i < len(latencies); {
		if latencies[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			bi++
		}
	}
	res := make([]Bucket, len(buckets))
	for i := 0; i < len(buckets); i++ {
		res[i] = Bucket{
			Mark:      buckets[i],
			Count:     counts[i],
			Frequency: float64(counts[i]) / float64(len(latencies)),
		}
	}
	return res
}

const (
	barChar = "■"
)

func histogramPrintString(buckets []Bucket) string {
	maxMark := 0.0
	maxCount := 0
	for _, b := range buckets {
		if v := b.Mark; v > maxMark {
			maxMark = v
		}
		if v := b.Count; v > maxCount {
			maxCount = v
		}
	}

	formatMark := func(mark float64) string {
		return fmt.Sprintf("%.3f", mark*1000)
	}
	formatCount := func(count int) string {
		return fmt.Sprintf("%v", count)
	}

	maxMarkLen := len(formatMark(maxMark))
	maxCountLen := len(formatCount(maxCount))
	res := new(bytes.Buffer)
	for i := 0; i < len(buckets); i++ {
		// Normalize bar lengths.
		var barLen int
		if maxCount > 0 {
			barLen = (buckets[i].Count*40 + maxCount/2) / maxCount
		}
		markStr := formatMark(buckets[i].Mark)
		countStr := formatCount(buckets[i].Count)
		res.WriteString(fmt.Sprintf(
			"  %v%s [%v]%s |%v\n",
			markStr,
			strings.Repeat(" ", maxMarkLen-len(markStr)),
			countStr,
			strings.Repeat(" ", maxCountLen-len(countStr)),
			strings.Repeat(barChar, barLen),
		))
	}

	return res.String()
}
