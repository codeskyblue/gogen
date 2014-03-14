package main

import (
	"github.com/astaxie/beego"
	"github.com/shxsun/gails/controllers"
)

func main() {
	beego.AutoRender = false
	beego.DirectoryIndex = true

	beego.Router("/api/user/new", &controllers.UserController{}, "post:Save")
	beego.Router("/api/user/all", &controllers.UserController{}, "get:All")
	beego.Router("/api/user/:id(\\d+)", &controllers.UserController{}) // GET + PUT  + DELETE

	beego.Run()
}
