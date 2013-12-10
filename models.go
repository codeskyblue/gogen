// models.go
package main

import (
	"fmt"
	"log"

	"github.com/lunny/xorm"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Engine *xorm.Engine
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	var err error
	Engine, err = xorm.NewEngine("sqlite3", "./test.db")
	checkError(err)
	err = Engine.Sync(new(Post))
	checkError(err)
}

type Post struct {
	id      int64  `json:"id" xorm:"pk"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (this *Post) Get(id int64) error {
	ok, err := Engine.Where("id=?", id).Get(this)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("can't found id(%d) for Post", id)
	}
	return nil
}

func (this *Post) Add() error {
	return nil
}

func (this *Post) Del() error {
	return nil
}
