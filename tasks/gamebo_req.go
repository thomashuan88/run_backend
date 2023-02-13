package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type Response struct {
	Message    string `json:"message"`
	Alive      bool   `json:"alive"`
	Runtime    string `json:"runtime"`
	MemUsageMB int    `json:"mem_usage_mb"`
}

var startTime time.Time
var memStats runtime.MemStats
var stopSignal = make(chan struct{})

func startHandler(w http.ResponseWriter, r *http.Request) {

	// Initialize stop channel
	stopSignal = make(chan struct{})
	// Start Goroutine process
	go longRunningProcess(stopSignal)

	// Set start time
	startTime = time.Now()

	// Return JSON response
	response := Response{"Process started", true, "", 0}
	json.NewEncoder(w).Encode(response)
}

func stopHandler(w http.ResponseWriter, r *http.Request) {
	// Stop Goroutine process
	stopSignal <- struct{}{}

	// Return JSON response
	response := Response{"Process stopped", false, "", 0}
	json.NewEncoder(w).Encode(response)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	// Get runtime stats
	durTime := time.Since(startTime).String()

	// Get memory usage stats
	runtime.ReadMemStats(&memStats)
	memUsage := int(memStats.Alloc) / 1024 / 1024

	// Check if process is still alive
	alive := isProcessAlive()

	// Return JSON response
	response := Response{"Process running", alive, durTime, memUsage}
	json.NewEncoder(w).Encode(response)
}

func longRunningProcess(stop chan struct{}) {
	// Your long running process logic here
	for {
		select {
		case <-stop:
			// Stop the long running process
			return
		default:
			// Do some work here
			fmt.Println("do something!")
			time.Sleep(2 * time.Second)
		}
	}
}

func isProcessAlive() bool {
	// Get the list of all Goroutines
	list := make([]byte, 1024*1024)
	n := runtime.Stack(list, true)

	// Loop through the list of Goroutines and check if the longRunningProcess Goroutine is present
	for _, line := range strings.Split(string(list[:n]), "\n") {
		if strings.Contains(line, "longRunningProcess") {
			return true
		}
	}

	return false
}

func Router() {
	http.HandleFunc("/gamebo_req/start", startHandler)
	http.HandleFunc("/gamebo_req/stop", stopHandler)
	http.HandleFunc("/gamebo_req/status", statusHandler)
	http.ListenAndServe(":8080", nil)
}
