package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func ConvertBytesToGBDecimal(bytes uint64) float64 {
	return float64(bytes) / (1000 * 1000 * 1000)
}

func main() {
	http.HandleFunc("/events", sseHandler)

	// Start the server on PORT 8080
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("[FATAL]: Unable to start server: %s", err.Error())
	}
}

func sseHandler(w http.ResponseWriter, r *http.Request) {
	// Setup SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Create timersfor memory and CPU updates
	memTicker := time.NewTicker(time.Second)
	defer memTicker.Stop()
	cpuTicker := time.NewTicker(time.Second)
	defer cpuTicker.Stop()

	clientGone := r.Context().Done()

	// Flush the response to ensure the client receives data continuously
	rc := http.NewResponseController(w)

	for {
		select {
		// handle disconnections
		case <-clientGone:
			fmt.Println("[ERROR]: Client disconnected")
			return

		// pass data from MEM ticker
		case <-memTicker.C:
			memoryStats, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("Failed too get memory stats: %s", err.Error())
				return
			}
			data := fmt.Sprintf(
				"[LOG] Total: %.2f GB, Used: %.2f GB, Percentage: %.2f%%\n",
				ConvertBytesToGBDecimal(memoryStats.Total),
				ConvertBytesToGBDecimal(memoryStats.Used),
				memoryStats.UsedPercent,
			)
			if _, err := fmt.Fprintf(w, "event:mem\ndata:%s\n\n", data); err != nil {
				log.Printf("Failed to write memory event: %s", err.Error())
				return
			}
			rc.Flush()

		// pass data from CPU ticker
		case <-cpuTicker.C:
			cpuStatus, err := cpu.Times(false)
			if err != nil {
				log.Printf("[ERROR]: Failed to get CPU status: %s", err.Error())
				return
			}
			data := fmt.Sprintf(
				"[LOG] User: %.2f, System: %.2f, Idle: %.2f",
				cpuStatus[0].User, cpuStatus[0].System, cpuStatus[0].Idle,
			)
			if _, err := fmt.Fprintf(w, "event:cpu\ndata:%s\n\n", data); err != nil {
				log.Printf("[ERROR]: Failed to write CPU event: %s", err.Error())
				return
			}
			rc.Flush()
		}
	}
}
