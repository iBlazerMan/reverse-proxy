package serverSelector

import (
	"errors"
	"log"
	"net/http"
	"net/url"
	"sync/atomic"

	"github.com/iBlazerMan/reverse-proxy/util"
)

type serverInstanceLc struct {
	serverInstance
	// Int64 is used here instead of Uint64 to allow Atomic.Add(-1) when processing a response.
	// currReq should never be negative.
	currReq atomic.Int64
}

type LeastConnection struct {
	defaultSelector
	servers []serverInstanceLc
}

func NewLeastConnection(serverUrls []string) *LeastConnection {
	lc := LeastConnection{}
	for _, urlString := range serverUrls {
		parsedUrl, err := url.Parse(urlString)
		if err != nil {
			log.Fatalf("Error parsing server url %q: %v", urlString, err)
		}
		lc.servers = append(lc.servers, serverInstanceLc{
			serverInstance: serverInstance{parsedUrl},
			currReq:        atomic.Int64{},
		})

	}
	return &lc
}

func (lc *LeastConnection) SelectServer() *url.URL {
	chosenServer := &lc.servers[0]
	minReqCount := chosenServer.currReq.Load()
	for i := 1; i < len(lc.servers); i++ {
		currReqCount := lc.servers[i].currReq.Load()
		if currReqCount < minReqCount {
			chosenServer = &lc.servers[i]
			minReqCount = currReqCount
		}
	}
	chosenServer.currReq.Add(1)

	return chosenServer.url
}

func (lc *LeastConnection) ModifyResponse(res *http.Response) error {
	// extract server address and decrement server's currReq count
	serverUrl, err := util.GetServerUrl(res.Request.Context())

	if err != nil {
		return err
	}

	for i := range lc.servers {
		if *(lc.servers[i].url) == *serverUrl {
			lc.servers[i].currReq.Add(-1)
			return nil
		}
	}

	return errors.New("no existing server's url match the response context's url, no server current request count changed")
}
