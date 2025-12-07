#!/bin/sh
set -e

# Quietly refresh man database if available
if command -v mandb >/dev/null 2>&1; then
    mandb -q || true
fi

exit 0

