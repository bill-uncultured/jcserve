package main

import (
	"github.com/bill-uncultured/jcserve/jchash"
	"log"
	"net/http"
	"time"
)

func hashHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	password := r.FormValue("password")
	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request - password field missing or invalid"))
		return
	}
	pwHash := jchash.HashPassword(password)
	w.Write([]byte(pwHash))
	time.Sleep(5 * time.Second)  // TODO: Make sure we're not blocking other threads
}

func startServer() {
	http.HandleFunc("/hash", hashHandler)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}

func main() {
	startServer()
}
