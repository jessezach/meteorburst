package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/validation"
	mapset "github.com/deckarep/golang-set"
	"github.com/gorilla/websocket"
)

// HomeController : controller
type HomeController struct {
	beego.Controller
}

// RequestDetails form details
type RequestDetails struct {
	URL     string `form:"url" valid:"Required"`
	Headers string `form:"headers"`
	Method  string `form:"method" valid:"Required"`
	Payload string `form:"payload"`
	Users   int    `form:"users" valid:"Required"`
}

// Get request
func (c *HomeController) Get() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	flash := beego.ReadFromRequest(&c.Controller)
	if _, ok := flash.Data["notice"]; ok {
		// Display settings successful
		c.Data["notice"] = true
	} else if _, ok = flash.Data["error"]; ok {
		c.Data["error"] = true
	}

	if !running {
		c.TplName = "home.tpl"
	} else {
		c.Data["users"] = users
		c.TplName = "burst.tpl"
	}
}

// Post request
func (c *HomeController) Post() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)

	command := c.GetString("command")

	if command == "start" {
		r := &RequestDetails{}
		c.ParseForm(r)
		flash := beego.NewFlash()

		valid := validation.Validation{}
		isValid, err := valid.Valid(r)

		if !isValid {
			for _, e := range valid.Errors {
				flash.Error("%#v %#v", e.Key, e.Message)
				break
			}
			c.Redirect("/", 302)
		}

		if err == nil {
			quit = make(chan bool)
			var headerList []string

			if r.Headers != "" {
				headerList = strings.Split(r.Headers, ";")
			}

			running = true
			users = r.Users
			setStartTime(time.Now().UnixNano() / int64(time.Millisecond))

			for i := 0; i < r.Users; i++ {
				log.Debug("Starting user %#v", i+1)
				go func() {
					for {
						select {
						case <-quit:
							log.Debug("Returning from go routine")
							return
						default:
							c.meteorBurst(r.URL, r.Method, r.Payload, headerList)
						}
					}
				}()
			}
		} else {
			flash.Error("%#v", err.Error())
		}

		flash.Store(&c.Controller)
		c.Redirect("/", 302)

	} else if command == "stop" {
		if quit != nil {
			close(quit)
			running = false
			users = 0
			setStartTime(0)
		}
		c.Data["json"] = "{'stopped': true}"
		c.ServeJSON()
	}
}

// MeteorBurst makes a REST call to the provided endpoint
func (c *HomeController) meteorBurst(url string, method string, payload string, headers []string) {
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
				req.Header.Add(header[0], header[1])
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

// Join creates a new websocket object for a new client and adds to subsriber list
func (c *HomeController) Join() {
	ws, err := websocket.Upgrade(c.Ctx.ResponseWriter, c.Ctx.Request, nil, 1024, 1024)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Ctx.ResponseWriter, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		beego.Error("Cannot setup WebSocket connection:", err)
		return
	}

	Join(ws)
	c.Data["success"] = true
	c.ServeJSON()
}

// broadcastWebSocket broadcasts messages to WebSocket users.
func broadcastWebSocket(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		beego.Error("Fail to marshal event:", err)
		return
	}

	for sub := subscribers.Front(); sub != nil; sub = sub.Next() {
		// Immediately send event to WebSocket users.
		ws := sub.Value.(Subscriber).Conn
		if ws != nil {
			if ws.WriteMessage(websocket.TextMessage, data) != nil {
				// User disconnected.
				unsubscribe <- ws
			}
		}
	}
}
