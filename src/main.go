package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var osExit = os.Exit

func main() {
	InitLogger()
	validateConfig()
	LoadTemplate(TEMPLATE_HTML)
	startServer(http.ListenAndServe)
}

// startServer initializes and runs the HTTP server.
// Accepts a function to allow injection in tests.
func startServer(listen func(addr string, handler http.Handler) error) {
	// Validate PORT
	if _, err := strconv.Atoi(PORT); err != nil {
		logger.Error("Invalid port number", "port", PORT, "err", err)
		osExit(1)
	}

	logger.Info("Starting HTTP server",
		"version", VERSION,
		"commit", COMMIT,
		"port", PORT,
	)

	mux := http.NewServeMux()
	mux.HandleFunc("/", LoggingMiddleware(root))
	mux.HandleFunc("/version", LoggingMiddleware(version))
	mux.HandleFunc("/healthz", LoggingMiddleware(healthz))
	mux.Handle("/metrics", promhttp.Handler())

	if err := listen(":"+PORT, mux); err != nil {
		logger.Error("Failed to start server", "err", err)
	}
}
