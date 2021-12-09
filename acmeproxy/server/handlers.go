package server

import (
	golog "log"
	"net/http"
	"os"
	"time"

	"github.com/codeskyblue/realip"
	"github.com/go-acme/lego/v4/challenge"
	"golang.org/x/net/context"
)

const (
	ModeDefault   string = "default"
	ModeRaw       string = "raw"
	ActionPresent string = "present"
	ActionCleanup string = "cleanup"
)

type providerSolved interface {
	challenge.Provider
	CreateRecord(fqdn, value string) error
	RemoveRecord(fqdn, value string) error
}

// message represents the JSON payload
// See https://github.com/go-acme/lego/tree/master/providers/dns/httpreq
type messageDefault struct {
	FQDN  string `json:"fqdn"`
	Value string `json:"value"`
}

// message represents the JSON payload
// See https://github.com/go-acme/lego/tree/master/providers/dns/httpreq
type messageRaw struct {
	Domain  string `json:"domain"`
	Token   string `json:"token"`
	KeyAuth string `json:"keyauth"`
}

// Incomingmessage represents the JSON payload of an incoming request
// Should be either FQDN,Value or Domain,Token,KeyAuth
// See https://github.com/go-acme/lego/tree/master/providers/dns/httpreq
type messageIncoming struct {
	messageDefault
	messageRaw
}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

// AuthenticatorInterface is the interface implemented by BasicAuth
// FIXME: is this deprecated?
type AuthenticatorInterface interface {
	// NewContext returns a new context carrying authentication
	// information extracted from the request.
	NewContext(ctx context.Context, r *http.Request) context.Context
}

// writeAccessLog Logs the Http Status for a request into fileHandler and returns a httphandler function which is a wrapper to log the requests.
func writeAccessLog(handle http.Handler, accessLogHandle *os.File) http.HandlerFunc {
	logger := golog.New(accessLogHandle, "", 0)
	return func(w http.ResponseWriter, request *http.Request) {
		writer := statusWriter{w, 0, 0}
		handle.ServeHTTP(&writer, request)
		end := time.Now()
		statusCode := writer.status
		length := writer.length
		if request.URL.RawQuery != "" {
			logger.Printf("%v %s %s \"%s %s%s%s %s\" %d %d \"%s\"", end.Format("2006/01/02 15:04:05"), request.Host, realip.FromRequest(request), request.Method, request.URL.Path, "?", request.URL.RawQuery, request.Proto, statusCode, length, request.Header.Get("User-Agent"))
		} else {
			logger.Printf("%v %s %s \"%s %s %s\" %d %d \"%s\"", end.Format("2006/01/02 15:04:05"), request.Host, realip.FromRequest(request), request.Method, request.URL.Path, request.Proto, statusCode, length, request.Header.Get("User-Agent"))
		}
	}
}
