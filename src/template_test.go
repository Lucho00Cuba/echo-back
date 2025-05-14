package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestLoadTemplateAndGetTemplate verifies that a valid HTML template can be parsed and retrieved.
func TestLoadTemplateAndGetTemplate(t *testing.T) {
	// Create a temporary HTML template file
	content := `<html><body><h1>{{.Title}}</h1></body></html>`
	tmpFile, err := os.CreateTemp("", "template-*.html")
	if err != nil {
		t.Fatalf("failed to create temp template: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			t.Logf("cleanup failed: %v", err)
		}
	}()

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("failed to write to temp template: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("failed to close temp file: %v", err)
	}

	// Load and retrieve the template
	LoadTemplate(tmpFile.Name())
	tmpl := GetTemplate()

	if tmpl == nil {
		t.Fatal("GetTemplate returned nil after LoadTemplate")
	}

	// Execute template with dummy data
	var sb strings.Builder
	err = tmpl.Execute(&sb, map[string]string{"Title": "Test Title"})
	if err != nil {
		t.Errorf("template execution failed: %v", err)
	}

	output := sb.String()
	if !strings.Contains(output, "Test Title") {
		t.Errorf("template output did not contain expected value: %s", output)
	}
}

// TestTriggerLoadTemplateFailure is invoked in a subprocess to trigger LoadTemplate failure.
func TestTriggerLoadTemplateFailure(t *testing.T) {
	if os.Getenv("TRIGGER_FAIL") == "1" {
		LoadTemplate("nonexistent-file.html")
	}
}

// TestLoadTemplateFailure validates that LoadTemplate logs fatal error on invalid file.
func TestLoadTemplateFailure(t *testing.T) {
	cmd := exec.Command(os.Args[0], "-test.run=TestTriggerLoadTemplateFailure")
	cmd.Env = append(os.Environ(), "TRIGGER_FAIL=1")

	output, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected LoadTemplate to exit with error, but got none. Output:\n%s", string(output))
	}

	if !strings.Contains(string(output), "failed to load template") {
		t.Errorf("expected log output to contain failure message. Got:\n%s", string(output))
	}
}
