package server

import (
	"fmt"
	auth "github.com/abbot/go-http-auth"
	"github.com/go-acme/lego/v4/challenge"
	"github.com/go-acme/lego/v4/providers/dns"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type Server struct {
	ProviderName   string
	HttpServer     *http.Server
	HtpasswdFile   string
	AllowedIPs     []string
	AllowedDomains []string
	AccessLogFile  string
	Port           int
	Provider       challenge.Provider
}

func NewServer(config *Config) (*Server, error) {
	httpServer := &http.Server{
		Addr: fmt.Sprintf(":%d", config.Port),
	}

	provider, err := dns.NewDNSChallengeProviderByName(config.ProviderName)
	if err != nil {
		return nil, fmt.Errorf("provider lookup of %s failed: %s", config.ProviderName, err.Error())
	}

	return &Server{
		HttpServer:     httpServer,
		HtpasswdFile:   config.HtpasswdFile,
		AllowedIPs:     config.AllowedIPs,
		AllowedDomains: config.AllowedDomains,
		AccessLogFile:  config.AccessLogFile,
		Port:           config.Port,
		Provider:       provider,
	}, nil
}

func (s Server) Run() {
	s.HttpServer.Handler = s.GetHandler()

	log.WithFields(log.Fields{
		"endpoint": fmt.Sprintf("listening on http port %d", s.Port),
		"addr":     s.HttpServer.Addr,
	}).Info("Starting acmeproxy")
	log.Fatal(s.HttpServer.ListenAndServe())
}

func (s *Server) GetHandler() http.Handler {
	// Define routes
	mux := http.NewServeMux()

	handlerPresent := ActionHandler(ActionPresent, s)
	handlerCleanup := ActionHandler(ActionCleanup, s)

	if len(s.HtpasswdFile) > 0 {
		authenticator := &auth.BasicAuth{
			Realm:   "Basic Realm",
			Secrets: auth.HtpasswdFileProvider(s.HtpasswdFile),
		}
		handlerPresent = AuthenticationHandler(handlerPresent, ActionPresent, authenticator)
		handlerCleanup = AuthenticationHandler(handlerCleanup, ActionCleanup, authenticator)
	}

	if len(s.AllowedIPs) > 0 {
		handlerPresent = FilterHandler(handlerPresent, ActionPresent, s)
		handlerCleanup = FilterHandler(handlerCleanup, ActionCleanup, s)
	}

	mux.Handle("/", HomeHandler())
	mux.Handle("/present", handlerPresent)
	mux.Handle("/cleanup", handlerCleanup)

	// Check if we need to write an access log
	var handler http.Handler
	if len(s.AccessLogFile) > 0 {
		accessLogHandle, err := os.OpenFile(s.AccessLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		defer accessLogHandle.Close()
		handler = writeAccessLog(mux, accessLogHandle)
	} else {
		handler = mux
	}

	return handler
}
