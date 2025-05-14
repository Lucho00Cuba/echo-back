package main

const (
	// HeaderServer is the name of the header that contains the server name.
	HeaderServer = "X-Server"

	// HeaderClientAddr is the name of the header that contains the client address.
	HeaderClientAddr = "X-Forwarded-For"

	// HeaderScheme is the name of the header that contains the scheme.
	HeaderScheme = "X-Scheme"

	// HeaderOriginalUri is URI that caused the error
	HeaderOriginalUri = "X-Original-URI"

	// HeaderHTTPCode is the name of the header used as the source of the HTTP status code to return.
	HeaderHTTPCode = "X-Code"

	// HeaderFormat name of the header that defines the format of the reply
	HeaderFormat = "X-Format"

	// HeaderNamespace is the name of the header that contains information about the Ingress namespace.
	HeaderNamespace = "X-Namespace"

	// HeaderIngressName is the name of the header that contains the matched Ingress.
	HeaderIngressName = "X-Ingress-Name"

	// HeaderServiceName is the name of the header that contains the matched Service in the Ingress.
	HeaderServiceName = "X-Service-Name"

	// HeaderServicePort is the name of the header that contains the matched Service port in the Ingress.
	HeaderServicePort = "X-Service-Port"

	// HeaderRequestId is a unique ID that identifies the request - same as for backend service.
	HeaderRequestId = "X-Request-ID"
)
