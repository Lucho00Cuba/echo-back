package main

import (
	"net/http/httptest"
	"testing"
)

// TestAddSecurityHeaders verifies that all expected security headers are correctly set
func TestAddSecurityHeaders(t *testing.T) {
	rec := httptest.NewRecorder()

	addSecurityHeaders(rec)

	headers := rec.Result().Header

	tests := map[string]string{
		"X-Content-Type-Options":      "nosniff",
		"X-Frame-Options":             "DENY",
		"Referrer-Policy":             "no-referrer",
		"Access-Control-Allow-Origin": "*",
	}

	for key, expected := range tests {
		got := headers.Get(key)
		if got != expected {
			t.Errorf("expected header %s to be %q, got %q", key, expected, got)
		}
	}
}
