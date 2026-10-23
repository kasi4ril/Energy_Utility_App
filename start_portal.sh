#!/usr/bin/env bash
# Starts Next.js Web Portal using configured Node engine on port 3000
set -e
PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
export PATH="$PROJECT_ROOT/bin:$PATH"

echo "=========================================================="
echo "⚡ NESCO Enterprise Utility Web Portal"
echo "🌐 URL: http://localhost:3000"
echo "=========================================================="
cd "$PROJECT_ROOT/apps/web-portal"
exec npm run dev

