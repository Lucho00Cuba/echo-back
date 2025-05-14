package main

import (
	"os"
	"strings"
	"testing"
)

// Test_LoadTemplateAndGetTemplate verifies that a valid HTML template can be parsed and retrieved.
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
