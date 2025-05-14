package main

import "testing"

// Test_getStatusInfo verifies that known status codes return correct text and emoji,
// and that unknown codes return the default fallback values.
func TestGetStatusInfo(t *testing.T) {
	tests := []struct {
		code          int
		expectedText  string
		expectedEmoji string
	}{
		{200, "OK", "😃"},
		{404, "Not Found", "👀"},
		{418, "I'm a teapot", "💻"},
		{999, defaultText, defaultEmoji}, // unknown code
		{-1, defaultText, defaultEmoji},  // invalid code
	}

	for _, tt := range tests {
		text, emoji := getStatusInfo(tt.code)
		if text != tt.expectedText || emoji != tt.expectedEmoji {
			t.Errorf("getStatusInfo(%d) = (%q, %q); want (%q, %q)",
				tt.code, text, emoji, tt.expectedText, tt.expectedEmoji)
		}
	}
}

// Test_sanitizeStatusCode ensures that invalid codes are redirected to 418,
// while valid codes are passed through.
func TestSanitizeStatusCode(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{200, 200},
		{404, 404},
		{418, 418},
		{99, 418},
		{600, 418},
		{-10, 418},
	}

	for _, tt := range tests {
		got := sanitizeStatusCode(tt.input)
		if got != tt.expected {
			t.Errorf("sanitizeStatusCode(%d) = %d; want %d", tt.input, got, tt.expected)
		}
	}
}
