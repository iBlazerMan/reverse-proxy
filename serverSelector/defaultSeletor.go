/*
The default behaviour of the ServerSelector's ModifyResponse and Handle Error functions, which is do nothing.
If the load balance algorithm used in the ServerSelector requires additional operations in the server's
returned response or error, shadow these methods and implement the logic accordingly. These handlers are
only scoped to load balance algorithm specific logics, as the general http response/error handling is done in
the ErrorHandler and ModifyResponse in proxy.go
*/
package serverSelector

import (
	"context"
	"errors"
	"log"
	"net/http"
)

type defaultSelector struct{}

func (*defaultSelector) ModifyResponse(*http.Response) error {
	return nil
}

func (ds *defaultSelector) HandleError(res http.ResponseWriter, req *http.Request, err error) {
	log.Printf("Proxy Error for %s: %v", req.URL.String(), err)

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		res.WriteHeader(http.StatusGatewayTimeout)
		res.Write([]byte("504 Gateway Timeout"))
	} else {
		res.WriteHeader(http.StatusBadGateway)
		res.Write([]byte("502 Bad Gateway"))
	}
}
