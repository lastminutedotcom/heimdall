package kubernetes

import (
	log "github.com/lastminutedotcom/heimdall/pkg/logging"
	"net/http"
)

var (
	Liveness = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})

	Readiness = http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	return
	})
)

func ConfigureDeployment(port string) {
	m := http.NewServeMux()
	log.Info("Starting kubernets probes")
	m.Handle("/readiness", Readiness)
	m.Handle("/liveness", Liveness)
	go func() {
		log.Fatal("%v", http.ListenAndServe(":"+port, m))
	}()
}
