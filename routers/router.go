package routers

import (
	"github.com/astaxie/beego"
	"github.com/jz-jess/meteorburst/controllers"
)

func init() {
	beego.Router("/", &controllers.HomeController{})
	beego.Router("ws/join", &controllers.HomeController{}, "get:Join")
}
