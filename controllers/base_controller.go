package controllers

import (
	"fmt"

	"github.com/astaxie/beego"
)

type APIController struct {
	beego.Controller
	err  error
	data interface{}
}

// key must exists, ^_^
func (this *APIController) MustString(key string) string {
	v := this.GetString(key)
	if v == "" {
		this.err = fmt.Errorf("require filed: %s", key)
		this.data = "orz!!"
		this.teminate()
	}
	return v
}

func (this *APIController) MustInt64(key string) int64 {
	v, err := this.GetInt(key)
	if err != nil {
		this.err = err
		this.teminate()
	}
	return v
}

func (this *APIController) teminate() {
	r := struct {
		Error interface{} `json:"error"`
		Data  interface{} `json:"data"`
	}{}
	if this.err != nil {
		r.Error = this.err.Error()
	}
	r.Data = this.data
	this.Data["json"] = r
	this.ServeJson()
	this.StopRun()
}

func (this *APIController) Finish() {
	this.teminate()
}
