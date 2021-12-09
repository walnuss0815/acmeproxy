package server

import (
	auth "github.com/abbot/go-http-auth"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func AuthenticationHandler(h http.Handler, action string, a AuthenticatorInterface) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := a.NewContext(r.Context(), r)
		r = r.WithContext(ctx)

		// Check authentication
		authInfo := auth.FromContext(r.Context())
		authInfo.UpdateHeaders(w.Header())
		if authInfo == nil || !authInfo.Authenticated {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			log.Warning("Unauthorized request")
			return
		}
		log.WithField("username", authInfo.Username).Info("Authorized")
		h.ServeHTTP(w, r)
	})
}
