package mgmt

import (
	"expvar"
	"fmt"
	"net/http"
	"time"

	"git01.bravofly.com/golang/appfw.git/pkg/logging"
	"git01.bravofly.com/golang/appfw.git/pkg/properties"
)

var (
	// router
	router = http.NewServeMux()

	// mgmtServer is the http server exposing health endpoints as well as metrics
	mgmtServer = &http.Server{
		Addr:         ":8082",
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		IdleTimeout:  10 * time.Second,
		Handler:      router,
	}
)

// This will provide the basic readiness and liveness endpoint served
func init() {
	// TODO(inge4pres) find a reasonable way of init() the readiness and liveness without setting them
	// router.Handle("/readiness", notReadyHandler())
	// router.Handle("/liveness", notReadyHandler())
	router.Handle("/configprops", configHandler)
	router.Handle("/metrics", expvar.Handler())

	go func() {
		if err := mgmtServer.ListenAndServe(); err != nil {
			logging.Fatal(fmt.Sprintf("error starting management server: %v", err), nil)
		}
	}()
}

func notReadyHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("KO"))
	})
}

var configHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	config, err := properties.ConfigJSON()
	if err != nil {
		http.Error(w, fmt.Sprintf("{\"error\": \"%v\"}", err), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(config)
})

// ConfigureReadiness lets you specify the handling function for readiness endpoint
func ConfigureReadiness(handler http.HandlerFunc) {
	router.HandleFunc("/readiness", handler)
}

// ConfigureLiveness lets you specify the handling function for liveness endpoint
func ConfigureLiveness(handler http.HandlerFunc) {
	router.HandleFunc("/liveness", handler)
}
