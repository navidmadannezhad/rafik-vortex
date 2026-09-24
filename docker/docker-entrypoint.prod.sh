#!/bin/sh
set -e

export GIN_MODE=release

echo "Starting code-reviewer (prod)..."
exec /app/code-reviewer
