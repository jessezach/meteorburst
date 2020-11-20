package controllers

import (
	"time"

	"github.com/astaxie/beego/logs"
)

type slaveRunner struct {
	r          *RequestDetails
	headerList []string
	usrs       []int
	dur        []int
	units      []string
}

func (sr slaveRunner) run() {
	switch sr.r.RampType {
	case "linear":
		sr.rampUpLinear()
	case "step":
		sr.rampUpSteps()
	default:
		sr.rampUpRegular()
	}
}

func (sr slaveRunner) rampUpLinear() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	timePerUser := float64(sr.r.RampTime) / float64(sr.r.Users)
	millis := timePerUser * 1000

	if slaves > users {
		usersPerSlave := 1

		for i := 0; i < users; i++ {
			request := &Request{MType: MSG, URL: sr.r.URL, Headers: sr.headerList,
				Method: sr.r.Method, Payload: sr.r.Payload, Users: usersPerSlave, Slave: i + 1}

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
				request := &Request{MType: MSG, URL: sr.r.URL, Headers: sr.headerList,
					Method: sr.r.Method, Payload: sr.r.Payload, Users: 1, Slave: i}

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

func (sr slaveRunner) rampUpSteps() {
	slaveNo := 1

	for i := 0; i < len(sr.usrs); i++ {
		usr := sr.usrs[i]
		d := sr.dur[i]
		unit := sr.units[i]

		request := &Request{MType: MSG, URL: sr.r.URL, Headers: sr.headerList,
			Method: sr.r.Method, Payload: sr.r.Payload, Users: usr, Slave: slaveNo}

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

func (sr slaveRunner) rampUpRegular() {
	usersPerSlave := users / slaves

	for i := 1; i <= slaves; i++ {
		if i == slaves {
			if usersPerSlave*slaves < users {
				diff := users - (usersPerSlave * slaves)
				usersPerSlave += diff
			}
		}

		request := &Request{MType: MSG, URL: sr.r.URL, Headers: sr.headerList,
			Method: sr.r.Method, Payload: sr.r.Payload, Users: usersPerSlave, Slave: i}

		write <- request
		totalUsersGenerated += usersPerSlave
	}
}

func (sr slaveRunner) execute(slave int, usersPerSlave int) {
	request := &Request{MType: MSG, URL: sr.r.URL, Headers: sr.headerList,
		Method: sr.r.Method, Payload: sr.r.Payload, Users: usersPerSlave, Slave: slave}

	write <- request
}
