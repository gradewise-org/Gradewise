#!/bin/bash

# Gradewise Development Data Cleanup Script
# 
# This script safely cleans PostgreSQL persistent data by:
# 1. Stopping Tilt (to avoid conflicts)
# 2. Cleaning PostgreSQL persistent volume claims
# 3. Providing instructions to restart development
#
# Usage: ./backend/scripts/clean-dev-data.sh

set -e  # Exit on any error

echo "🧹 Gradewise Development Data Cleanup"
echo "======================================"
echo ""

# Check if Tilt is running
TILT_PIDS=$(pgrep -f "tilt up" || true)
if [ ! -z "$TILT_PIDS" ]; then
    echo "⚠️  Tilt is currently running (PID: $TILT_PIDS) and must be stopped for safe cleanup."
    echo "📌 This script will kill the running Tilt process to stop all services."
    echo ""
    read -p "Continue? (y/N): " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Cleanup cancelled."
        exit 1
    fi
    
    echo "🛑 Killing Tilt process(es)..."
    pkill -f "tilt up" || true
    
    # Wait a moment for processes to terminate
    sleep 2
    
    # Check if any Tilt processes are still running
    if pgrep -f "tilt" > /dev/null; then
        echo "⚠️  Some Tilt processes are still running, forcing termination..."
        pkill -9 -f "tilt" || true
        sleep 1
    fi
    
    echo "✅ Tilt stopped successfully."
    echo ""
else
    echo "✅ Tilt is not running - safe to proceed with cleanup."
    echo ""
fi

echo "🗑️  Cleaning PostgreSQL development data..."
echo ""

# Delete PostgreSQL pod first (releases the PVC lock)
echo "📦 Deleting PostgreSQL pods..."
kubectl delete pod -l app=postgres --ignore-not-found=true

# Wait a moment for pod to terminate
echo "⏳ Waiting for PostgreSQL pod to terminate..."
kubectl wait --for=delete pod -l app=postgres --timeout=60s || true

# Delete the persistent volume claim
echo "🗄️  Deleting PostgreSQL persistent volume claim..."
kubectl delete pvc postgres-storage --ignore-not-found=true

echo ""
echo "✨ Development data cleanup completed!"
echo ""
echo "🚀 Next steps:"
echo "   1. Run 'tilt up' to restart your development environment"
echo "   2. PostgreSQL will start with a fresh, empty database"
echo "   3. All initialization scripts will run automatically"
echo ""
echo "💡 Tip: Your application will be available at http://localhost:3000" 