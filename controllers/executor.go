package controllers

import (
	"fmt"
	"strconv"
	"time"
)

type executor interface {
	run()
	rampUpRegular()
	rampUpLinear()
	rampUpSteps()
	execute(slaves int, usersPerSlave int)
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
