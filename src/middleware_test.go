package main

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestStatusRecorderWriteHeaderAndWrite ensures status and size are correctly recorded
func TestStatusRecorderWriteHeaderAndWrite(t *testing.T) {
	rec := httptest.NewRecorder()
	sr := &statusRecorder{ResponseWriter: rec, statusCode: 0}

	sr.WriteHeader(http.StatusTeapot)
	if sr.statusCode != http.StatusTeapot {
		t.Errorf("expected status %d, got %d", http.StatusTeapot, sr.statusCode)
	}

	body := []byte("test body")
	n, err := sr.Write(body)
	if err != nil {
		t.Errorf("unexpected write error: %v", err)
	}
	if n != len(body) {
		t.Errorf("expected write length %d, got %d", len(body), n)
	}
	if sr.size != len(body) {
		t.Errorf("expected recorded size %d, got %d", len(body), sr.size)
	}
}

// TestLoggingMiddleware verifies the middleware logs structured data and sets context
func TestLoggingMiddleware(t *testing.T) {
	// Setup a dummy handler
	var handlerCalled bool
	handler := LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		// Confirm logger is injected
		l := getRequestLogger(r)
		if l == nil {
			t.Error("expected logger in context, got nil")
		}
		w.WriteHeader(http.StatusAccepted)
		_, _ = w.Write([]byte("OK"))
	})

	// Prepare request
	req := httptest.NewRequest(http.MethodGet, "/test-uri", nil)
	req.Header.Set(HeaderRequestId, "req-123")
	req.Header.Set(HeaderClientAddr, "192.168.1.1")
	req.Header.Set(HeaderOriginalUri, "/test-uri")
	req.Header.Set(HeaderFormat, "application/json")

	// Capture logs (optional)
	var buf bytes.Buffer
	logger = slog.New(slog.NewTextHandler(&buf, nil))

	// Serve request
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if !handlerCalled {
		t.Error("handler was not called")
	}
	if rec.Code != http.StatusAccepted {
		t.Errorf("expected status %d, got %d", http.StatusAccepted, rec.Code)
	}
	if rec.Body.String() != "OK" {
		t.Errorf("unexpected response body: %s", rec.Body.String())
	}
}

// TestLoggingMiddlewarePanicRecovery ensures panic is caught and returns 500
func TestLoggingMiddlewarePanicRecovery(t *testing.T) {
	panicHandler := LoggingMiddleware(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	// Setup logger to avoid nil
	logger = slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))

	panicHandler.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 status on panic, got %d", rec.Code)
	}
}

// TestGetRequestLoggerDefault ensures fallback to global logger when missing in context
func TestGetRequestLoggerDefault(t *testing.T) {
	// fallback when context is empty
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	logger = slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))

	got := getRequestLogger(req)
	if got == nil {
		t.Error("expected fallback logger, got nil")
	}
}

// TestGetRequestLoggerContext ensures logger is fetched from context
func TestGetRequestLoggerContext(t *testing.T) {
	buf := &bytes.Buffer{}
	customLogger := slog.New(slog.NewTextHandler(buf, nil))

	ctx := context.WithValue(context.Background(), loggerKey, customLogger)
	req := httptest.NewRequest(http.MethodGet, "/", nil).WithContext(ctx)

	got := getRequestLogger(req)
	if got != customLogger {
		t.Error("expected custom logger from context")
	}
}
