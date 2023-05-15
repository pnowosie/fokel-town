package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
)

const (
	ServiceName = "merkle-service"
	// Version follows ZeroVer versioning schema, see: https://0ver.org
	Version = "0.0.1"
	HOST    = "localhost"
	PORT    = 4000
)

func main() {
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:  ServiceName,
		Level: hclog.LevelFromString("DEBUG"),
	})
	appLogger.Info("Simple HTTP service")

	host := flag.String("host", HOST, "host to listen on")
	port := flag.Int("port", PORT, "port to listen on")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)

	startTime := time.Now().Unix()
	appLogger.Info("Starting server", "addr", addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		response, _ := json.Marshal(ApiVersion{ServiceName, Version, time.Now().Unix() - startTime})
		w.Write(response)
	})

	http.ListenAndServe(addr, mux)
}

type ApiVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	UpTime  int64  `json:"uptime"`
}
