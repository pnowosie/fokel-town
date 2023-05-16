package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/pnowosie/fokeltown-merkle/internal"
)

const (
	ServiceName = "merkle-service"
	// Version follows ZeroVer versioning schema, see: https://0ver.org
	Version = "0.0.1"
	HOST    = "localhost"
	PORT    = 4000
)

type application struct {
	logger    hclog.Logger
	trie      internal.Trie
	startTime int64
}

func main() {
	appLogger := hclog.New(&hclog.LoggerOptions{
		Name:  ServiceName,
		Level: hclog.LevelFromString("DEBUG"),
	})
	appLogger.Info(ServiceName, "version", Version)

	// HTTP server configuration
	host := flag.String("host", HOST, "host to listen on")
	port := flag.Int("port", PORT, "port to listen on")
	flag.Parse()
	addr := fmt.Sprintf("%s:%d", *host, *port)

	// Run the application
	appLogger.Info("Starting server", "addr", addr)
	http.ListenAndServe(addr, newApp(appLogger).Routes())
}

func newApp(logger hclog.Logger) *application {
	trie := &internal.ThreadSafeTrie{Trie: &internal.MapIsNotATrie{}}
	return &application{logger: logger, trie: trie, startTime: time.Now().Unix()}
}
