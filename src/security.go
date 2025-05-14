package main

import "net/http"

// addSecurityHeaders sets common HTTP security and CORS headers
// to enhance protection against common web vulnerabilities
// and enable safe cross-origin requests.
//
// This includes:
//   - X-Content-Type-Options: prevents MIME-sniffing
//   - X-Frame-Options: disallows embedding in iframes (clickjacking protection)
//   - Referrer-Policy: avoids leaking referrer information
//   - Access-Control-Allow-Origin: enables cross-origin requests (CORS)
//
// params:
//   - w: the HTTP response writer to which the headers will be added
//
// returns:
//   - none (modifies the headers on the response writer)
func addSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Referrer-Policy", "no-referrer")
	w.Header().Set("Access-Control-Allow-Origin", "*")
}
