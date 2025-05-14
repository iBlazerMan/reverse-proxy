package serverSelector

import (
	"log"
	"net/url"
)

type SingleServer struct {
	defaultSelector
	serverUrl *url.URL
}

func NewSingleServer(serverUrl string) *SingleServer {
	parsedUrl, err := url.Parse(serverUrl)
	if err != nil {
		log.Fatalf("Error parsing server url %q: %v", serverUrl, err)
	}
	return &SingleServer{defaultSelector{}, parsedUrl}
}

func (ss *SingleServer) SelectServer() *url.URL {
	return ss.serverUrl
}
