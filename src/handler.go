package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// healthz handles the health check endpoint and responds with a JSON status.
//
// params:
//   - w: the response writer
//   - r: the request
//
// returns:
//   - the response writer
func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]interface{}{"api": map[string]string{"healthz": "ok"}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusOK)
}

// version handles the health check endpoint and responds with a JSON status.
//
// params:
//   - w: the response writer
//   - r: the request
//
// returns:
//   - the response writer
func version(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]Metadata{"api": Metadata{Name: NAME, Version: VERSION, Commit: COMMIT, Email: EMAIL}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// w.WriteHeader(http.StatusOK)
}

// root handles the root path and returns a JSON response containing metadata and request information.
//
// params:
//   - w: the response writer
//   - r: the request
//
// returns:
//   - the response writer
func root(w http.ResponseWriter, r *http.Request) {
	requestCount.Inc()
	start := time.Now()
	addSecurityHeaders(w)

	var data ApiResponse

	data.API.Metadata.Name = NAME
	data.API.Metadata.Version = VERSION
	data.API.Metadata.Commit = COMMIT
	data.API.Spec.Server = HOST
	data.API.Metadata.Email = EMAIL

	request := &data.API.Spec.Request
	request.Host = r.Host
	request.Method = r.Method
	// request.URI = r.RequestURI
	request.URI = r.Header.Get(HeaderOriginalUri)
	if request.URI == "" {
		request.URI = r.RequestURI
	}
	request.RequestID = r.Header.Get(HeaderRequestId)
	// request.ClientAddr = r.RemoteAddr
	request.ClientAddr = r.Header.Get(HeaderClientAddr)
	if request.ClientAddr == "" {
		request.ClientAddr = r.RemoteAddr
	}
	// request.Scheme = r.Proto
	request.Scheme = r.Header.Get(HeaderScheme)
	if request.Scheme == "" {
		request.Scheme = "http"
	}

	// https://go.dev/src/net/http/status.go
	code := http.StatusOK
	// code := http.StatusNotFound
	if headerCode := r.Header.Get(HeaderHTTPCode); headerCode != "" {
		if parsedCode, err := strconv.Atoi(headerCode); err == nil {
			code = parsedCode
		}
	}
	code = sanitizeStatusCode(code)
	requestStatusCount.WithLabelValues(strconv.Itoa(code)).Inc()

	response := &data.API.Spec.Response
	response.Status = code
	response.StatusText, response.StatusEmoji = getStatusInfo(code)

	// Set additional response information from headers
	response.ServiceName = r.Header.Get(HeaderServiceName)
	response.ServicePort = r.Header.Get(HeaderServicePort)
	response.IngressName = r.Header.Get(HeaderIngressName)
	response.Namespace = r.Header.Get(HeaderNamespace)

	if DEBUG {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("failed to read request body", "err", err)
		}

		if json.Valid(body) {
			var data map[string]interface{}
			if err := json.Unmarshal(body, &data); err != nil {
				request.Body = string(body)
			} else {
				request.Body = data
			}
		} else {
			request.Body = string(body)
		}

		// w.headers
		w.Header().Set(HeaderRequestId, r.Header.Get(HeaderRequestId))
		w.Header().Set(HeaderHTTPCode, r.Header.Get(HeaderHTTPCode))
		w.Header().Set(HeaderFormat, r.Header.Get(HeaderFormat))
		w.Header().Set(HeaderNamespace, r.Header.Get(HeaderNamespace))
		w.Header().Set(HeaderIngressName, r.Header.Get(HeaderIngressName))
		w.Header().Set(HeaderServiceName, r.Header.Get(HeaderServiceName))
		w.Header().Set(HeaderServicePort, r.Header.Get(HeaderServicePort))
		w.Header().Set(HeaderServer, HOST)
		// r.headers
		request.Headers = r.Header
		request.Headers.Add(HeaderServer, HOST)
	}

	// content_type
	format := r.Header.Get(HeaderFormat)
	if format == "" {
		format = r.Header.Get("Accept")
	}
	if format == "" {
		format = "text/html"
	}

	if strings.Contains(format, "text/html") {
		w.WriteHeader(code)
		// HTML
		err := GetTemplate().Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// JSON
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			logger.Error("failed to write response", "err", err)
		}
	}

	defer func() {
		duration := time.Since(start).Seconds()

		// Observe total duration
		requestDuration.Observe(duration)

		// Observe duration labeled by method and HTTP code
		requestDurationVec.WithLabelValues(r.Method, strconv.Itoa(code)).Observe(duration)
	}()
}
