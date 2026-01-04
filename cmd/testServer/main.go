package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const minRandDelay = 0.5
const maxRandDelay = 10.0

func handleRandDelay(port string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] Received %s", port, r.URL.Path)

		delay := minRandDelay + rand.Float64()*(maxRandDelay-minRandDelay)
		sleepDuration := time.Duration(delay * float64(time.Second))

		time.Sleep(sleepDuration)

		statusCode := http.StatusOK
		if statusCodeGiven := r.URL.Query().Get("statusCode"); statusCodeGiven != "" {
			var err error
			statusCode, err = strconv.Atoi(statusCodeGiven)
			if err != nil {
				statusCode = http.StatusOK
			}
		}

		msg := fmt.Sprintf("Response from server on port %s after %v with status code %d", port, sleepDuration, statusCode)
		w.WriteHeader(statusCode)
		w.Write([]byte(msg))
		log.Printf("[%s] Replied: %s", port, msg)
	}
}

func main() {
	testServerCount := flag.Int("count", 5, "Number of test servers to start")
	flag.Parse()

	currPort := 8081
	var wg sync.WaitGroup

	var serverUrls []string

	for i := 0; i < *testServerCount; i++ {
		port := fmt.Sprintf("%d", currPort+i)
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			addr := "127.0.0.1:" + p
			mux := http.NewServeMux()
			mux.HandleFunc("/randDelay", handleRandDelay(p))

			if err := http.ListenAndServe(addr, mux); err != nil {
				log.Printf("Server on port %s failed: %v", p, err)
			}
		}(port)
		serverUrls = append(serverUrls, "http://127.0.0.1:"+port)
	}

	log.Printf("Started %d test servers on addresses:\n%s", *testServerCount, strings.Join(serverUrls, ","))
	wg.Wait()
}
