package proxy

import (
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/iBlazerMan/reverse-proxy/serverSelector"
	"github.com/iBlazerMan/reverse-proxy/util"
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

		// add additional context to use when handling response
		ctx := util.WithServerUrl(req.Context(), serverUrl)
		*req = *req.WithContext(ctx)
	}
}

func makeErrorHandler(serverSelector serverSelector.ServerSelector) func(http.ResponseWriter, *http.Request, error) {
	return func(rw http.ResponseWriter, req *http.Request, err error) {
		serverSelector.HandleError(rw, req, err)

		// additional error handling logic here
	}
}

func makeModifyResponse(serverSelector serverSelector.ServerSelector) func(*http.Response) error {
	return func(res *http.Response) error {
		return serverSelector.ModifyResponse(res)

		// additional response modifier logic here
	}
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
		ErrorHandler:   makeErrorHandler(serverSelector),
		ModifyResponse: makeModifyResponse(serverSelector),
	}
}
