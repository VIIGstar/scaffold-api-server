package main

import (
	"fmt"
	"logur.dev/logur"
	"net/http"
)

func Middleware(h http.Handler, logger logur.LoggerFacade) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("middleware, path request: %v", r.URL))
		h.ServeHTTP(w, r)
	})
}
