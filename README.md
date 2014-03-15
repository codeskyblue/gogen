gails
=====

beego + xorm skeleton, for quickly REST-API develop.

## how to use
	go get github.com/shxsun/gails

	cd $GOPATH/src/
	mkdir test; cd test
	gails book name:string

	# create table user with column name,school
	gails user name:string school:string

	# start program
	bee run

### how to use this api
	curl localhost:8080/api/book/new -d name=parkour
	curl localhost:8080/api/book/all

## only work well on linux, mac
Even through I test on my machine for many times, but it is better to backup your code, before you use it.

## how to make a template
default template address is in <https://github.com/shxsun/gails-default>

* lines contails `XXXX` will be deleted.
* chars before `NNNN` will be deleted also include `NNNN` itself.
* use golang default `text/template`

A yaml file is needed. content is like

	notice: |
		This template is created by skyblue. in 2014/03/15
		You need to change setting in conf/app.conf,
		And create a mysql database before you use -bee run- to start it.
	file-rename:
	  user_router.go: "{{.Table}}_router.go"
	  controllers/user.go: controllers/{{.Table}}.go
	  models/user.go: models/{{.Table}}.go
	string-rename:
	  User: "{{.Table|title}}"

`title` is default function, which convert word `name` to `Name`.
