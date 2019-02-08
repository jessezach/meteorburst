package controllers

import (
	"container/list"
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

// Writer contains a tcp socket connection object
type Writer struct {
	Conn net.Conn
}

// Constants for tcp messages
const (
	CLOSED_CONNECTION = 1
	MSG               = 2
	STOP_TEST         = 3
)

var (
	slaves       = 0
	write        = make(chan *Request)
	stopClient   = make(chan string)
	removeWriter = make(chan net.Conn)
	writers      = list.New()
)

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

		// Client closed connection
		if r.MType == CLOSED_CONNECTION {
			removeWriter <- conn // Removes connection from writer list
			slaves--
			publish <- newEvent(SLAVE, strconv.Itoa(slaves))
			log.Debug("Connection closed by a client. Total slaves %v", slaves)
			return
		}
		resp, _ := strconv.Atoi(r.Content)
		response <- resp
	}
}

func writer() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	log.Debug("Started write goroutine")
	for {
		select {
		case c := <-write:
			msg, _ := json.Marshal(c)
			log.Debug("Sending msg to client %#v", string(msg))
			s := c.Slave
			counter := 1

			for wr := writers.Front(); wr != nil; wr = wr.Next() {
				if s == counter {
					_, err := wr.Value.(Writer).Conn.Write([]byte(msg))

					if err != nil {
						removeWriter <- wr.Value.(Writer).Conn
					}
					break
				}
				counter++
			}
		case <-stopClient:
			log.Debug("Sending stop message to all clients")
			msg, _ := json.Marshal(Stop{MType: STOP_TEST})

			for wr := writers.Front(); wr != nil; wr = wr.Next() {
				_, err := wr.Value.(Writer).Conn.Write([]byte(msg))

				if err != nil {
					removeWriter <- wr.Value.(Writer).Conn
				}
			}
		case cn := <-removeWriter:
			//Client closed connection hence stop stale writer connection
			for wr := writers.Front(); wr != nil; wr = wr.Next() {
				if wr.Value.(Writer).Conn == cn {
					writers.Remove(wr)
					wr.Value.(Writer).Conn.Close()
					log.Debug("Removed connection from writers")
				}
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
		publish <- newEvent(SLAVE, strconv.Itoa(slaves))

		log.Debug("Received connection from client %#v", conn.RemoteAddr().String())
		log.Debug("Total slaves %#v", slaves)
		go reader(conn)

		writers.PushBack(Writer{Conn: conn})
	}
}

func init() {
	go server()
	go writer()
}
