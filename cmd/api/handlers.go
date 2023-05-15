package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func (app *application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := apiVersion{
		Name:     ServiceName,
		Version:  Version,
		UpTime:   time.Now().Unix() - app.startTime,
		TrieRoot: app.trie.Root().String(),
	}
	jsonResponse, _ := json.Marshal(response)
	w.Write(jsonResponse)
}

type apiVersion struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	UpTime   int64  `json:"uptime"`
	TrieRoot string `json:"root"`
}
