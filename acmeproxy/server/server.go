package server

import (
	"fmt"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	port           uint
	ProviderName   string
	AllowedDomains []string

	Provider   challenge.Provider
	httpServer *http.Server
}

func NewServer(port uint, providerName string, allowedDomains []string) (*Server, error) {
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	provider, err := dns.NewDNSChallengeProviderByName(providerName)
	if err != nil {
		return nil, fmt.Errorf("provider lookup of %s failed: %s", providerName, err.Error())
	}

	return &Server{
		port:           port,
		ProviderName:   providerName,
		AllowedDomains: allowedDomains,
		Provider:       provider,
		httpServer:     httpServer,
	}, nil
}

func (s Server) Run() {
	s.httpServer.Handler = s.GetHandler()

	log.WithFields(log.Fields{
		"endpoint": fmt.Sprintf("listening on http port %d", s.port),
		"addr":     s.httpServer.Addr,
	}).Info("Starting acmeproxy")
	log.Fatal(s.httpServer.ListenAndServe())
}

func (s *Server) GetHandler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/present", ActionHandler(ActionPresent, s))
	mux.Handle("/cleanup", ActionHandler(ActionCleanup, s))

	return mux
}
