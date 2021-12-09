package server

import (
	"github.com/codeskyblue/realip"
	"github.com/orange-cloudfoundry/ipfiltering"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func FilterHandler(h http.Handler, action string, server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := realip.FromRequest(r)
		flog := log.WithFields(log.Fields{
			"prefix": action + ": " + ip,
			"ip":     ip,
		})

		//ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		f := ipfiltering.New(ipfiltering.Options{AllowedIPs: server.AllowedIPs, BlockByDefault: true, Logger: flog})
		if !f.Allowed(ip) {
			http.Error(w, "Requesting IP not in allowed-ips", http.StatusForbidden)
			flog.Warning("Access denied")
			return
		}
		//success!
		h.ServeHTTP(w, r)
	})
}
