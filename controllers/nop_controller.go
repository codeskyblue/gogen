package controllers

import "errors"

type NopController struct {
	APIController
}

func (this *NopController) NotExisted() {
	this.err = errors.New("not existed api")
}
