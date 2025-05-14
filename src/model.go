package main

import "net/http"

// ApiResponse represents the response structure of the API.
type ApiResponse struct {
	API struct {
		Metadata Metadata `json:"metadata"`
		Spec     Spec     `json:"spec"`
	} `json:"api"`
}

// Metadata contains metadata information of the API.
type Metadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Commit  string `json:"commit"`
	Email   string `json:"email"`
}

// Spec defines the specification of the API.
type Spec struct {
	Server   string   `json:"server"`
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

// Request represents the structure of an API request.
type Request struct {
	Host       string      `json:"host"`
	Method     string      `json:"method"`
	URI        string      `json:"uri"`
	ClientAddr string      `json:"client_addr"`
	Scheme     string      `json:"scheme"`
	RequestID  string      `json:"request_id"`
	Headers    http.Header `json:"headers,omitempty"`
	Body       interface{} `json:"body,omitempty"`
}

// Response represents the structure of an API response.
type Response struct {
	ServiceName string `json:"service_name"`
	ServicePort string `json:"service_port"`
	IngressName string `json:"ingress_name"`
	Namespace   string `json:"namespace"`
	Status      int    `json:"status"`
	StatusText  string `json:"status_text"`
	StatusEmoji string `json:"status_emoji"`
}
