package main

import (
	"github.com/astaxie/beego"
	"github.com/shxsun/gails/controllers"
)

var (
	urlbase = "/api"
)

func main() {
	beego.AutoRender = false
	beego.DirectoryIndex = true

	beego.Router("/api/template/new", &controllers.UserController{}, "post:Save")
	beego.Router("/api/template/all", &controllers.UserController{}, "get:All")
	beego.Router("/api/template/:id(\\d+)", &controllers.UserController{}) // GET + PUT  + DELETE

	beego.Run()
}
