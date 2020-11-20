package controllers

import (
	"container/list"
	"net"
	"time"

	"github.com/gorilla/websocket"
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

// RequestDetails form details
type RequestDetails struct {
	URL      string `form:"url" valid:"Required"`
	Headers  string `form:"headers"`
	Method   string `form:"method" valid:"Required"`
	Payload  string `form:"payload"`
	Users    int    `form:"users" valid:"Required"`
	Duration int    `form:"duration"`
	Format   string `form:"format"`
	RampType string `form:"ramp-type"`
	RampTime int    `form:"ramp"`
	RampStep string `form:"step"`
}

var quit chan bool
var running bool
var testStartTime int64

var (
	users               = 0
	timer               = time.NewTimer(time.Second)
	totalUsersGenerated = 0
	response            = make(chan int)
	responses           = []int{}
	totalRequests       = 0
	batchSize           = 0
	slaves              = 0
	write               = make(chan *Request)
	stopClient          = make(chan string)
	removeWriter        = make(chan net.Conn)
	writers             = list.New()
	// Channel for new join users.
	subscribe = make(chan Subscriber, 10)
	// Channel for exit users.
	unsubscribe = make(chan *websocket.Conn, 10)
	// Send events here to publish them.
	publish              = make(chan Event)
	subscribers          = list.New()
	responseStats        = make(map[string]int)
	responseStatsChannel = make(chan int)
	httpErrorChannel     = make(chan string)
)

// Constants for type of message
const (
	MESSAGE           = 2
	TOTAL             = 3
	P90               = 4
	P99               = 5
	P50               = 6
	RPS               = 7
	SLAVE             = 8
	STOPPED           = 9
	USERS             = 10
	ERROR             = 11
	CLOSED_CONNECTION = 1
	MSG               = 2
	STOP_TEST         = 3
	STATUS_CODE_STATS = 12
	HTTP_ERROR        = 13
)

func setStartTime(time int64) {
	testStartTime = time
}
