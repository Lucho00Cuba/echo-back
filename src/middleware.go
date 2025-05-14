package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

// statusRecorder is a wrapper around http.ResponseWriter that captures
// the HTTP response status code and the size of the response body.
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// WriteHeader overrides the default WriteHeader to record the status code.
//
// params:
//   - code: the HTTP status code to write
func (rec *statusRecorder) WriteHeader(code int) {
	rec.statusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

// Write overrides the default Write method to track the response size.
//
// params:
//   - b: the byte slice to write
//
// returns:
//   - the number of bytes written and any error encountered
func (rec *statusRecorder) Write(b []byte) (int, error) {
	size, err := rec.ResponseWriter.Write(b)
	rec.size += size
	return size, err
}

// contextKey is a custom type used to define context keys
// and avoid collisions with other context values.
type contextKey string

const loggerKey = contextKey("logger")

// getRequestLogger retrieves the contextual slog.Logger from the request.
// If not found, it returns the global logger.
//
// params:
//   - r: the HTTP request containing the context
//
// returns:
//   - the slog.Logger for this request
//
//nolint:unused
func getRequestLogger(r *http.Request) *slog.Logger {
	if l, ok := r.Context().Value(loggerKey).(*slog.Logger); ok {
		return l
	}
	return logger
}

// loggingMiddleware wraps an HTTP handler and adds structured logging.
// It enriches the logger with request metadata (method, URI, IP, request ID),
// captures status codes, measures duration, and recovers from panics.
//
// It also injects a contextual logger into the request context.
//
// params:
//   - next: the HTTP handler to wrap
//
// returns:
//   - an HTTP handler function with logging and panic recovery
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// Extract useful fields
		requestID := r.Header.Get(HeaderRequestId)
		clientIP := r.Header.Get(HeaderClientAddr)
		if clientIP == "" {
			clientIP = r.RemoteAddr
		}
		url := r.Header.Get(HeaderOriginalUri)
		if url == "" {
			url = r.RequestURI
		}
		format := r.Header.Get(HeaderFormat)
		if format == "" {
			format = r.Header.Get("Accept")
		}
		if format == "" {
			format = "text/html"
		}

		// Build request-scoped logger
		reqLogger := logger.With(
			slog.String("request_id", requestID),
			slog.String("client_ip", clientIP),
			slog.String("method", r.Method),
			slog.String("uri", url),
		)

		// Inject logger into context
		ctx := context.WithValue(r.Context(), loggerKey, reqLogger)
		r = r.WithContext(ctx)

		defer func() {
			if rec := recover(); rec != nil {
				reqLogger.Error("panic recovered", "error", rec)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		// Call next handler
		next(rec, r)

		duration := time.Since(start).Seconds()

		reqLogger.Info("request handled",
			slog.Int("status", rec.statusCode),
			slog.Float64("duration_sec", duration),
			slog.Int("size_bytes", rec.size),
			slog.String("format", format),
		)
	}
}
