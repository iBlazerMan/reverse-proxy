package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/iBlazerMan/reverse-proxy/config"
	"github.com/iBlazerMan/reverse-proxy/proxy"
	"github.com/iBlazerMan/reverse-proxy/serverSelector"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}

	var reverseProxy *httputil.ReverseProxy
	switch config.BalanceAlgorithm {
	case "SingleServer":
		reverseProxy = proxy.NewProxy(serverSelector.NewSingleServer(config.ServerAddresses[0]))
	case "RoundRobin":
		reverseProxy = proxy.NewProxy(serverSelector.NewRoundRobin(config.ServerAddresses))
		// case "LeastConnection":
		// 	reverseProxy = proxy.NewProxy()
	default:
		log.Fatalf("Undefined balance algorithm: %v", config.BalanceAlgorithm)
	}

	proxyServer := http.Server{
		Addr:         ":8080",
		Handler:      reverseProxy,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("Starting proxy")
	if err := proxyServer.ListenAndServe(); err != nil {
		log.Fatal("Proxy failed to start: ", err)
	}
}
