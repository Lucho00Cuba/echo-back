#!/usr/bin/env sh

set -e

entrypoint_log() {
    if [ -z "${NGINX_ENTRYPOINT_QUIET_LOGS:-}" ]; then
        echo "$@"
    fi
}

NGINX_PATH=/etc/nginx/ssl
KEYFILE=$NGINX_PATH/server.key
CERTFILE=$NGINX_PATH/server.crt

if [ ! -d "$NGINX_PATH" ]; then
    mkdir -p "$NGINX_PATH"
fi

if ! command -v openssl >/dev/null 2>&1; then
    entrypoint_log "openssl not found. Installing..."
    apk add --no-cache openssl
fi

DOMAINS="localhost api.example.io"
SAN=""
for item in $DOMAINS; do
    SAN="${SAN},DNS:${item}"
done
SAN="${SAN#,}" # remove leading comma

entrypoint_log "Generating cert for SANs: $SAN"

if [ ! -f "$KEYFILE" ] || [ ! -f "$CERTFILE" ]; then
    TMP_CONF=$(mktemp)

    cat > "$TMP_CONF" <<EOF
[req]
distinguished_name = req
req_extensions = v3_req
prompt = no

[req_distinguished_name]
CN = localhost

[v3_req]
keyUsage = digitalSignature, keyEncipherment
extendedKeyUsage = serverAuth
subjectAltName = $SAN
EOF

    openssl req -x509 -newkey rsa:4096 -days 365 -nodes \
        -keyout "$KEYFILE" -out "$CERTFILE" \
        -subj "/C=ES/ST=Madrid/O=Example/OU=API/CN=localhost" \
        -extensions v3_req -config "$TMP_CONF"

    rm -f "$TMP_CONF"
    echo "✅ Self-signed certificate generated."
else
    echo "📎 Certificate already exists."
fi
