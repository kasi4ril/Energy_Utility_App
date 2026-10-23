#!/usr/bin/env bash
# Starts NESCO Go Core API Gateway on port 8080
set -e
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================================="
echo "⚡ NESCO Core API Gateway"
echo "🌐 URL: http://localhost:8080"
echo "=========================================================="
cd "$PROJECT_ROOT/services/core-api"
exec go run cmd/server/main.go

