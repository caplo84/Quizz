#!/bin/bash

set -e  # Exit on any error

echo "🔍 Checking server health..."
if ! curl -s http://localhost:8080/api/v1/health | jq . > /dev/null 2>&1; then
    echo "❌ Server is not running or health check failed"
    echo "💡 Try: docker-compose up -d"
    exit 1
fi

echo "📊 Checking GitHub rate limit status..."
curl -s http://localhost:8080/api/admin/sync/github/status | jq .

echo ""
echo "🔄 Syncing ALL programming languages automatically..."
echo "   (JavaScript, Python, Java, React, Vue, Go, etc.)"
if curl -X POST http://localhost:8080/api/admin/sync/github \
    -H "Content-Type: application/json" \
    -s | jq .; then
    echo ""
    echo "✅ Sync completed! Checking status again..."
    curl -s http://localhost:8080/api/admin/sync/github/status | jq .
else
    echo "❌ Sync failed! Check logs with: docker-compose logs backend"
    exit 1
fi
