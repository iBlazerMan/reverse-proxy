package main

import (
	"log"

	"github.com/iBlazerMan/reverse-proxy/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Error loading config: ", err)
	}
}
