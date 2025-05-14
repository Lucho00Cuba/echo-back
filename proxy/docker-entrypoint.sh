#!/usr/bin/env sh
set -eu

# Log function
log() {
  echo "[entrypoint] $@"
}

# Run entrypoint scripts
if [ -d "/docker-entrypoint.d" ]; then
  log "Executing entrypoint scripts in /docker-entrypoint.d/"

  for script in /docker-entrypoint.d/*; do
    if [ -x "$script" ]; then
      log "Running $script"
      "$script"
    elif [ -f "$script" ]; then
      log "Sourcing $script"
      . "$script"
    fi
  done
fi

log "Starting NGINX..."

# Launch nginx
exec /usr/bin/nginx -c /etc/nginx/nginx.conf -g 'daemon off;'
