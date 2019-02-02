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

var quit chan bool
var running bool
var testStartTime int64

var (
	users               = 0
	timer               = time.NewTimer(time.Second)
	totalUsersGenerated = 0
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
	totalUsersGenerated = 0
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

func runLocal(r *RequestDetails, headerList []string) {
	switch r.RampType {
	case "linear":
		go rampUpLinear(r, headerList)
	case "step":
		go rampUpInSteps(r, headerList)
	}
}

func rampUpRegular(r *RequestDetails, headerList []string) {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	for i := 0; i < r.Users; i++ {
		log.Debug("Starting user %#v", i+1)
		go func() {
			for {
				select {
				case <-quit:
					log.Debug("Returning from go routine")
					return
				default:
					meteorBurst(r.URL, r.Method, r.Payload, headerList)
				}
			}
		}()
		totalUsersGenerated++
	}
	return
}

func rampUpLinear(r *RequestDetails, headerList []string) {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	timePerUser := float64(r.RampTime) / float64(r.Users)

	for i := 0; i < r.Users; i++ {
		log.Debug("Starting user %#v", i+1)
		go func() {
			for {
				select {
				case <-quit:
					log.Debug("Returning from go routine")
					return
				default:
					meteorBurst(r.URL, r.Method, r.Payload, headerList)
				}
			}
		}()
		totalUsersGenerated++
		time.Sleep(time.Second * time.Duration(timePerUser))
	}
	return
}

func rampUpInSteps(r *RequestDetails, headerList []string) {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	steps := strings.Split(r.RampStep, "\n")
	timeUnit := strings.ToLower(strings.Split(steps[0], ":")[1])
	unit := strings.TrimSpace(timeUnit)
	u := strings.TrimRight(unit, "\n")

	units := []string{"mins", "minutes", "minute", "seconds", "second", "sec"}

	if !contains(units, u) {
		publish <- newEvent(ERROR, "Bad Ramp up format")
		stopEverything()
		return
	}

	for _, step := range steps[1:] {
		stepList := strings.Split(step, ":")
		log.Debug(stepList[0])
		log.Debug(stepList[1])

		userCount, _ := strconv.Atoi(strings.TrimSpace(stepList[0]))
		d := strings.TrimSpace(stepList[1])
		dur, _ := strconv.Atoi(strings.TrimRight(d, "\n"))

		for i := 0; i < userCount; i++ {
			go func() {
				for {
					select {
					case <-quit:
						log.Debug("Returning from go routine")
						return
					default:
						meteorBurst(r.URL, r.Method, r.Payload, headerList)
					}
				}
			}()
			totalUsersGenerated++
		}
		log.Debug(strconv.Itoa(totalUsersGenerated))
		log.Debug("Going to sleep")
		log.Debug(strconv.Itoa(dur))

		if u == "seconds" || u == "second" || u == "sec" {
			log.Debug("Sleeping in seconds")
			time.Sleep(time.Second * time.Duration(dur))
		} else {
			log.Debug("Sleeping in minutes")
			time.Sleep(time.Minute * time.Duration(dur))
		}
		log.Debug("Awake now...")
	}
	return
}

func updateUsers() {
	for {
		select {
		case <-quit:
			return
		default:
			publish <- newEvent(USERS, strconv.Itoa(totalUsersGenerated))
			time.Sleep(time.Second * 1)
		}
	}
}

func contains(arr []string, unit string) bool {
	for _, value := range arr {
		if value == strings.TrimRight(unit, "\n") {
			return true
		}
	}
	return false
}
