package tasks

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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

func startHandler(c *gin.Context) {

	// Initialize stop channel
	stopSignal = make(chan struct{})
	// Start Goroutine process
	go longRunningProcess(stopSignal)

	// Set start time
	startTime = time.Now()

	// Return JSON response
	response := Response{"Process started", true, "", 0}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

func stopHandler(c *gin.Context) {
	// Stop Goroutine process
	stopSignal <- struct{}{}

	// Return JSON response
	response := Response{"Process stopped", false, "", 0}
	c.JSON(http.StatusOK, gin.H{"data": response})
}

func statusHandler(c *gin.Context) {
	// Check if process is still alive
	alive := isProcessAlive()

	if alive {
		durTime := time.Since(startTime).String()

		// Get memory usage stats
		runtime.ReadMemStats(&memStats)
		memUsage := int(memStats.Alloc) / 1024 / 1024

		// Return JSON response
		response := Response{"Process running", true, durTime, memUsage}
		c.JSON(http.StatusOK, gin.H{"data": response})
	} else {
		response := Response{"Process not running", false, "", 0}
		c.JSON(http.StatusOK, gin.H{"data": response})
	}

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

func SetupGameboRoutes(r *gin.Engine) {
	gamebo_req := r.Group("/api")
	gamebo_req.GET("/gamebo_req/start", startHandler)
	gamebo_req.GET("/gamebo_req/stop", stopHandler)
	gamebo_req.GET("/gamebo_req/status", statusHandler)
}
