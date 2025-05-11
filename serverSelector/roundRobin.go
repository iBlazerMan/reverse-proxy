package serverSelector

import (
	"log"
	"net/url"
	"sync/atomic"
)

type serverInstance struct {
	url *url.URL
}

type RoundRobin struct {
	servers []serverInstance
	index   atomic.Uint64
}

func NewRoundRobin(serverUrls []string) *RoundRobin {
	rr := RoundRobin{}
	for _, urlString := range serverUrls {
		parsedUrl, err := url.Parse(urlString)
		if err != nil {
			log.Fatalf("Error parsing server url %q: %v", urlString, err)
		}
		rr.servers = append(rr.servers, serverInstance{parsedUrl})

	}
	return &rr
}

func (rr *RoundRobin) SelectServer() *url.URL {
	// wraps back to 0 on overflow, and practically impossible for u_int64 overflow to occur
	// with a counter
	currIndex := rr.index.Add(1)
	return rr.servers[currIndex%uint64(len(rr.servers))].url
}
