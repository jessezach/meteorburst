package controllers

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	mapset "github.com/deckarep/golang-set"
)

type executor interface {
	run()
	rampUpRegular()
	rampUpLinear()
	rampUpSteps()
}

// MeteorBurst makes a REST call to the provided endpoint
func execute(url string, method string, payload string, headers []string) {
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
	totalUsersGenerated = 0
}

func updateUsers() {
	for {
		select {
		case <-quit:
			return
		default:
			sendMessage(newEvent(USERS, strconv.Itoa(totalUsersGenerated)))
			time.Sleep(time.Second * 1)
		}
	}
}

func timeKeeper(d int, format string) {
	select {
	case <-timer.C:
		if quit != nil && running {
			msg := fmt.Sprintf("Stopped after %d %s", d, format)
			sendMessage(newEvent(STOPPED, msg))
			time.Sleep(time.Millisecond * 1000)
			stopEverything()
		}
		return
	case <-quit:
		return
	}
}
