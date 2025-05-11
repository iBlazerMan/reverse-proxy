package serverSelector

import "net/url"

type ServerSelector interface {
	SelectServer() *url.URL
}
