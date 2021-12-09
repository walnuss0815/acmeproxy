package server

import (
	"encoding/json"
	"github.com/codeskyblue/realip"
	"github.com/go-acme/lego/v4/challenge/dns01"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func ActionHandler(action string, server *Server) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		alog := log.WithFields(log.Fields{
			"prefix": action + ": " + realip.FromRequest(r),
		})

		// Check if we're using POST
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			alog.WithField("method", r.Method).Error("Method not allowed")
			return
		}

		// Decode the JSON message
		incoming := &messageIncoming{}
		err := json.NewDecoder(r.Body).Decode(incoming)
		if err != nil {
			http.Error(w, "Bad JSON request", http.StatusBadRequest)
			alog.WithField("error", err.Error()).Error("Method not allowed")
			return
		}

		// Make sure domain and FQDN from the incoming message are correct
		incoming.FQDN = dns01.ToFqdn(incoming.FQDN)
		incoming.Domain = dns01.UnFqdn(incoming.Domain)

		// Check if we've received a message or messageRaw JSON
		// See https://github.com/go-acme/lego/tree/master/providers/dns/httpreq
		var mode string
		var checkDomain string
		if incoming.FQDN != "" && incoming.Value != "" {
			mode = ModeDefault
			checkDomain = dns01.UnFqdn(strings.TrimPrefix(incoming.FQDN, "_acme-challenge."))
			alog.WithFields(log.Fields{
				"fqdn":  incoming.FQDN,
				"value": incoming.Value,
			}).Debug("Received JSON payload (default mode)")
		} else if incoming.Domain != "" && (incoming.Token != "" || incoming.KeyAuth != "") {
			mode = ModeRaw
			checkDomain = incoming.Domain
			alog.WithFields(log.Fields{
				"domain":  incoming.Domain,
				"token":   incoming.Token,
				"keyAuth": incoming.KeyAuth,
			}).Debug("Received JSON payload (raw mode)")
		} else {
			http.Error(w, "Wrong JSON content", http.StatusBadRequest)
			alog.WithField("json", incoming).Error("Wrong JSON content")
			return
		}

		// Check if we are allowed to requests certificates for this domain
		var allowed = false
		for _, allowedDomain := range server.AllowedDomains {
			alog.WithFields(log.Fields{
				"checkDomain":   checkDomain,
				"allowedDomain": allowedDomain,
			}).Debug("Checking allowed domain")
			if checkDomain == allowedDomain || strings.HasSuffix(strings.SplitAfterN(checkDomain, ".", 2)[1], allowedDomain) {
				allowed = true
				break
			}
		}

		if !allowed {
			http.Error(w, "Requested domain not in allowed-domains", http.StatusInternalServerError)
			alog.WithFields(log.Fields{
				"domain":          checkDomain,
				"allowed-domains": server.AllowedDomains,
			}).Debug("Requested domain not in allowed-domains")
			return
		}

		// Check if this provider supports the selected mode
		// We assume that all providers support MODE_RAW (which is lego default)
		if mode == ModeDefault {
			provider, ok := server.Provider.(providerSolved)
			if ok {
				alog.WithFields(log.Fields{
					"provider": server.ProviderName,
					"mode":     mode,
				}).Debug("Provider supports requested mode")

				if action == ActionPresent {
					err = provider.CreateRecord(incoming.FQDN, incoming.Value)
				} else if action == ActionCleanup {
					err = provider.RemoveRecord(incoming.FQDN, incoming.Value)
				} else {
					alog.WithFields(log.Fields{
						"provider": server.ProviderName,
						"fqdn":     incoming.FQDN,
						"value":    incoming.Value,
						"mode":     mode,
						"error":    err.Error(),
					}).Error("Wrong action specified")
					http.Error(w, "Wrong action specified", http.StatusInternalServerError)
					return
				}

				if err != nil {
					alog.WithFields(log.Fields{
						"provider": server.ProviderName,
						"fqdn":     incoming.FQDN,
						"value":    incoming.Value,
						"mode":     mode,
						"error":    err.Error(),
					}).Errorf("Failed to update TXT record, %v", err)
					http.Error(w, "Failed to update TXT record", http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, "Provider does not support requested mode", http.StatusInternalServerError)
				alog.WithFields(log.Fields{
					"provider": server.ProviderName,
					"mode":     mode,
				}).Debug("Provider does not support requested mode")
				return
			}

			// Send back the original JSON to confirm success
			m := messageDefault{FQDN: incoming.FQDN, Value: incoming.Value}
			w.Header().Set("Content-Type", "application/json")
			returnErr := json.NewEncoder(w).Encode(m)
			if returnErr != nil {
				log.Error("Problem encoding return message")
			}

			// Succes!
			alog.WithFields(log.Fields{
				"provider": server.ProviderName,
				"fqdn":     incoming.FQDN,
				"value":    incoming.Value,
				"mode":     mode,
			}).Info("Sucessfully updated TXT record")
			// All lego providers should support raw mode
		} else if mode == ModeRaw {
			fqdn, value := dns01.GetRecord(incoming.Domain, incoming.KeyAuth)
			alog.WithFields(log.Fields{
				"provider": server.ProviderName,
				"mode":     mode,
			}).Debug("Provider supports requested mode")
			err = server.Provider.Present(incoming.Domain, incoming.Token, incoming.KeyAuth)
			if err != nil {
				alog.WithFields(log.Fields{
					"provider": server.ProviderName,
					"domain":   incoming.Domain,
					"fqdn":     fqdn,
					"token":    incoming.Token,
					"keyAuth":  incoming.KeyAuth,
					"value":    value,
					"mode":     mode,
				}).Errorf("Failed to update TXT record, %v", err)
				http.Error(w, "Failed to update TXT record", http.StatusInternalServerError)
				return
			}
			// Send back the original JSON to confirm success
			m := messageRaw{Domain: incoming.Domain, Token: incoming.Token, KeyAuth: incoming.KeyAuth}
			w.Header().Set("Content-Type", "application/json")
			returnErr := json.NewEncoder(w).Encode(m)
			if returnErr != nil {
				log.Error("Problem encoding return message")
			}

			// Succes!
			alog.WithFields(log.Fields{
				"provider": server.ProviderName,
				"domain":   incoming.Domain,
				"fqdn":     fqdn,
				"token":    incoming.Token,
				"keyAuth":  incoming.KeyAuth,
				"value":    value,
				"mode":     mode,
			}).Info("Sucessfully updated TXT record")
		} else {
			http.Error(w, "Unkown mode requested", http.StatusInternalServerError)
			alog.WithFields(log.Fields{
				"provider": server.ProviderName,
				"mode":     mode,
			}).Info("Unknown mode requested")
			return
		}

	})

}
