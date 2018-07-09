package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	//"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const hashUrl = "http://localhost:8080/hash"

func otherClientThings() {
	//client := http.Client{}
	var data io.Reader = nil
	//req, err := http.NewRequest("POST", hashUrl, data)
	req, _ := http.NewRequest("POST", hashUrl, data)
	form := url.Values{}
	form.Add("password", "angryMonkey")
	req.PostForm = form
	data = strings.NewReader(form.Encode())
}

func TestHashServer(t *testing.T) {
	startServer()

	var password, expectedHash string
	password = "angryMonkey"
	expectedHash = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	//resp, err := http.PostForm(hashUrl, url.Values{"password": {password}})
	resp, _ := http.PostForm(hashUrl, url.Values{"password": {password}})
	//body, err := ioutil.ReadAll(resp.Body)
	body, _ := ioutil.ReadAll(resp.Body)
	actualHash := string(body)
	if actualHash != expectedHash {
		t.Error(fmt.Sprintf("Wrong hash for %s: %s", password, actualHash))
	}
}
