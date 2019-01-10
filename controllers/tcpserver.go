package controllers

import (
	"encoding/json"
	"io"
	"net"
	"strconv"

	"github.com/astaxie/beego/logs"
)

// Stop object sent to clients
type Stop struct {
	MType int
}

// Resp received from clients
type Resp struct {
	MType   int
	Content string
}

var slaves = 0
var write = make(chan *Request)
var stopClient = make(chan string)

func reader(conn net.Conn) {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	log.Debug("Started read goroutine")
	for {
		d := json.NewDecoder(conn)
		r := &Resp{}
		err := d.Decode(r)

		if err == io.EOF {
			conn.Close()
			return
		}

		if r.MType == 1 {
			conn.Close()
			return
		}
		resp, _ := strconv.Atoi(r.Content)
		response <- resp
	}
}

func writer(conn net.Conn) {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	log.Debug("Started write goroutine")
	for {
		select {
		case c := <-write:
			msg, _ := json.Marshal(c)
			log.Debug("Sending msg to client %#v", string(msg))
			_, err := conn.Write([]byte(msg))

			if err != nil {
				conn.Close()
				slaves--
				return
			}
		case <-stopClient:
			log.Debug("Sending stop message to client %#v", conn.RemoteAddr())
			msg, _ := json.Marshal(Stop{MType: 3})
			_, err := conn.Write([]byte(msg))

			if err != nil {
				conn.Close()
				slaves--
				return
			}
		}
	}
}

func server() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	log.Debug("Launching server...")
	server, err := net.Listen("tcp", "0.0.0.0:8082")
	log.Debug("Listening on 0.0.0.0:8082")

	if err != nil {
		log.Error(err.Error())
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			log.Error(err.Error())
		}

		slaves++
		log.Debug("Received connection from client %#v", conn.RemoteAddr())
		log.Debug("Total slaves %#v", slaves)
		go reader(conn)
		go writer(conn)
	}
}

func init() {
	go server()
}
