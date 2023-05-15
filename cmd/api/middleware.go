package main

import "net/http"

// wrapWithApi wraps the route handler with standard HTTP headers
func (app *application) wrapWithApi(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.logger.Debug("request", r.Method, r.URL.RequestURI())
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
