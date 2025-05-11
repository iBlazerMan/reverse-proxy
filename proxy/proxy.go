package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/iBlazerMan/reverse-proxy/serverSelector"
)

func makeDirector(serverSelector serverSelector.ServerSelector) func(*http.Request) {
	return func(req *http.Request) {
		var serverUrl *url.URL = serverSelector.SelectServer()
		req.URL.Scheme = serverUrl.Scheme
		req.URL.Host = serverUrl.Host
		req.Host = serverUrl.Host
		req.URL.Path = serverUrl.Path + req.URL.Path

		// set forwarding header
		req.Header.Set("X-Forwarded-For", req.RemoteAddr)
		req.Header.Set("X-Forwarded-Host", req.Host)
		req.Header.Del("Connection")
	}
}

func errorHandler(rw http.ResponseWriter, req *http.Request, err error) {
	rw.WriteHeader(http.StatusBadGateway)
}

func NewProxy(serverSelector serverSelector.ServerSelector) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: makeDirector(serverSelector),
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
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
		ErrorHandler: errorHandler,
	}
}
