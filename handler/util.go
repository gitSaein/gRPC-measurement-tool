package handler

import (
	m "gRPC_measurement_tool/measure"
	"math"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func ShareRpsPerWorer(option m.Option, wno int, w m.Worker) m.Worker {

	perRps := int(ToFixed(float64(option.RPS/option.WorkerCnt), 0))
	alreadyGet := wno * perRps
	leftTotalRps := option.RPS - alreadyGet

	leftoverByDivide := leftTotalRps - perRps
	if leftoverByDivide > 0 && leftoverByDivide < perRps {
		w.RPS = perRps + leftoverByDivide
	} else {
		w.RPS = perRps
	}
	return w

}
