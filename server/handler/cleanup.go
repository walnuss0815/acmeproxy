package handler

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
	"github.com/walnuss0815/acmeproxy/v2/acme/proxy"
	"github.com/walnuss0815/acmeproxy/v2/server/models"
)

func Cleanup(proxy *(proxy.Proxy)) http.Handler {
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
			http.Error(w, "Unable to cleanup", http.StatusForbidden)
			log.WithField("error", err.Error()).Error("Unable to cleanup")
			return
		}

		err = proxy.Cleanup(request.FQDN, request.Value)
		if err != nil {
			http.Error(w, "Unable to present", http.StatusInternalServerError)
			log.WithField("error", err.Error()).Error("Unable to cleanup")
			return
		}

		// Send back the original JSON to confirm success
		w.Header().Set("Content-Type", "application/json")
		returnErr := json.NewEncoder(w).Encode(request)
		if returnErr != nil {
			log.Error("Problem encoding return message")
		}

		log.Infof("removed record for %s", request.FQDN)
	})
}
