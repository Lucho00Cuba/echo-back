package main

import (
	"os"
	"testing"
)

// Test_getEnv ensures getEnv returns correct values based on environment state.
func TestGetEnv(t *testing.T) {
	const key = "TEST_ENV_VAR"

	// case 1: value not set, should return default
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("failed to unset env: %v", err)
	}
	got := getEnv(key, "default")
	if got != "default" {
		t.Errorf("getEnv with unset var = %q; want 'default'", got)
	}

	// case 2: value set, should return value
	if err := os.Setenv(key, "actual"); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	got = getEnv(key, "default")
	if got != "actual" {
		t.Errorf("getEnv with set var = %q; want 'actual'", got)
	}
}

// Test_parseBool checks if parseBool works with valid and invalid values.
func TestParseBool(t *testing.T) {
	tests := map[string]bool{
		"true":  true,
		"TRUE":  true,
		"1":     true,
		"false": false,
		"FALSE": false,
		"0":     false,
		"foo":   false,
	}

	for input, expected := range tests {
		got := parseBool(input)
		if got != expected {
			t.Errorf("parseBool(%q) = %v; want %v", input, got, expected)
		}
	}
}

// Test_validateEmail verifies simple email format validation.
func TestValidateEmail(t *testing.T) {
	tests := map[string]bool{
		"user@example.com": true,
		"invalid@":         false,
		"@invalid.com":     false,
		"noatsign.com":     false,
		"user@domain.":     false,
		"user@.com":        false,
		"user@domain.com":  true,
	}

	for email, expected := range tests {
		if got := validateEmail(email); got != expected {
			t.Errorf("validateEmail(%q) = %v; want %v", email, got, expected)
		}
	}
}

// Test_validatePortRange ensures that only valid port strings pass.
func TestValidatePortRange(t *testing.T) {
	tests := map[string]bool{
		"80":    false, // reserved
		"1023":  false,
		"1024":  true,
		"3000":  true,
		"65535": true,
		"65536": false,
		"99999": false,
		"abc":   false,
		"":      false,
		"-100":  false,
	}

	for input, expected := range tests {
		if got := validatePortRange(input); got != expected {
			t.Errorf("validatePortRange(%q) = %v; want %v", input, got, expected)
		}
	}
}

// Test_fileExists tests file existence check.
func TestFileExists(t *testing.T) {
	// Create temporary file
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("cleanup failed: %v", err)
		}
	}()

	if !fileExists(tmpFile.Name()) {
		t.Errorf("fileExists(%q) = false; want true", tmpFile.Name())
	}

	if fileExists("/definitely/does/not/exist") {
		t.Errorf("fileExists(fakePath) = true; want false")
	}
}
