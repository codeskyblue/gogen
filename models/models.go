package models

import (
	"fmt"
	"log"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/lunny/xorm"
)

var (
	x *xorm.Engine
)

func init() {
	var err error
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	dataSource := beego.AppConfig.String("mysql")
	if dataSource == "" {
		dataSource = "root:@/trs"
	}
	x, err = xorm.NewEngine("mysql", dataSource)
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	ID   int64  `json:"id" xorm:"id pk autoincr"`
	Name string `json:"name"`
}

func InitDB() (err error) {
	if err = x.Sync(new(User)); err != nil {
		return
	}
	return
}

func init() {
	if err := InitDB(); err != nil {
		log.Fatal(err)
	}
}

// create
func (v *User) Create() error {
	_, err := x.Insert(v)
	return err
}

func (v *User) Update(id int64) error {
	affec, err := x.Id(id).Update(v)
	if err == nil && affec == 0 {
		err = fmt.Errorf("update user(id:%d) failed", id)
	}
	return err
}

func GetUser(id int64) (v *User, err error) {
	v = new(User)
	ok, err := x.Id(id).Get(v)
	if err == nil && !ok {
		err = fmt.Errorf("get user(id:%d) failed", id)
	}
	return
}

func (v *User) Save() error {
	if v.ID == 0 {
		return v.Create()
	}
	return v.Update(v.ID)
}

func AllUser() (vs []User, err error) {
	err = x.Find(&vs)
	return
}

func DelUser(id int64) error {
	affec, err := x.Id(id).Delete(new(User))
	if err == nil && affec == 0 {
		err = fmt.Errorf("user(id:%d) already deleted", id)
	}
	return err
}
