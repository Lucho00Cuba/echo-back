# Build stage
ARG GOLANG_VERSION=1.23.2
FROM golang:${GOLANG_VERSION} AS builder

ARG VERSION
ARG COMMIT

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/echo-back

# Copy the source code
COPY src/ .

# Run tests (optional for CI only — comment out for release build)
RUN go mod tidy && go mod verify && \
    go test -cover -v ./...

# Build the binary with reproducibility and metadata
RUN go build -trimpath -ldflags="-s -w -X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT}" \
    -o /go/bin/echo-back

# Final minimal image
FROM gcr.io/distroless/static:nonroot

ARG VERSION
ARG COMMIT
ARG BUILD_DATE

LABEL org.opencontainers.image.title="echo-back" \
      org.opencontainers.image.description="Minimal backend for Ingress NGINX errors" \
      org.opencontainers.image.authors="lomv0209@gmail.com" \
      org.opencontainers.image.url="https://github.com/Lucho00Cuba/echo-back" \
      org.opencontainers.image.source="https://github.com/Lucho00Cuba/echo-back" \
      org.opencontainers.image.revision=$COMMIT \
      org.opencontainers.image.version=$VERSION \
      org.opencontainers.image.created=$BUILD_DATE

WORKDIR /

COPY --from=builder /go/bin/echo-back /
COPY --from=builder /go/src/echo-back/templates /templates

EXPOSE 80 3000
USER nonroot:nonroot

CMD ["/echo-back"]
