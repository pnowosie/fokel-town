package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) Routes() http.Handler {
	router := httprouter.New()

	// define routes
	router.HandlerFunc(http.MethodGet, "/v0/health", app.HealthCheck)
	router.HandlerFunc(http.MethodGet, "/v0/user/:id", app.GetUser)
	router.HandlerFunc(http.MethodPut, "/v0/user", app.PutUser)

	return app.wrapWithApi(router)
}
