#!/usr/bin/env bash
# ==============================================================================
# NESCO Sovereign Enterprise Utility Platform - Master Launch Script
# ==============================================================================
# Automatically manages background daemons, cleans up stale port locks,
# mounts local Node/NPM engine, and launches Web Portal, Go Core API & Python AI Guardian.
# ==============================================================================

set -e
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PATH="$PROJECT_ROOT/bin:$PATH"

echo "=========================================================="
echo "⚡ NESCO Sovereign Enterprise Utility Ecosystem"
echo "=========================================================="

# 1. Clean up any stale processes from previous sessions
echo "Cleaning up any stale service locks..."
pkill -9 -f "cmd/server" 2>/dev/null || true
pkill -9 -f "app.main" 2>/dev/null || true
pkill -9 -f "next dev" 2>/dev/null || true
pkill -9 -f "next start" 2>/dev/null || true
for port in 3000 8080 8001; do
  fuser -k -9 "${port}/tcp" 2>/dev/null || true
done
sleep 1

# 2. Start AI Guardian Service (:8001)
echo "Starting AI Guardian Watchdog on port 8001..."
(cd "$PROJECT_ROOT/services/ai-guardian" && python3 -m app.main) > /tmp/nesco_ai_guardian.log 2>&1 &
AI_PID=$!

# 3. Start Go Core API Gateway (:8080)
echo "Starting Core API Gateway on port 8080..."
(cd "$PROJECT_ROOT/services/core-api" && go run ./cmd/server) > /tmp/nesco_core_api.log 2>&1 &
API_PID=$!

# 4. Start Next.js Web Portal (:3000)
echo "Starting Web Portal (Next.js) on port 3000..."
(cd "$PROJECT_ROOT/apps/web-portal" && npm run dev) > /tmp/nesco_web_portal.log 2>&1 &
WEB_PID=$!

# Trap Ctrl+C (SIGINT) and SIGTERM for graceful shutdown of all background services
cleanup() {
  echo ""
  echo "=========================================================="
  echo "🛑 Shutting down NESCO services..."
  kill -9 $AI_PID $API_PID $WEB_PID 2>/dev/null || true
  for port in 3000 8080 8001; do
    fuser -k -9 "${port}/tcp" 2>/dev/null || true
  done
  echo "✓ All services stopped cleanly."
  echo "=========================================================="
  exit 0
}
trap cleanup INT TERM

sleep 2

echo ""
echo "=========================================================="
echo "🚀 NESCO Utility Platform is LIVE & READY!"
echo "=========================================================="
echo "   🌐 Web Portal:        http://localhost:3000   (PID: $WEB_PID)"
echo "   ⚡ Core API Gateway:  http://localhost:8080   (PID: $API_PID)"
echo "   🤖 AI Guardian:       http://localhost:8001   (PID: $AI_PID)"
echo "=========================================================="
echo "Logs streaming to /tmp/nesco_*.log"
echo "Press Ctrl+C at any time to gracefully stop all services."
echo "=========================================================="
echo ""

wait
