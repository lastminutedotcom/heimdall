package kubernetes

import (
	"git01.bravofly.com/golang/appfw/pkg/http/mgmt"
	"net/http"
)

func Readiness() {
	mgmt.ConfigureReadiness(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
}
