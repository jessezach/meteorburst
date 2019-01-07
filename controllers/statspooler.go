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

			if length%100 == 0 {
				resp := getMean(responses)
				publish <- newEvent(MESSAGE, strconv.Itoa(resp))

				diff := (time.Now().UnixNano() / int64(time.Millisecond)) - testStartTime
				rps := rps(diff, length)
				publish <- newEvent(RPS, strconv.Itoa((int(rps))))

				totalRequests += 100
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

func getMean(slice []int) int {
	sum := 0
	for _, v := range slice {
		sum += v
	}
	return sum / len(slice)
}

func calcP(slice []int, metricType int, length int) int {
	var n int

	switch metricType {

	case 90:
		ten := (10 * length) / 100
		idx := (length - ten) - 1
		n = slice[idx]

	case 99:
		one := (1 * length) / 100
		idx := (length - one) - 1
		n = slice[idx]

	case 50:
		fifty := (50 * length) / 100
		idx := (length - fifty) - 1
		n = slice[idx]
	}
	return n
}

func rps(diff int64, length int) float64 {
	return float64(1000) / float64(diff) * float64(length)
}
