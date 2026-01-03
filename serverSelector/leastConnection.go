package serverSelector

import (
	"fmt"
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

	// DEBUG
	// for i := range lc.servers {
	// 	log.Printf("Server %s has %d current requests", lc.servers[i].url.String(), lc.servers[i].currReq.Load())
	// }
	// log.Printf("\n")

	return chosenServer.url
}

func (lc *LeastConnection) done(serverUrl *url.URL) error {
	// extract server address and decrement server's currReq count
	if serverUrl == nil {
		return fmt.Errorf("no server url provided to done function")
	}

	for i := range lc.servers {
		if lc.servers[i].url.String() == serverUrl.String() {
			lc.servers[i].currReq.Add(-1)
			return nil
		}
	}

	return fmt.Errorf("no existing server's url match the response context's url: %s", serverUrl.String())
}

func (lc *LeastConnection) ModifyResponse(res *http.Response) error {
	serverUrl, err := util.GetServerUrl(res.Request.Context())
	if err != nil {
		return err
	}
	return lc.done(serverUrl)
}

// still need to decrement the server's currReq count on error, then call downstream error handler
func (lc *LeastConnection) HandleError(rw http.ResponseWriter, req *http.Request, err error) {
	serverUrl, _ := util.GetServerUrl(req.Context())
	_ = lc.done(serverUrl)
	lc.defaultSelector.HandleError(rw, req, err)
}
