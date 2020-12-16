package server

import (
	"context"
	"net/http"

	"github.com/schwarzeni/Go-000/Week04/pkg/appgroup"
)

type HTTPServerConfig struct {
	Addr string
}

var _ appgroup.App = &HTTPServer{}

type HTTPServer struct {
	svc *http.Server
}

func NewHTTPServer(config *HTTPServerConfig, handler http.Handler) *HTTPServer {
	return &HTTPServer{svc: &http.Server{Addr: config.Addr, Handler: handler}}
}

func (h HTTPServer) Start(ctx context.Context) error { return h.svc.ListenAndServe() }

func (h HTTPServer) Stop(ctx context.Context) error { return h.svc.Shutdown(ctx) }
