package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

var (
	// VERSION is the current version of the application.
	VERSION string = "dev"
	// COMMIT is the commit hash of the current version.
	COMMIT string = "000"
	// NAME is the name of the application.
	NAME string = "echo-back"
	// EMAIL is the email of the application.
	EMAIL string = getEnv("EMAIL", "xe@xe.com")
	// PORT is the port the server listens on, defaulting to 3000.
	PORT string = getEnv("PORT", "3000")
	// HOST is the hostname of the server, defaulting to localhost.
	HOST string = getEnv("HOSTNAME", "localhost")
	// DEBUG indicates whether debugging is enabled.
	DEBUG, _ = strconv.ParseBool(getEnv("DEBUG", "true"))
	// TEMPLATE_HTML is the path to the HTML template file.
	TEMPLATE_HTML string = getEnv("TEMPLATE_HTML", "templates/simple.html")
)

// getEnv returns the value of the specified environment variable,
// or a fallback value if the variable is not set.
//
// params:
//   - key: the name of the environment variable
//   - defaultValue: the fallback value to use if the variable is not set
//
// returns:
//   - the value of the environment variable or the provided default
func getEnv(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}

// parseBool parses a string into a boolean value.
// It accepts values such as "1", "t", "true", "TRUE" (true)
// and "0", "f", "false", "FALSE" (false).
//
// If parsing fails, it returns false by default.
//
// params:
//   - val: the string to parse
//
// returns:
//   - the parsed boolean value, or false on error
func parseBool(val string) bool {
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return false
	}
	return parsed
}

// validateEmail checks if the email looks valid (contains @ and . after it)
//
// params:
//   - email: the email string to validate
//
// returns:
//   - true if the email has a basic valid format, false otherwise
func validateEmail(email string) bool {
	at := strings.Index(email, "@")
	dot := strings.LastIndex(email, ".")
	return at > 0 && dot > at+1 && dot < len(email)-1
}

// fileExists checks if a file exists at the given path.
//
// params:
//   - path: the path to the file
//
// returns:
//   - true if the file exists and is not a directory
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// validatePortRange checks if the provided port string is a valid TCP port number
// within the range 1024–65535 (excluding privileged or invalid ports).
//
// params:
//   - portStr: the port string to validate
//
// returns:
//   - true if the port is a valid number and within the safe range
func validatePortRange(portStr string) bool {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false
	}
	return port >= 1024 && port <= 65535
}

// validateConfig validates required runtime configuration values like
// EMAIL, TEMPLATE_HTML path, and PORT range. It logs fatal errors
// if any of the required values are missing or invalid.
//
// This function should be called during application startup (e.g., in init).
//
// panics:
//   - If EMAIL is not valid
//   - If the HTML template file does not exist
//   - If PORT is not within 1024–65535
func validateConfig() {
	if !validateEmail(EMAIL) {
		log.Fatalf("Invalid EMAIL configured: %s", EMAIL)
	}
	if !fileExists(TEMPLATE_HTML) {
		log.Fatalf("Template file does not exist: %s", TEMPLATE_HTML)
	}
	if !validatePortRange(PORT) {
		log.Fatalf("Invalid PORT configured: %s (must be between 1024-65535)", PORT)
	}
}
