package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

const hashUrl = "http://localhost:8080/hash"

func TestMain(m *testing.M) {
	go startServer()
	code := m.Run()
	server.Close()
	os.Exit(code)
}

func TestHashPass(t *testing.T) {
	var password, expectedHash string
	password = "angryMonkey"
	expectedHash = "ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q=="
	resp, err := http.PostForm(hashUrl, url.Values{"password": {password}})
	if err != nil {
		t.Error("Error connecting password request: ", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error("Error reading hash reponse: ", err)
	}
	actualHash := string(body)
	if actualHash != expectedHash {
		t.Error("Error: wrong hash:", actualHash)
	}
}

func TestNoPass(t *testing.T) {
	resp, err := http.PostForm(hashUrl, url.Values{})
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		t.Error("Error: got success when expected error: ", string(body))
	}
}


