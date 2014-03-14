package models

import (
	"github.com/astaxie/beego"
	"github.com/lunny/xorm"
)

var (
	x      *xorm.Engine
	tables = []interface{}{}
)

func InitDB() (err error) {
	dataSource := beego.AppConfig.String("mysql")
	if dataSource == "" {
		dataSource = "root:@/default"
	}
	x, err = xorm.NewEngine("mysql", dataSource)
	if err != nil {
		return
	}
	err = x.Sync(tables...)
	return
}
