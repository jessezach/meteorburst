package controllers

import (
	"sort"
	"strconv"
	"time"
)

var (
	response      = make(chan int)
	responses     = []int{}
	totalRequests = 0
	batchSize     = 0
)

func init() {
	go poolStats()
}

func poolStats() {
	for {
		select {
		case <-quit:
			responses = []int{}
			totalRequests = 0

		case r := <-response:
			responses = append(responses, r)
			length := len(responses)

			if batchSize > 0 && length%batchSize == 0 {
				resp := getMean(responses, length)
				publish <- newEvent(MESSAGE, strconv.Itoa(resp))

				diff := (time.Now().UnixNano() / int64(time.Millisecond)) - testStartTime
				rps := rps(diff, length)
				publish <- newEvent(RPS, strconv.Itoa((int(rps))))

				totalRequests += batchSize
				publish <- newEvent(TOTAL, strconv.Itoa(totalRequests))

				sort.Ints(responses)
				p90 := calcP(responses, 90, length)
				publish <- newEvent(P90, strconv.Itoa(p90))

				p99 := calcP(responses, 99, length)
				publish <- newEvent(P99, strconv.Itoa(p99))

				p50 := calcP(responses, 50, length)
				publish <- newEvent(P50, strconv.Itoa(p50))
			}
		}
	}
}

func getMean(slice []int, length int) int {
	sum := 0
	for _, v := range slice {
		sum += v
	}
	return sum / length
}

func calcP(slice []int, metricType int, length int) int {
	var percent int

	switch metricType {

	case 90:
		percent = 10

	case 99:
		percent = 1

	case 50:
		percent = 50
	}

	topEnd := (percent * length) / 100
	idx := (length - topEnd) - 1
	n := slice[idx]
	return n
}

func rps(diff int64, length int) float64 {
	return float64(1000) / float64(diff) * float64(length)
}
