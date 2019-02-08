# Meteor Burst
[![Go Report Card](https://goreportcard.com/badge/github.com/jz-jess/meteorburst)](https://goreportcard.com/report/github.com/jz-jess/meteorburst)

Meteor Burst is simplistic load testing tool. It can be used to quickly run basic loads tests against a REST API.
It provides a simple form based UI where you can provide the URL, Headers, Method, Payload and number of users to simulate. </br>
Meteor Burst will provide you with a realtime average response time, 99th percentile, 90th percentile and 50th percentile numbers.

## Installation
Make sure you have Go setup and `$GOPATH` added to path.</br>
Also make sure you have `$GOPATH/bin` add to path. 
[Check here](https://stackoverflow.com/questions/21001387/how-do-i-set-the-gopath-environment-variable-on-ubuntu-what-file-must-i-edit)
- Install bee tool.</br>
`$ go get -u github.com/beego/bee`

- Install code.</br>
`$ go get github.com/jz-jess/meteorburst`

- Inside project root, run server</br>
`$ cd jz-jess/meteorburst`</br>
`$ bee run`</br>

    OR

- Install go dep</br>
  `brew install dep` or `go get -u github.com/golang/dep/cmd/dep`

- Navigate to your Gopath/src and clone the repository </br>
```$ cd $GOPATH/src/``` </br>
```git clone https://github.com/jz-jess/meteorburst.git```

- Navigate to meteorburst and install dependencies</br>
`$ cd meteorburst`</br>
`$ dep ensure`

- Inside project root, run server</br>
`$ bee run`</br>
```______
| ___ \
| |_/ /  ___   ___
| ___ \ / _ \ / _ \
| |_/ /|  __/|  __/
\____/  \___| \___| v1.10.0
2019/01/06 00:31:24 INFO     ▶ 0001 Using 'meteorburst' as 'appname'
2019/01/06 00:31:24 INFO     ▶ 0002 Initializing watcher...
```

App should be running on `http://localhost:8080/`</br>
TCP server should be running on `http://0.0.0.0:8082/`

## How to use
- Navigate to `http://localhost:8080/`
- Fill the required fields with the details of the endpoint you want to load test.
![Alt text](/readme-images/home.png "Home screen")

- Press start
- Tests should begin and metrics should be visible
![Alt text](/readme-images/metrics.png "Metrics")
![Alt text](/readme-images/chart.png "Home screen")

- Press stop button whenever you want to stop the tests.
![Alt text](/readme-images/stop.png "Home screen")

## Duration
You can optionally add duration in minutes or seconds to a test. The tests will run for the specified duration and stop automatically after.</br>
![Alt text](/readme-images/duration.png "Duration")

## Ramp up
Most tests require a pattern of of load generation. You would want to generate load in a linear manner or a step by step manner. Meteor Burst provides ramping up of users in linear fashion or using step. Linear Ramp up duration has to be provided in seconds.</br>
![Alt text](/readme-images/linear.png "Linear")

Step Ramp up can be done as follow:</br>
- Select Step option in the Ramp up dropdown. A table should be displayed.</br>
- Add number of users, duration and unit (seconds, minutes). You can add as many steps as required.</br>
![Alt text](/readme-images/step.png "Step")

## Distributed Load testing
Incase you want to generate a high load. You can install meteor client and run them on different machines.Follow the steps below</br>
- Install meteor `go get github.com/jz-jess/meteor`</br>
- Run `meteor <server-ip>:<port>`, Ex: `meteor 0.0.0.0:8082`</br>
You can see the number of slaves connected in the UI. The tests will automatically run on slaves.</br>
![Alt text](/readme-images/slaves.png "Slaves")

Note: Clients will get disconnected when app server is closed.
## Note
Please star the repository if you find this useful.</br>
For any queries or issues, raise an issue or email me at iamjess988@gmail.com
