package qb

import (
	"net/http"
	"net/http/httptest"
)

func NewServerClientPair(fn http.HandlerFunc) (*httptest.Server, Client) {
	h := http.HandlerFunc(fn)
	server := httptest.NewServer(h)

	cfg := newTestConfig()
	cfg.viper.Set("realm-host", server.URL)

	client := NewClient(cfg)
	client.HTTPClient = server.Client()

	return server, client
}
