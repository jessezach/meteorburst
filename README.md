# Meteor Burst

Meteor Burst is simplistic load testing tool. It can be used to quickly run basic loads tests against a REST API.
It provides a simple form based UI where you can provide the URL, Headers, Method, Payload and number of users to simulate. </br>
Meteor Burst will provide you with a realtime average response time, 99th percentile, 90th percentile and 50th percentile numbers.

## Installation
Make sure you have Go setup and `$GOPATH` added to path.</br>
Also make sure you have `$GOPATH/bin` add to path. 
[Check here](https://stackoverflow.com/questions/21001387/how-do-i-set-the-gopath-environment-variable-on-ubuntu-what-file-must-i-edit)

- Install go dep</br>
  `brew install dep` or `go get -u github.com/golang/dep/cmd/dep`

- Navigate to your Gopath/src and clone the repository </br>
```$ cd $GOPATH/src/``` </br>
```git clone https://github.com/jz-jess/meteor-burst.git```

- Navigate to meteor-burst and install dependencies</br>
`$ cd meteor-burst`</br>
`$ dep ensure`

- Install bee tool.</br>
`$ go get -u github.com/beego/bee`

- Inside project root, run server</br>
`bee run`</br>
```______
| ___ \
| |_/ /  ___   ___
| ___ \ / _ \ / _ \
| |_/ /|  __/|  __/
\____/  \___| \___| v1.10.0
2019/01/06 00:31:24 INFO     ▶ 0001 Using 'meteor-burst' as 'appname'
2019/01/06 00:31:24 INFO     ▶ 0002 Initializing watcher...
```

App should be running on `http://localhost:8080/`

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

## Note
Do not refresh the page when tests are running, the data displayed is not stored and will be lost.</br>
For any queries or issues, raise an issue or email me at iamjess988@gmail.com
