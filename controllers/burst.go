package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	mapset "github.com/deckarep/golang-set"
)

var quit chan bool
var running bool
var testStartTime int64

var (
	users = 0
	timer = time.NewTimer(time.Second)
)

// Request struct for tcp message
type Request struct {
	MType   int
	URL     string
	Headers []string
	Method  string
	Payload string
	Users   int
}

// MeteorBurst makes a REST call to the provided endpoint
func meteorBurst(url string, method string, payload string, headers []string) {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	badResponses := mapset.NewSet()
	badResponses.Add(500)
	badResponses.Add(501)
	badResponses.Add(502)
	badResponses.Add(504)

	client := &http.Client{}

	req, err := http.NewRequest(method, url, bytes.NewBuffer([]byte(payload)))

	if err == nil {
		if len(headers) > 0 {
			for i := 0; i < len(headers); i++ {
				header := strings.Split(headers[i], ":")
				req.Header.Add(strings.TrimSpace(header[0]), strings.TrimSpace(header[1]))
			}
		}

		startTime := time.Now().UnixNano() / int64(time.Millisecond)
		_, err := client.Do(req)
		endTime := time.Now().UnixNano() / int64(time.Millisecond)

		responseTime := int(endTime - startTime)
		response <- responseTime
		if err != nil {
			log.Debug(err.Error())
		}
	} else {
		log.Debug(err.Error())
		return
	}
}

func stopEverything() {
	close(quit)
	running = false
	users = 0
	setStartTime(0)
	stopClient <- "stop"
	batchSize = 0
}

func timeKeeper(d int, format string) {
	select {
	case <-timer.C:
		if quit != nil && running {
			msg := fmt.Sprintf("Stopped after %d %s", d, format)
			publish <- newEvent(STOPPED, msg)
			time.Sleep(time.Millisecond * 1000)
			stopEverything()
		}
		return
	case <-quit:
		return
	}
}

func runOnSlaves(r *RequestDetails, headerList []string) {
	usersPerSlave := users / slaves

	request := &Request{MType: MSG, URL: r.URL, Headers: headerList,
		Method: r.Method, Payload: r.Payload, Users: usersPerSlave}

	write <- request
}
