package kubernetes

import (
	"log"
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
	m.Handle("/readiness", Readiness)
	m.Handle("/liveness", Liveness)
	go func() {
		log.Fatalf("%v", http.ListenAndServe(":"+port, m))
	}()
}
