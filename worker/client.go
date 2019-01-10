package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
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

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8082")

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	fmt.Println("Connected to 127.0.0.1:8082")
	for {
		d := json.NewDecoder(conn)
		msg := &Message{}
		err := d.Decode(msg)

		if err != nil {
			fmt.Print(err.Error())
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
	}
}
