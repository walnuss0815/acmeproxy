package server

import (
	"github.com/go-acme/lego/v4/challenge"
	"golang.org/x/net/context"
	"net/http"
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

// AuthenticatorInterface is the interface implemented by BasicAuth
// FIXME: is this deprecated?
type AuthenticatorInterface interface {
	// NewContext returns a new context carrying authentication
	// information extracted from the request.
	NewContext(ctx context.Context, r *http.Request) context.Context
}
