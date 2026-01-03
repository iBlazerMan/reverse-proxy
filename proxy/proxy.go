package proxy

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/iBlazerMan/reverse-proxy/serverSelector"
)

func makeErrorHandler(serverSelector serverSelector.ServerSelector) func(http.ResponseWriter, *http.Request, error) {
	return func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Printf("Proxy Error: Failed to reach backend %s: %v", req.URL.String(), err)
		serverSelector.HandleError(rw, req, err)

		// additional error handling logic here
	}
}

func makeModifyResponse(serverSelector serverSelector.ServerSelector) func(*http.Response) error {
	return func(res *http.Response) error {
		return serverSelector.ModifyResponse(res)
	}
}

func NewProxy(serverSelector serverSelector.ServerSelector) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			targetUrl := serverSelector.SelectServer()

			pr.SetURL(targetUrl)
			pr.Out.Host = targetUrl.Host
		},
		Transport: &http.Transport{
			Proxy: nil,
			// TODO: add user defined options here
			DialContext: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          100,
			MaxIdleConnsPerHost:   20,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   3 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		ErrorHandler:   makeErrorHandler(serverSelector),
		ModifyResponse: makeModifyResponse(serverSelector),
	}
}
