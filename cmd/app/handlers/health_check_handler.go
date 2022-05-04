package handlers

import "net/http"

type HealthCheckHandler struct{}

func NewHealthCheckHandler() HealthCheckHandler {
	return HealthCheckHandler{}
}

func (h HealthCheckHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
