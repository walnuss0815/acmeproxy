package server

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/walnuss0815/acmeproxy/v2/acme/proxy"
	"github.com/walnuss0815/acmeproxy/v2/server/handler"
)

type Server struct {
	port uint

	proxy      *proxy.Proxy
	httpServer *http.Server
}

func NewServer(port uint, proxy *(proxy.Proxy)) *Server {
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	return &Server{
		port:       port,
		proxy:      proxy,
		httpServer: httpServer,
	}
}

func (s *Server) Run() {
	s.httpServer.Handler = s.getHandler()

	log.WithFields(log.Fields{
		"endpoint": fmt.Sprintf("listening on http port %d", s.port),
		"addr":     s.httpServer.Addr,
	}).Info("Starting acmeproxy")
	log.Fatal(s.httpServer.ListenAndServe())
}

func (s *Server) getHandler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/present", handler.Present(s.proxy))
	mux.Handle("/cleanup", handler.Cleanup(s.proxy))

	return mux
}
