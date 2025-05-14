package util

import (
	"context"
	"errors"
	"net/url"
)

type contextKey string

const (
	ServerUrlKey contextKey = "serverUrl"
)

func WithServerUrl(ctx context.Context, serverUrl *url.URL) context.Context {
	return context.WithValue(ctx, ServerUrlKey, serverUrl)
}

func GetServerUrl(ctx context.Context) (*url.URL, error) {
	serverUrl, ok := ctx.Value(ServerUrlKey).(*url.URL)
	if !ok {
		return nil, errors.New("failed to retrieve server url from the context given")
	}
	return serverUrl, nil
}
