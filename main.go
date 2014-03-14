package main

import (
	"log"

	"github.com/astaxie/beego"
	"github.com/shxsun/gails/controllers"
	"github.com/shxsun/gails/models"
)

func main() {
	beego.AutoRender = false
	beego.DirectoryIndex = true

	if err := models.InitDB(); err != nil {
		log.Fatal(err)
	}
	beego.Router("/api/user/new", &controllers.UserController{}, "post:Save")
	beego.Router("/api/user/all", &controllers.UserController{}, "get:All")
	beego.Router("/api/user/:id(\\d+)", &controllers.UserController{}) // GET + PUT  + DELETE

	beego.Run()
}
