package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

const (
	TargetURL    = "http://localhost:8080/randDelay"
	SendInterval = 100 * time.Millisecond
)

func main() {
	ticker := time.NewTicker(SendInterval)
	defer ticker.Stop()

	log.Printf("Starting Load Test against %s", TargetURL)
	log.Printf("Sending a request every %v", SendInterval)

	requestID := 0

	for t := range ticker.C {
		requestID++
		go sendRequest(requestID, t)
	}
}

func sendRequest(id int, startTime time.Time) {
	// 1. Create a custom client with a timeout to prevent hanging forever
	client := http.Client{
		Timeout: 15 * time.Second,
	}

	log.Printf("[Req #%d] Sent", id)

	resp, err := client.Get(TargetURL)
	if err != nil {
		log.Printf("[Req #%d] FAILED: %v", id, err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	duration := time.Since(startTime)

	log.Printf("[Req #%d] Finished in %v | Status: %s | Body: %s", id, duration.Round(time.Millisecond), resp.Status, string(body))
}
