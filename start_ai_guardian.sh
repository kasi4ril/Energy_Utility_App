#!/usr/bin/env bash
# Starts NESCO Python AI Guardian Service on port 8001
set -e
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================================="
echo "⚡ NESCO AI Guardian Watchdog & @GridAssistant"
echo "🌐 URL: http://localhost:8001"
echo "=========================================================="
cd "$PROJECT_ROOT/services/ai-guardian"
exec python3 -m app.main

