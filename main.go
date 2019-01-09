package main

import (
	"github.com/astaxie/beego"
	_ "github.com/jz-jess/meteorburst/routers"
)

func main() {
	beego.Run()
}
