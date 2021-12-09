package server

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HomeHandler() http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusForbidden)
		log.Warning("Trying to access non-acmeproxy URL")
	})

}
