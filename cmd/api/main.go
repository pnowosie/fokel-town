package main

import (
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

type application struct {
	logger    hclog.Logger
	trie      Trie
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
	trie := &MapIsNotATrie{}

	appLogger.Info("Starting server", "addr", addr)
	http.ListenAndServe(addr, NewApp(appLogger, trie).Routes())
}

func NewApp(logger hclog.Logger, trie Trie) *application {
	return &application{logger: logger, trie: trie, startTime: time.Now().Unix()}
}
