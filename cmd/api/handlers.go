package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (app *application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(apiVersion{ServiceName, Version, time.Now().Unix() - app.startTime})
	w.Write(response)
}

type apiVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	UpTime  int64  `json:"uptime"`
}
