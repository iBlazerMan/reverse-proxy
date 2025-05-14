package serverSelector

import (
	"net/http"
	"net/url"
)

type ServerSelector interface {
	SelectServer() *url.URL
	ModifyResponse(res *http.Response) error
	HandleError(rw http.ResponseWriter, req *http.Request, err error)
}
