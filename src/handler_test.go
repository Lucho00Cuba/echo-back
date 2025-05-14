package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// setup sets defaults before running root handler
func setup() {
	// Set debug to true
	DEBUG = true
	// Enable logger
	InitLogger()
	// Set template
	tmp, _ := os.CreateTemp("", "tmpl.html")
	tmp.WriteString(`<html><body>{{.API.Metadata.Name}}</body></html>`)
	tmp.Close()
	TEMPLATE_HTML = tmp.Name()
	LoadTemplate(TEMPLATE_HTML)
}

type failingWriter struct {
	http.ResponseWriter
}

func (f *failingWriter) Write(p []byte) (int, error) {
	return 0, io.ErrClosedPipe
}

// TestHealthz ensures healthz endpoint returns expected response.
func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	healthz(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200; got %d", res.StatusCode)
	}
	var data map[string]map[string]string
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Errorf("Failed to decode JSON: %v", err)
	}
	if data["api"]["healthz"] != "ok" {
		t.Errorf("Expected healthz=ok; got %v", data)
	}
}

func TestHealthzWriteError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := &failingWriter{httptest.NewRecorder()}

	healthz(rec, req)
	// No panic should occur, just silently fail with 500
}

// TestVersion ensures version handler returns correct metadata.
func TestVersion(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := httptest.NewRecorder()

	version(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200; got %d", res.StatusCode)
	}
	var data map[string]Metadata
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		t.Errorf("Failed to decode JSON: %v", err)
	}
	if data["api"].Version != VERSION {
		t.Errorf("Expected version %s; got %s", VERSION, data["api"].Version)
	}
}

func TestVersionWriteError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rec := &failingWriter{httptest.NewRecorder()}

	version(rec, req)
}

// TestRoot_JSON ensures root returns valid JSON with DEBUG enabled
func TestRoot_JSON(t *testing.T) {
	setup()

	body := `{"hello":"world"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set(HeaderHTTPCode, "201")
	req.Header.Set("Accept", "application/json")
	req.Header.Set(HeaderFormat, "application/json")
	req.Header.Set(HeaderClientAddr, "1.2.3.4")
	req.Header.Set(HeaderRequestId, "test-id")

	rec := httptest.NewRecorder()

	root(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != 201 {
		t.Errorf("Expected 201; got %d", res.StatusCode)
	}
	if ct := res.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Errorf("Expected JSON response; got %s", ct)
	}
	var api ApiResponse
	if err := json.NewDecoder(res.Body).Decode(&api); err != nil {
		t.Errorf("Failed to parse JSON: %v", err)
	}
	if api.API.Spec.Response.Status != 201 {
		t.Errorf("Expected response status 201; got %d", api.API.Spec.Response.Status)
	}
}

// TestRoot_HTML ensures root returns HTML when Accept header is text/html
func TestRoot_HTML(t *testing.T) {
	setup()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept", "text/html")
	req.Header.Set(HeaderFormat, "text/html")

	rec := httptest.NewRecorder()

	root(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected 200; got %d", res.StatusCode)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(res.Body)
	html := buf.String()
	if !strings.Contains(html, NAME) {
		t.Errorf("Expected HTML response to contain name %s; got %s", NAME, html)
	}
}

type brokenReader struct{}

func (brokenReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func (brokenReader) Close() error {
	return nil
}

func TestRootBodyReadError(t *testing.T) {
	setup()

	req := httptest.NewRequest(http.MethodPost, "/", brokenReader{})
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	root(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Errorf("Expected 200; got %d", rec.Result().StatusCode)
	}
}

func TestRootJSONEncodeError(t *testing.T) {
	setup()
	body := `{"hello": "world"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Accept", "application/json")
	req.Header.Set(HeaderFormat, "application/json")

	rec := &failingWriter{httptest.NewRecorder()}
	root(rec, req)
}
