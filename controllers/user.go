package controllers

import "github.com/shxsun/gails/models"

type UserController struct {
	APIController
}

// POST /api/user/save
func (this *UserController) Save() {
	v := new(models.User)
	v.ID, _ = this.GetInt("id")
	v.Name = this.MustString("name")
	this.err = v.Save()
}

// PUT /api/user/:id
func (this *UserController) Put() {
	v := new(models.User)
	v.ID, _ = this.GetInt(":id")
	v.Name = this.MustString("name")
	this.err = v.Save()
}

// GET /api/user/all
func (this *UserController) All() {
	this.data, this.err = models.AllUser()
}

// DELELE /api/user/delete
func (this *UserController) Delete() {
	id := this.MustInt64(":id")
	this.err = models.DelUser(id)
}
