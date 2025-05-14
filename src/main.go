package main

import (
	"net/http"
	"os"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	InitLogger()
	validateConfig()
	LoadTemplate(TEMPLATE_HTML)
	startServer()
}

func startServer() {
	// Validate PORT
	if _, err := strconv.Atoi(PORT); err != nil {
		logger.Error("Invalid port number", "port", PORT, "err", err)
		os.Exit(1)
	}

	logger.Info("Starting HTTP server",
		"version", VERSION,
		"commit", COMMIT,
		"port", PORT,
	)

	http.HandleFunc("/", LoggingMiddleware(root))
	http.HandleFunc("/version", LoggingMiddleware(version))
	http.HandleFunc("/healthz", LoggingMiddleware(healthz))
	http.Handle("/metrics", promhttp.Handler())

	// Start server
	if err := http.ListenAndServe(":"+PORT, nil); err != nil {
		logger.Error("Failed to start server", "err", err)
	}
}
