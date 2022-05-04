package main

import (
	"github.com/codegangsta/negroni"
	"net/http"
	"scaffold-api-server/cmd/app/handlers"
)

type Route struct {
	Method      string
	BasePath    string
	Pattern     string
	Handler     http.Handler
	Timeout     int
	NoAuth      bool
	Middlewares []negroni.Handler
}

func InitRoutes() []Route {
	return []Route{
		{
			Method:   http.MethodGet,
			BasePath: "v1",
			Pattern:  "health",
			Handler:  handlers.NewHealthCheckHandler(),
			NoAuth:   true,
		},
	}
}
