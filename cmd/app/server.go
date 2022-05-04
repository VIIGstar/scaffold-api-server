package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"logur.dev/logur"
	"net/http"
	"time"
)

const DefaultAPITimeout = 3000

func initAppServer(config configuration) {
	const name = "app"
	apiLogger := logur.WithField(logger, "server", name)

	var handler http.Handler

	httpRouter := mux.NewRouter().StrictSlash(true)
	addRoutes(httpRouter)

	handler = httpRouter
	handler = Middleware(handler, logger)
	handler = CORSConfig(handler)

	apiLogger.Info("listening on address", map[string]interface{}{"address": config.App.HttpAddr})
	err := http.ListenAndServe(config.App.HttpAddr, handler)
	if err != nil {
		panic(err)
	}
}

func addRoutes(r *mux.Router) {
	routes := InitRoutes()
	for _, route := range routes {
		timeout := DefaultAPITimeout
		if route.Timeout > 0 {
			timeout = route.Timeout
		}

		var listHandler []negroni.Handler
		for _, middleware := range route.Middlewares {
			listHandler = append(listHandler, middleware)
		}
		handler := http.TimeoutHandler(route.Handler, time.Millisecond*time.Duration(timeout), "Timeout!")

		listHandler = append(listHandler, negroni.Wrap(handler))

		r.
			Methods(route.Method).
			Path(fmt.Sprintf("/%v/%v", route.BasePath, route.Pattern)).
			Handler(negroni.New(listHandler...))
	}
}

func CORSConfig(handler http.Handler) http.Handler {
	cors := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodPatch, http.MethodDelete}),
		handlers.AllowedHeaders([]string{"content-type"}),
	)

	return cors(handler)
}