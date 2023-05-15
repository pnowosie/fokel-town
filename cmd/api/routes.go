package main

import "net/http"

func (app *application) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v0/health", app.HealthCheck)
	return app.wrapWithApi(mux)
}
