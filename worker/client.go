package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Message struct {
	MType   int
	URL     string
	Headers []string
	Method  string
	Payload string
	Users   int
}

type Response struct {
	MType   int
	Content string
}

var quit chan bool

func getFireSignalsChannel() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
		syscall.SIGINT,  // Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGKILL, // "always fatal", "SIGKILL and SIGSTOP may not be caught by a program"
		syscall.SIGHUP,  // "terminal is disconnected"
	)
	return c
}

func killProcess(conn net.Conn) {
	exitChan := getFireSignalsChannel()
	<-exitChan
	fmt.Println("Interrupted. Exiting..")
	sendExitMsg(conn)
	os.Exit(1)
}

func main() {
	args := os.Args[1:]

	conn, err := net.Dial("tcp", args[0])

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else if len(args) == 0 {
		fmt.Println("Provide server host:<port>")
		os.Exit(1)
	}

	fmt.Printf("Connected to %v\n", args[0])

	go killProcess(conn)

	for {
		d := json.NewDecoder(conn)
		msg := &Message{}
		err := d.Decode(msg)

		if err != nil {
			fmt.Print(err.Error())
			sendExitMsg(conn)
			conn.Close()
			break
		}

		if msg.MType == 2 {
			quit = make(chan bool)

			fmt.Printf("Attacking url %v\n", msg.URL)
			for i := 0; i < msg.Users; i++ {
				go func() {
					for {
						select {
						case <-quit:
							return
						default:
							meteorBurst(msg.URL, msg.Payload, msg.Method, msg.Headers, conn)
						}
					}
				}()
			}
		} else if msg.MType == 3 {
			fmt.Println("Stopping attack..")
			close(quit)
		}
	}
}

func sendExitMsg(conn net.Conn) {
	fmt.Println("Exiting. Sending exit message to server..")
	resp := Response{MType: 1, Content: "dead"}
	m, _ := json.Marshal(resp)
	conn.Write([]byte(m))
}

func meteorBurst(url string, payload string, method string, headers []string, conn net.Conn) {
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
		client.Do(req)
		endTime := time.Now().UnixNano() / int64(time.Millisecond)

		responseTime := int(endTime - startTime)
		resp := Response{MType: 2, Content: strconv.Itoa(responseTime)}
		m, _ := json.Marshal(resp)
		conn.Write([]byte(m))
	} else {
		fmt.Println(err.Error())
	}
}
