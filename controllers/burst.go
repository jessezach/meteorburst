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
	Slave   int
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

func runOnSlaves(r *RequestDetails, headerList []string, usrs []int, dur []int, units []string) {
	switch r.RampType {
	case "linear":
		rampUpLinearSlaves(r, headerList)
	case "step":
		rampUpStepSlaves(r, headerList, usrs, dur, units)
	default:
		rampUpRegular(r, headerList)
	}
}

func rampUpStepSlaves(r *RequestDetails, headerList []string, usrs []int, dur []int, units []string) {
	slaveNo := 1

	for i := 0; i < len(usrs); i++ {
		usr := usrs[i]
		d := dur[i]
		unit := units[i]

		request := &Request{MType: MSG, URL: r.URL, Headers: headerList,
			Method: r.Method, Payload: r.Payload, Users: usr, Slave: slaveNo}

		write <- request
		totalUsersGenerated += usr

		if unit == "seconds" {
			time.Sleep(time.Second * time.Duration(d))
		} else {
			time.Sleep(time.Minute * time.Duration(d))
		}

		if slaveNo+1 > slaves {
			slaveNo = 1
		} else {
			slaveNo++
		}
	}
}

func rampUpLinearSlaves(r *RequestDetails, headerList []string) {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	timePerUser := float64(r.RampTime) / float64(r.Users)
	millis := timePerUser * 1000

	if slaves > users {
		usersPerSlave := 1

		for i := 0; i < users; i++ {
			request := &Request{MType: MSG, URL: r.URL, Headers: headerList,
				Method: r.Method, Payload: r.Payload, Users: usersPerSlave, Slave: i + 1}

			write <- request
			totalUsersGenerated++
			time.Sleep(time.Millisecond * time.Duration(millis))
		}
	} else {
		usersPerSlave := users / slaves

		log.Debug("Slaves less than user")

		for i := 1; i <= slaves; i++ {
			if i == slaves {
				if usersPerSlave*slaves < users {
					diff := users - (usersPerSlave * slaves)
					usersPerSlave = usersPerSlave + diff
				}
			}

			for j := 1; j <= usersPerSlave; j++ {
				request := &Request{MType: MSG, URL: r.URL, Headers: headerList,
					Method: r.Method, Payload: r.Payload, Users: 1, Slave: i}

				write <- request
				totalUsersGenerated++
				log.Debug("Going to Sleep")
				time.Sleep(time.Millisecond * time.Duration(millis))
				log.Debug("Awake now..")
			}
			log.Debug("Next slave")
		}
	}
}

func runLocal(r *RequestDetails, headerList []string, usrs []int, dur []int, units []string) {
	switch r.RampType {
	case "linear":
		rampUpLinear(r, headerList)
	case "step":
		rampUpInSteps(r, headerList, usrs, dur, units)
	default:
		rampUpRegular(r, headerList)
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
	millis := timePerUser * 1000

	for i := 0; i < r.Users; i++ {
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
		time.Sleep(time.Millisecond * time.Duration(millis))
	}
	return
}

func rampUpInSteps(r *RequestDetails, headerList []string, usrs []int, dur []int, units []string) {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	for i := 0; i < len(usrs); i++ {
		userCount := usrs[i]
		duration := dur[i]
		unit := units[i]

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

		log.Debug("Going to sleep")
		if unit == "seconds" {
			time.Sleep(time.Second * time.Duration(duration))
		} else {
			time.Sleep(time.Minute * time.Duration(duration))
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
