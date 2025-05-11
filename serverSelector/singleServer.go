package serverSelector

import (
	"log"
	"net/url"
)

type SingleServer struct {
	serverUrl *url.URL
}

func NewSingleServer(serverUrl string) *SingleServer {
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		log.Fatalf("Error parsing server url %q: %v", serverUrl, err)
	}
	return &SingleServer{parsedUrl}
}

func (ss *SingleServer) SelectServer() *url.URL {
	return ss.serverUrl
}
