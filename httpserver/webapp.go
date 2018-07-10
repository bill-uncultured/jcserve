// httpserver returns hashes for password strings via http post
package main

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"
	"github.com/bill-uncultured/jcserve/jchash"
)

// Global so we can call the shutdown from a handler
var server *http.Server
var shuttingDown bool

// Put the shared data with the mutex
type connCounter struct {
	current   int
	total	int
	mux sync.Mutex
}
// Global data & mutex, because http handlers only get 2 args.
var ccounter connCounter

// upCounter increments current & total shared counters
func upCounter() {
	ccounter.mux.Lock()
	ccounter.current += 1
	ccounter.total += 1
	// Could make the logs debug only
	log.Print("total conns:", ccounter.total)
	log.Print("current conns:", ccounter.current)
	ccounter.mux.Unlock()
}

// downCounter decrements current shared counter
func downCounter() {
	ccounter.mux.Lock()
	ccounter.current -= 1
	// Could make the logs debug only
	log.Print("current conns:", ccounter.current)
	ccounter.mux.Unlock()
}

// hashHandler is the http handler for /hash
func hashHandler(w http.ResponseWriter, r *http.Request) {
	if shuttingDown {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("503 - Service Unavailable - shutting down"))
		return
	}
	upCounter()
	r.ParseForm()
	password := r.FormValue("password")
	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request - password field missing or invalid"))
		return
	}
	time.Sleep(5 * time.Second)
	pwHash := jchash.HashPassword(password)
	w.Write([]byte(pwHash))
	downCounter()
}

// getShutdownWait: whether we should wait to shut down
func getShutdownWait(tries int) bool {
	shutdownTimeout := 15  // seconds
	ccounter.mux.Lock()
	wait := (tries <= shutdownTimeout) && (ccounter.current > 0)
	ccounter.mux.Unlock()
	return wait
}

// shutdownServer shuts down the server
func shutdownServer() {
	shuttingDown = true
	tries := 0
	wait := getShutdownWait(tries)
	for wait {
		time.Sleep(time.Second)
		wait = getShutdownWait(tries)
		tries += 1
	}
	server.Shutdown(context.Background())
}

// shutdownHandler is the http handler for /shutdown
func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutting down"))
	go shutdownServer()
}

// startServer starts the server
func startServer() {
	shuttingDown = false
	http.HandleFunc("/hash", hashHandler)
	http.HandleFunc("/shutdown", shutdownHandler)
	server = &http.Server {
		Addr: ":8080",
	}
	err := server.ListenAndServe()
	log.Fatal("Fatal error:", err)
}

func main() {
	startServer()
}
