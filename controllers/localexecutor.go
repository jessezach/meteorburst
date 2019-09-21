package controllers

import (
	"time"

	"github.com/astaxie/beego/logs"
)

type localRunner struct {
	r          *RequestDetails
	headerList []string
	usrs       []int
	dur        []int
	units      []string
}

func (lc localRunner) run() {
	switch lc.r.RampType {
	case "linear":
		lc.rampUpLinear()
	case "step":
		lc.rampUpSteps()
	default:
		lc.rampUpRegular()
	}
}

func (lc localRunner) rampUpLinear() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	timePerUser := float64(lc.r.RampTime) / float64(lc.r.Users)
	millis := timePerUser * 1000

	log.Debug("Total users %#v", lc.r.Users)
	for i := 0; i < lc.r.Users; i++ {
		go func() {
			for {
				select {
				case <-quit:
					log.Debug("Returning from go routine")
					return
				default:
					execute(lc.r.URL, lc.r.Method, lc.r.Payload, lc.headerList)
				}
			}
		}()
		totalUsersGenerated++
		time.Sleep(time.Millisecond * time.Duration(millis))
	}
	return
}

func (lc localRunner) rampUpSteps() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	for i := 0; i < len(lc.usrs); i++ {
		userCount := lc.usrs[i]
		duration := lc.dur[i]
		unit := lc.units[i]

		for i := 0; i < userCount; i++ {
			go func() {
				for {
					select {
					case <-quit:
						log.Debug("Returning from go routine")
						return
					default:
						execute(lc.r.URL, lc.r.Method, lc.r.Payload, lc.headerList)
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

func (lc localRunner) rampUpRegular() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	for i := 0; i < lc.r.Users; i++ {
		log.Debug("Starting user %#v", i+1)
		go func() {
			for {
				select {
				case <-quit:
					log.Debug("Returning from go routine")
					return
				default:
					execute(lc.r.URL, lc.r.Method, lc.r.Payload, lc.headerList)
				}
			}
		}()
		totalUsersGenerated++
	}
	return
}
