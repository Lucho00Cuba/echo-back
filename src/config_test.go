package main

import (
	"os"
	"os/exec"
	"strings"
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

// TestValidateConfigSuccess checks the happy path of validateConfig
func TestValidateConfigSuccess(t *testing.T) {
	// Create a valid temp template file
	tmpFile, err := os.CreateTemp("", "*.html")
	if err != nil {
		t.Fatalf("failed to create temp template: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	t.Setenv("EMAIL", "valid@example.com")
	t.Setenv("TEMPLATE_HTML", tmpFile.Name())
	t.Setenv("PORT", "3000")

	// Should not panic or exit
	validateConfig()
}

// TestValidateConfigFailures runs separate subprocesses to test fatal errors
func TestValidateConfigFailures(t *testing.T) {
	tests := []struct {
		name   string
		env    map[string]string
		expect string
	}{
		{
			name: "Invalid Email",
			env: map[string]string{
				"EMAIL":         "invalid",
				"TEMPLATE_HTML": "templates/simple.html", // dummy
				"PORT":          "3000",
			},
			expect: "Invalid EMAIL configured",
		},
		{
			name: "Missing Template File",
			env: map[string]string{
				"EMAIL":         "valid@example.com",
				"TEMPLATE_HTML": "/non/existent/file.html",
				"PORT":          "3000",
			},
			expect: "Template file does not exist",
		},
		{
			name: "Invalid Port",
			env: map[string]string{
				"EMAIL":         "valid@example.com",
				"TEMPLATE_HTML": "templates/simple.html", // dummy
				"PORT":          "99999",
			},
			expect: "Invalid PORT configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestHelperValidateConfig$")
			cmd.Env = append(os.Environ(), "GO_WANT_HELPER_VALIDATE=1")
			for k, v := range tt.env {
				cmd.Env = append(cmd.Env, k+"="+v)
			}
			output, err := cmd.CombinedOutput()
			if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() == 0 {
				t.Fatalf("expected failure exit for %s, got success", tt.name)
			}
			if !strings.Contains(string(output), tt.expect) {
				t.Errorf("expected output to contain %q, got: %s", tt.expect, output)
			}
		})
	}
}

// TestHelperValidateConfig runs validateConfig in a subprocess for fatal tests
func TestHelperValidateConfig(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_VALIDATE") != "1" {
		return
	}
	validateConfig()
}
