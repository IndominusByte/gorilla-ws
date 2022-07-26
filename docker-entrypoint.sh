#!/bin/sh
set -e

if [ "$BACKEND_STAGE" = 'production' ]; then
  # Create folder to store supervisor logs
  mkdir -p /var/log/supervisor

  # Run
  /usr/bin/supervisord -c /app/supervisord.conf
else
  make watch
fi
