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
// Could be done with OS signals instead, to avoid global
var server *http.Server

// Put the shared data with the mutex
type connCounter struct {
	current   int
	total	int
	mux sync.Mutex
}
// Global data & mutex, because http handlers only get 2 args.
// Maybe could add arguments to handlers somehow, or otherwise avoid global
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
	shutdownTimeout := 10  // seconds
	ccounter.mux.Lock()
	wait := (tries <= shutdownTimeout) && (ccounter.current > 0)
	ccounter.mux.Unlock()
	return wait
}

// shutdownHandler is the http handler for /shutdown
func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	// Problem 1: server.Shutdown aborts connections, despite docs to the contrary. Need different context?
	//   So we just watch the counters before shutting down.
	// Problem 2: We don't prevent new connections before stopping, so if things keep connecting, we'll
	//   never stop. To fix that, fix Problem 1, or stop accepting new connections, or just have a timeout
	//   and drop the hammer (not ideal).
	// Problem 3: The response never finishes. Maybe an OS signal based handler would take care of that.
	// Maybe we could do a notification system rather than polling.
	w.Write([]byte("Shutting down"))
	tries := 0
	wait := getShutdownWait(tries)
	for wait {
		time.Sleep(time.Second)
		wait = getShutdownWait(tries)
		tries += 1
	}
	server.Shutdown(context.Background())
}

// startServer starts the server
func startServer() {
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
