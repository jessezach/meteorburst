package controllers

import (
	"sort"
	"strconv"
	"time"

	"github.com/astaxie/beego/logs"
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
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	for {
		select {
		case <-quit:
			responses = []int{}
			totalRequests = 0

		case r := <-response:
			responses = append(responses, r)

			if len(responses)%100 == 0 {
				resp := getMean(responses)

				publish <- newEvent(MESSAGE, strconv.Itoa(resp))

				totalRequests += 100
				publish <- newEvent(TOTAL, strconv.Itoa(totalRequests))

				p90 := calcP(responses, 90)
				publish <- newEvent(P90, strconv.Itoa(p90))

				p99 := calcP(responses, 99)
				publish <- newEvent(P99, strconv.Itoa(p99))

				p50 := calcP(responses, 50)
				publish <- newEvent(P50, strconv.Itoa(p50))

				time.Sleep(time.Second * 5)
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

func calcP(slice []int, metricType int) int {
	sort.Ints(slice)
	var n int
	length := len(slice)

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
