package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-acme/lego/challenge/dns01"
	log "github.com/sirupsen/logrus"
	"github.com/walnuss0815/acmeproxy/v2/acme/proxy"
	"github.com/walnuss0815/acmeproxy/v2/server/models"
)

func Present(proxy *(proxy.Proxy)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if method post
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			log.WithField("method", r.Method).Error("Method not allowed")
			return
		}

		// Decode the JSON message
		incoming := &models.MessageIncoming{}
		err := json.NewDecoder(r.Body).Decode(incoming)
		if err != nil {
			http.Error(w, "Bad JSON request", http.StatusBadRequest)
			log.WithField("error", err.Error()).Error("Bad JSON request")
			return
		}

		request, err := decodeRequest(incoming)
		if err != nil {
			http.Error(w, "Bad JSON request", http.StatusBadRequest)
			log.WithField("error", err.Error()).Error("Bad JSON request")
			return
		}

		ok := proxy.CheckDomain(request.FQDN)
		if !ok {
			http.Error(w, "Unable to present", http.StatusForbidden)
			log.WithField("fqdn", request.FQDN).Error("Requested FQDN is not allowed")
			return
		}

		err = proxy.Present(request.FQDN, request.Value)
		if err != nil {
			http.Error(w, "Unable to present", http.StatusInternalServerError)
			log.WithField("error", err.Error()).Error("Unable to present")
			return
		}

		// Send back the original JSON to confirm success
		w.Header().Set("Content-Type", "application/json")
		returnErr := json.NewEncoder(w).Encode(request)
		if returnErr != nil {
			log.Error("Problem encoding return message")
		}

		log.Infof("new record for %s", request.FQDN)
	})
}

func decodeRequest(incoming *models.MessageIncoming) (models.MessageDefault, error) {
	// Make sure domain and FQDN from the incoming message are correct
	incoming.FQDN = dns01.ToFqdn(incoming.FQDN)
	incoming.Domain = dns01.UnFqdn(incoming.Domain)

	// Check if we've received a message or messageRaw JSON
	if incoming.FQDN != "" && incoming.Value != "" {
		return models.MessageDefault{
			FQDN:  incoming.FQDN,
			Value: incoming.Value,
		}, nil
	} else if incoming.Domain != "" && (incoming.Token != "" || incoming.KeyAuth != "") {
		fqdn, value := dns01.GetRecord(incoming.Domain, incoming.KeyAuth)

		return models.MessageDefault{
			FQDN:  fqdn,
			Value: value,
		}, nil
	}

	return models.MessageDefault{}, fmt.Errorf("wrong JSON content")
}
