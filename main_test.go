package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
)

const PORT = "8888"

type Response struct {
	Error error       `json:"error"`
	Data  interface{} `json:"data"`
}

func TestPost(t *testing.T) {
	resp, err := http.PostForm("http://localhost:"+PORT+"/api/user/new", url.Values{"name": {"skyblue"}})
	if err != nil {
		t.Error(err)
	}
	_ = resp
}

func TestGetAll(t *testing.T) {
	resp, err := http.Get("http://localhost:" + PORT + "/api/user/all")
	if err != nil {
		t.Error(err)
	}
	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}
